package queue

import (
	"time"

	"github.com/OhYee/blotter/api"
	"github.com/OhYee/blotter/api/pkg/user"
	"github.com/OhYee/blotter/mongo"
	"github.com/OhYee/blotter/register"
	"github.com/OhYee/rainbow/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func userValid(u *user.TypeDB) bool {
	return u.NintendoSwitch != "" &&
		u.NintendoSwitchName != "" &&
		u.AnimalCrossingName != "" &&
		u.AnimalCrossingIsland != ""
}

type CreateRequest struct {
	Max         int8   `json:"max"`
	Password    string `json:"password"`
	Description string `json:"description"`
}
type CreateResponse struct {
	api.SimpleResponse
	ID string `json:"id"`
}

// Create a queue
func Create(context register.HandleContext) (err error) {
	u := context.GetUser()
	if u == nil {
		context.Forbidden()
		return
	}

	args := new(CreateRequest)
	res := new(CreateResponse)
	context.RequestArgs(args)

	if !userValid(u) {
		res.Success = false
		res.Title = "创建失败"
		res.Content = "请在个人设置界面填写 Nintendo Switch、动森信息"
		err = context.ReturnJSON(res)
		return
	}

	if args.Max <= 0 || args.Max > 7 {
		res.Success = false
		res.Title = "创建失败"
		res.Content = "最大上岛人数应该为 1 ~ 7"
		err = context.ReturnJSON(res)
		return
	}

	if res.ID, err = create(u.ID, args.Password, args.Description); err != nil {
		res.Success = false
		res.Title = "队列创建失败"
		res.Content = err.Error()
	} else {
		res.Success = true
		res.Title = "队列创建成功"
	}

	err = context.ReturnJSON(res)
	return
}

func create(userID primitive.ObjectID, password string, description string) (id string, err error) {
	cnt, err := mongo.Find("blotter", "queue", bson.M{
		"leader":      userID,
		"finish_time": 0,
	}, nil, nil)
	if err != nil {
		return
	}
	if cnt != 0 {
		err = errors.New("您存在开启中的候机厅，请先关闭之前的候机厅再创建")
		return
	}

	ids, err := mongo.Add("blotter", "queue", nil, bson.M{
		"leader":      userID,
		"password":    password,
		"description": description,
		"create_time": time.Now().Unix(),
		"finish_time": 0,
		"queue":       []Member{},
	})
	if err != nil {
		return
	}

	id = ids[0].(primitive.ObjectID).Hex()
	return
}

type FinishRequest struct {
	ID string `json:"id"`
}
type FinishResponse api.SimpleResponse

func Finish(context register.HandleContext) (err error) {
	args := new(FinishRequest)
	res := new(FinishResponse)
	context.RequestArgs(args)

	u := context.GetUser()
	if u == nil {
		context.Forbidden()
		return
	}

	if err = finish(u.ID, args.ID); err != nil {
		res.Success = false
		res.Title = "排队完成失败"
		res.Content = err.Error()
	} else {
		res.Success = true
		res.Title = "排队完成成功"
	}

	err = context.ReturnJSON(res)
	return
}

func finish(userID primitive.ObjectID, ID string) (err error) {
	objID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		return
	}
	_, err = mongo.Update("blotter", "queue", bson.M{
		"_id":    objID,
		"leader": userID,
	}, bson.M{
		"$set": bson.M{
			"finish_time": time.Now().Unix(),
		},
	}, nil)
	return
}

type PushRequest struct {
	ID string `json:"id"`
}
type PushResponse api.SimpleResponse

func Push(context register.HandleContext) (err error) {
	u := context.GetUser()
	if u == nil {
		context.Forbidden()
		return
	}

	args := new(PushRequest)
	res := new(PushResponse)
	context.RequestArgs(args)

	if !userValid(u) {
		res.Success = false
		res.Title = "入队失败"
		res.Content = "请在个人设置界面填写 Nintendo Switch、动森信息"
		err = context.ReturnJSON(res)
		return
	}

	if err = push(u.ID, args.ID); err != nil {
		res.Success = false
		res.Title = "入队失败"
		res.Content = err.Error()
	} else {
		res.Success = true
		res.Title = "入队成功"
	}

	err = context.ReturnJSON(res)
	return
}

func push(userID primitive.ObjectID, ID string) (err error) {
	defer errors.Wrapper(&err)

	queueID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		return
	}

	cnt, err := mongo.Find("blotter", "queue", bson.M{
		"_id":         queueID,
		"finish_time": 0,
	}, nil, nil)
	if cnt == 0 {
		err = errors.New("队伍不存在或已结束")
		return
	}

	cnt, err = mongo.Find("blotter", "queue_members", bson.M{
		"user":     userID,
		"out_time": 0,
	}, nil, nil)
	if cnt != 0 {
		err = errors.New("您已经在队列中")
		return
	}

	ids, err := mongo.Add("blotter", "queue_members", nil, bson.M{
		"user":     userID,
		"in_time":  time.Now().Unix(),
		"out_time": 0,
		"status":   0,
	})
	if err != nil {
		return
	}

	memberID := ids[0].(primitive.ObjectID)
	_, err = mongo.Update("blotter", "queue", bson.M{
		"_id": queueID,
	}, bson.M{
		"$push": bson.M{
			"queue": memberID,
		},
	}, nil)

	return
}

type GetRequest struct {
	ID string `json:"id"`
}
type GetResponse struct {
	Queue *Queue `json:"queue"`
}

func Get(context register.HandleContext) (err error) {
	args := new(GetRequest)
	res := new(GetResponse)
	context.RequestArgs(args)

	if res.Queue, err = getQueue(args.ID); err != nil {
		return
	}

	if res == nil {
		context.PageNotFound()
		return
	}

	err = context.ReturnJSON(res)
	return
}

func getQueue(id string) (res *Queue, err error) {
	queueID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return
	}

	qs := make([]QueueDB, 0)

	cnt, err := mongo.Aggregate("blotter", "queue", []bson.M{
		{
			"$match": bson.M{
				"_id": queueID,
			},
		},
		{
			"$lookup": bson.M{
				"localField":   "queue",
				"foreignField": "_id",
				"from":         "queue_members",
				"as":           "queue",
			},
		},
	}, nil, &qs)

	if cnt > 0 {
		q := &qs[0]

		res = q.ToQueue()
	}

	return
}

type GetAllRequest struct {
	Offset int64 `json:"offset"`
	Number int64 `json:"number"`
	All    bool  `json:"all"`
}
type GetAllResponse struct {
	Total  int64    `json:"total"`
	Queues []*Queue `json:"queues"`
}

func GetAll(context register.HandleContext) (err error) {
	args := new(GetAllRequest)
	res := new(GetAllResponse)
	context.RequestArgs(args)

	if res.Total, res.Queues, err = getAllQueue(args.All, args.Offset, args.Number); err != nil {
		return
	}

	err = context.ReturnJSON(res)
	return
}

func getAllQueue(all bool, offset int64, number int64) (cnt int64, res []*Queue, err error) {
	pipeline := make([]bson.M, 0)
	if !all {
		pipeline = append(pipeline, bson.M{
			"$match": bson.M{
				"finish_time": 0,
			},
		})
	}
	if offset != 0 || number != 0 {
		pipeline = append(
			pipeline,
			mongo.AggregateOffset(offset, number)...,
		)
	}
	pipeline = append(
		pipeline,
		bson.M{
			"$lookup": bson.M{
				"localField":   "queue",
				"foreignField": "_id",
				"from":         "queue_members",
				"as":           "queue",
			},
		},
	)

	qs := make([]*QueueDB, 0)

	if cnt, err = mongo.Aggregate("blotter", "queue", pipeline, nil, &qs); err != nil {
		return
	}
	res = QueueDBsToQueues(qs)

	return
}