package mcp

import (
	"context"
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"

	sdk "github.com/modelcontextprotocol/go-sdk/mcp"

	"xcpng-center/internal/config"
)

func session(t *testing.T) *sdk.ClientSession { return sessionWith(t, false) }

func sessionWith(t *testing.T, allowWrite bool) *sdk.ClientSession {
	t.Helper()
	settings, err := config.Load(filepath.Join(t.TempDir(), "settings.toml"))
	if err != nil {
		t.Fatal(err)
	}
	settings.Servers = []config.Server{{Host: "xcp1", Port: 443, User: "root", Name: "pool one"}}
	m := sdk.NewServer(&sdk.Implementation{Name: "xcp-ng-center", Version: "test"}, nil)
	New(settings, "", Credentials{User: "root", Password: "secret-pw"}, allowWrite).Register(m)
	st, ct := sdk.NewInMemoryTransports()
	ctx := context.Background()
	if _, err := m.Connect(ctx, st, nil); err != nil {
		t.Fatal(err)
	}
	cs, err := sdk.NewClient(&sdk.Implementation{Name: "test", Version: "1"}, nil).Connect(ctx, ct, nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { cs.Close() })
	return cs
}

func TestToolsAreTheCSharpNamesAndReadOnly(t *testing.T) {
	cs := session(t)
	tools, err := cs.ListTools(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}
	names := map[string]bool{}
	for _, tool := range tools.Tools {
		names[tool.Name] = true
		if tool.Annotations == nil || !tool.Annotations.ReadOnlyHint {
			t.Errorf("%s must be read-only in this version", tool.Name)
		}
	}
	for _, want := range []string{"list_servers", "connect", "disconnect", "list_objects"} {
		if !names[want] {
			t.Errorf("missing tool %s (XcpTools.cs)", want)
		}
	}
}

func TestWriteToolsNeedAllowWriteAndConfirm(t *testing.T) {
	ctx := context.Background()
	tools, err := session(t).ListTools(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}
	for _, tool := range tools.Tools {
		if tool.Name == "vm_power" {
			t.Fatal("vm_power listed without --allow-write")
		}
	}
	cs := sessionWith(t, true)
	all, err := cs.ListTools(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}
	names := map[string]bool{}
	for _, tool := range all.Tools {
		names[tool.Name] = true
	}
	for _, want := range []string{"list_alerts", "vm_power", "host_power", "host_maintenance", "dismiss_alerts"} {
		if !names[want] {
			t.Errorf("missing tool %s (XcpTools.cs)", want)
		}
	}
	for _, call := range []*sdk.CallToolParams{
		{Name: "host_power", Arguments: map[string]any{"host": "h", "operation": "shutdown"}},
		{Name: "host_maintenance", Arguments: map[string]any{"host": "h", "enter": true}},
		{Name: "dismiss_alerts", Arguments: map[string]any{"alerts": []string{"a"}}},
	} {
		res, err := cs.CallTool(ctx, call)
		if text, _ := json.Marshal(res.Content); err != nil || !res.IsError || !strings.Contains(string(text), "confirm") {
			t.Errorf("%s without confirm: %s %v", call.Name, text, err)
		}
	}
	res, _ := cs.CallTool(ctx, &sdk.CallToolParams{Name: "host_power", Arguments: map[string]any{"host": "h", "operation": "halt", "confirm": true}})
	if text, _ := json.Marshal(res.Content); !res.IsError || !strings.Contains(string(text), "reboot or shutdown") {
		t.Errorf("host_power with an unknown operation: %s", text)
	}
	res, err = cs.CallTool(ctx, &sdk.CallToolParams{Name: "vm_power", Arguments: map[string]any{"vm": "u", "operation": "start"}})
	if err != nil {
		t.Fatal(err)
	}
	text, _ := json.Marshal(res.Content)
	if !res.IsError || !strings.Contains(string(text), `confirm`) {
		t.Errorf("vm_power without confirm: %s", text)
	}
	res, _ = cs.CallTool(ctx, &sdk.CallToolParams{Name: "vm_power", Arguments: map[string]any{"vm": "u", "operation": "start", "confirm": true}})
	if text, _ = json.Marshal(res.Content); !res.IsError || !strings.Contains(string(text), "no connected server") {
		t.Errorf("vm_power with no server: %s", text)
	}
}

func TestToolErrorsAreResultsWithoutSecrets(t *testing.T) {
	cs := session(t)
	ctx := context.Background()
	res, err := cs.CallTool(ctx, &sdk.CallToolParams{Name: "list_servers"})
	if err != nil {
		t.Fatal(err)
	}
	out, _ := json.Marshal(res.StructuredContent)
	if !strings.Contains(string(out), `"server":"xcp1:443"`) || !strings.Contains(string(out), `"connected":false`) {
		t.Errorf("list_servers: %s", out)
	}
	for _, call := range []*sdk.CallToolParams{
		{Name: "connect", Arguments: map[string]any{"host": "other"}},
		{Name: "list_objects", Arguments: map[string]any{"type": "vm"}},
		{Name: "list_objects", Arguments: map[string]any{"type": "nope"}},
	} {
		res, err := cs.CallTool(ctx, call)
		if err != nil {
			t.Fatal(err)
		}
		text := res.Content[0].(*sdk.TextContent).Text
		if !res.IsError || strings.Contains(text, "secret-pw") {
			t.Errorf("%s: want a tool error without the password, got %v %q", call.Name, res.IsError, text)
		}
	}
}
