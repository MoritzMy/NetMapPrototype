package arp_scan

import (
	"context"
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

func SendARPRequest(iface net.Interface, targetIP net.IP, fd int) bool {
	addrs, _ := iface.Addrs()

	for _, addr := range addrs {
		sourceIPNet, ok := addr.(*net.IPNet)
		if !ok {
			continue
		}

		if sourceIPNet.IP.To4() == nil {
			continue
		}

		req := arp.NewARPRequest(iface.HardwareAddr, sourceIPNet.IP, targetIP)
		b, err := proto.Marshal(&req)
		if err != nil {
			log.Println("error occurred while marshalling ARP request:", err)
			return false
		}
		err = eth.SendEthernetFrame(b, iface.Name, fd)
		if err != nil {
			log.Println("error occurred while sending ARP request:", err)
			return false
		}

		var hdr eth.EthernetHeader
		var pac arp.ARPRequest
		pac.EthernetHeader = &hdr
	}
	return true

}

func ScanNetwork(iface net.Interface) error {
	addrs, err := iface.Addrs()
	if err != nil {
		return err
	}

	fd, err := eth.CreateSocket(&iface)
	if err != nil {
		return err
	}
	defer syscall.Close(fd)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	ch := arp.ARPResponseListener(fd, ctx)

	go func() {
		for res := range ch {
			log.Println("Received ARP response from", res.IP, "with MAC", res.MAC)
		}
	}()

	var wg sync.WaitGroup
	sem := make(chan struct{}, 256) // Semaphore to limit concurrency

	for _, addr := range addrs {
		sourceIPNet, ok := addr.(*net.IPNet)
		if !ok || sourceIPNet.IP.To4() == nil {
			continue
		}

		for _, ip := range ip.ValidIpsInNetwork(sourceIPNet) {
			sem <- struct{}{} // Acquire semaphore
			ip := ip          // Capture range variable
			wg.Go(func() {
				defer func() { <-sem }() // Release semaphore
				SendARPRequest(iface, ip, fd)
			})
		}

	}

	wg.Wait()
	<-ctx.Done()

	return nil
}
