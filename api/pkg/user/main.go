package user

import (
	"crypto/sha512"
	"encoding/hex"
	"math/rand"
	"time"

	"github.com/OhYee/blotter/api/pkg/variable"
	"github.com/OhYee/blotter/mongo"
	"github.com/OhYee/blotter/register"
	"github.com/OhYee/goutils/bytes"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func PasswordHash(password string) string {
	hash := sha512.New()
	hash.Write([]byte(password))
	return hex.EncodeToString(hash.Sum(nil))
}

// Login using password
func Login(password string) bool {
	password = PasswordHash(password)

	m, err := variable.Get("password")
	if err != nil {
		return false
	}

	_password := ""
	m.SetString("password", &_password)

	if _password != password {
		return false
	}
	return true
}

// GenerateToken generate a valid token
func GenerateToken() (token string) {
	buf := bytes.NewBuffer()

	buf.Write(bytes.FromInt64(time.Now().Unix()))
	buf.Write(bytes.FromInt64(time.Now().UnixNano()))
	buf.Write(bytes.FromInt64(rand.Int63()))

	hash := sha512.New()
	hash.Write(buf.Bytes())

	token = hex.EncodeToString(hash.Sum(nil))

	mongo.Update("blotter", "variables", bson.M{"key": "token"},
		bson.M{"$set": bson.M{"value": token}}, options.Update().SetUpsert(true))
	return
}

// CheckToken check the token is valid
func CheckToken(token string) bool {
	m, err := variable.Get("token")
	if err != nil {
		return false
	}
	_token := ""
	m.SetString("token", &_token)
	if _token != token {
		return false
	}
	return true
}

// DeleteToken delete the token
func DeleteToken() {
	mongo.Remove("blotter", "variables", bson.M{"key": "token"}, nil)
}

func CheckUserPermission(context register.HandleContext) bool {
	httpContext, ok := context.(*register.HTTPContext)
	return ok && CheckToken(httpContext.GetCookie("token"))
}
