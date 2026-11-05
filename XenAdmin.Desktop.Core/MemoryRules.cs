using System;
using System.Collections.Generic;
using System.Globalization;
using System.Linq;
using XenAdmin.Actions;
using XenAdmin.Core;
using XenAPI;

namespace XenAdmin.Desktop
{
    /// <summary>
    /// VM memory: the Memory tab's Edit (VMMemoryControlsNoEdit), and the Memory Settings dialog (BallooningDialog with
    /// VMMemoryControlsBasic / VMMemoryControlsAdvanced; limits and confirmations of VMMemoryControlsEdit).
    /// </summary>
    public static class MemoryRules
    {
        /// <summary>VMMemoryControlsNoEdit: Edit for a halted VM, or a running one whose virtualization state is known.</summary>
        public static bool CanEdit(VM vm) =>
            vm != null && !vm.Locked && vm.Connection.IsConnected
            && (vm.power_state == vm_power_state.Halted
                || vm.power_state == vm_power_state.Running && !vm.GetVirtualizationStatus(out _).HasFlag(VM.VirtualizationStatus.Unknown));

        /// <summary>BallooningDialog: dynamic memory (minimum, maximum, static maximum) when licensed and supported.</summary>
        public static bool Advanced(VM vm) => !Helpers.FeatureForbidden(vm.Connection, Host.RestrictDMC) && vm.UsesBallooning();

        /// <summary>VMMemoryControlsBasic: "Automatically allocate memory within this range" needs ballooning and the license.</summary>
        public static bool DynamicAllowed(VM vm) => vm.SupportsBallooning() && !Helpers.FeatureForbidden(vm.Connection, Host.RestrictDMC);

        /// <summary>VMMemoryControlsBasic's "A memory range cannot be set" text for one VM (null when a range is allowed).</summary>
        public static string DynamicUnavailable(VM vm)
        {
            if (DynamicAllowed(vm))
                return null;
            if (Helpers.FeatureForbidden(vm.Connection, Host.RestrictDMC))
                return Messages.DMC_UNAVAILABLE_LICENSE_RESTRICTION;
            if (vm.is_a_template)
                return Messages.DMC_UNAVAILABLE_TEMPLATE;
            return vm.GetVirtualizationStatus(out _).HasFlag(VM.VirtualizationStatus.IoDriversInstalled)
                ? Messages.DMC_UNAVAILABLE_NOTSUPPORTED
                : vm.HasNewVirtualizationStates() ? Messages.DMC_UNAVAILABLE_NO_IO_NO_MGMNT : string.Format(Messages.DMC_UNAVAILABLE_NOTOOLS, BrandManager.VmTools);
        }

        /// <summary>BallooningDialog's rubric.</summary>
        public static string Rubric(VM vm) => Advanced(vm)
            ? vm.is_a_template ? Messages.BALLOONING_RUBRIC_ADVANCED_TEMPLATE : Messages.BALLOONING_RUBRIC_ADVANCED
            : vm.is_a_template ? Messages.BALLOONING_RUBRIC_TEMPLATE : Messages.BALLOONING_RUBRIC;

        /// <summary>MemorySpinner.Initialize: MB up to 1 GB of static maximum, GB above.</summary>
        public static bool InMegabytes(VM vm) => vm.memory_static_max <= Util.BINARY_GIGA;

        /// <summary>VMMemoryControlsEdit.GetMemoryRatio: the smallest dynamic minimum, as a fraction of the static maximum.</summary>
        public static double Ratio(VM vm)
        {
            var config = Helpers.GetOtherConfig(Helpers.GetPoolOfOne(vm.Connection));
            var key = vm.IsHVM() ? "memory-ratio-hvm" : "memory-ratio-pv";
            if (config != null && config.TryGetValue(key, out var v)
                && double.TryParse(v, NumberStyles.Float, CultureInfo.InvariantCulture, out var f) && f >= 0 && f <= 1)
                return f;
            return 0.25;
        }

        /// <summary>VMMemoryControlsEdit.CalcMaxDynMin: what the VM's host can still give a running VM (-1: no limit).</summary>
        public static double MaxDynMin(VM vm)
        {
            if (vm.power_state != vm_power_state.Running && vm.power_state != vm_power_state.Paused)
                return -1;
            var host = vm.Connection.Resolve(vm.resident_on);
            var metrics = host == null ? null : host.Connection.Resolve(host.metrics);
            if (metrics == null)
                return -1;
            double used = host.memory_overhead;
            foreach (var other in host.Connection.ResolveAll(host.resident_VMs))
            {
                used += other.memory_overhead;
                if (other.is_control_domain)
                    used += other.Connection.Resolve(other.metrics)?.memory_actual ?? 0;
                else if (other.opaque_ref != vm.opaque_ref)
                    used += other.memory_dynamic_min;
            }
            return Math.Max(0, metrics.memory_total - used);
        }

        /// <summary>VMMemoryControlsEdit.MemorySpinnerMax: the license limit, or the current static maximum if larger.</summary>
        public static long MemoryMax(VM vm) => Math.Max(vm.memory_static_max, vm.MaxMemAllowed());

        /// <summary>
        /// The ranges of the spinners (VMMemoryControlsAdvanced.SetSpinnerRanges): null when the values are valid, else
        /// the reason. Basic mode has one value: dynamic minimum = dynamic maximum = static maximum.
        /// </summary>
        public static string Validate(VM vm, long dynamicMin, long dynamicMax, long staticMax)
        {
            var maxDynMin = MaxDynMin(vm);
            var ratio = Ratio(vm);
            if (staticMax < Util.BINARY_MEGA || staticMax > MemoryMax(vm))
                return string.Format("The memory must be between 1 MB and {0}.", Util.MemorySizeStringSuitableUnits(MemoryMax(vm), true));
            if (maxDynMin >= 0 && staticMax * ratio > maxDynMin && staticMax > vm.memory_static_max)
                return "The server does not have enough free memory for this maximum.";
            if (dynamicMax > staticMax || dynamicMin > dynamicMax)
                return "The dynamic minimum must not be above the dynamic maximum, nor the dynamic maximum above the static maximum.";
            // VMMemoryControlsAdvanced.Spinners_ValueChanged: a new static maximum forces the minimum up to its ratio;
            // with the static maximum unchanged, an existing lower minimum stays allowed (DynMinSpinnerMin)
            var lowest = Math.Max(vm.memory_static_min, (long)(staticMax * ratio));
            if (dynamicMin < lowest && (staticMax != vm.memory_static_max || dynamicMin < vm.memory_dynamic_min))
                return string.Format("The dynamic minimum must be at least {0}.", Util.MemorySizeStringSuitableUnits(lowest, true));
            if (maxDynMin >= 0 && dynamicMin > maxDynMin && dynamicMin > vm.memory_dynamic_min)
                return "The server does not have enough free memory for this dynamic minimum.";
            return null;
        }

        /// <summary>
        /// VMMemoryControlsEdit.ChangeMemorySettings: an unchanged static maximum (to the MB) keeps the server's exact
        /// values, so rounding never resizes the VM.
        /// </summary>
        public static (long DynamicMin, long DynamicMax, long StaticMax) Settings(VM vm, long dynamicMin, long dynamicMax, long staticMax)
        {
            if (vm.memory_static_max / Util.BINARY_MEGA == staticMax / Util.BINARY_MEGA)
            {
                if (dynamicMin == staticMax) dynamicMin = vm.memory_static_max;
                if (dynamicMax == staticMax) dynamicMax = vm.memory_static_max;
                staticMax = vm.memory_static_max;
            }
            return (dynamicMin, dynamicMax, staticMax);
        }

        /// <summary>The warning before a running VM restarts to take a new static maximum (null when none is needed).</summary>
        public static string Confirmation(VM vm, long staticMax)
        {
            if (vm.power_state == vm_power_state.Halted || vm.memory_static_max == staticMax)
                return null;
            if (Advanced(vm))
                return Messages.CONFIRM_CHANGE_STATIC_MAX_SINGULAR;
            return vm.SupportsBallooning() && !Helpers.FeatureForbidden(vm, Host.RestrictDMC)
                ? Messages.CONFIRM_CHANGE_MEMORY_MAX_SINGULAR
                : Messages.CONFIRM_CHANGE_MEMORY_SINGULAR;
        }

        /// <summary>
        /// ChangeMemorySettingsAction (it shuts the VM down and starts it again when the static maximum changes). An
        /// HA-protected VM that is not agile is never unprotected silently: the restart stops instead (the WinForms
        /// dialog asks first; this app has no such dialog for actions yet, as VmOperations' AllowHaUnprotect = false).
        /// </summary>
        public static AsyncAction ChangeAction(VM vm, long dynamicMin, long dynamicMax, long staticMax) =>
            new ChangeMemorySettingsAction(vm, string.Format(Messages.ACTION_CHANGE_MEMORY_SETTINGS, vm.Name()),
                vm.memory_static_min, dynamicMin, dynamicMax, staticMax, (_, _) => throw new CancelledException(), null, false);

        /// <summary>The Memory tab's bar for a VM: the used memory and the dynamic range, as fractions of the static maximum.</summary>
        public static (double Used, double DynamicMin, double DynamicMax) VmBar(VM vm)
        {
            double max = Math.Max(1, vm.memory_static_max);
            var metrics = vm.Connection.Resolve(vm.metrics);
            var used = vm.power_state is vm_power_state.Running or vm_power_state.Paused && metrics != null ? metrics.memory_actual : 0;
            return (Math.Min(1, used / max), vm.memory_dynamic_min / max, vm.memory_dynamic_max / max);
        }

        /// <summary>HostShinyBar: the host's memory split into Xen and the control domain, each resident VM, and free memory.</summary>
        public static List<(string Name, long Bytes)> HostSegments(Host host)
        {
            var result = new List<(string, long)>();
            var metrics = host.Connection.Resolve(host.metrics);
            if (metrics == null)
                return result;
            long xen = host.memory_overhead;
            foreach (var vm in host.Connection.ResolveAll(host.resident_VMs).OrderBy(v => !v.is_control_domain).ThenBy(v => v.Name()))
            {
                var actual = vm.Connection.Resolve(vm.metrics)?.memory_actual ?? 0;
                if (vm.is_control_domain)
                    xen += actual + vm.memory_overhead;
                else
                    result.Add((vm.Name(), actual + vm.memory_overhead));
            }
            result.Insert(0, (BrandManager.ProductBrand + " (Xen and control domain)", xen));
            result.Add(("Free", Math.Max(0, metrics.memory_free)));
            return result;
        }
    }
}
