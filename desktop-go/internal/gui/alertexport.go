package gui

import (
	"os"

	"xcpng-center/internal/rules"
)

// ExportAlerts asks where to save the alerts and writes them as CSV (AlertSummaryPage's Export); the file, "" when
// cancelled.
func (w *Workspace) ExportAlerts(rows []rules.AlertCSV) (string, error) {
	if w.Files == nil {
		return "", errNoWindow
	}
	file, err := w.Files(true, "Export all", "alerts.csv", [][2]string{{"Comma separated values (*.csv)", "*.csv"}})
	if err != nil || file == "" {
		return "", err
	}
	return file, os.WriteFile(file, []byte(rules.AlertsCSV(rows)), 0o600)
}
