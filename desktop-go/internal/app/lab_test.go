package app

import (
	"context"
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
