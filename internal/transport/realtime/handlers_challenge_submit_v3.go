package realtime

import (
	"net/http"

	appchallenges "os-simulator-plan/internal/app/challenges"
	"os-simulator-plan/internal/sim"
)

type ChallengeSubmitV3Request struct {
	AttemptID string `json:"attempt_id"`
	LearnerID string `json:"learner_id,omitempty"`
}

type ChallengeSubmitV3Response struct {
	Version            string                 `json:"version"`
	SectionID          string                 `json:"section_id"`
	LessonID           string                 `json:"lesson_id"`
	LessonTitle        string                 `json:"lesson_title"`
	LessonObjective    string                 `json:"lesson_objective"`
	PartID             string                 `json:"part_id,omitempty"`
	PartTitle          string                 `json:"part_title,omitempty"`
	PartObjective      string                 `json:"part_objective,omitempty"`
	AttemptID          string                 `json:"attempt_id"`
	Passed             bool                   `json:"passed"`
	FeedbackKey        string                 `json:"feedback_key"`
	Hint               string                 `json:"hint,omitempty"`
	HintLevel          int                    `json:"hint_level,omitempty"`
	PassConditions     []string               `json:"pass_conditions,omitempty"`
	Output             LessonOutputDTO        `json:"output"`
	Analytics          CompletionAnalyticsDTO `json:"analytics"`
	ValidatorResults   []ValidatorResultDTO   `json:"validator_results,omitempty"`
	ActionCapabilities ActionCapabilities     `json:"action_capabilities"`
}

func (s *Server) handleChallengeSubmitV3(w http.ResponseWriter, r *http.Request) {
	if !requireMethod(w, r, http.MethodPost) {
		return
	}
	if s.cpuCurriculumErr != nil {
		respondError(w, r, http.StatusInternalServerError, "curriculum_unavailable", "curriculum unavailable")
		return
	}

	var req ChallengeSubmitV3Request
	if !decodeJSONBody(w, r, &req) {
		return
	}
	if req.AttemptID == "" {
		respondError(w, r, http.StatusBadRequest, "invalid_challenge_request", "attempt_id is required")
		return
	}

	learnerID := normalizeLearnerID(req.LearnerID)
	submit, err := s.challengeService.Submit(req.AttemptID, learnerID)
	if err != nil {
		respondChallengeServiceError(w, r, err)
		return
	}

	v3LessonID, ok := resolveV3LessonID(submit.Attempt.Prepared.LessonID)
	if !ok {
		respondError(w, r, http.StatusNotFound, "lesson_not_found", "lesson not found")
		return
	}
	section, lesson, ok := s.lessonV3ByID(v3LessonID)
	if !ok {
		respondError(w, r, http.StatusNotFound, "lesson_not_found", "lesson not found")
		return
	}

	out := ChallengeSubmitV3Response{
		Version:            s.cpuCurriculumV3.Version,
		SectionID:          section.ID,
		LessonID:           lesson.ID,
		LessonTitle:        lesson.Title,
		LessonObjective:    lesson.Objective,
		AttemptID:          submit.Attempt.AttemptID,
		Passed:             submit.Result.Passed,
		FeedbackKey:        submit.Result.FeedbackKey,
		Hint:               submit.Result.Hint,
		HintLevel:          submit.Result.HintLevel,
		PassConditions:     stagePassConditions(submit.Attempt.Prepared.Stage),
		Analytics:          convertAnalytics(submit.Analytics),
		ValidatorResults:   toValidatorResultDTO(appchallenges.BuildValidatorResults(submit.Result.ValidatorResults, submit.Attempt.Prepared.Stage, submit.Result.Output)),
		ActionCapabilities: classifyActionCapabilities(lesson.Challenge.Actions),
		Output: LessonOutputDTO{
			Tick:         submit.Result.Output.Metrics.TotalTicks,
			TraceHash:    sim.TraceHash(submit.Result.Output.Trace),
			TraceLength:  len(submit.Result.Output.Trace),
			Processes:    submit.Result.Output.Processes,
			Metrics:      submit.Result.Output.Metrics,
			Memory:       submit.Result.Output.Memory,
			FilesystemOK: submit.Result.Output.FilesystemOK,
		},
	}

	if part := lessonPartByStageIndex(lesson, submit.Attempt.Prepared.StageIndex); part != nil {
		out.PartID = part.ID
		out.PartTitle = part.Title
		out.PartObjective = part.Objective
	}

	qualityGates := evaluateV3QualityGates(out.LessonID, out.PartID, submit.Result.Output)
	if len(qualityGates) > 0 {
		out.ValidatorResults = append(out.ValidatorResults, qualityGatesToDTO(qualityGates)...)
		out.PassConditions = append(out.PassConditions, qualityGatePassConditions(qualityGates)...)
		if !allQualityGatesPassed(qualityGates) {
			out.Passed = false
			out.FeedbackKey = "v3.quality_not_met"
			out.Hint = firstQualityGateHint(qualityGates)
			out.HintLevel = 2
		}
	}

	respondJSON(w, http.StatusOK, out)
}
