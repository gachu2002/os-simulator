package realtime

import "os-simulator-plan/internal/sim"

const ProtocolVersion = "v1alpha1"

type SessionConfig struct {
	Seed            uint64   `json:"seed"`
	CheckpointEvery sim.Tick `json:"checkpoint_every"`
	Policy          string   `json:"policy"`
	Quantum         int      `json:"quantum"`
	Frames          int      `json:"frames"`
	TLBEntries      int      `json:"tlb_entries"`
	DiskLatency     sim.Tick `json:"disk_latency"`
	TerminalLatency sim.Tick `json:"terminal_latency"`
}

func (c SessionConfig) withDefaults() SessionConfig {
	if c.Seed == 0 {
		c.Seed = 1
	}
	if c.CheckpointEvery == 0 {
		c.CheckpointEvery = 5
	}
	if c.Policy == "" {
		c.Policy = sim.PolicyRR
	}
	if c.Quantum <= 0 {
		c.Quantum = 2
	}
	if c.Frames <= 0 {
		c.Frames = 8
	}
	if c.TLBEntries <= 0 {
		c.TLBEntries = 4
	}
	if c.DiskLatency == 0 {
		c.DiskLatency = 3
	}
	if c.TerminalLatency == 0 {
		c.TerminalLatency = 1
	}
	return c
}

type Command struct {
	Name    string `json:"name"`
	Count   int    `json:"count,omitempty"`
	Process string `json:"process,omitempty"`
	Program string `json:"program,omitempty"`
	Policy  string `json:"policy,omitempty"`
	Quantum int    `json:"quantum,omitempty"`
}

type CommandEnvelope struct {
	Type    string  `json:"type"`
	Command Command `json:"command"`
}

type SnapshotDTO struct {
	ProtocolVersion string                `json:"protocol_version"`
	SessionID       string                `json:"session_id"`
	Tick            sim.Tick              `json:"tick"`
	TraceHash       string                `json:"trace_hash"`
	TraceLength     int                   `json:"trace_length"`
	Processes       []sim.ProcessSnapshot `json:"processes"`
	Metrics         sim.SchedulingMetrics `json:"metrics"`
	Memory          sim.MemorySnapshot    `json:"memory"`
	LastCommand     string                `json:"last_command,omitempty"`
}

type Event struct {
	Type      string       `json:"type"`
	Sequence  uint64       `json:"sequence"`
	SessionID string       `json:"session_id"`
	Snapshot  *SnapshotDTO `json:"snapshot,omitempty"`
	Error     string       `json:"error,omitempty"`
}

type CreateSessionResponse struct {
	SessionID string       `json:"session_id"`
	Snapshot  *SnapshotDTO `json:"snapshot"`
}
