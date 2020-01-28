package register

import (
	"github.com/OhYee/blotter/output"
	"github.com/OhYee/rainbow/errors"
)

var (
	apiMap = make(map[string]HandleFunc)
)

// HandleFunc handle function type
type HandleFunc func(context *HandleContext) (err error)

// Register api
func Register(name string, f HandleFunc) {
	_, exist := apiMap[name]
	if exist {
		output.Log("API %s has existed, it will be replace by the new one", name)
	}
	apiMap[name] = f
}

// Call function
func Call(name string, context *HandleContext) (err error) {
	output.Debug("Call api %s, %s, %s", name, context.Request.URL.Path, context.Request.URL.Query())

	output.Debug("%+v", apiMap)
	api, exist := apiMap[name]
	if !exist {
		err = errors.New("Can not find api %s", name)
		return
	}
	err = api(context)
	return
}
