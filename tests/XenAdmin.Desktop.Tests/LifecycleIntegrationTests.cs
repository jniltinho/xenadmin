using System;
using System.Collections.Generic;
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

        private static string WaitForRef(XenAdmin.Network.XenConnection c, string name)
        {
            var deadline = DateTime.Now.AddSeconds(60);
            VM vm;
            while ((vm = c.Cache.VMs.FirstOrDefault(v => v.name_label == name)) == null && DateTime.Now < deadline)
                Thread.Sleep(250);
            Assert.IsNotNull(vm, name + " in the cache");
            return vm.opaque_ref;
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

        /// <summary>Memory tab: Edit (Memory Settings) on a halted throwaway VM, 256 MB to 384 MB.</summary>
        [TestMethod, TestCategory("Integration"), TestCategory("Write")]
        public void MemorySettingsOnThrowawayVm()
        {
            var host = Environment.GetEnvironmentVariable("XCP_HOST");
            if (string.IsNullOrEmpty(host) || Environment.GetEnvironmentVariable("XCP_ALLOW_WRITE") != "1")
                Assert.Inconclusive("XCP_HOST and XCP_ALLOW_WRITE=1 required");
            var c = ConnectionIntegrationTests.Connect(host);
            var session = c.DuplicateSession();
            string vmRef = null;
            try
            {
                vmRef = CreateThrowawayVm(c, session, VmName + "-memory");
                var vm = WaitFor(c, vmRef, v => !v.is_a_template && v.memory_static_max == 256L << 20);
                Assert.IsTrue(MemoryRules.CanEdit(vm), "a halted VM can be edited");
                Assert.IsNull(MemoryRules.Confirmation(vm, 384L << 20), "no restart warning for a halted VM");
                const long size = 384L << 20;
                var (dmin, dmax, smax) = (size, size, size); // fixed memory: what the Basic controls send
                Assert.IsNull(MemoryRules.Validate(vm, dmin, dmax, smax));
                Assert.IsNotNull(MemoryRules.Validate(vm, dmin, size + (1L << 20), smax), "dynamic maximum above the static maximum");
                var settings = MemoryRules.Settings(vm, dmin, dmax, smax);
                Run(MemoryRules.ChangeAction(vm, settings.DynamicMin, settings.DynamicMax, settings.StaticMax));
                WaitFor(c, vmRef, v => v.memory_static_max == size && v.memory_dynamic_max == size && v.memory_dynamic_min == size);
            }
            finally
            {
                if (vmRef != null)
                    VM.destroy(session, vmRef);
                c.EndConnect();
            }
        }

        /// <summary>VM menu: Convert to Template (ConvertVMToTemplateCommand) on a throwaway VM, then the template is destroyed.</summary>
        [TestMethod, TestCategory("Integration"), TestCategory("Write")]
        public void ConvertThrowawayVmToTemplate()
        {
            var host = Environment.GetEnvironmentVariable("XCP_HOST");
            if (string.IsNullOrEmpty(host) || Environment.GetEnvironmentVariable("XCP_ALLOW_WRITE") != "1")
                Assert.Inconclusive("XCP_HOST and XCP_ALLOW_WRITE=1 required");
            var c = ConnectionIntegrationTests.Connect(host);
            var session = c.DuplicateSession();
            string vmRef = null;
            try
            {
                vmRef = CreateThrowawayVm(c, session, VmName + "-template");
                var vm = WaitFor(c, vmRef, v => !v.is_a_template && VmRules.CanConvertToTemplate(v) && VmRules.CanTakeSnapshot(v));
                Run(VmRules.ConvertToTemplateAction(vm));
                vm = WaitFor(c, vmRef, v => v.is_a_template && v.other_config.TryGetValue("instant", out var i) && i == "true");
                Assert.IsFalse(VmRules.CanConvertToTemplate(vm), "a template cannot be converted again");
                Assert.IsFalse(VmRules.CanTakeSnapshot(vm), "templates are not snapshotted from the VM menu");

                // Templates menu on the new user template: Create VM From Selection, Copy, Delete Template
                Assert.IsTrue(VmRules.IsUserTemplate(vm), "a converted VM is a user template");
                Assert.IsTrue(VmRules.CanDeleteTemplate(vm));
                Assert.IsTrue(VmRules.CanCopyTemplate(vm), "allowed: " + string.Join(",", vm.allowed_operations));
                Run(VmRules.CopyAction(vm, false, null, VmName + "-template-copy", ""));
                var copy = WaitFor(c, WaitForRef(c, VmName + "-template-copy"), v => v.is_a_template);
                VM.destroy(session, copy.opaque_ref);
            }
            finally
            {
                if (vmRef != null)
                    VM.destroy(session, vmRef);
                foreach (var left in c.Cache.VMs.Where(v => v.name_label == VmName + "-template-copy").ToList())
                    try { VM.destroy(session, left.opaque_ref); } catch { }
                c.EndConnect();
            }
        }

        /// <summary>VM menu Export (XVA, verified) of a halted throwaway VM, then File > Import of that file to the default SR.</summary>
        [TestMethod, TestCategory("Integration"), TestCategory("Write")]
        public void ExportImportXvaOnThrowawayVm()
        {
            var host = Environment.GetEnvironmentVariable("XCP_HOST");
            if (string.IsNullOrEmpty(host) || Environment.GetEnvironmentVariable("XCP_ALLOW_WRITE") != "1")
                Assert.Inconclusive("XCP_HOST and XCP_ALLOW_WRITE=1 required");
            var c = ConnectionIntegrationTests.Connect(host);
            var session = c.DuplicateSession();
            var file = System.IO.Path.Combine(System.IO.Path.GetTempPath(), VmName + "-xva.xva");
            string vmRef = null;
            try
            {
                vmRef = CreateThrowawayVm(c, session, VmName + "-xva");
                var vm = WaitFor(c, vmRef, v => !v.is_a_template && XvaRules.CanExport(v));
                Assert.IsNull(XvaRules.ExportBlocker(vm), "a halted VM exports");
                var export = XvaRules.ExportAction(vm, file, true);
                Run(export);
                Assert.IsTrue(export.Succeeded, "export: " + export.Exception?.Message);
                Assert.IsTrue(new System.IO.FileInfo(file).Length > 0, "the XVA was written");
                Assert.IsNull(XvaRules.FileError(file));

                var srs = XvaRules.Targets(c, null);
                Assert.IsTrue(srs.Count > 0, "an SR can hold the import");
                Assert.IsNull(XvaRules.TargetError(c, null, srs[0]));
                var import = XvaRules.ImportAction(c, null, file, srs[0]);
                Run(import);
                Assert.IsTrue(import.Succeeded, "import: " + import.Exception?.Message);
                Assert.IsNotNull(import.VM, "the imported VM");
                Assert.AreNotEqual(vmRef, import.VM.opaque_ref);
                Assert.IsFalse(import.VM.is_a_template);
                var startAction = XvaRules.StartAfterImport(import);
                Assert.IsInstanceOfType(startAction, typeof(XenAdmin.Actions.VMActions.VMStartAction), "a halted import is started, not resumed");
            }
            finally
            {
                CleanupNamed(c, session, VmName + "-xva");
                try { System.IO.File.Delete(file); } catch { }
                c.EndConnect();
            }
        }

        /// <summary>VM > Manage vTPMs on a halted UEFI throwaway VM: add one, read its properties, remove it.</summary>
        [TestMethod, TestCategory("Integration"), TestCategory("Write")]
        public void VtpmAddRemoveOnThrowawayVm()
        {
            var host = Environment.GetEnvironmentVariable("XCP_HOST");
            if (string.IsNullOrEmpty(host) || Environment.GetEnvironmentVariable("XCP_ALLOW_WRITE") != "1")
                Assert.Inconclusive("XCP_HOST and XCP_ALLOW_WRITE=1 required");
            var c = ConnectionIntegrationTests.Connect(host);
            var session = c.DuplicateSession();
            string vmRef = null;
            try
            {
                vmRef = CreateThrowawayVm(c, session, VmName + "-vtpm");
                var bios = WaitFor(c, vmRef, v => v.IsHVM());
                Assert.IsFalse(VtpmRules.CanManage(bios), "BIOS VMs have no vTPM");
                if (XenAdmin.Core.Helpers.XapiEqualOrGreater_22_26_0(c) && !XenAdmin.Core.Helpers.FeatureForbidden(c, Host.RestrictVtpm))
                    Assert.IsNotNull(VtpmRules.Blocker(bios), "the reason names UEFI");
                var boot = new System.Collections.Generic.Dictionary<string, string>(bios.HVM_boot_params) { ["firmware"] = "uefi" };
                VM.set_HVM_boot_params(session, vmRef, boot);
                var vm = WaitFor(c, vmRef, v => v.IsDefaultBootModeUefi());
                if (!VtpmRules.CanManage(vm))
                    Assert.Inconclusive("vTPMs are not available on this server");
                Assert.AreEqual(0, VtpmRules.Vtpms(vm).Count);
                vm = WaitFor(c, vmRef, v => VtpmRules.AddBlocker(v) == null);
                Run(VtpmRules.AddAction(vm));
                vm = WaitFor(c, vmRef, v => v.VTPMs.Count == 1);
                var deadline = DateTime.Now.AddSeconds(60);
                VTPM vtpm;
                while ((vtpm = VtpmRules.Vtpms(vm).FirstOrDefault()) == null && DateTime.Now < deadline)
                    Thread.Sleep(250);
                Assert.IsNotNull(vtpm);
                Assert.IsNotNull(VtpmRules.Properties(vtpm).Protected);
                while (VtpmRules.RemoveBlocker(vm, vtpm = c.Resolve(new XenRef<VTPM>(vtpm.opaque_ref))) != null && DateTime.Now < deadline)
                    Thread.Sleep(250);
                Assert.IsNull(VtpmRules.RemoveBlocker(vm, vtpm), "a halted VM's vTPM can be removed");
                Run(VtpmRules.RemoveAction(vm, vtpm));
                WaitFor(c, vmRef, v => v.VTPMs.Count == 0);
            }
            finally
            {
                if (vmRef != null)
                {
                    try { foreach (var t in VM.get_VTPMs(session, vmRef)) VTPM.destroy(session, t.opaque_ref); } catch { }
                    try { VM.destroy(session, vmRef); } catch { }
                }
                c.EndConnect();
            }
        }

        /// <summary>Search > Save to the pool, listed in Saved Searches, then Delete (a throwaway gui_config entry).</summary>
        [TestMethod, TestCategory("Integration"), TestCategory("Write")]
        public void SaveAndDeleteSearchInPool()
        {
            var host = Environment.GetEnvironmentVariable("XCP_HOST");
            if (string.IsNullOrEmpty(host) || Environment.GetEnvironmentVariable("XCP_ALLOW_WRITE") != "1")
                Assert.Inconclusive("XCP_HOST and XCP_ALLOW_WRITE=1 required");
            var c = ConnectionIntegrationTests.Connect(host);
            var registered = !XenAdmin.ConnectionsManager.XenConnections.Contains(c);
            if (registered) XenAdmin.ConnectionsManager.XenConnections.Add(c);
            XenAdmin.XenSearch.Search search = null;
            try
            {
                XenAdmin.XenSearch.Search.InitSearch();
                search = SearchRules.Build(VmName + "-search", XenAdmin.XenSearch.ObjectTypes.VM, "", null, null);
                Run(SearchRules.SaveAction(search, c));
                var deadline = DateTime.Now.AddSeconds(60);
                while (NavigationViews.SavedSearches().All(s => s.UUID != search.UUID) && DateTime.Now < deadline)
                    Thread.Sleep(250);
                var saved = NavigationViews.SavedSearches().FirstOrDefault(s => s.UUID == search.UUID);
                Assert.IsNotNull(saved, "the saved search is listed");
                Assert.IsTrue(SearchRules.CanDelete(saved));
                Run(SearchRules.DeleteAction(saved, c));
                while (NavigationViews.SavedSearches().Any(s => s.UUID == search.UUID) && DateTime.Now < deadline)
                    Thread.Sleep(250);
                Assert.IsFalse(NavigationViews.SavedSearches().Any(s => s.UUID == search.UUID), "deleted");
            }
            finally
            {
                if (search?.UUID != null)
                    try { foreach (var p in c.Cache.Pools) Pool.remove_from_gui_config(c.DuplicateSession(), p.opaque_ref, XenAdmin.Core.BrandManager.ExtensionSearch + ".Search-" + search.UUID); } catch { }
                if (registered) XenAdmin.ConnectionsManager.XenConnections.Remove(c);
                c.EndConnect();
            }
        }

        /// <summary>Objects, Organization Views and Saved Searches (read only) on XCP_HOST.</summary>
        [TestMethod, TestCategory("Integration")]
        public void NavigationViewsReadOnly()
        {
            var host = Environment.GetEnvironmentVariable("XCP_HOST");
            if (string.IsNullOrEmpty(host))
                Assert.Inconclusive("XCP_HOST required");
            var c = ConnectionIntegrationTests.Connect(host);
            var registered = !XenAdmin.ConnectionsManager.XenConnections.Contains(c);
            if (registered) XenAdmin.ConnectionsManager.XenConnections.Add(c);
            try
            {
                var objects = NavigationViews.Build(NavView.Objects);
                Assert.AreEqual(Messages.VIEW_OBJECTS, objects.Name);
                var all = ResourceTree.Flatten(objects).Where(n => n.XenObject != null).Select(n => n.XenObject.opaque_ref).ToList();
                foreach (var h in c.Cache.Hosts)
                    CollectionAssert.Contains(all, h.opaque_ref, "Objects lists every server");
                Assert.IsTrue(objects.Children.All(g => g.GroupKey != null && g.XenObject == null), "the first level is the object types");
                var keys = ResourceTree.Flatten(objects).Select(n => n.Key).ToList();
                Assert.AreEqual(keys.Count, keys.Distinct().Count(), "node keys are unique");
                Assert.IsTrue(NavigationViews.Objects(objects).Any(), "the root's Search tab lists objects");

                foreach (var view in new[] { NavView.Tags, NavView.Folders, NavView.CustomFields, NavView.Vapps })
                    Assert.AreEqual(NavigationViews.RootOf(view).Name, NavigationViews.Build(view).Name);
                var tagged = c.Cache.VMs.Where(v => v.tags != null && v.tags.Length > 0).ToList();
                if (tagged.Count > 0)
                    Assert.IsTrue(NavigationViews.Build(NavView.Tags).Children.Count > 0, "a tag group per tag");

                XenAdmin.XenSearch.Search.InitSearch();
                var searches = NavigationViews.SavedSearches();
                Assert.IsTrue(searches.Any(s => s.DefaultSearch), "the default searches");
                var byOs = NavigationViews.Build(NavView.SavedSearch, searches.First(s => s.DefaultSearch));
                Assert.IsTrue(byOs.IsRoot);
            }
            finally
            {
                if (registered) XenAdmin.ConnectionsManager.XenConnections.Remove(c);
                c.EndConnect();
            }
        }

        /// <summary>New Network > SR-IOV (read only): the option is offered with a NIC to use, or shown disabled with the reason.</summary>
        [TestMethod, TestCategory("Integration")]
        public void SriovOptionReadOnly()
        {
            var host = Environment.GetEnvironmentVariable("XCP_HOST");
            if (string.IsNullOrEmpty(host))
                Assert.Inconclusive("XCP_HOST required");
            var c = ConnectionIntegrationTests.Connect(host);
            try
            {
                var (visible, blocker) = NetworkRules.SriovOption(c);
                Assert.AreEqual(XenAdmin.Core.Helpers.KolkataOrGreater(c), visible);
                var nics = c.Cache.Hosts.SelectMany(NetworkRules.SriovNics).ToList();
                if (visible && blocker == null)
                    Assert.IsTrue(nics.Count > 0, "SR-IOV offered: some NIC can enable it");
                if (nics.Count == 0 && visible)
                    Assert.IsNotNull(blocker, "no NIC can enable SR-IOV: the option says why");
                Assert.IsFalse(string.IsNullOrEmpty(NetworkRules.SriovWarning().Text));
            }
            finally
            {
                c.EndConnect();
            }
        }

        /// <summary>GPU and USB tabs (read only): the rules run on every server and the pool of XCP_HOST.</summary>
        [TestMethod, TestCategory("Integration")]
        public void GpuAndUsbTabsReadOnly()
        {
            var host = Environment.GetEnvironmentVariable("XCP_HOST");
            if (string.IsNullOrEmpty(host))
                Assert.Inconclusive("XCP_HOST required");
            var c = ConnectionIntegrationTests.Connect(host);
            try
            {
                var pool = XenAdmin.Core.Helpers.GetPoolOfOne(c);
                foreach (IXenObject o in c.Cache.Hosts.Cast<IXenObject>().Append(pool))
                {
                    var groups = HardwareRules.Groups(o);
                    Assert.AreEqual(HardwareRules.Pgpus(o).Count, groups.Sum(g => g.Count), "every GPU is in one row");
                    foreach (var g in groups)
                        Assert.IsTrue(HardwareRules.TypeRows(g).Count > 0, "a GPU row lists its vGPU types");
                    if (groups.Count == 0)
                        Assert.IsFalse(string.IsNullOrEmpty(HardwareRules.NoGpuText(o)));
                    Assert.IsFalse(string.IsNullOrEmpty(HardwareRules.PolicyText(HardwareRules.PlacementPolicy(c))));
                    Assert.IsNull(HardwareRules.PlacementAction(c, HardwareRules.PlacementPolicy(c)), "unchanged policy: nothing to save");
                    if (o is Host h)
                    {
                        Assert.AreEqual(h.PUSBs.Count > 0 && !XenAdmin.Core.Helpers.FeatureForbidden(c, Host.RestrictUsbPassthrough), Details.Tabs(h).Contains("USB"));
                        foreach (var p in HardwareRules.Pusbs(h))
                            Assert.IsNotNull(HardwareRules.PassthroughButton(p).Text);
                    }
                }
            }
            finally
            {
                c.EndConnect();
            }
        }

        /// <summary>vApps: New vApp wizard (create with a VM, startup order), Assign to vApp, Properties VMs, Delete.</summary>
        [TestMethod, TestCategory("Integration"), TestCategory("Write")]
        public void VappCreateAssignDeleteOnThrowawayVms()
        {
            var host = Environment.GetEnvironmentVariable("XCP_HOST");
            if (string.IsNullOrEmpty(host) || Environment.GetEnvironmentVariable("XCP_ALLOW_WRITE") != "1")
                Assert.Inconclusive("XCP_HOST and XCP_ALLOW_WRITE=1 required");
            var c = ConnectionIntegrationTests.Connect(host);
            var session = c.DuplicateSession();
            const string name = VmName + "-vapp";
            string aRef = null, bRef = null;
            try
            {
                foreach (var old in c.Cache.VM_appliances.Where(a => a.name_label == name).ToList())
                    VM_appliance.destroy(session, old.opaque_ref);
                aRef = CreateThrowawayVm(c, session, VmName + "-vapp-a");
                bRef = VM.clone(session, aRef, VmName + "-vapp-b").opaque_ref;
                var a = WaitFor(c, aRef, v => VappRules.CanAssign(v));
                var b = WaitFor(c, bRef, v => VappRules.CanAssign(v));
                Assert.IsTrue(VappRules.CanManage(c));
                Assert.IsNotNull(VappRules.NameError(" "));

                Run(VappRules.CreateAction(c, name, "throwaway", new List<VM> { a }));
                var deadline = DateTime.Now.AddSeconds(60);
                VM_appliance app;
                while ((app = c.Cache.VM_appliances.FirstOrDefault(x => x.name_label == name)) == null && DateTime.Now < deadline)
                    Thread.Sleep(250);
                Assert.IsNotNull(app, "the vApp");
                a = WaitFor(c, aRef, v => VappRules.Of(v)?.opaque_ref == app.opaque_ref);

                var changed = VappRules.ChangedStartup(new[] { (a, 3L, 7L) });
                Run(VappRules.StartupAction(c, changed));
                WaitFor(c, aRef, v => v.order == 3 && v.start_delay == 7);

                // Assign to vApp: b joins (a stays), then choosing the checked vApp takes b out again
                Assert.IsNull(VappRules.ChangeWarning(new[] { b }, app, name), "b has no vApp yet");
                Run(VappRules.AssignAction(app, new List<VM> { b }));
                WaitFor(c, bRef, v => VappRules.Of(v)?.opaque_ref == app.opaque_ref);
                Assert.AreEqual(2, c.Resolve(new XenRef<VM_appliance>(app.opaque_ref)).VMs.Count);
                Run(VappRules.AssignAction(c.Resolve(new XenRef<VM_appliance>(app.opaque_ref)), new List<VM> { c.Resolve(new XenRef<VM>(bRef)) }));
                WaitFor(c, bRef, v => VappRules.Of(v) == null);

                // Properties > Virtual Machines: only b
                Run(VappRules.SetMembersAction(c.Resolve(new XenRef<VM_appliance>(app.opaque_ref)), new[] { c.Resolve(new XenRef<VM>(bRef)) }));
                WaitFor(c, aRef, v => VappRules.Of(v) == null);
                WaitFor(c, bRef, v => VappRules.Of(v)?.opaque_ref == app.opaque_ref);

                var current = c.Resolve(new XenRef<VM_appliance>(app.opaque_ref));
                StringAssert.Contains(VappRules.DeleteConfirmation(new[] { current }), name);
                Run(VappRules.DeleteAction(c, new List<VM_appliance> { current }));
                WaitFor(c, bRef, v => VappRules.Of(v) == null);
                deadline = DateTime.Now.AddSeconds(60);
                while (c.Resolve(new XenRef<VM_appliance>(app.opaque_ref)) != null && DateTime.Now < deadline)
                    Thread.Sleep(250);
                Assert.IsNull(c.Resolve(new XenRef<VM_appliance>(app.opaque_ref)), "the vApp is gone");
                Assert.IsNotNull(c.Resolve(new XenRef<VM>(bRef)), "its VMs stay");
            }
            finally
            {
                foreach (var left in c.Cache.VM_appliances.Where(x => x.name_label == name).ToList())
                    try { VM_appliance.destroy(session, left.opaque_ref); } catch { }
                foreach (var r in new[] { aRef, bRef }.Where(r => r != null))
                    try { VM.destroy(session, r); } catch { }
                c.EndConnect();
            }
        }

        /// <summary>File > Import of an OVF appliance (one VM, a small VHD disk, one interface) to the pool, then cleanup.</summary>
        [TestMethod, TestCategory("Integration"), TestCategory("Write")]
        public void ImportOvfApplianceOnThrowawayVm()
        {
            var host = Environment.GetEnvironmentVariable("XCP_HOST");
            if (string.IsNullOrEmpty(host) || Environment.GetEnvironmentVariable("XCP_ALLOW_WRITE") != "1")
                Assert.Inconclusive("XCP_HOST and XCP_ALLOW_WRITE=1 required");
            var dir = System.IO.Path.Combine(System.IO.Path.GetTempPath(), Guid.NewGuid().ToString());
            System.IO.Directory.CreateDirectory(dir);
            const string name = VmName + "-ovf";
            var c = ConnectionIntegrationTests.Connect(host);
            var session = c.DuplicateSession();
            try
            {
                var ovf = LinuxAppTests.WriteSampleAppliance(dir, name);
                var vhd = System.IO.Path.Combine(dir, "disk1.vhd");

                Assert.IsTrue(OvfRules.IsAppliance(ovf));
                var (package, error, _) = OvfRules.Open(ovf);
                Assert.IsNull(error, error);
                var vms = OvfRules.Vms(package);
                Assert.AreEqual(1, vms.Count);
                Assert.AreEqual(name, vms[0].Name);
                Assert.AreEqual(0, OvfRules.Eulas(package).Count);

                var pool = XenAdmin.Core.Helpers.GetPoolOfOne(c);
                var sr = XvaRules.Targets(c, null).First();
                var net = OvfRules.Networks(c, null).First();
                Assert.IsNull(OvfRules.SpaceError(sr, OvfRules.RequiredSpace(package)), "8 MB fits");
                Assert.AreEqual(8UL << 20, OvfRules.RequiredSpace(package));
                var mappings = OvfRules.Mappings(package, pool, sr, net);
                Assert.AreEqual(1, mappings.Values.Single().Storage.Count, "the disk is mapped");
                var action = OvfRules.ImportAction(c, package, mappings, false, false);
                Run(action);
                Assert.IsTrue(action.Succeeded, "import: " + action.Exception?.Message);
                var vm = WaitForVm(c, name);
                Assert.AreEqual(1, c.ResolveAll(vm.VBDs).Count(b => b.type == vbd_type.Disk), "the VM has the disk");
                CleanupNamed(c, session, name);
                foreach (var app in c.Cache.VM_appliances.Where(a => a.name_label == name).ToList())
                    VM_appliance.destroy(session, app.opaque_ref);

                // the same appliance as an OVA (a tar with the descriptor first)
                var ova = System.IO.Path.Combine(System.IO.Path.GetTempPath(), Guid.NewGuid() + ".ova");
                try
                {
                    using (var tar = new System.Formats.Tar.TarWriter(System.IO.File.Create(ova), System.Formats.Tar.TarEntryFormat.Ustar, false))
                    {
                        tar.WriteEntry(ovf, System.IO.Path.GetFileName(ovf));
                        tar.WriteEntry(vhd, "disk1.vhd");
                    }
                    var (ovaPackage, ovaError, _) = OvfRules.Open(ova);
                    Assert.IsNull(ovaError, ovaError);
                    var ovaAction = OvfRules.ImportAction(c, ovaPackage, OvfRules.Mappings(ovaPackage, pool, sr, net), false, false);
                    Run(ovaAction);
                    Assert.IsTrue(ovaAction.Succeeded, "OVA import: " + ovaAction.Exception?.Message);
                    WaitForVm(c, name);
                }
                finally
                {
                    System.IO.File.Delete(ova);
                }
            }
            finally
            {
                CleanupNamed(c, session, name);
                // ImportApplianceAction also creates a vApp named after the package
                foreach (var app in c.Cache.VM_appliances.Where(a => a.name_label == name).ToList())
                    try { VM_appliance.destroy(session, app.opaque_ref); } catch { }
                try { System.IO.Directory.Delete(dir, true); } catch { }
                c.EndConnect();
            }
        }

        /// <summary>Properties pages Start Options, Advanced and Custom Fields on a throwaway VM (their save actions).</summary>
        [TestMethod, TestCategory("Integration"), TestCategory("Write")]
        public void StartOptionsAdvancedAndCustomFieldsOnThrowawayVm()
        {
            var host = Environment.GetEnvironmentVariable("XCP_HOST");
            if (string.IsNullOrEmpty(host) || Environment.GetEnvironmentVariable("XCP_ALLOW_WRITE") != "1")
                Assert.Inconclusive("XCP_HOST and XCP_ALLOW_WRITE=1 required");
            // as the app's Program: events reach the watcher that keeps the custom field definitions current
            XenAdmin.InvokeHelper.Initialize(new ImmediateSynchronizer());
            XenAdmin.OtherConfigAndTagsWatcher.RegisterEventHandlers();
            XenAdmin.CustomFields.CustomFieldsManager.GetCustomFields();
            var c = ConnectionIntegrationTests.Connect(host);
            var registered = !XenAdmin.ConnectionsManager.XenConnections.Contains(c);
            if (registered) XenAdmin.ConnectionsManager.XenConnections.Add(c);
            var session = c.DuplicateSession();
            var field = new XenAdmin.CustomFields.CustomFieldDefinition(VmName + "-field", XenAdmin.CustomFields.CustomFieldDefinition.Types.String);
            string vmRef = null;
            try
            {
                vmRef = CreateThrowawayVm(c, session, VmName + "-props");
                var vm = WaitFor(c, vmRef, v => !v.is_a_template);

                // Start Options (HA off): SetVMStartupOptionsAction
                Run(new SetVMStartupOptionsAction(c, new System.Collections.Generic.Dictionary<VM, VMStartupOptions> { [vm] = new VMStartupOptions(5, 30, vm.HARestartPriority()) }, true));
                WaitFor(c, vmRef, v => v.order == 5 && v.start_delay == 30);

                // Advanced, halted VM: the clone's multiplier goes through SaveChangesAction
                var clone = (VM)vm.Clone();
                clone.HVM_shadow_multiplier = 4.0;
                vm.Locked = true;
                Run(new SaveChangesAction(clone, true, vm));
                WaitFor(c, vmRef, v => Math.Abs(v.HVM_shadow_multiplier - 4.0) < 0.001);

                // Custom Fields: a field for the pool (stored in its gui_config), a value on the VM, then the field is removed
                Run(new AddCustomFieldAction(c, field));
                bool Defined() => XenAdmin.Core.Helpers.GetPoolOfOne(c).gui_config.TryGetValue(XenAdmin.CustomFields.CustomFieldsManager.CUSTOM_FIELD_BASE_KEY, out var xml) && xml.Contains(field.Name);
                var deadline = DateTime.Now.AddSeconds(30);
                while (!Defined() && DateTime.Now < deadline)
                    Thread.Sleep(250);
                Assert.IsTrue(Defined(), "the field is defined on the pool");
                var def = field;
                Run(new SaveCustomFieldsAction(c.Resolve(new XenRef<VM>(vmRef)), new System.Collections.Generic.List<XenAdmin.CustomFields.CustomField> { new XenAdmin.CustomFields.CustomField(def, "owner: test") }, true));
                deadline = DateTime.Now.AddSeconds(30);
                while (!"owner: test".Equals(XenAdmin.CustomFields.CustomFieldsManager.GetCustomFieldValue(c.Resolve(new XenRef<VM>(vmRef)), def)) && DateTime.Now < deadline)
                    Thread.Sleep(250);
                Assert.AreEqual("owner: test", XenAdmin.CustomFields.CustomFieldsManager.GetCustomFieldValue(c.Resolve(new XenRef<VM>(vmRef)), def));
                deadline = DateTime.Now.AddSeconds(30);
                while (XenAdmin.CustomFields.CustomFieldsManager.GetCustomFields(c).All(d => d.Name != field.Name) && DateTime.Now < deadline)
                    Thread.Sleep(250);
                Assert.IsTrue(XenAdmin.CustomFields.CustomFieldsManager.GetCustomFields(c).Any(d => d.Name == field.Name), "the Custom Fields page lists the new field");
                Run(new RemoveCustomFieldAction(c, def));
                deadline = DateTime.Now.AddSeconds(30);
                while (Defined() && DateTime.Now < deadline)
                    Thread.Sleep(250);
                Assert.IsFalse(Defined(), "the field is removed from the pool");
            }
            finally
            {
                try { new RemoveCustomFieldAction(c, field).RunSync(session); } catch { }
                if (vmRef != null)
                    try { VM.destroy(session, vmRef); } catch { }
                if (registered) XenAdmin.ConnectionsManager.XenConnections.Remove(c);
                c.EndConnect();
            }
        }

        /// <summary>Host Networking tab: create a single-server private network (New Network wizard) and remove it.</summary>
        [TestMethod, TestCategory("Integration"), TestCategory("Write")]
        public void HostNetworkPrivateCreateRemove()
        {
            var host = Environment.GetEnvironmentVariable("XCP_HOST");
            if (string.IsNullOrEmpty(host) || Environment.GetEnvironmentVariable("XCP_ALLOW_WRITE") != "1")
                Assert.Inconclusive("XCP_HOST and XCP_ALLOW_WRITE=1 required");
            var c = ConnectionIntegrationTests.Connect(host);
            var session = c.DuplicateSession();
            const string name = VmName + "-private-net";
            string netRef = null;
            try
            {
                foreach (var old in c.Cache.Networks.Where(n => n.name_label == name))
                    XenAPI.Network.destroy(session, old.opaque_ref);
                var coordinator = NetworkRules.WizardHost(c.Cache.Hosts[0]);
                Assert.IsNotNull(coordinator);
                var create = new NetworkAction(c, NetworkRules.NewNetwork(name, "created by the Linux New Network wizard test", false, null), true);
                Run(create);
                Assert.IsTrue(create.Succeeded, create.Exception?.ToString());
                var deadline = DateTime.Now.AddSeconds(60);
                XenAPI.Network net;
                while ((net = NetworkRules.HostNetworks(c).FirstOrDefault(n => n.name_label == name)) == null && DateTime.Now < deadline)
                    Thread.Sleep(250);
                Assert.IsNotNull(net, "the new network is listed on the Networking tab");
                netRef = net.opaque_ref;
                Assert.IsFalse(net.GetAutoPlug());
                Assert.IsTrue(Details.Networking(coordinator).Rows.Any(r => r[0] == name));
                Assert.IsTrue(NetworkRules.CanEditNetwork(coordinator, net));
                Assert.IsTrue(NetworkRules.CanRemoveNetwork(coordinator, net).Enabled, "a private network without physical NICs can be removed");

                var remove = new NetworkAction(c, net, false);
                Run(remove);
                Assert.IsTrue(remove.Succeeded, remove.Exception?.ToString());
                netRef = null;
            }
            finally
            {
                if (netRef != null)
                    try { XenAPI.Network.destroy(session, netRef); } catch { }
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

        /// <summary>Runs InvokeHelper callbacks on the calling thread (the app uses its UI dispatcher).</summary>
        private sealed class ImmediateSynchronizer : System.ComponentModel.ISynchronizeInvoke
        {
            public bool InvokeRequired => false;
            public IAsyncResult BeginInvoke(Delegate method, object[] args) { method.DynamicInvoke(args); return System.Threading.Tasks.Task.CompletedTask; }
            public object EndInvoke(IAsyncResult result) => null;
            public object Invoke(Delegate method, object[] args) => method.DynamicInvoke(args);
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
