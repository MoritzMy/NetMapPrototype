package icmp

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/MoritzMy/NetMap/proto"
)

type EchoICMPPacket struct {
	*ICMPHeader
	Identifier     uint16
	SequenceNumber uint16
	Payload        []byte
}

func (packet *EchoICMPPacket) GetHeaders() proto.Header {
	return packet.ICMPHeader
}

func (packet *EchoICMPPacket) SetHeaders(header proto.Header) {
	h, ok := header.(*ICMPHeader)
	if !ok {
		panic(fmt.Sprintf("can't convert %v to ICMP Header", header))
	}
	packet.ICMPHeader = h
}

func NewEchoICMPPacket(identifier uint16, sequenceNumber uint16, payload []byte) EchoICMPPacket {
	return EchoICMPPacket{
		ICMPHeader: &ICMPHeader{
			Type: echoType,
			Code: echoCode,
		},
		Identifier:     identifier,
		SequenceNumber: sequenceNumber,
		Payload:        payload,
	}
}

func (packet EchoICMPPacket) Equal(other_pkg EchoICMPPacket) bool {
	return packet.Type == other_pkg.Type &&
		packet.Code == other_pkg.Code &&
		packet.Identifier == other_pkg.Identifier &&
		packet.SequenceNumber == other_pkg.SequenceNumber &&
		bytes.Equal(packet.Payload, other_pkg.Payload)
}

func (packet EchoICMPPacket) String() string {
	const maxPreview = 16

	preview := packet.Payload
	if len(preview) > maxPreview {
		preview = preview[:maxPreview]
	}

	return fmt.Sprintf(
		"ICMP Echo (type=%d code=%d id=%d seq=%d payload_len=%d payload=%q)",
		packet.Type,
		packet.Code,
		packet.Identifier,
		packet.SequenceNumber,
		len(packet.Payload),
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

func (packet EchoICMPPacket) Clone() EchoICMPPacket {
	c := packet
	c.Payload = append([]byte(nil), packet.Payload...)
	return c
}
