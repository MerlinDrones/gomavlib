package frame

import (
	"bufio"
	"fmt"
	"io"
	"time"

	"github.com/merlindrones/gomavlib/pkg/dialect"
	"github.com/merlindrones/gomavlib/pkg/message"
)

const (
	bufferSize = 512 // frames cannot go beyond len(header) + 255 + len(check) + len(sig)
)

// 1st January 2015 GMT
var signatureReferenceDate = time.Date(2015, 0o1, 0o1, 0, 0, 0, 0, time.UTC)

// ReadError is the error returned in case of non-fatal parsing errors.
type ReadError struct {
	str string
}

func (e ReadError) Error() string {
	return e.str
}

func newError(format string, args ...interface{}) ReadError {
	return ReadError{
		str: fmt.Sprintf(format, args...),
	}
}

// ReaderConf is the configuration of a Reader.
type ReaderConf struct {
	// the underlying bytes reader.
	Reader io.Reader

	// (optional) the dialect which contains the messages that will be read.
	// If not provided, messages are decoded into the MessageRaw struct.
	DialectRW *dialect.ReadWriter

	// (optional) the secret key used to validate incoming frames.
	// Non-signed frames are discarded. This feature requires v2 frames.
	InKey *V2Key
}

// Reader is a Frame reader.
type Reader struct {
	conf                 ReaderConf
	br                   *bufio.Reader
	curReadSignatureTime uint64
}

// NewReader allocates a Reader.
func NewReader(conf ReaderConf) (*Reader, error) {
	if conf.Reader == nil {
		return nil, fmt.Errorf("Reader not provided")
	}

	return &Reader{
		conf: conf,
		br:   bufio.NewReaderSize(conf.Reader, bufferSize),
	}, nil
}

func (r *Reader) ReadFrom(buf *bufio.Reader) (*V2Frame, error) {
	r.br = bufio.NewReaderSize(buf, bufferSize)
	f, err := r.Read()
	return f, err
}

// Read reads a Frame from the reader.
// It must not be called by multiple routines in parallel.
func (r *Reader) Read() (*V2Frame, error) {
	magicByte, err := r.br.ReadByte()
	if err != nil {
		return nil, err
	}

	if magicByte != V2MagicByte {
		return nil, newError("invalid magic byte: %x", magicByte)
	}
	f := &V2Frame{}

	err = f.Decode(r.br)
	if err != nil {
		return nil, newError(err.Error())
	}

	if r.conf.InKey != nil {

		if sig := f.GenerateSignature(r.conf.InKey); *sig != *f.Signature {
			return nil, newError("wrong signature")
		}

		// in UDP, packet order is not guaranteed. Therefore, we accept frames
		// with a timestamp within 10 seconds with respect to the previous
		if r.curReadSignatureTime > 0 &&
			f.SignatureTimestamp < (r.curReadSignatureTime-(10*100000)) {
			return nil, newError("signature timestamp is too old")
		}

		if f.SignatureTimestamp > r.curReadSignatureTime {
			r.curReadSignatureTime = f.SignatureTimestamp
		}
	}

	// decode message if in dialect and validate checksum
	if r.conf.DialectRW != nil {
		if mp := r.conf.DialectRW.GetMessage(f.GetMessage().GetID()); mp != nil {
			if sum := f.GenerateChecksum(mp.CRCExtra()); sum != f.GetChecksum() {
				return nil, newError("wrong checksum, expected %.4x, got %.4x, message id is %d",
					sum, f.GetChecksum(), f.GetMessage().GetID())
			}

			msg, err := mp.Read(f.GetMessage().(*message.MessageRaw), true)
			if err != nil {
				return nil, newError(fmt.Sprintf("unable to decode message: %s", err.Error()))
			}

			f.Message = msg
		}
	}

	return f, nil
}
