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
)

type SRRecord struct {
	// Unique identifier/object reference
	UUID string `json:"uuid,omitempty"`
	// a human-readable name
	NameLabel string `json:"name_label,omitempty"`
	// a notes field containing human-readable description
	NameDescription string `json:"name_description,omitempty"`
	// list of the operations allowed in this state. This list is advisory only and the server state may have changed by the time this field is read by a client.
	AllowedOperations []StorageOperations `json:"allowed_operations,omitempty"`
	// links each of the running tasks using this object (by reference) to a current_operation enum which describes the nature of the task.
	CurrentOperations map[string]StorageOperations `json:"current_operations,omitempty"`
	// all virtual disks known to this storage repository
	VDIs []VDIRef `json:"VDIs,omitempty"`
	// describes how particular hosts can see this storage repository
	PBDs []PBDRef `json:"PBDs,omitempty"`
	// sum of virtual_sizes of all VDIs in this storage repository (in bytes)
	VirtualAllocation int `json:"virtual_allocation,omitempty"`
	// physical space currently utilised on this storage repository (in bytes). Note that for sparse disk formats, physical_utilisation may be less than virtual_allocation
	PhysicalUtilisation int `json:"physical_utilisation,omitempty"`
	// total physical size of the repository (in bytes)
	PhysicalSize int `json:"physical_size,omitempty"`
	// type of the storage repository
	Type string `json:"type,omitempty"`
	// the type of the SR&apos;s content, if required (for example, ISOs)
	ContentType string `json:"content_type,omitempty"`
	// true if this SR is (capable of being) shared between multiple hosts
	Shared bool `json:"shared,omitempty"`
	// additional configuration
	OtherConfig map[string]string `json:"other_config,omitempty"`
	// user-specified tags for categorization purposes
	Tags []string `json:"tags,omitempty"`
	// SM dependent data
	SmConfig map[string]string `json:"sm_config,omitempty"`
	// Binary blobs associated with this SR
	Blobs map[string]BlobRef `json:"blobs,omitempty"`
	// True if this SR is assigned to be the local cache for its host
	LocalCacheEnabled bool `json:"local_cache_enabled,omitempty"`
	// The disaster recovery task which introduced this SR
	IntroducedBy DRTaskRef `json:"introduced_by,omitempty"`
	// True if the SR is using aggregated local storage
	Clustered bool `json:"clustered,omitempty"`
	// True if this is the SR that contains the Tools ISO VDIs
	IsToolsSr bool `json:"is_tools_sr,omitempty"`
}

type SRRef string

// A storage repository
type sr struct{}

var SR sr

// GetAllRecords: Return a map of SR references to SR records for all SRs known to the system.
// Version: rio
func (sr) GetAllRecords(session *Session) (retval map[SRRef]SRRecord, err error) {
	method := "SR.get_all_records"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializeSRRefToSRRecordMap(method+" -> ", result)
	return
}

// GetAllRecords1: Return a map of SR references to SR records for all SRs known to the system.
// Version: rio
func (sr) GetAllRecords1(session *Session) (retval map[SRRef]SRRecord, err error) {
	method := "SR.get_all_records"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializeSRRefToSRRecordMap(method+" -> ", result)
	return
}

// GetAll: Return a list of all the SRs known to the system.
// Version: rio
func (sr) GetAll(session *Session) (retval []SRRef, err error) {
	method := "SR.get_all"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializeSRRefSet(method+" -> ", result)
	return
}

// GetAll1: Return a list of all the SRs known to the system.
// Version: rio
func (sr) GetAll1(session *Session) (retval []SRRef, err error) {
	method := "SR.get_all"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializeSRRefSet(method+" -> ", result)
	return
}

// ForgetDataSourceArchives: Forget the recorded statistics related to the specified data source
// Version: dundee
func (sr) ForgetDataSourceArchives(session *Session, sr SRRef, dataSource string) (err error) {
	method := "SR.forget_data_source_archives"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	srArg, err := serializeSRRef(fmt.Sprintf("%s(%s)", method, "sr"), sr)
	if err != nil {
		return
	}
	dataSourceArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "data_source"), dataSource)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, srArg, dataSourceArg)
	return
}

// ForgetDataSourceArchives3: Forget the recorded statistics related to the specified data source
// Version: dundee
func (sr) ForgetDataSourceArchives3(session *Session, sr SRRef, dataSource string) (err error) {
	method := "SR.forget_data_source_archives"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	srArg, err := serializeSRRef(fmt.Sprintf("%s(%s)", method, "sr"), sr)
	if err != nil {
		return
	}
	dataSourceArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "data_source"), dataSource)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, srArg, dataSourceArg)
	return
}

// QueryDataSource: Query the latest value of the specified data source
// Version: dundee
func (sr) QueryDataSource(session *Session, sr SRRef, dataSource string) (retval float64, err error) {
	method := "SR.query_data_source"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	srArg, err := serializeSRRef(fmt.Sprintf("%s(%s)", method, "sr"), sr)
	if err != nil {
		return
	}
	dataSourceArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "data_source"), dataSource)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, srArg, dataSourceArg)
	if err != nil {
		return
	}
	retval, err = deserializeFloat(method+" -> ", result)
	return
}

// QueryDataSource3: Query the latest value of the specified data source
// Version: dundee
func (sr) QueryDataSource3(session *Session, sr SRRef, dataSource string) (retval float64, err error) {
	method := "SR.query_data_source"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	srArg, err := serializeSRRef(fmt.Sprintf("%s(%s)", method, "sr"), sr)
	if err != nil {
		return
	}
	dataSourceArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "data_source"), dataSource)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, srArg, dataSourceArg)
	if err != nil {
		return
	}
	retval, err = deserializeFloat(method+" -> ", result)
	return
}

// RecordDataSource: Start recording the specified data source
// Version: dundee
func (sr) RecordDataSource(session *Session, sr SRRef, dataSource string) (err error) {
	method := "SR.record_data_source"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	srArg, err := serializeSRRef(fmt.Sprintf("%s(%s)", method, "sr"), sr)
	if err != nil {
		return
	}
	dataSourceArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "data_source"), dataSource)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, srArg, dataSourceArg)
	return
}

// RecordDataSource3: Start recording the specified data source
// Version: dundee
func (sr) RecordDataSource3(session *Session, sr SRRef, dataSource string) (err error) {
	method := "SR.record_data_source"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	srArg, err := serializeSRRef(fmt.Sprintf("%s(%s)", method, "sr"), sr)
	if err != nil {
		return
	}
	dataSourceArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "data_source"), dataSource)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, srArg, dataSourceArg)
	return
}

// GetDataSources:
// Version: dundee
func (sr) GetDataSources(session *Session, sr SRRef) (retval []DataSourceRecord, err error) {
	method := "SR.get_data_sources"
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
	retval, err = deserializeDataSourceRecordSet(method+" -> ", result)
	return
}

// GetDataSources2:
// Version: dundee
func (sr) GetDataSources2(session *Session, sr SRRef) (retval []DataSourceRecord, err error) {
	method := "SR.get_data_sources"
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
	retval, err = deserializeDataSourceRecordSet(method+" -> ", result)
	return
}

// DisableDatabaseReplication:
// Version: boston
func (sr) DisableDatabaseReplication(session *Session, sr SRRef) (err error) {
	method := "SR.disable_database_replication"
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

// AsyncDisableDatabaseReplication:
// Version: boston
func (sr) AsyncDisableDatabaseReplication(session *Session, sr SRRef) (retval TaskRef, err error) {
	method := "Async.SR.disable_database_replication"
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

// DisableDatabaseReplication2:
// Version: boston
func (sr) DisableDatabaseReplication2(session *Session, sr SRRef) (err error) {
	method := "SR.disable_database_replication"
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

// AsyncDisableDatabaseReplication2:
// Version: boston
func (sr) AsyncDisableDatabaseReplication2(session *Session, sr SRRef) (retval TaskRef, err error) {
	method := "Async.SR.disable_database_replication"
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

// EnableDatabaseReplication:
// Version: boston
func (sr) EnableDatabaseReplication(session *Session, sr SRRef) (err error) {
	method := "SR.enable_database_replication"
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

// AsyncEnableDatabaseReplication:
// Version: boston
func (sr) AsyncEnableDatabaseReplication(session *Session, sr SRRef) (retval TaskRef, err error) {
	method := "Async.SR.enable_database_replication"
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

// EnableDatabaseReplication2:
// Version: boston
func (sr) EnableDatabaseReplication2(session *Session, sr SRRef) (err error) {
	method := "SR.enable_database_replication"
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

// AsyncEnableDatabaseReplication2:
// Version: boston
func (sr) AsyncEnableDatabaseReplication2(session *Session, sr SRRef) (retval TaskRef, err error) {
	method := "Async.SR.enable_database_replication"
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

// AssertSupportsDatabaseReplication: Returns successfully if the given SR supports database replication. Otherwise returns an error to explain why not.
// Version: boston
func (sr) AssertSupportsDatabaseReplication(session *Session, sr SRRef) (err error) {
	method := "SR.assert_supports_database_replication"
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

// AsyncAssertSupportsDatabaseReplication: Returns successfully if the given SR supports database replication. Otherwise returns an error to explain why not.
// Version: boston
func (sr) AsyncAssertSupportsDatabaseReplication(session *Session, sr SRRef) (retval TaskRef, err error) {
	method := "Async.SR.assert_supports_database_replication"
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

// AssertSupportsDatabaseReplication2: Returns successfully if the given SR supports database replication. Otherwise returns an error to explain why not.
// Version: boston
func (sr) AssertSupportsDatabaseReplication2(session *Session, sr SRRef) (err error) {
	method := "SR.assert_supports_database_replication"
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

// AsyncAssertSupportsDatabaseReplication2: Returns successfully if the given SR supports database replication. Otherwise returns an error to explain why not.
// Version: boston
func (sr) AsyncAssertSupportsDatabaseReplication2(session *Session, sr SRRef) (retval TaskRef, err error) {
	method := "Async.SR.assert_supports_database_replication"
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

// AssertCanHostHaStatefile: Returns successfully if the given SR can host an HA statefile. Otherwise returns an error to explain why not
// Version: orlando
func (sr) AssertCanHostHaStatefile(session *Session, sr SRRef) (err error) {
	method := "SR.assert_can_host_ha_statefile"
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

// AsyncAssertCanHostHaStatefile: Returns successfully if the given SR can host an HA statefile. Otherwise returns an error to explain why not
// Version: orlando
func (sr) AsyncAssertCanHostHaStatefile(session *Session, sr SRRef) (retval TaskRef, err error) {
	method := "Async.SR.assert_can_host_ha_statefile"
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

// AssertCanHostHaStatefile2: Returns successfully if the given SR can host an HA statefile. Otherwise returns an error to explain why not
// Version: orlando
func (sr) AssertCanHostHaStatefile2(session *Session, sr SRRef) (err error) {
	method := "SR.assert_can_host_ha_statefile"
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

// AsyncAssertCanHostHaStatefile2: Returns successfully if the given SR can host an HA statefile. Otherwise returns an error to explain why not
// Version: orlando
func (sr) AsyncAssertCanHostHaStatefile2(session *Session, sr SRRef) (retval TaskRef, err error) {
	method := "Async.SR.assert_can_host_ha_statefile"
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

// SetPhysicalSize: Sets the SR&apos;s physical_size field
// Version: miami
func (sr) SetPhysicalSize(session *Session, self SRRef, value int) (err error) {
	method := "SR.set_physical_size"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeSRRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetPhysicalSize3: Sets the SR&apos;s physical_size field
// Version: miami
func (sr) SetPhysicalSize3(session *Session, self SRRef, value int) (err error) {
	method := "SR.set_physical_size"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeSRRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// CreateNewBlob: Create a placeholder for a named binary blob of data that is associated with this SR
// Version: tampa
func (sr) CreateNewBlob(session *Session, sr SRRef, name string, mimeType string, public bool) (retval BlobRef, err error) {
	method := "SR.create_new_blob"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	srArg, err := serializeSRRef(fmt.Sprintf("%s(%s)", method, "sr"), sr)
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
	result, err := session.client.sendCall(method, sessionIDArg, srArg, nameArg, mimeTypeArg, publicArg)
	if err != nil {
		return
	}
	retval, err = deserializeBlobRef(method+" -> ", result)
	return
}

// AsyncCreateNewBlob: Create a placeholder for a named binary blob of data that is associated with this SR
// Version: tampa
func (sr) AsyncCreateNewBlob(session *Session, sr SRRef, name string, mimeType string, public bool) (retval TaskRef, err error) {
	method := "Async.SR.create_new_blob"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	srArg, err := serializeSRRef(fmt.Sprintf("%s(%s)", method, "sr"), sr)
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
	result, err := session.client.sendCall(method, sessionIDArg, srArg, nameArg, mimeTypeArg, publicArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// CreateNewBlob5: Create a placeholder for a named binary blob of data that is associated with this SR
// Version: tampa
func (sr) CreateNewBlob5(session *Session, sr SRRef, name string, mimeType string, public bool) (retval BlobRef, err error) {
	method := "SR.create_new_blob"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	srArg, err := serializeSRRef(fmt.Sprintf("%s(%s)", method, "sr"), sr)
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
	result, err := session.client.sendCall(method, sessionIDArg, srArg, nameArg, mimeTypeArg, publicArg)
	if err != nil {
		return
	}
	retval, err = deserializeBlobRef(method+" -> ", result)
	return
}

// AsyncCreateNewBlob5: Create a placeholder for a named binary blob of data that is associated with this SR
// Version: tampa
func (sr) AsyncCreateNewBlob5(session *Session, sr SRRef, name string, mimeType string, public bool) (retval TaskRef, err error) {
	method := "Async.SR.create_new_blob"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	srArg, err := serializeSRRef(fmt.Sprintf("%s(%s)", method, "sr"), sr)
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
	result, err := session.client.sendCall(method, sessionIDArg, srArg, nameArg, mimeTypeArg, publicArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// CreateNewBlob4: Create a placeholder for a named binary blob of data that is associated with this SR
// Version: orlando
func (sr) CreateNewBlob4(session *Session, sr SRRef, name string, mimeType string) (retval BlobRef, err error) {
	method := "SR.create_new_blob"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	srArg, err := serializeSRRef(fmt.Sprintf("%s(%s)", method, "sr"), sr)
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
	result, err := session.client.sendCall(method, sessionIDArg, srArg, nameArg, mimeTypeArg)
	if err != nil {
		return
	}
	retval, err = deserializeBlobRef(method+" -> ", result)
	return
}

// AsyncCreateNewBlob4: Create a placeholder for a named binary blob of data that is associated with this SR
// Version: orlando
func (sr) AsyncCreateNewBlob4(session *Session, sr SRRef, name string, mimeType string) (retval TaskRef, err error) {
	method := "Async.SR.create_new_blob"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	srArg, err := serializeSRRef(fmt.Sprintf("%s(%s)", method, "sr"), sr)
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
	result, err := session.client.sendCall(method, sessionIDArg, srArg, nameArg, mimeTypeArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// SetNameDescription: Set the name description of the SR
// Version: rio
func (sr) SetNameDescription(session *Session, sr SRRef, value string) (err error) {
	method := "SR.set_name_description"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	srArg, err := serializeSRRef(fmt.Sprintf("%s(%s)", method, "sr"), sr)
	if err != nil {
		return
	}
	valueArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, srArg, valueArg)
	return
}

// AsyncSetNameDescription: Set the name description of the SR
// Version: rio
func (sr) AsyncSetNameDescription(session *Session, sr SRRef, value string) (retval TaskRef, err error) {
	method := "Async.SR.set_name_description"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	srArg, err := serializeSRRef(fmt.Sprintf("%s(%s)", method, "sr"), sr)
	if err != nil {
		return
	}
	valueArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, srArg, valueArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// SetNameDescription3: Set the name description of the SR
// Version: rio
func (sr) SetNameDescription3(session *Session, sr SRRef, value string) (err error) {
	method := "SR.set_name_description"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	srArg, err := serializeSRRef(fmt.Sprintf("%s(%s)", method, "sr"), sr)
	if err != nil {
		return
	}
	valueArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, srArg, valueArg)
	return
}

// AsyncSetNameDescription3: Set the name description of the SR
// Version: rio
func (sr) AsyncSetNameDescription3(session *Session, sr SRRef, value string) (retval TaskRef, err error) {
	method := "Async.SR.set_name_description"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	srArg, err := serializeSRRef(fmt.Sprintf("%s(%s)", method, "sr"), sr)
	if err != nil {
		return
	}
	valueArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, srArg, valueArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// SetNameLabel: Set the name label of the SR
// Version: rio
func (sr) SetNameLabel(session *Session, sr SRRef, value string) (err error) {
	method := "SR.set_name_label"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	srArg, err := serializeSRRef(fmt.Sprintf("%s(%s)", method, "sr"), sr)
	if err != nil {
		return
	}
	valueArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, srArg, valueArg)
	return
}

// AsyncSetNameLabel: Set the name label of the SR
// Version: rio
func (sr) AsyncSetNameLabel(session *Session, sr SRRef, value string) (retval TaskRef, err error) {
	method := "Async.SR.set_name_label"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	srArg, err := serializeSRRef(fmt.Sprintf("%s(%s)", method, "sr"), sr)
	if err != nil {
		return
	}
	valueArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, srArg, valueArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// SetNameLabel3: Set the name label of the SR
// Version: rio
func (sr) SetNameLabel3(session *Session, sr SRRef, value string) (err error) {
	method := "SR.set_name_label"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	srArg, err := serializeSRRef(fmt.Sprintf("%s(%s)", method, "sr"), sr)
	if err != nil {
		return
	}
	valueArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, srArg, valueArg)
	return
}

// AsyncSetNameLabel3: Set the name label of the SR
// Version: rio
func (sr) AsyncSetNameLabel3(session *Session, sr SRRef, value string) (retval TaskRef, err error) {
	method := "Async.SR.set_name_label"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	srArg, err := serializeSRRef(fmt.Sprintf("%s(%s)", method, "sr"), sr)
	if err != nil {
		return
	}
	valueArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, srArg, valueArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// SetShared: Sets the shared flag on the SR
// Version: rio
func (sr) SetShared(session *Session, sr SRRef, value bool) (err error) {
	method := "SR.set_shared"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	srArg, err := serializeSRRef(fmt.Sprintf("%s(%s)", method, "sr"), sr)
	if err != nil {
		return
	}
	valueArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, srArg, valueArg)
	return
}

// AsyncSetShared: Sets the shared flag on the SR
// Version: rio
func (sr) AsyncSetShared(session *Session, sr SRRef, value bool) (retval TaskRef, err error) {
	method := "Async.SR.set_shared"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	srArg, err := serializeSRRef(fmt.Sprintf("%s(%s)", method, "sr"), sr)
	if err != nil {
		return
	}
	valueArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, srArg, valueArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// SetShared3: Sets the shared flag on the SR
// Version: rio
func (sr) SetShared3(session *Session, sr SRRef, value bool) (err error) {
	method := "SR.set_shared"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	srArg, err := serializeSRRef(fmt.Sprintf("%s(%s)", method, "sr"), sr)
	if err != nil {
		return
	}
	valueArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, srArg, valueArg)
	return
}

// AsyncSetShared3: Sets the shared flag on the SR
// Version: rio
func (sr) AsyncSetShared3(session *Session, sr SRRef, value bool) (retval TaskRef, err error) {
	method := "Async.SR.set_shared"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	srArg, err := serializeSRRef(fmt.Sprintf("%s(%s)", method, "sr"), sr)
	if err != nil {
		return
	}
	valueArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, srArg, valueArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// ProbeExt: Perform a backend-specific scan, using the given device_config.  If the device_config is complete, then this will return a list of the SRs present of this type on the device, if any.  If the device_config is partial, then a backend-specific scan will be performed, returning results that will guide the user in improving the device_config.
// Version: kolkata
func (sr) ProbeExt(session *Session, host HostRef, deviceConfig map[string]string, typeKey string, smConfig map[string]string) (retval []ProbeResultRecord, err error) {
	method := "SR.probe_ext"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	deviceConfigArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "device_config"), deviceConfig)
	if err != nil {
		return
	}
	typeKeyArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "type"), typeKey)
	if err != nil {
		return
	}
	smConfigArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "sm_config"), smConfig)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, hostArg, deviceConfigArg, typeKeyArg, smConfigArg)
	if err != nil {
		return
	}
	retval, err = deserializeProbeResultRecordSet(method+" -> ", result)
	return
}

// AsyncProbeExt: Perform a backend-specific scan, using the given device_config.  If the device_config is complete, then this will return a list of the SRs present of this type on the device, if any.  If the device_config is partial, then a backend-specific scan will be performed, returning results that will guide the user in improving the device_config.
// Version: kolkata
func (sr) AsyncProbeExt(session *Session, host HostRef, deviceConfig map[string]string, typeKey string, smConfig map[string]string) (retval TaskRef, err error) {
	method := "Async.SR.probe_ext"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	deviceConfigArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "device_config"), deviceConfig)
	if err != nil {
		return
	}
	typeKeyArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "type"), typeKey)
	if err != nil {
		return
	}
	smConfigArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "sm_config"), smConfig)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, hostArg, deviceConfigArg, typeKeyArg, smConfigArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// ProbeExt5: Perform a backend-specific scan, using the given device_config.  If the device_config is complete, then this will return a list of the SRs present of this type on the device, if any.  If the device_config is partial, then a backend-specific scan will be performed, returning results that will guide the user in improving the device_config.
// Version: kolkata
func (sr) ProbeExt5(session *Session, host HostRef, deviceConfig map[string]string, typeKey string, smConfig map[string]string) (retval []ProbeResultRecord, err error) {
	method := "SR.probe_ext"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	deviceConfigArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "device_config"), deviceConfig)
	if err != nil {
		return
	}
	typeKeyArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "type"), typeKey)
	if err != nil {
		return
	}
	smConfigArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "sm_config"), smConfig)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, hostArg, deviceConfigArg, typeKeyArg, smConfigArg)
	if err != nil {
		return
	}
	retval, err = deserializeProbeResultRecordSet(method+" -> ", result)
	return
}

// AsyncProbeExt5: Perform a backend-specific scan, using the given device_config.  If the device_config is complete, then this will return a list of the SRs present of this type on the device, if any.  If the device_config is partial, then a backend-specific scan will be performed, returning results that will guide the user in improving the device_config.
// Version: kolkata
func (sr) AsyncProbeExt5(session *Session, host HostRef, deviceConfig map[string]string, typeKey string, smConfig map[string]string) (retval TaskRef, err error) {
	method := "Async.SR.probe_ext"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	deviceConfigArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "device_config"), deviceConfig)
	if err != nil {
		return
	}
	typeKeyArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "type"), typeKey)
	if err != nil {
		return
	}
	smConfigArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "sm_config"), smConfig)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, hostArg, deviceConfigArg, typeKeyArg, smConfigArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// Probe: Perform a backend-specific scan, using the given device_config.  If the device_config is complete, then this will return a list of the SRs present of this type on the device, if any.  If the device_config is partial, then a backend-specific scan will be performed, returning results that will guide the user in improving the device_config.
// Version: miami
func (sr) Probe(session *Session, host HostRef, deviceConfig map[string]string, typeKey string, smConfig map[string]string) (retval string, err error) {
	method := "SR.probe"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	deviceConfigArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "device_config"), deviceConfig)
	if err != nil {
		return
	}
	typeKeyArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "type"), typeKey)
	if err != nil {
		return
	}
	smConfigArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "sm_config"), smConfig)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, hostArg, deviceConfigArg, typeKeyArg, smConfigArg)
	if err != nil {
		return
	}
	retval, err = deserializeString(method+" -> ", result)
	return
}

// AsyncProbe: Perform a backend-specific scan, using the given device_config.  If the device_config is complete, then this will return a list of the SRs present of this type on the device, if any.  If the device_config is partial, then a backend-specific scan will be performed, returning results that will guide the user in improving the device_config.
// Version: miami
func (sr) AsyncProbe(session *Session, host HostRef, deviceConfig map[string]string, typeKey string, smConfig map[string]string) (retval TaskRef, err error) {
	method := "Async.SR.probe"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	deviceConfigArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "device_config"), deviceConfig)
	if err != nil {
		return
	}
	typeKeyArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "type"), typeKey)
	if err != nil {
		return
	}
	smConfigArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "sm_config"), smConfig)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, hostArg, deviceConfigArg, typeKeyArg, smConfigArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// Probe5: Perform a backend-specific scan, using the given device_config.  If the device_config is complete, then this will return a list of the SRs present of this type on the device, if any.  If the device_config is partial, then a backend-specific scan will be performed, returning results that will guide the user in improving the device_config.
// Version: miami
func (sr) Probe5(session *Session, host HostRef, deviceConfig map[string]string, typeKey string, smConfig map[string]string) (retval string, err error) {
	method := "SR.probe"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	deviceConfigArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "device_config"), deviceConfig)
	if err != nil {
		return
	}
	typeKeyArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "type"), typeKey)
	if err != nil {
		return
	}
	smConfigArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "sm_config"), smConfig)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, hostArg, deviceConfigArg, typeKeyArg, smConfigArg)
	if err != nil {
		return
	}
	retval, err = deserializeString(method+" -> ", result)
	return
}

// AsyncProbe5: Perform a backend-specific scan, using the given device_config.  If the device_config is complete, then this will return a list of the SRs present of this type on the device, if any.  If the device_config is partial, then a backend-specific scan will be performed, returning results that will guide the user in improving the device_config.
// Version: miami
func (sr) AsyncProbe5(session *Session, host HostRef, deviceConfig map[string]string, typeKey string, smConfig map[string]string) (retval TaskRef, err error) {
	method := "Async.SR.probe"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	deviceConfigArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "device_config"), deviceConfig)
	if err != nil {
		return
	}
	typeKeyArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "type"), typeKey)
	if err != nil {
		return
	}
	smConfigArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "sm_config"), smConfig)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, hostArg, deviceConfigArg, typeKeyArg, smConfigArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// Probe3: Perform a backend-specific scan, using the given device_config.  If the device_config is complete, then this will return a list of the SRs present of this type on the device, if any.  If the device_config is partial, then a backend-specific scan will be performed, returning results that will guide the user in improving the device_config.
// Version: rio
func (sr) Probe3(session *Session, host HostRef, deviceConfig map[string]string) (retval string, err error) {
	method := "SR.probe"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	deviceConfigArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "device_config"), deviceConfig)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, hostArg, deviceConfigArg)
	if err != nil {
		return
	}
	retval, err = deserializeString(method+" -> ", result)
	return
}

// AsyncProbe3: Perform a backend-specific scan, using the given device_config.  If the device_config is complete, then this will return a list of the SRs present of this type on the device, if any.  If the device_config is partial, then a backend-specific scan will be performed, returning results that will guide the user in improving the device_config.
// Version: rio
func (sr) AsyncProbe3(session *Session, host HostRef, deviceConfig map[string]string) (retval TaskRef, err error) {
	method := "Async.SR.probe"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	deviceConfigArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "device_config"), deviceConfig)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, hostArg, deviceConfigArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// Scan: Refreshes the list of VDIs associated with an SR
// Version: rio
func (sr) Scan(session *Session, sr SRRef) (err error) {
	method := "SR.scan"
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

// AsyncScan: Refreshes the list of VDIs associated with an SR
// Version: rio
func (sr) AsyncScan(session *Session, sr SRRef) (retval TaskRef, err error) {
	method := "Async.SR.scan"
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

// Scan2: Refreshes the list of VDIs associated with an SR
// Version: rio
func (sr) Scan2(session *Session, sr SRRef) (err error) {
	method := "SR.scan"
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

// AsyncScan2: Refreshes the list of VDIs associated with an SR
// Version: rio
func (sr) AsyncScan2(session *Session, sr SRRef) (retval TaskRef, err error) {
	method := "Async.SR.scan"
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

// GetSupportedTypes: Return a set of all the SR types supported by the system
// Version: rio
func (sr) GetSupportedTypes(session *Session) (retval []string, err error) {
	method := "SR.get_supported_types"
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

// GetSupportedTypes1: Return a set of all the SR types supported by the system
// Version: rio
func (sr) GetSupportedTypes1(session *Session) (retval []string, err error) {
	method := "SR.get_supported_types"
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

// Update: Refresh the fields on the SR object
// Version: symc
func (sr) Update(session *Session, sr SRRef) (err error) {
	method := "SR.update"
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

// AsyncUpdate: Refresh the fields on the SR object
// Version: symc
func (sr) AsyncUpdate(session *Session, sr SRRef) (retval TaskRef, err error) {
	method := "Async.SR.update"
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

// Update2: Refresh the fields on the SR object
// Version: symc
func (sr) Update2(session *Session, sr SRRef) (err error) {
	method := "SR.update"
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

// AsyncUpdate2: Refresh the fields on the SR object
// Version: symc
func (sr) AsyncUpdate2(session *Session, sr SRRef) (retval TaskRef, err error) {
	method := "Async.SR.update"
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

// Forget: Removing specified SR-record from database, without attempting to remove SR from disk
// Version: rio
//
// Errors:
// SR_HAS_PBD - The SR is still connected to a host via a PBD. It cannot be destroyed or forgotten.
func (sr) Forget(session *Session, sr SRRef) (err error) {
	method := "SR.forget"
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

// AsyncForget: Removing specified SR-record from database, without attempting to remove SR from disk
// Version: rio
//
// Errors:
// SR_HAS_PBD - The SR is still connected to a host via a PBD. It cannot be destroyed or forgotten.
func (sr) AsyncForget(session *Session, sr SRRef) (retval TaskRef, err error) {
	method := "Async.SR.forget"
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

// Forget2: Removing specified SR-record from database, without attempting to remove SR from disk
// Version: rio
//
// Errors:
// SR_HAS_PBD - The SR is still connected to a host via a PBD. It cannot be destroyed or forgotten.
func (sr) Forget2(session *Session, sr SRRef) (err error) {
	method := "SR.forget"
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

// AsyncForget2: Removing specified SR-record from database, without attempting to remove SR from disk
// Version: rio
//
// Errors:
// SR_HAS_PBD - The SR is still connected to a host via a PBD. It cannot be destroyed or forgotten.
func (sr) AsyncForget2(session *Session, sr SRRef) (retval TaskRef, err error) {
	method := "Async.SR.forget"
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

// Destroy: Destroy specified SR, removing SR-record from database and remove SR from disk. (In order to affect this operation the appropriate device_config is read from the specified SR&apos;s PBD on current host)
// Version: rio
//
// Errors:
// SR_HAS_PBD - The SR is still connected to a host via a PBD. It cannot be destroyed or forgotten.
func (sr) Destroy(session *Session, sr SRRef) (err error) {
	method := "SR.destroy"
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

// AsyncDestroy: Destroy specified SR, removing SR-record from database and remove SR from disk. (In order to affect this operation the appropriate device_config is read from the specified SR&apos;s PBD on current host)
// Version: rio
//
// Errors:
// SR_HAS_PBD - The SR is still connected to a host via a PBD. It cannot be destroyed or forgotten.
func (sr) AsyncDestroy(session *Session, sr SRRef) (retval TaskRef, err error) {
	method := "Async.SR.destroy"
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

// Destroy2: Destroy specified SR, removing SR-record from database and remove SR from disk. (In order to affect this operation the appropriate device_config is read from the specified SR&apos;s PBD on current host)
// Version: rio
//
// Errors:
// SR_HAS_PBD - The SR is still connected to a host via a PBD. It cannot be destroyed or forgotten.
func (sr) Destroy2(session *Session, sr SRRef) (err error) {
	method := "SR.destroy"
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

// AsyncDestroy2: Destroy specified SR, removing SR-record from database and remove SR from disk. (In order to affect this operation the appropriate device_config is read from the specified SR&apos;s PBD on current host)
// Version: rio
//
// Errors:
// SR_HAS_PBD - The SR is still connected to a host via a PBD. It cannot be destroyed or forgotten.
func (sr) AsyncDestroy2(session *Session, sr SRRef) (retval TaskRef, err error) {
	method := "Async.SR.destroy"
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

// Make: Create a new Storage Repository on disk. This call is deprecated: use SR.create instead.
// Version: miami
func (sr) Make(session *Session, host HostRef, deviceConfig map[string]string, physicalSize int, nameLabel string, nameDescription string, typeKey string, contentType string, smConfig map[string]string) (retval string, err error) {
	method := "SR.make"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	deviceConfigArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "device_config"), deviceConfig)
	if err != nil {
		return
	}
	physicalSizeArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "physical_size"), physicalSize)
	if err != nil {
		return
	}
	nameLabelArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "name_label"), nameLabel)
	if err != nil {
		return
	}
	nameDescriptionArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "name_description"), nameDescription)
	if err != nil {
		return
	}
	typeKeyArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "type"), typeKey)
	if err != nil {
		return
	}
	contentTypeArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "content_type"), contentType)
	if err != nil {
		return
	}
	smConfigArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "sm_config"), smConfig)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, hostArg, deviceConfigArg, physicalSizeArg, nameLabelArg, nameDescriptionArg, typeKeyArg, contentTypeArg, smConfigArg)
	if err != nil {
		return
	}
	retval, err = deserializeString(method+" -> ", result)
	return
}

// AsyncMake: Create a new Storage Repository on disk. This call is deprecated: use SR.create instead.
// Version: miami
func (sr) AsyncMake(session *Session, host HostRef, deviceConfig map[string]string, physicalSize int, nameLabel string, nameDescription string, typeKey string, contentType string, smConfig map[string]string) (retval TaskRef, err error) {
	method := "Async.SR.make"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	deviceConfigArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "device_config"), deviceConfig)
	if err != nil {
		return
	}
	physicalSizeArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "physical_size"), physicalSize)
	if err != nil {
		return
	}
	nameLabelArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "name_label"), nameLabel)
	if err != nil {
		return
	}
	nameDescriptionArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "name_description"), nameDescription)
	if err != nil {
		return
	}
	typeKeyArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "type"), typeKey)
	if err != nil {
		return
	}
	contentTypeArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "content_type"), contentType)
	if err != nil {
		return
	}
	smConfigArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "sm_config"), smConfig)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, hostArg, deviceConfigArg, physicalSizeArg, nameLabelArg, nameDescriptionArg, typeKeyArg, contentTypeArg, smConfigArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// Make9: Create a new Storage Repository on disk. This call is deprecated: use SR.create instead.
// Version: miami
func (sr) Make9(session *Session, host HostRef, deviceConfig map[string]string, physicalSize int, nameLabel string, nameDescription string, typeKey string, contentType string, smConfig map[string]string) (retval string, err error) {
	method := "SR.make"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	deviceConfigArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "device_config"), deviceConfig)
	if err != nil {
		return
	}
	physicalSizeArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "physical_size"), physicalSize)
	if err != nil {
		return
	}
	nameLabelArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "name_label"), nameLabel)
	if err != nil {
		return
	}
	nameDescriptionArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "name_description"), nameDescription)
	if err != nil {
		return
	}
	typeKeyArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "type"), typeKey)
	if err != nil {
		return
	}
	contentTypeArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "content_type"), contentType)
	if err != nil {
		return
	}
	smConfigArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "sm_config"), smConfig)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, hostArg, deviceConfigArg, physicalSizeArg, nameLabelArg, nameDescriptionArg, typeKeyArg, contentTypeArg, smConfigArg)
	if err != nil {
		return
	}
	retval, err = deserializeString(method+" -> ", result)
	return
}

// AsyncMake9: Create a new Storage Repository on disk. This call is deprecated: use SR.create instead.
// Version: miami
func (sr) AsyncMake9(session *Session, host HostRef, deviceConfig map[string]string, physicalSize int, nameLabel string, nameDescription string, typeKey string, contentType string, smConfig map[string]string) (retval TaskRef, err error) {
	method := "Async.SR.make"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	deviceConfigArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "device_config"), deviceConfig)
	if err != nil {
		return
	}
	physicalSizeArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "physical_size"), physicalSize)
	if err != nil {
		return
	}
	nameLabelArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "name_label"), nameLabel)
	if err != nil {
		return
	}
	nameDescriptionArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "name_description"), nameDescription)
	if err != nil {
		return
	}
	typeKeyArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "type"), typeKey)
	if err != nil {
		return
	}
	contentTypeArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "content_type"), contentType)
	if err != nil {
		return
	}
	smConfigArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "sm_config"), smConfig)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, hostArg, deviceConfigArg, physicalSizeArg, nameLabelArg, nameDescriptionArg, typeKeyArg, contentTypeArg, smConfigArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// Make8: Create a new Storage Repository on disk. This call is deprecated: use SR.create instead.
// Version: rio
func (sr) Make8(session *Session, host HostRef, deviceConfig map[string]string, physicalSize int, nameLabel string, nameDescription string, typeKey string, contentType string) (retval string, err error) {
	method := "SR.make"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	deviceConfigArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "device_config"), deviceConfig)
	if err != nil {
		return
	}
	physicalSizeArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "physical_size"), physicalSize)
	if err != nil {
		return
	}
	nameLabelArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "name_label"), nameLabel)
	if err != nil {
		return
	}
	nameDescriptionArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "name_description"), nameDescription)
	if err != nil {
		return
	}
	typeKeyArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "type"), typeKey)
	if err != nil {
		return
	}
	contentTypeArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "content_type"), contentType)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, hostArg, deviceConfigArg, physicalSizeArg, nameLabelArg, nameDescriptionArg, typeKeyArg, contentTypeArg)
	if err != nil {
		return
	}
	retval, err = deserializeString(method+" -> ", result)
	return
}

// AsyncMake8: Create a new Storage Repository on disk. This call is deprecated: use SR.create instead.
// Version: rio
func (sr) AsyncMake8(session *Session, host HostRef, deviceConfig map[string]string, physicalSize int, nameLabel string, nameDescription string, typeKey string, contentType string) (retval TaskRef, err error) {
	method := "Async.SR.make"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	deviceConfigArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "device_config"), deviceConfig)
	if err != nil {
		return
	}
	physicalSizeArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "physical_size"), physicalSize)
	if err != nil {
		return
	}
	nameLabelArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "name_label"), nameLabel)
	if err != nil {
		return
	}
	nameDescriptionArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "name_description"), nameDescription)
	if err != nil {
		return
	}
	typeKeyArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "type"), typeKey)
	if err != nil {
		return
	}
	contentTypeArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "content_type"), contentType)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, hostArg, deviceConfigArg, physicalSizeArg, nameLabelArg, nameDescriptionArg, typeKeyArg, contentTypeArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// Introduce: Introduce a new Storage Repository into the managed system
// Version: miami
func (sr) Introduce(session *Session, uuid string, nameLabel string, nameDescription string, typeKey string, contentType string, shared bool, smConfig map[string]string) (retval SRRef, err error) {
	method := "SR.introduce"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	uuidArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "uuid"), uuid)
	if err != nil {
		return
	}
	nameLabelArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "name_label"), nameLabel)
	if err != nil {
		return
	}
	nameDescriptionArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "name_description"), nameDescription)
	if err != nil {
		return
	}
	typeKeyArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "type"), typeKey)
	if err != nil {
		return
	}
	contentTypeArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "content_type"), contentType)
	if err != nil {
		return
	}
	sharedArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "shared"), shared)
	if err != nil {
		return
	}
	smConfigArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "sm_config"), smConfig)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, uuidArg, nameLabelArg, nameDescriptionArg, typeKeyArg, contentTypeArg, sharedArg, smConfigArg)
	if err != nil {
		return
	}
	retval, err = deserializeSRRef(method+" -> ", result)
	return
}

// AsyncIntroduce: Introduce a new Storage Repository into the managed system
// Version: miami
func (sr) AsyncIntroduce(session *Session, uuid string, nameLabel string, nameDescription string, typeKey string, contentType string, shared bool, smConfig map[string]string) (retval TaskRef, err error) {
	method := "Async.SR.introduce"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	uuidArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "uuid"), uuid)
	if err != nil {
		return
	}
	nameLabelArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "name_label"), nameLabel)
	if err != nil {
		return
	}
	nameDescriptionArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "name_description"), nameDescription)
	if err != nil {
		return
	}
	typeKeyArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "type"), typeKey)
	if err != nil {
		return
	}
	contentTypeArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "content_type"), contentType)
	if err != nil {
		return
	}
	sharedArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "shared"), shared)
	if err != nil {
		return
	}
	smConfigArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "sm_config"), smConfig)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, uuidArg, nameLabelArg, nameDescriptionArg, typeKeyArg, contentTypeArg, sharedArg, smConfigArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// Introduce8: Introduce a new Storage Repository into the managed system
// Version: miami
func (sr) Introduce8(session *Session, uuid string, nameLabel string, nameDescription string, typeKey string, contentType string, shared bool, smConfig map[string]string) (retval SRRef, err error) {
	method := "SR.introduce"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	uuidArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "uuid"), uuid)
	if err != nil {
		return
	}
	nameLabelArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "name_label"), nameLabel)
	if err != nil {
		return
	}
	nameDescriptionArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "name_description"), nameDescription)
	if err != nil {
		return
	}
	typeKeyArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "type"), typeKey)
	if err != nil {
		return
	}
	contentTypeArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "content_type"), contentType)
	if err != nil {
		return
	}
	sharedArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "shared"), shared)
	if err != nil {
		return
	}
	smConfigArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "sm_config"), smConfig)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, uuidArg, nameLabelArg, nameDescriptionArg, typeKeyArg, contentTypeArg, sharedArg, smConfigArg)
	if err != nil {
		return
	}
	retval, err = deserializeSRRef(method+" -> ", result)
	return
}

// AsyncIntroduce8: Introduce a new Storage Repository into the managed system
// Version: miami
func (sr) AsyncIntroduce8(session *Session, uuid string, nameLabel string, nameDescription string, typeKey string, contentType string, shared bool, smConfig map[string]string) (retval TaskRef, err error) {
	method := "Async.SR.introduce"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	uuidArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "uuid"), uuid)
	if err != nil {
		return
	}
	nameLabelArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "name_label"), nameLabel)
	if err != nil {
		return
	}
	nameDescriptionArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "name_description"), nameDescription)
	if err != nil {
		return
	}
	typeKeyArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "type"), typeKey)
	if err != nil {
		return
	}
	contentTypeArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "content_type"), contentType)
	if err != nil {
		return
	}
	sharedArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "shared"), shared)
	if err != nil {
		return
	}
	smConfigArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "sm_config"), smConfig)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, uuidArg, nameLabelArg, nameDescriptionArg, typeKeyArg, contentTypeArg, sharedArg, smConfigArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// Introduce7: Introduce a new Storage Repository into the managed system
// Version: rio
func (sr) Introduce7(session *Session, uuid string, nameLabel string, nameDescription string, typeKey string, contentType string, shared bool) (retval SRRef, err error) {
	method := "SR.introduce"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	uuidArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "uuid"), uuid)
	if err != nil {
		return
	}
	nameLabelArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "name_label"), nameLabel)
	if err != nil {
		return
	}
	nameDescriptionArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "name_description"), nameDescription)
	if err != nil {
		return
	}
	typeKeyArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "type"), typeKey)
	if err != nil {
		return
	}
	contentTypeArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "content_type"), contentType)
	if err != nil {
		return
	}
	sharedArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "shared"), shared)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, uuidArg, nameLabelArg, nameDescriptionArg, typeKeyArg, contentTypeArg, sharedArg)
	if err != nil {
		return
	}
	retval, err = deserializeSRRef(method+" -> ", result)
	return
}

// AsyncIntroduce7: Introduce a new Storage Repository into the managed system
// Version: rio
func (sr) AsyncIntroduce7(session *Session, uuid string, nameLabel string, nameDescription string, typeKey string, contentType string, shared bool) (retval TaskRef, err error) {
	method := "Async.SR.introduce"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	uuidArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "uuid"), uuid)
	if err != nil {
		return
	}
	nameLabelArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "name_label"), nameLabel)
	if err != nil {
		return
	}
	nameDescriptionArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "name_description"), nameDescription)
	if err != nil {
		return
	}
	typeKeyArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "type"), typeKey)
	if err != nil {
		return
	}
	contentTypeArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "content_type"), contentType)
	if err != nil {
		return
	}
	sharedArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "shared"), shared)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, uuidArg, nameLabelArg, nameDescriptionArg, typeKeyArg, contentTypeArg, sharedArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// Create: Create a new Storage Repository and introduce it into the managed system, creating both SR record and PBD record to attach it to current host (with specified device_config parameters)
// Version: miami
//
// Errors:
// SR_UNKNOWN_DRIVER - The SR could not be connected because the driver was not recognised.
func (sr) Create(session *Session, host HostRef, deviceConfig map[string]string, physicalSize int, nameLabel string, nameDescription string, typeKey string, contentType string, shared bool, smConfig map[string]string) (retval SRRef, err error) {
	method := "SR.create"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	deviceConfigArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "device_config"), deviceConfig)
	if err != nil {
		return
	}
	physicalSizeArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "physical_size"), physicalSize)
	if err != nil {
		return
	}
	nameLabelArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "name_label"), nameLabel)
	if err != nil {
		return
	}
	nameDescriptionArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "name_description"), nameDescription)
	if err != nil {
		return
	}
	typeKeyArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "type"), typeKey)
	if err != nil {
		return
	}
	contentTypeArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "content_type"), contentType)
	if err != nil {
		return
	}
	sharedArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "shared"), shared)
	if err != nil {
		return
	}
	smConfigArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "sm_config"), smConfig)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, hostArg, deviceConfigArg, physicalSizeArg, nameLabelArg, nameDescriptionArg, typeKeyArg, contentTypeArg, sharedArg, smConfigArg)
	if err != nil {
		return
	}
	retval, err = deserializeSRRef(method+" -> ", result)
	return
}

// AsyncCreate: Create a new Storage Repository and introduce it into the managed system, creating both SR record and PBD record to attach it to current host (with specified device_config parameters)
// Version: miami
//
// Errors:
// SR_UNKNOWN_DRIVER - The SR could not be connected because the driver was not recognised.
func (sr) AsyncCreate(session *Session, host HostRef, deviceConfig map[string]string, physicalSize int, nameLabel string, nameDescription string, typeKey string, contentType string, shared bool, smConfig map[string]string) (retval TaskRef, err error) {
	method := "Async.SR.create"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	deviceConfigArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "device_config"), deviceConfig)
	if err != nil {
		return
	}
	physicalSizeArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "physical_size"), physicalSize)
	if err != nil {
		return
	}
	nameLabelArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "name_label"), nameLabel)
	if err != nil {
		return
	}
	nameDescriptionArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "name_description"), nameDescription)
	if err != nil {
		return
	}
	typeKeyArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "type"), typeKey)
	if err != nil {
		return
	}
	contentTypeArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "content_type"), contentType)
	if err != nil {
		return
	}
	sharedArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "shared"), shared)
	if err != nil {
		return
	}
	smConfigArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "sm_config"), smConfig)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, hostArg, deviceConfigArg, physicalSizeArg, nameLabelArg, nameDescriptionArg, typeKeyArg, contentTypeArg, sharedArg, smConfigArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// Create10: Create a new Storage Repository and introduce it into the managed system, creating both SR record and PBD record to attach it to current host (with specified device_config parameters)
// Version: miami
//
// Errors:
// SR_UNKNOWN_DRIVER - The SR could not be connected because the driver was not recognised.
func (sr) Create10(session *Session, host HostRef, deviceConfig map[string]string, physicalSize int, nameLabel string, nameDescription string, typeKey string, contentType string, shared bool, smConfig map[string]string) (retval SRRef, err error) {
	method := "SR.create"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	deviceConfigArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "device_config"), deviceConfig)
	if err != nil {
		return
	}
	physicalSizeArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "physical_size"), physicalSize)
	if err != nil {
		return
	}
	nameLabelArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "name_label"), nameLabel)
	if err != nil {
		return
	}
	nameDescriptionArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "name_description"), nameDescription)
	if err != nil {
		return
	}
	typeKeyArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "type"), typeKey)
	if err != nil {
		return
	}
	contentTypeArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "content_type"), contentType)
	if err != nil {
		return
	}
	sharedArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "shared"), shared)
	if err != nil {
		return
	}
	smConfigArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "sm_config"), smConfig)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, hostArg, deviceConfigArg, physicalSizeArg, nameLabelArg, nameDescriptionArg, typeKeyArg, contentTypeArg, sharedArg, smConfigArg)
	if err != nil {
		return
	}
	retval, err = deserializeSRRef(method+" -> ", result)
	return
}

// AsyncCreate10: Create a new Storage Repository and introduce it into the managed system, creating both SR record and PBD record to attach it to current host (with specified device_config parameters)
// Version: miami
//
// Errors:
// SR_UNKNOWN_DRIVER - The SR could not be connected because the driver was not recognised.
func (sr) AsyncCreate10(session *Session, host HostRef, deviceConfig map[string]string, physicalSize int, nameLabel string, nameDescription string, typeKey string, contentType string, shared bool, smConfig map[string]string) (retval TaskRef, err error) {
	method := "Async.SR.create"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	deviceConfigArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "device_config"), deviceConfig)
	if err != nil {
		return
	}
	physicalSizeArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "physical_size"), physicalSize)
	if err != nil {
		return
	}
	nameLabelArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "name_label"), nameLabel)
	if err != nil {
		return
	}
	nameDescriptionArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "name_description"), nameDescription)
	if err != nil {
		return
	}
	typeKeyArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "type"), typeKey)
	if err != nil {
		return
	}
	contentTypeArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "content_type"), contentType)
	if err != nil {
		return
	}
	sharedArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "shared"), shared)
	if err != nil {
		return
	}
	smConfigArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "sm_config"), smConfig)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, hostArg, deviceConfigArg, physicalSizeArg, nameLabelArg, nameDescriptionArg, typeKeyArg, contentTypeArg, sharedArg, smConfigArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// Create9: Create a new Storage Repository and introduce it into the managed system, creating both SR record and PBD record to attach it to current host (with specified device_config parameters)
// Version: rio
//
// Errors:
// SR_UNKNOWN_DRIVER - The SR could not be connected because the driver was not recognised.
func (sr) Create9(session *Session, host HostRef, deviceConfig map[string]string, physicalSize int, nameLabel string, nameDescription string, typeKey string, contentType string, shared bool) (retval SRRef, err error) {
	method := "SR.create"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	deviceConfigArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "device_config"), deviceConfig)
	if err != nil {
		return
	}
	physicalSizeArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "physical_size"), physicalSize)
	if err != nil {
		return
	}
	nameLabelArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "name_label"), nameLabel)
	if err != nil {
		return
	}
	nameDescriptionArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "name_description"), nameDescription)
	if err != nil {
		return
	}
	typeKeyArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "type"), typeKey)
	if err != nil {
		return
	}
	contentTypeArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "content_type"), contentType)
	if err != nil {
		return
	}
	sharedArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "shared"), shared)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, hostArg, deviceConfigArg, physicalSizeArg, nameLabelArg, nameDescriptionArg, typeKeyArg, contentTypeArg, sharedArg)
	if err != nil {
		return
	}
	retval, err = deserializeSRRef(method+" -> ", result)
	return
}

// AsyncCreate9: Create a new Storage Repository and introduce it into the managed system, creating both SR record and PBD record to attach it to current host (with specified device_config parameters)
// Version: rio
//
// Errors:
// SR_UNKNOWN_DRIVER - The SR could not be connected because the driver was not recognised.
func (sr) AsyncCreate9(session *Session, host HostRef, deviceConfig map[string]string, physicalSize int, nameLabel string, nameDescription string, typeKey string, contentType string, shared bool) (retval TaskRef, err error) {
	method := "Async.SR.create"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	deviceConfigArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "device_config"), deviceConfig)
	if err != nil {
		return
	}
	physicalSizeArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "physical_size"), physicalSize)
	if err != nil {
		return
	}
	nameLabelArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "name_label"), nameLabel)
	if err != nil {
		return
	}
	nameDescriptionArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "name_description"), nameDescription)
	if err != nil {
		return
	}
	typeKeyArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "type"), typeKey)
	if err != nil {
		return
	}
	contentTypeArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "content_type"), contentType)
	if err != nil {
		return
	}
	sharedArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "shared"), shared)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, hostArg, deviceConfigArg, physicalSizeArg, nameLabelArg, nameDescriptionArg, typeKeyArg, contentTypeArg, sharedArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// RemoveFromSmConfig: Remove the given key and its corresponding value from the sm_config field of the given SR.  If the key is not in that Map, then do nothing.
// Version: miami
func (sr) RemoveFromSmConfig(session *Session, self SRRef, key string) (err error) {
	method := "SR.remove_from_sm_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeSRRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// RemoveFromSmConfig3: Remove the given key and its corresponding value from the sm_config field of the given SR.  If the key is not in that Map, then do nothing.
// Version: miami
func (sr) RemoveFromSmConfig3(session *Session, self SRRef, key string) (err error) {
	method := "SR.remove_from_sm_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeSRRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// RemoveFromSmConfig2: Remove the given key and its corresponding value from the sm_config field of the given SR.  If the key is not in that Map, then do nothing.
// Version: rio
func (sr) RemoveFromSmConfig2(session *Session, self SRRef) (err error) {
	method := "SR.remove_from_sm_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeSRRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// AddToSmConfig: Add the given key-value pair to the sm_config field of the given SR.
// Version: miami
func (sr) AddToSmConfig(session *Session, self SRRef, key string, value string) (err error) {
	method := "SR.add_to_sm_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeSRRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// AddToSmConfig4: Add the given key-value pair to the sm_config field of the given SR.
// Version: miami
func (sr) AddToSmConfig4(session *Session, self SRRef, key string, value string) (err error) {
	method := "SR.add_to_sm_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeSRRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// AddToSmConfig2: Add the given key-value pair to the sm_config field of the given SR.
// Version: rio
func (sr) AddToSmConfig2(session *Session, self SRRef) (err error) {
	method := "SR.add_to_sm_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeSRRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// SetSmConfig: Set the sm_config field of the given SR.
// Version: miami
func (sr) SetSmConfig(session *Session, self SRRef, value map[string]string) (err error) {
	method := "SR.set_sm_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeSRRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetSmConfig3: Set the sm_config field of the given SR.
// Version: miami
func (sr) SetSmConfig3(session *Session, self SRRef, value map[string]string) (err error) {
	method := "SR.set_sm_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeSRRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetSmConfig2: Set the sm_config field of the given SR.
// Version: rio
func (sr) SetSmConfig2(session *Session, self SRRef) (err error) {
	method := "SR.set_sm_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeSRRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// RemoveTags: Remove the given value from the tags field of the given SR.  If the value is not in that Set, then do nothing.
// Version: orlando
func (sr) RemoveTags(session *Session, self SRRef, value string) (err error) {
	method := "SR.remove_tags"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeSRRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// RemoveTags3: Remove the given value from the tags field of the given SR.  If the value is not in that Set, then do nothing.
// Version: orlando
func (sr) RemoveTags3(session *Session, self SRRef, value string) (err error) {
	method := "SR.remove_tags"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeSRRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// RemoveTags2: Remove the given value from the tags field of the given SR.  If the value is not in that Set, then do nothing.
// Version: rio
func (sr) RemoveTags2(session *Session, self SRRef) (err error) {
	method := "SR.remove_tags"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeSRRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// AddTags: Add the given value to the tags field of the given SR.  If the value is already in that Set, then do nothing.
// Version: orlando
func (sr) AddTags(session *Session, self SRRef, value string) (err error) {
	method := "SR.add_tags"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeSRRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// AddTags3: Add the given value to the tags field of the given SR.  If the value is already in that Set, then do nothing.
// Version: orlando
func (sr) AddTags3(session *Session, self SRRef, value string) (err error) {
	method := "SR.add_tags"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeSRRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// AddTags2: Add the given value to the tags field of the given SR.  If the value is already in that Set, then do nothing.
// Version: rio
func (sr) AddTags2(session *Session, self SRRef) (err error) {
	method := "SR.add_tags"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeSRRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// SetTags: Set the tags field of the given SR.
// Version: orlando
func (sr) SetTags(session *Session, self SRRef, value []string) (err error) {
	method := "SR.set_tags"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeSRRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetTags3: Set the tags field of the given SR.
// Version: orlando
func (sr) SetTags3(session *Session, self SRRef, value []string) (err error) {
	method := "SR.set_tags"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeSRRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetTags2: Set the tags field of the given SR.
// Version: rio
func (sr) SetTags2(session *Session, self SRRef) (err error) {
	method := "SR.set_tags"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeSRRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// RemoveFromOtherConfig: Remove the given key and its corresponding value from the other_config field of the given SR.  If the key is not in that Map, then do nothing.
// Version: rio
func (sr) RemoveFromOtherConfig(session *Session, self SRRef, key string) (err error) {
	method := "SR.remove_from_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeSRRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// RemoveFromOtherConfig3: Remove the given key and its corresponding value from the other_config field of the given SR.  If the key is not in that Map, then do nothing.
// Version: rio
func (sr) RemoveFromOtherConfig3(session *Session, self SRRef, key string) (err error) {
	method := "SR.remove_from_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeSRRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// AddToOtherConfig: Add the given key-value pair to the other_config field of the given SR.
// Version: rio
func (sr) AddToOtherConfig(session *Session, self SRRef, key string, value string) (err error) {
	method := "SR.add_to_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeSRRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// AddToOtherConfig4: Add the given key-value pair to the other_config field of the given SR.
// Version: rio
func (sr) AddToOtherConfig4(session *Session, self SRRef, key string, value string) (err error) {
	method := "SR.add_to_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeSRRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetOtherConfig: Set the other_config field of the given SR.
// Version: rio
func (sr) SetOtherConfig(session *Session, self SRRef, value map[string]string) (err error) {
	method := "SR.set_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeSRRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetOtherConfig3: Set the other_config field of the given SR.
// Version: rio
func (sr) SetOtherConfig3(session *Session, self SRRef, value map[string]string) (err error) {
	method := "SR.set_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeSRRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetIsToolsSr: Get the is_tools_sr field of the given SR.
// Version: rio
func (sr) GetIsToolsSr(session *Session, self SRRef) (retval bool, err error) {
	method := "SR.get_is_tools_sr"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeSRRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetIsToolsSr2: Get the is_tools_sr field of the given SR.
// Version: rio
func (sr) GetIsToolsSr2(session *Session, self SRRef) (retval bool, err error) {
	method := "SR.get_is_tools_sr"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeSRRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetClustered: Get the clustered field of the given SR.
// Version: rio
func (sr) GetClustered(session *Session, self SRRef) (retval bool, err error) {
	method := "SR.get_clustered"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeSRRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetClustered2: Get the clustered field of the given SR.
// Version: rio
func (sr) GetClustered2(session *Session, self SRRef) (retval bool, err error) {
	method := "SR.get_clustered"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeSRRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetIntroducedBy: Get the introduced_by field of the given SR.
// Version: rio
func (sr) GetIntroducedBy(session *Session, self SRRef) (retval DRTaskRef, err error) {
	method := "SR.get_introduced_by"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeSRRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeDRTaskRef(method+" -> ", result)
	return
}

// GetIntroducedBy2: Get the introduced_by field of the given SR.
// Version: rio
func (sr) GetIntroducedBy2(session *Session, self SRRef) (retval DRTaskRef, err error) {
	method := "SR.get_introduced_by"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeSRRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeDRTaskRef(method+" -> ", result)
	return
}

// GetLocalCacheEnabled: Get the local_cache_enabled field of the given SR.
// Version: rio
func (sr) GetLocalCacheEnabled(session *Session, self SRRef) (retval bool, err error) {
	method := "SR.get_local_cache_enabled"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeSRRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetLocalCacheEnabled2: Get the local_cache_enabled field of the given SR.
// Version: rio
func (sr) GetLocalCacheEnabled2(session *Session, self SRRef) (retval bool, err error) {
	method := "SR.get_local_cache_enabled"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeSRRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetBlobs: Get the blobs field of the given SR.
// Version: rio
func (sr) GetBlobs(session *Session, self SRRef) (retval map[string]BlobRef, err error) {
	method := "SR.get_blobs"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeSRRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetBlobs2: Get the blobs field of the given SR.
// Version: rio
func (sr) GetBlobs2(session *Session, self SRRef) (retval map[string]BlobRef, err error) {
	method := "SR.get_blobs"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeSRRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetSmConfig: Get the sm_config field of the given SR.
// Version: rio
func (sr) GetSmConfig(session *Session, self SRRef) (retval map[string]string, err error) {
	method := "SR.get_sm_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeSRRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetSmConfig2: Get the sm_config field of the given SR.
// Version: rio
func (sr) GetSmConfig2(session *Session, self SRRef) (retval map[string]string, err error) {
	method := "SR.get_sm_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeSRRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetTags: Get the tags field of the given SR.
// Version: rio
func (sr) GetTags(session *Session, self SRRef) (retval []string, err error) {
	method := "SR.get_tags"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeSRRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetTags2: Get the tags field of the given SR.
// Version: rio
func (sr) GetTags2(session *Session, self SRRef) (retval []string, err error) {
	method := "SR.get_tags"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeSRRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetOtherConfig: Get the other_config field of the given SR.
// Version: rio
func (sr) GetOtherConfig(session *Session, self SRRef) (retval map[string]string, err error) {
	method := "SR.get_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeSRRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetOtherConfig2: Get the other_config field of the given SR.
// Version: rio
func (sr) GetOtherConfig2(session *Session, self SRRef) (retval map[string]string, err error) {
	method := "SR.get_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeSRRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetShared: Get the shared field of the given SR.
// Version: rio
func (sr) GetShared(session *Session, self SRRef) (retval bool, err error) {
	method := "SR.get_shared"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeSRRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetShared2: Get the shared field of the given SR.
// Version: rio
func (sr) GetShared2(session *Session, self SRRef) (retval bool, err error) {
	method := "SR.get_shared"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeSRRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetContentType: Get the content_type field of the given SR.
// Version: rio
func (sr) GetContentType(session *Session, self SRRef) (retval string, err error) {
	method := "SR.get_content_type"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeSRRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetContentType2: Get the content_type field of the given SR.
// Version: rio
func (sr) GetContentType2(session *Session, self SRRef) (retval string, err error) {
	method := "SR.get_content_type"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeSRRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetType: Get the type field of the given SR.
// Version: rio
func (sr) GetType(session *Session, self SRRef) (retval string, err error) {
	method := "SR.get_type"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeSRRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetType2: Get the type field of the given SR.
// Version: rio
func (sr) GetType2(session *Session, self SRRef) (retval string, err error) {
	method := "SR.get_type"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeSRRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetPhysicalSize: Get the physical_size field of the given SR.
// Version: rio
func (sr) GetPhysicalSize(session *Session, self SRRef) (retval int, err error) {
	method := "SR.get_physical_size"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeSRRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetPhysicalSize2: Get the physical_size field of the given SR.
// Version: rio
func (sr) GetPhysicalSize2(session *Session, self SRRef) (retval int, err error) {
	method := "SR.get_physical_size"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeSRRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetPhysicalUtilisation: Get the physical_utilisation field of the given SR.
// Version: rio
func (sr) GetPhysicalUtilisation(session *Session, self SRRef) (retval int, err error) {
	method := "SR.get_physical_utilisation"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeSRRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetPhysicalUtilisation2: Get the physical_utilisation field of the given SR.
// Version: rio
func (sr) GetPhysicalUtilisation2(session *Session, self SRRef) (retval int, err error) {
	method := "SR.get_physical_utilisation"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeSRRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetVirtualAllocation: Get the virtual_allocation field of the given SR.
// Version: rio
func (sr) GetVirtualAllocation(session *Session, self SRRef) (retval int, err error) {
	method := "SR.get_virtual_allocation"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeSRRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetVirtualAllocation2: Get the virtual_allocation field of the given SR.
// Version: rio
func (sr) GetVirtualAllocation2(session *Session, self SRRef) (retval int, err error) {
	method := "SR.get_virtual_allocation"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeSRRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetPBDs: Get the PBDs field of the given SR.
// Version: rio
func (sr) GetPBDs(session *Session, self SRRef) (retval []PBDRef, err error) {
	method := "SR.get_PBDs"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeSRRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetPBDs2: Get the PBDs field of the given SR.
// Version: rio
func (sr) GetPBDs2(session *Session, self SRRef) (retval []PBDRef, err error) {
	method := "SR.get_PBDs"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeSRRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetVDIs: Get the VDIs field of the given SR.
// Version: rio
func (sr) GetVDIs(session *Session, self SRRef) (retval []VDIRef, err error) {
	method := "SR.get_VDIs"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeSRRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetVDIs2: Get the VDIs field of the given SR.
// Version: rio
func (sr) GetVDIs2(session *Session, self SRRef) (retval []VDIRef, err error) {
	method := "SR.get_VDIs"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeSRRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetCurrentOperations: Get the current_operations field of the given SR.
// Version: rio
func (sr) GetCurrentOperations(session *Session, self SRRef) (retval map[string]StorageOperations, err error) {
	method := "SR.get_current_operations"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeSRRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeStringToEnumStorageOperationsMap(method+" -> ", result)
	return
}

// GetCurrentOperations2: Get the current_operations field of the given SR.
// Version: rio
func (sr) GetCurrentOperations2(session *Session, self SRRef) (retval map[string]StorageOperations, err error) {
	method := "SR.get_current_operations"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeSRRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeStringToEnumStorageOperationsMap(method+" -> ", result)
	return
}

// GetAllowedOperations: Get the allowed_operations field of the given SR.
// Version: rio
func (sr) GetAllowedOperations(session *Session, self SRRef) (retval []StorageOperations, err error) {
	method := "SR.get_allowed_operations"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeSRRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeEnumStorageOperationsSet(method+" -> ", result)
	return
}

// GetAllowedOperations2: Get the allowed_operations field of the given SR.
// Version: rio
func (sr) GetAllowedOperations2(session *Session, self SRRef) (retval []StorageOperations, err error) {
	method := "SR.get_allowed_operations"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeSRRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeEnumStorageOperationsSet(method+" -> ", result)
	return
}

// GetNameDescription: Get the name/description field of the given SR.
// Version: rio
func (sr) GetNameDescription(session *Session, self SRRef) (retval string, err error) {
	method := "SR.get_name_description"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeSRRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetNameDescription2: Get the name/description field of the given SR.
// Version: rio
func (sr) GetNameDescription2(session *Session, self SRRef) (retval string, err error) {
	method := "SR.get_name_description"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeSRRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetNameLabel: Get the name/label field of the given SR.
// Version: rio
func (sr) GetNameLabel(session *Session, self SRRef) (retval string, err error) {
	method := "SR.get_name_label"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeSRRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetNameLabel2: Get the name/label field of the given SR.
// Version: rio
func (sr) GetNameLabel2(session *Session, self SRRef) (retval string, err error) {
	method := "SR.get_name_label"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeSRRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetUUID: Get the uuid field of the given SR.
// Version: rio
func (sr) GetUUID(session *Session, self SRRef) (retval string, err error) {
	method := "SR.get_uuid"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeSRRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetUUID2: Get the uuid field of the given SR.
// Version: rio
func (sr) GetUUID2(session *Session, self SRRef) (retval string, err error) {
	method := "SR.get_uuid"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeSRRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetByNameLabel: Get all the SR instances with the given label.
// Version: rio
func (sr) GetByNameLabel(session *Session, label string) (retval []SRRef, err error) {
	method := "SR.get_by_name_label"
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
	retval, err = deserializeSRRefSet(method+" -> ", result)
	return
}

// GetByNameLabel2: Get all the SR instances with the given label.
// Version: rio
func (sr) GetByNameLabel2(session *Session, label string) (retval []SRRef, err error) {
	method := "SR.get_by_name_label"
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
	retval, err = deserializeSRRefSet(method+" -> ", result)
	return
}

// GetByUUID: Get a reference to the SR instance with the specified UUID.
// Version: rio
func (sr) GetByUUID(session *Session, uuid string) (retval SRRef, err error) {
	method := "SR.get_by_uuid"
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
	retval, err = deserializeSRRef(method+" -> ", result)
	return
}

// GetByUUID2: Get a reference to the SR instance with the specified UUID.
// Version: rio
func (sr) GetByUUID2(session *Session, uuid string) (retval SRRef, err error) {
	method := "SR.get_by_uuid"
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
	retval, err = deserializeSRRef(method+" -> ", result)
	return
}

// GetRecord: Get a record containing the current state of the given SR.
// Version: rio
func (sr) GetRecord(session *Session, self SRRef) (retval SRRecord, err error) {
	method := "SR.get_record"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeSRRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeSRRecord(method+" -> ", result)
	return
}

// GetRecord2: Get a record containing the current state of the given SR.
// Version: rio
func (sr) GetRecord2(session *Session, self SRRef) (retval SRRecord, err error) {
	method := "SR.get_record"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeSRRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeSRRecord(method+" -> ", result)
	return
}
