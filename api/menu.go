package api

import (
	"github.com/OhYee/blotter/mongo"
	"github.com/OhYee/blotter/output"
	"github.com/OhYee/blotter/register"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func getMenu(context *register.HandleContext) (err error) {
	output.Debug("call Menu")
	res := make([]Menu, 0)
	err = mongo.Find("blotter", "pages", bson.M{}, options.Find().SetSort(bson.D{{"index", 1}}), &res)
	if err != nil {
		return
	}
	output.Debug("%+v", res)
	context.ReturnJSON(res)
	return
}
