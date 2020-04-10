package api

import (
	"github.com/OhYee/blotter/api/pkg/menu"
	"github.com/OhYee/blotter/api/pkg/user"
	"github.com/OhYee/blotter/register"
)

// Menus api, query all menus return []Menu
func Menus(context register.HandleContext) (err error) {
	res, err := menu.Get()
	if err != nil {
		return
	}
	context.ReturnJSON(res)
	return
}

// SetMenusRequest request of api SetMenus
type SetMenusRequest struct {
	Menus []menu.Type `json:"menus"`
}

// SetMenusResponse response of api SetMenus
type SetMenusResponse SimpleResponse

// SetMenus set menus data (method: POST)
func SetMenus(context register.HandleContext) (err error) {
	if !user.CheckUserPermission(context) {
		context.Forbidden()
		return
	}

	args := new(SetMenusRequest)
	res := new(SetMenusResponse)
	context.RequestArgs(args, "post")

	if err = menu.Set(args.Menus); err != nil {
		return
	}

	res.Success = true
	res.Title = "修改成功"

	err = context.ReturnJSON(res)
	return
}
