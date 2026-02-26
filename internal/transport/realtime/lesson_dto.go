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

type LessonOutputDTO struct {
	Tick         sim.Tick              `json:"tick"`
	TraceHash    string                `json:"trace_hash"`
	TraceLength  int                   `json:"trace_length"`
	Processes    []sim.ProcessSnapshot `json:"processes"`
	Metrics      sim.SchedulingMetrics `json:"metrics"`
	Memory       sim.MemorySnapshot    `json:"memory"`
	FilesystemOK bool                  `json:"filesystem_ok"`
}

type CompletionAnalyticsDTO struct {
	TotalStages     int     `json:"total_stages"`
	CompletedStages int     `json:"completed_stages"`
	AttemptedStages int     `json:"attempted_stages"`
	CompletionRate  float64 `json:"completion_rate"`
}
