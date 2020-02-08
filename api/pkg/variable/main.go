package variable

import (
	"github.com/OhYee/blotter/mongo"
	"github.com/OhYee/blotter/output"
	"go.mongodb.org/mongo-driver/bson"
)

// Get variables of keys
func Get(keys ...string) (res Variables, err error) {
	res = make(Variables)

	data := make([]Type, 0)
	_, err = mongo.Find(
		"blotter",
		"variables",
		bson.M{
			"key": bson.M{
				"$in": keys,
			},
		},
		nil,
		&data,
	)
	for _, d := range data {
		res[d.Key] = d.Value
	}
	output.Debug("%+v %+v", data, res)
	return
}
