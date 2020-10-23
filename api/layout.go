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
	Root     string           `json:"root"`
	Avatar   string           `json:"avatar"`
	Author   string           `json:"author"`
	From     string           `json:"from"`
	Head     string           `json:"head"`
	ADShow   string           `json:"ad_show"`
	ADInner  string           `json:"ad_inner"`
	ADText   string           `json:"ad_text"`
}

// Layout get site base info
func Layout(context register.HandleContext) (err error) {
	res := LayoutResponse{}

	res.Menus, err = menu.Get()
	if err != nil {
		return
	}

	m, err := variable.Get(
		"beian", "view", "blog_name", "email", "github", "qq", "grey", "root",
		"author", "avatar", "from", "head", "ad_inner", "ad_show", "ad_text",
	)
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
	if err = m.SetString("root", &res.Root); err != nil {
		return
	}
	if err = m.SetString("avatar", &res.Avatar); err != nil {
		return
	}
	if err = m.SetString("author", &res.Author); err != nil {
		return
	}
	if err = m.SetString("from", &res.From); err != nil {
		return
	}
	if err = m.SetString("head", &res.Head); err != nil {
		return
	}
	if err = m.SetString("ad_show", &res.ADShow); err != nil {
		return
	}
	if err = m.SetString("ad_inner", &res.ADInner); err != nil {
		return
	}
	if err = m.SetString("ad_text", &res.ADText); err != nil {
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

	go func() {
		mongo.Update(
			"blotter", "variables", bson.M{"key": "view"},
			bson.M{"$inc": bson.M{"value": 1}}, nil,
		)
	}()

	if args.URL != "" {
		go post.IncView(args.URL)
	}

	return
}
