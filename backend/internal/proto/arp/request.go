package arp

import (
	"encoding/binary"
	"errors"
	"fmt"
	"net"

	"github.com/MoritzMy/NetMap/backend/internal/proto"
	"github.com/MoritzMy/NetMap/backend/internal/proto/ethernet"
)

const (
	ARPetherType          = 0x0806
	HTYPEEthernet         = 1
	PTYPEIPv4             = 0x0800
	MACLength             = 6
	IPv4Length            = 4
	OPERRequest           = 1
	OPERResponse          = 2
	ARPRequestPayloadSize = 28
)

// Request represents an ARP Request Packet structure. For more information see RFC 826.
type Request struct {
	EthernetHeader *eth.EthernetHeader
	HTYPE          uint16
	PTYPE          uint16
	HLEN           uint8
	PLEN           uint8
	OPER           uint16
	SourceMAC      net.HardwareAddr
	SourceIP       net.IP
	TargetMAC      net.HardwareAddr
	TargetIP       net.IP
}

func (packet *Request) GetHeaders() proto.Header {
	return packet.EthernetHeader
}

func (packet *Request) SetHeaders(header proto.Header) {
	hdr, ok := header.(*eth.EthernetHeader)
	if !ok {
		panic("Wrong Header for ARP")
	}
	packet.EthernetHeader = hdr
}

func (packet *Request) Len() int {
	return 28
}

func NewARPRequest(sourceMAC net.HardwareAddr, sourceIP net.IP, targetIP net.IP) Request {
	dest := net.HardwareAddr{0xff, 0xff, 0xff, 0xff, 0xff, 0xff}
	targetMAC := net.HardwareAddr{0x00, 0x00, 0x00, 0x00, 0x00, 0x00}

	return Request{
		EthernetHeader: eth.NewEthernetHeader(dest, sourceMAC, ARPetherType),
		HTYPE:          HTYPEEthernet,
		PTYPE:          PTYPEIPv4,
		HLEN:           MACLength,
		PLEN:           IPv4Length,
		OPER:           OPERRequest,
		SourceMAC:      sourceMAC,
		SourceIP:       sourceIP.To4(),
		TargetMAC:      targetMAC,
		TargetIP:       targetIP.To4(),
	}

}

func (packet *Request) Marshal() ([]byte, error) {
	if len(packet.SourceMAC) != 6 || len(packet.TargetMAC) != 6 {
		return nil, errors.New(fmt.Sprintf("invalid MAC length : %v or %v are faulty", packet.SourceMAC, packet.TargetMAC))
	}
	if len(packet.SourceIP) != 4 || len(packet.TargetIP) != 4 {
		return nil, errors.New(fmt.Sprintf("invalid IP length: %v or %v are faulty", packet.SourceIP, packet.TargetIP))
	}

	b := make([]byte, 0, ARPRequestPayloadSize)
	b = binary.BigEndian.AppendUint16(b, packet.HTYPE)
	b = binary.BigEndian.AppendUint16(b, packet.PTYPE)
	b = append(b, packet.HLEN, packet.PLEN)
	b = binary.BigEndian.AppendUint16(b, packet.OPER)
	b = append(b, packet.SourceMAC...)
	b = append(b, packet.SourceIP...)
	b = append(b, packet.TargetMAC...)
	b = append(b, packet.TargetIP...)

	return b, nil
}

func (packet *Request) Unmarshal(b []byte) error {
	packet.HTYPE = binary.BigEndian.Uint16(b[0:2])
	packet.PTYPE = binary.BigEndian.Uint16(b[2:4])
	packet.HLEN = b[4]
	packet.PLEN = b[5]
	packet.OPER = binary.BigEndian.Uint16(b[6:8])
	packet.SourceMAC = b[8:14]
	packet.SourceIP = b[14:18]
	packet.TargetMAC = b[18:24]
	packet.TargetIP = b[24:28]

	return nil
}
