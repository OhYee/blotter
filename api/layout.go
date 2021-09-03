package api

import (
	"github.com/OhYee/blotter/api/pkg/friends"
	"github.com/OhYee/blotter/api/pkg/menu"
	"github.com/OhYee/blotter/api/pkg/post"
	"github.com/OhYee/blotter/api/pkg/variable"
	"github.com/OhYee/rainbow/log"

	"github.com/OhYee/blotter/mongo"
	"github.com/OhYee/blotter/register"
	"go.mongodb.org/mongo-driver/bson"
)

// LayoutResponse response of layout api
type LayoutResponse struct {
	Menus        []menu.Type      `json:"menus"`
	View         int64            `json:"view"`
	Beian        string           `json:"beian"`
	BlogName     string           `json:"blog_name"`
	Friends      []friends.Simple `json:"friends"`
	Email        string           `json:"email"`
	QQ           string           `json:"qq"`
	Github       string           `json:"github"`
	Grey         bool             `json:"grey"`
	Root         string           `json:"root"`
	Avatar       string           `json:"avatar"`
	Author       string           `json:"author"`
	From         string           `json:"from"`
	Head         string           `json:"head"`
	ADShow       string           `json:"ad_show"`
	ADInner      string           `json:"ad_inner"`
	ADText       string           `json:"ad_text"`
	Version      string           `json:"back_version"`
	EasterEgg    string           `json:"easter_egg"`
	Notification string           `json:"notification"`
}

// Layout get site base info
func Layout(context register.HandleContext) (err error) {
	res := LayoutResponse{}

	v, ok := context.GetContext("version")
	if !ok {
		v = "UNKNOWN"
	}
	vs, ok := v.(string)
	if !ok {
		vs = "UNKNOWN"
	}
	res.Version = vs

	res.Menus, err = menu.Get()
	if err != nil {
		return
	}

	m, err := variable.Get(
		"beian", "view", "blog_name", "email", "github", "qq", "grey", "root",
		"author", "avatar", "from", "head", "ad_inner", "ad_show", "ad_text",
		"easter_egg", "notification",
	)
	if err != nil {
		return
	}

	if err = m.SetInt64("view", &res.View); err != nil {
		log.Error.Println(err)
	}
	if err = m.SetString("beian", &res.Beian); err != nil {
		log.Error.Println(err)
	}
	if err = m.SetString("blog_name", &res.BlogName); err != nil {
		log.Error.Println(err)
	}
	if err = m.SetString("email", &res.Email); err != nil {
		log.Error.Println(err)
	}
	if err = m.SetString("github", &res.Github); err != nil {
		log.Error.Println(err)
	}
	if err = m.SetString("qq", &res.QQ); err != nil {
		log.Error.Println(err)
	}
	if res.Friends, err = friends.GetSimpleFriends(); err != nil {
		log.Error.Println(err)
	}
	if err = m.SetBool("grey", &res.Grey, false); err != nil {
		log.Error.Println(err)
	}
	if err = m.SetString("root", &res.Root); err != nil {
		log.Error.Println(err)
	}
	if err = m.SetString("avatar", &res.Avatar); err != nil {
		log.Error.Println(err)
	}
	if err = m.SetString("author", &res.Author); err != nil {
		log.Error.Println(err)
	}
	if err = m.SetString("from", &res.From); err != nil {
		log.Error.Println(err)
	}
	if err = m.SetString("head", &res.Head); err != nil {
		log.Error.Println(err)
	}
	if err = m.SetString("ad_show", &res.ADShow); err != nil {
		log.Error.Println(err)
	}
	if err = m.SetString("ad_inner", &res.ADInner); err != nil {
		log.Error.Println(err)
	}
	if err = m.SetString("ad_text", &res.ADText); err != nil {
		log.Error.Println(err)
	}
	if err = m.SetString("easter_egg", &res.EasterEgg); err != nil {
		log.Error.Println(err)
	}
	if err = m.SetString("notification", &res.Notification); err != nil {
		log.Error.Println(err)
	}

	err = context.ReturnJSON(res)
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
