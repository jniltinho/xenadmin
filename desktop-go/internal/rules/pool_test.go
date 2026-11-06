package rules

import (
	"strings"
	"testing"
)

func TestCanJoinPool(t *testing.T) {
	c := PoolFacts{Key: "a", Hosts: 1, CPUVendor: "AuthenticAMD", APIMajor: 2, APIMinor: 21, DBSchema: "5.0", PlatformVersion: "3.4.0"}
	m := c
	m.Key = "b"
	if CanJoinPool(c, c) != JoinWillBeCoordinator || CanJoinPool(m, c) != JoinAllowed {
		t.Fatal("basic")
	}
	for want, change := range map[JoinReason]func(*PoolFacts){
		JoinIsAPool:                func(p *PoolFacts) { p.PoolName = "x" },
		JoinDifferentCPUs:          func(p *PoolFacts) { p.CPUVendor = "GenuineIntel" },
		JoinDifferentServerVersion: func(p *PoolFacts) { p.PlatformVersion = "3.3.0" },
		JoinHasSharedStorage:       func(p *PoolFacts) { p.SharedStorage = true },
		JoinHasRunningVMs:          func(p *PoolFacts) { p.RunningVMs = true },
	} {
		x := m
		change(&x)
		if got := CanJoinPool(x, c); got != want || JoinReasonText(got) == "" {
			t.Errorf("%v: %v %q", want, got, JoinReasonText(got))
		}
	}
	if CanJoinPool(m, PoolFacts{Key: "a", Hosts: 1, CPUVendor: "AuthenticAMD", APIMajor: 2, APIMinor: 21, DBSchema: "5.0", PlatformVersion: "3.4.0", TLSVerification: true}) != JoinTLSOnlyOnCoordinator {
		t.Error("tls")
	}
}

func TestCPUWarning(t *testing.T) {
	pool := PoolFacts{PoolCPU: map[string]string{"features_hvm": "178bfbff-f6f83203"}}
	same := PoolFacts{Name: "a", CPU: map[string]string{"features_hvm": "178bfbff-f6f83203"}}
	weaker := PoolFacts{Name: "b", CPU: map[string]string{"features_hvm": "178bfbff-f6f80203"}}
	if CPUWarning([]PoolFacts{same}, pool) != "" {
		t.Error("same CPU")
	}
	if w := CPUWarning([]PoolFacts{weaker}, pool); !strings.Contains(w, "'b'") || !strings.Contains(w, "newer CPUs") {
		t.Errorf("weaker: %q", w)
	}
	if CanJoinPool(PoolFacts{Key: "x", ReadError: "timeout"}, PoolFacts{Key: "y"}) != JoinNotRead || CanBeCoordinator(JoinNotRead) {
		t.Error("unread server")
	}
}
