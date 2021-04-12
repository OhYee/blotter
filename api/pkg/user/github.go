package user

import (
	github "github.com/OhYee/auth_github"
	"github.com/OhYee/blotter/api/pkg/variable"
	"github.com/OhYee/blotter/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetGithubConnect() (conn *github.Connect) {
	conn = github.New("", "", "")

	var id, secret, redirect string
	v, err := variable.Get("github_id", "github_secret", "github_redirect")
	if err != nil {
		return
	}
	if v.SetString("github_id", &id) != nil {
		return
	}
	if v.SetString("github_secret", &secret) != nil {
		return
	}
	if v.SetString("github_redirect", &redirect) != nil {
		return
	}
	conn = github.New(id, secret, redirect)
	return
}

// GithubConnect connect github and return user data
func GithubConnect(code string, state string) (token string, res github.UserInfo, err error) {
	githubConn := GetGithubConnect()

	token, err = githubConn.Auth(code, state)
	if err != nil {
		return
	}
	res, err = githubConn.Info(token)
	if err != nil {
		return
	}
	_, err = mongo.Update("blotter", "users", bson.M{
		"github_id": res.ID,
	}, bson.M{
		"$set": bson.M{
			"github_token": token,
		},
	}, nil)
	return
}

func GetUserByGithubID(githubID int64) *TypeDB {
	users := make([]TypeDB, 0)
	cnt, err := mongo.Find("blotter", "users", bson.M{
		"github_id": githubID,
	}, nil, &users)
	if err == nil && cnt != 0 {
		return &users[0]
	}
	return nil
}

func NewUserFromGithubConnect(token string, info github.UserInfo) (u *TypeDB, err error) {
	objID := primitive.NewObjectID()
	username := objID.Hex()
	uu := GetUserByUsername(info.Name)
	if uu == nil {
		username = info.Name
	}
	u = &TypeDB{
		TypeBase: TypeBase{
			Username:       username,
			Avatar:         info.Avatar,
			Token:          "",
			Email:          info.Email,
			QQ:             "",
			NintendoSwitch: "",
			Permission:     0,
		},
		ID: objID,

		Password: "",

		GithubID:    info.ID,
		GithubToken: token,
	}

	_, err = mongo.Add("blotter", "users", nil, u)
	return
}

func (u *TypeDB) ConnectGithub(token string, info github.UserInfo) (err error) {
	u.GithubToken = token
	u.GithubID = info.ID
	_, err = mongo.Update("blotter", "users", bson.M{
		"_id": u.ID,
	}, bson.M{
		"$set": bson.M{
			"github_token": token,
			"github_id":    info.ID,
		},
	}, nil)
	return
}
