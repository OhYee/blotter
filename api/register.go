package api

import (
	"github.com/OhYee/blotter/register"
)

func Register() {
	register.Register(
		"friends",
		friends,
	)
}
