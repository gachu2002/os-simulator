package realtime

import (
	"encoding/json"
	"fmt"
	"net/http"
	"slices"

	"os-simulator-plan/internal/lessons"
	"os-simulator-plan/internal/sim"
)

func (s *Server) handleChallengeStart(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, r, http.StatusMethodNotAllowed, "method_not_allowed", "method not allowed")
		return
	}

	defer func() { _ = r.Body.Close() }()
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	var req ChallengeStartRequest
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&req); err != nil {
		respondError(w, r, http.StatusBadRequest, "invalid_body", "invalid JSON body")
		return
	}
	if req.LessonID == "" {
		respondError(w, r, http.StatusBadRequest, "invalid_challenge_request", "lesson_id is required")
		return
	}

	s.lessonMu.Lock()
	prepared, err := s.lessonEngine.PrepareStage(req.LessonID, req.StageIndex)
	s.lessonMu.Unlock()
	if err != nil {
		respondError(w, r, http.StatusBadRequest, "challenge_start_failed", err.Error())
		return
	}

	sessionCfg := SessionConfig{
		Seed:            prepared.Stage.Config.Seed,
		CheckpointEvery: 5,
		Policy:          prepared.Stage.Config.Policy,
		Quantum:         prepared.Stage.Config.Quantum,
		Frames:          prepared.Stage.Config.Frames,
		TLBEntries:      prepared.Stage.Config.TLBEntries,
		DiskLatency:     prepared.Stage.Config.DiskLatency,
		TerminalLatency: prepared.Stage.Config.TerminalLatency,
	}
	session, err := s.manager.Create(sessionCfg)
	if err != nil {
		respondError(w, r, http.StatusBadRequest, "invalid_session_config", err.Error())
		return
	}

	allowedCommands := challengeAllowedCommands(prepared)
	maxSteps, maxPolicyChanges := challengeLimits(prepared)

	if err := applyBootstrapCommands(session, prepared.Stage.Bootstrap); err != nil {
		respondError(w, r, http.StatusBadRequest, "challenge_start_failed", err.Error())
		return
	}

	session.SetChallengePolicy(NewChallengeCommandPolicy(
		allowedCommands,
		maxSteps,
		maxPolicyChanges,
	))

	attempt := s.challengeAttempts.Create(session.ID(), prepared)
	respondJSON(w, http.StatusOK, ChallengeStartResponse{
		AttemptID:       attempt.AttemptID,
		SessionID:       attempt.SessionID,
		LessonID:        prepared.LessonID,
		StageIndex:      prepared.StageIndex,
		StageTitle:      prepared.Stage.Title,
		Module:          prepared.Module,
		Objective:       challengeObjective(prepared),
		AllowedCommands: allowedCommands,
		Limits: ChallengeLimitsDTO{
			MaxSteps:         maxSteps,
			MaxPolicyChanges: maxPolicyChanges,
		},
	})
}

func (s *Server) handleChallengeGrade(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, r, http.StatusMethodNotAllowed, "method_not_allowed", "method not allowed")
		return
	}

	defer func() { _ = r.Body.Close() }()
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	var req ChallengeGradeRequest
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&req); err != nil {
		respondError(w, r, http.StatusBadRequest, "invalid_body", "invalid JSON body")
		return
	}
	if req.AttemptID == "" {
		respondError(w, r, http.StatusBadRequest, "invalid_challenge_request", "attempt_id is required")
		return
	}

	attempt, ok := s.challengeAttempts.Get(req.AttemptID)
	if !ok {
		respondError(w, r, http.StatusNotFound, "challenge_attempt_not_found", "challenge attempt not found")
		return
	}

	session, ok := s.manager.Get(attempt.SessionID)
	if !ok {
		respondError(w, r, http.StatusNotFound, "session_not_found", "session not found")
		return
	}

	output := session.StageOutput()
	s.lessonMu.Lock()
	result := s.lessonEngine.GradeStage(attempt.Prepared, output)
	analytics := s.lessonEngine.CompletionAnalytics()
	s.lessonMu.Unlock()

	respondJSON(w, http.StatusOK, ChallengeGradeResponse{
		AttemptID:   attempt.AttemptID,
		LessonID:    attempt.Prepared.LessonID,
		StageIndex:  attempt.Prepared.StageIndex,
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

func challengeObjective(prepared lessons.PreparedStage) string {
	if prepared.Stage.Objective != "" {
		return prepared.Stage.Objective
	}
	return prepared.Stage.Title
}

func challengeAllowedCommands(prepared lessons.PreparedStage) []string {
	if len(prepared.Stage.AllowedCmds) == 0 {
		return slices.Clone(defaultChallengeAllowedCommands)
	}
	return slices.Clone(prepared.Stage.AllowedCmds)
}

func challengeLimits(prepared lessons.PreparedStage) (int, int) {
	maxSteps := prepared.Stage.Limits.MaxSteps
	if maxSteps <= 0 {
		maxSteps = defaultChallengeMaxSteps
	}
	maxPolicyChanges := prepared.Stage.Limits.MaxPolicyChanges
	if maxPolicyChanges <= 0 {
		maxPolicyChanges = defaultChallengeMaxPolicyChanges
	}
	return maxSteps, maxPolicyChanges
}

func applyBootstrapCommands(session *Session, commands []sim.Command) error {
	for _, cmd := range commands {
		ev := session.Apply(Command{
			Name:    cmd.Name,
			Count:   cmd.Count,
			Process: cmd.Process,
			Program: cmd.Program,
			Policy:  cmd.Policy,
			Quantum: cmd.Quantum,
		})
		if ev.Type == "session.error" {
			return fmt.Errorf("bootstrap command failed: %s", ev.Error)
		}
	}
	return nil
}
