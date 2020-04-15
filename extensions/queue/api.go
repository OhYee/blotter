package queue

import (
	"github.com/OhYee/blotter/api"
	"github.com/OhYee/blotter/register"
)

type Type struct {
	ObjectID string `json:"_id" bson:"_id"`
	ID       string `json:"id" bson:"id"`
	Name     string `json:"name" bson:"name"`
	Time     int64  `json:"time" bson:"time"`
	Finish   bool   `json:"finish" bson:"finish"`
}

type GetRequest struct {
	ID string `json:"id"`
}

type GetResponse struct {
	Queue []Type `json:"queue"`
}

func Get(context register.HandleContext) (err error) {
	args := new(GetRequest)
	res := new(GetResponse)
	context.RequestArgs(args)

	if res.Queue, err = get(args.ID); err != nil {
		return
	}

	err = context.ReturnJSON(res)
	return
}

type PushRequest struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type PushResponse api.SimpleResponse

func Push(context register.HandleContext) (err error) {
	args := new(PushRequest)
	res := new(PushResponse)
	context.RequestArgs(args)

	if err = push(args.ID, args.Name); err != nil {
		return
	}

	res.Success = true
	res.Title = "排队成功"

	context.ReturnJSON(res)
	return
}

type PopRequest struct {
	ID string `json:"id"`
}

type PopResponse api.SimpleResponse

func Pop(context register.HandleContext) (err error) {
	args := new(PopRequest)
	res := new(PopResponse)
	context.RequestArgs(args)

	if err = pop(args.ID); err != nil {
		return
	}

	res.Success = true
	res.Title = "出队成功"

	context.ReturnJSON(res)

	return
}

type AdminRequest struct {
	ObjectID string `json:"_id"`
	ID       string `json:"id"`
	Type     string `json:"type"`
}

type AdminResponse api.SimpleResponse

func Admin(context register.HandleContext) (err error) {
	if !context.GetUser().HasPermission() {
		context.Forbidden()
		return
	}

	args := new(AdminRequest)
	res := new(AdminResponse)
	context.RequestArgs(args)

	if err = admin(args.ObjectID, args.ID, args.Type); err != nil {
		return
	}

	res.Success = true
	res.Title = "操作成功"

	context.ReturnJSON(res)

	return
}
