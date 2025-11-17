package frame

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadWriter(t *testing.T) {
	var buf bytes.Buffer
	_, err := NewReadWriter(ReadWriterConf{
		ReadWriter:  &buf,
		OutSystemID: 1,
	})
	require.NoError(t, err)
}

func TestReadWriterErrors(t *testing.T) {
	_, err := NewReadWriter(ReadWriterConf{
		ReadWriter:         nil,
		DialectRW:          nil,
		InKey:              nil,
		OutSystemID:        1,
		OutComponentID:     0,
		OutSignatureLinkID: 0,
		OutKey:             nil,
	})
	require.EqualError(t, err, "Reader not provided")
}
