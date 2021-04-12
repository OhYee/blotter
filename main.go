package main

import (
	"flag"

	"github.com/OhYee/blotter/api"
	"github.com/OhYee/blotter/http"
	"github.com/OhYee/blotter/output"
	"github.com/OhYee/blotter/register"
)

var (
	addr   = "127.0.0.1:50000"
	prefix = "/api/"
)

var (
	_version string
)

//go:generate /bin/bash ./generate.bash

func parseFlags() {
	flag.StringVar(&addr, "address", "127.0.0.1:50000", "listen address")
	flag.StringVar(&prefix, "prefix", "/api/", "api url prefix")

	flag.Parse()
}

func main() {
	parseFlags()

	register.SetContext("version", _version)

	api.Register()
	// queue.Register().Register("extensions/queue")
	// register.DebugApiMap()
	output.Log("Server will start at http://%s", addr)
	if err := http.Server(addr, prefix); err != nil {
		output.Err(err)
	}
}
