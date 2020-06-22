package linux

import (
	"net"
	"syscall"
	"unsafe"
)

// OriginalDest reads the SO_ORIGINAL_DST socket option to find the original
// destination of an incoming TCP connection redirected using the netfilter
// REDIRECT extension.
func OriginalDest(c *net.TCPConn) (*net.TCPAddr, error) {
	const SO_ORIGINAL_DST = 80

	f, err := c.File()
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var sockaddr syscall.RawSockaddrInet4
	socklen := uint32(syscall.SizeofSockaddrInet4)

	_, _, errno := syscall.Syscall6(syscall.SYS_GETSOCKOPT, f.Fd(),
		syscall.SOL_IP, SO_ORIGINAL_DST, uintptr(unsafe.Pointer(&sockaddr)),
		uintptr(unsafe.Pointer(&socklen)), 0)
	if errno != 0 {
		return nil, errno
	}

	return &net.TCPAddr{
		IP:   sockaddr.Addr[:],
		Port: int(ntohs(sockaddr.Port)),
	}, nil
}

// ntohs converts netshort from network byte order to host byte order
func ntohs(netshort uint16) uint16 {
	p := (*[2]byte)(unsafe.Pointer(&netshort))
	return uint16(p[0])<<8 + uint16(p[1])
}
