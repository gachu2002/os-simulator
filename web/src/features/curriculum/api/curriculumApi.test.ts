import { beforeEach, describe, expect, it, vi } from "vitest";

import { fetchJSON } from "../../../lib/http";
import { fetchCurriculumForLearner } from "./curriculumApi";

vi.mock("../../../lib/http", () => ({
  fetchJSON: vi.fn(),
}));

describe("fetchCurriculumForLearner", () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it("maps curriculum DTO to domain model", async () => {
    vi.mocked(fetchJSON).mockResolvedValue({
      sections: [
        {
          id: "cpu",
          title: "CPU Virtualization",
          subtitle: "Scheduling",
          order: 1,
          coming_soon: false,
          completed_stages: 2,
          total_stages: 3,
          lessons: [
            {
              id: "rr-basics",
              title: "Round Robin",
              module: "scheduler",
              section_id: "cpu",
              section_title: "CPU Virtualization",
              estimated_minutes: 20,
              chapter_refs: ["ch7"],
              stages: [
                {
                  index: 0,
                  id: "s0",
                  title: "Intro",
                  pass_conditions: ["throughput >= 1"],
                  allowed_commands: ["step"],
                  action_descriptions: [{ command: "step", description: "advance" }],
                  expected_visual_cues: ["ready queue drains"],
                  limits: {
                    max_steps: 20,
                    max_policy_changes: 1,
                    max_config_changes: 2,
                  },
                  completed: true,
                  unlocked: true,
                },
              ],
            },
          ],
        },
      ],
    });

    const out = await fetchCurriculumForLearner("http://localhost:8080", "learner-1");

    expect(fetchJSON).toHaveBeenCalledWith(
      "http://localhost:8080",
      "/curriculum?learner_id=learner-1",
    );
    expect(out[0].comingSoon).toBe(false);
    expect(out[0].completedStages).toBe(2);
    expect(out[0].lessons?.[0].estimatedMinutes).toBe(20);
    expect(out[0].lessons?.[0].stages[0].passConditions).toEqual(["throughput >= 1"]);
    expect(out[0].lessons?.[0].stages[0].limits?.maxPolicyChanges).toBe(1);
  });
});
