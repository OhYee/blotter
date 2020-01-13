package ms

import (
	"encoding/json"
	gb "github.com/OhYee/goutils/bytes"
	"github.com/OhYee/rainbow/errors"
	"github.com/xtaci/kcp-go"
	"io"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type any = interface{}

// HandleFunc handle function for api
type HandleFunc = func(request []byte) (response []byte, err error)

// HandleFuncWithServer handdle function for api with server argument
type HandleFuncWithServer = func(server *Server, request []byte) (response []byte, err error)

// Server object
type Server struct {
	info            *ServerInfo
	listener        net.Listener
	apiMap          map[string]HandleFunc
	subServerStatus map[string]Status // SubServerStatus status of this server
	deadTime        int64
	mutex           *sync.Mutex
	threadNumber    int
	errorCallback   func(threadID int, err error)
}

// NewServer initial the Server
func NewServer(serverInfo *ServerInfo, threadNumber int) (server *Server) {
	server = &Server{
		info:            serverInfo,
		apiMap:          make(map[string]HandleFunc),
		subServerStatus: make(map[string]Status),
		deadTime:        60,
		mutex:           new(sync.Mutex),
		threadNumber:    threadNumber,
	}
	server.Register("/heartbeat", server.handleHeartBeat)
	server.Register("/status", server.handleStatus)
	// server.Register("/api", server.handleAPI)
	return
}

// Register API function
func (server *Server) Register(address string, f HandleFunc) (err error) {
	server.mutex.Lock()
	if ff, exist := server.apiMap[address]; exist {
		err = errors.New(
			"Address %v has already registered by %v, can not register again",
			address, ff,
		)
		return
	}
	server.apiMap[address] = f
	server.mutex.Unlock()
	return
}

// RegisterWithServer register API function with server argument
func (server *Server) RegisterWithServer(address string, f HandleFuncWithServer) (err error) {
	outer := func(request []byte) (response []byte, err error) {
		response, err = f(server, request)
		return
	}

	if err = server.Register(address, outer); err != nil {
		return
	}
	return
}

// Start server listener
func (server *Server) Start() (err error) {
	server.mutex.Lock()
	listener, err := kcp.Listen(server.info.Address)
	server.mutex.Unlock()
	if err != nil {
		return
	}

	for i := 0; i < server.threadNumber; i++ {
		go func(threadID int) {
			for {
				if err := server.loop(listener); err != nil {
					server.errorCallback(threadID, err)
				}
			}
		}(i)
	}

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGUSR1, syscall.SIGUSR2)
	<-c

	return
}

func (server *Server) loop(listener net.Listener) (err error) {
	conn, err := server.listener.Accept()
	if err != nil {
		return
	}
	server.handle(conn)
	return
}

// Handle search handle function in API map
func (server *Server) handle(rw io.ReadWriter) (err error) {
	var address, request, response []byte
	var handleFunc HandleFunc
	var exist bool

	if address, err = gb.ReadWithLength32(rw); err != nil {
		return
	}

	server.mutex.Lock()
	if handleFunc, exist = server.apiMap[string(address)]; !exist {
		err = errors.New("No such API")
		return
	}
	server.mutex.Unlock()

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

	server.mutex.Lock()
	server.subServerStatus[status.Info.Address] = status
	// delete over-time status
	for k, v := range server.subServerStatus {
		if now-v.RecvTime >= server.deadTime {
			delete(server.subServerStatus, k)
		}
	}
	server.mutex.Unlock()

	if response, err = json.Marshal(map[string]string{"info": "ok"}); err != nil {
		return
	}

	return
}

func (server *Server) handleStatus(request []byte) (response []byte, err error) {
	server.mutex.Lock()
	response, err = json.Marshal(server.subServerStatus)
	server.mutex.Unlock()
	return
}

// func (server *Server) handleAPI(request []byte) (response []byte, err error) {
// 	server.mutex.Lock()
// 	response, err = json.Marshal(server.info.APIList)
// 	server.mutex.Unlock()
// 	return
// }
