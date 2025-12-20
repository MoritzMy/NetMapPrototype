package arp

import (
	"net"

	"github.com/MoritzMy/NetMap/proto/ethernet"
)

const (
	ARPetherType  = 0x0806
	HTYPEEthernet = 1
	PTYPEIPv4     = 0x0800
	MACLength     = 6
	IPv4Length    = 4
	OPERRequest   = 1
	OPERResponse  = 2
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
