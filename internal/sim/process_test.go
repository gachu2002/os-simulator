package sim

import "testing"

func TestCanTransitionMatrix(t *testing.T) {
	if !CanTransition(ProcStateNew, ProcStateReady) {
		t.Fatalf("new -> ready should be legal")
	}
	if !CanTransition(ProcStateReady, ProcStateRunning) {
		t.Fatalf("ready -> running should be legal")
	}
	if !CanTransition(ProcStateRunning, ProcStateBlocked) {
		t.Fatalf("running -> blocked should be legal")
	}
	if !CanTransition(ProcStateRunning, ProcStateTerminated) {
		t.Fatalf("running -> terminated should be legal")
	}
	if !CanTransition(ProcStateBlocked, ProcStateReady) {
		t.Fatalf("blocked -> ready should be legal")
	}

	if !CanTransition(ProcStateReady, ProcStateBlocked) {
		t.Fatalf("ready -> blocked should be legal")
	}
	if !CanTransition(ProcStateReady, ProcStateTerminated) {
		t.Fatalf("ready -> terminated should be legal")
	}
	if CanTransition(ProcStateBlocked, ProcStateRunning) {
		t.Fatalf("blocked -> running should be illegal")
	}
	if !CanTransition(ProcStateBlocked, ProcStateTerminated) {
		t.Fatalf("blocked -> terminated should be legal")
	}
}

func TestProcessLifecycleWithPseudoProgram(t *testing.T) {
	engine := NewEngine(99, 3)
	commands := []Command{
		{Name: "spawn", Process: "p1", Program: "COMPUTE 2; BLOCK 2; COMPUTE 1; EXIT"},
		{Name: "step", Count: 12},
	}

	if err := engine.ExecuteAll(commands); err != nil {
		t.Fatalf("execute commands failed: %v", err)
	}

	table := engine.ProcessTable()
	if len(table) != 1 {
		t.Fatalf("expected one process, got %d", len(table))
	}

	if table[0].State != ProcStateTerminated {
		t.Fatalf("expected terminated state, got %s", table[0].State)
	}

	if table[0].PC != 4 {
		t.Fatalf("expected pc=4 after completion, got %d", table[0].PC)
	}
}
