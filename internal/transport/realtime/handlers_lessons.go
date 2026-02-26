package realtime

import (
	"net/http"
	"strconv"
	"strings"

	"os-simulator-plan/internal/lessons"
)

func (s *Server) handleLessons(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, r, http.StatusMethodNotAllowed, "method_not_allowed", "method not allowed")
		return
	}

	learnerID := normalizeLearnerID(r.URL.Query().Get("learner_id"))
	engine := s.lessonEngineForLearner(learnerID)

	lessonsList := engine.Lessons()

	out := make([]LessonSummary, 0, len(lessonsList))
	for _, lesson := range lessonsList {
		stages := make([]LessonStageSummary, 0, len(lesson.Stages))
		for idx, stage := range lesson.Stages {
			status := engine.StageStatus(lesson.ID, stage)
			stages = append(stages, LessonStageSummary{
				Index:           idx,
				ID:              stage.ID,
				Title:           stage.Title,
				Theory:          stage.Hints.Concept,
				Objective:       stage.Objective,
				PassConditions:  stagePassConditions(stage),
				Prerequisites:   stage.Prerequisites,
				AllowedCommands: stage.AllowedCmds,
				Limits: ChallengeLimitsDTO{
					MaxSteps:         stage.Limits.MaxSteps,
					MaxPolicyChanges: stage.Limits.MaxPolicyChanges,
				},
				Attempts:  status.Attempts,
				Completed: status.Completed,
				Unlocked:  status.Unlocked,
			})
		}
		out = append(out, LessonSummary{ID: lesson.ID, Title: lesson.Title, Module: lesson.Module, Stages: stages})
	}

	respondJSON(w, http.StatusOK, LessonsResponse{Lessons: out})
}

func normalizeLearnerID(value string) string {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return "anonymous"
	}
	return trimmed
}

func stagePassConditions(stage lessons.Stage) []string {
	out := make([]string, 0, len(stage.Validators))
	for _, validator := range stage.Validators {
		out = append(out, describeValidator(validator))
	}
	return out
}

func describeValidator(v lessons.ValidatorSpec) string {
	switch v.Type {
	case "trace_contains_all":
		if len(v.Values) == 0 {
			return "Required trace events must appear."
		}
		return "Trace must contain: " + strings.Join(v.Values, ", ") + "."
	case "metric_eq":
		return "Metric " + v.Key + " must equal " + trimFloat(v.Number) + "."
	case "metric_lte":
		return "Metric " + v.Key + " must be <= " + trimFloat(v.Number) + "."
	case "fault_eq":
		return "Fault count " + v.Key + " must equal " + trimFloat(v.Number) + "."
	case "fault_lte":
		return "Fault count " + v.Key + " must be <= " + trimFloat(v.Number) + "."
	case "fs_ok":
		return "Filesystem invariants must hold."
	default:
		return "Check " + v.Name + " must pass."
	}
}

func trimFloat(value float64) string {
	if value == float64(int64(value)) {
		return strconv.FormatInt(int64(value), 10)
	}
	return strconv.FormatFloat(value, 'f', -1, 64)
}

func convertAnalytics(in lessons.CompletionAnalytics) CompletionAnalyticsDTO {
	return CompletionAnalyticsDTO{
		TotalStages:     in.TotalStages,
		CompletedStages: in.CompletedStages,
		AttemptedStages: in.AttemptedStages,
		CompletionRate:  in.CompletionRate,
	}
}
