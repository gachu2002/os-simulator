import { describe, expect, it } from "vitest";

import type { ChallengeGradeResponse } from "../lib/lessonApi";
import { snapshotFromChallengeGrade } from "./lessonSnapshot";

describe("snapshotFromChallengeGrade", () => {
  it("maps challenge grade output into snapshot dto", () => {
    const input: ChallengeGradeResponse = {
      attempt_id: "a-000001",
      lesson_id: "l01",
      stage_index: 1,
      passed: false,
      feedback_key: "validator.completed",
      output: {
        tick: 9,
        trace_hash: "def",
        trace_length: 21,
        processes: [],
        metrics: {
          policy: "rr",
          total_ticks: 9,
          completed_processes: 1,
          avg_response_time: 1,
          avg_turnaround_time: 4,
          throughput_per_100_ticks: 8,
          fairness_jain_index: 1,
          processes: [],
          gantt: [],
        },
        memory: {
          page_size: 4096,
          total_frames: 8,
          frames: [],
          tlb: [],
          faults: { not_present: 0, permission: 0, tlb_hit: 0, tlb_miss: 0 },
        },
        filesystem_ok: true,
      },
      analytics: {
        total_stages: 20,
        completed_stages: 1,
        attempted_stages: 2,
        completion_rate: 0.05,
      },
    };

    const out = snapshotFromChallengeGrade(input);
    expect(out.trace_hash).toBe("def");
    expect(out.tick).toBe(9);
    expect(out.session_id).toBe("challenge:a-000001");
    expect(out.last_command).toBe("challenge.grade.l01.stage.1");
  });
});
