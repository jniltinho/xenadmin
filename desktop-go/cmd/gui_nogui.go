//go:build nogui

package cmd

import (
	"context"
	"errors"
)

// runGUI in a build without the desktop app (tag nogui: tests and headless servers).
func runGUI(context.Context) error {
	return errors.New("this build has no desktop app (built with -tags nogui); use a subcommand, see --help")
}
