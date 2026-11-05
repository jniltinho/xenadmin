package app

import (
	"encoding/xml"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	xenapi "go/xenapi"

	"xcpng-center/internal/messages"
	"xcpng-center/internal/xapi"
)

// Alert is one row of Notifications > Alerts (NotificationRules.AlertRow, XenAdmin.Desktop.Core/NotificationRules.cs).
type Alert struct {
	UUID        string    `json:"uuid"`
	Priority    int       `json:"priority"` // 1 (highest) to 5; 6: unknown (AlertPriority)
	Title       string    `json:"title"`
	Description string    `json:"description"`
	AppliesTo   string    `json:"applies_to"`
	Time        time.Time `json:"time"`
}

// graphsOnly are Message.ShowOnGraphs: lifecycle events, not alerts; HA_POOL_OVERCOMMITTED is Message.IsSquelched.
var graphsOnly = map[string]bool{"VM_CLONED": true, "VM_CRASHED": true, "VM_REBOOTED": true, "VM_RESUMED": true,
	"VM_SHUTDOWN": true, "VM_STARTED": true, "VM_SUSPENDED": true, "HA_POOL_OVERCOMMITTED": true}

// namedTypes are the MessageAlert.Description types whose friendly text takes the object's name.
var namedTypes = map[string]bool{}

func init() {
	for _, t := range strings.Fields(`HA_HEARTBEAT_APPROACHING_TIMEOUT HA_HOST_FAILED HA_HOST_WAS_FENCED
		HA_PROTECTED_VM_RESTART_FAILED HA_STATEFILE_APPROACHING_TIMEOUT HA_STATEFILE_LOST
		HA_XAPI_HEALTHCHECK_APPROACHING_TIMEOUT LICENSE_DOES_NOT_SUPPORT_POOLING PBD_PLUG_FAILED_ON_SERVER_START VM_CLONED
		VM_CRASHED VM_REBOOTED VM_RESUMED VM_SHUTDOWN VM_STARTED VM_SUSPENDED METADATA_LUN_BROKEN METADATA_LUN_HEALTHY
		LICENSE_SERVER_UNREACHABLE LICENSE_SERVER_VERSION_OBSOLETE GRACE_LICENSE LICENSE_NOT_AVAILABLE LICENSE_EXPIRED
		LICENSE_SERVER_CONNECTED LICENSE_SERVER_UNAVAILABLE HOST_CLOCK_WENT_BACKWARDS POOL_CPU_FEATURES_UP
		POOL_CPU_FEATURES_DOWN HOST_CPU_FEATURES_UP HOST_CPU_FEATURES_DOWN VDI_CBT_RESIZE_FAILED VDI_CBT_SNAPSHOT_FAILED
		VDI_CBT_METADATA_INCONSISTENT CLUSTER_HOST_FENCING CLUSTER_HOST_ENABLE_FAILED VM_SECURE_BOOT_FAILED
		TLS_VERIFICATION_EMERGENCY_DISABLED LICENSE_EXPIRES_SOON`) {
		namedTypes[t] = true
	}
}

// pool is what the texts need about the connection: its name (Pool.Name: the coordinator's when unnamed), whether it
// is a pool rather than a standalone server (Helpers.IsPool), and the names of the objects messages refer to.
type pool struct {
	uuid, name string
	isPool     bool
	names      map[string]string // object uuid -> name
	hosts      map[string]bool   // host uuids
}

func readPool(s *xenapi.Session) (pool, error) {
	p := pool{names: map[string]string{}, hosts: map[string]bool{}}
	pools, err := xenapi.Pool.GetAllRecords(s)
	if err != nil {
		return p, err
	}
	hosts, err := xenapi.Host.GetAllRecords(s)
	if err != nil {
		return p, err
	}
	vms, err := xenapi.VM.GetAllRecords(s)
	if err != nil {
		return p, err
	}
	srs, err := xenapi.SR.GetAllRecords(s)
	if err != nil {
		return p, err
	}
	for ref, h := range hosts {
		p.names[h.UUID], p.hosts[h.UUID] = h.NameLabel, true
		for _, r := range pools {
			if r.Master == ref && p.name == "" {
				p.name = h.NameLabel
			}
		}
	}
	for _, r := range pools {
		p.uuid, p.isPool = r.UUID, r.NameLabel != "" || len(hosts) > 1
		if r.NameLabel != "" {
			p.name = r.NameLabel
		}
		p.names[r.UUID] = p.name
	}
	for _, v := range vms {
		p.names[v.UUID] = v.NameLabel
	}
	for _, r := range srs {
		p.names[r.UUID] = r.NameLabel
	}
	return p, nil
}

// ListAlerts is NotificationRules.AlertMessages + Row, newest first (the Alerts grid's default order).
// shortcut: certificate-expiry titles (CertificateAlert) and the names of VDIs, VMPPs and VMSSs are not resolved yet
// (the friendly name and the raw body are shown); port them with the other object kinds.
func ListAlerts(c *xapi.Connection) ([]Alert, error) {
	s := c.Session
	msgs, err := xenapi.Message.GetAllRecords(s)
	if err != nil {
		return nil, xapi.Sanitize(err)
	}
	p, err := readPool(s)
	if err != nil {
		return nil, xapi.Sanitize(err)
	}
	out := []Alert{}
	for _, m := range msgs {
		if graphsOnly[m.Name] {
			continue
		}
		prio := m.Priority
		if prio < 1 || prio > 5 {
			prio = 6
		}
		applies := p.name
		if applies == "" {
			applies = m.ObjUUID
		}
		out = append(out, Alert{UUID: m.UUID, Priority: prio, Title: alertTitle(m, p), Description: alertDescription(m, p),
			AppliesTo: applies, Time: m.Timestamp})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Time.After(out[j].Time) })
	return out, nil
}

// alertTitle is MessageAlert.Title: the friendly name, prefixed by the object's name unless it is the pool itself
// (or the host of a standalone server).
func alertTitle(m xenapi.MessageRecord, p pool) string {
	if m.Name == "FAILED_LOGIN_ATTEMPTS" && m.ObjUUID == p.uuid {
		return messages.Format(messages.Get("ALERT_FAILED_LOGIN_ATTEMPT_TITLE"), p.name)
	}
	title := messages.Friendly("Message.name-" + strings.ToLower(m.Name))
	if m.Name == "GFS2_CAPACITY" {
		return title
	}
	if title == "" {
		title = m.Name
	}
	if m.Cls != xenapi.ClsPool && (!p.hosts[m.ObjUUID] || p.isPool) {
		if name := p.names[m.ObjUUID]; name != "" {
			title = messages.Format(messages.Get("STRING_COLON_SPACE_STRING"), name, title)
		}
	}
	return title
}

// alertDescription is MessageAlert.Description: the friendly text filled with the object (and pool), else the body.
func alertDescription(m xenapi.MessageRecord, p pool) string {
	if m.Name == "FAILED_LOGIN_ATTEMPTS" {
		if text := failedLoginText(m.Body); text != "" {
			return text
		}
		return m.Body
	}
	name, found := p.names[m.ObjUUID]
	format := messages.Friendly("Message.body-" + strings.ToLower(m.Name))
	if plan, err := strconv.Atoi(m.Body); found && err == nil && m.Name == "HA_POOL_DROP_IN_PLAN_EXISTS_FOR" {
		suffix := "n"
		if plan == 0 || plan == 1 {
			suffix = strconv.Itoa(plan)
		}
		return messages.Format(messages.Friendly("Message.body-ha_pool_drop_in_plan_exists_for-"+suffix), name, plan)
	}
	if !found || format == "" {
		return m.Body
	}
	switch {
	case namedTypes[m.Name]:
		return messages.Format(format, name)
	case m.Name == "HOST_CLOCK_SKEW_DETECTED" || m.Name == "POOL_MASTER_TRANSITION":
		return messages.Format(format, name, p.name)
	case m.Name == "VBD_QOS_FAILED" || m.Name == "VCPU_QOS_FAILED" || m.Name == "VIF_QOS_FAILED":
		return messages.Format(format, "", name)
	case m.Name == "EXTAUTH_IN_POOL_IS_NON_HOMOGENEOUS":
		return messages.Format(format, p.name)
	}
	return m.Body
}

// failedLoginText is NotificationRules.FailedLoginText: the body's <total>, <unknown> and <known> sources.
func failedLoginText(body string) string {
	var doc struct {
		Total   int `xml:"total"`
		Unknown int `xml:"unknown"`
		Known   []struct {
			IP         string `xml:"ip"`
			Number     int    `xml:"number"`
			Date       string `xml:"date"`
			Username   string `xml:"username"`
			Originator string `xml:"originator"`
			UserAgent  string `xml:"useragent"`
		} `xml:"known"`
	}
	if xml.Unmarshal([]byte(body), &doc) != nil {
		return ""
	}
	if doc.Total == 0 && doc.Unknown > 0 {
		return messages.Format(messages.Get("ALERT_FAILED_LOGIN_ATTEMPT_DESCRIPTION_ALL_UNKNOWN"), doc.Unknown)
	}
	var lines []string
	if doc.Total > 0 {
		lines = append(lines, messages.Format(messages.Get("ALERT_FAILED_LOGIN_ATTEMPT_DESCRIPTION_TOTAL"), doc.Total))
	}
	if doc.Unknown > 0 {
		lines = append(lines, messages.Format(messages.Get("ALERT_FAILED_LOGIN_ATTEMPT_DESCRIPTION_UNKNOWN"), doc.Unknown))
	}
	if len(doc.Known) > 0 {
		lines = append(lines, messages.Get("ALERT_FAILED_LOGIN_ATTEMPT_DESCRIPTION_COMMON"))
	}
	for _, k := range doc.Known {
		var e []string
		add := func(key, v string) {
			if v != "" {
				e = append(e, messages.Format(messages.Get("ALERT_FAILED_LOGIN_ATTEMPT_OFFENDER_"+key), v))
			}
		}
		add("IP", k.IP)
		if k.Number > 0 {
			add("NUMBER", strconv.Itoa(k.Number))
		}
		if t, ok := iso8601(k.Date); ok {
			add("DATE", t.Local().Format("Jan 2, 2006 3:04 PM")) // Messages.DATEFORMAT_DMY_HM
		}
		add("USERNAME", k.Username)
		add("ORIGINATOR", k.Originator)
		add("USERAGENT", k.UserAgent)
		lines = append(lines, "- "+strings.Join(e, ", "))
	}
	return strings.Join(lines, "\n")
}

// iso8601 is Util.TryParseIso8601DateTime (Util.Iso8601DateFormats).
func iso8601(s string) (time.Time, bool) {
	for _, layout := range []string{"20060102T15:04:05Z", "2006-01-02T15:04:05Z"} {
		if t, err := time.Parse(layout, s); err == nil {
			return t, true
		}
	}
	return time.Time{}, false
}

// DismissAlerts is NotificationRules.Dismiss: Message.destroy for each, ignoring the ones already gone.
func DismissAlerts(c *xapi.Connection, uuids []string) error {
	var errs []error
	for _, u := range uuids {
		ref, err := xenapi.Message.GetByUUID(c.Session, u)
		if err == nil {
			err = xenapi.Message.Destroy(c.Session, ref)
		}
		var api *xenapi.APIError
		if errors.As(err, &api) && (api.Message == "UUID_INVALID" || api.Message == "HANDLE_INVALID") {
			continue
		}
		if err != nil {
			errs = append(errs, fmt.Errorf("%s: %w", u, xapi.Sanitize(err)))
		}
	}
	return errors.Join(errs...)
}
