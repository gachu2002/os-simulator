package sim

import "testing"

func TestControlActionsBlockUnblockKill(t *testing.T) {
	e := NewEngine(1, 0)

	if err := e.Execute(Command{Name: "spawn", Process: "p1", Program: "COMPUTE 2; EXIT"}); err != nil {
		t.Fatalf("spawn failed: %v", err)
	}

	if err := e.Execute(Command{Name: "block_process"}); err != nil {
		t.Fatalf("block_process failed: %v", err)
	}

	proc, ok := e.procs.Get(1)
	if !ok {
		t.Fatalf("missing process pid=1")
	}
	if proc.State != ProcStateBlocked {
		t.Fatalf("state=%s want=%s", proc.State, ProcStateBlocked)
	}

	if err := e.Execute(Command{Name: "unblock_process"}); err != nil {
		t.Fatalf("unblock_process failed: %v", err)
	}
	if proc.State != ProcStateReady {
		t.Fatalf("state=%s want=%s", proc.State, ProcStateReady)
	}

	if err := e.Execute(Command{Name: "kill_process", Process: "1"}); err != nil {
		t.Fatalf("kill_process failed: %v", err)
	}
	if proc.State != ProcStateTerminated {
		t.Fatalf("state=%s want=%s", proc.State, ProcStateTerminated)
	}

	if got := e.SchedulingMetrics().CompletedProcesses; got != 1 {
		t.Fatalf("completed_processes=%d want=1", got)
	}
}
