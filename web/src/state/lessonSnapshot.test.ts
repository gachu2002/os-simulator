import { describe, expect, it } from "vitest";

import type { LessonRunResponse } from "../lib/lessonApi";
import { snapshotFromLessonRun } from "./lessonSnapshot";

describe("snapshotFromLessonRun", () => {
  it("maps lesson run output into snapshot dto", () => {
    const input: LessonRunResponse = {
      lesson_id: "l01",
      stage_index: 0,
      passed: true,
      feedback_key: "stage.s1.passed",
      output: {
        tick: 12,
        trace_hash: "abc",
        trace_length: 30,
        processes: [],
        metrics: {
          policy: "rr",
          total_ticks: 12,
          completed_processes: 2,
          avg_response_time: 1,
          avg_turnaround_time: 5,
          throughput_per_100_ticks: 10,
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
        attempted_stages: 1,
        completion_rate: 0.05,
        attempt_coverage: 0.05,
        module_breakdown: [],
        weak_concepts: [],
        pilot_checklist: [],
        pilot_checklist_ok: false,
      },
    };

    const out = snapshotFromLessonRun(input);
    expect(out.trace_hash).toBe("abc");
    expect(out.tick).toBe(12);
    expect(out.session_id).toBe("lesson:l01");
    expect(out.last_command).toBe("lesson.run.l01.stage.0");
  });
});
