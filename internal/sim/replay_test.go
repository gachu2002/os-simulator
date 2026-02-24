package sim

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestReplayFromLogMatchesOriginalHash(t *testing.T) {
	commands := []Command{
		{Name: "schedule", Tick: 2, Kind: "work.io", Data: "pid=10"},
		{Name: "schedule", Tick: 2, Kind: "work.io", Data: "pid=11"},
		{Name: "schedule", Tick: 5, Kind: "work.cpu", Data: "pid=10"},
		{Name: "step", Count: 20},
	}

	originalEngine := NewEngine(42, 4)
	original, err := originalEngine.ReplayLog(commands)
	if err != nil {
		t.Fatalf("replay log failed: %v", err)
	}

	tempFile := filepath.Join(t.TempDir(), "replay.json")
	if err := WriteReplayLog(tempFile, original); err != nil {
		t.Fatalf("write replay log failed: %v", err)
	}

	loaded, err := ReadReplayLog(tempFile)
	if err != nil {
		t.Fatalf("read replay log failed: %v", err)
	}

	replayed, err := ReplayFromLog(loaded, 4)
	if err != nil {
		t.Fatalf("replay from log failed: %v", err)
	}

	if replayed.TraceHash != original.TraceHash {
		t.Fatalf("hash mismatch: replayed=%s original=%s", replayed.TraceHash, original.TraceHash)
	}
}

func TestGoldenTraceHash(t *testing.T) {
	engine := NewEngine(7, 5)
	log, err := engine.ReplayLog([]Command{{Name: "step", Count: 30}})
	if err != nil {
		t.Fatalf("replay log failed: %v", err)
	}

	goldenPath := filepath.Join("..", "..", "tests", "golden", "milestone01_seed7_steps30.hash")
	data, err := os.ReadFile(goldenPath)
	if err != nil {
		t.Fatalf("read golden hash failed: %v", err)
	}

	expected := strings.TrimSpace(string(data))
	if log.TraceHash != expected {
		t.Fatalf("golden hash mismatch: got=%s want=%s", log.TraceHash, expected)
	}
}
