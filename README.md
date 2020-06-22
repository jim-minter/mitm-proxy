# mitm-proxy

mitm-proxy proxies all incoming TCP connections.  Currently it drops all non-TLS
connections.  The remaining connections are handled based on user-configurable
rules.

* Run as root.  The proxy automatically enables IPv4 forwarding and its iptables
  rule.  SIGINT (^C) terminates the proxy and removes the iptables rule.

* Sample configuration in config.yml.  Supported handlers are `drop`, `mitm` and
  `raw`.  Regular expressions are automatically wrapped with `^` and `$`, i.e.
  they must match the entire hostname.

## Handlers

* drop: drops the connection.

* mitm: runs the connection through a TLS "man-in-the-middle" proxy.

* raw: runs the connection through a raw TCP proxy.
