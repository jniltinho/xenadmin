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

//nolint:gosec
const (
	//
	ErrorBootDevicePassthroughDisallowed = "BOOT_DEVICE_PASSTHROUGH_DISALLOWED"
	//
	ErrorCertificateLacksPurpose = "CERTIFICATE_LACKS_PURPOSE"
	//
	ErrorNotTrustedCertificate = "NOT_TRUSTED_CERTIFICATE"
	//
	ErrorNotAllowedWhenNtpIsEnabled = "NOT_ALLOWED_WHEN_NTP_IS_ENABLED"
	//
	ErrorInvalidNtpConfig = "INVALID_NTP_CONFIG"
	//
	ErrorSysprep = "SYSPREP"
	//
	ErrorTLSVerificationNotEnabledInPool = "TLS_VERIFICATION_NOT_ENABLED_IN_POOL"
	//
	ErrorHostDriverNoHardware = "HOST_DRIVER_NO_HARDWARE"
	//
	ErrorSetSSHAutoModePartiallyFailed = "SET_SSH_AUTO_MODE_PARTIALLY_FAILED"
	//
	ErrorSetConsoleTimeoutPartiallyFailed = "SET_CONSOLE_TIMEOUT_PARTIALLY_FAILED"
	//
	ErrorSetSSHTimeoutPartiallyFailed = "SET_SSH_TIMEOUT_PARTIALLY_FAILED"
	//
	ErrorDisableSSHPartiallyFailed = "DISABLE_SSH_PARTIALLY_FAILED"
	//
	ErrorEnableSSHPartiallyFailed = "ENABLE_SSH_PARTIALLY_FAILED"
	//
	ErrorDisableSSHFailed = "DISABLE_SSH_FAILED"
	//
	ErrorEnableSSHFailed = "ENABLE_SSH_FAILED"
	//
	ErrorTooManyGroups = "TOO_MANY_GROUPS"
	//
	ErrorIllegalInFipsMode = "ILLEGAL_IN_FIPS_MODE"
	//
	ErrorTelemetryNextCollectionTooLate = "TELEMETRY_NEXT_COLLECTION_TOO_LATE"
	//
	ErrorVtpmMaxAmountReached = "VTPM_MAX_AMOUNT_REACHED"
	//
	ErrorHostEvacuationIsRequired = "HOST_EVACUATION_IS_REQUIRED"
	//
	ErrorHostPendingMandatoryGuidanceNotEmpty = "HOST_PENDING_MANDATORY_GUIDANCE_NOT_EMPTY"
	//
	ErrorNoRepositoriesConfigured = "NO_REPOSITORIES_CONFIGURED"
	//
	ErrorInvalidUpdateSyncDay = "INVALID_UPDATE_SYNC_DAY"
	//
	ErrorApplyLivepatchFailed = "APPLY_LIVEPATCH_FAILED"
	//
	ErrorInvalidRepositoryDomainAllowlist = "INVALID_REPOSITORY_DOMAIN_ALLOWLIST"
	//
	ErrorInvalidRepositoryProxyCredential = "INVALID_REPOSITORY_PROXY_CREDENTIAL"
	//
	ErrorInvalidRepositoryProxyURL = "INVALID_REPOSITORY_PROXY_URL"
	//
	ErrorCannotRestartDeviceModel = "CANNOT_RESTART_DEVICE_MODEL"
	//
	ErrorUpdateinfoHashMismatch = "UPDATEINFO_HASH_MISMATCH"
	//
	ErrorApplyGuidanceFailed = "APPLY_GUIDANCE_FAILED"
	//
	ErrorApplyUpdatesFailed = "APPLY_UPDATES_FAILED"
	//
	ErrorApplyUpdatesInProgress = "APPLY_UPDATES_IN_PROGRESS"
	//
	ErrorGetUpdatesInProgress = "GET_UPDATES_IN_PROGRESS"
	//
	ErrorGetUpdatesFailed = "GET_UPDATES_FAILED"
	//
	ErrorInvalidRepomdXML = "INVALID_REPOMD_XML"
	//
	ErrorGetHostUpdatesFailed = "GET_HOST_UPDATES_FAILED"
	//
	ErrorInvalidUpdateinfoXML = "INVALID_UPDATEINFO_XML"
	//
	ErrorCreaterepoFailed = "CREATEREPO_FAILED"
	//
	ErrorBundleSyncFailed = "BUNDLE_SYNC_FAILED"
	//
	ErrorReposyncFailed = "REPOSYNC_FAILED"
	//
	ErrorSyncBundleInProgress = "SYNC_BUNDLE_IN_PROGRESS"
	//
	ErrorSyncUpdatesInProgress = "SYNC_UPDATES_IN_PROGRESS"
	//
	ErrorMultipleUpdateRepositoriesEnabled = "MULTIPLE_UPDATE_REPOSITORIES_ENABLED"
	//
	ErrorNoRepositoryEnabled = "NO_REPOSITORY_ENABLED"
	//
	ErrorRepositoryCleanupFailed = "REPOSITORY_CLEANUP_FAILED"
	//
	ErrorRepositoryIsInUse = "REPOSITORY_IS_IN_USE"
	//
	ErrorUpdateSyncingRemotePoolCoordinatorServiceFailed = "UPDATE_SYNCING_REMOTE_POOL_COORDINATOR_SERVICE_FAILED"
	//
	ErrorUpdateSyncingRemotePoolCoordinatorConnectionFailed = "UPDATE_SYNCING_REMOTE_POOL_COORDINATOR_CONNECTION_FAILED"
	//
	ErrorRepoShouldBeSingleOneEnabled = "REPO_SHOULD_BE_SINGLE_ONE_ENABLED"
	//
	ErrorCanNotPeriodicSyncUpdates = "CAN_NOT_PERIODIC_SYNC_UPDATES"
	//
	ErrorCanNotSyncUpdates = "CAN_NOT_SYNC_UPDATES"
	//
	ErrorBundleRepoNotEnabled = "BUNDLE_REPO_NOT_ENABLED"
	//
	ErrorBundleUnpackFailed = "BUNDLE_UNPACK_FAILED"
	//
	ErrorBundleRepositoryAlreadyExists = "BUNDLE_REPOSITORY_ALREADY_EXISTS"
	//
	ErrorRepositoryAlreadyExists = "REPOSITORY_ALREADY_EXISTS"
	//
	ErrorInvalidGpgkeyPath = "INVALID_GPGKEY_PATH"
	//
	ErrorBlockedRepoURL = "BLOCKED_REPO_URL"
	//
	ErrorInvalidBaseURL = "INVALID_BASE_URL"
	//
	ErrorConfigureRepositoriesInProgress = "CONFIGURE_REPOSITORIES_IN_PROGRESS"
	//
	ErrorCertRefreshInProgress = "CERT_REFRESH_IN_PROGRESS"
	//
	ErrorTLSVerificationEnableInProgress = "TLS_VERIFICATION_ENABLE_IN_PROGRESS"
	//
	ErrorPoolSecretRotationPending = "POOL_SECRET_ROTATION_PENDING"
	//
	ErrorDesignateNewMasterInProgress = "DESIGNATE_NEW_MASTER_IN_PROGRESS"
	//
	ErrorVcpuMaxNotCoresPerSocketMultiple = "VCPU_MAX_NOT_CORES_PER_SOCKET_MULTIPLE"
	//
	ErrorXenIncompatible = "XEN_INCOMPATIBLE"
	//
	ErrorNoClusterHostsReachable = "NO_CLUSTER_HOSTS_REACHABLE"
	//
	ErrorClusterHostNotJoined = "CLUSTER_HOST_NOT_JOINED"
	//
	ErrorPifNotAttachedToHost = "PIF_NOT_ATTACHED_TO_HOST"
	//
	ErrorInvalidClusterStack = "INVALID_CLUSTER_STACK"
	//
	ErrorClusterStackInUse = "CLUSTER_STACK_IN_USE"
	//
	ErrorClusterForceDestroyFailed = "CLUSTER_FORCE_DESTROY_FAILED"
	//
	ErrorNoCompatibleClusterHost = "NO_COMPATIBLE_CLUSTER_HOST"
	//
	ErrorClusterHostIsLast = "CLUSTER_HOST_IS_LAST"
	//
	ErrorClusterDoesNotHaveOneNode = "CLUSTER_DOES_NOT_HAVE_ONE_NODE"
	//
	ErrorClusteringDisabled = "CLUSTERING_DISABLED"
	//
	ErrorClusteringEnabled = "CLUSTERING_ENABLED"
	//
	ErrorClusterAlreadyExists = "CLUSTER_ALREADY_EXISTS"
	//
	ErrorClusterCreateInProgress = "CLUSTER_CREATE_IN_PROGRESS"
	//
	ErrorClusterHasNoCertificate = "CLUSTER_HAS_NO_CERTIFICATE"
	//
	ErrorVMHasVusbs = "VM_HAS_VUSBS"
	//
	ErrorPusbVdiConflict = "PUSB_VDI_CONFLICT"
	//
	ErrorTooManyVusbs = "TOO_MANY_VUSBS"
	//
	ErrorUsbAlreadyAttached = "USB_ALREADY_ATTACHED"
	//
	ErrorUsbGroupConflict = "USB_GROUP_CONFLICT"
	//
	ErrorUsbGroupContainsNoPusbs = "USB_GROUP_CONTAINS_NO_PUSBS"
	//
	ErrorUsbGroupContainsPusb = "USB_GROUP_CONTAINS_PUSB"
	//
	ErrorUsbGroupContainsVusb = "USB_GROUP_CONTAINS_VUSB"
	//
	ErrorExtensionProtocolFailure = "EXTENSION_PROTOCOL_FAILURE"
	//
	ErrorPvsProxyPresentOnHigherVifDevice = "PVS_PROXY_PRESENT_ON_HIGHER_VIF_DEVICE"
	//
	ErrorPvsVifMustBeFirstDevice = "PVS_VIF_MUST_BE_FIRST_DEVICE"
	//
	ErrorPvsServerAddressInUse = "PVS_SERVER_ADDRESS_IN_USE"
	//
	ErrorPvsProxyAlreadyPresent = "PVS_PROXY_ALREADY_PRESENT"
	//
	ErrorPvsCacheStorageIsInUse = "PVS_CACHE_STORAGE_IS_IN_USE"
	//
	ErrorPvsCacheStorageAlreadyPresent = "PVS_CACHE_STORAGE_ALREADY_PRESENT"
	//
	ErrorPvsSiteContainsServers = "PVS_SITE_CONTAINS_SERVERS"
	//
	ErrorPvsSiteContainsRunningProxies = "PVS_SITE_CONTAINS_RUNNING_PROXIES"
	//
	ErrorSuspendImageNotAccessible = "SUSPEND_IMAGE_NOT_ACCESSIBLE"
	//
	ErrorVMCallPluginRateLimit = "VM_CALL_PLUGIN_RATE_LIMIT"
	//
	ErrorUnimplementedInSmBackend = "UNIMPLEMENTED_IN_SM_BACKEND"
	//
	ErrorSrDoesNotSupportMigration = "SR_DOES_NOT_SUPPORT_MIGRATION"
	//
	ErrorTooManyStorageMigrates = "TOO_MANY_STORAGE_MIGRATES"
	//
	ErrorMirrorFailed = "MIRROR_FAILED"
	//
	ErrorVMHasCheckpoint = "VM_HAS_CHECKPOINT"
	//
	ErrorVMHasTooManySnapshots = "VM_HAS_TOO_MANY_SNAPSHOTS"
	//
	ErrorVdiNeedsVMForMigrate = "VDI_NEEDS_VM_FOR_MIGRATE"
	//
	ErrorVdiCopyFailed = "VDI_COPY_FAILED"
	//
	ErrorSuspendVdiReplacementIsNotIdentical = "SUSPEND_VDI_REPLACEMENT_IS_NOT_IDENTICAL"
	//
	ErrorVMToImportIsNotNewerVersion = "VM_TO_IMPORT_IS_NOT_NEWER_VERSION"
	//
	ErrorVMIsPartOfAnAppliance = "VM_IS_PART_OF_AN_APPLIANCE"
	//
	ErrorCannotDestroyDisasterRecoveryTask = "CANNOT_DESTROY_DISASTER_RECOVERY_TASK"
	//
	ErrorVMIncompatibleWithThisHost = "VM_INCOMPATIBLE_WITH_THIS_HOST"
	//
	ErrorCouldNotImportDatabase = "COULD_NOT_IMPORT_DATABASE"
	//
	ErrorNoMoreRedoLogsAllowed = "NO_MORE_REDO_LOGS_ALLOWED"
	//
	ErrorVdiContainsMetadataOfThisPool = "VDI_CONTAINS_METADATA_OF_THIS_POOL"
	//
	ErrorFeatureRequiresHvm = "FEATURE_REQUIRES_HVM"
	//
	ErrorCPUFeatureMaskingNotSupported = "CPU_FEATURE_MASKING_NOT_SUPPORTED"
	//
	ErrorInvalidFeatureString = "INVALID_FEATURE_STRING"
	//
	ErrorVMBiosStringsAlreadySet = "VM_BIOS_STRINGS_ALREADY_SET"
	//
	ErrorRedoLogIsEnabled = "REDO_LOG_IS_ENABLED"
	//
	ErrorCannotEnableRedoLog = "CANNOT_ENABLE_REDO_LOG"
	//
	ErrorSslVerifyError = "SSL_VERIFY_ERROR"
	//
	ErrorVMAssignedToSnapshotSchedule = "VM_ASSIGNED_TO_SNAPSHOT_SCHEDULE"
	//
	ErrorVmssHasVM = "VMSS_HAS_VM"
	//
	ErrorVMAssignedToProtectionPolicy = "VM_ASSIGNED_TO_PROTECTION_POLICY"
	//
	ErrorVmppArchiveMoreFrequentThanBackup = "VMPP_ARCHIVE_MORE_FREQUENT_THAN_BACKUP"
	//
	ErrorVmppHasVM = "VMPP_HAS_VM"
	//
	ErrorServerCertificateChainInvalid = "SERVER_CERTIFICATE_CHAIN_INVALID"
	//
	ErrorServerCertificateSignatureNotSupported = "SERVER_CERTIFICATE_SIGNATURE_NOT_SUPPORTED"
	//
	ErrorTrustedCertificateInvalid = "TRUSTED_CERTIFICATE_INVALID"
	//
	ErrorTrustedCertificateNotValidYet = "TRUSTED_CERTIFICATE_NOT_VALID_YET"
	//
	ErrorTrustedCertificateExpired = "TRUSTED_CERTIFICATE_EXPIRED"
	//
	ErrorCaCertificateExpired = "CA_CERTIFICATE_EXPIRED"
	//
	ErrorServerCertificateExpired = "SERVER_CERTIFICATE_EXPIRED"
	//
	ErrorCaCertificateNotValidYet = "CA_CERTIFICATE_NOT_VALID_YET"
	//
	ErrorServerCertificateNotValidYet = "SERVER_CERTIFICATE_NOT_VALID_YET"
	//
	ErrorServerCertificateKeyMismatch = "SERVER_CERTIFICATE_KEY_MISMATCH"
	//
	ErrorCaCertificateInvalid = "CA_CERTIFICATE_INVALID"
	//
	ErrorServerCertificateInvalid = "SERVER_CERTIFICATE_INVALID"
	//
	ErrorServerCertificateKeyRsaMultiNotSupported = "SERVER_CERTIFICATE_KEY_RSA_MULTI_NOT_SUPPORTED"
	//
	ErrorServerCertificateKeyRsaLengthNotSupported = "SERVER_CERTIFICATE_KEY_RSA_LENGTH_NOT_SUPPORTED"
	//
	ErrorServerCertificateKeyAlgorithmNotSupported = "SERVER_CERTIFICATE_KEY_ALGORITHM_NOT_SUPPORTED"
	//
	ErrorServerCertificateKeyInvalid = "SERVER_CERTIFICATE_KEY_INVALID"
	//
	ErrorCrlCorrupt = "CRL_CORRUPT"
	//
	ErrorCrlNameInvalid = "CRL_NAME_INVALID"
	//
	ErrorCrlAlreadyExists = "CRL_ALREADY_EXISTS"
	//
	ErrorCrlDoesNotExist = "CRL_DOES_NOT_EXIST"
	//
	ErrorCertificateLibraryCorrupt = "CERTIFICATE_LIBRARY_CORRUPT"
	//
	ErrorCertificateCorrupt = "CERTIFICATE_CORRUPT"
	//
	ErrorCertificateNameInvalid = "CERTIFICATE_NAME_INVALID"
	//
	ErrorTrustedCertificateAlreadyExists = "TRUSTED_CERTIFICATE_ALREADY_EXISTS"
	//
	ErrorCertificateAlreadyExists = "CERTIFICATE_ALREADY_EXISTS"
	//
	ErrorCertificateDoesNotExist = "CERTIFICATE_DOES_NOT_EXIST"
	//
	ErrorRbacPermissionDenied = "RBAC_PERMISSION_DENIED"
	//
	ErrorRoleAlreadyExists = "ROLE_ALREADY_EXISTS"
	//
	ErrorRoleNotFound = "ROLE_NOT_FOUND"
	//
	ErrorSubjectAlreadyExists = "SUBJECT_ALREADY_EXISTS"
	//
	ErrorAuthServiceError = "AUTH_SERVICE_ERROR"
	//
	ErrorSubjectCannotBeResolved = "SUBJECT_CANNOT_BE_RESOLVED"
	//
	ErrorPoolAuthDisableFailedInvalidAccount = "POOL_AUTH_DISABLE_FAILED_INVALID_ACCOUNT"
	//
	ErrorPoolAuthDisableFailedPermissionDenied = "POOL_AUTH_DISABLE_FAILED_PERMISSION_DENIED"
	//
	ErrorPoolAuthDisableFailedWrongCredentials = "POOL_AUTH_DISABLE_FAILED_WRONG_CREDENTIALS"
	//
	ErrorPoolAuthDisableFailed = "POOL_AUTH_DISABLE_FAILED"
	//
	ErrorAuthSetLdapsFailed = "AUTH_SET_LDAPS_FAILED"
	//
	ErrorPoolAuthEnableFailedDuplicateHostname = "POOL_AUTH_ENABLE_FAILED_DUPLICATE_HOSTNAME"
	//
	ErrorPoolAuthEnableFailedSetupTLSConnection = "POOL_AUTH_ENABLE_FAILED_SETUP_TLS_CONNECTION"
	//
	ErrorPoolAuthEnableFailedNoSupportEncryptType = "POOL_AUTH_ENABLE_FAILED_NO_SUPPORT_ENCRYPT_TYPE"
	//
	ErrorPoolAuthEnableFailedNoTrustedCerts = "POOL_AUTH_ENABLE_FAILED_NO_TRUSTED_CERTS"
	//
	ErrorPoolAuthEnableFailedInvalidTrustedCerts = "POOL_AUTH_ENABLE_FAILED_INVALID_TRUSTED_CERTS"
	//
	ErrorPoolAuthEnableFailedInvalidAccount = "POOL_AUTH_ENABLE_FAILED_INVALID_ACCOUNT"
	//
	ErrorPoolAuthEnableFailedInvalidOu = "POOL_AUTH_ENABLE_FAILED_INVALID_OU"
	//
	ErrorPoolAuthEnableFailedUnavailable = "POOL_AUTH_ENABLE_FAILED_UNAVAILABLE"
	//
	ErrorPoolAuthEnableFailedDomainLookupFailed = "POOL_AUTH_ENABLE_FAILED_DOMAIN_LOOKUP_FAILED"
	//
	ErrorPoolAuthEnableFailedPermissionDenied = "POOL_AUTH_ENABLE_FAILED_PERMISSION_DENIED"
	//
	ErrorPoolAuthEnableFailedWrongCredentials = "POOL_AUTH_ENABLE_FAILED_WRONG_CREDENTIALS"
	//
	ErrorPoolAuthEnableFailed = "POOL_AUTH_ENABLE_FAILED"
	//
	ErrorPoolAuthAlreadyEnabled = "POOL_AUTH_ALREADY_ENABLED"
	//
	ErrorAuthDisableFailedPermissionDenied = "AUTH_DISABLE_FAILED_PERMISSION_DENIED"
	//
	ErrorAuthDisableFailedWrongCredentials = "AUTH_DISABLE_FAILED_WRONG_CREDENTIALS"
	//
	ErrorAuthDisableFailed = "AUTH_DISABLE_FAILED"
	//
	ErrorAuthEnableFailedNoSupportEncryptType = "AUTH_ENABLE_FAILED_NO_SUPPORT_ENCRYPT_TYPE"
	//
	ErrorAuthEnableFailedNoTrustedCerts = "AUTH_ENABLE_FAILED_NO_TRUSTED_CERTS"
	//
	ErrorAuthEnableFailedInvalidTrustedCerts = "AUTH_ENABLE_FAILED_INVALID_TRUSTED_CERTS"
	//
	ErrorAuthEnableFailedInvalidAccount = "AUTH_ENABLE_FAILED_INVALID_ACCOUNT"
	//
	ErrorAuthEnableFailedInvalidOu = "AUTH_ENABLE_FAILED_INVALID_OU"
	//
	ErrorAuthEnableFailedUnavailable = "AUTH_ENABLE_FAILED_UNAVAILABLE"
	//
	ErrorAuthEnableFailedDomainLookupFailed = "AUTH_ENABLE_FAILED_DOMAIN_LOOKUP_FAILED"
	//
	ErrorAuthEnableFailedPermissionDenied = "AUTH_ENABLE_FAILED_PERMISSION_DENIED"
	//
	ErrorAuthEnableFailedWrongCredentials = "AUTH_ENABLE_FAILED_WRONG_CREDENTIALS"
	//
	ErrorAuthEnableFailed = "AUTH_ENABLE_FAILED"
	//
	ErrorAuthSetupTLSConnection = "AUTH_SETUP_TLS_CONNECTION"
	//
	ErrorAuthNoTrustedCerts = "AUTH_NO_TRUSTED_CERTS"
	//
	ErrorAuthInvalidTrustedCerts = "AUTH_INVALID_TRUSTED_CERTS"
	//
	ErrorAuthIsDisabled = "AUTH_IS_DISABLED"
	//
	ErrorAuthUnknownType = "AUTH_UNKNOWN_TYPE"
	//
	ErrorAuthAlreadyEnabled = "AUTH_ALREADY_ENABLED"
	//
	ErrorDomainBuilderError = "DOMAIN_BUILDER_ERROR"
	//
	ErrorSrNotAttached = "SR_NOT_ATTACHED"
	//
	ErrorSrAttached = "SR_ATTACHED"
	//
	ErrorXenapiPluginFailure = "XENAPI_PLUGIN_FAILURE"
	//
	ErrorXenapiMissingPlugin = "XENAPI_MISSING_PLUGIN"
	//
	ErrorNoLocalStorage = "NO_LOCAL_STORAGE"
	//
	ErrorXapiHookFailed = "XAPI_HOOK_FAILED"
	//
	ErrorSystemStatusMustUseTarOnOem = "SYSTEM_STATUS_MUST_USE_TAR_ON_OEM"
	//
	ErrorSystemStatusRetrievalFailed = "SYSTEM_STATUS_RETRIEVAL_FAILED"
	//
	ErrorHostEvacuateInProgress = "HOST_EVACUATE_IN_PROGRESS"
	//
	ErrorCannotEvacuateHost = "CANNOT_EVACUATE_HOST"
	//
	ErrorIncompatibleClusterStackActive = "INCOMPATIBLE_CLUSTER_STACK_ACTIVE"
	//
	ErrorIncompatibleStatefileSr = "INCOMPATIBLE_STATEFILE_SR"
	//
	ErrorHaOperationWouldBreakFailoverPlan = "HA_OPERATION_WOULD_BREAK_FAILOVER_PLAN"
	//
	ErrorHaConstraintViolationNetworkNotShared = "HA_CONSTRAINT_VIOLATION_NETWORK_NOT_SHARED"
	//
	ErrorHaConstraintViolationSrNotShared = "HA_CONSTRAINT_VIOLATION_SR_NOT_SHARED"
	//
	ErrorHaCannotChangeBondStatusOfMgmtIface = "HA_CANNOT_CHANGE_BOND_STATUS_OF_MGMT_IFACE"
	//
	ErrorHaFailedToFormLiveset = "HA_FAILED_TO_FORM_LIVESET"
	//
	ErrorHaHostCannotAccessStatefile = "HA_HOST_CANNOT_ACCESS_STATEFILE"
	//
	ErrorHaHeartbeatDaemonStartupFailed = "HA_HEARTBEAT_DAEMON_STARTUP_FAILED"
	//
	ErrorHaPoolIsEnabledButHostIsDisabled = "HA_POOL_IS_ENABLED_BUT_HOST_IS_DISABLED"
	//
	ErrorHaLostStatefile = "HA_LOST_STATEFILE"
	//
	ErrorHaNoPlan = "HA_NO_PLAN"
	//
	ErrorHaAbortNewMaster = "HA_ABORT_NEW_MASTER"
	//
	ErrorHaShouldBeFenced = "HA_SHOULD_BE_FENCED"
	//
	ErrorHaTooFewHosts = "HA_TOO_FEW_HOSTS"
	//
	ErrorHaHostCannotSeePeers = "HA_HOST_CANNOT_SEE_PEERS"
	//
	ErrorHaNotInstalled = "HA_NOT_INSTALLED"
	//
	ErrorHaDisableInProgress = "HA_DISABLE_IN_PROGRESS"
	//
	ErrorHaEnableInProgress = "HA_ENABLE_IN_PROGRESS"
	//
	ErrorHaNotEnabled = "HA_NOT_ENABLED"
	//
	ErrorHaIsEnabled = "HA_IS_ENABLED"
	//
	ErrorHaHostIsArmed = "HA_HOST_IS_ARMED"
	//
	ErrorBallooningTimeoutBeforeMigration = "BALLOONING_TIMEOUT_BEFORE_MIGRATION"
	//
	ErrorBallooningDisabled = "BALLOONING_DISABLED"
	//
	ErrorClientError = "CLIENT_ERROR"
	//
	ErrorDuplicateMacSeed = "DUPLICATE_MAC_SEED"
	//
	ErrorDuplicateVM = "DUPLICATE_VM"
	//
	ErrorXmlrpcUnmarshalFailure = "XMLRPC_UNMARSHAL_FAILURE"
	//
	ErrorFeatureRestricted = "FEATURE_RESTRICTED"
	//
	ErrorActivationWhileNotFree = "ACTIVATION_WHILE_NOT_FREE"
	//
	ErrorLicenseFileDeprecated = "LICENSE_FILE_DEPRECATED"
	//
	ErrorLicenseCheckoutError = "LICENSE_CHECKOUT_ERROR"
	//
	ErrorMissingConnectionDetails = "MISSING_CONNECTION_DETAILS"
	//
	ErrorInvalidEdition = "INVALID_EDITION"
	//
	ErrorV6dFailure = "V6D_FAILURE"
	//
	ErrorLicenseDoesNotSupportXha = "LICENSE_DOES_NOT_SUPPORT_XHA"
	//
	ErrorLicenseCannotDowngradeWhileInPool = "LICENSE_CANNOT_DOWNGRADE_WHILE_IN_POOL"
	//
	ErrorLicenseProcessingError = "LICENSE_PROCESSING_ERROR"
	//
	ErrorLicenseHostPoolMismatch = "LICENSE_HOST_POOL_MISMATCH"
	//
	ErrorLicenseDoesNotSupportPooling = "LICENSE_DOES_NOT_SUPPORT_POOLING"
	//
	ErrorLicenceRestriction = "LICENCE_RESTRICTION"
	//
	ErrorLicenseExpired = "LICENSE_EXPIRED"
	//
	ErrorRestoreScriptFailed = "RESTORE_SCRIPT_FAILED"
	//
	ErrorBackupScriptFailed = "BACKUP_SCRIPT_FAILED"
	//
	ErrorCannotFindStatePartition = "CANNOT_FIND_STATE_PARTITION"
	//
	ErrorNotAllowedOnOemEdition = "NOT_ALLOWED_ON_OEM_EDITION"
	//
	ErrorOnlyAllowedOnOemEdition = "ONLY_ALLOWED_ON_OEM_EDITION"
	//
	ErrorCannotFindOemBackupPartition = "CANNOT_FIND_OEM_BACKUP_PARTITION"
	//
	ErrorPatchApplyFailedBackupFilesExist = "PATCH_APPLY_FAILED_BACKUP_FILES_EXIST"
	//
	ErrorPatchApplyFailed = "PATCH_APPLY_FAILED"
	//
	ErrorPatchPrecheckFailedIsoMounted = "PATCH_PRECHECK_FAILED_ISO_MOUNTED"
	//
	ErrorUpdatePrecheckFailedOutOfSpace = "UPDATE_PRECHECK_FAILED_OUT_OF_SPACE"
	//
	ErrorPatchPrecheckFailedOutOfSpace = "PATCH_PRECHECK_FAILED_OUT_OF_SPACE"
	//
	ErrorPatchPrecheckFailedVMRunning = "PATCH_PRECHECK_FAILED_VM_RUNNING"
	//
	ErrorPatchPrecheckFailedWrongServerBuild = "PATCH_PRECHECK_FAILED_WRONG_SERVER_BUILD"
	//
	ErrorPatchPrecheckFailedWrongServerVersion = "PATCH_PRECHECK_FAILED_WRONG_SERVER_VERSION"
	//
	ErrorPatchPrecheckFailedPrerequisiteMissing = "PATCH_PRECHECK_FAILED_PREREQUISITE_MISSING"
	//
	ErrorPatchPrecheckFailedUnknownError = "PATCH_PRECHECK_FAILED_UNKNOWN_ERROR"
	//
	ErrorUpdatePrecheckFailedGpgkeyNotImported = "UPDATE_PRECHECK_FAILED_GPGKEY_NOT_IMPORTED"
	//
	ErrorUpdatePrecheckFailedWrongServerVersion = "UPDATE_PRECHECK_FAILED_WRONG_SERVER_VERSION"
	//
	ErrorUpdatePrecheckFailedConflictPresent = "UPDATE_PRECHECK_FAILED_CONFLICT_PRESENT"
	//
	ErrorUpdatePrecheckFailedPrerequisiteMissing = "UPDATE_PRECHECK_FAILED_PREREQUISITE_MISSING"
	//
	ErrorUpdatePrecheckFailedUnknownError = "UPDATE_PRECHECK_FAILED_UNKNOWN_ERROR"
	//
	ErrorUpdateApplyFailed = "UPDATE_APPLY_FAILED"
	//
	ErrorCouldNotUpdateIgmpSnoopingEverywhere = "COULD_NOT_UPDATE_IGMP_SNOOPING_EVERYWHERE"
	//
	ErrorUpdatePoolApplyFailed = "UPDATE_POOL_APPLY_FAILED"
	//
	ErrorUpdateAlreadyAppliedInPool = "UPDATE_ALREADY_APPLIED_IN_POOL"
	//
	ErrorUpdateAlreadyApplied = "UPDATE_ALREADY_APPLIED"
	//
	ErrorPatchAlreadyApplied = "PATCH_ALREADY_APPLIED"
	//
	ErrorCannotFetchPatch = "CANNOT_FETCH_PATCH"
	//
	ErrorCannotFindUpdate = "CANNOT_FIND_UPDATE"
	//
	ErrorCannotFindPatch = "CANNOT_FIND_PATCH"
	//
	ErrorUpdateIsApplied = "UPDATE_IS_APPLIED"
	//
	ErrorPatchIsApplied = "PATCH_IS_APPLIED"
	//
	ErrorUpdateAlreadyExists = "UPDATE_ALREADY_EXISTS"
	//
	ErrorPatchAlreadyExists = "PATCH_ALREADY_EXISTS"
	//
	ErrorInvalidPatchWithLog = "INVALID_PATCH_WITH_LOG"
	//
	ErrorInvalidUpdate = "INVALID_UPDATE"
	//
	ErrorInvalidPatch = "INVALID_PATCH"
	//
	ErrorOutOfSpace = "OUT_OF_SPACE"
	//
	ErrorTooBusy = "TOO_BUSY"
	//
	ErrorTooManyPendingTasks = "TOO_MANY_PENDING_TASKS"
	//
	ErrorTaskCancelled = "TASK_CANCELLED"
	//
	ErrorDefaultSrNotFound = "DEFAULT_SR_NOT_FOUND"
	//
	ErrorSrNotShared = "SR_NOT_SHARED"
	//
	ErrorWlbConnectionReset = "WLB_CONNECTION_RESET"
	//
	ErrorWlbURLInvalid = "WLB_URL_INVALID"
	//
	ErrorWlbInternalError = "WLB_INTERNAL_ERROR"
	//
	ErrorWlbXenserverMalformedResponse = "WLB_XENSERVER_MALFORMED_RESPONSE"
	//
	ErrorWlbXenserverAuthenticationFailed = "WLB_XENSERVER_AUTHENTICATION_FAILED"
	//
	ErrorWlbXenserverTimeout = "WLB_XENSERVER_TIMEOUT"
	//
	ErrorWlbXenserverUnknownHost = "WLB_XENSERVER_UNKNOWN_HOST"
	//
	ErrorWlbXenserverConnectionRefused = "WLB_XENSERVER_CONNECTION_REFUSED"
	//
	ErrorWlbMalformedResponse = "WLB_MALFORMED_RESPONSE"
	//
	ErrorWlbMalformedRequest = "WLB_MALFORMED_REQUEST"
	//
	ErrorWlbAuthenticationFailed = "WLB_AUTHENTICATION_FAILED"
	//
	ErrorWlbTimeout = "WLB_TIMEOUT"
	//
	ErrorWlbUnknownHost = "WLB_UNKNOWN_HOST"
	//
	ErrorWlbConnectionRefused = "WLB_CONNECTION_REFUSED"
	//
	ErrorWlbDisabled = "WLB_DISABLED"
	//
	ErrorWlbNotInitialized = "WLB_NOT_INITIALIZED"
	//
	ErrorPoolJoiningHostMustHaveOnlyOneIPOnClusteringNetwork = "POOL_JOINING_HOST_MUST_HAVE_ONLY_ONE_IP_ON_CLUSTERING_NETWORK"
	//
	ErrorPoolJoiningPoolCannotEnableClusteringOnVlanNetwork = "POOL_JOINING_POOL_CANNOT_ENABLE_CLUSTERING_ON_VLAN_NETWORK"
	//
	ErrorPoolJoiningSmFeaturesIncompatible = "POOL_JOINING_SM_FEATURES_INCOMPATIBLE"
	//
	ErrorPoolJoiningHostTrustedCertificatesConflict = "POOL_JOINING_HOST_TRUSTED_CERTIFICATES_CONFLICT"
	//
	ErrorPoolJoiningHostCaCertificatesConflict = "POOL_JOINING_HOST_CA_CERTIFICATES_CONFLICT"
	//
	ErrorPoolJoiningPoolBundleEmptyAfterImport = "POOL_JOINING_POOL_BUNDLE_EMPTY_AFTER_IMPORT"
	//
	ErrorPoolJoiningMasterCertificateNotInPoolBundle = "POOL_JOINING_MASTER_CERTIFICATE_NOT_IN_POOL_BUNDLE"
	//
	ErrorPoolJoiningHostTLSVerificationMismatch = "POOL_JOINING_HOST_TLS_VERIFICATION_MISMATCH"
	//
	ErrorPoolJoiningHostHasNetworkSriovs = "POOL_JOINING_HOST_HAS_NETWORK_SRIOVS"
	//
	ErrorPoolJoiningHostHasTunnels = "POOL_JOINING_HOST_HAS_TUNNELS"
	//
	ErrorPoolJoiningHostHasBonds = "POOL_JOINING_HOST_HAS_BONDS"
	//
	ErrorPoolJoiningHostHasNonManagementVlans = "POOL_JOINING_HOST_HAS_NON_MANAGEMENT_VLANS"
	//
	ErrorPoolJoiningHostManagementVlanDoesNotMatch = "POOL_JOINING_HOST_MANAGEMENT_VLAN_DOES_NOT_MATCH"
	//
	ErrorPoolJoiningHostMustOnlyHavePhysicalPifs = "POOL_JOINING_HOST_MUST_ONLY_HAVE_PHYSICAL_PIFS"
	//
	ErrorPoolJoiningHostMustHaveSameDBSchema = "POOL_JOINING_HOST_MUST_HAVE_SAME_DB_SCHEMA"
	//
	ErrorPoolJoiningHostMustHaveSameAPIVersion = "POOL_JOINING_HOST_MUST_HAVE_SAME_API_VERSION"
	//
	ErrorPoolJoiningHostMustHaveSameProductVersion = "POOL_JOINING_HOST_MUST_HAVE_SAME_PRODUCT_VERSION"
	//
	ErrorPoolJoiningExternalAuthMismatch = "POOL_JOINING_EXTERNAL_AUTH_MISMATCH"
	//
	ErrorPoolJoiningHostMustHavePhysicalManagementNic = "POOL_JOINING_HOST_MUST_HAVE_PHYSICAL_MANAGEMENT_NIC"
	//
	ErrorJoiningHostServiceFailed = "JOINING_HOST_SERVICE_FAILED"
	//
	ErrorJoiningHostConnectionFailed = "JOINING_HOST_CONNECTION_FAILED"
	//
	ErrorJoiningHostCannotBeMasterOfOtherHosts = "JOINING_HOST_CANNOT_BE_MASTER_OF_OTHER_HOSTS"
	//
	ErrorJoiningHostCannotHaveVmsWithCurrentOperations = "JOINING_HOST_CANNOT_HAVE_VMS_WITH_CURRENT_OPERATIONS"
	//
	ErrorJoiningHostCannotHaveRunningVms = "JOINING_HOST_CANNOT_HAVE_RUNNING_VMS"
	//
	ErrorJoiningHostCannotHaveRunningOrSuspendedVms = "JOINING_HOST_CANNOT_HAVE_RUNNING_OR_SUSPENDED_VMS"
	//
	ErrorJoiningHostCannotContainSharedSrs = "JOINING_HOST_CANNOT_CONTAIN_SHARED_SRS"
	//
	ErrorHostsNotHomogeneous = "HOSTS_NOT_HOMOGENEOUS"
	//
	ErrorHostsNotCompatible = "HOSTS_NOT_COMPATIBLE"
	//
	ErrorNotInEmergencyMode = "NOT_IN_EMERGENCY_MODE"
	//
	ErrorRestoreTargetMgmtIfNotInBackup = "RESTORE_TARGET_MGMT_IF_NOT_IN_BACKUP"
	//
	ErrorRestoreTargetMissingDevice = "RESTORE_TARGET_MISSING_DEVICE"
	//
	ErrorRestoreIncompatibleVersion = "RESTORE_INCOMPATIBLE_VERSION"
	//
	ErrorImportIncompatibleVersion = "IMPORT_INCOMPATIBLE_VERSION"
	//
	ErrorImportErrorUnexpectedFile = "IMPORT_ERROR_UNEXPECTED_FILE"
	//
	ErrorImportErrorAttachedDisksNotFound = "IMPORT_ERROR_ATTACHED_DISKS_NOT_FOUND"
	//
	ErrorImportErrorFailedToFindObject = "IMPORT_ERROR_FAILED_TO_FIND_OBJECT"
	//
	ErrorImportErrorCannotHandleChunked = "IMPORT_ERROR_CANNOT_HANDLE_CHUNKED"
	//
	ErrorImportErrorSomeChecksumsFailed = "IMPORT_ERROR_SOME_CHECKSUMS_FAILED"
	//
	ErrorImportErrorPrematureEOF = "IMPORT_ERROR_PREMATURE_EOF"
	//
	ErrorImportError = "IMPORT_ERROR"
	//
	ErrorVMPciBusFull = "VM_PCI_BUS_FULL"
	//
	ErrorNvidiaSriovMisconfigured = "NVIDIA_SRIOV_MISCONFIGURED"
	//
	ErrorNvidiaToolsError = "NVIDIA_TOOLS_ERROR"
	//
	ErrorVgpuGuestDriverLimit = "VGPU_GUEST_DRIVER_LIMIT"
	//
	ErrorVgpuSuspensionNotSupported = "VGPU_SUSPENSION_NOT_SUPPORTED"
	//
	ErrorVgpuDestinationIncompatible = "VGPU_DESTINATION_INCOMPATIBLE"
	//
	ErrorVgpuTypeNotCompatible = "VGPU_TYPE_NOT_COMPATIBLE"
	//
	ErrorVgpuTypeNotCompatibleWithRunningType = "VGPU_TYPE_NOT_COMPATIBLE_WITH_RUNNING_TYPE"
	//
	ErrorVgpuTypeNoLongerSupported = "VGPU_TYPE_NO_LONGER_SUPPORTED"
	//
	ErrorVgpuTypeNotSupported = "VGPU_TYPE_NOT_SUPPORTED"
	//
	ErrorVgpuTypeNotEnabled = "VGPU_TYPE_NOT_ENABLED"
	//
	ErrorPgpuInsufficientCapacityForVgpu = "PGPU_INSUFFICIENT_CAPACITY_FOR_VGPU"
	//
	ErrorPgpuNotCompatibleWithGpuGroup = "PGPU_NOT_COMPATIBLE_WITH_GPU_GROUP"
	//
	ErrorPgpuInUseByVM = "PGPU_IN_USE_BY_VM"
	//
	ErrorSessionNotRegistered = "SESSION_NOT_REGISTERED"
	//
	ErrorEventFromTokenParseFailure = "EVENT_FROM_TOKEN_PARSE_FAILURE"
	//
	ErrorEventSubscriptionParseFailure = "EVENT_SUBSCRIPTION_PARSE_FAILURE"
	//
	ErrorEventsLost = "EVENTS_LOST"
	//
	ErrorInvalidDevice = "INVALID_DEVICE"
	//
	ErrorGpuGroupContainsNoPgpus = "GPU_GROUP_CONTAINS_NO_PGPUS"
	//
	ErrorGpuGroupContainsPgpu = "GPU_GROUP_CONTAINS_PGPU"
	//
	ErrorGpuGroupContainsVgpu = "GPU_GROUP_CONTAINS_VGPU"
	//
	ErrorNetworkContainsVif = "NETWORK_CONTAINS_VIF"
	//
	ErrorNetworkContainsPif = "NETWORK_CONTAINS_PIF"
	//
	ErrorDeviceNotAttached = "DEVICE_NOT_ATTACHED"
	//
	ErrorDeviceAlreadyExists = "DEVICE_ALREADY_EXISTS"
	//
	ErrorDeviceAlreadyDetached = "DEVICE_ALREADY_DETACHED"
	//
	ErrorDeviceAlreadyAttached = "DEVICE_ALREADY_ATTACHED"
	//
	ErrorNotImplemented = "NOT_IMPLEMENTED"
	//
	ErrorPbdExists = "PBD_EXISTS"
	//
	ErrorSmPluginCommunicationFailure = "SM_PLUGIN_COMMUNICATION_FAILURE"
	//
	ErrorClusteredSrDegraded = "CLUSTERED_SR_DEGRADED"
	//
	ErrorSrIndestructible = "SR_INDESTRUCTIBLE"
	//
	ErrorSrNotSharable = "SR_NOT_SHARABLE"
	//
	ErrorSrOperationNotSupported = "SR_OPERATION_NOT_SUPPORTED"
	//
	ErrorSrDeviceInUse = "SR_DEVICE_IN_USE"
	//
	ErrorSrNotEmpty = "SR_NOT_EMPTY"
	//
	ErrorSrVdiLockingFailed = "SR_VDI_LOCKING_FAILED"
	//
	ErrorSrUnknownDriver = "SR_UNKNOWN_DRIVER"
	//
	ErrorSrBackendFailure = "SR_BACKEND_FAILURE"
	//
	ErrorSrHasMultiplePbds = "SR_HAS_MULTIPLE_PBDS"
	//
	ErrorSrHasNoPbds = "SR_HAS_NO_PBDS"
	//
	ErrorSrUUIDExists = "SR_UUID_EXISTS"
	//
	ErrorOperationPartiallyFailed = "OPERATION_PARTIALLY_FAILED"
	//
	ErrorCannotCreateStateFile = "CANNOT_CREATE_STATE_FILE"
	//
	ErrorVifNotInMap = "VIF_NOT_IN_MAP"
	//
	ErrorVdiIsEncrypted = "VDI_IS_ENCRYPTED"
	//
	ErrorVdiNoCbtMetadata = "VDI_NO_CBT_METADATA"
	//
	ErrorVdiCbtEnabled = "VDI_CBT_ENABLED"
	//
	ErrorVdiNotInMap = "VDI_NOT_IN_MAP"
	//
	ErrorVdiOnBootModeIncompatibleWithOperation = "VDI_ON_BOOT_MODE_INCOMPATIBLE_WITH_OPERATION"
	//
	ErrorVdiIoError = "VDI_IO_ERROR"
	//
	ErrorVdiNotManaged = "VDI_NOT_MANAGED"
	//
	ErrorVdiIncompatibleType = "VDI_INCOMPATIBLE_TYPE"
	//
	ErrorVdiMissing = "VDI_MISSING"
	//
	ErrorVdiContentIDMissing = "VDI_CONTENT_ID_MISSING"
	//
	ErrorVdiLocationMissing = "VDI_LOCATION_MISSING"
	//
	ErrorVdiHasRrds = "VDI_HAS_RRDS"
	//
	ErrorVdiNotAvailable = "VDI_NOT_AVAILABLE"
	//
	ErrorHostCdDriveEmpty = "HOST_CD_DRIVE_EMPTY"
	//
	ErrorDiskVbdMustBeReadwriteForHvm = "DISK_VBD_MUST_BE_READWRITE_FOR_HVM"
	//
	ErrorVMRequiresVusb = "VM_REQUIRES_VUSB"
	//
	ErrorVbdCdsMustBeReadonly = "VBD_CDS_MUST_BE_READONLY"
	//
	ErrorVdiIsNotIso = "VDI_IS_NOT_ISO"
	//
	ErrorVdiIsAPhysicalDevice = "VDI_IS_A_PHYSICAL_DEVICE"
	//
	ErrorVdiNotSparse = "VDI_NOT_SPARSE"
	//
	ErrorVdiTooLarge = "VDI_TOO_LARGE"
	//
	ErrorVdiTooSmall = "VDI_TOO_SMALL"
	//
	ErrorVdiReadonly = "VDI_READONLY"
	//
	ErrorVdiIsSharable = "VDI_IS_SHARABLE"
	//
	ErrorVdiInUse = "VDI_IN_USE"
	//
	ErrorSrUnhealthy = "SR_UNHEALTHY"
	//
	ErrorSrIsCacheSr = "SR_IS_CACHE_SR"
	//
	ErrorSrRequiresUpgrade = "SR_REQUIRES_UPGRADE"
	//
	ErrorSrHasPbd = "SR_HAS_PBD"
	//
	ErrorSrSuspendSpaceInsufficient = "SR_SUSPEND_SPACE_INSUFFICIENT"
	//
	ErrorSrSourceSpaceInsufficient = "SR_SOURCE_SPACE_INSUFFICIENT"
	//
	ErrorSrFull = "SR_FULL"
	//
	ErrorSrAttachFailed = "SR_ATTACH_FAILED"
	//
	ErrorObjectNolongerExists = "OBJECT_NOLONGER_EXISTS"
	//
	ErrorFailedToStartEmulator = "FAILED_TO_START_EMULATOR"
	//
	ErrorUnknownBootloader = "UNKNOWN_BOOTLOADER"
	//
	ErrorBootloaderFailed = "BOOTLOADER_FAILED"
	//
	ErrorProvisionFailedOutOfSpace = "PROVISION_FAILED_OUT_OF_SPACE"
	//
	ErrorRevertOnlyAllowedOnSnapshot = "REVERT_ONLY_ALLOWED_ON_SNAPSHOT"
	//
	ErrorProvisionOnlyAllowedOnTemplate = "PROVISION_ONLY_ALLOWED_ON_TEMPLATE"
	//
	ErrorNotSystemDomain = "NOT_SYSTEM_DOMAIN"
	//
	ErrorCannotResetControlDomain = "CANNOT_RESET_CONTROL_DOMAIN"
	//
	ErrorDomainExists = "DOMAIN_EXISTS"
	//
	ErrorVMPvDriversInUse = "VM_PV_DRIVERS_IN_USE"
	//
	ErrorVmsFailedToCooperate = "VMS_FAILED_TO_COOPERATE"
	//
	ErrorVMAttachedToMoreThanOneVdiWithTimeoffsetMarkedAsResetOnBoot = "VM_ATTACHED_TO_MORE_THAN_ONE_VDI_WITH_TIMEOFFSET_MARKED_AS_RESET_ON_BOOT"
	//
	ErrorVMHalted = "VM_HALTED"
	//
	ErrorVMRebooted = "VM_REBOOTED"
	//
	ErrorVMCrashed = "VM_CRASHED"
	//
	ErrorVMNotResidentHere = "VM_NOT_RESIDENT_HERE"
	//
	ErrorIllegalVbdDevice = "ILLEGAL_VBD_DEVICE"
	//
	ErrorVMDuplicateVbdDevice = "VM_DUPLICATE_VBD_DEVICE"
	//
	ErrorVMSuspendTimeout = "VM_SUSPEND_TIMEOUT"
	//
	ErrorVMShutdownTimeout = "VM_SHUTDOWN_TIMEOUT"
	//
	ErrorVMMemoryTargetWaitTimeout = "VM_MEMORY_TARGET_WAIT_TIMEOUT"
	//
	ErrorVMMemorySizeTooLow = "VM_MEMORY_SIZE_TOO_LOW"
	//
	ErrorVMCannotDeleteDefaultTemplate = "VM_CANNOT_DELETE_DEFAULT_TEMPLATE"
	//
	ErrorVMNonSuspendable = "VM_NON_SUSPENDABLE"
	//
	ErrorVMLacksFeature = "VM_LACKS_FEATURE"
	//
	ErrorVMOldPvDrivers = "VM_OLD_PV_DRIVERS"
	//
	ErrorVMFailedSuspendAcknowledgment = "VM_FAILED_SUSPEND_ACKNOWLEDGMENT"
	//
	ErrorVMFailedShutdownAcknowledgment = "VM_FAILED_SHUTDOWN_ACKNOWLEDGMENT"
	//
	ErrorVMMissingPvDrivers = "VM_MISSING_PV_DRIVERS"
	//
	ErrorVMMigrateContactRemoteServiceFailed = "VM_MIGRATE_CONTACT_REMOTE_SERVICE_FAILED"
	//
	ErrorVMMigrateFailed = "VM_MIGRATE_FAILED"
	//
	ErrorVMNoCrashdumpSr = "VM_NO_CRASHDUMP_SR"
	//
	ErrorVMNoSuspendSr = "VM_NO_SUSPEND_SR"
	//
	ErrorHostCannotAttachNetwork = "HOST_CANNOT_ATTACH_NETWORK"
	//
	ErrorVMHasNoSuspendVdi = "VM_HAS_NO_SUSPEND_VDI"
	//
	ErrorVMHasSriovVif = "VM_HAS_SRIOV_VIF"
	//
	ErrorVMHasVgpu = "VM_HAS_VGPU"
	//
	ErrorVMHasPciAttached = "VM_HAS_PCI_ATTACHED"
	//
	ErrorVMHostIncompatibleVirtualHardwarePlatformVersion = "VM_HOST_INCOMPATIBLE_VIRTUAL_HARDWARE_PLATFORM_VERSION"
	//
	ErrorVMHostIncompatibleVersion = "VM_HOST_INCOMPATIBLE_VERSION"
	//
	ErrorVMHostIncompatibleVersionMigrate = "VM_HOST_INCOMPATIBLE_VERSION_MIGRATE"
	//
	ErrorVMRequiresIommu = "VM_REQUIRES_IOMMU"
	//
	ErrorVMRequiresVgpu = "VM_REQUIRES_VGPU"
	//
	ErrorVMRequiresGpu = "VM_REQUIRES_GPU"
	//
	ErrorVMRequiresNetwork = "VM_REQUIRES_NETWORK"
	//
	ErrorVMRequiresVdi = "VM_REQUIRES_VDI"
	//
	ErrorVMRequiresSr = "VM_REQUIRES_SR"
	//
	ErrorVMUnsafeBoot = "VM_UNSAFE_BOOT"
	//
	ErrorVMCheckpointResumeFailed = "VM_CHECKPOINT_RESUME_FAILED"
	//
	ErrorVMCheckpointSuspendFailed = "VM_CHECKPOINT_SUSPEND_FAILED"
	//
	ErrorVMRevertFailed = "VM_REVERT_FAILED"
	//
	ErrorXenVssReqErrorCreatingSnapshotXMLString = "XEN_VSS_REQ_ERROR_CREATING_SNAPSHOT_XML_STRING"
	//
	ErrorXenVssReqErrorCreatingSnapshot = "XEN_VSS_REQ_ERROR_CREATING_SNAPSHOT"
	//
	ErrorXenVssReqErrorPreparingWriters = "XEN_VSS_REQ_ERROR_PREPARING_WRITERS"
	//
	ErrorXenVssReqErrorAddingVolumeToSnapsetFailed = "XEN_VSS_REQ_ERROR_ADDING_VOLUME_TO_SNAPSET_FAILED"
	//
	ErrorXenVssReqErrorStartSnapshotSetFailed = "XEN_VSS_REQ_ERROR_START_SNAPSHOT_SET_FAILED"
	//
	ErrorXenVssReqErrorNoVolumesSupported = "XEN_VSS_REQ_ERROR_NO_VOLUMES_SUPPORTED"
	//
	ErrorXenVssReqErrorProvNotLoaded = "XEN_VSS_REQ_ERROR_PROV_NOT_LOADED"
	//
	ErrorXenVssReqErrorInitFailed = "XEN_VSS_REQ_ERROR_INIT_FAILED"
	//
	ErrorVMSnapshotWithQuiesceNotSupported = "VM_SNAPSHOT_WITH_QUIESCE_NOT_SUPPORTED"
	//
	ErrorVMSnapshotWithQuiescePluginDeosNotRespond = "VM_SNAPSHOT_WITH_QUIESCE_PLUGIN_DEOS_NOT_RESPOND"
	//
	ErrorVMSnapshotWithQuiesceTimeout = "VM_SNAPSHOT_WITH_QUIESCE_TIMEOUT"
	//
	ErrorVMSnapshotWithQuiesceFailed = "VM_SNAPSHOT_WITH_QUIESCE_FAILED"
	//
	ErrorVMSnapshotFailed = "VM_SNAPSHOT_FAILED"
	//
	ErrorVMNoEmptyCdVbd = "VM_NO_EMPTY_CD_VBD"
	//
	ErrorVbdMissing = "VBD_MISSING"
	//
	ErrorVbdTrayLocked = "VBD_TRAY_LOCKED"
	//
	ErrorVbdIsEmpty = "VBD_IS_EMPTY"
	//
	ErrorVbdNotEmpty = "VBD_NOT_EMPTY"
	//
	ErrorVbdNotUnpluggable = "VBD_NOT_UNPLUGGABLE"
	//
	ErrorVbdNotRemovableMedia = "VBD_NOT_REMOVABLE_MEDIA"
	//
	ErrorOtherOperationInProgress = "OTHER_OPERATION_IN_PROGRESS"
	//
	ErrorVMIsSnapshot = "VM_IS_SNAPSHOT"
	//
	ErrorVMIsTemplate = "VM_IS_TEMPLATE"
	//
	ErrorVMBadPowerState = "VM_BAD_POWER_STATE"
	//
	ErrorVlanTagInvalid = "VLAN_TAG_INVALID"
	//
	ErrorBridgeNameExists = "BRIDGE_NAME_EXISTS"
	//
	ErrorBridgeNotAvailable = "BRIDGE_NOT_AVAILABLE"
	//
	ErrorPifTunnelStillExists = "PIF_TUNNEL_STILL_EXISTS"
	//
	ErrorIsTunnelAccessPif = "IS_TUNNEL_ACCESS_PIF"
	//
	ErrorTransportPifNotConfigured = "TRANSPORT_PIF_NOT_CONFIGURED"
	//
	ErrorOpenvswitchNotActive = "OPENVSWITCH_NOT_ACTIVE"
	//
	ErrorCouldNotFindNetworkInterfaceWithSpecifiedDeviceNameAndMacAddress = "COULD_NOT_FIND_NETWORK_INTERFACE_WITH_SPECIFIED_DEVICE_NAME_AND_MAC_ADDRESS"
	//
	ErrorDuplicatePifDeviceName = "DUPLICATE_PIF_DEVICE_NAME"
	//
	ErrorMacInvalid = "MAC_INVALID"
	//
	ErrorMacDoesNotExist = "MAC_DOES_NOT_EXIST"
	//
	ErrorMacStillExists = "MAC_STILL_EXISTS"
	//
	ErrorCannotPlugVif = "CANNOT_PLUG_VIF"
	//
	ErrorVifInUse = "VIF_IN_USE"
	//
	ErrorSlaveRequiresManagementInterface = "SLAVE_REQUIRES_MANAGEMENT_INTERFACE"
	//
	ErrorIncompatiblePifProperties = "INCOMPATIBLE_PIF_PROPERTIES"
	//
	ErrorCannotForgetSriovLogical = "CANNOT_FORGET_SRIOV_LOGICAL"
	//
	ErrorCannotChangePifProperties = "CANNOT_CHANGE_PIF_PROPERTIES"
	//
	ErrorCannotAddTunnelToVlanOnSriovLogical = "CANNOT_ADD_TUNNEL_TO_VLAN_ON_SRIOV_LOGICAL"
	//
	ErrorCannotAddTunnelToSriovLogical = "CANNOT_ADD_TUNNEL_TO_SRIOV_LOGICAL"
	//
	ErrorCannotAddTunnelToBondSlave = "CANNOT_ADD_TUNNEL_TO_BOND_SLAVE"
	//
	ErrorCannotAddVlanToBondSlave = "CANNOT_ADD_VLAN_TO_BOND_SLAVE"
	//
	ErrorCannotPlugBondSlave = "CANNOT_PLUG_BOND_SLAVE"
	//
	ErrorPifSriovStillExists = "PIF_SRIOV_STILL_EXISTS"
	//
	ErrorPifIsNotSriovCapable = "PIF_IS_NOT_SRIOV_CAPABLE"
	//
	ErrorPifUnmanaged = "PIF_UNMANAGED"
	//
	ErrorPifHasFcoeSrInUse = "PIF_HAS_FCOE_SR_IN_USE"
	//
	ErrorPifAllowsUnplug = "PIF_ALLOWS_UNPLUG"
	//
	ErrorPifDoesNotAllowUnplug = "PIF_DOES_NOT_ALLOW_UNPLUG"
	//
	ErrorPifNotPresent = "PIF_NOT_PRESENT"
	//
	ErrorRequiredPifIsUnplugged = "REQUIRED_PIF_IS_UNPLUGGED"
	//
	ErrorPifIncompatiblePrimaryAddressType = "PIF_INCOMPATIBLE_PRIMARY_ADDRESS_TYPE"
	//
	ErrorPifIsManagementInterface = "PIF_IS_MANAGEMENT_INTERFACE"
	//
	ErrorPifConfigurationError = "PIF_CONFIGURATION_ERROR"
	//
	ErrorPifBondMoreThanOneIP = "PIF_BOND_MORE_THAN_ONE_IP"
	//
	ErrorPifBondNeedsMoreMembers = "PIF_BOND_NEEDS_MORE_MEMBERS"
	//
	ErrorPifCannotBondCrossHost = "PIF_CANNOT_BOND_CROSS_HOST"
	//
	ErrorPifAlreadyBonded = "PIF_ALREADY_BONDED"
	//
	ErrorPifDeviceNotFound = "PIF_DEVICE_NOT_FOUND"
	//
	ErrorVlanInUse = "VLAN_IN_USE"
	//
	ErrorPifVlanStillExists = "PIF_VLAN_STILL_EXISTS"
	//
	ErrorPifVlanExists = "PIF_VLAN_EXISTS"
	//
	ErrorPifIsSriovLogical = "PIF_IS_SRIOV_LOGICAL"
	//
	ErrorPifIsVlan = "PIF_IS_VLAN"
	//
	ErrorPifIsNotPhysical = "PIF_IS_NOT_PHYSICAL"
	//
	ErrorPifIsPhysical = "PIF_IS_PHYSICAL"
	//
	ErrorCannotDestroySystemNetwork = "CANNOT_DESTROY_SYSTEM_NETWORK"
	//
	ErrorNetworkIncompatiblePurposes = "NETWORK_INCOMPATIBLE_PURPOSES"
	//
	ErrorNetworkUnmanaged = "NETWORK_UNMANAGED"
	//
	ErrorNetworkAlreadyConnected = "NETWORK_ALREADY_CONNECTED"
	//
	ErrorOperationBlocked = "OPERATION_BLOCKED"
	//
	ErrorOperationNotAllowed = "OPERATION_NOT_ALLOWED"
	//
	ErrorNetworkHasIncompatibleVlanOnSriovPifs = "NETWORK_HAS_INCOMPATIBLE_VLAN_ON_SRIOV_PIFS"
	//
	ErrorNetworkHasIncompatibleSriovPifs = "NETWORK_HAS_INCOMPATIBLE_SRIOV_PIFS"
	//
	ErrorNetworkIncompatibleWithTunnel = "NETWORK_INCOMPATIBLE_WITH_TUNNEL"
	//
	ErrorNetworkIncompatibleWithBond = "NETWORK_INCOMPATIBLE_WITH_BOND"
	//
	ErrorNetworkIncompatibleWithVlanOnSriov = "NETWORK_INCOMPATIBLE_WITH_VLAN_ON_SRIOV"
	//
	ErrorNetworkIncompatibleWithVlanOnBridge = "NETWORK_INCOMPATIBLE_WITH_VLAN_ON_BRIDGE"
	//
	ErrorNetworkIncompatibleWithSriov = "NETWORK_INCOMPATIBLE_WITH_SRIOV"
	//
	ErrorNetworkSriovDisableFailed = "NETWORK_SRIOV_DISABLE_FAILED"
	//
	ErrorNetworkSriovEnableFailed = "NETWORK_SRIOV_ENABLE_FAILED"
	//
	ErrorNetworkSriovAlreadyEnabled = "NETWORK_SRIOV_ALREADY_ENABLED"
	//
	ErrorNetworkSriovInsufficientCapacity = "NETWORK_SRIOV_INSUFFICIENT_CAPACITY"
	//
	ErrorDeviceDetachRejected = "DEVICE_DETACH_REJECTED"
	//
	ErrorDeviceDetachTimeout = "DEVICE_DETACH_TIMEOUT"
	//
	ErrorDeviceAttachTimeout = "DEVICE_ATTACH_TIMEOUT"
	//
	ErrorPifHasNoV6NetworkConfiguration = "PIF_HAS_NO_V6_NETWORK_CONFIGURATION"
	//
	ErrorPifHasNoNetworkConfiguration = "PIF_HAS_NO_NETWORK_CONFIGURATION"
	//
	ErrorAddressViolatesLockingConstraint = "ADDRESS_VIOLATES_LOCKING_CONSTRAINT"
	//
	ErrorInvalidCidrAddressSpecified = "INVALID_CIDR_ADDRESS_SPECIFIED"
	//
	ErrorInvalidIPAddressSpecified = "INVALID_IP_ADDRESS_SPECIFIED"
	//
	ErrorInterfaceHasNoIP = "INTERFACE_HAS_NO_IP"
	//
	ErrorHostBroken = "HOST_BROKEN"
	//
	ErrorHostXapiVersionHigherThanCoordinator = "HOST_XAPI_VERSION_HIGHER_THAN_COORDINATOR"
	//
	ErrorHostUnknownToMaster = "HOST_UNKNOWN_TO_MASTER"
	//
	ErrorHostMasterCannotTalkBack = "HOST_MASTER_CANNOT_TALK_BACK"
	//
	ErrorHostHasNoManagementIP = "HOST_HAS_NO_MANAGEMENT_IP"
	//
	ErrorHostStillBooting = "HOST_STILL_BOOTING"
	//
	ErrorHostItsOwnSlave = "HOST_ITS_OWN_SLAVE"
	//
	ErrorHostCannotSeeSr = "HOST_CANNOT_SEE_SR"
	//
	ErrorHostsFailedToDisableCaching = "HOSTS_FAILED_TO_DISABLE_CACHING"
	//
	ErrorHostsFailedToEnableCaching = "HOSTS_FAILED_TO_ENABLE_CACHING"
	//
	ErrorHostHasResidentVms = "HOST_HAS_RESIDENT_VMS"
	//
	ErrorHostNameInvalid = "HOST_NAME_INVALID"
	//
	ErrorHostIsSlave = "HOST_IS_SLAVE"
	//
	ErrorHostCannotDestroySelf = "HOST_CANNOT_DESTROY_SELF"
	//
	ErrorHostOffline = "HOST_OFFLINE"
	//
	ErrorNoHostsAvailable = "NO_HOSTS_AVAILABLE"
	//
	ErrorHostNotEnoughPcpus = "HOST_NOT_ENOUGH_PCPUS"
	//
	ErrorHostNotEnoughFreeMemory = "HOST_NOT_ENOUGH_FREE_MEMORY"
	//
	ErrorHostPowerOnModeDisabled = "HOST_POWER_ON_MODE_DISABLED"
	//
	ErrorHostIsLive = "HOST_IS_LIVE"
	//
	ErrorHostNotLive = "HOST_NOT_LIVE"
	//
	ErrorHostNotDisabled = "HOST_NOT_DISABLED"
	//
	ErrorHostDisabledIndefinitely = "HOST_DISABLED_INDEFINITELY"
	//
	ErrorHostDisabledUntilReboot = "HOST_DISABLED_UNTIL_REBOOT"
	//
	ErrorHostDisabled = "HOST_DISABLED"
	//
	ErrorHostCannotReadMetrics = "HOST_CANNOT_READ_METRICS"
	//
	ErrorHostInEmergencyMode = "HOST_IN_EMERGENCY_MODE"
	//
	ErrorHostInUse = "HOST_IN_USE"
	//
	ErrorVMIsUsingNestedVirt = "VM_IS_USING_NESTED_VIRT"
	//
	ErrorVMIsImmobile = "VM_IS_IMMOBILE"
	//
	ErrorVMIsProtected = "VM_IS_PROTECTED"
	//
	ErrorVMTooManyVcpus = "VM_TOO_MANY_VCPUS"
	//
	ErrorVMNoVcpus = "VM_NO_VCPUS"
	//
	ErrorVMHvmRequired = "VM_HVM_REQUIRED"
	//
	ErrorUUIDInvalid = "UUID_INVALID"
	//
	ErrorHandleInvalid = "HANDLE_INVALID"
	//
	ErrorNotSupportedDuringUpgrade = "NOT_SUPPORTED_DURING_UPGRADE"
	//
	ErrorTLSConnectionFailed = "TLS_CONNECTION_FAILED"
	//
	ErrorCannotContactHost = "CANNOT_CONTACT_HOST"
	//
	ErrorUserIsNotLocalSuperuser = "USER_IS_NOT_LOCAL_SUPERUSER"
	//
	ErrorChangePasswordRejected = "CHANGE_PASSWORD_REJECTED"
	//
	ErrorSessionInvalid = "SESSION_INVALID"
	//
	ErrorSessionAuthorizationFailed = "SESSION_AUTHORIZATION_FAILED"
	//
	ErrorSessionAuthenticationFailed = "SESSION_AUTHENTICATION_FAILED"
	//
	ErrorFieldTypeError = "FIELD_TYPE_ERROR"
	//
	ErrorMemoryConstraintViolationMaxpin = "MEMORY_CONSTRAINT_VIOLATION_MAXPIN"
	//
	ErrorMemoryConstraintViolationOrder = "MEMORY_CONSTRAINT_VIOLATION_ORDER"
	//
	ErrorMemoryConstraintViolation = "MEMORY_CONSTRAINT_VIOLATION"
	//
	ErrorInvalidValue = "INVALID_VALUE"
	//
	ErrorValueNotSupported = "VALUE_NOT_SUPPORTED"
	//
	ErrorMessageParameterCountMismatch = "MESSAGE_PARAMETER_COUNT_MISMATCH"
	//
	ErrorMessageMethodUnknown = "MESSAGE_METHOD_UNKNOWN"
	//
	ErrorLocationNotUnique = "LOCATION_NOT_UNIQUE"
	//
	ErrorDBUniquenessConstraintViolation = "DB_UNIQUENESS_CONSTRAINT_VIOLATION"
	//
	ErrorMapDuplicateKey = "MAP_DUPLICATE_KEY"
	//
	ErrorInternalError = "INTERNAL_ERROR"
	//
	ErrorPermissionDenied = "PERMISSION_DENIED"
	//
	ErrorMessageRemoved = "MESSAGE_REMOVED"
	//
	ErrorMessageDeprecated = "MESSAGE_DEPRECATED"
)
