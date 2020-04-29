package api

import (
	"fmt"
	"strings"

	"github.com/OhYee/blotter/api/pkg/user"
	"github.com/OhYee/blotter/register"
	"github.com/OhYee/goutils/condition"
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
		res.User = u.Desensitization(true)
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

	context.RequestArgs(args)

	if args.Username == "" {
		res = (*InfoResponse)(context.GetUser().Desensitization(true))
	} else {
		u := user.GetUserByUsername(args.Username)
		if u == nil {
			context.PageNotFound()
			return
		}
		res = (*InfoResponse)(u.Desensitization(!(context.GetUser() == nil || u.ID != context.GetUser().ID)))
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

type JumpToQQRequest struct {
	State string `json:"state"`
}

func JumpToQQ(context register.HandleContext) (err error) {
	args := new(JumpToQQRequest)
	context.RequestArgs(args)

	context.TemporarilyMoved(
		user.QQConn.LoginPage(
			condition.IfString(
				args.State == "",
				context.GetRequest().Header.Get("referer"),
				args.State,
			),
		),
	)

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

	var u *user.TypeDB

	token, openID, unionID, res, err := user.QQConnect(args.Code)
	if err != nil {
		return
	}

	switch args.State {
	case "connect":
		u = context.GetUser()
		if u == nil {
			context.ReturnText("You should login first\n你需要先登录")
			return
		}
		if uu := user.GetUserByUnionID(unionID); uu != nil {
			context.ReturnText(fmt.Sprintf("This QQ has connected to %s\n该 QQ 已绑定到 %s", uu.Username, uu.Username))
			return
		}
		if err = u.ConnectQQ(token, openID, unionID, res); err != nil {
			context.ReturnText(err.Error())
			return
		}
		context.ReturnText("Connect QQ successfully, refresh origin page\n绑定成功，请刷新原页面")
	case "avatar":
		u = user.GetUserByUnionID(unionID)
		if u == nil {
			context.ReturnText("This QQ is not connect to this site\n该 QQ 未在该网站绑定账号")
			return
		}
		if err = u.UpdateFields(map[string]string{"avatar": res.FigQQ}); err != nil {
			context.ReturnText(err.Error())
			return
		}
		context.ReturnText("Sync QQ avatar successfully, refresh origin page\nQQ 头像已更新，请刷新原页面")
	default:
		if u = user.GetUserByUnionID(unionID); u == nil {
			// New account
			if u, err = user.NewUserFromQQConnect(token, openID, unionID, res); err != nil {
				return
			}
		}
		u.GenerateToken()

		httpContext.SetCookie("token", u.Token)
		context.TemporarilyMoved(args.State)
	}

	return
}

type SetUserRequest struct {
	Avatar   string `json:"avatar"`
	Username string `json:"username"`
	NSID     string `json:"ns_id"`
	NSName   string `json:"ns_name"`
	ACName   string `json:"ac_name"`
	ACIsland string `json:"ac_island"`
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
		"avatar":    args.Avatar,
		"username":  args.Username,
		"ns_id":     args.NSID,
		"ns_name":   args.NSName,
		"ac_name":   args.ACName,
		"ac_island": args.ACIsland,
		"email":     args.Email,
		"qq":        args.QQ,
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

// CheckUsernameRequest request for CheckUsername api
type CheckUsernameRequest struct {
	Username string `json:"username"`
}

// CheckUsernameResponse response for CheckUsername api
type CheckUsernameResponse struct {
	Existed bool `json:"existed"`
}

// CheckUsername check if username is used
func CheckUsername(context register.HandleContext) (err error) {
	args := new(CheckUsernameRequest)
	res := new(CheckUsernameResponse)
	context.RequestArgs(args)

	res.Existed = user.GetUserByUsername(args.Username) != nil

	context.ReturnJSON(res)
	return
}

// RegisterUserRequest request for RegisterUser api
type RegisterUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// RegisterUserResponse response for RegisterUser api
type RegisterUserResponse SimpleResponse

// RegisterUser check if username is used
func RegisterUser(context register.HandleContext) (err error) {
	args := new(RegisterUserRequest)
	res := new(RegisterUserResponse)
	context.RequestArgs(args)

	if user.NewUser(args.Username, args.Password) != nil {
		res.Success = true
		res.Title = "注册成功"
	} else {
		res.Success = false
		res.Title = "注册失败"
		res.Content = "请检查用户名是否重复，以及网站是否崩溃"
	}

	err = context.ReturnJSON(res)
	return
}

func SyncQQAvatar(context register.HandleContext) (err error) {
	u := context.GetUser()
	if u == nil {
		context.TemporarilyMoved(user.QQConn.LoginPage("avatar"))
		return
	}
	res, err := user.QQConn.Info(u.QQToken, u.QQOpenID)
	if err != nil {
		context.TemporarilyMoved(user.QQConn.LoginPage("avatar"))
		return
	}
	if err = u.UpdateFields(map[string]string{"avatar": strings.Replace(res.FigQQ, "http://", "https://", 1)}); err != nil {
		context.ReturnText(err.Error())
		return
	}
	context.ReturnText("Sync QQ avatar successfully, refresh origin page\nQQ 头像已更新，请刷新原页面")
	return
}
