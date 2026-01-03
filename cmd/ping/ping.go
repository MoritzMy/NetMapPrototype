package ping

import (
	"context"
	"fmt"
	"net"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/MoritzMy/NetMap/cmd/arp_scan"
	icmp2 "github.com/MoritzMy/NetMap/internal/proto/icmp"
	"github.com/MoritzMy/NetMap/internal/proto/ip"
)

// Sweep performs a Ping Sweep over the given List of Network Adresses on the specified network interface.
func Sweep(iface net.Interface) error {
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
			replyChan := icmp2.PingReplyListener(pc, ctx)
			for reply := range replyChan {
				if _, loaded := seen.LoadOrStore(reply.String(), true); loaded {
					continue
				}
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
				err := icmp2.SendPing(pc, ip, id, 0)

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
