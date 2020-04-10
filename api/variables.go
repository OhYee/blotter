package api

import (
	"github.com/OhYee/blotter/api/pkg/user"
	"github.com/OhYee/blotter/api/pkg/variable"
	"github.com/OhYee/blotter/register"
)

// VariablesResponse response of Variables api
type VariablesResponse []variable.Type

// Variables get avatar of emial
func Variables(context register.HandleContext) (err error) {
	if !user.CheckUserPermission(context) {
		context.Forbidden()
		return
	}

	res, err := variable.GetAll()
	if err != nil {
		return
	}

	err = (context).ReturnJSON(res)
	return
}

// VariablesSetRequest request of VariablesSet api
type VariablesSetRequest struct {
	Data []variable.Type `json:"data"`
}

// VariablesSetResponse response of VariablesSet api
type VariablesSetResponse SimpleResponse

// VariablesSet get avatar of emial
func VariablesSet(context register.HandleContext) (err error) {
	if !user.CheckUserPermission(context) {
		context.Forbidden()
		return
	}

	args := new(VariablesSetRequest)
	res := new(VariablesSetResponse)
	context.RequestArgs(args, "post")

	err = variable.SetMany(args.Data...)
	if err != nil {
		return
	}

	res.Success = true
	res.Title = "修改成功"

	err = context.ReturnJSON(res)
	return
}
