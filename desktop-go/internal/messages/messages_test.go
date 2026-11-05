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
