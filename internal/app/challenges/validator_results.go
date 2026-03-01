package challenges

import (
	"strconv"
	"strings"

	"os-simulator-plan/internal/lessons"
)

type ValidatorResultView struct {
	Name     string
	Type     string
	Key      string
	Passed   bool
	Message  string
	Expected string
	Actual   string
}

func BuildValidatorResults(in []lessons.ValidationResult, stage lessons.Stage, output lessons.StageOutput) []ValidatorResultView {
	validatorByName := make(map[string]lessons.ValidatorSpec, len(stage.Validators))
	for _, spec := range stage.Validators {
		validatorByName[spec.Name] = spec
	}

	out := make([]ValidatorResultView, 0, len(in))
	for _, item := range in {
		expected, actual := expectedActual(item, validatorByName[item.Name], output)
		out = append(out, ValidatorResultView{
			Name:     item.Name,
			Type:     item.Type,
			Key:      item.Key,
			Passed:   item.Passed,
			Message:  item.Message,
			Expected: expected,
			Actual:   actual,
		})
	}
	return out
}

func expectedActual(result lessons.ValidationResult, spec lessons.ValidatorSpec, output lessons.StageOutput) (string, string) {
	switch result.Type {
	case "trace_contains_all":
		expected := "contains all: " + stringsJoinOrAny(spec.Values)
		if result.Passed {
			return expected, "all required events present"
		}
		if result.Message != "" {
			return expected, result.Message
		}
		return expected, "missing one or more required events"
	case "trace_order":
		expected := "ordered sequence: " + stringsJoinWithArrow(spec.Values)
		if result.Passed {
			return expected, "order satisfied"
		}
		if result.Message != "" {
			return expected, result.Message
		}
		return expected, "order not satisfied"
	case "trace_count_eq":
		kind := firstOr(spec.Values, "unknown")
		expected := strconv.Itoa(int(spec.Number))
		actual := strconv.Itoa(traceCount(output, kind))
		return expected, actual
	case "trace_count_lte":
		kind := firstOr(spec.Values, "unknown")
		expected := "<= " + strconv.Itoa(int(spec.Number))
		actual := strconv.Itoa(traceCount(output, kind))
		return expected, actual
	case "no_event":
		expected := "none of: " + stringsJoinOrAny(spec.Values)
		if result.Passed {
			return expected, "none present"
		}
		if result.Message != "" {
			return expected, result.Message
		}
		return expected, "forbidden event present"
	case "metric_eq":
		expected := trimFloat(spec.Number)
		got, ok := metricValue(output, spec.Key)
		if !ok {
			return expected, "n/a"
		}
		return expected, trimFloat(got)
	case "metric_gte":
		expected := ">= " + trimFloat(spec.Number)
		got, ok := metricValue(output, spec.Key)
		if !ok {
			return expected, "n/a"
		}
		return expected, trimFloat(got)
	case "metric_lte":
		expected := "<= " + trimFloat(spec.Number)
		got, ok := metricValue(output, spec.Key)
		if !ok {
			return expected, "n/a"
		}
		return expected, trimFloat(got)
	case "fault_eq":
		expected := strconv.Itoa(int(spec.Number))
		got, ok := faultValue(output, spec.Key)
		if !ok {
			return expected, "n/a"
		}
		return expected, strconv.Itoa(got)
	case "fault_lte":
		expected := "<= " + strconv.Itoa(int(spec.Number))
		got, ok := faultValue(output, spec.Key)
		if !ok {
			return expected, "n/a"
		}
		return expected, strconv.Itoa(got)
	case "fs_ok":
		if output.FilesystemOK {
			return "true", "true"
		}
		return "true", "false"
	default:
		if result.Passed {
			return "pass", "pass"
		}
		return "pass", "fail"
	}
}

func stringsJoinOrAny(values []string) string {
	if len(values) == 0 {
		return "any"
	}
	return strings.Join(values, ", ")
}

func stringsJoinWithArrow(values []string) string {
	if len(values) == 0 {
		return "any"
	}
	return strings.Join(values, " -> ")
}

func firstOr(values []string, fallback string) string {
	if len(values) == 0 {
		return fallback
	}
	return values[0]
}

func metricValue(output lessons.StageOutput, key string) (float64, bool) {
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

func traceCount(output lessons.StageOutput, kind string) int {
	total := 0
	for _, ev := range output.Trace {
		if ev.Kind == kind {
			total++
		}
	}
	return total
}

func faultValue(output lessons.StageOutput, key string) (int, bool) {
	switch key {
	case "not_present":
		return output.Memory.Faults.NotPresent, true
	case "permission":
		return output.Memory.Faults.Permission, true
	default:
		return 0, false
	}
}

func trimFloat(value float64) string {
	if value == float64(int64(value)) {
		return strconv.FormatInt(int64(value), 10)
	}
	return strconv.FormatFloat(value, 'f', -1, 64)
}
