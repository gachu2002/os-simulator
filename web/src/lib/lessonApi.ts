import type {
  MemorySnapshot,
  ProcessSnapshot,
  SchedulingMetrics,
} from "./types";
import { fetchJSON } from "./http";

export interface LessonStageSummary {
  index: number;
  id: string;
  title: string;
  objective?: string;
  goal?: string;
  pass_conditions?: string[];
  prerequisites?: string[];
  allowed_commands?: string[];
  action_descriptions?: ActionDescription[];
  expected_visual_cues?: string[];
  limits?: ChallengeLimits;
  attempts?: number;
  completed?: boolean;
  unlocked?: boolean;
}

export interface ActionDescription {
  command: string;
  description: string;
}

export interface LessonSummary {
  id: string;
  title: string;
  module: string;
  section_id?: string;
  section_title?: string;
  difficulty?: string;
  estimated_minutes?: number;
  chapter_refs?: string[];
  stages: LessonStageSummary[];
}

export interface CurriculumSection {
  id: string;
  title: string;
  subtitle?: string;
  order: number;
  coming_soon: boolean;
  lessons?: LessonSummary[];
  completed_stages?: number;
  total_stages?: number;
}

export interface CurriculumResponse {
  sections: CurriculumSection[];
}

export interface LessonLearnStage {
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

export interface LessonLearnSummary {
  id: string;
  title: string;
  module: string;
  section_id?: string;
  section_title?: string;
  difficulty?: string;
  estimated_minutes?: number;
  chapter_refs?: string[];
  stages: LessonLearnStage[];
}

export interface LessonLearnResponse {
  lesson: LessonLearnSummary;
}

export interface CompletionAnalytics {
  total_stages: number;
  completed_stages: number;
  attempted_stages: number;
  completion_rate: number;
}

export interface LessonRunOutput {
  tick: number;
  trace_hash: string;
  trace_length: number;
  processes: ProcessSnapshot[];
  metrics: SchedulingMetrics;
  memory: MemorySnapshot;
  filesystem_ok: boolean;
}

export interface ChallengeLimits {
  max_steps?: number;
  max_policy_changes?: number;
  max_config_changes?: number;
}

export interface ChallengeStartResponse {
  attempt_id: string;
  session_id: string;
  lesson_id: string;
  stage_index: number;
  stage_title: string;
  module: string;
  objective: string;
  goal?: string;
  pass_conditions?: string[];
  allowed_commands: string[];
  limits: ChallengeLimits;
}

export interface ChallengeGradeResponse {
  attempt_id: string;
  lesson_id: string;
  stage_index: number;
  passed: boolean;
  feedback_key: string;
  objective: string;
  goal?: string;
  pass_conditions?: string[];
  hint?: string;
  hint_level?: number;
  output: LessonRunOutput;
  analytics: CompletionAnalytics;
  validator_results?: Array<{
    name: string;
    type: string;
    key?: string;
    passed: boolean;
    message?: string;
    expected?: string;
    actual?: string;
  }>;
}

export async function fetchCurriculumForLearner(
  baseURL: string,
  learnerID: string,
): Promise<CurriculumSection[]> {
  const payload = await fetchJSON<CurriculumResponse>(
    baseURL,
    `/curriculum?learner_id=${encodeURIComponent(learnerID)}`,
  );
  return payload.sections;
}

export async function fetchLessonLearn(
  baseURL: string,
  lessonID: string,
  learnerID: string,
): Promise<LessonLearnSummary> {
  const payload = await fetchJSON<LessonLearnResponse>(
    baseURL,
    `/lessons/${encodeURIComponent(lessonID)}/learn?learner_id=${encodeURIComponent(learnerID)}`,
  );
  return payload.lesson;
}

export async function startChallenge(
  baseURL: string,
  lessonID: string,
  stageIndex: number,
  learnerID?: string,
): Promise<ChallengeStartResponse> {
  const body: Record<string, string | number> = {
    lesson_id: lessonID,
    stage_index: stageIndex,
  };
  if (learnerID && learnerID.trim() !== "") {
    body.learner_id = learnerID;
  }
  return fetchJSON<ChallengeStartResponse>(baseURL, "/challenges/start", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(body),
  });
}

export async function submitChallenge(
  baseURL: string,
  attemptID: string,
  learnerID?: string,
): Promise<ChallengeGradeResponse> {
  const body: Record<string, string> = { attempt_id: attemptID };
  if (learnerID && learnerID.trim() !== "") {
    body.learner_id = learnerID;
  }
  return fetchJSON<ChallengeGradeResponse>(baseURL, "/challenges/submit", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(body),
  });
}
