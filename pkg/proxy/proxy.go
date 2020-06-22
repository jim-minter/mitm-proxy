package proxy

import (
	"fmt"
	"net"
	"regexp"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/jim-minter/mitm-proxy/pkg/config"
	"github.com/jim-minter/mitm-proxy/pkg/conn"
	"github.com/jim-minter/mitm-proxy/pkg/linux"
	"github.com/jim-minter/mitm-proxy/pkg/signatures"
	"github.com/jim-minter/mitm-proxy/pkg/tls"
)

// Proxy describes a proxy
type Proxy struct {
	log *logrus.Entry
	db  *tls.DB
	l   *net.TCPListener

	rules []rule
}

type tcpConn interface {
	net.Conn
	CloseWrite() error
}

type rule struct {
	rx      *regexp.Regexp
	handler string
}

var handlers = map[string]func(p *Proxy, c1 tcpConn, raddr *net.TCPAddr) error{}

// NewProxy instantiates a new proxy
func NewProxy(log *logrus.Entry) (*Proxy, error) {
	var err error

	p := &Proxy{
		log: log,
	}

	p.db, err = tls.NewDB()
	if err != nil {
		return nil, err
	}

	for i, h := range config.Config.Handlers {
		rx, err := regexp.Compile("^" + h.Regexp + "$")
		if err != nil {
			return nil, fmt.Errorf("handlers[%d].regexp: %s", i, err)
		}

		if _, ok := handlers[strings.ToLower(h.Handler)]; !ok {
			return nil, fmt.Errorf("handlers[%d].handler: %q is unrecognised", i, h.Handler)
		}

		p.rules = append(p.rules, rule{
			rx:      rx,
			handler: strings.ToLower(h.Handler),
		})
	}

	p.rules = append(p.rules, rule{
		rx:      regexp.MustCompile(""),
		handler: "drop",
	})

	p.l, err = net.ListenTCP("tcp", &net.TCPAddr{Port: 3128})
	if err != nil {
		return nil, err
	}

	log.Print("listening")

	return p, nil
}

// Run runs the proxy
func (p *Proxy) Run() error {
	for {
		c, err := p.l.AcceptTCP()
		if err != nil {
			return err
		}

		go p.handle(c)
	}
}

func (p *Proxy) handle(c *net.TCPConn) error {
	defer c.Close()

	raddr, err := linux.OriginalDest(c)
	if err != nil {
		return err
	}

	rc := &conn.Recorder{TCPConn: c}

	serverName := signatures.TLSServerName(rc)

	rc.Rewind()
	rc.StopRecording()

	for _, r := range p.rules {
		if r.rx.MatchString(serverName) {
			p.log.Printf("%s->%s (%s): %s", c.RemoteAddr(), raddr, serverName, r.handler)
			return handlers[r.handler](p, rc, raddr)
		}
	}

	return nil
}
