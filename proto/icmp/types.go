package icmp

const (
	echoCode       = 0
	echoType       = 8
	maxPayload     = 56
	echoHeaderSize = 4
)

type Marshaler interface {
	Marshal() ([]byte, error)
}

type Unmarshaler interface {
	Unmarshal([]byte) error
}

type ICMPPacket interface {
	Marshaler
	Unmarshaler
	GetHeaders() *ICMPHeader
	SetHeaders(header ICMPHeader)
}
