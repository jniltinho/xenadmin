package rules

import (
	"testing"
	"time"
)

func TestCustomFields(t *testing.T) {
	in := []CustomField{{"Owner", "String"}, {"Expires <soon>", "Date"}}
	out := ParseCustomFields(CustomFieldsXML(in))
	if len(out) != 2 || out[0] != in[0] || out[1] != in[1] {
		t.Fatalf("%v", out)
	}
	d := time.Date(2026, 3, 4, 5, 6, 7, 0, time.UTC)
	if FieldDateText(d) != "2026-03-04 05:06:07Z" {
		t.Error(FieldDateText(d))
	}
	if got, ok := ParseFieldDate(FieldDateText(d)); !ok || !got.Equal(d) {
		t.Error(got, ok)
	}
	if ParseCustomFields("not xml") != nil {
		t.Error("bad XML")
	}
}
