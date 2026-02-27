package lessons

import (
	"fmt"
	"testing"

	"os-simulator-plan/internal/sim"
)

func TestDefaultCatalogHasTwentyEightLessons(t *testing.T) {
	e := NewEngine()
	if got := len(e.Lessons()); got != 28 {
		t.Fatalf("lesson count=%d want=28", got)
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
		"l06b-lottery-tradeoffs",
		"l06c-quantum-response-tuning",
		"l07-vm-fault-sequence",
		"l08-vm-pressure-repeat",
		"l09-vm-tlb-activity",
		"l10-vm-replacement-fifo",
		"l11-vm-mixed-access",
		"l11b-vm-locality-window",
		"l11c-vm-fault-burst-diagnose",
		"l12-irq-wakeup-read",
		"l13-terminal-write-irq",
		"l14-sleep-wakeup",
		"l15-mixed-blocking",
		"l15b-irq-interleave-order",
		"l15c-blocked-progress-guarantee",
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

func TestCompletionAnalytics(t *testing.T) {
	e := NewEngine()
	for _, id := range []string{"l01-sched-rr-basics", "l02-sched-fifo-baseline", "l03-sched-mlfq-balance", "l04-response-under-rr"} {
		if err := runLessonPath(e, id); err != nil {
			t.Fatalf("run lesson path %s failed: %v", id, err)
		}
	}

	a := e.CompletionAnalytics()
	if a.TotalStages != 84 {
		t.Fatalf("total stages=%d want=84", a.TotalStages)
	}
	if a.CompletedStages != 12 {
		t.Fatalf("completed stages=%d want=12", a.CompletedStages)
	}
	if a.AttemptedStages != 12 {
		t.Fatalf("attempted stages=%d want=12", a.AttemptedStages)
	}
}

func TestPrerequisiteGateBlocksOutOfOrderStage(t *testing.T) {
	e := NewEngine()
	if _, err := e.RunStage("l07-vm-fault-sequence", 0); err == nil {
		t.Fatalf("expected prerequisite failure for l07 stage s1")
	}
}

func TestPrepareStageReturnsLessonAndStageMetadata(t *testing.T) {
	e := NewEngine()
	prepared, err := e.PrepareStage("l01-sched-rr-basics", 0)
	if err != nil {
		t.Fatalf("prepare stage failed: %v", err)
	}
	if prepared.LessonID != "l01-sched-rr-basics" {
		t.Fatalf("lesson id=%q want=%q", prepared.LessonID, "l01-sched-rr-basics")
	}
	if prepared.Module != "cpu-virtualization" {
		t.Fatalf("module=%q want=%q", prepared.Module, "cpu-virtualization")
	}
	if prepared.Stage.ID != "s1" {
		t.Fatalf("stage id=%q want=%q", prepared.Stage.ID, "s1")
	}
	if prepared.Stage.Objective == "" {
		t.Fatalf("expected stage objective")
	}
	if len(prepared.Stage.AllowedCmds) == 0 {
		t.Fatalf("expected allowed challenge commands")
	}
	for _, cmd := range prepared.Stage.AllowedCmds {
		if cmd == "spawn" {
			t.Fatalf("spawn should be provided via bootstrap, not interactive controls")
		}
	}
	if prepared.Stage.Limits.MaxSteps <= 0 {
		t.Fatalf("expected challenge max steps")
	}
	if len(prepared.Stage.Bootstrap) == 0 {
		t.Fatalf("expected bootstrap commands for challenge stage")
	}
}

func TestPrepareStageRespectsPrerequisites(t *testing.T) {
	e := NewEngine()
	if _, err := e.PrepareStage("l07-vm-fault-sequence", 0); err == nil {
		t.Fatalf("expected prerequisite failure for prepare stage")
	}
}

func TestStageThreeAllowsModuleSpecificConfigActions(t *testing.T) {
	e := NewEngine()

	lesson := e.catalog["l11c-vm-fault-burst-diagnose"]
	stage := lesson.Stages[2]
	if !contains(stage.AllowedCmds, "set_frames") || !contains(stage.AllowedCmds, "set_tlb_entries") {
		t.Fatalf("expected memory stage s3 to allow frame/tlb configuration")
	}
	if stage.Limits.MaxConfigChanges <= 0 {
		t.Fatalf("expected positive config-change limit for memory stage s3")
	}

	ioLesson := e.catalog["l15c-blocked-progress-guarantee"]
	ioStage := ioLesson.Stages[2]
	if !contains(ioStage.AllowedCmds, "set_disk_latency") || !contains(ioStage.AllowedCmds, "set_terminal_latency") {
		t.Fatalf("expected concurrency stage s3 to allow device-latency tuning")
	}
	if ioStage.Limits.MaxConfigChanges <= 0 {
		t.Fatalf("expected positive config-change limit for concurrency stage s3")
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

func TestValidatorSpecificHintsOverrideStageHints(t *testing.T) {
	e := NewEngine()
	e.catalog["validator-hint"] = Lesson{
		ID:    "validator-hint",
		Title: "Validator Hint Lesson",
		Stages: []Stage{{
			ID:       "s1",
			Title:    "validator hint stage",
			Config:   SimConfig{Seed: 2, Policy: "rr", Quantum: 2, Frames: 8, TLBEntries: 4, DiskLatency: 3, TerminalLatency: 1},
			Commands: []sim.Command{{Name: "step", Count: 1}},
			Validators: []ValidatorSpec{
				{Name: "trace-required", Type: "trace_contains_all", Values: []string{"proc.dispatch", "proc.compute", "never.happens"}},
				{Name: "completion", Type: "metric_eq", Key: "completed_processes", Number: 1},
			},
			Hints: HintSet{Nudge: "stage-nudge", Concept: "stage-concept", Explicit: "stage-explicit"},
			ValidatorHints: []ValidatorHint{{
				Validator: "trace-required",
				Hints:     HintSet{Nudge: "validator-nudge", Concept: "validator-concept", Explicit: "validator-explicit"},
			}},
		}},
	}

	r1, _ := e.RunStage("validator-hint", 0)
	r2, _ := e.RunStage("validator-hint", 0)
	r3, _ := e.RunStage("validator-hint", 0)

	if r1.Passed || r2.Passed || r3.Passed {
		t.Fatalf("validator-hint lesson should not pass")
	}
	if r1.Hint != "validator-nudge" || r1.HintLevel != 1 {
		t.Fatalf("attempt1 validator hint mismatch: level=%d hint=%q", r1.HintLevel, r1.Hint)
	}
	if r2.Hint != "validator-concept" || r2.HintLevel != 2 {
		t.Fatalf("attempt2 validator hint mismatch: level=%d hint=%q", r2.HintLevel, r2.Hint)
	}
	if r3.Hint != "validator-explicit" || r3.HintLevel != 3 {
		t.Fatalf("attempt3 validator hint mismatch: level=%d hint=%q", r3.HintLevel, r3.Hint)
	}
}

func contains(items []string, target string) bool {
	for _, item := range items {
		if item == target {
			return true
		}
	}
	return false
}
