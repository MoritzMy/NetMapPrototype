package ping

import (
	"fmt"
	"net"
)

func ValidIpsInNetwork(addr *net.IPNet) []net.IP {
	ip := addr.IP
	subnet, size := addr.Mask.Size()

	bytes := []byte(ip)

	if ip.To4() == nil {
		return nil
	}

	ipv4Bytes := bytes[len(bytes)-4:]

	fmt.Println(ipv4Bytes)

	fmt.Println(ip, subnet, size)

	for i := 0; i < size-subnet; i++ {

	}

	return nil
}
