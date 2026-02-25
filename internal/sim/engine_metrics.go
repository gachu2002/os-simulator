package sim

import "sort"

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

		procMetrics = append(procMetrics, ProcessMetric{
			PID:          p.PID,
			Name:         p.Name,
			ResponseTime: resp,
			Turnaround:   turn,
			RunTicks:     st.runTicks,
			WaitTicks:    st.waitTicks,
		})

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
