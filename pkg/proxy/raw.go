package proxy

import (
	"io"
	"net"
)

func init() {
	handlers["raw"] = raw
}

func raw(p *Proxy, c1 tcpConn, raddr *net.TCPAddr) error {
	c2, err := net.DialTCP("tcp", nil, raddr)
	if err != nil {
		return err
	}
	defer c2.Close()

	return p.full(c1, c2)
}

func (p *Proxy) full(c1, c2 tcpConn) error {
	ch := make(chan error, 2)

	go func() {
		ch <- p.half(c2, c1)
	}()

	go func() {
		ch <- p.half(c1, c2)
	}()

	err := <-ch

	if err != nil {
		return err
	}

	return <-ch
}

func (p *Proxy) half(dst, src tcpConn) error {
	_, err := io.Copy(dst, src)
	if err != nil {
		return err
	}

	return dst.CloseWrite()
}
