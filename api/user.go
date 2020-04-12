package api

import (
	"fmt"

	"github.com/OhYee/blotter/api/pkg/user"
	"github.com/OhYee/blotter/register"
)

// ErrNotHTTP the api can only be called by HTTP request
var ErrNotHTTP = fmt.Errorf("Only can be called by HTTP request")

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
func Login(context register.HandleContext) (err error) {
	httpContext, ok := context.(*register.HTTPContext)
	if !ok {
		err = ErrNotHTTP
		return
	}

	args := new(LoginRequest)
	res := new(LoginResponse)

	context.RequestArgs(args)

	if args.Username == "" && user.Login(args.Password) {
		token := user.GenerateToken()
		httpContext.SetCookie("token", token)
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
func Info(context register.HandleContext) (err error) {
	httpContext, ok := context.(*register.HTTPContext)
	if !ok {
		err = ErrNotHTTP
		return
	}

	res := new(InfoResponse)
	token := httpContext.GetCookie("token")
	if user.CheckToken(token) {
		res.Token = token
	}
	context.ReturnJSON(res)
	return
}

// Logout the user
func Logout(context register.HandleContext) (err error) {
	httpContext, ok := context.(*register.HTTPContext)
	if !ok {
		err = ErrNotHTTP
		return
	}

	res := new(SimpleResponse)
	if user.CheckUserPermission(context) {
		user.DeleteToken()
		res.Success = true
		res.Title = "登出成功"
		res.Content = "Token已清除"
	} else {
		res.Success = false
		res.Title = "登出失败"
		res.Content = "Token验证错误"
	}
	httpContext.DeleteCookie("token")
	context.ReturnJSON(res)
	return
}

func JumpToQQ(context register.HandleContext) (err error) {
	context.TemporarilyMoved(user.QQConn.LoginPage(context.GetRequest().Header.Get("referer")))
	return
}

type QQRequest struct {
	Code  string `json:"code"`
	State string `json:"state"`
}
type QQResponse struct {
	Token string `json:"token"`
}

func QQ(context register.HandleContext) (err error) {
	args := new(QQRequest)
	// res := new(QQResponse)
	context.RequestArgs(args)

	_ = user.QQConnect(args.Code)

	context.TemporarilyMoved(args.State)

	// err = context.ReturnJSON(res)
	return
}
