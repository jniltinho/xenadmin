package rules

import "testing"

func TestValidateVIF(t *testing.T) {
	for _, c := range []struct {
		net, manual bool
		mac         string
		qos         bool
		limit       string
		ok          bool
	}{
		{true, false, "", false, "", true},
		{false, false, "", false, "", false},
		{true, true, "", false, "", false},
		{true, true, "aa:bb:cc:dd:ee:0f", false, "", true},
		{true, true, "aa:bb:cc:dd:ee", false, "", false},
		{true, false, "", true, "100", true},
		{true, false, "", true, "0", false},
		{true, false, "", true, "x", false},
	} {
		if _, ok := ValidateVIF(c.net, c.manual, c.mac, c.qos, c.limit); ok != c.ok {
			t.Errorf("%+v: %v", c, ok)
		}
	}
	if VIFDevice([]string{"0", "1", "3"}) != 2 || VIFDevice(nil) != 0 {
		t.Error("device")
	}
	if e, _, _ := VIFTooltips(false, true, true, false, true); e != "" {
		t.Error("no tooltip with drivers")
	}
	if _, r, g := VIFTooltips(true, false, false, false, false); r == "" || g == "" {
		t.Error("suspended")
	}
}
