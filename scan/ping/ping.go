package ping

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/MoritzMy/NetMap/proto"
	"github.com/MoritzMy/NetMap/proto/icmp"
	"github.com/MoritzMy/NetMap/proto/ip"
	"github.com/MoritzMy/NetMap/scan/arp_scan"
)

const (
	echoReplyType = 0
)

// Sweep performs a Ping Sweep over the given List of Network Adresses
func Sweep(iface net.Interface) error {
	count := 0
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
			replyChan := PingReplyListener(pc, ctx)
			for reply := range replyChan {
				fmt.Println("Host", reply, "is up!")
				count++
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
				err := SendPing(pc, ip, 0, 0)

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
	ctx.Done() // Stop listener

	fmt.Println(fmt.Sprintf("Ping Sweep complete, %d hosts are up!", count))
	return nil
}

func SendPing(conn net.PacketConn, dst net.IP, id, seq uint16) error {
	req := icmp.NewEchoICMPPacket(id, seq, []byte("ARE U UP?"))
	b, err := proto.Marshal(&req)
	if err != nil {
		return err
	}

	_, err = conn.WriteTo(b, &net.IPAddr{IP: dst})
	return err
}

func PingReplyListener(conn net.PacketConn, ctx context.Context) <-chan net.IP {
	ch := make(chan net.IP)
	buf := make([]byte, 200)

	go func() {
		defer close(ch)
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}

			_, addr, err := conn.ReadFrom(buf)
			if err != nil {
				continue
			}

			var packet icmp.EchoICMPPacket

			if err := proto.Unmarshal(buf, &packet); err != nil {
				continue
			}

			if packet.Type != echoReplyType {
				continue
			}

			ch <- addr.(*net.IPAddr).IP
		}
	}()
	return ch
}
