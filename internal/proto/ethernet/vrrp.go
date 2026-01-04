package eth

import (
	"fmt"
	"net"
)

func IsVRRPMulticastMAC(mac net.HardwareAddr) bool {
	fmt.Printf("MAC=%s len=%d bytes=% x\n", mac, len(mac), mac)

	if len(mac) != 6 {
		fmt.Printf("MAC=%s len=%d\n", mac, len(mac))
		return false
	}

	return mac[0] == 0x00 && mac[1] == 0x00 && mac[2] == 0x5e && mac[3] == 0x00 && mac[4] == 0x01 && mac[5] != 0x00
}
