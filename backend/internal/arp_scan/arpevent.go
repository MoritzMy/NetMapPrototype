package arp_scan

import "net"

type ARPEvent struct {
	IP      net.IP
	MAC     net.HardwareAddr
	Network *net.IPNet
	Source  string
}
