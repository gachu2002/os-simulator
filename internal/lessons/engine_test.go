package lessons

import (
	"testing"

	"os-simulator-plan/internal/sim"
)

func TestDefaultCatalogHasTwentyLessons(t *testing.T) {
	e := NewEngine()
	if got := len(e.Lessons()); got != 20 {
		t.Fatalf("lesson count=%d want=20", got)
	}
}

func TestScenarioLessonsPassWithExpectedFeedbackKeys(t *testing.T) {
	e := NewEngine()
	ids := []string{
		"l01-sched-rr-basics",
		"l04-response-under-rr",
		"l07-vm-fault-sequence",
		"l09-vm-tlb-activity",
		"l12-irq-wakeup-read",
		"l16-fs-open-traversal",
		"l20-fs-invariants",
	}
	for _, id := range ids {
		res, err := e.RunStage(id, 0)
		if err != nil {
			t.Fatalf("run stage %s failed: %v", id, err)
		}
		if !res.Passed {
			t.Fatalf("lesson %s should pass, got feedback=%s", id, res.FeedbackKey)
		}
		if res.FeedbackKey != "stage.s1.passed" {
			t.Fatalf("lesson %s feedback=%s want=stage.s1.passed", id, res.FeedbackKey)
		}
	}
}

func TestCompletionAnalyticsAndPilotChecklist(t *testing.T) {
	e := NewEngine()
	for _, id := range []string{"l01-sched-rr-basics", "l07-vm-fault-sequence", "l12-irq-wakeup-read", "l16-fs-open-traversal"} {
		if _, err := e.RunStage(id, 0); err != nil {
			t.Fatalf("run stage %s failed: %v", id, err)
		}
	}

	a := e.CompletionAnalytics()
	if a.TotalStages != 20 {
		t.Fatalf("total stages=%d want=20", a.TotalStages)
	}
	if a.CompletedStages != 4 {
		t.Fatalf("completed stages=%d want=4", a.CompletedStages)
	}
	if len(a.ModuleBreakdown) != 4 {
		t.Fatalf("module breakdown count=%d want=4", len(a.ModuleBreakdown))
	}
	if a.PilotChecklistOK {
		t.Fatalf("pilot checklist should not be fully complete for partial run")
	}
}

func TestHintProgressionLevels(t *testing.T) {
	e := NewEngine()
	e.catalog["fail-lesson"] = Lesson{
		ID:    "fail-lesson",
		Title: "Failing Lesson",
		Stages: []Stage{{
			ID:         "s1",
			Title:      "always fail",
			Config:     SimConfig{Seed: 1, Policy: "rr", Quantum: 2, Frames: 8, TLBEntries: 4, DiskLatency: 3, TerminalLatency: 1},
			Commands:   []sim.Command{{Name: "step", Count: 1}},
			Validators: []ValidatorSpec{{Name: "impossible", Type: "metric_eq", Key: "completed_processes", Number: 99}},
			Hints:      HintSet{Nudge: "nudge", Concept: "concept", Explicit: "explicit"},
		}},
	}

	r1, _ := e.RunStage("fail-lesson", 0)
	r2, _ := e.RunStage("fail-lesson", 0)
	r3, _ := e.RunStage("fail-lesson", 0)
	if r1.Passed || r2.Passed || r3.Passed {
		t.Fatalf("failing lesson should not pass")
	}
	if r1.HintLevel != 1 || r1.Hint != "nudge" {
		t.Fatalf("attempt1 hint mismatch: level=%d hint=%q", r1.HintLevel, r1.Hint)
	}
	if r2.HintLevel != 2 || r2.Hint != "concept" {
		t.Fatalf("attempt2 hint mismatch: level=%d hint=%q", r2.HintLevel, r2.Hint)
	}
	if r3.HintLevel != 3 || r3.Hint != "explicit" {
		t.Fatalf("attempt3 hint mismatch: level=%d hint=%q", r3.HintLevel, r3.Hint)
	}
}
