package queue

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/OhYee/blotter/register"
	"github.com/OhYee/rainbow/errors"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
} // use default options

type wsPool struct {
	pool  []*websocket.Conn
	mutex *sync.Mutex
}

func (p *wsPool) Add(ws *websocket.Conn) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.pool = append(p.pool, ws)
}

func (p *wsPool) Remove(ws *websocket.Conn) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	temp := make([]*websocket.Conn, 0)
	for _, w := range p.pool {
		if w != ws {
			temp = append(temp, w)
		}
	}
	p.pool = temp
}

var pool = &wsPool{
	pool:  make([]*websocket.Conn, 0),
	mutex: new(sync.Mutex),
}

type request struct {
	ID       int64       `json:"id"`
	Function string      `json:"function"`
	Args     interface{} `json:"args"`
}

type response struct {
	ID       int64       `json:"id"`
	Function string      `json:"function"`
	Success  bool        `json:"success"`
	Error    string      `json:"error"`
	Data     interface{} `json:"data"`
}

var wsFuncs = map[string]register.HandleFunc{
	// "get":  Get,
	// "pop":  Pop,
	// "push": Push,
}

func WebSocket(context register.HandleContext) (err error) {
	ws, err := upgrader.Upgrade(context.GetResponse(), context.GetRequest(), nil)
	if err != nil {
		errors.Wrapper(&err)
		return
	}
	defer ws.Close()

	if err = ws.SetReadDeadline(time.Now().Add(time.Minute * 5)); err != nil {
		return
	}

	pool.Add(ws)
	defer pool.Remove(ws)

	for {
		var mt int
		var b []byte
		mt, b, err = ws.ReadMessage()
		if err != nil {
			errors.Wrapper(&err)
			break
		}

		if mt != websocket.TextMessage {
			if err = ws.WriteMessage(mt, b); err != nil {
				break
			}
			continue
		}

		req := new(request)
		res := new(response)
		if err = json.Unmarshal(b, req); err != nil {
			errors.Wrapper(&err)
			return
		}
		res.ID = req.ID
		res.Function = req.Function

		wsContext := register.NewWebSocketContext(context.GetRequest(), context.GetResponse(), ws, req.Function, req.Args)
		if f, exist := wsFuncs[req.Function]; exist {
			err = f(wsContext)
			if err != nil {
				wsContext.ServerError(err)
			}
		} else {
			wsContext.NotImplemented()
		}

		if wsContext.Successed {
			res.Success = true
			res.Data = wsContext.Data
		} else {
			res.Success = false
			res.Error = wsContext.Error.Error()
		}

		if err = ws.WriteJSON(res); err != nil {
			errors.Wrapper(&err)
			break
		}

	}

	return
}
