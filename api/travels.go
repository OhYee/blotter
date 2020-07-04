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

type TravelsSetRequest struct {
	Travels []travels.Type `json:"travels"`
}

type TravelsSetResponse SimpleResponse

func TravelsSet(context register.HandleContext) (err error) {
	if !context.GetUser().HasPermission() {
		context.Forbidden()
		return
	}

	args := new(TravelsSetRequest)
	res := new(TravelsSetResponse)
	context.RequestArgs(args, "POST")

	output.Debug("%+v", args)

	if err = travels.Set(args.Travels); err != nil {
		res.Success = false
		res.Title = "游记设置失败"
		res.Content = errors.ShowStack(err)
	} else {
		res.Success = true
		res.Title = "游记设置成功"
	}

	err = context.ReturnJSON(res)
	return
}

type TravelsGetByURLRequest struct {
	URL string `json:"url"`
}

type TravelsGetByURLResponse struct {
	Exist  bool           `json:"exist"`
	Travel travels.Travel `json:"travel"`
}

func TravelsGetByURL(context register.HandleContext) (err error) {
	args := new(TravelsGetByURLRequest)
	res := new(TravelsGetByURLResponse)
	context.RequestArgs(args)

	if res.Exist, res.Travel, err = travels.GetByURL(args.URL); err != nil {
		context.ServerError(err)
	}

	err = context.ReturnJSON(res)
	return
}
