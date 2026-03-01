package realtime

import (
	"errors"

	appchallenges "os-simulator-plan/internal/app/challenges"
	"os-simulator-plan/internal/lessons"
	"os-simulator-plan/internal/sim"
)

type lessonEngineProviderAdapter struct {
	server *Server
}

func (p lessonEngineProviderAdapter) ForLearner(learnerID string) appchallenges.LessonEngine {
	return p.server.lessonEngineForLearner(learnerID)
}

type sessionStoreAdapter struct {
	manager *SessionManager
}

func (s sessionStoreAdapter) Create(cfg appchallenges.SessionConfig) (appchallenges.Session, error) {
	session, err := s.manager.Create(SessionConfig{
		Seed:            cfg.Seed,
		CheckpointEvery: cfg.CheckpointEvery,
		Policy:          cfg.Policy,
		Quantum:         cfg.Quantum,
		Frames:          cfg.Frames,
		TLBEntries:      cfg.TLBEntries,
		DiskLatency:     cfg.DiskLatency,
		TerminalLatency: cfg.TerminalLatency,
	})
	if err != nil {
		return nil, err
	}
	return sessionAdapter{session: session}, nil
}

func (s sessionStoreAdapter) Get(id string) (appchallenges.Session, bool) {
	session, ok := s.manager.Get(id)
	if !ok {
		return nil, false
	}
	return sessionAdapter{session: session}, true
}

type sessionAdapter struct {
	session *Session
}

func (s sessionAdapter) ID() string {
	return s.session.ID()
}

func (s sessionAdapter) ApplyBootstrapCommand(cmd sim.Command) error {
	ev := s.session.Apply(Command{
		Name:    cmd.Name,
		Count:   cmd.Count,
		Process: cmd.Process,
		Program: cmd.Program,
		Policy:  cmd.Policy,
		Quantum: cmd.Quantum,
	})
	if ev.Type == "session.error" {
		return errors.New(ev.Error)
	}
	return nil
}

func (s sessionAdapter) SetPolicy(policy appchallenges.CommandPolicy) {
	s.session.SetChallengePolicy(NewChallengeCommandPolicy(
		policy.AllowedCommands,
		policy.MaxSteps,
		policy.MaxPolicyChanges,
		policy.MaxConfigChanges,
	))
}

func (s sessionAdapter) StageOutput() lessons.StageOutput {
	return s.session.StageOutput()
}

type challengeAttemptStoreAdapter struct {
	store *ChallengeAttemptStore
}

func (s challengeAttemptStoreAdapter) Create(sessionID, learnerID string, prepared lessons.PreparedStage) appchallenges.Attempt {
	attempt := s.store.Create(sessionID, learnerID, prepared)
	return appchallenges.Attempt{AttemptID: attempt.AttemptID, SessionID: attempt.SessionID, LearnerID: attempt.LearnerID, Prepared: attempt.Prepared}
}

func (s challengeAttemptStoreAdapter) Get(attemptID string) (appchallenges.Attempt, bool) {
	attempt, ok := s.store.Get(attemptID)
	if !ok {
		return appchallenges.Attempt{}, false
	}
	return appchallenges.Attempt{AttemptID: attempt.AttemptID, SessionID: attempt.SessionID, LearnerID: attempt.LearnerID, Prepared: attempt.Prepared}, true
}
