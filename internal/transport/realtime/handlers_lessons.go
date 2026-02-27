package realtime

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"os-simulator-plan/internal/lessons"
)

func (s *Server) handleCurriculum(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, r, http.StatusMethodNotAllowed, "method_not_allowed", "method not allowed")
		return
	}

	learnerID := normalizeLearnerID(r.URL.Query().Get("learner_id"))
	engine := s.lessonEngineForLearner(learnerID)
	lessonsList := engine.Lessons()
	lessonSummaries := lessonSummariesForEngine(engine, lessonsList)

	sections := make([]CurriculumSection, 0, len(curriculumSectionOrder))
	for idx, spec := range curriculumSectionOrder {
		sectionLessons := make([]LessonSummary, 0)
		completedStages := 0
		totalStages := 0
		for _, lesson := range lessonSummaries {
			if lesson.SectionID != spec.ID {
				continue
			}
			sectionLessons = append(sectionLessons, lesson)
			totalStages += len(lesson.Stages)
			for _, stage := range lesson.Stages {
				if stage.Completed {
					completedStages++
				}
			}
		}
		sections = append(sections, CurriculumSection{
			ID:              spec.ID,
			Title:           spec.Title,
			Subtitle:        spec.Subtitle,
			Order:           idx + 1,
			ComingSoon:      spec.ComingSoon,
			Lessons:         sectionLessons,
			CompletedStages: completedStages,
			TotalStages:     totalStages,
		})
	}

	respondJSON(w, http.StatusOK, CurriculumResponse{Sections: sections})
}

func (s *Server) handleLessonLearn(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, r, http.StatusMethodNotAllowed, "method_not_allowed", "method not allowed")
		return
	}

	lessonID := strings.TrimSpace(chi.URLParam(r, "lessonID"))
	if lessonID == "" {
		respondError(w, r, http.StatusBadRequest, "invalid_lesson_id", "lesson id is required")
		return
	}

	learnerID := normalizeLearnerID(r.URL.Query().Get("learner_id"))
	engine := s.lessonEngineForLearner(learnerID)
	lessonsList := engine.Lessons()

	for _, lesson := range lessonsList {
		if lesson.ID != lessonID {
			continue
		}
		stages := make([]LessonLearnStage, 0, len(lesson.Stages))
		for idx, stage := range lesson.Stages {
			stages = append(stages, LessonLearnStage{
				Index:                 idx,
				ID:                    stage.ID,
				Title:                 stage.Title,
				CoreIdea:              stage.CoreIdea,
				MechanismSteps:        stage.MechanismSteps,
				WorkedExample:         stage.WorkedExample,
				CommonMistakes:        stage.CommonMistakes,
				PreChallengeChecklist: stage.PreChallengeChecklist,
				Objective:             stage.Objective,
				Goal:                  stage.Goal,
				Prerequisites:         stage.Prerequisites,
				ExpectedVisualCues:    stage.ExpectedVisualCues,
			})
		}
		respondJSON(w, http.StatusOK, LessonLearnResponse{Lesson: LessonLearnSummary{
			ID:               lesson.ID,
			Title:            lesson.Title,
			Module:           lesson.Module,
			SectionID:        lesson.SectionID,
			SectionTitle:     lesson.SectionTitle,
			Difficulty:       lesson.Difficulty,
			EstimatedMinutes: lesson.EstimatedMinutes,
			ChapterRefs:      lesson.ChapterRefs,
			Stages:           stages,
		}})
		return
	}

	respondError(w, r, http.StatusNotFound, "lesson_not_found", "lesson not found")
}

func lessonSummariesForEngine(engine *lessons.Engine, lessonsList []lessons.Lesson) []LessonSummary {

	out := make([]LessonSummary, 0, len(lessonsList))
	for _, lesson := range lessonsList {
		stages := make([]LessonStageSummary, 0, len(lesson.Stages))
		for idx, stage := range lesson.Stages {
			status := engine.StageStatus(lesson.ID, stage)
			stages = append(stages, LessonStageSummary{
				Index:              idx,
				ID:                 stage.ID,
				Title:              stage.Title,
				Objective:          stage.Objective,
				Goal:               stage.Goal,
				PassConditions:     stagePassConditions(stage),
				Prerequisites:      stage.Prerequisites,
				AllowedCommands:    stage.AllowedCmds,
				ActionDescriptions: convertActionDescriptions(stage.ActionDescriptions),
				ExpectedVisualCues: stage.ExpectedVisualCues,
				Limits: ChallengeLimitsDTO{
					MaxSteps:         stage.Limits.MaxSteps,
					MaxPolicyChanges: stage.Limits.MaxPolicyChanges,
					MaxConfigChanges: stage.Limits.MaxConfigChanges,
				},
				Attempts:  status.Attempts,
				Completed: status.Completed,
				Unlocked:  status.Unlocked,
			})
		}
		out = append(out, LessonSummary{
			ID:               lesson.ID,
			Title:            lesson.Title,
			Module:           lesson.Module,
			SectionID:        lesson.SectionID,
			SectionTitle:     lesson.SectionTitle,
			Difficulty:       lesson.Difficulty,
			EstimatedMinutes: lesson.EstimatedMinutes,
			ChapterRefs:      lesson.ChapterRefs,
			Stages:           stages,
		})
	}
	return out
}

type curriculumSectionSpec struct {
	ID         string
	Title      string
	Subtitle   string
	ComingSoon bool
}

var curriculumSectionOrder = []curriculumSectionSpec{
	{ID: "introduction", Title: "Introduction", Subtitle: "OSTEP setup and foundational framing", ComingSoon: true},
	{ID: "virtualization", Title: "Virtualization", Subtitle: "CPU and memory virtualization lessons"},
	{ID: "concurrency", Title: "Concurrency", Subtitle: "Threads, wakeups, and interrupt-driven progress"},
	{ID: "persistence", Title: "Persistence", Subtitle: "Storage and filesystem correctness"},
	{ID: "security", Title: "Security", Subtitle: "Authentication, access control, and protection", ComingSoon: true},
}

func convertActionDescriptions(items []lessons.ActionDescription) []LessonActionDescription {
	out := make([]LessonActionDescription, 0, len(items))
	for _, item := range items {
		out = append(out, LessonActionDescription{Command: item.Command, Description: item.Description})
	}
	return out
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
	case "trace_order":
		if len(v.Values) == 0 {
			return "Trace events must follow required order."
		}
		return "Trace order must include: " + strings.Join(v.Values, " -> ") + "."
	case "trace_count_eq":
		if len(v.Values) == 0 {
			return "Trace event count must equal required value."
		}
		return "Trace count for " + v.Values[0] + " must equal " + trimFloat(v.Number) + "."
	case "trace_count_lte":
		if len(v.Values) == 0 {
			return "Trace event count must be <= required value."
		}
		return "Trace count for " + v.Values[0] + " must be <= " + trimFloat(v.Number) + "."
	case "no_event":
		if len(v.Values) == 0 {
			return "Forbidden trace events must not appear."
		}
		return "Trace must not contain: " + strings.Join(v.Values, ", ") + "."
	case "metric_eq":
		return "Metric " + v.Key + " must equal " + trimFloat(v.Number) + "."
	case "metric_gte":
		return "Metric " + v.Key + " must be >= " + trimFloat(v.Number) + "."
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
