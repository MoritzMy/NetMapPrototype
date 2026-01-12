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

	graphing2 "github.com/MoritzMy/NetMap/backend/internal/graphing"
	"github.com/MoritzMy/NetMap/backend/internal/proto"
	"github.com/MoritzMy/NetMap/backend/internal/proto/arp"
	eth2 "github.com/MoritzMy/NetMap/backend/internal/proto/ethernet"
	"github.com/MoritzMy/NetMap/backend/internal/proto/ip"
)

// sendARPRequest constructs and sends an ARP request for the given target IP on the specified interface.
func sendARPRequest(iface net.Interface, targetIP net.IP, fd int) bool {
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
		err = eth2.SendEthernetFrame(b, iface.Name, fd) // Send ARP request
		if err != nil {
			log.Println("error occurred while sending ARP request:", err)
			return false
		}

	}
	return true

}

func ScanInterface(iface net.Interface, out chan<- ARPEvent) error {
	if SumBytes(iface.HardwareAddr) == 0 {
		return fmt.Errorf("interface %s has no MAC address, skipping ARP scan", iface.Name)
	}
	addrs, err := iface.Addrs()
	if err != nil {
		return err
	}

	fd, err := eth2.CreateSocket(&iface)
	if err != nil {
		return err
	}
	defer syscall.Close(fd)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ch := arp.ARPResponseListener(fd, ctx)

	var count atomic.Int64

	var wg sync.WaitGroup
	sem := make(chan struct{}, 256) // Semaphore to limit concurrency

	ticker := time.NewTicker(time.Millisecond * 10) // Throttle request rate
	defer ticker.Stop()

	for _, addr := range addrs {
		sourceIPNet, ok := addr.(*net.IPNet)
		if !ok || sourceIPNet.IP.To4() == nil {
			continue
		}

		go func(netw *net.IPNet) {
			for res := range ch {
				canonNetw := CanonicalIPNet(netw)

				out <- ARPEvent{
					IP:      res.IP,
					MAC:     res.MAC,
					Network: canonNetw,
					Source:  "arp",
				}
				count.Add(1)
			}
		}(sourceIPNet)

		for _, ip := range ip.ValidIpsInNetwork(sourceIPNet) {
			<-ticker.C        // Throttle
			sem <- struct{}{} // Acquire semaphore
			ip := ip          // Capture range variable
			wg.Go(func() {
				defer func() { <-sem }() // Release semaphore
				sendARPRequest(iface, ip, fd)
			})
		}

	}

	wg.Wait()

	drain := time.NewTimer(1 * time.Second) // Wait for late responses
	<-drain.C
	cancel() // Stop listener

	fmt.Println(fmt.Sprintf("%d ARP packets received", count.Load()))

	return nil
}

// SumBytes returns the sum of all byte values in the given slice.
func SumBytes(b []byte) int {
	sum := 0
	for _, byteVal := range b {
		sum += int(byteVal)
	}
	return sum
}

func RunARPScan(graph *graphing2.Graph) {
	ifaces, err := net.Interfaces()
	if err != nil {
		fmt.Println("Error getting network interfaces:", err)
		return
	}

	in := make(chan ARPEvent)

	go func() {
		for ev := range in {
			fmt.Printf("Discovered device - IP: %s, MAC: %s, Network: %s, Source: %s\n", ev.IP, ev.MAC, ev.Network, ev.Source)
			node := graph.GetOrCreateNode("ip:" + ev.IP.String())
			node.MAC = ev.MAC
			node.IP = ev.IP
			node.Protocols["arp"] = true
			netNode := graph.GetOrCreateNode("net:" + ev.Network.String())
			netNode.Type = graphing2.NodeNetwork
			graph.AddEdge(node.ID, netNode.ID, graphing2.EdgeMemberOf)
		}

		for node := range graph.Nodes {
			graph.GetOrCreateNode(node).EnrichNode() // Enrich nodes with additional information
		}

		graph.LinkNetworkToGateway()

	}()

	for _, iface := range ifaces {
		fmt.Printf("Starting ARP ScanInterface on interface %s\n", iface.Name)
		err := ScanInterface(iface, in)
		if err != nil {
			fmt.Printf("Error scanning network on interface %s: %v\n", iface.Name, err)
			continue
		}
	}
}
