package register

import (
	"fmt"
	"net/http"

	"github.com/OhYee/blotter/api/pkg/user"
	"github.com/gorilla/websocket"
	"github.com/mitchellh/mapstructure"
)

// WebSocketContext context of a api call
type WebSocketContext struct {
	Request   *http.Request
	Response  http.ResponseWriter
	ws        *websocket.Conn
	funcName  string
	args      interface{}
	Successed bool
	Data      interface{}
	Error     error
}

func assertWS() HandleContext {
	var a HandleContext = &WebSocketContext{}
	return a
}

func NewWebSocketContext(req *http.Request, rep http.ResponseWriter, ws *websocket.Conn, funcName string, args interface{}) (context *WebSocketContext) {
	context = &WebSocketContext{
		Request:  req,
		Response: rep,
		ws:       ws,
		funcName: funcName,
		args:     args,
	}
	return
}

// GetRequest implements HandleContext.GetRequest
func (context *WebSocketContext) GetRequest() *http.Request {
	return context.Request
}

// GetResponse implements HandleContext.GetResponse
func (context *WebSocketContext) GetResponse() http.ResponseWriter {
	return context.Response
}

// RequestArgs get request args of forms
func (context *WebSocketContext) RequestArgs(args interface{}, opts ...string) {
	mapstructure.Decode(context.args, args)
}

// ReturnJSON return json data
func (context *WebSocketContext) ReturnJSON(data interface{}) (err error) {
	context.Success()
	context.Data = data
	return
}

// ReturnText return json data
func (context *WebSocketContext) ReturnText(data string) (err error) {
	context.Success()
	context.Data = data
	return
}

// ReturnXML return json data
func (context *WebSocketContext) ReturnXML(data string) (err error) {
	context.Success()
	context.Data = data
	return
}
func (context *WebSocketContext) Write(b ...byte) (err error) {
	context.Success()
	context.Data = b
	return
}

// Success return 200 success
func (context *WebSocketContext) Success() {
	context.Successed = true
}

// PageNotFound return 404 page not found error
func (context *WebSocketContext) PageNotFound() {
	context.Successed = false
	context.Error = fmt.Errorf("Page not found %s", context.funcName)
}

// Forbidden return 403 Forbidden
func (context *WebSocketContext) Forbidden() {
	context.Successed = false
	context.Error = fmt.Errorf("Permission forbidden")
}

// NotImplemented return 501 Not Implemented
func (context *WebSocketContext) NotImplemented() {
	context.Successed = false
	context.Error = fmt.Errorf("%s not implemented", context.funcName)
}

// ServerError return 500 server error
func (context *WebSocketContext) ServerError(err error) {
	context.Successed = false
	context.Error = fmt.Errorf("Server Error: %s", err.Error())
}

// PermanentlyMoved to url (301)
func (context *WebSocketContext) PermanentlyMoved(url string) {
}

// TemporarilyMoved to url (302)
func (context *WebSocketContext) TemporarilyMoved(url string) {
}

// GetUser get current user
func (context *WebSocketContext) GetUser() *user.TypeDB {
	return nil
}

// GetContext get global context
func (context *WebSocketContext) GetContext(key string) (value interface{}, ok bool) {
	return GetContext(key)
}

// GetClientIP returns client ip
func (context *WebSocketContext) GetClientIP() string {
	clientIP := getIPFromHeader(&context.Request.Header, X_Real_IP)
	if clientIP == "" {
		clientIP = getIPFromHeader(&context.Request.Header, X_FORWARDED_FOR)
	}
	return clientIP
}
