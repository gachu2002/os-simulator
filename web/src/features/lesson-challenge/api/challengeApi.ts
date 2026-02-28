import type { ChallengeGrade, ChallengeStart } from "../../../entities/challenge/model";
import { fetchJSON } from "../../../lib/http";
import type {
  MemorySnapshot,
  ProcessSnapshot,
  SchedulingMetrics,
} from "../../../lib/types";

interface ChallengeStartDTO {
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
  limits: {
    max_steps?: number;
    max_policy_changes?: number;
    max_config_changes?: number;
  };
}

interface ChallengeGradeDTO {
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
  output: {
    tick: number;
    trace_hash: string;
    trace_length: number;
    processes: ProcessSnapshot[];
    metrics: SchedulingMetrics;
    memory: MemorySnapshot;
    filesystem_ok: boolean;
  };
  analytics: {
    total_stages: number;
    completed_stages: number;
    attempted_stages: number;
    completion_rate: number;
  };
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

export async function startChallenge(
  baseURL: string,
  lessonID: string,
  stageIndex: number,
  learnerID?: string,
): Promise<ChallengeStart> {
  const body: Record<string, string | number> = {
    lesson_id: lessonID,
    stage_index: stageIndex,
  };
  if (learnerID && learnerID.trim() !== "") {
    body.learner_id = learnerID;
  }
  const dto = await fetchJSON<ChallengeStartDTO>(baseURL, "/challenges/start", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(body),
  });
  return {
    attemptId: dto.attempt_id,
    sessionId: dto.session_id,
    lessonId: dto.lesson_id,
    stageIndex: dto.stage_index,
    stageTitle: dto.stage_title,
    module: dto.module,
    objective: dto.objective,
    goal: dto.goal,
    passConditions: dto.pass_conditions,
    allowedCommands: dto.allowed_commands,
    limits: {
      maxSteps: dto.limits.max_steps,
      maxPolicyChanges: dto.limits.max_policy_changes,
      maxConfigChanges: dto.limits.max_config_changes,
    },
  };
}

export async function submitChallenge(
  baseURL: string,
  attemptID: string,
  learnerID?: string,
): Promise<ChallengeGrade> {
  const body: Record<string, string> = { attempt_id: attemptID };
  if (learnerID && learnerID.trim() !== "") {
    body.learner_id = learnerID;
  }
  const dto = await fetchJSON<ChallengeGradeDTO>(baseURL, "/challenges/submit", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(body),
  });
  return {
    attemptId: dto.attempt_id,
    lessonId: dto.lesson_id,
    stageIndex: dto.stage_index,
    passed: dto.passed,
    feedbackKey: dto.feedback_key,
    objective: dto.objective,
    goal: dto.goal,
    passConditions: dto.pass_conditions,
    hint: dto.hint,
    hintLevel: dto.hint_level,
    output: {
      tick: dto.output.tick,
      traceHash: dto.output.trace_hash,
      traceLength: dto.output.trace_length,
      processes: dto.output.processes,
      metrics: dto.output.metrics,
      memory: dto.output.memory,
      filesystemOk: dto.output.filesystem_ok,
    },
    analytics: {
      totalStages: dto.analytics.total_stages,
      completedStages: dto.analytics.completed_stages,
      attemptedStages: dto.analytics.attempted_stages,
      completionRate: dto.analytics.completion_rate,
    },
    validatorResults: dto.validator_results,
  };
}
