package ip

import (
	"encoding/binary"
	"net"
)

func ValidIpsInNetwork(addr *net.IPNet) []net.IP {
	var hosts []net.IP

	baseAddr := addr.IP.Mask(addr.Mask)

	ip := addr.IP
	subnet, size := addr.Mask.Size()

	ip4 := ip.To4()

	if ip4 == nil {
		return nil
	}

	currAddr := append(net.IP(nil), baseAddr...)
	hostBits := size - subnet
	count := (1 << hostBits) - 2

	for i := 0; i < count; i++ {
		incrementIP(currAddr)

		if isNetworkIP(currAddr, subnet) || isBroadcastIP(currAddr, subnet) {
			continue
		}

		hosts = append(hosts, append(net.IP(nil), currAddr...))
	}

	return hosts
}

func incrementIP(ip net.IP) net.IP {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] != 0 {
			break
		}
	}
	return ip
}

func isNetworkIP(ip net.IP, prefixLen int) bool {
	hostBits := 32 - prefixLen

	ip4 := ip.To4()
	if ip4 == nil {
		return false
	}
	ipNumeric := binary.BigEndian.Uint32(ip4)

	if ipNumeric == ipNumeric>>hostBits<<hostBits {
		return true
	}

	return false
}

func isBroadcastIP(ip net.IP, prefixLen int) bool {
	ip4 := ip.To4()
	if ip4 == nil {
		return false
	}
	ipNumeric := binary.BigEndian.Uint32(ip4)

	hostBits := 32 - prefixLen

	hostMask := uint32((1 << hostBits) - 1)

	return ipNumeric&hostMask == hostMask
}
