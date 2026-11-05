package app

import (
	"context"
	"os"
	"path/filepath"
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
	if err != nil || len(kinds) != 4 || kinds[0].DefaultName == "" {
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
