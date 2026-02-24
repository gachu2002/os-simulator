package sim

import "testing"

func runTwoProcWorkload(t *testing.T, policy string, quantum int) SchedulingMetrics {
	t.Helper()
	e := NewEngine(123, 0)
	if err := e.SetSchedulingPolicy(policy, quantum); err != nil {
		t.Fatalf("set policy failed: %v", err)
	}
	commands := []Command{
		{Name: "spawn", Process: "p1", Program: "COMPUTE 4; EXIT"},
		{Name: "spawn", Process: "p2", Program: "COMPUTE 4; EXIT"},
		{Name: "step", Count: 20},
	}
	if err := e.ExecuteAll(commands); err != nil {
		t.Fatalf("execute failed: %v", err)
	}
	return e.SchedulingMetrics()
}

func TestKnownWorkloadMetrics_FIFOvsRR(t *testing.T) {
	fifo := runTwoProcWorkload(t, PolicyFIFO, 0)
	rr := runTwoProcWorkload(t, PolicyRR, 2)

	if fifo.CompletedProcesses != 2 || rr.CompletedProcesses != 2 {
		t.Fatalf("expected both workloads to complete two processes")
	}

	if fifo.AvgResponseTime <= rr.AvgResponseTime {
		t.Fatalf("expected fifo avg response (%f) > rr avg response (%f)", fifo.AvgResponseTime, rr.AvgResponseTime)
	}

	if len(fifo.Gantt) < 10 || len(rr.Gantt) < 10 {
		t.Fatalf("expected gantt slices for both policies")
	}

	if fifo.Gantt[0].PID != 1 || fifo.Gantt[1].PID != 1 || fifo.Gantt[2].PID != 1 || fifo.Gantt[3].PID != 1 {
		t.Fatalf("fifo should run p1 continuously at start")
	}

	if rr.Gantt[0].PID != 1 || rr.Gantt[1].PID != 1 || rr.Gantt[2].PID != 2 || rr.Gantt[3].PID != 2 {
		t.Fatalf("rr should alternate after quantum")
	}
}

func TestKnownWorkloadMetrics_MLFQ(t *testing.T) {
	mlfq := runTwoProcWorkload(t, PolicyMLFQ, 0)
	if mlfq.CompletedProcesses != 2 {
		t.Fatalf("expected mlfq to complete both processes")
	}
	if mlfq.FairnessJainIndex <= 0 {
		t.Fatalf("expected positive fairness index")
	}
	if len(mlfq.Processes) != 2 {
		t.Fatalf("expected two process metrics")
	}
}
