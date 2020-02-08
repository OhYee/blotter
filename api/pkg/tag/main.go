package tag

import (
	"github.com/OhYee/blotter/mongo"
	"go.mongodb.org/mongo-driver/bson"
)

// GetTags get all tags with count
func GetTags() (total int64, res []WithCount, err error) {
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

// SearchTags using keyword (will limited with number and offset)
func SearchTags(keyword string) (total int64, tags []Type, err error) {
	tags = make([]Type, 0)
	total, err = mongo.Find("blotter", "tags", bson.M{
		"$or": []bson.M{
			{"name": bson.M{"$regex": keyword}},
			{"short": bson.M{"$regex": keyword}},
		},
	}, nil, &tags)
	if err != nil {
		return
	}
	return
}
