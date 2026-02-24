import type {
  MemorySnapshot,
  ProcessSnapshot,
  SchedulingMetrics,
} from "./types";

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
  pilot_checklist: string[];
  pilot_checklist_ok: boolean;
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
  const response = await fetch(`${baseURL}/lessons`);
  if (!response.ok) {
    throw new Error(`load lessons failed with status ${response.status}`);
  }
  const payload = (await response.json()) as LessonsResponse;
  return payload.lessons;
}

export async function runLessonStage(
  baseURL: string,
  lessonID: string,
  stageIndex: number,
): Promise<LessonRunResponse> {
  const response = await fetch(`${baseURL}/lessons/run`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ lesson_id: lessonID, stage_index: stageIndex }),
  });
  if (!response.ok) {
    throw new Error(`run lesson failed with status ${response.status}`);
  }
  return (await response.json()) as LessonRunResponse;
}
