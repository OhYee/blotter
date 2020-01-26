package register

import (
	"encoding/json"
	"net/http"
	"reflect"
)

// HandleContext context of a api call
type HandleContext struct {
	Request  *http.Request
	Response http.ResponseWriter
}

// RequestArgs get request args
func (context *HandleContext) RequestArgs(args interface{}) {
	query := context.Request.URL.Query()

	t := reflect.TypeOf(args).Elem()
	v := reflect.ValueOf(args).Elem()
	num := v.NumField()
	for i := 0; i < num; i++ {
		fieldType := t.Field(i)
		v.Field(i).Set(reflect.ValueOf(query.Get(fieldType.Tag.Get("json"))))
	}
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

	return
}
