package register

import (
	"strings"

	"github.com/OhYee/blotter/output"
	"github.com/OhYee/rainbow/errors"
)

var (
	apiMap = make(map[string]HandleFunc)
)

// HandleFunc handle function type
type HandleFunc func(context HandleContext) (err error)

// Register api
func Register(name string, f HandleFunc) {
	name = strings.Trim(name, "/")
	_, exist := apiMap[name]
	if exist {
		output.Log("API %s has existed, it will be replace by the new one", name)
	}
	apiMap[name] = f
}

func DebugApiMap() {
	for name, _ := range apiMap {
		output.Debug("%+v", name)
	}
}

// Call function
func Call(name string, context *HTTPContext) (err error) {
	output.Log("%s:%s [%s] %s\nCall api %s, %s, %s [%s]",
		context.Request.Method,
		context.Request.Host,
		context.Request.Header.Get("nginx"),
		context.Request.UserAgent(),
		name,
		context.Request.URL.Path,
		context.Forms(),
		context.GetCookie("token"),
	)

	api, exist := apiMap[name]
	if !exist {
		err = errors.New("Can not find api %s", name)
		return
	}
	err = api(context)
	return
}
