import { useEffect } from "react";

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
  const {
    lessons,
    selectedLesson,
    selectedLessonID,
    selectedStageIndex,
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

  useEffect(() => {
    if (!onLiveSnapshot) {
      return;
    }
    if (attempt) {
      onLiveSnapshot(snapshot, `${attempt.lesson_id} step ${attempt.stage_index + 1}`);
      return;
    }
    onLiveSnapshot(null, "");
  }, [attempt, onLiveSnapshot, snapshot]);

  return (
    <section className="panel lesson-panel">
      <h2>Challenge Runner</h2>
      <div className="lesson-controls">
        <label>
          Challenge
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
          Step
          <select
            value={selectedStageIndex}
            disabled={!selectedLesson || isStartPending}
            onChange={(event) => setSelectedStageIndexState(Number(event.target.value))}
          >
            {(selectedLesson?.stages ?? []).map((stage) => (
              <option key={stage.id} value={stage.index}>
                {stage.title}
              </option>
            ))}
          </select>
        </label>

        <button type="button" disabled={isStartPending || !selectedLessonID} onClick={handleStart}>
          {isStartPending ? "Starting..." : "Start Challenge"}
        </button>
        <button
          type="button"
          disabled={isGradePending || !attempt?.attempt_id}
          onClick={handleGrade}
        >
          {isGradePending ? "Checking..." : "Check"}
        </button>
      </div>

      {attempt ? (
        <div className="lesson-summary">
          <span>attempt: {attempt.attempt_id}</span>
          <span>session: {attempt.session_id}</span>
          <span>objective: {attempt.objective}</span>
          <span>limits: steps {attempt.limits.max_steps ?? 0}</span>
          <span>policy edits: {attempt.limits.max_policy_changes ?? 0}</span>
        </div>
      ) : (
        <p className="empty">
          Pick a challenge and start it. Then interact with the simulator and click Check.
        </p>
      )}

      {attempt ? (
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
            <button
              type="button"
              disabled={!canSend || !isCommandAllowed("spawn")}
              onClick={() =>
                handleCommand({
                  name: "spawn",
                  process: "demo",
                  program: "COMPUTE 5; SYSCALL sleep 2; COMPUTE 3; EXIT",
                })
              }
            >
              Spawn Demo
            </button>
          </div>

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

          <p className="lesson-outcome">
            Completed steps: {runResult.analytics.completed_stages}/
            {runResult.analytics.total_stages} (
            {formatPercent(runResult.analytics.completion_rate)})
          </p>
        </>
      ) : null}
    </section>
  );
}

function formatPercent(value: number): string {
  return `${Math.round(value * 100)}%`;
}
