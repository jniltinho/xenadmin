package console

import (
	"context"
	"io"
	"log"
	"net/http"

	"github.com/coder/websocket"
)

// ServeHTTP is the WebSocket endpoint: the page must come from the webview (its Origin) and bring a live ticket;
// the console stream is then bridged both ways until either side closes.
func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/console" {
		http.NotFound(w, r)
		return
	}
	if origin := r.Header.Get("Origin"); !p.origins[origin] { // a missing Origin (not a browser page) is refused too
		log.Printf("console: refused a connection from origin %q", origin)
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	st := &stream{cancel: cancel}
	// consume the ticket and register the console in one step: a Revoke either finds the ticket or the stream
	p.mu.Lock()
	t, ok := p.take(r.URL.Query().Get("token"))
	if !ok || p.closed {
		p.mu.Unlock()
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	p.streams[t.owner] = append(p.streams[t.owner], st) // from now on Revoke stops it, even while dialing
	p.mu.Unlock()
	defer p.forget(t.owner, st)
	// the Origin was checked above against the webview's own, which the library's same-host check cannot know
	ws, err := websocket.Accept(w, r, &websocket.AcceptOptions{Subprotocols: []string{"binary"}, InsecureSkipVerify: true})
	if err != nil {
		return
	}
	conn, err := t.open(ctx)
	if err != nil {
		ws.Close(websocket.StatusInternalError, clip(err.Error()))
		return
	}
	p.mu.Lock()
	dead := st.dead
	st.conn = conn
	p.mu.Unlock()
	if dead {
		conn.Close()
		ws.Close(websocket.StatusGoingAway, "disconnected")
		return
	}
	ws.SetReadLimit(-1)                                       // RFB messages (clipboard, framebuffer updates) have no small upper bound
	nc := websocket.NetConn(ctx, ws, websocket.MessageBinary) // closed when ctx is cancelled (Revoke)
	done := make(chan struct{}, 2)
	go func() { io.Copy(conn, nc); done <- struct{}{} }() //nolint:errcheck // ends when either side closes
	go func() { io.Copy(nc, conn); done <- struct{}{} }() //nolint:errcheck
	<-done
	conn.Close()
	nc.Close()
}

func (p *Proxy) forget(owner string, st *stream) {
	p.mu.Lock()
	defer p.mu.Unlock()
	st.cancel()
	if st.conn != nil {
		st.conn.Close()
	}
	list := p.streams[owner]
	for i, s := range list {
		if s == st {
			p.streams[owner] = append(list[:i], list[i+1:]...)
			break
		}
	}
}

// clip keeps a close reason within the WebSocket limit (123 bytes).
func clip(s string) string {
	if len(s) > 120 {
		return s[:120]
	}
	return s
}
