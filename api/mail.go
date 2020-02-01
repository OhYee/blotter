package api

import (
	"fmt"
	"github.com/OhYee/blotter/mongo"
	"github.com/OhYee/blotter/output"
	"github.com/OhYee/blotter/register"
	"github.com/OhYee/rainbow/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/smtp"
	"strings"
)

func sendMail(host, username, user, password, subject, body string, to ...string) error {
	hp := strings.Split(host, ":")
	auth := smtp.PlainAuth("", user, password, hp[0])

	msg := []byte(fmt.Sprintf(
		"To: %s\r\nFrom: %s<%s>\r\nSubject: %s\r\nContent-Type: text/html;charset=UTF-8\r\n\r\n%s",
		strings.Join(to, ","),
		username, user,

		subject, body,
	))
	err := smtp.SendMail(host, auth, user, to, msg)
	return err
}

type MailRequest struct {
	Subject string `json:"subject"`
	Body    string `json:"body"`
	To      string `json:"to"`
}

func Mail(context *register.HandleContext) (err error) {
	args := new(MailRequest)
	res := new(APIResponse)
	context.RequestArgs(args)

	_, user, username, password, address, _, _, err := getSMTPData()
	if err != nil {
		return
	}

	err = sendMail(address, username, user, password, args.Subject, args.Body, strings.Split(args.To, ",")...)
	res.Success = err != nil
	res.Message = err.Error()
	context.ReturnJSON(res)
	return err
}

func getSMTPData() (email, user, username, password, address, root, blogName string, err error) {
	m, err := getVariables(
		"email", "smtp_user", "smtp_password", "smtp_address",
		"smtp_username", "root", "blog_name",
	)
	if err != nil {
		return
	}

	var set = func(s *string, name string) (err error) {
		var v interface{}
		var ok bool
		if v, ok = m[name]; !ok {
			err = errors.New("Can not get value of %s", name)
			return
		}
		if *s, ok = v.(string); !ok {
			err = errors.New("Value of %s is %s %T, not string", name, v, v)
			return
		}
		return err
	}

	if err = set(&email, "email"); err != nil {
		return
	}
	if err = set(&user, "smtp_user"); err != nil {
		return
	}
	if err = set(&username, "smtp_username"); err != nil {
		return
	}
	if err = set(&password, "smtp_password"); err != nil {
		return
	}
	if err = set(&address, "smtp_address"); err != nil {
		return
	}
	if err = set(&root, "root"); err != nil {
		return
	}
	if err = set(&blogName, "blog_name"); err != nil {
		return
	}
	return
}

func commentEmail(url string, html string, replyObjectID primitive.ObjectID) {
	email, user, username, password, address, root, blogName, err := getSMTPData()
	if err != nil {
		return
	}

	title := blogName
	to := []string{email}

	res := []struct {
		Email string `bson:"email"`
		Title string `bson:"title"`
		Recv  bool   `bson:"recv"`
	}{}
	cnt, err := mongo.Aggregate(
		"blotter", "comments",
		[]bson.M{
			{
				"$match": bson.M{
					"_id": replyObjectID,
				},
			},
			{
				"$set": bson.M{
					"url": func(path string) string {
						if len(path) > 3 {
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
		}, nil, &res)
	if err == nil && cnt > 0 && res[0].Recv {
		if res[0].Title != "" {
			title = res[0].Title
		}
		to = append(to, res[0].Email)
	}

	output.Debug("Send email to %+v", to)
	sendMail(
		address, username, user, password, "博客评论通知",
		fmt.Sprintf(
			"<html><body>您在<a href='%s'>《%s》</a>( %s )的评论收到一条回复<br><br>%s</body></html>",
			root+url, title, root+url, html,
		),
		to...,
	)
	return
}
