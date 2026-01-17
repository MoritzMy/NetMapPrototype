package ping

import (
	"context"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/MoritzMy/NetMap/backend/internal/arp_scan"
	"github.com/MoritzMy/NetMap/backend/internal/graphing"
	"github.com/MoritzMy/NetMap/backend/internal/proto/icmp"
	"github.com/MoritzMy/NetMap/backend/internal/proto/ip"
)

// SweepInterface performs a Ping SweepInterface over the given List of Network Adresses on the specified network interface.
func SweepInterface(iface net.Interface, out chan<- net.IP) error {
	var count atomic.Int64
	ticker := time.NewTicker(time.Millisecond * 10) // Throttle request rate
	defer ticker.Stop()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	addrs, err := iface.Addrs()

	if err != nil {
		return err
	}
	if arp_scan.SumBytes(iface.HardwareAddr) == 0 {
		return fmt.Errorf("interface %s has no MAC address, skipping ARP scan", iface.Name)
	}

	seen := sync.Map{}

	for _, addr := range addrs {
		if addr.(*net.IPNet).IP.To4() == nil {
			continue
		}

		sourceIPNet, ok := addr.(*net.IPNet)

		pc, err := net.ListenPacket("ip4:icmp", sourceIPNet.IP.String())
		if err != nil {
			fmt.Println(err)
			continue
		}
		defer pc.Close()

		go func() {
			replyChan := icmp.PingReplyListener(pc, ctx)
			for reply := range replyChan {
				if _, loaded := seen.LoadOrStore(reply.String(), true); loaded {
					continue
				}
				out <- reply
				fmt.Println("Host", reply, "is up!")
				count.Add(1)
			}
		}()

		if sourceIPNet.IP.IsLoopback() || sourceIPNet.IP.To4() == nil || !ok {
			continue
		}

		var wg sync.WaitGroup

		for _, ip := range ip.ValidIpsInNetwork(sourceIPNet) {
			ip := ip   // Otherwise Routines will use last IP
			<-ticker.C // Throttle

			wg.Go(func() {
				id := uint16(os.Getpid() & 0xffff)
				err := icmp.SendPing(pc, ip, id, 0)

				if err != nil {
					fmt.Println(err)
					return
				}
			})
		}

		wg.Wait()

	}

	drain := time.NewTimer(1 * time.Second) // Wait for late responses
	<-drain.C
	cancel() // Stop listener

	fmt.Println(fmt.Sprintf("Ping Sweep complete, %d hosts are up!", count.Load()))
	return nil
}

func RunICMPSweep(graph *graphing.Graph) {
	ifaces, err := net.Interfaces()
	if err != nil {
		fmt.Println("Error getting network interfaces:", err)
		return
	}
	in := make(chan net.IP)

	go func() {
		for ip := range in {
			fmt.Printf("Discovered host - IP: %s\n", ip)
			node := graph.GetOrCreateNode("ip:" + ip.String())
			node.Protocols["icmp"] = true
		}
		graph.LinkNetworkToGateway()
	}()

	for _, iface := range ifaces {
		if strings.HasPrefix(iface.Name, "docker") {
			continue
		}
		fmt.Printf("Starting ICMP Sweep on interface %s\n", iface.Name)
		if err := SweepInterface(iface, in); err != nil {
			fmt.Printf("Error during ICMP Sweep on interface %s: %v\n", iface.Name, err)
		}
	}
}
