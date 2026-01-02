package arp_scan

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"
	"syscall"
	"time"

	"github.com/MoritzMy/NetMap/proto"
	"github.com/MoritzMy/NetMap/proto/arp"
	eth "github.com/MoritzMy/NetMap/proto/ethernet"
	ip "github.com/MoritzMy/NetMap/proto/ip"
)

// SendARPRequest constructs and sends an ARP request for the given target IP on the specified interface.
func SendARPRequest(iface net.Interface, targetIP net.IP, fd int) bool {
	addrs, _ := iface.Addrs()

	for _, addr := range addrs {
		sourceIPNet, ok := addr.(*net.IPNet)
		if !ok || sourceIPNet.IP.To4() == nil {
			continue
		}

		req := arp.NewARPRequest(iface.HardwareAddr, sourceIPNet.IP, targetIP) // Create ARP request
		b, err := proto.Marshal(&req)
		if err != nil {
			log.Println("error occurred while marshalling ARP request:", err)
			return false
		}
		err = eth.SendEthernetFrame(b, iface.Name, fd) // Send ARP request
		if err != nil {
			log.Println("error occurred while sending ARP request:", err)
			return false
		}

	}
	return true

}

func ScanNetwork(iface net.Interface) error {
	if sumBytes(iface.HardwareAddr) == 0 {
		return fmt.Errorf("interface %s has no MAC address, skipping ARP scan", iface.Name)
	}
	addrs, err := iface.Addrs()
	if err != nil {
		return err
	}

	fd, err := eth.CreateSocket(&iface)
	if err != nil {
		return err
	}
	defer syscall.Close(fd)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ch := arp.ARPResponseListener(fd, ctx)

	go func() {
		for res := range ch {
			log.Println("Received ARP response from", res.IP, "with MAC", res.MAC)
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

	return nil
}

// sumBytes returns the sum of all byte values in the given slice.
func sumBytes(b []byte) int {
	sum := 0
	for _, byteVal := range b {
		sum += int(byteVal)
	}
	return sum
}
