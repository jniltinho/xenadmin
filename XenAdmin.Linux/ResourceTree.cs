using System;
using System.Collections.Generic;
using System.Linq;
using XenAdmin.Core;
using XenAdmin.Network;
using XenAPI;

namespace XenAdmin.Linux
{
    public class TreeNode
    {
        public string Name { get; set; }
        /// <summary>Status colour of the bullet drawn before the name.</summary>
        public string Color { get; set; } = "#808080";
        /// <summary>Kind and state, for the dimmed suffix and screen readers.</summary>
        public string Kind { get; set; }
        public IXenObject XenObject { get; set; }
        public IXenConnection Connection { get; set; }
        public List<TreeNode> Children { get; } = new List<TreeNode>();
        public string Label => $"{Name} ({Kind})";
        public string Key => XenObject?.opaque_ref ?? Connection?.HostnameWithPort ?? Name;
    }

    /// <summary>Builds the server → pool → host → VM/SR hierarchy from a connection cache.</summary>
    public static class ResourceTree
    {
        public static bool IsRealVm(VM vm) =>
            !vm.is_a_template && !vm.is_a_snapshot && !vm.is_control_domain && vm.Show(false);

        public const string Green = "#2E9E44", Gray = "#8A8A8A", Orange = "#E08A00", Red = "#C62828", Blue = "#1E6FD9";

        public static string PowerColor(vm_power_state state) => state switch
        {
            vm_power_state.Running => Green,
            vm_power_state.Suspended => Orange,
            vm_power_state.Paused => Orange,
            _ => Gray
        };

        public static string PowerText(vm_power_state state) => state switch
        {
            vm_power_state.Running => "em execução",
            vm_power_state.Suspended => "suspensa",
            vm_power_state.Paused => "pausada",
            _ => "desligada"
        };

        public static TreeNode Build(IXenConnection connection)
        {
            var root = new TreeNode { Connection = connection, Name = connection.Hostname, Color = connection.IsConnected ? Blue : Gray, Kind = "servidor" };
            if (!connection.IsConnected || !connection.CacheIsPopulated)
            {
                root.Kind = connection.IsConnected || connection.InProgress ? "conectando…" : "desconectado";
                return root;
            }

            var pool = Helpers.GetPool(connection);
            var top = root;
            if (pool != null)
            {
                top = new TreeNode { Connection = connection, XenObject = pool, Name = pool.Name(), Color = Blue, Kind = "pool" };
                root.Children.Add(top);
            }
            root.Name = pool?.Name() ?? connection.Cache.Hosts.FirstOrDefault()?.Name() ?? connection.Hostname;

            var hostNodes = new Dictionary<string, TreeNode>();
            foreach (var host in connection.Cache.Hosts.OrderBy(h => h.Name(), StringComparer.CurrentCultureIgnoreCase))
            {
                var node = new TreeNode { Connection = connection, XenObject = host, Name = host.Name(), Color = host.IsLive() ? Green : Red, Kind = host.IsLive() ? "host" : "host inativo" };
                hostNodes[host.opaque_ref] = node;
                top.Children.Add(node);
            }
            // a lone host without pool: VMs and SRs go straight under it
            var defaultParent = pool == null && hostNodes.Count == 1 ? hostNodes.Values.First() : top;

            foreach (var vm in connection.Cache.VMs.Where(IsRealVm).OrderBy(v => v.Name(), StringComparer.CurrentCultureIgnoreCase))
            {
                var parent = vm.resident_on != null && hostNodes.TryGetValue(vm.resident_on.opaque_ref, out var h) ? h : defaultParent;
                parent.Children.Add(new TreeNode { Connection = connection, XenObject = vm, Name = vm.Name(), Color = PowerColor(vm.power_state), Kind = "VM " + PowerText(vm.power_state) });
            }

            foreach (var sr in connection.Cache.SRs.Where(s => s.Show(false)).OrderBy(s => s.Name(), StringComparer.CurrentCultureIgnoreCase))
            {
                var home = sr.Home();
                var parent = home != null && hostNodes.TryGetValue(home.opaque_ref, out var h) ? h : defaultParent;
                parent.Children.Add(new TreeNode { Connection = connection, XenObject = sr, Name = sr.Name(), Color = "#6D4C9F", Kind = "storage" });
            }
            return root;
        }

        /// <summary>Keeps nodes whose name contains the text (case-insensitive) and their ancestors.</summary>
        public static TreeNode Filter(TreeNode node, string text)
        {
            if (string.IsNullOrWhiteSpace(text))
                return node;
            var kids = node.Children.Select(c => Filter(c, text)).Where(c => c != null).ToList();
            if (kids.Count == 0 && node.Name.IndexOf(text.Trim(), StringComparison.CurrentCultureIgnoreCase) < 0)
                return null;
            var copy = new TreeNode { Name = node.Name, Color = node.Color, Kind = node.Kind, XenObject = node.XenObject, Connection = node.Connection };
            copy.Children.AddRange(kids);
            return copy;
        }

        public static IEnumerable<(string, string)> Properties(IXenObject o)
        {
            switch (o)
            {
                case VM vm:
                    yield return ("Nome", vm.Name());
                    yield return ("Descrição", vm.Description());
                    yield return ("UUID", vm.uuid);
                    yield return ("Estado", vm.power_state.ToString());
                    yield return ("Memória", Util.MemorySizeStringSuitableUnits(vm.memory_dynamic_max, true));
                    yield return ("vCPUs", vm.VCPUs_at_startup.ToString());
                    yield return ("Host", vm.Connection?.Resolve(vm.resident_on)?.Name() ?? "—");
                    break;
                case Host host:
                    yield return ("Nome", host.Name());
                    yield return ("Descrição", host.Description());
                    yield return ("UUID", host.uuid);
                    yield return ("Endereço", host.address);
                    yield return ("Versão", host.ProductVersionText());
                    yield return ("Ativo", host.IsLive() ? "sim" : "não");
                    var metrics = host.Connection?.Resolve(host.metrics);
                    if (metrics != null)
                        yield return ("Memória livre/total", $"{Util.MemorySizeStringSuitableUnits(metrics.memory_free, true)} / {Util.MemorySizeStringSuitableUnits(metrics.memory_total, true)}");
                    break;
                case SR sr:
                    yield return ("Nome", sr.Name());
                    yield return ("Descrição", sr.Description());
                    yield return ("UUID", sr.uuid);
                    yield return ("Tipo", sr.type);
                    yield return ("Compartilhado", sr.shared ? "sim" : "não");
                    yield return ("Uso", $"{Util.DiskSizeString(sr.physical_utilisation)} / {Util.DiskSizeString(sr.physical_size)}");
                    break;
                case Pool pool:
                    yield return ("Nome", pool.Name());
                    yield return ("Descrição", pool.Description());
                    yield return ("UUID", pool.uuid);
                    yield return ("Coordenador", pool.Connection?.Resolve(pool.master)?.Name() ?? "—");
                    yield return ("HA", pool.ha_enabled ? "ativo" : "inativo");
                    break;
            }
        }
    }
}
