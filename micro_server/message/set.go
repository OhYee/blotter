package msg

import (
	"encoding/json"
	"github.com/OhYee/goutils/bytes"
	"io"
)

// Set message body
type Set struct {
	Pair map[string]interface{}
}

// NewSet initial a Set
func NewSet(m map[string]interface{}) *Set {
	return &Set{
		Pair: m,
	}
}

// NewSetFromBytes initial a Set from []byte
func NewSetFromBytes(r io.Reader) (s *Set, err error) {
	var b []byte
	if b, err = bytes.ReadBytesWithLength32(r); err != nil {
		return
	}

	m := make(map[string]interface{})
	if err = json.Unmarshal(b, &m); err != nil {
		return
	}

	s = NewSet(m)
	return
}

// ToBytes transfer Set to []byte
func (s *Set) ToBytes() []byte {
	b, _ := json.Marshal(s.Pair)
	return bytes.FromBytesWithLength32(b)
}

// ToMessage initial a Set message
func (s *Set) ToMessage() *Message {
	return NewMessage(MessageTypeSet, s)
}
