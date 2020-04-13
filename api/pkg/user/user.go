package user

import (
	"github.com/OhYee/blotter/mongo"
	qq "github.com/OhYee/qqconnect"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Type of user information
type Type struct {
	ID             primitive.ObjectID `json:"id" bson:"_id"`
	Username       string             `json:"username" bson:"username"`
	Password       string             `json:"password" bson:"password"`
	Avatar         string             `json:"avatar" bson:"avatar"`
	Token          string             `json:"token" bson:"token"`
	Email          string             `json:"email" bson:"email"`
	QQ             string             `json:"qq" bson:"qq"`
	NintendoSwitch string             `json:"ns" bson:"ns"`
	Permission     int64              `json:"permission" bson:"permission"`

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

func GetUserByUsername(username string) *Type {
	users := make([]Type, 0)
	cnt, err := mongo.Find("blotter", "users", bson.M{
		"username": username,
	}, nil, &users)
	if err == nil && cnt != 0 {
		return &users[0]
	}
	return nil
}

func GetUserByUnionID(unionID string) *Type {
	users := make([]Type, 0)
	mongo.Find("blotter", "user", bson.M{
		"qq_union_id": unionID,
	}, nil, &users)
	if len(users) != 0 {
		return &users[0]
	}
	return nil
}

func NewUserFromQQConnect(token string, openID string, unionID string, userInfo qq.UserInfo) (u *Type, err error) {
	objID := primitive.NewObjectID()
	u = &Type{
		ID:             objID,
		Username:       objID.Hex(),
		Password:       "",
		Avatar:         userInfo.FigQQ,
		Token:          GenerateToken(),
		Email:          "",
		QQ:             "",
		NintendoSwitch: "",
		Permission:     0,

		QQToken:   token,
		QQOpenID:  openID,
		QQUnionID: unionID,
	}
	_, err = mongo.Add("blotter", "users", nil, u)
	return
}

func (u *Type) HasPermission() bool {
	return u.Permission != 0
}

func (u *Type) ConnectQQ(openID string, unionID string, userinfo qq.UserInfo) (err error) {
	u.QQOpenID = openID
	u.QQUnionID = unionID
	_, err = mongo.Update("blotter", "users", bson.M{
		"_id": u.ID,
	}, bson.M{
		"$set": bson.M{
			"qq_open_id":  openID,
			"qq_union_id": unionID,
		},
	}, nil)
	return
}

// CheckPassword check password is right
func (u *Type) CheckPassword(password string) bool {
	return PasswordHash(u.Username, password) == u.Password
}

// GenerateToken generate token for this user
func (u *Type) GenerateToken() (token string) {
	token = GenerateToken()
	mongo.Update("blotter", "users", bson.M{
		"_id": u.ID,
	}, bson.M{
		"$set": bson.M{
			"token": token,
		},
	}, nil)
	return
}
