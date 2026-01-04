package _map

import (
	"net"
)

type NodeType int

const (
	NodeNetwork NodeType = iota
	NodeHost
	NodeGateway
	NodeUnknown
)

type Node struct {
	ID         string
	Type       NodeType
	IP         net.IP
	MAC        net.HardwareAddr
	Vendor     string
	Confidence float64
}

func newNode(ip net.IP) Node {
	return Node{
		ID:     "ip:" + ip.String(),
		Type:   NodeUnknown,
		IP:     ip,
		Vendor: "unknown",
	}
}
