package app

import (
	"os"
	"testing"

	xenapi "go/xenapi"

	"xcpng-center/internal/ops"
	"xcpng-center/internal/rules"
)

// TestLabVMHardware reads a VM's GPU and USB pages, the Attach USB and Add virtual GPU dialogs, the New Network
// wizard's SR-IOV choice and the NICs tab's SR-IOV column on the lab (read-only): the KVM lab has no GPU that does
// vGPU, no USB device and no SR-IOV NIC, so the pages are hidden and SR-IOV is not offered, as the C# rules say. The
// writes are only asked for things they must refuse before writing anything.
func TestLabVMHardware(t *testing.T) {
	c := labConnection(t)
	s := c.Session
	name := os.Getenv("XCP_LAB_VM")
	if name == "" {
		t.Skip("XCP_LAB_VM not set")
	}
	vms, err := ListObjects(c, "vms")
	if err != nil {
		t.Fatal(err)
	}
	uuid := ""
	for _, vm := range vms {
		if vm.Name == name {
			uuid = vm.UUID
		}
	}
	if uuid == "" {
		t.Fatalf("no VM %q", name)
	}

	// the GPU page: shown only when the VM can have a GPU and a group has GPUs with supported types
	groups, err := xenapi.GPUGroup.GetAllRecords(s)
	if err != nil {
		t.Fatal(err)
	}
	available := false
	for _, g := range groups {
		available = available || len(g.PGPUs) > 0 && len(g.SupportedVGPUTypes) > 0
	}
	gpu, err := ReadVMGPU(c, uuid)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("GPU page: %+v (%d GPU groups)", gpu, len(groups))
	if gpu.Shown != available || gpu.Input.GPUsAvailable != available {
		t.Errorf("GPU page shown %v, GPUs available %v; the pool's groups say %v", gpu.Shown, gpu.Input.GPUsAvailable, available)
	}
	if available {
		t.Log("this pool has GPUs that do vGPU: the no-GPU checks below are skipped")
	} else {
		if gpu.Shown || len(gpu.Rows) != 0 {
			t.Errorf("no GPU: the page is hidden and the VM has no vGPU: %+v", gpu)
		}
		picks, err := ReadVGPUChoices(c, uuid, nil)
		if err != nil || len(picks.Items) != 0 || picks.Picked != -1 {
			t.Errorf("Add virtual GPU offers %+v (%v)", picks, err)
		}
		if st := rules.GPUPage(gpu.Input); st.CanAdd {
			t.Errorf("Add enabled without GPUs: %+v", st)
		}
	}
	// refused before any write: no such GPU group or type (and a running VM)
	if err := ops.AssignVGPUs(c, uuid, nil, []ops.NewVGPU{{Group: "no-such-group", Type: "no-such-type"}}); err == nil {
		t.Error("a vGPU of no GPU group was accepted")
	} else {
		t.Logf("AssignVGPUs refused: %v", err)
	}

	// the USB page: shown only with a USB device in the pool
	pusbs, err := xenapi.PUSB.GetAllRecords(s)
	if err != nil {
		t.Fatal(err)
	}
	usb, err := ReadVMUSB(c, uuid, "")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("USB page: %+v (%d devices in the pool)", usb, len(pusbs))
	if len(pusbs) == 0 && usb.Shown {
		t.Error("the USB page is shown without a USB device in the pool")
	}
	attach, err := ReadUSBAttach(c, uuid)
	if err != nil {
		t.Fatal(err)
	}
	if len(pusbs) == 0 && len(attach.Hosts) != 0 || attach.Empty != "No devices available" {
		t.Errorf("Attach USB: %+v", attach)
	}
	if err := ops.AttachUSB(c, uuid, "no-such-device"); err == nil {
		t.Error("a USB device that does not exist was attached")
	}

	// the New Network wizard's SR-IOV Network and the NICs tab's SR-IOV column
	pifs, err := xenapi.PIF.GetAllRecords(s)
	if err != nil {
		t.Fatal(err)
	}
	capable := false
	for _, p := range pifs {
		capable = capable || rules.SriovCapable(p.Capabilities)
	}
	pools, err := ListObjects(c, "pools")
	if err != nil || len(pools) == 0 {
		t.Fatalf("pools %v %v", pools, err)
	}
	nn, err := NewNetworkOptions(c, "pools", pools[0].UUID)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("SR-IOV: offered %v (%q), NICs %+v, name %q", nn.Sriov, nn.SriovWhy, nn.SriovNICs, nn.SriovName)
	if !capable {
		if nn.Sriov || len(nn.SriovNICs) != 0 || nn.SriovWhy != "SR-IOV networks require an SR-IOV capable NIC." {
			t.Errorf("no SR-IOV NIC: %v %q %+v", nn.Sriov, nn.SriovWhy, nn.SriovNICs)
		}
		if err := ops.CreateSriov(c, ops.NewSriov{Name: "never made", Host: nn.Host, NIC: "no-such-nic"}); err == nil {
			t.Error("an SR-IOV network was made without an SR-IOV NIC")
		}
	}
	hosts, err := ListObjects(c, "hosts")
	if err != nil {
		t.Fatal(err)
	}
	for _, h := range hosts {
		nics, err := HostNicsTab(c, h.UUID)
		if err != nil {
			t.Fatal(err)
		}
		col := len(nics.Table.Headers) - 1
		if nics.Table.Headers[col] != "SR-IOV Capable" {
			t.Fatalf("NICs headers %v", nics.Table.Headers)
		}
		for _, r := range nics.Table.Rows {
			t.Logf("%s %s: SR-IOV %q", h.Name, r[0], r[col])
			if !capable && r[col] != "No" {
				t.Errorf("%s %s: SR-IOV %q", h.Name, r[0], r[col])
			}
		}
	}
}
