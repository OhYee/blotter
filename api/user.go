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
	User *user.Type `json:"user"`
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

	u := user.GetUserByUsername(args.Username)
	if u != nil && u.CheckPassword(args.Password) {
		u.Token = u.GenerateToken()
		res.Success = true
		res.Title = "登录成功"
		res.User = u
		httpContext.SetCookie("token", u.Token)
	} else {
		res.Success = false
		res.Title = "登录失败"
	}

	context.ReturnJSON(res)
	return
}

type InfoRequest struct {
	Username string `json:"username"`
}

// InfoResponse response of Info api
type InfoResponse user.Type

// Info get user token api
func Info(context register.HandleContext) (err error) {
	args := new(InfoRequest)
	res := new(InfoResponse)
	u := new(user.Type)

	context.RequestArgs(args)

	if args.Username == "" {
		u = context.GetUser()
		if u == nil {
			u = &user.Type{}
		}
	} else {
		u = user.GetUserByUsername(args.Username)
		if u == nil {
			context.PageNotFound()
			return
		}
		if context.GetUser() == nil || u.ID != context.GetUser().ID {
			if len(u.Email) > 2 {
				u.Email = fmt.Sprintf("%c******%c", u.Email[0], u.Email[len(u.Email)-1])
			}
		}
	}
	res = (*InfoResponse)(u)

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
	u := context.GetUser()
	if u != nil {
		u.ClearToken()
	}

	res.Success = true
	res.Title = "登出成功"
	res.Content = "Token已清除"
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
	httpContext, ok := context.(*register.HTTPContext)
	if !ok {
		err = ErrNotHTTP
		return
	}

	args := new(QQRequest)
	// res := new(QQResponse)
	context.RequestArgs(args)

	u, err := user.QQConnect(args.Code)
	if err != nil {
		return
	}

	httpContext.SetCookie("token", u.Token)
	context.TemporarilyMoved(args.State)

	// err = context.ReturnJSON(res)
	return
}

type SetUserRequest struct {
	Avatar   string `json:"avatar"`
	Username string `json:"username"`
	NS       string `json:"ns"`
	Email    string `json:"email"`
	QQ       string `json:"qq"`
	Password string `json:"password"`
}
type SetUserResponse SimpleResponse

func SetUser(context register.HandleContext) (err error) {
	args := new(SetUserRequest)
	res := new(SetUserResponse)
	context.RequestArgs(args)

	u := context.GetUser()
	if u == nil {
		context.Forbidden()
		return
	}

	if err = u.UpdateFields(map[string]string{
		"avatar":   args.Avatar,
		"username": args.Username,
		"ns":       args.NS,
		"email":    args.Email,
		"qq":       args.QQ,
	}); err != nil {
		return
	}

	if args.Password != "" {
		if err = u.ChangePassword(args.Username, args.Password); err != nil {
			return
		}
	}

	res.Success = true
	res.Title = "修改成功"

	err = context.ReturnJSON(res)
	return
}
