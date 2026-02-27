package realtime

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/websocket"
)

func TestTransportDeterministicTimelineHash(t *testing.T) {
	ts := httptest.NewServer(NewServer(NewSessionManager()).Handler())
	defer ts.Close()

	s1 := startChallengeSession(t, ts.URL, "learner-deterministic-1")
	s2 := startChallengeSession(t, ts.URL, "learner-deterministic-2")

	hash1 := runScenarioOverWS(t, ts.URL, s1.SessionID)
	hash2 := runScenarioOverWS(t, ts.URL, s2.SessionID)

	if hash1 != hash2 {
		t.Fatalf("expected deterministic hash match, got %s vs %s", hash1, hash2)
	}
}

func TestCommandValidationAndSequenceOrdering(t *testing.T) {
	ts := httptest.NewServer(NewServer(NewSessionManager()).Handler())
	defer ts.Close()

	s := startChallengeSession(t, ts.URL, "learner-sequence")
	conn := dialWS(t, ts.URL, s.SessionID)
	defer func() { _ = conn.Close() }()

	connected := mustReadEvent(t, conn)
	if connected.Type != "session.snapshot" {
		t.Fatalf("expected connected snapshot event, got %s", connected.Type)
	}
	if connected.Sequence == 0 {
		t.Fatalf("expected positive sequence, got %d", connected.Sequence)
	}

	mustWriteCommand(t, conn, Command{Name: "run", Count: 0})
	errEvent := mustReadEvent(t, conn)
	if errEvent.Type != "session.error" {
		t.Fatalf("expected session.error, got %s", errEvent.Type)
	}
	if errEvent.Sequence != connected.Sequence+1 {
		t.Fatalf("expected monotonic sequence after invalid command, got connected=%d error=%d", connected.Sequence, errEvent.Sequence)
	}

	mustWriteCommand(t, conn, Command{Name: "step", Count: 1})
	snapEvent := mustReadEvent(t, conn)
	if snapEvent.Type != "session.snapshot" {
		t.Fatalf("expected session.snapshot, got %s", snapEvent.Type)
	}
	if snapEvent.Sequence != errEvent.Sequence+1 {
		t.Fatalf("expected monotonic sequence after valid command, got error=%d snapshot=%d", errEvent.Sequence, snapEvent.Sequence)
	}
}

func TestCORSPreflightAndHeaders(t *testing.T) {
	t.Setenv("CORS_ALLOW_ORIGIN", "http://localhost:5173")

	ts := httptest.NewServer(NewServer(NewSessionManager()).Handler())
	defer ts.Close()

	req, err := http.NewRequest(http.MethodOptions, ts.URL+"/challenges/start", nil)
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

	req, err := http.NewRequest(http.MethodGet, ts.URL+"/challenges/start", nil)
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

func TestLegacyChallengeGradeEndpointRemoved(t *testing.T) {
	ts := httptest.NewServer(NewServer(NewSessionManager()).Handler())
	defer ts.Close()

	resp, err := http.Post(ts.URL+"/challenges/grade", "application/json", bytes.NewReader([]byte(`{"attempt_id":"a-1"}`)))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("status=%d want=%d", resp.StatusCode, http.StatusNotFound)
	}
}

func startChallengeSession(t *testing.T, baseURL, learnerID string) ChallengeStartResponse {
	t.Helper()
	b, err := json.Marshal(map[string]any{
		"lesson_id":   "l01-sched-rr-basics",
		"stage_index": 0,
		"learner_id":  learnerID,
	})
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	resp, err := http.Post(baseURL+"/challenges/start", "application/json", bytes.NewReader(b))
	if err != nil {
		t.Fatalf("challenge start request failed: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status %d got %d", http.StatusOK, resp.StatusCode)
	}
	var out ChallengeStartResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		t.Fatalf("decode challenge start response failed: %v", err)
	}
	if out.SessionID == "" {
		t.Fatalf("invalid challenge start response: %+v", out)
	}
	return out
}

func runScenarioOverWS(t *testing.T, baseURL, sessionID string) string {
	t.Helper()
	conn := dialWS(t, baseURL, sessionID)
	defer func() { _ = conn.Close() }()

	_ = mustReadEvent(t, conn)

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
