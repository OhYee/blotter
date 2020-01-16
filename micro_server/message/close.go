package msg

import ()

// Close message body
type Close struct {
}

// NewClose initial a Close
func NewClose() *Close {
	return new(Close)
}

// NewCloseFromBytes initial a Close from []byte
func NewCloseFromBytes() (close *Close, err error) {
	return new(Close), nil
}

// ToBytes transfer Close to []byte
func (close *Close) ToBytes() []byte {
	return []byte{}
}

// ToMessage initial a Close message
func (close *Close) ToMessage() *Message {
	return NewMessage(MessageTypeClose, close)
}
