package icmp

import "encoding/binary"

type ICMPHeader struct {
	Type     uint8
	Code     uint8
	Checksum uint16
}

// Marshal parses the ICMP Type and ICMP Code of the IPv4Packet and sets the Checksum Placeholder
func (packet ICMPHeader) Marshal() ([]byte, error) {
	b := make([]byte, 0)
	b = append(b, packet.Type, packet.Code)
	b = binary.BigEndian.AppendUint16(b, packet.Checksum)

	return b, nil
}

func (headers *ICMPHeader) Unmarshal(b []byte) error {
	headers.Type = b[0]
	headers.Code = b[1]
	headers.Checksum = binary.BigEndian.Uint16(b[2:4])
	return nil
}
