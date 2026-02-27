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

type ValidatorHint struct {
	Validator string
	Hints     HintSet
}

type ChallengeLimits struct {
	MaxSteps         int
	MaxPolicyChanges int
	MaxConfigChanges int
}

type ActionDescription struct {
	Command     string
	Description string
}

type Stage struct {
	ID                    string
	Title                 string
	CoreIdea              string
	MechanismSteps        []string
	WorkedExample         string
	CommonMistakes        []string
	PreChallengeChecklist []string
	Objective             string
	Goal                  string
	TheoryDetail          string
	Prerequisites         []string
	Config                SimConfig
	Commands              []sim.Command
	Bootstrap             []sim.Command
	AllowedCmds           []string
	ActionDescriptions    []ActionDescription
	ExpectedVisualCues    []string
	Limits                ChallengeLimits
	Validators            []ValidatorSpec
	Hints                 HintSet
	ValidatorHints        []ValidatorHint
}

type Lesson struct {
	ID               string
	Title            string
	Module           string
	SectionID        string
	SectionTitle     string
	Difficulty       string
	EstimatedMinutes int
	ChapterRefs      []string
	Stages           []Stage
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

type ValidationResult struct {
	Name    string
	Type    string
	Key     string
	Passed  bool
	Message string
}

type StageResult struct {
	Passed           bool
	FeedbackKey      string
	Hint             string
	HintLevel        int
	Output           StageOutput
	ValidatorResults []ValidationResult
}

type CompletionAnalytics struct {
	TotalStages     int
	CompletedStages int
	AttemptedStages int
	CompletionRate  float64
}
