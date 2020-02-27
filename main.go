package main

import (
	"github.com/OhYee/blotter/api"
	"github.com/OhYee/blotter/http"
	"github.com/OhYee/blotter/output"
)

const (
	addr   = "127.0.0.1:50000"
	prefix = "/api/"
)

func main() {
	api.Register()
	output.Log("Server will start at http://%s", addr)
	if err := http.Server(addr, prefix); err != nil {
		output.Err(err)
	}
}
