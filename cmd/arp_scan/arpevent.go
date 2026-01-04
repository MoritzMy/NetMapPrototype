package arp_scan

import "net"

type ARPEvent struct {
	IP      net.IP
	MAC     net.HardwareAddr
	Network *net.IPNet
	Source  string
}

func CanonicalIPNet(ipNet *net.IPNet) *net.IPNet {
	ip := ipNet.IP.Mask(ipNet.Mask)
	return &net.IPNet{
		IP:   ip,
		Mask: ipNet.Mask,
	}
}
