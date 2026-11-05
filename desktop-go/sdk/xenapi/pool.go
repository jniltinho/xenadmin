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

type PoolRecord struct {
	// Unique identifier/object reference
	UUID string `json:"uuid,omitempty"`
	// Short name
	NameLabel string `json:"name_label,omitempty"`
	// Description
	NameDescription string `json:"name_description,omitempty"`
	// The host that is pool master
	Master HostRef `json:"master,omitempty"`
	// Default SR for VDIs
	DefaultSR SRRef `json:"default_SR,omitempty"`
	// The SR in which VDIs for suspend images are created
	SuspendImageSR SRRef `json:"suspend_image_SR,omitempty"`
	// The SR in which VDIs for crash dumps are created
	CrashDumpSR SRRef `json:"crash_dump_SR,omitempty"`
	// additional configuration
	OtherConfig map[string]string `json:"other_config,omitempty"`
	// true if HA is enabled on the pool, false otherwise
	HaEnabled bool `json:"ha_enabled,omitempty"`
	// The current HA configuration
	HaConfiguration map[string]string `json:"ha_configuration,omitempty"`
	// HA statefile VDIs in use
	HaStatefiles []string `json:"ha_statefiles,omitempty"`
	// Number of host failures to tolerate before the Pool is declared to be overcommitted
	HaHostFailuresToTolerate int `json:"ha_host_failures_to_tolerate,omitempty"`
	// Number of future host failures we have managed to find a plan for. Once this reaches zero any future host failures will cause the failure of protected VMs.
	HaPlanExistsFor int `json:"ha_plan_exists_for,omitempty"`
	// If set to false then operations which would cause the Pool to become overcommitted will be blocked.
	HaAllowOvercommit bool `json:"ha_allow_overcommit,omitempty"`
	// True if the Pool is considered to be overcommitted i.e. if there exist insufficient physical resources to tolerate the configured number of host failures
	HaOvercommitted bool `json:"ha_overcommitted,omitempty"`
	// Binary blobs associated with this pool
	Blobs map[string]BlobRef `json:"blobs,omitempty"`
	// user-specified tags for categorization purposes
	Tags []string `json:"tags,omitempty"`
	// gui-specific configuration for pool
	GuiConfig map[string]string `json:"gui_config,omitempty"`
	// Configuration for the automatic health check feature
	HealthCheckConfig map[string]string `json:"health_check_config,omitempty"`
	// Url for the configured workload balancing host
	WlbURL string `json:"wlb_url,omitempty"`
	// Username for accessing the workload balancing host
	WlbUsername string `json:"wlb_username,omitempty"`
	// true if workload balancing is enabled on the pool, false otherwise
	WlbEnabled bool `json:"wlb_enabled,omitempty"`
	// true if communication with the WLB server should enforce TLS certificate verification.
	WlbVerifyCert bool `json:"wlb_verify_cert,omitempty"`
	// true a redo-log is to be used other than when HA is enabled, false otherwise
	RedoLogEnabled bool `json:"redo_log_enabled,omitempty"`
	// indicates the VDI to use for the redo-log other than when HA is enabled
	RedoLogVdi VDIRef `json:"redo_log_vdi,omitempty"`
	// address of the vswitch controller
	VswitchController string `json:"vswitch_controller,omitempty"`
	// Pool-wide restrictions currently in effect
	Restrictions map[string]string `json:"restrictions,omitempty"`
	// The set of currently known metadata VDIs for this pool
	MetadataVDIs []VDIRef `json:"metadata_VDIs,omitempty"`
	// The HA cluster stack that is currently in use. Only valid when HA is enabled.
	HaClusterStack string `json:"ha_cluster_stack,omitempty"`
	// list of the operations allowed in this state. This list is advisory only and the server state may have changed by the time this field is read by a client.
	AllowedOperations []PoolAllowedOperations `json:"allowed_operations,omitempty"`
	// links each of the running tasks using this object (by reference) to a current_operation enum which describes the nature of the task.
	CurrentOperations map[string]PoolAllowedOperations `json:"current_operations,omitempty"`
	// Pool-wide guest agent configuration information
	GuestAgentConfig map[string]string `json:"guest_agent_config,omitempty"`
	// Details about the physical CPUs on the pool
	CPUInfo map[string]string `json:"cpu_info,omitempty"`
	// This field was consulted when VM.create did not specify a value for &apos;has_vendor_device&apos;; VM.create now uses a simple default and no longer consults this value.
	PolicyNoVendorDevice bool `json:"policy_no_vendor_device,omitempty"`
	// The pool-wide flag to show if the live patching feauture is disabled or not.
	LivePatchingDisabled bool `json:"live_patching_disabled,omitempty"`
	// true if IGMP snooping is enabled in the pool, false otherwise.
	IgmpSnoopingEnabled bool `json:"igmp_snooping_enabled,omitempty"`
	// The UEFI certificates allowing Secure Boot
	UefiCertificates string `json:"uefi_certificates,omitempty"`
	// Custom UEFI certificates allowing Secure Boot
	CustomUefiCertificates string `json:"custom_uefi_certificates,omitempty"`
	// True if either a PSR is running or we are waiting for a PSR to be re-run
	IsPsrPending bool `json:"is_psr_pending,omitempty"`
	// True iff TLS certificate verification is enabled
	TLSVerificationEnabled bool `json:"tls_verification_enabled,omitempty"`
	// The set of currently enabled repositories
	Repositories []RepositoryRef `json:"repositories,omitempty"`
	// True if authentication by TLS client certificates is enabled
	ClientCertificateAuthEnabled bool `json:"client_certificate_auth_enabled,omitempty"`
	// The name (CN/SAN) that an incoming client certificate must have to allow authentication
	ClientCertificateAuthName string `json:"client_certificate_auth_name,omitempty"`
	// Url of the proxy used in syncing with the enabled repositories
	RepositoryProxyURL string `json:"repository_proxy_url,omitempty"`
	// Username for the authentication of the proxy used in syncing with the enabled repositories
	RepositoryProxyUsername string `json:"repository_proxy_username,omitempty"`
	// Password for the authentication of the proxy used in syncing with the enabled repositories
	RepositoryProxyPassword SecretRef `json:"repository_proxy_password,omitempty"`
	// Default behaviour during migration, True if stream compression should be used
	MigrationCompression bool `json:"migration_compression,omitempty"`
	// true if bias against pool master when scheduling vms is enabled, false otherwise
	CoordinatorBias bool `json:"coordinator_bias,omitempty"`
	// Maximum number of threads to use for PAM authentication
	LocalAuthMaxThreads int `json:"local_auth_max_threads,omitempty"`
	// Maximum number of threads to use for external (AD) authentication
	ExtAuthMaxThreads int `json:"ext_auth_max_threads,omitempty"`
	// Specifies whether external authentication caching is enabled for this pool or not
	ExtAuthCacheEnabled bool `json:"ext_auth_cache_enabled,omitempty"`
	// Maximum capacity of external authentication cache
	ExtAuthCacheSize int `json:"ext_auth_cache_size,omitempty"`
	// Specifies how long external authentication entries should be cached for (seconds)
	ExtAuthCacheExpiry int `json:"ext_auth_cache_expiry,omitempty"`
	// The UUID of the pool for identification of telemetry data
	TelemetryUUID SecretRef `json:"telemetry_uuid,omitempty"`
	// How often the telemetry collection will be carried out
	TelemetryFrequency TelemetryFrequency `json:"telemetry_frequency,omitempty"`
	// The earliest timestamp (in UTC) when the next round of telemetry collection can be carried out
	TelemetryNextCollection time.Time `json:"telemetry_next_collection,omitempty"`
	// time of the last update sychronization
	LastUpdateSync time.Time `json:"last_update_sync,omitempty"`
	// The frequency at which updates are synchronized from a remote CDN: daily or weekly.
	UpdateSyncFrequency UpdateSyncFrequency `json:"update_sync_frequency,omitempty"`
	// The day of the week the update synchronizations will be scheduled, based on pool&apos;s local timezone. Ignored when update_sync_frequency is daily
	UpdateSyncDay int `json:"update_sync_day,omitempty"`
	// Whether periodic update synchronization is enabled or not
	UpdateSyncEnabled bool `json:"update_sync_enabled,omitempty"`
	// The recommended pool properties for clients to respect for optimal performance. For example, max-vm-group=5
	Recommendations map[string]string `json:"recommendations,omitempty"`
	// Licensing data shared within the whole pool
	LicenseServer map[string]string `json:"license_server,omitempty"`
	// Indicates whether an HA-protected VM that is shut down from inside (not through the API) should be automatically rebooted when HA is enabled
	HaRebootVMOnInternalShutdown bool `json:"ha_reboot_vm_on_internal_shutdown,omitempty"`
	// When true, only one console connection per VM/host in the pool is accepted. Otherwise every connection for a VM/host&apos;s console is accepted. Note: when true, connection attempts via websocket will be rejected.
	LimitConsoleSessions bool `json:"limit_console_sessions,omitempty"`
	// The maximum time (in seconds) that a VM&apos;s console can be idle before it is automatically disconnected. The default value 0 means never timeout. This setting applies only to VM consoles; for host consoles, use the separate parameter &apos;host.console_idle_timeout&apos;.
	VMConsoleIdleTimeout int `json:"vm_console_idle_timeout,omitempty"`
	// When true, at VM.create time the pool automatically marks a newly created VM whose Secure Boot certificates are due to expire by setting VM.secureboot_certificates_state to update_on_boot (the same effect as VM.update_secureboot_certificates_on_boot), so the certificates are updated on the VM&apos;s next boot. This applies to VMs created by VM.create, including those created by operations that internally call it: the HTTP(s) PUT /import and /import_metadata calls and disaster recovery (VM.recover and VM_appliance.recover); it does not apply to VM.clone or VM.copy, which inherit the source VM&apos;s state.
	AutoUpdateVMSecurebootCertificates bool `json:"auto_update_vm_secureboot_certificates,omitempty"`
}

type PoolRef string

// Pool-wide information
type pool struct{}

var Pool pool

// GetAllRecords: Return a map of pool references to pool records for all pools known to the system.
// Version: rio
func (pool) GetAllRecords(session *Session) (retval map[PoolRef]PoolRecord, err error) {
	method := "pool.get_all_records"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializePoolRefToPoolRecordMap(method+" -> ", result)
	return
}

// GetAllRecords1: Return a map of pool references to pool records for all pools known to the system.
// Version: rio
func (pool) GetAllRecords1(session *Session) (retval map[PoolRef]PoolRecord, err error) {
	method := "pool.get_all_records"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializePoolRefToPoolRecordMap(method+" -> ", result)
	return
}

// GetAll: Return a list of all the pools known to the system.
// Version: rio
func (pool) GetAll(session *Session) (retval []PoolRef, err error) {
	method := "pool.get_all"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializePoolRefSet(method+" -> ", result)
	return
}

// GetAll1: Return a list of all the pools known to the system.
// Version: rio
func (pool) GetAll1(session *Session) (retval []PoolRef, err error) {
	method := "pool.get_all"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializePoolRefSet(method+" -> ", result)
	return
}

// SyncTrustedCertificatesFrom: Download trusted TLS certificates from a remote pool and install them in this pool. Certificates already present locally (matched by fingerprint and purpose) are skipped.
// Version: closed
func (pool) SyncTrustedCertificatesFrom(session *Session, self PoolRef, remotePool string, remoteSession SessionRef, remoteCertificate string, ca bool) (retval []CertificateRef, err error) {
	method := "pool.sync_trusted_certificates_from"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	remotePoolArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "remote_pool"), remotePool)
	if err != nil {
		return
	}
	remoteSessionArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "remote_session"), remoteSession)
	if err != nil {
		return
	}
	remoteCertificateArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "remote_certificate"), remoteCertificate)
	if err != nil {
		return
	}
	caArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "ca"), ca)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg, remotePoolArg, remoteSessionArg, remoteCertificateArg, caArg)
	if err != nil {
		return
	}
	retval, err = deserializeCertificateRefSet(method+" -> ", result)
	return
}

// AsyncSyncTrustedCertificatesFrom: Download trusted TLS certificates from a remote pool and install them in this pool. Certificates already present locally (matched by fingerprint and purpose) are skipped.
// Version: closed
func (pool) AsyncSyncTrustedCertificatesFrom(session *Session, self PoolRef, remotePool string, remoteSession SessionRef, remoteCertificate string, ca bool) (retval TaskRef, err error) {
	method := "Async.pool.sync_trusted_certificates_from"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	remotePoolArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "remote_pool"), remotePool)
	if err != nil {
		return
	}
	remoteSessionArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "remote_session"), remoteSession)
	if err != nil {
		return
	}
	remoteCertificateArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "remote_certificate"), remoteCertificate)
	if err != nil {
		return
	}
	caArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "ca"), ca)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg, remotePoolArg, remoteSessionArg, remoteCertificateArg, caArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// SyncTrustedCertificatesFrom6: Download trusted TLS certificates from a remote pool and install them in this pool. Certificates already present locally (matched by fingerprint and purpose) are skipped.
// Version: closed
func (pool) SyncTrustedCertificatesFrom6(session *Session, self PoolRef, remotePool string, remoteSession SessionRef, remoteCertificate string, ca bool) (retval []CertificateRef, err error) {
	method := "pool.sync_trusted_certificates_from"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	remotePoolArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "remote_pool"), remotePool)
	if err != nil {
		return
	}
	remoteSessionArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "remote_session"), remoteSession)
	if err != nil {
		return
	}
	remoteCertificateArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "remote_certificate"), remoteCertificate)
	if err != nil {
		return
	}
	caArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "ca"), ca)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg, remotePoolArg, remoteSessionArg, remoteCertificateArg, caArg)
	if err != nil {
		return
	}
	retval, err = deserializeCertificateRefSet(method+" -> ", result)
	return
}

// AsyncSyncTrustedCertificatesFrom6: Download trusted TLS certificates from a remote pool and install them in this pool. Certificates already present locally (matched by fingerprint and purpose) are skipped.
// Version: closed
func (pool) AsyncSyncTrustedCertificatesFrom6(session *Session, self PoolRef, remotePool string, remoteSession SessionRef, remoteCertificate string, ca bool) (retval TaskRef, err error) {
	method := "Async.pool.sync_trusted_certificates_from"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	remotePoolArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "remote_pool"), remotePool)
	if err != nil {
		return
	}
	remoteSessionArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "remote_session"), remoteSession)
	if err != nil {
		return
	}
	remoteCertificateArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "remote_certificate"), remoteCertificate)
	if err != nil {
		return
	}
	caArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "ca"), ca)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg, remotePoolArg, remoteSessionArg, remoteCertificateArg, caArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// UninstallTrustedCertificate: Uninstall a trusted TLS certificate, pool-wide.
// Version: closed
func (pool) UninstallTrustedCertificate(session *Session, self PoolRef, certificate CertificateRef) (err error) {
	method := "pool.uninstall_trusted_certificate"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	certificateArg, err := serializeCertificateRef(fmt.Sprintf("%s(%s)", method, "certificate"), certificate)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, certificateArg)
	return
}

// AsyncUninstallTrustedCertificate: Uninstall a trusted TLS certificate, pool-wide.
// Version: closed
func (pool) AsyncUninstallTrustedCertificate(session *Session, self PoolRef, certificate CertificateRef) (retval TaskRef, err error) {
	method := "Async.pool.uninstall_trusted_certificate"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	certificateArg, err := serializeCertificateRef(fmt.Sprintf("%s(%s)", method, "certificate"), certificate)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg, certificateArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// UninstallTrustedCertificate3: Uninstall a trusted TLS certificate, pool-wide.
// Version: closed
func (pool) UninstallTrustedCertificate3(session *Session, self PoolRef, certificate CertificateRef) (err error) {
	method := "pool.uninstall_trusted_certificate"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	certificateArg, err := serializeCertificateRef(fmt.Sprintf("%s(%s)", method, "certificate"), certificate)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, certificateArg)
	return
}

// AsyncUninstallTrustedCertificate3: Uninstall a trusted TLS certificate, pool-wide.
// Version: closed
func (pool) AsyncUninstallTrustedCertificate3(session *Session, self PoolRef, certificate CertificateRef) (retval TaskRef, err error) {
	method := "Async.pool.uninstall_trusted_certificate"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	certificateArg, err := serializeCertificateRef(fmt.Sprintf("%s(%s)", method, "certificate"), certificate)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg, certificateArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// InstallTrustedCertificate: Install a trusted TLS certificate, pool-wide.
// Version: closed
func (pool) InstallTrustedCertificate(session *Session, self PoolRef, ca bool, cert string, purpose []CertificatePurpose) (err error) {
	method := "pool.install_trusted_certificate"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	caArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "ca"), ca)
	if err != nil {
		return
	}
	certArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "cert"), cert)
	if err != nil {
		return
	}
	purposeArg, err := serializeEnumCertificatePurposeSet(fmt.Sprintf("%s(%s)", method, "purpose"), purpose)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, caArg, certArg, purposeArg)
	return
}

// AsyncInstallTrustedCertificate: Install a trusted TLS certificate, pool-wide.
// Version: closed
func (pool) AsyncInstallTrustedCertificate(session *Session, self PoolRef, ca bool, cert string, purpose []CertificatePurpose) (retval TaskRef, err error) {
	method := "Async.pool.install_trusted_certificate"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	caArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "ca"), ca)
	if err != nil {
		return
	}
	certArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "cert"), cert)
	if err != nil {
		return
	}
	purposeArg, err := serializeEnumCertificatePurposeSet(fmt.Sprintf("%s(%s)", method, "purpose"), purpose)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg, caArg, certArg, purposeArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// InstallTrustedCertificate5: Install a trusted TLS certificate, pool-wide.
// Version: closed
func (pool) InstallTrustedCertificate5(session *Session, self PoolRef, ca bool, cert string, purpose []CertificatePurpose) (err error) {
	method := "pool.install_trusted_certificate"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	caArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "ca"), ca)
	if err != nil {
		return
	}
	certArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "cert"), cert)
	if err != nil {
		return
	}
	purposeArg, err := serializeEnumCertificatePurposeSet(fmt.Sprintf("%s(%s)", method, "purpose"), purpose)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, caArg, certArg, purposeArg)
	return
}

// AsyncInstallTrustedCertificate5: Install a trusted TLS certificate, pool-wide.
// Version: closed
func (pool) AsyncInstallTrustedCertificate5(session *Session, self PoolRef, ca bool, cert string, purpose []CertificatePurpose) (retval TaskRef, err error) {
	method := "Async.pool.install_trusted_certificate"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	caArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "ca"), ca)
	if err != nil {
		return
	}
	certArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "cert"), cert)
	if err != nil {
		return
	}
	purposeArg, err := serializeEnumCertificatePurposeSet(fmt.Sprintf("%s(%s)", method, "purpose"), purpose)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg, caArg, certArg, purposeArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// SetSSHAutoMode: Set the SSH auto mode for all hosts in the pool
// Version: closed
func (pool) SetSSHAutoMode(session *Session, self PoolRef, value bool) (err error) {
	method := "pool.set_ssh_auto_mode"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// AsyncSetSSHAutoMode: Set the SSH auto mode for all hosts in the pool
// Version: closed
func (pool) AsyncSetSSHAutoMode(session *Session, self PoolRef, value bool) (retval TaskRef, err error) {
	method := "Async.pool.set_ssh_auto_mode"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetSSHAutoMode3: Set the SSH auto mode for all hosts in the pool
// Version: closed
func (pool) SetSSHAutoMode3(session *Session, self PoolRef, value bool) (err error) {
	method := "pool.set_ssh_auto_mode"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// AsyncSetSSHAutoMode3: Set the SSH auto mode for all hosts in the pool
// Version: closed
func (pool) AsyncSetSSHAutoMode3(session *Session, self PoolRef, value bool) (retval TaskRef, err error) {
	method := "Async.pool.set_ssh_auto_mode"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetConsoleIdleTimeout: Set the console idle timeout for all hosts in the pool
// Version: closed
func (pool) SetConsoleIdleTimeout(session *Session, self PoolRef, value int) (err error) {
	method := "pool.set_console_idle_timeout"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// AsyncSetConsoleIdleTimeout: Set the console idle timeout for all hosts in the pool
// Version: closed
func (pool) AsyncSetConsoleIdleTimeout(session *Session, self PoolRef, value int) (retval TaskRef, err error) {
	method := "Async.pool.set_console_idle_timeout"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetConsoleIdleTimeout3: Set the console idle timeout for all hosts in the pool
// Version: closed
func (pool) SetConsoleIdleTimeout3(session *Session, self PoolRef, value int) (err error) {
	method := "pool.set_console_idle_timeout"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// AsyncSetConsoleIdleTimeout3: Set the console idle timeout for all hosts in the pool
// Version: closed
func (pool) AsyncSetConsoleIdleTimeout3(session *Session, self PoolRef, value int) (retval TaskRef, err error) {
	method := "Async.pool.set_console_idle_timeout"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetSSHEnabledTimeout: Set the SSH enabled timeout for all hosts in the pool
// Version: closed
func (pool) SetSSHEnabledTimeout(session *Session, self PoolRef, value int) (err error) {
	method := "pool.set_ssh_enabled_timeout"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// AsyncSetSSHEnabledTimeout: Set the SSH enabled timeout for all hosts in the pool
// Version: closed
func (pool) AsyncSetSSHEnabledTimeout(session *Session, self PoolRef, value int) (retval TaskRef, err error) {
	method := "Async.pool.set_ssh_enabled_timeout"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetSSHEnabledTimeout3: Set the SSH enabled timeout for all hosts in the pool
// Version: closed
func (pool) SetSSHEnabledTimeout3(session *Session, self PoolRef, value int) (err error) {
	method := "pool.set_ssh_enabled_timeout"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// AsyncSetSSHEnabledTimeout3: Set the SSH enabled timeout for all hosts in the pool
// Version: closed
func (pool) AsyncSetSSHEnabledTimeout3(session *Session, self PoolRef, value int) (retval TaskRef, err error) {
	method := "Async.pool.set_ssh_enabled_timeout"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// DisableSSH: Disable SSH access on all hosts in the pool. It&apos;s a helper which calls host.disable_ssh for all the hosts in the pool.
// Version: closed
func (pool) DisableSSH(session *Session, self PoolRef) (err error) {
	method := "pool.disable_ssh"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// AsyncDisableSSH: Disable SSH access on all hosts in the pool. It&apos;s a helper which calls host.disable_ssh for all the hosts in the pool.
// Version: closed
func (pool) AsyncDisableSSH(session *Session, self PoolRef) (retval TaskRef, err error) {
	method := "Async.pool.disable_ssh"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// DisableSSH2: Disable SSH access on all hosts in the pool. It&apos;s a helper which calls host.disable_ssh for all the hosts in the pool.
// Version: closed
func (pool) DisableSSH2(session *Session, self PoolRef) (err error) {
	method := "pool.disable_ssh"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// AsyncDisableSSH2: Disable SSH access on all hosts in the pool. It&apos;s a helper which calls host.disable_ssh for all the hosts in the pool.
// Version: closed
func (pool) AsyncDisableSSH2(session *Session, self PoolRef) (retval TaskRef, err error) {
	method := "Async.pool.disable_ssh"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// EnableSSH: Enable SSH access on all hosts in the pool. It&apos;s a helper which calls host.enable_ssh for all the hosts in the pool.
// Version: closed
func (pool) EnableSSH(session *Session, self PoolRef) (err error) {
	method := "pool.enable_ssh"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// AsyncEnableSSH: Enable SSH access on all hosts in the pool. It&apos;s a helper which calls host.enable_ssh for all the hosts in the pool.
// Version: closed
func (pool) AsyncEnableSSH(session *Session, self PoolRef) (retval TaskRef, err error) {
	method := "Async.pool.enable_ssh"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// EnableSSH2: Enable SSH access on all hosts in the pool. It&apos;s a helper which calls host.enable_ssh for all the hosts in the pool.
// Version: closed
func (pool) EnableSSH2(session *Session, self PoolRef) (err error) {
	method := "pool.enable_ssh"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// AsyncEnableSSH2: Enable SSH access on all hosts in the pool. It&apos;s a helper which calls host.enable_ssh for all the hosts in the pool.
// Version: closed
func (pool) AsyncEnableSSH2(session *Session, self PoolRef) (retval TaskRef, err error) {
	method := "Async.pool.enable_ssh"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetGuestSecurebootReadiness:
// Version: closed
func (pool) GetGuestSecurebootReadiness(session *Session, self PoolRef) (retval PoolGuestSecurebootReadiness, err error) {
	method := "pool.get_guest_secureboot_readiness"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeEnumPoolGuestSecurebootReadiness(method+" -> ", result)
	return
}

// GetGuestSecurebootReadiness2:
// Version: closed
func (pool) GetGuestSecurebootReadiness2(session *Session, self PoolRef) (retval PoolGuestSecurebootReadiness, err error) {
	method := "pool.get_guest_secureboot_readiness"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeEnumPoolGuestSecurebootReadiness(method+" -> ", result)
	return
}

// SetExtAuthCacheExpiry:
// Version: closed
func (pool) SetExtAuthCacheExpiry(session *Session, self PoolRef, value int) (err error) {
	method := "pool.set_ext_auth_cache_expiry"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// AsyncSetExtAuthCacheExpiry:
// Version: closed
func (pool) AsyncSetExtAuthCacheExpiry(session *Session, self PoolRef, value int) (retval TaskRef, err error) {
	method := "Async.pool.set_ext_auth_cache_expiry"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetExtAuthCacheExpiry3:
// Version: closed
func (pool) SetExtAuthCacheExpiry3(session *Session, self PoolRef, value int) (err error) {
	method := "pool.set_ext_auth_cache_expiry"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// AsyncSetExtAuthCacheExpiry3:
// Version: closed
func (pool) AsyncSetExtAuthCacheExpiry3(session *Session, self PoolRef, value int) (retval TaskRef, err error) {
	method := "Async.pool.set_ext_auth_cache_expiry"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetExtAuthCacheSize:
// Version: closed
func (pool) SetExtAuthCacheSize(session *Session, self PoolRef, value int) (err error) {
	method := "pool.set_ext_auth_cache_size"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// AsyncSetExtAuthCacheSize:
// Version: closed
func (pool) AsyncSetExtAuthCacheSize(session *Session, self PoolRef, value int) (retval TaskRef, err error) {
	method := "Async.pool.set_ext_auth_cache_size"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetExtAuthCacheSize3:
// Version: closed
func (pool) SetExtAuthCacheSize3(session *Session, self PoolRef, value int) (err error) {
	method := "pool.set_ext_auth_cache_size"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// AsyncSetExtAuthCacheSize3:
// Version: closed
func (pool) AsyncSetExtAuthCacheSize3(session *Session, self PoolRef, value int) (retval TaskRef, err error) {
	method := "Async.pool.set_ext_auth_cache_size"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetExtAuthCacheEnabled:
// Version: closed
func (pool) SetExtAuthCacheEnabled(session *Session, self PoolRef, value bool) (err error) {
	method := "pool.set_ext_auth_cache_enabled"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// AsyncSetExtAuthCacheEnabled:
// Version: closed
func (pool) AsyncSetExtAuthCacheEnabled(session *Session, self PoolRef, value bool) (retval TaskRef, err error) {
	method := "Async.pool.set_ext_auth_cache_enabled"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetExtAuthCacheEnabled3:
// Version: closed
func (pool) SetExtAuthCacheEnabled3(session *Session, self PoolRef, value bool) (err error) {
	method := "pool.set_ext_auth_cache_enabled"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// AsyncSetExtAuthCacheEnabled3:
// Version: closed
func (pool) AsyncSetExtAuthCacheEnabled3(session *Session, self PoolRef, value bool) (retval TaskRef, err error) {
	method := "Async.pool.set_ext_auth_cache_enabled"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetExtAuthMaxThreads:
// Version: closed
func (pool) SetExtAuthMaxThreads(session *Session, self PoolRef, value int) (err error) {
	method := "pool.set_ext_auth_max_threads"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetExtAuthMaxThreads3:
// Version: closed
func (pool) SetExtAuthMaxThreads3(session *Session, self PoolRef, value int) (err error) {
	method := "pool.set_ext_auth_max_threads"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetLocalAuthMaxThreads:
// Version: closed
func (pool) SetLocalAuthMaxThreads(session *Session, self PoolRef, value int) (err error) {
	method := "pool.set_local_auth_max_threads"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetLocalAuthMaxThreads3:
// Version: closed
func (pool) SetLocalAuthMaxThreads3(session *Session, self PoolRef, value int) (err error) {
	method := "pool.set_local_auth_max_threads"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetUpdateSyncEnabled: enable or disable periodic update synchronization depending on the value
// Version: closed
func (pool) SetUpdateSyncEnabled(session *Session, self PoolRef, value bool) (err error) {
	method := "pool.set_update_sync_enabled"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// AsyncSetUpdateSyncEnabled: enable or disable periodic update synchronization depending on the value
// Version: closed
func (pool) AsyncSetUpdateSyncEnabled(session *Session, self PoolRef, value bool) (retval TaskRef, err error) {
	method := "Async.pool.set_update_sync_enabled"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetUpdateSyncEnabled3: enable or disable periodic update synchronization depending on the value
// Version: closed
func (pool) SetUpdateSyncEnabled3(session *Session, self PoolRef, value bool) (err error) {
	method := "pool.set_update_sync_enabled"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// AsyncSetUpdateSyncEnabled3: enable or disable periodic update synchronization depending on the value
// Version: closed
func (pool) AsyncSetUpdateSyncEnabled3(session *Session, self PoolRef, value bool) (retval TaskRef, err error) {
	method := "Async.pool.set_update_sync_enabled"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// ConfigureUpdateSync: Configure periodic update synchronization to synchronize updates from a remote CDN
// Version: closed
func (pool) ConfigureUpdateSync(session *Session, self PoolRef, updateSyncFrequency UpdateSyncFrequency, updateSyncDay int) (err error) {
	method := "pool.configure_update_sync"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	updateSyncFrequencyArg, err := serializeEnumUpdateSyncFrequency(fmt.Sprintf("%s(%s)", method, "update_sync_frequency"), updateSyncFrequency)
	if err != nil {
		return
	}
	updateSyncDayArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "update_sync_day"), updateSyncDay)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, updateSyncFrequencyArg, updateSyncDayArg)
	return
}

// AsyncConfigureUpdateSync: Configure periodic update synchronization to synchronize updates from a remote CDN
// Version: closed
func (pool) AsyncConfigureUpdateSync(session *Session, self PoolRef, updateSyncFrequency UpdateSyncFrequency, updateSyncDay int) (retval TaskRef, err error) {
	method := "Async.pool.configure_update_sync"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	updateSyncFrequencyArg, err := serializeEnumUpdateSyncFrequency(fmt.Sprintf("%s(%s)", method, "update_sync_frequency"), updateSyncFrequency)
	if err != nil {
		return
	}
	updateSyncDayArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "update_sync_day"), updateSyncDay)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg, updateSyncFrequencyArg, updateSyncDayArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// ConfigureUpdateSync4: Configure periodic update synchronization to synchronize updates from a remote CDN
// Version: closed
func (pool) ConfigureUpdateSync4(session *Session, self PoolRef, updateSyncFrequency UpdateSyncFrequency, updateSyncDay int) (err error) {
	method := "pool.configure_update_sync"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	updateSyncFrequencyArg, err := serializeEnumUpdateSyncFrequency(fmt.Sprintf("%s(%s)", method, "update_sync_frequency"), updateSyncFrequency)
	if err != nil {
		return
	}
	updateSyncDayArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "update_sync_day"), updateSyncDay)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, updateSyncFrequencyArg, updateSyncDayArg)
	return
}

// AsyncConfigureUpdateSync4: Configure periodic update synchronization to synchronize updates from a remote CDN
// Version: closed
func (pool) AsyncConfigureUpdateSync4(session *Session, self PoolRef, updateSyncFrequency UpdateSyncFrequency, updateSyncDay int) (retval TaskRef, err error) {
	method := "Async.pool.configure_update_sync"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	updateSyncFrequencyArg, err := serializeEnumUpdateSyncFrequency(fmt.Sprintf("%s(%s)", method, "update_sync_frequency"), updateSyncFrequency)
	if err != nil {
		return
	}
	updateSyncDayArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "update_sync_day"), updateSyncDay)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg, updateSyncFrequencyArg, updateSyncDayArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// ResetTelemetryUUID: Assign a new UUID to telemetry data.
// Version: closed
func (pool) ResetTelemetryUUID(session *Session, self PoolRef) (err error) {
	method := "pool.reset_telemetry_uuid"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// AsyncResetTelemetryUUID: Assign a new UUID to telemetry data.
// Version: closed
func (pool) AsyncResetTelemetryUUID(session *Session, self PoolRef) (retval TaskRef, err error) {
	method := "Async.pool.reset_telemetry_uuid"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// ResetTelemetryUUID2: Assign a new UUID to telemetry data.
// Version: closed
func (pool) ResetTelemetryUUID2(session *Session, self PoolRef) (err error) {
	method := "pool.reset_telemetry_uuid"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// AsyncResetTelemetryUUID2: Assign a new UUID to telemetry data.
// Version: closed
func (pool) AsyncResetTelemetryUUID2(session *Session, self PoolRef) (retval TaskRef, err error) {
	method := "Async.pool.reset_telemetry_uuid"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetTelemetryNextCollection: Set the timestamp for the next telemetry data collection.
// Version: closed
func (pool) SetTelemetryNextCollection(session *Session, self PoolRef, value time.Time) (err error) {
	method := "pool.set_telemetry_next_collection"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// AsyncSetTelemetryNextCollection: Set the timestamp for the next telemetry data collection.
// Version: closed
func (pool) AsyncSetTelemetryNextCollection(session *Session, self PoolRef, value time.Time) (retval TaskRef, err error) {
	method := "Async.pool.set_telemetry_next_collection"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetTelemetryNextCollection3: Set the timestamp for the next telemetry data collection.
// Version: closed
func (pool) SetTelemetryNextCollection3(session *Session, self PoolRef, value time.Time) (err error) {
	method := "pool.set_telemetry_next_collection"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// AsyncSetTelemetryNextCollection3: Set the timestamp for the next telemetry data collection.
// Version: closed
func (pool) AsyncSetTelemetryNextCollection3(session *Session, self PoolRef, value time.Time) (retval TaskRef, err error) {
	method := "Async.pool.set_telemetry_next_collection"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetHTTPSOnly: updates all the host firewalls in the pool to open or close port 80 depending on the value
// Version: closed
func (pool) SetHTTPSOnly(session *Session, self PoolRef, value bool) (err error) {
	method := "pool.set_https_only"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// AsyncSetHTTPSOnly: updates all the host firewalls in the pool to open or close port 80 depending on the value
// Version: closed
func (pool) AsyncSetHTTPSOnly(session *Session, self PoolRef, value bool) (retval TaskRef, err error) {
	method := "Async.pool.set_https_only"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetHTTPSOnly3: updates all the host firewalls in the pool to open or close port 80 depending on the value
// Version: closed
func (pool) SetHTTPSOnly3(session *Session, self PoolRef, value bool) (err error) {
	method := "pool.set_https_only"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// AsyncSetHTTPSOnly3: updates all the host firewalls in the pool to open or close port 80 depending on the value
// Version: closed
func (pool) AsyncSetHTTPSOnly3(session *Session, self PoolRef, value bool) (retval TaskRef, err error) {
	method := "Async.pool.set_https_only"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetCustomUefiCertificates: Set custom UEFI certificates for a pool and all its hosts. Need `allow-custom-uefi-certs` set to true in conf. If empty: default back to Pool.uefi_certificates
// Version: closed
func (pool) SetCustomUefiCertificates(session *Session, self PoolRef, value string) (err error) {
	method := "pool.set_custom_uefi_certificates"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// AsyncSetCustomUefiCertificates: Set custom UEFI certificates for a pool and all its hosts. Need `allow-custom-uefi-certs` set to true in conf. If empty: default back to Pool.uefi_certificates
// Version: closed
func (pool) AsyncSetCustomUefiCertificates(session *Session, self PoolRef, value string) (retval TaskRef, err error) {
	method := "Async.pool.set_custom_uefi_certificates"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetCustomUefiCertificates3: Set custom UEFI certificates for a pool and all its hosts. Need `allow-custom-uefi-certs` set to true in conf. If empty: default back to Pool.uefi_certificates
// Version: closed
func (pool) SetCustomUefiCertificates3(session *Session, self PoolRef, value string) (err error) {
	method := "pool.set_custom_uefi_certificates"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// AsyncSetCustomUefiCertificates3: Set custom UEFI certificates for a pool and all its hosts. Need `allow-custom-uefi-certs` set to true in conf. If empty: default back to Pool.uefi_certificates
// Version: closed
func (pool) AsyncSetCustomUefiCertificates3(session *Session, self PoolRef, value string) (retval TaskRef, err error) {
	method := "Async.pool.set_custom_uefi_certificates"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetUefiCertificates: Set the UEFI certificates for a pool and all its hosts. Deprecated: use set_custom_uefi_certificates instead
// Version: 22.16.0
func (pool) SetUefiCertificates(session *Session, self PoolRef, value string) (err error) {
	method := "pool.set_uefi_certificates"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// AsyncSetUefiCertificates: Set the UEFI certificates for a pool and all its hosts. Deprecated: use set_custom_uefi_certificates instead
// Version: 22.16.0
func (pool) AsyncSetUefiCertificates(session *Session, self PoolRef, value string) (retval TaskRef, err error) {
	method := "Async.pool.set_uefi_certificates"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetUefiCertificates3: Set the UEFI certificates for a pool and all its hosts. Deprecated: use set_custom_uefi_certificates instead
// Version: 22.16.0
func (pool) SetUefiCertificates3(session *Session, self PoolRef, value string) (err error) {
	method := "pool.set_uefi_certificates"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// AsyncSetUefiCertificates3: Set the UEFI certificates for a pool and all its hosts. Deprecated: use set_custom_uefi_certificates instead
// Version: 22.16.0
func (pool) AsyncSetUefiCertificates3(session *Session, self PoolRef, value string) (retval TaskRef, err error) {
	method := "Async.pool.set_uefi_certificates"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// DisableRepositoryProxy: Disable the proxy for RPM package repositories.
// Version: 21.4.0
func (pool) DisableRepositoryProxy(session *Session, self PoolRef) (err error) {
	method := "pool.disable_repository_proxy"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// AsyncDisableRepositoryProxy: Disable the proxy for RPM package repositories.
// Version: 21.4.0
func (pool) AsyncDisableRepositoryProxy(session *Session, self PoolRef) (retval TaskRef, err error) {
	method := "Async.pool.disable_repository_proxy"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// DisableRepositoryProxy2: Disable the proxy for RPM package repositories.
// Version: 21.4.0
func (pool) DisableRepositoryProxy2(session *Session, self PoolRef) (err error) {
	method := "pool.disable_repository_proxy"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// AsyncDisableRepositoryProxy2: Disable the proxy for RPM package repositories.
// Version: 21.4.0
func (pool) AsyncDisableRepositoryProxy2(session *Session, self PoolRef) (retval TaskRef, err error) {
	method := "Async.pool.disable_repository_proxy"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// ConfigureRepositoryProxy: Configure proxy for RPM package repositories.
// Version: 21.3.0
func (pool) ConfigureRepositoryProxy(session *Session, self PoolRef, url string, username string, password string) (err error) {
	method := "pool.configure_repository_proxy"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	urlArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "url"), url)
	if err != nil {
		return
	}
	usernameArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "username"), username)
	if err != nil {
		return
	}
	passwordArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "password"), password)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, urlArg, usernameArg, passwordArg)
	return
}

// AsyncConfigureRepositoryProxy: Configure proxy for RPM package repositories.
// Version: 21.3.0
func (pool) AsyncConfigureRepositoryProxy(session *Session, self PoolRef, url string, username string, password string) (retval TaskRef, err error) {
	method := "Async.pool.configure_repository_proxy"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	urlArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "url"), url)
	if err != nil {
		return
	}
	usernameArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "username"), username)
	if err != nil {
		return
	}
	passwordArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "password"), password)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg, urlArg, usernameArg, passwordArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// ConfigureRepositoryProxy5: Configure proxy for RPM package repositories.
// Version: 21.3.0
func (pool) ConfigureRepositoryProxy5(session *Session, self PoolRef, url string, username string, password string) (err error) {
	method := "pool.configure_repository_proxy"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	urlArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "url"), url)
	if err != nil {
		return
	}
	usernameArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "username"), username)
	if err != nil {
		return
	}
	passwordArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "password"), password)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, urlArg, usernameArg, passwordArg)
	return
}

// AsyncConfigureRepositoryProxy5: Configure proxy for RPM package repositories.
// Version: 21.3.0
func (pool) AsyncConfigureRepositoryProxy5(session *Session, self PoolRef, url string, username string, password string) (retval TaskRef, err error) {
	method := "Async.pool.configure_repository_proxy"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	urlArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "url"), url)
	if err != nil {
		return
	}
	usernameArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "username"), username)
	if err != nil {
		return
	}
	passwordArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "password"), password)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg, urlArg, usernameArg, passwordArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// DisableClientCertificateAuth: Disable client certificate authentication on the pool
// Version: 1.318.0
func (pool) DisableClientCertificateAuth(session *Session, self PoolRef) (err error) {
	method := "pool.disable_client_certificate_auth"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// AsyncDisableClientCertificateAuth: Disable client certificate authentication on the pool
// Version: 1.318.0
func (pool) AsyncDisableClientCertificateAuth(session *Session, self PoolRef) (retval TaskRef, err error) {
	method := "Async.pool.disable_client_certificate_auth"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// DisableClientCertificateAuth2: Disable client certificate authentication on the pool
// Version: 1.318.0
func (pool) DisableClientCertificateAuth2(session *Session, self PoolRef) (err error) {
	method := "pool.disable_client_certificate_auth"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// AsyncDisableClientCertificateAuth2: Disable client certificate authentication on the pool
// Version: 1.318.0
func (pool) AsyncDisableClientCertificateAuth2(session *Session, self PoolRef) (retval TaskRef, err error) {
	method := "Async.pool.disable_client_certificate_auth"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// EnableClientCertificateAuth: Enable client certificate authentication on the pool
// Version: 1.318.0
func (pool) EnableClientCertificateAuth(session *Session, self PoolRef, name string) (err error) {
	method := "pool.enable_client_certificate_auth"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	nameArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "name"), name)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, nameArg)
	return
}

// AsyncEnableClientCertificateAuth: Enable client certificate authentication on the pool
// Version: 1.318.0
func (pool) AsyncEnableClientCertificateAuth(session *Session, self PoolRef, name string) (retval TaskRef, err error) {
	method := "Async.pool.enable_client_certificate_auth"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	nameArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "name"), name)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg, nameArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// EnableClientCertificateAuth3: Enable client certificate authentication on the pool
// Version: 1.318.0
func (pool) EnableClientCertificateAuth3(session *Session, self PoolRef, name string) (err error) {
	method := "pool.enable_client_certificate_auth"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	nameArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "name"), name)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, nameArg)
	return
}

// AsyncEnableClientCertificateAuth3: Enable client certificate authentication on the pool
// Version: 1.318.0
func (pool) AsyncEnableClientCertificateAuth3(session *Session, self PoolRef, name string) (retval TaskRef, err error) {
	method := "Async.pool.enable_client_certificate_auth"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	nameArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "name"), name)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg, nameArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// CheckUpdateReadiness: Check if the pool is ready to be updated. If not, report the reasons.
// Version: 1.304.0
func (pool) CheckUpdateReadiness(session *Session, self PoolRef, requiresReboot bool) (retval [][]string, err error) {
	method := "pool.check_update_readiness"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	requiresRebootArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "requires_reboot"), requiresReboot)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg, requiresRebootArg)
	if err != nil {
		return
	}
	retval, err = deserializeStringSetSet(method+" -> ", result)
	return
}

// AsyncCheckUpdateReadiness: Check if the pool is ready to be updated. If not, report the reasons.
// Version: 1.304.0
func (pool) AsyncCheckUpdateReadiness(session *Session, self PoolRef, requiresReboot bool) (retval TaskRef, err error) {
	method := "Async.pool.check_update_readiness"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	requiresRebootArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "requires_reboot"), requiresReboot)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg, requiresRebootArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// CheckUpdateReadiness3: Check if the pool is ready to be updated. If not, report the reasons.
// Version: 1.304.0
func (pool) CheckUpdateReadiness3(session *Session, self PoolRef, requiresReboot bool) (retval [][]string, err error) {
	method := "pool.check_update_readiness"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	requiresRebootArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "requires_reboot"), requiresReboot)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg, requiresRebootArg)
	if err != nil {
		return
	}
	retval, err = deserializeStringSetSet(method+" -> ", result)
	return
}

// AsyncCheckUpdateReadiness3: Check if the pool is ready to be updated. If not, report the reasons.
// Version: 1.304.0
func (pool) AsyncCheckUpdateReadiness3(session *Session, self PoolRef, requiresReboot bool) (retval TaskRef, err error) {
	method := "Async.pool.check_update_readiness"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	requiresRebootArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "requires_reboot"), requiresReboot)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg, requiresRebootArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// SyncUpdates: Sync with the enabled repository
// Version: 25.7.0
func (pool) SyncUpdates(session *Session, self PoolRef, force bool, token string, tokenID string, username string, password string) (retval string, err error) {
	method := "pool.sync_updates"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	forceArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "force"), force)
	if err != nil {
		return
	}
	tokenArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "token"), token)
	if err != nil {
		return
	}
	tokenIDArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "token_id"), tokenID)
	if err != nil {
		return
	}
	usernameArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "username"), username)
	if err != nil {
		return
	}
	passwordArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "password"), password)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg, forceArg, tokenArg, tokenIDArg, usernameArg, passwordArg)
	if err != nil {
		return
	}
	retval, err = deserializeString(method+" -> ", result)
	return
}

// AsyncSyncUpdates: Sync with the enabled repository
// Version: 25.7.0
func (pool) AsyncSyncUpdates(session *Session, self PoolRef, force bool, token string, tokenID string, username string, password string) (retval TaskRef, err error) {
	method := "Async.pool.sync_updates"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	forceArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "force"), force)
	if err != nil {
		return
	}
	tokenArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "token"), token)
	if err != nil {
		return
	}
	tokenIDArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "token_id"), tokenID)
	if err != nil {
		return
	}
	usernameArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "username"), username)
	if err != nil {
		return
	}
	passwordArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "password"), password)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg, forceArg, tokenArg, tokenIDArg, usernameArg, passwordArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// SyncUpdates7: Sync with the enabled repository
// Version: 25.7.0
func (pool) SyncUpdates7(session *Session, self PoolRef, force bool, token string, tokenID string, username string, password string) (retval string, err error) {
	method := "pool.sync_updates"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	forceArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "force"), force)
	if err != nil {
		return
	}
	tokenArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "token"), token)
	if err != nil {
		return
	}
	tokenIDArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "token_id"), tokenID)
	if err != nil {
		return
	}
	usernameArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "username"), username)
	if err != nil {
		return
	}
	passwordArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "password"), password)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg, forceArg, tokenArg, tokenIDArg, usernameArg, passwordArg)
	if err != nil {
		return
	}
	retval, err = deserializeString(method+" -> ", result)
	return
}

// AsyncSyncUpdates7: Sync with the enabled repository
// Version: 25.7.0
func (pool) AsyncSyncUpdates7(session *Session, self PoolRef, force bool, token string, tokenID string, username string, password string) (retval TaskRef, err error) {
	method := "Async.pool.sync_updates"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	forceArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "force"), force)
	if err != nil {
		return
	}
	tokenArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "token"), token)
	if err != nil {
		return
	}
	tokenIDArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "token_id"), tokenID)
	if err != nil {
		return
	}
	usernameArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "username"), username)
	if err != nil {
		return
	}
	passwordArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "password"), password)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg, forceArg, tokenArg, tokenIDArg, usernameArg, passwordArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// SyncUpdates5: Sync with the enabled repository
// Version: 1.329.0
func (pool) SyncUpdates5(session *Session, self PoolRef, force bool, token string, tokenID string) (retval string, err error) {
	method := "pool.sync_updates"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	forceArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "force"), force)
	if err != nil {
		return
	}
	tokenArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "token"), token)
	if err != nil {
		return
	}
	tokenIDArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "token_id"), tokenID)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg, forceArg, tokenArg, tokenIDArg)
	if err != nil {
		return
	}
	retval, err = deserializeString(method+" -> ", result)
	return
}

// AsyncSyncUpdates5: Sync with the enabled repository
// Version: 1.329.0
func (pool) AsyncSyncUpdates5(session *Session, self PoolRef, force bool, token string, tokenID string) (retval TaskRef, err error) {
	method := "Async.pool.sync_updates"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	forceArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "force"), force)
	if err != nil {
		return
	}
	tokenArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "token"), token)
	if err != nil {
		return
	}
	tokenIDArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "token_id"), tokenID)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg, forceArg, tokenArg, tokenIDArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// RemoveRepository: Remove a repository from the enabled set
// Version: 1.301.0
func (pool) RemoveRepository(session *Session, self PoolRef, value RepositoryRef) (err error) {
	method := "pool.remove_repository"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeRepositoryRef(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	return
}

// AsyncRemoveRepository: Remove a repository from the enabled set
// Version: 1.301.0
func (pool) AsyncRemoveRepository(session *Session, self PoolRef, value RepositoryRef) (retval TaskRef, err error) {
	method := "Async.pool.remove_repository"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeRepositoryRef(fmt.Sprintf("%s(%s)", method, "value"), value)
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

// RemoveRepository3: Remove a repository from the enabled set
// Version: 1.301.0
func (pool) RemoveRepository3(session *Session, self PoolRef, value RepositoryRef) (err error) {
	method := "pool.remove_repository"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeRepositoryRef(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	return
}

// AsyncRemoveRepository3: Remove a repository from the enabled set
// Version: 1.301.0
func (pool) AsyncRemoveRepository3(session *Session, self PoolRef, value RepositoryRef) (retval TaskRef, err error) {
	method := "Async.pool.remove_repository"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeRepositoryRef(fmt.Sprintf("%s(%s)", method, "value"), value)
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

// AddRepository: Add a repository to the enabled set
// Version: 1.301.0
func (pool) AddRepository(session *Session, self PoolRef, value RepositoryRef) (err error) {
	method := "pool.add_repository"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeRepositoryRef(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	return
}

// AsyncAddRepository: Add a repository to the enabled set
// Version: 1.301.0
func (pool) AsyncAddRepository(session *Session, self PoolRef, value RepositoryRef) (retval TaskRef, err error) {
	method := "Async.pool.add_repository"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeRepositoryRef(fmt.Sprintf("%s(%s)", method, "value"), value)
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

// AddRepository3: Add a repository to the enabled set
// Version: 1.301.0
func (pool) AddRepository3(session *Session, self PoolRef, value RepositoryRef) (err error) {
	method := "pool.add_repository"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeRepositoryRef(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	return
}

// AsyncAddRepository3: Add a repository to the enabled set
// Version: 1.301.0
func (pool) AsyncAddRepository3(session *Session, self PoolRef, value RepositoryRef) (retval TaskRef, err error) {
	method := "Async.pool.add_repository"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeRepositoryRef(fmt.Sprintf("%s(%s)", method, "value"), value)
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

// SetRepositories: Set enabled set of repositories
// Version: 1.301.0
func (pool) SetRepositories(session *Session, self PoolRef, value []RepositoryRef) (err error) {
	method := "pool.set_repositories"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeRepositoryRefSet(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	return
}

// AsyncSetRepositories: Set enabled set of repositories
// Version: 1.301.0
func (pool) AsyncSetRepositories(session *Session, self PoolRef, value []RepositoryRef) (retval TaskRef, err error) {
	method := "Async.pool.set_repositories"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeRepositoryRefSet(fmt.Sprintf("%s(%s)", method, "value"), value)
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

// SetRepositories3: Set enabled set of repositories
// Version: 1.301.0
func (pool) SetRepositories3(session *Session, self PoolRef, value []RepositoryRef) (err error) {
	method := "pool.set_repositories"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeRepositoryRefSet(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	return
}

// AsyncSetRepositories3: Set enabled set of repositories
// Version: 1.301.0
func (pool) AsyncSetRepositories3(session *Session, self PoolRef, value []RepositoryRef) (retval TaskRef, err error) {
	method := "Async.pool.set_repositories"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeRepositoryRefSet(fmt.Sprintf("%s(%s)", method, "value"), value)
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

// RotateSecret:
// Version: stockholm_psr
//
// Errors:
// INTERNAL_ERROR - The server failed to handle your request, due to an internal error. The given message may give details useful for debugging the problem.
// HOST_IS_SLAVE - You cannot make regular API calls directly on a supporter. Please pass API calls via the coordinator host.
// CANNOT_CONTACT_HOST - Cannot forward messages because the server cannot be contacted. The server may be switched off or there may be network connectivity problems.
// HA_IS_ENABLED - The operation could not be performed because HA is enabled on the Pool
// NOT_SUPPORTED_DURING_UPGRADE - This operation is not supported during an upgrade.
func (pool) RotateSecret(session *Session) (err error) {
	method := "pool.rotate_secret"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg)
	return
}

// AsyncRotateSecret:
// Version: stockholm_psr
//
// Errors:
// INTERNAL_ERROR - The server failed to handle your request, due to an internal error. The given message may give details useful for debugging the problem.
// HOST_IS_SLAVE - You cannot make regular API calls directly on a supporter. Please pass API calls via the coordinator host.
// CANNOT_CONTACT_HOST - Cannot forward messages because the server cannot be contacted. The server may be switched off or there may be network connectivity problems.
// HA_IS_ENABLED - The operation could not be performed because HA is enabled on the Pool
// NOT_SUPPORTED_DURING_UPGRADE - This operation is not supported during an upgrade.
func (pool) AsyncRotateSecret(session *Session) (retval TaskRef, err error) {
	method := "Async.pool.rotate_secret"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// RotateSecret1:
// Version: stockholm_psr
//
// Errors:
// INTERNAL_ERROR - The server failed to handle your request, due to an internal error. The given message may give details useful for debugging the problem.
// HOST_IS_SLAVE - You cannot make regular API calls directly on a supporter. Please pass API calls via the coordinator host.
// CANNOT_CONTACT_HOST - Cannot forward messages because the server cannot be contacted. The server may be switched off or there may be network connectivity problems.
// HA_IS_ENABLED - The operation could not be performed because HA is enabled on the Pool
// NOT_SUPPORTED_DURING_UPGRADE - This operation is not supported during an upgrade.
func (pool) RotateSecret1(session *Session) (err error) {
	method := "pool.rotate_secret"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg)
	return
}

// AsyncRotateSecret1:
// Version: stockholm_psr
//
// Errors:
// INTERNAL_ERROR - The server failed to handle your request, due to an internal error. The given message may give details useful for debugging the problem.
// HOST_IS_SLAVE - You cannot make regular API calls directly on a supporter. Please pass API calls via the coordinator host.
// CANNOT_CONTACT_HOST - Cannot forward messages because the server cannot be contacted. The server may be switched off or there may be network connectivity problems.
// HA_IS_ENABLED - The operation could not be performed because HA is enabled on the Pool
// NOT_SUPPORTED_DURING_UPGRADE - This operation is not supported during an upgrade.
func (pool) AsyncRotateSecret1(session *Session) (retval TaskRef, err error) {
	method := "Async.pool.rotate_secret"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// RemoveFromGuestAgentConfig: Remove a key-value pair from the pool-wide guest agent configuration
// Version: dundee
func (pool) RemoveFromGuestAgentConfig(session *Session, self PoolRef, key string) (err error) {
	method := "pool.remove_from_guest_agent_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// AsyncRemoveFromGuestAgentConfig: Remove a key-value pair from the pool-wide guest agent configuration
// Version: dundee
func (pool) AsyncRemoveFromGuestAgentConfig(session *Session, self PoolRef, key string) (retval TaskRef, err error) {
	method := "Async.pool.remove_from_guest_agent_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	keyArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "key"), key)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg, keyArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// RemoveFromGuestAgentConfig3: Remove a key-value pair from the pool-wide guest agent configuration
// Version: dundee
func (pool) RemoveFromGuestAgentConfig3(session *Session, self PoolRef, key string) (err error) {
	method := "pool.remove_from_guest_agent_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// AsyncRemoveFromGuestAgentConfig3: Remove a key-value pair from the pool-wide guest agent configuration
// Version: dundee
func (pool) AsyncRemoveFromGuestAgentConfig3(session *Session, self PoolRef, key string) (retval TaskRef, err error) {
	method := "Async.pool.remove_from_guest_agent_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	keyArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "key"), key)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg, keyArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// AddToGuestAgentConfig: Add a key-value pair to the pool-wide guest agent configuration
// Version: dundee
func (pool) AddToGuestAgentConfig(session *Session, self PoolRef, key string, value string) (err error) {
	method := "pool.add_to_guest_agent_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// AsyncAddToGuestAgentConfig: Add a key-value pair to the pool-wide guest agent configuration
// Version: dundee
func (pool) AsyncAddToGuestAgentConfig(session *Session, self PoolRef, key string, value string) (retval TaskRef, err error) {
	method := "Async.pool.add_to_guest_agent_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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
	result, err := session.client.sendCall(method, sessionIDArg, selfArg, keyArg, valueArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// AddToGuestAgentConfig4: Add a key-value pair to the pool-wide guest agent configuration
// Version: dundee
func (pool) AddToGuestAgentConfig4(session *Session, self PoolRef, key string, value string) (err error) {
	method := "pool.add_to_guest_agent_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// AsyncAddToGuestAgentConfig4: Add a key-value pair to the pool-wide guest agent configuration
// Version: dundee
func (pool) AsyncAddToGuestAgentConfig4(session *Session, self PoolRef, key string, value string) (retval TaskRef, err error) {
	method := "Async.pool.add_to_guest_agent_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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
	result, err := session.client.sendCall(method, sessionIDArg, selfArg, keyArg, valueArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// HasExtension: Return true if the extension is available on the pool
// Version: dundee
func (pool) HasExtension(session *Session, self PoolRef, name string) (retval bool, err error) {
	method := "pool.has_extension"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	nameArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "name"), name)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg, nameArg)
	if err != nil {
		return
	}
	retval, err = deserializeBool(method+" -> ", result)
	return
}

// AsyncHasExtension: Return true if the extension is available on the pool
// Version: dundee
func (pool) AsyncHasExtension(session *Session, self PoolRef, name string) (retval TaskRef, err error) {
	method := "Async.pool.has_extension"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	nameArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "name"), name)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg, nameArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// HasExtension3: Return true if the extension is available on the pool
// Version: dundee
func (pool) HasExtension3(session *Session, self PoolRef, name string) (retval bool, err error) {
	method := "pool.has_extension"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	nameArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "name"), name)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg, nameArg)
	if err != nil {
		return
	}
	retval, err = deserializeBool(method+" -> ", result)
	return
}

// AsyncHasExtension3: Return true if the extension is available on the pool
// Version: dundee
func (pool) AsyncHasExtension3(session *Session, self PoolRef, name string) (retval TaskRef, err error) {
	method := "Async.pool.has_extension"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	nameArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "name"), name)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg, nameArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// SetIgmpSnoopingEnabled: Enable or disable IGMP Snooping on the pool.
// Version: inverness
func (pool) SetIgmpSnoopingEnabled(session *Session, self PoolRef, value bool) (err error) {
	method := "pool.set_igmp_snooping_enabled"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// AsyncSetIgmpSnoopingEnabled: Enable or disable IGMP Snooping on the pool.
// Version: inverness
func (pool) AsyncSetIgmpSnoopingEnabled(session *Session, self PoolRef, value bool) (retval TaskRef, err error) {
	method := "Async.pool.set_igmp_snooping_enabled"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetIgmpSnoopingEnabled3: Enable or disable IGMP Snooping on the pool.
// Version: inverness
func (pool) SetIgmpSnoopingEnabled3(session *Session, self PoolRef, value bool) (err error) {
	method := "pool.set_igmp_snooping_enabled"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// AsyncSetIgmpSnoopingEnabled3: Enable or disable IGMP Snooping on the pool.
// Version: inverness
func (pool) AsyncSetIgmpSnoopingEnabled3(session *Session, self PoolRef, value bool) (retval TaskRef, err error) {
	method := "Async.pool.set_igmp_snooping_enabled"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// DisableSslLegacy: Sets ssl_legacy false on each host, pool-master last. See Host.ssl_legacy and Host.set_ssl_legacy.
// Version: dundee
func (pool) DisableSslLegacy(session *Session, self PoolRef) (err error) {
	method := "pool.disable_ssl_legacy"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// AsyncDisableSslLegacy: Sets ssl_legacy false on each host, pool-master last. See Host.ssl_legacy and Host.set_ssl_legacy.
// Version: dundee
func (pool) AsyncDisableSslLegacy(session *Session, self PoolRef) (retval TaskRef, err error) {
	method := "Async.pool.disable_ssl_legacy"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// DisableSslLegacy2: Sets ssl_legacy false on each host, pool-master last. See Host.ssl_legacy and Host.set_ssl_legacy.
// Version: dundee
func (pool) DisableSslLegacy2(session *Session, self PoolRef) (err error) {
	method := "pool.disable_ssl_legacy"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// AsyncDisableSslLegacy2: Sets ssl_legacy false on each host, pool-master last. See Host.ssl_legacy and Host.set_ssl_legacy.
// Version: dundee
func (pool) AsyncDisableSslLegacy2(session *Session, self PoolRef) (retval TaskRef, err error) {
	method := "Async.pool.disable_ssl_legacy"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// EnableSslLegacy: Sets ssl_legacy true on each host, pool-master last. See Host.ssl_legacy and Host.set_ssl_legacy.
// Version: dundee
func (pool) EnableSslLegacy(session *Session, self PoolRef) (err error) {
	method := "pool.enable_ssl_legacy"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// AsyncEnableSslLegacy: Sets ssl_legacy true on each host, pool-master last. See Host.ssl_legacy and Host.set_ssl_legacy.
// Version: dundee
func (pool) AsyncEnableSslLegacy(session *Session, self PoolRef) (retval TaskRef, err error) {
	method := "Async.pool.enable_ssl_legacy"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// EnableSslLegacy2: Sets ssl_legacy true on each host, pool-master last. See Host.ssl_legacy and Host.set_ssl_legacy.
// Version: dundee
func (pool) EnableSslLegacy2(session *Session, self PoolRef) (err error) {
	method := "pool.enable_ssl_legacy"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// AsyncEnableSslLegacy2: Sets ssl_legacy true on each host, pool-master last. See Host.ssl_legacy and Host.set_ssl_legacy.
// Version: dundee
func (pool) AsyncEnableSslLegacy2(session *Session, self PoolRef) (retval TaskRef, err error) {
	method := "Async.pool.enable_ssl_legacy"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// ApplyEdition: Apply an edition to all hosts in the pool
// Version: clearwater
func (pool) ApplyEdition(session *Session, self PoolRef, edition string) (err error) {
	method := "pool.apply_edition"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	editionArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "edition"), edition)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, editionArg)
	return
}

// AsyncApplyEdition: Apply an edition to all hosts in the pool
// Version: clearwater
func (pool) AsyncApplyEdition(session *Session, self PoolRef, edition string) (retval TaskRef, err error) {
	method := "Async.pool.apply_edition"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	editionArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "edition"), edition)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg, editionArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// ApplyEdition3: Apply an edition to all hosts in the pool
// Version: clearwater
func (pool) ApplyEdition3(session *Session, self PoolRef, edition string) (err error) {
	method := "pool.apply_edition"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	editionArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "edition"), edition)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, editionArg)
	return
}

// AsyncApplyEdition3: Apply an edition to all hosts in the pool
// Version: clearwater
func (pool) AsyncApplyEdition3(session *Session, self PoolRef, edition string) (retval TaskRef, err error) {
	method := "Async.pool.apply_edition"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	editionArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "edition"), edition)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg, editionArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// GetLicenseState: This call returns the license state for the pool
// Version: clearwater
func (pool) GetLicenseState(session *Session, self PoolRef) (retval map[string]string, err error) {
	method := "pool.get_license_state"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// AsyncGetLicenseState: This call returns the license state for the pool
// Version: clearwater
func (pool) AsyncGetLicenseState(session *Session, self PoolRef) (retval TaskRef, err error) {
	method := "Async.pool.get_license_state"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetLicenseState2: This call returns the license state for the pool
// Version: clearwater
func (pool) GetLicenseState2(session *Session, self PoolRef) (retval map[string]string, err error) {
	method := "pool.get_license_state"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// AsyncGetLicenseState2: This call returns the license state for the pool
// Version: clearwater
func (pool) AsyncGetLicenseState2(session *Session, self PoolRef) (retval TaskRef, err error) {
	method := "Async.pool.get_license_state"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// DisableLocalStorageCaching: This call disables pool-wide local storage caching
// Version: cowley
func (pool) DisableLocalStorageCaching(session *Session, self PoolRef) (err error) {
	method := "pool.disable_local_storage_caching"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// AsyncDisableLocalStorageCaching: This call disables pool-wide local storage caching
// Version: cowley
func (pool) AsyncDisableLocalStorageCaching(session *Session, self PoolRef) (retval TaskRef, err error) {
	method := "Async.pool.disable_local_storage_caching"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// DisableLocalStorageCaching2: This call disables pool-wide local storage caching
// Version: cowley
func (pool) DisableLocalStorageCaching2(session *Session, self PoolRef) (err error) {
	method := "pool.disable_local_storage_caching"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// AsyncDisableLocalStorageCaching2: This call disables pool-wide local storage caching
// Version: cowley
func (pool) AsyncDisableLocalStorageCaching2(session *Session, self PoolRef) (retval TaskRef, err error) {
	method := "Async.pool.disable_local_storage_caching"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// EnableLocalStorageCaching: This call attempts to enable pool-wide local storage caching
// Version: cowley
func (pool) EnableLocalStorageCaching(session *Session, self PoolRef) (err error) {
	method := "pool.enable_local_storage_caching"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// AsyncEnableLocalStorageCaching: This call attempts to enable pool-wide local storage caching
// Version: cowley
func (pool) AsyncEnableLocalStorageCaching(session *Session, self PoolRef) (retval TaskRef, err error) {
	method := "Async.pool.enable_local_storage_caching"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// EnableLocalStorageCaching2: This call attempts to enable pool-wide local storage caching
// Version: cowley
func (pool) EnableLocalStorageCaching2(session *Session, self PoolRef) (err error) {
	method := "pool.enable_local_storage_caching"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// AsyncEnableLocalStorageCaching2: This call attempts to enable pool-wide local storage caching
// Version: cowley
func (pool) AsyncEnableLocalStorageCaching2(session *Session, self PoolRef) (retval TaskRef, err error) {
	method := "Async.pool.enable_local_storage_caching"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// TestArchiveTarget: This call tests if a location is valid
// Version: cowley
func (pool) TestArchiveTarget(session *Session, self PoolRef, config map[string]string) (retval string, err error) {
	method := "pool.test_archive_target"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	configArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "config"), config)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg, configArg)
	if err != nil {
		return
	}
	retval, err = deserializeString(method+" -> ", result)
	return
}

// TestArchiveTarget3: This call tests if a location is valid
// Version: cowley
func (pool) TestArchiveTarget3(session *Session, self PoolRef, config map[string]string) (retval string, err error) {
	method := "pool.test_archive_target"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	configArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "config"), config)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg, configArg)
	if err != nil {
		return
	}
	retval, err = deserializeString(method+" -> ", result)
	return
}

// SetVswitchController: Set the IP address of the vswitch controller.
// Version: midnight-ride
func (pool) SetVswitchController(session *Session, address string) (err error) {
	method := "pool.set_vswitch_controller"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	addressArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "address"), address)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, addressArg)
	return
}

// AsyncSetVswitchController: Set the IP address of the vswitch controller.
// Version: midnight-ride
func (pool) AsyncSetVswitchController(session *Session, address string) (retval TaskRef, err error) {
	method := "Async.pool.set_vswitch_controller"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	addressArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "address"), address)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, addressArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// SetVswitchController2: Set the IP address of the vswitch controller.
// Version: midnight-ride
func (pool) SetVswitchController2(session *Session, address string) (err error) {
	method := "pool.set_vswitch_controller"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	addressArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "address"), address)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, addressArg)
	return
}

// AsyncSetVswitchController2: Set the IP address of the vswitch controller.
// Version: midnight-ride
func (pool) AsyncSetVswitchController2(session *Session, address string) (retval TaskRef, err error) {
	method := "Async.pool.set_vswitch_controller"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	addressArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "address"), address)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, addressArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// DisableRedoLog: Disable the redo log if in use, unless HA is enabled.
// Version: midnight-ride
func (pool) DisableRedoLog(session *Session) (err error) {
	method := "pool.disable_redo_log"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg)
	return
}

// AsyncDisableRedoLog: Disable the redo log if in use, unless HA is enabled.
// Version: midnight-ride
func (pool) AsyncDisableRedoLog(session *Session) (retval TaskRef, err error) {
	method := "Async.pool.disable_redo_log"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// DisableRedoLog1: Disable the redo log if in use, unless HA is enabled.
// Version: midnight-ride
func (pool) DisableRedoLog1(session *Session) (err error) {
	method := "pool.disable_redo_log"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg)
	return
}

// AsyncDisableRedoLog1: Disable the redo log if in use, unless HA is enabled.
// Version: midnight-ride
func (pool) AsyncDisableRedoLog1(session *Session) (retval TaskRef, err error) {
	method := "Async.pool.disable_redo_log"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// EnableRedoLog: Enable the redo log on the given SR and start using it, unless HA is enabled.
// Version: midnight-ride
func (pool) EnableRedoLog(session *Session, sr SRRef) (err error) {
	method := "pool.enable_redo_log"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	srArg, err := serializeSRRef(fmt.Sprintf("%s(%s)", method, "sr"), sr)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, srArg)
	return
}

// AsyncEnableRedoLog: Enable the redo log on the given SR and start using it, unless HA is enabled.
// Version: midnight-ride
func (pool) AsyncEnableRedoLog(session *Session, sr SRRef) (retval TaskRef, err error) {
	method := "Async.pool.enable_redo_log"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	srArg, err := serializeSRRef(fmt.Sprintf("%s(%s)", method, "sr"), sr)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, srArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// EnableRedoLog2: Enable the redo log on the given SR and start using it, unless HA is enabled.
// Version: midnight-ride
func (pool) EnableRedoLog2(session *Session, sr SRRef) (err error) {
	method := "pool.enable_redo_log"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	srArg, err := serializeSRRef(fmt.Sprintf("%s(%s)", method, "sr"), sr)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, srArg)
	return
}

// AsyncEnableRedoLog2: Enable the redo log on the given SR and start using it, unless HA is enabled.
// Version: midnight-ride
func (pool) AsyncEnableRedoLog2(session *Session, sr SRRef) (retval TaskRef, err error) {
	method := "Async.pool.enable_redo_log"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	srArg, err := serializeSRRef(fmt.Sprintf("%s(%s)", method, "sr"), sr)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, srArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// EnableTLSVerification: Enable TLS server certificate verification
// Version: 1.290.0
func (pool) EnableTLSVerification(session *Session) (err error) {
	method := "pool.enable_tls_verification"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg)
	return
}

// EnableTLSVerification1: Enable TLS server certificate verification
// Version: 1.290.0
func (pool) EnableTLSVerification1(session *Session) (err error) {
	method := "pool.enable_tls_verification"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg)
	return
}

// CertificateSync: Copy the TLS CA certificates and CRLs of the master to all slaves.
// Version: george
func (pool) CertificateSync(session *Session) (err error) {
	method := "pool.certificate_sync"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg)
	return
}

// AsyncCertificateSync: Copy the TLS CA certificates and CRLs of the master to all slaves.
// Version: george
func (pool) AsyncCertificateSync(session *Session) (retval TaskRef, err error) {
	method := "Async.pool.certificate_sync"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// CertificateSync1: Copy the TLS CA certificates and CRLs of the master to all slaves.
// Version: george
func (pool) CertificateSync1(session *Session) (err error) {
	method := "pool.certificate_sync"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg)
	return
}

// AsyncCertificateSync1: Copy the TLS CA certificates and CRLs of the master to all slaves.
// Version: george
func (pool) AsyncCertificateSync1(session *Session) (retval TaskRef, err error) {
	method := "Async.pool.certificate_sync"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// CrlList: List the names of all installed TLS CA-issued Certificate Revocation Lists.
// Version: george
func (pool) CrlList(session *Session) (retval []string, err error) {
	method := "pool.crl_list"
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

// AsyncCrlList: List the names of all installed TLS CA-issued Certificate Revocation Lists.
// Version: george
func (pool) AsyncCrlList(session *Session) (retval TaskRef, err error) {
	method := "Async.pool.crl_list"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// CrlList1: List the names of all installed TLS CA-issued Certificate Revocation Lists.
// Version: george
func (pool) CrlList1(session *Session) (retval []string, err error) {
	method := "pool.crl_list"
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

// AsyncCrlList1: List the names of all installed TLS CA-issued Certificate Revocation Lists.
// Version: george
func (pool) AsyncCrlList1(session *Session) (retval TaskRef, err error) {
	method := "Async.pool.crl_list"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// CrlUninstall: Remove a pool-wide TLS CA-issued Certificate Revocation List.
// Version: george
func (pool) CrlUninstall(session *Session, name string) (err error) {
	method := "pool.crl_uninstall"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	nameArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "name"), name)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, nameArg)
	return
}

// AsyncCrlUninstall: Remove a pool-wide TLS CA-issued Certificate Revocation List.
// Version: george
func (pool) AsyncCrlUninstall(session *Session, name string) (retval TaskRef, err error) {
	method := "Async.pool.crl_uninstall"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	nameArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "name"), name)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, nameArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// CrlUninstall2: Remove a pool-wide TLS CA-issued Certificate Revocation List.
// Version: george
func (pool) CrlUninstall2(session *Session, name string) (err error) {
	method := "pool.crl_uninstall"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	nameArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "name"), name)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, nameArg)
	return
}

// AsyncCrlUninstall2: Remove a pool-wide TLS CA-issued Certificate Revocation List.
// Version: george
func (pool) AsyncCrlUninstall2(session *Session, name string) (retval TaskRef, err error) {
	method := "Async.pool.crl_uninstall"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	nameArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "name"), name)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, nameArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// CrlInstall: Install a TLS CA-issued Certificate Revocation List, pool-wide.
// Version: george
func (pool) CrlInstall(session *Session, name string, cert string) (err error) {
	method := "pool.crl_install"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	nameArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "name"), name)
	if err != nil {
		return
	}
	certArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "cert"), cert)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, nameArg, certArg)
	return
}

// AsyncCrlInstall: Install a TLS CA-issued Certificate Revocation List, pool-wide.
// Version: george
func (pool) AsyncCrlInstall(session *Session, name string, cert string) (retval TaskRef, err error) {
	method := "Async.pool.crl_install"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	nameArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "name"), name)
	if err != nil {
		return
	}
	certArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "cert"), cert)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, nameArg, certArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// CrlInstall3: Install a TLS CA-issued Certificate Revocation List, pool-wide.
// Version: george
func (pool) CrlInstall3(session *Session, name string, cert string) (err error) {
	method := "pool.crl_install"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	nameArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "name"), name)
	if err != nil {
		return
	}
	certArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "cert"), cert)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, nameArg, certArg)
	return
}

// AsyncCrlInstall3: Install a TLS CA-issued Certificate Revocation List, pool-wide.
// Version: george
func (pool) AsyncCrlInstall3(session *Session, name string, cert string) (retval TaskRef, err error) {
	method := "Async.pool.crl_install"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	nameArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "name"), name)
	if err != nil {
		return
	}
	certArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "cert"), cert)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, nameArg, certArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// UninstallCaCertificate: Remove a pool-wide TLS CA certificate.
// Version: 24.35.0
func (pool) UninstallCaCertificate(session *Session, name string, force bool) (err error) {
	method := "pool.uninstall_ca_certificate"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	nameArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "name"), name)
	if err != nil {
		return
	}
	forceArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "force"), force)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, nameArg, forceArg)
	return
}

// AsyncUninstallCaCertificate: Remove a pool-wide TLS CA certificate.
// Version: 24.35.0
func (pool) AsyncUninstallCaCertificate(session *Session, name string, force bool) (retval TaskRef, err error) {
	method := "Async.pool.uninstall_ca_certificate"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	nameArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "name"), name)
	if err != nil {
		return
	}
	forceArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "force"), force)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, nameArg, forceArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// UninstallCaCertificate3: Remove a pool-wide TLS CA certificate.
// Version: 24.35.0
func (pool) UninstallCaCertificate3(session *Session, name string, force bool) (err error) {
	method := "pool.uninstall_ca_certificate"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	nameArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "name"), name)
	if err != nil {
		return
	}
	forceArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "force"), force)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, nameArg, forceArg)
	return
}

// AsyncUninstallCaCertificate3: Remove a pool-wide TLS CA certificate.
// Version: 24.35.0
func (pool) AsyncUninstallCaCertificate3(session *Session, name string, force bool) (retval TaskRef, err error) {
	method := "Async.pool.uninstall_ca_certificate"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	nameArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "name"), name)
	if err != nil {
		return
	}
	forceArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "force"), force)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, nameArg, forceArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// UninstallCaCertificate2: Remove a pool-wide TLS CA certificate.
// Version: 1.290.0
func (pool) UninstallCaCertificate2(session *Session, name string) (err error) {
	method := "pool.uninstall_ca_certificate"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	nameArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "name"), name)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, nameArg)
	return
}

// AsyncUninstallCaCertificate2: Remove a pool-wide TLS CA certificate.
// Version: 1.290.0
func (pool) AsyncUninstallCaCertificate2(session *Session, name string) (retval TaskRef, err error) {
	method := "Async.pool.uninstall_ca_certificate"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	nameArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "name"), name)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, nameArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// InstallCaCertificate: Install a TLS CA certificate, pool-wide.
// Version: 1.290.0
func (pool) InstallCaCertificate(session *Session, name string, cert string) (err error) {
	method := "pool.install_ca_certificate"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	nameArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "name"), name)
	if err != nil {
		return
	}
	certArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "cert"), cert)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, nameArg, certArg)
	return
}

// AsyncInstallCaCertificate: Install a TLS CA certificate, pool-wide.
// Version: 1.290.0
func (pool) AsyncInstallCaCertificate(session *Session, name string, cert string) (retval TaskRef, err error) {
	method := "Async.pool.install_ca_certificate"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	nameArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "name"), name)
	if err != nil {
		return
	}
	certArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "cert"), cert)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, nameArg, certArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// InstallCaCertificate3: Install a TLS CA certificate, pool-wide.
// Version: 1.290.0
func (pool) InstallCaCertificate3(session *Session, name string, cert string) (err error) {
	method := "pool.install_ca_certificate"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	nameArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "name"), name)
	if err != nil {
		return
	}
	certArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "cert"), cert)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, nameArg, certArg)
	return
}

// AsyncInstallCaCertificate3: Install a TLS CA certificate, pool-wide.
// Version: 1.290.0
func (pool) AsyncInstallCaCertificate3(session *Session, name string, cert string) (retval TaskRef, err error) {
	method := "Async.pool.install_ca_certificate"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	nameArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "name"), name)
	if err != nil {
		return
	}
	certArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "cert"), cert)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, nameArg, certArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// CertificateList: List the names of all installed TLS CA certificates.
// Version: george
func (pool) CertificateList(session *Session) (retval []string, err error) {
	method := "pool.certificate_list"
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

// AsyncCertificateList: List the names of all installed TLS CA certificates.
// Version: george
func (pool) AsyncCertificateList(session *Session) (retval TaskRef, err error) {
	method := "Async.pool.certificate_list"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// CertificateList1: List the names of all installed TLS CA certificates.
// Version: george
func (pool) CertificateList1(session *Session) (retval []string, err error) {
	method := "pool.certificate_list"
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

// AsyncCertificateList1: List the names of all installed TLS CA certificates.
// Version: george
func (pool) AsyncCertificateList1(session *Session) (retval TaskRef, err error) {
	method := "Async.pool.certificate_list"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// CertificateUninstall: Remove a pool-wide TLS CA certificate.
// Version: george
func (pool) CertificateUninstall(session *Session, name string) (err error) {
	method := "pool.certificate_uninstall"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	nameArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "name"), name)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, nameArg)
	return
}

// AsyncCertificateUninstall: Remove a pool-wide TLS CA certificate.
// Version: george
func (pool) AsyncCertificateUninstall(session *Session, name string) (retval TaskRef, err error) {
	method := "Async.pool.certificate_uninstall"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	nameArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "name"), name)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, nameArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// CertificateUninstall2: Remove a pool-wide TLS CA certificate.
// Version: george
func (pool) CertificateUninstall2(session *Session, name string) (err error) {
	method := "pool.certificate_uninstall"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	nameArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "name"), name)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, nameArg)
	return
}

// AsyncCertificateUninstall2: Remove a pool-wide TLS CA certificate.
// Version: george
func (pool) AsyncCertificateUninstall2(session *Session, name string) (retval TaskRef, err error) {
	method := "Async.pool.certificate_uninstall"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	nameArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "name"), name)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, nameArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// CertificateInstall: Install a TLS CA certificate, pool-wide.
// Version: george
func (pool) CertificateInstall(session *Session, name string, cert string) (err error) {
	method := "pool.certificate_install"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	nameArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "name"), name)
	if err != nil {
		return
	}
	certArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "cert"), cert)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, nameArg, certArg)
	return
}

// AsyncCertificateInstall: Install a TLS CA certificate, pool-wide.
// Version: george
func (pool) AsyncCertificateInstall(session *Session, name string, cert string) (retval TaskRef, err error) {
	method := "Async.pool.certificate_install"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	nameArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "name"), name)
	if err != nil {
		return
	}
	certArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "cert"), cert)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, nameArg, certArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// CertificateInstall3: Install a TLS CA certificate, pool-wide.
// Version: george
func (pool) CertificateInstall3(session *Session, name string, cert string) (err error) {
	method := "pool.certificate_install"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	nameArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "name"), name)
	if err != nil {
		return
	}
	certArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "cert"), cert)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, nameArg, certArg)
	return
}

// AsyncCertificateInstall3: Install a TLS CA certificate, pool-wide.
// Version: george
func (pool) AsyncCertificateInstall3(session *Session, name string, cert string) (retval TaskRef, err error) {
	method := "Async.pool.certificate_install"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	nameArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "name"), name)
	if err != nil {
		return
	}
	certArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "cert"), cert)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, nameArg, certArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// SendTestPost: Send the given body to the given host and port, using HTTPS, and print the response.  This is used for debugging the SSL layer.
// Version: george
func (pool) SendTestPost(session *Session, host string, port int, body string) (retval string, err error) {
	method := "pool.send_test_post"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	portArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "port"), port)
	if err != nil {
		return
	}
	bodyArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "body"), body)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, hostArg, portArg, bodyArg)
	if err != nil {
		return
	}
	retval, err = deserializeString(method+" -> ", result)
	return
}

// AsyncSendTestPost: Send the given body to the given host and port, using HTTPS, and print the response.  This is used for debugging the SSL layer.
// Version: george
func (pool) AsyncSendTestPost(session *Session, host string, port int, body string) (retval TaskRef, err error) {
	method := "Async.pool.send_test_post"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	portArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "port"), port)
	if err != nil {
		return
	}
	bodyArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "body"), body)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, hostArg, portArg, bodyArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// SendTestPost4: Send the given body to the given host and port, using HTTPS, and print the response.  This is used for debugging the SSL layer.
// Version: george
func (pool) SendTestPost4(session *Session, host string, port int, body string) (retval string, err error) {
	method := "pool.send_test_post"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	portArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "port"), port)
	if err != nil {
		return
	}
	bodyArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "body"), body)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, hostArg, portArg, bodyArg)
	if err != nil {
		return
	}
	retval, err = deserializeString(method+" -> ", result)
	return
}

// AsyncSendTestPost4: Send the given body to the given host and port, using HTTPS, and print the response.  This is used for debugging the SSL layer.
// Version: george
func (pool) AsyncSendTestPost4(session *Session, host string, port int, body string) (retval TaskRef, err error) {
	method := "Async.pool.send_test_post"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	portArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "port"), port)
	if err != nil {
		return
	}
	bodyArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "body"), body)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, hostArg, portArg, bodyArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// RetrieveWlbRecommendations: Retrieves vm migrate recommendations for the pool from the workload balancing server
// Version: george
func (pool) RetrieveWlbRecommendations(session *Session) (retval map[VMRef][]string, err error) {
	method := "pool.retrieve_wlb_recommendations"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializeVMRefToStringSetMap(method+" -> ", result)
	return
}

// AsyncRetrieveWlbRecommendations: Retrieves vm migrate recommendations for the pool from the workload balancing server
// Version: george
func (pool) AsyncRetrieveWlbRecommendations(session *Session) (retval TaskRef, err error) {
	method := "Async.pool.retrieve_wlb_recommendations"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// RetrieveWlbRecommendations1: Retrieves vm migrate recommendations for the pool from the workload balancing server
// Version: george
func (pool) RetrieveWlbRecommendations1(session *Session) (retval map[VMRef][]string, err error) {
	method := "pool.retrieve_wlb_recommendations"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializeVMRefToStringSetMap(method+" -> ", result)
	return
}

// AsyncRetrieveWlbRecommendations1: Retrieves vm migrate recommendations for the pool from the workload balancing server
// Version: george
func (pool) AsyncRetrieveWlbRecommendations1(session *Session) (retval TaskRef, err error) {
	method := "Async.pool.retrieve_wlb_recommendations"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// RetrieveWlbConfiguration: Retrieves the pool optimization criteria from the workload balancing server
// Version: george
func (pool) RetrieveWlbConfiguration(session *Session) (retval map[string]string, err error) {
	method := "pool.retrieve_wlb_configuration"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializeStringToStringMap(method+" -> ", result)
	return
}

// AsyncRetrieveWlbConfiguration: Retrieves the pool optimization criteria from the workload balancing server
// Version: george
func (pool) AsyncRetrieveWlbConfiguration(session *Session) (retval TaskRef, err error) {
	method := "Async.pool.retrieve_wlb_configuration"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// RetrieveWlbConfiguration1: Retrieves the pool optimization criteria from the workload balancing server
// Version: george
func (pool) RetrieveWlbConfiguration1(session *Session) (retval map[string]string, err error) {
	method := "pool.retrieve_wlb_configuration"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializeStringToStringMap(method+" -> ", result)
	return
}

// AsyncRetrieveWlbConfiguration1: Retrieves the pool optimization criteria from the workload balancing server
// Version: george
func (pool) AsyncRetrieveWlbConfiguration1(session *Session) (retval TaskRef, err error) {
	method := "Async.pool.retrieve_wlb_configuration"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// SendWlbConfiguration: Sets the pool optimization criteria for the workload balancing server
// Version: george
func (pool) SendWlbConfiguration(session *Session, config map[string]string) (err error) {
	method := "pool.send_wlb_configuration"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	configArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "config"), config)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, configArg)
	return
}

// AsyncSendWlbConfiguration: Sets the pool optimization criteria for the workload balancing server
// Version: george
func (pool) AsyncSendWlbConfiguration(session *Session, config map[string]string) (retval TaskRef, err error) {
	method := "Async.pool.send_wlb_configuration"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	configArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "config"), config)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, configArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// SendWlbConfiguration2: Sets the pool optimization criteria for the workload balancing server
// Version: george
func (pool) SendWlbConfiguration2(session *Session, config map[string]string) (err error) {
	method := "pool.send_wlb_configuration"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	configArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "config"), config)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, configArg)
	return
}

// AsyncSendWlbConfiguration2: Sets the pool optimization criteria for the workload balancing server
// Version: george
func (pool) AsyncSendWlbConfiguration2(session *Session, config map[string]string) (retval TaskRef, err error) {
	method := "Async.pool.send_wlb_configuration"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	configArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "config"), config)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, configArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// DeconfigureWlb: Permanently deconfigures workload balancing monitoring on this pool
// Version: george
func (pool) DeconfigureWlb(session *Session) (err error) {
	method := "pool.deconfigure_wlb"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg)
	return
}

// AsyncDeconfigureWlb: Permanently deconfigures workload balancing monitoring on this pool
// Version: george
func (pool) AsyncDeconfigureWlb(session *Session) (retval TaskRef, err error) {
	method := "Async.pool.deconfigure_wlb"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// DeconfigureWlb1: Permanently deconfigures workload balancing monitoring on this pool
// Version: george
func (pool) DeconfigureWlb1(session *Session) (err error) {
	method := "pool.deconfigure_wlb"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg)
	return
}

// AsyncDeconfigureWlb1: Permanently deconfigures workload balancing monitoring on this pool
// Version: george
func (pool) AsyncDeconfigureWlb1(session *Session) (retval TaskRef, err error) {
	method := "Async.pool.deconfigure_wlb"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// InitializeWlb: Initializes workload balancing monitoring on this pool with the specified wlb server
// Version: george
func (pool) InitializeWlb(session *Session, wlbURL string, wlbUsername string, wlbPassword string, xenserverUsername string, xenserverPassword string) (err error) {
	method := "pool.initialize_wlb"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	wlbURLArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "wlb_url"), wlbURL)
	if err != nil {
		return
	}
	wlbUsernameArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "wlb_username"), wlbUsername)
	if err != nil {
		return
	}
	wlbPasswordArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "wlb_password"), wlbPassword)
	if err != nil {
		return
	}
	xenserverUsernameArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "xenserver_username"), xenserverUsername)
	if err != nil {
		return
	}
	xenserverPasswordArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "xenserver_password"), xenserverPassword)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, wlbURLArg, wlbUsernameArg, wlbPasswordArg, xenserverUsernameArg, xenserverPasswordArg)
	return
}

// AsyncInitializeWlb: Initializes workload balancing monitoring on this pool with the specified wlb server
// Version: george
func (pool) AsyncInitializeWlb(session *Session, wlbURL string, wlbUsername string, wlbPassword string, xenserverUsername string, xenserverPassword string) (retval TaskRef, err error) {
	method := "Async.pool.initialize_wlb"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	wlbURLArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "wlb_url"), wlbURL)
	if err != nil {
		return
	}
	wlbUsernameArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "wlb_username"), wlbUsername)
	if err != nil {
		return
	}
	wlbPasswordArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "wlb_password"), wlbPassword)
	if err != nil {
		return
	}
	xenserverUsernameArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "xenserver_username"), xenserverUsername)
	if err != nil {
		return
	}
	xenserverPasswordArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "xenserver_password"), xenserverPassword)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, wlbURLArg, wlbUsernameArg, wlbPasswordArg, xenserverUsernameArg, xenserverPasswordArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// InitializeWlb6: Initializes workload balancing monitoring on this pool with the specified wlb server
// Version: george
func (pool) InitializeWlb6(session *Session, wlbURL string, wlbUsername string, wlbPassword string, xenserverUsername string, xenserverPassword string) (err error) {
	method := "pool.initialize_wlb"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	wlbURLArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "wlb_url"), wlbURL)
	if err != nil {
		return
	}
	wlbUsernameArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "wlb_username"), wlbUsername)
	if err != nil {
		return
	}
	wlbPasswordArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "wlb_password"), wlbPassword)
	if err != nil {
		return
	}
	xenserverUsernameArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "xenserver_username"), xenserverUsername)
	if err != nil {
		return
	}
	xenserverPasswordArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "xenserver_password"), xenserverPassword)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, wlbURLArg, wlbUsernameArg, wlbPasswordArg, xenserverUsernameArg, xenserverPasswordArg)
	return
}

// AsyncInitializeWlb6: Initializes workload balancing monitoring on this pool with the specified wlb server
// Version: george
func (pool) AsyncInitializeWlb6(session *Session, wlbURL string, wlbUsername string, wlbPassword string, xenserverUsername string, xenserverPassword string) (retval TaskRef, err error) {
	method := "Async.pool.initialize_wlb"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	wlbURLArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "wlb_url"), wlbURL)
	if err != nil {
		return
	}
	wlbUsernameArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "wlb_username"), wlbUsername)
	if err != nil {
		return
	}
	wlbPasswordArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "wlb_password"), wlbPassword)
	if err != nil {
		return
	}
	xenserverUsernameArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "xenserver_username"), xenserverUsername)
	if err != nil {
		return
	}
	xenserverPasswordArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "xenserver_password"), xenserverPassword)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, wlbURLArg, wlbUsernameArg, wlbPasswordArg, xenserverUsernameArg, xenserverPasswordArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// DetectNonhomogeneousExternalAuth: This call asynchronously detects if the external authentication configuration in any slave is different from that in the master and raises appropriate alerts
// Version: george
func (pool) DetectNonhomogeneousExternalAuth(session *Session, pool PoolRef) (err error) {
	method := "pool.detect_nonhomogeneous_external_auth"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	poolArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "pool"), pool)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, poolArg)
	return
}

// DetectNonhomogeneousExternalAuth2: This call asynchronously detects if the external authentication configuration in any slave is different from that in the master and raises appropriate alerts
// Version: george
func (pool) DetectNonhomogeneousExternalAuth2(session *Session, pool PoolRef) (err error) {
	method := "pool.detect_nonhomogeneous_external_auth"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	poolArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "pool"), pool)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, poolArg)
	return
}

// ExternalAuthSetLdaps: This call enables or disables LDAPS for external authentication on all hosts in the pool
// Version: 26.16.0
func (pool) ExternalAuthSetLdaps(session *Session, pool PoolRef, ldaps bool, force bool) (err error) {
	method := "pool.external_auth_set_ldaps"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	poolArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "pool"), pool)
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
	_, err = session.client.sendCall(method, sessionIDArg, poolArg, ldapsArg, forceArg)
	return
}

// ExternalAuthSetLdaps4: This call enables or disables LDAPS for external authentication on all hosts in the pool
// Version: 26.16.0
func (pool) ExternalAuthSetLdaps4(session *Session, pool PoolRef, ldaps bool, force bool) (err error) {
	method := "pool.external_auth_set_ldaps"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	poolArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "pool"), pool)
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
	_, err = session.client.sendCall(method, sessionIDArg, poolArg, ldapsArg, forceArg)
	return
}

// DisableExternalAuth: This call disables external authentication on all the hosts of the pool
// Version: george
func (pool) DisableExternalAuth(session *Session, pool PoolRef, config map[string]string) (err error) {
	method := "pool.disable_external_auth"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	poolArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "pool"), pool)
	if err != nil {
		return
	}
	configArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "config"), config)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, poolArg, configArg)
	return
}

// DisableExternalAuth3: This call disables external authentication on all the hosts of the pool
// Version: george
func (pool) DisableExternalAuth3(session *Session, pool PoolRef, config map[string]string) (err error) {
	method := "pool.disable_external_auth"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	poolArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "pool"), pool)
	if err != nil {
		return
	}
	configArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "config"), config)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, poolArg, configArg)
	return
}

// EnableExternalAuth: This call enables external authentication on all the hosts of the pool
// Version: george
func (pool) EnableExternalAuth(session *Session, pool PoolRef, config map[string]string, serviceName string, authType string) (err error) {
	method := "pool.enable_external_auth"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	poolArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "pool"), pool)
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
	_, err = session.client.sendCall(method, sessionIDArg, poolArg, configArg, serviceNameArg, authTypeArg)
	return
}

// EnableExternalAuth5: This call enables external authentication on all the hosts of the pool
// Version: george
func (pool) EnableExternalAuth5(session *Session, pool PoolRef, config map[string]string, serviceName string, authType string) (err error) {
	method := "pool.enable_external_auth"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	poolArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "pool"), pool)
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
	_, err = session.client.sendCall(method, sessionIDArg, poolArg, configArg, serviceNameArg, authTypeArg)
	return
}

// CreateNewBlob: Create a placeholder for a named binary blob of data that is associated with this pool
// Version: tampa
func (pool) CreateNewBlob(session *Session, pool PoolRef, name string, mimeType string, public bool) (retval BlobRef, err error) {
	method := "pool.create_new_blob"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	poolArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "pool"), pool)
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
	result, err := session.client.sendCall(method, sessionIDArg, poolArg, nameArg, mimeTypeArg, publicArg)
	if err != nil {
		return
	}
	retval, err = deserializeBlobRef(method+" -> ", result)
	return
}

// AsyncCreateNewBlob: Create a placeholder for a named binary blob of data that is associated with this pool
// Version: tampa
func (pool) AsyncCreateNewBlob(session *Session, pool PoolRef, name string, mimeType string, public bool) (retval TaskRef, err error) {
	method := "Async.pool.create_new_blob"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	poolArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "pool"), pool)
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
	result, err := session.client.sendCall(method, sessionIDArg, poolArg, nameArg, mimeTypeArg, publicArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// CreateNewBlob5: Create a placeholder for a named binary blob of data that is associated with this pool
// Version: tampa
func (pool) CreateNewBlob5(session *Session, pool PoolRef, name string, mimeType string, public bool) (retval BlobRef, err error) {
	method := "pool.create_new_blob"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	poolArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "pool"), pool)
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
	result, err := session.client.sendCall(method, sessionIDArg, poolArg, nameArg, mimeTypeArg, publicArg)
	if err != nil {
		return
	}
	retval, err = deserializeBlobRef(method+" -> ", result)
	return
}

// AsyncCreateNewBlob5: Create a placeholder for a named binary blob of data that is associated with this pool
// Version: tampa
func (pool) AsyncCreateNewBlob5(session *Session, pool PoolRef, name string, mimeType string, public bool) (retval TaskRef, err error) {
	method := "Async.pool.create_new_blob"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	poolArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "pool"), pool)
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
	result, err := session.client.sendCall(method, sessionIDArg, poolArg, nameArg, mimeTypeArg, publicArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// CreateNewBlob4: Create a placeholder for a named binary blob of data that is associated with this pool
// Version: orlando
func (pool) CreateNewBlob4(session *Session, pool PoolRef, name string, mimeType string) (retval BlobRef, err error) {
	method := "pool.create_new_blob"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	poolArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "pool"), pool)
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
	result, err := session.client.sendCall(method, sessionIDArg, poolArg, nameArg, mimeTypeArg)
	if err != nil {
		return
	}
	retval, err = deserializeBlobRef(method+" -> ", result)
	return
}

// AsyncCreateNewBlob4: Create a placeholder for a named binary blob of data that is associated with this pool
// Version: orlando
func (pool) AsyncCreateNewBlob4(session *Session, pool PoolRef, name string, mimeType string) (retval TaskRef, err error) {
	method := "Async.pool.create_new_blob"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	poolArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "pool"), pool)
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
	result, err := session.client.sendCall(method, sessionIDArg, poolArg, nameArg, mimeTypeArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// SetHaHostFailuresToTolerate: Set the maximum number of host failures to consider in the HA VM restart planner
// Version: orlando
func (pool) SetHaHostFailuresToTolerate(session *Session, self PoolRef, value int) (err error) {
	method := "pool.set_ha_host_failures_to_tolerate"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// AsyncSetHaHostFailuresToTolerate: Set the maximum number of host failures to consider in the HA VM restart planner
// Version: orlando
func (pool) AsyncSetHaHostFailuresToTolerate(session *Session, self PoolRef, value int) (retval TaskRef, err error) {
	method := "Async.pool.set_ha_host_failures_to_tolerate"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetHaHostFailuresToTolerate3: Set the maximum number of host failures to consider in the HA VM restart planner
// Version: orlando
func (pool) SetHaHostFailuresToTolerate3(session *Session, self PoolRef, value int) (err error) {
	method := "pool.set_ha_host_failures_to_tolerate"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// AsyncSetHaHostFailuresToTolerate3: Set the maximum number of host failures to consider in the HA VM restart planner
// Version: orlando
func (pool) AsyncSetHaHostFailuresToTolerate3(session *Session, self PoolRef, value int) (retval TaskRef, err error) {
	method := "Async.pool.set_ha_host_failures_to_tolerate"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// HaComputeVMFailoverPlan: Return a VM failover plan assuming a given subset of hosts fail
// Version: orlando
func (pool) HaComputeVMFailoverPlan(session *Session, failedHosts []HostRef, failedVms []VMRef) (retval map[VMRef]map[string]string, err error) {
	method := "pool.ha_compute_vm_failover_plan"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	failedHostsArg, err := serializeHostRefSet(fmt.Sprintf("%s(%s)", method, "failed_hosts"), failedHosts)
	if err != nil {
		return
	}
	failedVmsArg, err := serializeVMRefSet(fmt.Sprintf("%s(%s)", method, "failed_vms"), failedVms)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, failedHostsArg, failedVmsArg)
	if err != nil {
		return
	}
	retval, err = deserializeVMRefToStringToStringMapMap(method+" -> ", result)
	return
}

// HaComputeVMFailoverPlan3: Return a VM failover plan assuming a given subset of hosts fail
// Version: orlando
func (pool) HaComputeVMFailoverPlan3(session *Session, failedHosts []HostRef, failedVms []VMRef) (retval map[VMRef]map[string]string, err error) {
	method := "pool.ha_compute_vm_failover_plan"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	failedHostsArg, err := serializeHostRefSet(fmt.Sprintf("%s(%s)", method, "failed_hosts"), failedHosts)
	if err != nil {
		return
	}
	failedVmsArg, err := serializeVMRefSet(fmt.Sprintf("%s(%s)", method, "failed_vms"), failedVms)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, failedHostsArg, failedVmsArg)
	if err != nil {
		return
	}
	retval, err = deserializeVMRefToStringToStringMapMap(method+" -> ", result)
	return
}

// HaComputeHypotheticalMaxHostFailuresToTolerate: Returns the maximum number of host failures we could tolerate before we would be unable to restart the provided VMs
// Version: orlando
func (pool) HaComputeHypotheticalMaxHostFailuresToTolerate(session *Session, configuration map[VMRef]string) (retval int, err error) {
	method := "pool.ha_compute_hypothetical_max_host_failures_to_tolerate"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	configurationArg, err := serializeVMRefToStringMap(fmt.Sprintf("%s(%s)", method, "configuration"), configuration)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, configurationArg)
	if err != nil {
		return
	}
	retval, err = deserializeInt(method+" -> ", result)
	return
}

// HaComputeHypotheticalMaxHostFailuresToTolerate2: Returns the maximum number of host failures we could tolerate before we would be unable to restart the provided VMs
// Version: orlando
func (pool) HaComputeHypotheticalMaxHostFailuresToTolerate2(session *Session, configuration map[VMRef]string) (retval int, err error) {
	method := "pool.ha_compute_hypothetical_max_host_failures_to_tolerate"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	configurationArg, err := serializeVMRefToStringMap(fmt.Sprintf("%s(%s)", method, "configuration"), configuration)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, configurationArg)
	if err != nil {
		return
	}
	retval, err = deserializeInt(method+" -> ", result)
	return
}

// HaComputeMaxHostFailuresToTolerate: Returns the maximum number of host failures we could tolerate before we would be unable to restart configured VMs
// Version: orlando
func (pool) HaComputeMaxHostFailuresToTolerate(session *Session) (retval int, err error) {
	method := "pool.ha_compute_max_host_failures_to_tolerate"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializeInt(method+" -> ", result)
	return
}

// HaComputeMaxHostFailuresToTolerate1: Returns the maximum number of host failures we could tolerate before we would be unable to restart configured VMs
// Version: orlando
func (pool) HaComputeMaxHostFailuresToTolerate1(session *Session) (retval int, err error) {
	method := "pool.ha_compute_max_host_failures_to_tolerate"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializeInt(method+" -> ", result)
	return
}

// HaFailoverPlanExists: Returns true if a VM failover plan exists for up to &apos;n&apos; host failures
// Version: orlando
func (pool) HaFailoverPlanExists(session *Session, n int) (retval bool, err error) {
	method := "pool.ha_failover_plan_exists"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	nArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "n"), n)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, nArg)
	if err != nil {
		return
	}
	retval, err = deserializeBool(method+" -> ", result)
	return
}

// HaFailoverPlanExists2: Returns true if a VM failover plan exists for up to &apos;n&apos; host failures
// Version: orlando
func (pool) HaFailoverPlanExists2(session *Session, n int) (retval bool, err error) {
	method := "pool.ha_failover_plan_exists"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	nArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "n"), n)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, nArg)
	if err != nil {
		return
	}
	retval, err = deserializeBool(method+" -> ", result)
	return
}

// HaPreventRestartsFor: When this call returns the VM restart logic will not run for the requested number of seconds. If the argument is zero then the restart thread is immediately unblocked
// Version: orlando-update-1
func (pool) HaPreventRestartsFor(session *Session, seconds int) (err error) {
	method := "pool.ha_prevent_restarts_for"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	secondsArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "seconds"), seconds)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, secondsArg)
	return
}

// HaPreventRestartsFor2: When this call returns the VM restart logic will not run for the requested number of seconds. If the argument is zero then the restart thread is immediately unblocked
// Version: orlando-update-1
func (pool) HaPreventRestartsFor2(session *Session, seconds int) (err error) {
	method := "pool.ha_prevent_restarts_for"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	secondsArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "seconds"), seconds)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, secondsArg)
	return
}

// DesignateNewMaster: Perform an orderly handover of the role of master to the referenced host.
// Version: miami
func (pool) DesignateNewMaster(session *Session, host HostRef) (err error) {
	method := "pool.designate_new_master"
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

// AsyncDesignateNewMaster: Perform an orderly handover of the role of master to the referenced host.
// Version: miami
func (pool) AsyncDesignateNewMaster(session *Session, host HostRef) (retval TaskRef, err error) {
	method := "Async.pool.designate_new_master"
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

// DesignateNewMaster2: Perform an orderly handover of the role of master to the referenced host.
// Version: miami
func (pool) DesignateNewMaster2(session *Session, host HostRef) (err error) {
	method := "pool.designate_new_master"
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

// AsyncDesignateNewMaster2: Perform an orderly handover of the role of master to the referenced host.
// Version: miami
func (pool) AsyncDesignateNewMaster2(session *Session, host HostRef) (retval TaskRef, err error) {
	method := "Async.pool.designate_new_master"
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

// SyncDatabase: Forcibly synchronise the database now
// Version: rio
func (pool) SyncDatabase(session *Session) (err error) {
	method := "pool.sync_database"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg)
	return
}

// AsyncSyncDatabase: Forcibly synchronise the database now
// Version: rio
func (pool) AsyncSyncDatabase(session *Session) (retval TaskRef, err error) {
	method := "Async.pool.sync_database"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// SyncDatabase1: Forcibly synchronise the database now
// Version: rio
func (pool) SyncDatabase1(session *Session) (err error) {
	method := "pool.sync_database"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg)
	return
}

// AsyncSyncDatabase1: Forcibly synchronise the database now
// Version: rio
func (pool) AsyncSyncDatabase1(session *Session) (retval TaskRef, err error) {
	method := "Async.pool.sync_database"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// DisableHa: Turn off High Availability mode
// Version: miami
func (pool) DisableHa(session *Session) (err error) {
	method := "pool.disable_ha"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg)
	return
}

// AsyncDisableHa: Turn off High Availability mode
// Version: miami
func (pool) AsyncDisableHa(session *Session) (retval TaskRef, err error) {
	method := "Async.pool.disable_ha"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// DisableHa1: Turn off High Availability mode
// Version: miami
func (pool) DisableHa1(session *Session) (err error) {
	method := "pool.disable_ha"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg)
	return
}

// AsyncDisableHa1: Turn off High Availability mode
// Version: miami
func (pool) AsyncDisableHa1(session *Session) (retval TaskRef, err error) {
	method := "Async.pool.disable_ha"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// EnableHa: Turn on High Availability mode
// Version: miami
func (pool) EnableHa(session *Session, heartbeatSrs []SRRef, configuration map[string]string) (err error) {
	method := "pool.enable_ha"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	heartbeatSrsArg, err := serializeSRRefSet(fmt.Sprintf("%s(%s)", method, "heartbeat_srs"), heartbeatSrs)
	if err != nil {
		return
	}
	configurationArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "configuration"), configuration)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, heartbeatSrsArg, configurationArg)
	return
}

// AsyncEnableHa: Turn on High Availability mode
// Version: miami
func (pool) AsyncEnableHa(session *Session, heartbeatSrs []SRRef, configuration map[string]string) (retval TaskRef, err error) {
	method := "Async.pool.enable_ha"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	heartbeatSrsArg, err := serializeSRRefSet(fmt.Sprintf("%s(%s)", method, "heartbeat_srs"), heartbeatSrs)
	if err != nil {
		return
	}
	configurationArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "configuration"), configuration)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, heartbeatSrsArg, configurationArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// EnableHa3: Turn on High Availability mode
// Version: miami
func (pool) EnableHa3(session *Session, heartbeatSrs []SRRef, configuration map[string]string) (err error) {
	method := "pool.enable_ha"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	heartbeatSrsArg, err := serializeSRRefSet(fmt.Sprintf("%s(%s)", method, "heartbeat_srs"), heartbeatSrs)
	if err != nil {
		return
	}
	configurationArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "configuration"), configuration)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, heartbeatSrsArg, configurationArg)
	return
}

// AsyncEnableHa3: Turn on High Availability mode
// Version: miami
func (pool) AsyncEnableHa3(session *Session, heartbeatSrs []SRRef, configuration map[string]string) (retval TaskRef, err error) {
	method := "Async.pool.enable_ha"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	heartbeatSrsArg, err := serializeSRRefSet(fmt.Sprintf("%s(%s)", method, "heartbeat_srs"), heartbeatSrs)
	if err != nil {
		return
	}
	configurationArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "configuration"), configuration)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, heartbeatSrsArg, configurationArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// CreateVLANFromPIF: Create a pool-wide VLAN by taking the PIF.
// Version: rio
//
// Errors:
// VLAN_TAG_INVALID - You tried to create a VLAN, but the tag you gave was invalid -- it must be between 0 and 4094. The parameter echoes the VLAN tag you gave.
func (pool) CreateVLANFromPIF(session *Session, pif PIFRef, network NetworkRef, vlan int) (retval []PIFRef, err error) {
	method := "pool.create_VLAN_from_PIF"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	pifArg, err := serializePIFRef(fmt.Sprintf("%s(%s)", method, "pif"), pif)
	if err != nil {
		return
	}
	networkArg, err := serializeNetworkRef(fmt.Sprintf("%s(%s)", method, "network"), network)
	if err != nil {
		return
	}
	vlanArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "VLAN"), vlan)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, pifArg, networkArg, vlanArg)
	if err != nil {
		return
	}
	retval, err = deserializePIFRefSet(method+" -> ", result)
	return
}

// AsyncCreateVLANFromPIF: Create a pool-wide VLAN by taking the PIF.
// Version: rio
//
// Errors:
// VLAN_TAG_INVALID - You tried to create a VLAN, but the tag you gave was invalid -- it must be between 0 and 4094. The parameter echoes the VLAN tag you gave.
func (pool) AsyncCreateVLANFromPIF(session *Session, pif PIFRef, network NetworkRef, vlan int) (retval TaskRef, err error) {
	method := "Async.pool.create_VLAN_from_PIF"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	pifArg, err := serializePIFRef(fmt.Sprintf("%s(%s)", method, "pif"), pif)
	if err != nil {
		return
	}
	networkArg, err := serializeNetworkRef(fmt.Sprintf("%s(%s)", method, "network"), network)
	if err != nil {
		return
	}
	vlanArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "VLAN"), vlan)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, pifArg, networkArg, vlanArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// CreateVLANFromPIF4: Create a pool-wide VLAN by taking the PIF.
// Version: rio
//
// Errors:
// VLAN_TAG_INVALID - You tried to create a VLAN, but the tag you gave was invalid -- it must be between 0 and 4094. The parameter echoes the VLAN tag you gave.
func (pool) CreateVLANFromPIF4(session *Session, pif PIFRef, network NetworkRef, vlan int) (retval []PIFRef, err error) {
	method := "pool.create_VLAN_from_PIF"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	pifArg, err := serializePIFRef(fmt.Sprintf("%s(%s)", method, "pif"), pif)
	if err != nil {
		return
	}
	networkArg, err := serializeNetworkRef(fmt.Sprintf("%s(%s)", method, "network"), network)
	if err != nil {
		return
	}
	vlanArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "VLAN"), vlan)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, pifArg, networkArg, vlanArg)
	if err != nil {
		return
	}
	retval, err = deserializePIFRefSet(method+" -> ", result)
	return
}

// AsyncCreateVLANFromPIF4: Create a pool-wide VLAN by taking the PIF.
// Version: rio
//
// Errors:
// VLAN_TAG_INVALID - You tried to create a VLAN, but the tag you gave was invalid -- it must be between 0 and 4094. The parameter echoes the VLAN tag you gave.
func (pool) AsyncCreateVLANFromPIF4(session *Session, pif PIFRef, network NetworkRef, vlan int) (retval TaskRef, err error) {
	method := "Async.pool.create_VLAN_from_PIF"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	pifArg, err := serializePIFRef(fmt.Sprintf("%s(%s)", method, "pif"), pif)
	if err != nil {
		return
	}
	networkArg, err := serializeNetworkRef(fmt.Sprintf("%s(%s)", method, "network"), network)
	if err != nil {
		return
	}
	vlanArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "VLAN"), vlan)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, pifArg, networkArg, vlanArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// ManagementReconfigure: Reconfigure the management network interface for all Hosts in the Pool
// Version: inverness
//
// Errors:
// HA_IS_ENABLED - The operation could not be performed because HA is enabled on the Pool
// PIF_NOT_PRESENT - This host has no PIF on the given network.
// CANNOT_PLUG_BOND_SLAVE - This PIF is a bond member and cannot be plugged.
// PIF_INCOMPATIBLE_PRIMARY_ADDRESS_TYPE - The primary address types are not compatible
// PIF_HAS_NO_NETWORK_CONFIGURATION - PIF has no IP configuration (mode currently set to &apos;none&apos;)
// PIF_HAS_NO_V6_NETWORK_CONFIGURATION - PIF has no IPv6 configuration (mode currently set to &apos;none&apos;)
func (pool) ManagementReconfigure(session *Session, network NetworkRef) (err error) {
	method := "pool.management_reconfigure"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	networkArg, err := serializeNetworkRef(fmt.Sprintf("%s(%s)", method, "network"), network)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, networkArg)
	return
}

// AsyncManagementReconfigure: Reconfigure the management network interface for all Hosts in the Pool
// Version: inverness
//
// Errors:
// HA_IS_ENABLED - The operation could not be performed because HA is enabled on the Pool
// PIF_NOT_PRESENT - This host has no PIF on the given network.
// CANNOT_PLUG_BOND_SLAVE - This PIF is a bond member and cannot be plugged.
// PIF_INCOMPATIBLE_PRIMARY_ADDRESS_TYPE - The primary address types are not compatible
// PIF_HAS_NO_NETWORK_CONFIGURATION - PIF has no IP configuration (mode currently set to &apos;none&apos;)
// PIF_HAS_NO_V6_NETWORK_CONFIGURATION - PIF has no IPv6 configuration (mode currently set to &apos;none&apos;)
func (pool) AsyncManagementReconfigure(session *Session, network NetworkRef) (retval TaskRef, err error) {
	method := "Async.pool.management_reconfigure"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	networkArg, err := serializeNetworkRef(fmt.Sprintf("%s(%s)", method, "network"), network)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, networkArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// ManagementReconfigure2: Reconfigure the management network interface for all Hosts in the Pool
// Version: inverness
//
// Errors:
// HA_IS_ENABLED - The operation could not be performed because HA is enabled on the Pool
// PIF_NOT_PRESENT - This host has no PIF on the given network.
// CANNOT_PLUG_BOND_SLAVE - This PIF is a bond member and cannot be plugged.
// PIF_INCOMPATIBLE_PRIMARY_ADDRESS_TYPE - The primary address types are not compatible
// PIF_HAS_NO_NETWORK_CONFIGURATION - PIF has no IP configuration (mode currently set to &apos;none&apos;)
// PIF_HAS_NO_V6_NETWORK_CONFIGURATION - PIF has no IPv6 configuration (mode currently set to &apos;none&apos;)
func (pool) ManagementReconfigure2(session *Session, network NetworkRef) (err error) {
	method := "pool.management_reconfigure"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	networkArg, err := serializeNetworkRef(fmt.Sprintf("%s(%s)", method, "network"), network)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, networkArg)
	return
}

// AsyncManagementReconfigure2: Reconfigure the management network interface for all Hosts in the Pool
// Version: inverness
//
// Errors:
// HA_IS_ENABLED - The operation could not be performed because HA is enabled on the Pool
// PIF_NOT_PRESENT - This host has no PIF on the given network.
// CANNOT_PLUG_BOND_SLAVE - This PIF is a bond member and cannot be plugged.
// PIF_INCOMPATIBLE_PRIMARY_ADDRESS_TYPE - The primary address types are not compatible
// PIF_HAS_NO_NETWORK_CONFIGURATION - PIF has no IP configuration (mode currently set to &apos;none&apos;)
// PIF_HAS_NO_V6_NETWORK_CONFIGURATION - PIF has no IPv6 configuration (mode currently set to &apos;none&apos;)
func (pool) AsyncManagementReconfigure2(session *Session, network NetworkRef) (retval TaskRef, err error) {
	method := "Async.pool.management_reconfigure"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	networkArg, err := serializeNetworkRef(fmt.Sprintf("%s(%s)", method, "network"), network)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, networkArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// CreateVLAN: Create PIFs, mapping a network to the same physical interface/VLAN on each host. This call is deprecated: use Pool.create_VLAN_from_PIF instead.
// Version: rio
//
// Errors:
// VLAN_TAG_INVALID - You tried to create a VLAN, but the tag you gave was invalid -- it must be between 0 and 4094. The parameter echoes the VLAN tag you gave.
func (pool) CreateVLAN(session *Session, device string, network NetworkRef, vlan int) (retval []PIFRef, err error) {
	method := "pool.create_VLAN"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	deviceArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "device"), device)
	if err != nil {
		return
	}
	networkArg, err := serializeNetworkRef(fmt.Sprintf("%s(%s)", method, "network"), network)
	if err != nil {
		return
	}
	vlanArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "VLAN"), vlan)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, deviceArg, networkArg, vlanArg)
	if err != nil {
		return
	}
	retval, err = deserializePIFRefSet(method+" -> ", result)
	return
}

// AsyncCreateVLAN: Create PIFs, mapping a network to the same physical interface/VLAN on each host. This call is deprecated: use Pool.create_VLAN_from_PIF instead.
// Version: rio
//
// Errors:
// VLAN_TAG_INVALID - You tried to create a VLAN, but the tag you gave was invalid -- it must be between 0 and 4094. The parameter echoes the VLAN tag you gave.
func (pool) AsyncCreateVLAN(session *Session, device string, network NetworkRef, vlan int) (retval TaskRef, err error) {
	method := "Async.pool.create_VLAN"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	deviceArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "device"), device)
	if err != nil {
		return
	}
	networkArg, err := serializeNetworkRef(fmt.Sprintf("%s(%s)", method, "network"), network)
	if err != nil {
		return
	}
	vlanArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "VLAN"), vlan)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, deviceArg, networkArg, vlanArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// CreateVLAN4: Create PIFs, mapping a network to the same physical interface/VLAN on each host. This call is deprecated: use Pool.create_VLAN_from_PIF instead.
// Version: rio
//
// Errors:
// VLAN_TAG_INVALID - You tried to create a VLAN, but the tag you gave was invalid -- it must be between 0 and 4094. The parameter echoes the VLAN tag you gave.
func (pool) CreateVLAN4(session *Session, device string, network NetworkRef, vlan int) (retval []PIFRef, err error) {
	method := "pool.create_VLAN"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	deviceArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "device"), device)
	if err != nil {
		return
	}
	networkArg, err := serializeNetworkRef(fmt.Sprintf("%s(%s)", method, "network"), network)
	if err != nil {
		return
	}
	vlanArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "VLAN"), vlan)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, deviceArg, networkArg, vlanArg)
	if err != nil {
		return
	}
	retval, err = deserializePIFRefSet(method+" -> ", result)
	return
}

// AsyncCreateVLAN4: Create PIFs, mapping a network to the same physical interface/VLAN on each host. This call is deprecated: use Pool.create_VLAN_from_PIF instead.
// Version: rio
//
// Errors:
// VLAN_TAG_INVALID - You tried to create a VLAN, but the tag you gave was invalid -- it must be between 0 and 4094. The parameter echoes the VLAN tag you gave.
func (pool) AsyncCreateVLAN4(session *Session, device string, network NetworkRef, vlan int) (retval TaskRef, err error) {
	method := "Async.pool.create_VLAN"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	deviceArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "device"), device)
	if err != nil {
		return
	}
	networkArg, err := serializeNetworkRef(fmt.Sprintf("%s(%s)", method, "network"), network)
	if err != nil {
		return
	}
	vlanArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "VLAN"), vlan)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, deviceArg, networkArg, vlanArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// RecoverSlaves: Instruct a pool master, M, to try and contact its slaves and, if slaves are in emergency mode, reset their master address to M.
// Version: rio
func (pool) RecoverSlaves(session *Session) (retval []HostRef, err error) {
	method := "pool.recover_slaves"
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

// AsyncRecoverSlaves: Instruct a pool master, M, to try and contact its slaves and, if slaves are in emergency mode, reset their master address to M.
// Version: rio
func (pool) AsyncRecoverSlaves(session *Session) (retval TaskRef, err error) {
	method := "Async.pool.recover_slaves"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// RecoverSlaves1: Instruct a pool master, M, to try and contact its slaves and, if slaves are in emergency mode, reset their master address to M.
// Version: rio
func (pool) RecoverSlaves1(session *Session) (retval []HostRef, err error) {
	method := "pool.recover_slaves"
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

// AsyncRecoverSlaves1: Instruct a pool master, M, to try and contact its slaves and, if slaves are in emergency mode, reset their master address to M.
// Version: rio
func (pool) AsyncRecoverSlaves1(session *Session) (retval TaskRef, err error) {
	method := "Async.pool.recover_slaves"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// EmergencyResetMaster: Instruct a slave already in a pool that the master has changed
// Version: rio
func (pool) EmergencyResetMaster(session *Session, masterAddress string) (err error) {
	method := "pool.emergency_reset_master"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	masterAddressArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "master_address"), masterAddress)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, masterAddressArg)
	return
}

// EmergencyResetMaster2: Instruct a slave already in a pool that the master has changed
// Version: rio
func (pool) EmergencyResetMaster2(session *Session, masterAddress string) (err error) {
	method := "pool.emergency_reset_master"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	masterAddressArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "master_address"), masterAddress)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, masterAddressArg)
	return
}

// EmergencyTransitionToMaster: Instruct host that&apos;s currently a slave to transition to being master
// Version: rio
func (pool) EmergencyTransitionToMaster(session *Session) (err error) {
	method := "pool.emergency_transition_to_master"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg)
	return
}

// EmergencyTransitionToMaster1: Instruct host that&apos;s currently a slave to transition to being master
// Version: rio
func (pool) EmergencyTransitionToMaster1(session *Session) (err error) {
	method := "pool.emergency_transition_to_master"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg)
	return
}

// Eject: Instruct a pool master to eject a host from the pool
// Version: rio
func (pool) Eject(session *Session, host HostRef) (err error) {
	method := "pool.eject"
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

// AsyncEject: Instruct a pool master to eject a host from the pool
// Version: rio
func (pool) AsyncEject(session *Session, host HostRef) (retval TaskRef, err error) {
	method := "Async.pool.eject"
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

// Eject2: Instruct a pool master to eject a host from the pool
// Version: rio
func (pool) Eject2(session *Session, host HostRef) (err error) {
	method := "pool.eject"
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

// AsyncEject2: Instruct a pool master to eject a host from the pool
// Version: rio
func (pool) AsyncEject2(session *Session, host HostRef) (retval TaskRef, err error) {
	method := "Async.pool.eject"
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

// JoinForce: Instruct host to join a new pool
// Version: rio
func (pool) JoinForce(session *Session, masterAddress string, masterUsername string, masterPassword string) (err error) {
	method := "pool.join_force"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	masterAddressArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "master_address"), masterAddress)
	if err != nil {
		return
	}
	masterUsernameArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "master_username"), masterUsername)
	if err != nil {
		return
	}
	masterPasswordArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "master_password"), masterPassword)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, masterAddressArg, masterUsernameArg, masterPasswordArg)
	return
}

// AsyncJoinForce: Instruct host to join a new pool
// Version: rio
func (pool) AsyncJoinForce(session *Session, masterAddress string, masterUsername string, masterPassword string) (retval TaskRef, err error) {
	method := "Async.pool.join_force"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	masterAddressArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "master_address"), masterAddress)
	if err != nil {
		return
	}
	masterUsernameArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "master_username"), masterUsername)
	if err != nil {
		return
	}
	masterPasswordArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "master_password"), masterPassword)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, masterAddressArg, masterUsernameArg, masterPasswordArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// JoinForce4: Instruct host to join a new pool
// Version: rio
func (pool) JoinForce4(session *Session, masterAddress string, masterUsername string, masterPassword string) (err error) {
	method := "pool.join_force"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	masterAddressArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "master_address"), masterAddress)
	if err != nil {
		return
	}
	masterUsernameArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "master_username"), masterUsername)
	if err != nil {
		return
	}
	masterPasswordArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "master_password"), masterPassword)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, masterAddressArg, masterUsernameArg, masterPasswordArg)
	return
}

// AsyncJoinForce4: Instruct host to join a new pool
// Version: rio
func (pool) AsyncJoinForce4(session *Session, masterAddress string, masterUsername string, masterPassword string) (retval TaskRef, err error) {
	method := "Async.pool.join_force"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	masterAddressArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "master_address"), masterAddress)
	if err != nil {
		return
	}
	masterUsernameArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "master_username"), masterUsername)
	if err != nil {
		return
	}
	masterPasswordArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "master_password"), masterPassword)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, masterAddressArg, masterUsernameArg, masterPasswordArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// Join: Instruct host to join a new pool
// Version: rio
//
// Errors:
// JOINING_HOST_CANNOT_CONTAIN_SHARED_SRS - The server joining the pool cannot contain any shared storage.
func (pool) Join(session *Session, masterAddress string, masterUsername string, masterPassword string) (err error) {
	method := "pool.join"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	masterAddressArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "master_address"), masterAddress)
	if err != nil {
		return
	}
	masterUsernameArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "master_username"), masterUsername)
	if err != nil {
		return
	}
	masterPasswordArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "master_password"), masterPassword)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, masterAddressArg, masterUsernameArg, masterPasswordArg)
	return
}

// AsyncJoin: Instruct host to join a new pool
// Version: rio
//
// Errors:
// JOINING_HOST_CANNOT_CONTAIN_SHARED_SRS - The server joining the pool cannot contain any shared storage.
func (pool) AsyncJoin(session *Session, masterAddress string, masterUsername string, masterPassword string) (retval TaskRef, err error) {
	method := "Async.pool.join"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	masterAddressArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "master_address"), masterAddress)
	if err != nil {
		return
	}
	masterUsernameArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "master_username"), masterUsername)
	if err != nil {
		return
	}
	masterPasswordArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "master_password"), masterPassword)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, masterAddressArg, masterUsernameArg, masterPasswordArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// Join4: Instruct host to join a new pool
// Version: rio
//
// Errors:
// JOINING_HOST_CANNOT_CONTAIN_SHARED_SRS - The server joining the pool cannot contain any shared storage.
func (pool) Join4(session *Session, masterAddress string, masterUsername string, masterPassword string) (err error) {
	method := "pool.join"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	masterAddressArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "master_address"), masterAddress)
	if err != nil {
		return
	}
	masterUsernameArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "master_username"), masterUsername)
	if err != nil {
		return
	}
	masterPasswordArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "master_password"), masterPassword)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, masterAddressArg, masterUsernameArg, masterPasswordArg)
	return
}

// AsyncJoin4: Instruct host to join a new pool
// Version: rio
//
// Errors:
// JOINING_HOST_CANNOT_CONTAIN_SHARED_SRS - The server joining the pool cannot contain any shared storage.
func (pool) AsyncJoin4(session *Session, masterAddress string, masterUsername string, masterPassword string) (retval TaskRef, err error) {
	method := "Async.pool.join"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	masterAddressArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "master_address"), masterAddress)
	if err != nil {
		return
	}
	masterUsernameArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "master_username"), masterUsername)
	if err != nil {
		return
	}
	masterPasswordArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "master_password"), masterPassword)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, masterAddressArg, masterUsernameArg, masterPasswordArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// SetAutoUpdateVMSecurebootCertificates: Set the auto_update_vm_secureboot_certificates field of the given pool.
// Version: rio
func (pool) SetAutoUpdateVMSecurebootCertificates(session *Session, self PoolRef, value bool) (err error) {
	method := "pool.set_auto_update_vm_secureboot_certificates"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetAutoUpdateVMSecurebootCertificates3: Set the auto_update_vm_secureboot_certificates field of the given pool.
// Version: rio
func (pool) SetAutoUpdateVMSecurebootCertificates3(session *Session, self PoolRef, value bool) (err error) {
	method := "pool.set_auto_update_vm_secureboot_certificates"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetAutoUpdateVMSecurebootCertificates2: Set the auto_update_vm_secureboot_certificates field of the given pool.
// Version: closed
func (pool) SetAutoUpdateVMSecurebootCertificates2(session *Session, value bool) (err error) {
	method := "pool.set_auto_update_vm_secureboot_certificates"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	valueArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, valueArg)
	return
}

// SetVMConsoleIdleTimeout: Set the vm_console_idle_timeout field of the given pool.
// Version: rio
func (pool) SetVMConsoleIdleTimeout(session *Session, self PoolRef, value int) (err error) {
	method := "pool.set_vm_console_idle_timeout"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetVMConsoleIdleTimeout3: Set the vm_console_idle_timeout field of the given pool.
// Version: rio
func (pool) SetVMConsoleIdleTimeout3(session *Session, self PoolRef, value int) (err error) {
	method := "pool.set_vm_console_idle_timeout"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetVMConsoleIdleTimeout2: Set the vm_console_idle_timeout field of the given pool.
// Version: closed
func (pool) SetVMConsoleIdleTimeout2(session *Session, value int) (err error) {
	method := "pool.set_vm_console_idle_timeout"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	valueArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, valueArg)
	return
}

// SetLimitConsoleSessions: Set the limit_console_sessions field of the given pool.
// Version: rio
func (pool) SetLimitConsoleSessions(session *Session, self PoolRef, value bool) (err error) {
	method := "pool.set_limit_console_sessions"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetLimitConsoleSessions3: Set the limit_console_sessions field of the given pool.
// Version: rio
func (pool) SetLimitConsoleSessions3(session *Session, self PoolRef, value bool) (err error) {
	method := "pool.set_limit_console_sessions"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetLimitConsoleSessions2: Set the limit_console_sessions field of the given pool.
// Version: closed
func (pool) SetLimitConsoleSessions2(session *Session, value bool) (err error) {
	method := "pool.set_limit_console_sessions"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	valueArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, valueArg)
	return
}

// SetHaRebootVMOnInternalShutdown: Set the ha_reboot_vm_on_internal_shutdown field of the given pool.
// Version: rio
func (pool) SetHaRebootVMOnInternalShutdown(session *Session, self PoolRef, value bool) (err error) {
	method := "pool.set_ha_reboot_vm_on_internal_shutdown"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetHaRebootVMOnInternalShutdown3: Set the ha_reboot_vm_on_internal_shutdown field of the given pool.
// Version: rio
func (pool) SetHaRebootVMOnInternalShutdown3(session *Session, self PoolRef, value bool) (err error) {
	method := "pool.set_ha_reboot_vm_on_internal_shutdown"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetHaRebootVMOnInternalShutdown2: Set the ha_reboot_vm_on_internal_shutdown field of the given pool.
// Version: closed
func (pool) SetHaRebootVMOnInternalShutdown2(session *Session, value bool) (err error) {
	method := "pool.set_ha_reboot_vm_on_internal_shutdown"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	valueArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, valueArg)
	return
}

// RemoveFromLicenseServer: Remove the given key and its corresponding value from the license_server field of the given pool.  If the key is not in that Map, then do nothing.
// Version: rio
func (pool) RemoveFromLicenseServer(session *Session, self PoolRef, key string) (err error) {
	method := "pool.remove_from_license_server"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// RemoveFromLicenseServer3: Remove the given key and its corresponding value from the license_server field of the given pool.  If the key is not in that Map, then do nothing.
// Version: rio
func (pool) RemoveFromLicenseServer3(session *Session, self PoolRef, key string) (err error) {
	method := "pool.remove_from_license_server"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// RemoveFromLicenseServer2: Remove the given key and its corresponding value from the license_server field of the given pool.  If the key is not in that Map, then do nothing.
// Version: closed
func (pool) RemoveFromLicenseServer2(session *Session, key string) (err error) {
	method := "pool.remove_from_license_server"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	keyArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "key"), key)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, keyArg)
	return
}

// AddToLicenseServer: Add the given key-value pair to the license_server field of the given pool.
// Version: rio
func (pool) AddToLicenseServer(session *Session, self PoolRef, key string, value string) (err error) {
	method := "pool.add_to_license_server"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// AddToLicenseServer4: Add the given key-value pair to the license_server field of the given pool.
// Version: rio
func (pool) AddToLicenseServer4(session *Session, self PoolRef, key string, value string) (err error) {
	method := "pool.add_to_license_server"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// AddToLicenseServer3: Add the given key-value pair to the license_server field of the given pool.
// Version: closed
func (pool) AddToLicenseServer3(session *Session, key string, value string) (err error) {
	method := "pool.add_to_license_server"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
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
	_, err = session.client.sendCall(method, sessionIDArg, keyArg, valueArg)
	return
}

// SetLicenseServer: Set the license_server field of the given pool.
// Version: rio
func (pool) SetLicenseServer(session *Session, self PoolRef, value map[string]string) (err error) {
	method := "pool.set_license_server"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetLicenseServer3: Set the license_server field of the given pool.
// Version: rio
func (pool) SetLicenseServer3(session *Session, self PoolRef, value map[string]string) (err error) {
	method := "pool.set_license_server"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetLicenseServer2: Set the license_server field of the given pool.
// Version: closed
func (pool) SetLicenseServer2(session *Session, value map[string]string) (err error) {
	method := "pool.set_license_server"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	valueArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, valueArg)
	return
}

// SetCoordinatorBias: Set the coordinator_bias field of the given pool.
// Version: rio
func (pool) SetCoordinatorBias(session *Session, self PoolRef, value bool) (err error) {
	method := "pool.set_coordinator_bias"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetCoordinatorBias3: Set the coordinator_bias field of the given pool.
// Version: rio
func (pool) SetCoordinatorBias3(session *Session, self PoolRef, value bool) (err error) {
	method := "pool.set_coordinator_bias"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetMigrationCompression: Set the migration_compression field of the given pool.
// Version: rio
func (pool) SetMigrationCompression(session *Session, self PoolRef, value bool) (err error) {
	method := "pool.set_migration_compression"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetMigrationCompression3: Set the migration_compression field of the given pool.
// Version: rio
func (pool) SetMigrationCompression3(session *Session, self PoolRef, value bool) (err error) {
	method := "pool.set_migration_compression"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetMigrationCompression2: Set the migration_compression field of the given pool.
// Version: closed
func (pool) SetMigrationCompression2(session *Session, value bool) (err error) {
	method := "pool.set_migration_compression"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	valueArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, valueArg)
	return
}

// SetIsPsrPending: Set the is_psr_pending field of the given pool.
// Version: stockholm_psr
func (pool) SetIsPsrPending(session *Session, self PoolRef, value bool) (err error) {
	method := "pool.set_is_psr_pending"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetIsPsrPending3: Set the is_psr_pending field of the given pool.
// Version: stockholm_psr
func (pool) SetIsPsrPending3(session *Session, self PoolRef, value bool) (err error) {
	method := "pool.set_is_psr_pending"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetIsPsrPending2: Set the is_psr_pending field of the given pool.
// Version: rio
func (pool) SetIsPsrPending2(session *Session, self PoolRef) (err error) {
	method := "pool.set_is_psr_pending"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// SetLivePatchingDisabled: Set the live_patching_disabled field of the given pool.
// Version: ely
func (pool) SetLivePatchingDisabled(session *Session, self PoolRef, value bool) (err error) {
	method := "pool.set_live_patching_disabled"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetLivePatchingDisabled3: Set the live_patching_disabled field of the given pool.
// Version: ely
func (pool) SetLivePatchingDisabled3(session *Session, self PoolRef, value bool) (err error) {
	method := "pool.set_live_patching_disabled"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetLivePatchingDisabled2: Set the live_patching_disabled field of the given pool.
// Version: rio
func (pool) SetLivePatchingDisabled2(session *Session, self PoolRef) (err error) {
	method := "pool.set_live_patching_disabled"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// SetPolicyNoVendorDevice: Set the policy_no_vendor_device field of the given pool.
// Version: dundee
func (pool) SetPolicyNoVendorDevice(session *Session, self PoolRef, value bool) (err error) {
	method := "pool.set_policy_no_vendor_device"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetPolicyNoVendorDevice3: Set the policy_no_vendor_device field of the given pool.
// Version: dundee
func (pool) SetPolicyNoVendorDevice3(session *Session, self PoolRef, value bool) (err error) {
	method := "pool.set_policy_no_vendor_device"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetPolicyNoVendorDevice2: Set the policy_no_vendor_device field of the given pool.
// Version: rio
func (pool) SetPolicyNoVendorDevice2(session *Session, self PoolRef) (err error) {
	method := "pool.set_policy_no_vendor_device"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// SetWlbVerifyCert: Set the wlb_verify_cert field of the given pool.
// Version: george
func (pool) SetWlbVerifyCert(session *Session, self PoolRef, value bool) (err error) {
	method := "pool.set_wlb_verify_cert"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetWlbVerifyCert3: Set the wlb_verify_cert field of the given pool.
// Version: george
func (pool) SetWlbVerifyCert3(session *Session, self PoolRef, value bool) (err error) {
	method := "pool.set_wlb_verify_cert"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetWlbVerifyCert2: Set the wlb_verify_cert field of the given pool.
// Version: rio
func (pool) SetWlbVerifyCert2(session *Session, self PoolRef) (err error) {
	method := "pool.set_wlb_verify_cert"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// SetWlbEnabled: Set the wlb_enabled field of the given pool.
// Version: george
func (pool) SetWlbEnabled(session *Session, self PoolRef, value bool) (err error) {
	method := "pool.set_wlb_enabled"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetWlbEnabled3: Set the wlb_enabled field of the given pool.
// Version: george
func (pool) SetWlbEnabled3(session *Session, self PoolRef, value bool) (err error) {
	method := "pool.set_wlb_enabled"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetWlbEnabled2: Set the wlb_enabled field of the given pool.
// Version: rio
func (pool) SetWlbEnabled2(session *Session, self PoolRef) (err error) {
	method := "pool.set_wlb_enabled"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// RemoveFromHealthCheckConfig: Remove the given key and its corresponding value from the health_check_config field of the given pool.  If the key is not in that Map, then do nothing.
// Version: dundee
func (pool) RemoveFromHealthCheckConfig(session *Session, self PoolRef, key string) (err error) {
	method := "pool.remove_from_health_check_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// RemoveFromHealthCheckConfig3: Remove the given key and its corresponding value from the health_check_config field of the given pool.  If the key is not in that Map, then do nothing.
// Version: dundee
func (pool) RemoveFromHealthCheckConfig3(session *Session, self PoolRef, key string) (err error) {
	method := "pool.remove_from_health_check_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// RemoveFromHealthCheckConfig2: Remove the given key and its corresponding value from the health_check_config field of the given pool.  If the key is not in that Map, then do nothing.
// Version: rio
func (pool) RemoveFromHealthCheckConfig2(session *Session, self PoolRef) (err error) {
	method := "pool.remove_from_health_check_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// AddToHealthCheckConfig: Add the given key-value pair to the health_check_config field of the given pool.
// Version: dundee
func (pool) AddToHealthCheckConfig(session *Session, self PoolRef, key string, value string) (err error) {
	method := "pool.add_to_health_check_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// AddToHealthCheckConfig4: Add the given key-value pair to the health_check_config field of the given pool.
// Version: dundee
func (pool) AddToHealthCheckConfig4(session *Session, self PoolRef, key string, value string) (err error) {
	method := "pool.add_to_health_check_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// AddToHealthCheckConfig2: Add the given key-value pair to the health_check_config field of the given pool.
// Version: rio
func (pool) AddToHealthCheckConfig2(session *Session, self PoolRef) (err error) {
	method := "pool.add_to_health_check_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// SetHealthCheckConfig: Set the health_check_config field of the given pool.
// Version: dundee
func (pool) SetHealthCheckConfig(session *Session, self PoolRef, value map[string]string) (err error) {
	method := "pool.set_health_check_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetHealthCheckConfig3: Set the health_check_config field of the given pool.
// Version: dundee
func (pool) SetHealthCheckConfig3(session *Session, self PoolRef, value map[string]string) (err error) {
	method := "pool.set_health_check_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetHealthCheckConfig2: Set the health_check_config field of the given pool.
// Version: rio
func (pool) SetHealthCheckConfig2(session *Session, self PoolRef) (err error) {
	method := "pool.set_health_check_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// RemoveFromGuiConfig: Remove the given key and its corresponding value from the gui_config field of the given pool.  If the key is not in that Map, then do nothing.
// Version: orlando
func (pool) RemoveFromGuiConfig(session *Session, self PoolRef, key string) (err error) {
	method := "pool.remove_from_gui_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// RemoveFromGuiConfig3: Remove the given key and its corresponding value from the gui_config field of the given pool.  If the key is not in that Map, then do nothing.
// Version: orlando
func (pool) RemoveFromGuiConfig3(session *Session, self PoolRef, key string) (err error) {
	method := "pool.remove_from_gui_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// RemoveFromGuiConfig2: Remove the given key and its corresponding value from the gui_config field of the given pool.  If the key is not in that Map, then do nothing.
// Version: rio
func (pool) RemoveFromGuiConfig2(session *Session, self PoolRef) (err error) {
	method := "pool.remove_from_gui_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// AddToGuiConfig: Add the given key-value pair to the gui_config field of the given pool.
// Version: orlando
func (pool) AddToGuiConfig(session *Session, self PoolRef, key string, value string) (err error) {
	method := "pool.add_to_gui_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// AddToGuiConfig4: Add the given key-value pair to the gui_config field of the given pool.
// Version: orlando
func (pool) AddToGuiConfig4(session *Session, self PoolRef, key string, value string) (err error) {
	method := "pool.add_to_gui_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// AddToGuiConfig2: Add the given key-value pair to the gui_config field of the given pool.
// Version: rio
func (pool) AddToGuiConfig2(session *Session, self PoolRef) (err error) {
	method := "pool.add_to_gui_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// SetGuiConfig: Set the gui_config field of the given pool.
// Version: orlando
func (pool) SetGuiConfig(session *Session, self PoolRef, value map[string]string) (err error) {
	method := "pool.set_gui_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetGuiConfig3: Set the gui_config field of the given pool.
// Version: orlando
func (pool) SetGuiConfig3(session *Session, self PoolRef, value map[string]string) (err error) {
	method := "pool.set_gui_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetGuiConfig2: Set the gui_config field of the given pool.
// Version: rio
func (pool) SetGuiConfig2(session *Session, self PoolRef) (err error) {
	method := "pool.set_gui_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// RemoveTags: Remove the given value from the tags field of the given pool.  If the value is not in that Set, then do nothing.
// Version: orlando
func (pool) RemoveTags(session *Session, self PoolRef, value string) (err error) {
	method := "pool.remove_tags"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// RemoveTags3: Remove the given value from the tags field of the given pool.  If the value is not in that Set, then do nothing.
// Version: orlando
func (pool) RemoveTags3(session *Session, self PoolRef, value string) (err error) {
	method := "pool.remove_tags"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// RemoveTags2: Remove the given value from the tags field of the given pool.  If the value is not in that Set, then do nothing.
// Version: rio
func (pool) RemoveTags2(session *Session, self PoolRef) (err error) {
	method := "pool.remove_tags"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// AddTags: Add the given value to the tags field of the given pool.  If the value is already in that Set, then do nothing.
// Version: orlando
func (pool) AddTags(session *Session, self PoolRef, value string) (err error) {
	method := "pool.add_tags"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// AddTags3: Add the given value to the tags field of the given pool.  If the value is already in that Set, then do nothing.
// Version: orlando
func (pool) AddTags3(session *Session, self PoolRef, value string) (err error) {
	method := "pool.add_tags"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// AddTags2: Add the given value to the tags field of the given pool.  If the value is already in that Set, then do nothing.
// Version: rio
func (pool) AddTags2(session *Session, self PoolRef) (err error) {
	method := "pool.add_tags"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// SetTags: Set the tags field of the given pool.
// Version: orlando
func (pool) SetTags(session *Session, self PoolRef, value []string) (err error) {
	method := "pool.set_tags"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetTags3: Set the tags field of the given pool.
// Version: orlando
func (pool) SetTags3(session *Session, self PoolRef, value []string) (err error) {
	method := "pool.set_tags"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetTags2: Set the tags field of the given pool.
// Version: rio
func (pool) SetTags2(session *Session, self PoolRef) (err error) {
	method := "pool.set_tags"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// SetHaAllowOvercommit: Set the ha_allow_overcommit field of the given pool.
// Version: orlando
func (pool) SetHaAllowOvercommit(session *Session, self PoolRef, value bool) (err error) {
	method := "pool.set_ha_allow_overcommit"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetHaAllowOvercommit3: Set the ha_allow_overcommit field of the given pool.
// Version: orlando
func (pool) SetHaAllowOvercommit3(session *Session, self PoolRef, value bool) (err error) {
	method := "pool.set_ha_allow_overcommit"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetHaAllowOvercommit2: Set the ha_allow_overcommit field of the given pool.
// Version: rio
func (pool) SetHaAllowOvercommit2(session *Session, self PoolRef) (err error) {
	method := "pool.set_ha_allow_overcommit"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// RemoveFromOtherConfig: Remove the given key and its corresponding value from the other_config field of the given pool.  If the key is not in that Map, then do nothing.
// Version: rio
func (pool) RemoveFromOtherConfig(session *Session, self PoolRef, key string) (err error) {
	method := "pool.remove_from_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// RemoveFromOtherConfig3: Remove the given key and its corresponding value from the other_config field of the given pool.  If the key is not in that Map, then do nothing.
// Version: rio
func (pool) RemoveFromOtherConfig3(session *Session, self PoolRef, key string) (err error) {
	method := "pool.remove_from_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// AddToOtherConfig: Add the given key-value pair to the other_config field of the given pool.
// Version: rio
func (pool) AddToOtherConfig(session *Session, self PoolRef, key string, value string) (err error) {
	method := "pool.add_to_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// AddToOtherConfig4: Add the given key-value pair to the other_config field of the given pool.
// Version: rio
func (pool) AddToOtherConfig4(session *Session, self PoolRef, key string, value string) (err error) {
	method := "pool.add_to_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetOtherConfig: Set the other_config field of the given pool.
// Version: rio
func (pool) SetOtherConfig(session *Session, self PoolRef, value map[string]string) (err error) {
	method := "pool.set_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetOtherConfig3: Set the other_config field of the given pool.
// Version: rio
func (pool) SetOtherConfig3(session *Session, self PoolRef, value map[string]string) (err error) {
	method := "pool.set_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetCrashDumpSR: Set the crash_dump_SR field of the given pool.
// Version: rio
func (pool) SetCrashDumpSR(session *Session, self PoolRef, value SRRef) (err error) {
	method := "pool.set_crash_dump_SR"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetCrashDumpSR3: Set the crash_dump_SR field of the given pool.
// Version: rio
func (pool) SetCrashDumpSR3(session *Session, self PoolRef, value SRRef) (err error) {
	method := "pool.set_crash_dump_SR"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetSuspendImageSR: Set the suspend_image_SR field of the given pool.
// Version: rio
func (pool) SetSuspendImageSR(session *Session, self PoolRef, value SRRef) (err error) {
	method := "pool.set_suspend_image_SR"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetSuspendImageSR3: Set the suspend_image_SR field of the given pool.
// Version: rio
func (pool) SetSuspendImageSR3(session *Session, self PoolRef, value SRRef) (err error) {
	method := "pool.set_suspend_image_SR"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetDefaultSR: Set the default_SR field of the given pool.
// Version: rio
func (pool) SetDefaultSR(session *Session, self PoolRef, value SRRef) (err error) {
	method := "pool.set_default_SR"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetDefaultSR3: Set the default_SR field of the given pool.
// Version: rio
func (pool) SetDefaultSR3(session *Session, self PoolRef, value SRRef) (err error) {
	method := "pool.set_default_SR"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetNameDescription: Set the name_description field of the given pool.
// Version: rio
func (pool) SetNameDescription(session *Session, self PoolRef, value string) (err error) {
	method := "pool.set_name_description"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetNameDescription3: Set the name_description field of the given pool.
// Version: rio
func (pool) SetNameDescription3(session *Session, self PoolRef, value string) (err error) {
	method := "pool.set_name_description"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetNameLabel: Set the name_label field of the given pool.
// Version: rio
func (pool) SetNameLabel(session *Session, self PoolRef, value string) (err error) {
	method := "pool.set_name_label"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetNameLabel3: Set the name_label field of the given pool.
// Version: rio
func (pool) SetNameLabel3(session *Session, self PoolRef, value string) (err error) {
	method := "pool.set_name_label"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetAutoUpdateVMSecurebootCertificates: Get the auto_update_vm_secureboot_certificates field of the given pool.
// Version: rio
func (pool) GetAutoUpdateVMSecurebootCertificates(session *Session, self PoolRef) (retval bool, err error) {
	method := "pool.get_auto_update_vm_secureboot_certificates"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetAutoUpdateVMSecurebootCertificates2: Get the auto_update_vm_secureboot_certificates field of the given pool.
// Version: rio
func (pool) GetAutoUpdateVMSecurebootCertificates2(session *Session, self PoolRef) (retval bool, err error) {
	method := "pool.get_auto_update_vm_secureboot_certificates"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetVMConsoleIdleTimeout: Get the vm_console_idle_timeout field of the given pool.
// Version: rio
func (pool) GetVMConsoleIdleTimeout(session *Session, self PoolRef) (retval int, err error) {
	method := "pool.get_vm_console_idle_timeout"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetVMConsoleIdleTimeout2: Get the vm_console_idle_timeout field of the given pool.
// Version: rio
func (pool) GetVMConsoleIdleTimeout2(session *Session, self PoolRef) (retval int, err error) {
	method := "pool.get_vm_console_idle_timeout"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetLimitConsoleSessions: Get the limit_console_sessions field of the given pool.
// Version: rio
func (pool) GetLimitConsoleSessions(session *Session, self PoolRef) (retval bool, err error) {
	method := "pool.get_limit_console_sessions"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetLimitConsoleSessions2: Get the limit_console_sessions field of the given pool.
// Version: rio
func (pool) GetLimitConsoleSessions2(session *Session, self PoolRef) (retval bool, err error) {
	method := "pool.get_limit_console_sessions"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetHaRebootVMOnInternalShutdown: Get the ha_reboot_vm_on_internal_shutdown field of the given pool.
// Version: rio
func (pool) GetHaRebootVMOnInternalShutdown(session *Session, self PoolRef) (retval bool, err error) {
	method := "pool.get_ha_reboot_vm_on_internal_shutdown"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetHaRebootVMOnInternalShutdown2: Get the ha_reboot_vm_on_internal_shutdown field of the given pool.
// Version: rio
func (pool) GetHaRebootVMOnInternalShutdown2(session *Session, self PoolRef) (retval bool, err error) {
	method := "pool.get_ha_reboot_vm_on_internal_shutdown"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetLicenseServer: Get the license_server field of the given pool.
// Version: rio
func (pool) GetLicenseServer(session *Session, self PoolRef) (retval map[string]string, err error) {
	method := "pool.get_license_server"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetLicenseServer2: Get the license_server field of the given pool.
// Version: rio
func (pool) GetLicenseServer2(session *Session, self PoolRef) (retval map[string]string, err error) {
	method := "pool.get_license_server"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetRecommendations: Get the recommendations field of the given pool.
// Version: rio
func (pool) GetRecommendations(session *Session, self PoolRef) (retval map[string]string, err error) {
	method := "pool.get_recommendations"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetRecommendations2: Get the recommendations field of the given pool.
// Version: rio
func (pool) GetRecommendations2(session *Session, self PoolRef) (retval map[string]string, err error) {
	method := "pool.get_recommendations"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetUpdateSyncEnabled: Get the update_sync_enabled field of the given pool.
// Version: rio
func (pool) GetUpdateSyncEnabled(session *Session, self PoolRef) (retval bool, err error) {
	method := "pool.get_update_sync_enabled"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetUpdateSyncEnabled2: Get the update_sync_enabled field of the given pool.
// Version: rio
func (pool) GetUpdateSyncEnabled2(session *Session, self PoolRef) (retval bool, err error) {
	method := "pool.get_update_sync_enabled"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetUpdateSyncDay: Get the update_sync_day field of the given pool.
// Version: rio
func (pool) GetUpdateSyncDay(session *Session, self PoolRef) (retval int, err error) {
	method := "pool.get_update_sync_day"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetUpdateSyncDay2: Get the update_sync_day field of the given pool.
// Version: rio
func (pool) GetUpdateSyncDay2(session *Session, self PoolRef) (retval int, err error) {
	method := "pool.get_update_sync_day"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetUpdateSyncFrequency: Get the update_sync_frequency field of the given pool.
// Version: rio
func (pool) GetUpdateSyncFrequency(session *Session, self PoolRef) (retval UpdateSyncFrequency, err error) {
	method := "pool.get_update_sync_frequency"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeEnumUpdateSyncFrequency(method+" -> ", result)
	return
}

// GetUpdateSyncFrequency2: Get the update_sync_frequency field of the given pool.
// Version: rio
func (pool) GetUpdateSyncFrequency2(session *Session, self PoolRef) (retval UpdateSyncFrequency, err error) {
	method := "pool.get_update_sync_frequency"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeEnumUpdateSyncFrequency(method+" -> ", result)
	return
}

// GetLastUpdateSync: Get the last_update_sync field of the given pool.
// Version: rio
func (pool) GetLastUpdateSync(session *Session, self PoolRef) (retval time.Time, err error) {
	method := "pool.get_last_update_sync"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetLastUpdateSync2: Get the last_update_sync field of the given pool.
// Version: rio
func (pool) GetLastUpdateSync2(session *Session, self PoolRef) (retval time.Time, err error) {
	method := "pool.get_last_update_sync"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetTelemetryNextCollection: Get the telemetry_next_collection field of the given pool.
// Version: rio
func (pool) GetTelemetryNextCollection(session *Session, self PoolRef) (retval time.Time, err error) {
	method := "pool.get_telemetry_next_collection"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetTelemetryNextCollection2: Get the telemetry_next_collection field of the given pool.
// Version: rio
func (pool) GetTelemetryNextCollection2(session *Session, self PoolRef) (retval time.Time, err error) {
	method := "pool.get_telemetry_next_collection"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetTelemetryFrequency: Get the telemetry_frequency field of the given pool.
// Version: rio
func (pool) GetTelemetryFrequency(session *Session, self PoolRef) (retval TelemetryFrequency, err error) {
	method := "pool.get_telemetry_frequency"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeEnumTelemetryFrequency(method+" -> ", result)
	return
}

// GetTelemetryFrequency2: Get the telemetry_frequency field of the given pool.
// Version: rio
func (pool) GetTelemetryFrequency2(session *Session, self PoolRef) (retval TelemetryFrequency, err error) {
	method := "pool.get_telemetry_frequency"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeEnumTelemetryFrequency(method+" -> ", result)
	return
}

// GetTelemetryUUID: Get the telemetry_uuid field of the given pool.
// Version: rio
func (pool) GetTelemetryUUID(session *Session, self PoolRef) (retval SecretRef, err error) {
	method := "pool.get_telemetry_uuid"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeSecretRef(method+" -> ", result)
	return
}

// GetTelemetryUUID2: Get the telemetry_uuid field of the given pool.
// Version: rio
func (pool) GetTelemetryUUID2(session *Session, self PoolRef) (retval SecretRef, err error) {
	method := "pool.get_telemetry_uuid"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeSecretRef(method+" -> ", result)
	return
}

// GetExtAuthCacheExpiry: Get the ext_auth_cache_expiry field of the given pool.
// Version: rio
func (pool) GetExtAuthCacheExpiry(session *Session, self PoolRef) (retval int, err error) {
	method := "pool.get_ext_auth_cache_expiry"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetExtAuthCacheExpiry2: Get the ext_auth_cache_expiry field of the given pool.
// Version: rio
func (pool) GetExtAuthCacheExpiry2(session *Session, self PoolRef) (retval int, err error) {
	method := "pool.get_ext_auth_cache_expiry"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetExtAuthCacheSize: Get the ext_auth_cache_size field of the given pool.
// Version: rio
func (pool) GetExtAuthCacheSize(session *Session, self PoolRef) (retval int, err error) {
	method := "pool.get_ext_auth_cache_size"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetExtAuthCacheSize2: Get the ext_auth_cache_size field of the given pool.
// Version: rio
func (pool) GetExtAuthCacheSize2(session *Session, self PoolRef) (retval int, err error) {
	method := "pool.get_ext_auth_cache_size"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetExtAuthCacheEnabled: Get the ext_auth_cache_enabled field of the given pool.
// Version: rio
func (pool) GetExtAuthCacheEnabled(session *Session, self PoolRef) (retval bool, err error) {
	method := "pool.get_ext_auth_cache_enabled"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetExtAuthCacheEnabled2: Get the ext_auth_cache_enabled field of the given pool.
// Version: rio
func (pool) GetExtAuthCacheEnabled2(session *Session, self PoolRef) (retval bool, err error) {
	method := "pool.get_ext_auth_cache_enabled"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetExtAuthMaxThreads: Get the ext_auth_max_threads field of the given pool.
// Version: rio
func (pool) GetExtAuthMaxThreads(session *Session, self PoolRef) (retval int, err error) {
	method := "pool.get_ext_auth_max_threads"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetExtAuthMaxThreads2: Get the ext_auth_max_threads field of the given pool.
// Version: rio
func (pool) GetExtAuthMaxThreads2(session *Session, self PoolRef) (retval int, err error) {
	method := "pool.get_ext_auth_max_threads"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetLocalAuthMaxThreads: Get the local_auth_max_threads field of the given pool.
// Version: rio
func (pool) GetLocalAuthMaxThreads(session *Session, self PoolRef) (retval int, err error) {
	method := "pool.get_local_auth_max_threads"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetLocalAuthMaxThreads2: Get the local_auth_max_threads field of the given pool.
// Version: rio
func (pool) GetLocalAuthMaxThreads2(session *Session, self PoolRef) (retval int, err error) {
	method := "pool.get_local_auth_max_threads"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetCoordinatorBias: Get the coordinator_bias field of the given pool.
// Version: rio
func (pool) GetCoordinatorBias(session *Session, self PoolRef) (retval bool, err error) {
	method := "pool.get_coordinator_bias"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetCoordinatorBias2: Get the coordinator_bias field of the given pool.
// Version: rio
func (pool) GetCoordinatorBias2(session *Session, self PoolRef) (retval bool, err error) {
	method := "pool.get_coordinator_bias"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetMigrationCompression: Get the migration_compression field of the given pool.
// Version: rio
func (pool) GetMigrationCompression(session *Session, self PoolRef) (retval bool, err error) {
	method := "pool.get_migration_compression"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetMigrationCompression2: Get the migration_compression field of the given pool.
// Version: rio
func (pool) GetMigrationCompression2(session *Session, self PoolRef) (retval bool, err error) {
	method := "pool.get_migration_compression"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetRepositoryProxyPassword: Get the repository_proxy_password field of the given pool.
// Version: rio
func (pool) GetRepositoryProxyPassword(session *Session, self PoolRef) (retval SecretRef, err error) {
	method := "pool.get_repository_proxy_password"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeSecretRef(method+" -> ", result)
	return
}

// GetRepositoryProxyPassword2: Get the repository_proxy_password field of the given pool.
// Version: rio
func (pool) GetRepositoryProxyPassword2(session *Session, self PoolRef) (retval SecretRef, err error) {
	method := "pool.get_repository_proxy_password"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeSecretRef(method+" -> ", result)
	return
}

// GetRepositoryProxyUsername: Get the repository_proxy_username field of the given pool.
// Version: rio
func (pool) GetRepositoryProxyUsername(session *Session, self PoolRef) (retval string, err error) {
	method := "pool.get_repository_proxy_username"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetRepositoryProxyUsername2: Get the repository_proxy_username field of the given pool.
// Version: rio
func (pool) GetRepositoryProxyUsername2(session *Session, self PoolRef) (retval string, err error) {
	method := "pool.get_repository_proxy_username"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetRepositoryProxyURL: Get the repository_proxy_url field of the given pool.
// Version: rio
func (pool) GetRepositoryProxyURL(session *Session, self PoolRef) (retval string, err error) {
	method := "pool.get_repository_proxy_url"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetRepositoryProxyURL2: Get the repository_proxy_url field of the given pool.
// Version: rio
func (pool) GetRepositoryProxyURL2(session *Session, self PoolRef) (retval string, err error) {
	method := "pool.get_repository_proxy_url"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetClientCertificateAuthName: Get the client_certificate_auth_name field of the given pool.
// Version: rio
func (pool) GetClientCertificateAuthName(session *Session, self PoolRef) (retval string, err error) {
	method := "pool.get_client_certificate_auth_name"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetClientCertificateAuthName2: Get the client_certificate_auth_name field of the given pool.
// Version: rio
func (pool) GetClientCertificateAuthName2(session *Session, self PoolRef) (retval string, err error) {
	method := "pool.get_client_certificate_auth_name"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetClientCertificateAuthEnabled: Get the client_certificate_auth_enabled field of the given pool.
// Version: rio
func (pool) GetClientCertificateAuthEnabled(session *Session, self PoolRef) (retval bool, err error) {
	method := "pool.get_client_certificate_auth_enabled"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetClientCertificateAuthEnabled2: Get the client_certificate_auth_enabled field of the given pool.
// Version: rio
func (pool) GetClientCertificateAuthEnabled2(session *Session, self PoolRef) (retval bool, err error) {
	method := "pool.get_client_certificate_auth_enabled"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetRepositories: Get the repositories field of the given pool.
// Version: rio
func (pool) GetRepositories(session *Session, self PoolRef) (retval []RepositoryRef, err error) {
	method := "pool.get_repositories"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeRepositoryRefSet(method+" -> ", result)
	return
}

// GetRepositories2: Get the repositories field of the given pool.
// Version: rio
func (pool) GetRepositories2(session *Session, self PoolRef) (retval []RepositoryRef, err error) {
	method := "pool.get_repositories"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeRepositoryRefSet(method+" -> ", result)
	return
}

// GetTLSVerificationEnabled: Get the tls_verification_enabled field of the given pool.
// Version: rio
func (pool) GetTLSVerificationEnabled(session *Session, self PoolRef) (retval bool, err error) {
	method := "pool.get_tls_verification_enabled"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetTLSVerificationEnabled2: Get the tls_verification_enabled field of the given pool.
// Version: rio
func (pool) GetTLSVerificationEnabled2(session *Session, self PoolRef) (retval bool, err error) {
	method := "pool.get_tls_verification_enabled"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetIsPsrPending: Get the is_psr_pending field of the given pool.
// Version: rio
func (pool) GetIsPsrPending(session *Session, self PoolRef) (retval bool, err error) {
	method := "pool.get_is_psr_pending"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetIsPsrPending2: Get the is_psr_pending field of the given pool.
// Version: rio
func (pool) GetIsPsrPending2(session *Session, self PoolRef) (retval bool, err error) {
	method := "pool.get_is_psr_pending"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetCustomUefiCertificates: Get the custom_uefi_certificates field of the given pool.
// Version: rio
func (pool) GetCustomUefiCertificates(session *Session, self PoolRef) (retval string, err error) {
	method := "pool.get_custom_uefi_certificates"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetCustomUefiCertificates2: Get the custom_uefi_certificates field of the given pool.
// Version: rio
func (pool) GetCustomUefiCertificates2(session *Session, self PoolRef) (retval string, err error) {
	method := "pool.get_custom_uefi_certificates"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetUefiCertificates: Get the uefi_certificates field of the given pool.
// Version: rio
func (pool) GetUefiCertificates(session *Session, self PoolRef) (retval string, err error) {
	method := "pool.get_uefi_certificates"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetUefiCertificates2: Get the uefi_certificates field of the given pool.
// Version: rio
func (pool) GetUefiCertificates2(session *Session, self PoolRef) (retval string, err error) {
	method := "pool.get_uefi_certificates"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetIgmpSnoopingEnabled: Get the igmp_snooping_enabled field of the given pool.
// Version: rio
func (pool) GetIgmpSnoopingEnabled(session *Session, self PoolRef) (retval bool, err error) {
	method := "pool.get_igmp_snooping_enabled"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetIgmpSnoopingEnabled2: Get the igmp_snooping_enabled field of the given pool.
// Version: rio
func (pool) GetIgmpSnoopingEnabled2(session *Session, self PoolRef) (retval bool, err error) {
	method := "pool.get_igmp_snooping_enabled"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetLivePatchingDisabled: Get the live_patching_disabled field of the given pool.
// Version: rio
func (pool) GetLivePatchingDisabled(session *Session, self PoolRef) (retval bool, err error) {
	method := "pool.get_live_patching_disabled"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetLivePatchingDisabled2: Get the live_patching_disabled field of the given pool.
// Version: rio
func (pool) GetLivePatchingDisabled2(session *Session, self PoolRef) (retval bool, err error) {
	method := "pool.get_live_patching_disabled"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetPolicyNoVendorDevice: Get the policy_no_vendor_device field of the given pool.
// Version: rio
func (pool) GetPolicyNoVendorDevice(session *Session, self PoolRef) (retval bool, err error) {
	method := "pool.get_policy_no_vendor_device"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetPolicyNoVendorDevice2: Get the policy_no_vendor_device field of the given pool.
// Version: rio
func (pool) GetPolicyNoVendorDevice2(session *Session, self PoolRef) (retval bool, err error) {
	method := "pool.get_policy_no_vendor_device"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetCPUInfo: Get the cpu_info field of the given pool.
// Version: rio
func (pool) GetCPUInfo(session *Session, self PoolRef) (retval map[string]string, err error) {
	method := "pool.get_cpu_info"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetCPUInfo2: Get the cpu_info field of the given pool.
// Version: rio
func (pool) GetCPUInfo2(session *Session, self PoolRef) (retval map[string]string, err error) {
	method := "pool.get_cpu_info"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetGuestAgentConfig: Get the guest_agent_config field of the given pool.
// Version: rio
func (pool) GetGuestAgentConfig(session *Session, self PoolRef) (retval map[string]string, err error) {
	method := "pool.get_guest_agent_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetGuestAgentConfig2: Get the guest_agent_config field of the given pool.
// Version: rio
func (pool) GetGuestAgentConfig2(session *Session, self PoolRef) (retval map[string]string, err error) {
	method := "pool.get_guest_agent_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetCurrentOperations: Get the current_operations field of the given pool.
// Version: rio
func (pool) GetCurrentOperations(session *Session, self PoolRef) (retval map[string]PoolAllowedOperations, err error) {
	method := "pool.get_current_operations"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeStringToEnumPoolAllowedOperationsMap(method+" -> ", result)
	return
}

// GetCurrentOperations2: Get the current_operations field of the given pool.
// Version: rio
func (pool) GetCurrentOperations2(session *Session, self PoolRef) (retval map[string]PoolAllowedOperations, err error) {
	method := "pool.get_current_operations"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeStringToEnumPoolAllowedOperationsMap(method+" -> ", result)
	return
}

// GetAllowedOperations: Get the allowed_operations field of the given pool.
// Version: rio
func (pool) GetAllowedOperations(session *Session, self PoolRef) (retval []PoolAllowedOperations, err error) {
	method := "pool.get_allowed_operations"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeEnumPoolAllowedOperationsSet(method+" -> ", result)
	return
}

// GetAllowedOperations2: Get the allowed_operations field of the given pool.
// Version: rio
func (pool) GetAllowedOperations2(session *Session, self PoolRef) (retval []PoolAllowedOperations, err error) {
	method := "pool.get_allowed_operations"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeEnumPoolAllowedOperationsSet(method+" -> ", result)
	return
}

// GetHaClusterStack: Get the ha_cluster_stack field of the given pool.
// Version: rio
func (pool) GetHaClusterStack(session *Session, self PoolRef) (retval string, err error) {
	method := "pool.get_ha_cluster_stack"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetHaClusterStack2: Get the ha_cluster_stack field of the given pool.
// Version: rio
func (pool) GetHaClusterStack2(session *Session, self PoolRef) (retval string, err error) {
	method := "pool.get_ha_cluster_stack"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetMetadataVDIs: Get the metadata_VDIs field of the given pool.
// Version: rio
func (pool) GetMetadataVDIs(session *Session, self PoolRef) (retval []VDIRef, err error) {
	method := "pool.get_metadata_VDIs"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeVDIRefSet(method+" -> ", result)
	return
}

// GetMetadataVDIs2: Get the metadata_VDIs field of the given pool.
// Version: rio
func (pool) GetMetadataVDIs2(session *Session, self PoolRef) (retval []VDIRef, err error) {
	method := "pool.get_metadata_VDIs"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeVDIRefSet(method+" -> ", result)
	return
}

// GetRestrictions: Get the restrictions field of the given pool.
// Version: rio
func (pool) GetRestrictions(session *Session, self PoolRef) (retval map[string]string, err error) {
	method := "pool.get_restrictions"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetRestrictions2: Get the restrictions field of the given pool.
// Version: rio
func (pool) GetRestrictions2(session *Session, self PoolRef) (retval map[string]string, err error) {
	method := "pool.get_restrictions"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetVswitchController: Get the vswitch_controller field of the given pool.
// Version: rio
func (pool) GetVswitchController(session *Session, self PoolRef) (retval string, err error) {
	method := "pool.get_vswitch_controller"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetVswitchController2: Get the vswitch_controller field of the given pool.
// Version: rio
func (pool) GetVswitchController2(session *Session, self PoolRef) (retval string, err error) {
	method := "pool.get_vswitch_controller"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetRedoLogVdi: Get the redo_log_vdi field of the given pool.
// Version: rio
func (pool) GetRedoLogVdi(session *Session, self PoolRef) (retval VDIRef, err error) {
	method := "pool.get_redo_log_vdi"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeVDIRef(method+" -> ", result)
	return
}

// GetRedoLogVdi2: Get the redo_log_vdi field of the given pool.
// Version: rio
func (pool) GetRedoLogVdi2(session *Session, self PoolRef) (retval VDIRef, err error) {
	method := "pool.get_redo_log_vdi"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeVDIRef(method+" -> ", result)
	return
}

// GetRedoLogEnabled: Get the redo_log_enabled field of the given pool.
// Version: rio
func (pool) GetRedoLogEnabled(session *Session, self PoolRef) (retval bool, err error) {
	method := "pool.get_redo_log_enabled"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetRedoLogEnabled2: Get the redo_log_enabled field of the given pool.
// Version: rio
func (pool) GetRedoLogEnabled2(session *Session, self PoolRef) (retval bool, err error) {
	method := "pool.get_redo_log_enabled"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetWlbVerifyCert: Get the wlb_verify_cert field of the given pool.
// Version: rio
func (pool) GetWlbVerifyCert(session *Session, self PoolRef) (retval bool, err error) {
	method := "pool.get_wlb_verify_cert"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetWlbVerifyCert2: Get the wlb_verify_cert field of the given pool.
// Version: rio
func (pool) GetWlbVerifyCert2(session *Session, self PoolRef) (retval bool, err error) {
	method := "pool.get_wlb_verify_cert"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetWlbEnabled: Get the wlb_enabled field of the given pool.
// Version: rio
func (pool) GetWlbEnabled(session *Session, self PoolRef) (retval bool, err error) {
	method := "pool.get_wlb_enabled"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetWlbEnabled2: Get the wlb_enabled field of the given pool.
// Version: rio
func (pool) GetWlbEnabled2(session *Session, self PoolRef) (retval bool, err error) {
	method := "pool.get_wlb_enabled"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetWlbUsername: Get the wlb_username field of the given pool.
// Version: rio
func (pool) GetWlbUsername(session *Session, self PoolRef) (retval string, err error) {
	method := "pool.get_wlb_username"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetWlbUsername2: Get the wlb_username field of the given pool.
// Version: rio
func (pool) GetWlbUsername2(session *Session, self PoolRef) (retval string, err error) {
	method := "pool.get_wlb_username"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetWlbURL: Get the wlb_url field of the given pool.
// Version: rio
func (pool) GetWlbURL(session *Session, self PoolRef) (retval string, err error) {
	method := "pool.get_wlb_url"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetWlbURL2: Get the wlb_url field of the given pool.
// Version: rio
func (pool) GetWlbURL2(session *Session, self PoolRef) (retval string, err error) {
	method := "pool.get_wlb_url"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetHealthCheckConfig: Get the health_check_config field of the given pool.
// Version: rio
func (pool) GetHealthCheckConfig(session *Session, self PoolRef) (retval map[string]string, err error) {
	method := "pool.get_health_check_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetHealthCheckConfig2: Get the health_check_config field of the given pool.
// Version: rio
func (pool) GetHealthCheckConfig2(session *Session, self PoolRef) (retval map[string]string, err error) {
	method := "pool.get_health_check_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetGuiConfig: Get the gui_config field of the given pool.
// Version: rio
func (pool) GetGuiConfig(session *Session, self PoolRef) (retval map[string]string, err error) {
	method := "pool.get_gui_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetGuiConfig2: Get the gui_config field of the given pool.
// Version: rio
func (pool) GetGuiConfig2(session *Session, self PoolRef) (retval map[string]string, err error) {
	method := "pool.get_gui_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetTags: Get the tags field of the given pool.
// Version: rio
func (pool) GetTags(session *Session, self PoolRef) (retval []string, err error) {
	method := "pool.get_tags"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetTags2: Get the tags field of the given pool.
// Version: rio
func (pool) GetTags2(session *Session, self PoolRef) (retval []string, err error) {
	method := "pool.get_tags"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetBlobs: Get the blobs field of the given pool.
// Version: rio
func (pool) GetBlobs(session *Session, self PoolRef) (retval map[string]BlobRef, err error) {
	method := "pool.get_blobs"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetBlobs2: Get the blobs field of the given pool.
// Version: rio
func (pool) GetBlobs2(session *Session, self PoolRef) (retval map[string]BlobRef, err error) {
	method := "pool.get_blobs"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetHaOvercommitted: Get the ha_overcommitted field of the given pool.
// Version: rio
func (pool) GetHaOvercommitted(session *Session, self PoolRef) (retval bool, err error) {
	method := "pool.get_ha_overcommitted"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetHaOvercommitted2: Get the ha_overcommitted field of the given pool.
// Version: rio
func (pool) GetHaOvercommitted2(session *Session, self PoolRef) (retval bool, err error) {
	method := "pool.get_ha_overcommitted"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetHaAllowOvercommit: Get the ha_allow_overcommit field of the given pool.
// Version: rio
func (pool) GetHaAllowOvercommit(session *Session, self PoolRef) (retval bool, err error) {
	method := "pool.get_ha_allow_overcommit"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetHaAllowOvercommit2: Get the ha_allow_overcommit field of the given pool.
// Version: rio
func (pool) GetHaAllowOvercommit2(session *Session, self PoolRef) (retval bool, err error) {
	method := "pool.get_ha_allow_overcommit"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetHaPlanExistsFor: Get the ha_plan_exists_for field of the given pool.
// Version: rio
func (pool) GetHaPlanExistsFor(session *Session, self PoolRef) (retval int, err error) {
	method := "pool.get_ha_plan_exists_for"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetHaPlanExistsFor2: Get the ha_plan_exists_for field of the given pool.
// Version: rio
func (pool) GetHaPlanExistsFor2(session *Session, self PoolRef) (retval int, err error) {
	method := "pool.get_ha_plan_exists_for"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetHaHostFailuresToTolerate: Get the ha_host_failures_to_tolerate field of the given pool.
// Version: rio
func (pool) GetHaHostFailuresToTolerate(session *Session, self PoolRef) (retval int, err error) {
	method := "pool.get_ha_host_failures_to_tolerate"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetHaHostFailuresToTolerate2: Get the ha_host_failures_to_tolerate field of the given pool.
// Version: rio
func (pool) GetHaHostFailuresToTolerate2(session *Session, self PoolRef) (retval int, err error) {
	method := "pool.get_ha_host_failures_to_tolerate"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetHaStatefiles: Get the ha_statefiles field of the given pool.
// Version: rio
func (pool) GetHaStatefiles(session *Session, self PoolRef) (retval []string, err error) {
	method := "pool.get_ha_statefiles"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetHaStatefiles2: Get the ha_statefiles field of the given pool.
// Version: rio
func (pool) GetHaStatefiles2(session *Session, self PoolRef) (retval []string, err error) {
	method := "pool.get_ha_statefiles"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetHaConfiguration: Get the ha_configuration field of the given pool.
// Version: rio
func (pool) GetHaConfiguration(session *Session, self PoolRef) (retval map[string]string, err error) {
	method := "pool.get_ha_configuration"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetHaConfiguration2: Get the ha_configuration field of the given pool.
// Version: rio
func (pool) GetHaConfiguration2(session *Session, self PoolRef) (retval map[string]string, err error) {
	method := "pool.get_ha_configuration"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetHaEnabled: Get the ha_enabled field of the given pool.
// Version: rio
func (pool) GetHaEnabled(session *Session, self PoolRef) (retval bool, err error) {
	method := "pool.get_ha_enabled"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetHaEnabled2: Get the ha_enabled field of the given pool.
// Version: rio
func (pool) GetHaEnabled2(session *Session, self PoolRef) (retval bool, err error) {
	method := "pool.get_ha_enabled"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetOtherConfig: Get the other_config field of the given pool.
// Version: rio
func (pool) GetOtherConfig(session *Session, self PoolRef) (retval map[string]string, err error) {
	method := "pool.get_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetOtherConfig2: Get the other_config field of the given pool.
// Version: rio
func (pool) GetOtherConfig2(session *Session, self PoolRef) (retval map[string]string, err error) {
	method := "pool.get_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetCrashDumpSR: Get the crash_dump_SR field of the given pool.
// Version: rio
func (pool) GetCrashDumpSR(session *Session, self PoolRef) (retval SRRef, err error) {
	method := "pool.get_crash_dump_SR"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetCrashDumpSR2: Get the crash_dump_SR field of the given pool.
// Version: rio
func (pool) GetCrashDumpSR2(session *Session, self PoolRef) (retval SRRef, err error) {
	method := "pool.get_crash_dump_SR"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetSuspendImageSR: Get the suspend_image_SR field of the given pool.
// Version: rio
func (pool) GetSuspendImageSR(session *Session, self PoolRef) (retval SRRef, err error) {
	method := "pool.get_suspend_image_SR"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetSuspendImageSR2: Get the suspend_image_SR field of the given pool.
// Version: rio
func (pool) GetSuspendImageSR2(session *Session, self PoolRef) (retval SRRef, err error) {
	method := "pool.get_suspend_image_SR"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetDefaultSR: Get the default_SR field of the given pool.
// Version: rio
func (pool) GetDefaultSR(session *Session, self PoolRef) (retval SRRef, err error) {
	method := "pool.get_default_SR"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetDefaultSR2: Get the default_SR field of the given pool.
// Version: rio
func (pool) GetDefaultSR2(session *Session, self PoolRef) (retval SRRef, err error) {
	method := "pool.get_default_SR"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetMaster: Get the master field of the given pool.
// Version: rio
func (pool) GetMaster(session *Session, self PoolRef) (retval HostRef, err error) {
	method := "pool.get_master"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeHostRef(method+" -> ", result)
	return
}

// GetMaster2: Get the master field of the given pool.
// Version: rio
func (pool) GetMaster2(session *Session, self PoolRef) (retval HostRef, err error) {
	method := "pool.get_master"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeHostRef(method+" -> ", result)
	return
}

// GetNameDescription: Get the name_description field of the given pool.
// Version: rio
func (pool) GetNameDescription(session *Session, self PoolRef) (retval string, err error) {
	method := "pool.get_name_description"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetNameDescription2: Get the name_description field of the given pool.
// Version: rio
func (pool) GetNameDescription2(session *Session, self PoolRef) (retval string, err error) {
	method := "pool.get_name_description"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetNameLabel: Get the name_label field of the given pool.
// Version: rio
func (pool) GetNameLabel(session *Session, self PoolRef) (retval string, err error) {
	method := "pool.get_name_label"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetNameLabel2: Get the name_label field of the given pool.
// Version: rio
func (pool) GetNameLabel2(session *Session, self PoolRef) (retval string, err error) {
	method := "pool.get_name_label"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetUUID: Get the uuid field of the given pool.
// Version: rio
func (pool) GetUUID(session *Session, self PoolRef) (retval string, err error) {
	method := "pool.get_uuid"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetUUID2: Get the uuid field of the given pool.
// Version: rio
func (pool) GetUUID2(session *Session, self PoolRef) (retval string, err error) {
	method := "pool.get_uuid"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetByUUID: Get a reference to the pool instance with the specified UUID.
// Version: rio
func (pool) GetByUUID(session *Session, uuid string) (retval PoolRef, err error) {
	method := "pool.get_by_uuid"
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
	retval, err = deserializePoolRef(method+" -> ", result)
	return
}

// GetByUUID2: Get a reference to the pool instance with the specified UUID.
// Version: rio
func (pool) GetByUUID2(session *Session, uuid string) (retval PoolRef, err error) {
	method := "pool.get_by_uuid"
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
	retval, err = deserializePoolRef(method+" -> ", result)
	return
}

// GetRecord: Get a record containing the current state of the given pool.
// Version: rio
func (pool) GetRecord(session *Session, self PoolRef) (retval PoolRecord, err error) {
	method := "pool.get_record"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializePoolRecord(method+" -> ", result)
	return
}

// GetRecord2: Get a record containing the current state of the given pool.
// Version: rio
func (pool) GetRecord2(session *Session, self PoolRef) (retval PoolRecord, err error) {
	method := "pool.get_record"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializePoolRecord(method+" -> ", result)
	return
}
