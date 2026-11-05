package ops

import (
	"errors"
	"strings"

	xenapi "go/xenapi"

	"xcpng-center/internal/messages"
	"xcpng-center/internal/rules"
	"xcpng-center/internal/xapi"
)

// ISOAlreadyAttached is SrRules.IsoAlreadyAttached: the pool already has an ISO library on this share (the same type
// and location in a PBD's device_config).
func ISOAlreadyAttached(c *xapi.Connection, kind rules.SRKind, path string) error {
	if !rules.SRIsIso(kind) {
		return nil
	}
	s := c.Session
	want := rules.SRDeviceConfig(kind, path, "", "", "", "")
	pbds, err := xenapi.PBD.GetAllRecords(s)
	if err != nil {
		return xapi.Sanitize(err)
	}
	for _, pbd := range pbds {
		if pbd.DeviceConfig["type"] != want["type"] || pbd.DeviceConfig["location"] != want["location"] {
			continue
		}
		if t, err := xenapi.SR.GetType(s, pbd.SR); err == nil && t == "iso" {
			key := "SMB_ISO_STORAGE_ALREADY_ATTACHED"
			if kind == rules.SRNfsIso {
				key = "NFS_ISO_ALREADY_ATTACHED"
			}
			return errors.New(messages.Format(messages.Get(key), PoolName(c)))
		}
	}
	return nil
}

// PoolName is Helpers.GetName of a connection: the pool's name, else its coordinator's, else the address.
func PoolName(c *xapi.Connection) string {
	s := c.Session
	if pools, err := xenapi.Pool.GetAllRecords(s); err == nil {
		for _, p := range pools {
			if p.NameLabel != "" {
				return p.NameLabel
			}
			if name, err := xenapi.Host.GetNameLabel(s, p.Master); err == nil {
				return name
			}
		}
	}
	return c.Host
}

// NewSR is what the New Storage wizard asks: the kind, the name and description (or the automatic one), the share and
// its options (NFS version "", "4" or "4.1"; SMB credentials, optional), and an existing SR to reattach.
type NewSR struct {
	Kind        rules.SRKind `json:"kind"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	AutoDesc    bool         `json:"auto_description"`
	Path        string       `json:"path"`
	Options     string       `json:"options"`
	NFSVersion  string       `json:"nfs_version"`
	User        string       `json:"user"`
	Password    string       `json:"password"`
	Existing    string       `json:"existing"` // the uuid of an SR a scan found: reattach it instead of creating one
	// iSCSI: the IQN and LUN the scans found, and CHAP (ChapUser "" for none)
	IQN          rules.IQN `json:"iqn"`
	LUN          rules.LUN `json:"lun"`
	ChapUser     string    `json:"chap_user"`
	ChapPassword string    `json:"chap_password"`
}

// CreateSR is NewSRWizard's SrCreateAction: a shared SR made on the coordinator from rules.SRDeviceConfig, with its
// password kept in a secret for the time of SR.create; a PBD the create left unplugged is plugged by hand (the SR is
// forgotten when that fails, unless the host is down); ISO libraries are auto-scanned; the pool's first shared SR
// that is not a library becomes its default SR (SetSrAsDefaultAction). It returns the new SR's UUID.
func CreateSR(c *xapi.Connection, spec NewSR) (uuid string, err error) {
	kind, name := spec.Kind, strings.TrimSpace(spec.Name)
	if !contains(rules.SRKinds, kind) || name == "" {
		return "", errors.New("enter a name for the new SR")
	}
	description, dconf, password, err := srSettings(c, spec)
	if err != nil {
		return "", err
	}
	if spec.Existing != "" && !rules.SRIsIso(kind) {
		return attachSR(c, kind, spec.Existing, name, description, dconf, password)
	}
	s := c.Session
	coordinator, err := coordinatorOf(s)
	if err != nil {
		return "", err
	}
	var submitted, created bool
	for _, key := range []string{"cifspassword", "password", "chappassword"} { // CreateSecret: the first one found
		if v, ok := dconf[key]; ok {
			secret, serr := xenapi.Secret.Create(s, xenapi.SecretRecord{Value: v, OtherConfig: map[string]string{}})
			if serr != nil {
				return "", xapi.Sanitize(serr, password)
			}
			defer func() {
				if createSettled(submitted, created, err) {
					_ = xenapi.Secret.Destroy(s, secret) // SR.create copied it; best effort, as the C# app
				}
			}()
			secretUUID, serr := xenapi.Secret.GetUUID(s, secret)
			if serr != nil {
				return "", xapi.Sanitize(serr, password)
			}
			delete(dconf, key)
			dconf[key+"_secret"] = secretUUID
			break
		}
	}
	contentType := ""
	if rules.SRIsIso(kind) {
		contentType = "iso"
	}
	submitted = true
	ref, err := waitRef(s)(xenapi.SR.AsyncCreate(s, coordinator, dconf, 0, name, description, rules.SRType(kind), contentType, true, map[string]string{}))
	if err != nil {
		return "", xapi.Sanitize(err, password)
	}
	created = true
	sr := xenapi.SRRef(ref)
	pbds, err := xenapi.SR.GetPBDs(s, sr)
	if err != nil {
		return "", xapi.Sanitize(err)
	}
	for _, pbd := range pbds {
		attached, err := xenapi.PBD.GetCurrentlyAttached(s, pbd)
		if err != nil {
			return "", xapi.Sanitize(err)
		}
		if attached {
			continue
		}
		if err := xenapi.PBD.Plug(s, pbd); err != nil {
			var apiErr *xenapi.APIError
			if errors.As(err, &apiErr) && (apiErr.Message == "HOST_OFFLINE" || apiErr.Message == "HOST_STILL_BOOTING") {
				continue // the host is down: its PBD plugs when it comes back
			}
			// SrCreateAction means to roll back here (its catch only does so for non-XAPI errors, which looks unintended)
			forgetSR(s, sr)
			return "", xapi.Sanitize(err, password)
		}
	}
	autoScan := "false"
	if rules.SRIsIso(kind) {
		autoScan = "true"
	}
	if err := xenapi.SR.SetOtherConfig(s, sr, map[string]string{"auto-scan": autoScan}); err != nil {
		return "", xapi.Sanitize(err)
	}
	if !rules.SRIsIso(kind) && firstSharedSR(s, sr) {
		if err := setDefaultSR(s, sr); err != nil {
			return "", err
		}
	}
	uuid, err = xenapi.SR.GetUUID(s, sr)
	return uuid, xapi.Sanitize(err)
}

// createSettled: SR.create's outcome is known (not sent, finished, or refused by XAPI), so its secret can go; after a
// lost reply or ErrUncertain it may still be running and need it (the secret is left on the server then).
func createSettled(submitted, created bool, err error) bool {
	var apiErr *xenapi.APIError
	return !submitted || created || errors.As(err, &apiErr) && !errors.Is(err, ErrUncertain)
}

// srSettings checks a NewSR and gives its description (the automatic one when asked), device_config and password:
// a share for NFS/SMB and the ISO libraries (no library twice), the scanned IQN and LUN for iSCSI (never a LUN an SR
// of this pool uses).
func srSettings(c *xapi.Connection, spec NewSR) (description string, dconf map[string]string, password string, err error) {
	description = spec.Description
	if spec.Kind == rules.SRIscsi {
		if spec.IQN.IQN == "" || spec.IQN.Address == "" || spec.LUN.ScsiID == "" && spec.LUN.Serial == "" {
			return "", nil, "", errors.New("scan the target host and select the target IQN and LUN")
		}
		if sr, err := LUNInUse(c, spec.LUN.ScsiID); err != nil {
			return "", nil, "", err
		} else if sr != "" && spec.Existing == "" {
			return "", nil, "", errors.New(messages.Format(messages.Get("NEWSR_LUN_IN_USE_ON_POOL"), sr, PoolName(c)))
		}
		if spec.AutoDesc {
			description = rules.ISCSIDescription(spec.IQN, spec.LUN)
		}
		return description, rules.ISCSIDeviceConfig(spec.IQN, spec.LUN, spec.ChapUser, spec.ChapPassword), spec.ChapPassword, nil
	}
	if !rules.SRValidLocation(spec.Kind, spec.Path) {
		return "", nil, "", errors.New("enter a valid share for the new SR")
	}
	if err := ISOAlreadyAttached(c, spec.Kind, spec.Path); err != nil {
		return "", nil, "", err
	}
	if spec.AutoDesc {
		description = rules.SRAutoDescription(spec.Kind, strings.TrimSpace(spec.Path))
	}
	return description, rules.SRDeviceConfig(spec.Kind, spec.Path, spec.Options, spec.NFSVersion, spec.User, spec.Password), spec.Password, nil
}

// LUNInUse is LVMoISCSI.SrUsingLunExists on one pool: the name (or uuid, when it has none) of an iSCSI or GFS2 SR
// whose PBD uses this SCSI id, "" for none.
func LUNInUse(c *xapi.Connection, scsiID string) (string, error) {
	if scsiID == "" {
		return "", nil
	}
	s := c.Session
	srs, err := xenapi.SR.GetAllRecords(s)
	if err != nil {
		return "", xapi.Sanitize(err)
	}
	for _, sr := range srs {
		if (sr.Type != "lvmoiscsi" && sr.Type != "gfs2") || len(sr.PBDs) == 0 {
			continue
		}
		pbd, err := xenapi.PBD.GetRecord(s, sr.PBDs[0])
		if err != nil {
			return "", xapi.Sanitize(err)
		}
		if pbd.DeviceConfig["SCSIid"] == scsiID {
			if sr.NameLabel == "" { // an unnamed SR is in use all the same
				return sr.UUID, nil
			}
			return sr.NameLabel, nil
		}
	}
	return "", nil
}

// setDefaultSR is SetSrAsDefaultAction: the pool's default, suspend image and crash dump SR.
func setDefaultSR(s *xenapi.Session, sr xenapi.SRRef) error {
	pools, err := xenapi.Pool.GetAll(s)
	if err != nil {
		return xapi.Sanitize(err)
	}
	for _, pool := range pools {
		for _, set := range []func(*xenapi.Session, xenapi.PoolRef, xenapi.SRRef) error{xenapi.Pool.SetDefaultSR, xenapi.Pool.SetSuspendImageSR, xenapi.Pool.SetCrashDumpSR} {
			if err := set(s, pool, sr); err != nil {
				return xapi.Sanitize(err)
			}
		}
	}
	return nil
}

// firstSharedSR is SrCreateAction.isFirstSharedNonISOSR: no other shared SR of the pool is virtual disk storage.
func firstSharedSR(s *xenapi.Session, sr xenapi.SRRef) bool {
	srs, err := xenapi.SR.GetAllRecords(s)
	if err != nil {
		return false
	}
	for ref, r := range srs {
		if ref != sr && r.Type != "iso" && r.Shared {
			return false
		}
	}
	return true
}

// forgetSR is SrCreateAction.ForgetFailedSR: unplug the SR's PBDs and forget it; failures are ignored.
func forgetSR(s *xenapi.Session, sr xenapi.SRRef) {
	if pbds, err := xenapi.SR.GetPBDs(s, sr); err == nil {
		for _, pbd := range pbds {
			if attached, err := xenapi.PBD.GetCurrentlyAttached(s, pbd); err == nil && attached {
				_ = xenapi.PBD.Unplug(s, pbd)
			}
		}
	}
	_ = xenapi.SR.Forget(s, sr)
}
