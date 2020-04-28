package api

import (
	"net/http"
	"time"

	"github.com/OhYee/blotter/api/pkg/notification"
	"github.com/OhYee/blotter/output"
	"github.com/OhYee/blotter/register"
	"github.com/OhYee/rainbow/errors"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
} // use default options

type WebSocketRequest struct {
	Token string `json:"token"`
}

func WebSocket(context register.HandleContext) (err error) {
	args := new(WebSocketRequest)
	context.RequestArgs(args)
	if args.Token == "" {
		args.Token = context.GetUser().Desensitization(false).ID
	}

	output.Debug("%+v", *args)

	closeChannel := make(chan bool)
	writeChannel := make(chan notification.WritePackage)

	ws, err := upgrader.Upgrade(context.GetResponse(), context.GetRequest(), nil)
	if err != nil {
		errors.Wrapper(&err)
		return
	}
	defer ws.Close()
	defer func() { closeChannel <- true }()

	notification.Hub.Set(args.Token, writeChannel)
	defer notification.Hub.Remove(args.Token)

	go func() {
		for {
			select {
			case message := <-writeChannel:
				ws.WriteMessage(message.MessageType, message.MessageData)
			case <-closeChannel:
				return
			}
		}
	}()

	go func() {
		for {
			time.Sleep(10 * time.Second)
			writeChannel <- notification.WritePackage{
				MessageType: websocket.PingMessage,
				MessageData: []byte{},
			}
		}
	}()

	for {
		var mt int
		var b []byte
		mt, b, err = ws.ReadMessage()
		if err != nil {
			errors.Wrapper(&err)
			break
		}
		output.Log("%v %s", mt, b)

		switch mt {
		case websocket.PingMessage:
			writeChannel <- notification.WritePackage{
				websocket.PongMessage,
				[]byte{},
			}
		case websocket.PongMessage:

		case websocket.CloseMessage:
			ws.Close()
			return
		}
		if err != nil {
			break
		}
	}
	return
}
