package config

import (
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"
)

// TestOptions: the window sends and reads the options by these JSON names, and they survive a save.
func TestOptions(t *testing.T) {
	var o Options
	if err := json.Unmarshal([]byte(`{"connection_timeout":7,"area_graphs":true,"proxy":"specified"}`), &o); err != nil || o.ConnectionTimeout != 7 || !o.AreaGraphs {
		t.Fatalf("%+v %v", o, err)
	}
	if b, _ := json.Marshal(o); !strings.Contains(string(b), `"connection_timeout":7`) {
		t.Fatal(string(b))
	}
	s, err := Load(filepath.Join(t.TempDir(), "settings.toml"))
	if err != nil {
		t.Fatal(err)
	}
	if err := s.SetOptions(o); err != nil {
		t.Fatal(err)
	}
	again, _ := Load(s.Path())
	if again.GetOptions() != o {
		t.Fatalf("%+v", again.GetOptions())
	}
}
