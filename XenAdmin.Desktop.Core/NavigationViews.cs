using System;
using System.Collections.Generic;
using System.Linq;
using XenAdmin.Core;
using XenAdmin.Model;
using XenAdmin.XenSearch;
using XenAPI;

namespace XenAdmin.Desktop
{
    /// <summary>The navigation pane's views besides Infrastructure (NavigationPane.NavigationMode).</summary>
    public enum NavView { Infrastructure, Objects, Tags, Folders, CustomFields, Vapps, SavedSearch }

    /// <summary>
    /// The Objects, Organization Views (Tags, Folders, Custom Fields, vApps) and Saved Searches trees of
    /// MainWindowTreeBuilder.CreateNewRootNode, built with the XenModel searches and groupings of OrganizationalView.
    /// </summary>
    public static class NavigationViews
    {
        /// <summary>OrganizationalView.OrgViewSearch: the search and grouping of each organization view.</summary>
        private static Search ViewSearch(NavView view)
        {
            QueryFilter filter;
            Grouping grouping;
            var types = ObjectTypes.AllExcFolders;
            switch (view)
            {
                case NavView.Objects:
                    filter = null;
                    grouping = new PropertyGrouping<ObjectTypes>(PropertyNames.type, null);
                    break;
                case NavView.Tags:
                    filter = new ListEmptyQuery<string>(PropertyNames.tags, false);
                    grouping = new PropertyGrouping<string>(PropertyNames.tags, null);
                    break;
                case NavView.Folders:
                    filter = new NullQuery<Folder>(PropertyNames.folder, false);
                    grouping = new FolderGrouping((Grouping)null);
                    types = ObjectTypes.AllIncFolders;
                    break;
                case NavView.CustomFields:
                    filter = new BooleanQuery(PropertyNames.has_custom_fields, true);
                    grouping = new AllCustomFieldsGrouping((Grouping)null);
                    break;
                case NavView.Vapps:
                    filter = new BooleanQuery(PropertyNames.in_any_appliance, true);
                    grouping = new XenModelObjectPropertyGrouping<VM_appliance>(PropertyNames.appliance, null);
                    break;
                default:
                    throw new ArgumentOutOfRangeException(nameof(view));
            }
            return new Search(new Query(new QueryScope(types), filter), grouping, "", null, null, new Sort[] { });
        }

        /// <summary>OrganizationalView.GroupingName / GetGroupIcon: the root of each view.</summary>
        public static (string Name, string Icon) RootOf(NavView view) => view switch
        {
            NavView.Objects => (Messages.VIEW_OBJECTS, "objects_16"),
            NavView.Tags => (Messages.TAGS, "000_Tag_h32bit_16"),
            NavView.Folders => (Messages.FOLDERS, "000_Folder_open_h32bit_16"),
            NavView.CustomFields => (Messages.CUSTOM_FIELDS, "000_Fields_h32bit_16"),
            NavView.Vapps => (Messages.VM_APPLIANCES, "000_VirtualAppliance_h32bit_16"),
            _ => (ResourceTree.RootName, "Logo")
        };

        /// <summary>The tree of an organization view, or of a saved search (MainWindowTreeBuilder.CreateNewRootNode).</summary>
        public static TreeNode Build(NavView view, Search saved = null)
        {
            TreeNode root;
            Search search;
            if (view == NavView.SavedSearch)
            {
                search = saved ?? throw new ArgumentNullException(nameof(saved));
                root = new TreeNode { Name = search.Name, Icon = search.DefaultSearch ? "000_defaultSpyglass_h32bit_16" : "000_Search_h32bit_16", IsRoot = true };
            }
            else
            {
                search = ViewSearch(view);
                var (name, icon) = RootOf(view);
                root = new TreeNode { Name = name, Icon = icon, IsRoot = true };
            }
            search.PopulateAdapters(new Acceptor(root, view + ":" + (saved?.UUID ?? "") + ":"));
            return root;
        }

        /// <summary>The searches of Saved Searches: the default ones and those saved in the pools, by name.</summary>
        public static List<Search> SavedSearches() =>
            Search.Searches.OrderBy(s => !s.DefaultSearch).ThenBy(s => s.Name, StringComparer.CurrentCultureIgnoreCase).ToList();

        /// <summary>The objects under a group node, with their depth (the Search tab of a group).</summary>
        public static IEnumerable<(IXenObject Object, int Depth)> Objects(TreeNode node, int depth = 0)
        {
            foreach (var child in node.Children)
            {
                if (child.XenObject != null && !(child.XenObject is Folder))
                {
                    yield return (child.XenObject, depth);
                    foreach (var x in Objects(child, depth + 1)) yield return x;
                }
                else
                    foreach (var x in Objects(child, depth)) yield return x;
            }
        }

        /// <summary>MainWindowTreeNodeGroupAcceptor: XenObjects become object nodes, other groups named group nodes.</summary>
        private sealed class Acceptor : IAcceptGroups
        {
            private readonly TreeNode _parent;
            private readonly string _path;

            public Acceptor(TreeNode parent, string path)
            {
                _parent = parent;
                _path = path;
            }

            public IAcceptGroups Add(Grouping grouping, object group, int indent)
            {
                if (group == null)
                    return null;
                TreeNode node;
                if (group is IXenObject o)
                {
                    // MainWindowTreeNodeGroupAcceptor.AddVMNode: vTPM templates are hidden where vTPMs are restricted
                    if (o is VM t && t.is_a_template && t.Connection.Cache.Hosts.Any(Host.RestrictVtpm)
                        && t.platform.TryGetValue("vtpm", out var v) && v.ToLower() == "true")
                        return null;
                    node = new TreeNode { Name = Name(o), Icon = IconNames.For(o), XenObject = o, Connection = o.Connection, GroupKey = _path + Segment(o.opaque_ref) };
                }
                else
                {
                    var name = group is DateTime d ? d.ToLocalTime().ToString("g") : grouping.GetGroupName(group);
                    // the key holds the grouping and the raw value, each segment escaped, so different groups never share one
                    var raw = grouping.GetType().Name + "=" + System.Convert.ToString(group, System.Globalization.CultureInfo.InvariantCulture);
                    node = new TreeNode { Name = name, Icon = GroupIcon(grouping, group), GroupKey = _path + Segment(raw) };
                }
                _parent.Children.Add(node);
                return new Acceptor(node, node.GroupKey + "/");
            }

            private static string Segment(string s) => Uri.EscapeDataString(s ?? "");

            public void FinishedInThisGroup(bool defaultExpand)
            {
            }

            private static string Name(IXenObject o) => o switch
            {
                VM vm when vm.IsHidden() => string.Format(Messages.X_HIDDEN, vm.Name()),
                SR sr => sr.IsHidden() ? string.Format(Messages.X_HIDDEN, sr.NameWithoutHost()) : sr.NameWithoutHost(),
                XenAPI.Network n when n.IsMember() => string.Format(Messages.NIC_BONDED_MEMBER, n.Name()),
                XenAPI.Network n when n.IsHidden() => string.Format(Messages.X_HIDDEN, n.Name()),
                VDI vdi => string.IsNullOrEmpty(vdi.Name()) ? Messages.NO_NAME : vdi.Name(),
                _ => Helpers.GetName(o).Ellipsise(1000)
            };

            /// <summary>Grouping.GetGroupIcon for the groups of these views.</summary>
            private static string GroupIcon(Grouping grouping, object group) => group switch
            {
                ObjectTypes t => t switch
                {
                    ObjectTypes.Pool => "000_PoolConnected_h32bit_16",
                    ObjectTypes.Server or ObjectTypes.DisconnectedServer => "000_TreeConnected_h32bit_16",
                    ObjectTypes.VM => "000_StartVM_h32bit_16",
                    ObjectTypes.Snapshot => "000_VMSnapShotDiskOnly_h32bit_16",
                    ObjectTypes.UserTemplate => "000_UserTemplate_h32bit_16",
                    ObjectTypes.DefaultTemplate => "000_VMTemplate_h32bit_16",
                    ObjectTypes.LocalSR or ObjectTypes.RemoteSR => "000_Storage_h32bit_16",
                    ObjectTypes.Network => "000_Network_h32bit_16",
                    ObjectTypes.VDI => "000_VirtualStorage_h32bit_16",
                    ObjectTypes.Appliance => "000_VirtualAppliance_h32bit_16",
                    _ => "000_Folder_open_h32bit_16"
                },
                string _ when grouping is PropertyGrouping<string> => "000_Tag_h32bit_16",
                _ when grouping is AllCustomFieldsGrouping || grouping.GetType().Name.Contains("CustomField") => "000_Fields_h32bit_16",
                _ => "000_Folder_open_h32bit_16"
            };
        }
    }
}
