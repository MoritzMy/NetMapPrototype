package ping

import "net"

type PingEvent struct {
	IP      net.IP
	Network *net.IPNet
}
