package xapi

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"go/xenapi"

	"xcpng-center/internal/trust"
)

type memStore map[string]string

func (m memStore) Fingerprint(k string) (string, bool) { v, ok := m[k]; return v, ok }
func (m memStore) Pin(k, fp string) error              { m[k] = fp; return nil }

// fakeXAPI answers session.login_with_password by user name ("bad" fails the login).
func fakeXAPI(t *testing.T, behaviour func(user string) (any, *xenapi.ResponseError)) *httptest.Server {
	return fakeXAPIWith(t, func(user, _ string) (any, *xenapi.ResponseError) { return behaviour(user) })
}

// fakeXAPIByHost answers the login by the Host the client asked for.
func fakeXAPIByHost(t *testing.T, behaviour func(host string) (any, *xenapi.ResponseError)) *httptest.Server {
	return fakeXAPIWith(t, func(_, host string) (any, *xenapi.ResponseError) { return behaviour(host) })
}

func fakeXAPIWith(t *testing.T, behaviour func(user, host string) (any, *xenapi.ResponseError)) *httptest.Server {
	return httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		var req struct {
			Method string `json:"method"`
			Params []any  `json:"params"`
		}
		_ = json.Unmarshal(body, &req)
		resp := map[string]any{"jsonrpc": "2.0", "id": 0}
		switch req.Method {
		case "session.login_with_password":
			result, rpcErr := behaviour(req.Params[0].(string), r.Host)
			if rpcErr != nil {
				resp["error"] = rpcErr
			} else {
				resp["result"] = result
			}
		case "session.logout":
			resp["result"] = ""
		// what the SDK reads after a login (setSessionDetails: the API version of the coordinator)
		case "pool.get_all":
			resp["result"] = []string{"OpaqueRef:pool"}
		case "pool.get_record":
			resp["result"] = map[string]any{"master": "OpaqueRef:host"}
		case "host.get_record":
			resp["result"] = map[string]any{"API_version_major": 2, "API_version_minor": 21, "software_version": map[string]string{"xapi": "25.33"}}
		default:
			resp["error"] = xenapi.ResponseError{Code: 1, Message: "MESSAGE_METHOD_UNKNOWN", Data: []any{req.Method}}
		}
		_ = json.NewEncoder(w).Encode(resp)
	}))
}

func split(t *testing.T, url string) (string, int) {
	h, p, err := ParseAddress(strings.TrimPrefix(url, "https://"))
	if err != nil {
		t.Fatal(err)
	}
	return h, p
}

func TestConnectTrustAndLogin(t *testing.T) {
	srv := fakeXAPI(t, func(user string) (any, *xenapi.ResponseError) {
		if user == "bad" {
			return nil, &xenapi.ResponseError{Code: 1, Message: "SESSION_AUTHENTICATION_FAILED", Data: []any{"bad", "Authentication failure"}}
		}
		return "OpaqueRef:0123456789abcdef-session", nil
	})
	defer srv.Close()
	host, port := split(t, srv.URL)
	key := HostKey(host, port)
	store := memStore{}
	p := &trust.Policy{Store: store, OperatorPins: map[string]string{}}

	// a self-signed certificate on first use needs the user
	_, err := Connect(context.Background(), p, host, port, "root", "s3cret")
	var untrusted *trust.UntrustedError
	if !errors.As(err, &untrusted) || untrusted.HostKey != key || untrusted.Previous != "" {
		t.Fatalf("want an untrusted certificate, got %v", err)
	}
	// accepted: it connects
	_ = store.Pin(key, untrusted.Fingerprint)
	conn, err := Connect(context.Background(), p, host, port, "root", "s3cret")
	if err != nil {
		t.Fatal(err)
	}
	_ = conn.Close()

	// a wrong password: a typed XAPI failure, without the password or the session reference in the text
	_, err = Connect(context.Background(), p, host, port, "bad", "s3cret")
	if err == nil || !strings.Contains(err.Error(), "SESSION_AUTHENTICATION_FAILED") || strings.Contains(err.Error(), "s3cret") {
		t.Fatalf("login failure: %v", err)
	}

	// the certificate changes: refused even though it was trusted before
	store[key] = "00"
	_, err = Connect(context.Background(), p, host, port, "root", "s3cret")
	if !errors.As(err, &untrusted) || untrusted.Previous != "00" {
		t.Fatalf("want a changed certificate, got %v", err)
	}
}

func TestCoordinatorRedirect(t *testing.T) {
	// one test server plays both: reached as 127.0.0.1 it is a pool member, as "localhost" the coordinator
	var srv *httptest.Server
	srv = fakeXAPIByHost(t, func(host string) (any, *xenapi.ResponseError) {
		if strings.HasPrefix(host, "localhost") {
			return "OpaqueRef:coordinator-session-0000", nil
		}
		return nil, &xenapi.ResponseError{Code: 1, Message: "HOST_IS_SLAVE", Data: []any{"localhost"}}
	})
	defer srv.Close()
	host, port := split(t, srv.URL)
	p := &trust.Policy{Store: memStore{}, OperatorPins: map[string]string{"*": trust.Fingerprint(srv.Certificate().Raw)}, RequirePin: true}
	conn, err := Connect(context.Background(), p, host, port, "root", "pw")
	if err != nil {
		t.Fatal(err)
	}
	if conn.Host != "localhost" {
		t.Fatalf("connected to %s, want the coordinator", conn.Host)
	}
}

func TestNoRedirects(t *testing.T) {
	srv := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "https://elsewhere.example/jsonrpc", http.StatusTemporaryRedirect)
	}))
	defer srv.Close()
	host, port := split(t, srv.URL)
	p := &trust.Policy{Store: memStore{}, OperatorPins: map[string]string{"*": trust.Fingerprint(srv.Certificate().Raw)}}
	_, err := Connect(context.Background(), p, host, port, "root", "pw")
	if err == nil {
		t.Fatal("a redirect must not log in")
	}
	if strings.Contains(err.Error(), "elsewhere") && !strings.Contains(err.Error(), "decode") && !strings.Contains(err.Error(), "Response missing") {
		t.Fatalf("followed the redirect: %v", err)
	}
}

func TestSanitize(t *testing.T) {
	err := Sanitize(errors.New("failed for OpaqueRef:0f3c2b1a-1111-2222-3333-444455556666 with pw hunter2"), "hunter2")
	if strings.Contains(err.Error(), "0f3c2b1a") || strings.Contains(err.Error(), "hunter2") {
		t.Error("session references are redacted")
	}
}
