package api

import (
	"github.com/OhYee/blotter/api/pkg/menu"
	"github.com/OhYee/blotter/register"
)

// Menus api, query all menus return []Menu
func Menus(context *register.HandleContext) (err error) {
	res, err := menu.Get()
	if err != nil {
		return
	}
	context.ReturnJSON(res)
	return
}
