package rules

import "testing"

func TestAddressMatches(t *testing.T) {
	for in, want := range map[string]string{"10.77": "10.0.0.77", "10.77.*": "10.77.*.*", "1": "0.0.0.1", "10.1.258": "10.1.1.2"} {
		if a, _ := ParseAddress(in, true, false); a.String() != want {
			t.Errorf("%q: %s", in, a.String())
		}
	}
	addr := func(s string, partial bool) Address {
		a, ok := ParseAddress(s, partial, true)
		if !ok {
			t.Fatalf("%q does not parse", s)
		}
		return a
	}
	for _, c := range []struct {
		query, value string
		want         bool
	}{
		{"10.77.58.11", "10.77.58.11", true},
		{"10.77.*", "10.77.58.11", true},
		{"10.77.", "10.77.58.11", true},
		{"10.78.*", "10.77.58.11", false},
		{"10.77.x.12", "10.77.58.12", true},
		{"fe80::1", "fe80::1", true},
		{"10.*", "fe80::1", false},
		{"server1", "server1", true},
		{"server1", "10.0.0.1", false},
	} {
		if got := addr(c.query, true).Matches(addr(c.value, false)); got != c.want {
			t.Errorf("%q vs %q = %v", c.query, c.value, got)
		}
	}
	if _, ok := ParseAddress("10.a.1", true, false); ok {
		t.Error("a bad partial address parsed")
	}
}
