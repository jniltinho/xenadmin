using System;
using System.Collections.Generic;
using System.IO;
using System.Linq;
using XenAdmin.Actions;
using XenAdmin.Core;
using XenAdmin.Network;
using XenAPI;

namespace XenAdmin.Desktop
{
    /// <summary>
    /// VM export and import as XVA: ExportVMCommand with the export wizard's XVA path (ExportVmAction) and the import
    /// wizard's XVA path (ImportVmAction on a server and SR). OVF/OVA appliances are not ported yet.
    /// </summary>
    public static class XvaRules
    {
        /// <summary>ExportVMCommand.CanRun: a VM (not a template) whose server allows the export.</summary>
        public static bool CanExport(VM vm) =>
            vm != null && !vm.is_a_template && !vm.is_a_snapshot && !vm.Locked && vm.allowed_operations != null
            && vm.allowed_operations.Contains(vm_operations.export);

        /// <summary>ExportVMCommand.GetCantRunReasonCore: a running VM must be shut down or suspended first.</summary>
        public static string ExportBlocker(VM vm) =>
            vm != null && vm.power_state == vm_power_state.Running ? Messages.MAINWINDOW_EXPORT_VM_AS_BACKUP_TOOLTIP : null;

        /// <summary>ExportVmAction (HTTP /export to the file; verify reads it back afterwards).</summary>
        public static AsyncAction ExportAction(VM vm, string filename, bool verify) =>
            new ExportVmAction(vm.Connection, vm.Home(), vm, filename, verify);

        /// <summary>The import wizard's file check: an XVA (optionally gzipped); OVF/OVA and disk images are not ported.</summary>
        public static string FileError(string path)
        {
            if (string.IsNullOrWhiteSpace(path))
                return "Select the file to import.";
            var name = path.ToLowerInvariant();
            if (!(name.EndsWith(".xva") || name.EndsWith(".xva.gz")))
                return "Only XVA files (.xva, .xva.gz) can be imported in this version.";
            return File.Exists(path) ? null : "The file does not exist.";
        }

        /// <summary>StoragePickerPage (SrPickerItem for an import): the SRs that can hold new disks on that server and are not full.</summary>
        public static List<SR> Targets(IXenConnection c, Host affinity) =>
            StorageRules.NewDiskTargets(c, affinity).Where(sr => sr.CanFitDisks(out _, Array.Empty<VDI>())).ToList();

        /// <summary>The chosen server and SR are still there (the connection may have dropped while the dialog was open).</summary>
        public static string TargetError(IXenConnection c, Host affinity, SR sr)
        {
            if (!c.IsConnected || Helpers.GetCoordinator(c) == null)
                return string.Format("The connection to '{0}' was lost.", Helpers.GetName(c));
            if (affinity != null && c.Resolve(new XenRef<Host>(affinity.opaque_ref)) == null)
                return string.Format("The server '{0}' is no longer in the pool.", affinity.Name());
            var current = c.Resolve(new XenRef<SR>(sr.opaque_ref));
            if (current == null || !Targets(c, affinity).Any(t => t.opaque_ref == current.opaque_ref))
                return "The storage repository is no longer available or is full.";
            return null;
        }

        /// <summary>
        /// StoragePickerPage.PageLeaveCore: ImportVmAction, told at once that the wizard is done (keep the VM's own
        /// networks) since this app has no network mapping page for XVA. The VM is started afterwards with
        /// StartAfterImport, as a tracked action, rather than by ImportVmAction's own untracked start.
        /// </summary>
        public static ImportVmAction ImportAction(IXenConnection c, Host affinity, string filename, SR sr)
        {
            var action = new ImportVmAction(c, affinity, filename, sr, (_, _) => throw new CancelledException(), null);
            action.EndWizard(false, null);
            return action;
        }

        /// <summary>"Start VM after import": resume a suspended import, start a halted one (null for a template or a failed import).</summary>
        public static AsyncAction StartAfterImport(ImportVmAction import)
        {
            var vm = import.Succeeded ? import.Connection.Resolve(new XenRef<VM>(import.VM?.opaque_ref)) : null;
            if (vm == null || vm.is_a_template)
                return null;
            return VmOperations.Create(vm, vm.power_state == vm_power_state.Suspended ? VmOp.Resume : VmOp.Start);
        }
    }
}
