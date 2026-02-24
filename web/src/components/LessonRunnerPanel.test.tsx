import { render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { afterEach, describe, expect, it, vi } from "vitest";

import { LessonRunnerPanel } from "./LessonRunnerPanel";

describe("LessonRunnerPanel", () => {
  afterEach(() => {
    vi.restoreAllMocks();
  });

  it("loads lessons and runs selected stage", async () => {
    const fetchMock = vi
      .spyOn(globalThis, "fetch")
      .mockResolvedValueOnce(
        jsonResponse({
          lessons: [
            {
              id: "l01",
              title: "CPU Basics",
              module: "cpu",
              stages: [{ index: 0, id: "s1", title: "Observe scheduler behavior" }],
            },
          ],
        }),
      )
      .mockResolvedValueOnce(
        jsonResponse({
          lesson_id: "l01",
          stage_index: 0,
          passed: true,
          feedback_key: "stage.s1.passed",
          output: {
            tick: 20,
            trace_hash: "abc",
            trace_length: 30,
            processes: [],
            metrics: {
              policy: "rr",
              total_ticks: 20,
              completed_processes: 2,
              avg_response_time: 1,
              avg_turnaround_time: 4,
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
            module_breakdown: [
              {
                module: "cpu",
                total_stages: 5,
                completed_stage: 1,
                completion_rate: 0.2,
              },
            ],
            pilot_checklist: [],
            pilot_checklist_ok: false,
          },
        }),
      );

    const user = userEvent.setup();
    render(<LessonRunnerPanel baseURL="http://localhost:8080" />);

    await waitFor(() => {
      expect(screen.getByText("cpu - CPU Basics")).toBeInTheDocument();
    });

    await user.click(screen.getByRole("button", { name: "Run Stage" }));

    await waitFor(() => {
      expect(screen.getByText("passed")).toBeInTheDocument();
      expect(screen.getByText("feedback: stage.s1.passed")).toBeInTheDocument();
      expect(screen.getByText("1/20 (5%)")).toBeInTheDocument();
    });

    expect(fetchMock).toHaveBeenCalledTimes(2);
  });
});

function jsonResponse(payload: unknown): Response {
  return {
    ok: true,
    status: 200,
    json: async () => payload,
  } as Response;
}
