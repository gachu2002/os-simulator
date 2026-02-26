package realtime

import (
	"net/http"

	"os-simulator-plan/internal/lessons"
)

func (s *Server) handleLessons(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, r, http.StatusMethodNotAllowed, "method_not_allowed", "method not allowed")
		return
	}

	s.lessonMu.Lock()
	lessonsList := s.lessonEngine.Lessons()
	s.lessonMu.Unlock()

	out := make([]LessonSummary, 0, len(lessonsList))
	for _, lesson := range lessonsList {
		stages := make([]LessonStageSummary, 0, len(lesson.Stages))
		for idx, stage := range lesson.Stages {
			stages = append(stages, LessonStageSummary{
				Index: idx,
				ID:    stage.ID,
				Title: stage.Title,
			})
		}
		out = append(out, LessonSummary{ID: lesson.ID, Title: lesson.Title, Module: lesson.Module, Stages: stages})
	}

	respondJSON(w, http.StatusOK, LessonsResponse{Lessons: out})
}

func convertAnalytics(in lessons.CompletionAnalytics) CompletionAnalyticsDTO {
	return CompletionAnalyticsDTO{
		TotalStages:     in.TotalStages,
		CompletedStages: in.CompletedStages,
		AttemptedStages: in.AttemptedStages,
		CompletionRate:  in.CompletionRate,
	}
}
