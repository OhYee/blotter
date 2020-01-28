package main

import (
	"github.com/OhYee/blotter/api"
	"github.com/OhYee/blotter/output"
	"github.com/OhYee/blotter/register"
	"net/http"
)

const (
	addr = "127.0.0.1:50000"
)

// Handle of blotter
type Handle struct {
}

func (handle Handle) ServeHTTP(rep http.ResponseWriter, req *http.Request) {
	output.Debug("connection begin")

	context := register.NewHandleContext(req, rep)

	// CORS
	context.AddHeader("Access-Control-Allow-Origin", "*")

	path := req.URL.Path
	if hasPrefix(path, "/api/") {
		err := register.Call(path[5:], context)
		if err != nil {
			context.ServerError(err)
		} else {
			context.Success()
		}
	} else {
		context.NotImplemented()
	}
	output.Debug("connection end")
}

func hasPrefix(s string, prefix string) bool {
	if len(s) < len(prefix) {
		return false
	}
	return s[0:len(prefix)] == prefix
}

func main() {
	api.Register()
	output.Log("Server will start at http://%s", addr)
	if err := http.ListenAndServe(addr, Handle{}); err != nil {
		output.Err(err)
	}
}
