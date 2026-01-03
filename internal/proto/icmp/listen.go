package icmp

import (
	"context"
	"net"

	"github.com/MoritzMy/NetMap/internal/proto"
)

const (
	echoReplyType = 0
)

func PingReplyListener(conn net.PacketConn, ctx context.Context) <-chan net.IP {
	ch := make(chan net.IP)
	buf := make([]byte, 200)

	go func() {
		defer close(ch)
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}

			_, addr, err := conn.ReadFrom(buf)
			if err != nil {
				continue
			}

			var packet EchoICMPPacket

			if err := proto.Unmarshal(buf, &packet); err != nil {
				continue
			}

			if packet.Type != echoReplyType {
				continue
			}

			ch <- addr.(*net.IPAddr).IP
		}
	}()
	return ch
}
