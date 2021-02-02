package main

import (
	"github.com/OhYee/blotter/api"
	"github.com/OhYee/blotter/http"
	"github.com/OhYee/blotter/output"
	"github.com/OhYee/blotter/register"
)

const (
	addr   = "127.0.0.1:50000"
	prefix = "/api/"
)

var (
	_version string
)

//go:generate /bin/bash ./build.bash
func main() {
	register.SetContext("version", _version)

	api.Register()
	// queue.Register().Register("extensions/queue")
	// register.DebugApiMap()
	output.Log("Server will start at http://%s", addr)
	if err := http.Server(addr, prefix); err != nil {
		output.Err(err)
	}
}
