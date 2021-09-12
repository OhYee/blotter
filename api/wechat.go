package api

// TODO: Remove

import (
	"github.com/OhYee/blotter/register"

	wechat "github.com/OhYee/gowechat"
)

var wc = wechat.Wechat{
	Token:  "ohyee_token",
	AESKey: "QTTdEebUiSH4FnPSH3OhY1ePjVmYv9UoAWnTZKvWg5Q",
}

// WechatCheckPermissionRequest request for WechatCheckPermission api
type WechatCheckPermissionRequest struct {
	Signature string `json:"signature"`
	Timestamp string `json:"timestamp"`
	Nonce     string `json:"nonce"`
	EchoStr   string `json:"echostr"`
}

// WechatCheckPermission check if username is used
func WechatCheckPermission(context register.HandleContext) (err error) {
	args := new(WechatCheckPermissionRequest)
	context.RequestArgs(args)

	if wc.CheckSignature(args.Signature, args.Timestamp, args.Nonce) {
		context.ReturnText(args.EchoStr)
	}

	return
}
