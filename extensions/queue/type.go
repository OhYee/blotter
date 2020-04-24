package queue

import (
	"github.com/OhYee/blotter/api/pkg/user"
	"github.com/OhYee/blotter/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Queue for animial crossing
type Queue struct {
	QueueDB `bson:",inline"`
	Leader  *user.Type `json:"leader" bson:"leader"`
	Queue   []*Member  `json:"queue" bson:"queue"`
}

type QueueDB struct {
	ID          primitive.ObjectID `json:"id" bson:"_id"`
	Leader      primitive.ObjectID `json:"leader" bson:"leader"`
	Password    string             `json:"password" bson:"paswsword"`
	Description string             `json:"description" bson:"description"`
	Max         int64              `json:"max" bson:"max"`
	CreateTime  int64              `json:"create_time" bson:"create_time"`
	FinishTime  int64              `json:"finish_time" bson:"finish_time"`
	Queue       []*MemberDB        `json:"queue" bson:"queue"`
}

// QueueDBsToQueues transfer []*QueueDB to []*Queue
func QueueDBsToQueues(queues []*QueueDB) (res []*Queue) {
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
		res = append(res, &Queue{
			QueueDB: *q,
			Leader:  userMap[q.Leader].Desensitization(false),
			Queue:   ms,
		})
	}

	return
}

func (q *QueueDB) ToQueue() *Queue {
	queues := QueueDBsToQueues([]*QueueDB{q})
	if len(queues) > 0 {
		return queues[0]
	}
	return nil
}

type Member struct {
	MemberDB `bson:",inline"`
	User     *user.Type `json:"user" bson:"user"`
}

type MemberDB struct {
	ID      primitive.ObjectID `json:"id" bson:"_id"`
	User    primitive.ObjectID `json:"user" bson:"user"`
	Queue   primitive.ObjectID `json:"queue" bson:"queue"`
	InTime  int64              `json:"in_time" bson:"in_time"`
	OutTime int64              `json:"out_time" bson:"out_time"`
	Status  int8               `json:"status" bson:"status"` // 0 in queue; 1 landed; 2 backed; 3 canceled
}

// ToMember transfer MemberDB to Member
func (m *MemberDB) ToMember(u *user.TypeDB) *Member {
	return &Member{
		MemberDB: *m,
		User:     u.Desensitization(false),
	}
}
