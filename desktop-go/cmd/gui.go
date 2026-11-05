//go:build !nogui

package cmd

import (
	"context"
	"errors"
)

// runGUI opens the desktop app. shortcut: the Wails v3 window (internal/gui, task 1.3) is not wired yet; until then
// the binary is the CLI and MCP server only.
func runGUI(context.Context) error {
	return errors.New("the desktop window is not available in this build yet; use a subcommand, see --help")
}
