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
			sourceIPNet, ok := addr.(*net.IPNet)

			if sourceIPNet.IP.IsLoopback() || sourceIPNet.IP.To4() == nil || !ok {
				continue
			}

			var wg sync.WaitGroup

			for _, ip := range ip.ValidIpsInNetwork(sourceIPNet) {
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

					if err := proto.Unmarshal(res.Data, &icmpResponse); err != nil {
						return
					}

					fmt.Println(fmt.Sprintf("Host %s is up!", ip.String()))
				}()
			}

			wg.Wait()
		}
	}
}
