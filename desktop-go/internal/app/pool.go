package app

import (
	xenapi "go/xenapi"

	"xcpng-center/internal/rules"
	"xcpng-center/internal/xapi"
)

// ReadPoolFacts reads what PoolJoinRules needs of a connected server: its pool, coordinator, SRs and VMs.
func ReadPoolFacts(c *xapi.Connection, key string) (rules.PoolFacts, error) {
	s := c.Session
	f := rules.PoolFacts{Key: key}
	pools, err := xenapi.Pool.GetAllRecords(s)
	if err != nil {
		return f, xapi.Sanitize(err)
	}
	hosts, err := xenapi.Host.GetAllRecords(s)
	if err != nil {
		return f, xapi.Sanitize(err)
	}
	f.Hosts = len(hosts)
	f.ServerPacks = map[string]map[string]string{}
	for _, h := range hosts {
		f.ServerPacks[h.NameLabel] = h.SoftwareVersion
	}
	for _, p := range pools {
		f.PoolName, f.HA, f.TLSVerification, f.PoolCPU = p.NameLabel, p.HaEnabled, p.TLSVerificationEnabled, p.CPUInfo
		if h, ok := hosts[p.Master]; ok {
			f.Name, f.CPUVendor, f.CPU, f.APIMajor, f.APIMinor = h.NameLabel, h.CPUInfo["vendor"], h.CPUInfo, h.APIVersionMajor, h.APIVersionMinor
			f.DBSchema, f.PlatformVersion, f.NetworkBackend = h.SoftwareVersion["db_schema"], h.SoftwareVersion["platform_version"], h.SoftwareVersion["network_backend"]
		}
	}
	srs, err := xenapi.SR.GetAllRecords(s)
	if err != nil {
		return f, xapi.Sanitize(err)
	}
	for _, sr := range srs {
		f.SharedStorage = f.SharedStorage || sr.Shared && !sr.IsToolsSr && sr.NameLabel != "XenServer Tools" // SR.IsToolsSR
	}
	vms, err := xenapi.VM.GetAllRecords(s)
	if err != nil {
		return f, xapi.Sanitize(err)
	}
	for _, vm := range vms {
		f.RunningVMs = f.RunningVMs || !vm.IsATemplate && !vm.IsASnapshot && !vm.IsControlDomain && vm.PowerState == xenapi.VMPowerStateRunning
	}
	clusters, err := xenapi.Cluster.GetAll(s)
	f.Clustered = err == nil && len(clusters) > 0
	return f, nil
}
