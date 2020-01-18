package ms

import (
	"encoding/json"
	// "github.com/OhYee/blotter/micro_server/message"
	msg "github.com/OhYee/blotter/micro_server/message"
	"github.com/OhYee/goutils/bytes"
	"github.com/OhYee/rainbow/errors"
	"github.com/xtaci/kcp-go"
	"io"
	"net"
	"os"
	"os/signal"
	"reflect"
	"sync"
	"syscall"
	"time"
)

type any = interface{}

// HandleFunc handle function for api
//
// Simple:
//     - `Sum(server *Server, req SumRequest, threadID int) (rep SumResponse)`
//     - `Sum(server *Server, req SumRequest) (rep SumResponse)`
//     - `Sum(req SumRequest) (rep SumResponse)`
type HandleFunc interface{}

// Server object
type Server struct {
	gateway         string
	info            *ServerInfo
	listener        net.Listener
	apiMap          map[string]HandleFunc
	subServerStatus map[string]*Status // SubServerStatus status of this server
	deadTime        int64
	mutex           *sync.Mutex
	threadNumber    int
	logCallback     func(threadID int, format string, args ...interface{})
	errorCallback   func(threadID int, err error)
	close           bool
}

// NewServer initial the Server
func NewServer(gateway string, serverInfo *ServerInfo, threadNumber int) (server *Server) {
	server = &Server{
		gateway:         gateway,
		info:            serverInfo,
		apiMap:          make(map[string]HandleFunc),
		subServerStatus: make(map[string]*Status),
		deadTime:        60,
		mutex:           new(sync.Mutex),
		threadNumber:    threadNumber,
		close:           true,
	}

	return
}

// IsClosed return the server is closed
func (server *Server) IsClosed() bool {
	server.mutex.Lock()
	defer server.mutex.Unlock()
	return server.close
}

// Register API function
func (server *Server) Register(address string, f HandleFunc) (err error) {
	value := reflect.ValueOf(f)
	numIn := value.Type().NumIn()
	numOut := value.Type().NumOut()
	if !((numIn == 1 || numIn == 2 || numIn == 3) && (numOut == 2)) || value.Type().Out(1).Elem().String() != "error" {
		err = errors.New(
			"%v can not be register, want function like:\n"+
				"    - Sum(req SumRequest) (rep Response, err error)\n"+
				"    - Sum(server *Server, req SumRequest) (rep Response, err error)\n",
			"    - Sum(server *Server, req SumRequest, threadID int) (rep SumResponse)\n",
			value.Interface(),
		)
		return
	}

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

// StartGateway start gateway server listener
func (server *Server) StartGateway() (err error) {
	server.Register("/heartbeat", HeartBeatHandle)
	err = server.startServer()
	return
}

// Start server listener
func (server *Server) Start() (err error) {
	if err = server.startServer(); err != nil {
		return
	}

	go func() {
		for !server.IsClosed() {
			time.After(time.Second * time.Duration(server.deadTime))
			if conn, err := kcp.Dial(server.gateway); err != nil {
				server.errorCallback(-1, err)
			} else {
				server.mutex.Lock()
				conn.Write([]byte{})
				server.mutex.Unlock()
			}
		}
	}()

	return
}

func (server *Server) startServer() (err error) {
	server.mutex.Lock()
	listener, err := kcp.Listen(server.info.Address)
	server.mutex.Unlock()
	if err != nil {
		return
	}

	for i := 0; i < server.threadNumber; i++ {
		go func(threadID int) {
			for !server.IsClosed() {
				if err := server.loop(threadID, listener); err != nil {
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

func (server *Server) loop(threadID int, listener net.Listener) (err error) {
	conn, err := server.listener.Accept()
	if err != nil {
		return
	}
	server.handle(threadID, conn)
	return
}

// Handle search handle function in API map
func (server *Server) handle(threadID int, rw io.ReadWriter) (err error) {
	var address string
	var handleFunc HandleFunc
	var exist bool
	var req, rep []byte

	// got the function address
	if address, err = bytes.ReadStringWithLength32(rw); err != nil {
		return
	}

	// search the function
	server.mutex.Lock()
	if handleFunc, exist = server.apiMap[string(address)]; !exist {
		err = errors.New("No such API")
		return
	}
	server.mutex.Unlock()

	// using reflect to call the function
	function := reflect.ValueOf(handleFunc)
	request := reflect.New(function.Type().In(0)).Interface()
	response := reflect.New(function.Type().Out(0)).Interface()

	// read []byte type request
	if req, err = bytes.ReadBytesWithLength32(rw); err != nil {
		return
	}
	// transfer request data to `request`
	if err = json.Unmarshal(req, &request); err != nil {
		return
	}

	var in = make([]reflect.Value, 0)
	switch n := function.Type().NumIn(); n {
	case 1:
		in = []reflect.Value{
			reflect.ValueOf(request),
		}
	case 2:
		in = []reflect.Value{
			reflect.ValueOf(server),
			reflect.ValueOf(request),
		}
	case 3:
		in = []reflect.Value{
			reflect.ValueOf(server),
			reflect.ValueOf(request),
			reflect.ValueOf(threadID),
		}
	default:
		err = errors.New(
			"Function with %d input arguments, want (request interface{}) or (server *Server, req interface{})",
			n,
		)
		return
	}

	// call the function
	out := function.Call(in)

	// got response and error data
	response = out[0].Interface()
	err = out[1].Interface().(error)

	// transfer response to []byte
	if rep, err = json.Marshal(response); err != nil {
		return
	}

	// write data to the connection
	if _, err = rw.Write(rep); err != nil {
		return
	}

	return
}

// Call api of server, req should be a pointer
func (server *Server) Call(serverName string, apiName string, req interface{}, rep interface{}) (err error) {
	server.mutex.Lock()
	status, exist := server.subServerStatus[serverName]
	server.mutex.Unlock()
	if !exist {
		err = errors.New("Do not have %s server", serverName)
		return
	}

	var request *msg.Request
	var response *msg.Response

	requestBytes, err := json.Marshal(req)
	if err != nil {
		return
	}
	request = msg.NewRequest(apiName, requestBytes)

	if response, err = server.Send(status.Info.Address, request); err != nil {
		return
	}

	err = json.Unmarshal(response.Arguments, rep)
	return
}

// Send a request to address
func (server *Server) Send(address string, req *msg.Request) (response *msg.Response, err error) {
	var conn net.Conn
	if conn, err = kcp.Dial(address); err != nil {
		return
	}
	if conn.SetDeadline(time.Now().Add(time.Second*5)) != nil {
		return
	}

	if _, err = conn.Write(req.ToBytes()); err != nil {
		return
	}

	close := make(chan bool, 1)
	for {
		select {
		case <-close:
			break
		default:
			message, err := msg.NewMessageFromBytes(conn)
			if err != nil {
				break
			}
			message.Handle(
				nil,
				func(req *msg.Response) (err error) {
					response = req
					return
				},
				func(rep *msg.Close) (err error) {
					close <- true
					return
				},
			)
		}

	}
	return
}
