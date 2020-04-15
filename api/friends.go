package api

import (
	"github.com/OhYee/blotter/api/pkg/friends"
	"github.com/OhYee/blotter/register"
)

// Friends API query all friends, return []Friend
func Friends(context register.HandleContext) (err error) {
	res, err := friends.GetFriends()
	if err != nil {
		return
	}
	context.ReturnJSON(res)
	return
}

// SetFriendsRequest request of api SetFriends
type SetFriendsRequest struct {
	Friends []friends.Friend `json:"friends"`
}

// SetFriendsResponse response of api SetFriends
type SetFriendsResponse SimpleResponse

// SetFriends set friends data (method: POST)
func SetFriends(context register.HandleContext) (err error) {
	if !context.GetUser().HasPermission() {
		context.Forbidden()
		return
	}

	args := new(SetFriendsRequest)
	res := new(SetFriendsResponse)
	context.RequestArgs(args, "post")

	if err = friends.SetFriends(args.Friends); err != nil {
		return
	}

	res.Success = true
	res.Title = "修改成功"

	err = context.ReturnJSON(res)
	return
}
