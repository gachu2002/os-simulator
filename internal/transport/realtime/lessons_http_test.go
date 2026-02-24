package realtime

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"os-simulator-plan/internal/lessons"
	"os-simulator-plan/internal/sim"
)

func TestLessonsListEndpoint(t *testing.T) {
	ts := httptest.NewServer(NewServer(NewSessionManager()).Handler())
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/lessons")
	if err != nil {
		t.Fatalf("get lessons failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status=%d want=%d", resp.StatusCode, http.StatusOK)
	}

	var out LessonsResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		t.Fatalf("decode lessons failed: %v", err)
	}
	if len(out.Lessons) != 20 {
		t.Fatalf("lesson count=%d want=20", len(out.Lessons))
	}
	if len(out.Lessons[0].Stages) == 0 {
		t.Fatalf("first lesson should include stage summaries")
	}
}

func TestLessonRunHintProgressionAndAnalytics(t *testing.T) {
	lessonEngine := lessons.NewEngine()
	lessonEngineCatalog := lessonEngine.Lessons()
	_ = lessonEngineCatalog
	lessonEngineFail := injectFailingLesson(lessonEngine)

	ts := httptest.NewServer(NewServerWithLessons(NewSessionManager(), lessonEngineFail).Handler())
	defer ts.Close()

	req := LessonRunRequest{LessonID: "fail-lesson", StageIndex: 0}
	first := runLesson(t, ts.URL, req)
	second := runLesson(t, ts.URL, req)
	third := runLesson(t, ts.URL, req)

	if first.Passed || second.Passed || third.Passed {
		t.Fatalf("failing lesson should not pass")
	}
	if first.HintLevel != 1 || first.Hint == "" {
		t.Fatalf("expected hint level 1 on first attempt")
	}
	if second.HintLevel != 2 {
		t.Fatalf("expected hint level 2 on second attempt, got %d", second.HintLevel)
	}
	if third.HintLevel != 3 {
		t.Fatalf("expected hint level 3 on third attempt, got %d", third.HintLevel)
	}
	if first.Output.TraceHash != second.Output.TraceHash || second.Output.TraceHash != third.Output.TraceHash {
		t.Fatalf("expected deterministic stage hash, got %s vs %s", first.Output.TraceHash, second.Output.TraceHash)
	}
	if third.Analytics.AttemptedStages < first.Analytics.AttemptedStages {
		t.Fatalf("attempted stages regressed: first=%d third=%d", first.Analytics.AttemptedStages, third.Analytics.AttemptedStages)
	}
}

func TestLessonRunValidationErrors(t *testing.T) {
	ts := httptest.NewServer(NewServer(NewSessionManager()).Handler())
	defer ts.Close()

	b, err := json.Marshal(LessonRunRequest{LessonID: "", StageIndex: 0})
	if err != nil {
		t.Fatalf("marshal request failed: %v", err)
	}
	resp, err := http.Post(ts.URL+"/lessons/run", "application/json", bytes.NewReader(b))
	if err != nil {
		t.Fatalf("post run failed: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("status=%d want=%d", resp.StatusCode, http.StatusBadRequest)
	}
}

func runLesson(t *testing.T, baseURL string, req LessonRunRequest) LessonRunResponse {
	t.Helper()
	b, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("marshal request failed: %v", err)
	}
	resp, err := http.Post(baseURL+"/lessons/run", "application/json", bytes.NewReader(b))
	if err != nil {
		t.Fatalf("post run failed: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status=%d want=%d", resp.StatusCode, http.StatusOK)
	}
	var out LessonRunResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		t.Fatalf("decode run response failed: %v", err)
	}
	return out
}

func injectFailingLesson(engine *lessons.Engine) *lessons.Engine {
	catalog := map[string]lessons.Lesson{}
	for _, lesson := range engine.Lessons() {
		catalog[lesson.ID] = lesson
	}
	catalog["fail-lesson"] = lessons.Lesson{
		ID:    "fail-lesson",
		Title: "Failing Lesson",
		Stages: []lessons.Stage{{
			ID:       "s1",
			Title:    "always fail",
			Config:   lessons.SimConfig{Seed: 1, Policy: "rr", Quantum: 2, Frames: 8, TLBEntries: 4, DiskLatency: 3, TerminalLatency: 1},
			Commands: []sim.Command{{Name: "step", Count: 1}},
			Validators: []lessons.ValidatorSpec{{
				Name:   "impossible",
				Type:   "metric_eq",
				Key:    "completed_processes",
				Number: 99,
			}},
			Hints: lessons.HintSet{Nudge: "nudge", Concept: "concept", Explicit: "explicit"},
		}},
	}
	return lessons.NewEngineWithCatalog(catalog)
}
