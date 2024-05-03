package json

import (
	"encoding/json"
	"test/zrpc/serialize"
)

var _ serialize.Serializer = &Serializer{}

type Serializer struct {
}

// Decode implements serialize.Serializer.
func (s *Serializer) Decode(data []byte, val any) error {
	return json.Unmarshal(data, val)
}

// Encode implements serialize.Serializer.
func (s *Serializer) Encode(val any) ([]byte, error) {
	return json.Marshal(val)
}

// Id implements serialize.Serializer.
func (s *Serializer) Id() uint8 {
	return 0
}
