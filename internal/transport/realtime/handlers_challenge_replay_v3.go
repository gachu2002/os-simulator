package realtime

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"os-simulator-plan/internal/sim"
)

type ChallengeReplayV3Response struct {
	Version            string                `json:"version"`
	SectionID          string                `json:"section_id"`
	LessonID           string                `json:"lesson_id"`
	LessonTitle        string                `json:"lesson_title"`
	LessonObjective    string                `json:"lesson_objective"`
	PartID             string                `json:"part_id,omitempty"`
	PartTitle          string                `json:"part_title,omitempty"`
	PartObjective      string                `json:"part_objective,omitempty"`
	AttemptID          string                `json:"attempt_id"`
	SessionID          string                `json:"session_id"`
	Trace              []sim.TraceEvent      `json:"trace"`
	TraceHash          string                `json:"trace_hash"`
	TraceLength        int                   `json:"trace_length"`
	Processes          []sim.ProcessSnapshot `json:"processes"`
	Metrics            sim.SchedulingMetrics `json:"metrics"`
	Memory             sim.MemorySnapshot    `json:"memory"`
	FilesystemOK       bool                  `json:"filesystem_ok"`
	ActionCapabilities ActionCapabilities    `json:"action_capabilities"`
}

func (s *Server) handleChallengeReplayV3(w http.ResponseWriter, r *http.Request) {
	if !requireMethod(w, r, http.MethodGet) {
		return
	}
	if s.cpuCurriculumErr != nil {
		respondError(w, r, http.StatusInternalServerError, "curriculum_unavailable", "curriculum unavailable")
		return
	}

	attemptID := strings.TrimSpace(chi.URLParam(r, "attemptID"))
	if attemptID == "" {
		respondError(w, r, http.StatusBadRequest, "invalid_challenge_request", "attempt_id is required")
		return
	}
	learnerID := normalizeLearnerID(r.URL.Query().Get("learner_id"))

	attempt, ok := s.challengeAttempts.Get(attemptID)
	if !ok {
		respondError(w, r, http.StatusNotFound, "challenge_attempt_not_found", "challenge attempt not found")
		return
	}
	if attempt.LearnerID != learnerID {
		respondError(w, r, http.StatusForbidden, "challenge_attempt_forbidden", "challenge attempt belongs to another learner")
		return
	}

	session, ok := s.manager.Get(attempt.SessionID)
	if !ok {
		respondError(w, r, http.StatusNotFound, "session_not_found", "session not found")
		return
	}

	v3LessonID, ok := resolveV3LessonID(attempt.Prepared.LessonID)
	if !ok {
		respondError(w, r, http.StatusNotFound, "lesson_not_found", "lesson not found")
		return
	}
	section, lesson, ok := s.lessonV3ByID(v3LessonID)
	if !ok {
		respondError(w, r, http.StatusNotFound, "lesson_not_found", "lesson not found")
		return
	}

	output := session.StageOutput()
	out := ChallengeReplayV3Response{
		Version:            s.cpuCurriculumV3.Version,
		SectionID:          section.ID,
		LessonID:           lesson.ID,
		LessonTitle:        lesson.Title,
		LessonObjective:    lesson.Objective,
		AttemptID:          attempt.AttemptID,
		SessionID:          attempt.SessionID,
		Trace:              output.Trace,
		TraceHash:          sim.TraceHash(output.Trace),
		TraceLength:        len(output.Trace),
		Processes:          output.Processes,
		Metrics:            output.Metrics,
		Memory:             output.Memory,
		FilesystemOK:       output.FilesystemOK,
		ActionCapabilities: classifyActionCapabilities(lesson.Challenge.Actions),
	}
	if part := lessonPartByStageIndex(lesson, attempt.Prepared.StageIndex); part != nil {
		out.PartID = part.ID
		out.PartTitle = part.Title
		out.PartObjective = part.Objective
	}

	respondJSON(w, http.StatusOK, out)
}
