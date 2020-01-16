package ms

import (
	"github.com/OhYee/goutils/bytes"
	"io"
)

// ServerInfo information of the sub-server
type ServerInfo struct {
	Address     string // Address of the sub-server, unique
	Name        string // Name of the sub-server, servers could use the same name for balanced
	Description string // Description of the sub-server
	// APIList     []APIInfo // APIList api of the server
}

// NewServerInfo initial the ServerInfo
func NewServerInfo(address string, name string, description string,
	/*apiList []APIInfo*/) *ServerInfo {
	return &ServerInfo{
		Address:     address,
		Name:        name,
		Description: description,
		// APIList:     apiList,
	}
}

// NewServerInfoFromBytes initial the ServerInfo from bytes
func NewServerInfoFromBytes(r io.Reader) (info *ServerInfo, err error) {
	var address, name, description string
	if address, err = bytes.ReadStringWithLength32(r); err != nil {
		return
	}
	if name, err = bytes.ReadStringWithLength32(r); err != nil {
		return
	}
	if description, err = bytes.ReadStringWithLength32(r); err != nil {
		return
	}
	// var size uint32
	// if size, err = gb.ReadUint32(r); err != nil {
	// 	return
	// }
	// apiList := make([]APIInfo, size)
	// for i := uint32(0); i < size; i++ {
	// 	var api APIInfo
	// 	if api, err = NewAPIInfoFromBytes(r); err != nil {
	// 		return
	// 	}
	// 	apiList[i] = api
	// }
	info = NewServerInfo(address, name, description,
	/*apiList*/)
	return
}

// ToBytes transfer ServerInfo to []byte
func (info *ServerInfo) ToBytes() []byte {
	buf := bytes.NewBuffer()
	bytes.FromStringWithLength32(info.Address)
	bytes.FromStringWithLength32(info.Name)
	bytes.FromStringWithLength32(info.Description)
	// buf.Write(gb.FromUint32(uint32(len(info.APIList))))
	// for _, api := range info.APIList {
	// 	buf.Write(api.ToBytes())
	// }
	return buf.Bytes()
}
