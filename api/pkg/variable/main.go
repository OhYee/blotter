package variable

import (
	"github.com/OhYee/blotter/mongo"
	"github.com/OhYee/goutils/transfer"
	"github.com/OhYee/rainbow/errors"
	"go.mongodb.org/mongo-driver/bson"
)

const databaseName = "variables"

var variablesFilter = bson.M{
	"key": bson.M{"$nin": []string{"token", "password"}},
}

// Get variables of keys
func Get(keys ...string) (res Variables, err error) {
	res = make(Variables)

	data := make([]map[string]interface{}, 0)
	_, err = mongo.Find(
		"blotter",
		databaseName,
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

// GetAll variables
func GetAll() (res BlotterVariables, err error) {
	defer errors.Wrapper(&err)

	temp := make([]map[string]interface{}, 0)
	if _, err = mongo.Find(
		"blotter",
		databaseName,
		variablesFilter,
		nil,
		&temp,
	); err != nil {
		return
	}
	res, err = NewBlotterVariables(temp)
	return
}

// SetMany variable
func SetMany(vars ...Type) (err error) {
	defer errors.Wrapper(&err)

	_, err = mongo.Remove(
		"blotter",
		databaseName,
		variablesFilter,
		nil,
	)
	if err != nil {
		return
	}

	_, err = mongo.Add(
		"blotter",
		databaseName,
		nil,
		transfer.ToInterfaceSlice(vars)...,
	)
	return
}
