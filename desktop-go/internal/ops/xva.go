package ops

import (
	"bufio"
	"cmp"
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	xenapi "go/xenapi"

	"xcpng-center/internal/messages"
	"xcpng-center/internal/rules"
	"xcpng-center/internal/trust"
	"xcpng-center/internal/xapi"
)

// XVA export and import (XvaRules: ExportVmAction, and the import wizard's XVA path, ImportVmAction). The session goes
// in a cookie, never in the URL; a server's certificate goes through the member trust as the performance graphs'.
// shortcut: no progress bar or Cancel yet. OVF/OVA appliances are in ovf.go.

// xvaClient is the HTTP client of a transfer to the server at address (the coordinator, or the SR's server): no limit
// on the whole transfer (an export takes hours), but a connection idle for transferIdle fails.
func xvaClient(s *xenapi.Session, policy *trust.Policy, address string, port int) (*http.Client, string) {
	key := xapi.HostKey(address, port)
	dialer := xapi.Dialer()
	return &http.Client{Transport: &http.Transport{Proxy: xapi.Proxy,
			TLSClientConfig: xapi.MemberTLS(policy, key, xapi.VouchedFingerprints(s, address)), TLSHandshakeTimeout: xapi.HandshakeTimeout(),
			ResponseHeaderTimeout: 30 * time.Minute, // HTTP_PUT_TIMEOUT: the server answers an import once it has read it
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				c, err := dialer.DialContext(ctx, network, addr)
				if err != nil {
					return nil, err
				}
				return &idleConn{Conn: c}, nil
			}},
			CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }},
		"https://" + key
}

// transferIdle is how long a transfer may make no progress.
const transferIdle = 10 * time.Minute

// idleConn fails a read or write that makes no progress for transferIdle.
type idleConn struct{ net.Conn }

func (c *idleConn) Read(p []byte) (int, error) {
	_ = c.Conn.SetReadDeadline(time.Now().Add(transferIdle))
	return c.Conn.Read(p)
}

func (c *idleConn) Write(p []byte) (int, error) {
	_ = c.Conn.SetWriteDeadline(time.Now().Add(transferIdle))
	return c.Conn.Write(p)
}

// ExportBlocker is XvaRules.CanExport / ExportBlocker: why the VM cannot be exported now ("" when it can).
func ExportBlocker(vm xenapi.VMRecord) string {
	switch {
	case vm.IsATemplate || vm.IsASnapshot || vm.IsControlDomain:
		return "Only VMs can be exported."
	case vm.PowerState == xenapi.VMPowerStateRunning:
		return messages.Get("MAINWINDOW_EXPORT_VM_AS_BACKUP_TOOLTIP")
	case len(vm.CurrentOperations) > 0 || !contains(vm.AllowedOperations, xenapi.VMOperationsExport):
		return "The VM cannot be exported now."
	}
	return ""
}

// ExportVM is ExportVmAction: GET /export of the VM into "<file>.tmp", flushed to disk, verified if asked (Export.verify),
// then renamed to the file; a failed or unverified export leaves no file behind.
func ExportVM(c *xapi.Connection, policy *trust.Policy, vmUUID, file string, verify bool) error {
	s := c.Session
	ref, err := xenapi.VM.GetByUUID(s, vmUUID)
	if err != nil {
		return xapi.Sanitize(err)
	}
	vm, err := xenapi.VM.GetRecord(s, ref)
	if err != nil {
		return xapi.Sanitize(err)
	}
	if why := ExportBlocker(vm); why != "" {
		return errors.New(why)
	}
	if !filepath.IsAbs(file) {
		return errors.New("choose where to save the export")
	}
	task, err := xenapi.Task.Create(s, "export", "Exporting "+vm.NameLabel+" to backup file")
	if err != nil {
		return xapi.Sanitize(err)
	}
	// a temporary file of its own beside the export, removed unless it becomes the export
	f, err := os.CreateTemp(filepath.Dir(file), "."+filepath.Base(file)+".*.tmp")
	if err != nil {
		return err
	}
	tmp := f.Name()
	defer os.Remove(tmp)
	err = download(s, policy, c.Host, c.Port, "/export?"+url.Values{"uuid": {vm.UUID}, "task_id": {string(task)}}.Encode(), f)
	if cerr := f.Close(); err == nil {
		err = cerr
	}
	if err = finishTask(s, task, err, nil); err != nil {
		return err
	}
	if verify { // "Verify export on completion": read the file back from the disk
		if err := verifyFile(tmp); err != nil {
			return err
		}
	}
	return os.Rename(tmp, file)
}

// finishTask waits for the transfer's task (its result read by done). After a failed transfer, a task the server never
// started (pending, no progress, both read) is destroyed, as ExportVmAction's CA-11100 case; one it started tells its
// own error; one that cannot be read is left alone (ErrUncertain).
func finishTask(s *xenapi.Session, task xenapi.TaskRef, transfer error, done func()) error {
	if transfer != nil {
		status, serr := xenapi.Task.GetStatus(s, task)
		progress, perr := xenapi.Task.GetProgress(s, task)
		if serr != nil || perr != nil {
			return fmt.Errorf("%w (%v)", ErrUncertain, transfer)
		}
		if status == xenapi.TaskStatusTypePending && progress == 0 {
			_ = xenapi.Task.Destroy(s, task)
			return transfer
		}
	}
	if err := pollThen(s, task, nil, done); err != nil {
		return xapi.Sanitize(err)
	}
	return transfer
}

// download is ExportVmAction.HttpGet: the export written to f while Export.verify checks it as it comes (the block
// checksums and the end of the archive), then flushed to the disk.
func download(s *xenapi.Session, policy *trust.Policy, host string, port int, path string, f *os.File) error {
	client, base := xvaClient(s, policy, host, port)
	req, err := http.NewRequest(http.MethodGet, base+path, nil)
	if err != nil {
		return err
	}
	req.AddCookie(&http.Cookie{Name: "session_id", Value: string(s.Ref())})
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("the export failed: %s", resp.Status)
	}
	w := bufio.NewWriterSize(f, 1<<20)
	tee := io.TeeReader(resp.Body, w)
	if err := rules.VerifyXVA(tee); err != nil {
		return err
	}
	if _, err := io.Copy(io.Discard, tee); err != nil { // what follows the end of the archive, kept as the server sent it
		return err
	}
	if err := w.Flush(); err != nil {
		return err
	}
	return f.Sync() // FlushFileBuffers / fsync before the file is called done
}

func verifyFile(file string) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()
	return rules.VerifyXVA(bufio.NewReaderSize(f, 1<<20))
}

// ImportFileError is XvaRules.FileError: an existing XVA file (optionally gzipped).
func ImportFileError(file string) string {
	name := strings.ToLower(strings.TrimSpace(file))
	switch {
	case name == "":
		return "Select the file to import."
	case !strings.HasSuffix(name, ".xva") && !strings.HasSuffix(name, ".xva.gz"):
		return "Only XVA files (.xva, .xva.gz) and OVF/OVA appliances (.ovf, .ova, .ova.gz, .ova.bz2) can be imported."
	}
	if st, err := os.Stat(file); err != nil || st.IsDir() {
		return "The file does not exist."
	}
	return ""
}

// ImportVM is ImportVmAction for an XVA: PUT /import to the SR's server (the coordinator for shared storage), then
// the C# follow-up: a name already used gets " (n)", the home server is set, and (start) the VM is started or resumed.
// It returns the new VM's uuid.
func ImportVM(c *xapi.Connection, policy *trust.Policy, file, hostUUID, srUUID string, start bool) (string, error) {
	if why := ImportFileError(file); why != "" {
		return "", errors.New(why)
	}
	s := c.Session
	sr, err := xenapi.SR.GetByUUID(s, srUUID)
	if err != nil {
		return "", xapi.Sanitize(err)
	}
	srRec, err := xenapi.SR.GetRecord(s, sr)
	if err != nil {
		return "", xapi.Sanitize(err)
	}
	if rules.SRFull(srRec.Type, srRec.PhysicalSize, srRec.PhysicalUtilisation) || srRec.ContentType == "iso" { // XvaRules.TargetError
		return "", errors.New("the storage repository is no longer available or is full")
	}
	var affinity xenapi.HostRef
	if hostUUID != "" {
		if affinity, err = xenapi.Host.GetByUUID(s, hostUUID); err != nil {
			return "", xapi.Sanitize(err)
		}
	}
	address, port := storageHost(c, srRec)
	task, err := xenapi.Task.Create(s, "put_import_task", address)
	if err != nil {
		return "", xapi.Sanitize(err)
	}
	var result string
	err = upload(s, policy, address, port, "/import?"+url.Values{"task_id": {string(task)}, "sr_id": {string(sr)}}.Encode(), file)
	if err = finishTask(s, task, err, func() { result, _ = xenapi.Task.GetResult(s, task) }); err != nil {
		return "", err
	}
	ref := xenapi.VMRef(resultRef(result))
	if ref == "" {
		return "", fmt.Errorf("%w (the import succeeded without saying which VM it made)", ErrUncertain)
	}
	vm, err := xenapi.VM.GetRecord(s, ref)
	if err != nil {
		return "", xapi.Sanitize(err)
	}
	if vm.IsATemplate && vm.IsDefaultTemplate && vm.OtherConfig["import_task"] != string(task) {
		return "", errors.New(messages.Format(messages.Get("IMPORT_TEMPLATE_ALREADY_EXISTS"), "XCP-ng"))
	}
	if name := defaultVMName(s, vm.NameLabel); name != vm.NameLabel {
		if err := xenapi.VM.SetNameLabel(s, ref, name); err != nil {
			return vm.UUID, xapi.Sanitize(err)
		}
	}
	if !vm.IsATemplate && affinity != "" {
		if err := xenapi.VM.SetAffinity(s, ref, affinity); err != nil {
			return vm.UUID, xapi.Sanitize(err)
		}
	}
	if start && !vm.IsATemplate { // XvaRules.StartAfterImport
		op := rules.VMStart
		if state, err := xenapi.VM.GetPowerState(s, ref); err == nil && state == xenapi.VMPowerStateSuspended {
			op = rules.VMResume
		}
		return vm.UUID, VMPower(c, vm.UUID, op)
	}
	return vm.UUID, nil
}

// storageHost is SR.GetStorageHost: the server a transfer to the SR goes to, a local SR's own server or the coordinator.
// A connection on another port than 443 (a forwarded one) keeps its own address: the server's is not reachable there.
func storageHost(c *xapi.Connection, sr xenapi.SRRecord) (string, int) {
	if !sr.Shared && len(sr.PBDs) == 1 && c.Port == 443 {
		if h, err := xenapi.PBD.GetHost(c.Session, sr.PBDs[0]); err == nil {
			if a, err := xenapi.Host.GetAddress(c.Session, h); err == nil && a != "" {
				return a, 443
			}
		}
	}
	return c.Host, c.Port
}

func upload(s *xenapi.Session, policy *trust.Policy, host string, port int, path, file string) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()
	st, err := f.Stat()
	if err != nil {
		return err
	}
	return put(s, policy, host, port, path, bufio.NewReaderSize(f, 1<<20), st.Size())
}

// exactReader gives left bytes of r and, with the last of them, fails unless r ends there: a stream that has more
// (or a decoder that finds an error past them, as a VMDK grain out of order) fails the transfer instead of being cut.
// The transport closes it when it is done with the body (Close); put waits for that, then checks the stream ended
// cleanly (a server may answer before the body is all sent).
type exactReader struct {
	r      io.Reader
	left   int64
	done   bool  // every byte read and the stream ended there
	failed error // the error a read returned
	closed chan struct{}
	once   sync.Once
}

func (e *exactReader) Read(p []byte) (int, error) {
	if e.left <= 0 {
		return 0, io.EOF
	}
	if int64(len(p)) > e.left {
		p = p[:e.left]
	}
	n, err := e.r.Read(p)
	if e.left -= int64(n); e.left == 0 {
		// only the plain io.EOF ends a stream: a decoder's "truncated: %w" of io.EOF is a failure
		if extra, perr := e.r.Read(make([]byte, 1)); extra > 0 || perr != nil && perr != io.EOF {
			e.failed = cmp.Or(perr, errors.New("the data is longer than its size"))
			return n, e.failed
		}
		e.done = true
		return n, io.EOF
	}
	if err == io.EOF {
		err = io.ErrUnexpectedEOF
	}
	if err != nil {
		e.failed = err
	}
	return n, err
}

func (e *exactReader) Close() error { e.once.Do(func() { close(e.closed) }); return nil }

// put is HTTPHelper.PUT: length bytes of body to the path, the session in a cookie.
func put(s *xenapi.Session, policy *trust.Policy, host string, port int, path string, body io.Reader, length int64) error {
	client, base := xvaClient(s, policy, host, port)
	er := &exactReader{r: body, left: length, closed: make(chan struct{})}
	req, err := http.NewRequest(http.MethodPut, base+path, er)
	if err != nil {
		return err
	}
	req.ContentLength = length
	req.AddCookie(&http.Cookie{Name: "session_id", Value: string(s.Ref())})
	resp, err := client.Do(req)
	if err != nil {
		var u *url.Error
		if errors.As(err, &u) { // the cause, without the URL and its task and disk references
			return u.Err
		}
		return err
	}
	defer resp.Body.Close()
	_, derr := io.Copy(io.Discard, io.LimitReader(resp.Body, 1<<16))
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("the import failed: %s", resp.Status)
	}
	// a body cut by a failed last write is the import task's failure (the callers' finishTask), not seen here
	select { // the transport closes the body once it stopped sending it
	case <-er.closed:
	case <-time.After(time.Minute):
		return errors.New("the upload did not finish")
	}
	switch {
	case er.failed != nil:
		return er.failed
	case !er.done:
		return errors.New("the server answered before the data was all sent")
	case derr != nil:
		return derr
	}
	return nil
}

// defaultVMName is ImportVmAction.DefaultVMName: when several VMs have the name, the first "name (n)" nobody uses.
func defaultVMName(s *xenapi.Session, name string) string {
	count := func(n string) int {
		refs, _ := xenapi.VM.GetByNameLabel(s, n)
		return len(refs)
	}
	if count(name) <= 1 {
		return name
	}
	for i := 1; ; i++ {
		if n := messages.Format(messages.Get("NEWVM_DEFAULTNAME"), name, i); count(n) == 0 {
			return n
		}
	}
}
