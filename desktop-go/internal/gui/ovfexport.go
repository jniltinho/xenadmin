package gui

import (
	"os"
	"path/filepath"

	"xcpng-center/internal/app"
	"xcpng-center/internal/ops"
	"xcpng-center/internal/xapi"
)

// ExportChoices are the Export wizard's VMs for the selected VM or vApp (uuids, either may be "").
func (w *Workspace) ExportChoices(key, vmUUID, vappUUID string) (out app.ExportChoices, err error) {
	err = w.with(key, func(c *xapi.Connection) error { out, err = app.ReadExportChoices(c, vmUUID, vappUUID); return err })
	return out, err
}

// ExportFolder is the folder the Export wizard starts in (ExportAppliancePage: the Downloads folder, else the home).
func (w *Workspace) ExportFolder() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	if st, err := os.Stat(filepath.Join(home, "Downloads")); err == nil && st.IsDir() {
		return filepath.Join(home, "Downloads")
	}
	return home
}

// PickExportFolder asks for the folder to export to (FOLDER_BROWSER_EXPORT_APPLIANCE); "" when cancelled.
func (w *Workspace) PickExportFolder(dir string) (string, error) {
	if w.Folder == nil {
		return "", errNoWindow
	}
	return w.Folder("Specify the folder where the exported file(s) will be saved.", dir)
}

// PickEulaFile asks for a EULA document to include (EXPORT_EULA_PAGE_FILETYPES); "" when cancelled.
func (w *Workspace) PickEulaFile() (string, error) {
	if w.Files == nil {
		return "", errNoWindow
	}
	return w.Files(false, "Add EULA", "", [][2]string{{"Normal text file (*.txt)", "*.txt"}, {"Rich text file (*.rtf)", "*.rtf"}})
}

// ApplianceTargetError checks where the package (or the XVA file) goes: "" when it can be made.
func (w *Workspace) ApplianceTargetError(dir, name string, xva bool) string {
	return ops.ApplianceTargetError(dir, name, xva)
}

// ExportAppliance exports VMs as an OVF/OVA package; the package's file.
func (w *Workspace) ExportAppliance(key string, e ops.ApplianceExport) (out string, err error) {
	err = w.with(key, func(c *xapi.Connection) error { out, err = ops.ExportAppliance(c, w.policy, e); return err })
	return out, err
}
