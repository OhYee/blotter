package api

import (
	"github.com/OhYee/blotter/api/pkg/friends"
	"github.com/OhYee/blotter/register"
)

// Friends API query all friends, return []Friend
func Friends(context *register.HandleContext) (err error) {
	res, err := friends.GetFriends()
	if err != nil {
		return
	}
	context.ReturnJSON(res)
	return
}
