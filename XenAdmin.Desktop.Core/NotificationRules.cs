using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Xml;
using XenAdmin.Actions;
using XenAdmin.Alerts;
using XenAdmin.Core;
using XenAdmin.Network;
using XenAPI;

namespace XenAdmin.Desktop
{
    /// <summary>
    /// Notifications > Alerts (MainWindow.MessageCollectionChanged, MessageAlert, AlertSummaryPage): the servers' messages
    /// that are not lifecycle events, with their friendly title and text, severity, server/pool and date.
    /// shortcut: no "fix" links (they open WinForms dialogs such as HA configuration or SR repair) and no alarm, VMSS,
    /// update-expiry or leaf-coalesce texts (their raw body is shown); port them with those features.
    /// </summary>
    public static class NotificationRules
    {
        public sealed record AlertRow(Message Message, AlertPriority Priority, string Title, string Description, string AppliesTo, DateTime Time);

        /// <summary>The messages shown as alerts: not lifecycle events (those are on the graphs) and not squelched.</summary>
        public static IEnumerable<Message> AlertMessages(IEnumerable<IXenConnection> connections) =>
            connections.Where(c => c.IsConnected).SelectMany(c => c.Cache.Messages).Where(m => !m.ShowOnGraphs() && !m.IsSquelched());

        /// <summary>MainWindow.UpdateAlertCount: counts without formatting.</summary>
        public static int Count(IEnumerable<IXenConnection> connections) => AlertMessages(connections).Count();

        /// <summary>One grid row; format only the rows shown (the texts resolve objects by uuid).</summary>
        /// <summary>FilterDatesToolStripDropDownButton's presets (the custom range of DateFilterDialog is not ported yet).</summary>
        public enum DateFilter { All, Last24Hours, Last3Days, Last7Days, Last30Days }

        public static string DateFilterText(DateFilter f) => f switch
        {
            DateFilter.Last24Hours => Messages.DATE_FILTER_LAST_24_HOURS,
            DateFilter.Last3Days => Messages.DATE_FILTER_LAST_3_DAYS,
            DateFilter.Last7Days => Messages.DATE_FILTER_LAST_7_DAYS,
            DateFilter.Last30Days => Messages.DATE_FILTER_LAST_30_DAYS,
            _ => Messages.FILTER_SHOW_ALL
        };

        /// <summary>DateFilterDialog.SetDays: the last 24 hours from now, or from the beginning of the day 3, 7 or 30 days ago.</summary>
        public static DateTime? DateFrom(DateFilter f, DateTime now) => f switch
        {
            DateFilter.Last24Hours => now.AddDays(-1),
            DateFilter.Last3Days => now.AddDays(-3).Date,
            DateFilter.Last7Days => now.AddDays(-7).Date,
            DateFilter.Last30Days => now.AddDays(-30).Date,
            _ => null
        };

        /// <summary>DateFilterDialog's start and end choices.</summary>
        public static readonly string[] DateBoundChoices = { "Not set", "24 hours ago", "3 days ago", "7 days ago", "30 days ago", "Custom" };

        public const int CustomBound = 5;

        /// <summary>
        /// DateFilterDialog.parseStartDate / parseEndDate: choice 0 is not set (null), 1 is 24 hours ago, 2..4 the
        /// beginning (start) or end (end) of the day 3, 7 or 30 days ago, 5 the custom date and time.
        /// </summary>
        public static DateTime? DateBound(int choice, bool start, DateTime now, DateTime custom) => choice switch
        {
            1 => now.AddDays(-1),
            2 => Day(now.AddDays(-3), start),
            3 => Day(now.AddDays(-7), start),
            4 => Day(now.AddDays(-30), start),
            CustomBound => custom,
            _ => null
        };

        private static DateTime Day(DateTime d, bool start) => start ? d.Date : d.Date.AddDays(1).AddTicks(-1);

        /// <summary>DateFilterDialog.validateInterval: the end must be after the start (unset bounds always pass).</summary>
        public static bool ValidRange(DateTime? from, DateTime? to) => from == null || to == null || from < to;

        /// <summary>FilterDatesToolStripDropDownButton.HideByDate.</summary>
        public static bool HideByDate(DateTime time, DateTime? from, DateTime? to) => from > time || to < time;

        /// <summary>AlertSummaryPage's CSV export: the header line.</summary>
        public static string CsvHeader() =>
            string.Join(",", Messages.TITLE, Messages.SEVERITY, Messages.DESCRIPTION, Messages.APPLIES_TO, Messages.TIMESTAMP);

        /// <summary>AlertExtensions.GetAlertDetailsCSVQuotes: every field quoted, quotes doubled, local time.</summary>
        public static string CsvLine(AlertRow r) =>
            string.Format("\"{0}\",\"{1}\",\"{2}\",\"{3}\",\"{4}\"",
                (r.Title ?? "").EscapeQuotes(), r.Priority.GetString().EscapeQuotes(), (r.Description ?? "").EscapeQuotes(),
                (r.AppliesTo ?? "").EscapeQuotes(), r.Time.ToString(Messages.DATEFORMAT_DMY_HM).EscapeQuotes());

        public static AlertRow Row(Message m) => new AlertRow(m, PriorityOf(m), Title(m), Description(m), AppliesTo(m), m.TimestampLocal());

        /// <summary>MessageAlert.HostUuid: the message's host, else the coordinator (the Filter by Server key).</summary>
        public static string HostUuid(Message m) =>
            (m.GetXenObject() as Host ?? Helpers.GetCoordinator(m.Connection))?.uuid;

        public static AlertPriority PriorityOf(Message m) =>
            m.priority >= 0 && m.priority <= int.MaxValue && Enum.IsDefined(typeof(AlertPriority), (int)m.priority) ? (AlertPriority)(int)m.priority : AlertPriority.Unknown;

        /// <summary>Images.GetIconFor(AlertPriority).</summary>
        public static string Icon(AlertPriority p) => p switch
        {
            AlertPriority.Priority1 => "alert1_16", AlertPriority.Priority2 => "alert2_16", AlertPriority.Priority3 => "alert3_16",
            AlertPriority.Priority4 => "alert4_16", AlertPriority.Priority5 => "alert5_16", _ => "alert6_16"
        };

        /// <summary>MessageAlert.Title: the friendly name, prefixed by the object's name unless it is the pool itself.</summary>
        public static string Title(Message m)
        {
            if (CertificateTitle(m) is string certificate)
                return certificate;
            if (m.Type == Message.MessageType.FAILED_LOGIN_ATTEMPTS
                && m.Connection?.Cache.Pools.FirstOrDefault(p => p.uuid == m.obj_uuid) is Pool loginPool)
                return string.Format(Messages.ALERT_FAILED_LOGIN_ATTEMPT_TITLE, loginPool.Name());
            if (m.name == "GFS2_CAPACITY")
                return Message.FriendlyName(m.name);
            var title = Message.FriendlyName(m.MessageTypeString());
            if (string.IsNullOrEmpty(title))
                title = m.name;
            if (m.cls != cls.Pool)
            {
                var o = m.GetXenObject();
                if (!(o is Host h) || Helpers.IsPool(h.Connection))
                {
                    var name = Helpers.GetName(o);
                    if (!string.IsNullOrEmpty(name))
                        title = string.Format(Messages.STRING_COLON_SPACE_STRING, name, title);
                }
            }
            return title;
        }

        // MessageAlert.Description: the types whose friendly text takes the object's name
        private static readonly HashSet<Message.MessageType> NamedTypes = new()
        {
            Message.MessageType.HA_HEARTBEAT_APPROACHING_TIMEOUT, Message.MessageType.HA_HOST_FAILED, Message.MessageType.HA_HOST_WAS_FENCED,
            Message.MessageType.HA_PROTECTED_VM_RESTART_FAILED, Message.MessageType.HA_STATEFILE_APPROACHING_TIMEOUT, Message.MessageType.HA_STATEFILE_LOST,
            Message.MessageType.HA_XAPI_HEALTHCHECK_APPROACHING_TIMEOUT, Message.MessageType.LICENSE_DOES_NOT_SUPPORT_POOLING,
            Message.MessageType.PBD_PLUG_FAILED_ON_SERVER_START, Message.MessageType.VM_CLONED, Message.MessageType.VM_CRASHED,
            Message.MessageType.VM_REBOOTED, Message.MessageType.VM_RESUMED, Message.MessageType.VM_SHUTDOWN, Message.MessageType.VM_STARTED,
            Message.MessageType.VM_SUSPENDED, Message.MessageType.METADATA_LUN_BROKEN, Message.MessageType.METADATA_LUN_HEALTHY,
            Message.MessageType.LICENSE_SERVER_UNREACHABLE, Message.MessageType.LICENSE_SERVER_VERSION_OBSOLETE, Message.MessageType.GRACE_LICENSE,
            Message.MessageType.LICENSE_NOT_AVAILABLE, Message.MessageType.LICENSE_EXPIRED, Message.MessageType.LICENSE_SERVER_CONNECTED,
            Message.MessageType.LICENSE_SERVER_UNAVAILABLE, Message.MessageType.HOST_CLOCK_WENT_BACKWARDS, Message.MessageType.POOL_CPU_FEATURES_UP,
            Message.MessageType.POOL_CPU_FEATURES_DOWN, Message.MessageType.HOST_CPU_FEATURES_UP, Message.MessageType.HOST_CPU_FEATURES_DOWN,
            Message.MessageType.VDI_CBT_RESIZE_FAILED, Message.MessageType.VDI_CBT_SNAPSHOT_FAILED, Message.MessageType.VDI_CBT_METADATA_INCONSISTENT,
            Message.MessageType.CLUSTER_HOST_FENCING, Message.MessageType.CLUSTER_HOST_ENABLE_FAILED, Message.MessageType.VM_SECURE_BOOT_FAILED,
            Message.MessageType.TLS_VERIFICATION_EMERGENCY_DISABLED, Message.MessageType.LICENSE_EXPIRES_SOON
        };

        /// <summary>MessageAlert.Description: the friendly text filled with the object (and pool), else the raw body.</summary>
        public static string Description(Message m)
        {
            var o = m.GetXenObject();
            if (CertificateText(m, o) is string certificate)
                return certificate;
            if (m.Type == Message.MessageType.FAILED_LOGIN_ATTEMPTS)
                return FailedLoginText(m.body) ?? m.body;
            var format = Message.FriendlyBody(m.MessageTypeString());
            try
            {
                if (o != null && m.Type is Message.MessageType.HA_POOL_DROP_IN_PLAN_EXISTS_FOR or Message.MessageType.HA_POOL_OVERCOMMITTED
                    && int.TryParse(m.body, out var plan))
                    return string.Format(Message.FriendlyBody("ha_pool_drop_in_plan_exists_for-" + (plan == 0 ? "0" : plan == 1 ? "1" : "n")),
                        Helpers.GetName(o), plan);
                if (o != null && !string.IsNullOrEmpty(format))
                {
                    if (NamedTypes.Contains(m.Type))
                        return string.Format(format, Helpers.GetName(o));
                    if (m.Type is Message.MessageType.HOST_CLOCK_SKEW_DETECTED or Message.MessageType.POOL_MASTER_TRANSITION
                        && Helpers.GetPoolOfOne(o.Connection) is Pool pool)
                        return string.Format(format, Helpers.GetName(o), pool.Name());
                    if (m.Type is Message.MessageType.VBD_QOS_FAILED or Message.MessageType.VCPU_QOS_FAILED or Message.MessageType.VIF_QOS_FAILED)
                        return string.Format(format, "", Helpers.GetName(o));
                    if (m.Type == Message.MessageType.EXTAUTH_IN_POOL_IS_NON_HOMOGENEOUS)
                        return string.Format(format, Helpers.GetName(Helpers.GetPoolOfOne(o.Connection)));
                }
            }
            catch (FormatException)
            {
                // a friendly text with other placeholders: the body says it as well
            }
            return m.body;
        }

        private static string GuiDate(DateTime? date) =>
            date.HasValue ? DateText(date.Value.ToLocalTime()) : "";

        // HelpersGUI.DateTimeToString(..., DATEFORMAT_DMY_HM, true) without WinForms
        private static string DateText(DateTime d) => d.ToString(Messages.DATEFORMAT_DMY_HM, System.Globalization.CultureInfo.CurrentCulture);

        private static XmlDocument Body(string body)
        {
            try
            {
                var doc = new XmlDocument();
                doc.LoadXml(body);
                return doc;
            }
            catch (XmlException)
            {
                return null;
            }
        }

        /// <summary>CertificateAlert: the kind of certificate ("ca", "internal", "host"; null for other messages) and the
        /// expiry date from the body's &lt;date&gt;.</summary>
        private static (string Kind, DateTime? Expiry) Certificate(Message m)
        {
            var t = m.Type.ToString();
            var kind = t.StartsWith("POOL_CA_CERTIFICATE_") ? "ca" : t.StartsWith("HOST_INTERNAL_CERTIFICATE_") ? "internal"
                : t.StartsWith("HOST_SERVER_CERTIFICATE_") ? "host" : null;
            if (kind == null)
                return (null, null);
            if (Body(m.body)?.GetElementsByTagName("date") is XmlNodeList dates && dates.Count > 0
                && Util.TryParseIso8601DateTime(dates[0].InnerText, out var d))
                return (kind, d);
            return (kind, null);
        }

        /// <summary>CertificateAlert.Title: expired, or expiring in days, hours or minutes from the message's time.</summary>
        public static string CertificateTitle(Message m)
        {
            var (kind, expiry) = Certificate(m);
            if (kind == null)
                return null;
            var o = m.GetXenObject();
            var name = o == null ? Messages.UNKNOWN_OBJECT : o.Name();
            var pool = Helpers.GetPoolOfOne(m.Connection)?.Name();
            var (expired, days, hours, minutes) = kind switch
            {
                "ca" => (Messages.CERTIFICATE_CA_ALERT_EXPIRED_TITLE, Messages.CERTIFICATE_CA_ALERT_EXPIRING_TITLE_DAYS,
                    Messages.CERTIFICATE_CA_ALERT_EXPIRING_TITLE_HOURS, Messages.CERTIFICATE_CA_ALERT_EXPIRING_TITLE_MINUTES),
                "internal" => (Messages.CERTIFICATE_HOST_INTERNAL_ALERT_EXPIRED_TITLE, Messages.CERTIFICATE_HOST_INTERNAL_ALERT_EXPIRING_TITLE_DAYS,
                    Messages.CERTIFICATE_HOST_INTERNAL_ALERT_EXPIRING_TITLE_HOURS, Messages.CERTIFICATE_HOST_INTERNAL_ALERT_EXPIRING_TITLE_MINUTES),
                _ => (Messages.CERTIFICATE_HOST_ALERT_EXPIRED_TITLE, Messages.CERTIFICATE_HOST_ALERT_EXPIRING_TITLE_DAYS,
                    Messages.CERTIFICATE_HOST_ALERT_EXPIRING_TITLE_HOURS, Messages.CERTIFICATE_HOST_ALERT_EXPIRING_TITLE_MINUTES)
            };
            // the CA texts also name the pool: {0} object, {1} pool, then the amount
            string Format(string f, double? n = null) => kind == "ca"
                ? (n.HasValue ? string.Format(f, name, pool, n) : string.Format(f, name, pool))
                : (n.HasValue ? string.Format(f, name, n) : string.Format(f, name));
            if (!m.Type.ToString().EndsWith("_EXPIRED") && expiry > m.timestamp)
            {
                var eta = expiry.Value - m.timestamp;
                if (eta.TotalDays >= 1) return Format(days, Math.Round(eta.TotalDays, MidpointRounding.AwayFromZero));
                if (eta.TotalHours >= 1) return Format(hours, Math.Round(eta.TotalHours, MidpointRounding.AwayFromZero));
                if (eta.TotalMinutes >= 1) return Format(minutes, Math.Round(eta.TotalMinutes, MidpointRounding.AwayFromZero));
            }
            return Format(expired);
        }

        /// <summary>CertificateAlert.Description: the expiry date from the body's &lt;date&gt;.</summary>
        public static string CertificateText(Message m, IXenObject o)
        {
            var (kind, expiry) = Certificate(m);
            if (kind == null)
                return null;
            var name = o == null ? Messages.UNKNOWN_OBJECT : o.Name();
            var expired = m.Type.ToString().EndsWith("_EXPIRED");
            var pool = Helpers.GetPoolOfOne(m.Connection)?.Name();
            return kind switch
            {
                "ca" => string.Format(expired ? Messages.CERTIFICATE_CA_ALERT_EXPIRED_DESCIRPTION : Messages.CERTIFICATE_CA_ALERT_EXPIRING_DESCRIPTION, name, pool, GuiDate(expiry)),
                "internal" => string.Format(expired ? Messages.CERTIFICATE_HOST_INTERNAL_ALERT_EXPIRED_DESCIRPTION : Messages.CERTIFICATE_HOST_INTERNAL_ALERT_EXPIRING_DESCRIPTION, name, GuiDate(expiry)),
                _ => string.Format(expired ? Messages.CERTIFICATE_HOST_ALERT_EXPIRED_DESCIRPTION : Messages.CERTIFICATE_HOST_ALERT_EXPIRING_DESCRIPTION, name, GuiDate(expiry))
            };
        }

        /// <summary>FailedLoginAttemptAlert.Description: totals and the most common sources; null if the body is not that XML.</summary>
        public static string FailedLoginText(string body)
        {
            var doc = Body(body);
            if (doc == null)
                return null;
            int Number(string tag) => doc.GetElementsByTagName(tag) is XmlNodeList n && n.Count > 0 && int.TryParse(n[0].InnerText, out var v) ? v : 0;
            int total = Number("total"), unknown = Number("unknown");
            if (total == 0 && unknown > 0)
                return string.Format(Messages.ALERT_FAILED_LOGIN_ATTEMPT_DESCRIPTION_ALL_UNKNOWN, unknown);
            var sb = new StringBuilder();
            if (total > 0)
                sb.AppendLine(string.Format(Messages.ALERT_FAILED_LOGIN_ATTEMPT_DESCRIPTION_TOTAL, total));
            if (unknown > 0)
                sb.AppendLine(string.Format(Messages.ALERT_FAILED_LOGIN_ATTEMPT_DESCRIPTION_UNKNOWN, unknown));
            var known = doc.GetElementsByTagName("known");
            if (known.Count > 0)
                sb.AppendLine(Messages.ALERT_FAILED_LOGIN_ATTEMPT_DESCRIPTION_COMMON);
            foreach (XmlNode node in known)
            {
                var entries = new List<string>();
                string Child(string tag) => node[tag]?.InnerText;
                if (!string.IsNullOrEmpty(Child("ip"))) entries.Add(string.Format(Messages.ALERT_FAILED_LOGIN_ATTEMPT_OFFENDER_IP, Child("ip")));
                if (int.TryParse(Child("number"), out var n) && n > 0) entries.Add(string.Format(Messages.ALERT_FAILED_LOGIN_ATTEMPT_OFFENDER_NUMBER, n));
                if (Util.TryParseIso8601DateTime(Child("date") ?? "", out var when)) entries.Add(string.Format(Messages.ALERT_FAILED_LOGIN_ATTEMPT_OFFENDER_DATE, GuiDate(when)));
                if (!string.IsNullOrEmpty(Child("username"))) entries.Add(string.Format(Messages.ALERT_FAILED_LOGIN_ATTEMPT_OFFENDER_USERNAME, Child("username")));
                if (!string.IsNullOrEmpty(Child("originator"))) entries.Add(string.Format(Messages.ALERT_FAILED_LOGIN_ATTEMPT_OFFENDER_ORIGINATOR, Child("originator")));
                if (!string.IsNullOrEmpty(Child("useragent"))) entries.Add(string.Format(Messages.ALERT_FAILED_LOGIN_ATTEMPT_OFFENDER_USERAGENT, Child("useragent")));
                sb.AppendLine($"- {string.Join(", ", entries)}");
            }
            return sb.Length > 0 ? sb.ToString().TrimEnd() : null;
        }

        /// <summary>MessageAlert.AppliesTo: the pool (or standalone server) the message came from.</summary>
        public static string AppliesTo(Message m)
        {
            var name = Helpers.GetName(Helpers.GetPoolOfOne(m.Connection));
            return !string.IsNullOrEmpty(name) ? name : m.obj_uuid;
        }

        /// <summary>MessageAlert.AllowedToDismiss: Message.destroy needs pool-operator rights unless logged in as root.</summary>
        public static bool CanDismiss(Message m)
        {
            var c = m.Connection;
            if (c?.Session == null) return false;
            if (c.Session.IsLocalSuperuser) return true;
            var roles = Role.ValidRoleList("Message.destroy", c);
            return roles.Any(r => c.Session.Roles.Contains(r));
        }

        /// <summary>Dismiss / Dismiss All: Message.destroy for each, ignoring ones already gone (MessageAlert.Dismiss).</summary>
        public static AsyncAction Dismiss(IXenConnection c, IReadOnlyCollection<Message> messages) =>
            new DelegatedAsyncAction(c, messages.Count == 1 ? "Dismissing alert" : $"Dismissing {messages.Count} alerts", "", "", session =>
            {
                foreach (var m in messages)
                {
                    try { Message.destroy(session, m.opaque_ref); }
                    catch (Failure f) when (f.ErrorDescription.Count > 0 && f.ErrorDescription[0] == Failure.HANDLE_INVALID) { }
                }
            }, true, "message.destroy");
    }
}
