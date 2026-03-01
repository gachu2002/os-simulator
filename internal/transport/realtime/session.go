package realtime

import (
	"fmt"
	"sync"
	"sync/atomic"

	"os-simulator-plan/internal/lessons"
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
	s := &Session{id: id, engine: e, cfg: cfg, runtime: cfg}
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
	runtime SessionConfig
	policy  *ChallengeCommandPolicy
	nextSeq uint64
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

func (s *Session) EmitError(message string) Event {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.errorEventLocked(fmt.Errorf("%s", message))
}

func (s *Session) SetChallengePolicy(policy ChallengeCommandPolicy) {
	s.mu.Lock()
	defer s.mu.Unlock()
	copyPolicy := policy.Clone()
	s.policy = &copyPolicy
}

func (s *Session) StageOutput() lessons.StageOutput {
	s.mu.Lock()
	defer s.mu.Unlock()
	return lessons.StageOutput{
		Trace:        s.engine.Trace(),
		Processes:    s.engine.ProcessTable(),
		Metrics:      s.engine.SchedulingMetrics(),
		Memory:       s.engine.MemoryView(),
		FilesystemOK: s.engine.ValidateFilesystem() == nil,
	}
}

func (s *Session) applyLocked(cmd Command) error {
	if s.policy != nil {
		if err := s.policy.Validate(cmd); err != nil {
			return err
		}
	}

	simCmd := sim.Command{Name: cmd.Name, Count: cmd.Count, Process: cmd.Process, Program: cmd.Program, Policy: cmd.Policy, Quantum: cmd.Quantum}
	switch cmd.Name {
	case "spawn":
		return s.applySpawn(cmd, simCmd)
	case "block_process", "unblock_process", "kill_process":
		return s.engine.Execute(simCmd)
	case "preempt_current_job", "choose_next_process":
		return s.engine.Execute(simCmd)
	case "step":
		return s.applyStep(simCmd)
	case "run":
		return s.applyRun(cmd, simCmd)
	case "pause":
		return nil
	case "policy":
		return s.applyPolicy(cmd, simCmd)
	case "set_frames":
		return s.applySetFrames(cmd)
	case "set_tlb_entries":
		return s.applySetTLBEntries(cmd)
	case "set_disk_latency":
		return s.applySetDiskLatency(cmd)
	case "set_terminal_latency":
		return s.applySetTerminalLatency(cmd)
	case "reset":
		return s.applyReset()
	default:
		return fmt.Errorf("unknown command %q", cmd.Name)
	}
}

func (s *Session) applySpawn(cmd Command, simCmd sim.Command) error {
	if cmd.Program == "" {
		return fmt.Errorf("spawn requires program")
	}
	if cmd.Process == "" {
		simCmd.Process = ""
	}
	return s.engine.Execute(simCmd)
}

func (s *Session) applyStep(simCmd sim.Command) error {
	if simCmd.Count == 0 {
		simCmd.Count = 1
	}
	if simCmd.Count < 0 {
		return fmt.Errorf("step count must be non-negative")
	}
	return s.engine.Execute(simCmd)
}

func (s *Session) applyRun(cmd Command, simCmd sim.Command) error {
	if cmd.Count <= 0 {
		return fmt.Errorf("run requires positive count")
	}
	simCmd.Name = "step"
	return s.engine.Execute(simCmd)
}

func (s *Session) applyPolicy(cmd Command, simCmd sim.Command) error {
	if cmd.Policy == "" {
		return fmt.Errorf("policy requires policy name")
	}
	if cmd.Quantum == 0 && cmd.Policy == sim.PolicyRR {
		simCmd.Quantum = 2
	}
	return s.engine.Execute(simCmd)
}

func (s *Session) applySetFrames(cmd Command) error {
	if cmd.Frames <= 0 {
		return fmt.Errorf("set_frames requires positive frames")
	}
	tlb := s.runtime.TLBEntries
	if tlb <= 0 {
		tlb = cmd.Frames
	}
	s.engine.ConfigureMemory(cmd.Frames, tlb)
	s.runtime.Frames = cmd.Frames
	s.runtime.TLBEntries = tlb
	return nil
}

func (s *Session) applySetTLBEntries(cmd Command) error {
	if cmd.TLBEntries <= 0 {
		return fmt.Errorf("set_tlb_entries requires positive tlb_entries")
	}
	frames := s.runtime.Frames
	if frames <= 0 {
		frames = cmd.TLBEntries
	}
	s.engine.ConfigureMemory(frames, cmd.TLBEntries)
	s.runtime.Frames = frames
	s.runtime.TLBEntries = cmd.TLBEntries
	return nil
}

func (s *Session) applySetDiskLatency(cmd Command) error {
	if cmd.DiskLatency <= 0 {
		return fmt.Errorf("set_disk_latency requires positive disk_latency")
	}
	s.engine.ConfigureDevices(cmd.DiskLatency, s.runtime.TerminalLatency)
	s.runtime.DiskLatency = cmd.DiskLatency
	return nil
}

func (s *Session) applySetTerminalLatency(cmd Command) error {
	if cmd.TerminalLatency <= 0 {
		return fmt.Errorf("set_terminal_latency requires positive terminal_latency")
	}
	s.engine.ConfigureDevices(s.runtime.DiskLatency, cmd.TerminalLatency)
	s.runtime.TerminalLatency = cmd.TerminalLatency
	return nil
}

func (s *Session) applyReset() error {
	e := sim.NewEngine(s.cfg.Seed, s.cfg.CheckpointEvery)
	e.ConfigureMemory(s.cfg.Frames, s.cfg.TLBEntries)
	e.ConfigureDevices(s.cfg.DiskLatency, s.cfg.TerminalLatency)
	if err := e.SetSchedulingPolicy(s.cfg.Policy, s.cfg.Quantum); err != nil {
		return err
	}
	s.engine = e
	s.runtime = s.cfg
	return nil
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
	if s.policy != nil {
		usage := s.policy.Usage()
		remainingSteps := s.policy.MaxSteps - usage.UsedSteps
		if remainingSteps < 0 {
			remainingSteps = 0
		}
		remainingPolicy := s.policy.MaxPolicyChanges - usage.UsedPolicyChanges
		if remainingPolicy < 0 {
			remainingPolicy = 0
		}
		remainingConfig := s.policy.MaxConfigChanges - usage.UsedConfigChanges
		if remainingConfig < 0 {
			remainingConfig = 0
		}
		snapshot.Challenge = &ChallengeStateDTO{
			MaxSteps:           s.policy.MaxSteps,
			MaxPolicyChanges:   s.policy.MaxPolicyChanges,
			MaxConfigChanges:   s.policy.MaxConfigChanges,
			UsedSteps:          usage.UsedSteps,
			UsedPolicyChanges:  usage.UsedPolicyChanges,
			UsedConfigChanges:  usage.UsedConfigChanges,
			RemainingSteps:     remainingSteps,
			RemainingPolicyOps: remainingPolicy,
			RemainingConfigOps: remainingConfig,
		}
	}
	return Event{
		Type:      "session.snapshot",
		Sequence:  s.nextSeq,
		SessionID: s.id,
		Snapshot:  snapshot,
	}
}
