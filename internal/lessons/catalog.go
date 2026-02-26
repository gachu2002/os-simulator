package lessons

import (
	"fmt"

	"os-simulator-plan/internal/sim"
)

type cpuLessonSpec struct {
	id    string
	seed  uint64
	p1    string
	p2    string
	steps int
}

type memoryLessonSpec struct {
	id      string
	seed    uint64
	frames  int
	program string
	steps   int
	faults  float64
}

type ioLessonSpec struct {
	id      string
	seed    uint64
	program string
	steps   int
}

func DefaultCatalog() map[string]Lesson {
	lessons := make([]Lesson, 0, 20)

	cpuSpecs := []cpuLessonSpec{
		{id: "l01-sched-rr-basics", seed: 11, p1: "COMPUTE 4; EXIT", p2: "COMPUTE 4; EXIT", steps: 20},
		{id: "l02-sched-fifo-baseline", seed: 12, p1: "COMPUTE 5; EXIT", p2: "COMPUTE 2; EXIT", steps: 20},
		{id: "l03-sched-mlfq-balance", seed: 13, p1: "COMPUTE 6; EXIT", p2: "COMPUTE 6; EXIT", steps: 24},
		{id: "l04-response-under-rr", seed: 14, p1: "COMPUTE 3; EXIT", p2: "COMPUTE 7; EXIT", steps: 22},
		{id: "l05-throughput-shared-cpu", seed: 15, p1: "COMPUTE 4; EXIT", p2: "COMPUTE 4; EXIT", steps: 18},
		{id: "l06-preemption-check", seed: 16, p1: "COMPUTE 5; EXIT", p2: "COMPUTE 5; EXIT", steps: 24},
	}
	appendCPULessons(&lessons, cpuSpecs)

	memorySpecs := []memoryLessonSpec{
		{id: "l07-vm-fault-sequence", seed: 21, frames: 2, program: "ACCESS 0x0 r; ACCESS 0x1000 r; ACCESS 0x2000 r; ACCESS 0x0 r; EXIT", steps: 12, faults: 4},
		{id: "l08-vm-pressure-repeat", seed: 22, frames: 2, program: "ACCESS 0x0 r; ACCESS 0x1000 r; ACCESS 0x2000 r; ACCESS 0x3000 r; EXIT", steps: 14, faults: 4},
		{id: "l09-vm-tlb-activity", seed: 23, frames: 3, program: "ACCESS 0x0 r; ACCESS 0x1000 r; ACCESS 0x0 r; ACCESS 0x1000 r; EXIT", steps: 12, faults: 2},
		{id: "l10-vm-replacement-fifo", seed: 24, frames: 2, program: "ACCESS 0x0 r; ACCESS 0x1000 r; ACCESS 0x2000 r; EXIT", steps: 10, faults: 3},
		{id: "l11-vm-mixed-access", seed: 25, frames: 2, program: "ACCESS 0x0 r; ACCESS 0x1000 w; ACCESS 0x2000 r; EXIT", steps: 10, faults: 3},
	}
	appendMemoryLessons(&lessons, memorySpecs, "l06-preemption-check:s3")

	concurrencySpecs := []ioLessonSpec{
		{id: "l12-irq-wakeup-read", seed: 31, program: "SYSCALL open /docs/readme.txt; SYSCALL read 4; COMPUTE 1; EXIT", steps: 14},
		{id: "l13-terminal-write-irq", seed: 32, program: "SYSCALL open /docs/readme.txt; SYSCALL write 3; COMPUTE 1; EXIT", steps: 10},
		{id: "l14-sleep-wakeup", seed: 33, program: "SYSCALL sleep 2; COMPUTE 1; EXIT", steps: 10},
		{id: "l15-mixed-blocking", seed: 34, program: "SYSCALL open /docs/readme.txt; SYSCALL read 3; SYSCALL sleep 2; EXIT", steps: 16},
	}
	appendConcurrencyLessons(&lessons, concurrencySpecs, "l11-vm-mixed-access:s3")

	persistenceSpecs := []ioLessonSpec{
		{id: "l16-fs-open-traversal", seed: 41, program: "SYSCALL open /docs/readme.txt; SYSCALL read 2; SYSCALL exit", steps: 12},
		{id: "l17-fs-read-blockmap", seed: 42, program: "SYSCALL open /docs/readme.txt; SYSCALL read 4; SYSCALL exit", steps: 14},
		{id: "l18-fs-write-blockmap", seed: 43, program: "SYSCALL open /docs/readme.txt; SYSCALL write 4; SYSCALL exit", steps: 14},
		{id: "l19-fs-read-write", seed: 44, program: "SYSCALL open /docs/readme.txt; SYSCALL read 4; SYSCALL write 3; SYSCALL exit", steps: 16},
		{id: "l20-fs-invariants", seed: 45, program: "SYSCALL open /docs/readme.txt; SYSCALL read 2; SYSCALL write 2; SYSCALL exit", steps: 16},
	}
	appendPersistenceLessons(&lessons, persistenceSpecs, "l15-mixed-blocking:s3")

	out := make(map[string]Lesson, len(lessons))
	for _, lesson := range lessons {
		out[lesson.ID] = lesson
	}
	return out
}

func appendCPULessons(out *[]Lesson, specs []cpuLessonSpec) {
	prereq := ""
	for _, spec := range specs {
		*out = append(*out, cpuLesson(spec, prereq))
		prereq = spec.id + ":s3"
	}
}

func appendMemoryLessons(out *[]Lesson, specs []memoryLessonSpec, modulePrereq string) {
	prereq := modulePrereq
	for _, spec := range specs {
		*out = append(*out, memoryLesson(spec, prereq))
		prereq = spec.id + ":s3"
	}
}

func appendConcurrencyLessons(out *[]Lesson, specs []ioLessonSpec, modulePrereq string) {
	prereq := modulePrereq
	for _, spec := range specs {
		*out = append(*out, concurrencyLesson(spec, prereq))
		prereq = spec.id + ":s3"
	}
}

func appendPersistenceLessons(out *[]Lesson, specs []ioLessonSpec, modulePrereq string) {
	prereq := modulePrereq
	for _, spec := range specs {
		*out = append(*out, persistenceLesson(spec, prereq))
		prereq = spec.id + ":s3"
	}
}

func baseConfig(seed uint64) SimConfig {
	return SimConfig{Seed: seed, Policy: sim.PolicyRR, Quantum: 2, Frames: 8, TLBEntries: 4, DiskLatency: 3, TerminalLatency: 1}
}

func cpuLesson(spec cpuLessonSpec, lessonPrereq string) Lesson {
	cfg := baseConfig(spec.seed)
	if spec.id == "l02-sched-fifo-baseline" {
		cfg.Policy = sim.PolicyFIFO
		cfg.Quantum = 0
	}
	if spec.id == "l03-sched-mlfq-balance" {
		cfg.Policy = sim.PolicyMLFQ
		cfg.Quantum = 0
	}

	applyCfg := cfg
	switch cfg.Policy {
	case sim.PolicyFIFO:
		applyCfg.Policy = sim.PolicyRR
		applyCfg.Quantum = 3
	case sim.PolicyMLFQ:
		applyCfg.Policy = sim.PolicyRR
		applyCfg.Quantum = 2
	default:
		applyCfg.Quantum = 1
	}

	commands := []sim.Command{
		{Name: "spawn", Process: "p1", Program: spec.p1},
		{Name: "spawn", Process: "p2", Program: spec.p2},
		{Name: "step", Count: spec.steps},
	}

	return Lesson{
		ID:     spec.id,
		Title:  fmt.Sprintf("CPU Virtualization %s", spec.id),
		Module: "cpu-virtualization",
		Stages: []Stage{
			{
				ID:            "s1",
				Title:         "Observe scheduler behavior",
				Prerequisites: prereqList(lessonPrereq),
				Config:        cfg,
				Commands:      commands,
				Validators: []ValidatorSpec{
					{Name: "gantt", Type: "trace_contains_all", Values: []string{"proc.dispatch", "proc.compute"}},
				},
				Hints: HintSet{Nudge: "Start by reading dispatch order, not just final metrics.", Concept: "Schedulers expose fairness and response tradeoffs in trace ordering.", Explicit: "Confirm proc.dispatch and proc.compute events appear in stable order, then compare process alternation."},
			},
			{
				ID:            "s2",
				Title:         "Diagnose fairness and completion",
				Prerequisites: []string{spec.id + ":s1"},
				Config:        cfg,
				Commands:      commands,
				Validators: []ValidatorSpec{
					{Name: "completed", Type: "metric_eq", Key: "completed_processes", Number: 2},
					{Name: "gantt", Type: "trace_contains_all", Values: []string{"proc.dispatch", "proc.compute"}},
				},
				Hints: HintSet{Nudge: "Track both dispatch frequency and process completion.", Concept: "Completion confirms the schedule sustained progress for all runnable jobs.", Explicit: "Validate completed_processes == 2 and reference dispatch ordering to justify fairness."},
			},
			{
				ID:            "s3",
				Title:         "Apply policy or quantum tuning",
				Prerequisites: []string{spec.id + ":s2"},
				Config:        applyCfg,
				Commands:      commands,
				Validators: []ValidatorSpec{
					{Name: "completed", Type: "metric_eq", Key: "completed_processes", Number: 2},
					{Name: "trace", Type: "trace_contains_all", Values: []string{"proc.dispatch", "proc.compute"}},
				},
				Hints: HintSet{Nudge: "Compare this tuned run with earlier stage ordering.", Concept: "Policy and quantum changes should alter behavior without breaking deterministic progression.", Explicit: "Run with tuned config, check completion, and contrast dispatch cadence against previous stages."},
			},
		},
	}
}

func memoryLesson(spec memoryLessonSpec, lessonPrereq string) Lesson {
	cfg := baseConfig(spec.seed)
	cfg.Frames = spec.frames
	cfg.TLBEntries = spec.frames

	applyCfg := cfg
	applyCfg.Frames = cfg.Frames + 1
	applyCfg.TLBEntries = cfg.TLBEntries + 1

	commands := []sim.Command{{Name: "spawn", Process: "vm", Program: spec.program}, {Name: "step", Count: spec.steps}}

	return Lesson{
		ID:     spec.id,
		Title:  fmt.Sprintf("Memory %s", spec.id),
		Module: "memory",
		Stages: []Stage{
			{
				ID:            "s1",
				Title:         "Observe translation and fault trace",
				Prerequisites: prereqList(lessonPrereq),
				Config:        cfg,
				Commands:      commands,
				Validators: []ValidatorSpec{
					{Name: "trace", Type: "trace_contains_all", Values: []string{"mem.fault", "mem.access"}},
				},
				Hints: HintSet{Nudge: "Count unique virtual pages versus available frames.", Concept: "Faults are deterministic from frame count and access ordering.", Explicit: "Verify mem.fault and mem.access appear, then map each access to frame pressure."},
			},
			{
				ID:            "s2",
				Title:         "Diagnose fault counts",
				Prerequisites: []string{spec.id + ":s1"},
				Config:        cfg,
				Commands:      commands,
				Validators: []ValidatorSpec{
					{Name: "faults", Type: "fault_eq", Key: "not_present", Number: spec.faults},
					{Name: "trace", Type: "trace_contains_all", Values: []string{"mem.fault", "mem.access"}},
				},
				Hints: HintSet{Nudge: "Re-check repeated virtual pages versus first-touch pages.", Concept: "Repeated accesses can hit after first translation while new VPNs fault.", Explicit: "Enumerate each ACCESS step, mark fault/hit, and match not_present count exactly."},
			},
			{
				ID:            "s3",
				Title:         "Apply frame configuration change",
				Prerequisites: []string{spec.id + ":s2"},
				Config:        applyCfg,
				Commands:      commands,
				Validators: []ValidatorSpec{
					{Name: "faults-lte", Type: "fault_lte", Key: "not_present", Number: spec.faults},
					{Name: "trace", Type: "trace_contains_all", Values: []string{"mem.fault", "mem.access"}},
				},
				Hints: HintSet{Nudge: "Higher frame count should reduce or preserve fault totals.", Concept: "Configuration shifts pressure characteristics while preserving deterministic behavior.", Explicit: "Compare baseline and tuned frame runs; validate not_present is less-than-or-equal and explain why."},
			},
		},
	}
}

func concurrencyLesson(spec ioLessonSpec, lessonPrereq string) Lesson {
	cfg := baseConfig(spec.seed)
	cfg.Policy = sim.PolicyFIFO
	cfg.Quantum = 0

	applyCfg := cfg
	applyCfg.DiskLatency = 2

	commands := []sim.Command{{Name: "spawn", Process: "c1", Program: spec.program}, {Name: "step", Count: spec.steps}}

	return Lesson{
		ID:     spec.id,
		Title:  fmt.Sprintf("Concurrency %s", spec.id),
		Module: "concurrency",
		Stages: []Stage{
			{
				ID:            "s1",
				Title:         "Observe block and wakeup flow",
				Prerequisites: prereqList(lessonPrereq),
				Config:        cfg,
				Commands:      commands,
				Validators: []ValidatorSpec{
					{Name: "wakeup", Type: "trace_contains_all", Values: []string{"proc.wakeup", "trap.return"}},
				},
				Hints: HintSet{Nudge: "Find where control leaves and returns to user work.", Concept: "Async completion and sleep both re-enter ready state through deterministic wakeups.", Explicit: "Trace proc.wakeup and trap.return events to explain unblock timing."},
			},
			{
				ID:            "s2",
				Title:         "Diagnose completion under blocking",
				Prerequisites: []string{spec.id + ":s1"},
				Config:        cfg,
				Commands:      commands,
				Validators: []ValidatorSpec{
					{Name: "wakeup", Type: "trace_contains_all", Values: []string{"proc.wakeup", "trap.return"}},
					{Name: "exit", Type: "metric_eq", Key: "completed_processes", Number: 1},
				},
				Hints: HintSet{Nudge: "Use both trace and metric outputs.", Concept: "A blocked process still completes once wakeups re-enable CPU progress.", Explicit: "Validate wakeup events and completed_processes == 1, then explain the path from block to exit."},
			},
			{
				ID:            "s3",
				Title:         "Apply device latency tuning",
				Prerequisites: []string{spec.id + ":s2"},
				Config:        applyCfg,
				Commands:      commands,
				Validators: []ValidatorSpec{
					{Name: "wakeup", Type: "trace_contains_all", Values: []string{"proc.wakeup", "trap.return"}},
					{Name: "exit", Type: "metric_eq", Key: "completed_processes", Number: 1},
				},
				Hints: HintSet{Nudge: "Latency tuning changes timing, not correctness guarantees.", Concept: "Deterministic IRQ delivery preserves outcome validity while shifting when events happen.", Explicit: "Compare tuned and baseline traces, then confirm wakeup and completion invariants still pass."},
			},
		},
	}
}

func persistenceLesson(spec ioLessonSpec, lessonPrereq string) Lesson {
	cfg := baseConfig(spec.seed)
	cfg.Policy = sim.PolicyFIFO
	cfg.Quantum = 0

	applyCfg := cfg
	applyCfg.DiskLatency = 2

	commands := []sim.Command{{Name: "spawn", Process: "fs", Program: spec.program}, {Name: "step", Count: spec.steps}}

	return Lesson{
		ID:     spec.id,
		Title:  fmt.Sprintf("Persistence %s", spec.id),
		Module: "persistence",
		Stages: []Stage{
			{
				ID:            "s1",
				Title:         "Observe path traversal",
				Prerequisites: prereqList(lessonPrereq),
				Config:        cfg,
				Commands:      commands,
				Validators: []ValidatorSpec{
					{Name: "fs-path", Type: "trace_contains_all", Values: []string{"fs.path"}},
				},
				Hints: HintSet{Nudge: "Confirm the path is absolute before IO operations.", Concept: "Path traversal resolves to a stable inode chain in deterministic runs.", Explicit: "Run open/read or open/write and verify fs.path before reasoning about data access."},
			},
			{
				ID:            "s2",
				Title:         "Diagnose block mapping",
				Prerequisites: []string{spec.id + ":s1"},
				Config:        cfg,
				Commands:      commands,
				Validators: []ValidatorSpec{
					{Name: "fs-path", Type: "trace_contains_all", Values: []string{"fs.path", "fs.blockmap"}},
				},
				Hints: HintSet{Nudge: "Look for fs.blockmap after path resolution succeeds.", Concept: "Filesystem IO uses resolved inode metadata to map block ids deterministically.", Explicit: "Verify fs.path then fs.blockmap appears and describe how read/write steps trigger mapping."},
			},
			{
				ID:            "s3",
				Title:         "Apply persistence invariants check",
				Prerequisites: []string{spec.id + ":s2"},
				Config:        applyCfg,
				Commands:      commands,
				Validators: []ValidatorSpec{
					{Name: "fs-path", Type: "trace_contains_all", Values: []string{"fs.path", "fs.blockmap"}},
					{Name: "exit", Type: "metric_eq", Key: "completed_processes", Number: 1},
					{Name: "fs-ok", Type: "fs_ok"},
				},
				Hints: HintSet{Nudge: "Invariant checks should remain true despite latency tuning.", Concept: "Deterministic persistence model separates timing from correctness invariants.", Explicit: "Validate fs_ok, completed_processes, and trace events, then contrast tuned vs baseline run characteristics."},
			},
		},
	}
}

func prereqList(prereq string) []string {
	if prereq == "" {
		return nil
	}
	return []string{prereq}
}
