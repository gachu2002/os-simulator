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
      version: "v3",
      sections: [
        {
          id: "virtualization-cpu",
          title: "CPU Virtualization",
          subtitle: "Scheduling",
          order: 1,
          lessons: [
            {
              id: "rr-basics",
              order: 1,
              title: "Round Robin",
              objective: "learn rr",
              challenge: {
                description: "run rr",
                actions: ["execute_instruction", "set_quantum"],
                visualizer: ["gantt-chart"],
                parts: [
                  {
                    id: "A",
                    title: "Part A",
                    objective: "phase a",
                    description: "part a desc",
                  },
                ],
              },
            },
            {
              id: "mlfq",
              order: 2,
              title: "MLFQ",
              objective: "learn mlfq",
              challenge: {
                description: "run mlfq",
                actions: ["step"],
                visualizer: ["queue"],
              },
            },
          ],
        },
      ],
    });

    const out = await fetchCurriculumForLearner("http://localhost:8080", "learner-1");

    expect(fetchJSON).toHaveBeenCalledWith(
      "http://localhost:8080",
      "/curriculum/v3?learner_id=learner-1",
    );
    expect(out[0].comingSoon).toBe(false);
    expect(out[0].completedStages).toBe(0);
    expect(out[0].lessons?.[0].sectionId).toBe("virtualization-cpu");
    expect(out[0].lessons?.[0].stages[0].id).toBe("A");
    expect(out[0].lessons?.[1].stages[0].id).toBe("core");
  });
});
