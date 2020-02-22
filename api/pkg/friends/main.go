package friends

import (
	"github.com/OhYee/blotter/mongo"
	"go.mongodb.org/mongo-driver/bson"
)

func GetFriends() (friends []Friend, err error) {
	friends = make([]Friend, 0)
	_, err = mongo.Find("blotter", "friends", bson.M{}, nil, &friends)
	return
}

func SetFriends(fs []Friend) (err error) {
	if _, err = mongo.Remove("blotter", "friends", bson.M{}, nil); err != nil {
		return
	}

	slice := make([]interface{}, len(fs))
	for idx, f := range fs {
		slice[idx] = WithIndex{Index: idx, Friend: f}
	}

	_, err = mongo.Add(
		"blotter", "friends", nil,
		slice...,
	)
	return
}
