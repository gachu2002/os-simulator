package sim

import "fmt"

func (e *Engine) wakeBlockedProcesses() {
	for _, snap := range e.procs.AllSnapshots() {
		proc, _ := e.procs.Get(snap.PID)
		if proc.State != ProcStateBlocked || proc.BlockedOn != "sleep" || proc.BlockedUntil > e.clock {
			continue
		}
		_ = proc.transition(ProcStateReady)
		proc.BlockedOn = ""
		e.scheduler.OnReady(proc.PID, true)
		e.trace = append(e.trace, TraceEvent{Tick: e.clock, Kind: "proc.wakeup", Data: fmt.Sprintf("pid=%d", proc.PID)})
	}
}

func (e *Engine) accumulateWaitTicks() {
	for _, snap := range e.procs.AllSnapshots() {
		if snap.State != ProcStateReady {
			continue
		}
		st := e.ensureStats(snap.PID)
		st.waitTicks++
	}
}

func (e *Engine) stepProcessCPU() {
	if e.runningPID != 0 {
		proc, ok := e.procs.Get(e.runningPID)
		if !ok || proc.State != ProcStateRunning {
			e.runningPID = 0
		}
	}

	if e.runningPID == 0 {
		pid, ok := e.scheduler.Next()
		if !ok {
			e.gantt = append(e.gantt, GanttSlice{Tick: e.clock, PID: 0})
			return
		}
		proc, _ := e.procs.Get(pid)
		_ = proc.transition(ProcStateRunning)
		e.scheduler.OnDispatch(pid)
		e.runningPID = pid
		st := e.ensureStats(pid)
		if !st.hasDispatched {
			st.hasDispatched = true
			st.firstDispatch = e.clock
		}
		e.trace = append(e.trace, TraceEvent{Tick: e.clock, Kind: "proc.dispatch", Data: fmt.Sprintf("pid=%d", pid)})
	}

	proc, _ := e.procs.Get(e.runningPID)
	e.gantt = append(e.gantt, GanttSlice{Tick: e.clock, PID: proc.PID})
	e.ensureStats(proc.PID).runTicks++
	if err := e.executeInstruction(proc); err != nil {
		e.trace = append(e.trace, TraceEvent{Tick: e.clock, Kind: "proc.error", Data: fmt.Sprintf("pid=%d err=%s", proc.PID, err.Error())})
		e.finishProcess(proc)
		e.runningPID = 0
		return
	}

	if e.runningPID == 0 {
		return
	}

	if e.scheduler.OnTick(proc.PID) {
		_ = proc.transition(ProcStateReady)
		e.scheduler.OnReady(proc.PID, false)
		e.trace = append(e.trace, TraceEvent{Tick: e.clock, Kind: "proc.preempt", Data: fmt.Sprintf("pid=%d", proc.PID)})
		e.runningPID = 0
	}
}

func (e *Engine) executeInstruction(proc *Process) error {
	if proc.ProgramIndex >= len(proc.Program) {
		e.finishProcess(proc)
		e.runningPID = 0
		return nil
	}

	inst := proc.Program[proc.ProgramIndex]
	proc.Trap.PC = uint64(proc.ProgramIndex)

	switch inst.Op {
	case "COMPUTE":
		if proc.Remaining == 0 {
			proc.Remaining = inst.Arg
		}
		proc.Remaining--
		e.trace = append(e.trace, TraceEvent{Tick: e.clock, Kind: "proc.compute", Data: fmt.Sprintf("pid=%d pc=%d remaining=%d", proc.PID, proc.ProgramIndex, proc.Remaining)})
		if proc.Remaining == 0 {
			proc.ProgramIndex++
		}
		return nil
	case "SYSCALL":
		proc.ProgramIndex++
		return e.executeSyscall(proc, inst.Syscall, inst.Arg, inst.ArgText)
	case "ACCESS":
		proc.ProgramIndex++
		pa, fault, err := e.memory.Access(proc.PID, inst.Addr, inst.Access)
		if fault != "" {
			e.trace = append(e.trace, TraceEvent{Tick: e.clock, Kind: "mem.fault", Data: fmt.Sprintf("pid=%d va=%d kind=%s", proc.PID, inst.Addr, fault)})
		}
		if err != nil {
			return err
		}
		e.trace = append(e.trace, TraceEvent{Tick: e.clock, Kind: "mem.access", Data: fmt.Sprintf("pid=%d va=%d pa=%d mode=%s", proc.PID, inst.Addr, pa, inst.Access)})
		return nil
	case "BLOCK":
		proc.ProgramIndex++
		return e.executeSyscall(proc, SysSleep, inst.Arg, "")
	case "EXIT":
		proc.ProgramIndex = len(proc.Program)
		return e.executeSyscall(proc, SysExit, 0, "")
	default:
		return fmt.Errorf("unknown op %q", inst.Op)
	}
}

func (e *Engine) finishProcess(proc *Process) {
	_ = proc.transition(ProcStateTerminated)
	e.scheduler.OnExit(proc.PID)
	st := e.ensureStats(proc.PID)
	st.completed = true
	st.completedAt = e.clock
	e.trace = append(e.trace, TraceEvent{Tick: e.clock, Kind: "proc.exit", Data: fmt.Sprintf("pid=%d", proc.PID)})
}

func (e *Engine) ensureStats(pid int) *procStats {
	if _, ok := e.stats[pid]; !ok {
		e.stats[pid] = &procStats{}
	}
	return e.stats[pid]
}
