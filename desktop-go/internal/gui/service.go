package gui

import (
	"context"
	"errors"
	"sort"
	"strings"
	"sync"

	"xcpng-center/internal/app"
	"xcpng-center/internal/config"
	"xcpng-center/internal/trust"
	"xcpng-center/internal/xapi"
)

// Workspace is the service the window calls (Wails binding): the server list, connecting with the certificate
// question, and the objects for the tree. It only adapts internal/app for the UI; passwords stay in memory for the
// session and are never written or returned.
type Workspace struct {
	settings *config.Settings
	policy   *trust.Policy
	mu       sync.Mutex
	conns    map[string]*xapi.Connection
	use      sync.RWMutex // listings hold it for reading; logging a session out holds it for writing
}

// NewWorkspace uses the settings file of the CLI (servers and pinned certificates are shared) and the CLI's trust
// policy (the operator's --trust-fingerprint applies to the window too).
func NewWorkspace(s *config.Settings, p *trust.Policy) *Workspace {
	return &Workspace{settings: s, policy: p, conns: map[string]*xapi.Connection{}}
}

// ServerRow is one server of the list.
type ServerRow struct {
	Key       string `json:"key"`
	Host      string `json:"host"`
	Port      int    `json:"port"`
	User      string `json:"user"`
	Name      string `json:"name"`
	Connected bool   `json:"connected"`
}

// Servers lists the configured servers.
func (w *Workspace) Servers() []ServerRow {
	servers := w.settings.ServerList()
	w.mu.Lock()
	defer w.mu.Unlock()
	rows := []ServerRow{}
	for _, e := range servers {
		key := xapi.HostKey(e.Host, e.Port)
		_, ok := w.conns[key]
		rows = append(rows, ServerRow{Key: key, Host: e.Host, Port: e.Port, User: e.User, Name: e.Name, Connected: ok})
	}
	sort.Slice(rows, func(i, j int) bool { return rows[i].Key < rows[j].Key })
	return rows
}

// ConnectResult tells the window what happened: connected, or a certificate the user must decide on.
type ConnectResult struct {
	Connected   bool   `json:"connected"`
	Key         string `json:"key,omitempty"`
	Fingerprint string `json:"fingerprint,omitempty"`
	Previous    string `json:"previous,omitempty"`
	Error       string `json:"error,omitempty"`
}

// Connect is Add New Server / Connect: it logs in, adds the server to the list, and reports an untrusted certificate
// instead of failing (the window shows the fingerprint and calls Trust, then Connect again).
func (w *Workspace) Connect(address, user, password string) ConnectResult {
	host, port, err := xapi.ParseAddress(address)
	if err != nil {
		return ConnectResult{Error: err.Error()}
	}
	if user == "" {
		user = "root"
	}
	conn, err := xapi.Connect(context.Background(), w.policy, host, port, user, password)
	var untrusted *trust.UntrustedError
	if errors.As(err, &untrusted) {
		return ConnectResult{Key: untrusted.HostKey, Fingerprint: untrusted.Fingerprint, Previous: untrusted.Previous}
	}
	if err != nil {
		return ConnectResult{Error: xapi.Redact(err.Error(), password)}
	}
	key := xapi.HostKey(conn.Host, conn.Port)
	w.mu.Lock()
	old := w.conns[key]
	w.conns[key] = conn
	w.mu.Unlock()
	w.logout(old)
	if err := w.settings.AddServer(config.Server{Host: conn.Host, Port: conn.Port, User: user}); err != nil {
		return ConnectResult{Connected: true, Key: key, Error: err.Error()}
	}
	return ConnectResult{Connected: true, Key: key}
}

// Trust pins a certificate the user accepted for a new server. A changed certificate is never replaced from here:
// the C# app also refuses it until the user removes the old pin deliberately.
func (w *Workspace) Trust(key, fingerprint string) error {
	if want, ok := w.policy.Expected(key); ok && !strings.EqualFold(want, trust.Normalize(fingerprint)) {
		return errors.New("this certificate is not the one expected for this server: it is not trusted from the window")
	}
	return w.settings.PinNew(key, trust.Normalize(fingerprint))
}

// Disconnect logs out of a server.
func (w *Workspace) Disconnect(key string) {
	w.mu.Lock()
	conn := w.conns[key]
	delete(w.conns, key)
	w.mu.Unlock()
	w.logout(conn)
}

// logout waits for the listings in progress (they use the session) before logging it out.
func (w *Workspace) logout(conn *xapi.Connection) {
	if conn == nil {
		return
	}
	w.use.Lock()
	defer w.use.Unlock()
	_ = conn.Close()
}

// Objects lists one kind (pools, hosts, vms, templates, srs, networks) of a connected server.
func (w *Workspace) Objects(key, kind string) ([]app.Object, error) {
	w.use.RLock()
	defer w.use.RUnlock()
	w.mu.Lock()
	conn := w.conns[key]
	w.mu.Unlock()
	if conn == nil {
		return nil, errors.New("not connected")
	}
	return app.ListObjects(conn, kind)
}

// Close logs out of every server (the window closes).
func (w *Workspace) Close() {
	w.mu.Lock()
	conns := w.conns
	w.conns = map[string]*xapi.Connection{}
	w.mu.Unlock()
	for _, c := range conns {
		w.logout(c)
	}
}
