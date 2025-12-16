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

func Ping(addr *net.IPNet) *icmp.EchoICMPPacket {
	var identifier uint16 = 0
	var sequenceNumber uint16 = 0

	req := icmp.NewEchoICMPPacket(identifier, sequenceNumber, []byte("ARE U UP?"))
	b, err := icmp.Marshal(&req)
	if err != nil {
		panic(err)
	}
	conn, err := net.Dial("ip4:icmp", addr.IP.String())

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

	fmt.Println(fmt.Sprintf("Wrote %d Bits to: %s", write, addr.IP.String()))

	cr := make([]byte, 84)

	err = conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	if err != nil {
		return nil
	}
	read, err := conn.Read(cr)
	if err != nil {
		panic(err)
	}

	fmt.Println(fmt.Sprintf("Recieved %d Bits from: %s", read, addr.IP.String()))

	var u icmp.EchoICMPPacket

	ipHeaderSize := NewIpv4VersionIHL(cr[0]).Size()

	// TODO, replace MagicNumberTM with the IHL from the IPHeader
	icmp.Unmarshal(cr[ipHeaderSize:], &u)
	fmt.Println(u)

	return &u

}
