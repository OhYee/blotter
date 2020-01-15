package msg

import (
	"bytes"
	"github.com/OhYee/goutils"
	"testing"
	"time"
)

func TestHeartBeat(t *testing.T) {
	connection := NewConnection(
		"test server",
		[]API{NewAPI(
			"/add",
			"a + b = c",
			map[string]Value{
				"a": NewValue("int", "number a"),
				"b": NewValue("int", "number b"),
			},
			map[string]Value{
				"c": NewValue("int", "sum of a and b"),
			},
		)},
		time.Now().Unix(),
		0,
	)
	b := NewHeartBeatMessage(connection).ToBytes()
	msg, err := NewMessageFromBytes(bytes.NewBuffer(b))
	if err != nil {
		t.Errorf("Error : %v", err)
	} else {
		msg.Handle(func(msg Message) error {
			conn, err := NewConnectionFromBytes(bytes.NewBuffer(msg.Body))
			if err != nil || !goutils.Equal(conn, connection) {
				t.Errorf("Want %v, but got %v %v", connection, conn, err)
			}
			return nil
		}, nil, nil)
	}
}
