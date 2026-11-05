//go:build !nogui

package cmd

import (
	"context"

	"xcpng-center/internal/gui"
)

// runGUI opens the desktop window (internal/gui, Wails v3).
func runGUI(context.Context) error {
	s, err := loadSettings()
	if err != nil {
		return err
	}
	return gui.Run(s, policy(s, false))
}
