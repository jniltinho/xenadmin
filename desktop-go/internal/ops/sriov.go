package ops

import (
	"errors"
	"slices"
	"strings"

	xenapi "go/xenapi"

	"xcpng-center/internal/rules"
	"xcpng-center/internal/xapi"
)

// SR-IOV networks: the New Network wizard's SR-IOV Network (XenAdmin/Wizards/NewNetworkWizard_Pages/NetWTypeSelect.cs,
// NetWSriovDetails.cs, NewNetworkWizard.CreateSRIOV) and CreateSriovAction (XenModel/Actions/Network).

// sriovCandidate is NetWSriovDetails' NIC test (rules.SriovCandidate) on a PIF.
func sriovCandidate(p xenapi.PIFRecord) bool {
	return rules.SriovCandidate(physical(p), len(p.BondMasterOf) > 0, rules.SriovCapable(p.Capabilities), len(p.SriovPhysicalPIFOf) > 0)
}

// sriovOption is NetWTypeSelect.PopulatePage's SR-IOV choice on this state.
func (d *netState) sriovOption() (bool, string) {
	experimental := false // Host.SriovNetworkDisabled: restricted and the network_sriov feature off, on a server
	if features, err := xenapi.Feature.GetAllRecords(d.s); err == nil {
		for _, f := range features {
			h, ok := d.hosts[f.Host]
			experimental = experimental || ok && strings.EqualFold(f.NameLabel, "network_sriov") && !f.Enabled &&
				featureForbidden(map[xenapi.HostRef]xenapi.HostRecord{f.Host: h}, "restrict_network_sriov")
		}
	}
	var capable, free bool
	for _, p := range d.pifs {
		capable = capable || rules.SriovCapable(p.Capabilities) // Pool.HasSriovNic
		free = free || physical(p) && rules.SriovCapable(p.Capabilities) && len(p.SriovPhysicalPIFOf) == 0
	}
	return rules.SriovOption(experimental, featureForbidden(d.hosts, "restrict_network_sriov"), capable, free)
}

// SriovChoices is the wizard's SR-IOV Network choice (enabled, else why not) and the NICs of the server it offers.
func SriovChoices(s *xenapi.Session, hostUUID string) (ok bool, why string, nics []xenapi.PIFRecord, err error) {
	d, err := readNet(s)
	if err != nil {
		return false, "", nil, err
	}
	ok, why = d.sriovOption()
	for _, p := range d.pifs {
		if d.hosts[p.Host].UUID == hostUUID && sriovCandidate(p) {
			nics = append(nics, p)
		}
	}
	return ok, why, nics, nil
}

// NewSriov is the wizard's SR-IOV network: its name (NetWName), the NIC of the wizard's server and whether new VMs get
// it (NetWSriovDetails).
type NewSriov struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Auto        bool   `json:"auto"`
	Host        string `json:"host"` // the wizard's server (NetworkRules.WizardHost)
	NIC         string `json:"nic"`  // the PIF uuid on that server
}

// CreateSriov is NewNetworkWizard.CreateSRIOV then CreateSriovAction: a network, then SR-IOV enabled on the NIC's
// device on every server where it can be (the coordinator first, which xapi needs); the network is destroyed when the
// first server fails. The servers' connections are disturbed and they may need a reboot.
func CreateSriov(c *xapi.Connection, spec NewSriov) error {
	s := c.Session
	d, err := readNet(s)
	if err != nil {
		return err
	}
	name := strings.TrimSpace(spec.Name)
	if name == "" {
		return errors.New("enter a name for the new network")
	}
	if ok, why := d.sriovOption(); !ok {
		return errors.New(why)
	}
	device := ""
	for _, p := range d.pifs {
		if p.UUID == spec.NIC && d.hosts[p.Host].UUID == spec.Host && sriovCandidate(p) {
			device = p.Device
		}
	}
	if device == "" {
		return errors.New("select a NIC that can enable SR-IOV")
	}
	var pifs []xenapi.PIFRef
	for r, p := range d.pifs {
		if p.Device == device && sriovCandidate(p) {
			pifs = append(pifs, r)
		}
	}
	slices.SortStableFunc(pifs, func(a, b xenapi.PIFRef) int { // the coordinator's first, then by server
		ca, cb := d.pifs[a].Host == d.coordinator, d.pifs[b].Host == d.coordinator
		if ca != cb {
			return map[bool]int{true: -1, false: 1}[ca]
		}
		return strings.Compare(string(d.pifs[a].Host), string(d.pifs[b].Host))
	})
	if err := permitted(s, "network.create", "network_sriov.create"); err != nil {
		return err
	}
	ref, err := xenapi.Network.Create(s, xenapi.NetworkRecord{NameLabel: name, NameDescription: spec.Description, Managed: true,
		MTU: rules.MTUDefault, OtherConfig: map[string]string{"automatic": map[bool]string{true: "true", false: "false"}[spec.Auto]}})
	if err != nil {
		return xapi.Sanitize(err)
	}
	for i, p := range pifs {
		if err := wait(s)(xenapi.NetworkSriov.AsyncCreate(s, p, ref)); err != nil {
			if api := (*xenapi.APIError)(nil); i == 0 && errors.As(err, &api) { // a lost reply may have enabled it: kept
				_ = xenapi.Network.Destroy(s, ref) // CreateSriovAction.DestroyNetwork: best effort
			}
			return xapi.Sanitize(err)
		}
	}
	return nil
}
