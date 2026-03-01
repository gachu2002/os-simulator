package realtime

import (
	"net/http"
	"strings"

	"os-simulator-plan/internal/sim"
)

type ChallengeActionV3Request struct {
	AttemptID       string   `json:"attempt_id"`
	LearnerID       string   `json:"learner_id,omitempty"`
	Action          string   `json:"action"`
	Count           int      `json:"count,omitempty"`
	Process         string   `json:"process,omitempty"`
	Program         string   `json:"program,omitempty"`
	Policy          string   `json:"policy,omitempty"`
	Quantum         int      `json:"quantum,omitempty"`
	Frames          int      `json:"frames,omitempty"`
	TLBEntries      int      `json:"tlb_entries,omitempty"`
	DiskLatency     sim.Tick `json:"disk_latency,omitempty"`
	TerminalLatency sim.Tick `json:"terminal_latency,omitempty"`
}

type ChallengeActionV3Response struct {
	AttemptID     string `json:"attempt_id"`
	SessionID     string `json:"session_id"`
	Action        string `json:"action"`
	MappedCommand string `json:"mapped_command"`
	Event         Event  `json:"event"`
}

func (s *Server) handleChallengeActionV3(w http.ResponseWriter, r *http.Request) {
	if !requireMethod(w, r, http.MethodPost) {
		return
	}
	var req ChallengeActionV3Request
	if !decodeJSONBody(w, r, &req) {
		return
	}
	if req.AttemptID == "" {
		respondError(w, r, http.StatusBadRequest, "invalid_challenge_request", "attempt_id is required")
		return
	}
	if strings.TrimSpace(req.Action) == "" {
		respondError(w, r, http.StatusBadRequest, "invalid_challenge_request", "action is required")
		return
	}

	attempt, ok := s.challengeAttempts.Get(req.AttemptID)
	if !ok {
		respondError(w, r, http.StatusNotFound, "challenge_attempt_not_found", "challenge attempt not found")
		return
	}
	learnerID := normalizeLearnerID(req.LearnerID)
	if attempt.LearnerID != learnerID {
		respondError(w, r, http.StatusForbidden, "challenge_attempt_forbidden", "challenge attempt belongs to another learner")
		return
	}
	session, ok := s.manager.Get(attempt.SessionID)
	if !ok {
		respondError(w, r, http.StatusNotFound, "session_not_found", "session not found")
		return
	}

	cmd, err := mapV3ActionToCommand(req)
	if err != nil {
		respondError(w, r, http.StatusBadRequest, "invalid_action", err.Error())
		return
	}

	respondJSON(w, http.StatusOK, ChallengeActionV3Response{
		AttemptID:     attempt.AttemptID,
		SessionID:     attempt.SessionID,
		Action:        req.Action,
		MappedCommand: cmd.Name,
		Event:         session.Apply(cmd),
	})
}

func unsupportedActionMessage(action string, note ActionCapabilityNote) string {
	msg := "action " + strings.TrimSpace(action) + " is planned"
	if note.Reason != "" {
		msg += ": " + note.Reason
	}
	if note.FallbackAction != "" {
		msg += ". fallback_action=" + note.FallbackAction
	}
	return msg
}
