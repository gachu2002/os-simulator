import { render, screen } from "@testing-library/react";
import { describe, expect, it } from "vitest";

import type { SnapshotDTO } from "../../../../lib/types";
import { LessonSpecificVisualization } from "./LessonSpecificVisualization";

const snapshot: SnapshotDTO = {
  protocol_version: "v1alpha1",
  session_id: "s-01",
  tick: 21,
  trace_hash: "abc",
  trace_length: 40,
  last_command: "run",
  processes: [
    { pid: 1, name: "alpha", state: "running", pc: 8 },
    { pid: 2, name: "beta", state: "ready", pc: 3 },
    { pid: 3, name: "gamma", state: "blocked", pc: 2 },
  ],
  metrics: {
    policy: "mlfq",
    quantum: 4,
    total_ticks: 21,
    completed_processes: 1,
    avg_response_time: 1.8,
    avg_turnaround_time: 5.2,
    throughput_per_100_ticks: 3.3,
    fairness_jain_index: 0.91,
    processes: [
      {
        pid: 1,
        name: "alpha",
        response_time: 1,
        turnaround: 8,
        run_ticks: 12,
        wait_ticks: 3,
      },
      {
        pid: 2,
        name: "beta",
        response_time: 2,
        turnaround: 6,
        run_ticks: 8,
        wait_ticks: 4,
      },
    ],
    gantt: [
      { tick: 20, pid: 1 },
      { tick: 21, pid: 2 },
    ],
  },
  memory: {
    page_size: 4096,
    total_frames: 8,
    frames: [{ frame: 0, pid: 1, vpn: 2 }],
    tlb: [{ slot: 0, pid: 1, vpn: 2, frame: 0 }],
    faults: {
      not_present: 0,
      permission: 0,
      tlb_hit: 4,
      tlb_miss: 1,
    },
  },
};

describe("lesson specific visualization", () => {
  it("renders MLFQ queue lens for lesson 6", () => {
    render(
      <LessonSpecificVisualization
        lessonID="l06-mlfq"
        snapshot={snapshot}
        lastLessonAction="trigger_priority_boost"
      />,
    );

    expect(screen.getByText("MLFQ Queue View")).toBeInTheDocument();
    expect(screen.getByText("Q0 (top)")).toBeInTheDocument();
    expect(screen.getByText("Fairness score: 0.91")).toBeInTheDocument();
  });

  it("renders proportional share bars for lesson 7", () => {
    render(
      <LessonSpecificVisualization
        lessonID="l07-lottery-stride"
        snapshot={snapshot}
        lastLessonAction="run_quanta"
      />,
    );

    expect(screen.getByText("Proportional Share View")).toBeInTheDocument();
    expect(screen.getByText("alpha (PID 1)")).toBeInTheDocument();
    expect(screen.getByText("60.00%")).toBeInTheDocument();
  });

  it("renders dual CPU layout for lesson 8", () => {
    render(
      <LessonSpecificVisualization
        lessonID="l08-multi-cpu-scheduling"
        snapshot={snapshot}
        lastLessonAction="migrate_job"
      />,
    );

    expect(screen.getByText("Multi-CPU Scheduler View")).toBeInTheDocument();
    expect(screen.getByText("CPU 0")).toBeInTheDocument();
    expect(screen.getByText("CPU 1")).toBeInTheDocument();
    expect(screen.getByText("Imbalance indicator: 1 | last action: migrate_job")).toBeInTheDocument();
  });
});
