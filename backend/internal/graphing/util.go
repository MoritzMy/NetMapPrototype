package graphing

import (
	"net"
	"strings"

	"github.com/MoritzMy/NetMap/backend/internal/proto/ip"
)

func CreateLocalHostNetworkNodes(g *Graph) {
	ifaces, err := net.Interfaces()
	if err != nil {
		panic(err)
	}

	for _, iface := range ifaces {
		if strings.HasPrefix(iface.Name, "docker") {
			continue
		}
		addrs, err := iface.Addrs()
		if err != nil {
			panic(err)
		}

		for _, addr := range addrs {
			netw := addr.(*net.IPNet)
			if netw.IP.IsLoopback() || netw.IP.To4() == nil {
				continue
			}

			localHost := g.GetOrCreateNode("ip:" + netw.IP.String())

			cannonNetw := ip.CanonicalIPNet(netw)

			netwNode := g.GetOrCreateNode("net:" + cannonNetw.String())

			g.AddEdge(localHost.ID, netwNode.ID, EdgeMemberOf)

		}
	}
}
