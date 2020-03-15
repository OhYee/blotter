package variable

import (
	"github.com/OhYee/blotter/mongo"
	"go.mongodb.org/mongo-driver/bson"
)

// Get variables of keys
func Get(keys ...string) (res Variables, err error) {
	res = make(Variables)

	data := make([]map[string]interface{}, 0)
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

	for _, d := range FromMapSliceToTypeSlice(data) {
		res[d.Key] = d.Value
	}
	return
}
