package ops

import "testing"

func TestTrimError(t *testing.T) {
	resp := `<?xml version="1.0" ?><trim_response><key_value_pair><key>errcode</key><value>NoSuchCode</value></key_value_pair>` +
		`<key_value_pair><key>errmsg</key><value>Trim on [x] not supported</value></key_value_pair></trim_response>`
	if got := TrimError(resp); got != "Trim on [x] not supported" {
		t.Errorf("errmsg: %q", got)
	}
	if got := TrimError("false"); got != "unknown error" {
		t.Errorf("not XML: %q", got)
	}
}
