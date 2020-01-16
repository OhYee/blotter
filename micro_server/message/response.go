package msg

import (
	"github.com/OhYee/goutils/bytes"
	"io"
)

// Response message body
type Response struct {
	Arguments []byte
}

// NewResponse initial a Response
func NewResponse(args []byte) *Response {
	return &Response{
		Arguments: args,
	}
}

// NewResponseFromBytes initial a Response from []byte
func NewResponseFromBytes(r io.Reader) (rep *Response, err error) {
	var args []byte
	if args, err = bytes.ReadBytesWithLength32(r); err != nil {
		return
	}
	rep = NewResponse(args)
	return
}

// ToBytes transfer Response to []byte
func (rep *Response) ToBytes() []byte {
	buf := bytes.NewBuffer()
	buf.Write(bytes.FromBytesWithLength32(rep.Arguments))
	return buf.Bytes()
}

// ToMessage initial a Response message
func (rep *Response) ToMessage() *Message {
	return NewMessage(MessageTypeResponse, rep)
}
