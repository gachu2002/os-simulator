package realtime

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type responseContract struct {
	RequiredFields  []string `json:"required_fields"`
	ForbiddenFields []string `json:"forbidden_fields"`
}

func TestChallengeStartV3Lifecycle(t *testing.T) {
	ts := httptest.NewServer(NewServer(NewSessionManager()).Handler())
	defer ts.Close()

	res := startChallengeV3(t, ts.URL, map[string]any{
		"lesson_id":  "l01-process-basics",
		"learner_id": "learner-v3",
	})

	if res["version"] != "v3" {
		t.Fatalf("version=%v want=v3", res["version"])
	}
	if res["lesson_id"] != "l01-process-basics" {
		t.Fatalf("lesson_id=%v want=l01-process-basics", res["lesson_id"])
	}
	actionCaps, ok := res["action_capabilities"].(map[string]any)
	if !ok {
		t.Fatalf("expected action_capabilities in start v3 response")
	}
	if len(actionCaps["supported_now"].([]any)) == 0 {
		t.Fatalf("expected supported_now actions in start v3 response")
	}
	actionNotes, ok := res["action_capability_notes"].(map[string]any)
	if !ok || len(actionNotes) == 0 {
		t.Fatalf("expected action_capability_notes in start v3 response")
	}
	if res["attempt_id"] == "" || res["session_id"] == "" {
		t.Fatalf("expected attempt_id and session_id")
	}
	assertResponseContract(t, "challenge_start_v3_contract.json", res)
}

func TestChallengeStartV3Validation(t *testing.T) {
	ts := httptest.NewServer(NewServer(NewSessionManager()).Handler())
	defer ts.Close()

	missingPart := postJSON(t, ts.URL+"/challenges/start/v3", map[string]any{"lesson_id": "l03-limited-direct-execution"})
	defer func() { _ = missingPart.Body.Close() }()
	if missingPart.StatusCode != http.StatusBadRequest {
		t.Fatalf("status=%d want=%d", missingPart.StatusCode, http.StatusBadRequest)
	}

	invalidPart := postJSON(t, ts.URL+"/challenges/start/v3", map[string]any{"lesson_id": "l03-limited-direct-execution", "part_id": "C"})
	defer func() { _ = invalidPart.Body.Close() }()
	if invalidPart.StatusCode != http.StatusBadRequest {
		t.Fatalf("status=%d want=%d", invalidPart.StatusCode, http.StatusBadRequest)
	}
}

func TestChallengeActionV3AppliesMappedStep(t *testing.T) {
	ts := httptest.NewServer(NewServer(NewSessionManager()).Handler())
	defer ts.Close()

	startRes := startChallengeV3(t, ts.URL, map[string]any{
		"lesson_id":  "l01-process-basics",
		"learner_id": "learner-action-v3",
	})

	actionRes := actionChallengeV3(t, ts.URL, map[string]any{
		"attempt_id": startRes["attempt_id"],
		"learner_id": "learner-action-v3",
		"action":     "execute_instruction",
		"count":      2,
	})

	if actionRes["mapped_command"] != "step" {
		t.Fatalf("mapped_command=%v want=step", actionRes["mapped_command"])
	}
	ev, ok := actionRes["event"].(map[string]any)
	if !ok {
		t.Fatalf("missing event in action response")
	}
	if ev["type"] != "session.snapshot" {
		t.Fatalf("event type=%v want=session.snapshot", ev["type"])
	}
}

func TestChallengeActionV3ExecAndWaitMappings(t *testing.T) {
	ts := httptest.NewServer(NewServer(NewSessionManager()).Handler())
	defer ts.Close()

	startRes := startChallengeV3(t, ts.URL, map[string]any{
		"lesson_id":  "l01-process-basics",
		"learner_id": "learner-action-v3-exec-wait",
	})

	execRes := actionChallengeV3(t, ts.URL, map[string]any{
		"attempt_id": startRes["attempt_id"],
		"learner_id": "learner-action-v3-exec-wait",
		"action":     "exec",
		"process":    "shell-child",
		"program":    "COMPUTE 1; EXIT",
	})
	if execRes["mapped_command"] != "spawn" {
		t.Fatalf("mapped_command=%v want=spawn", execRes["mapped_command"])
	}

	waitRes := actionChallengeV3(t, ts.URL, map[string]any{
		"attempt_id": startRes["attempt_id"],
		"learner_id": "learner-action-v3-exec-wait",
		"action":     "wait",
		"count":      2,
	})
	if waitRes["mapped_command"] != "run" {
		t.Fatalf("mapped_command=%v want=run", waitRes["mapped_command"])
	}
}

func TestChallengeActionV3RejectsUnsupportedAction(t *testing.T) {
	ts := httptest.NewServer(NewServer(NewSessionManager()).Handler())
	defer ts.Close()

	startRes := startChallengeV3(t, ts.URL, map[string]any{
		"lesson_id":  "l01-process-basics",
		"learner_id": "learner-action-v3-unsupported",
	})

	resp := postJSON(t, ts.URL+"/challenges/action/v3", map[string]any{
		"attempt_id": startRes["attempt_id"],
		"learner_id": "learner-action-v3-unsupported",
		"action":     "migrate_job",
	})
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("status=%d want=%d", resp.StatusCode, http.StatusBadRequest)
	}
	var out map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		t.Fatalf("decode error response failed: %v", err)
	}
	if out["code"] != "invalid_action" {
		t.Fatalf("code=%v want=invalid_action", out["code"])
	}
	msg, _ := out["message"].(string)
	if !strings.Contains(msg, "fallback_action=step") {
		t.Fatalf("message=%q want fallback guidance", msg)
	}
}

func TestChallengeActionV3BlockProcessMapping(t *testing.T) {
	ts := httptest.NewServer(NewServer(NewSessionManager()).Handler())
	defer ts.Close()

	startRes := startChallengeV3(t, ts.URL, map[string]any{
		"lesson_id":  "l01-process-basics",
		"learner_id": "learner-action-v3-block",
	})

	actionRes := actionChallengeV3(t, ts.URL, map[string]any{
		"attempt_id": startRes["attempt_id"],
		"learner_id": "learner-action-v3-block",
		"action":     "block_process",
	})

	if actionRes["mapped_command"] != "block_process" {
		t.Fatalf("mapped_command=%v want=block_process", actionRes["mapped_command"])
	}
}

func TestChallengeSubmitV3AndReplayV3(t *testing.T) {
	ts := httptest.NewServer(NewServer(NewSessionManager()).Handler())
	defer ts.Close()

	startRes := startChallengeV3(t, ts.URL, map[string]any{
		"lesson_id":  "l01-process-basics",
		"learner_id": "learner-submit-v3",
	})

	_ = actionChallengeV3(t, ts.URL, map[string]any{
		"attempt_id": startRes["attempt_id"],
		"learner_id": "learner-submit-v3",
		"action":     "run",
		"count":      8,
	})

	submitResp := postJSON(t, ts.URL+"/challenges/submit/v3", map[string]any{
		"attempt_id": startRes["attempt_id"],
		"learner_id": "learner-submit-v3",
	})
	defer func() { _ = submitResp.Body.Close() }()
	if submitResp.StatusCode != http.StatusOK {
		t.Fatalf("submit v3 status=%d want=%d", submitResp.StatusCode, http.StatusOK)
	}
	var submit map[string]any
	if err := json.NewDecoder(submitResp.Body).Decode(&submit); err != nil {
		t.Fatalf("decode submit v3 failed: %v", err)
	}
	if submit["version"] != "v3" {
		t.Fatalf("version=%v want=v3", submit["version"])
	}
	assertResponseContract(t, "challenge_submit_v3_contract.json", submit)

	replayURL := ts.URL + "/challenges/" + startRes["attempt_id"].(string) + "/replay/v3?learner_id=learner-submit-v3"
	replayResp, err := http.Get(replayURL)
	if err != nil {
		t.Fatalf("get replay v3 failed: %v", err)
	}
	defer func() { _ = replayResp.Body.Close() }()
	if replayResp.StatusCode != http.StatusOK {
		t.Fatalf("replay v3 status=%d want=%d", replayResp.StatusCode, http.StatusOK)
	}
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

func startChallengeV3(t *testing.T, baseURL string, payload map[string]any) map[string]any {
	t.Helper()
	resp := postJSON(t, baseURL+"/challenges/start/v3", payload)
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status=%d want=%d", resp.StatusCode, http.StatusOK)
	}
	var out map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		t.Fatalf("decode v3 start response failed: %v", err)
	}
	return out
}

func actionChallengeV3(t *testing.T, baseURL string, payload map[string]any) map[string]any {
	t.Helper()
	resp := postJSON(t, baseURL+"/challenges/action/v3", payload)
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status=%d want=%d", resp.StatusCode, http.StatusOK)
	}
	var out map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		t.Fatalf("decode v3 action response failed: %v", err)
	}
	return out
}

func assertResponseContract(t *testing.T, fileName string, payload map[string]any) {
	t.Helper()

	b, err := os.ReadFile(filepath.Join("..", "..", "..", "contracts", fileName))
	if err != nil {
		t.Fatalf("read contract fixture failed: %v", err)
	}

	var contract responseContract
	if err := json.Unmarshal(b, &contract); err != nil {
		t.Fatalf("unmarshal contract fixture failed: %v", err)
	}

	for _, field := range contract.RequiredFields {
		if _, ok := payload[field]; !ok {
			t.Fatalf("missing required field %q", field)
		}
	}

	for _, field := range contract.ForbiddenFields {
		if _, ok := payload[field]; ok {
			t.Fatalf("forbidden legacy field present %q", field)
		}
	}
}
