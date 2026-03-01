package sim

func (e *Engine) ExecuteAll(commands []Command) error {
	for _, cmd := range commands {
		if err := e.Execute(cmd); err != nil {
			return err
		}
	}
	return nil
}

func (e *Engine) ReplayLog(commands []Command) (ReplayLog, error) {
	if err := e.ExecuteAll(commands); err != nil {
		return ReplayLog{}, err
	}
	trace := e.Trace()
	return ReplayLog{Seed: e.seed, Commands: append([]Command(nil), commands...), Trace: trace, TraceHash: TraceHash(trace), Checkpoints: e.snapshots.Checkpoints()}, nil
}

func (e *Engine) Trace() []TraceEvent {
	out := make([]TraceEvent, len(e.trace))
	copy(out, e.trace)
	return out
}

func (e *Engine) ProcessTable() []ProcessSnapshot {
	return e.procs.AllSnapshots()
}

func (e *Engine) MemoryView() MemorySnapshot {
	return e.memory.Snapshot()
}

func (e *Engine) ValidateFilesystem() error {
	return e.fs.Invariants()
}
