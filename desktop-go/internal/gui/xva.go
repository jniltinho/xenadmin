package gui

import (
	"errors"

	xenapi "go/xenapi"

	"xcpng-center/internal/app"
	"xcpng-center/internal/ops"
	"xcpng-center/internal/rules"
	"xcpng-center/internal/xapi"
)

var errNoWindow = errors.New("no window to ask for a file")

// FileDialog asks for a file to save (save) or open, with a title, a suggested name and filters ("XVA (*.xva)" ->
// "*.xva"); "" when cancelled. The window sets it (gui.go); without a window it fails.
type FileDialog func(save bool, title, name string, filters [][2]string) (string, error)

// PickImportFile asks for the XVA file or OVF/OVA appliance to import.
func (w *Workspace) PickImportFile() (string, error) {
	if w.Files == nil {
		return "", errNoWindow
	}
	return w.Files(false, "Import VM", "", [][2]string{{"XVA, OVF and OVA (*.xva, *.xva.gz, *.ovf, *.ova, *.ova.gz, *.ova.bz2)",
		"*.xva;*.xva.gz;*.ovf;*.ova;*.ova.gz;*.ova.bz2"}, {"XVA (*.xva, *.xva.gz)", "*.xva;*.xva.gz"},
		{"OVF and OVA (*.ovf, *.ova, *.ova.gz, *.ova.bz2)", "*.ovf;*.ova;*.ova.gz;*.ova.bz2"}})
}

// IsAppliance tells whether the file is an OVF/OVA appliance (by its name), which has its own import dialog.
func (w *Workspace) IsAppliance(file string) bool { return rules.IsAppliance(file) }

// ApplianceInfo reads the appliance: its VMs, disks, interfaces, EULAs and manifest.
func (w *Workspace) ApplianceInfo(file string) (*rules.Appliance, error) {
	p, err := rules.OpenPackage(file)
	if err != nil {
		return nil, err
	}
	return p.Appliance, nil
}

// ApplianceMappingsError is why the mappings cannot be used for the appliance ApplianceInfo read ("" when they can):
// the dialog's check before it closes; ImportAppliance checks them again against the package itself.
func (w *Workspace) ApplianceMappingsError(key string, a rules.Appliance, maps []rules.ApplianceMapping) (out string, err error) {
	err = w.with(key, func(c *xapi.Connection) error { out, err = ops.ApplianceMappingsError(c, &a, maps); return err })
	return out, err
}

// ApplianceUnlockError is the dialog's check of an encrypted appliance's password (ImportSecurityPage): "" when it
// decrypts the package. The password is not kept.
func (w *Workspace) ApplianceUnlockError(file, password string) string {
	p, err := rules.OpenPackage(file)
	if err == nil {
		err = p.Unlock(password)
	}
	if err != nil {
		return err.Error()
	}
	return ""
}

// ImportAppliance imports an OVF/OVA appliance with a mapping per VM (server, SR per disk, network per interface),
// verifying its manifest, decrypting it with its password ("" when it is not encrypted) and starting it when asked;
// the first VM's uuid.
func (w *Workspace) ImportAppliance(key, file string, maps []rules.ApplianceMapping, verify, start bool, password string) (out string, err error) {
	err = w.with(key, func(c *xapi.Connection) error {
		out, err = ops.ImportAppliance(c, w.policy, file, maps, verify, start, password)
		return err
	})
	return out, err
}

// ExportBlocker is why the VM cannot be exported now ("" when it can).
func (w *Workspace) ExportBlocker(key, vmUUID string) (out string, err error) {
	err = w.with(key, func(c *xapi.Connection) error {
		ref, err := xenapi.VM.GetByUUID(c.Session, vmUUID)
		if err != nil {
			return xapi.Sanitize(err)
		}
		vm, err := xenapi.VM.GetRecord(c.Session, ref)
		if err != nil {
			return xapi.Sanitize(err)
		}
		out = ops.ExportBlocker(vm)
		return nil
	})
	return out, err
}

// ExportVM exports a halted or suspended VM to an XVA file, verified when asked.
func (w *Workspace) ExportVM(key, vmUUID, file string, verify bool) error {
	return w.with(key, func(c *xapi.Connection) error { return ops.ExportVM(c, w.policy, vmUUID, file, verify) })
}

// ImportOptions are the Import VM dialog's servers and SRs for a home server ("" for none).
func (w *Workspace) ImportOptions(key, hostUUID string) (out app.ImportOptions, err error) {
	err = w.with(key, func(c *xapi.Connection) error { out, err = app.ReadImportOptions(c, hostUUID); return err })
	return out, err
}

// ImportFileError checks the file to import ("" when it can be imported).
func (w *Workspace) ImportFileError(file string) string { return ops.ImportFileError(file) }

// ImportVM imports an XVA file to an SR (with an optional home server) and starts the VM when asked; the new VM's uuid.
func (w *Workspace) ImportVM(key, file, hostUUID, srUUID string, start bool) (out string, err error) {
	err = w.with(key, func(c *xapi.Connection) error {
		out, err = ops.ImportVM(c, w.policy, file, hostUUID, srUUID, start)
		return err
	})
	return out, err
}
