package sim

import "testing"

func TestEventQueueOrdersByTickThenSequence(t *testing.T) {
	q := NewEventQueue()
	q.Push(Event{Tick: 3, Sequence: 10, Kind: "event"})
	q.Push(Event{Tick: 1, Sequence: 4, Kind: "event"})
	q.Push(Event{Tick: 1, Sequence: 2, Kind: "event"})
	q.Push(Event{Tick: 2, Sequence: 1, Kind: "event"})

	ordered := []Event{
		{Tick: 1, Sequence: 2},
		{Tick: 1, Sequence: 4},
		{Tick: 2, Sequence: 1},
		{Tick: 3, Sequence: 10},
	}

	for i := range ordered {
		ev, ok := q.Pop()
		if !ok {
			t.Fatalf("missing event at index %d", i)
		}

		if ev.Tick != ordered[i].Tick || ev.Sequence != ordered[i].Sequence {
			t.Fatalf("event at index %d = (%d,%d), want (%d,%d)", i, ev.Tick, ev.Sequence, ordered[i].Tick, ordered[i].Sequence)
		}
	}
}
