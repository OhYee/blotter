package api

import (
	"github.com/OhYee/blotter/mongo"
	"github.com/OhYee/blotter/register"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TagsRequest struct {
	Keyword string `json:"keyword"`
	Number  int64  `json:"number"`
	Offset  int64  `json:"offset"`
}

type TagsResponse struct {
	Total int64 `json:"total" bson:"Total"`
	Tags  []Tag `json:"tags" bson:"tags"`
}

type TagsResponseWithCount struct {
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
	args := new(TagsRequest)
	var res interface{}

	context.RequestArgs(args)

	if len(args.Keyword) == 0 {
		resWithCount := new(TagsResponseWithCount)
		resWithCount.Total, resWithCount.Tags, err = getTags()
		res = resWithCount
	} else {
		resWithoutCount := new(TagsResponse)
		resWithoutCount.Total, resWithoutCount.Tags, err = tagSearch(args.Keyword, args.Number, args.Offset)
		res = resWithoutCount
	}
	if err != nil {
		return
	}
	err = context.ReturnJSON(res)
	if err != nil {
		return
	}
	return
}

func tagSearch(keyword string, number int64, offset int64) (total int64, tags []Tag, err error) {
	tags = make([]Tag, 0)
	total, err = mongo.Find("blotter", "tags", bson.M{
		"$or": []bson.M{
			{"name": bson.M{"$regex": keyword}},
			{"short": bson.M{"$regex": keyword}},
		},
	}, options.Find().SetLimit(number+offset).SetSkip(offset), &tags)
	if err != nil {
		return
	}
	return
}
