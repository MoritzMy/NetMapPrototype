package main

import (
	"fmt"
	"net"

	ping "github.com/MoritzMy/NetMap/sweep/ping"
)

func main() {

	ifaces, _ := net.Interfaces()

	for _, iface := range ifaces {
		addrs, _ := iface.Addrs()
		for _, addr := range addrs {

			// TODO: Currently pings self, need to get the network range

			// TODO: Fix Up ping.ValidateIP
			ping.ValidateIP(addr.(*net.IPNet))
			ipNet, ok := addr.(*net.IPNet)

			if !ok {
				continue
			}

			// Avoid the Loopback IP since it's not relevant for scan and any non IPv4 IPs
			if ipNet.IP.IsLoopback() || ipNet.IP.To4() == nil || !ipNet.IP.IsGlobalUnicast() || ipNet.IP.IsMulticast() {
				continue
			}

			res := ping.Ping(ipNet)

			fmt.Println(res)
		}
	}
}
