package sim

import "fmt"

func (e *Engine) executeSyscall(proc *Process, name string, arg int, argText string) error {
	proc.Trap.Mode = "kernel"
	proc.Trap.SyscallNo = syscallNumber(name)
	e.trace = append(e.trace, TraceEvent{Tick: e.clock, Kind: "trap.enter", Data: fmt.Sprintf("pid=%d sys=%s", proc.PID, name)})
	e.trace = append(e.trace, TraceEvent{Tick: e.clock, Kind: "trap.save", Data: fmt.Sprintf("pid=%d pc=%d", proc.PID, proc.Trap.PC)})
	e.trace = append(e.trace, TraceEvent{Tick: e.clock, Kind: "sys.dispatch", Data: fmt.Sprintf("pid=%d sys=%s", proc.PID, name)})

	result, err := e.dispatcher.Handle(proc, name, arg, argText)
	if err != nil {
		return err
	}

	if result.Blocked {
		if result.AsyncDevice != "" {
			req := e.devices.Submit(e.clock, proc.PID, result.FD, result.AsyncDevice, result.AsyncOp, result.AsyncBytes)
			e.Schedule(req.CompleteAt, IRQEventKind(req.Device), fmt.Sprintf("req=%d", req.ID))
			proc.BlockedOn = "io"
			proc.BlockedUntil = req.CompleteAt
			e.trace = append(e.trace, TraceEvent{Tick: e.clock, Kind: "io.submit", Data: fmt.Sprintf("req=%d pid=%d fd=%d device=%s op=%s n=%d done=%d", req.ID, req.PID, req.FD, req.Device, req.Op, req.Bytes, req.CompleteAt)})
		} else {
			proc.BlockedOn = "sleep"
			proc.BlockedUntil = e.clock + result.SleepTicks
			e.trace = append(e.trace, TraceEvent{Tick: e.clock, Kind: "sys.sleep", Data: fmt.Sprintf("pid=%d until=%d", proc.PID, proc.BlockedUntil)})
		}
		_ = proc.transition(ProcStateBlocked)
		e.scheduler.OnBlock(proc.PID)
		e.runningPID = 0
	}

	if name == SysRead {
		e.trace = append(e.trace, TraceEvent{Tick: e.clock, Kind: "sys.read", Data: fmt.Sprintf("pid=%d n=%d", proc.PID, arg)})
	}
	if name == SysWrite {
		e.trace = append(e.trace, TraceEvent{Tick: e.clock, Kind: "sys.write", Data: fmt.Sprintf("pid=%d n=%d", proc.PID, arg)})
	}
	if name == SysOpen {
		e.trace = append(e.trace, TraceEvent{Tick: e.clock, Kind: "sys.open", Data: fmt.Sprintf("pid=%d fd=%d path=%s", proc.PID, result.ReturnValue, result.Path)})
		e.trace = append(e.trace, TraceEvent{Tick: e.clock, Kind: "fs.path", Data: fmt.Sprintf("pid=%d path=%s traversal=%v", proc.PID, result.Path, result.Traversal)})
	}

	if result.Exit {
		e.finishProcess(proc)
		e.runningPID = 0
	}

	e.trace = append(e.trace, TraceEvent{Tick: e.clock, Kind: "trap.return", Data: fmt.Sprintf("pid=%d ret=%d", proc.PID, result.ReturnValue)})
	proc.Trap.Mode = "user"
	proc.Trap.SyscallNo = 0
	return nil
}

func syscallNumber(name string) int {
	switch name {
	case SysOpen:
		return 2
	case SysRead:
		return 3
	case SysWrite:
		return 4
	case SysSleep:
		return 5
	case SysExit:
		return 6
	default:
		return 0
	}
}
