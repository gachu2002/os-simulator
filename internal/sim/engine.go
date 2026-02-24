package sim

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/rand"
	"sort"
	"strconv"
	"strings"
)

type procStats struct {
	firstDispatch Tick
	hasDispatched bool
	runTicks      Tick
	waitTicks     Tick
	completedAt   Tick
	completed     bool
}

type Engine struct {
	seed         uint64
	clock        Tick
	nextSequence uint64
	queue        *EventQueue
	snapshots    *SnapshotManager
	procs        *ProcessTable
	scheduler    Scheduler
	dispatcher   SyscallDispatcher
	memory       *MemoryManager
	devices      *DeviceManager
	fs           *FileSystem
	runningPID   int
	trace        []TraceEvent
	gantt        []GanttSlice
	stats        map[int]*procStats
}

func NewEngine(seed uint64, checkpointEvery Tick) *Engine {
	scheduler, _ := NewScheduler(PolicyRR, 2)
	e := &Engine{
		seed:         seed,
		nextSequence: 1,
		queue:        NewEventQueue(),
		snapshots:    NewSnapshotManager(checkpointEvery),
		procs:        NewProcessTable(),
		scheduler:    scheduler,
		fs:           NewFileSystem(),
		memory:       NewMemoryManager(8, 4),
		devices:      NewDeviceManager(3, 1),
		stats:        map[int]*procStats{},
	}
	e.dispatcher = NewKernelDispatcher(e.fs)

	e.bootstrapFromSeed(seed)
	return e
}

func (e *Engine) bootstrapFromSeed(seed uint64) {
	rng := rand.New(rand.NewSource(int64(seed)))
	for i := 0; i < 4; i++ {
		offset := Tick(rng.Intn(8) + 1)
		data := fmt.Sprintf("slot=%d", rng.Intn(1000))
		e.Schedule(offset, "bootstrap.task", data)
	}
}

func (e *Engine) SetSchedulingPolicy(policy string, quantum int) error {
	scheduler, err := NewScheduler(policy, quantum)
	if err != nil {
		return err
	}
	e.scheduler = scheduler
	e.rebuildReadyQueues()
	e.trace = append(e.trace, TraceEvent{Tick: e.clock, Kind: "sched.policy", Data: fmt.Sprintf("policy=%s quantum=%d", policy, scheduler.Quantum())})
	return nil
}

func (e *Engine) rebuildReadyQueues() {
	for _, p := range e.procs.AllSnapshots() {
		if p.State == ProcStateReady {
			e.scheduler.OnReady(p.PID, false)
		}
	}
}

func (e *Engine) ConfigureMemory(totalFrames, tlbEntries int) {
	e.memory = NewMemoryManager(totalFrames, tlbEntries)
	for _, snap := range e.procs.AllSnapshots() {
		e.memory.EnsureProcess(snap.PID)
	}
}

func (e *Engine) ConfigureDevices(diskLatency, terminalLatency Tick) {
	e.devices = NewDeviceManager(diskLatency, terminalLatency)
}

func (e *Engine) Schedule(at Tick, kind, data string) {
	e.queue.Push(Event{Tick: at, Sequence: e.nextSequence, Kind: kind, Data: data})
	e.nextSequence++
}

func (e *Engine) Step() {
	e.clock++
	e.trace = append(e.trace, TraceEvent{Tick: e.clock, Kind: "clock.tick"})
	e.wakeBlockedProcesses()
	e.accumulateWaitTicks()
	e.stepProcessCPU()

	for {
		next, ok := e.queue.Peek()
		if !ok || next.Tick > e.clock {
			break
		}
		ev, _ := e.queue.Pop()
		e.trace = append(e.trace, ev)
		e.handleEvent(ev)
	}

	e.snapshots.MaybeCapture(Snapshot{Tick: e.clock, PendingEvents: e.queue.Len(), TraceLength: len(e.trace), Processes: e.procs.AllSnapshots(), Memory: e.memory.Snapshot()})
}

func (e *Engine) Run(count int) {
	for i := 0; i < count; i++ {
		e.Step()
	}
}

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
	default:
		return fmt.Errorf("unknown command %q", cmd.Name)
	}
}

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

func (e *Engine) handleEvent(ev Event) {
	if !strings.HasPrefix(ev.Kind, "irq.") {
		return
	}
	requestID, ok := parseRequestID(ev.Data)
	if !ok {
		return
	}
	req, ok := e.devices.Complete(requestID)
	if !ok {
		return
	}
	e.trace = append(e.trace, TraceEvent{Tick: e.clock, Kind: "irq.handle", Data: fmt.Sprintf("req=%d pid=%d device=%s", req.ID, req.PID, req.Device)})
	e.trace = append(e.trace, TraceEvent{Tick: e.clock, Kind: "io.complete", Data: fmt.Sprintf("req=%d pid=%d fd=%d op=%s n=%d", req.ID, req.PID, req.FD, req.Op, req.Bytes)})
	if of, ok := procOpenFile(e.procs, req.PID, req.FD); ok {
		if req.Op == SysRead {
			data, blocks, nextOffset, err := e.fs.ReadInode(of.InodeID, req.Bytes, of.Offset)
			if err == nil {
				of.Offset = nextOffset
				setProcOpenFile(e.procs, req.PID, req.FD, of)
				e.trace = append(e.trace, TraceEvent{Tick: e.clock, Kind: "fs.read", Data: fmt.Sprintf("pid=%d fd=%d bytes=%d", req.PID, req.FD, len(data))})
				e.trace = append(e.trace, TraceEvent{Tick: e.clock, Kind: "fs.blockmap", Data: fmt.Sprintf("pid=%d fd=%d blocks=%v", req.PID, req.FD, blocks)})
			}
		}
		if req.Op == SysWrite {
			payload := []byte(strings.Repeat("w", req.Bytes))
			written, blocks, nextOffset := e.fs.WriteInode(of.InodeID, payload, of.Offset)
			of.Offset = nextOffset
			setProcOpenFile(e.procs, req.PID, req.FD, of)
			e.trace = append(e.trace, TraceEvent{Tick: e.clock, Kind: "fs.write", Data: fmt.Sprintf("pid=%d fd=%d bytes=%d", req.PID, req.FD, written)})
			e.trace = append(e.trace, TraceEvent{Tick: e.clock, Kind: "fs.blockmap", Data: fmt.Sprintf("pid=%d fd=%d blocks=%v", req.PID, req.FD, blocks)})
		}
	}

	proc, ok := e.procs.Get(req.PID)
	if !ok || proc.State != ProcStateBlocked || proc.BlockedOn != "io" {
		return
	}
	_ = proc.transition(ProcStateReady)
	proc.BlockedOn = ""
	proc.BlockedUntil = 0
	e.scheduler.OnReady(proc.PID, true)
	e.trace = append(e.trace, TraceEvent{Tick: e.clock, Kind: "proc.wakeup", Data: fmt.Sprintf("pid=%d source=irq", proc.PID)})
}

func parseRequestID(data string) (int, bool) {
	if !strings.HasPrefix(data, "req=") {
		return 0, false
	}
	v, err := strconv.Atoi(strings.TrimPrefix(data, "req="))
	if err != nil {
		return 0, false
	}
	return v, true
}

func procOpenFile(pt *ProcessTable, pid, fd int) (OpenFile, bool) {
	proc, ok := pt.Get(pid)
	if !ok {
		return OpenFile{}, false
	}
	of, ok := proc.OpenFiles[fd]
	return of, ok
}

func setProcOpenFile(pt *ProcessTable, pid, fd int, of OpenFile) {
	proc, ok := pt.Get(pid)
	if !ok {
		return
	}
	proc.OpenFiles[fd] = of
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

	if !result.Blocked && !result.Exit {
		e.trace = append(e.trace, TraceEvent{Tick: e.clock, Kind: "trap.return", Data: fmt.Sprintf("pid=%d ret=%d", proc.PID, result.ReturnValue)})
		proc.Trap.Mode = "user"
		proc.Trap.SyscallNo = 0
		return nil
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

func (e *Engine) SchedulingMetrics() SchedulingMetrics {
	table := e.procs.AllSnapshots()
	metrics := SchedulingMetrics{
		Policy:     e.scheduler.Policy(),
		Quantum:    e.scheduler.Quantum(),
		TotalTicks: e.clock,
		Gantt:      append([]GanttSlice(nil), e.gantt...),
	}

	procMetrics := make([]ProcessMetric, 0, len(table))
	var totalResp Tick
	var totalTurn Tick
	var fairSum float64
	var fairSq float64
	for _, snap := range table {
		p, _ := e.procs.Get(snap.PID)
		st := e.ensureStats(snap.PID)
		resp := Tick(0)
		if st.hasDispatched {
			resp = st.firstDispatch - p.SpawnTick
		}
		turn := Tick(0)
		if st.completed {
			turn = st.completedAt - p.SpawnTick
			metrics.CompletedProcesses++
		}
		procMetrics = append(procMetrics, ProcessMetric{PID: p.PID, Name: p.Name, ResponseTime: resp, Turnaround: turn, RunTicks: st.runTicks, WaitTicks: st.waitTicks})
		totalResp += resp
		totalTurn += turn
		r := float64(st.runTicks)
		fairSum += r
		fairSq += r * r
	}
	sort.Slice(procMetrics, func(i, j int) bool { return procMetrics[i].PID < procMetrics[j].PID })
	metrics.Processes = procMetrics
	if len(procMetrics) > 0 {
		metrics.AvgResponseTime = float64(totalResp) / float64(len(procMetrics))
		metrics.AvgTurnaroundTime = float64(totalTurn) / float64(len(procMetrics))
	}
	if e.clock > 0 {
		metrics.ThroughputPer100Tick = float64(metrics.CompletedProcesses) * 100 / float64(e.clock)
	}
	if fairSq > 0 && len(procMetrics) > 0 {
		n := float64(len(procMetrics))
		metrics.FairnessJainIndex = (fairSum * fairSum) / (n * fairSq)
	}
	return metrics
}

func (e *Engine) MemoryView() MemorySnapshot {
	return e.memory.Snapshot()
}

func (e *Engine) ValidateFilesystem() error {
	return e.fs.Invariants()
}

func TraceHash(trace []TraceEvent) string {
	b := strings.Builder{}
	for _, ev := range trace {
		fmt.Fprintf(&b, "%d|%d|%s|%s\n", ev.Tick, ev.Sequence, ev.Kind, ev.Data)
	}
	sum := sha256.Sum256([]byte(b.String()))
	return hex.EncodeToString(sum[:])
}
