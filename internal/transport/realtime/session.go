package realtime

import (
	"fmt"
	"sync"
	"sync/atomic"

	"os-simulator-plan/internal/sim"
)

type SessionManager struct {
	nextID   atomic.Uint64
	mu       sync.RWMutex
	sessions map[string]*Session
}

func NewSessionManager() *SessionManager {
	return &SessionManager{sessions: map[string]*Session{}}
}

func (m *SessionManager) Create(cfg SessionConfig) (*Session, error) {
	cfg = cfg.withDefaults()
	e := sim.NewEngine(cfg.Seed, cfg.CheckpointEvery)
	e.ConfigureMemory(cfg.Frames, cfg.TLBEntries)
	e.ConfigureDevices(cfg.DiskLatency, cfg.TerminalLatency)
	if err := e.SetSchedulingPolicy(cfg.Policy, cfg.Quantum); err != nil {
		return nil, err
	}
	id := fmt.Sprintf("s-%06d", m.nextID.Add(1))
	s := &Session{id: id, engine: e, cfg: cfg}
	m.mu.Lock()
	m.sessions[id] = s
	m.mu.Unlock()
	return s, nil
}

func (m *SessionManager) Get(id string) (*Session, bool) {
	m.mu.RLock()
	s, ok := m.sessions[id]
	m.mu.RUnlock()
	return s, ok
}

type Session struct {
	id      string
	mu      sync.Mutex
	engine  *sim.Engine
	cfg     SessionConfig
	nextSeq uint64
	paused  bool
}

func (s *Session) ID() string {
	return s.id
}

func (s *Session) SnapshotEvent(lastCommand string) Event {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.snapshotEventLocked(lastCommand)
}

func (s *Session) Apply(cmd Command) Event {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := s.applyLocked(cmd); err != nil {
		return s.errorEventLocked(err)
	}
	return s.snapshotEventLocked(cmd.Name)
}

func (s *Session) applyLocked(cmd Command) error {
	simCmd := sim.Command{Name: cmd.Name, Count: cmd.Count, Process: cmd.Process, Program: cmd.Program, Policy: cmd.Policy, Quantum: cmd.Quantum}
	switch cmd.Name {
	case "spawn":
		if cmd.Program == "" {
			return fmt.Errorf("spawn requires program")
		}
		if cmd.Process == "" {
			simCmd.Process = ""
		}
		return s.engine.Execute(simCmd)
	case "step":
		if cmd.Count == 0 {
			simCmd.Count = 1
		}
		if simCmd.Count < 0 {
			return fmt.Errorf("step count must be non-negative")
		}
		return s.engine.Execute(simCmd)
	case "run":
		if cmd.Count <= 0 {
			return fmt.Errorf("run requires positive count")
		}
		s.paused = false
		simCmd.Name = "step"
		return s.engine.Execute(simCmd)
	case "pause":
		s.paused = true
		return nil
	case "policy":
		if cmd.Policy == "" {
			return fmt.Errorf("policy requires policy name")
		}
		if cmd.Quantum == 0 && cmd.Policy == sim.PolicyRR {
			simCmd.Quantum = 2
		}
		return s.engine.Execute(simCmd)
	case "reset":
		s.paused = false
		e := sim.NewEngine(s.cfg.Seed, s.cfg.CheckpointEvery)
		e.ConfigureMemory(s.cfg.Frames, s.cfg.TLBEntries)
		e.ConfigureDevices(s.cfg.DiskLatency, s.cfg.TerminalLatency)
		if err := e.SetSchedulingPolicy(s.cfg.Policy, s.cfg.Quantum); err != nil {
			return err
		}
		s.engine = e
		return nil
	default:
		return fmt.Errorf("unknown command %q", cmd.Name)
	}
}

func (s *Session) errorEventLocked(err error) Event {
	s.nextSeq++
	return Event{
		Type:      "session.error",
		Sequence:  s.nextSeq,
		SessionID: s.id,
		Error:     err.Error(),
	}
}

func (s *Session) snapshotEventLocked(lastCommand string) Event {
	s.nextSeq++
	trace := s.engine.Trace()
	snapshot := &SnapshotDTO{
		ProtocolVersion: ProtocolVersion,
		SessionID:       s.id,
		Tick:            s.engine.SchedulingMetrics().TotalTicks,
		TraceHash:       sim.TraceHash(trace),
		TraceLength:     len(trace),
		Processes:       s.engine.ProcessTable(),
		Metrics:         s.engine.SchedulingMetrics(),
		Memory:          s.engine.MemoryView(),
		LastCommand:     lastCommand,
	}
	return Event{
		Type:      "session.snapshot",
		Sequence:  s.nextSeq,
		SessionID: s.id,
		Snapshot:  snapshot,
	}
}
