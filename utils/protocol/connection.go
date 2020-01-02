package proto

import (
	"bytes"
	gb "github.com/OhYee/goutils/bytes"
	"io"
)

// Connection from a sub-server
type Connection struct {
	Description string // Description of the connection
	APIList     []API  // APIList list of the sub-server api
	SendTime    int64  // Send time of connection
	RecvTime    int64  // Receive time of connection
}

// NewConnection initial a Connection data
func NewConnection(d string, apis []API, sendTime int64, recvTime int64) Connection {
	return Connection{
		Description: d,
		APIList:     apis,
		SendTime:    sendTime,
		RecvTime:    recvTime,
	}
}

// NewConnectionFromBytes initial a API data from bytes
func NewConnectionFromBytes(r io.Reader) (conn Connection, err error) {
	d, err := gb.ReadWithLength32(r)
	apiSize, err := gb.ReadUint32(r)
	if err != nil {
		return
	}

	apiList := make([]API, apiSize)
	for i := uint32(0); i < apiSize; i++ {
		var api API
		if api, err = NewAPIFromBytes(r); err != nil {
			return
		}
		apiList[i] = api
	}

	var sendTime, recvTime int64
	if sendTime, err = gb.ReadInt64(r); err != nil {
		return
	}
	if recvTime, err = gb.ReadInt64(r); err != nil {
		return
	}

	conn = NewConnection(string(d), apiList, sendTime, recvTime)
	return
}

// ToBytes transfer Connection to []byte
func (conn Connection) ToBytes() []byte {
	buf := bytes.NewBuffer([]byte{})
	gb.WriteWithLength32(buf, gb.FromString(conn.Description))

	apiSize := len(conn.APIList)
	buf.Write(gb.FromUint32(uint32(apiSize)))
	for _, api := range conn.APIList {
		buf.Write(api.ToBytes())
	}

	buf.Write(gb.FromInt64(conn.SendTime))
	buf.Write(gb.FromInt64(conn.RecvTime))

	return buf.Bytes()
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

// NewAPIFromBytes initial a API data from bytes
func NewAPIFromBytes(r io.Reader) (api API, err error) {
	var url, description []byte

	url, err = gb.ReadWithLength32(r)
	if err != nil {
		return
	}

	description, err = gb.ReadWithLength32(r)
	if err != nil {
		return
	}

	inSize, err := gb.ReadUint32(r)
	if err != nil {
		return
	}
	in := make(map[string]Value)
	for i := 0; i < int(inSize); i++ {
		var key []byte
		var value Value

		key, err = gb.ReadWithLength32(r)
		if err != nil {
			return
		}
		value, err = NewValueFromBytes(r)
		if err != nil {
			return
		}
		in[string(key)] = value
	}

	outSize, err := gb.ReadUint32(r)
	if err != nil {
		return
	}
	out := make(map[string]Value)
	for i := 0; i < int(outSize); i++ {
		var key []byte
		var value Value

		key, err = gb.ReadWithLength32(r)
		if err != nil {
			return
		}
		value, err = NewValueFromBytes(r)
		if err != nil {
			return
		}
		out[string(key)] = value
	}

	api = NewAPI(string(url), string(description), in, out)
	return
}

// ToBytes transfer API to []byte
func (api API) ToBytes() []byte {
	buf := bytes.NewBuffer([]byte{})
	gb.WriteWithLength32(buf, gb.FromString(api.URL))
	gb.WriteWithLength32(buf, gb.FromString(api.Description))

	buf.Write(gb.FromInt32(int32(len(api.Input))))
	for key, value := range api.Input {
		gb.WriteWithLength32(buf, gb.FromString(key))
		buf.Write(value.ToBytes())
	}

	buf.Write(gb.FromInt32(int32(len(api.Output))))
	for key, value := range api.Output {
		gb.WriteWithLength32(buf, gb.FromString(key))
		buf.Write(value.ToBytes())
	}

	return buf.Bytes()
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
func NewValueFromBytes(r io.Reader) (v Value, err error) {
	t, err := gb.ReadWithLength32(r)
	if err != nil {
		return
	}
	d, err := gb.ReadWithLength32(r)
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
