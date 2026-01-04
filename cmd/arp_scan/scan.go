package arp_scan

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/MoritzMy/NetMap/internal/proto"
	arp2 "github.com/MoritzMy/NetMap/internal/proto/arp"
	eth2 "github.com/MoritzMy/NetMap/internal/proto/ethernet"
	"github.com/MoritzMy/NetMap/internal/proto/ip"
	"github.com/endobit/oui"
)

// SendARPRequest constructs and sends an ARP request for the given target IP on the specified interface.
func SendARPRequest(iface net.Interface, targetIP net.IP, fd int) bool {
	addrs, _ := iface.Addrs()

	for _, addr := range addrs {
		sourceIPNet, ok := addr.(*net.IPNet)
		if !ok || sourceIPNet.IP.To4() == nil {
			continue
		}

		req := arp2.NewARPRequest(iface.HardwareAddr, sourceIPNet.IP, targetIP) // Create ARP request
		b, err := proto.Marshal(&req)
		if err != nil {
			log.Println("error occurred while marshalling ARP request:", err)
			return false
		}
		err = eth2.SendEthernetFrame(b, iface.Name, fd) // Send ARP request
		if err != nil {
			log.Println("error occurred while sending ARP request:", err)
			return false
		}

	}
	return true

}

func Scan(iface net.Interface) ([]arp2.Reply, error) {
	if SumBytes(iface.HardwareAddr) == 0 {
		return nil, fmt.Errorf("interface %s has no MAC address, skipping ARP scan", iface.Name)
	}
	addrs, err := iface.Addrs()
	if err != nil {
		return nil, err
	}

	fd, err := eth2.CreateSocket(&iface)
	if err != nil {
		return nil, err
	}
	defer syscall.Close(fd)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ch := arp2.ARPResponseListener(fd, ctx)

	var count atomic.Int64

	replies := make([]arp2.Reply, 0)

	go func() {
		for res := range ch {
			if eth2.IsVRRPMulticastMAC(res.MAC) {
				fmt.Println("VRRP found")
			}
			log.Println("Received ARP response from", res.IP, "with MAC", res.MAC, "and Vendor", oui.VendorFromMAC(res.MAC))
			replies = append(replies, res)
			count.Add(1)
		}
	}()

	var wg sync.WaitGroup
	sem := make(chan struct{}, 256) // Semaphore to limit concurrency

	ticker := time.NewTicker(time.Millisecond * 10) // Throttle request rate
	defer ticker.Stop()

	for _, addr := range addrs {
		sourceIPNet, ok := addr.(*net.IPNet)
		if !ok || sourceIPNet.IP.To4() == nil {
			continue
		}

		for _, ip := range ip.ValidIpsInNetwork(sourceIPNet) {
			if ip.Equal(sourceIPNet.IP) {
				continue
			}
			<-ticker.C        // Throttle
			sem <- struct{}{} // Acquire semaphore
			ip := ip          // Capture range variable
			wg.Go(func() {
				defer func() { <-sem }() // Release semaphore
				SendARPRequest(iface, ip, fd)
			})
		}

	}

	wg.Wait()

	drain := time.NewTimer(1 * time.Second) // Wait for late responses
	<-drain.C
	cancel() // Stop listener

	fmt.Println(fmt.Sprintf("%d ARP packets received", count.Load()))

	return replies, nil
}

// SumBytes returns the sum of all byte values in the given slice.
func SumBytes(b []byte) int {
	sum := 0
	for _, byteVal := range b {
		sum += int(byteVal)
	}
	return sum
}
