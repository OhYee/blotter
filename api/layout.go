package api

import (
	"github.com/OhYee/blotter/register"
)

type LayoutResponse struct {
	Menus []Menu `json:"menus"`
	View  int    `json:"view"`
	Beian string `json:"beian"`
}

func Layout(context *register.HandleContext) (err error) {
	res := LayoutResponse{}

	res.Menus, err = getMenus()
	if err != nil {
		return
	}

	m := make(VariablesResponse)
	if m, err = getVariables("beian", "view"); err != nil {
		return
	}
	res.View = int(m["view"].(float64))
	res.Beian = m["beian"].(string)

	context.ReturnJSON(res)
	return
}
