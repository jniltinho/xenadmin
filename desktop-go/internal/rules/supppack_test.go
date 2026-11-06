package rules

import (
	"strings"
	"testing"
)

func TestPackWarning(t *testing.T) {
	pack := map[string]string{"xs:pv": "PV drivers, version 1.2, build 3, homogeneous", "platform_version": "3.4.0"}
	other := map[string]string{"xs:pv": "PV drivers, version 1.3, build 1, homogeneous"}
	pool := PoolFacts{ServerPacks: map[string]map[string]string{"a": pack}}
	if w := PackWarning([]PoolFacts{{ServerPacks: map[string]map[string]string{"b": pack}}}, pool); w != "" {
		t.Errorf("same pack: %q", w)
	}
	if w := PackWarning([]PoolFacts{{ServerPacks: map[string]map[string]string{"b": {}}}}, pool); !strings.Contains(w, "PV drivers (not installed on b)") {
		t.Errorf("missing: %q", w)
	}
	if w := PackWarning([]PoolFacts{{ServerPacks: map[string]map[string]string{"b": other}}}, pool); !strings.Contains(w, "not the same version") {
		t.Errorf("versions: %q", w)
	}
	if _, ok := parseSuppPack("xs:x", "X, version 1, build 2, nope"); ok {
		t.Error("a bad flag parses")
	}
}
