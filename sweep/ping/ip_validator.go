package ping

import (
	"fmt"
	"net"
)

func ValidateIP(addr *net.IPNet) bool {
	ip := addr.IP
	subnet, size := addr.Mask.Size()

	bytes := []byte(ip)

	if ip.To4() == nil {
		return false
	}

	fmt.Println(bytes)

	fmt.Println(ip, subnet, size)

	return true
}
