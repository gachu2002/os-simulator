package sim

import (
	"fmt"
	"sort"
)

type ProcState string

const (
	ProcStateNew        ProcState = "new"
	ProcStateReady      ProcState = "ready"
	ProcStateRunning    ProcState = "running"
	ProcStateBlocked    ProcState = "blocked"
	ProcStateTerminated ProcState = "terminated"
)

func CanTransition(from, to ProcState) bool {
	switch from {
	case ProcStateNew:
		return to == ProcStateReady
	case ProcStateReady:
		return to == ProcStateRunning || to == ProcStateBlocked || to == ProcStateTerminated
	case ProcStateRunning:
		return to == ProcStateReady || to == ProcStateBlocked || to == ProcStateTerminated
	case ProcStateBlocked:
		return to == ProcStateReady || to == ProcStateTerminated
	default:
		return false
	}
}

type TrapFrame struct {
	PC        uint64 `json:"pc"`
	SP        uint64 `json:"sp"`
	Mode      string `json:"mode"`
	SyscallNo int    `json:"syscall_no,omitempty"`
}

type Instruction struct {
	Op      string
	Arg     int
	ArgText string
	Syscall string
	Addr    uint64
	Access  AccessType
}

type Process struct {
	PID          int
	Name         string
	State        ProcState
	Trap         TrapFrame
	Program      []Instruction
	ProgramIndex int
	Remaining    int
	BlockedUntil Tick
	BlockedOn    string
	SpawnTick    Tick
	NextFD       int
	OpenFiles    map[int]OpenFile
}

func (p *Process) transition(to ProcState) error {
	if !CanTransition(p.State, to) {
		return fmt.Errorf("illegal transition %s -> %s", p.State, to)
	}
	p.State = to
	return nil
}

func (p *Process) snapshot() ProcessSnapshot {
	return ProcessSnapshot{
		PID:          p.PID,
		Name:         p.Name,
		State:        p.State,
		PC:           p.ProgramIndex,
		BlockedUntil: p.BlockedUntil,
	}
}

type ProcessTable struct {
	nextPID  int
	byPID    map[int]*Process
	pidOrder []int
}

func NewProcessTable() *ProcessTable {
	return &ProcessTable{
		nextPID: 1,
		byPID:   map[int]*Process{},
	}
}

func (pt *ProcessTable) Create(name string, program []Instruction, spawnTick Tick) *Process {
	pid := pt.nextPID
	pt.nextPID++

	proc := &Process{
		PID:       pid,
		Name:      name,
		State:     ProcStateNew,
		Trap:      TrapFrame{Mode: "user", SP: 0x1000},
		Program:   append([]Instruction(nil), program...),
		SpawnTick: spawnTick,
		NextFD:    3,
		OpenFiles: map[int]OpenFile{},
	}

	pt.byPID[pid] = proc
	pt.pidOrder = append(pt.pidOrder, pid)
	return proc
}

func (pt *ProcessTable) Get(pid int) (*Process, bool) {
	p, ok := pt.byPID[pid]
	return p, ok
}

func (pt *ProcessTable) AllSnapshots() []ProcessSnapshot {
	pids := append([]int(nil), pt.pidOrder...)
	sort.Ints(pids)
	out := make([]ProcessSnapshot, 0, len(pids))
	for _, pid := range pids {
		out = append(out, pt.byPID[pid].snapshot())
	}
	return out
}
