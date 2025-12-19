package ip

const (
	IHLHeaderByteIncrement = 4
	halfByte               = 4
	MinIPPacketSize        = 20
)

type Marshaler interface {
	Marshal() ([]byte, error)
}

type Unmarshaler interface {
	Unmarshal([]byte) error
}

type IPPacket interface {
	Marshaler
	Unmarshaler
	HeaderSize() int
	Size() int
	GetHeaders() *Header
	SetHeaders(header Header)
}
