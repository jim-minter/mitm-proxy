package conn

import (
	"net"
)

// WriteDiscarder is a connection wrapper that discards all writes.
type WriteDiscarder struct {
	net.Conn
}

func (*WriteDiscarder) Write(b []byte) (int, error) {
	return len(b), nil
}
