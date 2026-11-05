package console

import (
	"bufio"
	"context"
	"crypto/tls"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/coder/websocket"

	"xcpng-center/internal/trust"
)

type pins map[string]string

func (p pins) Fingerprint(k string) (string, bool) { v, ok := p[k]; return v, ok }
func (p pins) Pin(k, v string) error               { p[k] = v; return nil }

// a fake console host: TLS, answers CONNECT with 200 and sends the RFB greeting right after the headers
func consoleHost(t *testing.T) (*httptest.Server, string) {
	srv := httptest.NewUnstartedServer(nil)
	srv.Listener.Close()
	ln, _ := net.Listen("tcp", "127.0.0.1:0")
	srv.Listener = ln
	srv.StartTLS()
	srv.Listener.Close() // replace the HTTP server with a raw handler on the same certificate
	raw, _ := tls.Listen("tcp", "127.0.0.1:0", srv.TLS)
	go func() {
		for {
			c, err := raw.Accept()
			if err != nil {
				return
			}
			go func(c net.Conn) {
				r := bufio.NewReader(c)
				line, _ := r.ReadString('\n')
				for l, _ := r.ReadString('\n'); strings.TrimSpace(l) != ""; l, _ = r.ReadString('\n') {
					if strings.HasPrefix(l, "Cookie: ") && !strings.Contains(l, "session_id=S") {
						c.Close()
						return
					}
				}
				if !strings.HasPrefix(line, "CONNECT /console?uuid=x HTTP/1.0") {
					io.WriteString(c, "HTTP/1.0 404 Not Found\r\n\r\n")
					c.Close()
					return
				}
				io.WriteString(c, "HTTP/1.0 200 OK\r\n\r\nRFB 003.008\n") // greeting in the same packet as the headers
				io.Copy(c, r)                                             // echo
			}(c)
		}
	}()
	t.Cleanup(func() { raw.Close(); srv.Close() })
	return srv, raw.Addr().String()
}

func TestConsoleThroughTheProxy(t *testing.T) {
	srv, addr := consoleHost(t)
	fp := trust.Fingerprint(srv.Certificate().Raw)
	u, _ := url.Parse("https://" + addr + "/console?uuid=x")
	target := Target{URL: u, HostKey: addr}

	// an unknown certificate is refused; one the coordinator vouches for is accepted
	policy := &trust.Policy{Store: pins{}, OperatorPins: map[string]string{}}
	if _, err := Dial(context.Background(), policy, target, "S"); err == nil {
		t.Fatal("unvouched certificate accepted")
	}
	target.Vouched = []string{fp}
	p, err := NewProxy("wails://localhost")
	if err != nil {
		t.Fatal(err)
	}
	defer p.Close()
	open := func(ctx context.Context) (net.Conn, error) { return Dial(ctx, policy, target, "S") }
	wsURL, _ := p.Ticket("pool", open)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if _, _, err := websocket.Dial(ctx, wsURL, &websocket.DialOptions{HTTPHeader: http.Header{"Origin": {"https://evil.example"}}}); err == nil {
		t.Fatal("a hostile origin was accepted")
	}
	opts := &websocket.DialOptions{HTTPHeader: http.Header{"Origin": {"wails://localhost"}}, Subprotocols: []string{"binary"}}
	ws, _, err := websocket.Dial(ctx, wsURL, opts)
	if err != nil {
		t.Fatal(err)
	}
	_, greeting, err := ws.Read(ctx)
	if err != nil || string(greeting) != "RFB 003.008\n" {
		t.Fatalf("greeting %q %v", greeting, err)
	}
	ws.Write(ctx, websocket.MessageBinary, []byte("RFB 003.008\n"))
	if _, echo, err := ws.Read(ctx); err != nil || string(echo) != "RFB 003.008\n" {
		t.Fatalf("echo %q %v", echo, err)
	}
	// a ticket opens once
	if _, _, err := websocket.Dial(ctx, wsURL, opts); err == nil {
		t.Fatal("a ticket was used twice")
	}
	// disconnecting the server closes its consoles
	p.Revoke("pool")
	if _, _, err := ws.Read(ctx); err == nil {
		t.Fatal("console still open after Revoke")
	}
}

func TestRevokeStopsAConsoleStillDialing(t *testing.T) {
	p, err := NewProxy("wails://localhost")
	if err != nil {
		t.Fatal(err)
	}
	defer p.Close()
	dialing, cancelled := make(chan struct{}), make(chan struct{})
	wsURL, _ := p.Ticket("pool", func(ctx context.Context) (net.Conn, error) {
		close(dialing)
		<-ctx.Done() // a slow console host
		close(cancelled)
		return nil, ctx.Err()
	})
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	go websocket.Dial(ctx, wsURL, &websocket.DialOptions{HTTPHeader: http.Header{"Origin": {"wails://localhost"}}}) //nolint:errcheck
	<-dialing
	p.Revoke("pool")
	select {
	case <-cancelled:
	case <-ctx.Done():
		t.Fatal("Revoke did not stop the dial")
	}
	p.Close()
	if _, err := p.Ticket("pool", nil); err == nil {
		t.Fatal("a ticket after Close")
	}
}
