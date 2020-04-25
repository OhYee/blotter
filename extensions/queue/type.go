package queue

import (
	"github.com/OhYee/blotter/api/pkg/user"
	"github.com/OhYee/blotter/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type QueueBase struct {
	Password    string `json:"password" bson:"password"`
	Description string `json:"description" bson:"description"`
	Max         int64  `json:"max" bson:"max"`
	CreateTime  int64  `json:"create_time" bson:"create_time"`
	FinishTime  int64  `json:"finish_time" bson:"finish_time"`
}

// Queue for animial crossing
type Queue struct {
	QueueBase `bson:",inline"`

	ID     string     `json:"id" bson:"_id"`
	Leader *user.Type `json:"leader" bson:"leader"`
	Queue  []*Member  `json:"queue" bson:"queue"`
}

type QueueDB struct {
	QueueBase `bson:",inline"`

	ID     primitive.ObjectID `json:"id" bson:"_id"`
	Leader primitive.ObjectID `json:"leader" bson:"leader"`
	Queue  []*MemberDB        `json:"queue" bson:"queue"`
}

// queueDBsToQueues transfer []*QueueDB to []*Queue
func queueDBsToQueues(queues []*QueueDB, removePassword bool) (res []*Queue) {
	ids := make([]primitive.ObjectID, 0)

	for _, q := range queues {
		ids = append(ids, q.Leader)
		for _, m := range q.Queue {
			ids = append(ids, m.User)
		}
	}

	u := make([]*user.TypeDB, 0)
	if _, err := mongo.Find("blotter", "users", bson.M{
		"_id": bson.M{
			"$in": ids,
		},
	}, nil, &u); err != nil {
		return nil
	}
	userMap := make(map[primitive.ObjectID]*user.TypeDB)
	for _, uu := range u {
		userMap[uu.ID] = uu
	}

	for _, q := range queues {
		ids = append(ids, q.Leader)
		for _, m := range q.Queue {
			ids = append(ids, m.User)
		}
	}

	res = make([]*Queue, 0)

	for _, q := range queues {
		ms := make([]*Member, 0)
		for _, m := range q.Queue {
			u, exist := userMap[m.User]
			if exist {
				ms = append(ms, m.ToMember(u))
			}
		}

		queue := &Queue{
			QueueBase: q.QueueBase,
			ID:        q.ID.Hex(),
			Leader:    userMap[q.Leader].Desensitization(false),
			Queue:     ms,
		}
		if removePassword {
			queue.Password = ""
		}
		res = append(res, queue)
	}

	return
}

func (q *QueueDB) ToQueue() *Queue {
	queues := queueDBsToQueues([]*QueueDB{q}, false)
	if len(queues) > 0 {
		return queues[0]
	}
	return nil
}

// GetWaitingMembers get the members which are waiting for password
func (q *QueueDB) GetWaitingMembers() (res []*MemberDB) {
	max := q.Max
	res = make([]*MemberDB, 0)

	for _, member := range q.Queue {
		if member.OutTime == 0 && member.LandTime != 0 {
			// on island
			max--
		}
	}

	if max <= 0 {
		return
	}

	for _, member := range q.Queue {
		if member.OutTime == 0 && member.LandTime == 0 {
			// waiting
			res = append(res, member)
			max--
			if max <= 0 {
				return
			}
		}
	}

	return
}

type MemberBase struct {
	InTime   int64 `json:"in_time" bson:"in_time"`
	LandTime int64 `json:"land_time" bson:"land_time"`
	OutTime  int64 `json:"out_time" bson:"out_time"`
	// Status  int8  `json:"status" bson:"status"` // 0 in queue; 1 landed; 2 backed; 3 canceled
}

type Member struct {
	MemberBase `bson:",inline"`
	ID         string     `json:"id" bson:"_id"`
	User       *user.Type `json:"user" bson:"user"`
}

type MemberDB struct {
	MemberBase `bson:",inline"`
	ID         primitive.ObjectID `json:"id" bson:"_id"`
	User       primitive.ObjectID `json:"user" bson:"user"`
	Queue      primitive.ObjectID `json:"queue" bson:"queue"`
}

// ToMember transfer MemberDB to Member
func (m *MemberDB) ToMember(u *user.TypeDB) *Member {
	return &Member{
		MemberBase: m.MemberBase,
		ID:         m.ID.Hex(),
		User:       u.Desensitization(false),
	}
}

// NewMemberDB initial a MemberDB
func NewMemberDB(
	userID, queueID primitive.ObjectID,
	inTime, landTime, outTime int64,
) *MemberDB {
	return &MemberDB{
		ID:    primitive.NewObjectID(),
		User:  userID,
		Queue: queueID,
		MemberBase: MemberBase{
			InTime:   inTime,
			LandTime: landTime,
			OutTime:  outTime,
		},
	}
}
