import type { ChallengeGrade, ChallengeStart } from "../../../entities/challenge/model";
import {
  fromActionCapabilitiesDTO,
  fromActionCapabilityNotesDTO,
} from "../../../entities/challenge/actionCapabilities";
import { fetchJSON } from "../../../lib/http";
import type {
  MemorySnapshot,
  ProcessSnapshot,
  SchedulingMetrics,
  SessionEvent,
} from "../../../lib/types";
import type { SchedulerPolicy } from "../model/actionPresets";

interface ChallengeStartV3DTO {
  version: string;
  section_id: string;
  lesson_id: string;
  lesson_title: string;
  lesson_objective: string;
  part_id?: string;
  part_title?: string;
  part_objective?: string;
  attempt_id: string;
  session_id: string;
  allowed_commands: string[];
  limits: {
    max_steps?: number;
    max_policy_changes?: number;
    max_config_changes?: number;
  };
  action_capabilities?: {
    supported_now: string[];
    planned: string[];
  };
  action_capability_notes?: Record<
    string,
    {
      status: string;
      reason?: string;
      fallback_action?: string;
      mapped_command?: string;
    }
  >;
}

interface ChallengeActionV3DTO {
  attempt_id: string;
  session_id: string;
  action: string;
  mapped_command: string;
  event: SessionEvent;
}

interface ChallengeGradeDTO {
  attempt_id: string;
  lesson_id: string;
  part_id?: string;
  passed: boolean;
  feedback_key: string;
  lesson_objective: string;
  part_objective?: string;
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

interface ChallengeGradeV3DTO extends ChallengeGradeDTO {
  version: string;
  section_id: string;
  lesson_title: string;
  part_title?: string;
}


export async function startChallenge(
  baseURL: string,
  lessonID: string,
  stageID?: string,
  learnerID?: string,
): Promise<ChallengeStart> {
  return startChallengeV3(baseURL, lessonID, stageID, learnerID);
}

export async function startChallengeV3(
  baseURL: string,
  lessonID: string,
  stageID?: string,
  learnerID?: string,
): Promise<ChallengeStart> {
  const body: Record<string, string | number> = {
    lesson_id: lessonID,
  };
  const normalizedStageID = (stageID ?? "").trim().toUpperCase();
  if (normalizedStageID === "A" || normalizedStageID === "B") {
    body.part_id = normalizedStageID;
  }
  if (learnerID && learnerID.trim() !== "") {
    body.learner_id = learnerID;
  }
  const dto = await fetchJSON<ChallengeStartV3DTO>(baseURL, "/challenges/start/v3", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(body),
  });
  return {
    attemptId: dto.attempt_id,
    sessionId: dto.session_id,
    lessonId: dto.lesson_id,
    stageIndex: dto.part_id === "B" ? 1 : 0,
    stageTitle: dto.part_title ?? dto.lesson_title,
    module: dto.section_id,
    objective: dto.lesson_objective,
    goal: dto.part_objective,
    passConditions: dto.part_objective ? [dto.part_objective] : [dto.lesson_objective],
    allowedCommands: dto.allowed_commands,
    limits: {
      maxSteps: dto.limits.max_steps,
      maxPolicyChanges: dto.limits.max_policy_changes,
      maxConfigChanges: dto.limits.max_config_changes,
    },
    actionCapabilities: fromActionCapabilitiesDTO(dto.action_capabilities),
    actionCapabilityNotes: fromActionCapabilityNotesDTO(dto.action_capability_notes),
  };
}

export async function submitChallenge(
  baseURL: string,
  attemptID: string,
  learnerID?: string,
): Promise<ChallengeGrade> {
  return submitChallengeV3(baseURL, attemptID, learnerID);
}

export async function submitChallengeV3(
  baseURL: string,
  attemptID: string,
  learnerID?: string,
): Promise<ChallengeGrade> {
  const body: Record<string, string> = { attempt_id: attemptID };
  if (learnerID && learnerID.trim() !== "") {
    body.learner_id = learnerID;
  }
  const dto = await fetchJSON<ChallengeGradeV3DTO>(baseURL, "/challenges/submit/v3", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(body),
  });
  return {
    attemptId: dto.attempt_id,
    lessonId: dto.lesson_id,
    stageIndex: dto.part_id === "B" ? 1 : 0,
    passed: dto.passed,
    feedbackKey: dto.feedback_key,
    objective: dto.lesson_objective,
    goal: dto.part_objective,
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

export async function actionChallengeV3(
  baseURL: string,
  payload: {
    attemptID: string;
    learnerID?: string;
    action: string;
    count?: number;
    process?: string;
    program?: string;
    policy?: SchedulerPolicy;
    quantum?: number;
    frames?: number;
    tlbEntries?: number;
    diskLatency?: number;
    terminalLatency?: number;
  },
): Promise<ChallengeActionV3DTO> {
  const dto = await fetchJSON<ChallengeActionV3DTO>(baseURL, "/challenges/action/v3", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({
      attempt_id: payload.attemptID,
      learner_id: payload.learnerID,
      action: payload.action,
      count: payload.count,
      process: payload.process,
      program: payload.program,
      policy: payload.policy,
      quantum: payload.quantum,
      frames: payload.frames,
      tlb_entries: payload.tlbEntries,
      disk_latency: payload.diskLatency,
      terminal_latency: payload.terminalLatency,
    }),
  });
  return dto;
}
