package sim

import (
	"encoding/json"
	"fmt"
	"os"
)

type ReplayLog struct {
	Seed        uint64       `json:"seed"`
	Commands    []Command    `json:"commands"`
	Trace       []TraceEvent `json:"trace"`
	TraceHash   string       `json:"trace_hash"`
	Checkpoints []Snapshot   `json:"checkpoints,omitempty"`
}

func WriteReplayLog(path string, log ReplayLog) error {
	data, err := json.MarshalIndent(log, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal replay log: %w", err)
	}

	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("write replay log: %w", err)
	}

	return nil
}

func ReadReplayLog(path string) (ReplayLog, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return ReplayLog{}, fmt.Errorf("read replay log: %w", err)
	}

	var log ReplayLog
	if err := json.Unmarshal(data, &log); err != nil {
		return ReplayLog{}, fmt.Errorf("decode replay log: %w", err)
	}

	return log, nil
}

func ReplayFromLog(log ReplayLog, checkpointEvery Tick) (ReplayLog, error) {
	engine := NewEngine(log.Seed, checkpointEvery)
	return engine.ReplayLog(log.Commands)
}
