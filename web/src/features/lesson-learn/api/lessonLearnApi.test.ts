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
      lesson: {
        id: "rr-basics",
        title: "Round Robin",
        module: "scheduler",
        section_id: "cpu",
        section_title: "CPU",
        estimated_minutes: 15,
        chapter_refs: ["ch7"],
        stages: [
          {
            index: 0,
            id: "s0",
            title: "Core",
            core_idea: "time slicing",
            mechanism_steps: ["enqueue", "rotate"],
            worked_example: "P1 then P2",
            common_mistakes: ["starvation assumptions"],
            pre_challenge_checklist: ["watch queue"],
            expected_visual_cues: ["gantt alternates"],
          },
        ],
      },
    });

    const out = await fetchLessonLearn("http://localhost:8080", "rr-basics", "learner-1");

    expect(fetchJSON).toHaveBeenCalledWith(
      "http://localhost:8080",
      "/lessons/rr-basics/learn?learner_id=learner-1",
    );
    expect(out.estimatedMinutes).toBe(15);
    expect(out.chapterRefs).toEqual(["ch7"]);
    expect(out.stages[0].coreIdea).toBe("time slicing");
    expect(out.stages[0].mechanismSteps).toEqual(["enqueue", "rotate"]);
    expect(out.stages[0].preChallengeChecklist).toEqual(["watch queue"]);
    expect(out.stages[0].expectedVisualCues).toEqual(["gantt alternates"]);
  });
});
