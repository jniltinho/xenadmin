package rules

import (
	"strings"
	"testing"
)

func TestVappTexts(t *testing.T) {
	if q := VappDeleteQuestion("web", 0); !strings.Contains(q, "'web'") || strings.Contains(q, "VMs") {
		t.Error(q)
	}
	if q := VappDeleteQuestion("web", 3); !strings.Contains(q, "3 VMs") {
		t.Error(q)
	}
	if w := VappChangeWarning([]VappMove{{"a", ""}}, "web"); w != "" {
		t.Error("no vApp yet:", w)
	}
	if w := VappChangeWarning([]VappMove{{"b", "web"}}, "web"); w == "" { // another vApp of the same name
		t.Error("same name, another vApp")
	}
	if w := VappChangeWarning([]VappMove{{"a", "db"}}, "web"); !strings.Contains(w, "'a'") || !strings.Contains(w, "'db'") {
		t.Error(w)
	}
	if s := VappSummary("web", []string{"a", "b"}); !strings.Contains(s, "web") || !strings.Contains(s, "a, b") || strings.Contains(s, `\r`) {
		t.Error(s)
	}
}
