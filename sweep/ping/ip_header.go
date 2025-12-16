package ping

const (
	IHLHeaderByteIncrement = 4
	halfByte               = 4
	IPv6HeaderSize         = 40
)

type IPv4VersionIHL struct {
	version uint8
	IHL     uint8
}

func NewIpv4VersionIHL(b byte) IPv4VersionIHL {
	return IPv4VersionIHL{
		version: uint8(b >> halfByte),
		IHL:     uint8(b << halfByte >> halfByte),
	}
}

func (header IPv4VersionIHL) Size() int {
	if header.version != 4 {
		return IPv6HeaderSize // encountered IPv6 Header
	}
	return int(header.IHL) * IHLHeaderByteIncrement
}
