package proxy

import (
	"crypto/tls"
	"net"
)

func init() {
	handlers["mitm"] = mitm
}

func mitm(p *Proxy, c1 tcpConn, raddr *net.TCPAddr) error {
	var chi tls.ClientHelloInfo

	tc1 := tls.Server(c1, &tls.Config{
		GetCertificate: func(_chi *tls.ClientHelloInfo) (*tls.Certificate, error) {
			chi = *_chi

			key, cert, err := p.db.GetKeyAndCertificate(chi.ServerName)
			if err != nil {
				return nil, err
			}

			return &tls.Certificate{
				Certificate: [][]byte{
					cert.Raw,
					p.db.CACert.Raw,
				},
				PrivateKey: key,
			}, nil
		},
	})

	err := tc1.Handshake()
	if err != nil {
		return err
	}

	c2, err := net.DialTCP("tcp", nil, raddr)
	if err != nil {
		return err
	}
	defer c2.Close()

	tc2 := tls.Client(c2, &tls.Config{
		ServerName: chi.ServerName,
	})

	err = tc2.Handshake()
	if err != nil {
		return err
	}

	return p.full(tc1, tc2)
}
