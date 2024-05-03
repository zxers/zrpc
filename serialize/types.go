package serialize

type Serializer interface {
	Encode(any) ([]byte, error)
	Decode([]byte, any) error
	Id() uint8
}