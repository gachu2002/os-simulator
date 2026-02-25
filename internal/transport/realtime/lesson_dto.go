package realtime

import "os-simulator-plan/internal/sim"

type LessonStageSummary struct {
	Index int    `json:"index"`
	ID    string `json:"id"`
	Title string `json:"title"`
}

type LessonSummary struct {
	ID     string               `json:"id"`
	Title  string               `json:"title"`
	Module string               `json:"module"`
	Stages []LessonStageSummary `json:"stages"`
}

type LessonsResponse struct {
	Lessons []LessonSummary `json:"lessons"`
}

type LessonRunRequest struct {
	LessonID   string `json:"lesson_id"`
	StageIndex int    `json:"stage_index"`
}

type LessonOutputDTO struct {
	Tick         sim.Tick              `json:"tick"`
	TraceHash    string                `json:"trace_hash"`
	TraceLength  int                   `json:"trace_length"`
	Processes    []sim.ProcessSnapshot `json:"processes"`
	Metrics      sim.SchedulingMetrics `json:"metrics"`
	Memory       sim.MemorySnapshot    `json:"memory"`
	FilesystemOK bool                  `json:"filesystem_ok"`
}

type ModuleAnalyticsDTO struct {
	Module         string  `json:"module"`
	TotalStages    int     `json:"total_stages"`
	CompletedStage int     `json:"completed_stage"`
	CompletionRate float64 `json:"completion_rate"`
}

type CompletionAnalyticsDTO struct {
	TotalStages      int                  `json:"total_stages"`
	CompletedStages  int                  `json:"completed_stages"`
	AttemptedStages  int                  `json:"attempted_stages"`
	CompletionRate   float64              `json:"completion_rate"`
	AttemptCoverage  float64              `json:"attempt_coverage"`
	ModuleBreakdown  []ModuleAnalyticsDTO `json:"module_breakdown"`
	PilotChecklist   []string             `json:"pilot_checklist"`
	PilotChecklistOK bool                 `json:"pilot_checklist_ok"`
}

type LessonRunResponse struct {
	LessonID    string                 `json:"lesson_id"`
	StageIndex  int                    `json:"stage_index"`
	Passed      bool                   `json:"passed"`
	FeedbackKey string                 `json:"feedback_key"`
	Hint        string                 `json:"hint,omitempty"`
	HintLevel   int                    `json:"hint_level,omitempty"`
	Output      LessonOutputDTO        `json:"output"`
	Analytics   CompletionAnalyticsDTO `json:"analytics"`
}
