using System;
using System.Collections.Generic;
using System.Linq;
using XenAdmin.Actions;
using XenAdmin.Core;
using XenAdmin.Network;
using XenAdmin.XenSearch;
using XenAPI;

namespace XenAdmin.Desktop
{
    /// <summary>
    /// The Search tab's commands (SearchPage / SearchFor / QueryPanel): New and Edit Search, Save to a pool, Delete,
    /// Export and Import of .xensearch files. The editor covers the search scope, a full-text filter and one grouping.
    /// shortcut: no QueryElement tree of property conditions; an edited search keeps its existing filter and can add a
    /// text filter on top. Port QueryElement when users need field-by-field conditions.
    /// </summary>
    public static class SearchRules
    {
        /// <summary>The object types a search can return (SearchFor's type list), with their names.</summary>
        public static List<(ObjectTypes Type, string Name)> Types() =>
            PropertyAccessors.ObjectTypes_i18n.Select(kv => (kv.Value, kv.Key)).OrderBy(x => x.Key).ToList();

        /// <summary>The groupings offered (QueryPanel's Group by): none, or one property.</summary>
        public static List<(string Name, Func<Grouping> Make)> Groupings() => new List<(string, Func<Grouping>)>
        {
            (Messages.NONE, () => null),
            ("Type", () => new PropertyGrouping<ObjectTypes>(PropertyNames.type, null)),
            (Messages.POOL, () => new XenModelObjectPropertyGrouping<Pool>(PropertyNames.pool, null)),
            (Messages.SERVER, () => new XenModelObjectPropertyGrouping<Host>(PropertyNames.host, null)),
            ("Operating System", () => new PropertyGrouping<string>(PropertyNames.os_name, null)),
            ("Power State", () => new PropertyGrouping<vm_power_state>(PropertyNames.power_state, null)),
            (Messages.TAGS, () => new PropertyGrouping<string>(PropertyNames.tags, null)),
            (Messages.FOLDERS, () => new FolderGrouping((Grouping)null)),
            (Messages.VM_APPLIANCE, () => new XenModelObjectPropertyGrouping<VM_appliance>(PropertyNames.appliance, null)),
        };

        /// <summary>The grouping entry that matches a search's grouping (by type and property), else None.</summary>
        public static int GroupingIndex(Search search)
        {
            if (search?.Grouping == null)
                return 0;
            var list = Groupings();
            for (var i = 1; i < list.Count; i++)
                if (list[i].Make().Equals(search.Grouping))
                    return i;
            return 0;
        }

        /// <summary>The overview search (SearchPage's default): every type the Overview shows.</summary>
        public static Search Overview() => Search.SearchForAllTypes();

        public static string Validate(string name, ObjectTypes types) =>
            string.IsNullOrWhiteSpace(name) ? "Enter a name for the search."
            : types == ObjectTypes.None ? "Select at least one type of object to search for." : null;

        /// <summary>
        /// The search the editor describes: these types, the edited search's filter (if any) plus the text filter, the
        /// grouping; an edited search keeps its UUID, columns and sorting so saving replaces it.
        /// </summary>
        public static Search Build(string name, ObjectTypes types, string text, Grouping grouping, Search edited)
        {
            var filter = edited?.Query?.QueryFilter;
            if (!string.IsNullOrWhiteSpace(text))
            {
                var textFilter = Search.FullQueryFor(text.Trim());
                filter = filter == null ? textFilter : new GroupQuery(new[] { filter, textFilter }, GroupQuery.GroupQueryType.And);
            }
            // only a search saved in a pool keeps its identity (saving replaces it); others become new searches
            var uuid = IsSaved(edited) ? edited.UUID : null;
            var search = new Search(new Query(new QueryScope(types), filter), grouping, name.Trim(), uuid,
                edited?.Columns, edited?.Sorting ?? new Sort[] { });
            search.Connection = edited?.Connection;
            return search;
        }

        /// <summary>SearchPage.Save: the search goes into the pool's gui_config (where Saved Searches reads it).</summary>
        public static AsyncAction SaveAction(Search search, IXenConnection connection)
        {
            search.Connection = connection;
            return new DelegatedAsyncAction(connection, string.Format("Saving search '{0}'", search.Name), "Saving search", "Saved",
                _ => { if (!search.Save()) throw new Exception("The search could not be saved to the pool."); }, true);
        }

        /// <summary>A default search cannot be deleted; a saved one is removed from every pool.</summary>
        public static bool CanDelete(Search search) => IsSaved(search);

        /// <summary>A search saved in a connected pool (listed in Saved Searches and not a default one).</summary>
        public static bool IsSaved(Search search) =>
            search != null && !search.DefaultSearch && search.UUID != null && Search.Searches.Any(s => s.UUID == search.UUID && !s.DefaultSearch);

        public static AsyncAction DeleteAction(Search search, IXenConnection connection) =>
            new DelegatedAsyncAction(connection, string.Format("Deleting search '{0}'", search.Name), "Deleting search", "Deleted",
                _ => search.Delete(), true);

        /// <summary>SearchPage export: a .xensearch file with this search.</summary>
        public static void Export(Search search, string path) => search.Save(path);

        /// <summary>SearchPage import: the searches in a .xensearch file (null when it is not one).</summary>
        public static List<Search> Import(string path) =>
            Search.LoadFile(path)?.Select(s => new Search(s.Query, s.Grouping, s.Name, null, s.Columns, s.Sorting)).ToList(); // a new identity: never the original's
    }
}
