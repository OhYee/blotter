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

func Variables(context *register.HandleContext) (err error) {
	args := new(VariablesRequest)
	res := make(VariablesResponse)

	context.RequestArgs(args)
	keys := strings.Split(args.Keys, ",")

	data := make([]Variable, 0)
	var a int64
	a, err = mongo.Find(
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
	if err != nil {
		return
	}

	for _, d := range data {
		res[d.Key] = d.Value
	}
	output.Debug("%d %+v %t", a, data, args.Keys)
	err = context.ReturnJSON(res)

	return
}
