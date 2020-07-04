package api

import (
	"github.com/OhYee/blotter/api/pkg/travels"
	"github.com/OhYee/blotter/output"
	"github.com/OhYee/blotter/register"
	"github.com/OhYee/rainbow/errors"
)

type TravelsGetResponse struct {
	Total   int64          `json:"total"`
	Travels []travels.Type `json:"travels"`
}

func TravelsGet(context register.HandleContext) (err error) {
	res := new(TravelsGetResponse)

	if res.Total, res.Travels, err = travels.Get(); err != nil {
		context.ServerError(err)
	}

	err = context.ReturnJSON(res)
	return
}
