package api

import (
	"github.com/OhYee/blotter/mongo"
	"github.com/OhYee/blotter/output"
	"github.com/OhYee/blotter/register"
	"go.mongodb.org/mongo-driver/bson"
)

func Friends(context *register.HandleContext) (err error) {
	output.Debug("call friends")
	res := make([]Friend, 0)
	_,err = mongo.Find("blotter", "friends", bson.M{}, nil, &res)
	if err != nil {
		return
	}
	output.Debug("%+v", res)
	context.ReturnJSON(res)
	return
}
