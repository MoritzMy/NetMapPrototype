package arp

import "net"

type Reply struct {
	IP  net.IP
	MAC net.HardwareAddr
}
