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
      version: "v3",
      section_id: "virtualization-cpu",
      attempt_id: "a1",
      session_id: "s1",
      lesson_id: "rr-basics",
      lesson_title: "Round Robin",
      lesson_objective: "improve fairness",
      part_id: "B",
      part_title: "Part B",
      part_objective: "raise jain index",
      allowed_commands: ["step", "policy"],
      limits: {
        max_steps: 30,
        max_policy_changes: 2,
        max_config_changes: 2,
      },
    });

    const out = await startChallenge("http://localhost:8080", "rr-basics", "B", "learner-1");

    const callArgs = vi.mocked(fetchJSON).mock.calls[0];
    expect(callArgs[0]).toBe("http://localhost:8080");
    expect(callArgs[1]).toBe("/challenges/start/v3");
    expect(JSON.parse((callArgs[2] as { body?: unknown }).body as string)).toEqual({
      lesson_id: "rr-basics",
      part_id: "B",
      learner_id: "learner-1",
    });
    expect(out).toMatchObject({
      attemptId: "a1",
      sessionId: "s1",
      lessonId: "rr-basics",
      stageIndex: 1,
      passConditions: ["raise jain index"],
      limits: {
        maxSteps: 30,
        maxPolicyChanges: 2,
        maxConfigChanges: 2,
      },
    });
  });

  it("maps submitChallenge response", async () => {
    vi.mocked(fetchJSON).mockResolvedValueOnce({
      version: "v3",
      section_id: "virtualization-cpu",
      lesson_title: "Round Robin",
      lesson_objective: "improve fairness",
      part_id: "B",
      part_objective: "raise jain index",
      attempt_id: "a1",
      lesson_id: "rr-basics",
      passed: true,
      feedback_key: "ok",
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
    expect(callArgs[1]).toBe("/challenges/submit/v3");
    expect(JSON.parse((callArgs[2] as { body?: unknown }).body as string)).toEqual({
      attempt_id: "a1",
      learner_id: "learner-1",
    });
    expect(out.feedbackKey).toBe("ok");
    expect(out.stageIndex).toBe(1);
    expect(out.objective).toBe("improve fairness");
    expect(out.goal).toBe("raise jain index");
    expect(out.output.traceHash).toBe("abc");
    expect(out.output.filesystemOk).toBe(true);
    expect(out.analytics.completedStages).toBe(2);
    expect(out.validatorResults?.[0].name).toBe("fairness");
  });

  it("does not send part_id for non-part lessons", async () => {
    vi.mocked(fetchJSON).mockResolvedValueOnce({
      version: "v3",
      section_id: "virtualization-cpu",
      attempt_id: "a2",
      session_id: "s2",
      lesson_id: "l01-process-basics",
      lesson_title: "What is a Process?",
      lesson_objective: "Understand process state",
      allowed_commands: ["step"],
      limits: { max_steps: 20 },
    });

    await startChallenge("http://localhost:8080", "l01-process-basics", "core", "learner-1");

    const callArgs = vi.mocked(fetchJSON).mock.calls[0];
    expect(callArgs[1]).toBe("/challenges/start/v3");
    expect(JSON.parse((callArgs[2] as { body?: unknown }).body as string)).toEqual({
      lesson_id: "l01-process-basics",
      learner_id: "learner-1",
    });
  });

  it("maps submitChallenge response without part metadata", async () => {
    vi.mocked(fetchJSON).mockResolvedValueOnce({
      version: "v3",
      section_id: "virtualization-cpu",
      lesson_title: "What is a Process?",
      lesson_objective: "Understand process state",
      attempt_id: "a2",
      lesson_id: "l01-process-basics",
      passed: false,
      feedback_key: "validator.trace_contains_all",
      pass_conditions: ["Trace must contain: dispatch, wakeup."],
      output: {
        tick: 4,
        trace_hash: "def",
        trace_length: 4,
        processes: [],
        metrics: {
          policy: "rr",
          total_ticks: 4,
          completed_processes: 0,
          avg_response_time: 0,
          avg_turnaround_time: 0,
          throughput_per_100_ticks: 0,
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
            tlb_hit: 0,
            tlb_miss: 0,
          },
        },
        filesystem_ok: true,
      },
      analytics: {
        total_stages: 8,
        completed_stages: 0,
        attempted_stages: 1,
        completion_rate: 0,
      },
    });

    const out = await submitChallenge("http://localhost:8080", "a2", "learner-1");

    expect(out.stageIndex).toBe(0);
    expect(out.goal).toBeUndefined();
    expect(out.objective).toBe("Understand process state");
  });
});
