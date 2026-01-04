package graphing

import (
	"fmt"
	"net"

	eth "github.com/MoritzMy/NetMap/internal/proto/ethernet"
	"github.com/endobit/oui"
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

func newNode(id string) *Node {
	return &Node{
		ID:        id,
		Type:      NodeUnknown,
		Protocols: make(map[string]bool),
		Vendor:    "unknown",
	}
}

func (n *Node) EnrichNode() {
	n.Vendor = oui.VendorFromMAC(n.MAC)

	if eth.IsVRRPMulticastMAC(n.MAC) {
		fmt.Println("Found VRRP")
		n.Confidence += 1.0
	}

	if IsLikelyNetworkVendor(n.Vendor) {
		n.Confidence += 0.6
	}

	if n.Confidence > 1.0 && n.Type == NodeUnknown {
		n.Type = NodeGateway
	} else if n.Type != NodeNetwork {
		n.Type = NodeHost
	}

	fmt.Printf(n.String())

}
func IsLikelyNetworkVendor(vendor string) bool {
	switch vendor {
	case "Cisco", "Juniper", "Arista", "Palo Alto", "Fortinet", "MikroTik", "Ubiquiti", "Dell", "HP", "Aruba Networks":
		return true
	default:
		return false
	}
}

func (n *Node) String() string {
	return fmt.Sprintf("Node: %s\nIP: %v\nMAC:%v\nTYPE: %d", n.ID, n.IP, n.MAC, n.Type)
}
