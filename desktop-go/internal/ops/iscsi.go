package ops

import (
	"errors"
	"strconv"

	xenapi "go/xenapi"

	"xcpng-center/internal/messages"
	"xcpng-center/internal/rules"
	"xcpng-center/internal/xapi"
)

// ScanIQNs is ISCSIPopulateIQNsAction (LVMoISCSI): an SR.create the SM refuses on purpose, with the target only;
// its SR_BACKEND_FAILURE_96 carries the target's IQNs. The CHAP password goes only into that call.
func ScanIQNs(c *xapi.Connection, target string, port int, chapUser, chapPassword string) ([]rules.IQN, error) {
	dconf := map[string]string{"target": target, "port": strconv.Itoa(port)}
	if chapUser != "" {
		dconf["chapuser"], dconf["chappassword"] = chapUser, chapPassword
	}
	doc, err := iscsiProbeFailure(c, dconf, chapPassword, "SR_BACKEND_FAILURE_96")
	if err != nil {
		return nil, err
	}
	iqns, err := rules.ParseIQNs(doc)
	if err != nil {
		return nil, errors.New("the iSCSI target answered with something unexpected")
	}
	var out []rules.IQN
	for _, i := range iqns {
		if i.IQN != "" {
			out = append(out, i)
		}
	}
	if len(out) == 0 {
		return nil, errors.New(messages.Format(messages.Get("NEWSR_NO_IQNS_FOUND"), target))
	}
	return out, nil
}

// ScanLUNs is ISCSIPopulateLunsAction (LVMoISCSI): the same refused SR.create with the IQN, whose
// SR_BACKEND_FAILURE_107 (or _87) carries the LUNs.
func ScanLUNs(c *xapi.Connection, iqn rules.IQN, chapUser, chapPassword string) ([]rules.LUN, error) {
	dconf := map[string]string{"target": iqn.Address, "port": strconv.Itoa(iqn.Port), "targetIQN": iqn.IQN}
	if chapUser != "" {
		dconf["chapuser"], dconf["chappassword"] = chapUser, chapPassword
	}
	doc, err := iscsiProbeFailure(c, dconf, chapPassword, "SR_BACKEND_FAILURE_107", "SR_BACKEND_FAILURE_87")
	if err != nil {
		return nil, err
	}
	luns, err := rules.ParseLUNs(doc)
	if err != nil {
		return nil, errors.New("the iSCSI target answered with something unexpected")
	}
	if len(luns) == 0 {
		return nil, errors.New(messages.Format(messages.Get("NEWSR_NO_LUNS_FOUND"), iqn.Address))
	}
	return luns, nil
}

// iscsiProbeFailure runs the refused SR.create and returns the XML of the expected failure (its third parameter).
// Another failure is the scan's error (SrRules.IscsiError's texts for login, host and connection); should the SR be
// made after all, it is forgotten again.
func iscsiProbeFailure(c *xapi.Connection, dconf map[string]string, password string, expected ...string) (string, error) {
	s := c.Session
	coordinator, err := coordinatorOf(s)
	if err != nil {
		return "", err
	}
	ref, err := waitRef(s)(xenapi.SR.AsyncCreate(s, coordinator, dconf, 0, "__gui__", messages.Get("ISCSI_SHOULD_NO_BE_CREATED"),
		"lvmoiscsi", "user", true, map[string]string{}))
	if err == nil {
		forgetSR(s, xenapi.SRRef(ref))
		return "", errors.New("the iSCSI scan created an SR: it was forgotten again")
	}
	var api *xenapi.APIError
	if !errors.As(err, &api) {
		return "", xapi.Sanitize(err, password)
	}
	for _, code := range expected {
		if api.Message == code {
			if len(api.Data) < 3 || api.Data[2] == "" {
				return "", nil
			}
			return api.Data[2], nil
		}
	}
	if text := rules.ISCSIError(api.Message); text != "" {
		return "", errors.New(text)
	}
	return "", xapi.Sanitize(err, password)
}
