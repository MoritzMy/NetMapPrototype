package _map

import "net"

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
