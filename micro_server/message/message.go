package msg

import (
	"github.com/OhYee/goutils/bytes"
	"github.com/OhYee/rainbow/errors"
	"io"
	// "time"
)

//go:generate gcg ./template/data.json

// Message type of the protocol
type Message struct {
	Type MessageType
	Body bytes.Serializable
}

// NewMessage initial a Message
func NewMessage(t MessageType, body bytes.Serializable) *Message {
	return &Message{
		Type: t,
		Body: body,
	}
}

// NewMessageFromBytes initial a Message from []byte
func NewMessageFromBytes(r io.Reader) (msg *Message, err error) {
	var t uint8
	var body bytes.Serializable

	if t, err = bytes.ReadUint8(r); err != nil {
		return
	}

	switch MessageType(t) {
	case MessageTypeHeartBeat:
		body, err = NewHeartBeatFromBytes(r)
	default:
		err = errors.New("Unknow type message: %v", msg.Type)
	}

	if err == nil {
		msg = NewMessage(MessageType(t), body)
	}

	return
}

// ToBytes transfer Message to []byte
func (msg *Message) ToBytes() []byte {
	buf := bytes.NewBuffer()
	buf.Write(bytes.FromUint8(uint8(msg.Type)))
	buf.Write(bytes.FromBytesWithLength32(msg.Body.ToBytes()))
	return buf.Bytes()
}

// Handle of different message type
func (msg *Message) Handle(
	heartBeatHandle MessageTypeHeartBeatHandle,
	requestHandle MessageTypeRequestHandle,
	responseHandle MessageTypeResponseHandle,
	setHandle MessageTypeSetHandle,
	closeHandle MessageTypeCloseHandle,
) error {
	switch msg.Type {
	case MessageTypeHeartBeat:
		if heartBeatHandle != nil {
			return heartBeatHandle(msg.Body.(*HeartBeat))
		}
	default:
		return errors.New("Unknow type message: %v", msg.Type)
	}
	return nil
}
