package api

import (
	"github.com/OhYee/blotter/mongo"
	"github.com/OhYee/blotter/register"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func getMenus() (res []Menu, err error) {
	res = make([]Menu, 0)
	_, err = mongo.Find(
		"blotter",
		"pages",
		bson.M{},
		options.Find().SetSort(bson.M{"index": 1}),
		&res,
	)
	if err != nil {
		return
	}
	return
}

func Menus(context *register.HandleContext) (err error) {
	res, err := getMenus()
	if err != nil {
		return
	}
	context.ReturnJSON(res)
	return
}
