package sim

import "testing"

func TestSyscallTraceOrdering(t *testing.T) {
	e := NewEngine(55, 0)
	if err := e.SetSchedulingPolicy(PolicyFIFO, 0); err != nil {
		t.Fatalf("set policy failed: %v", err)
	}
	commands := []Command{
		{Name: "spawn", Process: "p1", Program: "SYSCALL open; SYSCALL read 4; SYSCALL write 4; SYSCALL exit"},
		{Name: "step", Count: 10},
	}
	if err := e.ExecuteAll(commands); err != nil {
		t.Fatalf("execute failed: %v", err)
	}

	trace := e.Trace()
	requireOrder := []string{
		"trap.enter",
		"trap.save",
		"sys.dispatch",
		"sys.open",
		"trap.return",
		"trap.enter",
		"trap.save",
		"sys.dispatch",
		"sys.read",
		"trap.return",
	}

	idx := 0
	for _, ev := range trace {
		if idx >= len(requireOrder) {
			break
		}
		if ev.Kind == requireOrder[idx] {
			idx++
		}
	}
	if idx != len(requireOrder) {
		t.Fatalf("trace did not contain expected syscall/trap sequence, matched %d/%d", idx, len(requireOrder))
	}
}

func TestSleepSyscallBlocksAndWakes(t *testing.T) {
	e := NewEngine(88, 0)
	if err := e.ExecuteAll([]Command{
		{Name: "spawn", Process: "p1", Program: "SYSCALL sleep 2; COMPUTE 1; EXIT"},
		{Name: "step", Count: 8},
	}); err != nil {
		t.Fatalf("execute failed: %v", err)
	}

	table := e.ProcessTable()
	if len(table) != 1 {
		t.Fatalf("expected one process")
	}
	if table[0].State != ProcStateTerminated {
		t.Fatalf("expected process terminated, got %s", table[0].State)
	}

	trace := e.Trace()
	hasSleep := false
	hasWake := false
	for _, ev := range trace {
		if ev.Kind == "sys.sleep" {
			hasSleep = true
		}
		if ev.Kind == "proc.wakeup" {
			hasWake = true
		}
	}
	if !hasSleep || !hasWake {
		t.Fatalf("expected sleep and wakeup events in trace")
	}
}
