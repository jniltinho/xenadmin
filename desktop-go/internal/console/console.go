// Package console opens a VM's graphical console (RFB over XAPI's HTTP CONNECT, as XSVNCScreen with
// HTTP.HttpConnectStream) and serves it to the window's noVNC through a loopback WebSocket that only a one-time
// token, bound to that VM and workspace connection, can open.
package console

import (
	"bufio"
	"context"
	"crypto/rand"
	"crypto/tls"
	"encoding/hex"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	xenapi "go/xenapi"

	"xcpng-center/internal/trust"
	"xcpng-center/internal/xapi"
)

// Target is where a VM's console is: its host and path, and the fingerprints its certificate may have.
type Target struct {
	URL     *url.URL
	HostKey string
	Vouched []string // the console host's certificates as the authenticated coordinator reports them
}

// Locate finds the RFB console of a running VM (VM.consoles, protocol rfb) and the certificates the coordinator
// reports for the host that serves it (CertificateTrust.PoolMemberFingerprints).
func Locate(c *xapi.Connection, vmUUID string) (Target, error) {
	s := c.Session
	vm, err := xenapi.VM.GetByUUID(s, vmUUID)
	if err != nil {
		return Target{}, xapi.Sanitize(err)
	}
	consoles, err := xenapi.VM.GetConsoles(s, vm)
	if err != nil {
		return Target{}, xapi.Sanitize(err)
	}
	for _, ref := range consoles {
		rec, err := xenapi.Console.GetRecord(s, ref)
		if err != nil {
			return Target{}, xapi.Sanitize(err)
		}
		if rec.Protocol != xenapi.ConsoleProtocolRfb {
			continue
		}
		u, err := url.Parse(rec.Location)
		if err != nil || u.Scheme != "https" || u.Hostname() == "" {
			return Target{}, fmt.Errorf("unexpected console location %q", rec.Location)
		}
		port := 443
		if p := u.Port(); p != "" {
			if _, err := fmt.Sscan(p, &port); err != nil {
				return Target{}, fmt.Errorf("unexpected console port %q", p)
			}
		}
		t := Target{URL: u, HostKey: xapi.HostKey(u.Hostname(), port)}
		t.Vouched = xapi.VouchedFingerprints(s, u.Hostname())
		return t, nil
	}
	return Target{}, errors.New("this VM has no graphical console now (is it running?)")
}

// Dial opens the console stream: TLS to the console's host, accepted when it is the pinned server (the policy, no
// prompt) or a pool member whose certificate the coordinator vouches for; then "CONNECT path" with the session
// cookie. Bytes the server sent after its HTTP answer stay in the returned connection.
func Dial(ctx context.Context, policy *trust.Policy, t Target, session string) (net.Conn, error) {
	conf := xapi.MemberTLS(policy, t.HostKey, t.Vouched)
	d := tls.Dialer{NetDialer: &net.Dialer{Timeout: 20 * time.Second}, Config: conf}
	raw, err := d.DialContext(ctx, "tcp", t.HostKey)
	if err != nil {
		return nil, err
	}
	conn := raw.(*tls.Conn)
	_ = conn.SetDeadline(time.Now().Add(30 * time.Second))
	// a revoked console stops waiting for the host's answer at once
	closeOnCancel := context.AfterFunc(ctx, func() { conn.Close() })
	path := t.URL.RequestURI()
	if _, err := fmt.Fprintf(conn, "CONNECT %s HTTP/1.0\r\nHost: %s\r\nCookie: session_id=%s\r\n\r\n", path, t.URL.Hostname(), session); err != nil {
		conn.Close()
		return nil, err
	}
	r := bufio.NewReader(conn)
	status, err := r.ReadString('\n')
	if err != nil {
		conn.Close()
		return nil, err
	}
	if f := strings.Fields(status); len(f) < 2 || f[1] != "200" {
		conn.Close()
		return nil, fmt.Errorf("the console was refused: %s", strings.TrimSpace(status))
	}
	for { // the rest of the HTTP answer's headers
		line, err := r.ReadString('\n')
		if err != nil {
			conn.Close()
			return nil, err
		}
		if strings.TrimSpace(line) == "" {
			break
		}
	}
	if !closeOnCancel() {
		return nil, ctx.Err() // cancelled during the handshake: the connection is closed
	}
	_ = conn.SetDeadline(time.Time{})
	return &buffered{Conn: conn, r: r}, nil
}

// buffered keeps the bytes bufio read past the HTTP headers (the start of the RFB handshake).
type buffered struct {
	net.Conn
	r *bufio.Reader
}

func (b *buffered) Read(p []byte) (int, error) { return b.r.Read(p) }

// Opener dials a console for a ticket (the workspace connects with its session and policy).
type Opener func(ctx context.Context) (net.Conn, error)

type ticket struct {
	owner   string // the workspace connection (server key): revoked when it disconnects
	open    Opener
	expires time.Time
}

// stream is one console being opened or open: revoking cancels its dial and closes it.
type stream struct {
	cancel context.CancelFunc
	conn   net.Conn // nil while dialing
	dead   bool     // revoked: a dial that finishes now is closed at once
}

// Proxy is the loopback WebSocket endpoint noVNC connects to.
type Proxy struct {
	mu      sync.Mutex
	srv     *http.Server
	ln      net.Listener
	closed  bool
	tickets map[string]ticket
	streams map[string][]*stream // owner -> consoles being opened or open
	origins map[string]bool
}

// NewProxy listens on 127.0.0.1 (a random port); origins are the webview's (wails://localhost and the like).
func NewProxy(origins ...string) (*Proxy, error) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return nil, err
	}
	p := &Proxy{ln: ln, tickets: map[string]ticket{}, streams: map[string][]*stream{}, origins: map[string]bool{}}
	for _, o := range origins {
		p.origins[o] = true
	}
	p.srv = &http.Server{Handler: p, ReadHeaderTimeout: 10 * time.Second}
	go p.srv.Serve(ln) //nolint:errcheck // ends when Close closes the server
	return p, nil
}

// Ticket returns a ws:// URL that opens this console once, within 30 seconds.
func (p *Proxy) Ticket(owner string, open Opener) (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	token := hex.EncodeToString(b)
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.closed {
		return "", errors.New("the console endpoint is closed")
	}
	now := time.Now()
	for t, k := range p.tickets { // expired ones
		if now.After(k.expires) {
			delete(p.tickets, t)
		}
	}
	p.tickets[token] = ticket{owner: owner, open: open, expires: now.Add(30 * time.Second)}
	return fmt.Sprintf("ws://%s/console?token=%s", p.ln.Addr(), token), nil
}

// Revoke drops the owner's unused tickets and stops its consoles, also the ones still dialing (the server was
// disconnected).
func (p *Proxy) Revoke(owner string) {
	p.mu.Lock()
	for t, k := range p.tickets {
		if k.owner == owner {
			delete(p.tickets, t)
		}
	}
	streams := p.streams[owner]
	delete(p.streams, owner)
	var conns []net.Conn
	for _, s := range streams {
		s.dead = true // a dial that finishes now closes its connection (ServeHTTP checks it under p.mu)
		s.cancel()
		if s.conn != nil {
			conns = append(conns, s.conn)
		}
	}
	p.mu.Unlock()
	for _, c := range conns {
		c.Close()
	}
}

// Close stops the endpoint, its tickets and every console; no ticket is issued afterwards.
func (p *Proxy) Close() {
	p.mu.Lock()
	p.closed = true
	p.tickets = map[string]ticket{}
	owners := make([]string, 0, len(p.streams))
	for o := range p.streams {
		owners = append(owners, o)
	}
	p.mu.Unlock()
	p.srv.Close() // the listener and the connections not yet upgraded (upgraded ones stop with their stream)
	for _, o := range owners {
		p.Revoke(o)
	}
}

// take consumes a ticket (p.mu held): a token works once, and only before it expires.
func (p *Proxy) take(token string) (ticket, bool) {
	t, ok := p.tickets[token]
	delete(p.tickets, token)
	return t, ok && time.Now().Before(t.expires)
}
