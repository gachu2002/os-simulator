package realtime

import "os-simulator-plan/internal/sim"

type LessonStageSummary struct {
	Index              int                       `json:"index"`
	ID                 string                    `json:"id"`
	Title              string                    `json:"title"`
	Theory             string                    `json:"theory,omitempty"`
	TheoryDetail       string                    `json:"theory_detail,omitempty"`
	Objective          string                    `json:"objective"`
	Goal               string                    `json:"goal,omitempty"`
	PassConditions     []string                  `json:"pass_conditions,omitempty"`
	Prerequisites      []string                  `json:"prerequisites,omitempty"`
	AllowedCommands    []string                  `json:"allowed_commands,omitempty"`
	ActionDescriptions []LessonActionDescription `json:"action_descriptions,omitempty"`
	ExpectedVisualCues []string                  `json:"expected_visual_cues,omitempty"`
	Limits             ChallengeLimitsDTO        `json:"limits"`
	Attempts           int                       `json:"attempts"`
	Completed          bool                      `json:"completed"`
	Unlocked           bool                      `json:"unlocked"`
}

type LessonActionDescription struct {
	Command     string `json:"command"`
	Description string `json:"description"`
}

type LessonSummary struct {
	ID               string               `json:"id"`
	Title            string               `json:"title"`
	Module           string               `json:"module"`
	SectionID        string               `json:"section_id,omitempty"`
	SectionTitle     string               `json:"section_title,omitempty"`
	Difficulty       string               `json:"difficulty,omitempty"`
	EstimatedMinutes int                  `json:"estimated_minutes,omitempty"`
	ChapterRefs      []string             `json:"chapter_refs,omitempty"`
	Stages           []LessonStageSummary `json:"stages"`
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
