package sim

import (
	"fmt"
	"strconv"
	"strings"
)

func (e *Engine) Execute(cmd Command) error {
	switch cmd.Name {
	case "step":
		if cmd.Count < 0 {
			return fmt.Errorf("step count must be non-negative")
		}
		e.Run(cmd.Count)
		return nil
	case "schedule":
		e.Schedule(cmd.Tick, cmd.Kind, cmd.Data)
		return nil
	case "spawn":
		program, err := ParseProgram(cmd.Program)
		if err != nil {
			return err
		}
		if cmd.Process == "" {
			cmd.Process = "proc-" + strconv.Itoa(e.procs.nextPID)
		}
		proc := e.procs.Create(cmd.Process, program, e.clock)
		if err := proc.transition(ProcStateReady); err != nil {
			return err
		}
		e.scheduler.OnReady(proc.PID, false)
		e.memory.EnsureProcess(proc.PID)
		e.stats[proc.PID] = &procStats{}
		e.trace = append(e.trace, TraceEvent{Tick: e.clock, Kind: "proc.spawn", Data: fmt.Sprintf("pid=%d name=%s", proc.PID, proc.Name)})
		return nil
	case "policy":
		return e.SetSchedulingPolicy(cmd.Policy, cmd.Quantum)
	case "block_process":
		return e.blockProcess(cmd.Process)
	case "unblock_process":
		return e.unblockProcess(cmd.Process)
	case "kill_process":
		return e.killProcess(cmd.Process)
	case "preempt_current_job":
		return e.preemptCurrentJob()
	case "choose_next_process":
		return e.chooseNextProcess(cmd.Process)
	default:
		return fmt.Errorf("unknown command %q", cmd.Name)
	}
}

func (e *Engine) preemptCurrentJob() error {
	if e.runningPID == 0 {
		return fmt.Errorf("no running process to preempt")
	}
	proc, ok := e.procs.Get(e.runningPID)
	if !ok {
		return fmt.Errorf("running process not found")
	}
	if proc.State != ProcStateRunning {
		return fmt.Errorf("running pid=%d is in state %s", proc.PID, proc.State)
	}
	if err := proc.transition(ProcStateReady); err != nil {
		return err
	}
	e.scheduler.OnReady(proc.PID, false)
	e.trace = append(e.trace, TraceEvent{Tick: e.clock, Kind: "proc.preempt.manual", Data: fmt.Sprintf("pid=%d", proc.PID)})
	e.runningPID = 0
	return nil
}

func (e *Engine) chooseNextProcess(target string) error {
	pid, proc, err := e.resolveReadyCandidate(target)
	if err != nil {
		return err
	}
	if e.runningPID != 0 && e.runningPID != pid {
		if running, ok := e.procs.Get(e.runningPID); ok && running.State == ProcStateRunning {
			if err := running.transition(ProcStateReady); err != nil {
				return err
			}
			e.scheduler.OnReady(running.PID, false)
		}
		e.runningPID = 0
	}

	if proc.State == ProcStateReady {
		e.scheduler.RemoveReady(pid)
		if err := proc.transition(ProcStateRunning); err != nil {
			return err
		}
		e.scheduler.OnDispatch(pid)
		e.runningPID = pid
		st := e.ensureStats(pid)
		if !st.hasDispatched {
			st.hasDispatched = true
			st.firstDispatch = e.clock
		}
	}
	e.trace = append(e.trace, TraceEvent{Tick: e.clock, Kind: "proc.choose_next", Data: fmt.Sprintf("pid=%d", pid)})
	return nil
}

func (e *Engine) resolveReadyCandidate(target string) (int, *Process, error) {
	if pid, proc, ok := e.resolveByTarget(target); ok {
		if proc.State == ProcStateReady || proc.State == ProcStateRunning {
			return pid, proc, nil
		}
		return 0, nil, fmt.Errorf("process pid=%d is in state %s, required ready or running", pid, proc.State)
	}
	for _, snap := range e.procs.AllSnapshots() {
		if snap.State == ProcStateReady {
			proc, _ := e.procs.Get(snap.PID)
			return proc.PID, proc, nil
		}
	}
	if e.runningPID != 0 {
		if proc, ok := e.procs.Get(e.runningPID); ok {
			return proc.PID, proc, nil
		}
	}
	return 0, nil, fmt.Errorf("no process available to choose as next")
}

func (e *Engine) blockProcess(target string) error {
	pid, proc, err := e.resolveBlockCandidate(target)
	if err != nil {
		return err
	}
	if err := proc.transition(ProcStateBlocked); err != nil {
		return err
	}
	proc.BlockedOn = "manual"
	proc.BlockedUntil = 0
	e.scheduler.OnBlock(pid)
	if e.runningPID == pid {
		e.runningPID = 0
	}
	e.trace = append(e.trace, TraceEvent{Tick: e.clock, Kind: "proc.block", Data: fmt.Sprintf("pid=%d", pid)})
	return nil
}

func (e *Engine) resolveBlockCandidate(target string) (int, *Process, error) {
	if pid, proc, ok := e.resolveByTarget(target); ok {
		if proc.State != ProcStateRunning && proc.State != ProcStateReady {
			return 0, nil, fmt.Errorf("process pid=%d is in state %s, required running or ready", pid, proc.State)
		}
		return pid, proc, nil
	}
	if e.runningPID != 0 {
		if proc, ok := e.procs.Get(e.runningPID); ok && proc.State == ProcStateRunning {
			return proc.PID, proc, nil
		}
	}
	for _, snap := range e.procs.AllSnapshots() {
		if snap.State == ProcStateReady {
			proc, _ := e.procs.Get(snap.PID)
			return proc.PID, proc, nil
		}
	}
	return 0, nil, fmt.Errorf("no process in required state running or ready")
}

func (e *Engine) unblockProcess(target string) error {
	pid, proc, err := e.resolveProcessForControl(target, ProcStateBlocked)
	if err != nil {
		return err
	}
	if err := proc.transition(ProcStateReady); err != nil {
		return err
	}
	proc.BlockedOn = ""
	proc.BlockedUntil = 0
	e.scheduler.OnReady(pid, true)
	e.trace = append(e.trace, TraceEvent{Tick: e.clock, Kind: "proc.unblock", Data: fmt.Sprintf("pid=%d", pid)})
	return nil
}

func (e *Engine) killProcess(target string) error {
	pid, proc, err := e.resolveAnyProcessForKill(target)
	if err != nil {
		return err
	}
	if proc.State == ProcStateTerminated {
		return fmt.Errorf("process pid=%d is already terminated", pid)
	}
	if err := proc.transition(ProcStateTerminated); err != nil {
		return err
	}
	e.scheduler.OnExit(pid)
	if e.runningPID == pid {
		e.runningPID = 0
	}
	st := e.ensureStats(pid)
	st.completed = true
	st.completedAt = e.clock
	e.trace = append(e.trace, TraceEvent{Tick: e.clock, Kind: "proc.kill", Data: fmt.Sprintf("pid=%d", pid)})
	return nil
}

func (e *Engine) resolveAnyProcessForKill(target string) (int, *Process, error) {
	if pid, proc, ok := e.resolveByTarget(target); ok {
		return pid, proc, nil
	}
	if e.runningPID != 0 {
		if proc, ok := e.procs.Get(e.runningPID); ok {
			return proc.PID, proc, nil
		}
	}
	for _, snap := range e.procs.AllSnapshots() {
		if snap.State == ProcStateReady || snap.State == ProcStateBlocked || snap.State == ProcStateRunning {
			proc, _ := e.procs.Get(snap.PID)
			return proc.PID, proc, nil
		}
	}
	return 0, nil, fmt.Errorf("no process available to kill")
}

func (e *Engine) resolveProcessForControl(target string, required ProcState) (int, *Process, error) {
	if pid, proc, ok := e.resolveByTarget(target); ok {
		if proc.State != required {
			return 0, nil, fmt.Errorf("process pid=%d is in state %s, required %s", pid, proc.State, required)
		}
		return pid, proc, nil
	}
	if required == ProcStateRunning && e.runningPID != 0 {
		if proc, ok := e.procs.Get(e.runningPID); ok && proc.State == required {
			return proc.PID, proc, nil
		}
	}
	for _, snap := range e.procs.AllSnapshots() {
		if snap.State == required {
			proc, _ := e.procs.Get(snap.PID)
			return proc.PID, proc, nil
		}
	}
	return 0, nil, fmt.Errorf("no process in required state %s", required)
}

func (e *Engine) resolveByTarget(target string) (int, *Process, bool) {
	target = strings.TrimSpace(target)
	if target == "" {
		return 0, nil, false
	}
	if pid, err := strconv.Atoi(target); err == nil {
		if proc, ok := e.procs.Get(pid); ok {
			return pid, proc, true
		}
	}
	for _, snap := range e.procs.AllSnapshots() {
		if snap.Name == target {
			proc, _ := e.procs.Get(snap.PID)
			return proc.PID, proc, true
		}
	}
	return 0, nil, false
}
