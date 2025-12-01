package main

import (
	"fmt"
	"net"
)

var (
	loopbackIPv4 = net.IPv4(127, 0, 0, 1).To4()
)

func main() {
	ifaces, _ := net.Interfaces()

	for _, iface := range ifaces {
		addrs, _ := iface.Addrs()
		for _, addr := range addrs {
			ipNet, ok := addr.(*net.IPNet)

			if !ok {
				continue
			}

			// Avoid the Loopback IP since it's not relevant for scan and any non IPv4 IPs
			if ipNet.Contains(loopbackIPv4) || ipNet.IP.To4() == nil {
				continue
			}
			ip4Addr, ip4Net, err := net.ParseCIDR(ipNet.String())

			if err != nil {
				panic(err)
			}

			fmt.Println(ip4Addr, ip4Net)
		}
	}
}
