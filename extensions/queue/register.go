package queue

import (
	"github.com/OhYee/blotter/register"
)

// Register api
func Register() *register.Extension {
	ext := register.NewExtension("queue")
	ext.PreRegister(
		"/",
		Get,
	)
	ext.PreRegister(
		"push",
		Push,
	)
	ext.PreRegister(
		"pop",
		Pop,
	)
	ext.PreRegister(
		"admin",
		Admin,
	)
	ext.PreRegister(
		"ws",
		WebSocket,
	)
	return ext
}
