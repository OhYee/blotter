package api

import (
	"github.com/OhYee/blotter/api/pkg/avatar"
	"github.com/OhYee/blotter/register"
)

const defaultAvatar = "https://github.githubassets.com/images/modules/logos_page/GitHub-Mark.png"

// AvatarRequest request of avatar api
type AvatarRequest struct {
	Email string `json:"email"`
}

// AvatarResponse response of avatar api
type AvatarResponse struct {
	Avatar string `json:"avatar"`
}

// Avatar get avatar of emial
func Avatar(context register.HandleContext) (err error) {
	args := new(AvatarRequest)
	res := new(AvatarResponse)
	context.RequestArgs(args)

	res.Avatar = avatar.Get(args.Email)

	context.ReturnJSON(res)
	return
}
