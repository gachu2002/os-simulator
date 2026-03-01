package realtime

import (
	"strconv"
	"strings"

	"os-simulator-plan/internal/lessons"
)

func normalizeLearnerID(value string) string {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return "anonymous"
	}
	return trimmed
}

func stagePassConditions(stage lessons.Stage) []string {
	out := make([]string, 0, len(stage.Validators))
	for _, validator := range stage.Validators {
		out = append(out, describeValidator(validator))
	}
	return out
}

func describeValidator(v lessons.ValidatorSpec) string {
	switch v.Type {
	case "trace_contains_all":
		if len(v.Values) == 0 {
			return "Required trace events must appear."
		}
		return "Trace must contain: " + strings.Join(v.Values, ", ") + "."
	case "trace_order":
		if len(v.Values) == 0 {
			return "Trace events must follow required order."
		}
		return "Trace order must include: " + strings.Join(v.Values, " -> ") + "."
	case "trace_count_eq":
		if len(v.Values) == 0 {
			return "Trace event count must equal required value."
		}
		return "Trace count for " + v.Values[0] + " must equal " + trimFloat(v.Number) + "."
	case "trace_count_lte":
		if len(v.Values) == 0 {
			return "Trace event count must be <= required value."
		}
		return "Trace count for " + v.Values[0] + " must be <= " + trimFloat(v.Number) + "."
	case "no_event":
		if len(v.Values) == 0 {
			return "Forbidden trace events must not appear."
		}
		return "Trace must not contain: " + strings.Join(v.Values, ", ") + "."
	case "metric_eq":
		return "Metric " + v.Key + " must equal " + trimFloat(v.Number) + "."
	case "metric_gte":
		return "Metric " + v.Key + " must be >= " + trimFloat(v.Number) + "."
	case "metric_lte":
		return "Metric " + v.Key + " must be <= " + trimFloat(v.Number) + "."
	case "fault_eq":
		return "Fault count " + v.Key + " must equal " + trimFloat(v.Number) + "."
	case "fault_lte":
		return "Fault count " + v.Key + " must be <= " + trimFloat(v.Number) + "."
	case "fs_ok":
		return "Filesystem invariants must hold."
	default:
		return "Check " + v.Name + " must pass."
	}
}

func trimFloat(value float64) string {
	if value == float64(int64(value)) {
		return strconv.FormatInt(int64(value), 10)
	}
	return strconv.FormatFloat(value, 'f', -1, 64)
}

func convertAnalytics(in lessons.CompletionAnalytics) CompletionAnalyticsDTO {
	return CompletionAnalyticsDTO{
		TotalStages:     in.TotalStages,
		CompletedStages: in.CompletedStages,
		AttemptedStages: in.AttemptedStages,
		CompletionRate:  in.CompletionRate,
	}
}
