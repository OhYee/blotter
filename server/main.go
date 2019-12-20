package main

import (
	"time"
	"github.com/OhYee/blotter/utils"
	"github.com/OhYee/blotter/utils/protocol"
	"github.com/OhYee/rainbow/color"
	"github.com/OhYee/rainbow/errors"
	"github.com/OhYee/rainbow/log"
	"net"
)

var (
	debugLogger = log.New()
	errLogger   = log.New().SetColor(color.New().SetFrontRed())
)

func server(conn net.Conn) (err error) {
	conn.SetDeadline(time.Now().Add(5 * time.Second))

	if !proto.VarifyHandshake(conn) {
		return errors.New("Handshake error")
	}
	if err = proto.SendHandshake(conn); err != nil {
		return err
	}

	return nil
}

func handler(id int, listener net.Listener) {
	debugLogger.Printf("Start thread %d\n", id)
	for {
		conn, err := listener.Accept()
		debugLogger.Printf("Thread %d receive a connection %v from %v\n", id, &conn, conn.RemoteAddr())
		if err != nil {
			errLogger.Println(err)
		} else {
			if err = server(conn); err != nil {
				errLogger.Println(err)
			}
		}
		debugLogger.Printf("Thread %d connection %v from %v closed\n", id, &conn, conn.RemoteAddr())
	}
}

func main() {
	laddr, verbose, err := utils.ParseFlag()
	if err != nil {
		errLogger.Println(err)
		return
	}
	if !verbose {
		debugLogger.SetOutputToNil()
	}
	debugLogger.Printf("Start server at %s\n", laddr)
	if err = utils.RunServer(laddr, 10, handler); err != nil {
		errLogger.Println(err)
	}
}
