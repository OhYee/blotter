package user

import (
	"github.com/OhYee/blotter/env"
	"github.com/OhYee/blotter/output"
	qq "github.com/OhYee/qqconnect"
)

var environments, _ = env.GetEnv(env.PWDFile(".env"))
var QQConn = qq.New(environments["APPID"], environments["APPKey"], environments["RedirectURI"])

type QQConnectToken struct {
	Token        string `json:"access_token"`
	Expire       int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

// QQConnect connect qq
func QQConnect(code string) (err error) {
	token, err := QQConn.Auth(code)
	if err != nil {
		return
	}

	_, openID, unionID, err := QQConn.OpenID(token)
	if err != nil {
		return
	}
	res, err := QQConn.Info(token, openID)
	if err != nil {
		return
	}

	output.Debug("%+v\n%+v\n%+v\n%+v", token, openID, unionID, res)

	return
}
