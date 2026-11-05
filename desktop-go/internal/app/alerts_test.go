package app

import (
	"strings"
	"testing"

	xenapi "go/xenapi"
)

func TestAlertTextsMirrorNotificationRules(t *testing.T) {
	p := pool{uuid: "P", name: "lab", isPool: false, names: map[string]string{"P": "lab", "H": "host1", "V": "web"}, hosts: map[string]bool{"H": true}}
	login := xenapi.MessageRecord{Name: "FAILED_LOGIN_ATTEMPTS", Cls: xenapi.ClsPool, ObjUUID: "P",
		Body: "<body><total>3</total><unknown>1</unknown><known><ip>10.0.0.9</ip><number>2</number><username>root</username></known></body>"}
	if got := alertTitle(login, p); got != "There has been a number of failed attempts to login to 'lab' within the past hour." {
		t.Errorf("title = %q", got)
	}
	want := "There has been a total of 3 failed login attempts.\n1 attempts originated from unknown sources.\n" +
		"The most common attempts originated from the following sources:\n- IP address: 10.0.0.9, Number of attempts: 2, Username: root"
	if got := alertDescription(login, p); got != want {
		t.Errorf("description = %q", got)
	}
	// a VM's alert takes its name; a standalone server's host does not
	vm := xenapi.MessageRecord{Name: "VM_CRASHED", Cls: xenapi.ClsVM, ObjUUID: "V"}
	if got := alertTitle(vm, p); !strings.HasPrefix(got, "web: ") {
		t.Errorf("VM title = %q", got)
	}
	host := xenapi.MessageRecord{Name: "HOST_CLOCK_SKEW_DETECTED", Cls: xenapi.ClsHost, ObjUUID: "H", Body: "raw"}
	if got := alertTitle(host, p); strings.HasPrefix(got, "host1: ") {
		t.Errorf("standalone host title = %q", got)
	}
	if got := alertDescription(host, p); !strings.Contains(got, "host1") || !strings.Contains(got, "lab") {
		t.Errorf("clock skew description = %q", got)
	}
	gfs2 := xenapi.MessageRecord{Name: "GFS2_CAPACITY", Cls: xenapi.ClsSR, ObjUUID: "V"}
	if got := alertTitle(gfs2, p); strings.HasPrefix(got, "web: ") {
		t.Errorf("GFS2 title = %q", got)
	}
	if _, ok := iso8601("2026-11-05T12:00:00Z"); !ok {
		t.Error("ISO 8601 with dashes")
	}
	unknown := xenapi.MessageRecord{Name: "SOMETHING_NEW", Body: "body text"}
	if alertTitle(unknown, p) != "SOMETHING_NEW" || alertDescription(unknown, p) != "body text" {
		t.Error("unknown messages show their name and body")
	}
}
