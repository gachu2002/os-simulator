package lessons

import (
	"fmt"
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
	orderedIDs := []string{
		"l01-sched-rr-basics",
		"l02-sched-fifo-baseline",
		"l03-sched-mlfq-balance",
		"l04-response-under-rr",
		"l05-throughput-shared-cpu",
		"l06-preemption-check",
		"l07-vm-fault-sequence",
		"l08-vm-pressure-repeat",
		"l09-vm-tlb-activity",
		"l10-vm-replacement-fifo",
		"l11-vm-mixed-access",
		"l12-irq-wakeup-read",
		"l13-terminal-write-irq",
		"l14-sleep-wakeup",
		"l15-mixed-blocking",
		"l16-fs-open-traversal",
	}
	for _, id := range orderedIDs {
		if err := runLessonPath(e, id); err != nil {
			t.Fatalf("run lesson path %s failed: %v", id, err)
		}
	}

	for _, stageID := range []string{"s1", "s2", "s3"} {
		if !e.progress.Get("l16-fs-open-traversal", stageID).Completed {
			t.Fatalf("expected stage %s completion recorded", stageID)
		}
	}
}

func TestCompletionAnalyticsAndPilotChecklist(t *testing.T) {
	e := NewEngine()
	for _, id := range []string{"l01-sched-rr-basics", "l02-sched-fifo-baseline", "l03-sched-mlfq-balance", "l04-response-under-rr"} {
		if err := runLessonPath(e, id); err != nil {
			t.Fatalf("run lesson path %s failed: %v", id, err)
		}
	}

	a := e.CompletionAnalytics()
	if a.TotalStages != 60 {
		t.Fatalf("total stages=%d want=60", a.TotalStages)
	}
	if a.CompletedStages != 12 {
		t.Fatalf("completed stages=%d want=12", a.CompletedStages)
	}
	if len(a.ModuleBreakdown) != 4 {
		t.Fatalf("module breakdown count=%d want=4", len(a.ModuleBreakdown))
	}
	if a.PilotChecklistOK {
		t.Fatalf("pilot checklist should not be fully complete for partial run")
	}
}

func TestPrerequisiteGateBlocksOutOfOrderStage(t *testing.T) {
	e := NewEngine()
	if _, err := e.RunStage("l07-vm-fault-sequence", 0); err == nil {
		t.Fatalf("expected prerequisite failure for l07 stage s1")
	}
}

func TestProgressPersistenceRoundTrip(t *testing.T) {
	store := &fakeProgressPersistence{snapshot: map[string]StageProgress{}}
	e := NewEngineWithCatalogAndPersistence(DefaultCatalog(), store)

	for _, lessonID := range []string{"l01-sched-rr-basics", "l02-sched-fifo-baseline", "l03-sched-mlfq-balance", "l04-response-under-rr", "l05-throughput-shared-cpu", "l06-preemption-check"} {
		for idx := 0; idx < 3; idx++ {
			if _, err := e.RunStage(lessonID, idx); err != nil {
				t.Fatalf("run stage %s[%d] failed: %v", lessonID, idx, err)
			}
		}
	}

	reloaded := NewEngineWithCatalogAndPersistence(DefaultCatalog(), store)
	if _, err := reloaded.RunStage("l07-vm-fault-sequence", 0); err != nil {
		t.Fatalf("expected l07 stage 0 to be unlocked after reload: %v", err)
	}
}

func runLessonPath(e *Engine, lessonID string) error {
	for idx := 0; idx < 3; idx++ {
		res, err := e.RunStage(lessonID, idx)
		if err != nil {
			return err
		}
		if !res.Passed {
			return fmt.Errorf("stage index %d failed with feedback %s", idx, res.FeedbackKey)
		}
	}
	return nil
}

type fakeProgressPersistence struct {
	snapshot map[string]StageProgress
}

func (f *fakeProgressPersistence) Load() (map[string]StageProgress, error) {
	out := make(map[string]StageProgress, len(f.snapshot))
	for key, stage := range f.snapshot {
		out[key] = stage
	}
	return out, nil
}

func (f *fakeProgressPersistence) Save(stages map[string]StageProgress) error {
	f.snapshot = make(map[string]StageProgress, len(stages))
	for key, stage := range stages {
		f.snapshot[key] = stage
	}
	return nil
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
