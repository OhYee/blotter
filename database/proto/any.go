package proto

import (
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
)

func FromInterfaceToAny(b []byte) (*any.Any, error) {
	return ptypes.MarshalAny(&any.Any{
		TypeUrl: "oyohyee.com",
		Value:   b,
	})
}
