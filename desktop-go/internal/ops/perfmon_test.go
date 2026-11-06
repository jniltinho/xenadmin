package ops

import "testing"

func TestGrants(t *testing.T) {
	for _, c := range []struct {
		perm, call string
		want       bool
	}{
		{"VM.add_to_other_config", "vm.add_to_other_config/key:perfmon", true},
		{"vm.add_to_other_config/key:perfmon", "vm.add_to_other_config/key:perfmon", true},
		{"vm.add_to_other_config/key:perf*", "vm.add_to_other_config/key:perfmon", true},
		{"vm.add_to_other_config/key:folder", "vm.add_to_other_config/key:perfmon", false},
		{"vm.add_to_other_config/key:perfmon", "vm.add_to_other_config", false},
		{"host.set_logging", "host.set_logging", true},
	} {
		if got := grants(c.perm, c.call); got != c.want {
			t.Errorf("grants(%q, %q) = %v", c.perm, c.call, got)
		}
	}
}
