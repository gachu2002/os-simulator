package sim

import "container/heap"

type eventHeap []Event

func (h eventHeap) Len() int { return len(h) }

func (h eventHeap) Less(i, j int) bool {
	if h[i].Tick != h[j].Tick {
		return h[i].Tick < h[j].Tick
	}

	return h[i].Sequence < h[j].Sequence
}

func (h eventHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *eventHeap) Push(x any) {
	*h = append(*h, x.(Event))
}

func (h *eventHeap) Pop() any {
	old := *h
	n := len(old)
	v := old[n-1]
	*h = old[:n-1]
	return v
}

type EventQueue struct {
	h eventHeap
}

func NewEventQueue() *EventQueue {
	q := &EventQueue{}
	heap.Init(&q.h)
	return q
}

func (q *EventQueue) Push(e Event) {
	heap.Push(&q.h, e)
}

func (q *EventQueue) Pop() (Event, bool) {
	if len(q.h) == 0 {
		return Event{}, false
	}

	v := heap.Pop(&q.h).(Event)
	return v, true
}

func (q *EventQueue) Peek() (Event, bool) {
	if len(q.h) == 0 {
		return Event{}, false
	}

	return q.h[0], true
}

func (q *EventQueue) Len() int {
	return len(q.h)
}
