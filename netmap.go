package main

import (
	"flag"
	"fmt"
	"net"

	"github.com/MoritzMy/NetMap/cmd/arp_scan"
	"github.com/MoritzMy/NetMap/cmd/ping"
	_map "github.com/MoritzMy/NetMap/internal/map"
)

func main() {
	arp := flag.Bool("arp-scan", false, "Run ARP Discovery Scan")
	icmp := flag.Bool("ping-sweep", false, "Run ICMP Sweep")

	flag.Parse()

	graph := _map.NewGraph()

	if *arp {
		runARPScan(graph)
	}

	if *icmp {
		runICMPSweep()
	}

	if !*arp && !*icmp {
		fmt.Println("Please specify a scan type. Use -h for help.")
	}

	fmt.Println(graph)
}

func runARPScan(graph *_map.Graph) {
	ifaces, err := net.Interfaces()
	if err != nil {
		fmt.Println("Error getting network interfaces:", err)
		return
	}

	in := make(chan arp_scan.ARPEvent)

	go func() {
		for ev := range in {
			fmt.Printf("Discovered device - IP: %s, MAC: %s, Network: %s, Source: %s\n", ev.IP, ev.MAC, ev.Network, ev.Source)
			node := graph.GetOrCreateNode("ip:" + ev.IP.String())
			node.Protocols["arp"] = true
			graph.AddEdge(node.ID, graph.GetOrCreateNode("net:"+ev.Network.String()).ID, _map.EdgeMemberOf)
		}
	}()

	for _, iface := range ifaces {
		fmt.Printf("Starting ARP Scan on interface %s\n", iface.Name)
		err := arp_scan.Scan(iface, in)
		if err != nil {
			fmt.Printf("Error scanning network on interface %s: %v\n", iface.Name, err)
			continue
		}
	}

}

func runICMPSweep() {
	ifaces, err := net.Interfaces()
	if err != nil {
		fmt.Println("Error getting network interfaces:", err)
		return
	}

	for _, iface := range ifaces {
		fmt.Printf("Starting ICMP Sweep on interface %s\n", iface.Name)
		if err := ping.Sweep(iface); err != nil {
			fmt.Printf("Error during ICMP Sweep on interface %s: %v\n", iface.Name, err)
		}
	}
}
