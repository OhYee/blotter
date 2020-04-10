package api

import (
	"github.com/OhYee/blotter/api/pkg/friends"
	"github.com/OhYee/blotter/api/pkg/menu"
	"github.com/OhYee/blotter/api/pkg/post"
	"github.com/OhYee/blotter/api/pkg/variable"

	"github.com/OhYee/blotter/mongo"
	"github.com/OhYee/blotter/register"
	"go.mongodb.org/mongo-driver/bson"
)

// LayoutResponse response of layout api
type LayoutResponse struct {
	Menus    []menu.Type      `json:"menus"`
	View     int              `json:"view"`
	Beian    string           `json:"beian"`
	BlogName string           `json:"blog_name"`
	Friends  []friends.Simple `json:"friends"`
	Email    string           `json:"email"`
	QQ       string           `json:"qq"`
	Github   string           `json:"github"`
	Grey     bool             `json:"grey"`
}

// Layout get site base info
func Layout(context register.HandleContext) (err error) {
	res := LayoutResponse{}

	res.Menus, err = menu.Get()
	if err != nil {
		return
	}

	m, err := variable.Get("beian", "view", "blog_name", "email", "github", "qq", "grey")
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
	if err = m.SetString("email", &res.Email); err != nil {
		return
	}
	if err = m.SetString("github", &res.Github); err != nil {
		return
	}
	if err = m.SetString("qq", &res.QQ); err != nil {
		return
	}
	if res.Friends, err = friends.GetSimpleFriends(); err != nil {
		return
	}
	if err = m.SetBool("grey", &res.Grey, false); err != nil {
		return
	}

	context.ReturnJSON(res)
	return
}

// ViewRequest request for inc api
type ViewRequest struct {
	URL string `json:"url"`
}

// View view number of the url
func View(context register.HandleContext) (err error) {
	args := PostRequest{}
	context.RequestArgs(&args)

	if args.URL == "" {
		go func() {
			mongo.Update(
				"blotter", "variables", bson.M{"key": "view"},
				bson.M{"$inc": bson.M{"value": 1}}, nil,
			)
		}()
	} else {
		go post.IncView(args.URL)
	}

	return
}
