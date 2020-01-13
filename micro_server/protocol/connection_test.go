package proto

import (
	"bytes"
	"github.com/OhYee/goutils"
	"testing"
	"time"
)

func TestConnect(t *testing.T) {
	conn := NewConnection("test server", []API{
		NewAPI("/add", "sum of a and b", map[string]Value{
			"a": NewValue("int", "number a"),
			"b": NewValue("int", "number b"),
		}, map[string]Value{
			"sum": NewValue("int", "sum of a and b"),
		}),
		NewAPI("/info", "get infomation", map[string]Value{}, map[string]Value{
			"data": NewValue("string", "information data"),
		}),
	}, time.Now().Unix(), time.Now().Unix())
	b := conn.ToBytes()

	buf := bytes.NewBuffer(b)

	conn2, err := NewConnectionFromBytes(buf)
	if err != nil {
		t.Errorf("got error: %v", err)
	} else if !goutils.Equal(conn, conn2) {
		t.Errorf("got %v and %v", conn, conn2)
	}

}
