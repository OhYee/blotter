package proto

import (
	"bytes"
	"encoding/json"
	gb "github.com/OhYee/goutils/bytes"
	"io"
)

type any = interface{}

// Request message body struct
type Request struct {
	URL       string
	Arguments map[string]any
}

// NewRequest initial a request
func NewRequest(url string, args map[string]any) Request {
	return Request{
		URL:       url,
		Arguments: args,
	}
}

// NewRequestFromBytes initial a request from []byte
func NewRequestFromBytes(r io.Reader) (req Request, err error) {
	var args map[string]any
	var urlBytes, argsBytes []byte

	if urlBytes, err = gb.ReadWithLength32(r); err != nil {
		return
	}
	if argsBytes, err = gb.ReadWithLength32(r); err != nil {
		return
	}
	if json.Unmarshal(argsBytes, &args) != nil {
		return
	}
	req = NewRequest(string(urlBytes), args)
	return
}

// ToBytes transfer Request to []byte
func (req Request) ToBytes() []byte {
	buf := bytes.NewBuffer([]byte{})
	gb.WriteWithLength32(buf, gb.FromString(req.URL))
	args, err := json.Marshal(req.Arguments)
	if err != nil {
		args = []byte{}
	}
	gb.WriteWithLength32(buf, args)
	return buf.Bytes()
}
