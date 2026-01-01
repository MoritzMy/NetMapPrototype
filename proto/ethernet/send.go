package eth

import (
	"net"
	"syscall"
)

func SendEthernetFrame(frame []byte, iface string, fd int) error {
	_, err := syscall.Write(fd, frame)
	return err
}

func CreateSocket(interf *net.Interface) (int, error) {
	ifIndex := interf.Index

	fd, err := syscall.Socket(
		syscall.AF_PACKET,
		syscall.SOCK_RAW,
		int(htons(syscall.ETH_P_ARP)))
	if err != nil {
		return 0, err
	}

	addr := syscall.SockaddrLinklayer{
		Protocol: htons(syscall.ETH_P_ARP),
		Ifindex:  ifIndex,
	}

	if err := syscall.Bind(fd, &addr); err != nil {
		syscall.Close(fd)
		return 0, err
	}

	return fd, nil
}

func htons(v uint16) uint16 {
	return (v << 8) | (v >> 8)
}
