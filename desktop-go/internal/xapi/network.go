package xapi

import (
	"net"
	"net/http"
	"net/url"
	"strconv"
	"sync/atomic"
	"time"
)

// Network is Tools > Options' proxy and connection timeout (OptionsDialog, HTTPHelper.ProxyStyle) for the
// connections to the servers.
type Network struct {
	Proxy   string // "direct" (the default), "system" (the operating system's settings) or "specified"
	Address string // the specified proxy server
	Port    int
	Bypass  bool          // "Bypass proxy server for all server connections": every connection of this app is to a server
	Timeout time.Duration // how long to wait for a server (dial and TLS handshake)
}

var network atomic.Pointer[Network]

// SetNetwork applies the options to the connections made from now on.
func SetNetwork(n Network) { network.Store(&n) }

func current() Network {
	if n := network.Load(); n != nil {
		return *n
	}
	return Network{Proxy: "direct", Timeout: 20 * time.Second}
}

// Proxy is the proxy of a request to a server (nil for none).
func Proxy(req *http.Request) (*url.URL, error) {
	n := current()
	switch n.Proxy {
	case "specified":
		if n.Bypass || n.Address == "" { // the bypass belongs to the specified proxy only (OptionsDialog)
			return nil, nil
		}
		return &url.URL{Scheme: "http", Host: net.JoinHostPort(n.Address, strconv.Itoa(n.Port))}, nil
	case "system":
		// shortcut: the system's proxy is read from HTTP_PROXY/HTTPS_PROXY/NO_PROXY (the Linux convention), not from the
		// Windows or macOS proxy settings; read those natively if users need it
		return http.ProxyFromEnvironment(req)
	}
	return nil, nil
}

// Dialer is the dialer of a connection to a server, with the options' timeout.
func Dialer() *net.Dialer {
	t := current().Timeout
	if t <= 0 {
		t = 20 * time.Second
	}
	return &net.Dialer{Timeout: t, KeepAlive: 30 * time.Second}
}

// HandshakeTimeout is the TLS handshake's limit (the options' timeout).
func HandshakeTimeout() time.Duration {
	if t := current().Timeout; t > 0 {
		return t
	}
	return 20 * time.Second
}

// RequestTimeout is the limit of one API call: the options' timeout, never under a minute (long calls such as
// event.from wait up to 30 s on the server).
func RequestTimeout() time.Duration { return max(time.Minute, HandshakeTimeout()) }
