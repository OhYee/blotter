package api

import (
	"github.com/OhYee/blotter/api/pkg/menu"
	"github.com/OhYee/blotter/api/pkg/variable"

	"github.com/OhYee/blotter/mongo"
	"github.com/OhYee/blotter/register"
	"go.mongodb.org/mongo-driver/bson"
)

// LayoutResponse response of layout api
type LayoutResponse struct {
	Menus    []menu.Type `json:"menus"`
	View     int         `json:"view"`
	Beian    string      `json:"beian"`
	BlogName string      `json:"blog_name"`
	// Token    string      `json:"token"`
}

// Layout get site base info
func Layout(context *register.HandleContext) (err error) {
	res := LayoutResponse{}

	res.Menus, err = menu.Get()
	if err != nil {
		return
	}

	m, err := variable.Get("beian", "view", "blog_name")
	if err != nil {
		return
	}
	res.View = int(m["view"].(float64))

	if err = m.SetString("beian", &res.Beian); err != nil {
		return
	}
	if err = m.SetString("blog_name", &res.BlogName); err != nil {
		return
	}
	// res.Token = context.GetCookie("token")

	go func() {
		mongo.Update(
			"blotter", "variables", bson.M{"key": "view"},
			bson.M{"$inc": bson.M{"value": 1}}, nil,
		)
	}()

	context.ReturnJSON(res)
	return
}
