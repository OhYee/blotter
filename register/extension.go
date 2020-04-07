package register

import (
	"fmt"
	"strings"

	"github.com/OhYee/blotter/output"
)

type Extension struct {
	name   string
	apiMap map[string]HandleFunc
}

func NewExtension(name string) *Extension {
	return &Extension{
		name:   name,
		apiMap: make(map[string]HandleFunc),
	}
}

func (ext *Extension) PreRegister(name string, f HandleFunc) {
	_, exist := ext.apiMap[name]
	if exist {
		output.Log("API %s has existed in %s, it will be replace by the new one", name, ext.name)
	}
	ext.apiMap[name] = f
}

func (ext *Extension) Register(prefix string) {
	for name, f := range ext.apiMap {
		Register(fmt.Sprintf("%s/%s", strings.Trim(prefix, "/"), strings.Trim(name, "/")), f)
	}
}
