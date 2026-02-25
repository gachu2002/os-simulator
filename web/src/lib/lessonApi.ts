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
  prompt?: string;
  difficulty?: string;
  estimated_minutes?: number;
  concept_tags?: string[];
  prerequisites?: string[];
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

export interface ModuleAnalytics {
  module: string;
  total_stages: number;
  completed_stage: number;
  completion_rate: number;
}

export interface CompletionAnalytics {
  total_stages: number;
  completed_stages: number;
  attempted_stages: number;
  completion_rate: number;
  attempt_coverage: number;
  module_breakdown: ModuleAnalytics[];
  weak_concepts: ConceptWeakness[];
  pilot_checklist: string[];
  pilot_checklist_ok: boolean;
}

export interface ConceptWeakness {
  concept: string;
  score: number;
  failed_attempts: number;
  high_hint_uses: number;
  affected_stages: number;
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

export interface LessonRunResponse {
  lesson_id: string;
  stage_index: number;
  passed: boolean;
  feedback_key: string;
  hint?: string;
  hint_level?: number;
  output: LessonRunOutput;
  analytics: CompletionAnalytics;
}

export interface LessonProgressResponse {
  analytics: CompletionAnalytics;
}

export async function fetchLessons(baseURL: string): Promise<LessonSummary[]> {
  const payload = await fetchJSON<LessonsResponse>(baseURL, "/lessons");
  return payload.lessons;
}

export async function runLessonStage(
  baseURL: string,
  lessonID: string,
  stageIndex: number,
): Promise<LessonRunResponse> {
  return fetchJSON<LessonRunResponse>(baseURL, "/lessons/run", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ lesson_id: lessonID, stage_index: stageIndex }),
  });
}

export async function fetchLessonProgress(
  baseURL: string,
): Promise<LessonProgressResponse> {
  return fetchJSON<LessonProgressResponse>(baseURL, "/lessons/progress");
}
