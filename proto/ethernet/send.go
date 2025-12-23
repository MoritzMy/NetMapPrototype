package eth

import "syscall"

func SendEthernetFrame(b []byte) {
	fd, err := syscall.Socket(
		syscall.AF_PACKET,
		syscall.SOCK_RAW,
		int(htons(syscall.ETH_P_ARP)))
}

func htons(v uint16) uint16 {
	return (v << 8) | (v >> 8)
}
