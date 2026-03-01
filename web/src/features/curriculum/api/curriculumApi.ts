import type { CurriculumSection, LessonSummary, LessonStageSummary } from "../../../entities/lesson/model";
import { fetchJSON } from "../../../lib/http";

interface CurriculumV3ResponseDTO {
  version: string;
  sections: CurriculumV3SectionDTO[];
}

interface CurriculumV3SectionDTO {
  id: string;
  title: string;
  subtitle?: string;
  order: number;
  lessons: LessonV3DTO[];
}

interface LessonV3DTO {
  id: string;
  title: string;
  order: number;
  objective: string;
  challenge: LessonChallengeV3DTO;
}

interface LessonChallengeV3DTO {
  description: string;
  actions: string[];
  visualizer: string[];
  parts?: Array<{
    id: string;
    title: string;
    objective: string;
    description: string;
  }>;
}

export async function fetchCurriculumForLearner(
  baseURL: string,
  learnerID: string,
): Promise<CurriculumSection[]> {
  const payload = await fetchJSON<CurriculumV3ResponseDTO>(
    baseURL,
    `/curriculum/v3?learner_id=${encodeURIComponent(learnerID)}`,
  );
  return payload.sections.map(mapSection);
}

function mapSection(dto: CurriculumV3SectionDTO): CurriculumSection {
  const lessons = dto.lessons
    .slice()
    .sort((a, b) => a.order - b.order)
    .map((item) => mapLessonSummary(item, dto.id, dto.title));
  const totalStages = lessons.reduce((sum, lesson) => sum + lesson.stages.length, 0);

  return {
    id: dto.id,
    title: dto.title,
    subtitle: dto.subtitle,
    order: dto.order,
    comingSoon: false,
    lessons,
    completedStages: 0,
    totalStages,
  };
}

function mapLessonSummary(dto: LessonV3DTO, sectionID: string, sectionTitle: string): LessonSummary {
  const stages = mapStagesFromV3Lesson(dto);
  return {
    id: dto.id,
    title: dto.title,
    module: sectionID,
    sectionId: sectionID,
    sectionTitle,
    stages,
  };
}

function mapStagesFromV3Lesson(dto: LessonV3DTO): LessonStageSummary[] {
  if ((dto.challenge.parts ?? []).length > 0) {
    return (dto.challenge.parts ?? []).map((part, index) => ({
      index,
      id: part.id,
      title: part.title,
      objective: part.objective,
      goal: part.description,
      passConditions: [dto.objective],
      unlocked: true,
      completed: false,
      actionDescriptions: dto.challenge.actions.map((action) => ({
        command: action,
        description: `Action: ${action}`,
      })),
      expectedVisualCues: dto.challenge.visualizer,
    }));
  }

  return [
    {
      index: 0,
      id: "core",
      title: "Core Challenge",
      objective: dto.objective,
      goal: dto.challenge.description,
      passConditions: [dto.objective],
      unlocked: true,
      completed: false,
      actionDescriptions: dto.challenge.actions.map((action) => ({
        command: action,
        description: `Action: ${action}`,
      })),
      expectedVisualCues: dto.challenge.visualizer,
    },
  ];
}
