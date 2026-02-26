package realtime

type ChallengeStartRequest struct {
	LessonID   string `json:"lesson_id"`
	StageIndex int    `json:"stage_index"`
}

type ChallengeLimitsDTO struct {
	MaxSteps         int `json:"max_steps,omitempty"`
	MaxPolicyChanges int `json:"max_policy_changes,omitempty"`
}

type ChallengeStartResponse struct {
	AttemptID       string             `json:"attempt_id"`
	SessionID       string             `json:"session_id"`
	LessonID        string             `json:"lesson_id"`
	StageIndex      int                `json:"stage_index"`
	StageTitle      string             `json:"stage_title"`
	Module          string             `json:"module"`
	Objective       string             `json:"objective"`
	AllowedCommands []string           `json:"allowed_commands"`
	Limits          ChallengeLimitsDTO `json:"limits"`
}

type ChallengeGradeRequest struct {
	AttemptID string `json:"attempt_id"`
}

type ChallengeGradeResponse struct {
	AttemptID   string                 `json:"attempt_id"`
	LessonID    string                 `json:"lesson_id"`
	StageIndex  int                    `json:"stage_index"`
	Passed      bool                   `json:"passed"`
	FeedbackKey string                 `json:"feedback_key"`
	Hint        string                 `json:"hint,omitempty"`
	HintLevel   int                    `json:"hint_level,omitempty"`
	Output      LessonOutputDTO        `json:"output"`
	Analytics   CompletionAnalyticsDTO `json:"analytics"`
}
