package lessons

import (
	"slices"

	"os-simulator-plan/internal/sim"
)

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
	content, err := loadDefaultCatalogContent()
	if err != nil {
		panic(err)
	}
	lessonContent, err := loadLessonContentV2()
	if err != nil {
		panic(err)
	}
	stageContent, err := loadLessonStageContentV2()
	if err != nil {
		panic(err)
	}

	lessons := make([]Lesson, 0, len(content.CPU)+len(content.Memory.Lessons)+len(content.Concurrency.Lessons)+len(content.Persistence.Lessons))

	appendCPULessons(&lessons, content.CPU)
	appendMemoryLessons(&lessons, content.Memory.Lessons, content.Memory.ModulePrerequisite)
	appendConcurrencyLessons(&lessons, content.Concurrency.Lessons, content.Concurrency.ModulePrerequisite)
	appendPersistenceLessons(&lessons, content.Persistence.Lessons, content.Persistence.ModulePrerequisite)

	for idx := range lessons {
		record, ok := lessonContent[lessons[idx].ID]
		if !ok {
			panic("missing lesson content for " + lessons[idx].ID)
		}
		overrides, ok := stageContent[lessons[idx].ID]
		if !ok {
			panic("missing stage content for " + lessons[idx].ID)
		}
		applyChallengeMetadata(&lessons[idx], record, overrides)
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

func applyChallengeMetadata(
	lesson *Lesson,
	lessonContent lessonContentRecord,
	stageContent map[string]lessonStageContentStage,
) {
	applyLessonMetadata(lesson, lessonContent)

	for idx := range lesson.Stages {
		stage := &lesson.Stages[idx]
		applyStageV2LearnContent(stage, lessonContent, stageContent)
		stage.Bootstrap = defaultBootstrapCommands(stage.Commands)
		stage.AllowedCmds = defaultAllowedCommandsForStage(lesson.Module, stage.ID)
		stage.Limits.MaxSteps = defaultMaxStepsForStage(lesson.Module, stage.Commands)
		if allowsPolicy(stage.AllowedCmds) && stage.ID == "s3" {
			stage.Limits.MaxPolicyChanges = 3
		} else if allowsPolicy(stage.AllowedCmds) {
			stage.Limits.MaxPolicyChanges = 1
		} else {
			stage.Limits.MaxPolicyChanges = 0
		}
		if allowsConfig(stage.AllowedCmds) && stage.ID == "s3" {
			stage.Limits.MaxConfigChanges = 2
		} else {
			stage.Limits.MaxConfigChanges = 0
		}
		stage.ActionDescriptions = defaultActionDescriptions(stage.AllowedCmds)
		stage.ExpectedVisualCues = defaultExpectedVisualCues(stage.Validators)
		if len(stage.ValidatorHints) == 0 {
			stage.ValidatorHints = defaultValidatorHints(lesson.Module, stage.Validators, stage.Hints)
		}
	}
}

func applyStageV2LearnContent(
	stage *Stage,
	lessonContent lessonContentRecord,
	stageContent map[string]lessonStageContentStage,
) {
	stage.CoreIdea = lessonContent.Learn.CoreIdea
	stage.MechanismSteps = slices.Clone(lessonContent.Learn.MechanismSteps)
	stage.WorkedExample = lessonContent.Learn.WorkedExample
	stage.CommonMistakes = slices.Clone(lessonContent.Learn.CommonMistakes)
	stage.PreChallengeChecklist = slices.Clone(lessonContent.Learn.PreChallengeChecklist)

	override, ok := stageContent[stage.ID]
	if !ok {
		panic("missing stage override for " + stage.ID)
	}
	if override.Objective == "" {
		panic("missing stage objective for " + stage.ID)
	}
	if override.Goal == "" {
		panic("missing stage goal for " + stage.ID)
	}
	if override.Hints.Nudge == "" || override.Hints.Concept == "" || override.Hints.Explicit == "" {
		panic("missing complete stage hints for " + stage.ID)
	}

	stage.Objective = override.Objective
	stage.Goal = override.Goal
	if override.TheoryDetail != "" {
		stage.TheoryDetail = override.TheoryDetail
	}
	stage.Hints = override.Hints
	if hints := toValidatorHints(override.ValidatorHints); len(hints) > 0 {
		stage.ValidatorHints = hints
	}
}

func applyLessonMetadata(lesson *Lesson, record lessonContentRecord) {
	if lesson.SectionID == "" {
		lesson.SectionID = defaultSectionID(lesson.Module)
	}
	if lesson.SectionTitle == "" {
		lesson.SectionTitle = defaultSectionTitle(lesson.SectionID)
	}
	lesson.Difficulty = record.Difficulty
	lesson.EstimatedMinutes = record.EstimatedMinutes
	lesson.ChapterRefs = slices.Clone(record.ChapterRefs)
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
		case "trace_order", "trace_count_eq", "trace_count_lte", "no_event":
			cue = "Trace ordering and event frequency match challenge constraints."
		case "metric_eq", "metric_lte", "metric_gte":
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

func defaultValidatorHints(module string, validators []ValidatorSpec, fallback HintSet) []ValidatorHint {
	out := make([]ValidatorHint, 0, len(validators))
	for _, validator := range validators {
		hints := fallback
		switch validator.Type {
		case "trace_contains_all":
			hints = HintSet{
				Nudge:    "Compare your trace timeline against the required events first.",
				Concept:  "Trace validators check mechanism evidence, not just final outcomes.",
				Explicit: "Run enough steps, then verify each required event appears before submitting.",
			}
		case "trace_order":
			hints = HintSet{
				Nudge:    "Check whether required events appear in the expected order.",
				Concept:  "Order failures usually indicate misunderstanding of control transfer flow.",
				Explicit: "Inspect event sequence and re-run until ordering matches the expected chain.",
			}
		case "trace_count_eq", "trace_count_lte":
			hints = HintSet{
				Nudge:    "Count relevant trace events before you submit.",
				Concept:  "Event frequency reflects policy and workload interaction.",
				Explicit: "Adjust steps or tuning actions, then re-check the exact event count target.",
			}
		case "no_event":
			hints = HintSet{
				Nudge:    "Look for forbidden events in the trace.",
				Concept:  "Some events signal invalid behavior for this lesson goal.",
				Explicit: "Reset and avoid the action path that triggers forbidden events.",
			}
		case "metric_eq", "metric_lte", "metric_gte":
			hints = HintSet{
				Nudge:    "Focus on the target metric value in the metrics panel.",
				Concept:  "Metrics capture performance tradeoffs; one change can improve one metric and hurt another.",
				Explicit: "Tune policy/config and run again until the metric threshold is satisfied.",
			}
		case "fault_eq", "fault_lte":
			hints = HintSet{
				Nudge:    "Track memory access order and fault counters together.",
				Concept:  "Fault outcomes are determined by access pattern, frame pressure, and translation reuse.",
				Explicit: "Enumerate accesses and compare expected vs observed not-present faults.",
			}
		case "fs_ok":
			hints = HintSet{
				Nudge:    "Check filesystem validity after each run.",
				Concept:  "Correct persistence behavior preserves invariants regardless of latency changes.",
				Explicit: "Reset, rerun carefully, and verify path/blockmap flow before final invariants check.",
			}
		}

		if module == "concurrency" && (validator.Type == "trace_order" || validator.Type == "trace_contains_all") {
			hints.Concept = "In concurrency lessons, trace ordering explains block, wakeup, and resume behavior."
		}

		out = append(out, ValidatorHint{Validator: validator.Name, Hints: hints})
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
