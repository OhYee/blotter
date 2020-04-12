package register

import (
	"net/http"

	"github.com/gorilla/schema"
)

var decoder = func() (decoder *schema.Decoder) {
	decoder = schema.NewDecoder()
	decoder.SetAliasTag("json")
	decoder.IgnoreUnknownKeys(true)
	return
}()

// HandleContext context of a api call
type HandleContext interface {
	GetRequest() *http.Request
	GetResponse() http.ResponseWriter
	// Write data to the context
	Write(b ...byte) (err error)
	// Request params
	RequestArgs(args interface{}, opts ...string)
	// ReturnJSON return json data
	ReturnJSON(data interface{}) (err error)
	// ReturnText return text data
	ReturnText(data string) (err error)
	// ReturnXML return xml data
	ReturnXML(data string) (err error)
	// Success return success
	Success()
	// PageNotFound return page not fount
	PageNotFound()
	// Forbidden return Forbidden
	Forbidden()
	// NotImplemented return Not Implemented
	NotImplemented()
	// ServerError return server error
	ServerError(err error)

	// PermanentlyMoved to url (301)
	PermanentlyMoved(url string)

	// TemporarilyMoved to url (302)
	TemporarilyMoved(url string)
}
