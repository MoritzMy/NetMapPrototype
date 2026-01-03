package arp

import (
	"context"
	"encoding/binary"
	"log"
	"net"
	"syscall"
)

// ARPResponseListener listens for ARP replies on the given file descriptor and sends them to the returned channel.
func ARPResponseListener(fd int, ctx context.Context) <-chan Reply {
	out := make(chan Reply, 256)

	go func() {
		defer close(out)

		buf := make([]byte, 128) // Buffer to hold incoming packets

		for {
			select {
			case <-ctx.Done():
				return
			default:
			}
			n, _, err := syscall.Recvfrom(fd, buf, 0)

			if err != nil {
				log.Printf("arp recv error: %v", err)
				continue
			}

			// Basic validation of ARP reply packet
			if n < 42 || !(buf[12] == 0x08 && buf[13] == 0x06 && binary.BigEndian.Uint16(buf[20:22]) == 2) {
				continue
			}
			select {
			case out <- Reply{
				IP:  net.IPv4(buf[28], buf[29], buf[30], buf[31]),
				MAC: net.HardwareAddr{buf[22], buf[23], buf[24], buf[25], buf[26], buf[27]},
			}:
			default:

			}
		}
	}()
	return out
}
