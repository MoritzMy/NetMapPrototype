package ping

import (
	"fmt"
	"net"
	"sync"

	"github.com/MoritzMy/NetMap/proto"
	"github.com/MoritzMy/NetMap/proto/icmp"
	"github.com/MoritzMy/NetMap/proto/ip"
)

// Sweep performs a Ping Sweep over the given List of Network Adresses
func Sweep(ifaces []net.Interface) {
	for _, iface := range ifaces {
		addrs, err := iface.Addrs()
		if err != nil {
			fmt.Println(err)
			continue
		}
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
					if res == nil {
						return
					}

					if err != nil {
						fmt.Println(err)
						return
					}
					var icmpResponse icmp.EchoICMPPacket

					fmt.Println(res.Data)

					if err := proto.Unmarshal(res.Data, &icmpResponse); err != nil {
						return
					}
					fmt.Println(fmt.Sprintf("%s\n%s", icmpResponse.String(), res.String()))
				}()
			}

			wg.Wait()
		}
	}
}
