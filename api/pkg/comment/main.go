package comment

import (
	"fmt"
	"time"

	"github.com/OhYee/blotter/output"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/OhYee/blotter/api/pkg/avatar"
	"github.com/OhYee/blotter/api/pkg/email"
	"github.com/OhYee/blotter/api/pkg/markdown"
	"github.com/OhYee/blotter/mongo"
)

var defaultObjectID = primitive.ObjectID{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

// Get comment of url
func Get(url string) (total int64, comments []TypeDB, err error) {
	comments = make([]TypeDB, 0)

	total, err = mongo.Find(
		"blotter",
		"comments",
		bson.M{
			"url": url,
		},
		options.Find().SetSort(bson.M{"time": 1}),
		&comments,
	)
	if err != nil {
		return
	}
	return
}

// MakeRelation make comment relation
func MakeRelation(_comments []TypeDB) (comments []*Type) {
	m := make(map[string]*Type)
	for _, cmdb := range _comments {
		cm := cmdb.ToComment()
		cm.Email = fmt.Sprintf("%c******%c", cm.Email[0], cm.Email[len(cm.Email)-1])

		if !cm.Show || cm.Ad {
			cm.Content = ""
		}

		m[cm.ID] = cm
		if parent, exist := m[cmdb.Reply.Hex()]; cmdb.Reply.Hex() != defaultObjectID.Hex() && exist {
			parent.Children = append(parent.Children, cm)
		}
	}

	comments = make([]*Type, 0)
	for _, cmdb := range _comments {
		if cmdb.Reply.Hex() == defaultObjectID.Hex() {
			comments = append(comments, m[cmdb.ID.Hex()])
		}
	}

	return
}

// Add a new comment
func Add(url string, reply string, email string, recv bool, raw string) (err error) {
	html, err := markdown.Render(raw)
	if err != nil {
		html = raw
	}

	replyObjectID, err := primitive.ObjectIDFromHex(reply)
	if err != nil {
		replyObjectID = defaultObjectID
	}

	_, err = mongo.Add("blotter", "comments", nil, TypeDB{
		ID:      primitive.NewObjectID(),
		Avatar:  avatar.Get(email),
		Email:   email,
		Reply:   replyObjectID,
		URL:     url,
		Recv:    recv,
		Raw:     raw,
		Content: html,
		Time:    time.Now().Unix(),
		Ad:      false,
		Show:    true,
	})
	if err != nil {
		return
	}
	go SendEmail(url, html, replyObjectID)

	return
}

// GetInfo get comment info
func GetInfo(url string, id primitive.ObjectID) Info {
	info := make([]Info, 0)
	cnt, err := mongo.Aggregate(
		"blotter", "comments",
		[]bson.M{
			{
				"$match": bson.M{
					"_id": id,
				},
			},
			{
				"$set": bson.M{
					"url": func(path string) string {
						if len(path) > 6 {
							path = path[6:]
						}
						return path
					}(url),
				},
			},
			{
				"$lookup": bson.M{
					"from":         "posts",
					"localField":   "url",
					"foreignField": "url",
					"as":           "posts",
				},
			},
			{
				"$set": bson.M{
					"title": "$posts.title",
					"size":  bson.M{"$size": "$posts"},
				},
			},
			{
				"$project": bson.M{
					"title": 1,
					"email": 1,
					"recv":  1,
					"size":  1,
				},
			},
			{
				"$set": bson.M{
					"title": bson.M{
						"$cond": bson.M{
							"if": bson.M{
								"$eq": []interface{}{"$size", 0},
							},
							"then": []interface{}{""},
							"else": "$title",
						},
					},
				},
			},
			{
				"$unwind": "$title",
			},
		}, nil, &info)
	if err == nil && cnt > 0 {
		return info[0]
	}
	return Info{}
}

// SendEmail for comment
func SendEmail(url string, html string, replyObjectID primitive.ObjectID) {
	emailAddr, user, username, password, address, root, blogName, err := email.GetSMTPData()
	if err != nil {
		return
	}

	info := GetInfo(url, replyObjectID)
	if info.Title == "" {
		info.Title = blogName
	}

	to := []string{emailAddr}
	if info.Recv {
		to = append(to, info.Email)
	}

	output.Debug("Send email to %+v", to)
	email.Send(
		address, username, user, password, "博客评论通知",
		fmt.Sprintf(
			"<html><body>您在<a href='%s'>《%s》</a>( %s )的评论收到一条回复<br><br>%s</body></html>",
			root+url, info.Title, root+url, html,
		),
		to...,
	)
	return
}
