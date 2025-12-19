package icmp

import (
	"encoding/binary"
	"fmt"

	ip2 "github.com/MoritzMy/NetMap/proto/ip"
)

type QuotedPacket struct {
	Header  ip2.Header
	Payload [8]byte
}

// Unmarshal takes a byte array of a Time Exceeded ICMP IPv4Packet and fills the fields of the Object that called
// the function. The given byte array must start right after the default ICMP Headers (start of the "Unused" Field)
func (packet *TimeExceededPacket) Unmarshal(b []byte) error {
	if len(b) < 4 {
		return fmt.Errorf("time exceeded packet too small to unmarshal")
	}
	packet.Unused = binary.BigEndian.Uint32(b[0:4])

	var ipPacket ip2.IPv4Packet
	packet.IPv4Packet = ip2.Unmarshal(b[4:], ipPacket)
	return nil
}

func (packet TimeExceededPacket) Marshal() ([]byte, error) {
	return nil, nil
}

func (packet *TimeExceededPacket) GetHeaders() *ICMPHeader {
	return &packet.ICMPHeader
}

func (packet *TimeExceededPacket) SetHeaders(header ICMPHeader) {
	packet.ICMPHeader = header
}

type TimeExceededPacket struct {
	ICMPHeader
	Unused uint32
	ip2.IPv4Packet
}
