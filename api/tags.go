package api

import (
	"github.com/OhYee/blotter/mongo"
	"github.com/OhYee/blotter/register"
	"go.mongodb.org/mongo-driver/bson"
)

type TagResponse struct {
	Total int64          `json:"total" bson:"Total"`
	Tags  []TagWithCount `json:"tags" bson:"tags"`
}

func getTags() (total int64, res []TagWithCount, err error) {
	total, err = mongo.Aggregate("blotter", "tags", []bson.M{
		{
			"$lookup": bson.M{
				"from":         "posts",
				"localField":   "_id",
				"foreignField": "tags",
				"as":           "posts",
			},
		},
		{
			"$set": bson.M{"count": bson.M{"$size": "$posts"}},
		},
		{
			"$sort": bson.M{"count": -1},
		},
	}, nil, &res)
	return
}

func Tags(context *register.HandleContext) (err error) {
	total, tags, err := getTags()
	if err != nil {
		return
	}
	res := TagResponse{
		Total: total,
		Tags:  tags,
	}
	err = context.ReturnJSON(res)
	return
}
