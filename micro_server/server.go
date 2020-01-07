package ms

// Server object
type Server struct {
	Info            *ServerInfo
	APIMap          map[string]interface{}
	SubServerStatus []Status // SubServerStatus status of this server
}



// NewServer initial the Server
func NewServer(serverInfo *ServerInfo) *Server {
	return &Server{
		Info:            serverInfo,
		APIMap:          make(map[string]interface{}),
		SubServerStatus: make([]Status, 0),
	}
}


