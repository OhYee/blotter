package api

import (
	"crypto/sha512"
	"encoding/hex"
	"math/rand"
	"time"

	"github.com/OhYee/goutils/bytes"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/OhYee/blotter/mongo"
	"github.com/OhYee/blotter/output"
	"github.com/OhYee/blotter/register"
)

func login(password string) bool {
	output.Debug("%+v", password)

	hash := sha512.New()
	hash.Write([]byte(password))
	password = hex.EncodeToString(hash.Sum(nil))

	m, err := getVariables("password")
	if err != nil {
		return false
	}
	_password, ok := m["password"]

	output.Debug("%+v %+v", _password, password)
	if !ok || _password != password {
		return false
	}
	return true
}

func generateToken() (token string) {
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

func checkToken(token string) bool {
	m, err := getVariables("token")
	if err != nil {
		return false
	}
	_token, ok := m["token"]
	if !ok || _token != token {
		return false
	}
	return true
}

func deleteToken() {
	mongo.Remove("blotter", "variables", bson.M{"key": "token"}, nil)
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	APIResponse
	InfoResponse
}

func Login(context *register.HandleContext) (err error) {
	args := new(LoginRequest)
	res := new(LoginResponse)

	context.RequestArgs(args)

	if args.Username == "" && login(args.Password) {
		token := generateToken()
		context.SetCookie("token", token)
		res.Success = true
		res.Title = "登录成功"
		res.Token = token
	} else {
		res.Success = false
		res.Title = "登录失败"
	}

	context.ReturnJSON(res)
	return
}

type InfoResponse struct {
	Token string `json:"token"`
}

func Info(context *register.HandleContext) (err error) {
	res := new(InfoResponse)
	token := context.GetCookie("token")
	if checkToken(token) {
		res.Token = token
	}
	context.ReturnJSON(res)
	return
}

func Logout(context *register.HandleContext) (err error) {
	res := new(APIResponse)
	if checkToken(context.GetCookie("token")) {
		deleteToken()
		res.Success = true
		res.Title = "登出成功"
		res.Content = "Token已清除"
	} else {
		res.Success = false
		res.Title = "登出失败"
		res.Content = "Token验证错误"
	}
	context.ReturnJSON(res)
	return
}
