package main

import (
	"github.com/OhYee/blotter/utils"
	"github.com/OhYee/rainbow/color"
	"github.com/OhYee/rainbow/log"
	"net"
)

var (
	debugLogger = log.New()
	errLogger   = log.New().SetColor(color.New().SetFrontRed())
)

func handler(id int, listener net.Listener) {
	debugLogger.Printf("Start thread %d\n", id)
	for {
		conn, err := listener.Accept()
		debugLogger.Printf("Thread %d receive a connection %v from %v\n", id, &conn, conn.RemoteAddr())
		if err != nil {
			errLogger.Println(err)
		} else {

		}
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
	if err = utils.RunServer(laddr, handler); err != nil {
		errLogger.Println(err)
	}
}
