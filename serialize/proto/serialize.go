package proto

import (
	"errors"
	"zrpc/serialize"

	"google.golang.org/protobuf/proto"
)

var _ serialize.Serializer = &Serializer{}

type Serializer struct {
}

// Decode implements serialize.Serializer.
func (s *Serializer) Decode(data []byte, val any) error {
	v, ok := val.(proto.Message)
	if !ok {
		return errors.New("类型不符")
	}
	return proto.Unmarshal(data, v)
}

// Encode implements serialize.Serializer.
func (s *Serializer) Encode(val any) ([]byte, error) {
	v, ok := val.(proto.Message)
	if !ok {
		return nil, errors.New("类型不符")
	}
	return proto.Marshal(v)
}

// Id implements serialize.Serializer.
func (s *Serializer) Id() uint8 {
	return 1
}
