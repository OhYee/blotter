package queue

import (
	"fmt"
	"strings"
	"time"

	"github.com/OhYee/blotter/api"
	"github.com/OhYee/blotter/api/pkg/user"
	"github.com/OhYee/blotter/api/pkg/variable"
	"github.com/OhYee/blotter/mongo"
	"github.com/OhYee/blotter/output"
	"github.com/OhYee/blotter/register"
	"github.com/OhYee/goutils/condition"
	"github.com/OhYee/rainbow/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func userValid(u *user.TypeDB) error {
	if u.QQUnionID == "" ||
		u.QQ == "" ||
		u.NintendoSwitch == "" ||
		u.NintendoSwitchName == "" ||
		u.AnimalCrossingName == "" ||
		u.AnimalCrossingIsland == "" {
		return errors.New("请在个人设置界面填写 QQ 号、Nintendo Switch、动森信息并绑定 QQ 互联")
	}
	if u.Black > time.Now().Unix() {
		return errors.New("您由于多次违规，已被拉黑")
	}
	return nil
}

type CreateAndUpdateRequest struct {
	ID          string `json:"id"` // Only need when updating
	Max         int8   `json:"max"`
	Password    string `json:"password"`
	Description string `json:"description"`
}
type CreateAndUpdateResponse struct {
	api.SimpleResponse
	ID string `json:"id"`
}

// CreateAndUpdate a queue
func CreateAndUpdate(context register.HandleContext) (err error) {
	u := context.GetUser()
	if u == nil {
		context.Forbidden()
		return
	}

	args := new(CreateAndUpdateRequest)
	res := new(CreateAndUpdateResponse)
	context.RequestArgs(args)

	if err = userValid(u); err != nil {
		res.Success = false
		res.Title = "创建失败"
		res.Content = err.Error()
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

	if args.ID != "" {
		if err = update(args.ID, u, args.Password, args.Description, args.Max); err != nil {
			res.Success = false
			res.Title = "修改信息失败"
			res.Content = err.Error()
		} else {
			res.Success = true
			res.Title = "修改信息成功"
		}

	} else {
		if res.ID, err = create(u.ID, args.Password, args.Description, args.Max); err != nil {
			res.Success = false
			res.Title = "队列创建失败"
			res.Content = err.Error()
		} else {
			res.Success = true
			res.Title = "队列创建成功"
		}

	}

	err = context.ReturnJSON(res)
	return
}

func create(userID primitive.ObjectID, password string, description string, max int8) (id string, err error) {
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
		"max":         max,
		"queue":       []Member{},
	})
	if err != nil {
		return
	}

	id = ids[0].(primitive.ObjectID).Hex()
	return
}

func update(queueID string, u *user.TypeDB, password string, description string, max int8) (err error) {
	if u == nil {
		return errors.New("user info is nil")
	}

	queueObjID, err := primitive.ObjectIDFromHex(queueID)
	if err != nil {
		return
	}

	condition := bson.M{
		"_id":         queueObjID,
		"finish_time": 0,
	}

	if !(u.Permission&1 == 1) {
		condition["leader"] = u.ID
	}

	if _, err = mongo.Update("blotter", "queue", condition, bson.M{
		"$set": bson.M{
			"password":    password,
			"description": description,
			"max":         max,
		},
	}, nil); err != nil {
		return
	}
	go boardcast(queueObjID, false)
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

	if err = finish(u, args.ID); err != nil {
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

func finish(u *user.TypeDB, ID string) (err error) {
	if u == nil {
		return errors.New("user info is nil")
	}

	objID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		return
	}

	condition := bson.M{
		"_id": objID,
	}
	if !(u.Permission&1 == 1) {
		condition["leader"] = u.ID
	}

	_, err = mongo.Update("blotter", "queue", condition, bson.M{
		"$set": bson.M{
			"finish_time": time.Now().Unix(),
		},
	}, nil)
	return
}

type InsertRequest struct {
	ID string `json:"id"`
}
type InsertResponse api.SimpleResponse

func Insert(context register.HandleContext) (err error) {
	u := context.GetUser()
	if u == nil {
		context.Forbidden()
		return
	}

	args := new(InsertRequest)
	res := new(InsertResponse)
	context.RequestArgs(args)

	if err = userValid(u); err != nil {
		res.Success = false
		res.Title = "入队失败"
		res.Content = err.Error()
		err = context.ReturnJSON(res)
		return
	}

	if err = insert(u.ID, args.ID); err != nil {
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

func insert(userID primitive.ObjectID, ID string) (err error) {
	defer errors.Wrapper(&err)

	queueID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		return
	}

	queues := make([]struct {
		Max int64 `bson:"max"`
	}, 0)
	cnt, err := mongo.Find("blotter", "queue", bson.M{
		"_id":         queueID,
		"finish_time": 0,
	}, nil, &queues)
	if err != nil || cnt == 0 {
		if err == nil {
			err = errors.New("队伍不存在或已结束")
		}
		return
	}

	members := make([]MemberDB, 0)
	if cnt, err = mongo.Find("blotter", "queue_members", bson.M{
		"queue": queueID,
		// "user":     userID,
		"out_time": 0,
	}, nil, &members); err != nil {
		return
	}
	for _, m := range members {
		if m.User == userID {
			err = errors.New("您已经在队列中")
			return
		}
	}

	ids, err := mongo.Add("blotter", "queue_members", nil, NewMemberDB(
		userID,
		queueID,
		time.Now().Unix(),
		0,
		0,
	))
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

	if cnt <= queues[0].Max {
		go boardcast(queueID, true)
	}

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

	if res.Queue, err = getQueue(context.GetUser(), args.ID); err != nil {
		return
	}

	if res.Queue == nil {
		context.PageNotFound()
		return
	}

	err = context.ReturnJSON(res)
	return
}

func getQueue(u *user.TypeDB, id string) (res *Queue, err error) {
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

		waitingMembers := q.GetWaitingMembers()
		output.Debug("%+v", waitingMembers)
		if !(u != nil &&
			(u.ID == q.Leader ||
				(len(waitingMembers) > 0 && u.ID == waitingMembers[0].User))) {
			q.Password = ""
		}

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
	res = queueDBsToQueues(qs, true)

	return
}

type LandRequest struct {
	QueueID  string `json:"queue_id"`
	MemberID string `json:"member_id"`
}
type LandResponse api.SimpleResponse

func Land(context register.HandleContext) (err error) {
	u := context.GetUser()
	if u == nil {
		context.Forbidden()
		return
	}

	args := new(LandRequest)
	res := new(LandResponse)
	context.RequestArgs(args)

	if err = userValid(u); err != nil {
		res.Success = false
		res.Title = "着陆失败"
		res.Content = err.Error()
		err = context.ReturnJSON(res)
		return
	}

	if err = landAndOut(u, args.QueueID, args.MemberID, "land"); err != nil {
		res.Success = false
		res.Title = "着陆失败"
		res.Content = err.Error()
	} else {
		res.Success = true
		res.Title = "着陆成功"
	}

	err = context.ReturnJSON(res)
	return
}

func landAndOut(u *user.TypeDB, queueID string, memberID string, op string) (err error) {
	defer errors.Wrapper(&err)
	if u == nil {
		return errors.New("user info is nil")
	}

	if op != "land" && op != "out" {
		err = errors.New("op must be \"land\" or \"out\"")
		return
	}
	fieldName := fmt.Sprintf("%s_time", op)

	queueObjID, err := primitive.ObjectIDFromHex(queueID)
	if err != nil {
		return
	}

	memberObjID, err := primitive.ObjectIDFromHex(memberID)
	if err != nil {
		return
	}

	queues := make([]struct {
		Leader primitive.ObjectID `bson:"_id"`
	}, 0)
	cnt, err := mongo.Find("blotter", "queue", bson.M{
		"_id":         queueObjID,
		"finish_time": 0,
	}, nil, &queues)
	if err != nil || cnt == 0 {
		if err == nil {
			err = errors.New("队伍不存在或已结束")
		}
		return
	}

	condition := bson.M{
		"_id":     memberObjID,
		"queue":   queueObjID,
		fieldName: 0,
	}
	if !(u.Permission&1 == 1 || u.ID == queues[0].Leader) {
		condition["user"] = u.ID
	}

	res, err := mongo.Update("blotter", "queue_members", condition, bson.M{
		"$set": bson.M{
			fieldName: time.Now().Unix(),
		},
	}, nil)

	if err == nil && res.ModifiedCount == 0 {
		err = errors.New("未找到符合的记录")
	}

	if op == "out" {
		go boardcast(queueObjID, false)
	}

	return
}

type OutRequest struct {
	QueueID  string `json:"queue_id"`
	MemberID string `json:"member_id"`
}
type OutResponse api.SimpleResponse

func Out(context register.HandleContext) (err error) {
	u := context.GetUser()
	if u == nil {
		context.Forbidden()
		return
	}

	args := new(OutRequest)
	res := new(OutResponse)
	context.RequestArgs(args)

	if err = userValid(u); err != nil {
		res.Success = false
		res.Title = "出队失败"
		res.Content = err.Error()
		err = context.ReturnJSON(res)
		return
	}

	if err = landAndOut(u, args.QueueID, args.MemberID, "out"); err != nil {
		res.Success = false
		res.Title = "出队失败"
		res.Content = err.Error()
	} else {
		res.Success = true
		res.Title = "出队成功"
	}

	err = context.ReturnJSON(res)
	return
}

func boardcast(queueObjID primitive.ObjectID, onlyOne bool) {
	notifications := make([]*boardcastType, 0)
	cnt, err := mongo.Aggregate("blotter", "queue", []bson.M{
		{"$match": bson.M{"_id": queueObjID}},
		{"$lookup": bson.M{"localField": "queue", "foreignField": "_id", "from": "queue_members", "as": "queue"}},
		{"$unwind": "$queue"},
		{"$set": bson.M{"queue.max": "$max"}},
		{"$set": bson.M{"queue.password": "$password"}},
		{"$replaceRoot": bson.M{"newRoot": "$queue"}},
		{"$match": bson.M{"out_time": 0}},
		{"$lookup": bson.M{"localField": "user", "foreignField": "_id", "from": "users", "as": "user"}},
		{"$unwind": "$user"},
		{"$project": bson.M{
			"queue":     "$queue",
			"password":  "$password",
			"max":       "$max",
			"user_id":   "$user._id",
			"username":  "$user.username",
			"email":     "$user.email",
			"qq":        "$user.qq",
			"ac_name":   "$user.ac_name",
			"ac_island": "$user.ac_island",
			"in_time":   "$in_time",
			"land_time": "$land_time",
			"out_time":  "$out_time",
		}},
	}, nil, &notifications)
	if err != nil {
		output.ErrOutput.Printf("%s\n", errors.ShowStack(err))
	}
	if cnt == 0 {
		return
	}

	root := ""
	v, err := variable.Get("root")
	if err != nil {
		output.ErrOutput.Printf("%s\n", errors.ShowStack(err))
	}
	v.SetString("root", &root)

	var landCount int8 = 0
	var status = 0
	for _, member := range notifications {
		if member.LandTime != 0 {
			landCount++
			if landCount >= member.Max {
				break
			}
		} else {
			if status == 0 {
				go member.notify(fmt.Sprintf(
					"您已被获准起飞！%s 请尽快起飞，并在着陆后点击对应按钮。队伍地址: %s",
					condition.IfString(
						member.Password == "",
						"",
						fmt.Sprintf("飞行密码是:%s.", strings.ToUpper(member.Password)),
					),
					fmt.Sprintf("%s/apps/queue/%s", strings.Trim(root, "/"), member.Queue),
				))
				if onlyOne {
					return
				}
			} else if status == 1 {
				go member.notify(fmt.Sprintf(
					"您即将起飞！请尽快前往机场做好起飞准备，等候进一步通知。队伍地址: %s",
					fmt.Sprintf("%s/apps/queue/%s", strings.Trim(root, "/"), member.Queue),
				))
			} else {
				break
			}
			status++
		}
	}

	return
}
