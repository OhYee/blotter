package main

import (
	"fmt"
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

	// CORS
	rep.Header().Set("Access-Control-Allow-Origin", "*")

	if hasPrefix(req.RequestURI, "/api/") {
		err := register.Call(req.RequestURI[5:], req, rep)
		if err != nil {
			output.Err(err)
			PageNotFound(rep, req)
		}
	} else {
		PageNotFound(rep, req)
	}
	output.Debug("connection end")
}

func hasPrefix(s string, prefix string) bool {
	if len(s) < len(prefix) {
		return false
	}
	return s[0:len(prefix)] == prefix
}

// PageNotFound return 404 page
func PageNotFound(rep http.ResponseWriter, req *http.Request) {
	output.Log("404 Page not found: %s", req.RequestURI)
	rep.WriteHeader(404)
	rep.Write([]byte(fmt.Sprintf("Page not found %s", req.RequestURI)))
}

// ServerError return 404 page
func ServerError(rep http.ResponseWriter, req *http.Request) {
	rep.WriteHeader(404)
	rep.Write([]byte(fmt.Sprintf("Page not found %s", req.RequestURI)))
}

func main() {
	api.Register()
	output.Log("Server will start at %s", addr)
	if err := http.ListenAndServe(addr, Handle{}); err != nil {
		output.Err(err)
	}
}
