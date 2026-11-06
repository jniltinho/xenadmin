package app

import (
	"sort"
	"time"

	xenapi "go/xenapi"

	"xcpng-center/internal/rules"
	"xcpng-center/internal/xapi"
)

// FieldValue is a custom field of the Custom Fields page and its value on the VM: the text, or a Date in RFC 3339
// universal time ("" for none).
type FieldValue struct {
	rules.CustomField
	Value string `json:"value"`
}

// VMCustomFields is the Custom Fields page of a VM (CustomFieldsPage): the pool's fields, sorted, with the VM's values.
func VMCustomFields(c *xapi.Connection, vmUUID string) ([]FieldValue, error) {
	s := c.Session
	ref, err := xenapi.VM.GetByUUID(s, vmUUID)
	if err != nil {
		return nil, xapi.Sanitize(err)
	}
	other, err := xenapi.VM.GetOtherConfig(s, ref)
	if err != nil {
		return nil, xapi.Sanitize(err)
	}
	pools, err := xenapi.Pool.GetAllRecords(s)
	if err != nil {
		return nil, xapi.Sanitize(err)
	}
	out := []FieldValue{}
	for _, p := range pools {
		for _, f := range rules.ParseCustomFields(p.GuiConfig[rules.CustomFieldsKey]) {
			v := other[rules.CustomFieldsKey+"."+f.Name]
			if f.Type == "Date" {
				t, ok := rules.ParseFieldDate(v)
				v = ""
				if ok {
					v = t.UTC().Format(time.RFC3339)
				}
			}
			out = append(out, FieldValue{f, v})
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out, nil
}
