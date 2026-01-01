package icmp

import "encoding/binary"

const fourBytes = 4

// ICMPHeader represents the ICMP Header structure. For more information see RFC 792.
type ICMPHeader struct {
	Type     uint8
	Code     uint8
	Checksum uint16
}

func NewICMPHeader() *ICMPHeader {
	return &ICMPHeader{
		Type:     0,
		Code:     0,
		Checksum: 0,
	}
}

func (headers *ICMPHeader) Len() int {
	return fourBytes
}

// Marshal parses the ICMP Type and ICMP Code of the IPv4Packet and sets the Checksum Placeholder
func (headers *ICMPHeader) Marshal(payload []byte) ([]byte, error) {
	b := make([]byte, 0)
	b = append(b, headers.Type, headers.Code)
	b = binary.BigEndian.AppendUint16(b, headers.Checksum)

	buf := append(b, payload...)

	cs := computeChecksum(buf)

	binary.BigEndian.PutUint16(b[2:4], cs)

	return b, nil
}

func (headers *ICMPHeader) Unmarshal(b []byte) error {
	headers.Type = b[0]
	headers.Code = b[1]
	headers.Checksum = binary.BigEndian.Uint16(b[2:4])
	return nil
}
