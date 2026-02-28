import { render, screen } from "@testing-library/react";
import { describe, expect, it } from "vitest";

import type { SnapshotDTO } from "../../../../lib/types";
import { MemoryPanel } from "./MemoryPanel";
import { ProcessMetricsPanel } from "./ProcessMetricsPanel";
import { ProcessQueuesPanel } from "./ProcessQueuesPanel";
import { SchedulerTimeline } from "./SchedulerTimeline";

const snapshot: SnapshotDTO = {
  protocol_version: "v1alpha1",
  session_id: "s-01",
  tick: 12,
  trace_hash: "abc",
  trace_length: 20,
  processes: [
    { pid: 1, name: "cpu", state: "running", pc: 2 },
    { pid: 2, name: "io", state: "blocked", pc: 1, blocked_until: 18 },
    { pid: 3, name: "next", state: "ready", pc: 0 },
  ],
  metrics: {
    policy: "rr",
    total_ticks: 12,
    completed_processes: 0,
    avg_response_time: 0,
    avg_turnaround_time: 0,
    throughput_per_100_ticks: 0,
    fairness_jain_index: 1,
    processes: [
      {
        pid: 1,
        name: "cpu",
        response_time: 2,
        turnaround: 8,
        run_ticks: 5,
        wait_ticks: 1,
      },
      {
        pid: 3,
        name: "next",
        response_time: 3,
        turnaround: 0,
        run_ticks: 1,
        wait_ticks: 4,
      },
    ],
    gantt: [
      { tick: 10, pid: 1 },
      { tick: 11, pid: 0 },
      { tick: 12, pid: 2 },
    ],
  },
  memory: {
    page_size: 4096,
    total_frames: 4,
    frames: [{ frame: 0, pid: 1, vpn: 2 }],
    tlb: [{ slot: 0, pid: 1, vpn: 2, frame: 0 }],
    faults: {
      not_present: 1,
      permission: 0,
      tlb_hit: 3,
      tlb_miss: 1,
    },
  },
};

describe("visualization panels", () => {
  it("renders scheduler cells", () => {
    render(<SchedulerTimeline snapshot={snapshot} />);
    expect(screen.getByRole("list", { name: "scheduler timeline" })).toBeInTheDocument();
    expect(screen.getByText("P2")).toBeInTheDocument();
    expect(screen.getByText("idle")).toBeInTheDocument();
  });

  it("renders memory rows and counters", () => {
    render(<MemoryPanel snapshot={snapshot} />);
    expect(screen.getByText("Frames: 4")).toBeInTheDocument();
    expect(screen.getByText("P1")).toBeInTheDocument();
    expect(screen.getByText("0x2")).toBeInTheDocument();
  });

  it("renders process queues", () => {
    render(<ProcessQueuesPanel snapshot={snapshot} />);
    expect(screen.getByText("Running")).toBeInTheDocument();
    expect(screen.getByText("P1 cpu")).toBeInTheDocument();
    expect(screen.getByText("P2 io (until 18)")).toBeInTheDocument();
    expect(screen.getByText("P3 next")).toBeInTheDocument();
  });

  it("renders process metrics table", () => {
    render(<ProcessMetricsPanel snapshot={snapshot} />);
    expect(screen.getByText("Process Metrics")).toBeInTheDocument();
    expect(screen.getByText("cpu")).toBeInTheDocument();
    expect(screen.getByText("running")).toBeInTheDocument();
    expect(screen.getByText("8")).toBeInTheDocument();
  });
});
