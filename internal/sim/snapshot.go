package sim

type Snapshot struct {
	Tick          Tick              `json:"tick"`
	PendingEvents int               `json:"pending_events"`
	TraceLength   int               `json:"trace_length"`
	Processes     []ProcessSnapshot `json:"processes,omitempty"`
	Memory        MemorySnapshot    `json:"memory"`
}

type SnapshotManager struct {
	interval    Tick
	checkpoints []Snapshot
}

func NewSnapshotManager(interval Tick) *SnapshotManager {
	return &SnapshotManager{interval: interval}
}

func (m *SnapshotManager) MaybeCapture(s Snapshot) {
	if m.interval == 0 {
		return
	}

	if s.Tick%m.interval != 0 {
		return
	}

	m.checkpoints = append(m.checkpoints, s)
}

func (m *SnapshotManager) Checkpoints() []Snapshot {
	out := make([]Snapshot, len(m.checkpoints))
	copy(out, m.checkpoints)
	return out
}
