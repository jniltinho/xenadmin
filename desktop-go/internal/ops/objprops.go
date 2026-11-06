package ops

import (
	"errors"
	"fmt"
	"regexp"
	"slices"
	"strings"
	"time"

	xenapi "go/xenapi"

	"xcpng-center/internal/rules"
	"xcpng-center/internal/xapi"
)

// The General and Custom Fields pages of any object's Properties (GeneralEditPage / GeneralEditPageAction,
// CustomFieldsPage / SaveCustomFieldsAction): name, description, folder (other_config "folder", Folders.Move), tags
// (Tags.AddTag/RemoveTag), a server's iSCSI IQN, and custom field values.

// object is one xenapi class seen through the calls these pages need.
type object struct {
	name, description string
	other             map[string]string
	tags              []string
	setName, setDesc  func(string) error
	addOther          func(k, v string) error
	removeOther       func(k string) error
	addTag, removeTag func(string) error
}

// readObject reads an object of a kind ("vms", "templates", "hosts", "pools", "srs", "networks", "vdis") by uuid.
func readObject(s *xenapi.Session, kind, uuid string) (*object, error) {
	switch kind {
	case "vms", "templates":
		ref, err := xenapi.VM.GetByUUID(s, uuid)
		if err != nil {
			return nil, err
		}
		r, err := xenapi.VM.GetRecord(s, ref)
		if err != nil {
			return nil, err
		}
		return &object{r.NameLabel, r.NameDescription, r.OtherConfig, r.Tags,
			func(v string) error { return xenapi.VM.SetNameLabel(s, ref, v) }, func(v string) error { return xenapi.VM.SetNameDescription(s, ref, v) },
			func(k, v string) error { return xenapi.VM.AddToOtherConfig(s, ref, k, v) }, func(k string) error { return xenapi.VM.RemoveFromOtherConfig(s, ref, k) },
			func(t string) error { return xenapi.VM.AddTags(s, ref, t) }, func(t string) error { return xenapi.VM.RemoveTags(s, ref, t) }}, nil
	case "hosts":
		ref, err := xenapi.Host.GetByUUID(s, uuid)
		if err != nil {
			return nil, err
		}
		r, err := xenapi.Host.GetRecord(s, ref)
		if err != nil {
			return nil, err
		}
		return &object{r.NameLabel, r.NameDescription, r.OtherConfig, r.Tags,
			func(v string) error { return xenapi.Host.SetNameLabel(s, ref, v) }, func(v string) error { return xenapi.Host.SetNameDescription(s, ref, v) },
			func(k, v string) error { return xenapi.Host.AddToOtherConfig(s, ref, k, v) }, func(k string) error { return xenapi.Host.RemoveFromOtherConfig(s, ref, k) },
			func(t string) error { return xenapi.Host.AddTags(s, ref, t) }, func(t string) error { return xenapi.Host.RemoveTags(s, ref, t) }}, nil
	case "pools":
		ref, err := xenapi.Pool.GetByUUID(s, uuid)
		if err != nil {
			return nil, err
		}
		r, err := xenapi.Pool.GetRecord(s, ref)
		if err != nil {
			return nil, err
		}
		return &object{r.NameLabel, r.NameDescription, r.OtherConfig, r.Tags,
			func(v string) error { return xenapi.Pool.SetNameLabel(s, ref, v) }, func(v string) error { return xenapi.Pool.SetNameDescription(s, ref, v) },
			func(k, v string) error { return xenapi.Pool.AddToOtherConfig(s, ref, k, v) }, func(k string) error { return xenapi.Pool.RemoveFromOtherConfig(s, ref, k) },
			func(t string) error { return xenapi.Pool.AddTags(s, ref, t) }, func(t string) error { return xenapi.Pool.RemoveTags(s, ref, t) }}, nil
	case "srs":
		ref, err := xenapi.SR.GetByUUID(s, uuid)
		if err != nil {
			return nil, err
		}
		r, err := xenapi.SR.GetRecord(s, ref)
		if err != nil {
			return nil, err
		}
		return &object{r.NameLabel, r.NameDescription, r.OtherConfig, r.Tags,
			func(v string) error { return xenapi.SR.SetNameLabel(s, ref, v) }, func(v string) error { return xenapi.SR.SetNameDescription(s, ref, v) },
			func(k, v string) error { return xenapi.SR.AddToOtherConfig(s, ref, k, v) }, func(k string) error { return xenapi.SR.RemoveFromOtherConfig(s, ref, k) },
			func(t string) error { return xenapi.SR.AddTags(s, ref, t) }, func(t string) error { return xenapi.SR.RemoveTags(s, ref, t) }}, nil
	case "networks":
		ref, err := xenapi.Network.GetByUUID(s, uuid)
		if err != nil {
			return nil, err
		}
		r, err := xenapi.Network.GetRecord(s, ref)
		if err != nil {
			return nil, err
		}
		return &object{r.NameLabel, r.NameDescription, r.OtherConfig, r.Tags,
			func(v string) error { return xenapi.Network.SetNameLabel(s, ref, v) }, func(v string) error { return xenapi.Network.SetNameDescription(s, ref, v) },
			func(k, v string) error { return xenapi.Network.AddToOtherConfig(s, ref, k, v) }, func(k string) error { return xenapi.Network.RemoveFromOtherConfig(s, ref, k) },
			func(t string) error { return xenapi.Network.AddTags(s, ref, t) }, func(t string) error { return xenapi.Network.RemoveTags(s, ref, t) }}, nil
	case "vdis": // a disk's Properties from the SR Storage tab (PropertiesDialog.PagesFor: General and Custom Fields)
		ref, err := xenapi.VDI.GetByUUID(s, uuid)
		if err != nil {
			return nil, err
		}
		r, err := xenapi.VDI.GetRecord(s, ref)
		if err != nil {
			return nil, err
		}
		return &object{r.NameLabel, r.NameDescription, r.OtherConfig, r.Tags,
			func(v string) error { return xenapi.VDI.SetNameLabel(s, ref, v) }, func(v string) error { return xenapi.VDI.SetNameDescription(s, ref, v) },
			func(k, v string) error { return xenapi.VDI.AddToOtherConfig(s, ref, k, v) }, func(k string) error { return xenapi.VDI.RemoveFromOtherConfig(s, ref, k) },
			func(t string) error { return xenapi.VDI.AddTags(s, ref, t) }, func(t string) error { return xenapi.VDI.RemoveTags(s, ref, t) }}, nil
	}
	return nil, fmt.Errorf("no properties for %q", kind)
}

// iqnRegex is Helpers.IqnRegex.
var iqnRegex = regexp.MustCompile(`^iqn\.\d{4}-\d{2}\.([a-zA-Z0-9][-_a-zA-Z0-9]*(\.[a-zA-Z0-9][-_a-zA-Z0-9]*)*)(:.+)?$`)

// keepEmptyFolders is FolderAction's EMPTY_FOLDERS bookkeeping in the pool's other_config: a folder the object leaves
// as its last content stays as an empty folder; the folder it moves into (and its parents) is no longer empty.
func keepEmptyFolders(s *xenapi.Session, from, to string) error {
	ref, p, _, err := poolOf(s)
	if err != nil {
		return err
	}
	var empty []string
	for _, f := range strings.Split(p.OtherConfig["EMPTY_FOLDERS"], ";") {
		if f != "" {
			empty = append(empty, f)
		}
	}
	before := strings.Join(empty, ";")
	if from != "" {
		inside, err := folderUsers(s, from)
		if err != nil {
			return err
		}
		if inside <= 1 && !slices.Contains(empty, from) {
			empty = append(empty, from)
		}
	}
	if to != "" {
		empty = slices.DeleteFunc(empty, func(f string) bool { return f == to || strings.HasPrefix(to, f+"/") })
	}
	after := strings.Join(empty, ";")
	if after == before {
		return nil
	}
	if _, ok := p.OtherConfig["EMPTY_FOLDERS"]; ok {
		if err := xenapi.Pool.RemoveFromOtherConfig(s, ref, "EMPTY_FOLDERS"); err != nil {
			return xapi.Sanitize(err)
		}
	}
	if after == "" {
		return nil
	}
	return xapi.Sanitize(xenapi.Pool.AddToOtherConfig(s, ref, "EMPTY_FOLDERS", after))
}

// folderUsers counts the objects in a folder or its subfolders (Folders.Children).
func folderUsers(s *xenapi.Session, folder string) (int, error) {
	n := 0
	in := func(other map[string]string) {
		if f := other["folder"]; f == folder || strings.HasPrefix(f, folder+"/") {
			n++
		}
	}
	vms, err := xenapi.VM.GetAllRecords(s)
	if err != nil {
		return 0, xapi.Sanitize(err)
	}
	for _, v := range vms {
		in(v.OtherConfig)
	}
	hosts, _ := xenapi.Host.GetAllRecords(s)
	for _, h := range hosts {
		in(h.OtherConfig)
	}
	srs, _ := xenapi.SR.GetAllRecords(s)
	for _, r := range srs {
		in(r.OtherConfig)
	}
	nets, _ := xenapi.Network.GetAllRecords(s)
	for _, x := range nets {
		in(x.OtherConfig)
	}
	pools, _ := xenapi.Pool.GetAllRecords(s)
	for _, p := range pools {
		in(p.OtherConfig)
	}
	vdis, _ := xenapi.VDI.GetAllRecords(s)
	for _, v := range vdis {
		in(v.OtherConfig)
	}
	return n, nil
}

// ObjectChanges are the General and Custom Fields pages' edits; nil fields are unchanged.
type ObjectChanges struct {
	Name        *string           `json:"name,omitempty"`
	Description *string           `json:"description,omitempty"`
	Folder      *string           `json:"folder,omitempty"` // "/a/b", "" for none
	Tags        []string          `json:"tags,omitempty"`   // the whole new list (with TagsSet)
	TagsSet     bool              `json:"tags_set,omitempty"`
	IQN         *string           `json:"iqn,omitempty"`    // a server's iSCSI IQN
	Fields      map[string]string `json:"fields,omitempty"` // custom fields: "" removes; a Date is RFC 3339
	// a server's remote syslog server ("" for local only) and multipathing
	LogDestination *string `json:"log_destination,omitempty"`
	Multipath      *bool   `json:"multipath,omitempty"`
}

// ValidFolder is a folder path as Folders keeps it: "/" then names separated by "/", none empty.
func ValidFolder(path string) bool {
	if path == "" {
		return true
	}
	if !strings.HasPrefix(path, "/") || strings.HasSuffix(path, "/") && path != "/" {
		return false
	}
	for _, p := range strings.Split(path[1:], "/") {
		if strings.TrimSpace(p) == "" {
			return false
		}
	}
	return path != "/"
}

// SaveObject applies the changes to an object (GeneralEditPageAction and SaveCustomFieldsAction).
func SaveObject(c *xapi.Connection, kind, uuid string, ch ObjectChanges) error {
	s := c.Session
	o, err := readObject(s, kind, uuid)
	if err != nil {
		return xapi.Sanitize(err)
	}
	// every value is checked before the first change, so a refused edit changes nothing
	if ch.Name != nil && strings.TrimSpace(*ch.Name) == "" { // GeneralEditPage: a pool loses its name only by Make into standalone server
		return errors.New("enter a name")
	}
	if ch.Folder != nil && !ValidFolder(*ch.Folder) {
		return errors.New("a folder is a path such as /Folder/Subfolder")
	}
	if ch.IQN != nil && kind == "hosts" && !iqnRegex.MatchString(strings.TrimSpace(*ch.IQN)) { // Helpers.ValidateIscsiIQN
		return errors.New("enter a valid iSCSI IQN, such as iqn.2026-01.com.example:server")
	}
	logTo := ""
	if ch.LogDestination != nil && kind == "hosts" {
		if logTo = strings.TrimSpace(*ch.LogDestination); logTo != "" && !syslogRegex.MatchString(logTo) {
			return errInvalidRemote()
		}
		if err := permitted(s, "host.set_logging", "host.syslog_reconfigure"); err != nil {
			return err
		}
	}
	if ch.Multipath != nil && kind == "hosts" {
		if err := multipathAllowed(s, uuid); err != nil {
			return err
		}
	}
	var fieldDefs []rules.CustomField
	if len(ch.Fields) > 0 {
		if _, fieldDefs, err = poolFields(s); err != nil {
			return err
		}
		for name, v := range ch.Fields {
			i := slices.IndexFunc(fieldDefs, func(f rules.CustomField) bool { return f.Name == name })
			if i < 0 {
				return errors.New("the custom field '" + name + "' no longer exists")
			}
			if _, err := time.Parse(time.RFC3339, v); v != "" && fieldDefs[i].Type == "Date" && err != nil {
				return errors.New("enter a valid date for '" + name + "'")
			}
		}
	}
	if ch.Name != nil {
		name := strings.TrimSpace(*ch.Name)
		if err := o.setName(name); err != nil {
			return xapi.Sanitize(err)
		}
	}
	if ch.Description != nil {
		if err := o.setDesc(*ch.Description); err != nil {
			return xapi.Sanitize(err)
		}
	}
	if ch.Folder != nil && *ch.Folder != o.other["folder"] {
		if err := keepEmptyFolders(s, o.other["folder"], *ch.Folder); err != nil {
			return err
		}
		if _, ok := o.other["folder"]; ok { // MoveToFolderAction / DeleteFolderAction
			if err := o.removeOther("folder"); err != nil {
				return xapi.Sanitize(err)
			}
		}
		if *ch.Folder != "" {
			if err := o.addOther("folder", *ch.Folder); err != nil {
				return xapi.Sanitize(err)
			}
		}
	}
	if ch.TagsSet {
		for _, t := range o.tags {
			if !slices.Contains(ch.Tags, t) {
				if err := o.removeTag(t); err != nil {
					return xapi.Sanitize(err)
				}
			}
		}
		for _, t := range ch.Tags {
			if t = strings.TrimSpace(t); t != "" && !slices.Contains(o.tags, t) {
				if err := o.addTag(t); err != nil {
					return xapi.Sanitize(err)
				}
			}
		}
	}
	if ch.IQN != nil && kind == "hosts" {
		ref, _ := xenapi.Host.GetByUUID(s, uuid)
		if err := xenapi.Host.SetIscsiIqn(s, ref, strings.TrimSpace(*ch.IQN)); err != nil { // Host.SetIscsiIqn (Kolkata on)
			return xapi.Sanitize(err)
		}
	}
	if ch.LogDestination != nil && kind == "hosts" {
		if err := setLogDestination(s, uuid, logTo); err != nil {
			return err
		}
	}
	if ch.Multipath != nil && kind == "hosts" {
		if err := setMultipath(s, uuid, *ch.Multipath); err != nil {
			return err
		}
	}
	if len(ch.Fields) > 0 {
		for name, v := range ch.Fields {
			i := slices.IndexFunc(fieldDefs, func(f rules.CustomField) bool { return f.Name == name })
			if v != "" && fieldDefs[i].Type == "Date" {
				t, _ := time.Parse(time.RFC3339, v) // checked above
				v = rules.FieldDateText(t)
			}
			key := rules.CustomFieldsKey + "." + name
			if _, ok := o.other[key]; ok {
				if err := o.removeOther(key); err != nil {
					return xapi.Sanitize(err)
				}
			}
			if v != "" {
				if err := o.addOther(key, v); err != nil {
					return xapi.Sanitize(err)
				}
			}
		}
	}
	return nil
}
