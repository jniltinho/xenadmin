package rules

import (
	"net"
	"strconv"
	"strings"
)

// The address search (XenCenterLib/ComparableAddress.cs, PartialIP.cs; XenSearch IPAddressQuery).

// Address is a ComparableAddress: an IP address, a partial IPv4 address ("10.77.*", "10.77"), or a name.
type Address struct {
	ip      net.IP
	partial []int // four segments, -1 for any
	name    string
}

// ParseAddress is ComparableAddress.TryParse.
func ParseAddress(s string, allowPartial, allowName bool) (Address, bool) {
	s = strings.TrimSpace(s)
	if s == "" {
		return Address{}, false
	}
	if ip := net.ParseIP(s); ip != nil {
		return Address{ip: ip}, true
	}
	if ip := shortIPv4(s); ip != nil { // .NET IPAddress.TryParse also reads "10.77" (10.0.0.77), before PartialIP
		return Address{ip: ip}, true
	}
	if allowPartial && strings.Contains(s, ".") { // PartialIP.TryParse
		segs := strings.Split(s, ".")
		p := []int{-1, -1, -1, -1}
		ok := len(segs) <= 4
		for i := 0; ok && i < len(segs); i++ {
			if n, err := strconv.Atoi(segs[i]); err == nil && n >= 0 && n <= 32767 {
				p[i] = n
			} else if !(segs[i] == "" || segs[i] == "*" || segs[i] == "_" || segs[i] == "x" || segs[i] == "n") {
				ok = false
			}
		}
		if ok {
			return Address{partial: p}, true
		}
	}
	if allowName {
		return Address{name: s}, true
	}
	return Address{}, false
}

// Matches is ComparableAddress.Equals: a partial address matches the IPv4 addresses in it, two addresses or two names
// are equal when they are the same.
func (a Address) Matches(b Address) bool {
	switch {
	case a.partial != nil && b.ip != nil:
		return partialMatch(a.partial, b.ip)
	case b.partial != nil && a.ip != nil:
		return partialMatch(b.partial, a.ip)
	case a.ip != nil && b.ip != nil:
		return a.ip.Equal(b.ip)
	case a.name != "" && b.name != "":
		return a.name == b.name
	}
	return false
}

func partialMatch(p []int, ip net.IP) bool {
	v4 := ip.To4()
	if v4 == nil {
		return false
	}
	for i, seg := range p {
		if seg > -1 && seg != int(v4[i]) {
			return false
		}
	}
	return true
}

// String is ComparableAddress.ToString: the address, a partial one with "*" for any segment, or the name.
func (a Address) String() string {
	switch {
	case a.ip != nil:
		return a.ip.String()
	case a.partial != nil:
		segs := make([]string, 4)
		for i, n := range a.partial {
			if segs[i] = "*"; n > -1 {
				segs[i] = strconv.Itoa(n)
			}
		}
		return strings.Join(segs, ".")
	}
	return a.name
}

// shortIPv4 is the inet_aton forms of IPv4 that .NET IPAddress.TryParse takes: one to three decimal numbers, the last
// one filling the remaining bytes ("10.77" is 10.0.0.77, "1" is 0.0.0.1). shortcut: no octal or hex parts.
func shortIPv4(s string) net.IP {
	parts := strings.Split(s, ".")
	if len(parts) > 3 {
		return nil
	}
	var b [4]byte
	for i, p := range parts {
		if p == "" || strings.Trim(p, "0123456789") != "" {
			return nil
		}
		n, err := strconv.ParseUint(p, 10, 32)
		if err != nil {
			return nil
		}
		if i < len(parts)-1 {
			if n > 255 {
				return nil
			}
			b[i] = byte(n)
			continue
		}
		rest := 4 - i // bytes the last number fills
		if n >= 1<<(8*rest) {
			return nil
		}
		for j := 3; j >= i; j-- {
			b[j], n = byte(n), n>>8
		}
	}
	return net.IPv4(b[0], b[1], b[2], b[3])
}
