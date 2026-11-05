//go:build !nogui

// Package gui is the desktop window: Wails v3 with the Vue/Tailwind SPA of frontend/ (built into frontend/dist and
// embedded). It only presents internal/app; no XAPI logic lives here.
package gui

import (
	"embed"
	"io/fs"

	"github.com/wailsapp/wails/v3/pkg/application"

	"xcpng-center/internal/config"
	"xcpng-center/internal/trust"
)

//go:embed all:dist
var assets embed.FS

// Run opens the main window and blocks until it closes.
func Run(s *config.Settings, p *trust.Policy) error {
	dist, err := fs.Sub(assets, "dist")
	if err != nil {
		return err
	}
	ws := NewWorkspace(s, p)
	defer ws.Close()
	icon, _ := fs.ReadFile(dist, "icons/app.png")
	app := application.New(application.Options{
		Name:     "XCP-ng Center",
		Icon:     icon,
		Linux:    application.LinuxOptions{ProgramName: "xcp-ng-center-go"}, // matches the AppImage's .desktop file
		Services: []application.Service{application.NewService(ws)},
		Assets:   application.AssetOptions{Handler: application.AssetFileServerFS(dist)},
		Mac:      application.MacOptions{ApplicationShouldTerminateAfterLastWindowClosed: true},
	})
	setWindowIcon(app, dist)
	app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:  "XCP-ng Center",
		Width:  1280,
		Height: 800,
		URL:    "/",
	})
	return app.Run()
}
