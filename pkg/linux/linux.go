package linux

import (
	"io/ioutil"
	"os/exec"

	"github.com/sirupsen/logrus"
)

// EnableIPForwarding enables IPv4 forwarding
func EnableIPForwarding(log *logrus.Entry) error {
	log.Print("enabling ip forwarding")

	return ioutil.WriteFile("/proc/sys/net/ipv4/ip_forward", []byte("1"), 0666)
}

// DisableIPTables disables the iptables rule used by proxy
func DisableIPTables(log *logrus.Entry) error {
	log.Print("disabling iptables rule")

	return iptables("-D")
}

// EnableIPTables enable the iptables rule used by proxy
func EnableIPTables(log *logrus.Entry) error {
	log.Print("enabling iptables rule")

	iptables("-D") // clean up if necessary, ignore failure

	return iptables("-I")
}

func iptables(verb string) error {
	return exec.Command("iptables", "-t", "nat", verb, "PREROUTING", "-p",
		"tcp", "!", "-d", interfaceIP.String(), "-j", "REDIRECT", "--to-port",
		"3128").Run()
}
