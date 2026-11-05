// Package archtest checks the layering (task 1.7): the core packages never import the UI, the CLI or MCP, so the
// desktop app, the CLI, the MCP server and a future web server share one core.
package archtest

import (
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

var core = []string{"xapi", "trust", "config", "app", "buildinfo", "model", "rules", "ops", "console"}

var forbidden = []string{"github.com/spf13/cobra", "github.com/modelcontextprotocol/", "github.com/wailsapp/", "xcpng-center/cmd", "xcpng-center/internal/mcp", "xcpng-center/internal/gui", "xcpng-center/internal/server"}

func TestCoreImportsNoFrontEnd(t *testing.T) {
	for _, pkg := range core {
		dir := filepath.Join("..", pkg)
		if _, err := os.Stat(dir); err != nil {
			continue // not written yet
		}
		files, _ := filepath.Glob(filepath.Join(dir, "*.go"))
		for _, f := range files {
			parsed, err := parser.ParseFile(token.NewFileSet(), f, nil, parser.ImportsOnly)
			if err != nil {
				t.Fatal(err)
			}
			for _, imp := range parsed.Imports {
				path, _ := strconv.Unquote(imp.Path.Value)
				for _, bad := range forbidden {
					if strings.HasPrefix(path, bad) {
						t.Errorf("%s imports %s: the core must not depend on a front end", f, path)
					}
				}
			}
		}
	}
}
