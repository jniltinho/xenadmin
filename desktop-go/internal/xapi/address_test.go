package xapi

import "testing"

func TestParseAddress(t *testing.T) {
	cases := []struct {
		in   string
		host string
		port int
		err  bool
	}{
		{"xcp1", "xcp1", 443, false},
		{" 10.0.0.1:8443 ", "10.0.0.1", 8443, false},
		{"2001:db8::11", "2001:db8::11", 443, false},
		{"[2001:db8::11]:4443", "2001:db8::11", 4443, false},
		{"[2001:db8::11]", "2001:db8::11", 443, false},
		{"xcp1:80", "", 0, true},
		{"xcp1:8080", "", 0, true},
		{"xcp1:0", "", 0, true},
		{"xcp1:abc", "", 0, true},
		{"", "", 0, true},
		{"[::1", "", 0, true},
	}
	for _, c := range cases {
		h, p, err := ParseAddress(c.in)
		if (err != nil) != c.err || h != c.host || p != c.port {
			t.Errorf("ParseAddress(%q) = %q, %d, %v", c.in, h, p, err)
		}
	}
	if HostKey("2001:db8::11", 443) != "[2001:db8::11]:443" {
		t.Error("HostKey brackets IPv6")
	}
}
