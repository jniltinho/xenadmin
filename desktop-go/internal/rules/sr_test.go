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

func TestParseSRList(t *testing.T) {
	r := "<value>&lt;?xml version=&quot;1.0&quot; ?&gt;\n&lt;SRlist&gt;\n\t&lt;SR&gt;\n\t\t&lt;UUID&gt;f3a2ebd8-58f3&lt;/UUID&gt;\n\t&lt;/SR&gt;\n\t&lt;SupportedVersions&gt;\n\t\t&lt;Version&gt;3&lt;/Version&gt;\n\t&lt;/SupportedVersions&gt;\n&lt;/SRlist&gt;\n</value>"
	srs, err := ParseSRList(r)
	if err != nil || len(srs) != 1 || srs[0].UUID != "f3a2ebd8-58f3" {
		t.Fatalf("%+v %v", srs, err)
	}
	if srs, err := ParseSRList("<value>&lt;SRlist/&gt;</value>"); err != nil || len(srs) != 0 {
		t.Fatalf("%+v %v", srs, err)
	}
}

func TestSRCommands(t *testing.T) {
	nfs := SRState{Type: "nfs", HasPBDs: true, Allowed: []string{"forget", "destroy"}, VDICreate: true, Shared: true, Hosts: 2}
	if got := SRCommands(nfs); len(got) != 4 {
		t.Errorf("%v", got)
	}
	nfs.RunningVMs, nfs.Default = true, true
	if got := SRCommands(nfs); len(got) != 0 {
		t.Errorf("busy default SR: %v", got)
	}
	if got := SRCommands(SRState{Type: "nfs", HasPBDs: true, Detached: true, Shared: true, Hosts: 2}); len(got) != 1 || got[0] != "repair" {
		t.Errorf("detached SR: %v", got)
	}
	local := SRState{Type: "ext", HasPBDs: true, VDICreate: true, Hosts: 2, Allowed: []string{"forget", "destroy"}}
	if got := SRCommands(local); len(got) != 1 || got[0] != "detach" {
		t.Errorf("local SR of a pool: %v", got)
	}
	if title, text, yes := SRConfirmation("destroy", "lab"); title == "" || text == "" || yes == "" {
		t.Error("destroy question")
	}
}

func TestSRIcon(t *testing.T) {
	for want, st := range map[string]SRState{
		"000_StorageDisabled_h32bit_16": {},
		"000_StorageBroken_h32bit_16":   {HasPBDs: true, Degraded: true, Default: true},
		"000_StorageDefault_h32bit_16":  {HasPBDs: true, Default: true},
		"000_Storage_h32bit_16":         {HasPBDs: true},
	} {
		if got := SRIcon(st); got != want {
			t.Errorf("%+v: %s", st, got)
		}
	}
}

func TestMultipathDown(t *testing.T) {
	if !MultipathDown("[1, 2]") || MultipathDown("[2L, 2L]") || MultipathDown("") {
		t.Error("path counts")
	}
}
