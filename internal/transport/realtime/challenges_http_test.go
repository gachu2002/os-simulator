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

func TestChallengeStartAndSubmitLifecycle(t *testing.T) {
	ts := httptest.NewServer(NewServer(NewSessionManager()).Handler())
	defer ts.Close()

	startReq := ChallengeStartRequest{LessonID: "l01-sched-rr-basics", StageIndex: 0}
	startRes := startChallenge(t, ts.URL, startReq)

	if startRes.AttemptID == "" {
		t.Fatalf("expected attempt id")
	}
	if startRes.SessionID == "" {
		t.Fatalf("expected session id")
	}
	if startRes.Objective == "" {
		t.Fatalf("expected objective")
	}
	if len(startRes.AllowedCommands) == 0 {
		t.Fatalf("expected allowed commands")
	}
	if startRes.Limits.MaxSteps <= 0 {
		t.Fatalf("expected challenge max steps")
	}

	firstSubmit := submitChallenge(t, ts.URL, ChallengeGradeRequest{AttemptID: startRes.AttemptID})
	if firstSubmit.Passed {
		t.Fatalf("expected first submit to fail without interaction")
	}

	conn := dialChallengeWS(t, ts.URL, startRes.SessionID)
	defer func() { _ = conn.Close() }()

	_ = mustReadEvent(t, conn)
	mustWriteCommand(t, conn, Command{Name: "step", Count: 8})
	_ = mustReadEvent(t, conn)

	secondSubmit := submitChallenge(t, ts.URL, ChallengeGradeRequest{AttemptID: startRes.AttemptID})
	if !secondSubmit.Passed {
		t.Fatalf("expected second submit to pass, feedback=%s", secondSubmit.FeedbackKey)
	}
	if len(secondSubmit.PassConditions) == 0 {
		t.Fatalf("expected pass conditions in submit response")
	}
	if len(secondSubmit.ValidatorResults) == 0 {
		t.Fatalf("expected validator results in submit response")
	}
	firstValidator := secondSubmit.ValidatorResults[0]
	if firstValidator.Expected == "" || firstValidator.Actual == "" {
		t.Fatalf("expected validator result to include expected/actual values")
	}
	if secondSubmit.Output.TraceHash == "" {
		t.Fatalf("expected trace hash")
	}
}

func TestChallengeStartAndSubmitValidationErrors(t *testing.T) {
	ts := httptest.NewServer(NewServer(NewSessionManager()).Handler())
	defer ts.Close()

	resp := postJSON(t, ts.URL+"/challenges/start", map[string]any{"lesson_id": "", "stage_index": 0})
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("start status=%d want=%d", resp.StatusCode, http.StatusBadRequest)
	}

	missing := submitChallengeRaw(t, ts.URL, ChallengeGradeRequest{AttemptID: "a-999999"})
	defer func() { _ = missing.Body.Close() }()
	if missing.StatusCode != http.StatusNotFound {
		t.Fatalf("submit status=%d want=%d", missing.StatusCode, http.StatusNotFound)
	}
}

func TestChallengeSessionCommandLimitsAreEnforced(t *testing.T) {
	ts := httptest.NewServer(NewServer(NewSessionManager()).Handler())
	defer ts.Close()

	startRes := startChallenge(t, ts.URL, ChallengeStartRequest{LessonID: "l01-sched-rr-basics", StageIndex: 0})

	conn := dialChallengeWS(t, ts.URL, startRes.SessionID)
	defer func() { _ = conn.Close() }()

	_ = mustReadEvent(t, conn)

	mustWriteCommand(t, conn, Command{Name: "step", Count: startRes.Limits.MaxSteps + 1})
	limitErr := mustReadEvent(t, conn)
	if limitErr.Type != "session.error" {
		t.Fatalf("event type=%s want=session.error", limitErr.Type)
	}
	if !strings.Contains(limitErr.Error, "step limit exceeded") {
		t.Fatalf("expected step limit error, got %q", limitErr.Error)
	}

	policyAllowed := false
	for _, name := range startRes.AllowedCommands {
		if name == "policy" {
			policyAllowed = true
			break
		}
	}
	if !policyAllowed || startRes.Limits.MaxPolicyChanges <= 0 {
		return
	}

	for idx := 0; idx < startRes.Limits.MaxPolicyChanges; idx++ {
		mustWriteCommand(t, conn, Command{Name: "policy", Policy: "rr", Quantum: 2})
		ev := mustReadEvent(t, conn)
		if ev.Type != "session.snapshot" {
			t.Fatalf("policy event type=%s want=session.snapshot", ev.Type)
		}
	}

	mustWriteCommand(t, conn, Command{Name: "policy", Policy: "rr", Quantum: 2})
	policyErr := mustReadEvent(t, conn)
	if policyErr.Type != "session.error" {
		t.Fatalf("event type=%s want=session.error", policyErr.Type)
	}
	if !strings.Contains(policyErr.Error, "policy change limit exceeded") {
		t.Fatalf("expected policy limit error, got %q", policyErr.Error)
	}
}

func TestChallengeAttemptIsScopedToLearner(t *testing.T) {
	ts := httptest.NewServer(NewServer(NewSessionManager()).Handler())
	defer ts.Close()

	startRes := startChallenge(t, ts.URL, ChallengeStartRequest{
		LessonID:   "l01-sched-rr-basics",
		StageIndex: 0,
		LearnerID:  "learner-a",
	})

	resp := submitChallengeRaw(t, ts.URL, ChallengeGradeRequest{AttemptID: startRes.AttemptID, LearnerID: "learner-b"})
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusForbidden {
		t.Fatalf("status=%d want=%d", resp.StatusCode, http.StatusForbidden)
	}
}

func startChallenge(t *testing.T, baseURL string, req ChallengeStartRequest) ChallengeStartResponse {
	t.Helper()
	resp := postJSON(t, baseURL+"/challenges/start", req)
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status=%d want=%d", resp.StatusCode, http.StatusOK)
	}
	var out ChallengeStartResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		t.Fatalf("decode start response failed: %v", err)
	}
	return out
}

func submitChallenge(t *testing.T, baseURL string, req ChallengeGradeRequest) ChallengeGradeResponse {
	t.Helper()
	resp := submitChallengeRaw(t, baseURL, req)
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status=%d want=%d", resp.StatusCode, http.StatusOK)
	}
	var out ChallengeGradeResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		t.Fatalf("decode submit response failed: %v", err)
	}
	return out
}

func submitChallengeRaw(t *testing.T, baseURL string, req ChallengeGradeRequest) *http.Response {
	t.Helper()
	return postJSON(t, baseURL+"/challenges/submit", req)
}

func postJSON(t *testing.T, url string, payload any) *http.Response {
	t.Helper()
	b, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("marshal request failed: %v", err)
	}
	resp, err := http.Post(url, "application/json", bytes.NewReader(b))
	if err != nil {
		t.Fatalf("post request failed: %v", err)
	}
	return resp
}

func dialChallengeWS(t *testing.T, baseURL, sessionID string) *websocket.Conn {
	t.Helper()
	wsURL := "ws" + strings.TrimPrefix(baseURL, "http") + "/ws/" + sessionID
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("dial websocket failed: %v", err)
	}
	return conn
}
