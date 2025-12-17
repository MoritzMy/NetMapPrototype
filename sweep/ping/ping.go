package ping

import (
	"fmt"
	"net"
	"time"

	"github.com/MoritzMy/NetMap/sweep/ping/icmp"
)

const (
	echoReplyType = 0
)

func Ping(addr net.IP) *IPPacket {
	var identifier uint16 = 0
	var sequenceNumber uint16 = 0

	req := icmp.NewEchoICMPPacket(identifier, sequenceNumber, []byte("ARE U UP?"))
	b, err := icmp.Marshal(&req)
	if err != nil {
		panic(err)
	}
	conn, err := net.Dial("ip4:icmp", addr.String())

	if err != nil {
		panic(err)
	}
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			panic(err)
		}
	}(conn)

	err = conn.SetWriteDeadline(time.Now().Add(2 * time.Second))
	if err != nil {
		return nil
	}
	write, err := conn.Write(b)
	if err != nil {
		return nil
	}

	fmt.Println(fmt.Sprintf("Wrote %d Bits to: %s", write, addr.String()))

	cr := make([]byte, 84)

	err = conn.SetReadDeadline(time.Now().Add(time.Second))
	if err != nil {
		return nil
	}
	read, err := conn.Read(cr)
	if err != nil {
		fmt.Println(fmt.Sprintf("%v is not responding", addr))
		return nil
	}
	fmt.Println(fmt.Sprintf("Recieved %d Bits from: %s", read, addr.String()))

	packet := NewIPPacket(cr)

	return packet

}
