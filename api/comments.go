package api

import (
	"fmt"
	"github.com/OhYee/blotter/mongo"
	"github.com/OhYee/blotter/output"
	"github.com/OhYee/blotter/register"
	gt "github.com/OhYee/goutils/time"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

var defaultObjectID = primitive.ObjectID{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

type CommentsRequest struct {
	URL string `json:"url"`
}
type CommentsResponse struct {
	Total    int64          `json:"total"`
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

	m := make(map[string]*CommentTime)

	for _, data := range comments {
		comment := data.ToCommentTime()
		comment.Email = fmt.Sprintf("%c******%c", comment.Email[0], comment.Email[len(comment.Email)-1])
		if !comment.Show || comment.Ad {
			comment.Content = ""
		}

		m[data.ID] = &comment
		if parent, exist := m[data.Reply]; data.Reply != defaultObjectID.Hex() && exist {
			parent.Children = append(parent.Children, &comment)
		}
	}

	for _, data := range comments {
		if data.Reply == defaultObjectID.Hex() {
			res.Comments = append(res.Comments, m[data.ID])
		}
	}

	err = context.ReturnJSON(res)

	return
}

func (comment CommentUnix) ToCommentTime() CommentTime {
	return CommentTime{
		ID:       comment.ID,
		Email:    comment.Email,
		Avatar:   comment.Avatar,
		Time:     gt.ToString(comment.Time),
		Content:  comment.Content,
		Children: make([]*CommentTime, 0),
		Ad:       comment.Ad,
		Show:     comment.Show,
		Recv:     comment.Recv,
	}
}

func (comment CommentTime) ToCommentUnix() CommentUnix {
	return CommentUnix{
		ID:      comment.ID,
		Email:   comment.Email,
		Avatar:  comment.Avatar,
		Time:    gt.FromString(comment.Time),
		Content: comment.Content,
		Ad:      comment.Ad,
		Show:    comment.Show,
		Recv:    comment.Recv,
	}
}

type CommentAddRequest struct {
	URL   string `json:"url"`
	Reply string `json:"reply"`
	Email string `json:"email"`
	Recv  bool   `json:"recv"`
	Raw   string `json:"raw"`
}

type CommentAddResponse struct{}

func CommentAdd(context *register.HandleContext) (err error) {
	args := new(CommentAddRequest)
	context.RequestArgs(args)

	html, err := RenderMarkdown(args.Raw)
	if err != nil {
		html = args.Raw
	}

	replyObjectID, err := primitive.ObjectIDFromHex(args.Reply)
	if err != nil {
		replyObjectID = defaultObjectID
	}
	c := Comment{
		ID:      primitive.NewObjectID(),
		Avatar:  getAvatar(args.Email),
		Email:   args.Email,
		Reply:   replyObjectID,
		URL:     args.URL,
		Recv:    args.Recv,
		Raw:     args.Raw,
		Content: html,
		Time:    time.Now().Unix(),
		Ad:      false,
		Show:    true,
	}

	_, err = mongo.Add("blotter", "comments", nil, c)
	if err != nil {
		return
	}

	go commentEmail(args.URL, html, replyObjectID)

	context.ReturnJSON(APIResponse{Success: true, Message: "评论发布成功"})
	return
}
