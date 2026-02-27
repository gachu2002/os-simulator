package lessons

import "testing"

func TestLoadLessonStageContentV2(t *testing.T) {
	content, err := loadLessonStageContentV2()
	if err != nil {
		t.Fatalf("load lesson stage content v2: %v", err)
	}
	if len(content) != 28 {
		t.Fatalf("stage content count=%d want=28", len(content))
	}
	lessonStages, ok := content["l01-sched-rr-basics"]
	if !ok {
		t.Fatalf("missing l01 stage content")
	}
	s1, ok := lessonStages["s1"]
	if !ok {
		t.Fatalf("missing l01 s1 content")
	}
	if s1.Objective == "" || s1.Goal == "" {
		t.Fatalf("expected non-empty objective and goal")
	}
}

func TestCatalogAppliesStageAuthoringOverrides(t *testing.T) {
	e := NewEngine()
	lesson := e.catalog["l01-sched-rr-basics"]
	stage := lesson.Stages[0]
	if stage.Objective == "" || stage.Goal == "" {
		t.Fatalf("expected stage objective and goal from authored content")
	}
	if stage.Hints.Nudge == "" {
		t.Fatalf("expected stage hints from authored content")
	}
}
