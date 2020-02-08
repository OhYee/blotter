package menu

import (
	"github.com/OhYee/blotter/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Get get all menus
func Get() (res []Type, err error) {
	res = make([]Type, 0)
	_, err = mongo.Find(
		"blotter",
		"pages",
		bson.M{},
		options.Find().SetSort(bson.M{"index": 1}),
		&res,
	)
	if err != nil {
		return
	}
	return
}
