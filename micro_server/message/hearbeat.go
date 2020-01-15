package msg

import (
	"bytes"
	"github.com/OhYee/blotter/micro_server"
	gb "github.com/OhYee/goutils/bytes"
	"time"
)

// HeartBeat message body
type HeartBeat struct {
	Info     *ms.ServerInfo
	SendTime int64
	RecvTime int64
}

// NewHeartBeatMessage initial a HeartBeat message
func NewHeartBeatMessage(info *ms.ServerInfo) Message {
	heartbeat := &HeartBeat{
		Info:     info,
		SendTime: time.Now().Unix(),
		RecvTime: 0,
	}
	return NewMessage(MessageTypeHeartBeat, heartbeat.ToBytes())
}

// ToBytes transfer HeartBeat to []byte
func (heartbeat *HeartBeat) ToBytes() []byte {
	buf := bytes.NewBuffer([]byte{})
	buf.Write(heartbeat.Info.ToBytes())
	buf.Write(gb.FromInt64(heartbeat.SendTime))
	buf.Write(gb.FromInt64(heartbeat.RecvTime))
	return buf.Bytes()
}
