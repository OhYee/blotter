package ms

import (
	"github.com/OhYee/goutils/bytes"
	"io"
)

// Status of the sub-server
type Status struct {
	Info     *ServerInfo
	SendTime int64
	RecvTime int64
}

// NewStatus initial Status
func NewStatus(serverInfo *ServerInfo, sendTime int64, recvTime int64) *Status {
	return &Status{
		Info:     serverInfo,
		SendTime: sendTime,
		RecvTime: recvTime,
	}
}

// NewStatusFromBytes initial Status from []byte
func NewStatusFromBytes(r io.Reader) (status *Status, err error) {
	var serverInfo *ServerInfo
	if serverInfo, err = NewServerInfoFromBytes(r); err != nil {
		return
	}
	var sendTime int64
	if sendTime, err = bytes.ReadInt64(r); err != nil {
		return
	}

	var recvTime int64
	if recvTime, err = bytes.ReadInt64(r); err != nil {
		return
	}

	status = &Status{
		Info:     serverInfo,
		SendTime: sendTime,
		RecvTime: recvTime,
	}
	return
}

// ToBytes transfer Status to []byte
func (status *Status) ToBytes() []byte {
	buf := bytes.NewBuffer()
	buf.Write(status.Info.ToBytes())
	buf.Write(bytes.FromInt64(status.SendTime))
	buf.Write(bytes.FromInt64(status.RecvTime))
	return buf.Bytes()
}
