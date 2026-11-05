// Package xapi is the connection core: server addresses, the HTTP client with the app's TLS trust, login and the
// pool coordinator redirect. It has no UI, CLI or MCP code.
package xapi

import (
	"errors"
	"net"
	"strconv"
	"strings"
)

// DefaultPort is the XAPI HTTPS port.
const DefaultPort = 443

// ErrPlainHTTP refuses ports 80 and 8080: XAPI answers plain HTTP there and the password would cross the network
// unencrypted (ServerAddress.IsPlainHttp in the C# app).
var ErrPlainHTTP = errors.New("ports 80 and 8080 are plain HTTP; connect over HTTPS (443)")

// ParseAddress splits the "Server:" field: "host[:port]", "[v6]:port" or a bare IPv6 address (port of
// XenAdmin.Desktop.Core/ServerAddress.cs).
func ParseAddress(text string) (host string, port int, err error) {
	host = strings.TrimSpace(text)
	port = DefaultPort
	bad := errors.New("enter a valid server name or address")
	switch {
	case strings.HasPrefix(host, "["):
		end := strings.Index(host, "]")
		if end < 0 {
			return "", 0, bad
		}
		rest := host[end+1:]
		host = host[1:end]
		if rest != "" {
			if rest[0] != ':' {
				return "", 0, bad
			}
			if port, err = strconv.Atoi(rest[1:]); err != nil {
				return "", 0, bad
			}
		}
	case strings.Count(host, ":") > 1 && net.ParseIP(host) != nil:
		// a bare IPv6 address has no port: "2001:db8::11" is not host "2001:db8:" port 11
	default:
		if i := strings.LastIndex(host, ":"); i >= 0 {
			if port, err = strconv.Atoi(host[i+1:]); err != nil {
				return "", 0, bad
			}
			host = host[:i]
		}
	}
	if host == "" || port < 1 || port > 65535 {
		return "", 0, bad
	}
	if port == 80 || port == 8080 {
		return "", 0, ErrPlainHTTP
	}
	return host, port, nil
}

// HostKey is how certificates are pinned: "host:port".
func HostKey(host string, port int) string { return net.JoinHostPort(host, strconv.Itoa(port)) }
