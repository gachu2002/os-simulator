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
	learnerID := normalizeLearnerID(req.LearnerID)
	engine := s.lessonEngineForLearner(learnerID)

	prepared, err := engine.PrepareStage(req.LessonID, req.StageIndex)
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
	maxSteps, maxPolicyChanges, maxConfigChanges := challengeLimits(prepared)

	if err := applyBootstrapCommands(session, prepared.Stage.Bootstrap); err != nil {
		respondError(w, r, http.StatusBadRequest, "challenge_start_failed", err.Error())
		return
	}

	session.SetChallengePolicy(NewChallengeCommandPolicy(
		allowedCommands,
		maxSteps,
		maxPolicyChanges,
		maxConfigChanges,
	))

	attempt := s.challengeAttempts.Create(session.ID(), learnerID, prepared)
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
			MaxConfigChanges: maxConfigChanges,
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
	learnerID := normalizeLearnerID(req.LearnerID)

	attempt, ok := s.challengeAttempts.Get(req.AttemptID)
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

	output := session.StageOutput()
	engine := s.lessonEngineForLearner(learnerID)
	result := engine.GradeStage(attempt.Prepared, output)
	analytics := engine.CompletionAnalytics()

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
		Analytics:        convertAnalytics(analytics),
		ValidatorResults: convertValidatorResults(result.ValidatorResults),
	})
}

func convertValidatorResults(in []lessons.ValidationResult) []ValidatorResultDTO {
	out := make([]ValidatorResultDTO, 0, len(in))
	for _, item := range in {
		out = append(out, ValidatorResultDTO{
			Name:    item.Name,
			Type:    item.Type,
			Key:     item.Key,
			Passed:  item.Passed,
			Message: item.Message,
		})
	}
	return out
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

func challengeLimits(prepared lessons.PreparedStage) (int, int, int) {
	maxSteps := prepared.Stage.Limits.MaxSteps
	if maxSteps <= 0 {
		maxSteps = defaultChallengeMaxSteps
	}

	allowed := challengeAllowedCommands(prepared)

	maxPolicyChanges := prepared.Stage.Limits.MaxPolicyChanges
	if maxPolicyChanges <= 0 && hasAllowedCommand(allowed, "policy") {
		maxPolicyChanges = defaultChallengeMaxPolicyChanges
	}

	maxConfigChanges := prepared.Stage.Limits.MaxConfigChanges
	if maxConfigChanges <= 0 && hasAnyAllowedCommands(allowed, "set_frames", "set_tlb_entries", "set_disk_latency", "set_terminal_latency") {
		maxConfigChanges = defaultChallengeMaxConfigChanges
	}
	return maxSteps, maxPolicyChanges, maxConfigChanges
}

func hasAllowedCommand(allowed []string, target string) bool {
	for _, item := range allowed {
		if item == target {
			return true
		}
	}
	return false
}

func hasAnyAllowedCommands(allowed []string, targets ...string) bool {
	for _, target := range targets {
		if hasAllowedCommand(allowed, target) {
			return true
		}
	}
	return false
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
