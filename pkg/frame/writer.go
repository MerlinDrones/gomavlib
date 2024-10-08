package frame

import (
	"fmt"
	"io"
	"time"

	"github.com/merlindrones/gomavlib/pkg/dialect"
	"github.com/merlindrones/gomavlib/pkg/message"
)

// WriterConf is the configuration of a Writer.
type WriterConf struct {
	// the underlying bytes writer.
	Writer io.Writer

	// (optional) the dialect which contains the messages that will be written.
	DialectRW *dialect.ReadWriter

	// Mavlink version used to encode messages.
	//OutVersion WriterOutVersion
	// the system id, added to every outgoing frame and used to identify this
	// node in the network.
	OutSystemID byte
	// (optional) the component id, added to every outgoing frame, defaults to 1.
	OutComponentID byte
	// (optional) the value to insert into the signature link id.
	// This feature requires v2 frames.
	OutSignatureLinkID byte
	// (optional) the secret key used to sign outgoing frames.
	// This feature requires v2 frames.
	OutKey *V2Key
}

// Writer is a Frame writer.
type Writer struct {
	conf                   WriterConf
	bw                     []byte
	curWriteSequenceNumber byte
}

// NewWriter allocates a Writer.
func NewWriter(conf WriterConf) (*Writer, error) {
	if conf.Writer == nil {
		return nil, fmt.Errorf("Writer not provided")
	}

	/*
		if conf.OutVersion == 0 {
			return nil, fmt.Errorf("OutVersion not provided")
		}
	*/
	if conf.OutSystemID < 1 {
		return nil, fmt.Errorf("OutSystemID must be greater than one")
	}
	if conf.OutComponentID < 1 {
		conf.OutComponentID = 1
	}

	/*
		if conf.OutKey != nil && conf.OutVersion != V2 {
			return nil, fmt.Errorf("OutKey requires V2 frames")
		}
	*/

	return &Writer{
		conf: conf,
		bw:   make([]byte, bufferSize),
	}, nil
}

// WriteMessage writes a Message.
// The Message is wrapped into a Frame whose fields are filled automatically.
// It must not be called by multiple routines in parallel.
func (w *Writer) WriteMessage(m message.Message) error {
	frame, err := w.FillFrameWithMessage(m)
	if err != nil {
		return err
	}
	return w.writeFrameInner(frame)
}

func (w *Writer) FillFrameWithMessage(m message.Message) (*V2Frame, error) {
	frame := &V2Frame{Message: m}
	fmt.Printf("Before: %v\n", frame.SystemID)
	err := w.fillFrame(frame)
	if err != nil {
		return nil, err
	}
	fmt.Printf("After: %v\n", frame.SystemID)
	return frame, nil
}

// WriteFrame writes a Frame.
// It must not be called by multiple routines in parallel.
// This function is intended only for routing pre-existing and/or filled frames.
func (w *Writer) WriteFrame(fr *V2Frame) error {
	if fr.GetMessage() == nil {
		return fmt.Errorf("message is nil")
	}

	// encode message if it is not already encoded
	if _, ok := fr.GetMessage().(*message.MessageRaw); !ok {
		if w.conf.DialectRW == nil {
			return fmt.Errorf("dialect is nil")
		}

		mp := w.conf.DialectRW.GetMessage(fr.GetMessage().GetID())
		if mp == nil {
			return fmt.Errorf("message is not in the dialect")
		}

		w.encodeMessageInFrame(fr, mp)
	}

	return w.writeFrameInner(fr)
}

/*
*
PRIVATE
*/
func (w *Writer) fillFrame(fr *V2Frame) (err error) {
	if fr.GetMessage() == nil {
		return fmt.Errorf("message is nil")
	}

	// fill SequenceNumber, SystemID, ComponentID, CompatibilityFlag, IncompatibilityFlag
	fr.SequenceNumber = w.curWriteSequenceNumber
	fr.SystemID = w.conf.OutSystemID
	fr.ComponentID = w.conf.OutComponentID

	fr.CompatibilityFlag = 0
	fr.IncompatibilityFlag = 0
	if w.conf.OutKey != nil {
		fr.IncompatibilityFlag |= V2FlagSigned
	}

	w.curWriteSequenceNumber++

	if w.conf.DialectRW == nil {
		return fmt.Errorf("dialect is nil")
	}

	mp := w.conf.DialectRW.GetMessage(fr.GetMessage().GetID())
	if mp == nil {
		return fmt.Errorf("message is not in the dialect")
	}

	// encode message if it is not already encoded
	if _, ok := fr.GetMessage().(*message.MessageRaw); !ok {
		w.encodeMessageInFrame(fr, mp)
	}

	// fill checksum
	fr.Checksum = fr.GenerateChecksum(mp.CRCExtra())

	// fill SignatureLinkID, SignatureTimestamp, Signature if v2
	if w.conf.OutKey != nil {
		fr.SignatureLinkID = w.conf.OutSignatureLinkID
		// Timestamp in 10 microsecond units since 1st January 2015 GMT time
		fr.SignatureTimestamp = uint64(time.Since(signatureReferenceDate)) / 10000
		fr.Signature = fr.GenerateSignature(w.conf.OutKey)
	}

	return nil
}

func (w *Writer) encodeMessageInFrame(fr Frame, mp *message.ReadWriter) {
	msgRaw := mp.Write(fr.GetMessage(), true)

	switch ff := fr.(type) {
	case *V2Frame:
		ff.Message = msgRaw
	}
}

func (w *Writer) writeFrameInner(fr *V2Frame) error {
	n, err := fr.EncodeTo(w.bw, fr.GetMessage().(*message.MessageRaw).Payload)
	if err != nil {
		return err
	}

	// do not check n, since io.Writer is not allowed to return n < len(buf)
	// without throwing an error
	_, err = w.conf.Writer.Write(w.bw[:n])
	return err
}
