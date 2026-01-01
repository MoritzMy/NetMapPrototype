package proto

type Marshaler interface {
	Marshal() ([]byte, error)
}

type HeaderMarshaler interface {
	Marshal([]byte) ([]byte, error)
}

type Unmarshaler interface {
	Unmarshal([]byte) error
}
type Header interface {
	HeaderMarshaler
	Unmarshaler
	Len() int
}

type Packet interface {
	Marshaler
	Unmarshaler
	GetHeaders() Header
	SetHeaders(Header)
	Len() int
}
