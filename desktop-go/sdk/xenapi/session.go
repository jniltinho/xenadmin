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
	"net/http"
	"time"
)

type SessionRecord struct {
	// Unique identifier/object reference
	UUID string `json:"uuid,omitempty"`
	// Currently connected host
	ThisHost HostRef `json:"this_host,omitempty"`
	// Currently connected user
	ThisUser UserRef `json:"this_user,omitempty"`
	// Timestamp for last time session was active
	LastActive time.Time `json:"last_active,omitempty"`
	// True if this session relates to a intra-pool login, false otherwise
	Pool bool `json:"pool,omitempty"`
	// additional configuration
	OtherConfig map[string]string `json:"other_config,omitempty"`
	// true iff this session was created using local superuser credentials
	IsLocalSuperuser bool `json:"is_local_superuser,omitempty"`
	// references the subject instance that created the session. If a session instance has is_local_superuser set, then the value of this field is undefined.
	Subject SubjectRef `json:"subject,omitempty"`
	// time when session was last validated
	ValidationTime time.Time `json:"validation_time,omitempty"`
	// the subject identifier of the user that was externally authenticated. If a session instance has is_local_superuser set, then the value of this field is undefined.
	AuthUserSid string `json:"auth_user_sid,omitempty"`
	// the subject name of the user that was externally authenticated. If a session instance has is_local_superuser set, then the value of this field is undefined.
	AuthUserName string `json:"auth_user_name,omitempty"`
	// list with all RBAC permissions for this session
	RbacPermissions []string `json:"rbac_permissions,omitempty"`
	// list of tasks created using the current session
	Tasks []TaskRef `json:"tasks,omitempty"`
	// references the parent session that created this session
	Parent SessionRef `json:"parent,omitempty"`
	// a key string provided by a API user to distinguish itself from other users sharing the same login name
	Originator string `json:"originator,omitempty"`
	// indicates whether this session was authenticated using a client certificate
	ClientCertificate bool `json:"client_certificate,omitempty"`
}

type SessionRef string

// A session
type Session struct {
	APIVersion  APIVersion
	client      *rpcClient
	ref         SessionRef
	XAPIVersion string
}

func NewSession(opts *ClientOpts) *Session {
	session, err := NewSessionE(opts)
	if err != nil {
		// XCP-NG PATCH (no log.Fatal): the upstream API has no error result; calls on this session fail instead
		return &Session{client: &rpcClient{endpoint: opts.URL, httpClient: &http.Client{Transport: failingTransport{err}}, headers: map[string]string{}}}
	}
	return session
}

// NewSessionE is NewSession with the configuration error returned (XCP-NG PATCH).
func NewSessionE(opts *ClientOpts) (*Session, error) {
	client, err := newJSONRPCClient(opts)
	if err != nil {
		return nil, err
	}
	return &Session{client: client}, nil
}

type failingTransport struct{ err error }

func (t failingTransport) RoundTrip(*http.Request) (*http.Response, error) { return nil, t.err }

// Ref is the session reference, for HTTP calls such as the console CONNECT (XCP-NG PATCH).
func (class *Session) Ref() SessionRef { return class.ref }

// LogoutSubjectIdentifier: Log out all sessions associated to a user subject-identifier, except the session associated with the context calling this function
// Version: george
func (class *Session) LogoutSubjectIdentifier(subjectIdentifier string) (err error) {
	method := "session.logout_subject_identifier"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), class.ref)
	if err != nil {
		return
	}
	subjectIdentifierArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "subject_identifier"), subjectIdentifier)
	if err != nil {
		return
	}
	_, err = class.client.sendCall(method, sessionIDArg, subjectIdentifierArg)
	return
}

// AsyncLogoutSubjectIdentifier: Log out all sessions associated to a user subject-identifier, except the session associated with the context calling this function
// Version: george
func (class *Session) AsyncLogoutSubjectIdentifier(subjectIdentifier string) (retval TaskRef, err error) {
	method := "Async.session.logout_subject_identifier"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), class.ref)
	if err != nil {
		return
	}
	subjectIdentifierArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "subject_identifier"), subjectIdentifier)
	if err != nil {
		return
	}
	result, err := class.client.sendCall(method, sessionIDArg, subjectIdentifierArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// LogoutSubjectIdentifier1: Log out all sessions associated to a user subject-identifier, except the session associated with the context calling this function
// Version: george
func (class *Session) LogoutSubjectIdentifier1(subjectIdentifier string) (err error) {
	method := "session.logout_subject_identifier"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), class.ref)
	if err != nil {
		return
	}
	subjectIdentifierArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "subject_identifier"), subjectIdentifier)
	if err != nil {
		return
	}
	_, err = class.client.sendCall(method, sessionIDArg, subjectIdentifierArg)
	return
}

// AsyncLogoutSubjectIdentifier1: Log out all sessions associated to a user subject-identifier, except the session associated with the context calling this function
// Version: george
func (class *Session) AsyncLogoutSubjectIdentifier1(subjectIdentifier string) (retval TaskRef, err error) {
	method := "Async.session.logout_subject_identifier"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), class.ref)
	if err != nil {
		return
	}
	subjectIdentifierArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "subject_identifier"), subjectIdentifier)
	if err != nil {
		return
	}
	result, err := class.client.sendCall(method, sessionIDArg, subjectIdentifierArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// GetAllSubjectIdentifiers: Return a list of all the user subject-identifiers of all existing sessions
// Version: george
func (class *Session) GetAllSubjectIdentifiers() (retval []string, err error) {
	method := "session.get_all_subject_identifiers"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), class.ref)
	if err != nil {
		return
	}
	result, err := class.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializeStringSet(method+" -> ", result)
	return
}

// AsyncGetAllSubjectIdentifiers: Return a list of all the user subject-identifiers of all existing sessions
// Version: george
func (class *Session) AsyncGetAllSubjectIdentifiers() (retval TaskRef, err error) {
	method := "Async.session.get_all_subject_identifiers"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), class.ref)
	if err != nil {
		return
	}
	result, err := class.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// GetAllSubjectIdentifiers0: Return a list of all the user subject-identifiers of all existing sessions
// Version: george
func (class *Session) GetAllSubjectIdentifiers0() (retval []string, err error) {
	method := "session.get_all_subject_identifiers"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), class.ref)
	if err != nil {
		return
	}
	result, err := class.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializeStringSet(method+" -> ", result)
	return
}

// AsyncGetAllSubjectIdentifiers0: Return a list of all the user subject-identifiers of all existing sessions
// Version: george
func (class *Session) AsyncGetAllSubjectIdentifiers0() (retval TaskRef, err error) {
	method := "Async.session.get_all_subject_identifiers"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), class.ref)
	if err != nil {
		return
	}
	result, err := class.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// LocalLogout: Log out of local session.
// Version: miami
func (class *Session) LocalLogout() (err error) {
	method := "session.local_logout"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), class.ref)
	if err != nil {
		return
	}
	_, err = class.client.sendCall(method, sessionIDArg)
	class.ref = ""
	return
}

// LocalLogout0: Log out of local session.
// Version: miami
func (class *Session) LocalLogout0() (err error) {
	method := "session.local_logout"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), class.ref)
	if err != nil {
		return
	}
	_, err = class.client.sendCall(method, sessionIDArg)
	class.ref = ""
	return
}

// CreateFromDBFile:
// Version: dundee
func (class *Session) CreateFromDBFile(filename string) (retval SessionRef, err error) {
	method := "session.create_from_db_file"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), class.ref)
	if err != nil {
		return
	}
	filenameArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "filename"), filename)
	if err != nil {
		return
	}
	result, err := class.client.sendCall(method, sessionIDArg, filenameArg)
	if err != nil {
		return
	}
	retval, err = deserializeSessionRef(method+" -> ", result)
	return
}

// AsyncCreateFromDBFile:
// Version: dundee
func (class *Session) AsyncCreateFromDBFile(filename string) (retval TaskRef, err error) {
	method := "Async.session.create_from_db_file"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), class.ref)
	if err != nil {
		return
	}
	filenameArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "filename"), filename)
	if err != nil {
		return
	}
	result, err := class.client.sendCall(method, sessionIDArg, filenameArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// CreateFromDBFile1:
// Version: dundee
func (class *Session) CreateFromDBFile1(filename string) (retval SessionRef, err error) {
	method := "session.create_from_db_file"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), class.ref)
	if err != nil {
		return
	}
	filenameArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "filename"), filename)
	if err != nil {
		return
	}
	result, err := class.client.sendCall(method, sessionIDArg, filenameArg)
	if err != nil {
		return
	}
	retval, err = deserializeSessionRef(method+" -> ", result)
	return
}

// AsyncCreateFromDBFile1:
// Version: dundee
func (class *Session) AsyncCreateFromDBFile1(filename string) (retval TaskRef, err error) {
	method := "Async.session.create_from_db_file"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), class.ref)
	if err != nil {
		return
	}
	filenameArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "filename"), filename)
	if err != nil {
		return
	}
	result, err := class.client.sendCall(method, sessionIDArg, filenameArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// SlaveLocalLoginWithPassword: Authenticate locally against a slave in emergency mode. Note the resulting sessions are only good for use on this host.
// Version: miami
func (class *Session) SlaveLocalLoginWithPassword(uname string, pwd string) (retval SessionRef, err error) {
	method := "session.slave_local_login_with_password"
	unameArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "uname"), uname)
	if err != nil {
		return
	}
	pwdArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "pwd"), pwd)
	if err != nil {
		return
	}
	result, err := class.client.sendCall(method, unameArg, pwdArg)
	if err != nil {
		return
	}
	retval, err = deserializeSessionRef(method+" -> ", result)
	if err != nil {
		return
	}
	class.ref = retval
	err = setSessionDetails(class)
	return
}

// SlaveLocalLoginWithPassword2: Authenticate locally against a slave in emergency mode. Note the resulting sessions are only good for use on this host.
// Version: miami
func (class *Session) SlaveLocalLoginWithPassword2(uname string, pwd string) (retval SessionRef, err error) {
	method := "session.slave_local_login_with_password"
	unameArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "uname"), uname)
	if err != nil {
		return
	}
	pwdArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "pwd"), pwd)
	if err != nil {
		return
	}
	result, err := class.client.sendCall(method, unameArg, pwdArg)
	if err != nil {
		return
	}
	retval, err = deserializeSessionRef(method+" -> ", result)
	if err != nil {
		return
	}
	class.ref = retval
	err = setSessionDetails(class)
	return
}

// ChangePassword: Change the account password; if your session is authenticated with root privileges then the old_pwd is validated and the new_pwd is set regardless
// Version: rio
func (class *Session) ChangePassword(oldPwd string, newPwd string) (err error) {
	method := "session.change_password"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), class.ref)
	if err != nil {
		return
	}
	oldPwdArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "old_pwd"), oldPwd)
	if err != nil {
		return
	}
	newPwdArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "new_pwd"), newPwd)
	if err != nil {
		return
	}
	_, err = class.client.sendCall(method, sessionIDArg, oldPwdArg, newPwdArg)
	return
}

// ChangePassword2: Change the account password; if your session is authenticated with root privileges then the old_pwd is validated and the new_pwd is set regardless
// Version: rio
func (class *Session) ChangePassword2(oldPwd string, newPwd string) (err error) {
	method := "session.change_password"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), class.ref)
	if err != nil {
		return
	}
	oldPwdArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "old_pwd"), oldPwd)
	if err != nil {
		return
	}
	newPwdArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "new_pwd"), newPwd)
	if err != nil {
		return
	}
	_, err = class.client.sendCall(method, sessionIDArg, oldPwdArg, newPwdArg)
	return
}

// Logout: Log out of a session
// Version: rio
func (class *Session) Logout() (err error) {
	method := "session.logout"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), class.ref)
	if err != nil {
		return
	}
	_, err = class.client.sendCall(method, sessionIDArg)
	class.ref = ""
	return
}

// Logout0: Log out of a session
// Version: rio
func (class *Session) Logout0() (err error) {
	method := "session.logout"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), class.ref)
	if err != nil {
		return
	}
	_, err = class.client.sendCall(method, sessionIDArg)
	class.ref = ""
	return
}

// LoginWithPassword: Attempt to authenticate the user, returning a session reference if successful
// Version: clearwater
//
// Errors:
// SESSION_AUTHENTICATION_FAILED - The credentials given by the user are incorrect, so access has been denied, and you have not been issued a session handle.
// HOST_IS_SLAVE - You cannot make regular API calls directly on a supporter. Please pass API calls via the coordinator host.
func (class *Session) LoginWithPassword(uname string, pwd string, version string, originator string) (retval SessionRef, err error) {
	method := "session.login_with_password"
	unameArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "uname"), uname)
	if err != nil {
		return
	}
	pwdArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "pwd"), pwd)
	if err != nil {
		return
	}
	versionArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "version"), version)
	if err != nil {
		return
	}
	originatorArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "originator"), originator)
	if err != nil {
		return
	}
	result, err := class.client.sendCall(method, unameArg, pwdArg, versionArg, originatorArg)
	if err != nil {
		return
	}
	retval, err = deserializeSessionRef(method+" -> ", result)
	if err != nil {
		return
	}
	class.ref = retval
	err = setSessionDetails(class)
	return
}

// LoginWithPassword4: Attempt to authenticate the user, returning a session reference if successful
// Version: clearwater
//
// Errors:
// SESSION_AUTHENTICATION_FAILED - The credentials given by the user are incorrect, so access has been denied, and you have not been issued a session handle.
// HOST_IS_SLAVE - You cannot make regular API calls directly on a supporter. Please pass API calls via the coordinator host.
func (class *Session) LoginWithPassword4(uname string, pwd string, version string, originator string) (retval SessionRef, err error) {
	method := "session.login_with_password"
	unameArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "uname"), uname)
	if err != nil {
		return
	}
	pwdArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "pwd"), pwd)
	if err != nil {
		return
	}
	versionArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "version"), version)
	if err != nil {
		return
	}
	originatorArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "originator"), originator)
	if err != nil {
		return
	}
	result, err := class.client.sendCall(method, unameArg, pwdArg, versionArg, originatorArg)
	if err != nil {
		return
	}
	retval, err = deserializeSessionRef(method+" -> ", result)
	if err != nil {
		return
	}
	class.ref = retval
	err = setSessionDetails(class)
	return
}

// LoginWithPassword3: Attempt to authenticate the user, returning a session reference if successful
// Version: miami
//
// Errors:
// SESSION_AUTHENTICATION_FAILED - The credentials given by the user are incorrect, so access has been denied, and you have not been issued a session handle.
// HOST_IS_SLAVE - You cannot make regular API calls directly on a supporter. Please pass API calls via the coordinator host.
func (class *Session) LoginWithPassword3(uname string, pwd string, version string) (retval SessionRef, err error) {
	method := "session.login_with_password"
	unameArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "uname"), uname)
	if err != nil {
		return
	}
	pwdArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "pwd"), pwd)
	if err != nil {
		return
	}
	versionArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "version"), version)
	if err != nil {
		return
	}
	result, err := class.client.sendCall(method, unameArg, pwdArg, versionArg)
	if err != nil {
		return
	}
	retval, err = deserializeSessionRef(method+" -> ", result)
	if err != nil {
		return
	}
	class.ref = retval
	err = setSessionDetails(class)
	return
}

// LoginWithPassword2: Attempt to authenticate the user, returning a session reference if successful
// Version: rio
//
// Errors:
// SESSION_AUTHENTICATION_FAILED - The credentials given by the user are incorrect, so access has been denied, and you have not been issued a session handle.
// HOST_IS_SLAVE - You cannot make regular API calls directly on a supporter. Please pass API calls via the coordinator host.
func (class *Session) LoginWithPassword2(uname string, pwd string) (retval SessionRef, err error) {
	method := "session.login_with_password"
	unameArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "uname"), uname)
	if err != nil {
		return
	}
	pwdArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "pwd"), pwd)
	if err != nil {
		return
	}
	result, err := class.client.sendCall(method, unameArg, pwdArg)
	if err != nil {
		return
	}
	retval, err = deserializeSessionRef(method+" -> ", result)
	if err != nil {
		return
	}
	class.ref = retval
	err = setSessionDetails(class)
	return
}

// RemoveFromOtherConfig: Remove the given key and its corresponding value from the other_config field of the given session.  If the key is not in that Map, then do nothing.
// Version: miami
func (class *Session) RemoveFromOtherConfig(self SessionRef, key string) (err error) {
	method := "session.remove_from_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), class.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	keyArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "key"), key)
	if err != nil {
		return
	}
	_, err = class.client.sendCall(method, sessionIDArg, selfArg, keyArg)
	return
}

// RemoveFromOtherConfig2: Remove the given key and its corresponding value from the other_config field of the given session.  If the key is not in that Map, then do nothing.
// Version: miami
func (class *Session) RemoveFromOtherConfig2(self SessionRef, key string) (err error) {
	method := "session.remove_from_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), class.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	keyArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "key"), key)
	if err != nil {
		return
	}
	_, err = class.client.sendCall(method, sessionIDArg, selfArg, keyArg)
	return
}

// RemoveFromOtherConfig1: Remove the given key and its corresponding value from the other_config field of the given session.  If the key is not in that Map, then do nothing.
// Version: rio
func (class *Session) RemoveFromOtherConfig1(self SessionRef) (err error) {
	method := "session.remove_from_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), class.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = class.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// AddToOtherConfig: Add the given key-value pair to the other_config field of the given session.
// Version: miami
func (class *Session) AddToOtherConfig(self SessionRef, key string, value string) (err error) {
	method := "session.add_to_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), class.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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
	_, err = class.client.sendCall(method, sessionIDArg, selfArg, keyArg, valueArg)
	return
}

// AddToOtherConfig3: Add the given key-value pair to the other_config field of the given session.
// Version: miami
func (class *Session) AddToOtherConfig3(self SessionRef, key string, value string) (err error) {
	method := "session.add_to_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), class.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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
	_, err = class.client.sendCall(method, sessionIDArg, selfArg, keyArg, valueArg)
	return
}

// AddToOtherConfig1: Add the given key-value pair to the other_config field of the given session.
// Version: rio
func (class *Session) AddToOtherConfig1(self SessionRef) (err error) {
	method := "session.add_to_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), class.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = class.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// SetOtherConfig: Set the other_config field of the given session.
// Version: miami
func (class *Session) SetOtherConfig(self SessionRef, value map[string]string) (err error) {
	method := "session.set_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), class.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = class.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	return
}

// SetOtherConfig2: Set the other_config field of the given session.
// Version: miami
func (class *Session) SetOtherConfig2(self SessionRef, value map[string]string) (err error) {
	method := "session.set_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), class.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = class.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	return
}

// SetOtherConfig1: Set the other_config field of the given session.
// Version: rio
func (class *Session) SetOtherConfig1(self SessionRef) (err error) {
	method := "session.set_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), class.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = class.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// GetClientCertificate: Get the client_certificate field of the given session.
// Version: rio
func (class *Session) GetClientCertificate(self SessionRef) (retval bool, err error) {
	method := "session.get_client_certificate"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), class.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := class.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeBool(method+" -> ", result)
	return
}

// GetClientCertificate1: Get the client_certificate field of the given session.
// Version: rio
func (class *Session) GetClientCertificate1(self SessionRef) (retval bool, err error) {
	method := "session.get_client_certificate"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), class.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := class.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeBool(method+" -> ", result)
	return
}

// GetOriginator: Get the originator field of the given session.
// Version: rio
func (class *Session) GetOriginator(self SessionRef) (retval string, err error) {
	method := "session.get_originator"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), class.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := class.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeString(method+" -> ", result)
	return
}

// GetOriginator1: Get the originator field of the given session.
// Version: rio
func (class *Session) GetOriginator1(self SessionRef) (retval string, err error) {
	method := "session.get_originator"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), class.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := class.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeString(method+" -> ", result)
	return
}

// GetParent: Get the parent field of the given session.
// Version: rio
func (class *Session) GetParent(self SessionRef) (retval SessionRef, err error) {
	method := "session.get_parent"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), class.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := class.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeSessionRef(method+" -> ", result)
	return
}

// GetParent1: Get the parent field of the given session.
// Version: rio
func (class *Session) GetParent1(self SessionRef) (retval SessionRef, err error) {
	method := "session.get_parent"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), class.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := class.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeSessionRef(method+" -> ", result)
	return
}

// GetTasks: Get the tasks field of the given session.
// Version: rio
func (class *Session) GetTasks(self SessionRef) (retval []TaskRef, err error) {
	method := "session.get_tasks"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), class.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := class.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRefSet(method+" -> ", result)
	return
}

// GetTasks1: Get the tasks field of the given session.
// Version: rio
func (class *Session) GetTasks1(self SessionRef) (retval []TaskRef, err error) {
	method := "session.get_tasks"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), class.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := class.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRefSet(method+" -> ", result)
	return
}

// GetRbacPermissions: Get the rbac_permissions field of the given session.
// Version: rio
func (class *Session) GetRbacPermissions(self SessionRef) (retval []string, err error) {
	method := "session.get_rbac_permissions"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), class.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := class.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeStringSet(method+" -> ", result)
	return
}

// GetRbacPermissions1: Get the rbac_permissions field of the given session.
// Version: rio
func (class *Session) GetRbacPermissions1(self SessionRef) (retval []string, err error) {
	method := "session.get_rbac_permissions"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), class.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := class.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeStringSet(method+" -> ", result)
	return
}

// GetAuthUserName: Get the auth_user_name field of the given session.
// Version: rio
func (class *Session) GetAuthUserName(self SessionRef) (retval string, err error) {
	method := "session.get_auth_user_name"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), class.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := class.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeString(method+" -> ", result)
	return
}

// GetAuthUserName1: Get the auth_user_name field of the given session.
// Version: rio
func (class *Session) GetAuthUserName1(self SessionRef) (retval string, err error) {
	method := "session.get_auth_user_name"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), class.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := class.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeString(method+" -> ", result)
	return
}

// GetAuthUserSid: Get the auth_user_sid field of the given session.
// Version: rio
func (class *Session) GetAuthUserSid(self SessionRef) (retval string, err error) {
	method := "session.get_auth_user_sid"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), class.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := class.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeString(method+" -> ", result)
	return
}

// GetAuthUserSid1: Get the auth_user_sid field of the given session.
// Version: rio
func (class *Session) GetAuthUserSid1(self SessionRef) (retval string, err error) {
	method := "session.get_auth_user_sid"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), class.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := class.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeString(method+" -> ", result)
	return
}

// GetValidationTime: Get the validation_time field of the given session.
// Version: rio
func (class *Session) GetValidationTime(self SessionRef) (retval time.Time, err error) {
	method := "session.get_validation_time"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), class.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := class.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeTime(method+" -> ", result)
	return
}

// GetValidationTime1: Get the validation_time field of the given session.
// Version: rio
func (class *Session) GetValidationTime1(self SessionRef) (retval time.Time, err error) {
	method := "session.get_validation_time"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), class.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := class.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeTime(method+" -> ", result)
	return
}

// GetSubject: Get the subject field of the given session.
// Version: rio
func (class *Session) GetSubject(self SessionRef) (retval SubjectRef, err error) {
	method := "session.get_subject"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), class.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := class.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeSubjectRef(method+" -> ", result)
	return
}

// GetSubject1: Get the subject field of the given session.
// Version: rio
func (class *Session) GetSubject1(self SessionRef) (retval SubjectRef, err error) {
	method := "session.get_subject"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), class.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := class.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeSubjectRef(method+" -> ", result)
	return
}

// GetIsLocalSuperuser: Get the is_local_superuser field of the given session.
// Version: rio
func (class *Session) GetIsLocalSuperuser(self SessionRef) (retval bool, err error) {
	method := "session.get_is_local_superuser"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), class.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := class.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeBool(method+" -> ", result)
	return
}

// GetIsLocalSuperuser1: Get the is_local_superuser field of the given session.
// Version: rio
func (class *Session) GetIsLocalSuperuser1(self SessionRef) (retval bool, err error) {
	method := "session.get_is_local_superuser"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), class.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := class.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeBool(method+" -> ", result)
	return
}

// GetOtherConfig: Get the other_config field of the given session.
// Version: rio
func (class *Session) GetOtherConfig(self SessionRef) (retval map[string]string, err error) {
	method := "session.get_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), class.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := class.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeStringToStringMap(method+" -> ", result)
	return
}

// GetOtherConfig1: Get the other_config field of the given session.
// Version: rio
func (class *Session) GetOtherConfig1(self SessionRef) (retval map[string]string, err error) {
	method := "session.get_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), class.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := class.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeStringToStringMap(method+" -> ", result)
	return
}

// GetPool: Get the pool field of the given session.
// Version: rio
func (class *Session) GetPool(self SessionRef) (retval bool, err error) {
	method := "session.get_pool"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), class.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := class.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeBool(method+" -> ", result)
	return
}

// GetPool1: Get the pool field of the given session.
// Version: rio
func (class *Session) GetPool1(self SessionRef) (retval bool, err error) {
	method := "session.get_pool"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), class.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := class.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeBool(method+" -> ", result)
	return
}

// GetLastActive: Get the last_active field of the given session.
// Version: rio
func (class *Session) GetLastActive(self SessionRef) (retval time.Time, err error) {
	method := "session.get_last_active"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), class.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := class.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeTime(method+" -> ", result)
	return
}

// GetLastActive1: Get the last_active field of the given session.
// Version: rio
func (class *Session) GetLastActive1(self SessionRef) (retval time.Time, err error) {
	method := "session.get_last_active"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), class.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := class.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeTime(method+" -> ", result)
	return
}

// GetThisUser: Get the this_user field of the given session.
// Version: rio
func (class *Session) GetThisUser(self SessionRef) (retval UserRef, err error) {
	method := "session.get_this_user"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), class.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := class.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeUserRef(method+" -> ", result)
	return
}

// GetThisUser1: Get the this_user field of the given session.
// Version: rio
func (class *Session) GetThisUser1(self SessionRef) (retval UserRef, err error) {
	method := "session.get_this_user"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), class.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := class.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeUserRef(method+" -> ", result)
	return
}

// GetThisHost: Get the this_host field of the given session.
// Version: rio
func (class *Session) GetThisHost(self SessionRef) (retval HostRef, err error) {
	method := "session.get_this_host"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), class.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := class.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeHostRef(method+" -> ", result)
	return
}

// GetThisHost1: Get the this_host field of the given session.
// Version: rio
func (class *Session) GetThisHost1(self SessionRef) (retval HostRef, err error) {
	method := "session.get_this_host"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), class.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := class.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeHostRef(method+" -> ", result)
	return
}

// GetUUID: Get the uuid field of the given session.
// Version: rio
func (class *Session) GetUUID(self SessionRef) (retval string, err error) {
	method := "session.get_uuid"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), class.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := class.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeString(method+" -> ", result)
	return
}

// GetUUID1: Get the uuid field of the given session.
// Version: rio
func (class *Session) GetUUID1(self SessionRef) (retval string, err error) {
	method := "session.get_uuid"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), class.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := class.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeString(method+" -> ", result)
	return
}

// GetByUUID: Get a reference to the session instance with the specified UUID.
// Version: rio
func (class *Session) GetByUUID(uuid string) (retval SessionRef, err error) {
	method := "session.get_by_uuid"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), class.ref)
	if err != nil {
		return
	}
	uuidArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "uuid"), uuid)
	if err != nil {
		return
	}
	result, err := class.client.sendCall(method, sessionIDArg, uuidArg)
	if err != nil {
		return
	}
	retval, err = deserializeSessionRef(method+" -> ", result)
	return
}

// GetByUUID1: Get a reference to the session instance with the specified UUID.
// Version: rio
func (class *Session) GetByUUID1(uuid string) (retval SessionRef, err error) {
	method := "session.get_by_uuid"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), class.ref)
	if err != nil {
		return
	}
	uuidArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "uuid"), uuid)
	if err != nil {
		return
	}
	result, err := class.client.sendCall(method, sessionIDArg, uuidArg)
	if err != nil {
		return
	}
	retval, err = deserializeSessionRef(method+" -> ", result)
	return
}

// GetRecord: Get a record containing the current state of the given session.
// Version: rio
func (class *Session) GetRecord(self SessionRef) (retval SessionRecord, err error) {
	method := "session.get_record"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), class.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := class.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeSessionRecord(method+" -> ", result)
	return
}

// GetRecord1: Get a record containing the current state of the given session.
// Version: rio
func (class *Session) GetRecord1(self SessionRef) (retval SessionRecord, err error) {
	method := "session.get_record"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), class.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := class.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeSessionRecord(method+" -> ", result)
	return
}
