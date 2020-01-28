package api

import (
	"github.com/OhYee/blotter/mongo"
	"github.com/OhYee/blotter/output"
	"github.com/OhYee/blotter/register"
	"go.mongodb.org/mongo-driver/bson"
)

type getPostRequest struct {
	URL string `json:"url"`
}

func getPost(context *register.HandleContext) (err error) {
	output.Debug("call friends")
	args := getPostRequest{}
	context.RequestArgs(&args)

	res := make([]map[string]interface{}, 0)
	err = mongo.Find("blotter", "posts", bson.M{"url": args.URL}, nil, &res)
	if err != nil {
		return
	}
	if len(res) > 0 {
		context.ReturnJSON(res[0])
	} else {
		context.Response.WriteHeader(404)
	}
	return
}
