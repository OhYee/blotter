package register

import (
	"fmt"
	"strings"
	"sync"

	"github.com/OhYee/blotter/output"
	"github.com/OhYee/rainbow/errors"
)

var (
	apiMap   = make(map[string]HandleFunc)
	ctx      = make(map[string]interface{})
	ctxMutex = new(sync.RWMutex)
)

func GetContext(key string) (value interface{}, ok bool) {
	ctxMutex.RLock()
	defer ctxMutex.RUnlock()
	value, ok = ctx[key]
	return
}

// SetContext set global context
func SetContext(key string, value interface{}) {
	ctxMutex.Lock()
	defer ctxMutex.Unlock()
	ctx[key] = value
	return
}

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
	defer func() {
		rec := recover()
		if rec != nil {
			err = fmt.Errorf("recover error: %+v", rec)
			output.ErrOutput.Println(err)
		}
	}()

	output.Log("%s:%s [%s] %s %s\nCall api %s, %s, %s [%s]",
		context.Request.Method,
		context.Request.Host,
		context.GetClientIP(),
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
