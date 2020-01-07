package ms

import (
	"bytes"
	gb "github.com/OhYee/goutils/bytes"
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
	var n, t, d []byte
	if n, err = gb.ReadWithLength32(r); err != nil {
		return
	}
	if t, err = gb.ReadWithLength32(r); err != nil {
		return
	}
	if d, err = gb.ReadWithLength32(r); err != nil {
		return
	}
	v = NewValue(string(n), string(t), string(d))
	return
}

// ToBytes transfer Value to []byte
func (v *Value) ToBytes() []byte {
	buf := bytes.NewBuffer([]byte{})
	gb.WriteWithLength32(buf, gb.FromString(v.Name))
	gb.WriteWithLength32(buf, gb.FromString(v.Type))
	gb.WriteWithLength32(buf, gb.FromString(v.Description))
	return buf.Bytes()
}

func readValueSlice(r io.Reader) (ret []Value, err error) {
	var size uint32
	if size, err = gb.ReadUint32(r); err != nil {
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
	w.Write(gb.FromInt32(int32(len(values))))
	for _, value := range values {
		w.Write(value.ToBytes())
	}
}
