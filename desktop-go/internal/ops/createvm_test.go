package ops

import "testing"

func TestCreateVMHelpers(t *testing.T) {
	for in, want := range map[string]string{"": "ncd", "cd": "ncd", "ncd": "ncd", "cdn": "ncd", "DN": "nd"} {
		if got := networkFirst(in); got != want {
			t.Errorf("networkFirst(%q) = %q", in, got)
		}
	}
	rec := `<restrictions><restriction field="has-vendor-device" value="true"/><restriction field="vcpus-max" max="32"/></restrictions>`
	if !hasRecommendation(rec, "has-vendor-device") || hasRecommendation(rec, "allow-gpu") {
		t.Error("hasRecommendation")
	}
}
