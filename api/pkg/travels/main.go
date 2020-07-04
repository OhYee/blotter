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

