package user

import (
	"github.com/OhYee/blotter/mongo"
	"go.mongodb.org/mongo-driver/bson"
)

// Type of user information
type Type struct {
	ID             string `json:"id" bson:"_id"`
	Username       string `json:"username" bson:"username"`
	Password       string `json:"password" bson:"password"`
	Avatar         string `json:"avatar" bson:"avatar"`
	Token          string `json:"token" bson:"token"`
	Email          string `json:"email" bson:"email"`
	QQ             string `json:"qq" bson:"qq"`
	NintendoSwitch string `json:"ns" bson:"ns"`
	Permission     int64  `json:"permission" bson:"permission"`

	QQToken   string `json:"qq_token" bson:"qq_token"`
	QQOpenID  string `json:"qq_open_id" bson:"qq_open_id"`
	QQUnionID string `json:"qq_union_id" bson:"qq_union_id"`
}

func GetUserByToken(token string) *Type {
	users := make([]Type, 0)
	mongo.Find("blotter", "user", bson.M{
		"token": token,
	}, nil, &users)
	if len(users) != 0 {
		return &users[0]
	}
	return nil
}

func (u *Type) HasPermission() bool {
	return u.Permission != 0
}
