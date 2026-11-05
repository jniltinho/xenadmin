//go:build !nogui && (!linux || !cgo)

package gui

import (
	"io/fs"

	"github.com/wailsapp/wails/v3/pkg/application"
)

// setWindowIcon: Windows and macOS take the icon of the executable or the bundle (Options.Icon is set as well).
func setWindowIcon(*application.App, fs.FS) {}
