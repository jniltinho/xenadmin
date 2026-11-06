package rules

import "testing"

func TestValidateMemory(t *testing.T) {
	gb := int64(1 << 30)
	m := MemoryVM{StaticMin: 256 * mib, StaticMax: 2 * gb, DynamicMin: 2 * gb, DynamicMax: 2 * gb, Running: true, Ratio: 0.25, MaxDynamicMin: 4 * gb, MaxAllowed: 1 << 40}
	if why := ValidateMemory(m, gb, 2*gb, 2*gb); why != "" {
		t.Errorf("a range within the limits: %s", why)
	}
	if ValidateMemory(m, 3*gb, 2*gb, 2*gb) == "" || ValidateMemory(m, 100*mib, 2*gb, 2*gb) == "" || ValidateMemory(m, 2*gb, 2*gb, 2<<40) == "" {
		t.Error("an invalid setting is accepted")
	}
	if ValidateMemory(m, 2*gb, 2*gb, 20*gb) == "" { // 20 GB * 0.25 > 4 GB free
		t.Error("more than the server has")
	}
	a, b, c := MemorySettings(m, 2*gb+10, 2*gb+10, 2*gb+10)
	if a != 2*gb || b != 2*gb || c != 2*gb {
		t.Errorf("the same MB keeps the bytes: %d %d %d", a, b, c)
	}
	if MemoryConfirmation(m, 4*gb) == "" || MemoryConfirmation(m, 2*gb) != "" {
		t.Error("confirmation")
	}
	if MemoryText(512*mib) != "512 MB" || MemoryText(3*gb/2) != "1.5 GB" || MemoryText(2*gb) != "2.0 GB" || MemoryText(1<<40) != "1024 GB" ||
		MemorySizeText(float64(2*gb), false) != "2 GB" || MemorySizeText(0, false) != "0" {
		t.Error("text", MemoryText(2*gb), MemoryText(1<<40))
	}
}
