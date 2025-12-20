package arp

import (
	"encoding/binary"
	"net"

	"github.com/MoritzMy/NetMap/proto/ethernet"
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

type ARPRequest struct {
	EthernetHeader eth.EthernetHeader
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

func NewARPRequest(sourceMAC net.HardwareAddr, sourceIP net.IP, targetIP net.IP) ARPRequest {
	dest := net.HardwareAddr{0xff, 0xff, 0xff, 0xff, 0xff, 0xff}
	targetMAC := net.HardwareAddr{0x00, 0x00, 0x00, 0x00, 0x00, 0x00}

	return ARPRequest{
		EthernetHeader: eth.NewEthernetHeader(dest, sourceMAC, ARPetherType),
		HTYPE:          HTYPEEthernet,
		PTYPE:          PTYPEIPv4,
		HLEN:           MACLength,
		PLEN:           IPv4Length,
		OPER:           OPERRequest,
		SourceMAC:      sourceMAC,
		SourceIP:       sourceIP,
		TargetMAC:      targetMAC,
		TargetIP:       targetIP,
	}

}

func (packet *ARPRequest) GetHeaders() *eth.EthernetHeader {
	return &packet.EthernetHeader
}

func (packet ARPRequest) Marshal() ([]byte, error) {
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

func (packet *ARPRequest) Unmarshal(b []byte) error {
	panic("implement me")
}
