package utils

import (
	"github.com/OhYee/blotter/utils/protocol"
	"time"
)

// GlobalVariable golbal variable
type GlobalVariable struct {
	SubServers []proto.Connection
	DeadTime   int64
}

// GV global variable of this server
var GV = &GlobalVariable{
	DeadTime: 120,
}

// RemoveUselessData remove useless global variable (dead connection etc.)
func (gv *GlobalVariable) RemoveUselessData() {
	connections := make([]proto.Connection, len(gv.SubServers))
	var length = 0
	now := time.Now().Unix()
	for _, s := range gv.SubServers {
		if now-s.RecvTime < gv.DeadTime {
			connections[length] = s
			length++
		}
	}
	gv.SubServers = connections
}

// AddConnection add a new sub server connection to the global variable
func (gv *GlobalVariable) AddConnection(connection proto.Connection) {

}
