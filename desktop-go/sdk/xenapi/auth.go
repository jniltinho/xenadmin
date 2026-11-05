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

type AuthRecord struct {
}

type AuthRef string

// Management of remote authentication services
type auth struct{}

var Auth auth

// GetGroupMembership: This calls queries the external directory service to obtain the transitively-closed set of groups that the the subject_identifier is member of.
// Version: george
func (auth) GetGroupMembership(session *Session, subjectIdentifier string) (retval []string, err error) {
	method := "auth.get_group_membership"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	subjectIdentifierArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "subject_identifier"), subjectIdentifier)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, subjectIdentifierArg)
	if err != nil {
		return
	}
	retval, err = deserializeStringSet(method+" -> ", result)
	return
}

// GetGroupMembership2: This calls queries the external directory service to obtain the transitively-closed set of groups that the the subject_identifier is member of.
// Version: george
func (auth) GetGroupMembership2(session *Session, subjectIdentifier string) (retval []string, err error) {
	method := "auth.get_group_membership"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	subjectIdentifierArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "subject_identifier"), subjectIdentifier)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, subjectIdentifierArg)
	if err != nil {
		return
	}
	retval, err = deserializeStringSet(method+" -> ", result)
	return
}

// GetSubjectInformationFromIdentifier: This call queries the external directory service to obtain the user information (for example, username, organization etc.) from the specified subject_identifier
// Version: george
func (auth) GetSubjectInformationFromIdentifier(session *Session, subjectIdentifier string) (retval map[string]string, err error) {
	method := "auth.get_subject_information_from_identifier"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	subjectIdentifierArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "subject_identifier"), subjectIdentifier)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, subjectIdentifierArg)
	if err != nil {
		return
	}
	retval, err = deserializeStringToStringMap(method+" -> ", result)
	return
}

// GetSubjectInformationFromIdentifier2: This call queries the external directory service to obtain the user information (for example, username, organization etc.) from the specified subject_identifier
// Version: george
func (auth) GetSubjectInformationFromIdentifier2(session *Session, subjectIdentifier string) (retval map[string]string, err error) {
	method := "auth.get_subject_information_from_identifier"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	subjectIdentifierArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "subject_identifier"), subjectIdentifier)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, subjectIdentifierArg)
	if err != nil {
		return
	}
	retval, err = deserializeStringToStringMap(method+" -> ", result)
	return
}

// GetSubjectIdentifier: This call queries the external directory service to obtain the subject_identifier as a string from the human-readable subject_name
// Version: george
func (auth) GetSubjectIdentifier(session *Session, subjectName string) (retval string, err error) {
	method := "auth.get_subject_identifier"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	subjectNameArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "subject_name"), subjectName)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, subjectNameArg)
	if err != nil {
		return
	}
	retval, err = deserializeString(method+" -> ", result)
	return
}

// GetSubjectIdentifier2: This call queries the external directory service to obtain the subject_identifier as a string from the human-readable subject_name
// Version: george
func (auth) GetSubjectIdentifier2(session *Session, subjectName string) (retval string, err error) {
	method := "auth.get_subject_identifier"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	subjectNameArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "subject_name"), subjectName)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, subjectNameArg)
	if err != nil {
		return
	}
	retval, err = deserializeString(method+" -> ", result)
	return
}
