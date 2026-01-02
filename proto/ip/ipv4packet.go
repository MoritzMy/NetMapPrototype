package ip

import (
	"fmt"

	"github.com/MoritzMy/NetMap/proto"
)

// IPv4Packet represents an IPv4 Packet with Header and Data sections.
type IPv4Packet struct {
	*Header
	Data []byte
}

func (packet *IPv4Packet) GetHeaders() proto.Header {
	if packet.Header == nil {
		var hdr Header
		packet.Header = &hdr
	}
	return packet.Header
}

func (packet *IPv4Packet) SetHeaders(header proto.Header) {
	h, ok := header.(*Header)
	if !ok {
		panic("inv")
	}
	packet.Header = h
}

func (packet *IPv4Packet) Len() int {
	return int(packet.Header.TotalLength)
}

func (packet *IPv4Packet) Marshal() ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (packet *IPv4Packet) Unmarshal(b []byte) error {
	packet.Data = b[:int(packet.Header.TotalLength)-packet.Header.VersionIHL.Size()]

	return nil
}

func (packet *IPv4Packet) String() string {
	return fmt.Sprintf(
		"%s\n"+
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
		packet.VersionIHL.String(), packet.ToS, packet.TotalLength, packet.Identification, packet.Flags, packet.Fragmentation, packet.TTL, packet.Protocol, packet.Checksum, packet.SourceAddress, packet.DestinationAddress, packet.Options)
}
