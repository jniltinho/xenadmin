package ops

import (
	"errors"
	"maps"
	"regexp"

	xenapi "go/xenapi"

	"xcpng-center/internal/messages"
	"xcpng-center/internal/xapi"
)

// A server's Log Destination and Multipathing pages (LogDestinationEditPage, HostMultipathPage, EditMultipathAction).

// syslogRegex is LogDestinationEditPage's check of the remote server: a host name or address.
var syslogRegex = regexp.MustCompile(`^[a-zA-Z0-9]([-a-zA-Z0-9]{0,61}[a-zA-Z0-9])?(\.[a-zA-Z0-9]([-a-zA-Z0-9]{0,61}[a-zA-Z0-9])?)*$`)

func errInvalidRemote() error { return errors.New(messages.Get("GENERAL_EDIT_INVALID_REMOTE")) }

// setLogDestination is Host.SetSysLogDestination (the logging map's syslog_destination; none for local only), then
// Host.syslog_reconfigure; nothing when it is unchanged (LogDestinationEditPage.HasChanged).
func setLogDestination(s *xenapi.Session, uuid, to string) error {
	ref, err := xenapi.Host.GetByUUID(s, uuid)
	if err != nil {
		return xapi.Sanitize(err)
	}
	logging, err := xenapi.Host.GetLogging(s, ref)
	if err != nil {
		return xapi.Sanitize(err)
	}
	if logging["syslog_destination"] == to { // already set (maybe by an earlier OK whose reconfigure failed): apply it
		return xapi.Sanitize(xenapi.Host.SyslogReconfigure(s, ref))
	}
	logging = maps.Clone(logging)
	if logging == nil {
		logging = map[string]string{}
	}
	if delete(logging, "syslog_destination"); to != "" {
		logging["syslog_destination"] = to
	}
	if err := xenapi.Host.SetLogging(s, ref, logging); err != nil {
		return xapi.Sanitize(err)
	}
	return xapi.Sanitize(xenapi.Host.SyslogReconfigure(s, ref))
}

// multipathAllowed is HostMultipathPage.MaintenanceMode and EditMultipathAction's RBAC check: the setting changes only
// in maintenance mode or while the server is down.
func multipathAllowed(s *xenapi.Session, uuid string) error {
	ref, err := xenapi.Host.GetByUUID(s, uuid)
	if err != nil {
		return xapi.Sanitize(err)
	}
	h, err := xenapi.Host.GetRecord(s, ref)
	if err != nil {
		return xapi.Sanitize(err)
	}
	live, err := xenapi.HostMetrics.GetLive(s, h.Metrics)
	if err != nil {
		return xapi.Sanitize(err)
	}
	if h.OtherConfig["MAINTENANCE_MODE"] != "true" && live {
		return errors.New("the server must be in maintenance mode to edit these settings") // HostMultipathPage.maintenanceWarningLabel
	}
	return permitted(s, "PBD.plug", "PBD.unplug", "host.set_multipathing", "host.remove_from_other_config", "host.add_to_other_config")
}

// setMultipath is EditMultipathAction: the plugged PBDs are unplugged, multipathing and its handle set, and the PBDs
// plugged again (a plug error is returned only when the change itself succeeded); nothing when it is unchanged.
func setMultipath(s *xenapi.Session, uuid string, on bool) (err error) {
	ref, err := xenapi.Host.GetByUUID(s, uuid)
	if err != nil {
		return xapi.Sanitize(err)
	}
	h, err := xenapi.Host.GetRecord(s, ref)
	if err != nil {
		return xapi.Sanitize(err)
	}
	if h.Multipathing == on {
		return nil
	}
	var plugged []xenapi.PBDRef
	defer func() {
		for _, p := range plugged {
			if perr := xenapi.PBD.Plug(s, p); perr != nil && err == nil {
				err = xapi.Sanitize(perr)
			}
		}
	}()
	for _, p := range h.PBDs {
		attached, err := xenapi.PBD.GetCurrentlyAttached(s, p)
		if err != nil {
			return xapi.Sanitize(err)
		}
		if !attached {
			continue
		}
		plugged = append(plugged, p)
		if err := xenapi.PBD.Unplug(s, p); err != nil {
			return xapi.Sanitize(err)
		}
	}
	if err := xenapi.Host.SetMultipathing(s, ref, on); err != nil { // Kolkata or greater: every XCP-ng
		return xapi.Sanitize(err)
	}
	if err := xenapi.Host.RemoveFromOtherConfig(s, ref, "multipathhandle"); err != nil {
		return xapi.Sanitize(err)
	}
	if on {
		return xapi.Sanitize(xenapi.Host.AddToOtherConfig(s, ref, "multipathhandle", "dmp")) // DEFAULT_MULTIPATH_HANDLE
	}
	return nil
}
