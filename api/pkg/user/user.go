package user

import (
	"errors"

	"github.com/OhYee/blotter/mongo"
	"github.com/OhYee/blotter/output"
	"github.com/OhYee/goutils/set"
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

func NewUser(username, password string) *Type {
	if u := GetUserByUsername(username); u != nil {
		return nil
	}
	u := &Type{
		ID:             primitive.NewObjectID(),
		Username:       username,
		Password:       PasswordHash(username, password),
		Avatar:         "/static/img/noimg.png",
		Token:          "",
		Email:          "",
		QQ:             "",
		NintendoSwitch: "",
		Permission:     0,
		QQToken:        "",
		QQOpenID:       "",
		QQUnionID:      "",
	}
	if _, err := mongo.Add("blotter", "users", nil, u); err != nil {
		output.Err(err)
		return nil
	}
	return u
}

func GetUserByToken(token string) *Type {
	if token == "" {
		return nil
	}

	users := make([]Type, 0)
	cnt, err := mongo.Find("blotter", "users", bson.M{
		"token": token,
	}, nil, &users)
	if err == nil && cnt != 0 {
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
	cnt, err := mongo.Find("blotter", "users", bson.M{
		"qq_union_id": unionID,
	}, nil, &users)
	if err == nil && cnt != 0 {
		return &users[0]
	}
	return nil
}

func NewUserFromQQConnect(token string, openID string, unionID string, userInfo qq.UserInfo) (u *Type, err error) {
	objID := primitive.NewObjectID()
	username := objID.Hex()
	uu := GetUserByUsername(userInfo.Nickname)
	if uu == nil {
		username = userInfo.Nickname
	}
	u = &Type{
		ID:             objID,
		Username:       username,
		Password:       "",
		Avatar:         userInfo.FigQQ,
		Token:          "",
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
	return u != nil && u.Permission != 0
}

func (u *Type) ConnectQQ(token string, openID string, unionID string, userinfo qq.UserInfo) (err error) {
	u.QQToken = token
	u.QQOpenID = openID
	u.QQUnionID = unionID
	_, err = mongo.Update("blotter", "users", bson.M{
		"_id": u.ID,
	}, bson.M{
		"$set": bson.M{
			"qq_token":    token,
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
func (u *Type) GenerateToken() string {
	u.Token = GenerateToken()
	u.updateToken(u.Token)
	return u.Token
}

// ClearToken delete token field
func (u *Type) ClearToken() {
	u.updateToken("")
	return
}

func (u *Type) updateToken(token string) {
	u.updateField("token", token)
	return
}

var validKeys = set.NewSet("username", "email", "avatar", "ns", "qq")

// updateField update user field
func (u *Type) updateField(key string, value interface{}) (err error) {
	_, err = mongo.Update("blotter", "users", bson.M{
		"_id": u.ID,
	}, bson.M{
		"$set": bson.M{
			key: value,
		},
	}, nil)
	return
}

// UpdateFields update user fields
func (u *Type) UpdateFields(fields map[string]string) (err error) {
	data := bson.M{}
	for key, value := range fields {
		if validKeys.Exist(key) {
			data[key] = value
		}
	}

	_, err = mongo.Update("blotter", "users", bson.M{
		"_id": u.ID,
	}, bson.M{
		"$set": data,
	}, nil)

	return
}

// ChangePassword change password with username and password plaintext
func (u *Type) ChangePassword(username, password string) (err error) {
	if password == "" {
		err = errors.New("Password can not be empty")
		return
	}

	_, err = mongo.Update("blotter", "users", bson.M{
		"_id": u.ID,
	}, bson.M{
		"$set": bson.M{"password": PasswordHash(u.Username, password)},
	}, nil)

	return

}
