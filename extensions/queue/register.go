package queue

import "github.com/OhYee/blotter/register"

// Register api
func Register() *register.Extension {
	ext := register.NewExtension("queue")
	ext.PreRegister(
		"create",
		CreateAndUpdate,
	)
	ext.PreRegister(
		"update",
		CreateAndUpdate,
	)
	ext.PreRegister(
		"finish",
		Finish,
	)
	ext.PreRegister(
		"insert",
		Insert,
	)
	ext.PreRegister(
		"land",
		Land,
	)
	ext.PreRegister(
		"out",
		Out,
	)
	ext.PreRegister(
		"get",
		Get,
	)
	ext.PreRegister(
		"get_all",
		GetAll,
	)

	// ext.PreRegister(
	// 	"pop",
	// 	Pop,
	// )
	// ext.PreRegister(
	// 	"admin",
	// 	Admin,
	// )
	// ext.PreRegister(
	// 	"ws",
	// 	WebSocket,
	// )
	return ext
}
