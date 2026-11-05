package rules

import "testing"

// The XML of the KVM lab's iSCSI target (tgt), as the SM returns it in the refused SR.create.
func TestParseIQNsAndLUNs(t *testing.T) {
	iqns, err := ParseIQNs("<?xml version=\"1.0\" ?>\n<iscsi-target-iqns>\n\t<TGT>\n\t\t<Index>1</Index>\n\t\t<IPAddress>10.77.58.20</IPAddress>\n\t\t<TargetIQN>*</TargetIQN>\n\t</TGT>\n\t<TGT>\n\t\t<Index>0</Index>\n\t\t<IPAddress>fe80::1</IPAddress>\n\t\t<Port>3261</Port>\n\t\t<TargetIQN>iqn.2026-11.org.xcp-ng.lab:sr</TargetIQN>\n\t</TGT>\n</iscsi-target-iqns>\n\n")
	if err != nil || len(iqns) != 2 || iqns[0].IQN != "iqn.2026-11.org.xcp-ng.lab:sr" || iqns[0].Address != "[fe80::1]" || iqns[0].Port != 3261 ||
		iqns[1].IQN != "*" || iqns[1].Port != DefaultISCSIPort || iqns[0].Text != "iqn.2026-11.org.xcp-ng.lab:sr ([fe80::1]:3261)" {
		t.Fatalf("%+v %v", iqns, err)
	}
	luns, err := ParseLUNs("<?xml version=\"1.0\" ?>\n<iscsi-target>\n\t<LUN>\n\t\t<vendor>IET</vendor>\n\t\t<serial>beaf11</serial>\n\t\t<LUNid>1</LUNid>\n\t\t<size>4294967296</size>\n\t\t<SCSIid>360000000000000000e00000000010001</SCSIid>\n\t</LUN>\n</iscsi-target>\n\n")
	if err != nil || len(luns) != 1 || luns[0].ScsiID != "360000000000000000e00000000010001" || luns[0].Text != "LUN 1: beaf11: 4.0 GB (IET)" {
		t.Fatalf("%+v %v", luns, err)
	}
	d := ISCSIDeviceConfig(iqns[0], luns[0], "u", "p")
	if d["SCSIid"] != luns[0].ScsiID || d["port"] != "3261" || d["chappassword"] != "p" {
		t.Errorf("%v", d)
	}
	if !naturalLess("lun2", "lun10") || naturalLess("b", "a") || !naturalLess("x١", "x٢") {
		t.Error("natural order")
	}
}
