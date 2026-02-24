package sim

import "testing"

func TestSyscallToIRQToWakeupFlowIsDeterministic(t *testing.T) {
	e := NewEngine(9, 0)
	e.ConfigureDevices(3, 1)
	if err := e.SetSchedulingPolicy(PolicyFIFO, 0); err != nil {
		t.Fatalf("set policy failed: %v", err)
	}
	commands := []Command{
		{Name: "spawn", Process: "io", Program: "SYSCALL open /docs/readme.txt; SYSCALL read 4; COMPUTE 1; EXIT"},
		{Name: "step", Count: 12},
	}
	if err := e.ExecuteAll(commands); err != nil {
		t.Fatalf("execute failed: %v", err)
	}

	required := []string{"sys.open", "fs.path", "io.submit", "trap.return", "irq.disk.complete", "irq.handle", "io.complete", "proc.wakeup", "proc.compute", "proc.exit"}
	trace := e.Trace()
	idx := 0
	for _, ev := range trace {
		if idx >= len(required) {
			break
		}
		if ev.Kind == required[idx] {
			idx++
		}
	}
	if idx != len(required) {
		t.Fatalf("trace missing expected syscall->irq->wakeup flow, matched %d/%d", idx, len(required))
	}

	table := e.ProcessTable()
	if len(table) != 1 || table[0].State != ProcStateTerminated {
		t.Fatalf("expected process to terminate after io completion")
	}
}

func TestWriteUsesTerminalIRQ(t *testing.T) {
	e := NewEngine(10, 0)
	e.ConfigureDevices(3, 1)
	if err := e.ExecuteAll([]Command{{Name: "spawn", Process: "tty", Program: "SYSCALL open /docs/readme.txt; SYSCALL write 3; EXIT"}, {Name: "step", Count: 6}}); err != nil {
		t.Fatalf("execute failed: %v", err)
	}

	hasTerminalIRQ := false
	for _, ev := range e.Trace() {
		if ev.Kind == "irq.terminal.complete" {
			hasTerminalIRQ = true
			break
		}
	}
	if !hasTerminalIRQ {
		t.Fatalf("expected terminal completion irq in trace")
	}
}
