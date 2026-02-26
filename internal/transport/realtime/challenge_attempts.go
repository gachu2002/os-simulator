package realtime

import (
	"fmt"
	"sync"
	"sync/atomic"

	"os-simulator-plan/internal/lessons"
)

var defaultChallengeAllowedCommands = []string{"step", "run", "pause", "policy", "reset"}

const (
	defaultChallengeMaxSteps         = 40
	defaultChallengeMaxPolicyChanges = 3
	defaultChallengeMaxConfigChanges = 2
)

type ChallengeAttempt struct {
	AttemptID string
	SessionID string
	LearnerID string
	Prepared  lessons.PreparedStage
}

type ChallengeAttemptStore struct {
	nextID   atomic.Uint64
	mu       sync.RWMutex
	attempts map[string]ChallengeAttempt
}

func NewChallengeAttemptStore() *ChallengeAttemptStore {
	return &ChallengeAttemptStore{attempts: map[string]ChallengeAttempt{}}
}

func (s *ChallengeAttemptStore) Create(sessionID, learnerID string, prepared lessons.PreparedStage) ChallengeAttempt {
	attemptID := fmt.Sprintf("a-%06d", s.nextID.Add(1))
	attempt := ChallengeAttempt{AttemptID: attemptID, SessionID: sessionID, LearnerID: learnerID, Prepared: prepared}
	s.mu.Lock()
	s.attempts[attemptID] = attempt
	s.mu.Unlock()
	return attempt
}

func (s *ChallengeAttemptStore) Get(attemptID string) (ChallengeAttempt, bool) {
	s.mu.RLock()
	attempt, ok := s.attempts[attemptID]
	s.mu.RUnlock()
	return attempt, ok
}
