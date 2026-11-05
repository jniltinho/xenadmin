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
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"
)

type Request struct {
	JSONRPC string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params,omitempty"`
	ID      int         `json:"id"`
}

type ResponseError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type Response struct {
	JSONRPC string         `json:"jsonrpc"`
	Result  interface{}    `json:"result,omitempty"`
	Error   *ResponseError `json:"error,omitempty"`
	ID      int            `json:"id"`
}

type rpcClient struct {
	endpoint   string
	httpClient *http.Client
	headers    map[string]string
	ctx        context.Context // XCP-NG PATCH (Session.WithContext): the calls end with it; nil is context.Background
}

func (client *rpcClient) newRequest(ctx context.Context, req interface{}) (*http.Request, error) {
	dataByte, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %w", err)
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, client.endpoint, bytes.NewReader(dataByte))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	request.Header.Set("Content-Type", "application/json; charset=utf-8")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("User-Agent", "XenAPI/"+APIVersionLatest.String())
	for k, v := range client.headers {
		request.Header.Set(k, v)
	}

	return request, nil
}

func convertUnhandledJSONData(jsonBytes []byte) []byte {
	jsonString := string(jsonBytes)
	re := regexp.MustCompile(`:[ ]*-Infinity`)
	jsonString = re.ReplaceAllString(jsonString, `:"-Inf"`)
	re = regexp.MustCompile(`:[ +]*Infinity`)
	jsonString = re.ReplaceAllString(jsonString, `:"+Inf"`)
	re = regexp.MustCompile(`:[ ]*NaN`)
	jsonString = re.ReplaceAllString(jsonString, `:"NaN"`)

	return []byte(jsonString)
}

func (client *rpcClient) call(ctx context.Context, methodName string, params ...interface{}) (*Response, error) {
	request := &Request{
		ID:      0,
		Method:  methodName,
		Params:  params,
		JSONRPC: "2.0",
	}

	httpRequest, err := client.newRequest(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("could not create request for %v() : %w", request.Method, err)
	}

	httpResponse, err := client.httpClient.Do(httpRequest)
	if err != nil {
		return nil, fmt.Errorf("call %v() on %v. Error making http request: %w", request.Method, httpRequest.URL.Redacted(), err)
	}
	defer httpResponse.Body.Close()

	var rpcResponse *Response
	body, err := io.ReadAll(httpResponse.Body)
	if err != nil {
		return nil, fmt.Errorf("call %v() on %v status code: %v. Could not read response body: %w", request.Method, httpRequest.URL.Redacted(), httpResponse.StatusCode, err)
	}
	body = convertUnhandledJSONData(body)
	err = json.Unmarshal(body, &rpcResponse)
	if err != nil && !errors.Is(err, io.EOF) {
		return nil, fmt.Errorf("call %v() on %v status code: %v. Could not decode response body: %w", request.Method, httpRequest.URL.Redacted(), httpResponse.StatusCode, err)
	}

	if rpcResponse == nil {
		return nil, fmt.Errorf("call %v() on %v status code: %v. Response missing", request.Method, httpRequest.URL.Redacted(), httpResponse.StatusCode)
	}

	return rpcResponse, nil
}

// XCP-NG PATCH (APIError): a XAPI failure keeps its code (e.g. HOST_IS_SLAVE) and parameters, so callers can
// react to it with errors.As instead of parsing a string.
type APIError struct {
	Code    int
	Message string
	Data    []string
}

func (e *APIError) Error() string {
	if len(e.Data) > 0 {
		return fmt.Sprintf("API error: code %d, message %s, data %v", e.Code, e.Message, e.Data)
	}
	return fmt.Sprintf("API error: code %d, message %s", e.Code, e.Message)
}

// ErrorCode is the XAPI failure name (the JSON-RPC message), e.g. SESSION_AUTHENTICATION_FAILED.
func (e *APIError) ErrorCode() string { return e.Message }

func (client *rpcClient) sendCall(methodName string, params ...interface{}) (result interface{}, err error) {
	ctx := client.ctx // XCP-NG PATCH (Session.WithContext)
	if ctx == nil {
		ctx = context.Background()
	}
	response, err := client.call(ctx, methodName, params...)
	if err != nil {
		return
	}

	if response.Error != nil {
		apiErr := &APIError{Code: response.Error.Code, Message: response.Error.Message}
		if list, ok := response.Error.Data.([]interface{}); ok {
			for _, d := range list {
				apiErr.Data = append(apiErr.Data, fmt.Sprint(d))
			}
		} else if response.Error.Data != nil {
			apiErr.Data = []string{fmt.Sprint(response.Error.Data)}
		}
		err = apiErr
		return
	}

	result = response.Result
	return
}

type SecureOpts struct {
	ServerCert string
	ClientCert string
	ClientKey  string
}

type ClientOpts struct {
	URL        string
	SecureOpts *SecureOpts
	Timeout    int
	Headers    map[string]string
	// XCP-NG PATCH (HTTPClient): the application's client (its TLS trust, no redirects, timeouts). When set, it is
	// used as is and SecureOpts/Timeout are ignored, so the SDK never opens a connection the app did not configure.
	HTTPClient *http.Client
}

// XCP-NG PATCH (no log.Fatal): configuration errors are returned instead of ending the process.
func newJSONRPCClient(opts *ClientOpts) (*rpcClient, error) {
	client := &rpcClient{
		endpoint:   fmt.Sprintf("%s%s", opts.URL, "/jsonrpc"),
		httpClient: &http.Client{},
		headers:    make(map[string]string),
	}
	for k, v := range opts.Headers {
		client.headers[k] = v
	}
	if opts.HTTPClient != nil {
		client.httpClient = opts.HTTPClient
		return client, nil
	}

	u, err := url.Parse(opts.URL)
	if err != nil {
		return nil, err
	}
	if strings.Compare(u.Scheme, "https") == 0 {
		skipVerify := true
		caCertPool := x509.NewCertPool()
		certs := []tls.Certificate{}
		if opts.SecureOpts != nil {
			if opts.SecureOpts.ServerCert != "" {
				skipVerify = false
				caCert, err := os.ReadFile(opts.SecureOpts.ServerCert)
				if err != nil {
					return nil, err
				}
				ok := caCertPool.AppendCertsFromPEM(caCert)
				if !ok {
					return nil, errors.New("failed to parse CA certificate")
				}
			}
			if opts.SecureOpts.ClientCert != "" || opts.SecureOpts.ClientKey != "" {
				if opts.SecureOpts.ClientCert == "" {
					return nil, errors.New("missing client certificate")
				}
				if opts.SecureOpts.ClientKey == "" {
					return nil, errors.New("missing client private key")
				}
				cert, err := tls.LoadX509KeyPair(opts.SecureOpts.ClientCert, opts.SecureOpts.ClientKey)
				if err != nil {
					return nil, err
				}
				certs = []tls.Certificate{cert}
			}
		}
		tlsConfig := &tls.Config{
			RootCAs:            caCertPool,
			Certificates:       certs,
			InsecureSkipVerify: skipVerify, // #nosec
			CipherSuites: []uint16{
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			},
			MinVersion:               tls.VersionTLS12,
			PreferServerCipherSuites: true,
		}
		transport := &http.Transport{
			TLSClientConfig: tlsConfig,
		}
		client.httpClient.Transport = transport
	}

	if opts.Timeout != 0 {
		client.httpClient.Timeout = time.Duration(opts.Timeout) * time.Second
	}

	return client, nil
}
