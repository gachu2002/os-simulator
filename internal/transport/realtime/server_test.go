package realtime

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gorilla/websocket"
)

func TestTransportDeterministicTimelineHash(t *testing.T) {
	ts := httptest.NewServer(NewServer(NewSessionManager()).Handler())
	defer ts.Close()

	body := map[string]any{"seed": 77, "policy": "rr", "quantum": 2}
	s1 := createSession(t, ts.URL, body)
	s2 := createSession(t, ts.URL, body)

	hash1 := runScenarioOverWS(t, ts.URL, s1.SessionID)
	hash2 := runScenarioOverWS(t, ts.URL, s2.SessionID)

	if hash1 != hash2 {
		t.Fatalf("expected deterministic hash match, got %s vs %s", hash1, hash2)
	}
}

func TestCommandValidationAndSequenceOrdering(t *testing.T) {
	ts := httptest.NewServer(NewServer(NewSessionManager()).Handler())
	defer ts.Close()

	s := createSession(t, ts.URL, map[string]any{"seed": 9})
	conn := dialWS(t, ts.URL, s.SessionID)
	defer conn.Close()

	connected := mustReadEvent(t, conn)
	if connected.Sequence != 2 {
		t.Fatalf("expected connected sequence 2, got %d", connected.Sequence)
	}

	mustWriteCommand(t, conn, Command{Name: "run", Count: 0})
	errEvent := mustReadEvent(t, conn)
	if errEvent.Type != "session.error" {
		t.Fatalf("expected session.error, got %s", errEvent.Type)
	}
	if errEvent.Sequence != 3 {
		t.Fatalf("expected sequence 3 after invalid command, got %d", errEvent.Sequence)
	}

	mustWriteCommand(t, conn, Command{Name: "step", Count: 1})
	snapEvent := mustReadEvent(t, conn)
	if snapEvent.Type != "session.snapshot" {
		t.Fatalf("expected session.snapshot, got %s", snapEvent.Type)
	}
	if snapEvent.Sequence != 4 {
		t.Fatalf("expected sequence 4 after valid command, got %d", snapEvent.Sequence)
	}
}

func TestCORSPreflightAndHeaders(t *testing.T) {
	os.Setenv("CORS_ALLOW_ORIGIN", "http://localhost:5173")
	t.Cleanup(func() {
		_ = os.Unsetenv("CORS_ALLOW_ORIGIN")
	})

	ts := httptest.NewServer(NewServer(NewSessionManager()).Handler())
	defer ts.Close()

	req, err := http.NewRequest(http.MethodOptions, ts.URL+"/sessions", nil)
	if err != nil {
		t.Fatalf("new request failed: %v", err)
	}
	req.Header.Set("Origin", "http://localhost:5173")
	req.Header.Set("Access-Control-Request-Method", "POST")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("options request failed: %v", err)
	}
	defer resp.Body.Close()

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

	req, err := http.NewRequest(http.MethodGet, ts.URL+"/sessions", nil)
	if err != nil {
		t.Fatalf("new request failed: %v", err)
	}
	req.Header.Set("X-Request-ID", "req-test-1")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()
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

func createSession(t *testing.T, baseURL string, payload map[string]any) CreateSessionResponse {
	t.Helper()
	b, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	resp, err := http.Post(baseURL+"/sessions", "application/json", bytes.NewReader(b))
	if err != nil {
		t.Fatalf("create session request failed: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected status %d got %d", http.StatusCreated, resp.StatusCode)
	}
	var out CreateSessionResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		t.Fatalf("decode create session response failed: %v", err)
	}
	if out.SessionID == "" || out.Snapshot == nil {
		t.Fatalf("invalid create session response: %+v", out)
	}
	return out
}

func runScenarioOverWS(t *testing.T, baseURL, sessionID string) string {
	t.Helper()
	conn := dialWS(t, baseURL, sessionID)
	defer conn.Close()

	_ = mustReadEvent(t, conn)

	mustWriteCommand(t, conn, Command{Name: "spawn", Process: "demo", Program: "COMPUTE 3; EXIT"})
	spawn := mustReadEvent(t, conn)
	if spawn.Snapshot == nil {
		t.Fatalf("spawn event missing snapshot")
	}

	mustWriteCommand(t, conn, Command{Name: "step", Count: 12})
	step := mustReadEvent(t, conn)
	if step.Snapshot == nil {
		t.Fatalf("step event missing snapshot")
	}
	return step.Snapshot.TraceHash
}

func dialWS(t *testing.T, baseURL, sessionID string) *websocket.Conn {
	t.Helper()
	wsURL := "ws" + strings.TrimPrefix(baseURL, "http") + "/ws/" + sessionID
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("dial websocket failed: %v", err)
	}
	return conn
}

func mustWriteCommand(t *testing.T, conn *websocket.Conn, cmd Command) {
	t.Helper()
	req := CommandEnvelope{Type: "command", Command: cmd}
	if err := conn.WriteJSON(req); err != nil {
		t.Fatalf("write command failed: %v", err)
	}
}

func mustReadEvent(t *testing.T, conn *websocket.Conn) Event {
	t.Helper()
	var ev Event
	if err := conn.ReadJSON(&ev); err != nil {
		t.Fatalf("read event failed: %v", err)
	}
	return ev
}
