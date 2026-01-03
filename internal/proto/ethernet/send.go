package eth

import (
	"net"
	"syscall"
)

// SendEthernetFrame sends a raw Ethernet frame through the specified network interface using the provided file descriptor.
func SendEthernetFrame(frame []byte, iface string, fd int) error {
	_, err := syscall.Write(fd, frame)
	return err
}

// CreateSocket creates a raw socket bound to the specified network interface for sending and receiving Ethernet frames.
func CreateSocket(interf *net.Interface) (int, error) {
	ifIndex := interf.Index

	// Create a raw socket
	fd, err := syscall.Socket(
		syscall.AF_PACKET,
		syscall.SOCK_RAW,
		int(htons(syscall.ETH_P_ARP)))
	if err != nil {
		return 0, err
	}

	// Bind the socket to the specified interface using SockaddrLinklayer
	addr := syscall.SockaddrLinklayer{
		Protocol: htons(syscall.ETH_P_ARP),
		Ifindex:  ifIndex,
	}

	// Bind the socket to the address structure
	if err := syscall.Bind(fd, &addr); err != nil {
		syscall.Close(fd)
		return 0, err
	}

	return fd, nil
}

// htons converts a 16-bit integer from host byte order to network byte order.
func htons(v uint16) uint16 {
	return (v << 8) | (v >> 8)
}
