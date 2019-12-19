package proto

import (
	"github.com/OhYee/goutils"
	"github.com/OhYee/goutils/bytes"
	"io"
)

// Handshake data
var Handshake = []byte{
	0xf6, 0xa6, 0x8f, 0x22, 0xa9, 0xd3, 0x0d, 0x0c,
	0x38, 0xeb, 0xe7, 0xa8, 0x49, 0x57, 0xd2, 0x09,
	0x15, 0xd2, 0x47, 0x3a, 0x25, 0x2d, 0xf5, 0x5e,
	0xbf, 0xd7, 0xa7, 0x78, 0x1d, 0x85, 0x09, 0xc9,
}

// SendHandshake Send the handshake data to the remote
func SendHandshake(w io.Writer) (err error) {
	_, err = w.Write(Handshake)
	return
}

// VarifyHandshake varify the handshake data
func VarifyHandshake(r io.Reader) (valid bool, err error) {
	var b []byte
	if b, err = bytes.ReadNBytes(r, 32); err != nil {
		return
	}
	valid = goutils.Equal(b, Handshake)
	return
}
