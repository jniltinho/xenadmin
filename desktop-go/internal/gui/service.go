package gui

import (
	"context"
	"errors"
	"fmt"
	"net"
	"sort"
	"strings"
	"sync"

	"xcpng-center/internal/app"
	"xcpng-center/internal/config"
	"xcpng-center/internal/console"
	"xcpng-center/internal/messages"
	"xcpng-center/internal/ops"
	"xcpng-center/internal/rules"
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
	console  *console.Proxy
	closed   bool // the window closed: no console endpoint is created any more
}

// webviewOrigins are the pages Wails serves the window from (internal/assetserver: Linux and macOS, Windows).
var webviewOrigins = []string{"wails://localhost", "http://wails.localhost"}

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
	// the name stays in the list, so a disconnected server keeps it in the tree (as the C# app's saved servers)
	if err := w.settings.AddServer(config.Server{Host: conn.Host, Port: conn.Port, User: user, Name: app.PoolName(conn)}); err != nil {
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

// Forget disconnects a server and removes it from the list (the certificate pin stays, as in the C# app).
func (w *Workspace) Forget(key string) error {
	host, port, err := xapi.ParseAddress(key)
	if err != nil {
		return err
	}
	w.Disconnect(key)
	return w.settings.RemoveServer(host, port)
}

// Properties reads a VM for its Properties dialog.
func (w *Workspace) Properties(key, vmUUID string) (out app.VMProperties, err error) {
	err = w.with(key, func(c *xapi.Connection) error { out, err = app.ReadVMProperties(c, vmUUID); return err })
	return out, err
}

// SaveProperties applies the Properties dialog's changes.
func (w *Workspace) SaveProperties(key, vmUUID string, changes ops.VMChanges) error {
	return w.with(key, func(c *xapi.Connection) error { return ops.SaveVMProperties(c, vmUUID, changes) })
}

// Disconnect logs out of a server. In-flight calls on it finish first, so no console ticket is issued after the
// revocation.
func (w *Workspace) Disconnect(key string) {
	w.mu.Lock()
	conn := w.conns[key]
	delete(w.conns, key)
	w.mu.Unlock()
	w.use.Lock()
	defer w.use.Unlock()
	w.mu.Lock()
	proxy := w.console // read after the in-flight calls ended: a first OpenConsole may have just created it
	w.mu.Unlock()
	if proxy != nil {
		proxy.Revoke(key)
	}
	if conn != nil {
		_ = conn.Close()
	}
}

// logout waits for the calls in progress (they use the session) before logging it out.
func (w *Workspace) logout(conn *xapi.Connection) {
	if conn == nil {
		return
	}
	w.use.Lock()
	defer w.use.Unlock()
	_ = conn.Close()
}

// Objects lists one kind (pools, hosts, vms, templates, srs, networks) of a connected server.
func (w *Workspace) Objects(key, kind string) (out []app.Object, err error) {
	err = w.with(key, func(c *xapi.Connection) error { out, err = app.ListObjects(c, kind); return err })
	return out, err
}

// Alerts lists the alerts of a connected server (Notifications > Alerts), newest first.
func (w *Workspace) Alerts(key string) (out []app.Alert, err error) {
	err = w.with(key, func(c *xapi.Connection) error { out, err = app.ListAlerts(c); return err })
	return out, err
}

// DismissAlerts removes alerts from a server for good (the window asked DismissConfirmation first).
func (w *Workspace) DismissAlerts(key string, uuids []string) error {
	return w.with(key, func(c *xapi.Connection) error { return app.DismissAlerts(c, uuids) })
}

// DismissConfirmation is the question of Dismiss (one alert), Dismiss Selected and Dismiss All (AlertSummaryPage).
func (w *Workspace) DismissConfirmation(count int, all bool) Confirm {
	text := messages.Get("ALERT_DISMISS_SELECTED_CONFIRM")
	switch {
	case all:
		text = messages.Get("ALERT_DISMISS_ALL_NO_FILTER_CONTINUE")
	case count == 1:
		text = messages.Get("ALERT_DISMISS_CONFIRM")
	}
	return Confirm{Title: "Dismiss Alerts", Text: strings.TrimSpace(text)}
}

// with runs fn on a connected server's session; logging out waits for it (Workspace.use).
func (w *Workspace) with(key string, fn func(*xapi.Connection) error) error {
	w.use.RLock()
	defer w.use.RUnlock()
	w.mu.Lock()
	conn := w.conns[key]
	w.mu.Unlock()
	if conn == nil {
		return errors.New("not connected")
	}
	return fn(conn)
}

// Confirm is the question a VM operation asks first (empty Title: none).
type Confirm struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}

// VMConfirmation is the confirmation of a VM operation (VmOperations.Confirmation).
func (w *Workspace) VMConfirmation(op string, haProtected bool) Confirm {
	if !rules.NeedsConfirmation(rules.VMOp(op)) {
		return Confirm{}
	}
	title, text := rules.Confirmation(rules.VMOp(op), haProtected)
	return Confirm{Title: title, Text: text}
}

// VMPower runs a VM power operation (the window asked VMConfirmation first).
func (w *Workspace) VMPower(key, uuid, op string) error {
	return w.with(key, func(c *xapi.Connection) error { return ops.VMPower(c, uuid, rules.VMOp(op)) })
}

// HostQuestion is what the window asks before a server command: Error (it cannot run, say why), else Title and
// Text (the confirmation); for Exit Maintenance Mode, Text is the restore question ("" when no VM can go back).
type HostQuestion struct {
	Error string `json:"error,omitempty"`
	Title string `json:"title,omitempty"`
	Text  string `json:"text,omitempty"`
	// Option is a check box shown with the question (Revert: take a snapshot first), "" for none
	Option string `json:"option,omitempty"`
}

// HostQuestion for op "shutdown", "reboot", "enter_maintenance" or "exit_maintenance".
func (w *Workspace) HostQuestion(key, uuid, op string) (q HostQuestion, err error) {
	err = w.with(key, func(c *xapi.Connection) error {
		if op == "exit_maintenance" {
			name, restore, err := ops.VMsToRestore(c, uuid)
			q = HostQuestion{Title: messages.Get("EXIT_MAINTENANCE_MODE_CONTEXT_MENU"), Text: rules.ExitMaintenanceQuestion(name, restore)}
			return err
		}
		h, err := ops.HostInfo(c, uuid)
		if err != nil {
			return err
		}
		switch op {
		case "shutdown", "reboot":
			q.Title, q.Text = rules.HostConfirmation(h, op == "reboot")
		case "enter_maintenance":
			if q.Error = rules.CannotEnterMaintenance(h); q.Error == "" {
				q.Title, q.Text = rules.EvacuateConfirmation(h)
			}
		default:
			return fmt.Errorf("unknown server command %q", op)
		}
		return nil
	})
	q.Title = strings.ReplaceAll(strings.ReplaceAll(q.Title, "&", ""), "...", "")
	return q, err
}

// HostCommand runs a server command after its question (restore: Exit Maintenance Mode puts the VMs back).
func (w *Workspace) HostCommand(key, uuid, op string, restore bool) error {
	return w.with(key, func(c *xapi.Connection) error {
		switch op {
		case "shutdown", "reboot":
			return ops.HostPower(c, uuid, op == "reboot")
		case "enter_maintenance":
			return ops.EnterMaintenance(c, uuid)
		case "exit_maintenance":
			return ops.ExitMaintenance(c, uuid, restore)
		}
		return fmt.Errorf("unknown server command %q", op)
	})
}

// Snapshots is the Snapshots tab of a VM.
func (w *Workspace) Snapshots(key, vmUUID string) (out app.VMSnapshots, err error) {
	err = w.with(key, func(c *xapi.Connection) error { out, err = app.ListSnapshots(c, vmUUID); return err })
	return out, err
}

// SnapshotQuestion is the question of Revert ("revert", with the snapshot-first option) or Delete ("delete").
func (w *Workspace) SnapshotQuestion(op, name string, snapshotAllowed bool) HostQuestion {
	if op == "revert" {
		text, option := rules.RevertQuestion(name, snapshotAllowed)
		return HostQuestion{Title: "XCP-ng Center", Text: text, Option: option}
	}
	return HostQuestion{Title: "XCP-ng Center", Text: rules.DeleteSnapshotQuestion(name)}
}

// TakeSnapshot takes a disk snapshot, or a disk and memory one.
func (w *Workspace) TakeSnapshot(key, vmUUID, name, description string, memory bool) error {
	return w.with(key, func(c *xapi.Connection) error { return ops.TakeSnapshot(c, vmUUID, name, description, memory) })
}

// RevertSnapshot reverts the VM to a snapshot (the window asked SnapshotQuestion first).
func (w *Workspace) RevertSnapshot(key, snapshotUUID string) error {
	return w.with(key, func(c *xapi.Connection) error { return ops.RevertSnapshot(c, snapshotUUID) })
}

// DeleteSnapshot deletes a snapshot and its disks (the window asked SnapshotQuestion first).
func (w *Workspace) DeleteSnapshot(key, snapshotUUID string) error {
	return w.with(key, func(c *xapi.Connection) error { return ops.DeleteSnapshot(c, snapshotUUID) })
}

// Storage is the Storage tab of a VM.
func (w *Workspace) Storage(key, vmUUID string) (out app.VMStorage, err error) {
	err = w.with(key, func(c *xapi.Connection) error { out, err = app.ListVMStorage(c, vmUUID); return err })
	return out, err
}

// Networking is the Networking tab of a VM.
func (w *Workspace) Networking(key, vmUUID string) (out app.Table, err error) {
	err = w.with(key, func(c *xapi.Connection) error { out, err = app.ListVMNetworking(c, vmUUID); return err })
	return out, err
}

// Memory is the Memory tab of a VM.
func (w *Workspace) Memory(key, vmUUID string) (out app.Table, err error) {
	err = w.with(key, func(c *xapi.Connection) error { out, err = app.ListVMMemory(c, vmUUID); return err })
	return out, err
}

// SetDiskActive activates or deactivates a VM's disk.
func (w *Workspace) SetDiskActive(key, vbdUUID string, active bool) error {
	return w.with(key, func(c *xapi.Connection) error { return ops.SetDiskActive(c, vbdUUID, active) })
}

// NewDiskOptions is what the Add Virtual Disk dialog offers.
func (w *Workspace) NewDiskOptions(key, vmUUID string) (out app.NewDisk, err error) {
	err = w.with(key, func(c *xapi.Connection) error { out, err = app.NewDiskOptions(c, vmUUID); return err })
	return out, err
}

// AddDisk creates a disk and attaches it to the VM; the notice says when the VM needs a restart first.
func (w *Workspace) AddDisk(key, vmUUID, name, description, srUUID string, size int) (notice string, err error) {
	err = w.with(key, func(c *xapi.Connection) error {
		notice, err = ops.AddDisk(c, vmUUID, name, description, srUUID, size)
		return err
	})
	return notice, err
}

// AttachableDisks lists the disks the Attach Disk dialog offers.
func (w *Workspace) AttachableDisks(key, vmUUID string) (out []app.ISO, err error) {
	err = w.with(key, func(c *xapi.Connection) error { out, err = app.AttachableDisks(c, vmUUID); return err })
	return out, err
}

// AttachDisk attaches an existing disk to the VM; the notice says when the VM needs a restart first.
func (w *Workspace) AttachDisk(key, vmUUID, vdiUUID string, readOnly bool) (notice string, err error) {
	err = w.with(key, func(c *xapi.Connection) error { notice, err = ops.AttachDisk(c, vmUUID, vdiUUID, readOnly); return err })
	return notice, err
}

// DiskQuestion is the Detach or Delete question of a disk.
func (w *Workspace) DiskQuestion(del, system bool) Confirm {
	title, text := rules.DiskQuestion(del, system)
	return Confirm{Title: title, Text: text}
}

// DetachDisk detaches a disk from its VM (the disk stays).
func (w *Workspace) DetachDisk(key, vbdUUID string) error {
	return w.with(key, func(c *xapi.Connection) error { return ops.DetachDisk(c, vbdUUID) })
}

// DeleteDisk deletes the disk of a VM's VBD (re-checked on the server).
func (w *Workspace) DeleteDisk(key, vmUUID, vbdUUID string) error {
	return w.with(key, func(c *xapi.Connection) error { return ops.DeleteDisk(c, vmUUID, vbdUUID) })
}

// CreateDVDDrive adds a DVD drive to a VM that has none; the notice says when it comes only after a reboot.
func (w *Workspace) CreateDVDDrive(key, vmUUID string) (notice string, err error) {
	err = w.with(key, func(c *xapi.Connection) error { notice, err = ops.CreateDVDDrive(c, vmUUID); return err })
	return notice, err
}

// ChangeISO puts an ISO in the DVD drive ("" ejects it).
func (w *Workspace) ChangeISO(key, dvdUUID, isoUUID string) error {
	return w.with(key, func(c *xapi.Connection) error { return ops.ChangeISO(c, dvdUUID, isoUUID) })
}

// Performance is the Performance tab of a VM or a server ("vms" or "hosts") for a zoom (10m, 1h, 1d, 1w, 1mo, 1y).
func (w *Workspace) Performance(key, kind, uuid, zoom string) (out app.Performance, err error) {
	err = w.with(key, func(c *xapi.Connection) error { out, err = app.ReadPerformance(c, w.policy, kind, uuid, zoom); return err })
	return out, err
}

// OpenConsole returns a one-time ws:// URL for the VM's graphical console (noVNC connects to it within 30 s); it is
// revoked, and an open console closed, when the server is disconnected.
func (w *Workspace) OpenConsole(key, vmUUID string) (ticket string, err error) {
	err = w.with(key, func(c *xapi.Connection) error {
		target, err := console.Locate(c, vmUUID)
		if err != nil {
			return err
		}
		w.mu.Lock()
		if w.closed {
			w.mu.Unlock()
			return errors.New("the window is closing")
		}
		if w.console == nil {
			if w.console, err = console.NewProxy(webviewOrigins...); err != nil {
				w.mu.Unlock()
				return err
			}
		}
		proxy := w.console
		w.mu.Unlock()
		session := string(c.Session.Ref())
		ticket, err = proxy.Ticket(key, func(ctx context.Context) (net.Conn, error) {
			return console.Dial(ctx, w.policy, target, session)
		})
		return err
	})
	return ticket, err
}

// Close logs out of every server (the window closes): calls in progress finish first, then the consoles stop.
func (w *Workspace) Close() {
	w.mu.Lock()
	w.closed = true
	conns := w.conns
	w.conns = map[string]*xapi.Connection{}
	w.mu.Unlock()
	w.use.Lock()
	defer w.use.Unlock()
	w.mu.Lock()
	proxy := w.console
	w.mu.Unlock()
	if proxy != nil {
		proxy.Close()
	}
	for _, c := range conns {
		_ = c.Close()
	}
}
