package queue

import (
	"time"

	"github.com/OhYee/blotter/api"
	"github.com/OhYee/blotter/api/pkg/user"
	"github.com/OhYee/blotter/mongo"
	"github.com/OhYee/blotter/register"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func Get(context *register.HandleContext) (err error) {
	args := new(GetRequest)
	res := new(GetResponse)
	context.RequestParams(args)

	res.Queue = make([]Type, 0)
	if _, err = mongo.Find("blotter", "queue", bson.M{
		"id": args.ID,
	}, nil, &res.Queue); err != nil {
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

func Push(context *register.HandleContext) (err error) {
	args := new(PushRequest)
	res := new(PushResponse)
	context.RequestParams(args)

	if cnt, err := mongo.Find("blotter", "queue", bson.M{
		"id":     args.ID,
		"name":   args.Name,
		"finish": false,
	}, nil, nil); cnt > 0 {
		res.Success = false
		res.Title = "您已在队列中，多次排队请等待下一轮"
		context.ReturnJSON(res)
		return err
	}

	data := Type{
		ID:     args.ID,
		Name:   args.Name,
		Finish: false,
		Time:   time.Now().Unix(),
	}
	if _, err = mongo.Add("blotter", "queue", nil, data); err != nil {
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

func Pop(context *register.HandleContext) (err error) {
	args := new(PopRequest)
	res := new(PopResponse)
	context.RequestParams(args)

	var cnt int64
	lst := make([]struct {
		ID primitive.ObjectID `bson:"_id"`
	}, 0)

	if cnt, err = mongo.Find("blotter", "queue", bson.M{
		"id":     args.ID,
		"finish": false,
	}, options.Find().SetSort(bson.M{"time": 1}).SetLimit(1), &lst); cnt > 0 {
		if _, err = mongo.Update("blotter", "queue", bson.M{
			"_id": lst[0].ID,
		}, bson.M{
			"$set": bson.M{"finish": true},
		}, nil); err != nil {
			return
		}
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

func Admin(context *register.HandleContext) (err error) {
	if !user.CheckToken(context.GetCookie("token")) {
		context.Forbidden()
		return
	}

	args := new(AdminRequest)
	res := new(AdminResponse)
	context.RequestParams(args)

	objectID, err := primitive.ObjectIDFromHex(args.ObjectID)
	if err != nil {
		res.Success = false
		res.Title = "Object ID 格式错误"
		context.ReturnJSON(res)
		return
	}

	switch args.Type {
	case "finish":
		_, err = mongo.Update("blotter", "queue",
			bson.M{"_id": objectID, "id": args.ID},
			bson.M{"$set": bson.M{"finish": true}},
			nil)
	case "unfinish":
		_, err = mongo.Update("blotter", "queue",
			bson.M{"_id": objectID, "id": args.ID},
			bson.M{"$set": bson.M{"finish": false}},
			nil)
	case "delete":
		_, err = mongo.Remove("blotter", "queue", bson.M{"_id": objectID, "id": args.ID}, nil)
	}
	if err != nil {
		return
	}

	res.Success = true
	res.Title = "操作成功"

	context.ReturnJSON(res)

	return
}
