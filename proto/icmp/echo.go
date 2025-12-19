package icmp

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type EchoICMPPacket struct {
	ICMPHeader
	Identifier     uint16
	SequenceNumber uint16
	Payload        []byte
}

func NewEchoICMPPacket(identifier uint16, sequenceNumber uint16, payload []byte) EchoICMPPacket {
	return EchoICMPPacket{
		ICMPHeader: ICMPHeader{
			Type: echoType,
			Code: echoCode,
		},
		Identifier:     identifier,
		SequenceNumber: sequenceNumber,
		Payload:        payload,
	}
}

func (req *EchoICMPPacket) GetHeaders() *ICMPHeader {
	return &req.ICMPHeader
}

func (req *EchoICMPPacket) SetHeaders(header ICMPHeader) {
	req.ICMPHeader = header
}

func (pkg *EchoICMPPacket) Equal(other_pkg EchoICMPPacket) bool {
	return pkg.Type == other_pkg.Type &&
		pkg.Code == other_pkg.Code &&
		pkg.Checksum == other_pkg.Checksum &&
		pkg.Identifier == other_pkg.Identifier &&
		pkg.SequenceNumber == other_pkg.SequenceNumber &&
		bytes.Equal(pkg.Payload, other_pkg.Payload)
}

func (p EchoICMPPacket) String() string {
	const maxPreview = 16

	preview := p.Payload
	if len(preview) > maxPreview {
		preview = preview[:maxPreview]
	}

	return fmt.Sprintf(
		"ICMP Echo (type=%d code=%d id=%d seq=%d payload_len=%d payload=%q)",
		p.Type,
		p.Code,
		p.Identifier,
		p.SequenceNumber,
		len(p.Payload),
		preview,
	)
}

func (packet EchoICMPPacket) Marshal() ([]byte, error) {
	if len(packet.Payload) > maxPayload {
		return nil, fmt.Errorf("marshal icmp request: payload size %d exceeds limit of %d Bytes", len(packet.Payload), maxPayload)
	}

	b := make([]byte, 0, echoHeaderSize+len(packet.Payload))
	b = binary.BigEndian.AppendUint16(b, packet.Identifier)
	b = binary.BigEndian.AppendUint16(b, packet.SequenceNumber)
	b = append(b, packet.Payload...)

	return b, nil
}

func (packet *EchoICMPPacket) Unmarshal(data []byte) error {
	packet.Identifier = binary.BigEndian.Uint16(data[0:2])
	packet.SequenceNumber = binary.BigEndian.Uint16(data[2:4])
	packet.Payload = data[4:]
	return nil
}
