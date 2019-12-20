package utils

import (
	"flag"
	"fmt"
	"github.com/OhYee/rainbow/errors"
	"github.com/xtaci/kcp-go"
	"net"
	"os"
	"os/signal"
	"regexp"
)

const (
	ipPattern   = `((2(5[0-5]|[0-4]\d))|[0-1]?\d{1,2})(\.((2(5[0-5]|[0-4]\d))|[0-1]?\d{1,2})){3}`
	portPattern = `(6(5(5(3[0-5]|[0-2]\d)|[0-4]\d{1,2})|[0-4]\d{1,3})|[0-5]\d{1,4}|\d{1,4})`
	fullPattern = ipPattern + ":" + portPattern
)

func ParseFlag() (laddr string, verbose bool, err error) {
	defer errors.NewErr(err)

	flag.StringVar(&laddr, "address", "127.0.0.1:12345", "Server address [ip]:[port]")
	flag.StringVar(&laddr, "a", "127.0.0.1:12345", "Server address [ip]:[port]")

	flag.BoolVar(&verbose, "verbose", false, "show debug information")
	flag.BoolVar(&verbose, "v", false, "show debug information")

	flag.Parse()

	var ipRe, portRe, fullRe *regexp.Regexp

	if ipRe, err = regexp.Compile(fmt.Sprintf("^%s$", ipPattern)); err != nil {
		return
	}
	if portRe, err = regexp.Compile(fmt.Sprintf("^%s$", portPattern)); err != nil {
		return
	}
	if fullRe, err = regexp.Compile(fmt.Sprintf("^%s$", fullPattern)); err != nil {
		return
	}
	switch {
	case ipRe.MatchString(laddr):
		laddr = fmt.Sprintf("%s:12345", laddr)
	case portRe.MatchString(laddr):
		laddr = fmt.Sprintf("127.0.0.1:%s", laddr)
	case fullRe.MatchString(laddr):
		break
	default:
		err = fmt.Errorf("Can not parse address %s", laddr)
		return
	}

	return
}

func RunServer(laddr string, threadNum int, handler func(id int, listener net.Listener)) (err error) {
	listener, err := kcp.Listen(laddr)
	if err != nil {
		return err
	}
	for i := 0; i < threadNum; i++ {
		go handler(i, listener)
	}
	c := make(chan os.Signal)
	signal.Notify(c)
	<-c
	return nil
}
