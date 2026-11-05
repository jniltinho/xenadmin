using System.Collections.Generic;
using System.Linq;
using System.Text.RegularExpressions;
using XenAdmin.Actions;
using XenAdmin.Core;
using XenAdmin.Network;
using XenAPI;

namespace XenAdmin.Desktop
{
    /// <summary>
    /// New Storage Repository wizard (NewSRWizard with the VHDoNFS, CifsFrontend and LVMoISCSI pages) for shared virtual
    /// disk storage.
    /// shortcut: iSCSI is LVM only (no GFS2: it needs clustering and a second probe); add with clustered pools.
    /// </summary>
    public static class SrRules
    {
        /// <summary>Virtual disk storage (NFS, SMB, iSCSI) and ISO libraries (NFS ISO, SMB ISO).</summary>
        public enum Kind { Nfs, Smb, Iscsi, NfsIso, SmbIso }

        public static bool IsIso(Kind kind) => kind is Kind.NfsIso or Kind.SmbIso;
        /// <summary>The share is written server:/path (NFS) or \\server\share (SMB).</summary>
        public static bool IsNfsPath(Kind kind) => kind is Kind.Nfs or Kind.NfsIso;

        public static SR.SRTypes TypeOf(Kind kind) => kind switch
        {
            Kind.Nfs => SR.SRTypes.nfs, Kind.Smb => SR.SRTypes.smb, Kind.Iscsi => SR.SRTypes.lvmoiscsi, _ => SR.SRTypes.iso
        };
        /// <summary>SrWizardType.ContentType: "iso" for the libraries, "" for virtual disk storage.</summary>
        public static string ContentType(Kind kind) => IsIso(kind) ? SR.Content_Type_ISO : "";
        public static string TypeName(Kind kind) => kind switch
        {
            Kind.Nfs => Messages.NEWSR_VHDONFS_TYPE_NAME, Kind.Smb => Messages.NEWSR_CIFS_TYPE_NAME, Kind.Iscsi => Messages.NEWSR_LVMOISCSI_TYPE_NAME,
            Kind.NfsIso => Messages.NEWSR_NFS_ISO_TYPE_NAME, _ => Messages.NEWSR_CIFS_ISO_TYPE_NAME
        };
        public static string Blurb(Kind kind) => kind switch
        {
            Kind.Nfs => Messages.NEWSR_VHDONFS_BLURB, Kind.Smb => Messages.NEWSR_CIFS_BLURB, Kind.Iscsi => Messages.NEWSR_LVMOISCSI_BLURB,
            Kind.NfsIso => Messages.NEWSR_NFS_ISO_BLURB, _ => Messages.NEWSR_CIFS_ISO_BLURB
        };

        /// <summary>SrWizardType.ResetSrName: "NFS virtual disk storage", "(1)" suffix when taken.</summary>
        public static string DefaultName(Kind kind, IXenConnection c)
        {
            var name = kind switch
            {
                Kind.Nfs => Messages.SRWIZARD_NFS_STORAGE, Kind.Smb => Messages.SRWIZARD_CIFS_STORAGE, Kind.Iscsi => Messages.SRWIZARD_ISCSI_STORAGE,
                Kind.NfsIso => Messages.SRWIZARD_NFS_LIBRARY, _ => Messages.SRWIZARD_CIFS_LIBRARY
            };
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
            path != null && (IsNfsPath(kind)
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

        public static string AutoDescription(Kind kind, string path) => string.Format(kind switch
        {
            Kind.Nfs => Messages.NEWSR_ACTION, Kind.NfsIso => Messages.NEWSR_ISO_DESCRIPTION, Kind.SmbIso => Messages.NEWSR_CIF_DESCRIPTION, _ => Messages.NEWSR_CIFS_ACTION
        }, path);

        /// <summary>NFS_ISO / CIFS_ISO DeviceConfig: location and type; NFSv4 on request; an SMB sub-folder goes to iso_path.</summary>
        public static Dictionary<string, string> IsoDeviceConfig(Kind kind, string path, bool nfs4, string user, string password)
        {
            if (kind == Kind.NfsIso)
            {
                var nfs = new Dictionary<string, string> { ["location"] = path.Trim(), ["type"] = "nfs_iso" };
                if (nfs4) nfs["nfsversion"] = "4";
                return nfs;
            }
            var d = new Dictionary<string, string> { ["location"] = path.Replace('\\', '/'), ["type"] = "cifs" };
            var bits = d["location"].Split('/');
            if (bits.Length > 4)
            {
                d["location"] = $"//{bits[2]}/{bits[3]}";
                d["iso_path"] = "/" + string.Join("/", bits, 4, bits.Length - 4);
            }
            if (user != null)
            {
                d["username"] = user;
                d["cifspassword"] = password;
            }
            return d;
        }

        /// <summary>NFS_ISO.IsSRNameUnique / CIFS_ISO.IsIsoStorageAlreadyAttached: this pool already has a library on this share.</summary>
        public static string IsoAlreadyAttached(IXenConnection c, Kind kind, string path)
        {
            var type = kind == Kind.NfsIso ? "nfs_iso" : "cifs";
            var location = kind == Kind.NfsIso ? path?.Trim() : IsoDeviceConfig(kind, path ?? "", false, null, null)["location"];
            var attached = c.Cache.SRs.Where(sr => sr.GetSRType(true) == SR.SRTypes.iso).SelectMany(sr => c.ResolveAll(sr.PBDs))
                .Any(pbd => pbd.device_config.TryGetValue("type", out var t) && t == type && pbd.device_config.TryGetValue("location", out var l) && l == location);
            return !attached ? null : string.Format(kind == Kind.NfsIso ? Messages.NFS_ISO_ALREADY_ATTACHED : Messages.SMB_ISO_STORAGE_ALREADY_ATTACHED, c.FriendlyName);
        }

        public static SrProbeAction Probe(IXenConnection c, Kind kind, Dictionary<string, string> dconf) =>
            new SrProbeAction(c, Helpers.GetCoordinator(c), TypeOf(kind), dconf);

        /// <summary>NewSRWizard.GetActions: create when no existing SR was picked, reattach a detached one of this pool, else introduce.</summary>
        public static AsyncAction Finish(IXenConnection c, Kind kind, string name, string description, Dictionary<string, string> dconf, string existingUuid)
        {
            if (string.IsNullOrEmpty(existingUuid))
                return new SrCreateAction(c, Helpers.GetCoordinator(c), name, description, TypeOf(kind), ContentType(kind), dconf, new Dictionary<string, string>());
            var detached = c.Cache.SRs.FirstOrDefault(sr => sr.uuid == existingUuid);
            return detached != null
                ? new SrReattachAction(detached, name, description, dconf)
                : new SrIntroduceAction(c, existingUuid, name, description, TypeOf(kind), ContentType(kind), dconf);
        }

        // ---------- iSCSI (LVMoISCSI page) ----------

        public const ushort DefaultIscsiPort = 3260; // Util.DEFAULT_ISCSI_PORT

        public static ISCSIPopulateIQNsAction ScanIqns(IXenConnection c, string host, ushort port, string chapUser, string chapPassword) =>
            new ISCSIPopulateIQNsAction(c, host, port, chapUser, chapPassword);

        public static ISCSIPopulateLunsAction ScanLuns(IXenConnection c, IScsiIqnInfo iqn, string chapUser, string chapPassword) =>
            new ISCSIPopulateLunsAction(c, iqn.IpAddress, iqn.Port, iqn.TargetIQN, chapUser, chapPassword);

        /// <summary>The IQN combo's text: "iqn... (address:port)".</summary>
        public static string IqnLabel(IScsiIqnInfo i) => $"{i.TargetIQN} ({i.IpAddress}:{i.Port})";

        /// <summary>The LUN combo's text: "LUN 1: serial: 4 GB (vendor)".</summary>
        public static string LunLabel(ISCSIInfo i)
        {
            var label = "LUN";
            if (i.LunID != -1) label += $" {i.LunID}";
            if (!string.IsNullOrEmpty(i.Serial)) label += $": {i.Serial}";
            if (i.Size >= 0) label += $": {Util.DiskSizeString(i.Size)}";
            if (!string.IsNullOrEmpty(i.Vendor)) label += $" ({i.Vendor})";
            return label;
        }

        /// <summary>LVMoISCSI.GetDeviceConfig: the IQN's own address (multi-homing, CA-11607), the LUN by SCSI id or serial.</summary>
        public static Dictionary<string, string> IscsiDeviceConfig(IScsiIqnInfo iqn, ISCSIInfo lun, string chapUser, string chapPassword)
        {
            var d = new Dictionary<string, string> { ["target"] = iqn.IpAddress, ["port"] = iqn.Port.ToString(), ["targetIQN"] = iqn.TargetIQN };
            if (lun.LunID == -1) d["LUNSerial"] = lun.Serial;
            else d["SCSIid"] = lun.ScsiID;
            if (chapUser != null)
            {
                d["chapuser"] = chapUser;
                d["chappassword"] = chapPassword;
            }
            return d;
        }

        public static string IscsiDescription(IScsiIqnInfo iqn, ISCSIInfo lun) =>
            string.Format(Messages.NEWSR_ISCSI_DESCRIPTION, iqn.IpAddress, iqn.TargetIQN, LunLabel(lun));

        /// <summary>LVMoISCSI.SrUsingLunExists: an iSCSI SR of any connected pool already on this LUN.</summary>
        public static SR LunInUse(ISCSIInfo lun) =>
            ConnectionsManager.XenConnectionsCopy.SelectMany(c => c.Cache.SRs).FirstOrDefault(sr =>
                sr.GetSRType(false) is SR.SRTypes.lvmoiscsi or SR.SRTypes.gfs2 && sr.PBDs.Count > 0
                && sr.Connection.Resolve(sr.PBDs[0]) is PBD pbd && pbd.device_config.TryGetValue("SCSIid", out var id) && id == lun.ScsiID);

        /// <summary>LVMoISCSI.HandleFailure: CHAP login, unknown host, no connection, else the server's message.</summary>
        public static string IscsiError(ActionBase action)
        {
            var code = action.Exception is Failure f && f.ErrorDescription.Count > 0 ? f.ErrorDescription[0] : null;
            return code switch
            {
                "SR_BACKEND_FAILURE_68" => Messages.LOGGING_IN_TO_THE_ISCSI_TARGET_FAILED,
                "SR_BACKEND_FAILURE_140" => Messages.INVALID_HOST,
                "SR_BACKEND_FAILURE_141" => Messages.SR_UNABLE_TO_CONNECT_TO_SCSI_TARGET,
                _ => action.Exception?.Message ?? Messages.SR_UNABLE_TO_CONNECT_TO_SCSI_TARGET
            };
        }

        // ---------- SR commands (DetachSRCommand, ForgetSRCommand, DestroySRCommand, SetAsDefaultSRCommand) ----------

        public static bool CanDetach(SR sr) => sr != null && !sr.IsDetached() && !sr.HasRunningVMs() && !sr.Locked;

        public static bool CanForget(SR sr) => sr != null && !sr.HasRunningVMs() && sr.CanCreateWithXenCenter()
                                               && sr.allowed_operations.Contains(storage_operations.forget) && !sr.Locked;

        public static bool CanDestroy(SR sr) => sr != null && !sr.HasRunningVMs() && sr.CanCreateWithXenCenter()
                                                && sr.allowed_operations.Contains(storage_operations.destroy) && !sr.Locked;

        public static bool CanSetDefault(SR sr) => sr != null && sr.HasPBDs() && !SR.IsDefaultSr(sr) && sr.SupportsVdiCreate()
                                                   && (sr.shared || sr.Connection.Cache.HostCount <= 1) && !sr.Locked;

        /// <summary>The commands' confirmation (title, text, yes button); Set as Default asks nothing.</summary>
        public static (string Title, string Text, string Yes) Confirmation(SR sr, string command) => command switch
        {
            "detach" => (Messages.MESSAGEBOX_DETACH_SR_CONTINUE_TITLE, Messages.MESSAGEBOX_DETACH_SR_CONTINUE, Messages.YES),
            "forget" => (Messages.MESSAGEBOX_FORGET_SR_CONTINUE_TITLE, string.Format(Messages.MESSAGEBOX_FORGET_SR_CONTINUE, sr.Name()), Messages.MESSAGEBOX_FORGET_SR_CONTINUE_YES_BUTTON_LABEL),
            "destroy" => (Messages.MESSAGEBOX_DESTROY_SR_CONTINUE_TITLE, string.Format(Messages.MESSAGEBOX_DESTROY_SR_CONTINUE, sr.Name()), Messages.MESSAGEBOX_DESTROY_SR_YES_BUTTON_LABEL),
            _ => (null, null, null)
        };

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
