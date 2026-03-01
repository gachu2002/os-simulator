package challenges

import (
	"testing"

	"os-simulator-plan/internal/lessons"
	"os-simulator-plan/internal/sim"
)

func TestBuildValidatorResultsIncludesExpectedActualValues(t *testing.T) {
	stage := lessons.Stage{
		Validators: []lessons.ValidatorSpec{
			{Name: "trace-eq", Type: "trace_count_eq", Values: []string{"proc.dispatch"}, Number: 2},
			{Name: "fs", Type: "fs_ok"},
		},
	}
	results := []lessons.ValidationResult{
		{Name: "trace-eq", Type: "trace_count_eq", Passed: true},
		{Name: "fs", Type: "fs_ok", Passed: false},
	}
	output := lessons.StageOutput{
		Trace: []sim.TraceEvent{
			{Kind: "proc.dispatch"},
			{Kind: "proc.dispatch"},
		},
		FilesystemOK: false,
	}

	views := BuildValidatorResults(results, stage, output)
	if len(views) != 2 {
		t.Fatalf("views len=%d want=2", len(views))
	}
	if views[0].Expected != "2" || views[0].Actual != "2" {
		t.Fatalf("trace expected/actual=%q/%q want=2/2", views[0].Expected, views[0].Actual)
	}
	if views[1].Expected != "true" || views[1].Actual != "false" {
		t.Fatalf("fs expected/actual=%q/%q want=true/false", views[1].Expected, views[1].Actual)
	}
}

func TestBuildValidatorResultsUsesFallbackForUnknownValidatorSpec(t *testing.T) {
	stage := lessons.Stage{}
	results := []lessons.ValidationResult{{Name: "missing", Type: "trace_contains_all", Passed: false}}

	views := BuildValidatorResults(results, stage, lessons.StageOutput{})
	if len(views) != 1 {
		t.Fatalf("views len=%d want=1", len(views))
	}
	if views[0].Expected != "contains all: any" {
		t.Fatalf("expected=%q want=%q", views[0].Expected, "contains all: any")
	}
}
