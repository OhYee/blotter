package user

import (
	"strings"

	qq "github.com/OhYee/auth_qq"
	"github.com/OhYee/blotter/env"
	"github.com/OhYee/blotter/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var environments, _ = env.GetEnv(env.PWDFile(".env"))

// QQConn qq connect object
var QQConn = qq.New(environments["APPID"], environments["APPKey"], environments["RedirectURI"])

// QQConnect connect qq and return user data
func QQConnect(code string) (token, openID, unionID string, res qq.UserInfo, err error) {
	token, err = QQConn.Auth(code)
	if err != nil {
		return
	}

	_, openID, unionID, err = QQConn.OpenID(token)
	if err != nil {
		return
	}
	res, err = QQConn.Info(token, openID)
	if err != nil {
		return
	}

	return
}

func GetUserByQQUnionID(unionID string) *TypeDB {
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
		TypeBase: TypeBase{
			Username:       username,
			Avatar:         strings.Replace(userInfo.FigQQ, "http://", "https://", 1),
			Token:          "",
			Email:          "",
			QQ:             "",
			NintendoSwitch: "",
			Permission:     0,
		},
		ID: objID,

		Password: "",

		QQToken:   token,
		QQOpenID:  openID,
		QQUnionID: unionID,
	}

	_, err = mongo.Add("blotter", "users", nil, u)
	return
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
