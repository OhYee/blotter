package api

import (
	"github.com/OhYee/blotter/register"
)

func Register() {
	register.Register(
		"friends",
		Friends,
	)
	register.Register(
		"menu",
		Menu,
	)
	register.Register(
		"post",
		Post,
	)
	register.Register(
		"posts",
		Posts,
	)
	register.Register(
		"markdown",
		Markdown,
	)
	register.Register(
		"comments",
		Comments,
	)
	register.Register(
		"variables",
		Variables,
	)
}
