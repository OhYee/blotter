package msg

import (
	"github.com/OhYee/goutils/bytes"
	"io"
)

// Request message body
type Request struct {
	Name      string
	Arguments []byte
}

// NewRequest initial a Request
func NewRequest(name string, args []byte) *Request {
	return &Request{
		Name:      name,
		Arguments: args,
	}
}

// NewRequestFromBytes initial a Request from []byte
func NewRequestFromBytes(r io.Reader) (req *Request, err error) {
	var name string
	var args []byte
	if name, err = bytes.ReadStringWithLength32(r); err != nil {
		return
	}
	if args, err = bytes.ReadBytesWithLength32(r); err != nil {
		return
	}
	req = NewRequest(name, args)
	return
}

// ToBytes transfer Request to []byte
func (req *Request) ToBytes() []byte {
	buf := bytes.NewBuffer()
	buf.Write(bytes.FromStringWithLength32(req.Name))
	buf.Write(bytes.FromBytesWithLength32(req.Arguments))
	return buf.Bytes()
}

// ToMessage initial a Request message
func (req *Request) ToMessage() *Message {
	return NewMessage(MessageTypeRequest, req)
}
