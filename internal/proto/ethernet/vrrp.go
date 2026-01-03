package eth

import "net"

func IsVRRPMulticastMAC(mac net.HardwareAddr) bool {
	return len(mac) == 6 && mac[0] == 0x00 && mac[1] == 0x00 && mac[2] == 0x5e && mac[3] == 0x00 && mac[4] == 0x01 && mac[5] != 0
}
