import { describe, expect, it } from "vitest";

import type { SnapshotDTO } from "../lib/types";
import {
  selectFrameRows,
  selectProcessMetricRows,
  selectProcessQueues,
  selectTimelineWindow,
} from "./selectors";

const snapshot: SnapshotDTO = {
  protocol_version: "v1alpha1",
  session_id: "s-01",
  tick: 9,
  trace_hash: "hash",
  trace_length: 10,
  processes: [
    { pid: 2, name: "io", state: "blocked", pc: 1, blocked_until: 15 },
    { pid: 1, name: "cpu", state: "running", pc: 2 },
    { pid: 3, name: "done", state: "terminated", pc: 4 },
    { pid: 4, name: "next", state: "ready", pc: 0 },
  ],
  metrics: {
    policy: "rr",
    total_ticks: 9,
    completed_processes: 1,
    avg_response_time: 1,
    avg_turnaround_time: 7,
    throughput_per_100_ticks: 11,
    fairness_jain_index: 1,
    processes: [
      {
        pid: 1,
        name: "cpu",
        response_time: 1,
        turnaround: 7,
        run_ticks: 4,
        wait_ticks: 2,
      },
      {
        pid: 4,
        name: "next",
        response_time: 2,
        turnaround: 0,
        run_ticks: 1,
        wait_ticks: 3,
      },
    ],
    gantt: [
      { tick: 1, pid: 1 },
      { tick: 2, pid: 1 },
      { tick: 3, pid: 0 },
    ],
  },
  memory: {
    page_size: 4096,
    total_frames: 2,
    frames: [
      { frame: 1, pid: 2, vpn: 3 },
      { frame: 0, vpn: 0 },
    ],
    tlb: [],
    faults: {
      not_present: 2,
      permission: 0,
      tlb_hit: 4,
      tlb_miss: 1,
    },
  },
};

describe("selectors", () => {
  it("selects timeline labels", () => {
    const out = selectTimelineWindow(snapshot, 2);
    expect(out).toEqual([
      { tick: 2, pid: 1, label: "P1" },
      { tick: 3, pid: 0, label: "idle" },
    ]);
  });

  it("sorts frame rows by frame index", () => {
    const out = selectFrameRows(snapshot);
    expect(out[0]).toEqual({ frame: 0, owner: "free", vpn: "0x0" });
    expect(out[1]).toEqual({ frame: 1, owner: "P2", vpn: "0x3" });
  });

  it("groups process queues by state", () => {
    const out = selectProcessQueues(snapshot);
    expect(out.running).toEqual(["P1 cpu"]);
    expect(out.ready).toEqual(["P4 next"]);
    expect(out.blocked).toEqual(["P2 io (until 15)"]);
    expect(out.terminated).toEqual(["P3 done"]);
  });

  it("maps process metrics with state and pc", () => {
    const out = selectProcessMetricRows(snapshot);
    expect(out[0]).toEqual({
      pid: 1,
      name: "cpu",
      state: "running",
      pc: 2,
      responseTime: 1,
      turnaround: 7,
      runTicks: 4,
      waitTicks: 2,
    });
  });
});
