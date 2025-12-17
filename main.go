package main

import (
	"fmt"
	"net"
	"sync"

	ping "github.com/MoritzMy/NetMap/sweep/ping"
	"github.com/MoritzMy/NetMap/sweep/ping/icmp"
)

func main() {
	ifaces, _ := net.Interfaces()

	for _, iface := range ifaces {
		addrs, _ := iface.Addrs()
		for _, addr := range addrs {
			if addr.(*net.IPNet).IP.IsLoopback() {
				continue
			}

			fmt.Println(ping.ValidIpsInNetwork(addr.(*net.IPNet)))

			var wg sync.WaitGroup

			for _, ip := range ping.ValidIpsInNetwork(addr.(*net.IPNet)) {
				ip := ip // Otherwise Routines will use last IP

				wg.Add(1)
				go func() {
					defer wg.Done()

					res, err := ping.Ping(ip)
					if res == nil || err != nil {
						return
					}
					var icmpResponse icmp.EchoICMPPacket

					icmp.Unmarshal(res.Data, &icmpResponse)

					fmt.Println(icmpResponse.String(), "\n", res.String())
				}()
			}

			wg.Wait()
		}
	}
}
