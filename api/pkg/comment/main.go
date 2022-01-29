package comment

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/OhYee/blotter/output"
	pool "github.com/OhYee/blotter/utils/goroutine_pool"
	"github.com/OhYee/blotter/utils/lru"
	"github.com/OhYee/rainbow/errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/OhYee/blotter/api/pkg/avatar"
	"github.com/OhYee/blotter/api/pkg/email"
	"github.com/OhYee/blotter/api/pkg/markdown"
	"github.com/OhYee/blotter/api/pkg/post"
	"github.com/OhYee/blotter/api/pkg/serverchan"
	"github.com/OhYee/blotter/mongo"
)

// Get comment of url
func Get(url string) (total int64, comments []TypeDB, err error) {
	defer errors.Wrapper(&err)

	comments = make([]TypeDB, 0)

	total, err = mongo.Find(
		DatabaseName,
		CollectionName,
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

// GetAdmin get comments for admin page
func GetAdmin(offset int64, number int64) (total int64, comments []Admin, err error) {
	defer errors.Wrapper(&err)

	commentsDB := make([]AdminDB, 0)

	pipeline := []bson.M{
		{
			"$set": bson.M{
				"link": bson.M{
					"$arrayElemAt": []interface{}{
						bson.M{"$split": []string{"$url", "/post/"}}, 1,
					},
				},
			},
		},
		{
			"$lookup": bson.M{
				"from":         "posts",
				"localField":   "link",
				"foreignField": "url",
				"as":           "post",
			},
		},
		{
			"$lookup": bson.M{
				"from":         "comments",
				"localField":   "reply",
				"foreignField": "_id",
				"as":           "reply_comment",
			},
		},
		{
			"$set": bson.M{"reply_comment": bson.M{"$arrayElemAt": []interface{}{"$reply_comment", 0}}},
		},
		{
			"$set": bson.M{"post": bson.M{"$arrayElemAt": []interface{}{"$post", 0}}},
		},
		{
			"$set": bson.M{
				"title": "$post.title",
			},
		},
		{
			"$project": bson.M{
				"post": 0,
			},
		},
		{
			"$sort": bson.M{"time": -1},
		},
	}
	if number > 0 {
		pipeline = append(pipeline, mongo.AggregateOffset(offset, number)...)
	}
	total, err = mongo.Aggregate(
		DatabaseName,
		CollectionName,
		pipeline,
		nil,
		&commentsDB,
	)
	if err != nil {
		return
	}

	comments = make([]Admin, len(commentsDB))
	for idx, commentDB := range commentsDB {
		comments[idx] = *commentDB.ToAdmin()
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
	if antiShake(url, email, raw) {
		// shake!
		return ErrShake
	}

	html, err := markdown.Render(raw, false)
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
	go SendEmail(url, raw, html, replyObjectID)

	return
}

// GetInfo get comment info
func GetInfo(url string, id primitive.ObjectID) (info Info, title string) {
	// query title
	if strings.HasPrefix(url, "/post/") {
		url = url[6:]
	}
	p, err := post.GetPublicFieldPost(url)
	if err == nil {
		title = p.Title
	}

	// query reply info
	infoArr := make([]Info, 0)
	_, err = mongo.Aggregate("blotter", "comments", []bson.M{
		{
			"$match": bson.M{
				"_id": id,
			},
		},
	}, nil, &infoArr)
	if err == nil && len(infoArr) > 0 {
		info = infoArr[0]
	}
	return
}

// SendEmail for comment
func SendEmail(url string, raw string, html string, replyObjectID primitive.ObjectID) {
	emailAddr, user, username, password, address, ssl, root, blogName, err := email.GetSMTPData()
	if err != nil {
		return
	}

	info, title := GetInfo(url, replyObjectID)
	if title == "" {
		title = blogName
	}

	to := []string{emailAddr}
	if info.Recv {
		to = append(to, info.Email)
	}

	go serverchan.Notify("新评论提醒", fmt.Sprintf("%s 在 [%s](%s) 发布了一条评论\n\n\n\n%s", info.Email, title, root+url, raw))

	output.Debug("Send email to %+v", to)
	err = email.Send(
		address, username, user, password, ssl, "博客评论通知",
		fmt.Sprintf(
			"<html><body>您在<a href='%s'>《%s》</a>( %s )的评论收到一条回复<br><br>%s</body></html>",
			root+url, title, root+url, html,
		),
		to...,
	)
	if err != nil {
		output.Err(err)
	}
	return
}

// Set comment state by id
func Set(id string, ad bool, show bool, recv bool) (err error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return
	}
	_, err = mongo.Update("blotter", "comments", bson.M{"_id": objectID}, bson.M{
		"$set": bson.M{
			"recv": recv,
			"ad":   ad,
			"show": show,
		},
	}, nil)
	return
}

// Delete comment state by id
func Delete(id string) (err error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return
	}
	_, err = mongo.Remove("blotter", "comments", bson.M{"_id": objectID}, nil)
	return
}

var shakeMap = lru.NewMap().WithExpired()

func antiShake(url, email, raw string) bool {
	h := sha256.New()
	h.Write([]byte(fmt.Sprintf("%s|%s|%s", url, email, raw)))
	key := hex.EncodeToString(h.Sum([]byte{}))

	if _, exists := shakeMap.Get(key); exists {
		return true
	}

	shakeMap.PutWithExpired(key, struct{}{}, 5*time.Minute)
	return false
}

func UpdateAvatar() (int, int) {
	_, comments, err := GetAdmin(0, -1)
	if err != nil {
		return 0, 0
	}
	wg := sync.WaitGroup{}
	var success int64
	for _, c := range comments {
		wg.Add(1)
		func(c Admin) {
			pool.Do(func() {
				defer wg.Done()
				if err := c.UpdateAvatar(); err == nil {
					atomic.AddInt64(&success, 1)
				} else {
					output.ErrOutput.Println(err)
				}
			})
		}(c)
	}
	wg.Wait()
	return int(success), len(comments)
}
