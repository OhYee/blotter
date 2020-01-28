package api

import (
	"github.com/OhYee/blotter/mongo"
	"github.com/OhYee/blotter/output"
	"github.com/OhYee/blotter/register"
	"go.mongodb.org/mongo-driver/bson"
	// "go.mongodb.org/mongo-driver/mongo/options"
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

type PostsRequest struct {
	Number int64  `json:"number"`
	Offset int64  `json:"offset"`
	Type   string `json:"type"`
	Arg    string `json:"arg"`
}

type PostsResponse struct {
	Total int        `json:"total"`
	Posts []PostCard `json:"posts"`
}

func posts(context *register.HandleContext) (err error) {
	args := PostsRequest{}
	context.RequestArgs(&args)

	output.Debug("%+v", args)

	res := PostsResponse{}
	res.Posts = make([]PostCard, 10)
	switch args.Type {
	case "index":
		fallthrough
	default:
		res.Total, err = mongo.Aggregate("blotter", "posts", []bson.M{
			{"$sort": bson.M{"publish_time": -1}},
			{
				"$lookup": bson.M{
					"from":         "tags",
					"localField":   "tags",
					"foreignField": "_id",
					"as":           "tags",
				},
			},
			{"$limit": args.Offset + args.Number},
			{"$skip": args.Offset},
		}, nil, &res.Posts)
	}
	if err != nil {
		return
	}
	err = context.ReturnJSON(res)
	return
}
