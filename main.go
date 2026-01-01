package main

import (
	"fmt"
	"net"

	"github.com/MoritzMy/NetMap/scan/arp_scan"
	"github.com/MoritzMy/NetMap/scan/ping"
)

func main() {
	ifaces, _ := net.Interfaces()
	fmt.Println(ifaces)
	arp_scan.ScanNetwork(ifaces[1])

	ping.Sweep(ifaces)

	return
}
