package user

import (
	"errors"
	"fmt"

	"github.com/OhYee/blotter/mongo"
	"github.com/OhYee/blotter/output"
	"github.com/OhYee/goutils/set"
	qq "github.com/OhYee/qqconnect"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TypeDB of user information
type TypeDB struct {
	ID                   primitive.ObjectID `json:"id" bson:"_id"`
	Username             string             `json:"username" bson:"username"`
	Password             string             `json:"password" bson:"password"`
	Avatar               string             `json:"avatar" bson:"avatar"`
	Token                string             `json:"token" bson:"token"`
	Email                string             `json:"email" bson:"email"`
	QQ                   string             `json:"qq" bson:"qq"`
	NintendoSwitch       string             `json:"ns_id" bson:"ns_id"`
	NintendoSwitchName   string             `json:"ns_name" bson:"ns_name"`
	AnimalCrossingName   string             `json:"ac_name" bson:"ac_name"`
	AnimalCrossingIsland string             `json:"ac_island" bson:"ac_island"`

	Permission int64 `json:"permission" bson:"permission"`

	QQToken   string `json:"qq_token" bson:"qq_token"`
	QQOpenID  string `json:"qq_open_id" bson:"qq_open_id"`
	QQUnionID string `json:"qq_union_id" bson:"qq_union_id"`
}

type Type struct {
	ID       string `json:"id" bson:"_id"`
	Username string `json:"username" bson:"username"`
	Avatar   string `json:"avatar" bson:"avatar"`
	Email    string `json:"email" bson:"email"`
	QQ       string `json:"qq" bson:"qq"`

	Token string `json:"token" bson:"token"`

	NintendoSwitch       string `json:"ns_id" bson:"ns_id"`
	NintendoSwitchName   string `json:"ns_name" bson:"ns_name"`
	AnimalCrossingName   string `json:"ac_name" bson:"ac_name"`
	AnimalCrossingIsland string `json:"ac_island" bson:"ac_island"`

	Permission int64 `json:"permission" bson:"permission"`

	QQConnected bool `json:"qq_connected" bson:"qq_connected"`

	Existed bool `json:"existed" bson:"existed"`
	Self    bool `json:"self" bson:"self"`
}

func NewUser(username, password string) *TypeDB {
	if u := GetUserByUsername(username); u != nil {
		return nil
	}
	u := &TypeDB{
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

func GetUserByToken(token string) *TypeDB {
	if token == "" {
		return nil
	}

	users := make([]TypeDB, 0)
	cnt, err := mongo.Find("blotter", "users", bson.M{
		"token": token,
	}, nil, &users)
	if err == nil && cnt != 0 {
		return &users[0]
	}
	return nil
}

func GetUserByUsername(username string) *TypeDB {
	users := make([]TypeDB, 0)
	cnt, err := mongo.Find("blotter", "users", bson.M{
		"username": bson.M{"$regex": fmt.Sprintf("^%s$", username), "$options": "i"},
	}, nil, &users)
	if err == nil && cnt != 0 {
		return &users[0]
	}
	return nil
}

func GetUserByUnionID(unionID string) *TypeDB {
	users := make([]TypeDB, 0)
	cnt, err := mongo.Find("blotter", "users", bson.M{
		"qq_union_id": unionID,
	}, nil, &users)
	if err == nil && cnt != 0 {
		return &users[0]
	}
	return nil
}

func NewUserFromQQConnect(token string, openID string, unionID string, userInfo qq.UserInfo) (u *TypeDB, err error) {
	objID := primitive.NewObjectID()
	username := objID.Hex()
	uu := GetUserByUsername(userInfo.Nickname)
	if uu == nil {
		username = userInfo.Nickname
	}
	u = &TypeDB{
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

func (u *TypeDB) HasPermission() bool {
	return u != nil && u.Permission != 0
}

func (u *TypeDB) ConnectQQ(token string, openID string, unionID string, userinfo qq.UserInfo) (err error) {
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
func (u *TypeDB) CheckPassword(password string) bool {
	return PasswordHash(u.Username, password) == u.Password
}

// GenerateToken generate token for this user
func (u *TypeDB) GenerateToken() string {
	u.Token = GenerateToken()
	u.updateToken(u.Token)
	return u.Token
}

// ClearToken delete token field
func (u *TypeDB) ClearToken() {
	u.updateToken("")
	return
}

func (u *TypeDB) updateToken(token string) {
	u.updateField("token", token)
	return
}

var validKeys = set.NewSet(
	"username", "email", "avatar", "ns_id",
	"ns_name", "ac_name", "ac_island", "qq",
)

// updateField update user field
func (u *TypeDB) updateField(key string, value interface{}) (err error) {
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
func (u *TypeDB) UpdateFields(fields map[string]string) (err error) {
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
func (u *TypeDB) ChangePassword(username, password string) (err error) {
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

// Desensitization data desensitization
func (u *TypeDB) Desensitization(self bool) (uu *Type) {
	if u == nil {
		return &Type{
			ID:                   "000000000000000000000000",
			Username:             "",
			Avatar:               "",
			Email:                "",
			QQ:                   "",
			Token:                "",
			NintendoSwitch:       "",
			NintendoSwitchName:   "",
			AnimalCrossingName:   "",
			AnimalCrossingIsland: "",
			Permission:           0,
			QQConnected:          false,
			Existed:              false,
			Self:                 false,
		}
	}
	uu = &Type{
		ID:                   u.ID.Hex(),
		Username:             u.Username,
		Avatar:               u.Avatar,
		Email:                u.Email,
		QQ:                   u.QQ,
		Token:                u.Token,
		NintendoSwitch:       u.NintendoSwitch,
		NintendoSwitchName:   u.NintendoSwitchName,
		AnimalCrossingName:   u.AnimalCrossingName,
		AnimalCrossingIsland: u.AnimalCrossingIsland,
		Permission:           u.Permission,
		QQConnected:          u.QQUnionID != "",
		Existed:              true,
		Self:                 self,
	}
	if !self {
		if len(uu.Email) > 2 {
			uu.Email = fmt.Sprintf("%c******%c", uu.Email[0], uu.Email[len(uu.Email)-1])
		}
		uu.QQ = ""
		uu.Token = ""
	}
	return
}
