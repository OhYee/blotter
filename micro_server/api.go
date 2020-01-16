package ms

import (
	"github.com/OhYee/goutils/bytes"
	"io"
)

// APIInfo information
type APIInfo struct {
	Address     string  // Address of the api
	Description string  // Description of api
	Input       []Value // Input arguments
	Output      []Value // Output argument
}

// NewAPIInfo initial a APIInfo
func NewAPIInfo(address string, d string, in []Value, out []Value) APIInfo {
	return APIInfo{
		Address:     address,
		Description: d,
		Input:       in,
		Output:      out,
	}
}

// NewAPIInfoFromBytes initial a APIInfo data from bytes
func NewAPIInfoFromBytes(r io.Reader) (apiInfo APIInfo, err error) {
	var url, description string

	if url, err = bytes.ReadStringWithLength32(r); err != nil {
		return
	}

	if description, err = bytes.ReadStringWithLength32(r); err != nil {
		return
	}

	var in, out []Value
	if in, err = readValueSlice(r); err != nil {
		return
	}
	if out, err = readValueSlice(r); err != nil {
		return
	}

	apiInfo = NewAPIInfo(string(url), string(description), in, out)
	return
}

// ToBytes transfer API to []byte
func (api APIInfo) ToBytes() []byte {
	buf := bytes.NewBuffer()
	buf.Write(bytes.FromStringWithLength32(api.Address))
	buf.Write(bytes.FromStringWithLength32(api.Description))

	writeValueSlice(buf, api.Input)
	writeValueSlice(buf, api.Output)

	return buf.Bytes()
}
