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
	case MessageTypeRequest:
		var req *Request
		if req, err = NewRequestFromBytes(r); err != nil {
			return
		}
		msg = req.ToMessage()
	case MessageTypeResponse:
		var rep *Response
		if rep, err = NewResponseFromBytes(r); err != nil {
			return
		}
		msg = rep.ToMessage()
	case MessageTypeClose:
		var close *Close
		if close, err = NewCloseFromBytes(r); err != nil {
			return
		}
		msg = close.ToMessage()
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
	requestHandle MessageTypeRequestHandle,
	responseHandle MessageTypeResponseHandle,
	closeHandle MessageTypeCloseHandle,
) error {
	switch msg.Type {
	case MessageTypeRequest:
		if requestHandle != nil {
			return requestHandle(msg.Body.(*Request))
		}
	case MessageTypeResponse:
		if responseHandle != nil {
			return responseHandle(msg.Body.(*Response))
		}
	case MessageTypeClose:
		if closeHandle != nil {
			return closeHandle(msg.Body.(*Close))
		}
	default:
		return errors.New("Unknow type message: %v", msg.Type)
	}
	return nil
}
