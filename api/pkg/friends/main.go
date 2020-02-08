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
