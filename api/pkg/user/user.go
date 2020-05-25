package user

import (
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"math/rand"
	"time"

	"github.com/OhYee/blotter/mongo"
	"github.com/OhYee/blotter/output"
	"github.com/OhYee/goutils/bytes"
	"github.com/OhYee/goutils/set"
	"github.com/OhYee/rainbow/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TypeBase struct {
	Username string `json:"username" bson:"username"`
	Avatar   string `json:"avatar" bson:"avatar"`
	Email    string `json:"email" bson:"email"`
	QQ       string `json:"qq" bson:"qq"`
	Black    int64  `json:"black" bson:"black"`

	Token string `json:"token" bson:"token"`

	NintendoSwitch       string `json:"ns_id" bson:"ns_id"`
	NintendoSwitchName   string `json:"ns_name" bson:"ns_name"`
	AnimalCrossingName   string `json:"ac_name" bson:"ac_name"`
	AnimalCrossingIsland string `json:"ac_island" bson:"ac_island"`

	Permission int64 `json:"permission" bson:"permission"`
}

// TypeDB of user information
type TypeDB struct {
	TypeBase `bson:",inline"`

	ID primitive.ObjectID `json:"id" bson:"_id"`

	Password  string `json:"password" bson:"password"`
	QQToken   string `json:"qq_token" bson:"qq_token"`
	QQOpenID  string `json:"qq_open_id" bson:"qq_open_id"`
	QQUnionID string `json:"qq_union_id" bson:"qq_union_id"`

	GithubID    int64  `json:"github_id" bson:"github_id"`
	GithubToken string `json:"github_token" bson:"github_token"`
}

type Type struct {
	TypeBase `bson:",inline"`

	ID string `json:"id" bson:"_id"`

	QQConnected bool `json:"qq_connected" bson:"qq_connected"`

	GithubConnected bool `json:"github_connected" bson:"github_connected"`

	Existed bool `json:"existed" bson:"existed"`
	Self    bool `json:"self" bson:"self"`
}

func NewUser(username, password string) *TypeDB {
	if u := GetUserByUsername(username); u != nil {
		return nil
	}
	uid := primitive.NewObjectID()
	u := &TypeDB{
		TypeBase: TypeBase{
			Username:       username,
			Avatar:         "/static/img/noimg.png",
			Token:          "",
			Email:          "",
			QQ:             "",
			NintendoSwitch: "",
			Permission:     0,
		},
		ID:       uid,
		Password: PasswordHash(username, password, uid.Hex()),

		QQToken:   "",
		QQOpenID:  "",
		QQUnionID: "",

		GithubID:    0,
		GithubToken: "",
	}
	if _, err := mongo.Add("blotter", "users", nil, u); err != nil {
		output.Err(err)
		return nil
	}
	return u
}

// Desensitization data desensitization
func (u *TypeDB) Desensitization(self bool) (uu *Type) {
	if u == nil {
		return &Type{
			TypeBase: TypeBase{
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
			},
			ID:              "000000000000000000000000",
			QQConnected:     false,
			GithubConnected: false,
			Existed:         false,
			Self:            false,
		}
	}
	uu = &Type{
		TypeBase:        u.TypeBase,
		ID:              u.ID.Hex(),
		QQConnected:     u.QQUnionID != "",
		GithubConnected: u.GithubID != 0,
		Existed:         true,
		Self:            self,
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
	if err != nil {
		output.Err(err)
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

func (u *TypeDB) HasPermission() bool {
	return u != nil && u.Permission != 0
}

// CheckPassword check password is right
func (u *TypeDB) CheckPassword(password string) bool {
	return PasswordHash(u.Username, password, u.ID.Hex()) == u.Password
}

// GenerateToken generate token for this user
func (u *TypeDB) GenerateToken() string {
	buf := bytes.NewBuffer()

	buf.Write(bytes.FromString(u.ID.Hex()))
	buf.Write(bytes.FromInt64(time.Now().UnixNano()))
	buf.Write(bytes.FromInt64(rand.Int63()))

	hash := sha512.New()
	hash.Write(buf.Bytes())

	u.Token = hex.EncodeToString(hash.Sum(nil))

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
	"github_token",
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

// ChangePassword change password with  password plaintext
func (u *TypeDB) ChangePassword(password string) (err error) {
	if password == "" {
		err = errors.New("Password can not be empty")
		return
	}

	_, err = mongo.Update("blotter", "users", bson.M{
		"_id": u.ID,
	}, bson.M{
		"$set": bson.M{"password": PasswordHash(u.Username, password, u.ID.Hex())},
	}, nil)

	return

}

// GetUsers get user list
func GetUsers(
	offset int64, number int64,
	sortField string, sortType int,
	searchWord string,
) (total int64, users []TypeDB, err error) {
	pipeline := []bson.M{}
	users = make([]TypeDB, 0)

	if number != 0 {
		pipeline = append(pipeline, mongo.AggregateOffset(offset, number)...)
	}
	if searchWord != "" {
		pipeline = append(pipeline, bson.M{
			"$match": bson.M{
				"username": bson.M{"$regex": searchWord, "$options": "i"},
			},
		})
	}
	if sortField != "" {
		pipeline = append(pipeline, bson.M{
			"$sort": bson.M{
				sortField: sortType,
			},
		})
	}

	total, err = mongo.Aggregate(
		"blotter", "users",
		pipeline, nil,
		&users,
	)
	return
}
