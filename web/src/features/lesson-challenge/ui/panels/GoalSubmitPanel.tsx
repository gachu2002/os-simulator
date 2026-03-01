import type { ChallengeGrade, ValidatorResult } from "../../../../entities/challenge/model";
import type { LessonStageSummary } from "../../../../entities/lesson/model";
import { Badge } from "../../../../components/ui/badge";

interface GoalSubmitPanelProps {
  selectedStage: LessonStageSummary | null;
  attemptGoal?: string;
  attemptPassConditions?: string[];
  result: ChallengeGrade | null;
}

export function GoalSubmitPanel({
  selectedStage,
  attemptGoal,
  attemptPassConditions,
  result,
}: GoalSubmitPanelProps) {
  const passConditions = attemptPassConditions ?? selectedStage?.passConditions ?? [];
  const validatorResults = result?.validatorResults ?? [];
  const coreResults = validatorResults.filter((item) => item.type !== "v3_quality");
  const qualityResults = validatorResults.filter((item) => item.type === "v3_quality");
  const failedQualityResults = qualityResults.filter((item) => !item.passed);
  const nextSteps = buildNextSteps(result, failedQualityResults.map((item) => item.name));
  const qualityGateProgress = buildQualityGateProgress(passConditions, qualityResults);

  return (
    <section className="rounded-lg border border-slate-200 bg-white p-3">
      <h3 className="text-sm font-semibold text-slate-900">Result</h3>
      <p className="mt-2 text-sm text-slate-600">
        Goal: {attemptGoal ?? selectedStage?.goal ?? selectedStage?.objective ?? "Pass all checks."}
      </p>

      {passConditions.map((item) => (
        <p key={item} className="mt-2 text-sm text-slate-600">
          - {item}
        </p>
      ))}

      {qualityGateProgress.length > 0 ? (
        <div className="mt-3 rounded-md border border-slate-200 bg-slate-50 p-2">
          <p className="text-xs font-semibold uppercase tracking-wide text-slate-700">Quality Gate Progress</p>
          <div className="mt-2 flex flex-wrap gap-2">
            {qualityGateProgress.map((item) => (
              <span key={item.name} className="inline-flex items-center gap-1 rounded-full border border-slate-200 bg-white px-2 py-1 text-xs text-slate-700">
                <span>{toReadableGateName(item.name)}</span>
                <Badge variant={item.status === "passed" ? "success" : item.status === "failed" ? "destructive" : "secondary"}>
                  {item.status === "passed" ? "passed" : item.status === "failed" ? "needs work" : "not tried"}
                </Badge>
              </span>
            ))}
          </div>
        </div>
      ) : null}

      {result ? (
        <>
          <div className="mt-3 flex flex-wrap gap-2 text-sm text-slate-600">
            <Badge variant={result.passed ? "success" : "destructive"}>
              {result.passed ? "passed" : "failed"}
            </Badge>
            <span>feedback: {result.feedbackKey}</span>
          </div>

          {coreResults.length > 0 ? (
            <div className="mt-3 rounded-md border border-slate-200 bg-slate-50 p-2">
              <p className="text-xs font-semibold uppercase tracking-wide text-slate-700">Core Validators</p>
              {coreResults.map((item) => (
                <p key={`${item.name}.${item.type}`} className="mt-1 text-sm text-slate-600">
                  - {item.passed ? "PASS" : "FAIL"}: {item.name} | expected {item.expected ?? "n/a"}, actual {item.actual ?? "n/a"}
                </p>
              ))}
            </div>
          ) : null}

          {qualityResults.length > 0 ? (
            <div className="mt-3 rounded-md border border-sky-200 bg-sky-50 p-2">
              <p className="text-xs font-semibold uppercase tracking-wide text-sky-800">Challenge Quality Gates</p>
              {qualityResults.map((item) => (
                <p key={`${item.name}.${item.type}`} className="mt-1 text-sm text-slate-700">
                  - {item.passed ? "PASS" : "FAIL"}: {item.name} | expected {item.expected ?? "n/a"}, actual {item.actual ?? "n/a"}
                </p>
              ))}
            </div>
          ) : null}

          {nextSteps.length > 0 ? (
            <div className="mt-3 rounded-md border border-amber-200 bg-amber-50 p-2">
              <p className="text-xs font-semibold uppercase tracking-wide text-amber-900">What To Do Next</p>
              {nextSteps.map((item) => (
                <p key={item} className="mt-1 text-sm text-amber-900">
                  - {item}
                </p>
              ))}
            </div>
          ) : null}

          {!result.passed && result.hint ? (
            <p className="mt-2 text-sm text-orange-700">
              Hint L{result.hintLevel ?? 0}: {result.hint}
            </p>
          ) : null}
        </>
      ) : (
        <p className="mt-2 text-sm text-slate-600">
          Start challenge actions, then submit to get pass/fail results.
        </p>
      )}
    </section>
  );
}

interface QualityGateProgressItem {
  name: string;
  status: "passed" | "failed" | "pending";
}

function buildQualityGateProgress(
  passConditions: string[],
  qualityResults: ValidatorResult[],
): QualityGateProgressItem[] {
  const namesFromConditions = passConditions
    .map(extractQualityGateName)
    .filter((item): item is string => Boolean(item));

  const resultNames = (qualityResults ?? []).map((item) => item.name);
  const allNames = Array.from(new Set([...namesFromConditions, ...resultNames]));

  return allNames.map((name) => {
    const matched = (qualityResults ?? []).find((item) => item.name === name);
    if (!matched) {
      return { name, status: "pending" };
    }
    return { name, status: matched.passed ? "passed" : "failed" };
  });
}

function extractQualityGateName(condition: string): string | null {
  const normalized = condition.trim();
  const delimiterIndex = normalized.indexOf(":");
  if (delimiterIndex <= 0) {
    return null;
  }
  const candidate = normalized.slice(0, delimiterIndex).trim();
  if (!candidate.includes("_")) {
    return null;
  }
  return candidate;
}

function toReadableGateName(name: string): string {
  return name
    .split("_")
    .map((part) => part.charAt(0).toUpperCase() + part.slice(1))
    .join(" ");
}

function buildNextSteps(result: ChallengeGrade | null, failedQualityGateNames: string[]): string[] {
  if (!result || failedQualityGateNames.length === 0) {
    return [];
  }

  const lessonHintMap: Record<string, Record<string, string>> = {
    "l01-process-basics": {
      state_transitions: "Use Block Process then Unblock Process and verify lane movement in visualization.",
      process_lifecycle_end: "Terminate at least one process (kill or natural exit) before submitting.",
    },
    "l02-process-api-fork-exec-wait": {
      family_growth: "Create parent-child flow with fork/exec so family tree has multiple processes.",
      child_completion: "Run long enough for child completion and verify completion evidence.",
    },
    "l03-limited-direct-execution": {
      trap_round_trip: "Do execute_instruction -> issue_trap -> handle_syscall -> return_from_trap in order.",
      forced_switch: "Trigger timer/preemption and observe process switch in timeline.",
      manual_selection: "Use Choose Next Process at least once after preemption.",
    },
    "l04-cpu-scheduling-basics": {
      policy_comparison: "Change scheduling policy at least once and compare results on same workload.",
      completion_evidence: "Run until multiple jobs complete so metrics comparison is meaningful.",
      timeline_evidence: "Collect a longer run window before submitting.",
    },
    "l05-round-robin": {
      rr_tuning: "Adjust RR quantum at least once and rerun to compare outcomes.",
      interactive_window: "Run more ticks to expose response/turnaround tradeoff.",
      throughput_floor: "Increase quantum or workload stability to recover throughput.",
    },
    "l06-mlfq": {
      mlfq_activity: "Step longer until you observe queue pressure/preemption behavior.",
      mlfq_progress: "Collect more runtime to reveal demotion and boost effects.",
      fairness_window: "Continue run and tune actions until fairness improves.",
      mlfq_runtime: "Use a longer timeline before evaluating gaming outcomes.",
    },
    "l07-lottery-stride": {
      share_window: "Run more quanta to let proportional share converge.",
      share_signal: "Adjust ticket split and compare fairness trend over longer run.",
      workload_depth: "Use at least three active jobs for 50/30/20 target checks.",
    },
    "l08-multi-cpu-scheduling": {
      dispatch_depth: "Run more steps to observe enough dispatch/migration activity.",
      load_balance_signal: "Balance queues and avoid long idle streaks on one CPU.",
      utilization_window: "Collect a longer utilization timeline before submit.",
    },
  };

  const perLesson = lessonHintMap[result.lessonId] ?? {};
  const mapped = failedQualityGateNames
    .map((gateName) => perLesson[gateName])
    .filter((item): item is string => Boolean(item));

  if (mapped.length > 0) {
    return mapped;
  }

  if (result.hint) {
    return [result.hint];
  }
  return ["Review failed quality gates and gather stronger evidence before submitting."];
}
