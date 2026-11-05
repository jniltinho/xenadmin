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

func session(t *testing.T) *sdk.ClientSession {
	t.Helper()
	settings, err := config.Load(filepath.Join(t.TempDir(), "settings.toml"))
	if err != nil {
		t.Fatal(err)
	}
	settings.Servers = []config.Server{{Host: "xcp1", Port: 443, User: "root", Name: "pool one"}}
	m := sdk.NewServer(&sdk.Implementation{Name: "xcp-ng-center", Version: "test"}, nil)
	New(settings, "", Credentials{User: "root", Password: "secret-pw"}).Register(m)
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
