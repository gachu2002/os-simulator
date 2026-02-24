import { describe, expect, it } from "vitest";

import {
  initialSessionState,
  sessionReducer,
  type SessionState,
} from "./sessionReducer";

function seededState(): SessionState {
  return {
    ...initialSessionState,
    sessionID: "s-000001",
  };
}

describe("sessionReducer", () => {
  it("updates snapshot and log on ordered snapshot event", () => {
    const state = seededState();
    const next = sessionReducer(state, {
      type: "event.received",
      event: {
        type: "session.snapshot",
        sequence: 2,
        session_id: "s-000001",
        snapshot: {
          protocol_version: "v1alpha1",
          session_id: "s-000001",
          tick: 5,
          trace_hash: "abc",
          trace_length: 8,
          processes: [],
          metrics: {
            policy: "rr",
            total_ticks: 5,
            completed_processes: 0,
            avg_response_time: 0,
            avg_turnaround_time: 0,
            throughput_per_100_ticks: 0,
            fairness_jain_index: 0,
            processes: [],
            gantt: [],
          },
          memory: {
            page_size: 4096,
            total_frames: 8,
            frames: [],
            tlb: [],
            faults: {
              not_present: 0,
              permission: 0,
              tlb_hit: 0,
              tlb_miss: 0,
            },
          },
          last_command: "step",
        },
      },
    });

    expect(next.lastSequence).toBe(2);
    expect(next.snapshot?.tick).toBe(5);
    expect(next.logs).toHaveLength(1);
    expect(next.logs[0].detail).toBe("step");
  });

  it("ignores out-of-order events for deterministic state updates", () => {
    const state: SessionState = {
      ...seededState(),
      lastSequence: 4,
    };

    const next = sessionReducer(state, {
      type: "event.received",
      event: {
        type: "session.snapshot",
        sequence: 3,
        session_id: "s-000001",
      },
    });

    expect(next).toEqual(state);
  });
});
