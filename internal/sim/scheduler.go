package sim

import "fmt"

const (
	PolicyFIFO = "fifo"
	PolicyRR   = "rr"
	PolicyMLFQ = "mlfq"
)

type Scheduler interface {
	Policy() string
	Quantum() int
	OnReady(pid int, fromWake bool)
	RemoveReady(pid int)
	OnDispatch(pid int)
	OnTick(pid int) bool
	OnBlock(pid int)
	OnExit(pid int)
	Next() (int, bool)
}

func NewScheduler(policy string, quantum int) (Scheduler, error) {
	switch policy {
	case PolicyFIFO:
		return newFIFOScheduler(), nil
	case PolicyRR:
		if quantum <= 0 {
			quantum = 2
		}
		return newRRScheduler(quantum), nil
	case PolicyMLFQ:
		return newMLFQScheduler([]int{1, 2, 4}), nil
	default:
		return nil, fmt.Errorf("unknown scheduling policy %q", policy)
	}
}

type queueState struct {
	queue   []int
	inQueue map[int]bool
}

func newQueueState() queueState {
	return queueState{inQueue: map[int]bool{}}
}

func (s *queueState) Enqueue(pid int) {
	if s.inQueue[pid] {
		return
	}
	s.queue = append(s.queue, pid)
	s.inQueue[pid] = true
}

func (s *queueState) Remove(pid int) {
	if !s.inQueue[pid] {
		return
	}
	delete(s.inQueue, pid)
	for i, p := range s.queue {
		if p == pid {
			s.queue = append(s.queue[:i], s.queue[i+1:]...)
			return
		}
	}
}

func (s *queueState) Dequeue() (int, bool) {
	if len(s.queue) == 0 {
		return 0, false
	}
	p := s.queue[0]
	s.queue = s.queue[1:]
	delete(s.inQueue, p)
	return p, true
}

type fifoScheduler struct {
	ready queueState
}

func newFIFOScheduler() *fifoScheduler {
	return &fifoScheduler{ready: newQueueState()}
}

func (s *fifoScheduler) Policy() string          { return PolicyFIFO }
func (s *fifoScheduler) Quantum() int            { return 0 }
func (s *fifoScheduler) OnReady(pid int, _ bool) { s.ready.Enqueue(pid) }
func (s *fifoScheduler) RemoveReady(pid int)     { s.ready.Remove(pid) }
func (s *fifoScheduler) OnDispatch(_ int)        {}
func (s *fifoScheduler) OnTick(_ int) bool       { return false }
func (s *fifoScheduler) OnBlock(pid int)         { s.ready.Remove(pid) }
func (s *fifoScheduler) OnExit(pid int)          { s.ready.Remove(pid) }
func (s *fifoScheduler) Next() (int, bool)       { return s.ready.Dequeue() }

type rrScheduler struct {
	ready    queueState
	quantum  int
	runTicks map[int]int
}

func newRRScheduler(quantum int) *rrScheduler {
	return &rrScheduler{ready: newQueueState(), quantum: quantum, runTicks: map[int]int{}}
}

func (s *rrScheduler) Policy() string { return PolicyRR }
func (s *rrScheduler) Quantum() int   { return s.quantum }
func (s *rrScheduler) OnReady(pid int, _ bool) {
	s.ready.Enqueue(pid)
}
func (s *rrScheduler) RemoveReady(pid int) { s.ready.Remove(pid) }
func (s *rrScheduler) OnDispatch(pid int)  { s.runTicks[pid] = 0 }
func (s *rrScheduler) OnTick(pid int) bool {
	s.runTicks[pid]++
	if s.runTicks[pid] >= s.quantum {
		s.runTicks[pid] = 0
		return true
	}
	return false
}
func (s *rrScheduler) OnBlock(pid int) {
	s.ready.Remove(pid)
	delete(s.runTicks, pid)
}
func (s *rrScheduler) OnExit(pid int) {
	s.ready.Remove(pid)
	delete(s.runTicks, pid)
}
func (s *rrScheduler) Next() (int, bool) { return s.ready.Dequeue() }

type mlfqScheduler struct {
	queues   []queueState
	levels   map[int]int
	runTicks map[int]int
	quanta   []int
}

func newMLFQScheduler(quanta []int) *mlfqScheduler {
	qs := make([]queueState, len(quanta))
	for i := range qs {
		qs[i] = newQueueState()
	}
	return &mlfqScheduler{queues: qs, levels: map[int]int{}, runTicks: map[int]int{}, quanta: quanta}
}

func (s *mlfqScheduler) Policy() string { return PolicyMLFQ }
func (s *mlfqScheduler) Quantum() int   { return 0 }

func (s *mlfqScheduler) OnReady(pid int, fromWake bool) {
	if _, ok := s.levels[pid]; !ok {
		s.levels[pid] = 0
	}
	if fromWake {
		s.levels[pid] = 0
	}
	lvl := s.levels[pid]
	s.queues[lvl].Enqueue(pid)
}

func (s *mlfqScheduler) RemoveReady(pid int) {
	for i := range s.queues {
		s.queues[i].Remove(pid)
	}
}

func (s *mlfqScheduler) OnDispatch(pid int) {
	s.runTicks[pid] = 0
}

func (s *mlfqScheduler) OnTick(pid int) bool {
	lvl := s.levels[pid]
	s.runTicks[pid]++
	if s.runTicks[pid] >= s.quanta[lvl] {
		s.runTicks[pid] = 0
		if lvl < len(s.quanta)-1 {
			s.levels[pid] = lvl + 1
		}
		return true
	}
	return false
}

func (s *mlfqScheduler) OnBlock(pid int) {
	for i := range s.queues {
		s.queues[i].Remove(pid)
	}
	delete(s.runTicks, pid)
}

func (s *mlfqScheduler) OnExit(pid int) {
	s.OnBlock(pid)
	delete(s.levels, pid)
}

func (s *mlfqScheduler) Next() (int, bool) {
	for i := range s.queues {
		if pid, ok := s.queues[i].Dequeue(); ok {
			return pid, true
		}
	}
	return 0, false
}
