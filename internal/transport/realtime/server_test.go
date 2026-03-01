package realtime

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCORSPreflightAndHeaders(t *testing.T) {
	t.Setenv("CORS_ALLOW_ORIGIN", "http://localhost:5173")

	ts := httptest.NewServer(NewServer(NewSessionManager()).Handler())
	defer ts.Close()

	req, err := http.NewRequest(http.MethodOptions, ts.URL+"/challenges/start/v3", nil)
	if err != nil {
		t.Fatalf("new request failed: %v", err)
	}
	req.Header.Set("Origin", "http://localhost:5173")
	req.Header.Set("Access-Control-Request-Method", "POST")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("options request failed: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusNoContent {
		t.Fatalf("status=%d want=%d", resp.StatusCode, http.StatusNoContent)
	}
	if got := resp.Header.Get("Access-Control-Allow-Origin"); got != "http://localhost:5173" {
		t.Fatalf("allow-origin=%q want=%q", got, "http://localhost:5173")
	}
}

func TestErrorEnvelopeIncludesRequestID(t *testing.T) {
	ts := httptest.NewServer(NewServer(NewSessionManager()).Handler())
	defer ts.Close()

	req, err := http.NewRequest(http.MethodGet, ts.URL+"/challenges/start/v3", nil)
	if err != nil {
		t.Fatalf("new request failed: %v", err)
	}
	req.Header.Set("X-Request-ID", "req-test-1")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Fatalf("status=%d want=%d", resp.StatusCode, http.StatusMethodNotAllowed)
	}

	var out map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		t.Fatalf("decode failed: %v", err)
	}
	if out["code"] != "method_not_allowed" {
		t.Fatalf("code=%v want=method_not_allowed", out["code"])
	}
	if out["request_id"] != "req-test-1" {
		t.Fatalf("request_id=%v want=req-test-1", out["request_id"])
	}
}
