package proxy

import (
	"net"
)

func init() {
	handlers["drop"] = drop
}

func drop(p *Proxy, c1 tcpConn, raddr *net.TCPAddr) error {
	return nil
}
