package ops

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	xenapi "go/xenapi"

	"xcpng-center/internal/xapi"
)

// A task that succeeded but whose result cannot be read is uncertain (the object exists), never a plain failure.
func TestWaitResultUnreadableIsUncertain(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Method string `json:"method"`
			ID     any    `json:"id"`
		}
		_ = json.NewDecoder(r.Body).Decode(&req)
		resp := map[string]any{"jsonrpc": "2.0", "id": req.ID}
		switch req.Method {
		case "task.get_status":
			resp["result"] = "success"
		case "task.get_result":
			resp["error"] = xenapi.ResponseError{Code: 1, Message: "HANDLE_INVALID"}
		default:
			resp["result"] = ""
		}
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer srv.Close()
	every, retries := pollEvery, pollRetries
	pollEvery, pollRetries = 0, 2
	defer func() { pollEvery, pollRetries = every, retries }()
	s := xenapi.NewSession(&xenapi.ClientOpts{URL: srv.URL})
	if _, err := waitResult(s)(xenapi.TaskRef("OpaqueRef:task"), nil); !errors.Is(err, ErrUncertain) {
		t.Fatalf("got %v", err)
	}
}

func TestCreateSettled(t *testing.T) {
	refused := xapi.Sanitize(&xenapi.APIError{Message: "SR_BACKEND_FAILURE_111"})
	for _, c := range []struct {
		submitted, created bool
		err                error
		want               bool
	}{
		{false, false, errors.New("secret"), true}, // never sent
		{true, true, refused, true},                // made, failed afterwards
		{true, false, refused, true},               // refused by XAPI
		{true, false, errors.New("connection reset"), false},
		{true, false, ErrUncertain, false},
	} {
		if got := createSettled(c.submitted, c.created, c.err); got != c.want {
			t.Errorf("%+v: %v", c, got)
		}
	}
}
