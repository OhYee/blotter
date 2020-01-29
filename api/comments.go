package api

import (
	"fmt"
	"github.com/OhYee/blotter/mongo"
	"github.com/OhYee/blotter/output"
	"github.com/OhYee/blotter/register"
	"github.com/OhYee/goutils/time"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CommentsRequest struct {
	URL string `json:"url"`
}
type CommentsResponse struct {
	Total    int64            `json:"total"`
	Comments []*CommentTime `json:"comments"`
}

func Comments(context *register.HandleContext) (err error) {
	output.Debug("call friends")
	args := new(CommentsRequest)
	res := new(CommentsResponse)
	res.Comments = make([]*CommentTime, 0)

	context.RequestArgs(args)

	comments := make([]*CommentUnix, 0)
	res.Total, err = mongo.Find(
		"blotter",
		"comments",
		bson.M{
			"url": args.URL,
		},
		options.Find().SetSort(bson.M{"time": 1}),
		&comments,
	)
	if err != nil {
		return
	}

	m := make(map[int]*CommentTime)
	for _, data := range comments {
		output.Debug("%+v", data)
		comment := data.ToCommentTime()
		comment.Mail = fmt.Sprintf("%c******%c",comment.Mail[0], comment.Mail[len(comment.Mail)-1])
		m[data.ID] = &comment
		if parent, exist := m[data.Reply]; data.Reply != -1 && exist {
			parent.Children = append(parent.Children, &comment)
		}
	}
	for _, data := range comments {
		if data.Reply == -1 {
			res.Comments = append(res.Comments, m[data.ID])
		}
	}
	err = context.ReturnJSON(res)

	return
}

func (comment CommentUnix) ToCommentTime() CommentTime {
	return CommentTime{
		ID:       comment.ID,
		Mail:     comment.Mail,
		Avatar:   comment.Avatar,
		Time:     time.ToString(comment.Time),
		Content:  comment.Content,
		Children: make([]*CommentTime, 0),
	}
}

func (comment CommentTime) ToCommentUnix() CommentUnix {
	return CommentUnix{
		ID:      comment.ID,
		Mail:    comment.Mail,
		Avatar:  comment.Avatar,
		Time:    time.FromString(comment.Time),
		Content: comment.Content,
	}
}
