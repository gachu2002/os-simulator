package lessons

import "os-simulator-plan/internal/sim"

type SimConfig struct {
	Seed            uint64
	Policy          string
	Quantum         int
	Frames          int
	TLBEntries      int
	DiskLatency     sim.Tick
	TerminalLatency sim.Tick
}

type ValidatorSpec struct {
	Name   string
	Type   string
	Key    string
	Number float64
	Values []string
}

type HintSet struct {
	Nudge    string
	Concept  string
	Explicit string
}

type ChallengeLimits struct {
	MaxSteps         int
	MaxPolicyChanges int
}

type Stage struct {
	ID            string
	Title         string
	Objective     string
	Prerequisites []string
	Config        SimConfig
	Commands      []sim.Command
	Bootstrap     []sim.Command
	AllowedCmds   []string
	Limits        ChallengeLimits
	Validators    []ValidatorSpec
	Hints         HintSet
}

type Lesson struct {
	ID     string
	Title  string
	Module string
	Stages []Stage
}

type PreparedStage struct {
	LessonID    string
	LessonTitle string
	Module      string
	StageIndex  int
	Stage       Stage
}

type StageOutput struct {
	Trace        []sim.TraceEvent
	Processes    []sim.ProcessSnapshot
	Metrics      sim.SchedulingMetrics
	Memory       sim.MemorySnapshot
	FilesystemOK bool
}

type StageResult struct {
	Passed      bool
	FeedbackKey string
	Hint        string
	HintLevel   int
	Output      StageOutput
}

type CompletionAnalytics struct {
	TotalStages     int
	CompletedStages int
	AttemptedStages int
	CompletionRate  float64
}
