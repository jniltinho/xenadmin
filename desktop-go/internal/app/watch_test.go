package app

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	xenapi "go/xenapi"
)

// The first batch is a change (the baseline); later batches with events are changes; EVENTS_LOST starts over from a
// new baseline (one more change); cancelling stops the loop, also a call that is waiting.
func TestWatch(t *testing.T) {
	var calls atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			ID     any   `json:"id"`
			Params []any `json:"params"`
		}
		_ = json.NewDecoder(r.Body).Decode(&req)
		resp := map[string]any{"jsonrpc": "2.0", "id": req.ID}
		event := map[string]any{"id": 1, "class": "vm", "operation": "mod", "ref": "OpaqueRef:a", "obj_uuid": "u", "timestamp": "20260101T00:00:00Z"}
		switch calls.Add(1) {
		case 1:
			resp["result"] = map[string]any{"token": "t1", "events": []any{event}, "valid_ref_counts": map[string]any{}}
		case 2:
			resp["result"] = map[string]any{"token": "t2", "events": []any{event}, "valid_ref_counts": map[string]any{}}
		case 3:
			resp["error"] = map[string]any{"code": 1, "message": "EVENTS_LOST", "data": []any{}}
		default:
			<-r.Context().Done() // a long event.from: only cancelling the client ends it
			return
		}
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer srv.Close()
	ctx, cancel := context.WithCancel(context.Background())
	s := xenapi.NewSession(&xenapi.ClientOpts{URL: srv.URL}).WithContext(ctx)
	var changes atomic.Int32
	done := make(chan struct{})
	go func() { Watch(ctx, s, func() { changes.Add(1) }); close(done) }()
	deadline := time.Now().Add(5 * time.Second)
	for calls.Load() < 4 && time.Now().Before(deadline) {
		time.Sleep(10 * time.Millisecond)
	}
	cancel()
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("Watch did not stop")
	}
	if n := changes.Load(); n != 2 { // the baseline and batch 2; EVENTS_LOST's new baseline waits in call 4
		t.Fatalf("%d changes", n)
	}
}
