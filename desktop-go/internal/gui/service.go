package gui

import (
	"context"
	"errors"
	"fmt"
	xenapi "go/xenapi"
	"net"
	"sort"
	"strings"
	"sync"
	"time"

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
	closed   bool                // the window closed: no console endpoint is created any more
	watchers map[string]*watcher // a connection's event.from loop (app.Watch)
	// the password of each connection, as the C# Connection.Password: only for pool.join (New Pool), never returned
	passwords map[string]string
	pending   map[string]bool // a change notice is scheduled for this server
	// Emit tells the window that a server's objects changed (the Wails event "xcp:changed" with the server's key)
	Emit func(key string)
	// Files is the window's file dialog (export and import)
	Files FileDialog
	// hostBusy counts the server commands running on each host (by uuid)
	hostBusy map[string]int
}

// watcher is a running app.Watch: cancel cuts its call, done closes when it has returned.
type watcher struct {
	cancel context.CancelFunc
	done   chan struct{}
}

// stop ends the loop and waits for it, so that no call of it is left when the session is logged out.
func (wt *watcher) stop() {
	if wt != nil {
		wt.cancel()
		<-wt.done
	}
}

// changeDelay coalesces the event batches of one operation (a VM start changes the VM, its metrics, VBDs, VIFs...)
// into one reload of the window.
const changeDelay = 700 * time.Millisecond

// webviewOrigins are the pages Wails serves the window from (internal/assetserver: Linux and macOS, Windows).
var webviewOrigins = []string{"wails://localhost", "http://wails.localhost"}

// NewWorkspace uses the settings file of the CLI (servers and pinned certificates are shared) and the CLI's trust
// policy (the operator's --trust-fingerprint applies to the window too).
func NewWorkspace(s *config.Settings, p *trust.Policy) *Workspace {
	applyNetwork(s.GetOptions())
	return &Workspace{settings: s, policy: p, conns: map[string]*xapi.Connection{}, watchers: map[string]*watcher{}, pending: map[string]bool{}, passwords: map[string]string{}}
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
	// Tools > Options: "Warn me when a new certificate is found" off: a first certificate is pinned without asking, for
	// the server and then the coordinator it redirects to; a changed certificate always asks
	for tries := 0; tries < 3 && errors.As(err, &untrusted) && untrusted.Previous == "" && w.settings.GetOptions().NoNewCertificateWarning; tries++ {
		if perr := w.settings.PinNew(untrusted.HostKey, trust.Normalize(untrusted.Fingerprint)); perr != nil {
			break
		}
		conn, err = xapi.Connect(context.Background(), w.policy, host, port, user, password)
	}
	if errors.As(err, &untrusted) {
		return ConnectResult{Key: untrusted.HostKey, Fingerprint: untrusted.Fingerprint, Previous: untrusted.Previous}
	}
	if err != nil {
		return ConnectResult{Error: xapi.Redact(err.Error(), password)}
	}
	key := xapi.HostKey(conn.Host, conn.Port)
	ctx, cancel := context.WithCancel(context.Background())
	wt := &watcher{cancel: cancel, done: make(chan struct{})}
	session := conn.Session.WithContext(ctx) // its own copy: logging conn out never races with the loop
	w.mu.Lock()
	old, oldWatch := w.conns[key], w.watchers[key]
	w.conns[key], w.watchers[key], w.passwords[key] = conn, wt, password
	w.mu.Unlock()
	oldWatch.stop() // the old session's loop, before that session is logged out
	w.logout(old)
	go func() { defer close(wt.done); app.Watch(ctx, session, func() { w.changed(key) }) }()
	// the name stays in the list, so a disconnected server keeps it in the tree (as the C# app's saved servers)
	if err := w.settings.AddServer(config.Server{Host: conn.Host, Port: conn.Port, User: user, Name: app.PoolName(conn)}); err != nil {
		return ConnectResult{Connected: true, Key: key, Error: err.Error()}
	}
	return ConnectResult{Connected: true, Key: key}
}

// reconnect logs in again with the connection's password kept in memory, trying for three minutes while the
// coordinator's toolstack starts.
// It stops when the connection is no longer the one restarted (Disconnect, Forget, a new login) or the window closed.
func (w *Workspace) reconnect(key, user string) error {
	w.mu.Lock()
	password, old := w.passwords[key], w.conns[key]
	w.mu.Unlock()
	// the agent restarts after the call returns: wait for the old session to fail first (a login before that would die
	// with it), for at most a minute
	for deadline := time.Now().Add(time.Minute); old != nil && time.Now().Before(deadline); time.Sleep(2 * time.Second) {
		w.use.RLock() // a Disconnect or Close logging the session out waits for this call
		w.mu.Lock()
		gone := w.closed || w.conns[key] != old
		w.mu.Unlock()
		var err error
		if !gone {
			_, err = xenapi.Pool.GetAll(old.Session)
		}
		w.use.RUnlock()
		if gone {
			return nil
		}
		if err != nil {
			break
		}
	}
	var last string
	for deadline := time.Now().Add(3 * time.Minute); time.Now().Before(deadline); time.Sleep(5 * time.Second) {
		w.mu.Lock()
		current, closed := w.conns[key], w.closed
		w.mu.Unlock()
		if closed || current != old {
			return nil
		}
		r := w.Connect(key, user, password)
		if r.Connected {
			w.changed(key)
			return nil
		}
		last = r.Error
	}
	return errors.New("the server did not come back after the toolstack restart: " + last)
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

// MemoryEdit reads a VM for its Memory Settings dialog.
func (w *Workspace) MemoryEdit(key, vmUUID string) (out app.MemoryEdit, err error) {
	err = w.with(key, func(c *xapi.Connection) error { out, err = app.ReadMemoryEdit(c, vmUUID); return err })
	return out, err
}

// MemoryQuestion checks the Memory Settings values: an error when they cannot be set, else the question to ask first.
func (w *Workspace) MemoryQuestion(key, vmUUID string, dynMin, dynMax, staticMax int64) (out string, err error) {
	err = w.with(key, func(c *xapi.Connection) error {
		out, err = app.MemoryQuestion(c, vmUUID, dynMin, dynMax, staticMax)
		return err
	})
	return out, err
}

// ChangeMemory applies the Memory Settings dialog (a new static maximum restarts a running VM).
func (w *Workspace) ChangeMemory(key, vmUUID string, dynMin, dynMax, staticMax int64) error {
	return w.with(key, func(c *xapi.Connection) error { return app.ChangeMemory(c, vmUUID, dynMin, dynMax, staticMax) })
}

// CustomFields are the Custom Fields page of a VM: the pool's fields and the VM's values.
func (w *Workspace) CustomFields(key, vmUUID string) (out []app.FieldValue, err error) {
	err = w.with(key, func(c *xapi.Connection) error { out, err = app.VMCustomFields(c, vmUUID); return err })
	return out, err
}

// SaveCustomFields sets the VM's custom field values ("" removes one).
func (w *Workspace) SaveCustomFields(key, vmUUID string, values map[string]string) error {
	return w.with(key, func(c *xapi.Connection) error { return ops.SaveVMCustomFields(c, vmUUID, values) })
}

// AddCustomField adds a field (kind "String" or "Date") to the pool.
func (w *Workspace) AddCustomField(key, name, kind string) error {
	return w.with(key, func(c *xapi.Connection) error { return ops.AddCustomField(c, name, kind) })
}

// RemoveCustomField deletes a field of the pool and its values on every object.
func (w *Workspace) RemoveCustomField(key, name string) error {
	return w.with(key, func(c *xapi.Connection) error { return ops.RemoveCustomField(c, name) })
}

// HostTab is the Memory, Storage, Networking or NICs tab of a server or pool ("hosts", "pools"), or an SR's Storage tab.
func (w *Workspace) HostTab(key, kind, uuid, tab string) (out app.Table, err error) {
	err = w.with(key, func(c *xapi.Connection) error { out, err = app.HostTab(c, kind, uuid, tab); return err })
	return out, err
}

// ConvertQuestion is Convert to Template's question (ConvertVMToTemplateCommand).
func (w *Workspace) ConvertQuestion(name string) Confirm {
	return Confirm{Title: messages.Get("CONVERT_TO_TEMPLATE"), Text: messages.Format(messages.Get("CONVERT_TEMPLATE_DIALOG_TEXT"), ops.Ellipsise(name, 25))}
}

// ConvertToTemplate turns a VM into a template.
func (w *Workspace) ConvertToTemplate(key, vmUUID string) error {
	return w.with(key, func(c *xapi.Connection) error { return ops.ConvertToTemplate(c, vmUUID) })
}

// Overview is the live columns of the Search overview for a connection's servers and VMs (by uuid).
func (w *Workspace) Overview(key string) (out map[string]app.OverviewRow, err error) {
	err = w.with(key, func(c *xapi.Connection) error { out, err = app.Overview(c, w.policy); return err })
	return out, err
}

// Disconnect logs out of a server. In-flight calls on it finish first, so no console ticket is issued after the
// revocation.
func (w *Workspace) Disconnect(key string) {
	w.mu.Lock()
	conn, wt := w.conns[key], w.watchers[key]
	delete(w.conns, key)
	delete(w.watchers, key)
	delete(w.passwords, key)
	w.mu.Unlock()
	wt.stop()
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

// changed schedules one change notice for a server: the batches that come within changeDelay make one reload.
func (w *Workspace) changed(key string) {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.pending[key] {
		return
	}
	w.pending[key] = true
	time.AfterFunc(changeDelay, func() {
		w.mu.Lock()
		delete(w.pending, key)
		_, connected := w.conns[key]
		emit := w.Emit
		w.mu.Unlock()
		if connected && emit != nil {
			emit(key)
		}
	})
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

var errNotConnected = errors.New("not connected")

// with runs fn on a connected server's session; logging out waits for it (Workspace.use).
func (w *Workspace) with(key string, fn func(*xapi.Connection) error) error {
	w.use.RLock()
	defer w.use.RUnlock()
	w.mu.Lock()
	conn := w.conns[key]
	w.mu.Unlock()
	if conn == nil {
		return errNotConnected
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
	// Option is a check box shown with the question (Revert: take a snapshot first), or the Yes button's label of
	// Remove From Pool; "" for none
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
		switch op {
		case "restart_toolstack":
			q.Title, q.Text, err = ops.RestartToolstackQuestion(c, uuid)
			return err
		case "remove_from_pool":
			q.Title, q.Text, q.Option, err = ops.RemoveFromPoolQuestion(c, uuid)
			if err != nil {
				q.Error, err = err.Error(), nil
			}
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

// HostCommand runs a server command after its question (restore: Exit Maintenance Mode puts the VMs back). Restart
// Toolstack waits for no other command of ours on that server (HelpersGUI.HasActiveHostAction), and on the
// coordinator connects again afterwards (RestartToolstackAction: the session does not survive the restart).
func (w *Workspace) HostCommand(key, uuid, op string, restore bool) error {
	w.mu.Lock()
	if w.hostBusy == nil {
		w.hostBusy = map[string]int{}
	}
	if op == "restart_toolstack" && w.hostBusy[uuid] > 0 {
		w.mu.Unlock()
		return errors.New("another command is running on this server: wait for it to finish")
	}
	w.hostBusy[uuid]++
	w.mu.Unlock()
	defer func() { w.mu.Lock(); w.hostBusy[uuid]--; w.mu.Unlock() }()
	if op == "restart_toolstack" {
		var coordinator bool
		var user string
		err := w.with(key, func(c *xapi.Connection) error {
			user = c.User
			h, err := ops.HostInfo(c, uuid)
			if err != nil { // never restart without knowing whether the session goes with it
				return err
			}
			coordinator = h.Coordinator
			return ops.RestartToolstack(c, uuid)
		})
		if err != nil || !coordinator {
			return err
		}
		return w.reconnect(key, user)
	}
	return w.with(key, func(c *xapi.Connection) error {
		switch op {
		case "shutdown", "reboot":
			return ops.HostPower(c, uuid, op == "reboot")
		case "enter_maintenance":
			return ops.EnterMaintenance(c, uuid)
		case "exit_maintenance":
			return ops.ExitMaintenance(c, uuid, restore)
		case "restart_toolstack":
			return ops.RestartToolstack(c, uuid)
		case "remove_from_pool":
			return ops.RemoveFromPool(c, uuid)
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
func (w *Workspace) Networking(key, vmUUID string) (out app.VMNetworking, err error) {
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

// VIFNetworks are the networks the Add Virtual Interface dialog offers for a VM.
func (w *Workspace) VIFNetworks(key, vmUUID string) (out []app.NetworkPick, err error) {
	err = w.with(key, func(c *xapi.Connection) error { out, err = app.VIFNetworks(c, vmUUID); return err })
	return out, err
}

// VIFQuestion is VIFDialog_FormClosing's question when a real VM of any connected server already has this MAC (an
// empty Confirm when none).
func (w *Workspace) VIFQuestion(vifUUID, mac string) (Confirm, error) {
	for _, key := range w.connectedKeys() {
		var vm string
		err := w.with(key, func(c *xapi.Connection) (err error) { vm, err = app.DuplicateMAC(c, vifUUID, mac); return err })
		if err != nil {
			return Confirm{}, err
		}
		if vm != "" {
			return Confirm{Title: messages.Get("PROBLEM_MAC_ADDRESS_IS_DUPLICATE_TITLE"),
				Text: messages.Format(messages.Get("PROBLEM_MAC_ADDRESS_IS_DUPLICATE"), mac, vm)}, nil
		}
	}
	return Confirm{}, nil
}

// VIFDeleteQuestion is Remove Interface's question.
func (w *Workspace) VIFDeleteQuestion() Confirm {
	return Confirm{Title: messages.Get("MESSAGEBOX_VIF_DELETE_TITLE"), Text: messages.Get("MESSAGEBOX_VIF_DELETE")}
}

// SaveVIF adds an interface to a VM (vifUUID "") or changes one; the text when the change waits for a restart
// (VIF_HOTPLUG_FAILED_MESSAGE), else "".
func (w *Workspace) SaveVIF(key, vmUUID, vifUUID string, set ops.VIFSettings) (string, error) {
	var reboot bool
	err := w.with(key, func(c *xapi.Connection) (err error) {
		if vifUUID == "" {
			reboot, err = ops.CreateVIF(c, vmUUID, set)
		} else {
			reboot, err = ops.UpdateVIF(c, vifUUID, set)
		}
		return err
	})
	if reboot {
		return messages.Get("VIF_HOTPLUG_FAILED_MESSAGE"), err
	}
	return "", err
}

// DeleteVIF removes a VM's interface.
func (w *Workspace) DeleteVIF(key, vifUUID string) error {
	return w.with(key, func(c *xapi.Connection) error { return ops.DeleteVIF(c, vifUUID) })
}

// SetVIFActive activates or deactivates a running VM's interface.
func (w *Workspace) SetVIFActive(key, vifUUID string, active bool) error {
	return w.with(key, func(c *xapi.Connection) error { return ops.SetVIFActive(c, vifUUID, active) })
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
	err = w.with(key, func(c *xapi.Connection) error {
		out, err = app.ReadPerformance(c, w.policy, kind, uuid, zoom)
		return err
	})
	return out, err
}

// GraphDataSources are the data sources Graph Details offers for a VM or a server.
func (w *Workspace) GraphDataSources(key, kind, uuid string) (out []app.DataSource, err error) {
	err = w.with(key, func(c *xapi.Connection) error { out, err = app.GraphDataSources(c, kind, uuid); return err })
	return out, err
}

// SaveGraphs saves the Performance graphs of a VM or a server (none: the default graphs) made from base (the layout
// shown) and data source colours.
func (w *Workspace) SaveGraphs(key, kind, uuid string, base, graphs []rules.DesignedGraph, colours map[string]string) error {
	return w.with(key, func(c *xapi.Connection) error { return ops.SaveGraphs(c, kind, uuid, base, graphs, colours) })
}

// GraphTexts are the Performance tab's dialog titles and questions about the graph named name (GraphDetailsDialog,
// PerformancePage.DeleteGraph and RestoreDefaultGraphs).
func (w *Workspace) GraphTexts(name string) map[string]string {
	return map[string]string{"new_title": messages.Get("GRAPHS_NEW_TITLE"), "edit_title": messages.Format(messages.Get("GRAPHS_EDIT_TITLE"), name),
		"delete": messages.Format(messages.Get("DELETE_GRAPH_MESSAGE"), name), "restore": strings.TrimSpace(messages.Get("GRAPHS_RESTORE_DEFAULT_MESSAGE"))}
}

// EnableDataSource turns recording on for a data source of a VM or a server (text: its name, for the error).
func (w *Workspace) EnableDataSource(key, kind, uuid, name, text string) error {
	return w.with(key, func(c *xapi.Connection) error { return ops.EnableDataSource(c, kind, uuid, name, text) })
}

// NewVMOptions is what the New VM wizard offers.
func (w *Workspace) NewVMOptions(key string) (out app.NewVMOptions, err error) {
	err = w.with(key, func(c *xapi.Connection) error { out, err = app.ReadNewVMOptions(c); return err })
	return out, err
}

// NewVMDefaults are the wizard's defaults for a template and an optional home server.
func (w *Workspace) NewVMDefaults(key, templateUUID, homeUUID string) (out app.NewVMDefaults, err error) {
	err = w.with(key, func(c *xapi.Connection) error {
		out, err = app.ReadNewVMDefaults(c, templateUUID, homeUUID)
		return err
	})
	return out, err
}

// CreateVM creates (and optionally starts) the VM the wizard describes: the new VM's uuid, and the start error when
// only starting it failed (the VM exists then; the wizard must not create it again).
func (w *Workspace) CreateVM(key string, spec ops.NewVM) (out ops.Created, err error) {
	err = w.with(key, func(c *xapi.Connection) error { out, err = ops.CreateVM(c, spec); return err })
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
	watchers := w.watchers
	w.watchers = map[string]*watcher{}
	w.passwords = map[string]string{}
	w.mu.Unlock()
	for _, wt := range watchers {
		wt.stop()
	}
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

// NewSRKinds is the New Storage wizard's Type page.
func (w *Workspace) NewSRKinds(key string) (out []app.NewSRKind, err error) {
	err = w.with(key, func(c *xapi.Connection) error { out, err = app.ReadNewSRKinds(c); return err })
	return out, err
}

// CheckSRLocation is the Location page's check before Finish: "" when the share can be used.
func (w *Workspace) CheckSRLocation(key, kind, path string) (problem string, err error) {
	if !rules.SRValidLocation(rules.SRKind(kind), path) {
		if kind == string(rules.SRNfs) || kind == string(rules.SRNfsIso) {
			return "Enter the share as server:/path.", nil
		}
		return `Enter the share as \\server\share.`, nil
	}
	err = w.with(key, func(c *xapi.Connection) error {
		if e := ops.ISOAlreadyAttached(c, rules.SRKind(kind), path); e != nil {
			problem = e.Error()
		}
		return nil
	})
	return problem, err
}

// CreateSR creates the SR the wizard describes and returns its uuid.
func (w *Workspace) CreateSR(key string, spec ops.NewSR) (uuid string, err error) {
	if spec.Kind == rules.SRIscsi && spec.Existing == "" { // Format: never a LUN an SR of any connected pool uses
		if text, err := w.LUNInUse(spec.LUN.ScsiID); err != nil {
			return "", err
		} else if text != "" {
			return "", errors.New(text)
		}
	}
	if spec.Existing != "" {
		if pool, err := w.srInUseElsewhere(key, spec.Existing); err != nil {
			return "", err
		} else if pool != "" {
			return "", fmt.Errorf("this SR is already in use by %s", pool)
		}
	}
	err = w.with(key, func(c *xapi.Connection) error { uuid, err = ops.CreateSR(c, spec); return err })
	return uuid, err
}

// srInUseElsewhere is SrWizardHelpers.SrInUse over the other connected pools: the name of one that has the SR attached
// ("" for none). A pool that cannot answer stops the attach (two pools must never share an SR).
func (w *Workspace) srInUseElsewhere(key, uuid string) (string, error) {
	w.mu.Lock()
	var keys []string
	for k := range w.conns {
		if k != key {
			keys = append(keys, k)
		}
	}
	w.mu.Unlock()
	for _, k := range keys {
		var pool string
		err := w.with(k, func(c *xapi.Connection) error {
			used, err := ops.SRInUse(c, uuid)
			if used {
				pool = ops.PoolName(c)
			}
			return err
		})
		if err != nil && !errors.Is(err, errNotConnected) { // disconnected meanwhile: it holds nothing
			return "", fmt.Errorf("cannot check whether %s uses this SR: %w", k, err)
		}
		if pool != "" {
			return pool, nil
		}
	}
	return "", nil
}

// ProbeSR is the Location page's Scan: the SRs already on the share.
func (w *Workspace) ProbeSR(key string, spec ops.NewSR) (out []ops.FoundSR, err error) {
	if err = w.with(key, func(c *xapi.Connection) error { out, err = ops.ProbeSR(c, spec); return err }); err != nil {
		return nil, err
	}
	for i := range out {
		if out[i].InUse == "" {
			if out[i].InUse, err = w.srInUseElsewhere(key, out[i].UUID); err != nil {
				return nil, err
			}
		}
	}
	return out, nil
}

// SRAttachWarning is the question before reattaching an existing NFS SR ("" for none).
func (w *Workspace) SRAttachWarning(kind, uuid string) string {
	return ops.MultiPoolWarning(rules.SRKind(kind), uuid)
}

// SRQuestion is the confirmation of an SR command (an empty title: none, for Set as Default).
func (w *Workspace) SRQuestion(command, name string) HostQuestion {
	title, text, yes := rules.SRConfirmation(command, name)
	return HostQuestion{Title: title, Text: text, Option: yes}
}

// SRCommand runs an SR command of the Storage menu (set_default, detach, forget, destroy).
func (w *Workspace) SRCommand(key, uuid, command string) error {
	return w.with(key, func(c *xapi.Connection) error { return ops.SRCommand(c, uuid, command) })
}

// VMDeleteChoices fills the Delete VM dialog.
func (w *Workspace) VMDeleteChoices(key, uuid string) (out app.VMDelete, err error) {
	err = w.with(key, func(c *xapi.Connection) error { out, err = app.ReadVMDelete(c, uuid); return err })
	return out, err
}

// DeleteVM deletes the VM with the disks (VBD uuids) and snapshots chosen in the dialog.
func (w *Workspace) DeleteVM(key, uuid string, vbds, snapshots []string) error {
	return w.with(key, func(c *xapi.Connection) error { return ops.DeleteVM(c, uuid, vbds, snapshots) })
}

// VMCopyChoices fills the Copy VM dialog.
func (w *Workspace) VMCopyChoices(key, uuid string) (out app.VMCopy, err error) {
	err = w.with(key, func(c *xapi.Connection) error { out, err = app.ReadVMCopy(c, uuid); return err })
	return out, err
}

// CopyVM clones (or, with full, copies to the SR) a shut down VM and returns the copy's uuid.
func (w *Workspace) CopyVM(key, uuid, name, description string, full bool, sr string) (out string, err error) {
	err = w.with(key, func(c *xapi.Connection) error {
		out, err = ops.CopyVM(c, uuid, name, description, full, sr)
		return err
	})
	return out, err
}

// VMTargets are the servers of the Migrate/Start/Resume on Server submenus, the home server first.
func (w *Workspace) VMTargets(key, uuid string) (out []app.Target, err error) {
	err = w.with(key, func(c *xapi.Connection) error { out, err = app.ReadTargets(c, uuid); return err })
	return out, err
}

// VMOnHost migrates the VM to a server, or starts or resumes it there (op: migrate, start_on, resume_on).
func (w *Workspace) VMOnHost(key, uuid, host, op string) error {
	return w.with(key, func(c *xapi.Connection) error { return ops.VMOnHost(c, uuid, host, op) })
}

// ScanIQNs is the iSCSI page's Scan Target Host.
func (w *Workspace) ScanIQNs(key, target string, port int, chapUser, chapPassword string) (out []rules.IQN, err error) {
	err = w.with(key, func(c *xapi.Connection) error {
		out, err = ops.ScanIQNs(c, target, port, chapUser, chapPassword)
		return err
	})
	return out, err
}

// ScanLUNs lists the LUNs of the chosen IQN.
func (w *Workspace) ScanLUNs(key string, iqn rules.IQN, chapUser, chapPassword string) (out []rules.LUN, err error) {
	err = w.with(key, func(c *xapi.Connection) error { out, err = ops.ScanLUNs(c, iqn, chapUser, chapPassword); return err })
	return out, err
}

// LUNInUse is the iSCSI page's check before Finish (LVMoISCSI.SrUsingLunExists over every connected pool): the
// message when an SR already uses the LUN, "" otherwise.
func (w *Workspace) LUNInUse(scsiID string) (string, error) {
	w.mu.Lock()
	var keys []string
	for k := range w.conns {
		keys = append(keys, k)
	}
	w.mu.Unlock()
	for _, k := range keys {
		var text string
		err := w.with(k, func(c *xapi.Connection) error {
			sr, err := ops.LUNInUse(c, scsiID)
			if sr != "" {
				text = messages.Format(messages.Get("NEWSR_LUN_IN_USE_ON_POOL"), sr, ops.PoolName(c))
			}
			return err
		})
		if err != nil && !errors.Is(err, errNotConnected) {
			return "", err
		}
		if text != "" {
			return text, nil
		}
	}
	return "", nil
}

// LUNQuestion is LVMoISCSI's question after probing the LUN (LVMoIsciWarningDialog): Format only, when no SR is on
// it, else Reattach or Format.
func (w *Workspace) LUNQuestion(found bool) string {
	kind := messages.Friendly("Label-SR.SRTypes-lvmoiscsi")
	warning := "To prevent data loss you must ensure that the LUN is not in use by any other system, including XCP-ng hosts that are not connected to XCP-ng Center."
	format := messages.Format(messages.Get("LVMOHBA_WARNING_DIALOG_FORMAT_INFO"), kind)
	if !found {
		return messages.Get("LVMOHBA_WARNING_DIALOG_HEADER_NO_EXISTING_SRS") + "\n\n" + format + "\n\n" + warning
	}
	return messages.Format(messages.Get("LVMOHBA_WARNING_DIALOG_HEADER_FOUND_EXISTING_SR"), kind) + "\n\nReattach: " +
		messages.Format(messages.Get("LVMOHBA_WARNING_DIALOG_REATTACH_INFO"), kind) + "\nFormat: " + format + "\n\n" + warning
}

// PoolCandidate is a connected server of the New Pool dialog: whether it can be the coordinator and, with this
// coordinator, why it cannot join (PoolRules.Candidates).
type PoolCandidate struct {
	Key           string `json:"key"`
	Name          string `json:"name"`
	Reason        string `json:"reason"`
	CanCoordinate bool   `json:"can_coordinate"`
	CanJoin       bool   `json:"can_join"`
	IsCoordinator bool   `json:"is_coordinator"`
	reason        rules.JoinReason
}

// PoolCandidates lists the connected servers for the New Pool dialog with coordinator as the coordinator ("" for the
// first that can be one), in PoolJoinRules' order (coordinator, allowed, then by reason), then by name.
func (w *Workspace) PoolCandidates(coordinator string) ([]PoolCandidate, error) {
	facts, err := w.poolFacts()
	if err != nil {
		return nil, err
	}
	if _, ok := facts[coordinator]; !ok {
		coordinator = ""
		for _, k := range sortedKeys(facts) {
			if rules.CanBeCoordinator(rules.CanJoinPool(facts[k], facts[k])) {
				coordinator = k
				break
			}
		}
	}
	out := []PoolCandidate{}
	for _, k := range sortedKeys(facts) {
		r := rules.CanJoinPool(facts[k], facts[coordinator])
		if coordinator == "" {
			r = rules.JoinIsAPool
		}
		out = append(out, PoolCandidate{Key: k, Name: facts[k].Name, Reason: rules.JoinReasonText(r), reason: r,
			CanCoordinate: rules.CanBeCoordinator(rules.CanJoinPool(facts[k], facts[k])), CanJoin: r == rules.JoinAllowed,
			IsCoordinator: r == rules.JoinWillBeCoordinator})
	}
	sort.SliceStable(out, func(i, j int) bool {
		if out[i].reason != out[j].reason {
			return out[i].reason < out[j].reason
		}
		return strings.ToLower(out[i].Name) < strings.ToLower(out[j].Name)
	})
	return out, nil
}

// PoolResult is what CreatePool did: Created once the coordinator's pool has its name (a member that failed after that
// leaves a pool to finish with the other members), and the error, if any.
type PoolResult struct {
	Created bool   `json:"created"`
	Error   string `json:"error,omitempty"`
}

// PoolWarning is the question before creating the pool (PoolRules.Warnings' CPU levelling), "" for none.
func (w *Workspace) PoolWarning(coordinator string, members []string) (string, error) {
	facts, err := w.poolFacts()
	if err != nil {
		return "", err
	}
	var list []rules.PoolFacts
	for _, k := range append([]string{coordinator}, members...) {
		if f, ok := facts[k]; !ok || f.ReadError != "" { // a chosen server must be read: its warning may not be skipped
			return "", fmt.Errorf("%s could not be read: try again", k)
		}
	}
	for _, m := range members {
		list = append(list, facts[m])
	}
	// PoolRules.Warnings: the supplemental packs, then the CPU levelling
	return strings.Trim(rules.PackWarning(list, facts[coordinator])+"\n\n"+rules.CPUWarning(list, facts[coordinator]), "\n"), nil
}

// CreatePool is CreatePoolAction: the coordinator's pool gets the name and description, then each member joins it
// with pool.join and leaves the server list (it is part of the pool now). Every member is checked again first; the
// first member that fails stops the rest, as the C# action.
func (w *Workspace) CreatePool(coordinator string, members []string, name, description string) PoolResult {
	fail := func(err error) PoolResult { return PoolResult{Error: err.Error()} }
	if strings.TrimSpace(name) == "" {
		return fail(errors.New(messages.Get("POOL_NAME_EMPTY")))
	}
	facts, err := w.poolFacts()
	if err != nil {
		return fail(err)
	}
	cf, ok := facts[coordinator]
	if !ok || !rules.CanBeCoordinator(rules.CanJoinPool(cf, cf)) {
		return fail(errors.New(messages.Get("NO_ELIGIBLE_COORDINATOR")))
	}
	for _, m := range members {
		mf, ok := facts[m]
		if !ok {
			return fail(fmt.Errorf("%s is not connected", m))
		}
		if r := rules.CanJoinPool(mf, cf); r != rules.JoinAllowed {
			return fail(fmt.Errorf("%s: %s", mf.Name, rules.JoinReasonText(r)))
		}
	}
	w.mu.Lock()
	coord := w.conns[coordinator]
	password := w.passwords[coordinator]
	w.mu.Unlock()
	if coord == nil {
		return fail(errNotConnected)
	}
	var named bool
	if err := w.with(coordinator, func(c *xapi.Connection) error {
		var err error
		named, err = ops.NamePool(c, name, description)
		return err
	}); err != nil {
		if named { // the server is a pool now: the dialog cannot start over with it
			return PoolResult{Created: true, Error: fmt.Sprintf("The pool was created, but it is not ready for members: %v", err)}
		}
		return fail(err)
	}
	for i, m := range members {
		err := w.with(m, func(c *xapi.Connection) error { return ops.JoinPool(c, coord.Host, coord.User, password) })
		if err != nil {
			text := fmt.Sprintf("The pool was created, but %s did not join it: %v", facts[m].Name, err)
			if rest := len(members) - i - 1; rest > 0 {
				text += fmt.Sprintf(" (%d more servers were not added)", rest)
			}
			return PoolResult{Created: true, Error: text}
		}
		_ = w.Forget(m) // a member now: the coordinator's connection shows it
	}
	return PoolResult{Created: true}
}

// poolFacts reads PoolJoinRules' facts of every connected server.
func (w *Workspace) poolFacts() (map[string]rules.PoolFacts, error) {
	w.mu.Lock()
	keys := make([]string, 0, len(w.conns))
	for k := range w.conns {
		keys = append(keys, k)
	}
	w.mu.Unlock()
	facts := map[string]rules.PoolFacts{}
	for _, k := range keys {
		var f rules.PoolFacts
		err := w.with(k, func(c *xapi.Connection) error { var err error; f, err = app.ReadPoolFacts(c, k); return err })
		if errors.Is(err, errNotConnected) {
			continue
		}
		if err != nil { // only this server is out of the choice; the others stay usable
			f = rules.PoolFacts{Key: k, Name: k, ReadError: err.Error()}
		}
		facts[k] = f
	}
	return facts, nil
}

func sortedKeys(m map[string]rules.PoolFacts) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// Views lists the Organization Views and the Saved Searches (the default ones, then those saved in the pools, which
// are listed but not run yet: they use the XenSearch query language).
func (w *Workspace) Views() (organization, searches []app.ViewInfo, err error) {
	organization, searches = app.OrganizationViews(), app.DefaultSearches()
	seen := map[string]bool{}
	for _, k := range w.connectedKeys() {
		var list []app.ViewInfo
		err := w.with(k, func(c *xapi.Connection) error { var err error; list, err = app.PoolSearches(c); return err })
		if err != nil && !errors.Is(err, errNotConnected) {
			return nil, nil, err
		}
		for _, s := range list {
			if !seen[s.ID] {
				seen[s.ID] = true
				searches = append(searches, s)
			}
		}
	}
	return organization, searches, nil
}

// View is one view over every connected server (OrganizationalView / a default search / a search saved in a pool,
// which runs on every server, as the C# app), the servers' groups joined.
func (w *Workspace) View(id string) ([]*app.ViewNode, error) {
	saved := ""
	if uuid, ok := strings.CutPrefix(id, "saved:"); ok {
		for _, k := range w.connectedKeys() {
			_ = w.with(k, func(c *xapi.Connection) error { var err error; saved, err = app.PoolSearch(c, uuid); return err })
			if saved != "" {
				break
			}
		}
		if saved == "" {
			return nil, errors.New("this search is not saved in a connected pool any more")
		}
	}
	return w.runView(id, saved)
}

// runView reads a view, or runs a search when saved is its XML, on every connected server.
func (w *Workspace) runView(id, saved string) ([]*app.ViewNode, error) {
	var parts [][]*app.ViewNode
	for _, k := range w.connectedKeys() {
		var nodes []*app.ViewNode
		err := w.with(k, func(c *xapi.Connection) error {
			var err error
			if saved != "" {
				nodes, err = app.RunSearch(c, k, saved)
			} else {
				nodes, err = app.ReadView(c, k, id)
			}
			return err
		})
		if errors.Is(err, errNotConnected) {
			continue
		}
		if err != nil {
			return nil, err
		}
		parts = append(parts, nodes)
	}
	return app.MergeViews(parts...), nil
}

// connectedKeys are the servers connected now, sorted.
func (w *Workspace) connectedKeys() []string {
	w.mu.Lock()
	defer w.mu.Unlock()
	keys := make([]string, 0, len(w.conns))
	for k := range w.conns {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
