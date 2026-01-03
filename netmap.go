package main

import (
	"flag"
	"fmt"
	"net"

	"github.com/MoritzMy/NetMap/cmd/scan/arp_scan"
	"github.com/MoritzMy/NetMap/cmd/scan/ping"
)

func main() {
	arp := flag.Bool("arp-scan", false, "Run ARP Discovery Scan")
	icmp := flag.Bool("ping-sweep", false, "Run ICMP Sweep")

	flag.Parse()

	if *arp {
		runARPScan()
	}

	if *icmp {
		runICMPSweep()
	}

	if !*arp && !*icmp {
		fmt.Println("Please specify a scan type. Use -h for help.")
	}
}

func runARPScan() {
	ifaces, err := net.Interfaces()
	if err != nil {
		fmt.Println("Error getting network interfaces:", err)
		return
	}

	for _, iface := range ifaces {
		fmt.Printf("Starting ARP Scan on interface %s\n", iface.Name)
		if err := arp_scan.Scan(iface); err != nil {
			fmt.Printf("Error scanning network on interface %s: %v\n", iface.Name, err)
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
