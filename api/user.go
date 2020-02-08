package api

import (
	"github.com/OhYee/blotter/api/pkg/user"
	"github.com/OhYee/blotter/register"
)

// LoginRequest request of login api
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginResponse response of login api
type LoginResponse struct {
	SimpleResponse
	InfoResponse
}

// Login try to login
func Login(context *register.HandleContext) (err error) {
	args := new(LoginRequest)
	res := new(LoginResponse)

	context.RequestArgs(args)

	if args.Username == "" && user.Login(args.Password) {
		token := user.GenerateToken()
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

// InfoResponse response of Info api
type InfoResponse struct {
	Token string `json:"token"`
}

// Info get user token api
func Info(context *register.HandleContext) (err error) {
	res := new(InfoResponse)
	token := context.GetCookie("token")
	if user.CheckToken(token) {
		res.Token = token
	}
	context.ReturnJSON(res)
	return
}

// Logout the user
func Logout(context *register.HandleContext) (err error) {
	res := new(SimpleResponse)
	if user.CheckToken(context.GetCookie("token")) {
		user.DeleteToken()
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
