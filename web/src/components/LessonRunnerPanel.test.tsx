import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import type { ReactElement } from "react";
import { afterEach, describe, expect, it, vi } from "vitest";

import { LessonRunnerPanel } from "./LessonRunnerPanel";

vi.mock("../lib/ws", () => {
  return {
    connectSessionSocket: vi.fn(() => ({
      sendCommand: vi.fn(),
      close: vi.fn(),
    })),
  };
});

describe("LessonRunnerPanel", () => {
  afterEach(() => {
    vi.restoreAllMocks();
  });

  it("starts and grades a challenge attempt", async () => {
    const fetchMock = vi
      .spyOn(globalThis, "fetch")
      .mockResolvedValueOnce(
        jsonResponse({
          lessons: [
            {
              id: "l01",
              title: "CPU Basics",
              module: "cpu",
              stages: [
                {
                  index: 0,
                  id: "s1",
                  title: "Observe scheduler behavior",
                },
              ],
            },
          ],
        }),
      )
      .mockResolvedValueOnce(
        jsonResponse({
          attempt_id: "a-000001",
          session_id: "s-000001",
          lesson_id: "l01",
          stage_index: 0,
          stage_title: "Observe scheduler behavior",
          module: "cpu",
          objective: "Observe scheduler behavior",
          allowed_commands: ["step", "run", "reset", "pause", "spawn", "policy"],
          limits: { max_steps: 40, max_policy_changes: 3 },
        }),
      )
      .mockResolvedValueOnce(
        jsonResponse({
          attempt_id: "a-000001",
          lesson_id: "l01",
          stage_index: 0,
          passed: false,
          feedback_key: "validator.completed",
          hint: "Try stepping until completion.",
          hint_level: 1,
          output: {
            tick: 20,
            trace_hash: "abc",
            trace_length: 30,
            processes: [],
            metrics: {
              policy: "rr",
              total_ticks: 20,
              completed_processes: 1,
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
            total_stages: 20,
            completed_stages: 1,
            attempted_stages: 1,
            completion_rate: 0.05,
          },
        }),
      );

    const user = userEvent.setup();
    renderWithQuery(<LessonRunnerPanel baseURL="http://localhost:8080" />);

    await waitFor(() => {
      expect(screen.getByText("cpu - CPU Basics")).toBeInTheDocument();
    });

    await user.click(screen.getByRole("button", { name: "Start Challenge" }));

    await waitFor(() => {
      expect(screen.getByText("objective: Observe scheduler behavior")).toBeInTheDocument();
    });

    await user.click(screen.getByRole("button", { name: "Check" }));

    await waitFor(() => {
      expect(screen.getByText("failed")).toBeInTheDocument();
      expect(screen.getByText("result: validator.completed")).toBeInTheDocument();
      expect(screen.getByText("Hint L1: Try stepping until completion.")).toBeInTheDocument();
      expect(screen.getByText("Completed steps: 1/20 (5%)")).toBeInTheDocument();
    });

    expect(fetchMock).toHaveBeenCalledTimes(3);
  });
});

function jsonResponse(payload: unknown): Response {
  return {
    ok: true,
    status: 200,
    json: async () => payload,
  } as Response;
}

function renderWithQuery(ui: ReactElement) {
  const queryClient = new QueryClient({
    defaultOptions: {
      queries: { retry: false },
      mutations: { retry: false },
    },
  });
  return render(
    <QueryClientProvider client={queryClient}>{ui}</QueryClientProvider>,
  );
}
