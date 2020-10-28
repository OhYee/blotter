package api

import (
	"encoding/json"
	"time"

	"github.com/OhYee/blotter/api/pkg/markdown"
	"github.com/OhYee/blotter/output"
	"github.com/OhYee/blotter/register"
	"github.com/OhYee/rainbow/errors"
	"github.com/gorilla/websocket"
)

// MarkdownRequest request of markdown api
type MarkdownRequest struct {
	Source string `json:"source"`
}

// MarkdownResponse response of markdown api
type MarkdownResponse struct {
	HTML string `json:"html"`
}

// Markdown render markdown to html
func Markdown(context register.HandleContext) (err error) {
	args := new(MarkdownRequest)
	res := new(MarkdownResponse)
	context.RequestArgs(args, "post")

	if res.HTML, err = markdown.Render(args.Source, true); err != nil {
		return
	}

	context.ReturnJSON(res)
	return
}

// Markdown render markdown to html
func MarkdownWS(context register.HandleContext) (err error) {
	ws, err := upgrader.Upgrade(context.GetResponse(), context.GetRequest(), nil)
	if err != nil {
		errors.Wrapper(&err)
		return
	}
	closeChannel := make(chan bool)
	defer ws.Close()

	ws.SetCloseHandler(func(code int, text string) error {
		closeChannel <- true
		return nil
	})

	go func() {
		for {
			time.Sleep(10 * time.Second)
			ws.WriteMessage(websocket.PingMessage, []byte{})
		}
	}()

	go func() {
		for {
			var args MarkdownRequest
			var res MarkdownResponse

			t, b, err := ws.ReadMessage()
			if err != nil {
				output.Err(err)
				closeChannel <- true
				break
			}
			switch t {
			case websocket.TextMessage:
				err = json.Unmarshal(b, &args)
				if err != nil {
					output.Err(err)
				} else {
					if res.HTML, err = markdown.Render(args.Source, true); err != nil {
						output.Err(err)
					}
					ws.WriteJSON(res)
				}
			case websocket.PingMessage:
				ws.WriteMessage(websocket.PongMessage, []byte{})
			}

		}
	}()

	<-closeChannel

	return
}
