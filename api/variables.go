package api

import (
	"github.com/OhYee/blotter/mongo"
	"github.com/OhYee/blotter/output"
	"github.com/OhYee/blotter/register"
	"go.mongodb.org/mongo-driver/bson"
	"strings"
)

type VariablesRequest struct {
	Keys string `json:"keys"`
}
type VariablesResponse map[string]interface{}

func getVariables(keys ...string) (res map[string]interface{}, err error) {
	res = make(VariablesResponse)

	data := make([]Variable, 0)
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
	output.Debug("%+v", data)
	for _, d := range data {
		res[d.Key] = d.Value
	}
	return
}

func Variables(context *register.HandleContext) (err error) {
	args := VariablesRequest{}
	context.RequestArgs(&args)
	res, err := getVariables(strings.Split(args.Keys, ",")...)
	if err != nil {
		return
	}
	err = context.ReturnJSON(res)
	return
}
