package main

import (
	"github.com/OhYee/blotter/utils"
	"github.com/OhYee/blotter/utils/protocol"
	"github.com/xtaci/kcp-go"
	"net"
	"testing"
	"time"
)

func Test_Client(t *testing.T) {
	laddr := "127.0.0.1:55555"

	go func() {
		if err := utils.RunServer(laddr, 10, handler); err != nil {
			t.Errorf("%v\n", err)
			t.FailNow()
		}
	}()

	var err error
	var conn net.Conn

	conn, err = kcp.Dial(laddr)
	conn.SetDeadline(time.Now().Add(5 * time.Second))
	if err != nil {
		t.Errorf("%v\n", err)
	}
	if err = proto.SendHandshake(conn); err != nil {
		t.Errorf("%v\n", err)
	}
	if valid := proto.VarifyHandshake(conn); valid == false {
		t.Errorf("Client handshark error\n")
	}
	conn.Close()

	conn, err = kcp.Dial(laddr)
	conn.SetDeadline(time.Now().Add(5 * time.Second))
	if err != nil {
		t.Errorf("%v\n", err)
	}
	conn.Write([]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
	if proto.VarifyHandshake(conn) == true {
		t.Errorf("Unexcepted connection\n")
	}
	conn.Close()

}
