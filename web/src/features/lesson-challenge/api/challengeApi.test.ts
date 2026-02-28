import { beforeEach, describe, expect, it, vi } from "vitest";

import { fetchJSON } from "../../../lib/http";
import { startChallenge, submitChallenge } from "./challengeApi";

vi.mock("../../../lib/http", () => ({
  fetchJSON: vi.fn(),
}));

describe("challengeApi", () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it("maps startChallenge response and sends learner id", async () => {
    vi.mocked(fetchJSON).mockResolvedValueOnce({
      attempt_id: "a1",
      session_id: "s1",
      lesson_id: "rr-basics",
      stage_index: 1,
      stage_title: "Fairness",
      module: "scheduler",
      objective: "improve fairness",
      goal: "raise jain index",
      pass_conditions: ["fairness_jain_index >= 0.9"],
      allowed_commands: ["step", "policy"],
      limits: {
        max_steps: 30,
        max_policy_changes: 2,
        max_config_changes: 2,
      },
    });

    const out = await startChallenge("http://localhost:8080", "rr-basics", 1, "learner-1");

    const callArgs = vi.mocked(fetchJSON).mock.calls[0];
    expect(callArgs[0]).toBe("http://localhost:8080");
    expect(callArgs[1]).toBe("/challenges/start");
    expect(JSON.parse((callArgs[2] as { body?: unknown }).body as string)).toEqual({
      lesson_id: "rr-basics",
      stage_index: 1,
      learner_id: "learner-1",
    });
    expect(out).toMatchObject({
      attemptId: "a1",
      sessionId: "s1",
      lessonId: "rr-basics",
      stageIndex: 1,
      passConditions: ["fairness_jain_index >= 0.9"],
      limits: {
        maxSteps: 30,
        maxPolicyChanges: 2,
        maxConfigChanges: 2,
      },
    });
  });

  it("maps submitChallenge response", async () => {
    vi.mocked(fetchJSON).mockResolvedValueOnce({
      attempt_id: "a1",
      lesson_id: "rr-basics",
      stage_index: 1,
      passed: true,
      feedback_key: "ok",
      objective: "improve fairness",
      pass_conditions: ["fairness_jain_index >= 0.9"],
      output: {
        tick: 12,
        trace_hash: "abc",
        trace_length: 12,
        processes: [],
        metrics: {
          policy: "rr",
          total_ticks: 12,
          completed_processes: 1,
          avg_response_time: 2,
          avg_turnaround_time: 6,
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
          faults: {
            not_present: 0,
            permission: 0,
            tlb_hit: 1,
            tlb_miss: 1,
          },
        },
        filesystem_ok: true,
      },
      analytics: {
        total_stages: 3,
        completed_stages: 2,
        attempted_stages: 2,
        completion_rate: 66.7,
      },
      validator_results: [
        {
          name: "fairness",
          type: "metric",
          passed: true,
        },
      ],
    });

    const out = await submitChallenge("http://localhost:8080", "a1", "learner-1");

    const callArgs = vi.mocked(fetchJSON).mock.calls[0];
    expect(callArgs[1]).toBe("/challenges/submit");
    expect(JSON.parse((callArgs[2] as { body?: unknown }).body as string)).toEqual({
      attempt_id: "a1",
      learner_id: "learner-1",
    });
    expect(out.feedbackKey).toBe("ok");
    expect(out.output.traceHash).toBe("abc");
    expect(out.output.filesystemOk).toBe(true);
    expect(out.analytics.completedStages).toBe(2);
    expect(out.validatorResults?.[0].name).toBe("fairness");
  });
});
