package msg

import (
	"bytes"
	"github.com/OhYee/blotter/micro_server"
	gb "github.com/OhYee/goutils/bytes"
	"io"
	"time"
)

// HeartBeat message body
type HeartBeat struct {
	Info     *ms.ServerInfo
	SendTime int64
	RecvTime int64
}

// NewHeartBeat initial a HeartBeat message
func NewHeartBeat(info *ms.ServerInfo) *HeartBeat {
	return &HeartBeat{
		Info:     info,
		SendTime: time.Now().Unix(),
		RecvTime: 0,
	}
}

// NewHeartBeatFromBytes initial a HeartBeat from bytes
func NewHeartBeatFromBytes(r io.Reader) (heartbeat *HeartBeat, err error) {
	var info *ms.ServerInfo
	var sendTime, recvTime int64
	if info, err = ms.NewServerInfoFromBytes(r); err != nil {
		return
	}
	if sendTime, err = gb.ReadInt64(r); err != nil {
		return
	}
	if recvTime, err = gb.ReadInt64(r); err != nil {
		return
	}
	heartbeat = &HeartBeat{
		Info:     info,
		SendTime: sendTime,
		RecvTime: recvTime,
	}
	return
}

// ToBytes transfer HeartBeat to []byte
func (heartbeat *HeartBeat) ToBytes() []byte {
	buf := bytes.NewBuffer([]byte{})
	buf.Write(heartbeat.Info.ToBytes())
	buf.Write(gb.FromInt64(heartbeat.SendTime))
	buf.Write(gb.FromInt64(heartbeat.RecvTime))
	return buf.Bytes()
}

// ToMessage initial a HeartBeat message
func (heartbeat *HeartBeat) ToMessage() *Message {
	return NewMessage(MessageTypeHeartBeat, heartbeat)
}
