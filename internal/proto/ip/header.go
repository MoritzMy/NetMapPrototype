package ip

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

// Header represents an IPv4 Header structure. For more information see RFC 791.
type Header struct {
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
}

func (header *Header) Marshal(bytes []byte) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (header *Header) Len() int {
	return header.VersionIHL.Size()
}

// Unmarshal provides a way to unmarshal the IP IPv4Packet Payload separately from the Headers
func (header *Header) Unmarshal(b []byte) error {
	if len(b) < MinIPPacketSize {
		return fmt.Errorf("packet smaller than the minimum IP packet size")
	}

	versionIHL := NewIpv4VersionIHL(b[0])

	if len(b) < versionIHL.Size() {
		return fmt.Errorf("packet smaller than defined in the IHL")
	}

	totalLength := binary.BigEndian.Uint16(b[2:4])

	header.VersionIHL = versionIHL
	header.ToS = b[1]
	header.TotalLength = totalLength
	header.Identification = binary.BigEndian.Uint16(b[4:6])
	header.Flags = b[6] >> 5
	header.Fragmentation = binary.BigEndian.Uint16(b[6:8]) << 3 >> 3
	header.TTL = b[8]
	header.Protocol = b[9]
	header.Checksum = binary.BigEndian.Uint16(b[10:12])
	header.SourceAddress = net.IPv4(b[12], b[13], b[14], b[15])
	header.DestinationAddress = net.IPv4(b[16], b[17], b[18], b[19])
	header.Options = b[20:versionIHL.Size()]

	return nil
}
