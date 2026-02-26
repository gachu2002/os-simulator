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
