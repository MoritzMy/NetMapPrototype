package ip

import (
	"fmt"
)

type IPv4Packet struct {
	Header
	Data []byte
}

func (packet *IPv4Packet) Marshal() ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (packet *IPv4Packet) Unmarshal(b []byte) error {
	packet.Data = b[packet.HeaderSize():packet.Size()]

	return nil
}

func (packet *IPv4Packet) GetHeaders() *Header {
	return &packet.Header
}

func (packet IPv4Packet) SetHeaders(header Header) {
	packet.Header = header
}

func (packet *IPv4Packet) HeaderSize() int {
	return int(packet.Header.VersionIHL.IHL)
}

func (packet *IPv4Packet) Size() int {
	return int(packet.Header.TotalLength)
}

func (packet *IPv4Packet) String() string {
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
