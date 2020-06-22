package linux

import (
	"net"
	"syscall"
	"unsafe"
)

func init() {
	var err error
	interfaceIP, err = getInterfaceIP("eth0")
	if err != nil {
		panic(err)
	}
}

var interfaceIP net.IP

type ifreq struct {
	ifrName [16]byte
	ifrAddr syscall.RawSockaddrInet4
	pad     [8]byte
}

func getInterfaceIP(iface string) (net.IP, error) {
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
	if err != nil {
		return nil, err
	}
	defer syscall.Close(fd)

	var ifreq ifreq
	copy(ifreq.ifrName[:], []byte(iface))

	_, _, errno := syscall.Syscall6(syscall.SYS_IOCTL, uintptr(fd),
		syscall.SIOCGIFADDR, uintptr(unsafe.Pointer(&ifreq)), 0, 0, 0)
	if errno != 0 {
		return nil, errno
	}

	return net.IP(ifreq.ifrAddr.Addr[:]), nil
}
