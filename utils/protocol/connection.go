package proto

import (
	"bytes"
	gb "github.com/OhYee/goutils/bytes"
	"time"
)

// Connection from a sub-server
type Connection struct {
	Description string    // Description of the connection
	APIList     []API     // APIList list of the sub-server api
	KeepAlive   time.Time // KeepAlive the last connect time of the sub-server
}

// NewConnection initial a Connection data
func NewConnection(d string, apis []API) Connection {
	return Connection{
		Description: d,
		APIList:     apis,
		KeepAlive:   time.Now(),
	}
}

// API of the sub-server
type API struct {
	URL         string           // URL of the api
	Description string           // Description of the api
	Input       map[string]Value // Input arguments types
	Output      map[string]Value // Output data types
}

// NewAPI initial a API data
func NewAPI(url string, d string, in map[string]Value, out map[string]Value) API {
	return API{
		URL:         url,
		Description: d,
		Input:       in,
		Output:      out,
	}
}

// Value of the input/output
type Value struct {
	Type        string // Type of the value
	Description string // Description of the value
}

// NewValue initial a value data
func NewValue(t string, d string) Value {
	return Value{
		Type:        t,
		Description: d,
	}
}

// NewValueFromBytes initial a value data from bytes
func NewValueFromBytes(b []byte) (v Value, err error) {
	buf := bytes.NewBuffer(b)
	t, err := gb.ReadWithLength32(buf)
	if err != nil {
		return
	}
	d, err := gb.ReadWithLength32(buf)
	if err != nil {
		return
	}
	v = NewValue(string(t), string(d))
	return
}

// ToBytes transfer Value to []byte
func (v *Value) ToBytes() []byte {
	buf := bytes.NewBuffer([]byte{})
	gb.WriteWithLength32(buf, gb.FromString(v.Type))
	gb.WriteWithLength32(buf, gb.FromString(v.Description))
	return buf.Bytes()
}
