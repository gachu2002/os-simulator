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
	case "metric_eq":
		got, ok := metricValue(output, v.Key)
		if !ok || got != v.Number {
			return false, fmt.Sprintf("metric %s got=%v want=%v", v.Key, got, v.Number)
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
	default:
		return 0, false
	}
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
