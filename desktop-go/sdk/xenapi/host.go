/*
 * Copyright (c) Cloud Software Group, Inc.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions
 * are met:
 *
 *   1) Redistributions of source code must retain the above copyright
 *      notice, this list of conditions and the following disclaimer.
 *
 *   2) Redistributions in binary form must reproduce the above
 *      copyright notice, this list of conditions and the following
 *      disclaimer in the documentation and/or other materials
 *      provided with the distribution.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
 * "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
 * LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS
 * FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE
 * COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT,
 * INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
 * (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
 * SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION)
 * HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT,
 * STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
 * ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED
 * OF THE POSSIBILITY OF SUCH DAMAGE.
 */

package xenapi

import (
	"fmt"
	"time"
)

type HostRecord struct {
	// Unique identifier/object reference
	UUID string `json:"uuid,omitempty"`
	// a human-readable name
	NameLabel string `json:"name_label,omitempty"`
	// a notes field containing human-readable description
	NameDescription string `json:"name_description,omitempty"`
	// Virtualization memory overhead (bytes).
	MemoryOverhead int `json:"memory_overhead,omitempty"`
	// list of the operations allowed in this state. This list is advisory only and the server state may have changed by the time this field is read by a client.
	AllowedOperations []HostAllowedOperations `json:"allowed_operations,omitempty"`
	// links each of the running tasks using this object (by reference) to a current_operation enum which describes the nature of the task.
	CurrentOperations map[string]HostAllowedOperations `json:"current_operations,omitempty"`
	// major version number
	APIVersionMajor int `json:"API_version_major,omitempty"`
	// minor version number
	APIVersionMinor int `json:"API_version_minor,omitempty"`
	// identification of vendor
	APIVersionVendor string `json:"API_version_vendor,omitempty"`
	// details of vendor implementation
	APIVersionVendorImplementation map[string]string `json:"API_version_vendor_implementation,omitempty"`
	// True if the host is currently enabled
	Enabled bool `json:"enabled,omitempty"`
	// version strings
	SoftwareVersion map[string]string `json:"software_version,omitempty"`
	// additional configuration
	OtherConfig map[string]string `json:"other_config,omitempty"`
	// Xen capabilities
	Capabilities []string `json:"capabilities,omitempty"`
	// The CPU configuration on this host.  May contain keys such as &quot;nr_nodes&quot;, &quot;sockets_per_node&quot;, &quot;cores_per_socket&quot;, or &quot;threads_per_core&quot;
	CPUConfiguration map[string]string `json:"cpu_configuration,omitempty"`
	// Scheduler policy currently in force on this host
	SchedPolicy string `json:"sched_policy,omitempty"`
	// a list of the bootloaders installed on the machine
	SupportedBootloaders []string `json:"supported_bootloaders,omitempty"`
	// list of VMs currently resident on host
	ResidentVMs []VMRef `json:"resident_VMs,omitempty"`
	// logging configuration
	Logging map[string]string `json:"logging,omitempty"`
	// physical network interfaces
	PIFs []PIFRef `json:"PIFs,omitempty"`
	// The SR in which VDIs for suspend images are created
	SuspendImageSr SRRef `json:"suspend_image_sr,omitempty"`
	// The SR in which VDIs for crash dumps are created
	CrashDumpSr SRRef `json:"crash_dump_sr,omitempty"`
	// Set of host crash dumps
	Crashdumps []HostCrashdumpRef `json:"crashdumps,omitempty"`
	// Set of host patches
	Patches []HostPatchRef `json:"patches,omitempty"`
	// Set of updates
	Updates []PoolUpdateRef `json:"updates,omitempty"`
	// physical blockdevices
	PBDs []PBDRef `json:"PBDs,omitempty"`
	// The physical CPUs on this host
	HostCPUs []HostCPURef `json:"host_CPUs,omitempty"`
	// Details about the physical CPUs on this host
	CPUInfo map[string]string `json:"cpu_info,omitempty"`
	// The hostname of this host
	Hostname string `json:"hostname,omitempty"`
	// The address by which this host can be contacted from any other host in the pool
	Address string `json:"address,omitempty"`
	// metrics associated with this host
	Metrics HostMetricsRef `json:"metrics,omitempty"`
	// State of the current license
	LicenseParams map[string]string `json:"license_params,omitempty"`
	// The set of statefiles accessible from this host
	HaStatefiles []string `json:"ha_statefiles,omitempty"`
	// The set of hosts visible via the network from this host
	HaNetworkPeers []string `json:"ha_network_peers,omitempty"`
	// Binary blobs associated with this host
	Blobs map[string]BlobRef `json:"blobs,omitempty"`
	// user-specified tags for categorization purposes
	Tags []string `json:"tags,omitempty"`
	// type of external authentication service configured; empty if none configured.
	ExternalAuthType string `json:"external_auth_type,omitempty"`
	// name of external authentication service configured; empty if none configured.
	ExternalAuthServiceName string `json:"external_auth_service_name,omitempty"`
	// configuration specific to external authentication service
	ExternalAuthConfiguration map[string]string `json:"external_auth_configuration,omitempty"`
	// Product edition
	Edition string `json:"edition,omitempty"`
	// Contact information of the license server
	LicenseServer map[string]string `json:"license_server,omitempty"`
	// BIOS strings
	BiosStrings map[string]string `json:"bios_strings,omitempty"`
	// The power on mode
	PowerOnMode string `json:"power_on_mode,omitempty"`
	// The power on config
	PowerOnConfig map[string]string `json:"power_on_config,omitempty"`
	// The SR that is used as a local cache
	LocalCacheSr SRRef `json:"local_cache_sr,omitempty"`
	// Information about chipset features
	ChipsetInfo map[string]string `json:"chipset_info,omitempty"`
	// List of PCI devices in the host
	PCIs []PCIRef `json:"PCIs,omitempty"`
	// List of physical GPUs in the host
	PGPUs []PGPURef `json:"PGPUs,omitempty"`
	// List of physical USBs in the host
	PUSBs []PUSBRef `json:"PUSBs,omitempty"`
	// Allow SSLv3 protocol and ciphersuites as used by older server versions. This controls both incoming and outgoing connections. When this is set to a different value, the host immediately restarts its SSL/TLS listening service; typically this takes less than a second but existing connections to it will be broken. API login sessions will remain valid.
	SslLegacy bool `json:"ssl_legacy,omitempty"`
	// VCPUs params to apply to all resident guests
	GuestVCPUsParams map[string]string `json:"guest_VCPUs_params,omitempty"`
	// indicates whether the host is configured to output its console to a physical display device
	Display HostDisplay `json:"display,omitempty"`
	// The set of versions of the virtual hardware platform that the host can offer to its guests
	VirtualHardwarePlatformVersions []int `json:"virtual_hardware_platform_versions,omitempty"`
	// The control domain (domain 0)
	ControlDomain VMRef `json:"control_domain,omitempty"`
	// List of updates which require reboot
	UpdatesRequiringReboot []PoolUpdateRef `json:"updates_requiring_reboot,omitempty"`
	// List of features available on this host
	Features []FeatureRef `json:"features,omitempty"`
	// The initiator IQN for the host
	IscsiIqn string `json:"iscsi_iqn,omitempty"`
	// Specifies whether multipathing is enabled
	Multipathing bool `json:"multipathing,omitempty"`
	// The UEFI certificates allowing Secure Boot
	UefiCertificates string `json:"uefi_certificates,omitempty"`
	// List of certificates installed in the host
	Certificates []CertificateRef `json:"certificates,omitempty"`
	// List of all available product editions
	Editions []string `json:"editions,omitempty"`
	// The set of pending mandatory guidances after applying updates, which must be applied, otherwise there may be, for example, VM failures
	PendingGuidances []UpdateGuidances `json:"pending_guidances,omitempty"`
	// True if this host has TLS verifcation enabled
	TLSVerificationEnabled bool `json:"tls_verification_enabled,omitempty"`
	// Date and time when the last software update was applied
	LastSoftwareUpdate time.Time `json:"last_software_update,omitempty"`
	// Reflects whether port 80 is open (false) or not (true)
	HTTPSOnly bool `json:"https_only,omitempty"`
	// Default as &apos;unknown&apos;, &apos;yes&apos; if the host is up to date with updates synced from remote CDN, otherwise &apos;no&apos;
	LatestSyncedUpdatesApplied LatestSyncedUpdatesAppliedState `json:"latest_synced_updates_applied,omitempty"`
	// NUMA-aware VM memory and vCPU placement policy
	NumaAffinityPolicy HostNumaAffinityPolicy `json:"numa_affinity_policy,omitempty"`
	// The set of pending recommended guidances after applying updates, which most users should follow to make the updates effective, but if not followed, will not cause a failure
	PendingGuidancesRecommended []UpdateGuidances `json:"pending_guidances_recommended,omitempty"`
	// The set of pending full guidances after applying updates, which a user should follow to make some updates, for example, specific hardware drivers or CPU features, fully effective, but the &apos;average user&apos; doesn&apos;t need to
	PendingGuidancesFull []UpdateGuidances `json:"pending_guidances_full,omitempty"`
	// The SHA256 checksum of updateinfo of the most recently applied update on the host
	LastUpdateHash string `json:"last_update_hash,omitempty"`
	// True if SSH access is enabled for the host
	SSHEnabled bool `json:"ssh_enabled,omitempty"`
	// The timeout in seconds after which SSH access will be automatically disabled (0 means never), this setting will be applied every time the SSH is enabled by XAPI
	SSHEnabledTimeout int `json:"ssh_enabled_timeout,omitempty"`
	// The time in UTC after which the SSH access will be automatically disabled
	SSHExpiry time.Time `json:"ssh_expiry,omitempty"`
	// The timeout in seconds after which idle console will be automatically terminated (0 means never)
	ConsoleIdleTimeout int `json:"console_idle_timeout,omitempty"`
	// Reflects whether SSH auto mode is enabled for the host
	SSHAutoMode bool `json:"ssh_auto_mode,omitempty"`
	// The maximum C-state that the host is allowed to enter, &quot;&quot; means unlimited; &quot;N&quot; means limit to CN; &quot;N,M&quot; means limit to CN with max sub cstate M.
	MaxCstate string `json:"max_cstate,omitempty"`
	// Whether the host has booted in secure boot mode
	SecureBoot bool `json:"secure_boot,omitempty"`
	// Indicates NTP servers are assigned by DHCP, or configured by user, or the factory servers, or NTP is disabled
	NtpMode HostNtpMode `json:"ntp_mode,omitempty"`
	// Custom NTP servers configured by users, used in Custom NTP mode
	NtpCustomServers []string `json:"ntp_custom_servers,omitempty"`
	// The time zone identifier as defined in the IANA Time Zone Database
	Timezone string `json:"timezone,omitempty"`
}

type HostRef string

// A physical host
type host struct{}

var Host host

// GetAllRecords: Return a map of host references to host records for all hosts known to the system.
// Version: rio
func (host) GetAllRecords(session *Session) (retval map[HostRef]HostRecord, err error) {
	method := "host.get_all_records"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializeHostRefToHostRecordMap(method+" -> ", result)
	return
}

// GetAllRecords1: Return a map of host references to host records for all hosts known to the system.
// Version: rio
func (host) GetAllRecords1(session *Session) (retval map[HostRef]HostRecord, err error) {
	method := "host.get_all_records"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializeHostRefToHostRecordMap(method+" -> ", result)
	return
}

// GetAll: Return a list of all the hosts known to the system.
// Version: rio
func (host) GetAll(session *Session) (retval []HostRef, err error) {
	method := "host.get_all"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializeHostRefSet(method+" -> ", result)
	return
}

// GetAll1: Return a list of all the hosts known to the system.
// Version: rio
func (host) GetAll1(session *Session) (retval []HostRef, err error) {
	method := "host.get_all"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializeHostRefSet(method+" -> ", result)
	return
}

// SetServertime: Set the host&apos;s system clock when NTP is disabled.
// Version: closed
func (host) SetServertime(session *Session, self HostRef, value time.Time) (err error) {
	method := "host.set_servertime"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeTime(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	return
}

// AsyncSetServertime: Set the host&apos;s system clock when NTP is disabled.
// Version: closed
func (host) AsyncSetServertime(session *Session, self HostRef, value time.Time) (retval TaskRef, err error) {
	method := "Async.host.set_servertime"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeTime(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// SetServertime3: Set the host&apos;s system clock when NTP is disabled.
// Version: closed
func (host) SetServertime3(session *Session, self HostRef, value time.Time) (err error) {
	method := "host.set_servertime"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeTime(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	return
}

// AsyncSetServertime3: Set the host&apos;s system clock when NTP is disabled.
// Version: closed
func (host) AsyncSetServertime3(session *Session, self HostRef, value time.Time) (retval TaskRef, err error) {
	method := "Async.host.set_servertime"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeTime(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// GetNtpSynchronized: Returns true if the system clock on the host is synchronized with the NTP servers.
// Version: closed
func (host) GetNtpSynchronized(session *Session, self HostRef) (retval bool, err error) {
	method := "host.get_ntp_synchronized"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeBool(method+" -> ", result)
	return
}

// AsyncGetNtpSynchronized: Returns true if the system clock on the host is synchronized with the NTP servers.
// Version: closed
func (host) AsyncGetNtpSynchronized(session *Session, self HostRef) (retval TaskRef, err error) {
	method := "Async.host.get_ntp_synchronized"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// GetNtpSynchronized2: Returns true if the system clock on the host is synchronized with the NTP servers.
// Version: closed
func (host) GetNtpSynchronized2(session *Session, self HostRef) (retval bool, err error) {
	method := "host.get_ntp_synchronized"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeBool(method+" -> ", result)
	return
}

// AsyncGetNtpSynchronized2: Returns true if the system clock on the host is synchronized with the NTP servers.
// Version: closed
func (host) AsyncGetNtpSynchronized2(session *Session, self HostRef) (retval TaskRef, err error) {
	method := "Async.host.get_ntp_synchronized"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// ListTimezones: List all available timezones on the host.
// Version: closed
func (host) ListTimezones(session *Session, self HostRef) (retval []string, err error) {
	method := "host.list_timezones"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeStringSet(method+" -> ", result)
	return
}

// AsyncListTimezones: List all available timezones on the host.
// Version: closed
func (host) AsyncListTimezones(session *Session, self HostRef) (retval TaskRef, err error) {
	method := "Async.host.list_timezones"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// ListTimezones2: List all available timezones on the host.
// Version: closed
func (host) ListTimezones2(session *Session, self HostRef) (retval []string, err error) {
	method := "host.list_timezones"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeStringSet(method+" -> ", result)
	return
}

// AsyncListTimezones2: List all available timezones on the host.
// Version: closed
func (host) AsyncListTimezones2(session *Session, self HostRef) (retval TaskRef, err error) {
	method := "Async.host.list_timezones"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// SetTimezone: Set the host&apos;s timezone.
// Version: closed
func (host) SetTimezone(session *Session, self HostRef, value string) (err error) {
	method := "host.set_timezone"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	return
}

// AsyncSetTimezone: Set the host&apos;s timezone.
// Version: closed
func (host) AsyncSetTimezone(session *Session, self HostRef, value string) (retval TaskRef, err error) {
	method := "Async.host.set_timezone"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// SetTimezone3: Set the host&apos;s timezone.
// Version: closed
func (host) SetTimezone3(session *Session, self HostRef, value string) (err error) {
	method := "host.set_timezone"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	return
}

// AsyncSetTimezone3: Set the host&apos;s timezone.
// Version: closed
func (host) AsyncSetTimezone3(session *Session, self HostRef, value string) (retval TaskRef, err error) {
	method := "Async.host.set_timezone"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// GetNtpServersStatus: Get the NTP servers status on the host
// Version: closed
func (host) GetNtpServersStatus(session *Session, self HostRef) (retval map[string]string, err error) {
	method := "host.get_ntp_servers_status"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeStringToStringMap(method+" -> ", result)
	return
}

// AsyncGetNtpServersStatus: Get the NTP servers status on the host
// Version: closed
func (host) AsyncGetNtpServersStatus(session *Session, self HostRef) (retval TaskRef, err error) {
	method := "Async.host.get_ntp_servers_status"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// GetNtpServersStatus2: Get the NTP servers status on the host
// Version: closed
func (host) GetNtpServersStatus2(session *Session, self HostRef) (retval map[string]string, err error) {
	method := "host.get_ntp_servers_status"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeStringToStringMap(method+" -> ", result)
	return
}

// AsyncGetNtpServersStatus2: Get the NTP servers status on the host
// Version: closed
func (host) AsyncGetNtpServersStatus2(session *Session, self HostRef) (retval TaskRef, err error) {
	method := "Async.host.get_ntp_servers_status"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// SetNtpCustomServers: Set the custom NTP servers for the host
// Version: closed
func (host) SetNtpCustomServers(session *Session, self HostRef, value []string) (err error) {
	method := "host.set_ntp_custom_servers"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeStringSet(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	return
}

// AsyncSetNtpCustomServers: Set the custom NTP servers for the host
// Version: closed
func (host) AsyncSetNtpCustomServers(session *Session, self HostRef, value []string) (retval TaskRef, err error) {
	method := "Async.host.set_ntp_custom_servers"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeStringSet(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// SetNtpCustomServers3: Set the custom NTP servers for the host
// Version: closed
func (host) SetNtpCustomServers3(session *Session, self HostRef, value []string) (err error) {
	method := "host.set_ntp_custom_servers"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeStringSet(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	return
}

// AsyncSetNtpCustomServers3: Set the custom NTP servers for the host
// Version: closed
func (host) AsyncSetNtpCustomServers3(session *Session, self HostRef, value []string) (retval TaskRef, err error) {
	method := "Async.host.set_ntp_custom_servers"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeStringSet(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// SetNtpMode: Set the NTP mode for the host
// Version: closed
func (host) SetNtpMode(session *Session, self HostRef, value HostNtpMode) (err error) {
	method := "host.set_ntp_mode"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeEnumHostNtpMode(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	return
}

// AsyncSetNtpMode: Set the NTP mode for the host
// Version: closed
func (host) AsyncSetNtpMode(session *Session, self HostRef, value HostNtpMode) (retval TaskRef, err error) {
	method := "Async.host.set_ntp_mode"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeEnumHostNtpMode(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// SetNtpMode3: Set the NTP mode for the host
// Version: closed
func (host) SetNtpMode3(session *Session, self HostRef, value HostNtpMode) (err error) {
	method := "host.set_ntp_mode"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeEnumHostNtpMode(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	return
}

// AsyncSetNtpMode3: Set the NTP mode for the host
// Version: closed
func (host) AsyncSetNtpMode3(session *Session, self HostRef, value HostNtpMode) (retval TaskRef, err error) {
	method := "Async.host.set_ntp_mode"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeEnumHostNtpMode(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// SetMaxCstate: Sets xen&apos;s max-cstate on a host. See: https://xenbits.xen.org/docs/unstable/misc/xen-command-line.html#max_cstate-x86. &quot;&quot; means unlimited; &quot;N&quot; means limit to CN; &quot;N,M&quot; means limit to CN with max sub cstate M. Note: Only C0, C1, unlimited are supported currently.
// Version: closed
func (host) SetMaxCstate(session *Session, self HostRef, value string) (err error) {
	method := "host.set_max_cstate"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	return
}

// AsyncSetMaxCstate: Sets xen&apos;s max-cstate on a host. See: https://xenbits.xen.org/docs/unstable/misc/xen-command-line.html#max_cstate-x86. &quot;&quot; means unlimited; &quot;N&quot; means limit to CN; &quot;N,M&quot; means limit to CN with max sub cstate M. Note: Only C0, C1, unlimited are supported currently.
// Version: closed
func (host) AsyncSetMaxCstate(session *Session, self HostRef, value string) (retval TaskRef, err error) {
	method := "Async.host.set_max_cstate"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// SetMaxCstate3: Sets xen&apos;s max-cstate on a host. See: https://xenbits.xen.org/docs/unstable/misc/xen-command-line.html#max_cstate-x86. &quot;&quot; means unlimited; &quot;N&quot; means limit to CN; &quot;N,M&quot; means limit to CN with max sub cstate M. Note: Only C0, C1, unlimited are supported currently.
// Version: closed
func (host) SetMaxCstate3(session *Session, self HostRef, value string) (err error) {
	method := "host.set_max_cstate"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	return
}

// AsyncSetMaxCstate3: Sets xen&apos;s max-cstate on a host. See: https://xenbits.xen.org/docs/unstable/misc/xen-command-line.html#max_cstate-x86. &quot;&quot; means unlimited; &quot;N&quot; means limit to CN; &quot;N,M&quot; means limit to CN with max sub cstate M. Note: Only C0, C1, unlimited are supported currently.
// Version: closed
func (host) AsyncSetMaxCstate3(session *Session, self HostRef, value string) (retval TaskRef, err error) {
	method := "Async.host.set_max_cstate"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// GetTrackedUserAgents: Get the (name, version) list of tracked user agents on this host. If different versions of the same name are seen, keep the last-seen version. The oldest entry will be removed if reach the max num. Note that the list is cleared after host/XAPI restart
// Version: closed
func (host) GetTrackedUserAgents(session *Session, self HostRef) (retval map[string]string, err error) {
	method := "host.get_tracked_user_agents"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeStringToStringMap(method+" -> ", result)
	return
}

// AsyncGetTrackedUserAgents: Get the (name, version) list of tracked user agents on this host. If different versions of the same name are seen, keep the last-seen version. The oldest entry will be removed if reach the max num. Note that the list is cleared after host/XAPI restart
// Version: closed
func (host) AsyncGetTrackedUserAgents(session *Session, self HostRef) (retval TaskRef, err error) {
	method := "Async.host.get_tracked_user_agents"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// GetTrackedUserAgents2: Get the (name, version) list of tracked user agents on this host. If different versions of the same name are seen, keep the last-seen version. The oldest entry will be removed if reach the max num. Note that the list is cleared after host/XAPI restart
// Version: closed
func (host) GetTrackedUserAgents2(session *Session, self HostRef) (retval map[string]string, err error) {
	method := "host.get_tracked_user_agents"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeStringToStringMap(method+" -> ", result)
	return
}

// AsyncGetTrackedUserAgents2: Get the (name, version) list of tracked user agents on this host. If different versions of the same name are seen, keep the last-seen version. The oldest entry will be removed if reach the max num. Note that the list is cleared after host/XAPI restart
// Version: closed
func (host) AsyncGetTrackedUserAgents2(session *Session, self HostRef) (retval TaskRef, err error) {
	method := "Async.host.get_tracked_user_agents"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// SetSSHAutoMode: Set the SSH auto mode for the host
// Version: closed
func (host) SetSSHAutoMode(session *Session, self HostRef, value bool) (err error) {
	method := "host.set_ssh_auto_mode"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	return
}

// AsyncSetSSHAutoMode: Set the SSH auto mode for the host
// Version: closed
func (host) AsyncSetSSHAutoMode(session *Session, self HostRef, value bool) (retval TaskRef, err error) {
	method := "Async.host.set_ssh_auto_mode"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// SetSSHAutoMode3: Set the SSH auto mode for the host
// Version: closed
func (host) SetSSHAutoMode3(session *Session, self HostRef, value bool) (err error) {
	method := "host.set_ssh_auto_mode"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	return
}

// AsyncSetSSHAutoMode3: Set the SSH auto mode for the host
// Version: closed
func (host) AsyncSetSSHAutoMode3(session *Session, self HostRef, value bool) (retval TaskRef, err error) {
	method := "Async.host.set_ssh_auto_mode"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// SetConsoleIdleTimeout: Set the console idle timeout for the host
// Version: closed
func (host) SetConsoleIdleTimeout(session *Session, self HostRef, value int) (err error) {
	method := "host.set_console_idle_timeout"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	return
}

// AsyncSetConsoleIdleTimeout: Set the console idle timeout for the host
// Version: closed
func (host) AsyncSetConsoleIdleTimeout(session *Session, self HostRef, value int) (retval TaskRef, err error) {
	method := "Async.host.set_console_idle_timeout"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// SetConsoleIdleTimeout3: Set the console idle timeout for the host
// Version: closed
func (host) SetConsoleIdleTimeout3(session *Session, self HostRef, value int) (err error) {
	method := "host.set_console_idle_timeout"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	return
}

// AsyncSetConsoleIdleTimeout3: Set the console idle timeout for the host
// Version: closed
func (host) AsyncSetConsoleIdleTimeout3(session *Session, self HostRef, value int) (retval TaskRef, err error) {
	method := "Async.host.set_console_idle_timeout"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// SetSSHEnabledTimeout: Set the SSH service enabled timeout for the host
// Version: closed
func (host) SetSSHEnabledTimeout(session *Session, self HostRef, value int) (err error) {
	method := "host.set_ssh_enabled_timeout"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	return
}

// AsyncSetSSHEnabledTimeout: Set the SSH service enabled timeout for the host
// Version: closed
func (host) AsyncSetSSHEnabledTimeout(session *Session, self HostRef, value int) (retval TaskRef, err error) {
	method := "Async.host.set_ssh_enabled_timeout"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// SetSSHEnabledTimeout3: Set the SSH service enabled timeout for the host
// Version: closed
func (host) SetSSHEnabledTimeout3(session *Session, self HostRef, value int) (err error) {
	method := "host.set_ssh_enabled_timeout"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	return
}

// AsyncSetSSHEnabledTimeout3: Set the SSH service enabled timeout for the host
// Version: closed
func (host) AsyncSetSSHEnabledTimeout3(session *Session, self HostRef, value int) (retval TaskRef, err error) {
	method := "Async.host.set_ssh_enabled_timeout"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// DisableSSH: Disable SSH access on the host. It will stop the service sshd only if it is running. It will also disable the service sshd only if it is enabled. A newly joined host in the pool or an ejected host from the pool would keep the original status.
// Version: closed
func (host) DisableSSH(session *Session, self HostRef) (err error) {
	method := "host.disable_ssh"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// AsyncDisableSSH: Disable SSH access on the host. It will stop the service sshd only if it is running. It will also disable the service sshd only if it is enabled. A newly joined host in the pool or an ejected host from the pool would keep the original status.
// Version: closed
func (host) AsyncDisableSSH(session *Session, self HostRef) (retval TaskRef, err error) {
	method := "Async.host.disable_ssh"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// DisableSSH2: Disable SSH access on the host. It will stop the service sshd only if it is running. It will also disable the service sshd only if it is enabled. A newly joined host in the pool or an ejected host from the pool would keep the original status.
// Version: closed
func (host) DisableSSH2(session *Session, self HostRef) (err error) {
	method := "host.disable_ssh"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// AsyncDisableSSH2: Disable SSH access on the host. It will stop the service sshd only if it is running. It will also disable the service sshd only if it is enabled. A newly joined host in the pool or an ejected host from the pool would keep the original status.
// Version: closed
func (host) AsyncDisableSSH2(session *Session, self HostRef) (retval TaskRef, err error) {
	method := "Async.host.disable_ssh"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// EnableSSH: Enable SSH access on the host. It will start the service sshd only if it is not running. It will also enable the service sshd only if it is not enabled. A newly joined host in the pool or an ejected host from the pool would keep the original status.
// Version: closed
func (host) EnableSSH(session *Session, self HostRef) (err error) {
	method := "host.enable_ssh"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// AsyncEnableSSH: Enable SSH access on the host. It will start the service sshd only if it is not running. It will also enable the service sshd only if it is not enabled. A newly joined host in the pool or an ejected host from the pool would keep the original status.
// Version: closed
func (host) AsyncEnableSSH(session *Session, self HostRef) (retval TaskRef, err error) {
	method := "Async.host.enable_ssh"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// EnableSSH2: Enable SSH access on the host. It will start the service sshd only if it is not running. It will also enable the service sshd only if it is not enabled. A newly joined host in the pool or an ejected host from the pool would keep the original status.
// Version: closed
func (host) EnableSSH2(session *Session, self HostRef) (err error) {
	method := "host.enable_ssh"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// AsyncEnableSSH2: Enable SSH access on the host. It will start the service sshd only if it is not running. It will also enable the service sshd only if it is not enabled. A newly joined host in the pool or an ejected host from the pool would keep the original status.
// Version: closed
func (host) AsyncEnableSSH2(session *Session, self HostRef) (retval TaskRef, err error) {
	method := "Async.host.enable_ssh"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// EmergencyClearMandatoryGuidance: Clear the pending mandatory guidance on this host
// Version: closed
func (host) EmergencyClearMandatoryGuidance(session *Session) (err error) {
	method := "host.emergency_clear_mandatory_guidance"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg)
	return
}

// EmergencyClearMandatoryGuidance1: Clear the pending mandatory guidance on this host
// Version: closed
func (host) EmergencyClearMandatoryGuidance1(session *Session) (err error) {
	method := "host.emergency_clear_mandatory_guidance"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg)
	return
}

// ApplyRecommendedGuidances: apply all recommended guidances both on the host and on all HVM VMs on the host after updates are applied on the host
// Version: closed
func (host) ApplyRecommendedGuidances(session *Session, self HostRef) (err error) {
	method := "host.apply_recommended_guidances"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// AsyncApplyRecommendedGuidances: apply all recommended guidances both on the host and on all HVM VMs on the host after updates are applied on the host
// Version: closed
func (host) AsyncApplyRecommendedGuidances(session *Session, self HostRef) (retval TaskRef, err error) {
	method := "Async.host.apply_recommended_guidances"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// ApplyRecommendedGuidances2: apply all recommended guidances both on the host and on all HVM VMs on the host after updates are applied on the host
// Version: closed
func (host) ApplyRecommendedGuidances2(session *Session, self HostRef) (err error) {
	method := "host.apply_recommended_guidances"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// AsyncApplyRecommendedGuidances2: apply all recommended guidances both on the host and on all HVM VMs on the host after updates are applied on the host
// Version: closed
func (host) AsyncApplyRecommendedGuidances2(session *Session, self HostRef) (retval TaskRef, err error) {
	method := "Async.host.apply_recommended_guidances"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// SetHTTPSOnly: updates the host firewall to open or close port 80 depending on the value
// Version: closed
func (host) SetHTTPSOnly(session *Session, self HostRef, value bool) (err error) {
	method := "host.set_https_only"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	return
}

// AsyncSetHTTPSOnly: updates the host firewall to open or close port 80 depending on the value
// Version: closed
func (host) AsyncSetHTTPSOnly(session *Session, self HostRef, value bool) (retval TaskRef, err error) {
	method := "Async.host.set_https_only"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// SetHTTPSOnly3: updates the host firewall to open or close port 80 depending on the value
// Version: closed
func (host) SetHTTPSOnly3(session *Session, self HostRef, value bool) (err error) {
	method := "host.set_https_only"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	return
}

// AsyncSetHTTPSOnly3: updates the host firewall to open or close port 80 depending on the value
// Version: closed
func (host) AsyncSetHTTPSOnly3(session *Session, self HostRef, value bool) (retval TaskRef, err error) {
	method := "Async.host.set_https_only"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// RescanDrivers: Scan the host and update its driver information.
// Version: closed
func (host) RescanDrivers(session *Session, self HostRef) (err error) {
	method := "host.rescan_drivers"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// AsyncRescanDrivers: Scan the host and update its driver information.
// Version: closed
func (host) AsyncRescanDrivers(session *Session, self HostRef) (retval TaskRef, err error) {
	method := "Async.host.rescan_drivers"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// RescanDrivers2: Scan the host and update its driver information.
// Version: closed
func (host) RescanDrivers2(session *Session, self HostRef) (err error) {
	method := "host.rescan_drivers"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// AsyncRescanDrivers2: Scan the host and update its driver information.
// Version: closed
func (host) AsyncRescanDrivers2(session *Session, self HostRef) (retval TaskRef, err error) {
	method := "Async.host.rescan_drivers"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// ApplyUpdates: apply updates from current enabled repository on a host
// Version: 1.301.0
func (host) ApplyUpdates(session *Session, self HostRef, hash string) (retval [][]string, err error) {
	method := "host.apply_updates"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	hashArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "hash"), hash)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg, hashArg)
	if err != nil {
		return
	}
	retval, err = deserializeStringSetSet(method+" -> ", result)
	return
}

// AsyncApplyUpdates: apply updates from current enabled repository on a host
// Version: 1.301.0
func (host) AsyncApplyUpdates(session *Session, self HostRef, hash string) (retval TaskRef, err error) {
	method := "Async.host.apply_updates"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	hashArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "hash"), hash)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg, hashArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// ApplyUpdates3: apply updates from current enabled repository on a host
// Version: 1.301.0
func (host) ApplyUpdates3(session *Session, self HostRef, hash string) (retval [][]string, err error) {
	method := "host.apply_updates"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	hashArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "hash"), hash)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg, hashArg)
	if err != nil {
		return
	}
	retval, err = deserializeStringSetSet(method+" -> ", result)
	return
}

// AsyncApplyUpdates3: apply updates from current enabled repository on a host
// Version: 1.301.0
func (host) AsyncApplyUpdates3(session *Session, self HostRef, hash string) (retval TaskRef, err error) {
	method := "Async.host.apply_updates"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	hashArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "hash"), hash)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg, hashArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// EmergencyReenableTLSVerification: Reenable TLS verification for this host only, and only after it was emergency disabled
// Version: 1.298.0
func (host) EmergencyReenableTLSVerification(session *Session) (err error) {
	method := "host.emergency_reenable_tls_verification"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg)
	return
}

// EmergencyReenableTLSVerification1: Reenable TLS verification for this host only, and only after it was emergency disabled
// Version: 1.298.0
func (host) EmergencyReenableTLSVerification1(session *Session) (err error) {
	method := "host.emergency_reenable_tls_verification"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg)
	return
}

// EmergencyDisableTLSVerification: Disable TLS verification for this host only
// Version: 1.290.0
func (host) EmergencyDisableTLSVerification(session *Session) (err error) {
	method := "host.emergency_disable_tls_verification"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg)
	return
}

// EmergencyDisableTLSVerification1: Disable TLS verification for this host only
// Version: 1.290.0
func (host) EmergencyDisableTLSVerification1(session *Session) (err error) {
	method := "host.emergency_disable_tls_verification"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg)
	return
}

// SetNumaAffinityPolicy: Set VM placement NUMA affinity policy
// Version: closed
func (host) SetNumaAffinityPolicy(session *Session, self HostRef, value HostNumaAffinityPolicy) (err error) {
	method := "host.set_numa_affinity_policy"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeEnumHostNumaAffinityPolicy(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	return
}

// AsyncSetNumaAffinityPolicy: Set VM placement NUMA affinity policy
// Version: closed
func (host) AsyncSetNumaAffinityPolicy(session *Session, self HostRef, value HostNumaAffinityPolicy) (retval TaskRef, err error) {
	method := "Async.host.set_numa_affinity_policy"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeEnumHostNumaAffinityPolicy(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// SetNumaAffinityPolicy3: Set VM placement NUMA affinity policy
// Version: closed
func (host) SetNumaAffinityPolicy3(session *Session, self HostRef, value HostNumaAffinityPolicy) (err error) {
	method := "host.set_numa_affinity_policy"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeEnumHostNumaAffinityPolicy(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	return
}

// AsyncSetNumaAffinityPolicy3: Set VM placement NUMA affinity policy
// Version: closed
func (host) AsyncSetNumaAffinityPolicy3(session *Session, self HostRef, value HostNumaAffinityPolicy) (retval TaskRef, err error) {
	method := "Async.host.set_numa_affinity_policy"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeEnumHostNumaAffinityPolicy(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// GetSchedGran: Gets xen&apos;s sched-gran on a host
// Version: 1.271.0
func (host) GetSchedGran(session *Session, self HostRef) (retval HostSchedGran, err error) {
	method := "host.get_sched_gran"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeEnumHostSchedGran(method+" -> ", result)
	return
}

// AsyncGetSchedGran: Gets xen&apos;s sched-gran on a host
// Version: 1.271.0
func (host) AsyncGetSchedGran(session *Session, self HostRef) (retval TaskRef, err error) {
	method := "Async.host.get_sched_gran"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// GetSchedGran2: Gets xen&apos;s sched-gran on a host
// Version: 1.271.0
func (host) GetSchedGran2(session *Session, self HostRef) (retval HostSchedGran, err error) {
	method := "host.get_sched_gran"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeEnumHostSchedGran(method+" -> ", result)
	return
}

// AsyncGetSchedGran2: Gets xen&apos;s sched-gran on a host
// Version: 1.271.0
func (host) AsyncGetSchedGran2(session *Session, self HostRef) (retval TaskRef, err error) {
	method := "Async.host.get_sched_gran"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// SetSchedGran: Sets xen&apos;s sched-gran on a host. See: https://xenbits.xen.org/docs/unstable/misc/xen-command-line.html#sched-gran-x86
// Version: 1.271.0
func (host) SetSchedGran(session *Session, self HostRef, value HostSchedGran) (err error) {
	method := "host.set_sched_gran"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeEnumHostSchedGran(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	return
}

// AsyncSetSchedGran: Sets xen&apos;s sched-gran on a host. See: https://xenbits.xen.org/docs/unstable/misc/xen-command-line.html#sched-gran-x86
// Version: 1.271.0
func (host) AsyncSetSchedGran(session *Session, self HostRef, value HostSchedGran) (retval TaskRef, err error) {
	method := "Async.host.set_sched_gran"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeEnumHostSchedGran(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// SetSchedGran3: Sets xen&apos;s sched-gran on a host. See: https://xenbits.xen.org/docs/unstable/misc/xen-command-line.html#sched-gran-x86
// Version: 1.271.0
func (host) SetSchedGran3(session *Session, self HostRef, value HostSchedGran) (err error) {
	method := "host.set_sched_gran"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeEnumHostSchedGran(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	return
}

// AsyncSetSchedGran3: Sets xen&apos;s sched-gran on a host. See: https://xenbits.xen.org/docs/unstable/misc/xen-command-line.html#sched-gran-x86
// Version: 1.271.0
func (host) AsyncSetSchedGran3(session *Session, self HostRef, value HostSchedGran) (retval TaskRef, err error) {
	method := "Async.host.set_sched_gran"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeEnumHostSchedGran(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// SetUefiCertificates: Sets the UEFI certificates on a host
// Version: quebec
func (host) SetUefiCertificates(session *Session, host HostRef, value string) (err error) {
	method := "host.set_uefi_certificates"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	valueArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, hostArg, valueArg)
	return
}

// AsyncSetUefiCertificates: Sets the UEFI certificates on a host
// Version: quebec
func (host) AsyncSetUefiCertificates(session *Session, host HostRef, value string) (retval TaskRef, err error) {
	method := "Async.host.set_uefi_certificates"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	valueArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, hostArg, valueArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// SetUefiCertificates3: Sets the UEFI certificates on a host
// Version: quebec
func (host) SetUefiCertificates3(session *Session, host HostRef, value string) (err error) {
	method := "host.set_uefi_certificates"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	valueArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, hostArg, valueArg)
	return
}

// AsyncSetUefiCertificates3: Sets the UEFI certificates on a host
// Version: quebec
func (host) AsyncSetUefiCertificates3(session *Session, host HostRef, value string) (retval TaskRef, err error) {
	method := "Async.host.set_uefi_certificates"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	valueArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, hostArg, valueArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// SetMultipathing: Specifies whether multipathing is enabled
// Version: kolkata
func (host) SetMultipathing(session *Session, host HostRef, value bool) (err error) {
	method := "host.set_multipathing"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	valueArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, hostArg, valueArg)
	return
}

// AsyncSetMultipathing: Specifies whether multipathing is enabled
// Version: kolkata
func (host) AsyncSetMultipathing(session *Session, host HostRef, value bool) (retval TaskRef, err error) {
	method := "Async.host.set_multipathing"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	valueArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, hostArg, valueArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// SetMultipathing3: Specifies whether multipathing is enabled
// Version: kolkata
func (host) SetMultipathing3(session *Session, host HostRef, value bool) (err error) {
	method := "host.set_multipathing"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	valueArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, hostArg, valueArg)
	return
}

// AsyncSetMultipathing3: Specifies whether multipathing is enabled
// Version: kolkata
func (host) AsyncSetMultipathing3(session *Session, host HostRef, value bool) (retval TaskRef, err error) {
	method := "Async.host.set_multipathing"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	valueArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, hostArg, valueArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// SetIscsiIqn: Sets the initiator IQN for the host
// Version: kolkata
func (host) SetIscsiIqn(session *Session, host HostRef, value string) (err error) {
	method := "host.set_iscsi_iqn"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	valueArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, hostArg, valueArg)
	return
}

// AsyncSetIscsiIqn: Sets the initiator IQN for the host
// Version: kolkata
func (host) AsyncSetIscsiIqn(session *Session, host HostRef, value string) (retval TaskRef, err error) {
	method := "Async.host.set_iscsi_iqn"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	valueArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, hostArg, valueArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// SetIscsiIqn3: Sets the initiator IQN for the host
// Version: kolkata
func (host) SetIscsiIqn3(session *Session, host HostRef, value string) (err error) {
	method := "host.set_iscsi_iqn"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	valueArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, hostArg, valueArg)
	return
}

// AsyncSetIscsiIqn3: Sets the initiator IQN for the host
// Version: kolkata
func (host) AsyncSetIscsiIqn3(session *Session, host HostRef, value string) (retval TaskRef, err error) {
	method := "Async.host.set_iscsi_iqn"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	valueArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, hostArg, valueArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// SetSslLegacy: Enable/disable SSLv3 for interoperability with older server versions. When this is set to a different value, the host immediately restarts its SSL/TLS listening service; typically this takes less than a second but existing connections to it will be broken. API login sessions will remain valid.
// Version: dundee
func (host) SetSslLegacy(session *Session, self HostRef, value bool) (err error) {
	method := "host.set_ssl_legacy"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	return
}

// AsyncSetSslLegacy: Enable/disable SSLv3 for interoperability with older server versions. When this is set to a different value, the host immediately restarts its SSL/TLS listening service; typically this takes less than a second but existing connections to it will be broken. API login sessions will remain valid.
// Version: dundee
func (host) AsyncSetSslLegacy(session *Session, self HostRef, value bool) (retval TaskRef, err error) {
	method := "Async.host.set_ssl_legacy"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// SetSslLegacy3: Enable/disable SSLv3 for interoperability with older server versions. When this is set to a different value, the host immediately restarts its SSL/TLS listening service; typically this takes less than a second but existing connections to it will be broken. API login sessions will remain valid.
// Version: dundee
func (host) SetSslLegacy3(session *Session, self HostRef, value bool) (err error) {
	method := "host.set_ssl_legacy"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	return
}

// AsyncSetSslLegacy3: Enable/disable SSLv3 for interoperability with older server versions. When this is set to a different value, the host immediately restarts its SSL/TLS listening service; typically this takes less than a second but existing connections to it will be broken. API login sessions will remain valid.
// Version: dundee
func (host) AsyncSetSslLegacy3(session *Session, self HostRef, value bool) (retval TaskRef, err error) {
	method := "Async.host.set_ssl_legacy"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// DisableDisplay: Disable console output to the physical display device next time this host boots
// Version: cream
func (host) DisableDisplay(session *Session, host HostRef) (retval HostDisplay, err error) {
	method := "host.disable_display"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, hostArg)
	if err != nil {
		return
	}
	retval, err = deserializeEnumHostDisplay(method+" -> ", result)
	return
}

// AsyncDisableDisplay: Disable console output to the physical display device next time this host boots
// Version: cream
func (host) AsyncDisableDisplay(session *Session, host HostRef) (retval TaskRef, err error) {
	method := "Async.host.disable_display"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, hostArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// DisableDisplay2: Disable console output to the physical display device next time this host boots
// Version: cream
func (host) DisableDisplay2(session *Session, host HostRef) (retval HostDisplay, err error) {
	method := "host.disable_display"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, hostArg)
	if err != nil {
		return
	}
	retval, err = deserializeEnumHostDisplay(method+" -> ", result)
	return
}

// AsyncDisableDisplay2: Disable console output to the physical display device next time this host boots
// Version: cream
func (host) AsyncDisableDisplay2(session *Session, host HostRef) (retval TaskRef, err error) {
	method := "Async.host.disable_display"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, hostArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// EnableDisplay: Enable console output to the physical display device next time this host boots
// Version: cream
func (host) EnableDisplay(session *Session, host HostRef) (retval HostDisplay, err error) {
	method := "host.enable_display"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, hostArg)
	if err != nil {
		return
	}
	retval, err = deserializeEnumHostDisplay(method+" -> ", result)
	return
}

// AsyncEnableDisplay: Enable console output to the physical display device next time this host boots
// Version: cream
func (host) AsyncEnableDisplay(session *Session, host HostRef) (retval TaskRef, err error) {
	method := "Async.host.enable_display"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, hostArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// EnableDisplay2: Enable console output to the physical display device next time this host boots
// Version: cream
func (host) EnableDisplay2(session *Session, host HostRef) (retval HostDisplay, err error) {
	method := "host.enable_display"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, hostArg)
	if err != nil {
		return
	}
	retval, err = deserializeEnumHostDisplay(method+" -> ", result)
	return
}

// AsyncEnableDisplay2: Enable console output to the physical display device next time this host boots
// Version: cream
func (host) AsyncEnableDisplay2(session *Session, host HostRef) (retval TaskRef, err error) {
	method := "Async.host.enable_display"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, hostArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// DeclareDead: Declare that a host is dead. This is a dangerous operation, and should only be called if the administrator is absolutely sure the host is definitely dead
// Version: clearwater
func (host) DeclareDead(session *Session, host HostRef) (err error) {
	method := "host.declare_dead"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, hostArg)
	return
}

// AsyncDeclareDead: Declare that a host is dead. This is a dangerous operation, and should only be called if the administrator is absolutely sure the host is definitely dead
// Version: clearwater
func (host) AsyncDeclareDead(session *Session, host HostRef) (retval TaskRef, err error) {
	method := "Async.host.declare_dead"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, hostArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// DeclareDead2: Declare that a host is dead. This is a dangerous operation, and should only be called if the administrator is absolutely sure the host is definitely dead
// Version: clearwater
func (host) DeclareDead2(session *Session, host HostRef) (err error) {
	method := "host.declare_dead"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, hostArg)
	return
}

// AsyncDeclareDead2: Declare that a host is dead. This is a dangerous operation, and should only be called if the administrator is absolutely sure the host is definitely dead
// Version: clearwater
func (host) AsyncDeclareDead2(session *Session, host HostRef) (retval TaskRef, err error) {
	method := "Async.host.declare_dead"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, hostArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// MigrateReceive: Prepare to receive a VM, returning a token which can be passed to VM.migrate.
// Version: tampa
func (host) MigrateReceive(session *Session, host HostRef, network NetworkRef, options map[string]string) (retval map[string]string, err error) {
	method := "host.migrate_receive"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	networkArg, err := serializeNetworkRef(fmt.Sprintf("%s(%s)", method, "network"), network)
	if err != nil {
		return
	}
	optionsArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "options"), options)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, hostArg, networkArg, optionsArg)
	if err != nil {
		return
	}
	retval, err = deserializeStringToStringMap(method+" -> ", result)
	return
}

// AsyncMigrateReceive: Prepare to receive a VM, returning a token which can be passed to VM.migrate.
// Version: tampa
func (host) AsyncMigrateReceive(session *Session, host HostRef, network NetworkRef, options map[string]string) (retval TaskRef, err error) {
	method := "Async.host.migrate_receive"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	networkArg, err := serializeNetworkRef(fmt.Sprintf("%s(%s)", method, "network"), network)
	if err != nil {
		return
	}
	optionsArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "options"), options)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, hostArg, networkArg, optionsArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// MigrateReceive4: Prepare to receive a VM, returning a token which can be passed to VM.migrate.
// Version: tampa
func (host) MigrateReceive4(session *Session, host HostRef, network NetworkRef, options map[string]string) (retval map[string]string, err error) {
	method := "host.migrate_receive"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	networkArg, err := serializeNetworkRef(fmt.Sprintf("%s(%s)", method, "network"), network)
	if err != nil {
		return
	}
	optionsArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "options"), options)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, hostArg, networkArg, optionsArg)
	if err != nil {
		return
	}
	retval, err = deserializeStringToStringMap(method+" -> ", result)
	return
}

// AsyncMigrateReceive4: Prepare to receive a VM, returning a token which can be passed to VM.migrate.
// Version: tampa
func (host) AsyncMigrateReceive4(session *Session, host HostRef, network NetworkRef, options map[string]string) (retval TaskRef, err error) {
	method := "Async.host.migrate_receive"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	networkArg, err := serializeNetworkRef(fmt.Sprintf("%s(%s)", method, "network"), network)
	if err != nil {
		return
	}
	optionsArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "options"), options)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, hostArg, networkArg, optionsArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// DisableLocalStorageCaching: Disable the use of a local SR for caching purposes
// Version: cowley
func (host) DisableLocalStorageCaching(session *Session, host HostRef) (err error) {
	method := "host.disable_local_storage_caching"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, hostArg)
	return
}

// DisableLocalStorageCaching2: Disable the use of a local SR for caching purposes
// Version: cowley
func (host) DisableLocalStorageCaching2(session *Session, host HostRef) (err error) {
	method := "host.disable_local_storage_caching"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, hostArg)
	return
}

// EnableLocalStorageCaching: Enable the use of a local SR for caching purposes
// Version: cowley
func (host) EnableLocalStorageCaching(session *Session, host HostRef, sr SRRef) (err error) {
	method := "host.enable_local_storage_caching"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	srArg, err := serializeSRRef(fmt.Sprintf("%s(%s)", method, "sr"), sr)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, hostArg, srArg)
	return
}

// EnableLocalStorageCaching3: Enable the use of a local SR for caching purposes
// Version: cowley
func (host) EnableLocalStorageCaching3(session *Session, host HostRef, sr SRRef) (err error) {
	method := "host.enable_local_storage_caching"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	srArg, err := serializeSRRef(fmt.Sprintf("%s(%s)", method, "sr"), sr)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, hostArg, srArg)
	return
}

// ResetCPUFeatures: Remove the feature mask, such that after a reboot all features of the CPU are enabled.
// Version: midnight-ride
func (host) ResetCPUFeatures(session *Session, host HostRef) (err error) {
	method := "host.reset_cpu_features"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, hostArg)
	return
}

// ResetCPUFeatures2: Remove the feature mask, such that after a reboot all features of the CPU are enabled.
// Version: midnight-ride
func (host) ResetCPUFeatures2(session *Session, host HostRef) (err error) {
	method := "host.reset_cpu_features"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, hostArg)
	return
}

// SetCPUFeatures: Set the CPU features to be used after a reboot, if the given features string is valid.
// Version: midnight-ride
func (host) SetCPUFeatures(session *Session, host HostRef, features string) (err error) {
	method := "host.set_cpu_features"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	featuresArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "features"), features)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, hostArg, featuresArg)
	return
}

// SetCPUFeatures3: Set the CPU features to be used after a reboot, if the given features string is valid.
// Version: midnight-ride
func (host) SetCPUFeatures3(session *Session, host HostRef, features string) (err error) {
	method := "host.set_cpu_features"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	featuresArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "features"), features)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, hostArg, featuresArg)
	return
}

// SetPowerOnMode: Set the power-on-mode, host, user and password
// Version: cowley
func (host) SetPowerOnMode(session *Session, self HostRef, powerOnMode string, powerOnConfig map[string]string) (err error) {
	method := "host.set_power_on_mode"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	powerOnModeArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "power_on_mode"), powerOnMode)
	if err != nil {
		return
	}
	powerOnConfigArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "power_on_config"), powerOnConfig)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, powerOnModeArg, powerOnConfigArg)
	return
}

// AsyncSetPowerOnMode: Set the power-on-mode, host, user and password
// Version: cowley
func (host) AsyncSetPowerOnMode(session *Session, self HostRef, powerOnMode string, powerOnConfig map[string]string) (retval TaskRef, err error) {
	method := "Async.host.set_power_on_mode"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	powerOnModeArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "power_on_mode"), powerOnMode)
	if err != nil {
		return
	}
	powerOnConfigArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "power_on_config"), powerOnConfig)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg, powerOnModeArg, powerOnConfigArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// SetPowerOnMode4: Set the power-on-mode, host, user and password
// Version: cowley
func (host) SetPowerOnMode4(session *Session, self HostRef, powerOnMode string, powerOnConfig map[string]string) (err error) {
	method := "host.set_power_on_mode"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	powerOnModeArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "power_on_mode"), powerOnMode)
	if err != nil {
		return
	}
	powerOnConfigArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "power_on_config"), powerOnConfig)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, powerOnModeArg, powerOnConfigArg)
	return
}

// AsyncSetPowerOnMode4: Set the power-on-mode, host, user and password
// Version: cowley
func (host) AsyncSetPowerOnMode4(session *Session, self HostRef, powerOnMode string, powerOnConfig map[string]string) (retval TaskRef, err error) {
	method := "Async.host.set_power_on_mode"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	powerOnModeArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "power_on_mode"), powerOnMode)
	if err != nil {
		return
	}
	powerOnConfigArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "power_on_config"), powerOnConfig)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg, powerOnModeArg, powerOnConfigArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// RefreshPackInfo: Refresh the list of installed Supplemental Packs.
// Version: midnight-ride
func (host) RefreshPackInfo(session *Session, host HostRef) (err error) {
	method := "host.refresh_pack_info"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, hostArg)
	return
}

// AsyncRefreshPackInfo: Refresh the list of installed Supplemental Packs.
// Version: midnight-ride
func (host) AsyncRefreshPackInfo(session *Session, host HostRef) (retval TaskRef, err error) {
	method := "Async.host.refresh_pack_info"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, hostArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// RefreshPackInfo2: Refresh the list of installed Supplemental Packs.
// Version: midnight-ride
func (host) RefreshPackInfo2(session *Session, host HostRef) (err error) {
	method := "host.refresh_pack_info"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, hostArg)
	return
}

// AsyncRefreshPackInfo2: Refresh the list of installed Supplemental Packs.
// Version: midnight-ride
func (host) AsyncRefreshPackInfo2(session *Session, host HostRef) (retval TaskRef, err error) {
	method := "Async.host.refresh_pack_info"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, hostArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// ApplyEdition: Change to another edition, or reactivate the current edition after a license has expired. This may be subject to the successful checkout of an appropriate license.
// Version: clearwater
func (host) ApplyEdition(session *Session, host HostRef, edition string, force bool) (err error) {
	method := "host.apply_edition"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	editionArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "edition"), edition)
	if err != nil {
		return
	}
	forceArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "force"), force)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, hostArg, editionArg, forceArg)
	return
}

// ApplyEdition4: Change to another edition, or reactivate the current edition after a license has expired. This may be subject to the successful checkout of an appropriate license.
// Version: clearwater
func (host) ApplyEdition4(session *Session, host HostRef, edition string, force bool) (err error) {
	method := "host.apply_edition"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	editionArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "edition"), edition)
	if err != nil {
		return
	}
	forceArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "force"), force)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, hostArg, editionArg, forceArg)
	return
}

// ApplyEdition3: Change to another edition, or reactivate the current edition after a license has expired. This may be subject to the successful checkout of an appropriate license.
// Version: midnight-ride
func (host) ApplyEdition3(session *Session, host HostRef, edition string) (err error) {
	method := "host.apply_edition"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	editionArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "edition"), edition)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, hostArg, editionArg)
	return
}

// ResetServerCertificate: Delete the current TLS server certificate and replace by a new, self-signed one. This should only be used with extreme care.
// Version: 1.290.0
func (host) ResetServerCertificate(session *Session, host HostRef) (err error) {
	method := "host.reset_server_certificate"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, hostArg)
	return
}

// AsyncResetServerCertificate: Delete the current TLS server certificate and replace by a new, self-signed one. This should only be used with extreme care.
// Version: 1.290.0
func (host) AsyncResetServerCertificate(session *Session, host HostRef) (retval TaskRef, err error) {
	method := "Async.host.reset_server_certificate"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, hostArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// ResetServerCertificate2: Delete the current TLS server certificate and replace by a new, self-signed one. This should only be used with extreme care.
// Version: 1.290.0
func (host) ResetServerCertificate2(session *Session, host HostRef) (err error) {
	method := "host.reset_server_certificate"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, hostArg)
	return
}

// AsyncResetServerCertificate2: Delete the current TLS server certificate and replace by a new, self-signed one. This should only be used with extreme care.
// Version: 1.290.0
func (host) AsyncResetServerCertificate2(session *Session, host HostRef) (retval TaskRef, err error) {
	method := "Async.host.reset_server_certificate"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, hostArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// EmergencyResetServerCertificate: Delete the current TLS server certificate and replace by a new, self-signed one. This should only be used with extreme care.
// Version: stockholm
func (host) EmergencyResetServerCertificate(session *Session) (err error) {
	method := "host.emergency_reset_server_certificate"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg)
	return
}

// EmergencyResetServerCertificate1: Delete the current TLS server certificate and replace by a new, self-signed one. This should only be used with extreme care.
// Version: stockholm
func (host) EmergencyResetServerCertificate1(session *Session) (err error) {
	method := "host.emergency_reset_server_certificate"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg)
	return
}

// InstallServerCertificate: Install the TLS server certificate.
// Version: stockholm
func (host) InstallServerCertificate(session *Session, host HostRef, certificate string, privateKey string, certificateChain string) (err error) {
	method := "host.install_server_certificate"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	certificateArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "certificate"), certificate)
	if err != nil {
		return
	}
	privateKeyArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "private_key"), privateKey)
	if err != nil {
		return
	}
	certificateChainArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "certificate_chain"), certificateChain)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, hostArg, certificateArg, privateKeyArg, certificateChainArg)
	return
}

// AsyncInstallServerCertificate: Install the TLS server certificate.
// Version: stockholm
func (host) AsyncInstallServerCertificate(session *Session, host HostRef, certificate string, privateKey string, certificateChain string) (retval TaskRef, err error) {
	method := "Async.host.install_server_certificate"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	certificateArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "certificate"), certificate)
	if err != nil {
		return
	}
	privateKeyArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "private_key"), privateKey)
	if err != nil {
		return
	}
	certificateChainArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "certificate_chain"), certificateChain)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, hostArg, certificateArg, privateKeyArg, certificateChainArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// InstallServerCertificate5: Install the TLS server certificate.
// Version: stockholm
func (host) InstallServerCertificate5(session *Session, host HostRef, certificate string, privateKey string, certificateChain string) (err error) {
	method := "host.install_server_certificate"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	certificateArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "certificate"), certificate)
	if err != nil {
		return
	}
	privateKeyArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "private_key"), privateKey)
	if err != nil {
		return
	}
	certificateChainArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "certificate_chain"), certificateChain)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, hostArg, certificateArg, privateKeyArg, certificateChainArg)
	return
}

// AsyncInstallServerCertificate5: Install the TLS server certificate.
// Version: stockholm
func (host) AsyncInstallServerCertificate5(session *Session, host HostRef, certificate string, privateKey string, certificateChain string) (retval TaskRef, err error) {
	method := "Async.host.install_server_certificate"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	certificateArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "certificate"), certificate)
	if err != nil {
		return
	}
	privateKeyArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "private_key"), privateKey)
	if err != nil {
		return
	}
	certificateChainArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "certificate_chain"), certificateChain)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, hostArg, certificateArg, privateKeyArg, certificateChainArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// RefreshServerCertificate: Replace the internal self-signed host certficate with a new one.
// Version: 1.307.0
func (host) RefreshServerCertificate(session *Session, host HostRef) (err error) {
	method := "host.refresh_server_certificate"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, hostArg)
	return
}

// AsyncRefreshServerCertificate: Replace the internal self-signed host certficate with a new one.
// Version: 1.307.0
func (host) AsyncRefreshServerCertificate(session *Session, host HostRef) (retval TaskRef, err error) {
	method := "Async.host.refresh_server_certificate"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, hostArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// RefreshServerCertificate2: Replace the internal self-signed host certficate with a new one.
// Version: 1.307.0
func (host) RefreshServerCertificate2(session *Session, host HostRef) (err error) {
	method := "host.refresh_server_certificate"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, hostArg)
	return
}

// AsyncRefreshServerCertificate2: Replace the internal self-signed host certficate with a new one.
// Version: 1.307.0
func (host) AsyncRefreshServerCertificate2(session *Session, host HostRef) (retval TaskRef, err error) {
	method := "Async.host.refresh_server_certificate"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, hostArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// GetServerCertificate: Get the installed server public TLS certificate.
// Version: george
func (host) GetServerCertificate(session *Session, host HostRef) (retval string, err error) {
	method := "host.get_server_certificate"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, hostArg)
	if err != nil {
		return
	}
	retval, err = deserializeString(method+" -> ", result)
	return
}

// AsyncGetServerCertificate: Get the installed server public TLS certificate.
// Version: george
func (host) AsyncGetServerCertificate(session *Session, host HostRef) (retval TaskRef, err error) {
	method := "Async.host.get_server_certificate"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, hostArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// GetServerCertificate2: Get the installed server public TLS certificate.
// Version: george
func (host) GetServerCertificate2(session *Session, host HostRef) (retval string, err error) {
	method := "host.get_server_certificate"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, hostArg)
	if err != nil {
		return
	}
	retval, err = deserializeString(method+" -> ", result)
	return
}

// AsyncGetServerCertificate2: Get the installed server public TLS certificate.
// Version: george
func (host) AsyncGetServerCertificate2(session *Session, host HostRef) (retval TaskRef, err error) {
	method := "Async.host.get_server_certificate"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, hostArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// RetrieveWlbEvacuateRecommendations: Retrieves recommended host migrations to perform when evacuating the host from the wlb server. If a VM cannot be migrated from the host the reason is listed instead of a recommendation.
// Version: george
func (host) RetrieveWlbEvacuateRecommendations(session *Session, self HostRef) (retval map[VMRef][]string, err error) {
	method := "host.retrieve_wlb_evacuate_recommendations"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeVMRefToStringSetMap(method+" -> ", result)
	return
}

// AsyncRetrieveWlbEvacuateRecommendations: Retrieves recommended host migrations to perform when evacuating the host from the wlb server. If a VM cannot be migrated from the host the reason is listed instead of a recommendation.
// Version: george
func (host) AsyncRetrieveWlbEvacuateRecommendations(session *Session, self HostRef) (retval TaskRef, err error) {
	method := "Async.host.retrieve_wlb_evacuate_recommendations"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// RetrieveWlbEvacuateRecommendations2: Retrieves recommended host migrations to perform when evacuating the host from the wlb server. If a VM cannot be migrated from the host the reason is listed instead of a recommendation.
// Version: george
func (host) RetrieveWlbEvacuateRecommendations2(session *Session, self HostRef) (retval map[VMRef][]string, err error) {
	method := "host.retrieve_wlb_evacuate_recommendations"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeVMRefToStringSetMap(method+" -> ", result)
	return
}

// AsyncRetrieveWlbEvacuateRecommendations2: Retrieves recommended host migrations to perform when evacuating the host from the wlb server. If a VM cannot be migrated from the host the reason is listed instead of a recommendation.
// Version: george
func (host) AsyncRetrieveWlbEvacuateRecommendations2(session *Session, self HostRef) (retval TaskRef, err error) {
	method := "Async.host.retrieve_wlb_evacuate_recommendations"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// ExternalAuthSetLdaps: This call enables or disables LDAPS for external authentication on the host
// Version: 26.16.0
func (host) ExternalAuthSetLdaps(session *Session, host HostRef, ldaps bool, force bool) (err error) {
	method := "host.external_auth_set_ldaps"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	ldapsArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "ldaps"), ldaps)
	if err != nil {
		return
	}
	forceArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "force"), force)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, hostArg, ldapsArg, forceArg)
	return
}

// ExternalAuthSetLdaps4: This call enables or disables LDAPS for external authentication on the host
// Version: 26.16.0
func (host) ExternalAuthSetLdaps4(session *Session, host HostRef, ldaps bool, force bool) (err error) {
	method := "host.external_auth_set_ldaps"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	ldapsArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "ldaps"), ldaps)
	if err != nil {
		return
	}
	forceArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "force"), force)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, hostArg, ldapsArg, forceArg)
	return
}

// DisableExternalAuth: This call disables external authentication on the local host
// Version: 26.2.0-next
func (host) DisableExternalAuth(session *Session, host HostRef, config map[string]string, force bool) (err error) {
	method := "host.disable_external_auth"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	configArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "config"), config)
	if err != nil {
		return
	}
	forceArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "force"), force)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, hostArg, configArg, forceArg)
	return
}

// DisableExternalAuth4: This call disables external authentication on the local host
// Version: 26.2.0-next
func (host) DisableExternalAuth4(session *Session, host HostRef, config map[string]string, force bool) (err error) {
	method := "host.disable_external_auth"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	configArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "config"), config)
	if err != nil {
		return
	}
	forceArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "force"), force)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, hostArg, configArg, forceArg)
	return
}

// DisableExternalAuth3: This call disables external authentication on the local host
// Version: george
func (host) DisableExternalAuth3(session *Session, host HostRef, config map[string]string) (err error) {
	method := "host.disable_external_auth"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	configArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "config"), config)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, hostArg, configArg)
	return
}

// EnableExternalAuth: This call enables external authentication on a host
// Version: george
func (host) EnableExternalAuth(session *Session, host HostRef, config map[string]string, serviceName string, authType string) (err error) {
	method := "host.enable_external_auth"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	configArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "config"), config)
	if err != nil {
		return
	}
	serviceNameArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "service_name"), serviceName)
	if err != nil {
		return
	}
	authTypeArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "auth_type"), authType)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, hostArg, configArg, serviceNameArg, authTypeArg)
	return
}

// EnableExternalAuth5: This call enables external authentication on a host
// Version: george
func (host) EnableExternalAuth5(session *Session, host HostRef, config map[string]string, serviceName string, authType string) (err error) {
	method := "host.enable_external_auth"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	configArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "config"), config)
	if err != nil {
		return
	}
	serviceNameArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "service_name"), serviceName)
	if err != nil {
		return
	}
	authTypeArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "auth_type"), authType)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, hostArg, configArg, serviceNameArg, authTypeArg)
	return
}

// GetServerLocaltime: This call queries the host&apos;s clock for the current time in the host&apos;s local timezone
// Version: cowley
func (host) GetServerLocaltime(session *Session, host HostRef) (retval time.Time, err error) {
	method := "host.get_server_localtime"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, hostArg)
	if err != nil {
		return
	}
	retval, err = deserializeTime(method+" -> ", result)
	return
}

// GetServerLocaltime2: This call queries the host&apos;s clock for the current time in the host&apos;s local timezone
// Version: cowley
func (host) GetServerLocaltime2(session *Session, host HostRef) (retval time.Time, err error) {
	method := "host.get_server_localtime"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, hostArg)
	if err != nil {
		return
	}
	retval, err = deserializeTime(method+" -> ", result)
	return
}

// GetServertime: This call queries the host&apos;s clock for the current time
// Version: orlando
func (host) GetServertime(session *Session, host HostRef) (retval time.Time, err error) {
	method := "host.get_servertime"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, hostArg)
	if err != nil {
		return
	}
	retval, err = deserializeTime(method+" -> ", result)
	return
}

// GetServertime2: This call queries the host&apos;s clock for the current time
// Version: orlando
func (host) GetServertime2(session *Session, host HostRef) (retval time.Time, err error) {
	method := "host.get_servertime"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, hostArg)
	if err != nil {
		return
	}
	retval, err = deserializeTime(method+" -> ", result)
	return
}

// CallExtension: Call an API extension on this host
// Version: ely
func (host) CallExtension(session *Session, host HostRef, call string) (retval string, err error) {
	method := "host.call_extension"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	callArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "call"), call)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, hostArg, callArg)
	if err != nil {
		return
	}
	retval, err = deserializeString(method+" -> ", result)
	return
}

// CallExtension3: Call an API extension on this host
// Version: ely
func (host) CallExtension3(session *Session, host HostRef, call string) (retval string, err error) {
	method := "host.call_extension"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	callArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "call"), call)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, hostArg, callArg)
	if err != nil {
		return
	}
	retval, err = deserializeString(method+" -> ", result)
	return
}

// HasExtension: Return true if the extension is available on the host
// Version: ely
func (host) HasExtension(session *Session, host HostRef, name string) (retval bool, err error) {
	method := "host.has_extension"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	nameArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "name"), name)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, hostArg, nameArg)
	if err != nil {
		return
	}
	retval, err = deserializeBool(method+" -> ", result)
	return
}

// AsyncHasExtension: Return true if the extension is available on the host
// Version: ely
func (host) AsyncHasExtension(session *Session, host HostRef, name string) (retval TaskRef, err error) {
	method := "Async.host.has_extension"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	nameArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "name"), name)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, hostArg, nameArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// HasExtension3: Return true if the extension is available on the host
// Version: ely
func (host) HasExtension3(session *Session, host HostRef, name string) (retval bool, err error) {
	method := "host.has_extension"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	nameArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "name"), name)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, hostArg, nameArg)
	if err != nil {
		return
	}
	retval, err = deserializeBool(method+" -> ", result)
	return
}

// AsyncHasExtension3: Return true if the extension is available on the host
// Version: ely
func (host) AsyncHasExtension3(session *Session, host HostRef, name string) (retval TaskRef, err error) {
	method := "Async.host.has_extension"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	nameArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "name"), name)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, hostArg, nameArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// CallPlugin: Call an API plugin on this host
// Version: orlando
func (host) CallPlugin(session *Session, host HostRef, plugin string, fn string, args map[string]string) (retval string, err error) {
	method := "host.call_plugin"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	pluginArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "plugin"), plugin)
	if err != nil {
		return
	}
	fnArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "fn"), fn)
	if err != nil {
		return
	}
	argsArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "args"), args)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, hostArg, pluginArg, fnArg, argsArg)
	if err != nil {
		return
	}
	retval, err = deserializeString(method+" -> ", result)
	return
}

// AsyncCallPlugin: Call an API plugin on this host
// Version: orlando
func (host) AsyncCallPlugin(session *Session, host HostRef, plugin string, fn string, args map[string]string) (retval TaskRef, err error) {
	method := "Async.host.call_plugin"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	pluginArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "plugin"), plugin)
	if err != nil {
		return
	}
	fnArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "fn"), fn)
	if err != nil {
		return
	}
	argsArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "args"), args)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, hostArg, pluginArg, fnArg, argsArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// CallPlugin5: Call an API plugin on this host
// Version: orlando
func (host) CallPlugin5(session *Session, host HostRef, plugin string, fn string, args map[string]string) (retval string, err error) {
	method := "host.call_plugin"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	pluginArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "plugin"), plugin)
	if err != nil {
		return
	}
	fnArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "fn"), fn)
	if err != nil {
		return
	}
	argsArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "args"), args)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, hostArg, pluginArg, fnArg, argsArg)
	if err != nil {
		return
	}
	retval, err = deserializeString(method+" -> ", result)
	return
}

// AsyncCallPlugin5: Call an API plugin on this host
// Version: orlando
func (host) AsyncCallPlugin5(session *Session, host HostRef, plugin string, fn string, args map[string]string) (retval TaskRef, err error) {
	method := "Async.host.call_plugin"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	pluginArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "plugin"), plugin)
	if err != nil {
		return
	}
	fnArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "fn"), fn)
	if err != nil {
		return
	}
	argsArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "args"), args)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, hostArg, pluginArg, fnArg, argsArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// CreateNewBlob: Create a placeholder for a named binary blob of data that is associated with this host
// Version: tampa
func (host) CreateNewBlob(session *Session, host HostRef, name string, mimeType string, public bool) (retval BlobRef, err error) {
	method := "host.create_new_blob"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	nameArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "name"), name)
	if err != nil {
		return
	}
	mimeTypeArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "mime_type"), mimeType)
	if err != nil {
		return
	}
	publicArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "public"), public)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, hostArg, nameArg, mimeTypeArg, publicArg)
	if err != nil {
		return
	}
	retval, err = deserializeBlobRef(method+" -> ", result)
	return
}

// AsyncCreateNewBlob: Create a placeholder for a named binary blob of data that is associated with this host
// Version: tampa
func (host) AsyncCreateNewBlob(session *Session, host HostRef, name string, mimeType string, public bool) (retval TaskRef, err error) {
	method := "Async.host.create_new_blob"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	nameArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "name"), name)
	if err != nil {
		return
	}
	mimeTypeArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "mime_type"), mimeType)
	if err != nil {
		return
	}
	publicArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "public"), public)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, hostArg, nameArg, mimeTypeArg, publicArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// CreateNewBlob5: Create a placeholder for a named binary blob of data that is associated with this host
// Version: tampa
func (host) CreateNewBlob5(session *Session, host HostRef, name string, mimeType string, public bool) (retval BlobRef, err error) {
	method := "host.create_new_blob"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	nameArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "name"), name)
	if err != nil {
		return
	}
	mimeTypeArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "mime_type"), mimeType)
	if err != nil {
		return
	}
	publicArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "public"), public)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, hostArg, nameArg, mimeTypeArg, publicArg)
	if err != nil {
		return
	}
	retval, err = deserializeBlobRef(method+" -> ", result)
	return
}

// AsyncCreateNewBlob5: Create a placeholder for a named binary blob of data that is associated with this host
// Version: tampa
func (host) AsyncCreateNewBlob5(session *Session, host HostRef, name string, mimeType string, public bool) (retval TaskRef, err error) {
	method := "Async.host.create_new_blob"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	nameArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "name"), name)
	if err != nil {
		return
	}
	mimeTypeArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "mime_type"), mimeType)
	if err != nil {
		return
	}
	publicArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "public"), public)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, hostArg, nameArg, mimeTypeArg, publicArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// CreateNewBlob4: Create a placeholder for a named binary blob of data that is associated with this host
// Version: orlando
func (host) CreateNewBlob4(session *Session, host HostRef, name string, mimeType string) (retval BlobRef, err error) {
	method := "host.create_new_blob"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	nameArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "name"), name)
	if err != nil {
		return
	}
	mimeTypeArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "mime_type"), mimeType)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, hostArg, nameArg, mimeTypeArg)
	if err != nil {
		return
	}
	retval, err = deserializeBlobRef(method+" -> ", result)
	return
}

// AsyncCreateNewBlob4: Create a placeholder for a named binary blob of data that is associated with this host
// Version: orlando
func (host) AsyncCreateNewBlob4(session *Session, host HostRef, name string, mimeType string) (retval TaskRef, err error) {
	method := "Async.host.create_new_blob"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	nameArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "name"), name)
	if err != nil {
		return
	}
	mimeTypeArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "mime_type"), mimeType)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, hostArg, nameArg, mimeTypeArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// BackupRrds: This causes the RRDs to be backed up to the master
// Version: orlando
func (host) BackupRrds(session *Session, host HostRef, delay float64) (err error) {
	method := "host.backup_rrds"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	delayArg, err := serializeFloat(fmt.Sprintf("%s(%s)", method, "delay"), delay)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, hostArg, delayArg)
	return
}

// BackupRrds3: This causes the RRDs to be backed up to the master
// Version: orlando
func (host) BackupRrds3(session *Session, host HostRef, delay float64) (err error) {
	method := "host.backup_rrds"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	delayArg, err := serializeFloat(fmt.Sprintf("%s(%s)", method, "delay"), delay)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, hostArg, delayArg)
	return
}

// SyncData: This causes the synchronisation of the non-database data (messages, RRDs and so on) stored on the master to be synchronised with the host
// Version: orlando
func (host) SyncData(session *Session, host HostRef) (err error) {
	method := "host.sync_data"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, hostArg)
	return
}

// SyncData2: This causes the synchronisation of the non-database data (messages, RRDs and so on) stored on the master to be synchronised with the host
// Version: orlando
func (host) SyncData2(session *Session, host HostRef) (err error) {
	method := "host.sync_data"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, hostArg)
	return
}

// ComputeMemoryOverhead: Computes the virtualization memory overhead of a host.
// Version: midnight-ride
func (host) ComputeMemoryOverhead(session *Session, host HostRef) (retval int, err error) {
	method := "host.compute_memory_overhead"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, hostArg)
	if err != nil {
		return
	}
	retval, err = deserializeInt(method+" -> ", result)
	return
}

// AsyncComputeMemoryOverhead: Computes the virtualization memory overhead of a host.
// Version: midnight-ride
func (host) AsyncComputeMemoryOverhead(session *Session, host HostRef) (retval TaskRef, err error) {
	method := "Async.host.compute_memory_overhead"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, hostArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// ComputeMemoryOverhead2: Computes the virtualization memory overhead of a host.
// Version: midnight-ride
func (host) ComputeMemoryOverhead2(session *Session, host HostRef) (retval int, err error) {
	method := "host.compute_memory_overhead"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, hostArg)
	if err != nil {
		return
	}
	retval, err = deserializeInt(method+" -> ", result)
	return
}

// AsyncComputeMemoryOverhead2: Computes the virtualization memory overhead of a host.
// Version: midnight-ride
func (host) AsyncComputeMemoryOverhead2(session *Session, host HostRef) (retval TaskRef, err error) {
	method := "Async.host.compute_memory_overhead"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, hostArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// ComputeFreeMemory: Computes the amount of free memory on the host.
// Version: orlando
func (host) ComputeFreeMemory(session *Session, host HostRef) (retval int, err error) {
	method := "host.compute_free_memory"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, hostArg)
	if err != nil {
		return
	}
	retval, err = deserializeInt(method+" -> ", result)
	return
}

// AsyncComputeFreeMemory: Computes the amount of free memory on the host.
// Version: orlando
func (host) AsyncComputeFreeMemory(session *Session, host HostRef) (retval TaskRef, err error) {
	method := "Async.host.compute_free_memory"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, hostArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// ComputeFreeMemory2: Computes the amount of free memory on the host.
// Version: orlando
func (host) ComputeFreeMemory2(session *Session, host HostRef) (retval int, err error) {
	method := "host.compute_free_memory"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, hostArg)
	if err != nil {
		return
	}
	retval, err = deserializeInt(method+" -> ", result)
	return
}

// AsyncComputeFreeMemory2: Computes the amount of free memory on the host.
// Version: orlando
func (host) AsyncComputeFreeMemory2(session *Session, host HostRef) (retval TaskRef, err error) {
	method := "Async.host.compute_free_memory"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, hostArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// SetHostnameLive: Sets the host name to the specified string.  Both the API and lower-level system hostname are changed immediately.
// Version: miami
//
// Errors:
// HOST_NAME_INVALID - The server name is invalid.
func (host) SetHostnameLive(session *Session, host HostRef, hostname string) (err error) {
	method := "host.set_hostname_live"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	hostnameArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "hostname"), hostname)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, hostArg, hostnameArg)
	return
}

// SetHostnameLive3: Sets the host name to the specified string.  Both the API and lower-level system hostname are changed immediately.
// Version: miami
//
// Errors:
// HOST_NAME_INVALID - The server name is invalid.
func (host) SetHostnameLive3(session *Session, host HostRef, hostname string) (err error) {
	method := "host.set_hostname_live"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	hostnameArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "hostname"), hostname)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, hostArg, hostnameArg)
	return
}

// ShutdownAgent: Shuts the agent down after a 10 second pause. WARNING: this is a dangerous operation. Any operations in progress will be aborted, and unrecoverable data loss may occur. The caller is responsible for ensuring that there are no operations in progress when this method is called.
// Version: orlando
func (host) ShutdownAgent(session *Session) (err error) {
	method := "host.shutdown_agent"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg)
	return
}

// ShutdownAgent1: Shuts the agent down after a 10 second pause. WARNING: this is a dangerous operation. Any operations in progress will be aborted, and unrecoverable data loss may occur. The caller is responsible for ensuring that there are no operations in progress when this method is called.
// Version: orlando
func (host) ShutdownAgent1(session *Session) (err error) {
	method := "host.shutdown_agent"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg)
	return
}

// RestartAgent: Restarts the agent after a 10 second pause. WARNING: this is a dangerous operation. Any operations in progress will be aborted, and unrecoverable data loss may occur. The caller is responsible for ensuring that there are no operations in progress when this method is called.
// Version: rio
func (host) RestartAgent(session *Session, host HostRef) (err error) {
	method := "host.restart_agent"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, hostArg)
	return
}

// AsyncRestartAgent: Restarts the agent after a 10 second pause. WARNING: this is a dangerous operation. Any operations in progress will be aborted, and unrecoverable data loss may occur. The caller is responsible for ensuring that there are no operations in progress when this method is called.
// Version: rio
func (host) AsyncRestartAgent(session *Session, host HostRef) (retval TaskRef, err error) {
	method := "Async.host.restart_agent"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, hostArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// RestartAgent2: Restarts the agent after a 10 second pause. WARNING: this is a dangerous operation. Any operations in progress will be aborted, and unrecoverable data loss may occur. The caller is responsible for ensuring that there are no operations in progress when this method is called.
// Version: rio
func (host) RestartAgent2(session *Session, host HostRef) (err error) {
	method := "host.restart_agent"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, hostArg)
	return
}

// AsyncRestartAgent2: Restarts the agent after a 10 second pause. WARNING: this is a dangerous operation. Any operations in progress will be aborted, and unrecoverable data loss may occur. The caller is responsible for ensuring that there are no operations in progress when this method is called.
// Version: rio
func (host) AsyncRestartAgent2(session *Session, host HostRef) (retval TaskRef, err error) {
	method := "Async.host.restart_agent"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, hostArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// GetSystemStatusCapabilities:
// Version: miami
func (host) GetSystemStatusCapabilities(session *Session, host HostRef) (retval string, err error) {
	method := "host.get_system_status_capabilities"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, hostArg)
	if err != nil {
		return
	}
	retval, err = deserializeString(method+" -> ", result)
	return
}

// GetSystemStatusCapabilities2:
// Version: miami
func (host) GetSystemStatusCapabilities2(session *Session, host HostRef) (retval string, err error) {
	method := "host.get_system_status_capabilities"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, hostArg)
	if err != nil {
		return
	}
	retval, err = deserializeString(method+" -> ", result)
	return
}

// GetManagementInterface: Returns the management interface for the specified host
// Version: tampa
func (host) GetManagementInterface(session *Session, host HostRef) (retval PIFRef, err error) {
	method := "host.get_management_interface"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, hostArg)
	if err != nil {
		return
	}
	retval, err = deserializePIFRef(method+" -> ", result)
	return
}

// AsyncGetManagementInterface: Returns the management interface for the specified host
// Version: tampa
func (host) AsyncGetManagementInterface(session *Session, host HostRef) (retval TaskRef, err error) {
	method := "Async.host.get_management_interface"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, hostArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// GetManagementInterface2: Returns the management interface for the specified host
// Version: tampa
func (host) GetManagementInterface2(session *Session, host HostRef) (retval PIFRef, err error) {
	method := "host.get_management_interface"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, hostArg)
	if err != nil {
		return
	}
	retval, err = deserializePIFRef(method+" -> ", result)
	return
}

// AsyncGetManagementInterface2: Returns the management interface for the specified host
// Version: tampa
func (host) AsyncGetManagementInterface2(session *Session, host HostRef) (retval TaskRef, err error) {
	method := "Async.host.get_management_interface"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, hostArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// ManagementDisable: Disable the management network interface
// Version: miami
func (host) ManagementDisable(session *Session) (err error) {
	method := "host.management_disable"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg)
	return
}

// ManagementDisable1: Disable the management network interface
// Version: miami
func (host) ManagementDisable1(session *Session) (err error) {
	method := "host.management_disable"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg)
	return
}

// LocalManagementReconfigure: Reconfigure the management network interface. Should only be used if Host.management_reconfigure is impossible because the network configuration is broken.
// Version: miami
func (host) LocalManagementReconfigure(session *Session, inter string) (err error) {
	method := "host.local_management_reconfigure"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	interArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "interface"), inter)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, interArg)
	return
}

// LocalManagementReconfigure2: Reconfigure the management network interface. Should only be used if Host.management_reconfigure is impossible because the network configuration is broken.
// Version: miami
func (host) LocalManagementReconfigure2(session *Session, inter string) (err error) {
	method := "host.local_management_reconfigure"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	interArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "interface"), inter)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, interArg)
	return
}

// ManagementReconfigure: Reconfigure the management network interface
// Version: miami
func (host) ManagementReconfigure(session *Session, pif PIFRef) (err error) {
	method := "host.management_reconfigure"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	pifArg, err := serializePIFRef(fmt.Sprintf("%s(%s)", method, "pif"), pif)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, pifArg)
	return
}

// AsyncManagementReconfigure: Reconfigure the management network interface
// Version: miami
func (host) AsyncManagementReconfigure(session *Session, pif PIFRef) (retval TaskRef, err error) {
	method := "Async.host.management_reconfigure"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	pifArg, err := serializePIFRef(fmt.Sprintf("%s(%s)", method, "pif"), pif)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, pifArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// ManagementReconfigure2: Reconfigure the management network interface
// Version: miami
func (host) ManagementReconfigure2(session *Session, pif PIFRef) (err error) {
	method := "host.management_reconfigure"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	pifArg, err := serializePIFRef(fmt.Sprintf("%s(%s)", method, "pif"), pif)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, pifArg)
	return
}

// AsyncManagementReconfigure2: Reconfigure the management network interface
// Version: miami
func (host) AsyncManagementReconfigure2(session *Session, pif PIFRef) (retval TaskRef, err error) {
	method := "Async.host.management_reconfigure"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	pifArg, err := serializePIFRef(fmt.Sprintf("%s(%s)", method, "pif"), pif)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, pifArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// SyslogReconfigure: Re-configure syslog logging
// Version: miami
func (host) SyslogReconfigure(session *Session, host HostRef) (err error) {
	method := "host.syslog_reconfigure"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, hostArg)
	return
}

// AsyncSyslogReconfigure: Re-configure syslog logging
// Version: miami
func (host) AsyncSyslogReconfigure(session *Session, host HostRef) (retval TaskRef, err error) {
	method := "Async.host.syslog_reconfigure"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, hostArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// SyslogReconfigure2: Re-configure syslog logging
// Version: miami
func (host) SyslogReconfigure2(session *Session, host HostRef) (err error) {
	method := "host.syslog_reconfigure"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, hostArg)
	return
}

// AsyncSyslogReconfigure2: Re-configure syslog logging
// Version: miami
func (host) AsyncSyslogReconfigure2(session *Session, host HostRef) (retval TaskRef, err error) {
	method := "Async.host.syslog_reconfigure"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, hostArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// Evacuate: Migrate all VMs off of this host, where possible.
// Version: 23.27.0
func (host) Evacuate(session *Session, host HostRef, network NetworkRef, evacuateBatchSize int) (err error) {
	method := "host.evacuate"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	networkArg, err := serializeNetworkRef(fmt.Sprintf("%s(%s)", method, "network"), network)
	if err != nil {
		return
	}
	evacuateBatchSizeArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "evacuate_batch_size"), evacuateBatchSize)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, hostArg, networkArg, evacuateBatchSizeArg)
	return
}

// AsyncEvacuate: Migrate all VMs off of this host, where possible.
// Version: 23.27.0
func (host) AsyncEvacuate(session *Session, host HostRef, network NetworkRef, evacuateBatchSize int) (retval TaskRef, err error) {
	method := "Async.host.evacuate"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	networkArg, err := serializeNetworkRef(fmt.Sprintf("%s(%s)", method, "network"), network)
	if err != nil {
		return
	}
	evacuateBatchSizeArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "evacuate_batch_size"), evacuateBatchSize)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, hostArg, networkArg, evacuateBatchSizeArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// Evacuate4: Migrate all VMs off of this host, where possible.
// Version: 23.27.0
func (host) Evacuate4(session *Session, host HostRef, network NetworkRef, evacuateBatchSize int) (err error) {
	method := "host.evacuate"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	networkArg, err := serializeNetworkRef(fmt.Sprintf("%s(%s)", method, "network"), network)
	if err != nil {
		return
	}
	evacuateBatchSizeArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "evacuate_batch_size"), evacuateBatchSize)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, hostArg, networkArg, evacuateBatchSizeArg)
	return
}

// AsyncEvacuate4: Migrate all VMs off of this host, where possible.
// Version: 23.27.0
func (host) AsyncEvacuate4(session *Session, host HostRef, network NetworkRef, evacuateBatchSize int) (retval TaskRef, err error) {
	method := "Async.host.evacuate"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	networkArg, err := serializeNetworkRef(fmt.Sprintf("%s(%s)", method, "network"), network)
	if err != nil {
		return
	}
	evacuateBatchSizeArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "evacuate_batch_size"), evacuateBatchSize)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, hostArg, networkArg, evacuateBatchSizeArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// Evacuate3: Migrate all VMs off of this host, where possible.
// Version: 1.297.0
func (host) Evacuate3(session *Session, host HostRef, network NetworkRef) (err error) {
	method := "host.evacuate"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	networkArg, err := serializeNetworkRef(fmt.Sprintf("%s(%s)", method, "network"), network)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, hostArg, networkArg)
	return
}

// AsyncEvacuate3: Migrate all VMs off of this host, where possible.
// Version: 1.297.0
func (host) AsyncEvacuate3(session *Session, host HostRef, network NetworkRef) (retval TaskRef, err error) {
	method := "Async.host.evacuate"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	networkArg, err := serializeNetworkRef(fmt.Sprintf("%s(%s)", method, "network"), network)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, hostArg, networkArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// Evacuate2: Migrate all VMs off of this host, where possible.
// Version: miami
func (host) Evacuate2(session *Session, host HostRef) (err error) {
	method := "host.evacuate"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, hostArg)
	return
}

// AsyncEvacuate2: Migrate all VMs off of this host, where possible.
// Version: miami
func (host) AsyncEvacuate2(session *Session, host HostRef) (retval TaskRef, err error) {
	method := "Async.host.evacuate"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, hostArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// GetUncooperativeResidentVMs: Return a set of VMs which are not co-operating with the host&apos;s memory control system
// Version: midnight-ride
func (host) GetUncooperativeResidentVMs(session *Session, self HostRef) (retval []VMRef, err error) {
	method := "host.get_uncooperative_resident_VMs"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeVMRefSet(method+" -> ", result)
	return
}

// AsyncGetUncooperativeResidentVMs: Return a set of VMs which are not co-operating with the host&apos;s memory control system
// Version: midnight-ride
func (host) AsyncGetUncooperativeResidentVMs(session *Session, self HostRef) (retval TaskRef, err error) {
	method := "Async.host.get_uncooperative_resident_VMs"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// GetUncooperativeResidentVMs2: Return a set of VMs which are not co-operating with the host&apos;s memory control system
// Version: midnight-ride
func (host) GetUncooperativeResidentVMs2(session *Session, self HostRef) (retval []VMRef, err error) {
	method := "host.get_uncooperative_resident_VMs"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeVMRefSet(method+" -> ", result)
	return
}

// AsyncGetUncooperativeResidentVMs2: Return a set of VMs which are not co-operating with the host&apos;s memory control system
// Version: midnight-ride
func (host) AsyncGetUncooperativeResidentVMs2(session *Session, self HostRef) (retval TaskRef, err error) {
	method := "Async.host.get_uncooperative_resident_VMs"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// GetVmsWhichPreventEvacuation: Return a set of VMs which prevent the host being evacuated, with per-VM error codes
// Version: orlando
func (host) GetVmsWhichPreventEvacuation(session *Session, self HostRef) (retval map[VMRef][]string, err error) {
	method := "host.get_vms_which_prevent_evacuation"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeVMRefToStringSetMap(method+" -> ", result)
	return
}

// AsyncGetVmsWhichPreventEvacuation: Return a set of VMs which prevent the host being evacuated, with per-VM error codes
// Version: orlando
func (host) AsyncGetVmsWhichPreventEvacuation(session *Session, self HostRef) (retval TaskRef, err error) {
	method := "Async.host.get_vms_which_prevent_evacuation"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// GetVmsWhichPreventEvacuation2: Return a set of VMs which prevent the host being evacuated, with per-VM error codes
// Version: orlando
func (host) GetVmsWhichPreventEvacuation2(session *Session, self HostRef) (retval map[VMRef][]string, err error) {
	method := "host.get_vms_which_prevent_evacuation"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeVMRefToStringSetMap(method+" -> ", result)
	return
}

// AsyncGetVmsWhichPreventEvacuation2: Return a set of VMs which prevent the host being evacuated, with per-VM error codes
// Version: orlando
func (host) AsyncGetVmsWhichPreventEvacuation2(session *Session, self HostRef) (retval TaskRef, err error) {
	method := "Async.host.get_vms_which_prevent_evacuation"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// AssertCanEvacuate: Check this host can be evacuated.
// Version: miami
func (host) AssertCanEvacuate(session *Session, host HostRef) (err error) {
	method := "host.assert_can_evacuate"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, hostArg)
	return
}

// AsyncAssertCanEvacuate: Check this host can be evacuated.
// Version: miami
func (host) AsyncAssertCanEvacuate(session *Session, host HostRef) (retval TaskRef, err error) {
	method := "Async.host.assert_can_evacuate"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, hostArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// AssertCanEvacuate2: Check this host can be evacuated.
// Version: miami
func (host) AssertCanEvacuate2(session *Session, host HostRef) (err error) {
	method := "host.assert_can_evacuate"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, hostArg)
	return
}

// AsyncAssertCanEvacuate2: Check this host can be evacuated.
// Version: miami
func (host) AsyncAssertCanEvacuate2(session *Session, host HostRef) (retval TaskRef, err error) {
	method := "Async.host.assert_can_evacuate"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, hostArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// ForgetDataSourceArchives: Forget the recorded statistics related to the specified data source
// Version: orlando
func (host) ForgetDataSourceArchives(session *Session, host HostRef, dataSource string) (err error) {
	method := "host.forget_data_source_archives"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	dataSourceArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "data_source"), dataSource)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, hostArg, dataSourceArg)
	return
}

// ForgetDataSourceArchives3: Forget the recorded statistics related to the specified data source
// Version: orlando
func (host) ForgetDataSourceArchives3(session *Session, host HostRef, dataSource string) (err error) {
	method := "host.forget_data_source_archives"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	dataSourceArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "data_source"), dataSource)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, hostArg, dataSourceArg)
	return
}

// QueryDataSource: Query the latest value of the specified data source
// Version: orlando
func (host) QueryDataSource(session *Session, host HostRef, dataSource string) (retval float64, err error) {
	method := "host.query_data_source"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	dataSourceArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "data_source"), dataSource)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, hostArg, dataSourceArg)
	if err != nil {
		return
	}
	retval, err = deserializeFloat(method+" -> ", result)
	return
}

// QueryDataSource3: Query the latest value of the specified data source
// Version: orlando
func (host) QueryDataSource3(session *Session, host HostRef, dataSource string) (retval float64, err error) {
	method := "host.query_data_source"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	dataSourceArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "data_source"), dataSource)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, hostArg, dataSourceArg)
	if err != nil {
		return
	}
	retval, err = deserializeFloat(method+" -> ", result)
	return
}

// RecordDataSource: Start recording the specified data source
// Version: orlando
func (host) RecordDataSource(session *Session, host HostRef, dataSource string) (err error) {
	method := "host.record_data_source"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	dataSourceArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "data_source"), dataSource)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, hostArg, dataSourceArg)
	return
}

// RecordDataSource3: Start recording the specified data source
// Version: orlando
func (host) RecordDataSource3(session *Session, host HostRef, dataSource string) (err error) {
	method := "host.record_data_source"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	dataSourceArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "data_source"), dataSource)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, hostArg, dataSourceArg)
	return
}

// GetDataSources:
// Version: orlando
func (host) GetDataSources(session *Session, host HostRef) (retval []DataSourceRecord, err error) {
	method := "host.get_data_sources"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, hostArg)
	if err != nil {
		return
	}
	retval, err = deserializeDataSourceRecordSet(method+" -> ", result)
	return
}

// GetDataSources2:
// Version: orlando
func (host) GetDataSources2(session *Session, host HostRef) (retval []DataSourceRecord, err error) {
	method := "host.get_data_sources"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, hostArg)
	if err != nil {
		return
	}
	retval, err = deserializeDataSourceRecordSet(method+" -> ", result)
	return
}

// EmergencyHaDisable: This call disables HA on the local host. This should only be used with extreme care.
// Version: ely
func (host) EmergencyHaDisable(session *Session, soft bool) (err error) {
	method := "host.emergency_ha_disable"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	softArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "soft"), soft)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, softArg)
	return
}

// EmergencyHaDisable2: This call disables HA on the local host. This should only be used with extreme care.
// Version: ely
func (host) EmergencyHaDisable2(session *Session, soft bool) (err error) {
	method := "host.emergency_ha_disable"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	softArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "soft"), soft)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, softArg)
	return
}

// PowerOn: Attempt to power-on the host (if the capability exists).
// Version: orlando
func (host) PowerOn(session *Session, host HostRef) (err error) {
	method := "host.power_on"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, hostArg)
	return
}

// AsyncPowerOn: Attempt to power-on the host (if the capability exists).
// Version: orlando
func (host) AsyncPowerOn(session *Session, host HostRef) (retval TaskRef, err error) {
	method := "Async.host.power_on"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, hostArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// PowerOn2: Attempt to power-on the host (if the capability exists).
// Version: orlando
func (host) PowerOn2(session *Session, host HostRef) (err error) {
	method := "host.power_on"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, hostArg)
	return
}

// AsyncPowerOn2: Attempt to power-on the host (if the capability exists).
// Version: orlando
func (host) AsyncPowerOn2(session *Session, host HostRef) (retval TaskRef, err error) {
	method := "Async.host.power_on"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, hostArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// Destroy: Destroy specified host record in database
// Version: rio
func (host) Destroy(session *Session, self HostRef) (err error) {
	method := "host.destroy"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// AsyncDestroy: Destroy specified host record in database
// Version: rio
func (host) AsyncDestroy(session *Session, self HostRef) (retval TaskRef, err error) {
	method := "Async.host.destroy"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// Destroy2: Destroy specified host record in database
// Version: rio
func (host) Destroy2(session *Session, self HostRef) (err error) {
	method := "host.destroy"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// AsyncDestroy2: Destroy specified host record in database
// Version: rio
func (host) AsyncDestroy2(session *Session, self HostRef) (retval TaskRef, err error) {
	method := "Async.host.destroy"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// LicenseRemove: Remove any license file from the specified host, and switch that host to the unlicensed edition
// Version: indigo
func (host) LicenseRemove(session *Session, host HostRef) (err error) {
	method := "host.license_remove"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, hostArg)
	return
}

// AsyncLicenseRemove: Remove any license file from the specified host, and switch that host to the unlicensed edition
// Version: indigo
func (host) AsyncLicenseRemove(session *Session, host HostRef) (retval TaskRef, err error) {
	method := "Async.host.license_remove"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, hostArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// LicenseRemove2: Remove any license file from the specified host, and switch that host to the unlicensed edition
// Version: indigo
func (host) LicenseRemove2(session *Session, host HostRef) (err error) {
	method := "host.license_remove"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, hostArg)
	return
}

// AsyncLicenseRemove2: Remove any license file from the specified host, and switch that host to the unlicensed edition
// Version: indigo
func (host) AsyncLicenseRemove2(session *Session, host HostRef) (retval TaskRef, err error) {
	method := "Async.host.license_remove"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, hostArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// LicenseAdd: Apply a new license to a host
// Version: indigo
//
// Errors:
// LICENSE_PROCESSING_ERROR - There was an error processing your license. Please contact your support representative.
func (host) LicenseAdd(session *Session, host HostRef, contents string) (err error) {
	method := "host.license_add"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	contentsArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "contents"), contents)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, hostArg, contentsArg)
	return
}

// AsyncLicenseAdd: Apply a new license to a host
// Version: indigo
//
// Errors:
// LICENSE_PROCESSING_ERROR - There was an error processing your license. Please contact your support representative.
func (host) AsyncLicenseAdd(session *Session, host HostRef, contents string) (retval TaskRef, err error) {
	method := "Async.host.license_add"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	contentsArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "contents"), contents)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, hostArg, contentsArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// LicenseAdd3: Apply a new license to a host
// Version: indigo
//
// Errors:
// LICENSE_PROCESSING_ERROR - There was an error processing your license. Please contact your support representative.
func (host) LicenseAdd3(session *Session, host HostRef, contents string) (err error) {
	method := "host.license_add"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	contentsArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "contents"), contents)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, hostArg, contentsArg)
	return
}

// AsyncLicenseAdd3: Apply a new license to a host
// Version: indigo
//
// Errors:
// LICENSE_PROCESSING_ERROR - There was an error processing your license. Please contact your support representative.
func (host) AsyncLicenseAdd3(session *Session, host HostRef, contents string) (retval TaskRef, err error) {
	method := "Async.host.license_add"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	contentsArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "contents"), contents)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, hostArg, contentsArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// LicenseApply: Apply a new license to a host
// Version: rio
//
// Errors:
// LICENSE_PROCESSING_ERROR - There was an error processing your license. Please contact your support representative.
func (host) LicenseApply(session *Session, host HostRef, contents string) (err error) {
	method := "host.license_apply"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	contentsArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "contents"), contents)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, hostArg, contentsArg)
	return
}

// AsyncLicenseApply: Apply a new license to a host
// Version: rio
//
// Errors:
// LICENSE_PROCESSING_ERROR - There was an error processing your license. Please contact your support representative.
func (host) AsyncLicenseApply(session *Session, host HostRef, contents string) (retval TaskRef, err error) {
	method := "Async.host.license_apply"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	contentsArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "contents"), contents)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, hostArg, contentsArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// LicenseApply3: Apply a new license to a host
// Version: rio
//
// Errors:
// LICENSE_PROCESSING_ERROR - There was an error processing your license. Please contact your support representative.
func (host) LicenseApply3(session *Session, host HostRef, contents string) (err error) {
	method := "host.license_apply"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	contentsArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "contents"), contents)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, hostArg, contentsArg)
	return
}

// AsyncLicenseApply3: Apply a new license to a host
// Version: rio
//
// Errors:
// LICENSE_PROCESSING_ERROR - There was an error processing your license. Please contact your support representative.
func (host) AsyncLicenseApply3(session *Session, host HostRef, contents string) (retval TaskRef, err error) {
	method := "Async.host.license_apply"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	contentsArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "contents"), contents)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, hostArg, contentsArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// ListMethods: List all supported methods
// Version: rio
func (host) ListMethods(session *Session) (retval []string, err error) {
	method := "host.list_methods"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializeStringSet(method+" -> ", result)
	return
}

// ListMethods1: List all supported methods
// Version: rio
func (host) ListMethods1(session *Session) (retval []string, err error) {
	method := "host.list_methods"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializeStringSet(method+" -> ", result)
	return
}

// BugreportUpload: Run xen-bugtool --yestoall and upload the output to support
// Version: rio
func (host) BugreportUpload(session *Session, host HostRef, url string, options map[string]string) (err error) {
	method := "host.bugreport_upload"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	urlArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "url"), url)
	if err != nil {
		return
	}
	optionsArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "options"), options)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, hostArg, urlArg, optionsArg)
	return
}

// AsyncBugreportUpload: Run xen-bugtool --yestoall and upload the output to support
// Version: rio
func (host) AsyncBugreportUpload(session *Session, host HostRef, url string, options map[string]string) (retval TaskRef, err error) {
	method := "Async.host.bugreport_upload"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	urlArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "url"), url)
	if err != nil {
		return
	}
	optionsArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "options"), options)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, hostArg, urlArg, optionsArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// BugreportUpload4: Run xen-bugtool --yestoall and upload the output to support
// Version: rio
func (host) BugreportUpload4(session *Session, host HostRef, url string, options map[string]string) (err error) {
	method := "host.bugreport_upload"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	urlArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "url"), url)
	if err != nil {
		return
	}
	optionsArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "options"), options)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, hostArg, urlArg, optionsArg)
	return
}

// AsyncBugreportUpload4: Run xen-bugtool --yestoall and upload the output to support
// Version: rio
func (host) AsyncBugreportUpload4(session *Session, host HostRef, url string, options map[string]string) (retval TaskRef, err error) {
	method := "Async.host.bugreport_upload"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	urlArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "url"), url)
	if err != nil {
		return
	}
	optionsArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "options"), options)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, hostArg, urlArg, optionsArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// SendDebugKeys: Inject the given string as debugging keys into Xen
// Version: rio
func (host) SendDebugKeys(session *Session, host HostRef, keys string) (err error) {
	method := "host.send_debug_keys"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	keysArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "keys"), keys)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, hostArg, keysArg)
	return
}

// AsyncSendDebugKeys: Inject the given string as debugging keys into Xen
// Version: rio
func (host) AsyncSendDebugKeys(session *Session, host HostRef, keys string) (retval TaskRef, err error) {
	method := "Async.host.send_debug_keys"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	keysArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "keys"), keys)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, hostArg, keysArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// SendDebugKeys3: Inject the given string as debugging keys into Xen
// Version: rio
func (host) SendDebugKeys3(session *Session, host HostRef, keys string) (err error) {
	method := "host.send_debug_keys"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	keysArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "keys"), keys)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, hostArg, keysArg)
	return
}

// AsyncSendDebugKeys3: Inject the given string as debugging keys into Xen
// Version: rio
func (host) AsyncSendDebugKeys3(session *Session, host HostRef, keys string) (retval TaskRef, err error) {
	method := "Async.host.send_debug_keys"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	keysArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "keys"), keys)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, hostArg, keysArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// GetLog: Get the host&apos;s log file
// Version: rio
func (host) GetLog(session *Session, host HostRef) (retval string, err error) {
	method := "host.get_log"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, hostArg)
	if err != nil {
		return
	}
	retval, err = deserializeString(method+" -> ", result)
	return
}

// AsyncGetLog: Get the host&apos;s log file
// Version: rio
func (host) AsyncGetLog(session *Session, host HostRef) (retval TaskRef, err error) {
	method := "Async.host.get_log"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, hostArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// GetLog2: Get the host&apos;s log file
// Version: rio
func (host) GetLog2(session *Session, host HostRef) (retval string, err error) {
	method := "host.get_log"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, hostArg)
	if err != nil {
		return
	}
	retval, err = deserializeString(method+" -> ", result)
	return
}

// AsyncGetLog2: Get the host&apos;s log file
// Version: rio
func (host) AsyncGetLog2(session *Session, host HostRef) (retval TaskRef, err error) {
	method := "Async.host.get_log"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, hostArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// DmesgClear: Get the host xen dmesg, and clear the buffer.
// Version: rio
func (host) DmesgClear(session *Session, host HostRef) (retval string, err error) {
	method := "host.dmesg_clear"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, hostArg)
	if err != nil {
		return
	}
	retval, err = deserializeString(method+" -> ", result)
	return
}

// AsyncDmesgClear: Get the host xen dmesg, and clear the buffer.
// Version: rio
func (host) AsyncDmesgClear(session *Session, host HostRef) (retval TaskRef, err error) {
	method := "Async.host.dmesg_clear"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, hostArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// DmesgClear2: Get the host xen dmesg, and clear the buffer.
// Version: rio
func (host) DmesgClear2(session *Session, host HostRef) (retval string, err error) {
	method := "host.dmesg_clear"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, hostArg)
	if err != nil {
		return
	}
	retval, err = deserializeString(method+" -> ", result)
	return
}

// AsyncDmesgClear2: Get the host xen dmesg, and clear the buffer.
// Version: rio
func (host) AsyncDmesgClear2(session *Session, host HostRef) (retval TaskRef, err error) {
	method := "Async.host.dmesg_clear"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, hostArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// Dmesg: Get the host xen dmesg.
// Version: rio
func (host) Dmesg(session *Session, host HostRef) (retval string, err error) {
	method := "host.dmesg"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, hostArg)
	if err != nil {
		return
	}
	retval, err = deserializeString(method+" -> ", result)
	return
}

// AsyncDmesg: Get the host xen dmesg.
// Version: rio
func (host) AsyncDmesg(session *Session, host HostRef) (retval TaskRef, err error) {
	method := "Async.host.dmesg"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, hostArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// Dmesg2: Get the host xen dmesg.
// Version: rio
func (host) Dmesg2(session *Session, host HostRef) (retval string, err error) {
	method := "host.dmesg"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, hostArg)
	if err != nil {
		return
	}
	retval, err = deserializeString(method+" -> ", result)
	return
}

// AsyncDmesg2: Get the host xen dmesg.
// Version: rio
func (host) AsyncDmesg2(session *Session, host HostRef) (retval TaskRef, err error) {
	method := "Async.host.dmesg"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, hostArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// Reboot: Reboot the host. (This function can only be called if there are no currently running VMs on the host and it is disabled.)
// Version: rio
func (host) Reboot(session *Session, host HostRef) (err error) {
	method := "host.reboot"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, hostArg)
	return
}

// AsyncReboot: Reboot the host. (This function can only be called if there are no currently running VMs on the host and it is disabled.)
// Version: rio
func (host) AsyncReboot(session *Session, host HostRef) (retval TaskRef, err error) {
	method := "Async.host.reboot"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, hostArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// Reboot2: Reboot the host. (This function can only be called if there are no currently running VMs on the host and it is disabled.)
// Version: rio
func (host) Reboot2(session *Session, host HostRef) (err error) {
	method := "host.reboot"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, hostArg)
	return
}

// AsyncReboot2: Reboot the host. (This function can only be called if there are no currently running VMs on the host and it is disabled.)
// Version: rio
func (host) AsyncReboot2(session *Session, host HostRef) (retval TaskRef, err error) {
	method := "Async.host.reboot"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, hostArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// Shutdown: Shutdown the host. (This function can only be called if there are no currently running VMs on the host and it is disabled.)
// Version: rio
func (host) Shutdown(session *Session, host HostRef) (err error) {
	method := "host.shutdown"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, hostArg)
	return
}

// AsyncShutdown: Shutdown the host. (This function can only be called if there are no currently running VMs on the host and it is disabled.)
// Version: rio
func (host) AsyncShutdown(session *Session, host HostRef) (retval TaskRef, err error) {
	method := "Async.host.shutdown"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, hostArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// Shutdown2: Shutdown the host. (This function can only be called if there are no currently running VMs on the host and it is disabled.)
// Version: rio
func (host) Shutdown2(session *Session, host HostRef) (err error) {
	method := "host.shutdown"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, hostArg)
	return
}

// AsyncShutdown2: Shutdown the host. (This function can only be called if there are no currently running VMs on the host and it is disabled.)
// Version: rio
func (host) AsyncShutdown2(session *Session, host HostRef) (retval TaskRef, err error) {
	method := "Async.host.shutdown"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, hostArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// Enable: Puts the host into a state in which new VMs can be started.
// Version: rio
func (host) Enable(session *Session, host HostRef) (err error) {
	method := "host.enable"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, hostArg)
	return
}

// AsyncEnable: Puts the host into a state in which new VMs can be started.
// Version: rio
func (host) AsyncEnable(session *Session, host HostRef) (retval TaskRef, err error) {
	method := "Async.host.enable"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, hostArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// Enable2: Puts the host into a state in which new VMs can be started.
// Version: rio
func (host) Enable2(session *Session, host HostRef) (err error) {
	method := "host.enable"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, hostArg)
	return
}

// AsyncEnable2: Puts the host into a state in which new VMs can be started.
// Version: rio
func (host) AsyncEnable2(session *Session, host HostRef) (retval TaskRef, err error) {
	method := "Async.host.enable"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, hostArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// Disable: Puts the host into a state in which no new VMs can be started. Currently active VMs on the host continue to execute.
// Version: 25.31.0
func (host) Disable(session *Session, host HostRef, autoEnable bool) (err error) {
	method := "host.disable"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	autoEnableArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "auto_enable"), autoEnable)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, hostArg, autoEnableArg)
	return
}

// AsyncDisable: Puts the host into a state in which no new VMs can be started. Currently active VMs on the host continue to execute.
// Version: 25.31.0
func (host) AsyncDisable(session *Session, host HostRef, autoEnable bool) (retval TaskRef, err error) {
	method := "Async.host.disable"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	autoEnableArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "auto_enable"), autoEnable)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, hostArg, autoEnableArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// Disable3: Puts the host into a state in which no new VMs can be started. Currently active VMs on the host continue to execute.
// Version: 25.31.0
func (host) Disable3(session *Session, host HostRef, autoEnable bool) (err error) {
	method := "host.disable"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	autoEnableArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "auto_enable"), autoEnable)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, hostArg, autoEnableArg)
	return
}

// AsyncDisable3: Puts the host into a state in which no new VMs can be started. Currently active VMs on the host continue to execute.
// Version: 25.31.0
func (host) AsyncDisable3(session *Session, host HostRef, autoEnable bool) (retval TaskRef, err error) {
	method := "Async.host.disable"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	autoEnableArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "auto_enable"), autoEnable)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, hostArg, autoEnableArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// Disable2: Puts the host into a state in which no new VMs can be started. Currently active VMs on the host continue to execute.
// Version: rio
func (host) Disable2(session *Session, host HostRef) (err error) {
	method := "host.disable"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, hostArg)
	return
}

// AsyncDisable2: Puts the host into a state in which no new VMs can be started. Currently active VMs on the host continue to execute.
// Version: rio
func (host) AsyncDisable2(session *Session, host HostRef) (retval TaskRef, err error) {
	method := "Async.host.disable"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, hostArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// SetDisplay: Set the display field of the given host.
// Version: cream
func (host) SetDisplay(session *Session, self HostRef, value HostDisplay) (err error) {
	method := "host.set_display"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeEnumHostDisplay(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	return
}

// SetDisplay3: Set the display field of the given host.
// Version: cream
func (host) SetDisplay3(session *Session, self HostRef, value HostDisplay) (err error) {
	method := "host.set_display"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeEnumHostDisplay(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	return
}

// SetDisplay2: Set the display field of the given host.
// Version: rio
func (host) SetDisplay2(session *Session, self HostRef) (err error) {
	method := "host.set_display"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// RemoveFromGuestVCPUsParams: Remove the given key and its corresponding value from the guest_VCPUs_params field of the given host.  If the key is not in that Map, then do nothing.
// Version: tampa
func (host) RemoveFromGuestVCPUsParams(session *Session, self HostRef, key string) (err error) {
	method := "host.remove_from_guest_VCPUs_params"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	keyArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "key"), key)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, keyArg)
	return
}

// RemoveFromGuestVCPUsParams3: Remove the given key and its corresponding value from the guest_VCPUs_params field of the given host.  If the key is not in that Map, then do nothing.
// Version: tampa
func (host) RemoveFromGuestVCPUsParams3(session *Session, self HostRef, key string) (err error) {
	method := "host.remove_from_guest_VCPUs_params"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	keyArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "key"), key)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, keyArg)
	return
}

// RemoveFromGuestVCPUsParams2: Remove the given key and its corresponding value from the guest_VCPUs_params field of the given host.  If the key is not in that Map, then do nothing.
// Version: rio
func (host) RemoveFromGuestVCPUsParams2(session *Session, self HostRef) (err error) {
	method := "host.remove_from_guest_VCPUs_params"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// AddToGuestVCPUsParams: Add the given key-value pair to the guest_VCPUs_params field of the given host.
// Version: tampa
func (host) AddToGuestVCPUsParams(session *Session, self HostRef, key string, value string) (err error) {
	method := "host.add_to_guest_VCPUs_params"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	keyArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "key"), key)
	if err != nil {
		return
	}
	valueArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, keyArg, valueArg)
	return
}

// AddToGuestVCPUsParams4: Add the given key-value pair to the guest_VCPUs_params field of the given host.
// Version: tampa
func (host) AddToGuestVCPUsParams4(session *Session, self HostRef, key string, value string) (err error) {
	method := "host.add_to_guest_VCPUs_params"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	keyArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "key"), key)
	if err != nil {
		return
	}
	valueArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, keyArg, valueArg)
	return
}

// AddToGuestVCPUsParams2: Add the given key-value pair to the guest_VCPUs_params field of the given host.
// Version: rio
func (host) AddToGuestVCPUsParams2(session *Session, self HostRef) (err error) {
	method := "host.add_to_guest_VCPUs_params"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// SetGuestVCPUsParams: Set the guest_VCPUs_params field of the given host.
// Version: tampa
func (host) SetGuestVCPUsParams(session *Session, self HostRef, value map[string]string) (err error) {
	method := "host.set_guest_VCPUs_params"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	return
}

// SetGuestVCPUsParams3: Set the guest_VCPUs_params field of the given host.
// Version: tampa
func (host) SetGuestVCPUsParams3(session *Session, self HostRef, value map[string]string) (err error) {
	method := "host.set_guest_VCPUs_params"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	return
}

// SetGuestVCPUsParams2: Set the guest_VCPUs_params field of the given host.
// Version: rio
func (host) SetGuestVCPUsParams2(session *Session, self HostRef) (err error) {
	method := "host.set_guest_VCPUs_params"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// RemoveFromLicenseServer: Remove the given key and its corresponding value from the license_server field of the given host.  If the key is not in that Map, then do nothing.
// Version: midnight-ride
func (host) RemoveFromLicenseServer(session *Session, self HostRef, key string) (err error) {
	method := "host.remove_from_license_server"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	keyArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "key"), key)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, keyArg)
	return
}

// RemoveFromLicenseServer3: Remove the given key and its corresponding value from the license_server field of the given host.  If the key is not in that Map, then do nothing.
// Version: midnight-ride
func (host) RemoveFromLicenseServer3(session *Session, self HostRef, key string) (err error) {
	method := "host.remove_from_license_server"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	keyArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "key"), key)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, keyArg)
	return
}

// RemoveFromLicenseServer2: Remove the given key and its corresponding value from the license_server field of the given host.  If the key is not in that Map, then do nothing.
// Version: rio
func (host) RemoveFromLicenseServer2(session *Session, self HostRef) (err error) {
	method := "host.remove_from_license_server"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// AddToLicenseServer: Add the given key-value pair to the license_server field of the given host.
// Version: midnight-ride
func (host) AddToLicenseServer(session *Session, self HostRef, key string, value string) (err error) {
	method := "host.add_to_license_server"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	keyArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "key"), key)
	if err != nil {
		return
	}
	valueArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, keyArg, valueArg)
	return
}

// AddToLicenseServer4: Add the given key-value pair to the license_server field of the given host.
// Version: midnight-ride
func (host) AddToLicenseServer4(session *Session, self HostRef, key string, value string) (err error) {
	method := "host.add_to_license_server"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	keyArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "key"), key)
	if err != nil {
		return
	}
	valueArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, keyArg, valueArg)
	return
}

// AddToLicenseServer2: Add the given key-value pair to the license_server field of the given host.
// Version: rio
func (host) AddToLicenseServer2(session *Session, self HostRef) (err error) {
	method := "host.add_to_license_server"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// SetLicenseServer: Set the license_server field of the given host.
// Version: midnight-ride
func (host) SetLicenseServer(session *Session, self HostRef, value map[string]string) (err error) {
	method := "host.set_license_server"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	return
}

// SetLicenseServer3: Set the license_server field of the given host.
// Version: midnight-ride
func (host) SetLicenseServer3(session *Session, self HostRef, value map[string]string) (err error) {
	method := "host.set_license_server"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	return
}

// SetLicenseServer2: Set the license_server field of the given host.
// Version: rio
func (host) SetLicenseServer2(session *Session, self HostRef) (err error) {
	method := "host.set_license_server"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// RemoveTags: Remove the given value from the tags field of the given host.  If the value is not in that Set, then do nothing.
// Version: orlando
func (host) RemoveTags(session *Session, self HostRef, value string) (err error) {
	method := "host.remove_tags"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	return
}

// RemoveTags3: Remove the given value from the tags field of the given host.  If the value is not in that Set, then do nothing.
// Version: orlando
func (host) RemoveTags3(session *Session, self HostRef, value string) (err error) {
	method := "host.remove_tags"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	return
}

// RemoveTags2: Remove the given value from the tags field of the given host.  If the value is not in that Set, then do nothing.
// Version: rio
func (host) RemoveTags2(session *Session, self HostRef) (err error) {
	method := "host.remove_tags"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// AddTags: Add the given value to the tags field of the given host.  If the value is already in that Set, then do nothing.
// Version: orlando
func (host) AddTags(session *Session, self HostRef, value string) (err error) {
	method := "host.add_tags"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	return
}

// AddTags3: Add the given value to the tags field of the given host.  If the value is already in that Set, then do nothing.
// Version: orlando
func (host) AddTags3(session *Session, self HostRef, value string) (err error) {
	method := "host.add_tags"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	return
}

// AddTags2: Add the given value to the tags field of the given host.  If the value is already in that Set, then do nothing.
// Version: rio
func (host) AddTags2(session *Session, self HostRef) (err error) {
	method := "host.add_tags"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// SetTags: Set the tags field of the given host.
// Version: orlando
func (host) SetTags(session *Session, self HostRef, value []string) (err error) {
	method := "host.set_tags"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeStringSet(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	return
}

// SetTags3: Set the tags field of the given host.
// Version: orlando
func (host) SetTags3(session *Session, self HostRef, value []string) (err error) {
	method := "host.set_tags"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeStringSet(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	return
}

// SetTags2: Set the tags field of the given host.
// Version: rio
func (host) SetTags2(session *Session, self HostRef) (err error) {
	method := "host.set_tags"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// SetAddress: Set the address field of the given host.
// Version: rio
func (host) SetAddress(session *Session, self HostRef, value string) (err error) {
	method := "host.set_address"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	return
}

// SetAddress3: Set the address field of the given host.
// Version: rio
func (host) SetAddress3(session *Session, self HostRef, value string) (err error) {
	method := "host.set_address"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	return
}

// SetHostname: Set the hostname field of the given host.
// Version: rio
func (host) SetHostname(session *Session, self HostRef, value string) (err error) {
	method := "host.set_hostname"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	return
}

// SetHostname3: Set the hostname field of the given host.
// Version: rio
func (host) SetHostname3(session *Session, self HostRef, value string) (err error) {
	method := "host.set_hostname"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	return
}

// SetCrashDumpSr: Set the crash_dump_sr field of the given host.
// Version: rio
func (host) SetCrashDumpSr(session *Session, self HostRef, value SRRef) (err error) {
	method := "host.set_crash_dump_sr"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeSRRef(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	return
}

// SetCrashDumpSr3: Set the crash_dump_sr field of the given host.
// Version: rio
func (host) SetCrashDumpSr3(session *Session, self HostRef, value SRRef) (err error) {
	method := "host.set_crash_dump_sr"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeSRRef(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	return
}

// SetSuspendImageSr: Set the suspend_image_sr field of the given host.
// Version: rio
func (host) SetSuspendImageSr(session *Session, self HostRef, value SRRef) (err error) {
	method := "host.set_suspend_image_sr"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeSRRef(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	return
}

// SetSuspendImageSr3: Set the suspend_image_sr field of the given host.
// Version: rio
func (host) SetSuspendImageSr3(session *Session, self HostRef, value SRRef) (err error) {
	method := "host.set_suspend_image_sr"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeSRRef(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	return
}

// RemoveFromLogging: Remove the given key and its corresponding value from the logging field of the given host.  If the key is not in that Map, then do nothing.
// Version: rio
func (host) RemoveFromLogging(session *Session, self HostRef, key string) (err error) {
	method := "host.remove_from_logging"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	keyArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "key"), key)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, keyArg)
	return
}

// RemoveFromLogging3: Remove the given key and its corresponding value from the logging field of the given host.  If the key is not in that Map, then do nothing.
// Version: rio
func (host) RemoveFromLogging3(session *Session, self HostRef, key string) (err error) {
	method := "host.remove_from_logging"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	keyArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "key"), key)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, keyArg)
	return
}

// AddToLogging: Add the given key-value pair to the logging field of the given host.
// Version: rio
func (host) AddToLogging(session *Session, self HostRef, key string, value string) (err error) {
	method := "host.add_to_logging"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	keyArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "key"), key)
	if err != nil {
		return
	}
	valueArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, keyArg, valueArg)
	return
}

// AddToLogging4: Add the given key-value pair to the logging field of the given host.
// Version: rio
func (host) AddToLogging4(session *Session, self HostRef, key string, value string) (err error) {
	method := "host.add_to_logging"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	keyArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "key"), key)
	if err != nil {
		return
	}
	valueArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, keyArg, valueArg)
	return
}

// SetLogging: Set the logging field of the given host.
// Version: rio
func (host) SetLogging(session *Session, self HostRef, value map[string]string) (err error) {
	method := "host.set_logging"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	return
}

// SetLogging3: Set the logging field of the given host.
// Version: rio
func (host) SetLogging3(session *Session, self HostRef, value map[string]string) (err error) {
	method := "host.set_logging"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	return
}

// RemoveFromOtherConfig: Remove the given key and its corresponding value from the other_config field of the given host.  If the key is not in that Map, then do nothing.
// Version: rio
func (host) RemoveFromOtherConfig(session *Session, self HostRef, key string) (err error) {
	method := "host.remove_from_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	keyArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "key"), key)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, keyArg)
	return
}

// RemoveFromOtherConfig3: Remove the given key and its corresponding value from the other_config field of the given host.  If the key is not in that Map, then do nothing.
// Version: rio
func (host) RemoveFromOtherConfig3(session *Session, self HostRef, key string) (err error) {
	method := "host.remove_from_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	keyArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "key"), key)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, keyArg)
	return
}

// AddToOtherConfig: Add the given key-value pair to the other_config field of the given host.
// Version: rio
func (host) AddToOtherConfig(session *Session, self HostRef, key string, value string) (err error) {
	method := "host.add_to_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	keyArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "key"), key)
	if err != nil {
		return
	}
	valueArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, keyArg, valueArg)
	return
}

// AddToOtherConfig4: Add the given key-value pair to the other_config field of the given host.
// Version: rio
func (host) AddToOtherConfig4(session *Session, self HostRef, key string, value string) (err error) {
	method := "host.add_to_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	keyArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "key"), key)
	if err != nil {
		return
	}
	valueArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, keyArg, valueArg)
	return
}

// SetOtherConfig: Set the other_config field of the given host.
// Version: rio
func (host) SetOtherConfig(session *Session, self HostRef, value map[string]string) (err error) {
	method := "host.set_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	return
}

// SetOtherConfig3: Set the other_config field of the given host.
// Version: rio
func (host) SetOtherConfig3(session *Session, self HostRef, value map[string]string) (err error) {
	method := "host.set_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	return
}

// SetNameDescription: Set the name/description field of the given host.
// Version: rio
func (host) SetNameDescription(session *Session, self HostRef, value string) (err error) {
	method := "host.set_name_description"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	return
}

// SetNameDescription3: Set the name/description field of the given host.
// Version: rio
func (host) SetNameDescription3(session *Session, self HostRef, value string) (err error) {
	method := "host.set_name_description"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	return
}

// SetNameLabel: Set the name/label field of the given host.
// Version: rio
func (host) SetNameLabel(session *Session, self HostRef, value string) (err error) {
	method := "host.set_name_label"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	return
}

// SetNameLabel3: Set the name/label field of the given host.
// Version: rio
func (host) SetNameLabel3(session *Session, self HostRef, value string) (err error) {
	method := "host.set_name_label"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	return
}

// GetTimezone: Get the timezone field of the given host.
// Version: rio
func (host) GetTimezone(session *Session, self HostRef) (retval string, err error) {
	method := "host.get_timezone"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeString(method+" -> ", result)
	return
}

// GetTimezone2: Get the timezone field of the given host.
// Version: rio
func (host) GetTimezone2(session *Session, self HostRef) (retval string, err error) {
	method := "host.get_timezone"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeString(method+" -> ", result)
	return
}

// GetNtpCustomServers: Get the ntp_custom_servers field of the given host.
// Version: rio
func (host) GetNtpCustomServers(session *Session, self HostRef) (retval []string, err error) {
	method := "host.get_ntp_custom_servers"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeStringSet(method+" -> ", result)
	return
}

// GetNtpCustomServers2: Get the ntp_custom_servers field of the given host.
// Version: rio
func (host) GetNtpCustomServers2(session *Session, self HostRef) (retval []string, err error) {
	method := "host.get_ntp_custom_servers"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeStringSet(method+" -> ", result)
	return
}

// GetNtpMode: Get the ntp_mode field of the given host.
// Version: rio
func (host) GetNtpMode(session *Session, self HostRef) (retval HostNtpMode, err error) {
	method := "host.get_ntp_mode"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeEnumHostNtpMode(method+" -> ", result)
	return
}

// GetNtpMode2: Get the ntp_mode field of the given host.
// Version: rio
func (host) GetNtpMode2(session *Session, self HostRef) (retval HostNtpMode, err error) {
	method := "host.get_ntp_mode"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeEnumHostNtpMode(method+" -> ", result)
	return
}

// GetSecureBoot: Get the secure_boot field of the given host.
// Version: rio
func (host) GetSecureBoot(session *Session, self HostRef) (retval bool, err error) {
	method := "host.get_secure_boot"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeBool(method+" -> ", result)
	return
}

// GetSecureBoot2: Get the secure_boot field of the given host.
// Version: rio
func (host) GetSecureBoot2(session *Session, self HostRef) (retval bool, err error) {
	method := "host.get_secure_boot"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeBool(method+" -> ", result)
	return
}

// GetMaxCstate: Get the max_cstate field of the given host.
// Version: rio
func (host) GetMaxCstate(session *Session, self HostRef) (retval string, err error) {
	method := "host.get_max_cstate"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeString(method+" -> ", result)
	return
}

// GetMaxCstate2: Get the max_cstate field of the given host.
// Version: rio
func (host) GetMaxCstate2(session *Session, self HostRef) (retval string, err error) {
	method := "host.get_max_cstate"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeString(method+" -> ", result)
	return
}

// GetSSHAutoMode: Get the ssh_auto_mode field of the given host.
// Version: rio
func (host) GetSSHAutoMode(session *Session, self HostRef) (retval bool, err error) {
	method := "host.get_ssh_auto_mode"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeBool(method+" -> ", result)
	return
}

// GetSSHAutoMode2: Get the ssh_auto_mode field of the given host.
// Version: rio
func (host) GetSSHAutoMode2(session *Session, self HostRef) (retval bool, err error) {
	method := "host.get_ssh_auto_mode"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeBool(method+" -> ", result)
	return
}

// GetConsoleIdleTimeout: Get the console_idle_timeout field of the given host.
// Version: rio
func (host) GetConsoleIdleTimeout(session *Session, self HostRef) (retval int, err error) {
	method := "host.get_console_idle_timeout"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeInt(method+" -> ", result)
	return
}

// GetConsoleIdleTimeout2: Get the console_idle_timeout field of the given host.
// Version: rio
func (host) GetConsoleIdleTimeout2(session *Session, self HostRef) (retval int, err error) {
	method := "host.get_console_idle_timeout"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeInt(method+" -> ", result)
	return
}

// GetSSHExpiry: Get the ssh_expiry field of the given host.
// Version: rio
func (host) GetSSHExpiry(session *Session, self HostRef) (retval time.Time, err error) {
	method := "host.get_ssh_expiry"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeTime(method+" -> ", result)
	return
}

// GetSSHExpiry2: Get the ssh_expiry field of the given host.
// Version: rio
func (host) GetSSHExpiry2(session *Session, self HostRef) (retval time.Time, err error) {
	method := "host.get_ssh_expiry"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeTime(method+" -> ", result)
	return
}

// GetSSHEnabledTimeout: Get the ssh_enabled_timeout field of the given host.
// Version: rio
func (host) GetSSHEnabledTimeout(session *Session, self HostRef) (retval int, err error) {
	method := "host.get_ssh_enabled_timeout"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeInt(method+" -> ", result)
	return
}

// GetSSHEnabledTimeout2: Get the ssh_enabled_timeout field of the given host.
// Version: rio
func (host) GetSSHEnabledTimeout2(session *Session, self HostRef) (retval int, err error) {
	method := "host.get_ssh_enabled_timeout"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeInt(method+" -> ", result)
	return
}

// GetSSHEnabled: Get the ssh_enabled field of the given host.
// Version: rio
func (host) GetSSHEnabled(session *Session, self HostRef) (retval bool, err error) {
	method := "host.get_ssh_enabled"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeBool(method+" -> ", result)
	return
}

// GetSSHEnabled2: Get the ssh_enabled field of the given host.
// Version: rio
func (host) GetSSHEnabled2(session *Session, self HostRef) (retval bool, err error) {
	method := "host.get_ssh_enabled"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeBool(method+" -> ", result)
	return
}

// GetLastUpdateHash: Get the last_update_hash field of the given host.
// Version: rio
func (host) GetLastUpdateHash(session *Session, self HostRef) (retval string, err error) {
	method := "host.get_last_update_hash"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeString(method+" -> ", result)
	return
}

// GetLastUpdateHash2: Get the last_update_hash field of the given host.
// Version: rio
func (host) GetLastUpdateHash2(session *Session, self HostRef) (retval string, err error) {
	method := "host.get_last_update_hash"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeString(method+" -> ", result)
	return
}

// GetPendingGuidancesFull: Get the pending_guidances_full field of the given host.
// Version: rio
func (host) GetPendingGuidancesFull(session *Session, self HostRef) (retval []UpdateGuidances, err error) {
	method := "host.get_pending_guidances_full"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeEnumUpdateGuidancesSet(method+" -> ", result)
	return
}

// GetPendingGuidancesFull2: Get the pending_guidances_full field of the given host.
// Version: rio
func (host) GetPendingGuidancesFull2(session *Session, self HostRef) (retval []UpdateGuidances, err error) {
	method := "host.get_pending_guidances_full"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeEnumUpdateGuidancesSet(method+" -> ", result)
	return
}

// GetPendingGuidancesRecommended: Get the pending_guidances_recommended field of the given host.
// Version: rio
func (host) GetPendingGuidancesRecommended(session *Session, self HostRef) (retval []UpdateGuidances, err error) {
	method := "host.get_pending_guidances_recommended"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeEnumUpdateGuidancesSet(method+" -> ", result)
	return
}

// GetPendingGuidancesRecommended2: Get the pending_guidances_recommended field of the given host.
// Version: rio
func (host) GetPendingGuidancesRecommended2(session *Session, self HostRef) (retval []UpdateGuidances, err error) {
	method := "host.get_pending_guidances_recommended"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeEnumUpdateGuidancesSet(method+" -> ", result)
	return
}

// GetNumaAffinityPolicy: Get the numa_affinity_policy field of the given host.
// Version: rio
func (host) GetNumaAffinityPolicy(session *Session, self HostRef) (retval HostNumaAffinityPolicy, err error) {
	method := "host.get_numa_affinity_policy"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeEnumHostNumaAffinityPolicy(method+" -> ", result)
	return
}

// GetNumaAffinityPolicy2: Get the numa_affinity_policy field of the given host.
// Version: rio
func (host) GetNumaAffinityPolicy2(session *Session, self HostRef) (retval HostNumaAffinityPolicy, err error) {
	method := "host.get_numa_affinity_policy"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeEnumHostNumaAffinityPolicy(method+" -> ", result)
	return
}

// GetLatestSyncedUpdatesApplied: Get the latest_synced_updates_applied field of the given host.
// Version: rio
func (host) GetLatestSyncedUpdatesApplied(session *Session, self HostRef) (retval LatestSyncedUpdatesAppliedState, err error) {
	method := "host.get_latest_synced_updates_applied"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeEnumLatestSyncedUpdatesAppliedState(method+" -> ", result)
	return
}

// GetLatestSyncedUpdatesApplied2: Get the latest_synced_updates_applied field of the given host.
// Version: rio
func (host) GetLatestSyncedUpdatesApplied2(session *Session, self HostRef) (retval LatestSyncedUpdatesAppliedState, err error) {
	method := "host.get_latest_synced_updates_applied"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeEnumLatestSyncedUpdatesAppliedState(method+" -> ", result)
	return
}

// GetHTTPSOnly: Get the https_only field of the given host.
// Version: rio
func (host) GetHTTPSOnly(session *Session, self HostRef) (retval bool, err error) {
	method := "host.get_https_only"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeBool(method+" -> ", result)
	return
}

// GetHTTPSOnly2: Get the https_only field of the given host.
// Version: rio
func (host) GetHTTPSOnly2(session *Session, self HostRef) (retval bool, err error) {
	method := "host.get_https_only"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeBool(method+" -> ", result)
	return
}

// GetLastSoftwareUpdate: Get the last_software_update field of the given host.
// Version: rio
func (host) GetLastSoftwareUpdate(session *Session, self HostRef) (retval time.Time, err error) {
	method := "host.get_last_software_update"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeTime(method+" -> ", result)
	return
}

// GetLastSoftwareUpdate2: Get the last_software_update field of the given host.
// Version: rio
func (host) GetLastSoftwareUpdate2(session *Session, self HostRef) (retval time.Time, err error) {
	method := "host.get_last_software_update"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeTime(method+" -> ", result)
	return
}

// GetTLSVerificationEnabled: Get the tls_verification_enabled field of the given host.
// Version: rio
func (host) GetTLSVerificationEnabled(session *Session, self HostRef) (retval bool, err error) {
	method := "host.get_tls_verification_enabled"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeBool(method+" -> ", result)
	return
}

// GetTLSVerificationEnabled2: Get the tls_verification_enabled field of the given host.
// Version: rio
func (host) GetTLSVerificationEnabled2(session *Session, self HostRef) (retval bool, err error) {
	method := "host.get_tls_verification_enabled"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeBool(method+" -> ", result)
	return
}

// GetPendingGuidances: Get the pending_guidances field of the given host.
// Version: rio
func (host) GetPendingGuidances(session *Session, self HostRef) (retval []UpdateGuidances, err error) {
	method := "host.get_pending_guidances"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeEnumUpdateGuidancesSet(method+" -> ", result)
	return
}

// GetPendingGuidances2: Get the pending_guidances field of the given host.
// Version: rio
func (host) GetPendingGuidances2(session *Session, self HostRef) (retval []UpdateGuidances, err error) {
	method := "host.get_pending_guidances"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeEnumUpdateGuidancesSet(method+" -> ", result)
	return
}

// GetEditions: Get the editions field of the given host.
// Version: rio
func (host) GetEditions(session *Session, self HostRef) (retval []string, err error) {
	method := "host.get_editions"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeStringSet(method+" -> ", result)
	return
}

// GetEditions2: Get the editions field of the given host.
// Version: rio
func (host) GetEditions2(session *Session, self HostRef) (retval []string, err error) {
	method := "host.get_editions"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeStringSet(method+" -> ", result)
	return
}

// GetCertificates: Get the certificates field of the given host.
// Version: rio
func (host) GetCertificates(session *Session, self HostRef) (retval []CertificateRef, err error) {
	method := "host.get_certificates"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeCertificateRefSet(method+" -> ", result)
	return
}

// GetCertificates2: Get the certificates field of the given host.
// Version: rio
func (host) GetCertificates2(session *Session, self HostRef) (retval []CertificateRef, err error) {
	method := "host.get_certificates"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeCertificateRefSet(method+" -> ", result)
	return
}

// GetUefiCertificates: Get the uefi_certificates field of the given host.
// Version: rio
func (host) GetUefiCertificates(session *Session, self HostRef) (retval string, err error) {
	method := "host.get_uefi_certificates"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeString(method+" -> ", result)
	return
}

// GetUefiCertificates2: Get the uefi_certificates field of the given host.
// Version: rio
func (host) GetUefiCertificates2(session *Session, self HostRef) (retval string, err error) {
	method := "host.get_uefi_certificates"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeString(method+" -> ", result)
	return
}

// GetMultipathing: Get the multipathing field of the given host.
// Version: rio
func (host) GetMultipathing(session *Session, self HostRef) (retval bool, err error) {
	method := "host.get_multipathing"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeBool(method+" -> ", result)
	return
}

// GetMultipathing2: Get the multipathing field of the given host.
// Version: rio
func (host) GetMultipathing2(session *Session, self HostRef) (retval bool, err error) {
	method := "host.get_multipathing"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeBool(method+" -> ", result)
	return
}

// GetIscsiIqn: Get the iscsi_iqn field of the given host.
// Version: rio
func (host) GetIscsiIqn(session *Session, self HostRef) (retval string, err error) {
	method := "host.get_iscsi_iqn"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeString(method+" -> ", result)
	return
}

// GetIscsiIqn2: Get the iscsi_iqn field of the given host.
// Version: rio
func (host) GetIscsiIqn2(session *Session, self HostRef) (retval string, err error) {
	method := "host.get_iscsi_iqn"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeString(method+" -> ", result)
	return
}

// GetFeatures: Get the features field of the given host.
// Version: rio
func (host) GetFeatures(session *Session, self HostRef) (retval []FeatureRef, err error) {
	method := "host.get_features"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeFeatureRefSet(method+" -> ", result)
	return
}

// GetFeatures2: Get the features field of the given host.
// Version: rio
func (host) GetFeatures2(session *Session, self HostRef) (retval []FeatureRef, err error) {
	method := "host.get_features"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeFeatureRefSet(method+" -> ", result)
	return
}

// GetUpdatesRequiringReboot: Get the updates_requiring_reboot field of the given host.
// Version: rio
func (host) GetUpdatesRequiringReboot(session *Session, self HostRef) (retval []PoolUpdateRef, err error) {
	method := "host.get_updates_requiring_reboot"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializePoolUpdateRefSet(method+" -> ", result)
	return
}

// GetUpdatesRequiringReboot2: Get the updates_requiring_reboot field of the given host.
// Version: rio
func (host) GetUpdatesRequiringReboot2(session *Session, self HostRef) (retval []PoolUpdateRef, err error) {
	method := "host.get_updates_requiring_reboot"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializePoolUpdateRefSet(method+" -> ", result)
	return
}

// GetControlDomain: Get the control_domain field of the given host.
// Version: rio
func (host) GetControlDomain(session *Session, self HostRef) (retval VMRef, err error) {
	method := "host.get_control_domain"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeVMRef(method+" -> ", result)
	return
}

// GetControlDomain2: Get the control_domain field of the given host.
// Version: rio
func (host) GetControlDomain2(session *Session, self HostRef) (retval VMRef, err error) {
	method := "host.get_control_domain"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeVMRef(method+" -> ", result)
	return
}

// GetVirtualHardwarePlatformVersions: Get the virtual_hardware_platform_versions field of the given host.
// Version: rio
func (host) GetVirtualHardwarePlatformVersions(session *Session, self HostRef) (retval []int, err error) {
	method := "host.get_virtual_hardware_platform_versions"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeIntSet(method+" -> ", result)
	return
}

// GetVirtualHardwarePlatformVersions2: Get the virtual_hardware_platform_versions field of the given host.
// Version: rio
func (host) GetVirtualHardwarePlatformVersions2(session *Session, self HostRef) (retval []int, err error) {
	method := "host.get_virtual_hardware_platform_versions"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeIntSet(method+" -> ", result)
	return
}

// GetDisplay: Get the display field of the given host.
// Version: rio
func (host) GetDisplay(session *Session, self HostRef) (retval HostDisplay, err error) {
	method := "host.get_display"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeEnumHostDisplay(method+" -> ", result)
	return
}

// GetDisplay2: Get the display field of the given host.
// Version: rio
func (host) GetDisplay2(session *Session, self HostRef) (retval HostDisplay, err error) {
	method := "host.get_display"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeEnumHostDisplay(method+" -> ", result)
	return
}

// GetGuestVCPUsParams: Get the guest_VCPUs_params field of the given host.
// Version: rio
func (host) GetGuestVCPUsParams(session *Session, self HostRef) (retval map[string]string, err error) {
	method := "host.get_guest_VCPUs_params"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeStringToStringMap(method+" -> ", result)
	return
}

// GetGuestVCPUsParams2: Get the guest_VCPUs_params field of the given host.
// Version: rio
func (host) GetGuestVCPUsParams2(session *Session, self HostRef) (retval map[string]string, err error) {
	method := "host.get_guest_VCPUs_params"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeStringToStringMap(method+" -> ", result)
	return
}

// GetSslLegacy: Get the ssl_legacy field of the given host.
// Version: rio
func (host) GetSslLegacy(session *Session, self HostRef) (retval bool, err error) {
	method := "host.get_ssl_legacy"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeBool(method+" -> ", result)
	return
}

// GetSslLegacy2: Get the ssl_legacy field of the given host.
// Version: rio
func (host) GetSslLegacy2(session *Session, self HostRef) (retval bool, err error) {
	method := "host.get_ssl_legacy"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeBool(method+" -> ", result)
	return
}

// GetPUSBs: Get the PUSBs field of the given host.
// Version: rio
func (host) GetPUSBs(session *Session, self HostRef) (retval []PUSBRef, err error) {
	method := "host.get_PUSBs"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializePUSBRefSet(method+" -> ", result)
	return
}

// GetPUSBs2: Get the PUSBs field of the given host.
// Version: rio
func (host) GetPUSBs2(session *Session, self HostRef) (retval []PUSBRef, err error) {
	method := "host.get_PUSBs"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializePUSBRefSet(method+" -> ", result)
	return
}

// GetPGPUs: Get the PGPUs field of the given host.
// Version: rio
func (host) GetPGPUs(session *Session, self HostRef) (retval []PGPURef, err error) {
	method := "host.get_PGPUs"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializePGPURefSet(method+" -> ", result)
	return
}

// GetPGPUs2: Get the PGPUs field of the given host.
// Version: rio
func (host) GetPGPUs2(session *Session, self HostRef) (retval []PGPURef, err error) {
	method := "host.get_PGPUs"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializePGPURefSet(method+" -> ", result)
	return
}

// GetPCIs: Get the PCIs field of the given host.
// Version: rio
func (host) GetPCIs(session *Session, self HostRef) (retval []PCIRef, err error) {
	method := "host.get_PCIs"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializePCIRefSet(method+" -> ", result)
	return
}

// GetPCIs2: Get the PCIs field of the given host.
// Version: rio
func (host) GetPCIs2(session *Session, self HostRef) (retval []PCIRef, err error) {
	method := "host.get_PCIs"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializePCIRefSet(method+" -> ", result)
	return
}

// GetChipsetInfo: Get the chipset_info field of the given host.
// Version: rio
func (host) GetChipsetInfo(session *Session, self HostRef) (retval map[string]string, err error) {
	method := "host.get_chipset_info"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeStringToStringMap(method+" -> ", result)
	return
}

// GetChipsetInfo2: Get the chipset_info field of the given host.
// Version: rio
func (host) GetChipsetInfo2(session *Session, self HostRef) (retval map[string]string, err error) {
	method := "host.get_chipset_info"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeStringToStringMap(method+" -> ", result)
	return
}

// GetLocalCacheSr: Get the local_cache_sr field of the given host.
// Version: rio
func (host) GetLocalCacheSr(session *Session, self HostRef) (retval SRRef, err error) {
	method := "host.get_local_cache_sr"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeSRRef(method+" -> ", result)
	return
}

// GetLocalCacheSr2: Get the local_cache_sr field of the given host.
// Version: rio
func (host) GetLocalCacheSr2(session *Session, self HostRef) (retval SRRef, err error) {
	method := "host.get_local_cache_sr"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeSRRef(method+" -> ", result)
	return
}

// GetPowerOnConfig: Get the power_on_config field of the given host.
// Version: rio
func (host) GetPowerOnConfig(session *Session, self HostRef) (retval map[string]string, err error) {
	method := "host.get_power_on_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeStringToStringMap(method+" -> ", result)
	return
}

// GetPowerOnConfig2: Get the power_on_config field of the given host.
// Version: rio
func (host) GetPowerOnConfig2(session *Session, self HostRef) (retval map[string]string, err error) {
	method := "host.get_power_on_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeStringToStringMap(method+" -> ", result)
	return
}

// GetPowerOnMode: Get the power_on_mode field of the given host.
// Version: rio
func (host) GetPowerOnMode(session *Session, self HostRef) (retval string, err error) {
	method := "host.get_power_on_mode"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeString(method+" -> ", result)
	return
}

// GetPowerOnMode2: Get the power_on_mode field of the given host.
// Version: rio
func (host) GetPowerOnMode2(session *Session, self HostRef) (retval string, err error) {
	method := "host.get_power_on_mode"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeString(method+" -> ", result)
	return
}

// GetBiosStrings: Get the bios_strings field of the given host.
// Version: rio
func (host) GetBiosStrings(session *Session, self HostRef) (retval map[string]string, err error) {
	method := "host.get_bios_strings"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeStringToStringMap(method+" -> ", result)
	return
}

// GetBiosStrings2: Get the bios_strings field of the given host.
// Version: rio
func (host) GetBiosStrings2(session *Session, self HostRef) (retval map[string]string, err error) {
	method := "host.get_bios_strings"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeStringToStringMap(method+" -> ", result)
	return
}

// GetLicenseServer: Get the license_server field of the given host.
// Version: rio
func (host) GetLicenseServer(session *Session, self HostRef) (retval map[string]string, err error) {
	method := "host.get_license_server"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeStringToStringMap(method+" -> ", result)
	return
}

// GetLicenseServer2: Get the license_server field of the given host.
// Version: rio
func (host) GetLicenseServer2(session *Session, self HostRef) (retval map[string]string, err error) {
	method := "host.get_license_server"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeStringToStringMap(method+" -> ", result)
	return
}

// GetEdition: Get the edition field of the given host.
// Version: rio
func (host) GetEdition(session *Session, self HostRef) (retval string, err error) {
	method := "host.get_edition"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeString(method+" -> ", result)
	return
}

// GetEdition2: Get the edition field of the given host.
// Version: rio
func (host) GetEdition2(session *Session, self HostRef) (retval string, err error) {
	method := "host.get_edition"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeString(method+" -> ", result)
	return
}

// GetExternalAuthConfiguration: Get the external_auth_configuration field of the given host.
// Version: rio
func (host) GetExternalAuthConfiguration(session *Session, self HostRef) (retval map[string]string, err error) {
	method := "host.get_external_auth_configuration"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeStringToStringMap(method+" -> ", result)
	return
}

// GetExternalAuthConfiguration2: Get the external_auth_configuration field of the given host.
// Version: rio
func (host) GetExternalAuthConfiguration2(session *Session, self HostRef) (retval map[string]string, err error) {
	method := "host.get_external_auth_configuration"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeStringToStringMap(method+" -> ", result)
	return
}

// GetExternalAuthServiceName: Get the external_auth_service_name field of the given host.
// Version: rio
func (host) GetExternalAuthServiceName(session *Session, self HostRef) (retval string, err error) {
	method := "host.get_external_auth_service_name"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeString(method+" -> ", result)
	return
}

// GetExternalAuthServiceName2: Get the external_auth_service_name field of the given host.
// Version: rio
func (host) GetExternalAuthServiceName2(session *Session, self HostRef) (retval string, err error) {
	method := "host.get_external_auth_service_name"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeString(method+" -> ", result)
	return
}

// GetExternalAuthType: Get the external_auth_type field of the given host.
// Version: rio
func (host) GetExternalAuthType(session *Session, self HostRef) (retval string, err error) {
	method := "host.get_external_auth_type"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeString(method+" -> ", result)
	return
}

// GetExternalAuthType2: Get the external_auth_type field of the given host.
// Version: rio
func (host) GetExternalAuthType2(session *Session, self HostRef) (retval string, err error) {
	method := "host.get_external_auth_type"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeString(method+" -> ", result)
	return
}

// GetTags: Get the tags field of the given host.
// Version: rio
func (host) GetTags(session *Session, self HostRef) (retval []string, err error) {
	method := "host.get_tags"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeStringSet(method+" -> ", result)
	return
}

// GetTags2: Get the tags field of the given host.
// Version: rio
func (host) GetTags2(session *Session, self HostRef) (retval []string, err error) {
	method := "host.get_tags"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeStringSet(method+" -> ", result)
	return
}

// GetBlobs: Get the blobs field of the given host.
// Version: rio
func (host) GetBlobs(session *Session, self HostRef) (retval map[string]BlobRef, err error) {
	method := "host.get_blobs"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeStringToBlobRefMap(method+" -> ", result)
	return
}

// GetBlobs2: Get the blobs field of the given host.
// Version: rio
func (host) GetBlobs2(session *Session, self HostRef) (retval map[string]BlobRef, err error) {
	method := "host.get_blobs"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeStringToBlobRefMap(method+" -> ", result)
	return
}

// GetHaNetworkPeers: Get the ha_network_peers field of the given host.
// Version: rio
func (host) GetHaNetworkPeers(session *Session, self HostRef) (retval []string, err error) {
	method := "host.get_ha_network_peers"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeStringSet(method+" -> ", result)
	return
}

// GetHaNetworkPeers2: Get the ha_network_peers field of the given host.
// Version: rio
func (host) GetHaNetworkPeers2(session *Session, self HostRef) (retval []string, err error) {
	method := "host.get_ha_network_peers"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeStringSet(method+" -> ", result)
	return
}

// GetHaStatefiles: Get the ha_statefiles field of the given host.
// Version: rio
func (host) GetHaStatefiles(session *Session, self HostRef) (retval []string, err error) {
	method := "host.get_ha_statefiles"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeStringSet(method+" -> ", result)
	return
}

// GetHaStatefiles2: Get the ha_statefiles field of the given host.
// Version: rio
func (host) GetHaStatefiles2(session *Session, self HostRef) (retval []string, err error) {
	method := "host.get_ha_statefiles"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeStringSet(method+" -> ", result)
	return
}

// GetLicenseParams: Get the license_params field of the given host.
// Version: rio
func (host) GetLicenseParams(session *Session, self HostRef) (retval map[string]string, err error) {
	method := "host.get_license_params"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeStringToStringMap(method+" -> ", result)
	return
}

// GetLicenseParams2: Get the license_params field of the given host.
// Version: rio
func (host) GetLicenseParams2(session *Session, self HostRef) (retval map[string]string, err error) {
	method := "host.get_license_params"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeStringToStringMap(method+" -> ", result)
	return
}

// GetMetrics: Get the metrics field of the given host.
// Version: rio
func (host) GetMetrics(session *Session, self HostRef) (retval HostMetricsRef, err error) {
	method := "host.get_metrics"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeHostMetricsRef(method+" -> ", result)
	return
}

// GetMetrics2: Get the metrics field of the given host.
// Version: rio
func (host) GetMetrics2(session *Session, self HostRef) (retval HostMetricsRef, err error) {
	method := "host.get_metrics"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeHostMetricsRef(method+" -> ", result)
	return
}

// GetAddress: Get the address field of the given host.
// Version: rio
func (host) GetAddress(session *Session, self HostRef) (retval string, err error) {
	method := "host.get_address"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeString(method+" -> ", result)
	return
}

// GetAddress2: Get the address field of the given host.
// Version: rio
func (host) GetAddress2(session *Session, self HostRef) (retval string, err error) {
	method := "host.get_address"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeString(method+" -> ", result)
	return
}

// GetHostname: Get the hostname field of the given host.
// Version: rio
func (host) GetHostname(session *Session, self HostRef) (retval string, err error) {
	method := "host.get_hostname"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeString(method+" -> ", result)
	return
}

// GetHostname2: Get the hostname field of the given host.
// Version: rio
func (host) GetHostname2(session *Session, self HostRef) (retval string, err error) {
	method := "host.get_hostname"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeString(method+" -> ", result)
	return
}

// GetCPUInfo: Get the cpu_info field of the given host.
// Version: rio
func (host) GetCPUInfo(session *Session, self HostRef) (retval map[string]string, err error) {
	method := "host.get_cpu_info"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeStringToStringMap(method+" -> ", result)
	return
}

// GetCPUInfo2: Get the cpu_info field of the given host.
// Version: rio
func (host) GetCPUInfo2(session *Session, self HostRef) (retval map[string]string, err error) {
	method := "host.get_cpu_info"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeStringToStringMap(method+" -> ", result)
	return
}

// GetHostCPUs: Get the host_CPUs field of the given host.
// Version: rio
func (host) GetHostCPUs(session *Session, self HostRef) (retval []HostCPURef, err error) {
	method := "host.get_host_CPUs"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeHostCPURefSet(method+" -> ", result)
	return
}

// GetHostCPUs2: Get the host_CPUs field of the given host.
// Version: rio
func (host) GetHostCPUs2(session *Session, self HostRef) (retval []HostCPURef, err error) {
	method := "host.get_host_CPUs"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeHostCPURefSet(method+" -> ", result)
	return
}

// GetPBDs: Get the PBDs field of the given host.
// Version: rio
func (host) GetPBDs(session *Session, self HostRef) (retval []PBDRef, err error) {
	method := "host.get_PBDs"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializePBDRefSet(method+" -> ", result)
	return
}

// GetPBDs2: Get the PBDs field of the given host.
// Version: rio
func (host) GetPBDs2(session *Session, self HostRef) (retval []PBDRef, err error) {
	method := "host.get_PBDs"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializePBDRefSet(method+" -> ", result)
	return
}

// GetUpdates: Get the updates field of the given host.
// Version: rio
func (host) GetUpdates(session *Session, self HostRef) (retval []PoolUpdateRef, err error) {
	method := "host.get_updates"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializePoolUpdateRefSet(method+" -> ", result)
	return
}

// GetUpdates2: Get the updates field of the given host.
// Version: rio
func (host) GetUpdates2(session *Session, self HostRef) (retval []PoolUpdateRef, err error) {
	method := "host.get_updates"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializePoolUpdateRefSet(method+" -> ", result)
	return
}

// GetPatches: Get the patches field of the given host.
// Version: rio
func (host) GetPatches(session *Session, self HostRef) (retval []HostPatchRef, err error) {
	method := "host.get_patches"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeHostPatchRefSet(method+" -> ", result)
	return
}

// GetPatches2: Get the patches field of the given host.
// Version: rio
func (host) GetPatches2(session *Session, self HostRef) (retval []HostPatchRef, err error) {
	method := "host.get_patches"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeHostPatchRefSet(method+" -> ", result)
	return
}

// GetCrashdumps: Get the crashdumps field of the given host.
// Version: rio
func (host) GetCrashdumps(session *Session, self HostRef) (retval []HostCrashdumpRef, err error) {
	method := "host.get_crashdumps"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeHostCrashdumpRefSet(method+" -> ", result)
	return
}

// GetCrashdumps2: Get the crashdumps field of the given host.
// Version: rio
func (host) GetCrashdumps2(session *Session, self HostRef) (retval []HostCrashdumpRef, err error) {
	method := "host.get_crashdumps"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeHostCrashdumpRefSet(method+" -> ", result)
	return
}

// GetCrashDumpSr: Get the crash_dump_sr field of the given host.
// Version: rio
func (host) GetCrashDumpSr(session *Session, self HostRef) (retval SRRef, err error) {
	method := "host.get_crash_dump_sr"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeSRRef(method+" -> ", result)
	return
}

// GetCrashDumpSr2: Get the crash_dump_sr field of the given host.
// Version: rio
func (host) GetCrashDumpSr2(session *Session, self HostRef) (retval SRRef, err error) {
	method := "host.get_crash_dump_sr"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeSRRef(method+" -> ", result)
	return
}

// GetSuspendImageSr: Get the suspend_image_sr field of the given host.
// Version: rio
func (host) GetSuspendImageSr(session *Session, self HostRef) (retval SRRef, err error) {
	method := "host.get_suspend_image_sr"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeSRRef(method+" -> ", result)
	return
}

// GetSuspendImageSr2: Get the suspend_image_sr field of the given host.
// Version: rio
func (host) GetSuspendImageSr2(session *Session, self HostRef) (retval SRRef, err error) {
	method := "host.get_suspend_image_sr"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeSRRef(method+" -> ", result)
	return
}

// GetPIFs: Get the PIFs field of the given host.
// Version: rio
func (host) GetPIFs(session *Session, self HostRef) (retval []PIFRef, err error) {
	method := "host.get_PIFs"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializePIFRefSet(method+" -> ", result)
	return
}

// GetPIFs2: Get the PIFs field of the given host.
// Version: rio
func (host) GetPIFs2(session *Session, self HostRef) (retval []PIFRef, err error) {
	method := "host.get_PIFs"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializePIFRefSet(method+" -> ", result)
	return
}

// GetLogging: Get the logging field of the given host.
// Version: rio
func (host) GetLogging(session *Session, self HostRef) (retval map[string]string, err error) {
	method := "host.get_logging"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeStringToStringMap(method+" -> ", result)
	return
}

// GetLogging2: Get the logging field of the given host.
// Version: rio
func (host) GetLogging2(session *Session, self HostRef) (retval map[string]string, err error) {
	method := "host.get_logging"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeStringToStringMap(method+" -> ", result)
	return
}

// GetResidentVMs: Get the resident_VMs field of the given host.
// Version: rio
func (host) GetResidentVMs(session *Session, self HostRef) (retval []VMRef, err error) {
	method := "host.get_resident_VMs"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeVMRefSet(method+" -> ", result)
	return
}

// GetResidentVMs2: Get the resident_VMs field of the given host.
// Version: rio
func (host) GetResidentVMs2(session *Session, self HostRef) (retval []VMRef, err error) {
	method := "host.get_resident_VMs"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeVMRefSet(method+" -> ", result)
	return
}

// GetSupportedBootloaders: Get the supported_bootloaders field of the given host.
// Version: rio
func (host) GetSupportedBootloaders(session *Session, self HostRef) (retval []string, err error) {
	method := "host.get_supported_bootloaders"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeStringSet(method+" -> ", result)
	return
}

// GetSupportedBootloaders2: Get the supported_bootloaders field of the given host.
// Version: rio
func (host) GetSupportedBootloaders2(session *Session, self HostRef) (retval []string, err error) {
	method := "host.get_supported_bootloaders"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeStringSet(method+" -> ", result)
	return
}

// GetSchedPolicy: Get the sched_policy field of the given host.
// Version: rio
func (host) GetSchedPolicy(session *Session, self HostRef) (retval string, err error) {
	method := "host.get_sched_policy"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeString(method+" -> ", result)
	return
}

// GetSchedPolicy2: Get the sched_policy field of the given host.
// Version: rio
func (host) GetSchedPolicy2(session *Session, self HostRef) (retval string, err error) {
	method := "host.get_sched_policy"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeString(method+" -> ", result)
	return
}

// GetCPUConfiguration: Get the cpu_configuration field of the given host.
// Version: rio
func (host) GetCPUConfiguration(session *Session, self HostRef) (retval map[string]string, err error) {
	method := "host.get_cpu_configuration"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeStringToStringMap(method+" -> ", result)
	return
}

// GetCPUConfiguration2: Get the cpu_configuration field of the given host.
// Version: rio
func (host) GetCPUConfiguration2(session *Session, self HostRef) (retval map[string]string, err error) {
	method := "host.get_cpu_configuration"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeStringToStringMap(method+" -> ", result)
	return
}

// GetCapabilities: Get the capabilities field of the given host.
// Version: rio
func (host) GetCapabilities(session *Session, self HostRef) (retval []string, err error) {
	method := "host.get_capabilities"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeStringSet(method+" -> ", result)
	return
}

// GetCapabilities2: Get the capabilities field of the given host.
// Version: rio
func (host) GetCapabilities2(session *Session, self HostRef) (retval []string, err error) {
	method := "host.get_capabilities"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeStringSet(method+" -> ", result)
	return
}

// GetOtherConfig: Get the other_config field of the given host.
// Version: rio
func (host) GetOtherConfig(session *Session, self HostRef) (retval map[string]string, err error) {
	method := "host.get_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeStringToStringMap(method+" -> ", result)
	return
}

// GetOtherConfig2: Get the other_config field of the given host.
// Version: rio
func (host) GetOtherConfig2(session *Session, self HostRef) (retval map[string]string, err error) {
	method := "host.get_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeStringToStringMap(method+" -> ", result)
	return
}

// GetSoftwareVersion: Get the software_version field of the given host.
// Version: rio
func (host) GetSoftwareVersion(session *Session, self HostRef) (retval map[string]string, err error) {
	method := "host.get_software_version"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeStringToStringMap(method+" -> ", result)
	return
}

// GetSoftwareVersion2: Get the software_version field of the given host.
// Version: rio
func (host) GetSoftwareVersion2(session *Session, self HostRef) (retval map[string]string, err error) {
	method := "host.get_software_version"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeStringToStringMap(method+" -> ", result)
	return
}

// GetEnabled: Get the enabled field of the given host.
// Version: rio
func (host) GetEnabled(session *Session, self HostRef) (retval bool, err error) {
	method := "host.get_enabled"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeBool(method+" -> ", result)
	return
}

// GetEnabled2: Get the enabled field of the given host.
// Version: rio
func (host) GetEnabled2(session *Session, self HostRef) (retval bool, err error) {
	method := "host.get_enabled"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeBool(method+" -> ", result)
	return
}

// GetAPIVersionVendorImplementation: Get the API_version/vendor_implementation field of the given host.
// Version: rio
func (host) GetAPIVersionVendorImplementation(session *Session, self HostRef) (retval map[string]string, err error) {
	method := "host.get_API_version_vendor_implementation"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeStringToStringMap(method+" -> ", result)
	return
}

// GetAPIVersionVendorImplementation2: Get the API_version/vendor_implementation field of the given host.
// Version: rio
func (host) GetAPIVersionVendorImplementation2(session *Session, self HostRef) (retval map[string]string, err error) {
	method := "host.get_API_version_vendor_implementation"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeStringToStringMap(method+" -> ", result)
	return
}

// GetAPIVersionVendor: Get the API_version/vendor field of the given host.
// Version: rio
func (host) GetAPIVersionVendor(session *Session, self HostRef) (retval string, err error) {
	method := "host.get_API_version_vendor"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeString(method+" -> ", result)
	return
}

// GetAPIVersionVendor2: Get the API_version/vendor field of the given host.
// Version: rio
func (host) GetAPIVersionVendor2(session *Session, self HostRef) (retval string, err error) {
	method := "host.get_API_version_vendor"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeString(method+" -> ", result)
	return
}

// GetAPIVersionMinor: Get the API_version/minor field of the given host.
// Version: rio
func (host) GetAPIVersionMinor(session *Session, self HostRef) (retval int, err error) {
	method := "host.get_API_version_minor"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeInt(method+" -> ", result)
	return
}

// GetAPIVersionMinor2: Get the API_version/minor field of the given host.
// Version: rio
func (host) GetAPIVersionMinor2(session *Session, self HostRef) (retval int, err error) {
	method := "host.get_API_version_minor"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeInt(method+" -> ", result)
	return
}

// GetAPIVersionMajor: Get the API_version/major field of the given host.
// Version: rio
func (host) GetAPIVersionMajor(session *Session, self HostRef) (retval int, err error) {
	method := "host.get_API_version_major"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeInt(method+" -> ", result)
	return
}

// GetAPIVersionMajor2: Get the API_version/major field of the given host.
// Version: rio
func (host) GetAPIVersionMajor2(session *Session, self HostRef) (retval int, err error) {
	method := "host.get_API_version_major"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeInt(method+" -> ", result)
	return
}

// GetCurrentOperations: Get the current_operations field of the given host.
// Version: rio
func (host) GetCurrentOperations(session *Session, self HostRef) (retval map[string]HostAllowedOperations, err error) {
	method := "host.get_current_operations"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeStringToEnumHostAllowedOperationsMap(method+" -> ", result)
	return
}

// GetCurrentOperations2: Get the current_operations field of the given host.
// Version: rio
func (host) GetCurrentOperations2(session *Session, self HostRef) (retval map[string]HostAllowedOperations, err error) {
	method := "host.get_current_operations"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeStringToEnumHostAllowedOperationsMap(method+" -> ", result)
	return
}

// GetAllowedOperations: Get the allowed_operations field of the given host.
// Version: rio
func (host) GetAllowedOperations(session *Session, self HostRef) (retval []HostAllowedOperations, err error) {
	method := "host.get_allowed_operations"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeEnumHostAllowedOperationsSet(method+" -> ", result)
	return
}

// GetAllowedOperations2: Get the allowed_operations field of the given host.
// Version: rio
func (host) GetAllowedOperations2(session *Session, self HostRef) (retval []HostAllowedOperations, err error) {
	method := "host.get_allowed_operations"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeEnumHostAllowedOperationsSet(method+" -> ", result)
	return
}

// GetMemoryOverhead: Get the memory/overhead field of the given host.
// Version: rio
func (host) GetMemoryOverhead(session *Session, self HostRef) (retval int, err error) {
	method := "host.get_memory_overhead"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeInt(method+" -> ", result)
	return
}

// GetMemoryOverhead2: Get the memory/overhead field of the given host.
// Version: rio
func (host) GetMemoryOverhead2(session *Session, self HostRef) (retval int, err error) {
	method := "host.get_memory_overhead"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeInt(method+" -> ", result)
	return
}

// GetNameDescription: Get the name/description field of the given host.
// Version: rio
func (host) GetNameDescription(session *Session, self HostRef) (retval string, err error) {
	method := "host.get_name_description"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeString(method+" -> ", result)
	return
}

// GetNameDescription2: Get the name/description field of the given host.
// Version: rio
func (host) GetNameDescription2(session *Session, self HostRef) (retval string, err error) {
	method := "host.get_name_description"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeString(method+" -> ", result)
	return
}

// GetNameLabel: Get the name/label field of the given host.
// Version: rio
func (host) GetNameLabel(session *Session, self HostRef) (retval string, err error) {
	method := "host.get_name_label"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeString(method+" -> ", result)
	return
}

// GetNameLabel2: Get the name/label field of the given host.
// Version: rio
func (host) GetNameLabel2(session *Session, self HostRef) (retval string, err error) {
	method := "host.get_name_label"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeString(method+" -> ", result)
	return
}

// GetUUID: Get the uuid field of the given host.
// Version: rio
func (host) GetUUID(session *Session, self HostRef) (retval string, err error) {
	method := "host.get_uuid"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeString(method+" -> ", result)
	return
}

// GetUUID2: Get the uuid field of the given host.
// Version: rio
func (host) GetUUID2(session *Session, self HostRef) (retval string, err error) {
	method := "host.get_uuid"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeString(method+" -> ", result)
	return
}

// GetByNameLabel: Get all the host instances with the given label.
// Version: rio
func (host) GetByNameLabel(session *Session, label string) (retval []HostRef, err error) {
	method := "host.get_by_name_label"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	labelArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "label"), label)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, labelArg)
	if err != nil {
		return
	}
	retval, err = deserializeHostRefSet(method+" -> ", result)
	return
}

// GetByNameLabel2: Get all the host instances with the given label.
// Version: rio
func (host) GetByNameLabel2(session *Session, label string) (retval []HostRef, err error) {
	method := "host.get_by_name_label"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	labelArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "label"), label)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, labelArg)
	if err != nil {
		return
	}
	retval, err = deserializeHostRefSet(method+" -> ", result)
	return
}

// GetByUUID: Get a reference to the host instance with the specified UUID.
// Version: rio
func (host) GetByUUID(session *Session, uuid string) (retval HostRef, err error) {
	method := "host.get_by_uuid"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	uuidArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "uuid"), uuid)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, uuidArg)
	if err != nil {
		return
	}
	retval, err = deserializeHostRef(method+" -> ", result)
	return
}

// GetByUUID2: Get a reference to the host instance with the specified UUID.
// Version: rio
func (host) GetByUUID2(session *Session, uuid string) (retval HostRef, err error) {
	method := "host.get_by_uuid"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	uuidArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "uuid"), uuid)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, uuidArg)
	if err != nil {
		return
	}
	retval, err = deserializeHostRef(method+" -> ", result)
	return
}

// GetRecord: Get a record containing the current state of the given host.
// Version: rio
func (host) GetRecord(session *Session, self HostRef) (retval HostRecord, err error) {
	method := "host.get_record"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeHostRecord(method+" -> ", result)
	return
}

// GetRecord2: Get a record containing the current state of the given host.
// Version: rio
func (host) GetRecord2(session *Session, self HostRef) (retval HostRecord, err error) {
	method := "host.get_record"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeHostRecord(method+" -> ", result)
	return
}
