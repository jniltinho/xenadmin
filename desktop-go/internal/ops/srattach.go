package ops

import (
	"errors"
	"fmt"

	xenapi "go/xenapi"

	"xcpng-center/internal/messages"
	"xcpng-center/internal/rules"
	"xcpng-center/internal/xapi"
)

// FoundSR is an SR a scan found on the share; InUse names this pool when the SR is attached here already
// (SrWizardHelpers.SrInUse; the wizard does not offer it then).
type FoundSR struct {
	rules.SRInfo
	InUse string `json:"in_use,omitempty"`
}

// ProbeSR is the Scan button of the NFS and SMB pages (VHDoNFS / CIFS_SR.buttonScan): SR.probe on the coordinator for
// the SRs already on the share. NFS also asks for the server's versions (probeversion), as the C# app.
func ProbeSR(c *xapi.Connection, spec NewSR) ([]FoundSR, error) {
	if (spec.Kind != rules.SRNfs && spec.Kind != rules.SRSmb) || !rules.SRValidLocation(spec.Kind, spec.Path) {
		return nil, errors.New("enter a valid share to scan")
	}
	s := c.Session
	coordinator, err := coordinatorOf(s)
	if err != nil {
		return nil, err
	}
	dconf := rules.SRDeviceConfig(spec.Kind, spec.Path, spec.Options, "", spec.User, spec.Password)
	if spec.Kind == rules.SRNfs {
		dconf["probeversion"] = ""
	}
	result, err := waitResult(s)(xenapi.SR.AsyncProbe(s, coordinator, dconf, rules.SRType(spec.Kind), map[string]string{}))
	if err != nil {
		return nil, xapi.Sanitize(err, spec.Password)
	}
	srs, err := rules.ParseSRList(result)
	if err != nil {
		return nil, err
	}
	out := make([]FoundSR, 0, len(srs))
	for _, sr := range srs {
		f := FoundSR{SRInfo: sr}
		if used, err := SRInUse(c, sr.UUID); err != nil {
			return nil, err
		} else if used {
			f.InUse = PoolName(c)
		}
		out = append(out, f)
	}
	return out, nil
}

// MultiPoolWarning is NewSRWizard's question before attaching an existing NFS SR (NEWSR_MULTI_POOL_WARNING); SMB has
// none (SrRules.AttachWarning).
func MultiPoolWarning(kind rules.SRKind, uuid string) string {
	if kind != rules.SRNfs || uuid == "" {
		return ""
	}
	return messages.Format(messages.Get("NEWSR_MULTI_POOL_WARNING"), "XCP-ng Center", uuid)
}

// SRInUse is SrWizardHelpers.SrInUse on one pool: it knows an SR with this uuid and the SR has PBDs (SR.HasPBDs).
// The window asks every connected pool, so that two pools never share an SR.
func SRInUse(c *xapi.Connection, uuid string) (bool, error) {
	ref, err := xenapi.SR.GetByUUID(c.Session, uuid)
	var apiErr *xenapi.APIError
	if errors.As(err, &apiErr) && apiErr.Message == "UUID_INVALID" {
		return false, nil
	}
	if err != nil {
		return false, xapi.Sanitize(err)
	}
	pbds, err := xenapi.SR.GetPBDs(c.Session, ref)
	if err != nil {
		return false, xapi.Sanitize(err)
	}
	return len(pbds) > 0, nil
}

func coordinatorOf(s *xenapi.Session) (xenapi.HostRef, error) {
	pools, err := xenapi.Pool.GetAllRecords(s)
	if err != nil {
		return "", xapi.Sanitize(err)
	}
	for _, p := range pools {
		return p.Master, nil
	}
	return "", errors.New("the pool has no coordinator")
}

// attachSR is SrRules.Finish for a scanned SR: SrReattachAction when the pool knows it detached (PBDs for every host,
// then its name and description), else SrIntroduceAction (forget a stale record, SR.introduce, the PBDs, the default
// SR when it is the first shared one). A failure half way unplugs and destroys the PBDs made (RollbackFailedReattach).
func attachSR(c *xapi.Connection, kind rules.SRKind, uuid, name, description string, dconf map[string]string, password string) (_ string, err error) {
	s := c.Session
	if used, err := SRInUse(c, uuid); err != nil {
		return "", err
	} else if used {
		return "", fmt.Errorf("this SR is already in use by %s", PoolName(c))
	}
	hosts, err := xenapi.Host.GetAll(s)
	if err != nil {
		return "", xapi.Sanitize(err)
	}
	sr, err := xenapi.SR.GetByUUID(s, uuid)
	reattach := err == nil
	if !reattach {
		ref, err := waitRef(s)(xenapi.SR.AsyncIntroduce(s, uuid, name, description, rules.SRType(kind), "", true, map[string]string{}))
		if err != nil {
			return "", xapi.Sanitize(err, password)
		}
		sr = xenapi.SRRef(ref)
	}
	var made []xenapi.PBDRef
	defer func() {
		if err != nil && !errors.Is(err, ErrUncertain) {
			removePBDs(s, made) // RollbackFailedReattach: only what this attempt made
		}
	}()
	for _, host := range hosts {
		pbd, err := waitRef(s)(xenapi.PBD.AsyncCreate(s, xenapi.PBDRecord{Host: host, SR: sr, DeviceConfig: dconf, OtherConfig: map[string]string{}}))
		if err != nil {
			return "", xapi.Sanitize(err, password)
		}
		made = append(made, xenapi.PBDRef(pbd))
		if err := wait(s)(xenapi.PBD.AsyncPlug(s, xenapi.PBDRef(pbd))); err != nil {
			return "", xapi.Sanitize(err, password)
		}
	}
	if reattach {
		if err := xenapi.SR.SetNameLabel(s, sr, name); err != nil {
			return "", xapi.Sanitize(err)
		}
		if err := xenapi.SR.SetNameDescription(s, sr, description); err != nil {
			return "", xapi.Sanitize(err)
		}
	} else if firstSharedSR(s, sr) {
		if err := setDefaultSR(s, sr); err != nil {
			return "", err
		}
	}
	return uuid, nil
}

// removePBDs is DetachSrAction with destroyPbds, on the given PBDs: unplug and destroy them; failures are ignored.
func removePBDs(s *xenapi.Session, pbds []xenapi.PBDRef) {
	for _, pbd := range pbds {
		if attached, err := xenapi.PBD.GetCurrentlyAttached(s, pbd); err == nil && attached {
			_ = xenapi.PBD.Unplug(s, pbd)
		}
		_ = xenapi.PBD.Destroy(s, pbd)
	}
}
