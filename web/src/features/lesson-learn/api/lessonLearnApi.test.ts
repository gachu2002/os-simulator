import { beforeEach, describe, expect, it, vi } from "vitest";

import { fetchJSON } from "../../../lib/http";
import { fetchLessonLearn } from "./lessonLearnApi";

vi.mock("../../../lib/http", () => ({
  fetchJSON: vi.fn(),
}));

describe("fetchLessonLearn", () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it("maps learn DTO fields to camelCase", async () => {
    vi.mocked(fetchJSON).mockResolvedValue({
      version: "v3",
      section_id: "virtualization-cpu",
      lesson: {
        id: "rr-basics",
        title: "Round Robin",
        objective: "Learn time slicing",
        theory: {
          concepts: ["time slicing", "quantum", "fairness"],
        },
        challenge: {
          description: "Run processes with rr",
          actions: ["execute_instruction", "set_quantum"],
          visualizer: ["gantt", "queue"],
          parts: [
            {
              id: "A",
              title: "Part A",
              objective: "observe alternation",
              description: "Observe process alternation",
            },
          ],
        },
      },
    });

    const out = await fetchLessonLearn("http://localhost:8080", "rr-basics", "learner-1");

    expect(fetchJSON).toHaveBeenCalledWith(
      "http://localhost:8080",
      "/lessons/rr-basics/learn/v3?learner_id=learner-1",
    );
    expect(out.sectionId).toBe("virtualization-cpu");
    expect(out.stages[0].coreIdea).toBe("time slicing");
    expect(out.stages[0].mechanismSteps).toEqual(["quantum", "fairness"]);
    expect(out.stages[0].preChallengeChecklist).toEqual(["execute_instruction", "set_quantum"]);
    expect(out.stages[0].expectedVisualCues).toEqual(["gantt", "queue"]);
  });
});
