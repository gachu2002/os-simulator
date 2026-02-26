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
}

export interface LessonSummary {
  id: string;
  title: string;
  module: string;
  stages: LessonStageSummary[];
}

export interface LessonsResponse {
  lessons: LessonSummary[];
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
}

export interface ChallengeStartResponse {
  attempt_id: string;
  session_id: string;
  lesson_id: string;
  stage_index: number;
  stage_title: string;
  module: string;
  objective: string;
  allowed_commands: string[];
  limits: ChallengeLimits;
}

export interface ChallengeGradeResponse {
  attempt_id: string;
  lesson_id: string;
  stage_index: number;
  passed: boolean;
  feedback_key: string;
  hint?: string;
  hint_level?: number;
  output: LessonRunOutput;
  analytics: CompletionAnalytics;
}

export async function fetchLessons(baseURL: string): Promise<LessonSummary[]> {
  const payload = await fetchJSON<LessonsResponse>(baseURL, "/lessons");
  return payload.lessons;
}

export async function startChallenge(
  baseURL: string,
  lessonID: string,
  stageIndex: number,
): Promise<ChallengeStartResponse> {
  return fetchJSON<ChallengeStartResponse>(baseURL, "/challenges/start", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ lesson_id: lessonID, stage_index: stageIndex }),
  });
}

export async function gradeChallenge(
  baseURL: string,
  attemptID: string,
): Promise<ChallengeGradeResponse> {
  return fetchJSON<ChallengeGradeResponse>(baseURL, "/challenges/grade", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ attempt_id: attemptID }),
  });
}
