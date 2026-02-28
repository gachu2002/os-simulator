import type { CurriculumSection, LessonSummary, LessonStageSummary } from "../../../entities/lesson/model";
import { fetchJSON } from "../../../lib/http";

interface CurriculumResponseDTO {
  sections: CurriculumSectionDTO[];
}

interface CurriculumSectionDTO {
  id: string;
  title: string;
  subtitle?: string;
  order: number;
  coming_soon: boolean;
  lessons?: LessonSummaryDTO[];
  completed_stages?: number;
  total_stages?: number;
}

interface LessonSummaryDTO {
  id: string;
  title: string;
  module: string;
  section_id?: string;
  section_title?: string;
  difficulty?: string;
  estimated_minutes?: number;
  chapter_refs?: string[];
  stages: LessonStageSummaryDTO[];
}

interface LessonStageSummaryDTO {
  index: number;
  id: string;
  title: string;
  objective?: string;
  goal?: string;
  pass_conditions?: string[];
  prerequisites?: string[];
  allowed_commands?: string[];
  action_descriptions?: Array<{ command: string; description: string }>;
  expected_visual_cues?: string[];
  limits?: {
    max_steps?: number;
    max_policy_changes?: number;
    max_config_changes?: number;
  };
  attempts?: number;
  completed?: boolean;
  unlocked?: boolean;
}

export async function fetchCurriculumForLearner(
  baseURL: string,
  learnerID: string,
): Promise<CurriculumSection[]> {
  const payload = await fetchJSON<CurriculumResponseDTO>(
    baseURL,
    `/curriculum?learner_id=${encodeURIComponent(learnerID)}`,
  );
  return payload.sections.map(mapSection);
}

function mapSection(dto: CurriculumSectionDTO): CurriculumSection {
  return {
    id: dto.id,
    title: dto.title,
    subtitle: dto.subtitle,
    order: dto.order,
    comingSoon: dto.coming_soon,
    lessons: dto.lessons?.map(mapLessonSummary),
    completedStages: dto.completed_stages,
    totalStages: dto.total_stages,
  };
}

function mapLessonSummary(dto: LessonSummaryDTO): LessonSummary {
  return {
    id: dto.id,
    title: dto.title,
    module: dto.module,
    sectionId: dto.section_id,
    sectionTitle: dto.section_title,
    difficulty: dto.difficulty,
    estimatedMinutes: dto.estimated_minutes,
    chapterRefs: dto.chapter_refs,
    stages: dto.stages.map(mapLessonStageSummary),
  };
}

function mapLessonStageSummary(dto: LessonStageSummaryDTO): LessonStageSummary {
  return {
    index: dto.index,
    id: dto.id,
    title: dto.title,
    objective: dto.objective,
    goal: dto.goal,
    passConditions: dto.pass_conditions,
    prerequisites: dto.prerequisites,
    allowedCommands: dto.allowed_commands,
    actionDescriptions: dto.action_descriptions,
    expectedVisualCues: dto.expected_visual_cues,
    limits: dto.limits
      ? {
          maxSteps: dto.limits.max_steps,
          maxPolicyChanges: dto.limits.max_policy_changes,
          maxConfigChanges: dto.limits.max_config_changes,
        }
      : undefined,
    attempts: dto.attempts,
    completed: dto.completed,
    unlocked: dto.unlocked,
  };
}
