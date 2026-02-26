package lessons

import (
	"fmt"
	"sort"
	"strings"

	"os-simulator-plan/internal/sim"
)

type StageProgress struct {
	Attempts  int
	Completed bool
}

type ProgressStore struct {
	stages map[string]*StageProgress
}

func NewProgressStore() *ProgressStore {
	return &ProgressStore{stages: map[string]*StageProgress{}}
}

func (p *ProgressStore) key(lessonID, stageID string) string {
	return lessonID + ":" + stageID
}

func (p *ProgressStore) Record(lessonID, stageID string, passed bool) StageProgress {
	k := p.key(lessonID, stageID)
	if _, ok := p.stages[k]; !ok {
		p.stages[k] = &StageProgress{}
	}
	sp := p.stages[k]
	sp.Attempts++
	if passed {
		sp.Completed = true
	}
	return *sp
}

func (p *ProgressStore) Get(lessonID, stageID string) StageProgress {
	if v, ok := p.stages[p.key(lessonID, stageID)]; ok {
		return *v
	}
	return StageProgress{}
}

type Engine struct {
	catalog  map[string]Lesson
	progress *ProgressStore
}

func NewEngine() *Engine {
	return NewEngineWithCatalog(DefaultCatalog())
}

func NewEngineWithCatalog(catalog map[string]Lesson) *Engine {
	copyCatalog := make(map[string]Lesson, len(catalog))
	for id, lesson := range catalog {
		copyCatalog[id] = lesson
	}
	return &Engine{catalog: copyCatalog, progress: NewProgressStore()}
}

func (e *Engine) Lessons() []Lesson {
	out := make([]Lesson, 0, len(e.catalog))
	ids := make([]string, 0, len(e.catalog))
	for id := range e.catalog {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	for _, id := range ids {
		out = append(out, e.catalog[id])
	}
	return out
}

func (e *Engine) RunStage(lessonID string, stageIndex int) (StageResult, error) {
	lesson, ok := e.catalog[lessonID]
	if !ok {
		return StageResult{}, fmt.Errorf("lesson %q not found", lessonID)
	}
	if stageIndex < 0 || stageIndex >= len(lesson.Stages) {
		return StageResult{}, fmt.Errorf("invalid stage index %d", stageIndex)
	}
	stage := lesson.Stages[stageIndex]
	for _, prereq := range stage.Prerequisites {
		if !e.isPrerequisiteCompleted(prereq) {
			return StageResult{}, fmt.Errorf("prerequisite %q not completed", prereq)
		}
	}

	output, err := executeStage(stage)
	if err != nil {
		return StageResult{}, err
	}

	for _, v := range stage.Validators {
		ok, _ := validate(output, v)
		if !ok {
			prog := e.progress.Record(lesson.ID, stage.ID, false)
			hintLevel, hint := hintForAttempt(stage.Hints, prog.Attempts)
			return StageResult{Passed: false, FeedbackKey: "validator." + v.Name, Hint: hint, HintLevel: hintLevel, Output: output}, nil
		}
	}

	e.progress.Record(lesson.ID, stage.ID, true)
	return StageResult{Passed: true, FeedbackKey: "stage." + stage.ID + ".passed", Output: output}, nil
}

func (e *Engine) isPrerequisiteCompleted(key string) bool {
	lessonID, stageID, ok := strings.Cut(key, ":")
	if !ok || lessonID == "" || stageID == "" {
		return false
	}
	return e.progress.Get(lessonID, stageID).Completed
}

func executeStage(stage Stage) (StageOutput, error) {
	engine := sim.NewEngine(stage.Config.Seed, 0)
	engine.ConfigureMemory(stage.Config.Frames, stage.Config.TLBEntries)
	engine.ConfigureDevices(stage.Config.DiskLatency, stage.Config.TerminalLatency)
	if err := engine.SetSchedulingPolicy(stage.Config.Policy, stage.Config.Quantum); err != nil {
		return StageOutput{}, err
	}
	if err := engine.ExecuteAll(stage.Commands); err != nil {
		return StageOutput{}, err
	}
	fsOK := engine.ValidateFilesystem() == nil
	return StageOutput{
		Trace:        engine.Trace(),
		Processes:    engine.ProcessTable(),
		Metrics:      engine.SchedulingMetrics(),
		Memory:       engine.MemoryView(),
		FilesystemOK: fsOK,
	}, nil
}

func hintForAttempt(h HintSet, attempts int) (int, string) {
	if attempts <= 1 {
		return 1, h.Nudge
	}
	if attempts == 2 {
		return 2, h.Concept
	}
	return 3, h.Explicit
}

func (e *Engine) CompletionAnalytics() CompletionAnalytics {
	total := 0
	completed := 0
	attempted := 0

	for _, lesson := range e.catalog {
		for _, stage := range lesson.Stages {
			total++
			prog := e.progress.Get(lesson.ID, stage.ID)
			if prog.Completed {
				completed++
			}
			if prog.Attempts > 0 {
				attempted++
			}
		}
	}

	analytics := CompletionAnalytics{TotalStages: total, CompletedStages: completed, AttemptedStages: attempted}
	if total > 0 {
		analytics.CompletionRate = float64(completed) / float64(total)
	}
	return analytics
}
