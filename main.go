package main

import (
	"fmt"
	"net"

	ping "github.com/MoritzMy/NetMap/sweep/ping"
	"github.com/MoritzMy/NetMap/sweep/ping/icmp"
)

func main() {
	ifaces, _ := net.Interfaces()

	for _, iface := range ifaces {
		addrs, _ := iface.Addrs()
		for _, addr := range addrs {
			for _, ip := range ping.ValidIpsInNetwork(addr.(*net.IPNet)) {
				// Avoid the Loopback IP since it's not relevant for scan and any non IPv4 IPs
				if ip.IsLoopback() || ip.To4() == nil || !ip.IsGlobalUnicast() || ip.IsMulticast() {
					continue
				}

				res := ping.Ping(ip)

				var icmpResponse icmp.EchoICMPPacket

				icmp.Unmarshal(res.Data, &icmpResponse)

				fmt.Println(icmpResponse, res.String())
			}
		}
	}
}
