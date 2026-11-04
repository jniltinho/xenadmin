using System.Collections.Generic;
using System.Linq;
using System.Text.RegularExpressions;
using XenAdmin.Actions;
using XenAdmin.Core;
using XenAdmin.Network;
using XenAPI;

namespace XenAdmin.Desktop
{
    /// <summary>New Storage Repository wizard (NewSRWizard with the VHDoNFS and CifsFrontend pages) for shared virtual disk storage.</summary>
    public static class SrRules
    {
        public enum Kind { Nfs, Smb }

        public static SR.SRTypes TypeOf(Kind kind) => kind == Kind.Nfs ? SR.SRTypes.nfs : SR.SRTypes.smb;
        public static string TypeName(Kind kind) => kind == Kind.Nfs ? Messages.NEWSR_VHDONFS_TYPE_NAME : Messages.NEWSR_CIFS_TYPE_NAME;
        public static string Blurb(Kind kind) => kind == Kind.Nfs ? Messages.NEWSR_VHDONFS_BLURB : Messages.NEWSR_CIFS_BLURB;

        /// <summary>SrWizardType.ResetSrName: "NFS virtual disk storage", "(1)" suffix when taken.</summary>
        public static string DefaultName(Kind kind, IXenConnection c)
        {
            var name = kind == Kind.Nfs ? Messages.SRWIZARD_NFS_STORAGE : Messages.SRWIZARD_CIFS_STORAGE;
            var taken = c.Cache.SRs.Select(sr => sr.Name()).ToList();
            if (!taken.Contains(name)) return name;
            for (var i = 1; ; i++)
            {
                var n = string.Format(Messages.NEWVM_DEFAULTNAME, name, i);
                if (!taken.Contains(n)) return n;
            }
        }

        /// <summary>SrWizardHelpers.ValidateNfsSharename / ValidateCifsSharename: server:/path, or \\server\share.</summary>
        public static bool ValidLocation(Kind kind, string path) =>
            path != null && (kind == Kind.Nfs
                ? Regex.IsMatch(path.Trim(), @"^[^/:*?<>|\\]+:/[^:*?<>|\\]+$")
                : Regex.IsMatch(path, @"^\\(\\[^\/:*?<>|\\]+\\{0,1})+$"));

        /// <summary>The pages' DeviceConfig. nfsVersion is "" (3), "4" or "4.1"; credentials only when given (SMB).</summary>
        public static Dictionary<string, string> DeviceConfig(Kind kind, string path, string options, string nfsVersion, string user, string password, bool forProbe = false)
        {
            var parts = (kind == Kind.Nfs ? path.Trim() : path).Split(':');
            var d = new Dictionary<string, string> { ["server"] = parts[0] };
            if (parts.Length > 1) d["serverpath"] = parts[1];
            if (kind == Kind.Nfs)
            {
                d["options"] = options ?? "";
                if (forProbe) d["probeversion"] = ""; // asks the server for its NFS versions
                else if (!string.IsNullOrEmpty(nfsVersion)) d["nfsversion"] = nfsVersion;
            }
            else if ((user ?? "").Trim().Length > 0 || (password ?? "").Trim().Length > 0)
            {
                d["username"] = user;
                d["password"] = password;
            }
            return d;
        }

        public static string AutoDescription(Kind kind, string path) =>
            string.Format(kind == Kind.Nfs ? Messages.NEWSR_ACTION : Messages.NEWSR_CIFS_ACTION, path);

        public static SrProbeAction Probe(IXenConnection c, Kind kind, Dictionary<string, string> dconf) =>
            new SrProbeAction(c, Helpers.GetCoordinator(c), TypeOf(kind), dconf);

        /// <summary>NewSRWizard.GetActions: create when no existing SR was picked, reattach a detached one of this pool, else introduce.</summary>
        public static AsyncAction Finish(IXenConnection c, Kind kind, string name, string description, Dictionary<string, string> dconf, string existingUuid)
        {
            if (string.IsNullOrEmpty(existingUuid))
                return new SrCreateAction(c, Helpers.GetCoordinator(c), name, description, TypeOf(kind), "", dconf, new Dictionary<string, string>());
            var detached = c.Cache.SRs.FirstOrDefault(sr => sr.uuid == existingUuid);
            return detached != null
                ? new SrReattachAction(detached, name, description, dconf)
                : new SrIntroduceAction(c, existingUuid, name, description, TypeOf(kind), "", dconf);
        }

        /// <summary>SrWizardHelpers.SrInUse: an attached SR with this uuid on any connected pool (two pools must never share it).</summary>
        public static SR InUse(string uuid) =>
            ConnectionsManager.XenConnectionsCopy.SelectMany(c => c.Cache.SRs).FirstOrDefault(sr => sr.uuid == uuid && sr.HasPBDs());

        /// <summary>NewSRWizard asks before attaching an existing NFS SR (ShowIntroducePrompt / ShowReattachWarning); SMB does not.</summary>
        public static string AttachWarning(Kind kind, IXenConnection c, string uuid) =>
            kind != Kind.Nfs || string.IsNullOrEmpty(uuid) ? null
                : string.Format(Messages.NEWSR_MULTI_POOL_WARNING, BrandManager.BrandConsole, c.Cache.SRs.FirstOrDefault(sr => sr.uuid == uuid)?.Name() ?? uuid);

        /// <summary>NewSRWizard: a reattach that failed half way leaves PBDs behind; unplug and destroy them.</summary>
        public static AsyncAction RollbackFailedReattach(IXenConnection c, string uuid) =>
            c.Cache.SRs.FirstOrDefault(sr => sr.uuid == uuid) is SR sr && sr.HasPBDs() ? new DetachSrAction(sr, true) : null;
    }
}
