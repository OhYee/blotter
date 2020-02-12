package register

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/OhYee/blotter/output"
	"github.com/gorilla/schema"
)

var decoder = func() (decoder *schema.Decoder) {
	decoder = schema.NewDecoder()
	decoder.SetAliasTag("json")
	decoder.IgnoreUnknownKeys(true)
	return
}()

type httpHeader struct {
	key   string
	value string
}

// HandleContext context of a api call
type HandleContext struct {
	Request  *http.Request
	Response http.ResponseWriter
	buf      *bytes.Buffer
	header   []httpHeader
}

// NewHandleContext initial a handle context object
func NewHandleContext(req *http.Request, rep http.ResponseWriter) *HandleContext {
	return &HandleContext{
		Request:  req,
		Response: rep,
		buf:      bytes.NewBuffer([]byte{}),
		header:   make([]httpHeader, 0),
	}
}

// RequestArgs get request args
func (context *HandleContext) RequestArgs(args interface{}) {
	query := context.Forms()

	err := decoder.Decode(args, query)
	output.Debug("query %+v args: %+v err %+v", query, args, err)
}

func (context *HandleContext) Forms() url.Values {
	context.Request.ParseForm()
	return context.Request.Form
}

func (context *HandleContext) SetCookie(key string, value string) {
	cookie := http.Cookie{Name: key, Value: value, Path: "/"}
	cookie.Expires.After(time.Now().Add(time.Hour * 24 * 7))
	http.SetCookie(context.Response, &cookie)
}

func (context *HandleContext) GetCookie(key string) (value string) {
	if cookie, err := context.Request.Cookie(key); err == nil {
		value = cookie.Value
	}
	return
}

func (context *HandleContext) DeleteCookie(key string) {
	cookie := http.Cookie{Name: key, Value: "", MaxAge: -1}
	cookie.Expires.After(time.Now().Add(time.Second * 1))
	http.SetCookie(context.Response, &cookie)
}

// ReturnJSON return json data
func (context *HandleContext) ReturnJSON(data interface{}) (err error) {
	var b []byte
	if b, err = json.Marshal(data); err != nil {
		return
	}

	context.AddHeader("Content-Type", "application/json")
	context.Write(b...)
	return
}

func (context *HandleContext) Write(b ...byte) {
	context.buf.Write(b)
}

// AddHeader add a header in response
func (context *HandleContext) AddHeader(key string, value string) {
	context.header = append(context.header, httpHeader{key, value})
}

func (context *HandleContext) writeHeaderWithCode(code int) {
	for _, header := range context.header {
		context.Response.Header().Add(header.key, header.value)
	}
	context.Response.WriteHeader(code)
}

// Success return 200 success
func (context *HandleContext) Success() {
	context.writeHeaderWithCode(200)
	context.Response.Write(context.buf.Bytes())
}

// PageNotFound return 404 page not found error
func (context *HandleContext) PageNotFound() {
	output.ErrOutput.Printf("404 Page not Found: %s\n", context.Request.RequestURI)
	context.writeHeaderWithCode(404)
	context.Response.Write([]byte(fmt.Sprintf("Page not found %s", context.Request.RequestURI)))
}

// Forbidden return 403 Forbidden
func (context *HandleContext) Forbidden() {
	output.ErrOutput.Printf("403 Forbidden: %s\n", context.Request.RequestURI)
	context.writeHeaderWithCode(403)
	context.Response.Write([]byte(fmt.Sprintf("Forbidden %s", context.Request.RequestURI)))
}

// NotImplemented return 501 Not Implemented
func (context *HandleContext) NotImplemented() {
	output.ErrOutput.Printf("501 Page not Found: %s\n", context.Request.RequestURI)
	context.writeHeaderWithCode(501)
	context.Response.Write([]byte(fmt.Sprintf("Can not solve request %s", context.Request.RequestURI)))
}

// ServerError return 500 server error
func (context *HandleContext) ServerError(err error) {
	output.ErrOutput.Printf("500 Server Error: %s\n", err.Error())
	context.writeHeaderWithCode(500)
	context.Response.Write([]byte(fmt.Sprintf("Server error %s", err.Error())))
}
