package travels

import (
	"github.com/OhYee/blotter/mongo"
	"github.com/OhYee/blotter/output"
	"github.com/OhYee/goutils/transfer"
	"github.com/OhYee/rainbow/errors"
	"go.mongodb.org/mongo-driver/bson"
)

func Get() (total int64, res []Type, err error) {
	defer errors.Wrapper(&err)
	res = make([]Type, 0)
	total, err = mongo.Find("blotter", "travels", bson.M{}, nil, &res)
	return
}

func Set(travels []Type) (err error) {
	output.Debug("%+v", travels)
	defer errors.Wrapper(&err)
	if _, err = mongo.Remove("blotter", "travels", bson.M{}, nil); err != nil {
		return
	}
	_, err = mongo.Add("blotter", "travels", nil, transfer.ToInterfaceSlice(travels)...)
	return
}

func GetByURL(url string) (exist bool, res Travel, err error) {
	var total int64
	results := make([]Travel, 0)

	if total, err = mongo.Aggregate("blotter", "travels", []bson.M{
		{"$unwind": "$travels"},
		{"$project": bson.M{
			"name": 1,
			"lng":  1,
			"lat":  1,
			"zoom": 1,
			"time": "$travels.time",
			"link": "$travels.link",
		}},
		{"$match": bson.M{"link": url}},
	}, nil, &results); err != nil {
		return
	}
	if total > 0 {
		exist = true
		res = results[0]
	}
	return
}
