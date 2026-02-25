package lessons

import (
	"cmp"
	"fmt"
	"sort"
	"strings"

	"os-simulator-plan/internal/sim"
)

type StageProgress struct {
	Attempts         int  `json:"attempts"`
	Completed        bool `json:"completed"`
	HighestHintLevel int  `json:"highest_hint_level"`
}

type ProgressStore struct {
	stages map[string]*StageProgress
}

type ProgressPersistence interface {
	Load() (map[string]StageProgress, error)
	Save(stages map[string]StageProgress) error
}

func NewProgressStore() *ProgressStore {
	return &ProgressStore{stages: map[string]*StageProgress{}}
}

func NewProgressStoreFromSnapshot(snapshot map[string]StageProgress) *ProgressStore {
	stages := make(map[string]*StageProgress, len(snapshot))
	for key, stage := range snapshot {
		copy := stage
		stages[key] = &copy
	}
	return &ProgressStore{stages: stages}
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

func (p *ProgressStore) Snapshot() map[string]StageProgress {
	out := make(map[string]StageProgress, len(p.stages))
	for key, stage := range p.stages {
		out[key] = *stage
	}
	return out
}

func (p *ProgressStore) SetHighestHintLevel(lessonID, stageID string, hintLevel int) {
	k := p.key(lessonID, stageID)
	if _, ok := p.stages[k]; !ok {
		p.stages[k] = &StageProgress{}
	}
	if hintLevel > p.stages[k].HighestHintLevel {
		p.stages[k].HighestHintLevel = hintLevel
	}
}

type Engine struct {
	catalog      map[string]Lesson
	progress     *ProgressStore
	persistence  ProgressPersistence
	persistError error
}

func NewEngine() *Engine {
	return newEngine(DefaultCatalog(), nil)
}

func NewEngineWithCatalog(catalog map[string]Lesson) *Engine {
	return newEngine(catalog, nil)
}

func NewEngineWithCatalogAndPersistence(catalog map[string]Lesson, persistence ProgressPersistence) *Engine {
	return newEngine(catalog, persistence)
}

func newEngine(catalog map[string]Lesson, persistence ProgressPersistence) *Engine {
	copyCatalog := make(map[string]Lesson, len(catalog))
	for id, lesson := range catalog {
		copyCatalog[id] = lesson
	}
	progress := NewProgressStore()
	var persistErr error
	if persistence != nil {
		loaded, err := persistence.Load()
		if err != nil {
			persistErr = fmt.Errorf("load progress: %w", err)
		} else {
			progress = NewProgressStoreFromSnapshot(loaded)
		}
	}
	return &Engine{catalog: copyCatalog, progress: progress, persistence: persistence, persistError: persistErr}
}

func (e *Engine) Lessons() []Lesson {
	out := make([]Lesson, 0, len(e.catalog))
	ids := make([]string, 0, len(e.catalog))
	for id := range e.catalog {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	for _, id := range ids {
		lesson := e.catalog[id]
		out = append(out, lesson)
	}
	return out
}

func (e *Engine) RunStage(lessonID string, stageIndex int) (StageResult, error) {
	if e.persistError != nil {
		return StageResult{}, e.persistError
	}
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
			e.progress.SetHighestHintLevel(lesson.ID, stage.ID, hintLevel)
			if err := e.persistProgress(); err != nil {
				return StageResult{}, err
			}
			return StageResult{Passed: false, FeedbackKey: "validator." + v.Name, Hint: hint, HintLevel: hintLevel, Output: output}, nil
		}
	}

	e.progress.Record(lesson.ID, stage.ID, true)
	if err := e.persistProgress(); err != nil {
		return StageResult{}, err
	}
	return StageResult{Passed: true, FeedbackKey: "stage." + stage.ID + ".passed", Output: output}, nil
}

func (e *Engine) persistProgress() error {
	if e.persistence == nil {
		return nil
	}
	if err := e.persistence.Save(e.progress.Snapshot()); err != nil {
		return fmt.Errorf("persist progress: %w", err)
	}
	return nil
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
	modules := map[string]*ModuleAnalytics{}
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

			if _, ok := modules[lesson.Module]; !ok {
				modules[lesson.Module] = &ModuleAnalytics{Module: lesson.Module}
			}
			m := modules[lesson.Module]
			m.TotalStages++
			if prog.Completed {
				m.CompletedStage++
			}
		}
	}

	modOut := make([]ModuleAnalytics, 0, len(modules))
	for _, m := range modules {
		if m.TotalStages > 0 {
			m.CompletionRate = float64(m.CompletedStage) / float64(m.TotalStages)
		}
		modOut = append(modOut, *m)
	}
	sort.Slice(modOut, func(i, j int) bool { return modOut[i].Module < modOut[j].Module })

	analytics := CompletionAnalytics{TotalStages: total, CompletedStages: completed, AttemptedStages: attempted, ModuleBreakdown: modOut}
	if total > 0 {
		analytics.CompletionRate = float64(completed) / float64(total)
		analytics.AttemptCoverage = float64(attempted) / float64(total)
	}

	checklist := []string{}
	if attempted == total {
		checklist = append(checklist, "pilot.attempt-coverage.ok")
	}
	if completed >= total/2 {
		checklist = append(checklist, "pilot.completion-flow.ok")
	}
	if len(modOut) >= 4 {
		checklist = append(checklist, "pilot.module-coverage.ok")
	}
	analytics.PilotChecklist = checklist
	analytics.PilotChecklistOK = len(checklist) == 3
	analytics.WeakConcepts = weakConcepts(e.catalog, e.progress)
	return analytics
}

func weakConcepts(catalog map[string]Lesson, progress *ProgressStore) []ConceptWeakness {
	agg := map[string]*ConceptWeakness{}
	for _, lesson := range catalog {
		for _, stage := range lesson.Stages {
			if len(stage.ConceptTags) == 0 {
				continue
			}
			prog := progress.Get(lesson.ID, stage.ID)
			if prog.Attempts == 0 {
				continue
			}
			failedAttempts := prog.Attempts
			if prog.Completed && failedAttempts > 0 {
				failedAttempts--
			}
			highHintUses := 0
			if prog.HighestHintLevel >= 2 {
				highHintUses = 1
			}
			if failedAttempts == 0 && highHintUses == 0 {
				continue
			}
			for _, tag := range stage.ConceptTags {
				entry, ok := agg[tag]
				if !ok {
					entry = &ConceptWeakness{Concept: tag}
					agg[tag] = entry
				}
				entry.FailedAttempts += failedAttempts
				entry.HighHintUses += highHintUses
				entry.AffectedStages++
			}
		}
	}
	out := make([]ConceptWeakness, 0, len(agg))
	for _, item := range agg {
		item.Score = float64(item.FailedAttempts) + float64(item.HighHintUses)*0.5
		out = append(out, *item)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Score != out[j].Score {
			return out[i].Score > out[j].Score
		}
		if out[i].FailedAttempts != out[j].FailedAttempts {
			return out[i].FailedAttempts > out[j].FailedAttempts
		}
		return cmp.Less(out[i].Concept, out[j].Concept)
	})
	return out
}
