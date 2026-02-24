package lessons

import (
	"fmt"

	"os-simulator-plan/internal/sim"
)

func DefaultCatalog() map[string]Lesson {
	lessons := []Lesson{}

	lessons = append(lessons,
		cpuLesson("l01-sched-rr-basics", 11, "COMPUTE 4; EXIT", "COMPUTE 4; EXIT", 20),
		cpuLesson("l02-sched-fifo-baseline", 12, "COMPUTE 5; EXIT", "COMPUTE 2; EXIT", 20),
		cpuLesson("l03-sched-mlfq-balance", 13, "COMPUTE 6; EXIT", "COMPUTE 6; EXIT", 24),
		cpuLesson("l04-response-under-rr", 14, "COMPUTE 3; EXIT", "COMPUTE 7; EXIT", 22),
		cpuLesson("l05-throughput-shared-cpu", 15, "COMPUTE 4; EXIT", "COMPUTE 4; EXIT", 18),
		cpuLesson("l06-preemption-check", 16, "COMPUTE 5; EXIT", "COMPUTE 5; EXIT", 24),
	)

	lessons = append(lessons,
		memoryLesson("l07-vm-fault-sequence", 21, 2, "ACCESS 0x0 r; ACCESS 0x1000 r; ACCESS 0x2000 r; ACCESS 0x0 r; EXIT", 12, 4),
		memoryLesson("l08-vm-pressure-repeat", 22, 2, "ACCESS 0x0 r; ACCESS 0x1000 r; ACCESS 0x2000 r; ACCESS 0x3000 r; EXIT", 14, 4),
		memoryLesson("l09-vm-tlb-activity", 23, 3, "ACCESS 0x0 r; ACCESS 0x1000 r; ACCESS 0x0 r; ACCESS 0x1000 r; EXIT", 12, 2),
		memoryLesson("l10-vm-replacement-fifo", 24, 2, "ACCESS 0x0 r; ACCESS 0x1000 r; ACCESS 0x2000 r; EXIT", 10, 3),
		memoryLesson("l11-vm-mixed-access", 25, 2, "ACCESS 0x0 r; ACCESS 0x1000 w; ACCESS 0x2000 r; EXIT", 10, 3),
	)

	lessons = append(lessons,
		concurrencyLesson("l12-irq-wakeup-read", 31, "SYSCALL open /docs/readme.txt; SYSCALL read 4; COMPUTE 1; EXIT", 14),
		concurrencyLesson("l13-terminal-write-irq", 32, "SYSCALL open /docs/readme.txt; SYSCALL write 3; COMPUTE 1; EXIT", 10),
		concurrencyLesson("l14-sleep-wakeup", 33, "SYSCALL sleep 2; COMPUTE 1; EXIT", 10),
		concurrencyLesson("l15-mixed-blocking", 34, "SYSCALL open /docs/readme.txt; SYSCALL read 3; SYSCALL sleep 2; EXIT", 16),
	)

	lessons = append(lessons,
		persistenceLesson("l16-fs-open-traversal", 41, "SYSCALL open /docs/readme.txt; SYSCALL read 2; SYSCALL exit", 12),
		persistenceLesson("l17-fs-read-blockmap", 42, "SYSCALL open /docs/readme.txt; SYSCALL read 4; SYSCALL exit", 14),
		persistenceLesson("l18-fs-write-blockmap", 43, "SYSCALL open /docs/readme.txt; SYSCALL write 4; SYSCALL exit", 14),
		persistenceLesson("l19-fs-read-write", 44, "SYSCALL open /docs/readme.txt; SYSCALL read 4; SYSCALL write 3; SYSCALL exit", 16),
		persistenceLesson("l20-fs-invariants", 45, "SYSCALL open /docs/readme.txt; SYSCALL read 2; SYSCALL write 2; SYSCALL exit", 16),
	)

	out := make(map[string]Lesson, len(lessons))
	for _, lesson := range lessons {
		out[lesson.ID] = lesson
	}
	return out
}

func baseConfig(seed uint64) SimConfig {
	return SimConfig{Seed: seed, Policy: sim.PolicyRR, Quantum: 2, Frames: 8, TLBEntries: 4, DiskLatency: 3, TerminalLatency: 1}
}

func cpuLesson(id string, seed uint64, p1, p2 string, steps int) Lesson {
	cfg := baseConfig(seed)
	if id == "l02-sched-fifo-baseline" {
		cfg.Policy = sim.PolicyFIFO
		cfg.Quantum = 0
	}
	if id == "l03-sched-mlfq-balance" {
		cfg.Policy = sim.PolicyMLFQ
		cfg.Quantum = 0
	}
	return Lesson{
		ID:     id,
		Title:  fmt.Sprintf("CPU Virtualization %s", id),
		Module: "cpu-virtualization",
		Stages: []Stage{{
			ID:     "s1",
			Title:  "Observe scheduler behavior",
			Config: cfg,
			Commands: []sim.Command{
				{Name: "spawn", Process: "p1", Program: p1},
				{Name: "spawn", Process: "p2", Program: p2},
				{Name: "step", Count: steps},
			},
			Validators: []ValidatorSpec{
				{Name: "completed", Type: "metric_eq", Key: "completed_processes", Number: 2},
				{Name: "gantt", Type: "trace_contains_all", Values: []string{"proc.dispatch", "proc.compute"}},
			},
			Hints: HintSet{Nudge: "Compare how each policy interleaves runnable jobs.", Concept: "CPU virtualization lessons focus on fairness and response tradeoffs.", Explicit: "Run two CPU-bound jobs and inspect dispatch/compute ordering and metrics."},
		}},
	}
}

func memoryLesson(id string, seed uint64, frames int, program string, steps int, faults float64) Lesson {
	cfg := baseConfig(seed)
	cfg.Frames = frames
	cfg.TLBEntries = frames
	return Lesson{
		ID:     id,
		Title:  fmt.Sprintf("Memory %s", id),
		Module: "memory",
		Stages: []Stage{{
			ID:     "s1",
			Title:  "Trigger deterministic translation/faults",
			Config: cfg,
			Commands: []sim.Command{
				{Name: "spawn", Process: "vm", Program: program},
				{Name: "step", Count: steps},
			},
			Validators: []ValidatorSpec{
				{Name: "faults", Type: "fault_eq", Key: "not_present", Number: faults},
				{Name: "trace", Type: "trace_contains_all", Values: []string{"mem.fault", "mem.access"}},
			},
			Hints: HintSet{Nudge: "Use more virtual pages than available frames.", Concept: "Fault behavior is deterministic with fixed frame count and access order.", Explicit: "Execute ACCESS instructions across multiple VPNs and verify fault counts."},
		}},
	}
}

func concurrencyLesson(id string, seed uint64, program string, steps int) Lesson {
	cfg := baseConfig(seed)
	cfg.Policy = sim.PolicyFIFO
	cfg.Quantum = 0
	return Lesson{
		ID:     id,
		Title:  fmt.Sprintf("Concurrency %s", id),
		Module: "concurrency",
		Stages: []Stage{{
			ID:       "s1",
			Title:    "Follow block and wakeup flow",
			Config:   cfg,
			Commands: []sim.Command{{Name: "spawn", Process: "c1", Program: program}, {Name: "step", Count: steps}},
			Validators: []ValidatorSpec{
				{Name: "wakeup", Type: "trace_contains_all", Values: []string{"proc.wakeup", "trap.return"}},
				{Name: "exit", Type: "metric_eq", Key: "completed_processes", Number: 1},
			},
			Hints: HintSet{Nudge: "Watch what causes blocked processes to re-enter ready state.", Concept: "Asynchronous completion and sleep both resume work through deterministic wakeups.", Explicit: "Run a blocking syscall path and verify wakeup plus eventual process exit."},
		}},
	}
}

func persistenceLesson(id string, seed uint64, program string, steps int) Lesson {
	cfg := baseConfig(seed)
	cfg.Policy = sim.PolicyFIFO
	cfg.Quantum = 0
	return Lesson{
		ID:     id,
		Title:  fmt.Sprintf("Persistence %s", id),
		Module: "persistence",
		Stages: []Stage{{
			ID:       "s1",
			Title:    "Resolve path and map file blocks",
			Config:   cfg,
			Commands: []sim.Command{{Name: "spawn", Process: "fs", Program: program}, {Name: "step", Count: steps}},
			Validators: []ValidatorSpec{
				{Name: "fs-trace", Type: "trace_contains_all", Values: []string{"fs.path", "fs.blockmap"}},
				{Name: "fs-ok", Type: "fs_ok"},
			},
			Hints: HintSet{Nudge: "Open an absolute path before issuing file IO.", Concept: "Path traversal picks an inode and IO maps to deterministic block IDs.", Explicit: "Run open/read/write on /docs/readme.txt and verify fs.path and fs.blockmap traces."},
		}},
	}
}
