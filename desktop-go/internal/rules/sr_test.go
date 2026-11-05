package rules

import "testing"

func TestSRRulesMirrorTheCSharp(t *testing.T) {
	if !SRValidLocation(SRNfs, " 10.0.0.1:/export ") || SRValidLocation(SRNfs, "10.0.0.1/export") || !SRValidLocation(SRSmb, `\\srv\share`) || SRValidLocation(SRSmb, "//srv/share") {
		t.Error("locations")
	}
	if SRDefaultName(SRNfs, map[string]bool{"NFS virtual disk storage": true}) != "NFS virtual disk storage (1)" {
		t.Error("default name")
	}
	d := SRDeviceConfig(SRNfs, "srv:/export", "", "4.1", "", "")
	if d["server"] != "srv" || d["serverpath"] != "/export" || d["nfsversion"] != "4.1" {
		t.Errorf("%v", d)
	}
	iso := SRDeviceConfig(SRSmbIso, `\\srv\share\isos\linux`, "", "", "u", "p")
	if iso["location"] != "//srv/share" || iso["iso_path"] != "/isos/linux" || iso["cifspassword"] != "p" {
		t.Errorf("%v", iso)
	}
	if SRAutoDescription(SRNfs, "srv:/export") != "NFS SR [srv:/export]" || SRType(SRNfsIso) != "iso" {
		t.Error("description or type")
	}
}
