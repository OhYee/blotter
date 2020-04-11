package queue

import (
	"fmt"
	"time"

	"github.com/OhYee/blotter/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func get(id string) (q []Type, err error) {
	q = make([]Type, 0)
	if _, err = mongo.Find("blotter", "queue", bson.M{
		"id": id,
	}, nil, &q); err != nil {
		return
	}
	return
}

func push(id string, name string) (err error) {
	var cnt int64

	if cnt, err = mongo.Find("blotter", "queue", bson.M{
		"id":     id,
		"name":   name,
		"finish": false,
	}, nil, nil); cnt > 0 {
		err = fmt.Errorf("您已在队列中，多次排队请等待下一轮")
		return
	}

	if _, err = mongo.Add("blotter", "queue", nil, bson.M{
		"id":     id,
		"name":   name,
		"finish": false,
		"time":   time.Now().Unix(),
	}); err != nil {
		return
	}

	return
}

func pop(id string) (err error) {
	var cnt int64

	lst := make([]struct {
		ID primitive.ObjectID `bson:"_id"`
	}, 0)

	if cnt, err = mongo.Find("blotter", "queue", bson.M{
		"id":     id,
		"finish": false,
	}, options.Find().SetSort(bson.M{"time": 1}).SetLimit(1), &lst); cnt > 0 {
		if _, err = mongo.Update("blotter", "queue", bson.M{
			"_id": lst[0].ID,
		}, bson.M{
			"$set": bson.M{"finish": true},
		}, nil); err != nil {
			return
		}
	}
	return
}

func admin(objID string, id string, t string) (err error) {
	objectID, err := primitive.ObjectIDFromHex(objID)
	if err != nil {
		err = fmt.Errorf("Object ID 格式错误")
		return
	}

	switch t {
	case "finish":
		_, err = mongo.Update("blotter", "queue",
			bson.M{"_id": objectID, "id": id},
			bson.M{"$set": bson.M{"finish": true}},
			nil)
	case "unfinish":
		_, err = mongo.Update("blotter", "queue",
			bson.M{"_id": objectID, "id": id},
			bson.M{"$set": bson.M{"finish": false}},
			nil)
	case "delete":
		_, err = mongo.Remove("blotter", "queue", bson.M{"_id": objectID, "id": id}, nil)
	}

	return
}
