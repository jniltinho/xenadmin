using System;
using System.Linq;
using System.Threading;
using Microsoft.VisualStudio.TestTools.UnitTesting;
using XenAdmin.Actions;
using XenAdmin.Desktop;
using XenAPI;

namespace XenAdmin.Desktop.Tests
{
    /// <summary>
    /// Creates a throwaway diskless VM on XCP_HOST, drives it through the same actions the app uses, then destroys it.
    /// Runs only with XCP_ALLOW_WRITE=1.
    /// </summary>
    [TestClass]
    public class LifecycleIntegrationTests
    {
        private const string VmName = "xcpng-center-linux-test";

        [ClassInitialize]
        public static void Init(TestContext ctx) => ConnectionIntegrationTests.Init(ctx);

        private static void Run(AsyncAction action)
        {
            var done = new ManualResetEventSlim();
            action.Completed += _ => done.Set();
            action.RunAsync();
            Assert.IsTrue(done.Wait(TimeSpan.FromMinutes(3)), action.Title + " timed out");
        }

        private static VM WaitFor(XenAdmin.Network.XenConnection c, string vmRef, Func<VM, bool> cond)
        {
            var deadline = DateTime.Now.AddSeconds(60);
            VM vm;
            while (!((vm = c.Resolve(new XenRef<VM>(vmRef))) != null && cond(vm)) && DateTime.Now < deadline)
                Thread.Sleep(250);
            Assert.IsTrue(vm != null && cond(vm), "cache did not reflect the expected VM state");
            return vm;
        }

        /// <summary>
        /// Saves properties the way PropertiesDialog does (SaveChangesAction on a clone + page actions)
        /// and checks the result through the cache, then destroys the VM.
        /// </summary>
        [TestMethod, TestCategory("Integration"), TestCategory("Write")]
        public void PropertiesRoundTripOnThrowawayVm()
        {
            var host = Environment.GetEnvironmentVariable("XCP_HOST");
            if (string.IsNullOrEmpty(host) || Environment.GetEnvironmentVariable("XCP_ALLOW_WRITE") != "1")
                Assert.Inconclusive("XCP_HOST and XCP_ALLOW_WRITE=1 required");
            var c = ConnectionIntegrationTests.Connect(host);
            var session = c.DuplicateSession();
            string vmRef = null;
            try
            {
                vmRef = CreateThrowawayVm(c, session, VmName + "-props");
                var vm = WaitFor(c, vmRef, v => !v.is_a_template);

                var copy = (VM)vm.Clone();
                copy.name_label = VmName + "-renamed";
                copy.name_description = "edited by the Linux Properties dialog test";
                copy.SetAutoPowerOn(true);
                copy.SetBootOrder("DC");
                copy.SetVcpuWeight(1024);
                copy.SetCoresPerSocket(2);
                var actions = new System.Collections.Generic.List<AsyncAction>
                {
                    new SaveChangesAction(copy, true, vm),
                    new XenAdmin.Actions.GeneralEditPageAction(vm, copy, "/linux-test", new System.Collections.Generic.List<string> { "linux", "test" }, true),
                    new XenAdmin.Actions
.ChangeVCPUSettingsAction(copy, 2, 2)
                };
                var all = new MultipleAction(c, "Update properties", "Updating", "Updated", actions);
                vm.Locked = true; // as PropertiesDialog/ShowProperties do: SaveChanges requires the original locked
                Run(all);
                vm.Locked = false;
                Assert.IsTrue(all.Succeeded, all.Exception?.ToString());

                vm = WaitFor(c, vmRef, v => v.name_label == VmName + "-renamed" && v.VCPUs_max == 2 && v.VCPUs_at_startup == 2 && v.Path == "/linux-test");
                System.Console.WriteLine($"name={vm.name_label} vcpus={vm.VCPUs_max}/{vm.VCPUs_at_startup} path='{vm.Path}' folder='{(vm.other_config.TryGetValue("folder", out var f) ? f : null)}' tags={string.Join(",", vm.tags)} boot={vm.GetBootOrder()} weight={vm.GetVcpuWeight()} cores={vm.GetCoresPerSocket()} auto={vm.GetAutoPowerOn()}");
                Assert.AreEqual("edited by the Linux Properties dialog test", vm.name_description);
                Assert.IsTrue(vm.GetAutoPowerOn());
                Assert.AreEqual("DC", vm.GetBootOrder().ToUpper());
                Assert.AreEqual(1024, vm.GetVcpuWeight());
                Assert.AreEqual(2, vm.GetCoresPerSocket());
                CollectionAssert.AreEquivalent(new[] { "linux", "test" }, XenAdmin.Model.Tags.GetTags(vm));
                CollectionAssert.AreEquivalent(new[] { (2L, 1L), (1L, 2L) }, XenAdmin.Desktop.VmRules.ValidTopologies(2, vm.MaxCoresPerSocket()));
            }
            finally
            {
                if (vmRef != null)
                {
                    try { XenAdmin.Model.Folders.Unfolder(session, c.Resolve(new XenRef<VM>(vmRef))); } catch { }
                    VM.destroy(session, vmRef);
                }
                c.EndConnect();
            }
        }

        /// <summary>Take snapshot, revert, delete — the Snapshots tab actions — on a throwaway VM.</summary>
        [TestMethod, TestCategory("Integration"), TestCategory("Write")]
        public void SnapshotTakeRevertDelete()
        {
            var host = Environment.GetEnvironmentVariable("XCP_HOST");
            if (string.IsNullOrEmpty(host) || Environment.GetEnvironmentVariable("XCP_ALLOW_WRITE") != "1")
                Assert.Inconclusive("XCP_HOST and XCP_ALLOW_WRITE=1 required");
            var c = ConnectionIntegrationTests.Connect(host);
            var session = c.DuplicateSession();
            string vmRef = null;
            try
            {
                vmRef = CreateThrowawayVm(c, session, VmName + "-snap");
                // allowed_operations of a fresh VM reaches the cache in a later event
                var vm = WaitFor(c, vmRef, v => !v.is_a_template && v.allowed_operations.Contains(vm_operations.snapshot));

                var take = new VMSnapshotCreateAction(vm, "linux-test-snapshot", "taken by the Linux Snapshots tab test", SnapshotType.DISK, null);
                Run(take);
                Assert.IsTrue(take.Succeeded, take.Exception?.ToString());
                // the description is set after the snapshot task completes, so wait for it too
                vm = WaitFor(c, vmRef, v => v.snapshots.Count == 1 && c.Resolve(v.snapshots[0])?.name_description == "taken by the Linux Snapshots tab test");
                var snap = c.Resolve(vm.snapshots[0]);
                Assert.AreEqual("linux-test-snapshot", snap.name_label);
                Assert.AreEqual("taken by the Linux Snapshots tab test", snap.name_description);
                Assert.AreEqual(XenAdmin.Desktop.VmRules.SnapshotTypeText(snap), XenAdmin.Messages.DISKS_ONLY);
                Assert.IsTrue(XenAdmin.Desktop.VmRules.SnapshotSize(snap) >= 0);
                System.Console.WriteLine($"snapshot {snap.name_label} size={XenAdmin.Desktop.VmRules.SnapshotSize(snap)} allowed={string.Join(",", snap.allowed_operations)}");

                var revert = new VMSnapshotRevertAction(snap);
                Run(revert);
                Assert.IsTrue(revert.Succeeded, revert.Exception?.ToString());

                var delete = new VMSnapshotDeleteAction(snap);
                Run(delete);
                Assert.IsTrue(delete.Succeeded, delete.Exception?.ToString());
                WaitFor(c, vmRef, v => v.snapshots.Count == 0);
            }
            finally
            {
                if (vmRef != null)
                {
                    foreach (var s in VM.get_snapshots(session, vmRef))
                        try { VM.destroy(session, s.opaque_ref); } catch { }
                    VM.destroy(session, vmRef);
                }
                c.EndConnect();
            }
        }

        /// <summary>Networking tab: add an interface, edit its QoS limit and MAC, remove it (VIF actions of NetworkList).</summary>
        [TestMethod, TestCategory("Integration"), TestCategory("Write")]
        public void NetworkTabInterfaces()
        {
            var host = Environment.GetEnvironmentVariable("XCP_HOST");
            if (string.IsNullOrEmpty(host) || Environment.GetEnvironmentVariable("XCP_ALLOW_WRITE") != "1")
                Assert.Inconclusive("XCP_HOST and XCP_ALLOW_WRITE=1 required");
            var c = ConnectionIntegrationTests.Connect(host);
            var session = c.DuplicateSession();
            string vmRef = null;
            try
            {
                vmRef = CreateThrowawayVm(c, session, VmName + "-net");
                var vm = WaitFor(c, vmRef, v => !v.is_a_template);
                Assert.IsTrue(NetworkRules.CanAdd(vm) && !NetworkRules.AtMaxInterfaces(vm));
                var net = NetworkRules.Networks(c, vm).First();

                var add = new CreateVIFAction(vm, NetworkRules.Settings(vm, null, net, "", VIF.GetDeviceId(vm), false, ""));
                Run(add);
                Assert.IsTrue(add.Succeeded, add.Exception?.ToString());
                vm = WaitFor(c, vmRef, v => v.VIFs.Count == 1 && c.Resolve(v.VIFs[0]) != null);
                // a second interface: the edit below targets device 1, which UpdateVIFAction destroys and recreates
                var second = new CreateVIFAction(vm, NetworkRules.Settings(vm, null, net, "", VIF.GetDeviceId(vm), false, ""));
                Run(second);
                Assert.IsTrue(second.Succeeded, second.Exception?.ToString());
                vm = WaitFor(c, vmRef, v => v.VIFs.Count == 2 && c.ResolveAll(v.VIFs).Count == 2);
                var vif = c.ResolveAll(vm.VIFs).Single(v => v.device == "1");
                Assert.AreEqual(net.opaque_ref, vif.network.opaque_ref);
                Assert.IsFalse(string.IsNullOrEmpty(vif.MAC), "auto-generated MAC");
                Assert.IsTrue(NetworkRules.CanEdit(vm, vif) && NetworkRules.CanRemove(vm, vif), "a detached VIF of a halted VM can be edited and removed");
                Assert.IsFalse(NetworkRules.CanToggle(vm, vif), "a halted VM cannot plug interfaces");

                const string mac = "02:00:5e:10:00:01";
                Assert.IsNull(NetworkRules.Validate(net, true, mac, true, "1024"));
                Assert.IsTrue(NetworkRules.Changed(vif, net, mac, 1, true, "1024"));
                var edit = new UpdateVIFAction(vm, vif, NetworkRules.Settings(vm, vif, net, mac, 1, true, "1024"));
                Run(edit);
                Assert.IsTrue(edit.Succeeded, edit.Exception?.ToString());
                vm = WaitFor(c, vmRef, v => v.VIFs.Count == 2 && c.ResolveAll(v.VIFs).Any(e => e.device == "1" && e.MAC == mac && NetworkRules.Limit(e) == "1024"));
                vif = c.ResolveAll(vm.VIFs).Single(v => v.device == "1");
                CollectionAssert.AreEqual(new[] { "1", mac, "1024" }, Details.Networking(vm).Rows[1].Take(3).ToArray(), "device 1 keeps its row after the edit");
                Assert.IsFalse(NetworkRules.Changed(vif, net, mac, 1, true, "1024"), "unchanged edit is a cancel");

                var remove = new DeleteVIFAction(vif);
                Run(remove);
                Assert.IsTrue(remove.Succeeded, remove.Exception?.ToString());
                WaitFor(c, vmRef, v => v.VIFs.Count == 1);
            }
            finally
            {
                if (vmRef != null)
                {
                    foreach (var v in VM.get_VIFs(session, vmRef))
                        try { VIF.destroy(session, v.opaque_ref); } catch { }
                    VM.destroy(session, vmRef);
                }
                c.EndConnect();
            }
        }

        /// <summary>SR Storage tab: add an unattached disk on the SR, rescan, then delete it (SrStoragePage).</summary>
        [TestMethod, TestCategory("Integration"), TestCategory("Write")]
        public void SrStorageTabAddRescanDelete()
        {
            var host = Environment.GetEnvironmentVariable("XCP_HOST");
            if (string.IsNullOrEmpty(host) || Environment.GetEnvironmentVariable("XCP_ALLOW_WRITE") != "1")
                Assert.Inconclusive("XCP_HOST and XCP_ALLOW_WRITE=1 required");
            var c = ConnectionIntegrationTests.Connect(host);
            var session = c.DuplicateSession();
            const string diskName = VmName + "-sr-disk";
            string vdiRef = null;
            try
            {
                foreach (var old in c.Cache.VDIs.Where(v => v.name_label == diskName))
                    VDI.destroy(session, old.opaque_ref);
                var sr = StorageRules.NewDiskTargets(c, null).First(s => !s.IsToolsSR() && StorageRules.CanAddDisk(s));
                Assert.IsTrue(StorageRules.CanRescan(sr, false).Enabled);
                Assert.IsFalse(StorageRules.CanRescan(sr, true).Enabled);

                var vdi = new VDI { Connection = c, SR = new XenRef<SR>(sr), virtual_size = 64L * 1024 * 1024, name_label = diskName, name_description = "", type = vdi_type.user, sharable = false, read_only = false };
                var create = new CreateDiskAction(vdi);
                Run(create);
                Assert.IsTrue(create.Succeeded, create.Exception?.ToString());
                var deadline = DateTime.Now.AddSeconds(60);
                VDI disk;
                while ((disk = StorageRules.SrDisks(sr).FirstOrDefault(v => v.name_label == diskName)) == null && DateTime.Now < deadline)
                    Thread.Sleep(250);
                Assert.IsNotNull(disk, "the new disk is listed on the SR Storage tab");
                vdiRef = disk.opaque_ref;
                Assert.IsTrue(Details.Storage(sr).Rows.Any(r => r[0] == diskName));
                Assert.IsTrue(StorageRules.CanEditDisk(sr, disk));

                var rescan = new SrRefreshAction(sr);
                Run(rescan);
                Assert.IsTrue(rescan.Succeeded, rescan.Exception?.ToString());

                deadline = DateTime.Now.AddSeconds(30);
                while (!StorageRules.CanDelete(disk = c.Resolve(new XenRef<VDI>(vdiRef)), true) && DateTime.Now < deadline)
                    Thread.Sleep(250);
                Assert.IsTrue(StorageRules.CanDelete(disk, true), "an unattached user disk can be deleted");
                Assert.AreEqual(Messages.MESSAGEBOX_DELETE_VD, StorageRules.DeleteConfirmation(disk).Text);
                var delete = StorageRules.DeleteActions(disk).Single();
                Run(delete);
                Assert.IsTrue(delete.Succeeded, delete.Exception?.ToString());
                vdiRef = null;
            }
            finally
            {
                if (vdiRef != null)
                    try { VDI.destroy(session, vdiRef); } catch { }
                c.EndConnect();
            }
        }

        /// <summary>Storage tab: DVD drive + ISO + eject, add disk, detach, re-attach, delete.</summary>
        [TestMethod, TestCategory("Integration"), TestCategory("Write")]
        public void StorageTabDiskAndDvd()
        {
            var host = Environment.GetEnvironmentVariable("XCP_HOST");
            if (string.IsNullOrEmpty(host) || Environment.GetEnvironmentVariable("XCP_ALLOW_WRITE") != "1")
                Assert.Inconclusive("XCP_HOST and XCP_ALLOW_WRITE=1 required");
            var c = ConnectionIntegrationTests.Connect(host);
            var session = c.DuplicateSession();
            string vmRef = null, vdiRef = null;
            try
            {
                vmRef = CreateThrowawayVm(c, session, VmName + "-disk");
                var vm = WaitFor(c, vmRef, v => !v.is_a_template);

                // DVD drive
                Run(new CreateCdDriveAction(vm));
                vm = WaitFor(c, vmRef, v => c.ResolveAll(v.VBDs).Any(b => b.IsCDROM()));
                var cd = c.ResolveAll(vm.VBDs).First(b => b.IsCDROM());
                var iso = c.Cache.SRs.Where(sr => sr.content_type == SR.Content_Type_ISO).SelectMany(sr => c.ResolveAll(sr.VDIs)).FirstOrDefault();
                if (iso != null)
                {
                    var insert = new ChangeVMISOAction(c, vm, iso, cd);
                    Run(insert);
                    Assert.IsTrue(insert.Succeeded, insert.Exception?.ToString());
                    WaitFor(c, vmRef, v => c.ResolveAll(v.VBDs).Any(b => b.IsCDROM() && !b.empty && b.VDI.opaque_ref == iso.opaque_ref));
                    var eject = new ChangeVMISOAction(c, vm, null, c.ResolveAll(vm.VBDs).First(b => b.IsCDROM()));
                    Run(eject);
                    Assert.IsTrue(eject.Succeeded, eject.Exception?.ToString());
                    WaitFor(c, vmRef, v => c.ResolveAll(v.VBDs).Any(b => b.IsCDROM() && b.empty));
                    System.Console.WriteLine($"DVD: inserted and ejected '{iso.Name()}'");
                }

                // add a 64 MB disk the way the Add Virtual Disk dialog does
                var home = vm.Home() ?? c.Cache.Hosts[0];
                var sr0 = c.Cache.SRs.First(sr => sr.SupportsVdiCreate() && sr.HasPBDs() && !sr.IsBroken() && sr.content_type != SR.Content_Type_ISO && sr.CanBeSeenFrom(home) && !sr.IsToolsSR());
                var vdi = new VDI { Connection = c, SR = new XenRef<SR>(sr0), virtual_size = 64L * 1024 * 1024, name_label = "linux-test-disk", name_description = "", type = vdi_type.user };
                vdi.SetVmHint(vm.uuid);
                var vbd = new VBD { Connection = c, device = "", empty = false, type = vbd_type.Disk, mode = vbd_mode.RW, unpluggable = true };
                vbd.SetIsOwner(true);
                var add = new MultipleAction(c, "Add disk", "", "", new System.Collections.Generic.List<AsyncAction> { new CreateDiskAction(vdi, vbd, vm), new VbdCreateAndPlugAction(vm, vbd, vdi.Name(), false) }, false, false, true);
                Run(add);
                Assert.IsTrue(add.Succeeded, add.Exception?.ToString());
                vm = WaitFor(c, vmRef, v => c.ResolveAll(v.VBDs).Any(b => b.type == vbd_type.Disk && c.Resolve(b.VDI)?.name_label == "linux-test-disk"));
                var disk = c.ResolveAll(vm.VBDs).Select(b => c.Resolve(b.VDI)).First(v => v?.name_label == "linux-test-disk");
                vdiRef = disk.opaque_ref;
                var diskVbd = c.ResolveAll(vm.VBDs).First(b => b.VDI.opaque_ref == vdiRef);
                Assert.IsFalse(XenAdmin.Desktop.StorageRules.CanActivate(vm, diskVbd), "halted VM: nothing to plug");
                Assert.IsTrue(XenAdmin.Desktop.StorageRules.CanDetach(vm, disk));
                Assert.IsTrue(XenAdmin.Desktop.StorageRules.CanDelete(disk), "unattached user disk on a halted VM can be deleted");
                System.Console.WriteLine($"disk {disk.Name()} {disk.SizeText()} on {sr0.Name()} as {diskVbd.userdevice}");

                // detach, then attach again as the Attach Disk dialog does
                var detach = new DetachVirtualDiskAction(disk, vm, true);
                Run(detach);
                Assert.IsTrue(detach.Succeeded, detach.Exception?.ToString());
                vm = WaitFor(c, vmRef, v => c.ResolveAll(v.VBDs).All(b => b.VDI.opaque_ref != vdiRef));
                var ud = VM.get_allowed_VBD_devices(session, vmRef)[0];
                var again = new VBD { VDI = new XenRef<VDI>(vdiRef), VM = new XenRef<VM>(vmRef), bootable = ud == "0", device = "", empty = false, userdevice = ud, type = vbd_type.Disk, mode = vbd_mode.RW, unpluggable = true };
                var attach = new VbdCreateAndPlugAction(vm, again, disk.Name(), false);
                Run(attach);
                Assert.IsTrue(attach.Succeeded, attach.Exception?.ToString());
                WaitFor(c, vmRef, v => c.ResolveAll(v.VBDs).Any(b => b.VDI.opaque_ref == vdiRef));

                var destroy = new DestroyDiskAction(c.Resolve(new XenRef<VDI>(vdiRef)));
                Run(destroy);
                Assert.IsTrue(destroy.Succeeded, destroy.Exception?.ToString());
                WaitFor(c, vmRef, v => c.ResolveAll(v.VBDs).All(b => b.VDI.opaque_ref != vdiRef));
                Assert.IsNull(c.Resolve(new XenRef<VDI>(vdiRef)) is VDI gone && !gone.managed ? null : c.Resolve(new XenRef<VDI>(vdiRef)), "disk destroyed");
                vdiRef = null;
            }
            finally
            {
                if (vdiRef != null)
                    try { VDI.destroy(session, vdiRef); } catch { }
                if (vmRef != null)
                    VM.destroy(session, vmRef);
                c.EndConnect();
            }
        }

        /// <summary>New VM wizard (CreateVMAction) -> reboot -> shut down -> fast clone + full copy -> delete all with their disks.</summary>
        [TestMethod, TestCategory("Integration"), TestCategory("Write")]
        public void NewVmCloneCopyDelete()
        {
            var host = Environment.GetEnvironmentVariable("XCP_HOST");
            if (string.IsNullOrEmpty(host) || Environment.GetEnvironmentVariable("XCP_ALLOW_WRITE") != "1")
                Assert.Inconclusive("XCP_HOST and XCP_ALLOW_WRITE=1 required");
            var c = ConnectionIntegrationTests.Connect(host);
            var session = c.DuplicateSession();
            var name = VmName + "-newvm";
            try
            {
                CleanupNamed(c, session, VmName + "-newvm");
                // a default HVM template that provisions a disk (e.g. a Debian/Ubuntu template)
                var template = c.Cache.VMs.Where(v => v.is_a_template && v.DefaultTemplate() && v.IsHVM() && v.ProvisionXml()?.ChildNodes.Count > 0)
                    .OrderBy(v => v.name_label.Contains("Debian") ? 0 : 1).First();
                System.Console.WriteLine("template: " + template.name_label);
                var disks = XenAdmin.Desktop.NewVmDefaults.DisksFor(c, template, name, null);
                Assert.IsTrue(disks.Count > 0, "default template provisions a disk");
                foreach (var d in disks) d.Disk.virtual_size = 64L * 1024 * 1024; // small, as requested
                var vifs = XenAdmin.Desktop.NewVmDefaults.VifsFor(c, template);
                var mem = Math.Max(256L * 1024 * 1024, template.memory_static_min);
                var create = XenAdmin.Desktop.NewVmDefaults.Build(c, template, name, "created by the Linux New VM wizard test", XenAdmin.Actions.VMActions
.InstallMethod.CD, null, "",
                    VmBootMode.Bios, null, 1, 1, mem, disks, vifs, false);
                Run(create);
                Assert.IsTrue(create.Succeeded, create.Exception?.ToString());
                var vm = WaitFor(c, create.VM.opaque_ref, v => v.name_label == name);
                // "Start the new VM automatically" runs this same action after creation
                var start = VmOperations.Create(vm, VmOp.Start);
                Run(start);
                Assert.IsTrue(start.Succeeded, "start: " + start.Exception?.Message);
                vm = WaitFor(c, create.VM.opaque_ref, v => v.power_state == vm_power_state.Running);
                Assert.AreEqual(name, vm.name_label);
                Assert.AreEqual(1, c.ResolveAll(vm.VBDs).Count(b => b.type == vbd_type.Disk), "one provisioned disk");
                Assert.AreEqual(vifs.Count, vm.VIFs.Count);
                System.Console.WriteLine($"created {vm.name_label}: {vm.VCPUs_max} vCPU, {vm.memory_static_max / 1048576} MB, disks={vm.VBDs.Count}, vifs={vm.VIFs.Count}, running on {c.Resolve(vm.resident_on)?.Name()}");

                // reboot (no guest tools in an empty VM, so the forced variant) and shut down
                Assert.IsTrue(VmOperations.IsAllowed(vm, VmOp.ForceReboot));
                var reboot = VmOperations.Create(vm, VmOp.ForceReboot);
                Run(reboot);
                Assert.IsTrue(reboot.Succeeded, reboot.Exception?.ToString());
                vm = WaitFor(c, vm.opaque_ref, v => v.power_state == vm_power_state.Running);
                var stop = VmOperations.Create(vm, VmOp.ForceShutdown);
                Run(stop);
                vm = WaitFor(c, vm.opaque_ref, v => v.power_state == vm_power_state.Halted);
                Assert.IsTrue(XenAdmin.Desktop.VmRules.CanCopy(vm) && XenAdmin.Desktop.VmRules.CanDeleteVm(vm));

                // Copy VM: fast clone and full copy
                var clone = new XenAdmin.Actions.VMActions.VMCloneAction(vm, name + "-clone", "fast clone");
                Run(clone);
                Assert.IsTrue(clone.Succeeded, clone.Exception?.ToString());
                var sr = c.Resolve(c.Resolve(c.ResolveAll(vm.VBDs).First(b => b.type == vbd_type.Disk).VDI).SR);
                var copy = new XenAdmin.Actions.VMActions.VMCopyAction(vm, vm.GetStorageHost(false), sr, name + "-copy", "full copy");
                Run(copy);
                Assert.IsTrue(copy.Succeeded, copy.Exception?.ToString());
                var copies = new[] { name + "-clone", name + "-copy" }.Select(n => WaitForVm(c, n)).ToList();
                foreach (var cp in copies)
                    Assert.AreEqual(1, c.ResolveAll(cp.VBDs).Count(b => b.type == vbd_type.Disk), cp.name_label + " has its own disk");
                System.Console.WriteLine("copies: " + string.Join(", ", copies.Select(v => v.name_label)));

                // Delete VM with its disks, as ConfirmVMDeleteDialog does by default (owned disks checked)
                foreach (var target in copies.Append(vm))
                {
                    var owned = c.ResolveAll(target.VBDs).Where(b => b.type == vbd_type.Disk && b.GetIsOwner()).ToList();
                    var vdis = owned.Select(b => b.VDI.opaque_ref).ToList();
                    var del = new XenAdmin.Actions.VMActions.VMDestroyAction(target, owned, new System.Collections.Generic.List<VM>());
                    Run(del);
                    Assert.IsTrue(del.Succeeded, del.Exception?.ToString());
                    var deadline = DateTime.Now.AddSeconds(30);
                    while ((c.Resolve(new XenRef<VM>(target.opaque_ref)) != null || vdis.Any(r => c.Resolve(new XenRef<VDI>(r)) != null)) && DateTime.Now < deadline)
                        Thread.Sleep(250);
                    Assert.IsNull(c.Resolve(new XenRef<VM>(target.opaque_ref)), target.name_label + " deleted");
                    Assert.IsTrue(vdis.All(r => c.Resolve(new XenRef<VDI>(r)) == null), target.name_label + " disks deleted");
                }
            }
            finally
            {
                CleanupNamed(c, session, VmName + "-newvm");
                c.EndConnect();
            }
        }

        private static VM WaitForVm(XenAdmin.Network.XenConnection c, string name)
        {
            var deadline = DateTime.Now.AddSeconds(60);
            VM vm;
            while ((vm = c.Cache.VMs.FirstOrDefault(v => v.name_label == name && !v.is_a_template)) == null && DateTime.Now < deadline)
                Thread.Sleep(250);
            Assert.IsNotNull(vm, name + " not found");
            return vm;
        }

        /// <summary>Removes leftovers (VMs whose name starts with the prefix, plus their owned disks).</summary>
        private static void CleanupNamed(XenAdmin.Network.XenConnection c, Session session, string prefix)
        {
            foreach (var old in c.Cache.VMs.Where(v => v.name_label.StartsWith(prefix) && !v.is_a_template && !v.is_a_snapshot).ToList())
            {
                try { if (old.power_state != vm_power_state.Halted) VM.hard_shutdown(session, old.opaque_ref); } catch { }
                foreach (var b in c.ResolveAll(old.VBDs).Where(b => b.type == vbd_type.Disk && b.GetIsOwner()))
                    try { VDI.destroy(session, b.VDI.opaque_ref); } catch { }
                try { VM.destroy(session, old.opaque_ref); } catch { }
            }
        }

        private static string CreateThrowawayVm(XenAdmin.Network.XenConnection c, Session session, string name)
        {
            foreach (var old in c.Cache.VMs.Where(v => v.name_label.StartsWith(VmName) && !v.is_a_template))
            {
                if (old.power_state != vm_power_state.Halted) VM.hard_shutdown(session, old.opaque_ref);
                VM.destroy(session, old.opaque_ref);
            }
            var template = c.Cache.VMs.First(v => v.is_a_template && v.name_label == "Other install media");
            var vmRef = VM.clone(session, template.opaque_ref, name).opaque_ref;
            VM.set_is_a_template(session, vmRef, false);
            const long mem = 256L * 1024 * 1024;
            VM.set_memory_limits(session, vmRef, mem, mem, mem, mem);
            VM.set_VCPUs_max(session, vmRef, 1);
            VM.set_VCPUs_at_startup(session, vmRef, 1);
            return vmRef;
        }

        [TestMethod, TestCategory("Integration"), TestCategory("Write")]
        public void PowerCycleThrowawayVm()
        {
            var host = Environment.GetEnvironmentVariable("XCP_HOST");
            if (string.IsNullOrEmpty(host) || Environment.GetEnvironmentVariable("XCP_ALLOW_WRITE") != "1")
                Assert.Inconclusive("XCP_HOST and XCP_ALLOW_WRITE=1 required");

            var c = ConnectionIntegrationTests.Connect(host);
            var session = c.DuplicateSession();
            string vmRef = null;
            try
            {
                // leftovers from an aborted run
                foreach (var old in c.Cache.VMs.Where(v => v.name_label == VmName && !v.is_a_template))
                {
                    if (old.power_state != vm_power_state.Halted) VM.hard_shutdown(session, old.opaque_ref);
                    VM.destroy(session, old.opaque_ref);
                }

                var template = c.Cache.VMs.First(v => v.is_a_template && v.name_label == "Other install media");
                vmRef = VM.clone(session, template.opaque_ref, VmName).opaque_ref;
                VM.set_is_a_template(session, vmRef, false);
                const long mem = 256L * 1024 * 1024;
                VM.set_memory_limits(session, vmRef, mem, mem, mem, mem);
                VM.set_VCPUs_max(session, vmRef, 1);
                VM.set_VCPUs_at_startup(session, vmRef, 1);

                var vm = WaitFor(c, vmRef, v => v.power_state == vm_power_state.Halted && !v.is_a_template);
                Assert.IsTrue(VmOperations.IsAllowed(vm, VmOp.Start));
                Assert.IsFalse(VmOperations.IsAllowed(vm, VmOp.ForceShutdown), "halted VM cannot be force-shutdown");

                Run(VmOperations.Create(vm, VmOp.Start));
                vm = WaitFor(c, vmRef, v => v.power_state == vm_power_state.Running);
                Assert.IsNotNull(c.Resolve(vm.resident_on), "running VM must have a resident host");
                System.Console.WriteLine("allowed while running: " + string.Join(",", vm.allowed_operations));

                // expected failure: a diskless VM without guest tools cannot be suspended
                if (!VmOperations.IsAllowed(vm, VmOp.Suspend))
                {
                    var suspend = VmOperations.Create(vm, VmOp.Suspend);
                    Run(suspend);
                    Assert.IsFalse(suspend.Succeeded, "suspend should fail without guest tools");
                    Assert.IsInstanceOfType(suspend.Exception, typeof(Failure));
                    System.Console.WriteLine("expected suspend failure: " + suspend.Exception.Message);
                }

                Assert.IsTrue(VmOperations.IsAllowed(vm, VmOp.ForceReboot));
                var reboot = VmOperations.Create(vm, VmOp.ForceReboot);
                Run(reboot);
                Assert.IsTrue(reboot.Succeeded, reboot.Exception?.Message);
                vm = WaitFor(c, vmRef, v => v.power_state == vm_power_state.Running);

                var shutdown = VmOperations.Create(vm, VmOp.ForceShutdown);
                Run(shutdown);
                Assert.IsTrue(shutdown.Succeeded, shutdown.Exception?.Message);
                WaitFor(c, vmRef, v => v.power_state == vm_power_state.Halted);
            }
            finally
            {
                if (vmRef != null)
                {
                    try { if (VM.get_power_state(session, vmRef) != vm_power_state.Halted) VM.hard_shutdown(session, vmRef); } catch { }
                    VM.destroy(session, vmRef);
                }
                c.EndConnect();
            }
        }
    }
}
