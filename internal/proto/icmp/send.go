package icmp

import (
	"net"

	"github.com/MoritzMy/NetMap/internal/proto"
)

func SendPing(conn net.PacketConn, dst net.IP, id, seq uint16) error {
	req := NewEchoICMPPacket(id, seq, []byte("ARE U UP?"))
	b, err := proto.Marshal(&req)
	if err != nil {
		return err
	}

	_, err = conn.WriteTo(b, &net.IPAddr{IP: dst})
	return err
}
