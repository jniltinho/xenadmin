package app

import (
	"context"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	xenapi "go/xenapi"

	"xcpng-center/internal/config"
	"xcpng-center/internal/ops"
	"xcpng-center/internal/rules"
	"xcpng-center/internal/trust"
	"xcpng-center/internal/xapi"
)

// labConnection logs in to XCP_LAB_HOST (a throwaway lab server) whose certificate is XCP_TRUST_FINGERPRINT.
// Without XCP_LAB_HOST and XCP_PASS the lab tests are skipped.
func labConnection(t *testing.T) *xapi.Connection {
	host, pass := os.Getenv("XCP_LAB_HOST"), os.Getenv("XCP_PASS")
	if host == "" || pass == "" {
		t.Skip("XCP_LAB_HOST and XCP_PASS not set")
	}
	s, err := config.Load(filepath.Join(t.TempDir(), "settings.toml"))
	if err != nil {
		t.Fatal(err)
	}
	h, p, err := xapi.ParseAddress(host)
	if err != nil {
		t.Fatal(err)
	}
	c, err := xapi.Connect(context.Background(), &trust.Policy{Store: s, OperatorPins: map[string]string{"*": os.Getenv("XCP_TRUST_FINGERPRINT")}}, h, p, "root", pass)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { c.Close() })
	return c
}

// labPolicy is the trust of the lab's certificates for the HTTP transfers (the coordinator's pinned, members vouched).
func labPolicy(t *testing.T) *trust.Policy {
	s, err := config.Load(filepath.Join(t.TempDir(), "settings.toml"))
	if err != nil {
		t.Fatal(err)
	}
	h, p, err := xapi.ParseAddress(os.Getenv("XCP_LAB_HOST"))
	if err != nil {
		t.Fatal(err)
	}
	// only the coordinator is pinned: the members are trusted through it (VouchedFingerprints), as in the window
	return &trust.Policy{Store: s, OperatorPins: map[string]string{xapi.HostKey(h, p): os.Getenv("XCP_TRUST_FINGERPRINT")}}
}

// TestLabVMTabs reads the Storage, Networking and Memory tabs of every VM of the lab (read-only).
func TestLabVMTabs(t *testing.T) {
	c := labConnection(t)
	vms, err := ListObjects(c, "vms")
	if err != nil {
		t.Fatal(err)
	}
	for _, vm := range vms {
		st, err := ListVMStorage(c, vm.UUID)
		if err != nil {
			t.Fatalf("%s storage: %v", vm.Name, err)
		}
		t.Logf("%s: disks %v, dvd %q iso %q, isos %v", vm.Name, st.Table.Rows, st.DVDName, st.ISO, st.ISOs)
		if _, err := ListVMNetworking(c, vm.UUID); err != nil {
			t.Fatalf("%s networking: %v", vm.Name, err)
		}
		if _, err := ListVMMemory(c, vm.UUID); err != nil {
			t.Fatalf("%s memory: %v", vm.Name, err)
		}
	}
}

// TestLabChangeISO inserts and ejects an ISO in the DVD drive of the lab's throwaway VM (XCP_LAB_VM, writes only
// with XCP_ALLOW_WRITE=1).
func TestLabChangeISO(t *testing.T) {
	c := labConnection(t)
	name := os.Getenv("XCP_LAB_VM")
	if name == "" || os.Getenv("XCP_ALLOW_WRITE") != "1" {
		t.Skip("XCP_LAB_VM and XCP_ALLOW_WRITE=1 not set")
	}
	vms, _ := ListObjects(c, "vms")
	for _, vm := range vms {
		if vm.Name != name {
			continue
		}
		st, err := ListVMStorage(c, vm.UUID)
		if err != nil || st.DVD == "" || len(st.ISOs) < 2 {
			t.Fatalf("no DVD drive or ISO: %+v %v", st, err)
		}
		iso := st.ISOs[len(st.ISOs)-1]
		if err := ops.ChangeISO(c, st.DVD, iso.VDI); err != nil {
			t.Fatal(err)
		}
		if st, _ = ListVMStorage(c, vm.UUID); st.ISO != iso.VDI {
			t.Fatalf("inserted %q, drive holds %q", iso.VDI, st.ISO)
		}
		if err := ops.ChangeISO(c, st.DVD, ""); err != nil {
			t.Fatal(err)
		}
		if st, _ = ListVMStorage(c, vm.UUID); st.ISO != "" {
			t.Fatalf("still holds %q after eject", st.ISO)
		}
		return
	}
	t.Fatalf("no VM named %q", name)
}

// TestLabDiskLifecycle adds a 64 MB disk to the throwaway VM, detaches it and deletes it.
func TestLabDiskLifecycle(t *testing.T) {
	c := labConnection(t)
	name := os.Getenv("XCP_LAB_VM")
	if name == "" || os.Getenv("XCP_ALLOW_WRITE") != "1" {
		t.Skip("XCP_LAB_VM and XCP_ALLOW_WRITE=1 not set")
	}
	vms, _ := ListObjects(c, "vms")
	for _, vm := range vms {
		if vm.Name != name {
			continue
		}
		opts, err := NewDiskOptions(c, vm.UUID)
		if err != nil || opts.DefaultSR == "" || opts.Name == "" {
			t.Fatalf("options %+v %v", opts, err)
		}
		notice, err := ops.AddDisk(c, vm.UUID, "go-lab-disk", "lab test", opts.DefaultSR, 64<<20)
		if err != nil {
			t.Fatal(err)
		}
		st, _ := ListVMStorage(c, vm.UUID)
		var d Disk
		for _, x := range st.Disks {
			if x.Name == "go-lab-disk" {
				d = x
			}
		}
		if d.VBD == "" || (notice == "" && !d.Attached) || d.System {
			t.Fatalf("added disk %+v (notice %q)", d, notice)
		}
		if !d.CanDetach || d.CanDelete != !d.Attached {
			t.Fatalf("rules for the new disk: %+v", d)
		}
		if d.Attached {
			if err := ops.SetDiskActive(c, d.VBD, false); err != nil {
				t.Fatal(err)
			}
		}
		if err := ops.DeleteDisk(c, vm.UUID, d.VBD); err != nil {
			t.Fatal(err)
		}
		st, _ = ListVMStorage(c, vm.UUID)
		for _, x := range st.Disks {
			if x.Name == "go-lab-disk" {
				t.Fatal("the disk is still attached")
			}
		}
		// Detach keeps the disk: it leaves the VM and stays on the SR (removed here by the test)
		if _, err := ops.AddDisk(c, vm.UUID, "go-lab-detach", "", opts.DefaultSR, 64<<20); err != nil {
			t.Fatal(err)
		}
		st, _ = ListVMStorage(c, vm.UUID)
		for _, x := range st.Disks {
			if x.Name == "go-lab-detach" {
				d = x
			}
		}
		if err := ops.DetachDisk(c, d.VBD); err != nil {
			t.Fatal(err)
		}
		// the detached disk is offered by Attach Disk, and attaches again, then goes for good (read-only is refused by XAPI for HVM guests)
		choices, err := AttachableDisks(c, vm.UUID)
		found := false
		for _, ch := range choices {
			found = found || ch.VDI == d.VDI
		}
		if err != nil || !found {
			t.Fatalf("the detached disk is not offered: %v %v", choices, err)
		}
		if _, err := ops.AttachDisk(c, vm.UUID, d.VDI, false); err != nil {
			t.Fatal(err)
		}
		if _, err := ops.AttachDisk(c, vm.UUID, d.VDI, false); err == nil {
			t.Fatal("the same disk attached twice")
		}
		st, _ = ListVMStorage(c, vm.UUID)
		for _, x := range st.Disks {
			if x.VDI == d.VDI {
				d = x
			}
		}
		if d.Attached {
			_ = ops.SetDiskActive(c, d.VBD, false)
		}
		if err := ops.DetachDisk(c, d.VBD); err != nil {
			t.Fatal(err)
		}
		vdi, err := xenapi.VDI.GetByUUID(c.Session, d.VDI)
		if err != nil {
			t.Fatalf("the detached disk is gone: %v", err)
		}
		if vbds, _ := xenapi.VDI.GetVBDs(c.Session, vdi); len(vbds) != 0 {
			t.Fatal("the detached disk still has a VBD")
		}
		if err := xenapi.VDI.Destroy(c.Session, vdi); err != nil {
			t.Fatal(err)
		}
		return
	}
	t.Fatalf("no VM named %q", name)
}

// TestLabProperties reads the throwaway VM's Properties and changes its description and home server, then restores
// them.
func TestLabProperties(t *testing.T) {
	c := labConnection(t)
	name := os.Getenv("XCP_LAB_VM")
	if name == "" || os.Getenv("XCP_ALLOW_WRITE") != "1" {
		t.Skip("XCP_LAB_VM and XCP_ALLOW_WRITE=1 not set")
	}
	vms, _ := ListObjects(c, "vms")
	for _, vm := range vms {
		if vm.Name != name {
			continue
		}
		p, err := ReadVMProperties(c, vm.UUID)
		if err != nil || len(p.VCPUChoices) == 0 || len(p.Servers) == 0 {
			t.Fatalf("%+v %v", p, err)
		}
		desc, home := "go lab test", p.Servers[0].UUID
		if err := ops.SaveVMProperties(c, vm.UUID, ops.VMChanges{Description: &desc, HomeServer: &home}); err != nil {
			t.Fatal(err)
		}
		q, _ := ReadVMProperties(c, vm.UUID)
		if q.Description != desc || q.HomeServer != home {
			t.Fatalf("saved %+v", q)
		}
		// a running VM keeps its vCPU weight and topology (only a halted one changes them)
		prio := q.Priority + 1
		if err := ops.SaveVMProperties(c, vm.UUID, ops.VMChanges{Priority: &prio}); err == nil && q.Running {
			t.Fatal("vCPU priority changed on a running VM")
		}
		if err := ops.SaveVMProperties(c, vm.UUID, ops.VMChanges{Description: &p.Description, HomeServer: &p.HomeServer}); err != nil {
			t.Fatal(err)
		}
		return
	}
	t.Fatalf("no VM named %q", name)
}

// TestLabCreateVM creates a VM from "Other install media" with the lab's Alpine ISO and a 1 GB disk, starts it, then
// destroys it (XCP_ALLOW_WRITE=1; the VM is named xcpng-center-go-newvm).
func TestLabCreateVM(t *testing.T) {
	c := labConnection(t)
	if os.Getenv("XCP_ALLOW_WRITE") != "1" {
		t.Skip("XCP_ALLOW_WRITE=1 not set")
	}
	opts, err := ReadNewVMOptions(c)
	if err != nil {
		t.Fatal(err)
	}
	var tpl, iso string
	for _, x := range opts.Templates {
		if x.Name == "Other install media" {
			tpl = x.UUID
		}
	}
	for _, x := range opts.ISOs {
		if strings.Contains(x.Text, "alpine") {
			iso = x.VDI
		}
	}
	if tpl == "" || iso == "" || opts.DefaultSR == "" {
		t.Fatalf("template %q iso %q sr %q", tpl, iso, opts.DefaultSR)
	}
	d, err := ReadNewVMDefaults(c, tpl, "")
	if err != nil || !d.HVM || d.Install != "cd" {
		t.Fatalf("%+v %v", d, err)
	}
	spec := ops.NewVM{Template: tpl, Name: "xcpng-center-go-newvm", Install: "cd", ISO: iso, BootMode: d.BootMode, VCPUs: 1, Cores: 1,
		Memory: 256 << 20, Networks: d.Networks, Start: true,
		Disks: []ops.NewDisk{{Position: "0", Name: "xcpng-center-go-newvm 0", Size: 1 << 30, SR: opts.DefaultSR}}}
	created, err := ops.CreateVM(c, spec)
	uuid := created.UUID
	if uuid != "" {
		defer func() {
			ref, _ := xenapi.VM.GetByUUID(c.Session, uuid)
			rec, _ := xenapi.VM.GetRecord(c.Session, ref)
			if rec.PowerState != xenapi.VMPowerStateHalted {
				_ = xenapi.VM.HardShutdown(c.Session, ref)
			}
			for _, b := range rec.VBDs {
				if vbd, err := xenapi.VBD.GetRecord(c.Session, b); err == nil && vbd.OtherConfig["owner"] != "" {
					_ = xenapi.VDI.Destroy(c.Session, vbd.VDI)
				}
			}
			_ = xenapi.VM.Destroy(c.Session, ref)
		}()
	}
	if err != nil || created.StartError != "" {
		t.Fatal(err, created.StartError)
	}
	vms, _ := ListObjects(c, "vms")
	for _, v := range vms {
		if v.UUID == uuid {
			if v.Name != spec.Name || v.PowerState != "Running" {
				t.Fatalf("created %+v", v)
			}
			st, _ := ListVMStorage(c, uuid)
			if len(st.Disks) != 1 || st.ISO != iso {
				t.Fatalf("storage %+v", st)
			}
			return
		}
	}
	t.Fatal("the new VM is not listed")
}

// TestLabCreateSR creates an NFS ISO library on XCP_LAB_NFS/iso (a second one on that share is refused) and an NFS SR
// on XCP_LAB_NFS (server:/path, XCP_ALLOW_WRITE=1): shared, plugged and not the default (the lab has a shared SR
// already); then forgets the library and destroys the SR.
func TestLabCreateSR(t *testing.T) {
	c := labConnection(t)
	share := os.Getenv("XCP_LAB_NFS")
	if share == "" || os.Getenv("XCP_ALLOW_WRITE") != "1" {
		t.Skip("XCP_LAB_NFS and XCP_ALLOW_WRITE=1 not set")
	}
	s := c.Session
	iso := ops.NewSR{Kind: rules.SRNfsIso, Name: "xcpng-center-go-test-iso", AutoDesc: true, Path: share + "/iso"}
	isoUUID, err := ops.CreateSR(c, iso)
	if isoUUID != "" {
		defer func() {
			ref, _ := xenapi.SR.GetByUUID(s, isoUUID)
			pbds, _ := xenapi.SR.GetPBDs(s, ref)
			for _, p := range pbds {
				_ = xenapi.PBD.Unplug(s, p)
			}
			if err := xenapi.SR.Forget(s, ref); err != nil {
				t.Error(err)
			}
		}()
	}
	if err != nil {
		t.Fatal(err)
	}
	if ref, _ := xenapi.SR.GetByUUID(s, isoUUID); true {
		if rec, _ := xenapi.SR.GetRecord(s, ref); rec.OtherConfig["auto-scan"] != "true" || rec.ContentType != "iso" {
			t.Fatalf("%+v", rec)
		}
	}
	if _, err := ops.CreateSR(c, iso); err == nil || !strings.Contains(err.Error(), "already") {
		t.Fatalf("ISO library attached twice: %v", err)
	}
	kinds, err := ReadNewSRKinds(c)
	if err != nil || len(kinds) != len(rules.SRKinds) || kinds[0].DefaultName == "" {
		t.Fatalf("%+v %v", kinds, err)
	}
	pools, _ := xenapi.Pool.GetAllRecords(s)
	uuid, err := ops.CreateSR(c, ops.NewSR{Kind: rules.SRNfs, Name: "xcpng-center-go-test-sr", AutoDesc: true, Path: share, NFSVersion: "4"})
	if uuid != "" {
		defer func() {
			ref, _ := xenapi.SR.GetByUUID(s, uuid)
			pbds, _ := xenapi.SR.GetPBDs(s, ref)
			for _, p := range pbds {
				_ = xenapi.PBD.Unplug(s, p)
			}
			if err := xenapi.SR.Destroy(s, ref); err != nil {
				t.Error(err)
			}
		}()
	}
	if err != nil {
		t.Fatal(err)
	}
	ref, _ := xenapi.SR.GetByUUID(s, uuid)
	rec, err := xenapi.SR.GetRecord(s, ref)
	if err != nil || !rec.Shared || rec.Type != "nfs" || rec.NameDescription != "NFS SR ["+share+"]" || rec.OtherConfig["auto-scan"] != "false" {
		t.Fatalf("%+v %v", rec, err)
	}
	for _, p := range rec.PBDs {
		if pbd, _ := xenapi.PBD.GetRecord(s, p); !pbd.CurrentlyAttached || pbd.DeviceConfig["nfsversion"] != "4" {
			t.Fatalf("PBD %+v", pbd)
		}
	}
	for ref, p := range pools {
		if now, _ := xenapi.Pool.GetDefaultSR(s, ref); now != p.DefaultSR {
			t.Fatal("the default SR changed")
		}
	}
}

// TestLabCreateSMBSR creates an SMB SR on XCP_LAB_SMB (\\server\share, user xcp with the lab password, which goes
// through a secret) and destroys it (XCP_ALLOW_WRITE=1); no secret is left behind.
func TestLabCreateSMBSR(t *testing.T) {
	c := labConnection(t)
	share := os.Getenv("XCP_LAB_SMB")
	if share == "" || os.Getenv("XCP_ALLOW_WRITE") != "1" {
		t.Skip("XCP_LAB_SMB and XCP_ALLOW_WRITE=1 not set")
	}
	s := c.Session
	secrets, _ := xenapi.Secret.GetAll(s)
	uuid, err := ops.CreateSR(c, ops.NewSR{Kind: rules.SRSmb, Name: "xcpng-center-go-test-smb", AutoDesc: true, Path: share, User: "xcp", Password: os.Getenv("XCP_PASS")})
	if uuid != "" {
		defer func() {
			ref, _ := xenapi.SR.GetByUUID(s, uuid)
			pbds, _ := xenapi.SR.GetPBDs(s, ref)
			for _, p := range pbds {
				_ = xenapi.PBD.Unplug(s, p)
			}
			if err := xenapi.SR.Destroy(s, ref); err != nil {
				t.Error(err)
			}
		}()
	}
	if err != nil {
		t.Fatal(err)
	}
	ref, _ := xenapi.SR.GetByUUID(s, uuid)
	if rec, _ := xenapi.SR.GetRecord(s, ref); rec.Type != "smb" || !rec.Shared || len(rec.PBDs) != 2 {
		t.Fatalf("%+v", rec)
	}
	// the PBDs keep their own copy of the secret (SR.create); the one given to SR.create is gone
	pbds, _ := xenapi.SR.GetPBDs(s, ref)
	pbd, _ := xenapi.PBD.GetRecord(s, pbds[0])
	if pbd.DeviceConfig["password"] != "" || pbd.DeviceConfig["password_secret"] == "" {
		t.Fatal("the password is in the device_config")
	}
	after, _ := xenapi.Secret.GetAll(s)
	if len(after) > len(secrets)+len(pbds) {
		t.Fatalf("secrets: %d before, %d after", len(secrets), len(after))
	}
}

// TestLabReattachSR creates an SR on XCP_LAB_NFS or, with XCP_LAB_SMB, on that SMB share (XCP_ALLOW_WRITE=1), detaches
// it, finds it with a scan and reattaches it (SrReattachAction), forgets it and attaches it again (SrIntroduceAction),
// then destroys it. On SMB the password never stays in a PBD's device_config.
func TestLabReattachSR(t *testing.T) {
	c := labConnection(t)
	spec := ops.NewSR{Kind: rules.SRNfs, Name: "xcpng-center-go-test-reattach", AutoDesc: true, Path: os.Getenv("XCP_LAB_NFS")}
	if smb := os.Getenv("XCP_LAB_SMB"); smb != "" {
		spec = ops.NewSR{Kind: rules.SRSmb, Name: spec.Name, AutoDesc: true, Path: smb, User: "xcp", Password: os.Getenv("XCP_PASS")}
	}
	if spec.Path == "" || os.Getenv("XCP_ALLOW_WRITE") != "1" {
		t.Skip("XCP_LAB_NFS or XCP_LAB_SMB, and XCP_ALLOW_WRITE=1, not set")
	}
	s := c.Session
	uuid, err := ops.CreateSR(c, spec)
	if err != nil {
		t.Fatal(err)
	}
	unplug := func() xenapi.SRRef {
		ref, err := xenapi.SR.GetByUUID(s, uuid)
		if err != nil {
			t.Fatal(err)
		}
		pbds, _ := xenapi.SR.GetPBDs(s, ref)
		for _, p := range pbds {
			_ = xenapi.PBD.Unplug(s, p)
			_ = xenapi.PBD.Destroy(s, p)
		}
		return ref
	}
	defer func() {
		if ref, err := xenapi.SR.GetByUUID(s, uuid); err == nil {
			pbds, _ := xenapi.SR.GetPBDs(s, ref)
			for _, p := range pbds {
				_ = xenapi.PBD.Unplug(s, p)
			}
			if len(pbds) == 0 { // SR.destroy needs a plugged PBD: attach it once more
				spec.Existing = uuid
				_, _ = ops.CreateSR(c, spec)
				pbds, _ = xenapi.SR.GetPBDs(s, ref)
				for _, p := range pbds {
					_ = xenapi.PBD.Unplug(s, p)
				}
			}
			if err := xenapi.SR.Destroy(s, ref); err != nil {
				t.Error(err)
			}
		}
	}()
	attach := func(path string) {
		found, err := ops.ProbeSR(c, spec)
		if err != nil {
			t.Fatal(err)
		}
		ok := false
		for _, f := range found {
			ok = ok || f.UUID == uuid && f.InUse == ""
		}
		if !ok {
			t.Fatalf("%s: the scan did not offer %s: %+v", path, uuid, found)
		}
		spec.Existing, spec.Name = uuid, "xcpng-center-go-test-reattach-"+path
		if _, err := ops.CreateSR(c, spec); err != nil {
			t.Fatal(path, err)
		}
		ref, _ := xenapi.SR.GetByUUID(s, uuid)
		rec, _ := xenapi.SR.GetRecord(s, ref)
		hosts, _ := xenapi.Host.GetAll(s)
		if rec.NameLabel != spec.Name || len(rec.PBDs) != len(hosts) {
			t.Fatalf("%s: %+v", path, rec)
		}
		for _, p := range rec.PBDs {
			if pbd, _ := xenapi.PBD.GetRecord(s, p); !pbd.CurrentlyAttached || pbd.DeviceConfig["password"] != "" {
				t.Fatalf("%s: PBD %v %v", path, pbd.CurrentlyAttached, pbd.DeviceConfig["password"] != "")
			}
		}
		if found, _ := ops.ProbeSR(c, spec); len(found) == 0 {
			t.Fatal("no SR found after the attach")
		} else {
			for _, f := range found {
				if f.UUID == uuid && f.InUse == "" {
					t.Fatal("an attached SR is offered again")
				}
			}
		}
		if _, err := ops.CreateSR(c, spec); err == nil {
			t.Fatal("attached twice")
		}
	}
	unplug()
	attach("reattach")
	if err := xenapi.SR.Forget(s, unplug()); err != nil {
		t.Fatal(err)
	}
	attach("introduce")
}

// TestLabSRCommands creates an NFS SR on XCP_LAB_NFS (XCP_ALLOW_WRITE=1) and runs the Storage menu on it: Detach (the
// SR turns broken and offers Repair, Forget and Destroy), Repair, Detach again, then Destroy; an NFS ISO library on
// XCP_LAB_NFS/iso is forgotten.
func TestLabSRCommands(t *testing.T) {
	c := labConnection(t)
	share := os.Getenv("XCP_LAB_NFS")
	if share == "" || os.Getenv("XCP_ALLOW_WRITE") != "1" {
		t.Skip("XCP_LAB_NFS and XCP_ALLOW_WRITE=1 not set")
	}
	find := func(uuid string) Object {
		srs, err := ListObjects(c, "srs")
		if err != nil {
			t.Fatal(err)
		}
		for _, o := range srs {
			if o.UUID == uuid {
				return o
			}
		}
		return Object{}
	}
	uuid, err := ops.CreateSR(c, ops.NewSR{Kind: rules.SRNfs, Name: "xcpng-center-go-test-commands", Path: share})
	if err != nil {
		t.Fatal(err)
	}
	o := find(uuid)
	if o.Icon != "000_Storage_h32bit_16" || strings.Join(o.Allowed, ",") != "set_default,detach" {
		t.Fatalf("new SR %+v", o)
	}
	if err := ops.SRCommand(c, uuid, "destroy"); err == nil {
		t.Fatal("an attached SR was destroyed")
	}
	if err := ops.SRCommand(c, uuid, "detach"); err != nil {
		t.Fatal(err)
	}
	if o = find(uuid); o.Icon != "000_StorageBroken_h32bit_16" || strings.Join(o.Allowed, ",") != "repair,set_default,forget,destroy" {
		t.Fatalf("detached SR %+v", o)
	}
	if err := ops.SRCommand(c, uuid, "repair"); err != nil {
		t.Fatal(err)
	}
	if o = find(uuid); o.Icon != "000_Storage_h32bit_16" {
		t.Fatalf("repaired SR %+v", o)
	}
	if err := ops.SRCommand(c, uuid, "detach"); err != nil {
		t.Fatal(err)
	}
	if err := ops.SRCommand(c, uuid, "destroy"); err != nil {
		t.Fatal(err)
	}
	if o = find(uuid); o.UUID != "" {
		t.Fatalf("destroyed SR still there: %+v", o)
	}
	iso, err := ops.CreateSR(c, ops.NewSR{Kind: rules.SRNfsIso, Name: "xcpng-center-go-test-commands-iso", Path: share + "/iso"})
	if err != nil {
		t.Fatal(err)
	}
	if err := ops.SRCommand(c, iso, "detach"); err != nil {
		t.Fatal(err)
	}
	if err := ops.SRCommand(c, iso, "forget"); err != nil {
		t.Fatal(err)
	}
	if o = find(iso); o.UUID != "" {
		t.Fatalf("forgotten SR still there: %+v", o)
	}
}

// TestLabCopyDeleteVM creates a halted VM with a 1 GB disk (XCP_ALLOW_WRITE=1), fast-clones it and fully copies it to
// the default SR, then deletes the three with their disks (a disk of another VM is refused).
func TestLabCopyDeleteVM(t *testing.T) {
	c := labConnection(t)
	if os.Getenv("XCP_ALLOW_WRITE") != "1" {
		t.Skip("XCP_ALLOW_WRITE=1 not set")
	}
	s := c.Session
	opts, err := ReadNewVMOptions(c)
	if err != nil {
		t.Fatal(err)
	}
	var tpl string
	for _, x := range opts.Templates {
		if x.Name == "Other install media" {
			tpl = x.UUID
		}
	}
	d, err := ReadNewVMDefaults(c, tpl, "")
	if err != nil {
		t.Fatal(err)
	}
	created, err := ops.CreateVM(c, ops.NewVM{Template: tpl, Name: "xcpng-center-go-test-copy", Install: "cd", BootMode: d.BootMode, VCPUs: 1, Cores: 1,
		Memory: 256 << 20, Disks: []ops.NewDisk{{Position: "0", Name: "xcpng-center-go-test-copy 0", Size: 1 << 30, SR: opts.DefaultSR}}})
	if err != nil {
		t.Fatal(err)
	}
	vms := []string{created.UUID}
	defer func() { // whatever the test left
		for _, u := range vms {
			if ref, err := xenapi.VM.GetByUUID(s, u); err == nil {
				del, _ := ReadVMDelete(c, u)
				var vbds []string
				for _, x := range del.Disks {
					vbds = append(vbds, x.UUID)
				}
				if err := ops.DeleteVM(c, u, vbds, nil); err != nil {
					t.Error(err)
					_ = xenapi.VM.Destroy(s, ref)
				}
			}
		}
	}()
	cp, err := ReadVMCopy(c, created.UUID)
	if err != nil || !cp.Clone || !cp.Copy || cp.Name != "Copy of xcpng-center-go-test-copy" || len(cp.SRs) == 0 {
		t.Fatalf("%+v %v", cp, err)
	}
	clone, err := ops.CopyVM(c, created.UUID, "xcpng-center-go-test-copy clone", "fast", false, "")
	if clone != "" {
		vms = append(vms, clone)
	}
	if err != nil {
		t.Fatal(err)
	}
	full, err := ops.CopyVM(c, created.UUID, "xcpng-center-go-test-copy full", "full", true, cp.DefaultSR)
	if full != "" {
		vms = append(vms, full)
	}
	if err != nil {
		t.Fatal(err)
	}
	if ref, _ := xenapi.VM.GetByUUID(s, full); true {
		if desc, _ := xenapi.VM.GetNameDescription(s, ref); desc != "full" {
			t.Fatalf("description %q", desc)
		}
	}
	other, _ := ReadVMDelete(c, clone)
	if len(other.Disks) != 1 || !other.Disks[0].Checked {
		t.Fatalf("clone disks %+v", other)
	}
	if err := ops.DeleteVM(c, full, []string{other.Disks[0].UUID}, nil); err == nil {
		t.Fatal("deleted a disk of another VM")
	}
	for _, u := range vms {
		del, err := ReadVMDelete(c, u)
		if err != nil || len(del.Disks) != 1 {
			t.Fatalf("%s: %+v %v", u, del, err)
		}
		vbd, _ := xenapi.VBD.GetByUUID(s, del.Disks[0].UUID)
		vdi, _ := xenapi.VBD.GetVDI(s, vbd)
		if err := ops.DeleteVM(c, u, []string{del.Disks[0].UUID}, nil); err != nil {
			t.Fatal(u, err)
		}
		if _, err := xenapi.VDI.GetUUID(s, vdi); err == nil {
			t.Fatalf("%s: the disk was not deleted", u)
		}
	}
	vms = nil
}

// TestLabMigrateVM reads the targets of the running XCP_LAB_VM (its own server is "Current server"), live-migrates it
// to the other server and back (XCP_ALLOW_WRITE=1).
func TestLabMigrateVM(t *testing.T) {
	c := labConnection(t)
	name := os.Getenv("XCP_LAB_VM")
	if name == "" || os.Getenv("XCP_ALLOW_WRITE") != "1" {
		t.Skip("XCP_LAB_VM and XCP_ALLOW_WRITE=1 not set")
	}
	s := c.Session
	refs, err := xenapi.VM.GetByNameLabel(s, name)
	if err != nil || len(refs) != 1 {
		t.Fatalf("%s: %v %v", name, refs, err)
	}
	uuid, _ := xenapi.VM.GetUUID(s, refs[0])
	resident := func() string {
		h, _ := xenapi.VM.GetResidentOn(s, refs[0])
		u, _ := xenapi.Host.GetUUID(s, h)
		return u
	}
	from := resident()
	targets, err := ReadTargets(c, uuid)
	if err != nil || len(targets) < 3 || !targets[0].Home {
		t.Fatalf("%+v %v", targets, err)
	}
	var to string
	for _, x := range targets[1:] {
		if x.UUID == from && (x.Enabled || !strings.HasSuffix(x.Text, "(Current server)")) {
			t.Fatalf("own server %+v", x)
		}
		if x.UUID != from && x.Enabled {
			to = x.UUID
		}
	}
	if to == "" {
		t.Fatalf("no server to migrate to: %+v", targets)
	}
	if err := ops.VMOnHost(c, uuid, to, "migrate"); err != nil {
		t.Fatal(err)
	}
	if resident() != to {
		t.Fatal("the VM did not move")
	}
	if err := ops.VMOnHost(c, uuid, from, "migrate"); err != nil {
		t.Fatal(err)
	}
	if err := ops.VMOnHost(c, uuid, from, "start_on"); err == nil {
		t.Fatal("started a running VM")
	}
}

// TestLabISCSI scans XCP_LAB_ISCSI (the lab's iSCSI target, XCP_ALLOW_WRITE=1), formats its first LUN as an SR, checks
// that the LUN is then in use, forgets the SR, finds it again with the LUN probe and reattaches it, then destroys it.
func TestLabISCSI(t *testing.T) {
	c := labConnection(t)
	target := os.Getenv("XCP_LAB_ISCSI")
	if target == "" || os.Getenv("XCP_ALLOW_WRITE") != "1" {
		t.Skip("XCP_LAB_ISCSI and XCP_ALLOW_WRITE=1 not set")
	}
	s := c.Session
	iqns, err := ops.ScanIQNs(c, target, rules.DefaultISCSIPort, "", "")
	if err != nil || len(iqns) == 0 || iqns[0].IQN == "*" {
		t.Fatalf("%+v %v", iqns, err)
	}
	luns, err := ops.ScanLUNs(c, iqns[0], "", "")
	if err != nil || len(luns) == 0 {
		t.Fatalf("%+v %v", luns, err)
	}
	if used, _ := ops.LUNInUse(c, luns[0].ScsiID); used != "" {
		t.Skipf("the lab LUN is used by %s", used)
	}
	spec := ops.NewSR{Kind: rules.SRIscsi, Name: "xcpng-center-go-test-iscsi", AutoDesc: true, IQN: iqns[0], LUN: luns[0]}
	if found, err := ops.ProbeSR(c, spec); err != nil || len(found) != 0 {
		t.Fatalf("a fresh LUN: %+v %v", found, err)
	}
	uuid, err := ops.CreateSR(c, spec)
	if uuid != "" {
		defer func() {
			if ref, err := xenapi.SR.GetByUUID(s, uuid); err == nil {
				pbds, _ := xenapi.SR.GetPBDs(s, ref)
				for _, p := range pbds {
					_ = xenapi.PBD.Unplug(s, p)
				}
				if err := xenapi.SR.Destroy(s, ref); err != nil {
					t.Error(err)
				}
			}
		}()
	}
	if err != nil {
		t.Fatal(err)
	}
	ref, _ := xenapi.SR.GetByUUID(s, uuid)
	if rec, _ := xenapi.SR.GetRecord(s, ref); rec.Type != "lvmoiscsi" || !rec.Shared || !strings.HasPrefix(rec.NameDescription, "iSCSI SR [") {
		t.Fatalf("%+v", rec)
	}
	if used, _ := ops.LUNInUse(c, luns[0].ScsiID); used != spec.Name {
		t.Fatalf("in use: %q", used)
	}
	if _, err := ops.CreateSR(c, spec); err == nil {
		t.Fatal("formatted a LUN in use")
	}
	if err := ops.SRCommand(c, uuid, "detach"); err != nil {
		t.Fatal(err)
	}
	if err := ops.SRCommand(c, uuid, "forget"); err != nil {
		t.Fatal(err)
	}
	found, err := ops.ProbeSR(c, spec)
	if err != nil || len(found) != 1 || found[0].UUID != uuid {
		t.Fatalf("the LUN probe: %+v %v", found, err)
	}
	spec.Existing = uuid
	if _, err := ops.CreateSR(c, spec); err != nil {
		t.Fatal(err)
	}
	if ref, err = xenapi.SR.GetByUUID(s, uuid); err != nil {
		t.Fatal(err)
	}
	for _, p := range func() []xenapi.PBDRef { r, _ := xenapi.SR.GetPBDs(s, ref); return r }() {
		if pbd, _ := xenapi.PBD.GetRecord(s, p); !pbd.CurrentlyAttached {
			t.Fatal("a PBD of the reattached SR is not attached")
		}
	}
}

// TestLabViews puts a tag, a folder and a custom field on XCP_LAB_VM (XCP_ALLOW_WRITE=1), reads the Organization Views
// and the default searches, then removes them again.
func TestLabViews(t *testing.T) {
	c := labConnection(t)
	name := os.Getenv("XCP_LAB_VM")
	if name == "" || os.Getenv("XCP_ALLOW_WRITE") != "1" {
		t.Skip("XCP_LAB_VM and XCP_ALLOW_WRITE=1 not set")
	}
	s := c.Session
	refs, _ := xenapi.VM.GetByNameLabel(s, name)
	if len(refs) != 1 {
		t.Fatal(name)
	}
	vm := refs[0]
	pools, _ := xenapi.Pool.GetAll(s)
	gc, err := xenapi.Pool.GetGuiConfig(s, pools[0])
	if err != nil {
		t.Fatal(err)
	}
	rec, err := xenapi.VM.GetRecord(s, vm)
	if err != nil {
		t.Fatal(err)
	}
	oldFields, hadFields := gc["XenCenter.CustomFields"]
	oldFolder, hadFolder := rec.OtherConfig["folder"]
	oldValue, hadValue := rec.OtherConfig["XenCenter.CustomFields.go-test-field"]
	hadTag := slices.Contains(rec.Tags, "xcpng-center-go-test-tag")
	restored := func(err error) {
		if err != nil {
			t.Errorf("the lab's metadata was not restored: %v", err)
		}
	}
	defer func() { // only what the test added goes; what was there comes back
		if !hadTag {
			restored(xenapi.VM.RemoveTags(s, vm, "xcpng-center-go-test-tag"))
		}
		restored(xenapi.VM.RemoveFromOtherConfig(s, vm, "folder"))
		if hadFolder {
			restored(xenapi.VM.AddToOtherConfig(s, vm, "folder", oldFolder))
		}
		restored(xenapi.VM.RemoveFromOtherConfig(s, vm, "XenCenter.CustomFields.go-test-field"))
		if hadValue {
			restored(xenapi.VM.AddToOtherConfig(s, vm, "XenCenter.CustomFields.go-test-field", oldValue))
		}
		restored(xenapi.Pool.RemoveFromGuiConfig(s, pools[0], "XenCenter.CustomFields"))
		if hadFields {
			restored(xenapi.Pool.AddToGuiConfig(s, pools[0], "XenCenter.CustomFields", oldFields))
		}
	}()
	if err := xenapi.VM.AddTags(s, vm, "xcpng-center-go-test-tag"); err != nil {
		t.Fatal(err)
	}
	_ = xenapi.VM.RemoveFromOtherConfig(s, vm, "folder")
	if err := xenapi.VM.AddToOtherConfig(s, vm, "folder", "/go-test/inner"); err != nil {
		t.Fatal(err)
	}
	_ = xenapi.Pool.RemoveFromGuiConfig(s, pools[0], "XenCenter.CustomFields")
	fields := `<CustomFieldDefinitions><CustomFieldDefinition name="go-test-field" type="String" defaultValue="" /></CustomFieldDefinitions>`
	if hadFields { // the lab's own fields stay: the test's goes in their list
		fields = strings.Replace(oldFields, "</CustomFieldDefinitions>", `<CustomFieldDefinition name="go-test-field" type="String" defaultValue="" /></CustomFieldDefinitions>`, 1)
	}
	if err := xenapi.Pool.AddToGuiConfig(s, pools[0], "XenCenter.CustomFields", fields); err != nil {
		t.Fatal(err)
	}
	_ = xenapi.VM.RemoveFromOtherConfig(s, vm, "XenCenter.CustomFields.go-test-field")
	if err := xenapi.VM.AddToOtherConfig(s, vm, "XenCenter.CustomFields.go-test-field", "value-1"); err != nil {
		t.Fatal(err)
	}
	find := func(nodes []*ViewNode, path ...string) *ViewNode {
		for _, p := range path {
			var next *ViewNode
			for _, n := range nodes {
				if n.Name == p {
					next = n
				}
			}
			if next == nil {
				return nil
			}
			nodes = next.Children
			if p == path[len(path)-1] {
				return next
			}
		}
		return nil
	}
	for view, path := range map[string][]string{
		"tags": {"xcpng-center-go-test-tag", name}, "folders": {"go-test", "inner", name},
		"custom_fields": {"go-test-field", "value-1", name}, "vms_by_power_state": {"Running", name}, "vms_by_os": {name},
	} {
		got, err := ReadView(c, "lab", view)
		if err != nil {
			t.Fatal(view, err)
		}
		if view == "vms_by_os" { // whatever the guest reports: the VM is somewhere one level down
			found := false
			for _, g := range got {
				found = found || find(g.Children, name) != nil
			}
			if !found {
				t.Errorf("%s: %s not found", view, name)
			}
			continue
		}
		if n := find(got, path...); n == nil || n.Object == nil || n.Server != "lab" {
			t.Errorf("%s: %v not found", view, path)
		}
	}
	if searches := DefaultSearches(); len(searches) != 7 {
		t.Errorf("%d default searches", len(searches))
	}
}

// TestLabPause pauses the running XCP_LAB_VM and unpauses it (XCP_ALLOW_WRITE=1).
func TestLabPause(t *testing.T) {
	c := labConnection(t)
	name := os.Getenv("XCP_LAB_VM")
	if name == "" || os.Getenv("XCP_ALLOW_WRITE") != "1" {
		t.Skip("XCP_LAB_VM and XCP_ALLOW_WRITE=1 not set")
	}
	refs, err := xenapi.VM.GetByNameLabel(c.Session, name)
	if err != nil || len(refs) != 1 {
		t.Fatalf("%s: %d VMs, %v", name, len(refs), err)
	}
	uuid, err := xenapi.VM.GetUUID(c.Session, refs[0])
	if err != nil {
		t.Fatal(err)
	}
	state := func() xenapi.VMPowerState { p, _ := xenapi.VM.GetPowerState(c.Session, refs[0]); return p }
	t.Cleanup(func() { // never leave the lab's VM paused
		if state() == xenapi.VMPowerStatePaused {
			_ = xenapi.VM.Unpause(c.Session, refs[0])
		}
	})
	if err := ops.VMPower(c, uuid, rules.VMPause); err != nil || state() != xenapi.VMPowerStatePaused {
		t.Fatal(err, state())
	}
	if err := ops.VMPower(c, uuid, rules.VMUnpause); err != nil || state() != xenapi.VMPowerStateRunning {
		t.Fatal(err, state())
	}
}

// TestLabBootAndAdvanced changes the running XCP_LAB_VM's autostart, boot order and shadow memory multiplier (live)
// through the Properties pages, reads them back and puts them back (XCP_ALLOW_WRITE=1).
func TestLabBootAndAdvanced(t *testing.T) {
	c := labConnection(t)
	name := os.Getenv("XCP_LAB_VM")
	if name == "" || os.Getenv("XCP_ALLOW_WRITE") != "1" {
		t.Skip("XCP_LAB_VM and XCP_ALLOW_WRITE=1 not set")
	}
	refs, err := xenapi.VM.GetByNameLabel(c.Session, name)
	if err != nil || len(refs) != 1 {
		t.Fatalf("%s: %v", name, err)
	}
	uuid, _ := xenapi.VM.GetUUID(c.Session, refs[0])
	before, err := ReadVMProperties(c, uuid)
	if err != nil || !before.HVM || before.BootOrder == "" || len(before.BootDevices) != 3 {
		t.Fatalf("%+v %v", before, err)
	}
	t.Cleanup(func() {
		auto, order, shadow := before.AutoPowerOn, before.BootOrder, before.ShadowMultiplier
		if err := ops.SaveVMProperties(c, uuid, ops.VMChanges{AutoPowerOn: &auto, BootOrder: &order, ShadowMultiplier: &shadow}); err != nil {
			t.Errorf("not restored: %v", err)
		}
	})
	auto, order, shadow := !before.AutoPowerOn, "NDC", before.ShadowMultiplier+0.5
	if err := ops.SaveVMProperties(c, uuid, ops.VMChanges{AutoPowerOn: &auto, BootOrder: &order, ShadowMultiplier: &shadow}); err != nil {
		t.Fatal(err)
	}
	after, _ := ReadVMProperties(c, uuid)
	// a live change: Xen rounds the shadow memory to whole MB, so the multiplier read back is close, not equal
	if after.AutoPowerOn != auto || after.BootOrder != "NDC" || math.Abs(after.ShadowMultiplier-shadow) > 0.2 {
		t.Fatalf("%v %q %v", after.AutoPowerOn, after.BootOrder, after.ShadowMultiplier)
	}
	bad := "CC"
	if err := ops.SaveVMProperties(c, uuid, ops.VMChanges{BootOrder: &bad}); err == nil {
		t.Fatal("a bad boot order was saved")
	}
}

// TestLabSavedSearch saves a search in the pool's gui_config as the C# app does (running VMs named after XCP_LAB_VM,
// grouped by power state), lists and runs it, then removes it (XCP_ALLOW_WRITE=1).
func TestLabSavedSearch(t *testing.T) {
	c := labConnection(t)
	name := os.Getenv("XCP_LAB_VM")
	if name == "" || os.Getenv("XCP_ALLOW_WRITE") != "1" {
		t.Skip("XCP_LAB_VM and XCP_ALLOW_WRITE=1 not set")
	}
	s := c.Session
	pools, _ := xenapi.Pool.GetAll(s)
	key := searchPrefix + "go-test-search"
	text := `<Search uuid="go-test-search" name="Go test search" major_version="2" minor_version="0"><Query><QueryScope><VM /></QueryScope>` +
		`<GroupQuery type="And"><EnumPropertyQuery property="power_state" equals="yes" query="Running" />` +
		`<StringPropertyQuery property="label" casesensitive="no" query="` + name + `" type="exactmatch" /></GroupQuery></Query>` +
		`<PropertyGrouping property="power_state" /></Search>`
	if err := xenapi.Pool.AddToGuiConfig(s, pools[0], key, text); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := xenapi.Pool.RemoveFromGuiConfig(s, pools[0], key); err != nil {
			t.Errorf("not removed: %v", err)
		}
	})
	list, err := PoolSearches(c)
	if err != nil || len(list) == 0 {
		t.Fatalf("%+v %v", list, err)
	}
	for _, v := range list {
		if v.ID == "saved:go-test-search" && (!v.Supported || v.Name != "Go test search") {
			t.Fatalf("%+v", v)
		}
	}
	saved, err := PoolSearch(c, "go-test-search")
	if err != nil || saved != text {
		t.Fatalf("%q %v", saved, err)
	}
	nodes, err := RunSearch(c, "lab", saved)
	if err != nil || len(nodes) != 1 || nodes[0].Name != "Running" || len(nodes[0].Children) != 1 || nodes[0].Children[0].Name != name {
		t.Fatalf("%+v %v", nodes, err)
	}
}

// TestLabSearchEditor builds a search with the editor (VMs named XCP_LAB_VM, by power state), saves it to the lab
// pool, reads it back through the saved-search list, runs it, exports and imports it, and deletes it.
func TestLabSearchEditor(t *testing.T) {
	c := labConnection(t)
	name := os.Getenv("XCP_LAB_VM")
	if name == "" || os.Getenv("XCP_ALLOW_WRITE") != "1" {
		t.Skip("XCP_LAB_VM and XCP_ALLOW_WRITE=1 not set")
	}
	const title = "xcpng-center-go-test-search"
	text, err := BuildSearch(SearchEdit{Name: title, Types: []string{"VM"}, Text: name, Grouping: 5})
	if err != nil {
		t.Fatal(err)
	}
	uuid := SearchUUID(text)
	if err := ops.SaveSearch(c, uuid, text); err != nil {
		t.Fatal(err)
	}
	deleted := false
	t.Cleanup(func() {
		if !deleted {
			if err := ops.DeleteSearch(c, uuid); err != nil {
				t.Errorf("not removed: %v", err)
			}
		}
	})
	list, err := PoolSearches(c)
	if err != nil || !slices.ContainsFunc(list, func(v ViewInfo) bool { return v.ID == "saved:"+uuid && v.Name == title && v.Supported }) {
		t.Fatalf("%+v %v", list, err)
	}
	saved, err := PoolSearch(c, uuid)
	if err != nil || saved != text {
		t.Fatalf("%q %v", saved, err)
	}
	nodes, err := RunSearch(c, "lab", saved)
	found := false
	for _, g := range nodes { // a power state group holds the VM (its name contains the filter text)
		for _, n := range g.Children {
			found = found || n.Name == name
		}
	}
	if err != nil || !found {
		t.Fatalf("%+v %v", nodes, err)
	}
	// editing the saved search keeps its uuid: saving replaces it
	again, err := BuildSearch(SearchEdit{Name: title, Types: []string{"VM", "Server"}, Grouping: 5, Edited: saved, Saved: true})
	if err != nil || SearchUUID(again) != uuid {
		t.Fatalf("%v", err)
	}
	if err := ops.SaveSearch(c, uuid, again); err != nil {
		t.Fatal(err)
	}
	if got, _ := PoolSearch(c, uuid); got != again {
		t.Fatalf("not replaced: %q", got)
	}
	file, err := SearchFile(again)
	if err != nil {
		t.Fatal(err)
	}
	imported, err := ImportSearches([]byte(file))
	if err != nil || len(imported) != 1 || SearchUUID(imported[0]) == uuid || SearchName(imported[0]) != title {
		t.Fatalf("%v %v", imported, err)
	}
	if err := ops.DeleteSearch(c, uuid); err != nil {
		t.Fatal(err)
	}
	deleted = true
	if got, _ := PoolSearch(c, uuid); got != "" {
		t.Fatalf("not deleted: %q", got)
	}
}

// TestLabMemory reads the Memory Settings of XCP_LAB_VM (running, no guest tools: a fixed memory), refuses a value
// over the limits, raises the memory by 64 MB (the VM restarts on its server) and restores it (XCP_ALLOW_WRITE=1).
func TestLabMemory(t *testing.T) {
	c := labConnection(t)
	name := os.Getenv("XCP_LAB_VM")
	if name == "" || os.Getenv("XCP_ALLOW_WRITE") != "1" {
		t.Skip("XCP_LAB_VM and XCP_ALLOW_WRITE=1 not set")
	}
	refs, err := xenapi.VM.GetByNameLabel(c.Session, name)
	if err != nil || len(refs) != 1 {
		t.Fatalf("%s: %v", name, err)
	}
	uuid, _ := xenapi.VM.GetUUID(c.Session, refs[0])
	host, _ := xenapi.VM.GetResidentOn(c.Session, refs[0])
	e, err := ReadMemoryEdit(c, uuid)
	if err != nil || !e.CanEdit || e.Advanced || e.Dynamic || e.Unavailable == "" || e.Rubric == "" || e.Max < e.StaticMax {
		t.Fatalf("%+v %v", e, err)
	}
	if _, err := MemoryQuestion(c, uuid, 0, 0, e.Max+1); err == nil {
		t.Fatal("more than the maximum is accepted")
	}
	more := e.StaticMax + 64<<20
	if q, err := MemoryQuestion(c, uuid, more, more, more); err != nil || q == "" {
		t.Fatalf("no question before the restart: %q %v", q, err)
	}
	t.Cleanup(func() {
		if err := ChangeMemory(c, uuid, e.DynamicMin, e.DynamicMax, e.StaticMax); err != nil {
			t.Errorf("not restored: %v", err)
		}
	})
	if err := ChangeMemory(c, uuid, more, more, more); err != nil {
		t.Fatal(err)
	}
	after, _ := xenapi.VM.GetRecord(c.Session, refs[0])
	if int64(after.MemoryStaticMax) != more || int64(after.MemoryDynamicMin) != more || after.PowerState != xenapi.VMPowerStateRunning || after.ResidentOn != host {
		t.Fatalf("%d %d %s %s", after.MemoryStaticMax, after.MemoryDynamicMin, after.PowerState, after.ResidentOn)
	}
}

// TestLabVIF adds an interface to XCP_LAB_VM (running: plugged at once), changes its MAC and QoS limit (the VIF is
// recreated on its device), deactivates and activates it, and removes it (XCP_ALLOW_WRITE=1).
func TestLabVIF(t *testing.T) {
	c := labConnection(t)
	name := os.Getenv("XCP_LAB_VM")
	if name == "" || os.Getenv("XCP_ALLOW_WRITE") != "1" {
		t.Skip("XCP_LAB_VM and XCP_ALLOW_WRITE=1 not set")
	}
	refs, err := xenapi.VM.GetByNameLabel(c.Session, name)
	if err != nil || len(refs) != 1 {
		t.Fatalf("%s: %v", name, err)
	}
	uuid, _ := xenapi.VM.GetUUID(c.Session, refs[0])
	before, err := ListVMNetworking(c, uuid)
	if err != nil || !before.CanAdd {
		t.Fatalf("%+v %v", before, err)
	}
	nets, err := VIFNetworks(c, uuid)
	if err != nil || len(nets) == 0 {
		t.Fatalf("%v %v", nets, err)
	}
	if _, err := ops.CreateVIF(c, uuid, ops.VIFSettings{Network: nets[0].UUID, MAC: "zz"}); err == nil {
		t.Fatal("a bad MAC is accepted")
	}
	added := func() *VIF {
		n, _ := ListVMNetworking(c, uuid)
		for _, v := range n.VIFs {
			if !slices.ContainsFunc(before.VIFs, func(b VIF) bool { return b.Device == v.Device }) {
				return &v
			}
		}
		return nil
	}
	if _, err := ops.CreateVIF(c, uuid, ops.VIFSettings{Network: nets[0].UUID}); err != nil {
		t.Fatal(err)
	}
	v := added()
	if v == nil || v.MAC == "" {
		t.Fatalf("not added: %+v", v)
	}
	t.Cleanup(func() {
		if v := added(); v != nil {
			if err := ops.DeleteVIF(c, v.UUID); err != nil {
				t.Errorf("not removed: %v", err)
			}
		}
	})
	if dup, err := DuplicateMAC(c, "", v.MAC); err != nil || dup == "" {
		t.Fatalf("the MAC of %s is not a duplicate: %q %v", name, dup, err)
	}
	mac := "02:00:5e:10:20:30"
	if _, err := ops.UpdateVIF(c, v.UUID, ops.VIFSettings{Network: nets[0].UUID, MAC: mac, QoS: true, Limit: "1000"}); err != nil {
		t.Fatal(err)
	}
	u := added()
	if u == nil || u.Device != v.Device || u.MAC != mac || !u.QoS || u.Limit != "1000" {
		t.Fatalf("not updated: %+v", u)
	}
	if u.Attached && u.CanToggle {
		if err := ops.SetVIFActive(c, u.UUID, false); err != nil {
			t.Fatal(err)
		}
		if w := added(); w.Attached {
			t.Fatal("still attached")
		}
		if err := ops.SetVIFActive(c, u.UUID, true); err != nil {
			t.Fatal(err)
		}
	}
	if err := ops.DeleteVIF(c, u.UUID); err != nil || added() != nil {
		t.Fatalf("not removed: %v", err)
	}
}

// TestLabStartOptions changes the start order and delay of XCP_LAB_VM and restores them (XCP_ALLOW_WRITE=1).
func TestLabStartOptions(t *testing.T) {
	c := labConnection(t)
	name := os.Getenv("XCP_LAB_VM")
	if name == "" || os.Getenv("XCP_ALLOW_WRITE") != "1" {
		t.Skip("XCP_LAB_VM and XCP_ALLOW_WRITE=1 not set")
	}
	refs, err := xenapi.VM.GetByNameLabel(c.Session, name)
	if err != nil || len(refs) != 1 {
		t.Fatalf("%s: %v", name, err)
	}
	uuid, _ := xenapi.VM.GetUUID(c.Session, refs[0])
	before, err := ReadVMProperties(c, uuid)
	if err != nil || before.HAEnabled {
		t.Fatalf("%+v %v", before.HAEnabled, err)
	}
	t.Cleanup(func() {
		if err := ops.SaveVMProperties(c, uuid, ops.VMChanges{Order: &before.Order, StartDelay: &before.StartDelay}); err != nil {
			t.Errorf("not restored: %v", err)
		}
	})
	order, delay := before.Order+3, before.StartDelay+7
	if err := ops.SaveVMProperties(c, uuid, ops.VMChanges{Order: &order, StartDelay: &delay}); err != nil {
		t.Fatal(err)
	}
	if after, _ := ReadVMProperties(c, uuid); after.Order != order || after.StartDelay != delay {
		t.Fatalf("%d %d", after.Order, after.StartDelay)
	}
	bad := -1
	if err := ops.SaveVMProperties(c, uuid, ops.VMChanges{Order: &bad}); err == nil {
		t.Fatal("a negative order is saved")
	}
}

// TestLabCustomFields adds a text and a date field to the pool, sets them on XCP_LAB_VM, reads them back and deletes the
// fields with their values (XCP_ALLOW_WRITE=1).
func TestLabCustomFields(t *testing.T) {
	c := labConnection(t)
	name := os.Getenv("XCP_LAB_VM")
	if name == "" || os.Getenv("XCP_ALLOW_WRITE") != "1" {
		t.Skip("XCP_LAB_VM and XCP_ALLOW_WRITE=1 not set")
	}
	refs, err := xenapi.VM.GetByNameLabel(c.Session, name)
	if err != nil || len(refs) != 1 {
		t.Fatalf("%s: %v", name, err)
	}
	uuid, _ := xenapi.VM.GetUUID(c.Session, refs[0])
	text, date := "xcpng-center-go-test owner", "xcpng-center-go-test expiry"
	t.Cleanup(func() {
		for _, f := range []string{text, date} {
			if err := ops.RemoveCustomField(c, f); err != nil {
				t.Errorf("%s not removed: %v", f, err)
			}
		}
		if other, _ := xenapi.VM.GetOtherConfig(c.Session, refs[0]); other[rules.CustomFieldsKey+"."+text] != "" {
			t.Error("the value stays on the VM")
		}
	})
	if err := ops.AddCustomField(c, text, "String"); err != nil {
		t.Fatal(err)
	}
	if err := ops.AddCustomField(c, date, "Date"); err != nil {
		t.Fatal(err)
	}
	if err := ops.AddCustomField(c, text, "String"); err == nil {
		t.Fatal("the same field twice")
	}
	if err := ops.SaveVMCustomFields(c, uuid, map[string]string{text: "lab", date: "2026-03-04T05:06:07Z"}); err != nil {
		t.Fatal(err)
	}
	other, _ := xenapi.VM.GetOtherConfig(c.Session, refs[0])
	if other[rules.CustomFieldsKey+"."+date] != "2026-03-04 05:06:07Z" {
		t.Fatalf("the date is not stored as the C# app does: %q", other[rules.CustomFieldsKey+"."+date])
	}
	got := map[string]string{}
	fields, err := VMCustomFields(c, uuid)
	for _, f := range fields {
		got[f.Name] = f.Value
	}
	if err != nil || got[text] != "lab" || got[date] != "2026-03-04T05:06:07Z" {
		t.Fatalf("%v %v", got, err)
	}
	if err := ops.SaveVMCustomFields(c, uuid, map[string]string{text: ""}); err != nil {
		t.Fatal(err)
	}
	if other, _ := xenapi.VM.GetOtherConfig(c.Session, refs[0]); other[rules.CustomFieldsKey+"."+text] != "" {
		t.Fatal("an empty value is not removed")
	}
}

// TestLabHostTabs reads the Memory, Storage, Networking and NICs tabs of the coordinator and the pool, and the
// Storage tab of the default SR.
func TestLabHostTabs(t *testing.T) {
	c := labConnection(t)
	s := c.Session
	pools, _ := xenapi.Pool.GetAllRecords(s)
	var master xenapi.HostRef
	var poolUUID string
	var sr xenapi.SRRef
	for _, p := range pools {
		master, poolUUID, sr = p.Master, p.UUID, p.DefaultSR
	}
	hostUUID, _ := xenapi.Host.GetUUID(s, master)
	srUUID, _ := xenapi.SR.GetUUID(s, sr)
	for _, c2 := range []struct{ kind, uuid, tab string }{
		{"hosts", hostUUID, "memory"}, {"hosts", hostUUID, "storage"}, {"hosts", hostUUID, "networking"}, {"hosts", hostUUID, "nics"},
		{"pools", poolUUID, "memory"}, {"pools", poolUUID, "storage"}, {"pools", poolUUID, "networking"}, {"srs", srUUID, "storage"},
	} {
		tab, err := HostTab(c, c2.kind, c2.uuid, c2.tab)
		if err != nil || len(tab.Rows) == 0 {
			t.Fatalf("%s %s: %+v %v", c2.kind, c2.tab, tab, err)
		}
		for _, r := range tab.Rows {
			if len(r) != len(tab.Headers) {
				t.Fatalf("%s %s: row %v for %v", c2.kind, c2.tab, r, tab.Headers)
			}
		}
		t.Logf("%s %s: %v", c2.kind, c2.tab, tab.Rows)
	}
	if _, err := HostTab(c, "pools", poolUUID, "nics"); err == nil {
		t.Error("a pool has no NICs tab")
	}
}

// TestLabExportImport makes a small halted VM (a 16 MB disk on the default SR), exports it to an XVA with verification,
// imports it back (named "<name> (1)"), and destroys both VMs and their disks (XCP_ALLOW_WRITE=1).
func TestLabExportImport(t *testing.T) {
	c := labConnection(t)
	if os.Getenv("XCP_ALLOW_WRITE") != "1" {
		t.Skip("XCP_ALLOW_WRITE=1 not set")
	}
	s := c.Session
	pools, _ := xenapi.Pool.GetAllRecords(s)
	var sr xenapi.SRRef
	for _, p := range pools {
		sr = p.DefaultSR
	}
	srUUID, _ := xenapi.SR.GetUUID(s, sr)
	tmpls, err := xenapi.VM.GetByNameLabel(s, "Other install media")
	if err != nil || len(tmpls) == 0 {
		t.Fatalf("no template: %v", err)
	}
	name := "xcpng-center-go-test-export"
	vm, err := xenapi.VM.Clone(s, tmpls[0], name)
	if err != nil {
		t.Fatal(err)
	}
	destroy := func(vm xenapi.VMRef) {
		vbds, _ := xenapi.VM.GetVBDs(s, vm)
		for _, b := range vbds {
			if vdi, err := xenapi.VBD.GetVDI(s, b); err == nil && vdi != "OpaqueRef:NULL" {
				if typ, _ := xenapi.VBD.GetType(s, b); typ == xenapi.VbdTypeDisk {
					_ = xenapi.VDI.Destroy(s, vdi)
				}
			}
		}
		if err := xenapi.VM.Destroy(s, vm); err != nil {
			t.Errorf("not destroyed: %v", err)
		}
	}
	t.Cleanup(func() { destroy(vm) })
	if err := xenapi.VM.SetIsATemplate(s, vm, false); err != nil {
		t.Fatal(err)
	}
	vdi, err := xenapi.VDI.Create(s, xenapi.VDIRecord{NameLabel: name, SR: sr, VirtualSize: 16 << 20, Type: xenapi.VdiTypeUser, OtherConfig: map[string]string{}, SmConfig: map[string]string{}})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := xenapi.VBD.Create(s, xenapi.VBDRecord{VM: vm, VDI: vdi, Userdevice: "0", Bootable: true, Mode: xenapi.VbdModeRW, Type: xenapi.VbdTypeDisk,
		OtherConfig: map[string]string{}, QosAlgorithmParams: map[string]string{}}); err != nil {
		t.Fatal(err)
	}
	uuid, _ := xenapi.VM.GetUUID(s, vm)
	file := filepath.Join(t.TempDir(), name+".xva")
	if err := ops.ExportVM(c, labPolicy(t), uuid, file, true); err != nil {
		t.Fatal(err)
	}
	if st, err := os.Stat(file); err != nil || st.Size() < 1<<20 {
		t.Fatalf("no export: %v", err)
	}
	if left, _ := os.ReadDir(filepath.Dir(file)); len(left) != 1 {
		t.Fatalf("the temporary file is left: %v", left)
	}
	if ops.ImportFileError(file+".txt") == "" || ops.ImportFileError(file) != "" {
		t.Fatal("file checks")
	}
	imported, err := ops.ImportVM(c, labPolicy(t), file, "", srUUID, false)
	if err != nil {
		t.Fatal(err)
	}
	ref, err := xenapi.VM.GetByUUID(s, imported)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { destroy(ref) })
	if got, _ := xenapi.VM.GetNameLabel(s, ref); got != name+" (1)" {
		t.Fatalf("the imported VM is named %q", got)
	}
}

// TestLabServerAndTemplateCommands converts a throwaway VM to a template (then destroys it), restarts the toolstack of
// the pool member and checks that Remove From Pool refuses the coordinator and a member running a VM (XCP_ALLOW_WRITE=1).
func TestLabServerAndTemplateCommands(t *testing.T) {
	c := labConnection(t)
	if os.Getenv("XCP_ALLOW_WRITE") != "1" {
		t.Skip("XCP_ALLOW_WRITE=1 not set")
	}
	s := c.Session
	opts, err := ReadNewVMOptions(c)
	if err != nil {
		t.Fatal(err)
	}
	var tpl string
	for _, x := range opts.Templates {
		if x.Name == "Other install media" {
			tpl = x.UUID
		}
	}
	d, err := ReadNewVMDefaults(c, tpl, "")
	if err != nil {
		t.Fatal(err)
	}
	created, err := ops.CreateVM(c, ops.NewVM{Template: tpl, Name: "xcpng-center-go-test-template", Install: "cd", BootMode: d.BootMode, VCPUs: 1, Cores: 1,
		Memory: 256 << 20, Disks: []ops.NewDisk{{Position: "0", Name: "xcpng-center-go-test-template 0", Size: 64 << 20, SR: opts.DefaultSR}}})
	if err != nil {
		t.Fatal(err)
	}
	uuid := created.UUID
	vm, _ := xenapi.VM.GetByUUID(s, uuid)
	t.Cleanup(func() { // whatever the test left: a template's disks go with it
		if _, err := xenapi.VM.GetByUUID(s, uuid); err != nil {
			return
		}
		vbds, _ := xenapi.VM.GetVBDs(s, vm)
		for _, b := range vbds {
			if typ, _ := xenapi.VBD.GetType(s, b); typ == xenapi.VbdTypeDisk {
				if vdi, err := xenapi.VBD.GetVDI(s, b); err == nil {
					_ = xenapi.VDI.Destroy(s, vdi)
				}
			}
		}
		if err := xenapi.VM.Destroy(s, vm); err != nil {
			t.Errorf("not destroyed: %v", err)
		}
	})
	if err := ops.ConvertToTemplate(c, uuid); err != nil {
		ops, _ := xenapi.VM.GetAllowedOperations(s, vm)
		t.Fatal(err, ops)
	}
	rec, _ := xenapi.VM.GetRecord(s, vm)
	if !rec.IsATemplate || rec.OtherConfig["instant"] != "true" {
		t.Fatalf("not a template: %v %v", rec.IsATemplate, rec.OtherConfig["instant"])
	}
	if err := ops.ConvertToTemplate(c, uuid); err == nil {
		t.Fatal("a template converted again")
	}
	if objs, _ := ListObjects(c, "templates"); !slices.ContainsFunc(objs, func(o Object) bool { return o.UUID == uuid && o.UserTemplate }) {
		t.Fatal("the template is not a user template")
	}
	// Templates > Copy... (a fast clone of a user template is a template)
	copied, err := ops.CopyVM(c, uuid, "xcpng-center-go-test-template copy", "", false, "")
	if err != nil {
		t.Fatal(err)
	}
	cref, _ := xenapi.VM.GetByUUID(s, copied)
	if tpl, _ := xenapi.VM.GetIsATemplate(s, cref); !tpl {
		t.Error("the copy of a template is not a template")
	}
	if cdel, err := ReadVMDelete(c, copied); err == nil {
		var cvbds []string
		for _, x := range cdel.Disks {
			cvbds = append(cvbds, x.UUID)
		}
		if err := ops.DeleteVM(c, copied, cvbds, nil); err != nil {
			t.Fatal(err)
		}
	}
	// Delete Template..., with its disk
	del, err := ReadVMDelete(c, uuid)
	if err != nil {
		t.Fatal(err)
	}
	var vbds []string
	for _, x := range del.Disks {
		vbds = append(vbds, x.UUID)
	}
	if err := ops.DeleteVM(c, uuid, vbds, nil); err != nil {
		t.Fatal(err)
	}
	if _, err := xenapi.VM.GetByUUID(s, uuid); err == nil {
		t.Fatal("the template is still there")
	}
	pools, _ := xenapi.Pool.GetAllRecords(s)
	var master xenapi.HostRef
	for _, p := range pools {
		master = p.Master
	}
	hosts, _ := xenapi.Host.GetAllRecords(s)
	for r, h := range hosts {
		_, _, _, err := ops.RemoveFromPoolQuestion(c, h.UUID)
		if err == nil {
			t.Fatalf("%s may be removed (resident VMs %d)", h.NameLabel, len(h.ResidentVMs))
		}
		if r == master {
			continue
		}
		if title, text, err := ops.RestartToolstackQuestion(c, h.UUID); err != nil || title == "" || !strings.Contains(text, h.NameLabel) {
			t.Fatalf("%q %q %v", title, text, err)
		}
		if err := ops.RestartToolstack(c, h.UUID); err != nil {
			t.Fatal(err)
		}
	}
}

// TestLabOverview reads the live columns of the Search overview: every server and the running VM have figures.
func TestLabOverview(t *testing.T) {
	c := labConnection(t)
	rows, err := Overview(c, labPolicy(t))
	if err != nil {
		t.Fatal(err)
	}
	hosts, _ := xenapi.Host.GetAllRecords(c.Session)
	for _, h := range hosts {
		r := rows[h.UUID]
		// "-" for the CPU while a restarted toolstack has no samples yet, as the C# app
		if !strings.Contains(r.CPU, "% of") && r.CPU != "-" || !strings.Contains(r.Memory, " of ") || r.Uptime == "-" || r.Address == "-" {
			t.Fatalf("%s: %+v", h.NameLabel, r)
		}
		t.Logf("%s: %+v", h.NameLabel, r)
	}
	if name := os.Getenv("XCP_LAB_VM"); name != "" {
		refs, _ := xenapi.VM.GetByNameLabel(c.Session, name)
		uuid, _ := xenapi.VM.GetUUID(c.Session, refs[0])
		r := rows[uuid]
		if !strings.Contains(r.CPU, "% of") || r.Uptime == "-" || r.Disks == "-" {
			t.Fatalf("%s: %+v", name, r)
		}
		t.Logf("%s: %+v", name, r)
	}
}

// TestLabPoolCommands rotates the pool secret and checks that a pool of two servers cannot be made into a standalone
// server (XCP_ALLOW_WRITE=1).
func TestLabPoolCommands(t *testing.T) {
	c := labConnection(t)
	if os.Getenv("XCP_ALLOW_WRITE") != "1" {
		t.Skip("XCP_ALLOW_WRITE=1 not set")
	}
	if why := ops.MakeStandaloneBlocker(c); why == "" {
		t.Fatal("a pool of two servers may be made standalone")
	}
	if err := ops.MakeStandalone(c); err == nil {
		t.Fatal("made standalone")
	}
	if why := ops.RotateSecretBlocker(c); why != "" {
		t.Fatal(why)
	}
	if err := ops.RotateSecret(c); err != nil {
		t.Fatal(err)
	}
	if _, err := xenapi.Pool.GetAllRecords(c.Session); err != nil { // the session goes on
		t.Fatal(err)
	}
}

// TestLabNetworks bonds NIC 1 and NIC 2 of every lab server (never NIC 0, the management interface) and deletes the
// bond, creates a VLAN network on NIC 1 and a private network, changes and removes them, and rescans the NICs; the
// lab is left with the networks and bonds it had (writes only with XCP_ALLOW_WRITE=1).
func TestLabNetworks(t *testing.T) {
	c := labConnection(t)
	if os.Getenv("XCP_ALLOW_WRITE") != "1" {
		t.Skip("XCP_ALLOW_WRITE=1 not set")
	}
	s := c.Session
	pools, _ := xenapi.Pool.GetAllRecords(s)
	var master xenapi.HostRef
	var poolUUID string
	for _, p := range pools {
		master, poolUUID = p.Master, p.UUID
	}
	hostUUID, _ := xenapi.Host.GetUUID(s, master)
	hosts, _ := xenapi.Host.GetAll(s)
	count := func() (nets, bonds int) {
		n, _ := xenapi.Network.GetAll(s)
		b, _ := xenapi.Bond.GetAll(s)
		return len(n), len(b)
	}
	nets0, bonds0 := count()
	byName := func(name string) (string, xenapi.NetworkRecord) {
		all, _ := xenapi.Network.GetAllRecords(s)
		for _, n := range all {
			if n.NameLabel == name {
				return n.UUID, n
			}
		}
		return "", xenapi.NetworkRecord{}
	}
	for _, n := range []string{"Bond 1+2", "xcpng-center-go-test-vlan", "xcpng-center-go-test-vlan2", "xcpng-center-go-test-private", "xcpng-center-go-test-private2"} {
		if u, _ := byName(n); u != "" {
			t.Fatalf("the lab already has a network %q", n)
		}
	}
	goHome := func() {}
	made := map[string]bool{} // the networks this test made (true: a bond's)
	t.Cleanup(func() {        // whatever a failed step left
		for u, bond := range made {
			if _, err := xenapi.Network.GetByUUID(s, u); err != nil {
				continue
			}
			if bond {
				_ = ops.DestroyBond(c, u)
				goHome()
			} else {
				_ = ops.RemoveNetwork(c, u)
			}
		}
		if n, b := count(); n != nets0 || b != bonds0 {
			t.Errorf("the lab has %d networks and %d bonds, it had %d and %d", n, b, nets0, bonds0)
		}
	})

	// the bond: NIC 1 + NIC 2 of the coordinator, on every server
	opts, err := NewNetworkOptions(c, "hosts", hostUUID)
	if err != nil {
		t.Fatal(err)
	}
	// the VIFs on the member networks move to the bond, and to the primary member's network (NIC 1's) after it: only
	// the throwaway VM's may be there, and they are moved back home
	var members []string
	home := map[xenapi.VIFRef]xenapi.NetworkRef{}
	for _, n := range opts.Bond.NICs {
		if n.InBond {
			t.Fatalf("%s is already bonded", n.Name)
		}
		if n.Name != "NIC 1" && n.Name != "NIC 2" {
			continue
		}
		members = append(members, n.UUID)
		p, _ := xenapi.PIF.GetByUUID(s, n.UUID)
		net, _ := xenapi.PIF.GetNetwork(s, p)
		vifs, _ := xenapi.Network.GetVIFs(s, net)
		for _, v := range vifs {
			vm, _ := xenapi.VIF.GetVM(s, v)
			if name, _ := xenapi.VM.GetNameLabel(s, vm); name != os.Getenv("XCP_LAB_VM") || name == "" {
				t.Fatalf("%s's network has a VIF of %q, not of XCP_LAB_VM", n.Name, name)
			}
			home[v] = net
		}
	}
	goHome = func() {
		for v, net := range home {
			if now, err := xenapi.VIF.GetNetwork(s, v); err == nil && now != net {
				if err := xenapi.VIF.Move(s, v, net); err != nil {
					t.Errorf("VIF %s back to its network: %v", v, err)
				}
			}
		}
	}
	if len(members) != 2 || opts.Bond.Limit != 4 || !opts.Bond.LACP {
		t.Fatalf("bond choices %+v", opts.Bond)
	}
	if w, err := ops.BondCreateWarning(c, members); err != nil || w != nil {
		t.Fatalf("create warning %+v %v", w, err)
	}
	if err := ops.CreateBond(c, ops.BondSpec{Members: members[:1], Mode: "active-backup", MTU: 1500}); err == nil {
		t.Fatal("a bond of one NIC")
	}
	if err := ops.CreateBond(c, ops.BondSpec{Members: members, Mode: "lacp", Hashing: "src_mac", MTU: 1500}); err != nil {
		t.Fatal(err)
	}
	bondNet, bn := byName("Bond 1+2")
	made[bondNet] = true
	if bondNet == "" || bn.OtherConfig["automatic"] != "false" || bn.OtherConfig["XenCenterCreateInProgress"] != "" || len(bn.PIFs) != len(hosts) {
		t.Fatalf("bond network %+v", bn)
	}
	bonds, _ := xenapi.Bond.GetAllRecords(s)
	if len(bonds) != bonds0+len(hosts) {
		t.Fatalf("%d bonds", len(bonds))
	}
	for _, b := range bonds {
		if dev, _ := xenapi.PIF.GetDevice(s, b.PrimarySlave); b.Mode != xenapi.BondModeLacp || b.Properties["hashing_algorithm"] != "src_mac" || dev != "eth1" {
			t.Errorf("bond %+v (primary member %s)", b, dev)
		}
	}
	nics, err := HostNicsTab(c, hostUUID)
	if err != nil {
		t.Fatal(err)
	}
	if !slices.ContainsFunc(nics.NICs, func(n HostNic) bool { return n.BondNetwork == bondNet }) || !nics.CanBond {
		t.Fatalf("NICs tab %+v", nics)
	}
	t.Logf("NICs: %v", nics.Table.Rows)
	tab, err := HostNetworkingTab(c, "pools", poolUUID)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("networks: %v", tab.Table.Rows)
	if i := slices.Index(tab.Table.Keys, bondNet); i < 0 || !tab.Networks[i].Bond || !tab.Networks[i].CanRemove || !tab.Networks[i].CanEdit {
		t.Fatalf("networking tab %+v", tab)
	}
	if len(tab.Table.Keys) != 2 { // the member networks are hidden: network 0 and the bond
		t.Errorf("networks %v", tab.Table.Rows)
	}
	if w, err := ops.BondDeleteWarning(c, bondNet); err != nil || w.Block || w.Text != "Are you sure that you want to delete Bond 1+2?" {
		t.Fatalf("delete warning %+v %v", w, err)
	}
	if err := ops.DestroyBond(c, bondNet); err != nil {
		t.Fatal(err)
	}
	if n, b := count(); n != nets0 || b != bonds0 {
		t.Fatalf("after the bond: %d networks, %d bonds", n, b)
	}
	goHome()
	for v, want := range home {
		if net, err := xenapi.VIF.GetNetwork(s, v); err != nil || net != want {
			t.Errorf("VIF %s is on %s, not back on %s: %v", v, net, want, err)
		}
	}

	// a VLAN network on NIC 1 of every server, then a private one
	opts, _ = NewNetworkOptions(c, "pools", poolUUID)
	i := slices.IndexFunc(opts.NICs, func(n WizardNIC) bool { return n.Name == "NIC 1" })
	if i < 0 || opts.Host != hostUUID || opts.ExternalName != "New Network" || opts.MinVLAN != 0 {
		t.Fatalf("wizard choices %+v", opts)
	}
	nic := opts.NICs[i]
	vlan := rules.VLANFor(3901, nic.VLANs)
	spec := ops.NewNetwork{External: true, Name: "xcpng-center-go-test-vlan", Host: opts.Host, NIC: nic.UUID, VLAN: fmt.Sprint(vlan), MTU: "1500"}
	bad := spec
	bad.VLAN = "-1"
	if err := ops.CreateNetwork(c, bad); err == nil {
		t.Fatal("VLAN -1")
	}
	if err := ops.CreateNetwork(c, spec); err != nil {
		t.Fatal(err)
	}
	vlanNet, vn := byName(spec.Name)
	made[vlanNet] = false
	if err := ops.CreateNetwork(c, spec); err == nil || !strings.Contains(err.Error(), "already in use") {
		t.Fatalf("the same VLAN twice: %v", err)
	}
	if len(vn.PIFs) != len(hosts) || vn.OtherConfig["automatic"] != "false" {
		t.Fatalf("VLAN network %+v", vn)
	}
	edit, err := ReadNetworkEdit(c, vlanNet)
	if err != nil || !edit.MTUState.Enabled || edit.MTU != 1500 || edit.Auto {
		t.Fatalf("VLAN network properties %+v %v", edit, err)
	}
	if err := SaveNetwork(c, vlanNet, ops.NetworkChanges{Name: "x", MTU: "100", NIC: edit.NIC, VLAN: fmt.Sprint(edit.VLAN)}); err == nil {
		t.Fatal("MTU 100")
	}
	// shortcut: the MTU itself is not changed on the lab (a VLAN's MTU can raise NIC 1's), so the unplug and plug of
	// UnplugPlugNetworkAction is not run here; upgrade when the lab has a NIC to spare.
	if err := SaveNetwork(c, vlanNet, ops.NetworkChanges{Name: " xcpng-center-go-test-vlan2 ", Description: "d", Auto: true, MTU: "1500", NIC: edit.NIC, VLAN: fmt.Sprint(edit.VLAN)}); err != nil {
		t.Fatal(err)
	}
	if _, n := byName("xcpng-center-go-test-vlan2"); n.NameDescription != "d" || n.OtherConfig["automatic"] != "true" {
		t.Fatalf("renamed %+v", n)
	}
	tab, _ = HostNetworkingTab(c, "hosts", hostUUID)
	if i := slices.Index(tab.Table.Keys, vlanNet); i < 0 || !tab.Networks[i].CanRemove || tab.Networks[i].Bond {
		t.Fatalf("VLAN row %+v", tab)
	}
	if err := ops.RemoveNetwork(c, vlanNet); err != nil {
		t.Fatal(err)
	}
	if err := ops.CreateNetwork(c, ops.NewNetwork{Name: "xcpng-center-go-test-private", Auto: true}); err != nil {
		t.Fatal(err)
	}
	priv, pn := byName("xcpng-center-go-test-private")
	made[priv] = false
	if len(pn.PIFs) != 0 || pn.OtherConfig["automatic"] != "true" {
		t.Fatalf("private network %+v", pn)
	}
	if edit, err := ReadNetworkEdit(c, priv); err != nil || edit.MTUState.Enabled || !edit.Auto {
		t.Fatalf("private network properties %+v %v", edit, err)
	}
	if err := SaveNetwork(c, priv, ops.NetworkChanges{Name: "xcpng-center-go-test-private2", Auto: false}); err != nil {
		t.Fatal(err)
	}
	if err := ops.RemoveNetwork(c, priv); err != nil {
		t.Fatal(err)
	}
	if tab, _ := HostNetworkingTab(c, "hosts", hostUUID); !slices.ContainsFunc(tab.Networks, func(n HostNetwork) bool { return !n.CanRemove && n.RemoveTip != "" }) {
		t.Errorf("the physical NICs' networks can be removed: %+v", tab.Networks)
	}
	if err := ops.RescanPIFs(c, hostUUID); err != nil {
		t.Fatal(err)
	}
	if n, b := count(); n != nets0 || b != bonds0 {
		t.Fatalf("after the networks: %d networks, %d bonds", n, b)
	}
}

// TestLabNetworkSettings is a network's Network Settings page (EditNetworkPage): a throwaway private network is made
// external on a NIC that is no management interface, its VLAN is changed and it is made private again, the PIFs of
// every server checked after each step; then a throwaway bond of two spare NICs (when the lab has them) changes its
// mode and LACP hashing. Everything made is destroyed, also on a failure (XCP_ALLOW_WRITE=1).
func TestLabNetworkSettings(t *testing.T) {
	c := labConnection(t)
	if os.Getenv("XCP_ALLOW_WRITE") != "1" {
		t.Skip("XCP_ALLOW_WRITE=1 not set")
	}
	s := c.Session
	hosts, _ := xenapi.Host.GetAll(s)
	const name = "xcpng-center-go-test-net-settings"
	byName := func(name string) (string, xenapi.NetworkRecord) {
		all, _ := xenapi.Network.GetAllRecords(s)
		for _, n := range all {
			if n.NameLabel == name {
				return n.UUID, n
			}
		}
		return "", xenapi.NetworkRecord{}
	}
	if u, _ := byName(name); u != "" {
		t.Fatalf("the lab already has a network %q", name)
	}
	if err := ops.CreateNetwork(c, ops.NewNetwork{Name: name}); err != nil {
		t.Fatal(err)
	}
	net, _ := byName(name)
	t.Cleanup(func() {
		if _, err := xenapi.Network.GetByUUID(s, net); err == nil {
			if err := ops.RemoveNetwork(c, net); err != nil {
				t.Errorf("remove %s: %v", name, err)
			}
		}
	})
	edit, err := ReadNetworkEdit(c, net)
	if err != nil || !edit.NICState.Shown || !edit.NICState.Enabled || edit.NIC != "" || edit.Bond || edit.MTUState.Enabled || !edit.VLANMTU.Enabled {
		t.Fatalf("private network settings %+v %v", edit, err)
	}
	// a NIC that is no management interface (primary or secondary, an IP on it) on any server
	spare := func(n ops.SettingsNIC) bool {
		pifs, _ := xenapi.PIF.GetAllRecords(s)
		dev := ""
		for _, p := range pifs {
			if p.UUID == n.UUID {
				dev = p.Device
			}
		}
		if dev == "" || n.Name == "NIC 0" {
			return false
		}
		for _, p := range pifs {
			if p.Device == dev && (p.Management || p.IPConfigurationMode != xenapi.IPConfigurationModeNone) {
				return false
			}
		}
		return true
	}
	i := slices.IndexFunc(edit.NICs, func(n ops.SettingsNIC) bool { return n.Name == "NIC 1" && spare(n) })
	if i < 0 {
		t.Fatalf("no spare NIC 1 in %+v", edit.NICs)
	}
	nic := edit.NICs[i]
	nicRef, _ := xenapi.PIF.GetByUUID(s, nic.UUID)
	device, _ := xenapi.PIF.GetDevice(s, nicRef)
	pifsAre := func(step string, vlan int) {
		_, n := byName(name)
		if vlan < 0 {
			if len(n.PIFs) != 0 {
				t.Fatalf("%s: %d PIFs on a private network", step, len(n.PIFs))
			}
			return
		}
		on := map[xenapi.HostRef]bool{}
		for _, r := range n.PIFs {
			p, err := xenapi.PIF.GetRecord(s, r)
			if err != nil || p.VLAN != vlan || p.Device != device || on[p.Host] {
				t.Fatalf("%s: PIF %+v %v", step, p, err)
			}
			on[p.Host] = true
		}
		if len(on) != len(hosts) {
			t.Fatalf("%s: PIFs on %d of %d servers", step, len(on), len(hosts))
		}
	}
	ch := ops.NetworkChanges{Name: name, MTU: "1500", NIC: nic.UUID, VLAN: fmt.Sprint(rules.VLANFor(3999, nic.VLANs))}
	bad := ch
	bad.VLAN = "4095"
	if err := SaveNetwork(c, net, bad); err == nil || !strings.Contains(err.Error(), "4094") {
		t.Fatalf("VLAN 4095: %v", err)
	}
	pifsAre("a refused VLAN", -1)
	bad.VLAN, bad.NIC = "5", "no-such-nic"
	if err := SaveNetwork(c, net, bad); err == nil {
		t.Fatal("an unknown NIC")
	}
	if err := SaveNetwork(c, net, ch); err != nil {
		t.Fatal(err)
	}
	vlan := rules.VLANFor(3999, nic.VLANs)
	pifsAre("private to external", vlan)
	edit, err = ReadNetworkEdit(c, net)
	if err != nil || edit.NIC != nic.UUID || edit.VLAN != vlan || !edit.NICState.Enabled || !edit.MTUState.Enabled {
		t.Fatalf("external network settings %+v %v", edit, err)
	}
	i = slices.IndexFunc(edit.NICs, func(n ops.SettingsNIC) bool { return n.UUID == nic.UUID })
	if !slices.Contains(edit.NICs[i].VLANs, vlan) { // as NetworkAction's check; the unchanged VLAN is not checked
		t.Errorf("the network's own VLAN is not in use: %v", edit.NICs[i].VLANs)
	}
	if err := SaveNetwork(c, net, ops.NetworkChanges{Name: name, MTU: "1500", NIC: nic.UUID, VLAN: fmt.Sprint(vlan), Description: "d"}); err != nil {
		t.Fatalf("a save with the same NIC and VLAN: %v", err)
	}
	pifsAre("the same NIC and VLAN", vlan)
	ch.VLAN = fmt.Sprint(rules.VLANFor(3998, edit.NICs[i].VLANs))
	if err := SaveNetwork(c, net, ch); err != nil {
		t.Fatal(err)
	}
	pifsAre("a new VLAN", rules.VLANFor(3998, edit.NICs[i].VLANs))
	ch.NIC, ch.VLAN = "", ""
	if err := SaveNetwork(c, net, ch); err != nil {
		t.Fatal(err)
	}
	pifsAre("external to private", -1)

	// the bond mode, on a throwaway bond of the first two spare NICs (never NIC 0, the management interface)
	var members []string
	for _, n := range edit.NICs {
		if !strings.HasPrefix(n.Name, "Bond") && spare(n) && len(members) < 2 {
			p, _ := xenapi.PIF.GetByUUID(s, n.UUID)
			pn, _ := xenapi.PIF.GetNetwork(s, p)
			if vifs, _ := xenapi.Network.GetVIFs(s, pn); len(vifs) == 0 {
				members = append(members, n.UUID)
			} else {
				t.Logf("%s: its network has %d VIFs", n.Name, len(vifs))
			}
		}
	}
	if len(members) < 2 {
		t.Logf("no two spare NICs for a bond in %+v: the bond mode is not tried", edit.NICs)
		return
	}
	bonds0, _ := xenapi.Bond.GetAll(s)
	if err := ops.CreateBond(c, ops.BondSpec{Members: members, Mode: "balance-slb", MTU: 1500}); err != nil {
		t.Fatal(err)
	}
	var bondNet string
	t.Cleanup(func() {
		if bondNet != "" {
			if err := ops.DestroyBond(c, bondNet); err != nil {
				t.Errorf("delete the bond: %v", err)
			}
		}
		if b, _ := xenapi.Bond.GetAll(s); len(b) != len(bonds0) {
			t.Errorf("%d bonds, the lab had %d", len(b), len(bonds0))
		}
	})
	all, _ := xenapi.Network.GetAllRecords(s)
	for _, n := range all {
		for _, r := range n.PIFs {
			if p, _ := xenapi.PIF.GetRecord(s, r); len(p.BondMasterOf) > 0 {
				if b, _ := xenapi.Bond.GetSlaves(s, p.BondMasterOf[0]); len(b) == 2 {
					sl0, _ := xenapi.PIF.GetUUID(s, b[0])
					sl1, _ := xenapi.PIF.GetUUID(s, b[1])
					if slices.Contains(members, sl0) && slices.Contains(members, sl1) {
						bondNet = n.UUID
					}
				}
			}
		}
	}
	if bondNet == "" {
		t.Fatal("the bond network is not found")
	}
	modes := func(step, mode, hashing string) {
		nr, _ := xenapi.Network.GetByUUID(s, bondNet)
		pifs, _ := xenapi.Network.GetPIFs(s, nr)
		for _, r := range pifs {
			p, _ := xenapi.PIF.GetRecord(s, r)
			b, err := xenapi.Bond.GetRecord(s, p.BondMasterOf[0])
			if err != nil || string(b.Mode) != mode || hashing != "" && b.Properties["hashing_algorithm"] != hashing {
				t.Fatalf("%s: bond %+v %v", step, b, err)
			}
		}
	}
	be, err := ReadNetworkEdit(c, bondNet)
	if err != nil || !be.Bond || be.Mode != "balance-slb" || be.NICState.Shown {
		t.Fatalf("bond network settings %+v %v", be, err)
	}
	bc := ops.NetworkChanges{Name: be.Name, MTU: fmt.Sprint(be.MTU), Auto: be.Auto, Mode: "lacp", Hashing: ""}
	if err := SaveNetwork(c, bondNet, bc); err == nil {
		t.Fatal("LACP without a hashing algorithm")
	}
	modes("a refused LACP", "balance-slb", "")
	bc.Mode = "active-backup"
	if err := SaveNetwork(c, bondNet, bc); err != nil {
		t.Fatal(err)
	}
	modes("active-passive", "active-backup", "")
	if be.LACP {
		bc.Mode, bc.Hashing = "lacp", "tcpudp_ports"
		if err := SaveNetwork(c, bondNet, bc); err != nil {
			t.Fatal(err)
		}
		modes("LACP", "lacp", "tcpudp_ports")
		bc.Hashing = "src_mac"
		if err := SaveNetwork(c, bondNet, bc); err != nil {
			t.Fatal(err)
		}
		modes("LACP on source MAC", "lacp", "src_mac")
	}
}

// TestLabObjectProperties changes the description, folder and tags of the default SR, a network and the pool through
// the General page and restores them (XCP_ALLOW_WRITE=1).
func TestLabObjectProperties(t *testing.T) {
	c := labConnection(t)
	if os.Getenv("XCP_ALLOW_WRITE") != "1" {
		t.Skip("XCP_ALLOW_WRITE=1 not set")
	}
	s := c.Session
	pools, _ := xenapi.Pool.GetAllRecords(s)
	var poolUUID, srUUID string
	for _, p := range pools {
		poolUUID = p.UUID
		srUUID, _ = xenapi.SR.GetUUID(s, p.DefaultSR)
	}
	nets, _ := xenapi.Network.GetAllRecords(s)
	var netUUID string
	for _, n := range nets {
		if len(n.PIFs) > 0 {
			netUUID = n.UUID
			break
		}
	}
	poolRef, _ := xenapi.Pool.GetByUUID(s, poolUUID)
	other, _ := xenapi.Pool.GetOtherConfig(s, poolRef)
	empty, hadEmpty := other["EMPTY_FOLDERS"]
	t.Cleanup(func() { // the folders the test emptied stay as empty folders (FolderAction): give the pool its list back
		_ = xenapi.Pool.RemoveFromOtherConfig(s, poolRef, "EMPTY_FOLDERS")
		if hadEmpty {
			_ = xenapi.Pool.AddToOtherConfig(s, poolRef, "EMPTY_FOLDERS", empty)
		}
	})
	for _, o := range []struct{ kind, uuid string }{{"srs", srUUID}, {"networks", netUUID}, {"pools", poolUUID}} {
		before, err := ReadObjectProperties(c, o.kind, o.uuid)
		if err != nil {
			t.Fatal(o.kind, err)
		}
		restore := ops.ObjectChanges{Description: &before.Description, Folder: &before.Folder, Tags: before.Tags, TagsSet: true}
		desc, folder := before.Description+" (go test)", "/xcpng-center-go-test/sub"
		tags := append(append([]string{}, before.Tags...), "xcpng-center-go-test")
		if err := ops.SaveObject(c, o.kind, o.uuid, ops.ObjectChanges{Description: &desc, Folder: &folder, Tags: tags, TagsSet: true}); err != nil {
			t.Fatal(o.kind, err)
		}
		after, _ := ReadObjectProperties(c, o.kind, o.uuid)
		if after.Description != desc || after.Folder != folder || !slices.Contains(after.Tags, "xcpng-center-go-test") ||
			!slices.Contains(after.AllTags, "xcpng-center-go-test") || !slices.Contains(after.Folders, folder) {
			t.Errorf("%s: %+v", o.kind, after)
		}
		bad := "no/slash"
		if err := ops.SaveObject(c, o.kind, o.uuid, ops.ObjectChanges{Folder: &bad}); err == nil {
			t.Errorf("%s: a bad folder is saved", o.kind)
		}
		if err := ops.SaveObject(c, o.kind, o.uuid, restore); err != nil {
			t.Fatal(o.kind, "not restored:", err)
		}
		back, _ := ReadObjectProperties(c, o.kind, o.uuid)
		if back.Description != before.Description || back.Folder != before.Folder || len(back.Tags) != len(before.Tags) {
			t.Fatalf("%s not restored: %+v", o.kind, back)
		}
		if !slices.Contains(back.Folders, folder) { // left empty, kept in EMPTY_FOLDERS
			t.Errorf("%s: the emptied folder is lost: %v", o.kind, back.Folders)
		}
		empty := ""
		if err := ops.SaveObject(c, o.kind, o.uuid, ops.ObjectChanges{Description: &empty, Name: &empty}); err == nil {
			t.Errorf("%s: a blank name is saved", o.kind)
		}
	}
}

// TestLabVapps creates a vApp with XCP_LAB_VM, renames it, takes the VM out and back with Assign to vApp, and deletes
// the vApp (the VM stays, in no vApp) (XCP_ALLOW_WRITE=1).
func TestLabVapps(t *testing.T) {
	c := labConnection(t)
	name := os.Getenv("XCP_LAB_VM")
	if name == "" || os.Getenv("XCP_ALLOW_WRITE") != "1" {
		t.Skip("XCP_LAB_VM and XCP_ALLOW_WRITE=1 not set")
	}
	s := c.Session
	refs, _ := xenapi.VM.GetByNameLabel(s, name)
	if len(refs) != 1 {
		t.Fatalf("%s not found", name)
	}
	vmUUID, _ := xenapi.VM.GetUUID(s, refs[0])
	if a, _ := xenapi.VM.GetAppliance(s, refs[0]); a != "OpaqueRef:NULL" && a != "" {
		t.Skip("the lab VM is already in a vApp")
	}
	uuid, err := ops.CreateVapp(c, ops.VappSpec{Name: "xcpng-center-go-test vApp", Description: "test", VMs: []string{vmUUID}})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if _, err := xenapi.VMAppliance.GetByUUID(s, uuid); err == nil {
			if err := ops.DeleteVapp(c, uuid); err != nil {
				t.Errorf("not deleted: %v", err)
			}
		}
	})
	d, err := ReadVapps(c)
	if err != nil {
		t.Fatal(err)
	}
	i := slices.IndexFunc(d.Vapps, func(a Vapp) bool { return a.UUID == uuid })
	if i < 0 || len(d.Vapps[i].VMs) != 1 || d.Vapps[i].VMs[0].UUID != vmUUID {
		t.Fatalf("%+v", d.Vapps)
	}
	if err := ops.SaveVapp(c, uuid, ops.VappSpec{Name: "xcpng-center-go-test vApp 2", VMs: []string{vmUUID}}); err != nil {
		t.Fatal(err)
	}
	if err := ops.AssignVapp(c, uuid, vmUUID); err != nil { // already in it: leaves
		t.Fatal(err)
	}
	if vms, _ := xenapi.VMAppliance.GetVMs(s, mustVapp(t, c, uuid)); len(vms) != 0 {
		t.Fatal("the VM did not leave the vApp")
	}
	if err := ops.AssignVapp(c, uuid, vmUUID); err != nil { // joins again
		t.Fatal(err)
	}
	if err := ops.SaveVapp(c, uuid, ops.VappSpec{Name: " "}); err == nil {
		t.Fatal("a blank name is saved")
	}
	if err := ops.DeleteVapp(c, uuid); err != nil {
		t.Fatal(err)
	}
	if a, _ := xenapi.VM.GetAppliance(s, refs[0]); a != "OpaqueRef:NULL" {
		t.Fatalf("the VM is still in a vApp: %s", a)
	}
}

func mustVapp(t *testing.T, c *xapi.Connection, uuid string) xenapi.VMApplianceRef {
	ref, err := xenapi.VMAppliance.GetByUUID(c.Session, uuid)
	if err != nil {
		t.Fatal(err)
	}
	return ref
}

// TestLabVtpm makes a small UEFI VM, adds a vTPM, removes it and deletes the VM; the lab VM (BIOS) is refused
// (XCP_ALLOW_WRITE=1).
func TestLabVtpm(t *testing.T) {
	c := labConnection(t)
	if os.Getenv("XCP_ALLOW_WRITE") != "1" {
		t.Skip("XCP_ALLOW_WRITE=1 not set")
	}
	s := c.Session
	opts, err := ReadNewVMOptions(c)
	if err != nil {
		t.Fatal(err)
	}
	var tpl string
	for _, x := range opts.Templates {
		if x.Name == "Other install media" {
			tpl = x.UUID
		}
	}
	created, err := ops.CreateVM(c, ops.NewVM{Template: tpl, Name: "xcpng-center-go-test-vtpm", Install: "cd", BootMode: "uefi", VCPUs: 1, Cores: 1,
		Memory: 256 << 20, Disks: []ops.NewDisk{{Position: "0", Name: "xcpng-center-go-test-vtpm 0", Size: 64 << 20, SR: opts.DefaultSR}}})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		del, _ := ReadVMDelete(c, created.UUID)
		var vbds []string
		for _, x := range del.Disks {
			vbds = append(vbds, x.UUID)
		}
		if err := ops.DeleteVM(c, created.UUID, vbds, nil); err != nil {
			t.Errorf("not deleted: %v", err)
		}
	})
	v, err := ops.ReadVtpms(c, created.UUID)
	if err != nil || v.Blocker != "" || v.AddBlocker != "" || len(v.Vtpms) != 0 {
		t.Fatalf("%+v %v", v, err)
	}
	if err := ops.AddVtpm(c, created.UUID); err != nil {
		t.Fatal(err)
	}
	v, _ = ops.ReadVtpms(c, created.UUID)
	if len(v.Vtpms) != 1 || v.Vtpms[0].RemoveBlocker != "" {
		t.Fatalf("%+v", v)
	}
	if err := ops.RemoveVtpm(c, created.UUID, v.Vtpms[0].UUID); err != nil {
		t.Fatal(err)
	}
	if v, _ = ops.ReadVtpms(c, created.UUID); len(v.Vtpms) != 0 {
		t.Fatalf("not removed: %+v", v)
	}
	if name := os.Getenv("XCP_LAB_VM"); name != "" {
		refs, _ := xenapi.VM.GetByNameLabel(s, name)
		uuid, _ := xenapi.VM.GetUUID(s, refs[0])
		if v, _ := ops.ReadVtpms(c, uuid); v.Blocker == "" {
			t.Error("a BIOS VM may have vTPMs")
		}
	}
}

// TestLabHA turns HA on with the throwaway VM (XCP_LAB_VM) protected and the shared NFS SR as heartbeat, changes its
// priority, order, delay and the failures to tolerate (the HA dialog and the Start Options page), then turns HA off
// (XCP_ALLOW_WRITE=1). HA ends off and the VM's and the pool's settings as they were, even when the test fails.
func TestLabHA(t *testing.T) {
	c := labConnection(t)
	name := os.Getenv("XCP_LAB_VM")
	if name == "" || os.Getenv("XCP_ALLOW_WRITE") != "1" {
		t.Skip("XCP_LAB_VM and XCP_ALLOW_WRITE=1 not set")
	}
	s := c.Session
	vms, err := xenapi.VM.GetByNameLabel(s, name)
	if err != nil || len(vms) != 1 {
		t.Fatalf("VM %q: %v %v", name, vms, err)
	}
	ref := vms[0]
	vm, err := xenapi.VM.GetRecord(s, ref)
	if err != nil {
		t.Fatal(err)
	}
	pools, err := xenapi.Pool.GetAllRecords(s)
	if err != nil {
		t.Fatal(err)
	}
	var poolRef xenapi.PoolRef
	var pool xenapi.PoolRecord
	for r, p := range pools {
		poolRef, pool = r, p
	}
	if pool.HaEnabled {
		t.Fatal("HA is already on in the lab: this test only runs from HA off")
	}
	t.Cleanup(func() { // HA off, then the VM's and the pool's own settings back
		if on, _ := xenapi.Pool.GetHaEnabled(s, poolRef); on {
			if err := ops.DisableHA(c); err != nil {
				t.Errorf("cleanup: HA is still on: %v", err)
			}
		}
		for _, err := range []error{xenapi.VM.SetHaRestartPriority(s, ref, vm.HaRestartPriority), xenapi.VM.SetOrder(s, ref, vm.Order),
			xenapi.VM.SetStartDelay(s, ref, vm.StartDelay), xenapi.Pool.SetHaHostFailuresToTolerate(s, poolRef, pool.HaHostFailuresToTolerate)} {
			if err != nil {
				t.Errorf("cleanup: %v", err)
			}
		}
	})

	// refused with HA off
	if err := ops.DisableHA(c); err == nil || !strings.Contains(err.Error(), "HA is not enabled") {
		t.Errorf("Disable with HA off: %v", err)
	}
	if err := ops.SetHAPriorities(c, []ops.HAStartup{{VM: vm.UUID, Priority: rules.HARestart}}, 0); err == nil {
		t.Error("priorities with HA off")
	}
	if err := ops.SaveVMProperties(c, vm.UUID, ops.VMChanges{HARestartPriority: new(string)}); err == nil {
		t.Error("a restart priority on the Start Options page with HA off")
	}

	cfg, err := ReadHAConfig(c, true)
	if err != nil || cfg.Blocker != "" || cfg.HAEnabled {
		t.Fatalf("HA config: %+v %v", cfg, err)
	}
	i := slices.IndexFunc(cfg.VMs, func(v HAPlanVM) bool { return v.UUID == vm.UUID })
	if i < 0 {
		t.Fatalf("%s is not in the plan: %+v", name, cfg.VMs)
	}
	t.Logf("plan: %+v (hosts %d)", cfg.VMs[i], cfg.Hosts)
	if cfg.VMs[i].NotAgile != "" {
		t.Fatalf("%s is not agile: %s", name, cfg.VMs[i].NotAgile)
	}
	srs, err := ops.HeartbeatSRs(c)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("heartbeat SRs: %+v", srs)
	hb := slices.IndexFunc(srs, func(h ops.HeartbeatSR) bool { return h.Name == "lab-nfs" && h.Reason == "" })
	if hb < 0 {
		t.Fatal("lab-nfs cannot hold the heartbeat")
	}
	plan := []ops.HAStartup{{VM: vm.UUID, Priority: rules.HARestart, Order: 2, StartDelay: 10}}
	max, err := ops.HAMaxFailures(c, map[string]string{vm.UUID: rules.HARestart})
	if err != nil || max < 0 {
		t.Fatalf("max failures: %d %v", max, err)
	}
	t.Logf("max failures with %s restarted: %d", name, max)
	if err := ops.EnableHA(c, srs[hb].UUID, plan, max+1); err == nil || !strings.Contains(err.Error(), "HA cannot be guaranteed") {
		t.Fatalf("overcommitted plan: %v", err)
	}
	if on, _ := xenapi.Pool.GetHaEnabled(s, poolRef); on {
		t.Fatal("HA on after a refused plan")
	}
	ntol := min(1, max)
	if err := ops.EnableHA(c, srs[hb].UUID, plan, ntol); err != nil {
		p, _ := xenapi.Pool.GetRecord(s, poolRef)
		t.Fatalf("enable HA: %v (HA on: %v)", err, p.HaEnabled)
	}
	cfg, err = ReadHAConfig(c, true)
	if err != nil || !cfg.HAEnabled || cfg.Blocker != "" || cfg.Ntol != ntol {
		t.Fatalf("after enabling: %+v %v", cfg, err)
	}
	if v, _ := xenapi.VM.GetRecord(s, ref); v.HaRestartPriority != rules.HARestart || v.Order != 2 || v.StartDelay != 10 {
		t.Errorf("after enabling: priority %q order %d delay %d", v.HaRestartPriority, v.Order, v.StartDelay)
	}
	if err := ops.EnableHA(c, srs[hb].UUID, plan, ntol); err == nil {
		t.Error("enable HA twice")
	}

	// the restart priorities dialog: best effort, another order and delay, no failure tolerated
	if err := ops.SetHAPriorities(c, []ops.HAStartup{{VM: vm.UUID, Priority: rules.HABestEffort, Order: 3, StartDelay: 5}}, 0); err != nil {
		t.Fatal(err)
	}
	if v, _ := xenapi.VM.GetRecord(s, ref); v.HaRestartPriority != rules.HABestEffort || v.Order != 3 || v.StartDelay != 5 {
		t.Errorf("after the priorities: %q %d %d", v.HaRestartPriority, v.Order, v.StartDelay)
	}
	if n, _ := xenapi.Pool.GetHaHostFailuresToTolerate(s, poolRef); n != 0 {
		t.Errorf("failures to tolerate %d, want 0", n)
	}
	// the Start Options page: restart again, the pool's failures to tolerate kept
	restart, order := rules.HARestart, 4
	if err := ops.SaveVMProperties(c, vm.UUID, ops.VMChanges{HARestartPriority: &restart, Order: &order}); err != nil {
		t.Fatal(err)
	}
	if v, _ := xenapi.VM.GetRecord(s, ref); v.HaRestartPriority != rules.HARestart || v.Order != 4 || v.StartDelay != 5 {
		t.Errorf("after Start Options: %q %d %d", v.HaRestartPriority, v.Order, v.StartDelay)
	}
	if props, err := ReadVMProperties(c, vm.UUID); err != nil || !props.HAEnabled || props.HARestartPriority != rules.HARestart || len(props.HAPriorities) != 3 {
		t.Errorf("Start Options read: %+v %v", props.HAPriorities, err)
	}

	if why := ops.DisableHABlocker(c); why != "" {
		t.Fatal(why)
	}
	if err := ops.DisableHA(c); err != nil {
		t.Fatal(err)
	}
	if on, _ := xenapi.Pool.GetHaEnabled(s, poolRef); on {
		t.Fatal("HA still on")
	}
}

func TestLabMemoryBars(t *testing.T) {
	c := labConnection(t)
	pools, _ := xenapi.Pool.GetAllRecords(c.Session)
	for _, p := range pools {
		b, err := ReadMemoryBars(c, "pools", p.UUID)
		if err != nil || len(b.Hosts) == 0 {
			t.Fatal(b, err)
		}
		for _, h := range b.Hosts {
			// Xen and the control domain first, free memory last, each with a size
			if n := len(h.Segments); n < 2 || h.Segments[0].Bytes <= 0 || h.Segments[n-1].Name != "Free" {
				t.Fatalf("%s: %+v", h.Name, h.Segments)
			}
			t.Logf("%s: %+v", h.Name, h.Segments)
		}
	}
	if name := os.Getenv("XCP_LAB_VM"); name != "" {
		refs, _ := xenapi.VM.GetByNameLabel(c.Session, name)
		uuid, _ := xenapi.VM.GetUUID(c.Session, refs[0])
		b, err := ReadMemoryBars(c, "vms", uuid)
		if err != nil || b.DynamicMax <= 0 || b.DynamicMax > 1 || b.UsedText == "" {
			t.Fatal(b, err)
		}
		t.Logf("%s: %+v", name, b)
	}
}

// TestLabGeneral reads the General tab of the pool, its servers, the lab VM and the default SR (read-only) and checks
// the key rows of each section (Details.General).
func TestLabGeneral(t *testing.T) {
	c := labConnection(t)
	s := c.Session
	pools, err := xenapi.Pool.GetAllRecords(s)
	if err != nil {
		t.Fatal(err)
	}
	var pool xenapi.PoolRecord
	for _, p := range pools {
		pool = p
	}
	type want struct{ section, row string }
	read := func(name, kind, uuid string, keys ...want) {
		secs, err := ReadGeneral(c, kind, uuid)
		if err != nil {
			t.Fatalf("%s: %v", name, err)
		}
		rows := map[want]string{}
		for _, sec := range secs {
			for _, r := range sec.Rows {
				rows[want{sec.Title, r[0]}] = r[1]
				t.Logf("%s / %s / %s: %q", name, sec.Title, r[0], r[1])
			}
		}
		for _, k := range keys {
			if rows[k] == "" {
				t.Errorf("%s: no %s / %s", name, k.section, k.row)
			}
		}
	}
	read("pool", "pools", pool.UUID, want{"General", "Name"}, want{"General", "Tags"}, want{"General", "Number of Sockets"},
		want{"General", "Version"}, want{"General", "UUID"}, want{"Management Interfaces", "DNS hostname"}, want{"Management Interfaces", "Management interface"})
	hosts, err := xenapi.Host.GetAllRecords(s)
	if err != nil {
		t.Fatal(err)
	}
	for _, h := range hosts {
		read(h.NameLabel, "hosts", h.UUID, want{"General", "Pool coordinator"}, want{"General", "Enabled"}, want{"General", "Server uptime"},
			want{"General", "Log destination"}, want{"Certificates", "Server identity certificate"}, want{"Version Details", "Version"},
			want{"Management Interfaces", "Management interface"}, want{"Memory", "Server"}, want{"CPUs", fmt.Sprintf("CPU 0 - %d", len(h.HostCPUs)-1)})
	}
	if name := os.Getenv("XCP_LAB_VM"); name != "" {
		vms, err := ListObjects(c, "vms")
		if err != nil {
			t.Fatal(err)
		}
		i := slices.IndexFunc(vms, func(o Object) bool { return o.Name == name })
		if i < 0 {
			t.Fatalf("no VM %s", name)
		}
		read(name, "vms", vms[i].UUID, want{"General", "Operating System"}, want{"General", "Virtualization mode"}, want{"General", "UUID"},
			want{"Boot Options", "Auto-boot"}, want{"CPUs", "Virtual CPUs"}, want{"CPUs", "Topology"})
	}
	if sr, err := xenapi.SR.GetUUID(s, pool.DefaultSR); err == nil {
		read("default SR", "srs", sr, want{"General", "Type"}, want{"General", "Size"}, want{"Status", "State"}, want{"Multipathing", "Multipath capable"})
	} else {
		t.Errorf("default SR: %v", err)
	}
}

// TestLabHardware reads the GPU tab of the pool and its servers (the lab's emulated GPUs do no vGPU: the "no GPU" text)
// and the USB tab of the servers that show one (read-only).
func TestLabHardware(t *testing.T) {
	c := labConnection(t)
	for _, kind := range []string{"pools", "hosts"} {
		objs, err := ListObjects(c, kind)
		if err != nil {
			t.Fatal(err)
		}
		for _, o := range objs {
			if !o.GPUTab {
				t.Fatalf("%s %s: no GPU tab (restrict_gpu is false on the lab)", kind, o.Name)
			}
			g, err := ReadGPU(c, kind, o.UUID)
			if err != nil || len(g.Groups) == 0 && g.NoGPU == "" {
				t.Fatal(o.Name, g, err)
			}
			t.Logf("%s: %+v", o.Name, g)
			if o.USBTab {
				u, err := ReadUSB(c, o.UUID)
				if err != nil || len(u.Devices) == 0 {
					t.Fatal(o.Name, u, err)
				}
				t.Logf("%s: %+v", o.Name, u.Table)
			}
		}
	}
}

// TestLabHostLogAndMultipath sets lab-2's remote syslog server and puts it back (Log Destination), and checks that
// Multipathing is refused outside maintenance mode (entering it would move the lab VM).
func TestLabHostLogAndMultipath(t *testing.T) {
	if os.Getenv("XCP_ALLOW_WRITE") != "1" {
		t.Skip("XCP_ALLOW_WRITE=1 not set")
	}
	c := labConnection(t)
	refs, err := xenapi.Host.GetByNameLabel(c.Session, "xcp-lab-2")
	if err != nil || len(refs) == 0 {
		t.Fatal(err)
	}
	uuid, _ := xenapi.Host.GetUUID(c.Session, refs[0])
	p, err := ReadObjectProperties(c, "hosts", uuid)
	if err != nil || p.LogDestination == nil || p.Multipath == nil {
		t.Fatal(p, err)
	}
	was := *p.LogDestination
	t.Cleanup(func() {
		if err := ops.SaveObject(c, "hosts", uuid, ops.ObjectChanges{LogDestination: &was}); err != nil {
			t.Error("restore:", err)
		}
	})
	bad := "-bad-"
	if err := ops.SaveObject(c, "hosts", uuid, ops.ObjectChanges{LogDestination: &bad}); err == nil {
		t.Fatal("an invalid server name was taken")
	}
	to := "syslog.lab.test"
	if err := ops.SaveObject(c, "hosts", uuid, ops.ObjectChanges{LogDestination: &to}); err != nil {
		t.Fatal(err)
	}
	if p, _ := ReadObjectProperties(c, "hosts", uuid); p.LogDestination == nil || *p.LogDestination != to {
		t.Fatal(p.LogDestination)
	}
	if !p.MultipathEditable {
		on := !*p.Multipath
		if err := ops.SaveObject(c, "hosts", uuid, ops.ObjectChanges{Multipath: &on}); err == nil {
			t.Fatal("multipathing changed outside maintenance mode")
		}
	}
}

// TestLabSrDisks drives the SR Storage tab's buttons on lab-nfs (XCP_ALLOW_WRITE=1): Add a 32 MB disk attached to no
// VM, Rescan, Properties, Move to the coordinator's local SR and back, Delete; and reads the pool's Storage tab. It
// touches only the disks it makes (their names start with xcpng-center-go-test-srdisk).
func TestLabSrDisks(t *testing.T) {
	c := labConnection(t)
	if os.Getenv("XCP_ALLOW_WRITE") != "1" {
		t.Skip("XCP_ALLOW_WRITE=1 not set")
	}
	s := c.Session
	const prefix = "xcpng-center-go-test-srdisk"
	t.Cleanup(func() { // whatever a failed step left, on any SR
		vdis, _ := xenapi.VDI.GetAllRecords(s)
		for r, v := range vdis {
			if strings.HasPrefix(v.NameLabel, prefix) && len(v.VBDs) == 0 {
				_ = xenapi.VDI.Destroy(s, r)
			}
		}
	})
	srs, err := xenapi.SR.GetAllRecords(s)
	if err != nil {
		t.Fatal(err)
	}
	pools, err := xenapi.Pool.GetAllRecords(s)
	if err != nil || len(pools) != 1 {
		t.Fatal(pools, err)
	}
	var pool xenapi.PoolRecord
	for _, p := range pools {
		pool = p
	}
	nfs, local := "", ""
	for _, sr := range srs {
		if sr.NameLabel == "lab-nfs" {
			nfs = sr.UUID
		}
		if sr.NameLabel == "Local storage" && !sr.Shared && len(sr.PBDs) == 1 && local == "" {
			if h, err := xenapi.PBD.GetHost(s, sr.PBDs[0]); err == nil && h == pool.Master {
				local = sr.UUID
			}
		}
	}
	if nfs == "" || local == "" {
		t.Fatalf("no lab-nfs (%q) or local SR on the coordinator (%q)", nfs, local)
	}

	opts, err := NewSRDiskOptions(c, nfs)
	if err != nil || opts.DefaultSR != nfs || opts.Name == "" {
		t.Fatalf("Add's defaults: %+v %v", opts, err)
	}
	uuid, err := ops.CreateSRDisk(c, nfs, prefix, "lab test", 32<<20)
	if err != nil {
		t.Fatal(err)
	}
	if err := ops.RescanSRs(c, nfs); err != nil {
		t.Fatal(err)
	}
	tab, err := SRStorageTab(c, nfs)
	if err != nil {
		t.Fatal(err)
	}
	i := slices.Index(tab.Table.Keys, uuid)
	if i < 0 || !tab.CanAdd || !tab.CanRescan {
		t.Fatalf("the new disk is not in the tab, or Add/Rescan disabled: %+v", tab)
	}
	if d := tab.Disks[i]; !d.CanEdit || !d.CanDelete || !d.CanMove || d.DeleteTitle == "" {
		t.Fatalf("buttons of an unattached disk: %+v", d)
	}

	name, desc := prefix+"-renamed", "renamed by the lab test"
	if err := ops.SaveObject(c, "vdis", uuid, ops.ObjectChanges{Name: &name, Description: &desc}); err != nil {
		t.Fatal(err)
	}
	if p, err := ReadObjectProperties(c, "vdis", uuid); err != nil || p.Name != name || p.Description != desc {
		t.Fatalf("Properties after saving: %+v %v", p, err)
	}

	targets, err := MoveDiskTargets(c, uuid)
	if err != nil {
		t.Fatal(err)
	}
	reason := func(sr string) string {
		for _, u := range targets.SRs {
			if u.UUID == sr {
				return u.Reason
			}
		}
		return "not listed"
	}
	if targets.Migrate || reason(nfs) != "Current location" || reason(local) != "" {
		t.Fatalf("Move targets: %+v", targets)
	}
	moved, err := ops.MoveSRDisk(c, uuid, local)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := xenapi.VDI.GetByUUID(s, uuid); err == nil {
		t.Error("the moved disk is still on lab-nfs")
	}
	ref, _ := xenapi.VDI.GetByUUID(s, moved)
	v, err := xenapi.VDI.GetRecord(s, ref)
	if err != nil || srs[v.SR].UUID != local {
		t.Fatalf("after Move: %+v %v", v, err)
	}
	if v.NameLabel != name || v.VirtualSize != 32<<20 {
		t.Errorf("the copy is %q, %d bytes", v.NameLabel, v.VirtualSize)
	}
	back, err := ops.MoveSRDisk(c, moved, nfs)
	if err != nil {
		t.Fatal(err)
	}

	if err := ops.DeleteSRDisk(c, back); err != nil {
		t.Fatal(err)
	}
	if _, err := xenapi.VDI.GetByUUID(s, back); err == nil {
		t.Error("the disk is still there after Delete")
	}

	ps, err := PhysicalStorageTab(c, "pools", pool.UUID)
	if err != nil {
		t.Fatal(err)
	}
	if !slices.Contains(ps.Table.Keys, nfs) || !ps.CanNew || ps.Pool == "" {
		t.Fatalf("the pool's Storage tab: %+v", ps)
	}
	t.Logf("pool %q: trim %+v", ps.Pool, ps.SRs)
}

// TestLabDR formats the lab's iSCSI LUN as an SR (XCP_LAB_ISCSI, XCP_ALLOW_WRITE=1: XAPI keeps pool metadata on block
// SRs only, not on lab-nfs), turns DR on for it, reads the metadata back from it as the wizard does (the storage page,
// the current pool and XCP_LAB_VM found), runs a test failover's prechecks (the VM already exists here and is not
// older: refused) and the recovery's refusal of the current pool, then turns DR off and destroys the SR. Nothing is
// recovered.
func TestLabDR(t *testing.T) {
	c := labConnection(t)
	name, target := os.Getenv("XCP_LAB_VM"), os.Getenv("XCP_LAB_ISCSI")
	if name == "" || os.Getenv("XCP_ALLOW_WRITE") != "1" {
		t.Skip("XCP_LAB_VM and XCP_ALLOW_WRITE=1 not set")
	}
	s := c.Session
	vms, err := xenapi.VM.GetByNameLabel(s, name)
	if err != nil || len(vms) != 1 {
		t.Fatalf("VM %q: %v %v", name, vms, err)
	}
	vmUUID, _ := xenapi.VM.GetUUID(s, vms[0])
	_, pool, _, err := ops.PoolOf(s)
	if err != nil {
		t.Fatal(err)
	}
	before, err := ReadDRSRs(c)
	if err != nil {
		t.Fatal(err)
	}
	if slices.ContainsFunc(before.SRs, func(r DRSR) bool { return r.Enabled }) {
		t.Fatalf("DR is already on in the lab: %+v", before.SRs)
	}
	// XAPI keeps pool metadata on block SRs only (iSCSI, HBA): lab-nfs is refused (SR_OPERATION_NOT_SUPPORTED)
	if slices.ContainsFunc(before.SRs, func(r DRSR) bool { return r.Name == "lab-nfs" }) {
		t.Error("lab-nfs listed as able to hold the metadata")
	}
	if target == "" {
		t.Skip("XCP_LAB_ISCSI not set (DR needs an iSCSI SR: the lab's NFS SR cannot hold the metadata)")
	}
	iqns, err := ops.ScanIQNs(c, target, rules.DefaultISCSIPort, "", "")
	if err != nil || len(iqns) == 0 || iqns[0].IQN == "*" {
		t.Fatalf("%+v %v", iqns, err)
	}
	luns, err := ops.ScanLUNs(c, iqns[0], "", "")
	if err != nil || len(luns) == 0 {
		t.Fatalf("%+v %v", luns, err)
	}
	if used, _ := ops.LUNInUse(c, luns[0].ScsiID); used != "" {
		t.Skipf("the lab LUN is used by %s", used)
	}
	srUUID, err := ops.CreateSR(c, ops.NewSR{Kind: rules.SRIscsi, Name: "xcpng-center-go-test-dr", AutoDesc: true, IQN: iqns[0], LUN: luns[0]})
	if srUUID != "" {
		t.Cleanup(func() { // after DR is off (cleanups run last first): the SR goes, as TestLabISCSI's
			if ref, err := xenapi.SR.GetByUUID(s, srUUID); err == nil {
				pbds, _ := xenapi.SR.GetPBDs(s, ref)
				for _, p := range pbds {
					_ = xenapi.PBD.Unplug(s, p)
				}
				if err := xenapi.SR.Destroy(s, ref); err != nil {
					t.Errorf("cleanup: %v", err)
				}
			}
		})
	}
	if err != nil {
		t.Fatal(err)
	}
	if before, err = ReadDRSRs(c); err != nil {
		t.Fatal(err)
	}
	i := slices.IndexFunc(before.SRs, func(r DRSR) bool { return r.UUID == srUUID })
	if i < 0 || before.SRs[i].Enabled || before.SRs[i].HasMetadata {
		t.Fatalf("the iSCSI SR cannot hold the metadata: %+v", before.SRs)
	}
	dr := before.SRs[i]
	t.Cleanup(func() { // DR off everywhere, as found
		if err := ops.ConfigureDR(c, nil); err != nil {
			t.Errorf("cleanup: %v", err)
		}
		if after, err := ReadDRSRs(c); err != nil || slices.ContainsFunc(after.SRs, func(r DRSR) bool { return r.Enabled }) {
			t.Errorf("cleanup: DR still on: %+v %v", after.SRs, err)
		}
	})
	if before.Title != fmt.Sprintf("Configure DR on '%s'", pool.NameLabel) || before.MaxSRs != 8 {
		t.Errorf("dialog: %q %d", before.Title, before.MaxSRs)
	}
	if err := ops.ConfigureDR(c, []string{"no-such-sr"}); err == nil {
		t.Error("an unknown SR was accepted")
	}

	if err := ops.ConfigureDR(c, []string{dr.UUID}); err != nil {
		t.Fatal(err)
	}
	on, err := ReadDRSRs(c)
	if err != nil {
		t.Fatal(err)
	}
	if i := slices.IndexFunc(on.SRs, func(r DRSR) bool { return r.UUID == dr.UUID }); i < 0 || !on.SRs[i].Enabled || !on.SRs[i].HasMetadata {
		t.Fatalf("DR on lab-nfs: %+v", on.SRs)
	}

	// the wizard's VMs and vApps page skips a metadata VDI in use, as this pool's own is (attached for its redo log)
	pools, err := LoadDRMetadata(c, []string{dr.UUID})
	if err != nil || len(pools) != 0 {
		t.Fatalf("metadata in use: %+v %v", pools, err)
	}
	// read from that VDI as VdiLoadMetadataAction does, the metadata is this pool's, marked current, with the lab VM
	srRef, _ := xenapi.SR.GetByUUID(s, dr.UUID)
	vdis, err := ops.MetadataVDIs(s, nil)
	if err != nil || len(vdis) != 0 {
		t.Errorf("metadata VDIs on no SR: %v %v", vdis, err)
	}
	var md DRPool
	all, _ := xenapi.SR.GetVDIs(s, srRef)
	for _, v := range all {
		if typ, _ := xenapi.VDI.GetType(s, v); typ != xenapi.VdiTypeMetadata {
			continue
		}
		ok, err := ops.WithMetadata(s, v, func(_ *xenapi.Session, m ops.DRMetadata) error { md = drPool(m, pool.UUID); return nil })
		if err != nil || !ok {
			t.Fatalf("metadata: %v %v", ok, err)
		}
	}
	if md.UUID != pool.UUID || !md.Current || md.Name != pool.NameLabel {
		t.Fatalf("metadata: %+v", md)
	}
	if !slices.ContainsFunc(md.VMs, func(o DRObject) bool { return o.UUID == vmUUID && o.Name == name }) {
		t.Fatalf("%s is not in the metadata: %+v", name, md.VMs)
	}
	t.Logf("metadata of %s on VDI %s: %d vApps, %d VMs", md.Name, md.VDI, len(md.VApps), len(md.VMs))

	// a test failover's prechecks: the VM is here already and not older
	choice := []ops.DRChoice{{VDI: md.VDI, VMs: []string{vmUUID}}}
	checks, err := DRPrechecks(c, rules.DRTest, choice, nil)
	if err != nil {
		t.Fatal(err)
	}
	for _, ch := range checks {
		t.Logf("precheck: %+v", ch)
	}
	if !rules.DRBlocked(checks) || !slices.ContainsFunc(checks, func(ch rules.DRCheck) bool {
		return ch.Error && strings.Contains(ch.Text, fmt.Sprintf("The VM '%s' already exists in this pool and is not older", name))
	}) {
		t.Errorf("prechecks: %+v", checks)
	}
	if _, err := DRPrechecks(c, rules.DRTest, []ops.DRChoice{{VDI: md.VDI, VMs: []string{"no-such-vm"}}}, nil); err == nil {
		t.Error("a VM outside the metadata was accepted")
	}
	// the recovery refuses the current pool before recovering anything (md.Current was checked above)
	if rec, err := ops.RecoverDR(c, rules.DRTest, choice, rules.DRStartNone); err == nil || !strings.Contains(err.Error(), "is this pool") || !rec.Empty() {
		t.Errorf("recover onto the current pool: %+v %v", rec, err)
	}

	if err := ops.ConfigureDR(c, nil); err != nil {
		t.Fatal(err)
	}
	off, err := ReadDRSRs(c)
	if err != nil {
		t.Fatal(err)
	}
	if slices.ContainsFunc(off.SRs, func(r DRSR) bool { return r.Enabled }) {
		t.Errorf("DR still on: %+v", off.SRs)
	}
}

// TestLabPerfmonAlerts sets the Alerts page of the lab VM, lab-2 (with its control domain memory alert) and lab-nfs,
// reads them back, and puts every perfmon key back as it was.
func TestLabPerfmonAlerts(t *testing.T) {
	name := os.Getenv("XCP_LAB_VM")
	if name == "" || os.Getenv("XCP_ALLOW_WRITE") != "1" {
		t.Skip("XCP_LAB_VM and XCP_ALLOW_WRITE=1 not set")
	}
	c := labConnection(t)
	s := c.Session
	vms, _ := xenapi.VM.GetByNameLabel(s, name)
	hosts, _ := xenapi.Host.GetByNameLabel(s, "xcp-lab-2")
	srs, _ := xenapi.SR.GetByNameLabel(s, "lab-nfs")
	if len(vms) == 0 || len(hosts) == 0 || len(srs) == 0 {
		t.Fatal("lab objects missing")
	}
	dom0, _ := xenapi.Host.GetControlDomain(s, hosts[0])
	restore := func(get func() (map[string]string, error), rm func() error, add func(string) error) {
		other, err := get()
		if err != nil {
			t.Fatal(err)
		}
		was, had := other["perfmon"]
		t.Cleanup(func() {
			_ = rm()
			if had {
				if err := add(was); err != nil {
					t.Error("restore:", err)
				}
			}
		})
	}
	restore(func() (map[string]string, error) { return xenapi.VM.GetOtherConfig(s, vms[0]) },
		func() error { return xenapi.VM.RemoveFromOtherConfig(s, vms[0], "perfmon") },
		func(v string) error { return xenapi.VM.AddToOtherConfig(s, vms[0], "perfmon", v) })
	restore(func() (map[string]string, error) { return xenapi.VM.GetOtherConfig(s, dom0) },
		func() error { return xenapi.VM.RemoveFromOtherConfig(s, dom0, "perfmon") },
		func(v string) error { return xenapi.VM.AddToOtherConfig(s, dom0, "perfmon", v) })
	restore(func() (map[string]string, error) { return xenapi.Host.GetOtherConfig(s, hosts[0]) },
		func() error { return xenapi.Host.RemoveFromOtherConfig(s, hosts[0], "perfmon") },
		func(v string) error { return xenapi.Host.AddToOtherConfig(s, hosts[0], "perfmon", v) })
	restore(func() (map[string]string, error) { return xenapi.SR.GetOtherConfig(s, srs[0]) },
		func() error { return xenapi.SR.RemoveFromOtherConfig(s, srs[0], "perfmon") },
		func(v string) error { return xenapi.SR.AddToOtherConfig(s, srs[0], "perfmon", v) })

	vmUUID, _ := xenapi.VM.GetUUID(s, vms[0])
	hostUUID, _ := xenapi.Host.GetUUID(s, hosts[0])
	srUUID, _ := xenapi.SR.GetUUID(s, srs[0])
	for _, c2 := range []struct {
		kind, uuid string
		on         []ops.PerfmonChange
	}{
		{"vms", vmUUID, []ops.PerfmonChange{{Name: "cpu_usage", Level: 80, Period: 2}, {Name: "disk_usage", Level: 2048, Period: 5}}},
		{"hosts", hostUUID, []ops.PerfmonChange{{Name: "network_usage", Level: 500, Period: 3}, {Name: "mem_usage", Level: 90, Period: 4}}},
		{"srs", srUUID, []ops.PerfmonChange{{Name: "sr_io_throughput_total_per_host", Level: 4096, Period: 1}}},
	} {
		if err := ops.SavePerfmon(c, c2.kind, c2.uuid, c2.on, 30); err != nil {
			t.Fatal(c2.kind, err)
		}
		p, err := ReadPerfmon(c, c2.kind, c2.uuid)
		if err != nil || p.Interval != 30 {
			t.Fatal(c2.kind, p, err)
		}
		on := map[string]PerfmonSetting{}
		for _, a := range p.Alerts {
			if a.Enabled {
				on[a.Name] = a
			}
		}
		if len(on) != len(c2.on) {
			t.Fatal(c2.kind, p.Alerts)
		}
		for _, w := range c2.on {
			if a := on[w.Name]; a.Level != w.Level || a.Period != w.Period {
				t.Fatalf("%s %s: %+v", c2.kind, w.Name, a)
			}
		}
		t.Logf("%s: %+v", c2.kind, p)
	}
	if err := ops.SavePerfmon(c, "vms", vmUUID, nil, 7); err == nil {
		t.Fatal("an interval that is not a multiple of 5 was taken")
	}
}

// TestLabAddressSearch runs searches by address (IPAddressQuery, read-only): a partial address finds both lab servers,
// another subnet none.
func TestLabAddressSearch(t *testing.T) {
	c := labConnection(t)
	run := func(address string) map[string]bool {
		text := `<Search Name="by address" major_version="2" minor_version="1" show_expanded="yes"><Query>` +
			`<QueryScope><Server /></QueryScope><IPAddressQuery property="ip_address" contains="yes" address="` + address + `" />` +
			`</Query></Search>`
		nodes, err := RunSearch(c, "lab", text)
		if err != nil {
			t.Fatal(err)
		}
		found := map[string]bool{}
		var walk func([]*ViewNode)
		walk = func(ns []*ViewNode) {
			for _, n := range ns {
				if n.Object != nil {
					found[n.Object.Name] = true
				}
				walk(n.Children)
			}
		}
		walk(nodes)
		return found
	}
	if got := run("10.77.58.*"); !got["xcp-lab-1"] || !got["xcp-lab-2"] {
		t.Fatal(got)
	}
	if got := run("10.77.58.12"); !got["xcp-lab-2"] || got["xcp-lab-1"] {
		t.Fatal(got)
	}
	if got := run("192.0.2.*"); len(got) != 0 {
		t.Fatal(got)
	}
}

// TestLabDateSearch runs searches by start time (DatePropertyQuery, read-only) over the lab VM.
func TestLabDateSearch(t *testing.T) {
	name := os.Getenv("XCP_LAB_VM")
	if name == "" {
		t.Skip("XCP_LAB_VM not set")
	}
	c := labConnection(t)
	found := func(kind, day string) bool {
		text := `<Search Name="by date" major_version="2" minor_version="1" show_expanded="yes"><Query><QueryScope><VM /></QueryScope>` +
			`<DatePropertyQuery property="start_time" query="` + day + `" type="` + kind + `" /></Query></Search>`
		nodes, err := RunSearch(c, "lab", text)
		if err != nil {
			t.Fatal(err)
		}
		hit := false
		var walk func([]*ViewNode)
		walk = func(ns []*ViewNode) {
			for _, n := range ns {
				hit = hit || n.Object != nil && n.Object.Name == name
				walk(n.Children)
			}
		}
		walk(nodes)
		return hit
	}
	if !found("after", "20000101") || found("before", "20000101") || found("exact", "20000101") {
		t.Fatal("the lab VM's start time does not match as expected")
	}
	if err := SearchSupported(`<Search Name="x" major_version="2"><Query><QueryScope><VM /></QueryScope>` +
		`<DatePropertyQuery property="start_time" query="20000101" type="someday" /></Query></Search>`); err == nil {
		t.Fatal("an unknown date query type was taken")
	}
}

// TestLabCustomFieldSearch runs searches by custom field (CustomFieldQuery, CustomFieldDateQuery) over the lab VM,
// with the two values set on it for the test and removed after.
func TestLabCustomFieldSearch(t *testing.T) {
	name := os.Getenv("XCP_LAB_VM")
	if name == "" || os.Getenv("XCP_ALLOW_WRITE") != "1" {
		t.Skip("XCP_LAB_VM and XCP_ALLOW_WRITE=1 not set")
	}
	c := labConnection(t)
	refs, _ := xenapi.VM.GetByNameLabel(c.Session, name)
	if len(refs) == 0 {
		t.Fatal("no lab VM")
	}
	keys := map[string]string{"XenCenter.CustomFields.go-test-text": "Owner Alice", "XenCenter.CustomFields.go-test-date": "2026-01-15 10:00:00Z"}
	for k, v := range keys {
		_ = xenapi.VM.RemoveFromOtherConfig(c.Session, refs[0], k)
		if err := xenapi.VM.AddToOtherConfig(c.Session, refs[0], k, v); err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() { _ = xenapi.VM.RemoveFromOtherConfig(c.Session, refs[0], k) })
	}
	found := func(filter string) bool {
		text := `<Search Name="by field" major_version="2" minor_version="1" show_expanded="yes"><Query><QueryScope><VM /></QueryScope>` +
			filter + `</Query></Search>`
		nodes, err := RunSearch(c, "lab", text)
		if err != nil {
			t.Fatal(err)
		}
		hit := false
		var walk func([]*ViewNode)
		walk = func(ns []*ViewNode) {
			for _, n := range ns {
				hit = hit || n.Object != nil && n.Object.Name == name
				walk(n.Children)
			}
		}
		walk(nodes)
		return hit
	}
	text := func(q, kind string) string {
		return `<CustomFieldQuery query="` + q + `" type="` + kind + `"><CustomFieldDefinition name="go-test-text" type="String" defaultValue="" /></CustomFieldQuery>`
	}
	date := func(q, kind string) string {
		return `<CustomFieldDateQuery query="` + q + `" type="` + kind + `"><CustomFieldDefinition name="go-test-date" type="Date" defaultValue="" /></CustomFieldDateQuery>`
	}
	if !found(text("alice", "contains")) || found(text("bob", "contains")) || !found(text("owner", "startswith")) {
		t.Fatal("custom field text search")
	}
	if !found(date("20260115", "exact")) || found(date("20260116", "exact")) || !found(date("20260101", "after")) {
		t.Fatal("custom field date search")
	}
}

// TestLabObjectSearch runs searches by the objects a VM belongs to (XenModelObjectPropertyQuery on its server,
// XenModelObjectListContains(Name)Query on its networks), read-only.
func TestLabObjectSearch(t *testing.T) {
	name := os.Getenv("XCP_LAB_VM")
	if name == "" {
		t.Skip("XCP_LAB_VM not set")
	}
	c := labConnection(t)
	s := c.Session
	refs, _ := xenapi.VM.GetByNameLabel(s, name)
	if len(refs) == 0 {
		t.Fatal("no lab VM")
	}
	vm, _ := xenapi.VM.GetRecord(s, refs[0])
	hostUUID, _ := xenapi.Host.GetUUID(s, vm.ResidentOn)
	if len(vm.VIFs) == 0 {
		t.Skip("the lab VM has no VIF")
	}
	netRef, _ := xenapi.VIF.GetNetwork(s, vm.VIFs[0])
	netUUID, _ := xenapi.Network.GetUUID(s, netRef)
	netName, _ := xenapi.Network.GetNameLabel(s, netRef)
	found := func(filter string) bool {
		nodes, err := RunSearch(c, "lab", `<Search Name="by object" major_version="2" minor_version="1" show_expanded="yes"><Query>`+
			`<QueryScope><VM /></QueryScope>`+filter+`</Query></Search>`)
		if err != nil {
			t.Fatal(err)
		}
		hit := false
		var walk func([]*ViewNode)
		walk = func(ns []*ViewNode) {
			for _, n := range ns {
				hit = hit || n.Object != nil && n.Object.Name == name
				walk(n.Children)
			}
		}
		walk(nodes)
		return hit
	}
	if !found(`<XenModelObjectListContainsQuery property="host" contains="yes" uuid="`+hostUUID+`" />`) ||
		found(`<XenModelObjectListContainsQuery property="host" contains="no" uuid="`+hostUUID+`" />`) {
		t.Fatal("search by server")
	}
	if !found(`<XenModelObjectListContainsQuery property="networks" contains="yes" uuid="`+netUUID+`" />`) ||
		!found(`<XenModelObjectListContainsNameQuery property="networks" contains="yes" query="`+netName+`" type="exactmatch" />`) ||
		found(`<XenModelObjectListContainsNameQuery property="networks" contains="yes" query="no such network" type="exactmatch" />`) {
		t.Fatal("search by network")
	}
}
