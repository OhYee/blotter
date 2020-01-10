package ms

import (
	"encoding/json"
	gb "github.com/OhYee/goutils/bytes"
	"github.com/OhYee/rainbow/errors"
	"io"
	"sync"
	"time"
)

type any = interface{}

// HandleFunc handle function for api
type HandleFunc = func(request []byte) (response []byte, err error)

// Server object
type Server struct {
	Info            *ServerInfo
	APIMap          map[string]HandleFunc
	SubServerStatus map[string]Status // SubServerStatus status of this server
	DeadTime        int64
	Mutex           *sync.Mutex
}

// NewServer initial the Server
func NewServer(serverInfo *ServerInfo) *Server {
	return &Server{
		Info:            serverInfo,
		APIMap:          make(map[string]HandleFunc),
		SubServerStatus: make(map[string]Status),
		DeadTime:        60,
		Mutex:           new(sync.Mutex),
	}
}

// Handle search handle function in API map
func (server *Server) Handle(rw io.ReadWriter) (err error) {
	var address, request, response []byte
	var handleFunc HandleFunc
	var exist bool

	if address, err = gb.ReadWithLength32(rw); err != nil {
		return
	}

	server.Mutex.Lock()
	if handleFunc, exist = server.APIMap[string(address)]; !exist {
		err = errors.New("No such API")
		return
	}
	server.Mutex.Unlock()

	if request, err = gb.ReadWithLength32(rw); err != nil {
		return
	}
	if response, err = handleFunc(request); err != nil {
		return
	}
	if _, err = rw.Write(response); err != nil {
		return
	}
	return
}

func (server *Server) handleHeartBeat(request []byte) (response []byte, err error) {
	var status Status
	now := time.Now().Unix()

	if err = json.Unmarshal(request, &status); err != nil {
		return
	}
	status.RecvTime = now

	server.Mutex.Lock()
	server.SubServerStatus[status.Info.Address] = status
	// delete over-time status
	for k, v := range server.SubServerStatus {
		if now-v.RecvTime >= server.DeadTime {
			delete(server.SubServerStatus, k)
		}
	}
	server.Mutex.Unlock()

	if response, err = json.Marshal(map[string]string{"info": "ok"}); err != nil {
		return
	}

	return
}
