package ops

import (
	"errors"
	"slices"
	"strings"

	xenapi "go/xenapi"

	"xcpng-center/internal/rules"
	"xcpng-center/internal/xapi"
)

// poolFields reads the pool and its custom field definitions.
func poolFields(s *xenapi.Session) (xenapi.PoolRef, []rules.CustomField, error) {
	pools, err := xenapi.Pool.GetAllRecords(s)
	if err != nil {
		return "", nil, xapi.Sanitize(err)
	}
	for ref, p := range pools {
		return ref, rules.ParseCustomFields(p.GuiConfig[rules.CustomFieldsKey]), nil
	}
	return "", nil, errors.New("the server has no pool")
}

// saveFields is CustomFieldsManager.SaveCustomFields: the pool's gui_config entry replaced.
func saveFields(s *xenapi.Session, pool xenapi.PoolRef, fields []rules.CustomField) error {
	_ = xenapi.Pool.RemoveFromGuiConfig(s, pool, rules.CustomFieldsKey)
	return xapi.Sanitize(xenapi.Pool.AddToGuiConfig(s, pool, rules.CustomFieldsKey, rules.CustomFieldsXML(fields)))
}

// AddCustomField is AddCustomFieldAction (CustomFieldsDialog's Add): a new field for every object of the pool.
func AddCustomField(c *xapi.Connection, name, kind string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return errors.New("enter a name for the field")
	}
	if kind != "String" && kind != "Date" {
		return errors.New("unknown field type")
	}
	s := c.Session
	pool, fields, err := poolFields(s)
	if err != nil {
		return err
	}
	if slices.ContainsFunc(fields, func(f rules.CustomField) bool { return f.Name == name }) {
		return errors.New("a custom field with this name already exists")
	}
	return saveFields(s, pool, append(fields, rules.CustomField{Name: name, Type: kind}))
}

// RemoveCustomField is RemoveCustomFieldAction: the values on the pools, servers, VMs, SRs, disks and networks, then
// the definition (kept until every value is gone, so a failed removal can be retried).
func RemoveCustomField(c *xapi.Connection, name string) error {
	s := c.Session
	pool, fields, err := poolFields(s)
	if err != nil {
		return err
	}
	i := slices.IndexFunc(fields, func(f rules.CustomField) bool { return f.Name == name })
	if i < 0 {
		return nil
	}
	key := rules.CustomFieldsKey + "." + name
	if err := removeFieldValues(s, key); err != nil {
		return err
	}
	return saveFields(s, pool, slices.Delete(fields, i, i+1))
}

// removeFieldValues removes the other_config key from every object type that may hold a custom field.
func removeFieldValues(s *xenapi.Session, key string) error {
	pools, err := xenapi.Pool.GetAllRecords(s)
	if err != nil {
		return xapi.Sanitize(err)
	}
	hosts, err := xenapi.Host.GetAllRecords(s)
	if err != nil {
		return xapi.Sanitize(err)
	}
	vms, err := xenapi.VM.GetAllRecords(s)
	if err != nil {
		return xapi.Sanitize(err)
	}
	srs, err := xenapi.SR.GetAllRecords(s)
	if err != nil {
		return xapi.Sanitize(err)
	}
	vdis, err := xenapi.VDI.GetAllRecords(s)
	if err != nil {
		return xapi.Sanitize(err)
	}
	nets, err := xenapi.Network.GetAllRecords(s)
	if err != nil {
		return xapi.Sanitize(err)
	}
	var errs []error
	remove := func(other map[string]string, f func() error) {
		if _, ok := other[key]; ok {
			errs = append(errs, xapi.Sanitize(f()))
		}
	}
	for r, o := range pools {
		remove(o.OtherConfig, func() error { return xenapi.Pool.RemoveFromOtherConfig(s, r, key) })
	}
	for r, o := range hosts {
		remove(o.OtherConfig, func() error { return xenapi.Host.RemoveFromOtherConfig(s, r, key) })
	}
	for r, o := range vms {
		remove(o.OtherConfig, func() error { return xenapi.VM.RemoveFromOtherConfig(s, r, key) })
	}
	for r, o := range srs {
		remove(o.OtherConfig, func() error { return xenapi.SR.RemoveFromOtherConfig(s, r, key) })
	}
	for r, o := range vdis {
		remove(o.OtherConfig, func() error { return xenapi.VDI.RemoveFromOtherConfig(s, r, key) })
	}
	for r, o := range nets {
		remove(o.OtherConfig, func() error { return xenapi.Network.RemoveFromOtherConfig(s, r, key) })
	}
	return errors.Join(errs...)
}

// SaveVMCustomFields is SaveCustomFieldsAction for a VM: each field's value set, or removed when empty (SaveObject).
func SaveVMCustomFields(c *xapi.Connection, vmUUID string, values map[string]string) error {
	return SaveObject(c, "vms", vmUUID, ObjectChanges{Fields: values})
}
