package register

import (
	"net/http"

	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"time"

	"github.com/OhYee/blotter/output"
)

type httpHeader struct {
	key   string
	value string
}

// HTTPContext context of a api call
type HTTPContext struct {
	Request  *http.Request
	Response http.ResponseWriter
	buf      *bytes.Buffer
	header   []httpHeader
}

func assertHTTP() HandleContext {
	var a HandleContext = &HTTPContext{}
	return a
}

// NewHTTPContext initial a handle context object
func NewHTTPContext(req *http.Request, rep http.ResponseWriter) *HTTPContext {
	return &HTTPContext{
		Request:  req,
		Response: rep,
		buf:      bytes.NewBuffer([]byte{}),
		header:   make([]httpHeader, 0),
	}
}

func (context *HTTPContext) GetRequest() *http.Request {
	return context.Request
}

func (context *HTTPContext) GetResponse() http.ResponseWriter {
	return context.Response
}

// RequestArgs get request args of forms
func (context *HTTPContext) RequestArgs(args interface{}, opts ...string) {
	if len(opts) > 0 && (opts[0] == "post" || opts[0] == "POST") {
		b, err := ioutil.ReadAll(context.Request.Body)
		if err != nil {
			output.Err(err)
			return
		}
		if err = json.Unmarshal(b, args); err != nil {
			output.Err(err)
			return
		}
	} else {
		query := context.Forms()

		err := decoder.Decode(args, query)
		if err != nil {
			output.Err(err)
		}
	}

}

func (context *HTTPContext) Forms() url.Values {
	context.Request.ParseForm()
	return context.Request.Form
}

func (context *HTTPContext) SetCookie(key string, value string) {
	cookie := http.Cookie{Name: key, Value: value, Path: "/"}
	cookie.Expires.After(time.Now().Add(time.Hour * 24 * 7))
	http.SetCookie(context.Response, &cookie)
}

func (context *HTTPContext) GetCookie(key string) (value string) {
	if cookie, err := context.Request.Cookie(key); err == nil {
		value = cookie.Value
	}
	return
}

func (context *HTTPContext) DeleteCookie(key string) {
	cookie := http.Cookie{Name: key, Value: "", MaxAge: -1}
	cookie.Expires.After(time.Now().Add(time.Second * 1))
	http.SetCookie(context.Response, &cookie)
}

// ReturnJSON return json data
func (context *HTTPContext) ReturnJSON(data interface{}) (err error) {
	var b []byte
	if b, err = json.Marshal(data); err != nil {
		return
	}

	context.AddHeader("Content-Type", "application/json")
	context.Write(b...)
	return
}

// ReturnText return json data
func (context *HTTPContext) ReturnText(data string) (err error) {
	context.AddHeader("Content-Type", "text/plain")
	context.Write([]byte(data)...)
	return
}

// ReturnXML return json data
func (context *HTTPContext) ReturnXML(data string) (err error) {
	context.AddHeader("Content-Type", "application/xml")
	context.Write([]byte(data)...)
	return
}
func (context *HTTPContext) Write(b ...byte) (err error) {
	_, err = context.buf.Write(b)
	return
}

// AddHeader add a header in response
func (context *HTTPContext) AddHeader(key string, value string) {
	context.header = append(context.header, httpHeader{key, value})
}

func (context *HTTPContext) writeHeaderWithCode(code int) {
	for _, header := range context.header {
		context.Response.Header().Add(header.key, header.value)
	}
	context.Response.WriteHeader(code)
}

// Success return 200 success
func (context *HTTPContext) Success() {
	context.writeHeaderWithCode(200)
	context.Response.Write(context.buf.Bytes())
}

// PageNotFound return 404 page not found error
func (context *HTTPContext) PageNotFound() {
	output.ErrOutput.Printf("404 Page not Found: %s\n", context.Request.RequestURI)
	context.writeHeaderWithCode(404)
	context.Response.Write([]byte(fmt.Sprintf("Page not found %s", context.Request.RequestURI)))
}

// Forbidden return 403 Forbidden
func (context *HTTPContext) Forbidden() {
	output.ErrOutput.Printf("403 Forbidden: %s\n", context.Request.RequestURI)
	context.writeHeaderWithCode(403)
	context.Response.Write([]byte(fmt.Sprintf("Forbidden %s", context.Request.RequestURI)))
}

// NotImplemented return 501 Not Implemented
func (context *HTTPContext) NotImplemented() {
	output.ErrOutput.Printf("501 Page not Found: %s\n", context.Request.RequestURI)
	context.writeHeaderWithCode(501)
	context.Response.Write([]byte(fmt.Sprintf("Can not solve request %s", context.Request.RequestURI)))
}

// ServerError return 500 server error
func (context *HTTPContext) ServerError(err error) {
	output.ErrOutput.Printf("500 Server Error: %s\n", err.Error())
	context.writeHeaderWithCode(500)
	context.Response.Write([]byte(fmt.Sprintf("Server error %s", err.Error())))
}

// PermanentlyMoved to url (301)
func (context *HTTPContext) PermanentlyMoved(url string) {
	http.Redirect(context.Response, context.Request, url, 301)
}

// TemporarilyMoved to url (302)
func (context *HTTPContext) TemporarilyMoved(url string) {
	http.Redirect(context.Response, context.Request, url, 302)
}
