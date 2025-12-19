package ping

import (
	"fmt"
	"net"
	"sync"

	"github.com/MoritzMy/NetMap/proto/icmp"
	"github.com/MoritzMy/NetMap/proto/ip"
)

//TODO: Seperate the Network Interface IP extraction Logic from Ping into seperate file

// Sweep performs a Ping Sweep over the active NetworkInterfaces
func Sweep() {
	ifaces, _ := net.Interfaces()

	for _, iface := range ifaces {
		addrs, _ := iface.Addrs()
		for _, addr := range addrs {
			if addr.(*net.IPNet).IP.IsLoopback() {
				continue
			}

			fmt.Println(ip.ValidIpsInNetwork(addr.(*net.IPNet)))

			var wg sync.WaitGroup

			for _, ip := range ip.ValidIpsInNetwork(addr.(*net.IPNet)) {
				ip := ip // Otherwise Routines will use last IP

				wg.Add(1)
				go func() {
					defer wg.Done()

					res, err := Ping(ip)
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
