package realtime

import (
	"fmt"

	"os-simulator-plan/internal/lessons"
	"os-simulator-plan/internal/sim"
)

type qualityGate struct {
	Name     string
	Expected string
	Actual   string
	Passed   bool
	Message  string
}

func evaluateV3QualityGates(lessonID, partID string, output lessons.StageOutput) []qualityGate {
	switch lessonID {
	case "l01-process-basics":
		return []qualityGate{
			metricGate("state_transitions", ">= 1 block and >= 1 unblock", fmt.Sprintf("block=%d unblock=%d", traceCount(output.Trace, "proc.block"), traceCount(output.Trace, "proc.unblock")), traceCount(output.Trace, "proc.block") >= 1 && traceCount(output.Trace, "proc.unblock") >= 1, "show blocked-to-ready lifecycle explicitly"),
			metricGate("process_lifecycle_end", ">= 1 proc.kill or proc.exit", fmt.Sprintf("kill=%d exit=%d", traceCount(output.Trace, "proc.kill"), traceCount(output.Trace, "proc.exit")), traceCount(output.Trace, "proc.kill")+traceCount(output.Trace, "proc.exit") >= 1, "demonstrate full process lifecycle including termination"),
		}
	case "l02-process-api-fork-exec-wait":
		return []qualityGate{
			metricGate("family_growth", ">= 2 process creations", fmt.Sprintf("proc.spawn=%d", traceCount(output.Trace, "proc.spawn")), traceCount(output.Trace, "proc.spawn") >= 2, "fork/exec challenge should produce parent-child execution flow"),
			metricGate("child_completion", ">= 1 process exit", fmt.Sprintf("proc.exit=%d completed=%d", traceCount(output.Trace, "proc.exit"), output.Metrics.CompletedProcesses), traceCount(output.Trace, "proc.exit") >= 1 && output.Metrics.CompletedProcesses >= 1, "complete and reap child lifecycle before submit"),
		}
	case "l03-limited-direct-execution":
		if partID == "A" {
			return []qualityGate{
				metricGate("trap_round_trip", ">=1 trap.enter, >=1 sys.dispatch, >=1 trap.return", fmt.Sprintf("enter=%d dispatch=%d return=%d", traceCount(output.Trace, "trap.enter"), traceCount(output.Trace, "sys.dispatch"), traceCount(output.Trace, "trap.return")), traceCount(output.Trace, "trap.enter") >= 1 && traceCount(output.Trace, "sys.dispatch") >= 1 && traceCount(output.Trace, "trap.return") >= 1, "execute full syscall trap lifecycle"),
			}
		}
		return []qualityGate{
			metricGate("forced_switch", ">=1 preempt and >=2 dispatch events", fmt.Sprintf("preempt=%d dispatch=%d", traceCount(output.Trace, "proc.preempt")+traceCount(output.Trace, "proc.preempt.manual"), traceCount(output.Trace, "proc.dispatch")), traceCount(output.Trace, "proc.preempt")+traceCount(output.Trace, "proc.preempt.manual") >= 1 && traceCount(output.Trace, "proc.dispatch") >= 2, "show forced control return and context switch"),
			metricGate("manual_selection", ">=1 choose_next decision", fmt.Sprintf("choose_next=%d", traceCount(output.Trace, "proc.choose_next")), traceCount(output.Trace, "proc.choose_next") >= 1, "explicitly choose next process at least once"),
		}
	case "l04-cpu-scheduling-basics":
		return []qualityGate{
			metricGate("policy_comparison", ">= 2 policy events (initial + at least one change)", fmt.Sprintf("sched.policy=%d", traceCount(output.Trace, "sched.policy")), traceCount(output.Trace, "sched.policy") >= 2, "switch scheduling policy at least once to compare behavior"),
			metricGate("completion_evidence", ">= 2 completed processes", fmt.Sprintf("completed=%d", output.Metrics.CompletedProcesses), output.Metrics.CompletedProcesses >= 2, "run long enough to compare policy outcomes on more than one job"),
			metricGate("timeline_evidence", ">= 12 total ticks", fmt.Sprintf("ticks=%d", output.Metrics.TotalTicks), output.Metrics.TotalTicks >= 12, "collect enough timeline evidence before submitting"),
		}
	case "l05-round-robin":
		return []qualityGate{
			metricGate("rr_tuning", "policy=rr and >= 2 policy events", fmt.Sprintf("policy=%s sched.policy=%d", output.Metrics.Policy, traceCount(output.Trace, "sched.policy")), output.Metrics.Policy == sim.PolicyRR && traceCount(output.Trace, "sched.policy") >= 2, "set and test at least one new RR quantum before submit"),
			metricGate("interactive_window", ">= 16 total ticks", fmt.Sprintf("ticks=%d", output.Metrics.TotalTicks), output.Metrics.TotalTicks >= 16, "run enough quanta to expose response/turnaround tension"),
			metricGate("throughput_floor", "throughput_per_100_ticks >= 2", fmt.Sprintf("throughput=%.2f", output.Metrics.ThroughputPer100Tick), output.Metrics.ThroughputPer100Tick >= 2, "quantum choice currently hurts throughput too much"),
		}
	case "l06-mlfq":
		if partID == "A" {
			return []qualityGate{
				metricGate("mlfq_activity", ">= 1 preempt event", fmt.Sprintf("preempt=%d", traceCount(output.Trace, "proc.preempt")+traceCount(output.Trace, "proc.preempt.manual")), traceCount(output.Trace, "proc.preempt")+traceCount(output.Trace, "proc.preempt.manual") >= 1, "step long enough to observe queue/quantum pressure"),
				metricGate("mlfq_progress", ">= 20 total ticks", fmt.Sprintf("ticks=%d", output.Metrics.TotalTicks), output.Metrics.TotalTicks >= 20, "collect enough queue evolution to reason about demotion behavior"),
			}
		}
		return []qualityGate{
			metricGate("fairness_window", "fairness_jain_index >= 0.70", fmt.Sprintf("fairness=%.2f", output.Metrics.FairnessJainIndex), output.Metrics.FairnessJainIndex >= 0.70, "tune run/step sequence until fairness improves"),
			metricGate("mlfq_runtime", ">= 24 total ticks", fmt.Sprintf("ticks=%d", output.Metrics.TotalTicks), output.Metrics.TotalTicks >= 24, "run longer before evaluating gaming/fairness effects"),
		}
	case "l07-lottery-stride":
		return []qualityGate{
			metricGate("share_window", ">= 24 total ticks", fmt.Sprintf("ticks=%d", output.Metrics.TotalTicks), output.Metrics.TotalTicks >= 24, "run enough quanta to measure proportional-share behavior"),
			metricGate("share_signal", "fairness_jain_index >= 0.75", fmt.Sprintf("fairness=%.2f", output.Metrics.FairnessJainIndex), output.Metrics.FairnessJainIndex >= 0.75, "distribution is still too uneven; adjust run and workload"),
			metricGate("workload_depth", ">= 3 tracked processes", fmt.Sprintf("process_metrics=%d", len(output.Metrics.Processes)), len(output.Metrics.Processes) >= 3, "use at least three jobs to validate target shares"),
		}
	case "l08-multi-cpu-scheduling":
		return []qualityGate{
			metricGate("dispatch_depth", ">= 4 dispatch events", fmt.Sprintf("dispatch=%d", traceCount(output.Trace, "proc.dispatch")), traceCount(output.Trace, "proc.dispatch") >= 4, "step further to observe multi-CPU scheduling pressure"),
			metricGate("load_balance_signal", "fairness_jain_index >= 0.70", fmt.Sprintf("fairness=%.2f", output.Metrics.FairnessJainIndex), output.Metrics.FairnessJainIndex >= 0.70, "current run appears imbalanced; keep balancing workload"),
			metricGate("utilization_window", ">= 20 total ticks", fmt.Sprintf("ticks=%d", output.Metrics.TotalTicks), output.Metrics.TotalTicks >= 20, "collect sufficient utilization timeline before submit"),
		}
	default:
		return nil
	}
}

func qualityGatesToDTO(gates []qualityGate) []ValidatorResultDTO {
	out := make([]ValidatorResultDTO, 0, len(gates))
	for _, gate := range gates {
		out = append(out, ValidatorResultDTO{
			Name:     gate.Name,
			Type:     "v3_quality",
			Passed:   gate.Passed,
			Message:  gate.Message,
			Expected: gate.Expected,
			Actual:   gate.Actual,
		})
	}
	return out
}

func allQualityGatesPassed(gates []qualityGate) bool {
	for _, gate := range gates {
		if !gate.Passed {
			return false
		}
	}
	return true
}

func qualityGatePassConditions(gates []qualityGate) []string {
	out := make([]string, 0, len(gates))
	for _, gate := range gates {
		out = append(out, fmt.Sprintf("%s: %s", gate.Name, gate.Expected))
	}
	return out
}

func firstQualityGateHint(gates []qualityGate) string {
	for _, gate := range gates {
		if !gate.Passed {
			if gate.Message != "" {
				return gate.Message
			}
			return "complete required challenge-quality checks before submitting"
		}
	}
	return ""
}

func metricGate(name, expected, actual string, passed bool, message string) qualityGate {
	return qualityGate{Name: name, Expected: expected, Actual: actual, Passed: passed, Message: message}
}

func traceCount(trace []sim.TraceEvent, kind string) int {
	total := 0
	for _, ev := range trace {
		if ev.Kind == kind {
			total++
		}
	}
	return total
}
