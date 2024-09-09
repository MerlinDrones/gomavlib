package x25

import (
	"encoding/binary"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestX25Params(t *testing.T) {
	h := New()
	assert.Equal(t, 2, h.Size())
	assert.Equal(t, 1, h.BlockSize())
}

func TestX25Hash(t *testing.T) {
	for _, ca := range []struct {
		name string
		in   []byte
		out  uint16
	}{
		{
			"empty",
			[]byte{},
			0xFFFF,
		},
		{
			"0x01",
			[]byte("\x01"),
			0x1e0e,
		},
		{
			"hello world",
			[]byte("hello world"),
			0x51f9,
		},
		{
			"reference string",
			[]byte("The quick brown fox jumps over the lazy dog"),
			0x6ca7,
		},
	} {
		t.Run(ca.name, func(t *testing.T) {
			h := New()
			h.Write(ca.in)
			out1 := h.Sum16()
			assert.Equal(t, ca.out, out1)
			out2 := binary.LittleEndian.Uint16(h.Sum(nil))
			assert.Equal(t, ca.out, out2)
		})
	}
}
