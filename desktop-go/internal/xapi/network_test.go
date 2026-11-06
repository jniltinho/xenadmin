package xapi

import (
	"net/http"
	"testing"
	"time"
)

func TestProxy(t *testing.T) {
	defer SetNetwork(Network{Proxy: "direct", Timeout: 20 * time.Second})
	req, _ := http.NewRequest(http.MethodGet, "https://10.0.0.1/jsonrpc", nil)
	SetNetwork(Network{Proxy: "specified", Address: "proxy.lan", Port: 3128})
	if u, _ := Proxy(req); u == nil || u.Host != "proxy.lan:3128" {
		t.Fatalf("specified: %v", u)
	}
	SetNetwork(Network{Proxy: "specified", Address: "proxy.lan", Port: 3128, Bypass: true})
	if u, _ := Proxy(req); u != nil {
		t.Fatal("bypass")
	}
	t.Setenv("HTTPS_PROXY", "http://env.lan:8080")
	SetNetwork(Network{Proxy: "system", Bypass: true}) // the bypass is the specified proxy's only
	if u, _ := Proxy(req); u == nil || u.Host != "env.lan:8080" {
		t.Fatalf("system: %v", u)
	}
	SetNetwork(Network{Proxy: "direct"})
	if u, _ := Proxy(req); u != nil {
		t.Fatal("direct")
	}
	SetNetwork(Network{Proxy: "direct", Timeout: 2 * time.Minute})
	if RequestTimeout() != 2*time.Minute || HandshakeTimeout() != 2*time.Minute {
		t.Fatal("timeout")
	}
}
