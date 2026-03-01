package sim

import (
	"fmt"
	"math/rand"
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
