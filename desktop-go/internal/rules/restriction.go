package rules

import (
	"encoding/xml"
	"strconv"
)

// Restriction reads a restriction ("vcpus-max", "memory-static-max", ...) of the VM's recommendations (VM.GetMaxRestrictionValue).
// shortcut: the value of a matching template with the same reference_label (CP-44766) is not looked up.
func Restriction(recommendations, field, attr string, def int64) int64 {
	var doc struct {
		R []struct {
			Field string `xml:"field,attr"`
			Max   string `xml:"max,attr"`
			Min   string `xml:"min,attr"`
			Value string `xml:"value,attr"`
		} `xml:"restriction"`
	}
	if xml.Unmarshal([]byte(recommendations), &doc) != nil {
		return def
	}
	for _, r := range doc.R {
		if r.Field != field {
			continue
		}
		v := map[string]string{"max": r.Max, "min": r.Min, "value": r.Value}[attr]
		if n, err := strconv.ParseInt(v, 10, 64); err == nil {
			return n
		}
	}
	return def
}
