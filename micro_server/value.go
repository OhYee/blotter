package ms

import (
	"github.com/OhYee/goutils/bytes"
	"io"
)

// Value of arguments
type Value struct {
	Name        string // Name of the value
	Type        string // Type of the value
	Description string // Description of the value
}

// NewValue initial a value data
func NewValue(n string, t string, d string) Value {
	return Value{
		Name:        n,
		Type:        t,
		Description: d,
	}
}

// NewValueFromBytes initial a value data from bytes
func NewValueFromBytes(r io.Reader) (v Value, err error) {
	var n, t, d string
	if n, err = bytes.ReadStringWithLength32(r); err != nil {
		return
	}
	if t, err = bytes.ReadStringWithLength32(r); err != nil {
		return
	}
	if d, err = bytes.ReadStringWithLength32(r); err != nil {
		return
	}
	v = NewValue(n, t, d)
	return
}

// ToBytes transfer Value to []byte
func (v *Value) ToBytes() []byte {
	buf := bytes.NewBuffer()
	buf.Write(bytes.FromStringWithLength32(v.Name))
	buf.Write(bytes.FromStringWithLength32(v.Type))
	buf.Write(bytes.FromStringWithLength32(v.Description))
	return buf.Bytes()
}

func readValueSlice(r io.Reader) (ret []Value, err error) {
	var size uint32
	if size, err = bytes.ReadUint32(r); err != nil {
		return
	}

	ret = make([]Value, size)
	for i := 0; i < int(size); i++ {
		var value Value
		if value, err = NewValueFromBytes(r); err != nil {
			return
		}
		ret[i] = value
	}

	return
}

func writeValueSlice(w io.Writer, values []Value) {
	w.Write(bytes.FromInt32(int32(len(values))))
	for _, value := range values {
		w.Write(value.ToBytes())
	}
}
