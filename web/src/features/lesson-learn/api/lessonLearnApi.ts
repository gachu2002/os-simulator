import type { LessonLearnStage, LessonLearnSummary } from "../../../entities/lesson/model";
import { fetchJSON } from "../../../lib/http";

interface LessonLearnResponseDTO {
  lesson: LessonLearnSummaryDTO;
}

interface LessonLearnSummaryDTO {
  id: string;
  title: string;
  module: string;
  section_id?: string;
  section_title?: string;
  difficulty?: string;
  estimated_minutes?: number;
  chapter_refs?: string[];
  stages: LessonLearnStageDTO[];
}

interface LessonLearnStageDTO {
  index: number;
  id: string;
  title: string;
  core_idea?: string;
  mechanism_steps?: string[];
  worked_example?: string;
  common_mistakes?: string[];
  pre_challenge_checklist?: string[];
  objective?: string;
  goal?: string;
  prerequisites?: string[];
  expected_visual_cues?: string[];
}

export async function fetchLessonLearn(
  baseURL: string,
  lessonID: string,
  learnerID: string,
): Promise<LessonLearnSummary> {
  const payload = await fetchJSON<LessonLearnResponseDTO>(
    baseURL,
    `/lessons/${encodeURIComponent(lessonID)}/learn?learner_id=${encodeURIComponent(learnerID)}`,
  );
  return mapLesson(payload.lesson);
}

function mapLesson(dto: LessonLearnSummaryDTO): LessonLearnSummary {
  return {
    id: dto.id,
    title: dto.title,
    module: dto.module,
    sectionId: dto.section_id,
    sectionTitle: dto.section_title,
    difficulty: dto.difficulty,
    estimatedMinutes: dto.estimated_minutes,
    chapterRefs: dto.chapter_refs,
    stages: dto.stages.map(mapStage),
  };
}

function mapStage(dto: LessonLearnStageDTO): LessonLearnStage {
  return {
    index: dto.index,
    id: dto.id,
    title: dto.title,
    coreIdea: dto.core_idea,
    mechanismSteps: dto.mechanism_steps,
    workedExample: dto.worked_example,
    commonMistakes: dto.common_mistakes,
    preChallengeChecklist: dto.pre_challenge_checklist,
    objective: dto.objective,
    goal: dto.goal,
    prerequisites: dto.prerequisites,
    expectedVisualCues: dto.expected_visual_cues,
  };
}
