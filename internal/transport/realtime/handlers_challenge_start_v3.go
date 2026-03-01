package realtime

import (
	"net/http"
	"strings"

	contentv3 "os-simulator-plan/internal/content/v3"
)

type ChallengeStartV3Request struct {
	LessonID  string `json:"lesson_id"`
	PartID    string `json:"part_id,omitempty"`
	LearnerID string `json:"learner_id,omitempty"`
}

type ChallengeStartV3Response struct {
	Version               string                          `json:"version"`
	SectionID             string                          `json:"section_id"`
	LessonID              string                          `json:"lesson_id"`
	LessonTitle           string                          `json:"lesson_title"`
	LessonObjective       string                          `json:"lesson_objective"`
	PartID                string                          `json:"part_id,omitempty"`
	PartTitle             string                          `json:"part_title,omitempty"`
	PartObjective         string                          `json:"part_objective,omitempty"`
	AttemptID             string                          `json:"attempt_id"`
	SessionID             string                          `json:"session_id"`
	AllowedCommands       []string                        `json:"allowed_commands"`
	Limits                ChallengeLimitsDTO              `json:"limits"`
	ChallengeDescriptor   string                          `json:"challenge_description"`
	Visualizer            []string                        `json:"visualizer"`
	Actions               []string                        `json:"actions"`
	ActionCapabilities    ActionCapabilities              `json:"action_capabilities"`
	ActionCapabilityNotes map[string]ActionCapabilityNote `json:"action_capability_notes,omitempty"`
}

func (s *Server) handleChallengeStartV3(w http.ResponseWriter, r *http.Request) {
	if !requireMethod(w, r, http.MethodPost) {
		return
	}
	if s.cpuCurriculumErr != nil {
		respondError(w, r, http.StatusInternalServerError, "curriculum_unavailable", "curriculum unavailable")
		return
	}

	var req ChallengeStartV3Request
	if !decodeJSONBody(w, r, &req) {
		return
	}
	if req.LessonID == "" {
		respondError(w, r, http.StatusBadRequest, "invalid_challenge_request", "lesson_id is required")
		return
	}
	if !isActiveCPULesson(req.LessonID) {
		respondError(w, r, http.StatusBadRequest, "invalid_challenge_request", "lesson_id is not part of active section")
		return
	}

	section, lesson, ok := s.lessonV3ByID(req.LessonID)
	if !ok {
		respondError(w, r, http.StatusNotFound, "lesson_not_found", "lesson not found")
		return
	}
	engineLessonID, ok := resolveEngineLessonID(req.LessonID)
	if !ok {
		respondError(w, r, http.StatusBadRequest, "invalid_challenge_request", "lesson mapping not found")
		return
	}

	stageIndex, part, err := resolveLessonPartToStage(lesson, req.PartID)
	if err != nil {
		respondError(w, r, http.StatusBadRequest, "invalid_challenge_part", err.Error())
		return
	}

	learnerID := normalizeLearnerID(req.LearnerID)
	start, err := s.challengeService.Start(engineLessonID, stageIndex, learnerID)
	if err != nil {
		respondChallengeServiceError(w, r, err)
		return
	}
	session, ok := s.manager.Get(start.Attempt.SessionID)
	if !ok {
		respondError(w, r, http.StatusNotFound, "session_not_found", "session not found")
		return
	}
	start.AllowedCommands = applyV3ActionPolicy(session, start.AllowedCommands, lesson.Challenge.Actions, start.MaxSteps, start.MaxPolicy, start.MaxConfig)

	out := ChallengeStartV3Response{
		Version:               s.cpuCurriculumV3.Version,
		SectionID:             section.ID,
		LessonID:              lesson.ID,
		LessonTitle:           lesson.Title,
		LessonObjective:       lesson.Objective,
		AttemptID:             start.Attempt.AttemptID,
		SessionID:             start.Attempt.SessionID,
		AllowedCommands:       start.AllowedCommands,
		ChallengeDescriptor:   lesson.Challenge.Description,
		Visualizer:            append([]string(nil), lesson.Challenge.Visualizer...),
		Actions:               append([]string(nil), lesson.Challenge.Actions...),
		ActionCapabilities:    classifyActionCapabilities(lesson.Challenge.Actions),
		ActionCapabilityNotes: buildActionCapabilityNotes(lesson.Challenge.Actions),
		Limits: ChallengeLimitsDTO{
			MaxSteps:         start.MaxSteps,
			MaxPolicyChanges: start.MaxPolicy,
			MaxConfigChanges: start.MaxConfig,
		},
	}
	if part != nil {
		out.PartID = part.ID
		out.PartTitle = part.Title
		out.PartObjective = part.Objective
	}

	respondJSON(w, http.StatusOK, out)
}

func resolveLessonPartToStage(lesson contentv3.Lesson, partID string) (int, *contentv3.ChallengePart, error) {
	if len(lesson.Challenge.Parts) == 0 {
		if strings.TrimSpace(partID) != "" {
			return 0, nil, &appErr{msg: "part_id is not supported for this lesson"}
		}
		return 0, nil, nil
	}

	trimmed := strings.TrimSpace(partID)
	if trimmed == "" {
		return 0, nil, &appErr{msg: "part_id is required for this lesson"}
	}
	for idx := range lesson.Challenge.Parts {
		if lesson.Challenge.Parts[idx].ID == trimmed {
			return idx, &lesson.Challenge.Parts[idx], nil
		}
	}
	return 0, nil, &appErr{msg: "part_id is invalid for this lesson"}
}

type appErr struct{ msg string }

func (e *appErr) Error() string { return e.msg }
