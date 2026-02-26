package realtime

import (
	"encoding/json"
	"net/http"

	"os-simulator-plan/internal/lessons"
	"os-simulator-plan/internal/sim"
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

func (s *Server) handleLessonRun(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, r, http.StatusMethodNotAllowed, "method_not_allowed", "method not allowed")
		return
	}

	defer func() { _ = r.Body.Close() }()
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	var req LessonRunRequest
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&req); err != nil {
		respondError(w, r, http.StatusBadRequest, "invalid_body", "invalid JSON body")
		return
	}
	if req.LessonID == "" {
		respondError(w, r, http.StatusBadRequest, "invalid_lesson_request", "lesson_id is required")
		return
	}

	s.lessonMu.Lock()
	result, err := s.lessonEngine.RunStage(req.LessonID, req.StageIndex)
	if err != nil {
		s.lessonMu.Unlock()
		respondError(w, r, http.StatusBadRequest, "lesson_run_failed", err.Error())
		return
	}
	analytics := s.lessonEngine.CompletionAnalytics()
	s.lessonMu.Unlock()

	respondJSON(w, http.StatusOK, LessonRunResponse{
		LessonID:    req.LessonID,
		StageIndex:  req.StageIndex,
		Passed:      result.Passed,
		FeedbackKey: result.FeedbackKey,
		Hint:        result.Hint,
		HintLevel:   result.HintLevel,
		Output: LessonOutputDTO{
			Tick:         result.Output.Metrics.TotalTicks,
			TraceHash:    sim.TraceHash(result.Output.Trace),
			TraceLength:  len(result.Output.Trace),
			Processes:    result.Output.Processes,
			Metrics:      result.Output.Metrics,
			Memory:       result.Output.Memory,
			FilesystemOK: result.Output.FilesystemOK,
		},
		Analytics: convertAnalytics(analytics),
	})
}

func convertAnalytics(in lessons.CompletionAnalytics) CompletionAnalyticsDTO {
	return CompletionAnalyticsDTO{
		TotalStages:     in.TotalStages,
		CompletedStages: in.CompletedStages,
		AttemptedStages: in.AttemptedStages,
		CompletionRate:  in.CompletionRate,
	}
}
