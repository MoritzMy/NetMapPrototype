package arp_scan

import (
	"net"

	"github.com/MoritzMy/NetMap/proto/arp"
	eth "github.com/MoritzMy/NetMap/proto/ethernet"
	"github.com/MoritzMy/NetMap/scan"
)

func ARPScan(adress scan.InterfaceAdress, targetIP net.IP) {
	req := arp.NewARPRequest(adress.MAC, adress.IPs[0], net.IP{0x0a, 0xfe, 0xf0, 0x13})
	b, err := arp.Marshal(req)
	if err != nil {
		panic(err)
	}
	eth.SendEthernetFrame(b)

}

func SendEthernetFrame() {

}
