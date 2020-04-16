package user

import (
	"github.com/OhYee/blotter/env"
	qq "github.com/OhYee/qqconnect"
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
