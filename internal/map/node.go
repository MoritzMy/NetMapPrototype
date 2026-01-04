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
	Protocols  map[string]bool
	Vendor     string
	Confidence float64
}

func newNode(id string) Node {
	return Node{
		ID:        id,
		Type:      NodeUnknown,
		Protocols: make(map[string]bool),
		Vendor:    "unknown",
	}
}
