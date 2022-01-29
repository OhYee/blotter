package http

import (
	"net/http"
	"strings"

	"github.com/OhYee/blotter/register"
)

// Handle of blotter
type Handle struct {
	Prefix string
}

func (handle Handle) ServeHTTP(rep http.ResponseWriter, req *http.Request) {
	context := register.NewHTTPContext(req, rep)

	// CORS
	context.AddHeader("Access-Control-Allow-Origin", req.Header.Get("Origin"))
	context.AddHeader("Access-Control-Allow-Headers", "*")
	context.AddHeader("Access-Control-Allow-Methods", "POST, GET, OPTIONS, DELETE")
	context.AddHeader("Access-Control-Allow-Credentials", "true")
	if context.Request.Method == "OPTIONS" {
		context.Success()
		return
	}

	path := req.URL.Path
	if strings.HasPrefix(path, handle.Prefix) {
		err := register.Call(path[len(handle.Prefix):], context)
		if err != nil {
			context.ServerError(err)
		} else {
			context.Success()
		}
	} else {
		context.NotImplemented()
	}
}

// Server start http server
func Server(addr string, prefix string) error {
	return http.ListenAndServe(addr, Handle{Prefix: prefix})
}
