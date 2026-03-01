import type { LessonLearnStage, LessonLearnSummary } from "../../../entities/lesson/model";
import { fetchJSON } from "../../../lib/http";

interface LessonLearnV3ResponseDTO {
  version: string;
  section_id: string;
  lesson: LessonV3DTO;
}

interface LessonV3DTO {
  id: string;
  title: string;
  objective: string;
  theory: {
    concepts: string[];
  };
  challenge: {
    description: string;
    actions: string[];
    visualizer: string[];
    parts?: LessonPartV3DTO[];
  };
}

interface LessonPartV3DTO {
  id: string;
  title: string;
  objective: string;
  description: string;
}

export async function fetchLessonLearn(
  baseURL: string,
  lessonID: string,
  learnerID: string,
): Promise<LessonLearnSummary> {
  const payload = await fetchJSON<LessonLearnV3ResponseDTO>(
    baseURL,
    `/lessons/${encodeURIComponent(lessonID)}/learn/v3?learner_id=${encodeURIComponent(learnerID)}`,
  );
  return mapLesson(payload.lesson, payload.section_id);
}

function mapLesson(dto: LessonV3DTO, sectionID: string): LessonLearnSummary {
  return {
    id: dto.id,
    title: dto.title,
    module: sectionID,
    sectionId: sectionID,
    sectionTitle: "Virtualization - CPU",
    stages: mapStages(dto),
  };
}

function mapStages(dto: LessonV3DTO): LessonLearnStage[] {
  const coreIdea = dto.theory.concepts[0] ?? dto.objective;
  const mechanismSteps = dto.theory.concepts.slice(1);

  if ((dto.challenge.parts ?? []).length > 0) {
    return (dto.challenge.parts ?? []).map((part, index) =>
      mapStage({
        index,
        id: part.id,
        title: part.title,
        coreIdea,
        mechanismSteps,
        workedExample: part.description,
        objective: part.objective,
        goal: dto.challenge.description,
        expectedVisualCues: dto.challenge.visualizer,
        preChallengeChecklist: dto.challenge.actions,
      }),
    );
  }

  return [
    mapStage({
      index: 0,
      id: "core",
      title: dto.title,
      coreIdea,
      mechanismSteps,
      workedExample: dto.challenge.description,
      objective: dto.objective,
      goal: dto.challenge.description,
      expectedVisualCues: dto.challenge.visualizer,
      preChallengeChecklist: dto.challenge.actions,
    }),
  ];
}

function mapStage(dto: {
  index: number;
  id: string;
  title: string;
  coreIdea: string;
  mechanismSteps: string[];
  workedExample: string;
  objective: string;
  goal: string;
  preChallengeChecklist: string[];
  expectedVisualCues: string[];
}): LessonLearnStage {
  return {
    index: dto.index,
    id: dto.id,
    title: dto.title,
    coreIdea: dto.coreIdea,
    mechanismSteps: dto.mechanismSteps,
    workedExample: dto.workedExample,
    commonMistakes: [],
    preChallengeChecklist: dto.preChallengeChecklist,
    objective: dto.objective,
    goal: dto.goal,
    prerequisites: [],
    expectedVisualCues: dto.expectedVisualCues,
  };
}
