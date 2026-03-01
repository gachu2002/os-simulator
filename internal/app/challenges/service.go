package challenges

import (
	"fmt"
	"slices"

	"os-simulator-plan/internal/lessons"
	"os-simulator-plan/internal/sim"
)

var defaultAllowedCommands = []string{"step", "run", "pause", "policy", "reset"}

const (
	defaultMaxSteps         = 40
	defaultMaxPolicyChanges = 3
	defaultMaxConfigChanges = 2
)

type Error struct {
	HTTPStatus int
	Code       string
	Message    string
}

func (e *Error) Error() string {
	return e.Message
}

type LessonEngine interface {
	PrepareStage(lessonID string, stageIndex int) (lessons.PreparedStage, error)
	GradeStage(prepared lessons.PreparedStage, output lessons.StageOutput) lessons.StageResult
	CompletionAnalytics() lessons.CompletionAnalytics
}

type LessonEngineProvider interface {
	ForLearner(learnerID string) LessonEngine
}

type SessionConfig struct {
	Seed            uint64
	CheckpointEvery sim.Tick
	Policy          string
	Quantum         int
	Frames          int
	TLBEntries      int
	DiskLatency     sim.Tick
	TerminalLatency sim.Tick
}

type CommandPolicy struct {
	AllowedCommands  []string
	MaxSteps         int
	MaxPolicyChanges int
	MaxConfigChanges int
}

type Session interface {
	ID() string
	ApplyBootstrapCommand(cmd sim.Command) error
	SetPolicy(policy CommandPolicy)
	StageOutput() lessons.StageOutput
}

type SessionStore interface {
	Create(cfg SessionConfig) (Session, error)
	Get(id string) (Session, bool)
}

type Attempt struct {
	AttemptID string
	SessionID string
	LearnerID string
	Prepared  lessons.PreparedStage
}

type AttemptStore interface {
	Create(sessionID, learnerID string, prepared lessons.PreparedStage) Attempt
	Get(attemptID string) (Attempt, bool)
}

type StartResult struct {
	Attempt         Attempt
	AllowedCommands []string
	MaxSteps        int
	MaxPolicy       int
	MaxConfig       int
	Objective       string
}

type SubmitResult struct {
	Attempt   Attempt
	Result    lessons.StageResult
	Analytics lessons.CompletionAnalytics
	Objective string
}

type Service struct {
	engines  LessonEngineProvider
	sessions SessionStore
	attempts AttemptStore
}

func NewService(engines LessonEngineProvider, sessions SessionStore, attempts AttemptStore) *Service {
	return &Service{engines: engines, sessions: sessions, attempts: attempts}
}

func (s *Service) Start(lessonID string, stageIndex int, learnerID string) (StartResult, error) {
	engine := s.engines.ForLearner(learnerID)
	prepared, err := engine.PrepareStage(lessonID, stageIndex)
	if err != nil {
		return StartResult{}, &Error{HTTPStatus: 400, Code: "challenge_start_failed", Message: err.Error()}
	}

	session, err := s.sessions.Create(toSessionConfig(prepared))
	if err != nil {
		return StartResult{}, &Error{HTTPStatus: 400, Code: "invalid_session_config", Message: err.Error()}
	}

	allowed := allowedCommands(prepared)
	maxSteps, maxPolicy, maxConfig := limits(prepared, allowed)

	for _, cmd := range prepared.Stage.Bootstrap {
		if err := session.ApplyBootstrapCommand(cmd); err != nil {
			return StartResult{}, &Error{HTTPStatus: 400, Code: "challenge_start_failed", Message: fmt.Sprintf("bootstrap command failed: %s", err.Error())}
		}
	}

	session.SetPolicy(CommandPolicy{
		AllowedCommands:  allowed,
		MaxSteps:         maxSteps,
		MaxPolicyChanges: maxPolicy,
		MaxConfigChanges: maxConfig,
	})

	attempt := s.attempts.Create(session.ID(), learnerID, prepared)
	return StartResult{
		Attempt:         attempt,
		AllowedCommands: allowed,
		MaxSteps:        maxSteps,
		MaxPolicy:       maxPolicy,
		MaxConfig:       maxConfig,
		Objective:       objective(prepared),
	}, nil
}

func (s *Service) Submit(attemptID string, learnerID string) (SubmitResult, error) {
	attempt, ok := s.attempts.Get(attemptID)
	if !ok {
		return SubmitResult{}, &Error{HTTPStatus: 404, Code: "challenge_attempt_not_found", Message: "challenge attempt not found"}
	}
	if attempt.LearnerID != learnerID {
		return SubmitResult{}, &Error{HTTPStatus: 403, Code: "challenge_attempt_forbidden", Message: "challenge attempt belongs to another learner"}
	}

	session, ok := s.sessions.Get(attempt.SessionID)
	if !ok {
		return SubmitResult{}, &Error{HTTPStatus: 404, Code: "session_not_found", Message: "session not found"}
	}

	engine := s.engines.ForLearner(learnerID)
	result := engine.GradeStage(attempt.Prepared, session.StageOutput())
	analytics := engine.CompletionAnalytics()

	return SubmitResult{
		Attempt:   attempt,
		Result:    result,
		Analytics: analytics,
		Objective: objective(attempt.Prepared),
	}, nil
}

func toSessionConfig(prepared lessons.PreparedStage) SessionConfig {
	return SessionConfig{
		Seed:            prepared.Stage.Config.Seed,
		CheckpointEvery: 5,
		Policy:          prepared.Stage.Config.Policy,
		Quantum:         prepared.Stage.Config.Quantum,
		Frames:          prepared.Stage.Config.Frames,
		TLBEntries:      prepared.Stage.Config.TLBEntries,
		DiskLatency:     prepared.Stage.Config.DiskLatency,
		TerminalLatency: prepared.Stage.Config.TerminalLatency,
	}
}

func objective(prepared lessons.PreparedStage) string {
	if prepared.Stage.Objective != "" {
		return prepared.Stage.Objective
	}
	return prepared.Stage.Title
}

func allowedCommands(prepared lessons.PreparedStage) []string {
	if len(prepared.Stage.AllowedCmds) == 0 {
		return slices.Clone(defaultAllowedCommands)
	}
	return slices.Clone(prepared.Stage.AllowedCmds)
}

func limits(prepared lessons.PreparedStage, allowed []string) (int, int, int) {
	maxSteps := prepared.Stage.Limits.MaxSteps
	if maxSteps <= 0 {
		maxSteps = defaultMaxSteps
	}

	maxPolicyChanges := prepared.Stage.Limits.MaxPolicyChanges
	if maxPolicyChanges <= 0 && hasAllowedCommand(allowed, "policy") {
		maxPolicyChanges = defaultMaxPolicyChanges
	}

	maxConfigChanges := prepared.Stage.Limits.MaxConfigChanges
	if maxConfigChanges <= 0 && hasAnyAllowedCommands(allowed, "set_frames", "set_tlb_entries", "set_disk_latency", "set_terminal_latency") {
		maxConfigChanges = defaultMaxConfigChanges
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
