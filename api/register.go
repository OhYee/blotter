package api

import (
	"github.com/OhYee/blotter/register"
)

func Register() {
	register.Register(
		"friends",
		friends,
	)
	register.Register(
		"menu",
		getMenu,
	)
	register.Register(
		"post",
		getPost,
	)
	register.Register(
		"posts",
		posts,
	)
}
