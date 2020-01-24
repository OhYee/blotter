package api

import (
	"github.com/OhYee/blotter/mongo"
	"github.com/OhYee/blotter/output"
	"github.com/OhYee/blotter/register"
	"go.mongodb.org/mongo-driver/bson"
)

type FriendsPostType struct {
	Title string `json:"title"`
	Link  string `json:"link"`
}

type FriendsType struct {
	Image       string            `json:"image"`
	Link        string            `json:"link"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Posts       []FriendsPostType `json:"posts"`
}

func friends(context register.HandleContext) (err error) {
	output.Debug("call friends")
	res := make([]FriendsType, 0)
	err = mongo.Query("blotter", "friends", bson.M{}, &res)
	if err != nil {
		return
	}
	output.Debug("%+v", res)
	context.ReturnJSON(res)
	return
}
