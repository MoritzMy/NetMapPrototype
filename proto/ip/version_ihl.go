package ip

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
		return 0 // encountered non IPv4 Header
	}
	return int(header.IHL) * IHLHeaderByteIncrement
}
