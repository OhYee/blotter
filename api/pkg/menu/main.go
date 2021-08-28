package menu

import (
	"github.com/OhYee/blotter/mongo"
	"github.com/OhYee/blotter/output"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const DatabaseName = "pages"

func init() {
	if exists, _ := mongo.CollectionExists("blotter", DatabaseName); exists == false {
		initMenus()
	}
}

func initMenus() {
	output.LogOutput.Printf("Initial database %s", DatabaseName)
	mongo.Add(
		"blotter",
		DatabaseName,
		nil,

		WithIndex{
			Index: 0,
			Type: Type{
				Icon: "home",
				Name: "首页",
				Link: "/",
			},
		},
		WithIndex{
			Index: 1,
			Type: Type{
				Icon: "archive",
				Name: "归档",
				Link: "/archives",
			},
		},
		WithIndex{
			Index: 2,
			Type: Type{
				Icon: "tag",
				Name: "标签",
				Link: "/tags",
			},
		},
		WithIndex{
			Index: 3,
			Type: Type{
				Icon: "comments",
				Name: "评论区",
				Link: "/comments",
			},
		},
		WithIndex{
			Index: 4,
			Type: Type{
				Icon: "idcard",
				Name: "关于",
				Link: "/about",
			},
		},
		WithIndex{
			Index: 5,
			Type: Type{
				Icon: "link",
				Name: "优秀博客订阅",
				Link: "/friends",
			},
		},
		WithIndex{
			Index: 6,
			Type: Type{
				Icon: "github",
				Name: "Github: OhYee",
				Link: "https://github.com/OhYee",
			},
		},
	)
}

// Get get all menus
func Get() (res []Type, err error) {
	res = make([]Type, 0)
	_, err = mongo.Find(
		"blotter",
		DatabaseName,
		bson.M{},
		options.Find().SetSort(bson.M{"index": 1}),
		&res,
	)
	if err != nil {
		return
	}
	return
}

func Set(menus []Type) (err error) {
	if _, err = mongo.Remove("blotter", DatabaseName, bson.M{}, nil); err != nil {
		return
	}

	slice := make([]interface{}, len(menus))
	for idx, menu := range menus {
		slice[idx] = WithIndex{Index: idx, Type: menu}
	}

	_, err = mongo.Add(
		"blotter", DatabaseName, nil,
		slice...,
	)
	return
}
