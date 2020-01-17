package ms

import (
	// "github.com/OhYee/blotter/micro_server"
	"time"
)

// HeartBeatRequest request of HeartBeatHandle
type HeartBeatRequest struct {
	Name        string
	Address     string
	Description string
	SendTime    int64
}

// HeartBeatResponse response of HeartBeatHandle
type HeartBeatResponse struct{}

// HeartBeatHandle handle of heartbeat
func HeartBeatHandle(server *Server, req HeartBeatRequest, threadID int) (rep HeartBeatResponse) {
	info := NewServerInfo(
		req.Address,
		req.Name,
		req.Description,
	)

	now := time.Now().Unix()

	server.mutex.Lock()
	server.subServerStatus[req.Name] = NewStatus(info, req.SendTime, now)
	for name, status := range server.subServerStatus {
		if now-status.RecvTime > server.deadTime {
			delete(server.subServerStatus, name)
		}
	}
	server.mutex.Unlock()

	server.logCallback(
		threadID, "got heartbeat from %s(%s) at %s",
		req.Name, req.Address, time.Unix(now, 0).Format("2006-01-02 15:04:05"))

	return
}
