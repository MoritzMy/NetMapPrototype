package graphing

import (
	"fmt"
	"net"
	"strings"

	"github.com/MoritzMy/NetMap/backend/internal/proto/ethernet"
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
	ID         string           `json:"id"`
	Type       NodeType         `json:"type"`
	IP         net.IP           `json:"ip"`
	MAC        net.HardwareAddr `json:"mac"`
	Protocols  map[string]bool  `json:"protocols"`
	Vendor     string           `json:"vendor"`
	Confidence float64          `json:"confidence"`
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

	n.Confidence += float64(NetworkVendorScore(n.Vendor)) * 0.2

	if n.Confidence >= 1.0 && n.Type != NodeHost {
		n.Type = NodeGateway
	} else if n.Type != NodeNetwork {
		n.Type = NodeHost
	}

	fmt.Printf(n.String())

}

var networkVendorWeights = map[string]int{
	"palo alto":  5,
	"fortinet":   5,
	"checkpoint": 5,
	"sophos":     5,
	"sonicwall":  5,
	"watchguard": 5,
	"untangle":   5,
	"pfSense":    5,
	"opnsense":   5,

	"cisco":    4,
	"juniper":  4,
	"arista":   4,
	"vyos":     4,
	"mikrotik": 4,

	"aruba":            3,
	"aruba networks":   3,
	"hewlett":          3,
	"hpe":              3,
	"dell":             3,
	"extreme networks": 3,
	"brocade":          3,
	"ruckus":           3,

	"ubiquiti":          2,
	"ubiquiti networks": 2,
	"netgate":           2,
	"meraki":            2,
	"peplink":           2,
	"cradlepoint":       2,
}

func NetworkVendorScore(vendor string) int {
	v := strings.ToLower(vendor)
	score := 0

	for token, weight := range networkVendorWeights {
		if strings.Contains(v, token) {
			score += weight
			break
		}
	}
	return score
}

func (n *Node) String() string {
	return fmt.Sprintf("Node: %s\nIP: %v\nMAC:%v\nTYPE: %d", n.ID, n.IP, n.MAC, n.Type)
}
