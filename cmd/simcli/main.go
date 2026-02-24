package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"runtime/pprof"
	"runtime/trace"

	"os-simulator-plan/internal/lessons"
	"os-simulator-plan/internal/sim"
)

func main() {
	seed := flag.Uint64("seed", 1, "seed for deterministic bootstrap")
	steps := flag.Int("steps", 25, "number of simulation steps")
	checkpointEvery := flag.Uint64("checkpoint-every", 5, "snapshot checkpoint interval in ticks")
	policy := flag.String("policy", sim.PolicyRR, "scheduler policy: fifo|rr|mlfq")
	quantum := flag.Int("quantum", 2, "rr quantum (ticks)")
	frames := flag.Int("frames", 8, "total physical frames")
	tlbEntries := flag.Int("tlb-entries", 4, "tlb entry count")
	diskLatency := flag.Uint64("disk-latency", 3, "disk device latency in ticks")
	terminalLatency := flag.Uint64("terminal-latency", 1, "terminal device latency in ticks")
	comparePolicies := flag.Bool("compare-policies", false, "run same scenario with fifo, rr, and mlfq")
	logFile := flag.String("log-file", "", "optional replay log output path")
	replayFile := flag.String("replay-file", "", "optional replay log input path")
	processName := flag.String("process", "", "optional process name for pseudo-program run")
	program := flag.String("program", "", "pseudo-program, e.g. 'COMPUTE 3; BLOCK 2; EXIT'")
	showProcessTable := flag.Bool("show-process-table", true, "print final process table as JSON")
	showMetrics := flag.Bool("show-metrics", true, "print scheduling metrics as JSON")
	showMemory := flag.Bool("show-memory", true, "print memory snapshot as JSON")
	checkFS := flag.Bool("check-fs", true, "validate filesystem invariants")
	lessonID := flag.String("lesson-id", "", "run lesson stage by lesson id")
	lessonStage := flag.Int("lesson-stage", 0, "lesson stage index")
	runLessonPack := flag.Bool("run-lesson-pack", false, "run all lesson stages and print completion analytics")
	cpuProfile := flag.String("cpu-profile", "", "write CPU profile to file")
	traceFile := flag.String("trace-file", "", "write runtime trace to file")
	emitObservability := flag.Bool("emit-observability", false, "print structured observability event")
	flag.Parse()

	stopProfile, err := startProfiling(*cpuProfile, *traceFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "simcli profiling setup failed: %v\n", err)
		os.Exit(1)
	}
	defer stopProfile()

	if *runLessonPack {
		if err := runLessonPackAnalytics(); err != nil {
			fmt.Fprintf(os.Stderr, "simcli lesson pack failed: %v\n", err)
			os.Exit(1)
		}
		return
	}

	if *lessonID != "" {
		if err := runLesson(*lessonID, *lessonStage); err != nil {
			fmt.Fprintf(os.Stderr, "simcli lesson failed: %v\n", err)
			os.Exit(1)
		}
		return
	}

	if *replayFile != "" {
		if err := runReplay(*replayFile, sim.Tick(*checkpointEvery)); err != nil {
			fmt.Fprintf(os.Stderr, "simcli replay failed: %v\n", err)
			os.Exit(1)
		}
		return
	}

	if *comparePolicies {
		if err := runPolicyComparison(*seed, *steps, sim.Tick(*checkpointEvery), *processName, *program, *frames, *tlbEntries, sim.Tick(*diskLatency), sim.Tick(*terminalLatency)); err != nil {
			fmt.Fprintf(os.Stderr, "simcli compare failed: %v\n", err)
			os.Exit(1)
		}
		return
	}

	if err := runFresh(*seed, *steps, sim.Tick(*checkpointEvery), *policy, *quantum, *frames, *tlbEntries, sim.Tick(*diskLatency), sim.Tick(*terminalLatency), *logFile, *processName, *program, *showProcessTable, *showMetrics, *showMemory, *checkFS, *emitObservability); err != nil {
		fmt.Fprintf(os.Stderr, "simcli failed: %v\n", err)
		os.Exit(1)
	}
}

func runFresh(seed uint64, steps int, checkpointEvery sim.Tick, policy string, quantum int, frames int, tlbEntries int, diskLatency, terminalLatency sim.Tick, logFile, processName, program string, showProcessTable, showMetrics, showMemory, checkFS, emitObservability bool) error {
	engine := sim.NewEngine(seed, checkpointEvery)
	engine.ConfigureMemory(frames, tlbEntries)
	engine.ConfigureDevices(diskLatency, terminalLatency)
	if err := engine.SetSchedulingPolicy(policy, quantum); err != nil {
		return err
	}
	commands := make([]sim.Command, 0, 2)
	if program != "" {
		name := processName
		if name == "" {
			name = "demo"
		}
		commands = append(commands, sim.Command{Name: "spawn", Process: name, Program: program})
	}
	commands = append(commands, sim.Command{Name: "step", Count: steps})

	log, err := engine.ReplayLog(commands)
	if err != nil {
		return err
	}

	fmt.Printf("mode=fresh seed=%d steps=%d hash=%s trace_events=%d\n", seed, steps, log.TraceHash, len(log.Trace))
	if showProcessTable {
		table, err := json.Marshal(engine.ProcessTable())
		if err != nil {
			return err
		}
		fmt.Printf("process_table=%s\n", string(table))
	}
	if showMetrics {
		metrics, err := json.Marshal(engine.SchedulingMetrics())
		if err != nil {
			return err
		}
		fmt.Printf("metrics=%s\n", string(metrics))
	}
	if showMemory {
		memory, err := json.Marshal(engine.MemoryView())
		if err != nil {
			return err
		}
		fmt.Printf("memory=%s\n", string(memory))
	}
	if checkFS {
		if err := engine.ValidateFilesystem(); err != nil {
			return err
		}
		fmt.Printf("filesystem=ok\n")
	}
	if emitObservability {
		o := map[string]any{
			"event":               "sim.run.summary",
			"seed":                seed,
			"policy":              policy,
			"quantum":             quantum,
			"trace_hash":          log.TraceHash,
			"trace_events":        len(log.Trace),
			"completed_processes": engine.SchedulingMetrics().CompletedProcesses,
			"fault_not_present":   engine.MemoryView().Faults.NotPresent,
			"fault_permission":    engine.MemoryView().Faults.Permission,
		}
		b, err := json.Marshal(o)
		if err != nil {
			return err
		}
		fmt.Printf("observability=%s\n", string(b))
	}

	if logFile != "" {
		if err := sim.WriteReplayLog(logFile, log); err != nil {
			return err
		}
		fmt.Printf("wrote replay log to %s\n", logFile)
	}

	return nil
}

func startProfiling(cpuProfilePath, tracePath string) (func(), error) {
	var cpuFile *os.File
	var traceOut *os.File

	if cpuProfilePath != "" {
		f, err := os.Create(cpuProfilePath)
		if err != nil {
			return nil, err
		}
		if err := pprof.StartCPUProfile(f); err != nil {
			_ = f.Close()
			return nil, err
		}
		cpuFile = f
	}

	if tracePath != "" {
		f, err := os.Create(tracePath)
		if err != nil {
			if cpuFile != nil {
				pprof.StopCPUProfile()
				_ = cpuFile.Close()
			}
			return nil, err
		}
		if err := trace.Start(f); err != nil {
			_ = f.Close()
			if cpuFile != nil {
				pprof.StopCPUProfile()
				_ = cpuFile.Close()
			}
			return nil, err
		}
		traceOut = f
	}

	return func() {
		if traceOut != nil {
			trace.Stop()
			_ = traceOut.Close()
		}
		if cpuFile != nil {
			pprof.StopCPUProfile()
			_ = cpuFile.Close()
		}
	}, nil
}

func runPolicyComparison(seed uint64, steps int, checkpointEvery sim.Tick, processName, program string, frames int, tlbEntries int, diskLatency, terminalLatency sim.Tick) error {
	if program == "" {
		return fmt.Errorf("-compare-policies requires -program")
	}
	name := processName
	if name == "" {
		name = "demo"
	}
	for _, policy := range []string{sim.PolicyFIFO, sim.PolicyRR, sim.PolicyMLFQ} {
		engine := sim.NewEngine(seed, checkpointEvery)
		engine.ConfigureMemory(frames, tlbEntries)
		engine.ConfigureDevices(diskLatency, terminalLatency)
		if err := engine.SetSchedulingPolicy(policy, 2); err != nil {
			return err
		}
		commands := []sim.Command{
			{Name: "spawn", Process: name, Program: program},
			{Name: "spawn", Process: name + "-2", Program: program},
			{Name: "step", Count: steps},
		}
		if err := engine.ExecuteAll(commands); err != nil {
			return err
		}
		metrics, err := json.Marshal(engine.SchedulingMetrics())
		if err != nil {
			return err
		}
		fmt.Printf("comparison=%s\n", string(metrics))
	}
	return nil
}

func runReplay(replayFile string, checkpointEvery sim.Tick) error {
	log, err := sim.ReadReplayLog(replayFile)
	if err != nil {
		return err
	}

	replayed, err := sim.ReplayFromLog(log, checkpointEvery)
	if err != nil {
		return err
	}

	if replayed.TraceHash != log.TraceHash {
		return fmt.Errorf("trace hash mismatch: replay=%s log=%s", replayed.TraceHash, log.TraceHash)
	}

	fmt.Printf("mode=replay seed=%d hash=%s status=ok\n", replayed.Seed, replayed.TraceHash)
	return nil
}

func runLesson(lessonID string, stageIndex int) error {
	eng := lessons.NewEngine()
	res, err := eng.RunStage(lessonID, stageIndex)
	if err != nil {
		return err
	}
	out, err := json.Marshal(res)
	if err != nil {
		return err
	}
	fmt.Printf("lesson=%s\n", string(out))
	return nil
}

func runLessonPackAnalytics() error {
	eng := lessons.NewEngine()
	for _, lesson := range eng.Lessons() {
		for i := range lesson.Stages {
			if _, err := eng.RunStage(lesson.ID, i); err != nil {
				return err
			}
		}
	}
	analytics, err := json.Marshal(eng.CompletionAnalytics())
	if err != nil {
		return err
	}
	fmt.Printf("lesson_pack=%s\n", string(analytics))
	return nil
}
