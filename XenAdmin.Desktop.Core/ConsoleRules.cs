using System.Collections.Generic;
using System.Linq;
using XenAdmin.Actions;
using XenAdmin.Core;
using XenAPI;

namespace XenAdmin.Desktop
{
    /// <summary>The console tab's "Launch SSH Console" and Remote Desktop buttons (VNCTabView, XSVNCScreen.PollPort).</summary>
    public static class ConsoleRules
    {
        /// <summary>VNCTabView.IsSSHConsoleSupported: not for Windows; for dom0 only while the host is live.</summary>
        public static bool SshSupported(VM vm)
        {
            if (vm.IsWindows())
                return false;
            if (vm.IsControlDomainZero(out var host))
                return vm.Connection.Resolve(host.metrics)?.live == true;
            return true;
        }

        public static bool SshVisible(VM vm) => SshSupported(vm) && vm.power_state != vm_power_state.Halted;

        /// <summary>CanStartSSHConsole: running, with an address reported by the guest agent.</summary>
        public static bool SshEnabled(VM vm) => SshSupported(vm) && vm.power_state == vm_power_state.Running && !string.IsNullOrEmpty(vm.IPAddressForSSH());

        /// <summary>
        /// buttonSSH_Click with OpenSSH: the session's subject (else root) at the VM's address. The address comes from the
        /// guest, so it must be a literal IP and the user a plain name; null otherwise (nothing may reach a command line).
        /// </summary>
        public static string SshTarget(VM vm)
        {
            var subject = vm.Connection.Session?.SessionSubject is XenRef<Subject> r ? vm.Connection.Resolve(r) : null;
            return SshTarget(subject?.SubjectName ?? "root", vm.IPAddressForSSH());
        }

        public static string SshTarget(string user, string address) =>
            IsIpAddress(address) && IsPlainUser(user) ? $"{user}@{address}" : null;

        public static bool IsIpAddress(string address) =>
            !string.IsNullOrEmpty(address) && System.Net.IPAddress.TryParse(address.Trim('[', ']'), out _) && !address.Any(char.IsWhiteSpace);

        // AD names may be DOMAIN\user or user@domain; no leading '-', whitespace, quotes or shell characters
        private static bool IsPlainUser(string user) =>
            !string.IsNullOrEmpty(user) && user[0] != '-' && user.All(c => char.IsLetterOrDigit(c) || "._-\\@$".IndexOf(c) >= 0);

        /// <summary>The console switch is hidden for dom0 and for HVM guests without the RDP feature (Linux guests).</summary>
        public static bool RdpVisible(VM vm) => !vm.IsControlDomainZero(out _) && vm.IsHVM() && vm.HasRDP();

        /// <summary>CanEnableRDP: the guest can turn RDP on and it is off.</summary>
        public static bool CanTurnOnRdp(VM vm) => vm.RDPControlEnabled() && !vm.RDPEnabled();

        public static string RdpLabel(VM vm) => CanTurnOnRdp(vm) ? Messages.VNC_RDESKTOP_TURN_ON : Messages.VNC_RDESKTOP;

        public static bool RdpEnabled(VM vm) => vm.power_state == vm_power_state.Running && (vm.CanUseRDP() || CanTurnOnRdp(vm));

        /// <summary>XSVNCScreen.PollPort's order: IPv4 on connected networks, IPv6, then networks without a PIF.</summary>
        public static List<string> RdpAddresses(VM vm)
        {
            var metrics = vm.Connection.Resolve(vm.guest_metrics);
            if (metrics?.networks == null)
                return new List<string>();
            List<string> v4 = new(), v6 = new(), v4NoPif = new(), v6NoPif = new();
            var host = vm.Connection.Resolve(vm.resident_on);
            foreach (var vif in vm.Connection.ResolveAll(vm.VIFs))
            {
                var pif = Helpers.FindPIF(vm.Connection.Resolve(vif.network), host);
                foreach (var info in metrics.networks.Where(n => n.Key.StartsWith($"{vif.device}/ip")))
                {
                    var ipv4 = info.Key.EndsWith("ip") || info.Key.Contains("ipv4");
                    if (!ipv4 && !info.Key.Contains("ipv6"))
                        continue;
                    var address = ipv4 ? info.Value : $"[{info.Value}]";
                    if (pif == null)
                        (ipv4 ? v4NoPif : v6NoPif).Add(address);
                    else if (pif.LinkStatus() == PIF.LinkState.Connected)
                        (ipv4 ? v4 : v6).Add(address);
                }
            }
            return v4.Distinct().Concat(v6).Concat(v4NoPif).Concat(v6NoPif).Where(IsIpAddress).ToList();
        }

        /// <summary>"Turn on Remote Desktop" (FORCE_ENABLE_RDP): asks the guest agent to enable RDP.</summary>
        public static AsyncAction TurnOnRdp(VM vm) =>
            new DelegatedAsyncAction(vm.Connection, Messages.VNC_RDESKTOP_TURN_ON.Replace("&", ""), "", "",
                session => VM.call_plugin(session, vm.opaque_ref, "guest-agent-operation", "request-rdp-on", new Dictionary<string, string>()), true);
    }
}
