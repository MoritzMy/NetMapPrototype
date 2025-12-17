package ping

import (
	"encoding/binary"
	"fmt"
	"net"
)

const (
	IHLHeaderByteIncrement = 4
	halfByte               = 4
	MinIPPacketSize        = 20
)

type IPPacket struct {
	VersionIHL         IPv4VersionIHL
	ToS                uint8
	TotalLength        uint16
	Identification     uint16
	Flags              uint8 // 3 Bit
	Fragmentation      uint16
	TTL                uint8
	Protocol           uint8
	Checksum           uint16
	SourceAddress      net.IP
	DestinationAddress net.IP
	Options            []byte
	Data               []byte
}

func NewIPPacket(b []byte) *IPPacket {

	if len(b) < MinIPPacketSize {
		return nil
	}

	versionIHL := NewIpv4VersionIHL(b[0])

	if len(b) < versionIHL.Size() {
		return nil
	}

	totalLength := binary.BigEndian.Uint16(b[2:4])

	return &IPPacket{
		VersionIHL:         versionIHL,
		ToS:                b[1],
		TotalLength:        totalLength,
		Identification:     binary.BigEndian.Uint16(b[4:6]),
		Flags:              b[6] >> 5,
		Fragmentation:      binary.BigEndian.Uint16(b[6:8]) << 3 >> 3,
		TTL:                b[8],
		Protocol:           b[9],
		Checksum:           binary.BigEndian.Uint16(b[10:12]),
		SourceAddress:      net.IPv4(b[12], b[13], b[14], b[15]),
		DestinationAddress: net.IPv4(b[16], b[17], b[18], b[19]),
		Options:            b[20:versionIHL.Size()],
		Data:               b[versionIHL.Size():totalLength],
	}
}

func (packet *IPPacket) String() string {
	return fmt.Sprintf(
		"VersionIHL: %d\n"+
			"ToS: %d\n"+
			"TotalLength: %d\n"+
			"Identification: %d\n"+
			"Flags: %d\n"+
			"Fragmentation: %d\n"+
			"TTL: %d\n"+"Protocol: %d\n"+
			"Checksum: %d\n"+
			"SourceAddress: %s\n"+
			"DestinationAddress: %s\n"+
			"Options:%x\n\n",
		packet.VersionIHL, packet.ToS, packet.TotalLength, packet.Identification, packet.Flags, packet.Fragmentation, packet.TTL, packet.Protocol, packet.Checksum, packet.SourceAddress, packet.DestinationAddress, packet.Options)
}

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
