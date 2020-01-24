package register

import (
	"encoding/json"
	"github.com/OhYee/blotter/output"
	"net/http"
)

// HandleContext context of a api call
type HandleContext struct {
	Request  *http.Request
	Response http.ResponseWriter
}

// RequestArgs get request args
func (context *HandleContext) RequestArgs(args interface{}) {
	// json.Unmarshal(context.Request.)
}

// ReturnJSON return json data
func (context *HandleContext) ReturnJSON(data interface{}) (err error) {
	var b []byte
	if b, err = json.Marshal(data); err != nil {
		return
	}
	
	// set header first, then write status code, finially write body
	context.Response.Header().Add("Content-Type", "application/json")
	context.Response.WriteHeader(200)
	context.Response.Write(b)

	output.Debug("Write json %+v", b)
	return
}
