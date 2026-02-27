package lessons

import (
	"fmt"
)

func validate(output StageOutput, v ValidatorSpec) (bool, string) {
	switch v.Type {
	case "trace_contains_all":
		seen := map[string]bool{}
		for _, ev := range output.Trace {
			seen[ev.Kind] = true
		}
		for _, need := range v.Values {
			if !seen[need] {
				return false, fmt.Sprintf("missing trace event %s", need)
			}
		}
		return true, ""
	case "trace_order":
		if len(v.Values) < 2 {
			return false, "trace_order requires at least two event kinds"
		}
		pos := 0
		for _, ev := range output.Trace {
			if ev.Kind == v.Values[pos] {
				pos++
				if pos == len(v.Values) {
					return true, ""
				}
			}
		}
		return false, fmt.Sprintf("trace order not satisfied: %v", v.Values)
	case "trace_count_eq":
		if len(v.Values) == 0 {
			return false, "trace_count_eq requires event kind in values"
		}
		got := traceCount(output, v.Values[0])
		want := int(v.Number)
		if got != want {
			return false, fmt.Sprintf("trace count %s got=%d want=%d", v.Values[0], got, want)
		}
		return true, ""
	case "trace_count_lte":
		if len(v.Values) == 0 {
			return false, "trace_count_lte requires event kind in values"
		}
		got := traceCount(output, v.Values[0])
		want := int(v.Number)
		if got > want {
			return false, fmt.Sprintf("trace count %s got=%d want<=%d", v.Values[0], got, want)
		}
		return true, ""
	case "no_event":
		seen := map[string]bool{}
		for _, ev := range output.Trace {
			seen[ev.Kind] = true
		}
		for _, forbidden := range v.Values {
			if seen[forbidden] {
				return false, fmt.Sprintf("forbidden trace event %s present", forbidden)
			}
		}
		return true, ""
	case "metric_eq":
		got, ok := metricValue(output, v.Key)
		if !ok || got != v.Number {
			return false, fmt.Sprintf("metric %s got=%v want=%v", v.Key, got, v.Number)
		}
		return true, ""
	case "metric_gte":
		got, ok := metricValue(output, v.Key)
		if !ok || got < v.Number {
			return false, fmt.Sprintf("metric %s got=%v want>=%v", v.Key, got, v.Number)
		}
		return true, ""
	case "metric_lte":
		got, ok := metricValue(output, v.Key)
		if !ok || got > v.Number {
			return false, fmt.Sprintf("metric %s got=%v want<=%v", v.Key, got, v.Number)
		}
		return true, ""
	case "fault_eq":
		got, ok := faultValue(output, v.Key)
		if !ok || got != int(v.Number) {
			return false, fmt.Sprintf("fault %s got=%v want=%v", v.Key, got, int(v.Number))
		}
		return true, ""
	case "fault_lte":
		got, ok := faultValue(output, v.Key)
		if !ok || got > int(v.Number) {
			return false, fmt.Sprintf("fault %s got=%v want<=%v", v.Key, got, int(v.Number))
		}
		return true, ""
	case "fs_ok":
		if !output.FilesystemOK {
			return false, "filesystem invariants failed"
		}
		return true, ""
	default:
		return false, fmt.Sprintf("unknown validator type %s", v.Type)
	}
}

func metricValue(output StageOutput, key string) (float64, bool) {
	switch key {
	case "completed_processes":
		return float64(output.Metrics.CompletedProcesses), true
	case "avg_response_time":
		return output.Metrics.AvgResponseTime, true
	case "avg_turnaround_time":
		return output.Metrics.AvgTurnaroundTime, true
	case "throughput_per_100_ticks":
		return output.Metrics.ThroughputPer100Tick, true
	case "fairness_jain_index":
		return output.Metrics.FairnessJainIndex, true
	case "total_ticks":
		return float64(output.Metrics.TotalTicks), true
	default:
		return 0, false
	}
}

func traceCount(output StageOutput, kind string) int {
	total := 0
	for _, ev := range output.Trace {
		if ev.Kind == kind {
			total++
		}
	}
	return total
}

func faultValue(output StageOutput, key string) (int, bool) {
	switch key {
	case "not_present":
		return output.Memory.Faults.NotPresent, true
	case "permission":
		return output.Memory.Faults.Permission, true
	default:
		return 0, false
	}
}
