package ping

import (
	"net"
	"time"

	"github.com/MoritzMy/NetMap/proto"
	"github.com/MoritzMy/NetMap/proto/icmp"
	"github.com/MoritzMy/NetMap/proto/ip"
)

const (
	echoReplyType = 0
)

func Ping(addr net.IP) (*ip.IPv4Packet, error) {
	var identifier uint16 = 0
	var sequenceNumber uint16 = 0

	// Create ICMP Echo Request packet

	req := icmp.NewEchoICMPPacket(identifier, sequenceNumber, []byte("ARE U UP?"))
	b, err := proto.Marshal(&req)
	if err != nil {
		panic(err)
	}
	conn, err := net.Dial("ip4:icmp", addr.String())

	// Ensure the connection is closed after we're done

	if err != nil {
		panic(err)
	}
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			panic(err)
		}
	}(conn)

	// Set write deadline and send the packet

	err = conn.SetWriteDeadline(time.Now().Add(2 * time.Second))
	if err != nil {
		return nil, err
	}
	_, err = conn.Write(b)
	if err != nil {
		return nil, err
	}

	cr := make([]byte, 200)

	// Set read deadline and read the response

	err = conn.SetReadDeadline(time.Now().Add(time.Second))
	if err != nil {
		return nil, err
	}
	_, err = conn.Read(cr)
	if err != nil {
		return nil, err
	}

	// Unmarshal the response into an IPv4Packet

	var packet ip.IPv4Packet

	if err := proto.Unmarshal(cr, &packet); err != nil {
		return nil, err
	}

	return &packet, nil

}
