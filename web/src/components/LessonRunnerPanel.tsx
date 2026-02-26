import { useEffect, useState } from "react";

import { useLessonRunner } from "../hooks/useLessonRunner";
import type { ChallengeGradeResponse } from "../lib/lessonApi";
import type { SnapshotDTO } from "../lib/types";

interface LessonRunnerPanelProps {
  baseURL: string;
  onGradeResult?: (result: ChallengeGradeResponse) => void;
  onLiveSnapshot?: (snapshot: SnapshotDTO | null, title: string) => void;
}

export function LessonRunnerPanel({
  baseURL,
  onGradeResult,
  onLiveSnapshot,
}: LessonRunnerPanelProps) {
  const [viewMode, setViewMode] = useState<"learn" | "exercise">("learn");
  const {
    lessons,
    selectedLesson,
    selectedLessonID,
    selectedStageIndex,
    selectedStage,
    runResult,
    attempt,
    policy,
    quantum,
    snapshot,
    liveError,
    canSend,
    errorMessage,
    isLessonsLoading,
    isStartPending,
    isGradePending,
    setPolicy,
    setQuantum,
    setSelectedStageIndexState,
    handleLessonChange,
    handleStart,
    handleCommand,
    handleGrade,
    isCommandAllowed,
  } = useLessonRunner({ baseURL, onGradeResult });

  const challengeState = snapshot?.challenge;
  const remainingSteps =
    challengeState?.remaining_steps ?? attempt?.limits.max_steps ?? 0;
  const remainingPolicyChanges =
    challengeState?.remaining_policy_changes ??
    attempt?.limits.max_policy_changes ??
    0;

  useEffect(() => {
    if (!onLiveSnapshot) {
      return;
    }
    if (attempt) {
      onLiveSnapshot(snapshot, `${attempt.lesson_id} stage ${attempt.stage_index + 1}`);
      return;
    }
    onLiveSnapshot(null, "");
  }, [attempt, onLiveSnapshot, snapshot]);

  return (
    <section className="panel lesson-panel">
      <h2>Challenge Exercise</h2>
      <div className="mode-nav challenge-subnav" role="tablist" aria-label="Challenge views">
        <button
          type="button"
          role="tab"
          aria-selected={viewMode === "learn"}
          className={viewMode === "learn" ? "mode-link active" : "mode-link"}
          onClick={() => setViewMode("learn")}
        >
          Learn
        </button>
        <button
          type="button"
          role="tab"
          aria-selected={viewMode === "exercise"}
          className={viewMode === "exercise" ? "mode-link active" : "mode-link"}
          onClick={() => setViewMode("exercise")}
        >
          Exercise
        </button>
      </div>
      <div className="lesson-controls">
        <label>
          Lesson
          <select
            value={selectedLessonID}
            disabled={isLessonsLoading || lessons.length === 0 || isStartPending}
            onChange={(event) => handleLessonChange(event.target.value)}
          >
            {lessons.map((lesson) => (
              <option key={lesson.id} value={lesson.id}>
                {lesson.module} - {lesson.title}
              </option>
            ))}
          </select>
        </label>

        <label>
          Stage
          <select
            value={selectedStageIndex}
            disabled={!selectedLesson || isStartPending}
            onChange={(event) => setSelectedStageIndexState(Number(event.target.value))}
          >
            {(selectedLesson?.stages ?? []).map((stage) => (
              <option key={stage.id} value={stage.index} disabled={stage.unlocked === false}>
                {stage.title}{stage.unlocked === false ? " (locked)" : ""}
              </option>
            ))}
          </select>
        </label>

        <button
          type="button"
          disabled={isStartPending || !selectedLessonID || selectedStage?.unlocked === false}
          onClick={handleStart}
        >
          {isStartPending ? "Starting..." : "Start Stage"}
        </button>
        <button
          type="button"
          disabled={isGradePending || !attempt?.attempt_id}
          onClick={handleGrade}
        >
          {isGradePending ? "Checking..." : "Check Result"}
        </button>
      </div>

      {attempt ? (
        <div className="lesson-summary">
          <span>attempt: {attempt.attempt_id}</span>
          <span>session: {attempt.session_id}</span>
          <span>objective: {attempt.objective}</span>
          <span>limits: steps {attempt.limits.max_steps ?? 0}</span>
          <span>policy edits: {attempt.limits.max_policy_changes ?? 0}</span>
          <span>steps left: {remainingSteps}</span>
          <span>policy edits left: {remainingPolicyChanges}</span>
        </div>
      ) : (
        <p className="empty">
          Pick a lesson stage, review Learn, then run Exercise and use Check Result.
        </p>
      )}

      {selectedStage ? (
        <div className="lesson-summary">
          <span>selected objective: {selectedStage.objective ?? selectedStage.title}</span>
          <span>status: {selectedStage.completed ? "passed" : selectedStage.unlocked === false ? "locked" : "ready"}</span>
          <span>attempts: {selectedStage.attempts ?? 0}</span>
          <span>prerequisites: {(selectedStage.prerequisites ?? []).length}</span>
        </div>
      ) : null}

      {viewMode === "learn" && selectedStage ? (
        <section className="lesson-learn-block" role="tabpanel">
          <h3>Theory</h3>
          <p className="lesson-outcome">
            {selectedStage.theory ??
              "Study trace order and outcome metrics before running the exercise."}
          </p>
          <h3>What You Need To Achieve</h3>
          <p className="lesson-outcome">{selectedStage.objective ?? selectedStage.title}</p>
          {(selectedStage.pass_conditions ?? []).length ? (
            <div className="lesson-validator-results">
              {(selectedStage.pass_conditions ?? []).map((item) => (
                <p key={item} className="lesson-outcome">
                  - {item}
                </p>
              ))}
            </div>
          ) : null}
          <p className="lesson-outcome">
            Allowed actions: {(selectedStage.allowed_commands ?? []).join(", ") || "step, run, pause, reset"}
          </p>
          <p className="lesson-outcome">
            Limits: steps {selectedStage.limits?.max_steps ?? 0}, policy edits {selectedStage.limits?.max_policy_changes ?? 0}
          </p>
          <p className="challenge-note">
            Move to Exercise when ready, run actions, then click Check to validate against deterministic rules.
          </p>
        </section>
      ) : null}

      {viewMode === "exercise" && attempt ? (
        <>
          <div className="control-row">
            <button
              type="button"
              disabled={!canSend || !isCommandAllowed("run")}
              onClick={() => handleCommand({ name: "run", count: 8 })}
            >
              Run 8
            </button>
            <button
              type="button"
              disabled={!canSend || !isCommandAllowed("step")}
              onClick={() => handleCommand({ name: "step", count: 1 })}
            >
              Step
            </button>
            <button
              type="button"
              disabled={!canSend || !isCommandAllowed("pause")}
              onClick={() => handleCommand({ name: "pause" })}
            >
              Pause
            </button>
            <button
              type="button"
              disabled={!canSend || !isCommandAllowed("reset")}
              onClick={() => handleCommand({ name: "reset" })}
            >
              Reset
            </button>
          </div>

          {isCommandAllowed("policy") ? (
            <div className="control-row">
              <label>
                Policy
                <select
                  value={policy}
                  disabled={!canSend || !isCommandAllowed("policy")}
                  onChange={(event) =>
                    setPolicy(event.target.value as "fifo" | "rr" | "mlfq")
                  }
                >
                  <option value="fifo">FIFO</option>
                  <option value="rr">RR</option>
                  <option value="mlfq">MLFQ</option>
                </select>
              </label>
              <label>
                Quantum
                <input
                  type="number"
                  min={1}
                  max={16}
                  value={quantum}
                  disabled={!canSend || !isCommandAllowed("policy") || policy !== "rr"}
                  onChange={(event) => setQuantum(Number(event.target.value))}
                />
              </label>
              <button
                type="button"
                disabled={!canSend || !isCommandAllowed("policy")}
                onClick={() =>
                  handleCommand({
                    name: "policy",
                    policy,
                    quantum: policy === "rr" ? quantum : 0,
                  })
                }
              >
                Apply Policy
              </button>
            </div>
          ) : null}
        </>
      ) : null}

      {errorMessage ? <p className="error">{errorMessage}</p> : null}
      {liveError ? <p className="error">{liveError}</p> : null}

      {runResult ? (
        <>
          <div className="lesson-summary">
            <span className={runResult.passed ? "badge pass" : "badge fail"}>
              {runResult.passed ? "passed" : "failed"}
            </span>
            <span>result: {runResult.feedback_key}</span>
            <span>trace hash: {runResult.output.trace_hash}</span>
            <span>trace length: {runResult.output.trace_length}</span>
          </div>

          {!runResult.passed && runResult.hint ? (
            <p className="hint">
              Hint L{runResult.hint_level ?? 0}: {runResult.hint}
            </p>
          ) : null}

          {!runResult.passed && (runResult.hint_level ?? 0) < 3 ? (
            <div className="control-row">
              <button type="button" disabled={isGradePending} onClick={handleGrade}>
                {isGradePending ? "Checking..." : "Get Next Hint"}
              </button>
              <span className="hint-progress">
                Hint progression: L{runResult.hint_level ?? 0} {"->"} L3
              </span>
            </div>
          ) : null}

          <p className="lesson-outcome">
            Completed stages: {runResult.analytics.completed_stages}/
            {runResult.analytics.total_stages} (
            {formatPercent(runResult.analytics.completion_rate)})
          </p>

          {runResult.validator_results?.length ? (
            <div className="validator-groups">
              <section className="validator-group">
                <h3>Failed Checks</h3>
                {runResult.validator_results.filter((item) => !item.passed).length ? (
                  runResult.validator_results
                    .filter((item) => !item.passed)
                    .map((item) => (
                      <p key={`${item.name}.${item.type}`} className="lesson-outcome">
                        - {describeCheck(item)}
                      </p>
                    ))
                ) : (
                  <p className="lesson-outcome">- None</p>
                )}
              </section>

              <section className="validator-group">
                <h3>Passed Checks</h3>
                {runResult.validator_results.filter((item) => item.passed).length ? (
                  runResult.validator_results
                    .filter((item) => item.passed)
                    .map((item) => (
                      <p key={`${item.name}.${item.type}`} className="lesson-outcome">
                        - {describeCheck(item)}
                      </p>
                    ))
                ) : (
                  <p className="lesson-outcome">- None</p>
                )}
              </section>
            </div>
          ) : null}
        </>
      ) : null}
    </section>
  );
}

function formatPercent(value: number): string {
  return `${Math.round(value * 100)}%`;
}

function describeCheck(item: {
  name: string;
  type: string;
  key?: string;
  passed: boolean;
  message?: string;
}): string {
  const keyText = item.key ? ` (${item.key})` : "";
  if (item.message && item.message.trim() !== "") {
    return humanizeMessage(item.message);
  }
  if (item.type === "trace_contains_all") {
    return `Required trace events appeared${keyText}.`;
  }
  if (item.type === "metric_eq") {
    return `Required metric matched expected value${keyText}.`;
  }
  if (item.type === "fault_eq" || item.type === "fault_lte") {
    return `Required fault count condition passed${keyText}.`;
  }
  if (item.type === "fs_ok") {
    return "Filesystem invariants held.";
  }
  return `${item.name} check passed${keyText}.`;
}

function humanizeMessage(message: string): string {
  if (message.startsWith("missing trace event ")) {
    const eventName = message.replace("missing trace event ", "");
    return `Trace is missing event '${eventName}'.`;
  }
  if (message.startsWith("metric ")) {
    return message
      .replace("metric ", "Metric '")
      .replace(" got=", "' is ")
      .replace(" want=", ", expected ")
      .replace(" want<=", ", expected <= ");
  }
  if (message.startsWith("fault ")) {
    return message
      .replace("fault ", "Fault '")
      .replace(" got=", "' is ")
      .replace(" want=", ", expected ")
      .replace(" want<=", ", expected <= ");
  }
  return message;
}
