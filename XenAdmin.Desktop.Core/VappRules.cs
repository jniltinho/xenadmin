using System.Collections.Generic;
using System.Linq;
using XenAdmin.Actions;
using XenAdmin.Core;
using XenAdmin.Network;
using XenAPI;

namespace XenAdmin.Desktop
{
    /// <summary>
    /// vApps: VMAppliancesDialog (Manage vApps), NewVMApplianceWizard, VappStartCommand / VappShutDownCommand and
    /// AssignGroupToolStripMenuItem&lt;VM_appliance&gt; (VM > Assign to vApp), with the vApp strings of VMGroup&lt;T&gt;.
    /// Exporting and importing a vApp (OVF) are not ported.
    /// </summary>
    public static class VappRules
    {
        /// <summary>VMGroupCommand&lt;VM_appliance&gt;: a connected pool or standalone server.</summary>
        public static bool CanManage(IXenConnection c) => c != null && c.IsConnected && Helpers.GetPoolOfOne(c) != null;

        /// <summary>AssignVMsToGroup.CanRun: a real, unlocked VM where vApps are not restricted by the license.</summary>
        public static bool CanAssign(VM vm) =>
            vm != null && vm.IsRealVm() && !vm.Locked && vm.Connection.IsConnected && !Helpers.FeatureForbidden(vm.Connection, Host.RestrictVMAppliances);

        public static List<VM_appliance> Appliances(IXenConnection c) => c.Cache.VM_appliances.OrderBy(a => a).ToList();

        /// <summary>NewVMGroupVMsPage: the real VMs of the pool (hidden VMs stay hidden), in the usual VM order.</summary>
        public static List<VM> Candidates(IXenConnection c) => c.Cache.VMs.Where(v => v.IsRealVm() && v.Show(false)).OrderBy(v => v).ToList();

        public static VM_appliance Of(VM vm) => vm.Connection.Resolve(vm.appliance);

        public static bool CanStart(VM_appliance a) => a != null && a.allowed_operations.Contains(vm_appliance_operation.start);

        public static bool CanShutDown(VM_appliance a) =>
            a != null && (a.allowed_operations.Contains(vm_appliance_operation.clean_shutdown)
                          || a.allowed_operations.Contains(vm_appliance_operation.hard_shutdown));

        public static AsyncAction StartAction(VM_appliance a) => new StartApplianceAction(a, false);

        public static AsyncAction ShutDownAction(VM_appliance a) => new ShutDownApplianceAction(a);

        /// <summary>NewVMApplianceNamePage.EnableNext.</summary>
        public static string NameError(string name) => string.IsNullOrWhiteSpace(name) ? "Enter a name for the vApp." : null;

        /// <summary>VMAppliancesDialog.toolStripButtonDelete_Click: the confirmation for these vApps and their VMs.</summary>
        public static string DeleteConfirmation(IList<VM_appliance> apps)
        {
            var vms = apps.Sum(a => a.VMs.Count);
            if (apps.Count == 1)
                return vms == 1
                    ? string.Format(Messages.CONFIRM_DELETE_VM_APPLIANCE_1, apps[0].Name().Ellipsise(120), vms)
                    : string.Format(vms == 0 ? Messages.CONFIRM_DELETE_VM_APPLIANCE_0 : Messages.CONFIRM_DELETE_VM_APPLIANCE, apps[0].Name().Ellipsise(120), vms);
            return vms == 1
                ? string.Format(Messages.CONFIRM_DELETE_VM_APPLIANCES_1, vms)
                : string.Format(vms == 0 ? Messages.CONFIRM_DELETE_VM_APPLIANCES_0 : Messages.CONFIRM_DELETE_VM_APPLIANCES, vms);
        }

        public static AsyncAction DeleteAction(IXenConnection c, List<VM_appliance> apps) => new DestroyVMApplianceAction(c, apps);

        /// <summary>
        /// AssignGroupToVMCommand.ChangesOK: the warning when some VMs already belong to another vApp (null: no warning).
        /// </summary>
        public static string ChangeWarning(IList<VM> vms, VM_appliance target, string targetName)
        {
            var moving = vms.Where(vm => Of(vm) is VM_appliance old && (target == null || old.opaque_ref != target.opaque_ref)).ToList();
            if (moving.Count == 0)
                return null;
            return moving.Count == 1
                ? string.Format(Messages.CHANGE_VMS_APPLIANCE_WARNING, moving[0].Name().Ellipsise(250), Of(moving[0]).Name().Ellipsise(250), targetName.Ellipsise(250))
                : string.Format(Messages.CHANGE_VMS_APPLIANCES_WARNING, targetName.Ellipsise(250));
        }

        /// <summary>
        /// AssignGroupToVMCommand.RunCore: one VM already in that vApp leaves it (the menu item is checked);
        /// otherwise the VMs join it, keeping its other members.
        /// </summary>
        public static AsyncAction AssignAction(VM_appliance a, IList<VM> vms)
        {
            if (vms.Count == 1 && a.VMs.Any(r => r.opaque_ref == vms[0].opaque_ref))
                return new RemoveVMsFromVMApplianceAction(a, new List<XenRef<VM>> { new XenRef<VM>(vms[0].opaque_ref) });
            var refs = vms.Select(v => new XenRef<VM>(v.opaque_ref)).Concat(a.VMs).GroupBy(r => r.opaque_ref).Select(g => g.First()).ToList();
            return new AssignVMsToVMApplianceAction(a, refs, false);
        }

        /// <summary>The vApp's VMs as the Properties VMs page leaves them (members not in the list leave the vApp).</summary>
        public static AsyncAction SetMembersAction(VM_appliance a, IEnumerable<VM> vms) =>
            new AssignVMsToVMApplianceAction(a, vms.Select(v => new XenRef<VM>(v.opaque_ref)).ToList(), true);

        /// <summary>NewVMApplianceWizard.FinishWizard: create the vApp with its VMs.</summary>
        public static AsyncAction CreateAction(IXenConnection c, string name, string description, List<VM> vms) =>
            new CreateVMApplianceAction(new VM_appliance { name_label = name.Trim(), name_description = description ?? "", Connection = c }, vms);

        /// <summary>NewVMApplianceVMOrderAndDelaysPage.GetChangedSettings: only the VMs whose order or delay changed.</summary>
        public static Dictionary<VM, VMStartupOptions> ChangedStartup(IEnumerable<(VM Vm, long Order, long Delay)> rows) =>
            rows.Where(r => r.Vm.order != r.Order || r.Vm.start_delay != r.Delay)
                .ToDictionary(r => r.Vm, r => new VMStartupOptions(r.Order, r.Delay));

        /// <summary>SetVMStartupOptionsAction for the changed rows (null when nothing changed).</summary>
        public static AsyncAction StartupAction(IXenConnection c, Dictionary<VM, VMStartupOptions> changed) =>
            changed.Count == 0 ? null : new SetVMStartupOptionsAction(c, changed, false);

        /// <summary>NewVMApplianceWizard.GetSummary.</summary>
        public static string Summary(string name, IEnumerable<VM> vms) =>
            string.Format(Messages.VM_APPLIANCE_SUMMARY.Replace("\\n", "\n").Replace("\\r", "\r"), name, string.Join(", ", vms.Select(v => v.Name())));
    }
}
