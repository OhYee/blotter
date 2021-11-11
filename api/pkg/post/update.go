package post

import (
	"fmt"

	"github.com/OhYee/blotter/mongo"
	"github.com/OhYee/blotter/utils/initial"
	"go.mongodb.org/mongo-driver/bson"
)

// 自动更新数据库

// 更新 published 字段为 status 字段
func updatePublishedToStatus() {
	total, err := mongo.Find("blotter", "posts", bson.M{
		"$or": []bson.M{
			{"published": true},
			{"published": false},
		},
	}, nil, nil)
	if err == nil && total > 0 {
		fmt.Println(mongo.Update("blotter", "posts", bson.M{
			"published": true,
		}, bson.M{
			"$set":   bson.M{"status": 2},
			"$unset": bson.M{"published": ""},
		}, nil))
		mongo.Update("blotter", "posts", bson.M{
			"published": false,
		}, bson.M{
			"$set":   bson.M{"status": 0},
			"$unset": bson.M{"published": ""},
		}, nil)
		fmt.Println("update", total)
	}
}

func init() {
	initial.Register(updatePublishedToStatus)
}
