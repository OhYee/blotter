package ms

import (
	"bytes"
	gb "github.com/OhYee/goutils/bytes"
	"io"
	"time"
)

// Status of the sub-server
type Status struct {
	Info     *ServerInfo
	SendTime int64
	RecvTime int64
}

// NewStatus initial Status
func NewStatus(info *ServerInfo) *Status {
	return &Status{
		Info:     info,
		SendTime: time.Now().Unix(),
		RecvTime: 0,
	}
}

// NewStatusFromBytes initial Status from []byte
func NewStatusFromBytes(r io.Reader) (status *Status, err error) {
	var info *ServerInfo
	if info, err = NewServerInfoFromBytes(r); err != nil {
		return
	}
	var sendTime int64
	if sendTime, err = gb.ReadInt64(r); err != nil {
		return
	}
	status = &Status{
		Info:     info,
		SendTime: sendTime,
		RecvTime: time.Now().Unix(),
	}
	return
}

// ToBytes transfer Status to []byte
func (status *Status) ToBytes() []byte {
	buf := bytes.NewBuffer([]byte{})
	buf.Write(status.Info.ToBytes())
	buf.Write(gb.FromInt64(status.SendTime))
	return buf.Bytes()
}
