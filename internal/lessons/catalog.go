package lessons

import "os-simulator-plan/internal/sim"

type cpuLessonSpec struct {
	id    string
	title string
	seed  uint64
	p1    string
	p2    string
	steps int
}

type memoryLessonSpec struct {
	id      string
	title   string
	seed    uint64
	frames  int
	program string
	steps   int
	faults  float64
}

type ioLessonSpec struct {
	id      string
	title   string
	seed    uint64
	program string
	steps   int
}

func DefaultCatalog() map[string]Lesson {
	lessons := make([]Lesson, 0, 28)

	cpuSpecs := []cpuLessonSpec{
		{id: "l01-sched-rr-basics", title: "Round Robin Dispatch Basics", seed: 11, p1: "COMPUTE 4; EXIT", p2: "COMPUTE 4; EXIT", steps: 20},
		{id: "l02-sched-fifo-baseline", title: "FIFO Baseline and Ordering", seed: 12, p1: "COMPUTE 5; EXIT", p2: "COMPUTE 2; EXIT", steps: 20},
		{id: "l03-sched-mlfq-balance", title: "MLFQ Fairness and Balance", seed: 13, p1: "COMPUTE 6; EXIT", p2: "COMPUTE 6; EXIT", steps: 24},
		{id: "l04-response-under-rr", title: "Response Time under RR", seed: 14, p1: "COMPUTE 3; EXIT", p2: "COMPUTE 7; EXIT", steps: 22},
		{id: "l05-throughput-shared-cpu", title: "Throughput on Shared CPU", seed: 15, p1: "COMPUTE 4; EXIT", p2: "COMPUTE 4; EXIT", steps: 18},
		{id: "l06-preemption-check", title: "Preemption Behavior Check", seed: 16, p1: "COMPUTE 5; EXIT", p2: "COMPUTE 5; EXIT", steps: 24},
		{id: "l06b-lottery-tradeoffs", title: "Fairness Tradeoffs in Shared CPU", seed: 17, p1: "COMPUTE 8; EXIT", p2: "COMPUTE 2; EXIT", steps: 26},
		{id: "l06c-quantum-response-tuning", title: "Quantum Tuning and Response", seed: 18, p1: "COMPUTE 6; EXIT", p2: "COMPUTE 6; EXIT", steps: 24},
	}
	appendCPULessons(&lessons, cpuSpecs)

	memorySpecs := []memoryLessonSpec{
		{id: "l07-vm-fault-sequence", title: "Page Fault Sequence", seed: 21, frames: 2, program: "ACCESS 0x0 r; ACCESS 0x1000 r; ACCESS 0x2000 r; ACCESS 0x0 r; EXIT", steps: 12, faults: 4},
		{id: "l08-vm-pressure-repeat", title: "Frame Pressure with Repeated Access", seed: 22, frames: 2, program: "ACCESS 0x0 r; ACCESS 0x1000 r; ACCESS 0x2000 r; ACCESS 0x3000 r; EXIT", steps: 14, faults: 4},
		{id: "l09-vm-tlb-activity", title: "TLB Hit and Miss Activity", seed: 23, frames: 3, program: "ACCESS 0x0 r; ACCESS 0x1000 r; ACCESS 0x0 r; ACCESS 0x1000 r; EXIT", steps: 12, faults: 2},
		{id: "l10-vm-replacement-fifo", title: "FIFO Page Replacement", seed: 24, frames: 2, program: "ACCESS 0x0 r; ACCESS 0x1000 r; ACCESS 0x2000 r; EXIT", steps: 10, faults: 3},
		{id: "l11-vm-mixed-access", title: "Mixed Read/Write Access", seed: 25, frames: 2, program: "ACCESS 0x0 r; ACCESS 0x1000 w; ACCESS 0x2000 r; EXIT", steps: 10, faults: 3},
		{id: "l11b-vm-locality-window", title: "Locality Window and Reuse", seed: 26, frames: 3, program: "ACCESS 0x0 r; ACCESS 0x1000 r; ACCESS 0x2000 r; ACCESS 0x1000 r; ACCESS 0x0 r; EXIT", steps: 14, faults: 3},
		{id: "l11c-vm-fault-burst-diagnose", title: "Fault Burst Diagnosis", seed: 27, frames: 2, program: "ACCESS 0x0 r; ACCESS 0x1000 r; ACCESS 0x2000 r; ACCESS 0x3000 r; ACCESS 0x0 r; EXIT", steps: 16, faults: 5},
	}
	appendMemoryLessons(&lessons, memorySpecs, "l06c-quantum-response-tuning:s3")

	concurrencySpecs := []ioLessonSpec{
		{id: "l12-irq-wakeup-read", title: "Read Syscall IRQ Wakeup", seed: 31, program: "SYSCALL open /docs/readme.txt; SYSCALL read 4; COMPUTE 1; EXIT", steps: 14},
		{id: "l13-terminal-write-irq", title: "Terminal Write Interrupt Path", seed: 32, program: "SYSCALL open /docs/readme.txt; SYSCALL write 3; COMPUTE 1; EXIT", steps: 10},
		{id: "l14-sleep-wakeup", title: "Sleep and Wakeup Timing", seed: 33, program: "SYSCALL sleep 2; COMPUTE 1; EXIT", steps: 10},
		{id: "l15-mixed-blocking", title: "Mixed Blocking Workload", seed: 34, program: "SYSCALL open /docs/readme.txt; SYSCALL read 3; SYSCALL sleep 2; EXIT", steps: 16},
		{id: "l15b-irq-interleave-order", title: "IRQ Interleave Event Order", seed: 35, program: "SYSCALL open /docs/readme.txt; SYSCALL read 3; SYSCALL write 2; COMPUTE 1; EXIT", steps: 16},
		{id: "l15c-blocked-progress-guarantee", title: "Blocked Progress Guarantee", seed: 36, program: "SYSCALL sleep 2; SYSCALL read 2; COMPUTE 2; EXIT", steps: 18},
	}
	appendConcurrencyLessons(&lessons, concurrencySpecs, "l11c-vm-fault-burst-diagnose:s3")

	persistenceSpecs := []ioLessonSpec{
		{id: "l16-fs-open-traversal", title: "Filesystem Path Traversal", seed: 41, program: "SYSCALL open /docs/readme.txt; SYSCALL read 2; SYSCALL exit", steps: 12},
		{id: "l17-fs-read-blockmap", title: "Read Path Block Mapping", seed: 42, program: "SYSCALL open /docs/readme.txt; SYSCALL read 4; SYSCALL exit", steps: 14},
		{id: "l18-fs-write-blockmap", title: "Write Path Block Mapping", seed: 43, program: "SYSCALL open /docs/readme.txt; SYSCALL write 4; SYSCALL exit", steps: 14},
		{id: "l19-fs-read-write", title: "Read/Write Sequence", seed: 44, program: "SYSCALL open /docs/readme.txt; SYSCALL read 4; SYSCALL write 3; SYSCALL exit", steps: 16},
		{id: "l20-fs-invariants", title: "Filesystem Invariants", seed: 45, program: "SYSCALL open /docs/readme.txt; SYSCALL read 2; SYSCALL write 2; SYSCALL exit", steps: 16},
		{id: "l20b-fs-latency-compare", title: "Persistence Latency Comparison", seed: 46, program: "SYSCALL open /docs/readme.txt; SYSCALL read 3; SYSCALL write 3; SYSCALL exit", steps: 18},
		{id: "l20c-fs-trace-audit", title: "Filesystem Trace Audit", seed: 47, program: "SYSCALL open /docs/readme.txt; SYSCALL read 2; SYSCALL write 2; SYSCALL read 1; SYSCALL exit", steps: 18},
	}
	appendPersistenceLessons(&lessons, persistenceSpecs, "l15c-blocked-progress-guarantee:s3")

	for idx := range lessons {
		applyChallengeMetadata(&lessons[idx])
	}

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
		Title:  spec.title,
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
		Title:  spec.title,
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
		Title:  spec.title,
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
		Title:  spec.title,
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

func applyChallengeMetadata(lesson *Lesson) {
	applyLessonMetadata(lesson)

	for idx := range lesson.Stages {
		stage := &lesson.Stages[idx]
		if stage.Objective == "" {
			stage.Objective = defaultObjectiveForStage(lesson.Module, stage.ID)
		}
		if stage.Goal == "" {
			stage.Goal = stage.Objective
		}
		if stage.TheoryDetail == "" {
			stage.TheoryDetail = stage.Hints.Concept
		}
		if len(stage.Bootstrap) == 0 {
			stage.Bootstrap = defaultBootstrapCommands(stage.Commands)
		}
		if len(stage.AllowedCmds) == 0 {
			stage.AllowedCmds = defaultAllowedCommandsForStage(lesson.Module, stage.ID)
		}
		if stage.Limits.MaxSteps <= 0 {
			stage.Limits.MaxSteps = defaultMaxStepsForStage(lesson.Module, stage.Commands)
		}
		if stage.Limits.MaxPolicyChanges <= 0 {
			if allowsPolicy(stage.AllowedCmds) && stage.ID == "s3" {
				stage.Limits.MaxPolicyChanges = 3
			} else if allowsPolicy(stage.AllowedCmds) {
				stage.Limits.MaxPolicyChanges = 1
			} else {
				stage.Limits.MaxPolicyChanges = 0
			}
		}
		if stage.Limits.MaxConfigChanges <= 0 {
			if allowsConfig(stage.AllowedCmds) && stage.ID == "s3" {
				stage.Limits.MaxConfigChanges = 2
			} else {
				stage.Limits.MaxConfigChanges = 0
			}
		}
		if len(stage.ActionDescriptions) == 0 {
			stage.ActionDescriptions = defaultActionDescriptions(stage.AllowedCmds)
		}
		if len(stage.ExpectedVisualCues) == 0 {
			stage.ExpectedVisualCues = defaultExpectedVisualCues(stage.Validators)
		}
	}
}

func applyLessonMetadata(lesson *Lesson) {
	if lesson.SectionID == "" {
		lesson.SectionID = defaultSectionID(lesson.Module)
	}
	if lesson.SectionTitle == "" {
		lesson.SectionTitle = defaultSectionTitle(lesson.SectionID)
	}
	if lesson.Difficulty == "" {
		lesson.Difficulty = defaultDifficulty(lesson.ID)
	}
	if lesson.EstimatedMinutes <= 0 {
		lesson.EstimatedMinutes = 20
	}
	if len(lesson.ChapterRefs) == 0 {
		lesson.ChapterRefs = defaultChapterRefs(lesson.Module)
	}
}

func defaultSectionID(module string) string {
	switch module {
	case "cpu-virtualization", "memory":
		return "virtualization"
	case "concurrency":
		return "concurrency"
	case "persistence":
		return "persistence"
	default:
		return "core"
	}
}

func defaultSectionTitle(sectionID string) string {
	switch sectionID {
	case "virtualization":
		return "Virtualization"
	case "concurrency":
		return "Concurrency"
	case "persistence":
		return "Persistence"
	default:
		return "OSTEP Core"
	}
}

func defaultDifficulty(lessonID string) string {
	switch lessonID {
	case "l01-sched-rr-basics", "l02-sched-fifo-baseline", "l07-vm-fault-sequence", "l12-irq-wakeup-read", "l16-fs-open-traversal":
		return "foundation"
	case "l03-sched-mlfq-balance", "l10-vm-replacement-fifo", "l15-mixed-blocking", "l20-fs-invariants":
		return "advanced"
	default:
		return "intermediate"
	}
}

func defaultChapterRefs(module string) []string {
	switch module {
	case "cpu-virtualization":
		return []string{"cpu-intro", "cpu-mechanisms", "cpu-sched", "cpu-sched-mlfq"}
	case "memory":
		return []string{"vm-intro", "vm-mechanism", "vm-paging", "vm-tlbs"}
	case "concurrency":
		return []string{"threads-intro", "threads-locks", "threads-cv", "threads-sema"}
	case "persistence":
		return []string{"file-devices", "file-intro", "file-implementation", "file-journaling"}
	default:
		return nil
	}
}

func defaultActionDescriptions(allowed []string) []ActionDescription {
	if len(allowed) == 0 {
		return nil
	}
	out := make([]ActionDescription, 0, len(allowed))
	for _, cmd := range allowed {
		desc := "Run a supported challenge action."
		switch cmd {
		case "step":
			desc = "Advance the simulator by one deterministic tick."
		case "run":
			desc = "Advance multiple ticks quickly while preserving deterministic ordering."
		case "pause":
			desc = "Pause command execution and inspect current simulator state."
		case "reset":
			desc = "Reset to the stage bootstrap state and replay from the start."
		case "policy":
			desc = "Change scheduler policy or quantum (when enabled for this stage)."
		case "set_frames":
			desc = "Change available physical frames to observe page-fault pressure changes."
		case "set_tlb_entries":
			desc = "Tune TLB capacity and compare hit/miss behavior in the memory panel."
		case "set_disk_latency":
			desc = "Tune disk IRQ latency to compare block/wakeup timing in trace output."
		case "set_terminal_latency":
			desc = "Tune terminal IRQ latency and inspect how completion timing shifts."
		}
		out = append(out, ActionDescription{Command: cmd, Description: desc})
	}
	return out
}

func defaultExpectedVisualCues(validators []ValidatorSpec) []string {
	if len(validators) == 0 {
		return nil
	}
	seen := make(map[string]struct{}, len(validators))
	out := make([]string, 0, len(validators))
	for _, validator := range validators {
		cue := ""
		switch validator.Type {
		case "trace_contains_all":
			cue = "Trace timeline shows the required event sequence."
		case "metric_eq", "metric_lte":
			cue = "Metrics panel satisfies required numeric thresholds."
		case "fault_eq", "fault_lte":
			cue = "Memory panel fault counters match the expected condition."
		case "fs_ok":
			cue = "Filesystem state remains valid after command execution."
		}
		if cue == "" {
			continue
		}
		if _, ok := seen[cue]; ok {
			continue
		}
		seen[cue] = struct{}{}
		out = append(out, cue)
	}
	return out
}

func defaultAllowedCommandsForStage(module, stageID string) []string {
	switch module {
	case "cpu-virtualization":
		if stageID == "s3" {
			return []string{"step", "run", "pause", "policy", "reset"}
		}
		return []string{"step", "run", "pause", "reset"}
	case "memory":
		if stageID == "s3" {
			return []string{"step", "run", "pause", "set_frames", "set_tlb_entries", "reset"}
		}
		return []string{"step", "run", "pause", "reset"}
	case "concurrency", "persistence":
		if stageID == "s3" {
			return []string{"step", "run", "pause", "set_disk_latency", "set_terminal_latency", "reset"}
		}
		return []string{"step", "run", "pause", "reset"}
	default:
		return []string{"step", "run", "pause", "reset"}
	}
}

func defaultMaxStepsForStage(module string, commands []sim.Command) int {
	planned := 0
	for _, cmd := range commands {
		switch cmd.Name {
		case "step", "run":
			planned += cmd.Count
		}
	}
	if planned == 0 {
		switch module {
		case "memory", "concurrency", "persistence":
			return 24
		default:
			return 28
		}
	}
	if planned+4 < 8 {
		return 8
	}
	return planned + 4
}

func defaultBootstrapCommands(commands []sim.Command) []sim.Command {
	out := make([]sim.Command, 0, len(commands))
	for _, cmd := range commands {
		switch cmd.Name {
		case "spawn", "policy":
			out = append(out, cmd)
		}
	}
	if len(out) == 0 {
		return nil
	}
	return out
}

func defaultObjectiveForStage(module, stageID string) string {
	prefix := "Run the challenge"
	switch module {
	case "cpu-virtualization":
		prefix = "Run the scheduler workload"
	case "memory":
		prefix = "Run the memory workload"
	case "concurrency":
		prefix = "Run the blocking/IRQ workload"
	case "persistence":
		prefix = "Run the filesystem workload"
	}

	suffix := "and satisfy all stage checks."
	switch stageID {
	case "s1":
		suffix = "until the expected trace events appear."
	case "s2":
		suffix = "and verify the required outcome metrics."
	case "s3":
		suffix = "and apply tuning within limits before checking."
	}
	return prefix + " " + suffix
}

func allowsPolicy(allowed []string) bool {
	for _, cmd := range allowed {
		if cmd == "policy" {
			return true
		}
	}
	return false
}

func allowsConfig(allowed []string) bool {
	for _, cmd := range allowed {
		switch cmd {
		case "set_frames", "set_tlb_entries", "set_disk_latency", "set_terminal_latency":
			return true
		}
	}
	return false
}
