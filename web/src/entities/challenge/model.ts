import type { MemorySnapshot, ProcessSnapshot, SchedulingMetrics } from "../../lib/types";
import type {
  ActionCapabilities,
  ActionCapabilityNotes,
} from "./actionCapabilities";

export interface CompletionAnalytics {
  totalStages: number;
  completedStages: number;
  attemptedStages: number;
  completionRate: number;
}

export interface LessonRunOutput {
  tick: number;
  traceHash: string;
  traceLength: number;
  processes: ProcessSnapshot[];
  metrics: SchedulingMetrics;
  memory: MemorySnapshot;
  filesystemOk: boolean;
}

export interface ChallengeLimits {
  maxSteps?: number;
  maxPolicyChanges?: number;
  maxConfigChanges?: number;
}

export interface ChallengeStart {
  attemptId: string;
  sessionId: string;
  lessonId: string;
  stageIndex: number;
  stageTitle: string;
  module: string;
  objective: string;
  goal?: string;
  passConditions?: string[];
  allowedCommands: string[];
  limits: ChallengeLimits;
  actionCapabilities?: ActionCapabilities;
  actionCapabilityNotes?: ActionCapabilityNotes;
}

export interface ValidatorResult {
  name: string;
  type: string;
  key?: string;
  passed: boolean;
  message?: string;
  expected?: string;
  actual?: string;
}

export interface ChallengeGrade {
  attemptId: string;
  lessonId: string;
  stageIndex: number;
  passed: boolean;
  feedbackKey: string;
  objective: string;
  goal?: string;
  passConditions?: string[];
  hint?: string;
  hintLevel?: number;
  output: LessonRunOutput;
  analytics: CompletionAnalytics;
  validatorResults?: ValidatorResult[];
}
