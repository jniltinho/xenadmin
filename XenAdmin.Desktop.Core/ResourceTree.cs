using System;
using System.Collections.Generic;
using System.Linq;
using XenAdmin.Core;
using XenAdmin.Network;
using XenAPI;

namespace XenAdmin.Desktop
{
    public class TreeNode
    {
        public string Name { get; set; }
        /// <summary>File name (no extension) of the 16px icon from XenAdmin/Images.</summary>
        public string Icon { get; set; }
        public IXenObject XenObject { get; set; }
        public IXenConnection Connection { get; set; }
        public bool IsRoot { get; set; }
        public List<TreeNode> Children { get; } = new List<TreeNode>();
        /// <summary>Expanded nodes, kept by key across tree rebuilds (only the root starts expanded, as in WinForms).</summary>
        public static readonly HashSet<string> ExpandedKeys = new HashSet<string> { "<root>" };
        public bool IsExpanded
        {
            get => ExpandedKeys.Contains(Key);
            set { if (value) ExpandedKeys.Add(Key); else ExpandedKeys.Remove(Key); }
        }
        public string Key => IsRoot ? "<root>" : XenObject?.opaque_ref ?? Connection?.HostnameWithPort ?? Name;
    }

    /// <summary>Infrastructure view, laid out like the WinForms navigation tree.</summary>
    public static class ResourceTree
    {
        public const string RootName = "XCP-ng Center";

        public static bool IsRealVm(VM vm) =>
            !vm.is_a_template && !vm.is_a_snapshot && !vm.is_control_domain && vm.Show(false);

        public static bool IsUserTemplate(VM vm) =>
            vm.is_a_template && !vm.is_a_snapshot && !vm.DefaultTemplate() && vm.Show(false);

        /// <summary>Pool name remembered from the last connection, else the address (XenConnection.Name).</summary>
        private static string Name(IXenConnection c) => c is XenConnection xc ? xc.Name : c.Hostname;

        public static bool IsShownSr(SR sr) => sr.Show(false) && !sr.IsToolsSR();

        private static TreeNode Node(IXenObject o, string name = null) =>
            new TreeNode { XenObject = o, Connection = o.Connection, Name = name ?? o.Name(), Icon = IconNames.For(o) };

        public static TreeNode Root(IEnumerable<IXenConnection> connections)
        {
            var root = new TreeNode { Name = RootName, Icon = IconNames.Home, IsRoot = true };
            root.Children.AddRange(connections.Select(Build).OrderBy(n => n.Name, StringComparer.CurrentCultureIgnoreCase));
            return root;
        }

        public static TreeNode Build(IXenConnection connection)
        {
            if (!connection.IsConnected || !connection.CacheIsPopulated)
                return new TreeNode { Connection = connection, Name = Name(connection), Icon = IconNames.For(connection) };

            var pool = Helpers.GetPool(connection);
            var hosts = connection.Cache.Hosts.OrderBy(h => h.Name(), StringComparer.CurrentCultureIgnoreCase).ToList();
            TreeNode top;
            var hostNodes = new Dictionary<string, TreeNode>();
            if (pool != null)
            {
                top = Node(pool);
                foreach (var host in hosts)
                    top.Children.Add(hostNodes[host.opaque_ref] = Node(host));
            }
            else
            {
                var host = hosts.FirstOrDefault();
                if (host == null)
                    return new TreeNode { Connection = connection, Name = Name(connection), Icon = IconNames.For(connection) };
                top = hostNodes[host.opaque_ref] = Node(host);
            }

            TreeNode ParentFor(Host home) =>
                home != null && hostNodes.TryGetValue(home.opaque_ref, out var n) ? n : top;

            // under each host: its VMs first, then its storage (as in the WinForms tree)
            foreach (var vm in connection.Cache.VMs.Where(IsRealVm).OrderBy(v => v.Name(), StringComparer.CurrentCultureIgnoreCase))
                ParentFor(vm.Home()).Children.Add(Node(vm));
            // user templates sit with the pool (default templates stay hidden, as with DefaultTemplatesVisible=false)
            foreach (var t in connection.Cache.VMs.Where(IsUserTemplate).OrderBy(v => v.Name(), StringComparer.CurrentCultureIgnoreCase))
                top.Children.Add(Node(t));
            foreach (var sr in connection.Cache.SRs.Where(IsShownSr).OrderBy(s => s.NameWithoutHost(), StringComparer.CurrentCultureIgnoreCase))
                ParentFor(sr.Home()).Children.Add(Node(sr, sr.NameWithoutHost()));
            return top;
        }

        /// <summary>Keeps nodes whose name contains the text (case-insensitive) and their ancestors.</summary>
        public static TreeNode Filter(TreeNode node, string text)
        {
            if (string.IsNullOrWhiteSpace(text))
                return node;
            var kids = node.Children.Select(c => Filter(c, text)).Where(c => c != null).ToList();
            if (kids.Count == 0 && !node.IsRoot && node.Name.IndexOf(text.Trim(), StringComparison.CurrentCultureIgnoreCase) < 0)
                return null;
            var copy = new TreeNode { Name = node.Name, Icon = node.Icon, XenObject = node.XenObject, Connection = node.Connection, IsRoot = node.IsRoot };
            copy.Children.AddRange(kids);
            return copy;
        }

        public static IEnumerable<TreeNode> Flatten(TreeNode n) => new[] { n }.Concat(n.Children.SelectMany(Flatten));
    }
}
