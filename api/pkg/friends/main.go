package friends

import (
	"github.com/OhYee/blotter/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// GetSimpleFriends get simple friend type
func GetSimpleFriends() (friends []Simple, err error) {
	friends = make([]Simple, 0)
	opts := options.Find().SetProjection(bson.M{
		"name": 1,
		"link": 1,
	})
	_, err = mongo.Find("blotter", "friends", bson.M{}, opts, &friends)
	return
}

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

func SetFriendPosts(url string, posts []FriendPost) (err error) {
	_, err = mongo.Update(
		"blotter",
		"friends",
		bson.M{
			"link": url,
		},
		bson.M{
			"$set": bson.M{
				"posts": posts,
				"error": len(posts) == 0,
			},
		},
		nil,
	)
	return
}
