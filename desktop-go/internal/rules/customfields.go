package rules

import (
	"encoding/xml"
	"strings"
	"time"
)

// CustomField is a CustomFieldDefinition: its name and type ("String" or "Date").
type CustomField struct {
	Name string `json:"name" xml:"name,attr"`
	Type string `json:"type" xml:"type,attr"`
}

// CustomFieldsKey is CustomFieldsManager.CUSTOM_FIELD_BASE_KEY: the pool's gui_config key of the definitions, and the
// prefix of each value's other_config key (GetCustomFieldKey: "XenCenter.CustomFields.<name>").
const CustomFieldsKey = "XenCenter.CustomFields"

// ParseCustomFields reads the definitions as CustomFieldsCache: <CustomFieldDefinitions><CustomFieldDefinition
// name=".." type="String"/>...
func ParseCustomFields(text string) []CustomField {
	var defs struct {
		Fields []CustomField `xml:"CustomFieldDefinition"`
	}
	if v := strings.TrimSpace(text); v == "" || xml.Unmarshal([]byte(v), &defs) != nil {
		return nil
	}
	return defs.Fields
}

// CustomFieldsXML is CustomFieldsManager.GetCustomFieldDefinitionXML (with the unused defaultValue of CA-37473).
func CustomFieldsXML(fields []CustomField) string {
	type def struct {
		XMLName xml.Name `xml:"CustomFieldDefinition"`
		Name    string   `xml:"name,attr"`
		Type    string   `xml:"type,attr"`
		Default string   `xml:"defaultValue,attr"`
	}
	doc := struct {
		XMLName xml.Name `xml:"CustomFieldDefinitions"`
		Fields  []def
	}{}
	for _, f := range fields {
		doc.Fields = append(doc.Fields, def{Name: f.Name, Type: f.Type})
	}
	out, _ := xml.Marshal(doc)
	return string(out)
}

// ParseFieldDate is CustomFieldsManager.ParseValue for a Date: the "u" format the C# app writes, or the usual ISO forms
// (a time without a zone is taken as UTC, as DateTime.ToLocalTime takes an unspecified kind).
// shortcut: not every form DateTime.TryParse takes (culture-dependent dates such as 15/01/2026).
func ParseFieldDate(v string) (time.Time, bool) {
	for _, layout := range []string{"2006-01-02 15:04:05Z", time.RFC3339, "2006-01-02T15:04:05", "2006-01-02 15:04:05", "2006-01-02"} {
		if t, err := time.Parse(layout, v); err == nil {
			return t, true
		}
	}
	return time.Time{}, false
}

// FieldDateText is CustomField.ValueAsInvariantString for a Date: universal time, "u" format.
func FieldDateText(t time.Time) string { return t.UTC().Format("2006-01-02 15:04:05Z") }
