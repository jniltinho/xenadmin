package messages

import "testing"

func TestFormatAndGeneratedTexts(t *testing.T) {
	if got := Format("{0}: {1}", "db {1}", "title"); got != "db {1}: title" {
		t.Errorf("Format = %q", got)
	}
	if Get("CONFIRM_SHUTDOWN_VM") != "Are you sure you want to shut down the selected VM?" || Friendly("Message.name-vm_crashed") == "" {
		t.Error("generated texts missing: run go generate ./internal/messages")
	}
}

func TestFailureText(t *testing.T) {
	if text, short := FailureText("VM_REQUIRES_SR", []string{"OpaqueRef:vm", "OpaqueRef:sr"}); text != "This VM needs storage that cannot be seen from that server" || short != "Cannot see required storage" {
		t.Errorf("%q %q", text, short)
	}
	if text, _ := FailureText("NOT_A_CODE", []string{" a ", "", "b"}); text != "NOT_A_CODE - a - b" {
		t.Errorf("%q", text)
	}
	if text, _ := FailureText("SR_BACKEND_FAILURE", []string{"TransportException", `{"error": "Unable to connect"}`}); text != "Unable to connect" {
		t.Errorf("%q", text)
	}
}
