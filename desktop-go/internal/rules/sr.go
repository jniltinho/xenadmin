package rules

import (
	"regexp"
	"strings"

	"xcpng-center/internal/messages"
)

// SRKind is a type of the New Storage wizard (SrRules.Kind): virtual disk storage on NFS or SMB, or an ISO library.
type SRKind string

const (
	SRNfs    SRKind = "nfs"
	SRSmb    SRKind = "smb"
	SRNfsIso SRKind = "nfs_iso"
	SRSmbIso SRKind = "smb_iso"
)

// SRKinds in the wizard's order, with their names and blurbs (SrRules.TypeName / Blurb).
var SRKinds = []SRKind{SRNfs, SRSmb, SRNfsIso, SRSmbIso}

func srKey(k SRKind, nfs, smb, nfsIso, smbIso string) string {
	return map[SRKind]string{SRNfs: nfs, SRSmb: smb, SRNfsIso: nfsIso, SRSmbIso: smbIso}[k]
}

// SRTypeName is the wizard's label of a kind.
func SRTypeName(k SRKind) string {
	return messages.Get(srKey(k, "NEWSR_VHDONFS_TYPE_NAME", "NEWSR_CIFS_TYPE_NAME", "NEWSR_NFS_ISO_TYPE_NAME", "NEWSR_CIFS_ISO_TYPE_NAME"))
}

// SRBlurb is the wizard's description of a kind.
func SRBlurb(k SRKind) string {
	return messages.Get(srKey(k, "NEWSR_VHDONFS_BLURB", "NEWSR_CIFS_BLURB", "NEWSR_NFS_ISO_BLURB", "NEWSR_CIFS_ISO_BLURB"))
}

// SRIsIso: an ISO library ("iso" content, type iso), not virtual disk storage.
func SRIsIso(k SRKind) bool { return k == SRNfsIso || k == SRSmbIso }

// SRType is the XAPI SR type of a kind (SrRules.TypeOf).
func SRType(k SRKind) string { return srKey(k, "nfs", "smb", "iso", "iso") }

// SRDefaultName is SrRules.DefaultName: "NFS virtual disk storage", then "(1)", "(2)"... when taken.
func SRDefaultName(k SRKind, taken map[string]bool) string {
	return UniqueName(messages.Get(srKey(k, "SRWIZARD_NFS_STORAGE", "SRWIZARD_CIFS_STORAGE", "SRWIZARD_NFS_LIBRARY", "SRWIZARD_CIFS_LIBRARY")), taken)
}

var (
	nfsPath = regexp.MustCompile(`^[^/:*?<>|\\]+:/[^:*?<>|\\]+$`)
	smbPath = regexp.MustCompile(`^\\(\\[^/:*?<>|\\]+\\{0,1})+$`)
)

// SRValidLocation is SrRules.ValidLocation: server:/path for NFS, \\server\share for SMB.
func SRValidLocation(k SRKind, path string) bool {
	if k == SRNfs || k == SRNfsIso {
		return nfsPath.MatchString(strings.TrimSpace(path))
	}
	return smbPath.MatchString(path)
}

// SRAutoDescription is SrRules.AutoDescription ("NFS SR [server:/path]").
func SRAutoDescription(k SRKind, path string) string {
	return messages.Format(messages.Get(srKey(k, "NEWSR_ACTION", "NEWSR_CIFS_ACTION", "NEWSR_ISO_DESCRIPTION", "NEWSR_CIF_DESCRIPTION")), path)
}

// SRDeviceConfig is SrRules.DeviceConfig / IsoDeviceConfig: the device_config of the new SR. nfsVersion is "" (3),
// "4" or "4.1"; credentials only when given.
func SRDeviceConfig(k SRKind, path, options, nfsVersion, user, password string) map[string]string {
	switch k {
	case SRNfsIso:
		d := map[string]string{"location": strings.TrimSpace(path), "type": "nfs_iso"}
		if nfsVersion != "" {
			d["nfsversion"] = "4"
		}
		return d
	case SRSmbIso:
		d := map[string]string{"location": strings.ReplaceAll(path, `\`, "/"), "type": "cifs"}
		if bits := strings.Split(d["location"], "/"); len(bits) > 4 { // a sub-folder of the share goes to iso_path
			d["location"] = "//" + bits[2] + "/" + bits[3]
			d["iso_path"] = "/" + strings.Join(bits[4:], "/")
		}
		if user != "" {
			d["username"], d["cifspassword"] = user, password
		}
		return d
	}
	if k == SRNfs {
		path = strings.TrimSpace(path)
	}
	parts := strings.SplitN(path, ":", 2)
	d := map[string]string{"server": parts[0]}
	if len(parts) > 1 {
		d["serverpath"] = parts[1]
	}
	if k == SRNfs {
		d["options"] = options
		if nfsVersion != "" {
			d["nfsversion"] = nfsVersion
		}
	} else if strings.TrimSpace(user) != "" || strings.TrimSpace(password) != "" {
		d["username"], d["password"] = user, password
	}
	return d
}
