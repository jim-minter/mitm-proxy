package signatures

import (
	"crypto/tls"
	"errors"

	"github.com/jim-minter/mitm-proxy/pkg/conn"
)

// TLSServerName attempts to read a ClientHello message on the connection c and
// returns the server name requested via SNI if found.
func TLSServerName(c *conn.Recorder) string {
	var chi tls.ClientHelloInfo

	tlsc := tls.Server(&conn.WriteDiscarder{Conn: c}, &tls.Config{
		GetCertificate: func(_chi *tls.ClientHelloInfo) (*tls.Certificate, error) {
			chi = *_chi
			return nil, errors.New("")
		},
	})

	tlsc.Handshake()

	return chi.ServerName
}
