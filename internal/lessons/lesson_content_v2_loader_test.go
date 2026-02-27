package lessons

import "testing"

func TestLoadLessonContentV2(t *testing.T) {
	content, err := loadLessonContentV2()
	if err != nil {
		t.Fatalf("load lesson content v2: %v", err)
	}
	if len(content) != 28 {
		t.Fatalf("lesson content count=%d want=28", len(content))
	}
	record, ok := content["l01-sched-rr-basics"]
	if !ok {
		t.Fatalf("missing l01 content")
	}
	if record.Learn.CoreIdea == "" || len(record.Learn.MechanismSteps) < 4 {
		t.Fatalf("l01 content missing learn essentials")
	}
}

func TestCatalogStagesIncludeV2LearnContent(t *testing.T) {
	e := NewEngine()
	lesson := e.catalog["l01-sched-rr-basics"]
	if len(lesson.Stages) == 0 {
		t.Fatalf("expected lesson stages")
	}
	stage := lesson.Stages[0]
	if stage.CoreIdea == "" {
		t.Fatalf("expected core idea on stage")
	}
	if len(stage.MechanismSteps) == 0 {
		t.Fatalf("expected mechanism steps on stage")
	}
	if len(stage.CommonMistakes) < 3 {
		t.Fatalf("expected common mistakes on stage")
	}
	if len(stage.PreChallengeChecklist) < 3 {
		t.Fatalf("expected pre challenge checklist on stage")
	}
}
