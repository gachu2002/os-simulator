package lessons

import (
	"testing"

	"os-simulator-plan/internal/sim"
)

func TestValidateSupportsExtendedTraceValidators(t *testing.T) {
	output := StageOutput{
		Trace: []sim.TraceEvent{
			{Kind: "trap.enter"},
			{Kind: "proc.dispatch"},
			{Kind: "proc.compute"},
			{Kind: "trap.return"},
			{Kind: "proc.compute"},
		},
	}

	cases := []struct {
		name string
		spec ValidatorSpec
		want bool
	}{
		{name: "trace order pass", spec: ValidatorSpec{Name: "ordered", Type: "trace_order", Values: []string{"trap.enter", "proc.dispatch", "trap.return"}}, want: true},
		{name: "trace order fail", spec: ValidatorSpec{Name: "ordered-fail", Type: "trace_order", Values: []string{"trap.return", "trap.enter"}}, want: false},
		{name: "trace count eq pass", spec: ValidatorSpec{Name: "count", Type: "trace_count_eq", Values: []string{"proc.compute"}, Number: 2}, want: true},
		{name: "trace count lte pass", spec: ValidatorSpec{Name: "count-lte", Type: "trace_count_lte", Values: []string{"proc.compute"}, Number: 3}, want: true},
		{name: "no event pass", spec: ValidatorSpec{Name: "no-panic", Type: "no_event", Values: []string{"kernel.panic"}}, want: true},
		{name: "no event fail", spec: ValidatorSpec{Name: "no-dispatch", Type: "no_event", Values: []string{"proc.dispatch"}}, want: false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, _ := validate(output, tc.spec)
			if got != tc.want {
				t.Fatalf("validate(%s)=%v want=%v", tc.spec.Type, got, tc.want)
			}
		})
	}
}

func TestValidateSupportsMetricGTEAndExtendedMetrics(t *testing.T) {
	output := StageOutput{
		Metrics: sim.SchedulingMetrics{
			CompletedProcesses:   2,
			AvgResponseTime:      2,
			AvgTurnaroundTime:    7,
			ThroughputPer100Tick: 11.5,
			FairnessJainIndex:    0.95,
			TotalTicks:           18,
		},
	}

	cases := []struct {
		name string
		spec ValidatorSpec
		want bool
	}{
		{name: "metric gte pass", spec: ValidatorSpec{Name: "throughput", Type: "metric_gte", Key: "throughput_per_100_ticks", Number: 10}, want: true},
		{name: "metric gte fail", spec: ValidatorSpec{Name: "fairness", Type: "metric_gte", Key: "fairness_jain_index", Number: 0.99}, want: false},
		{name: "metric lte pass total ticks", spec: ValidatorSpec{Name: "ticks", Type: "metric_lte", Key: "total_ticks", Number: 20}, want: true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, _ := validate(output, tc.spec)
			if got != tc.want {
				t.Fatalf("validate(%s)=%v want=%v", tc.spec.Name, got, tc.want)
			}
		})
	}
}
