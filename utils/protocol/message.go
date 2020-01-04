package proto

import (
	"bytes"
	gb "github.com/OhYee/goutils/bytes"
	"github.com/OhYee/rainbow/errors"
	"io"
	"time"
)

// Message type of the protocol
type Message struct {
	Type MessageType
	Body []byte
}

// MessageHandle handle of different type messages
type MessageHandle func(msg Message) error

// NewMessage initial a Message
func NewMessage(t MessageType, body []byte) Message {
	return Message{
		Type: t,
		Body: body,
	}
}

// NewMessageFromBytes initial a Message from []byte
func NewMessageFromBytes(r io.Reader) (msg Message, err error) {
	var t uint8
	var b []byte
	if t, err = gb.ReadUint8(r); err != nil {
		return
	}
	if b, err = gb.ReadWithLength32(r); err != nil {
		return
	}

	msg = NewMessage(MessageType(t), b)
	return
}

// NewHeartBeatMessage initial a HeartBeat message
func NewHeartBeatMessage(conn Connection) Message {
	conn.SendTime = time.Now().Unix()
	conn.RecvTime = 0
	return NewMessage(MessageTypeHeartBeat, conn.ToBytes())
}

// ToBytes transfer Message to []byte
func (msg Message) ToBytes() []byte {
	buf := bytes.NewBuffer([]byte{})
	buf.Write(gb.FromUint8(uint8(msg.Type)))
	gb.WriteWithLength32(buf, msg.Body)
	return buf.Bytes()
}

// Handle of different message type
func (msg Message) Handle(heartBeatHandle MessageHandle,
	requestHandle MessageHandle, responseHandle MessageHandle) error {
	switch msg.Type {
	case MessageTypeHeartBeat:
		if heartBeatHandle != nil {
			return heartBeatHandle(msg)
		}
	case MessageTypeRequest:
		if requestHandle != nil {
			return requestHandle(msg)
		}
	case MessageTypeResponse:
		if responseHandle != nil {
			return responseHandle(msg)
		}
	default:
		return errors.New("Unknow type message: %v", msg.Type)
	}
	return nil
}
