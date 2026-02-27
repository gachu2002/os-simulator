import { useEffect, useState } from "react";

import { useLessonRunner } from "../hooks/useLessonRunner";
import type { ChallengeGradeResponse } from "../lib/lessonApi";
import type { SnapshotDTO } from "../lib/types";

interface LessonRunnerPanelProps {
  baseURL: string;
  onGradeResult?: (result: ChallengeGradeResponse) => void;
  onLiveSnapshot?: (snapshot: SnapshotDTO | null, title: string) => void;
  preferredLessonID?: string;
  preferredStageIndex?: number;
}

export function LessonRunnerPanel({
  baseURL,
  onGradeResult,
  onLiveSnapshot,
  preferredLessonID,
  preferredStageIndex,
}: LessonRunnerPanelProps) {
  const [viewMode, setViewMode] = useState<"learn" | "exercise">("learn");
  const [frames, setFrames] = useState(8);
  const [tlbEntries, setTLBEntries] = useState(4);
  const [diskLatency, setDiskLatency] = useState(3);
  const [terminalLatency, setTerminalLatency] = useState(1);
  const {
    selectedLesson,
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
    handleStart,
    handleCommand,
    handleGrade,
    isCommandAllowed,
  } = useLessonRunner({
    baseURL,
    onGradeResult,
    preferredLessonID,
    preferredStageIndex,
  });

  const challengeState = snapshot?.challenge;
  const remainingSteps =
    challengeState?.remaining_steps ?? attempt?.limits.max_steps ?? 0;
  const remainingPolicyChanges =
    challengeState?.remaining_policy_changes ??
    attempt?.limits.max_policy_changes ??
    0;
  const remainingConfigChanges =
    challengeState?.remaining_config_changes ??
    attempt?.limits.max_config_changes ??
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
      <h2>{selectedLesson?.title ?? "Challenge"}</h2>
      <p className="lesson-outcome">
        {selectedStage?.title ?? "Select a challenge from the overview page."}
      </p>
      <p className="subtitle">
        {selectedStage?.objective ?? "Review theory, then run the exercise and check results."}
      </p>
      <div className="mode-nav challenge-subnav" role="tablist" aria-label="Challenge views">
        <button
          type="button"
          role="tab"
          aria-selected={viewMode === "learn"}
          className={viewMode === "learn" ? "mode-link active btn btn-secondary" : "mode-link btn btn-secondary"}
          onClick={() => setViewMode("learn")}
        >
          Learn
        </button>
        <button
          type="button"
          role="tab"
          aria-selected={viewMode === "exercise"}
          className={viewMode === "exercise" ? "mode-link active btn btn-primary" : "mode-link btn btn-primary"}
          onClick={() => setViewMode("exercise")}
        >
          Exercise
        </button>
      </div>
      <div className="lesson-controls">
        <button
          type="button"
          className="btn btn-primary"
          disabled={isStartPending || isLessonsLoading || !selectedLesson || selectedStage?.unlocked === false}
          onClick={handleStart}
        >
          {isStartPending ? "Starting..." : "Start Stage"}
        </button>
        <button
          type="button"
          className="btn btn-success"
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
          <span>config edits: {attempt.limits.max_config_changes ?? 0}</span>
          <span>steps left: {remainingSteps}</span>
          <span>policy edits left: {remainingPolicyChanges}</span>
          <span>config edits left: {remainingConfigChanges}</span>
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
            {selectedStage.theory_detail ??
              selectedStage.theory ??
              "Study trace order and outcome metrics before running the exercise."}
          </p>
        </section>
      ) : null}

      {viewMode === "exercise" && attempt ? (
        <>
          <h3>Goal and Result To Achieve</h3>
          <p className="lesson-outcome">
            Goal: {selectedStage?.goal ?? selectedStage?.objective ?? selectedStage?.title ?? "Run actions and pass checks."}
          </p>

          {(selectedStage?.pass_conditions ?? []).length ? (
            <div className="lesson-validator-results">
              <h3>Pass Conditions</h3>
              {(selectedStage?.pass_conditions ?? []).map((item) => (
                <p key={item} className="lesson-outcome">
                  - {item}
                </p>
              ))}
            </div>
          ) : null}

          {(selectedStage?.action_descriptions ?? []).length ? (
            <div className="lesson-validator-results">
              <h3>Actions You Can Do</h3>
              {(selectedStage?.action_descriptions ?? []).map((item) => (
                <p key={item.command} className="lesson-outcome">
                  - {item.command}: {item.description}
                </p>
              ))}
            </div>
          ) : null}

          {(selectedStage?.expected_visual_cues ?? []).length ? (
            <div className="lesson-validator-results">
              <h3>Expected Visual Result</h3>
              {(selectedStage?.expected_visual_cues ?? []).map((item) => (
                <p key={item} className="lesson-outcome">
                  - {item}
                </p>
              ))}
            </div>
          ) : null}

          <p className="lesson-outcome">
            Limits: steps {selectedStage?.limits?.max_steps ?? 0}, policy edits {selectedStage?.limits?.max_policy_changes ?? 0}, config edits {selectedStage?.limits?.max_config_changes ?? 0}
          </p>

          <div className="control-row">
            <button
              type="button"
              className="btn btn-secondary"
              disabled={!canSend || !isCommandAllowed("run")}
              onClick={() => handleCommand({ name: "run", count: 8 })}
            >
              Run 8
            </button>
            <button
              type="button"
              className="btn btn-secondary"
              disabled={!canSend || !isCommandAllowed("step")}
              onClick={() => handleCommand({ name: "step", count: 1 })}
            >
              Step
            </button>
            <button
              type="button"
              className="btn btn-ghost"
              disabled={!canSend || !isCommandAllowed("pause")}
              onClick={() => handleCommand({ name: "pause" })}
            >
              Pause
            </button>
            <button
              type="button"
              className="btn btn-danger"
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
                className="btn btn-primary"
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

          {isCommandAllowed("set_frames") || isCommandAllowed("set_tlb_entries") ? (
            <div className="control-row">
              {isCommandAllowed("set_frames") ? (
                <>
                  <label>
                    Frames
                    <input
                      type="number"
                      min={1}
                      max={64}
                      value={frames}
                      disabled={!canSend || !isCommandAllowed("set_frames")}
                      onChange={(event) => setFrames(Number(event.target.value))}
                    />
                  </label>
                  <button
                    type="button"
                    className="btn btn-primary"
                    disabled={!canSend || !isCommandAllowed("set_frames")}
                    onClick={() => handleCommand({ name: "set_frames", frames })}
                  >
                    Apply Frames
                  </button>
                </>
              ) : null}

              {isCommandAllowed("set_tlb_entries") ? (
                <>
                  <label>
                    TLB entries
                    <input
                      type="number"
                      min={1}
                      max={64}
                      value={tlbEntries}
                      disabled={!canSend || !isCommandAllowed("set_tlb_entries")}
                      onChange={(event) => setTLBEntries(Number(event.target.value))}
                    />
                  </label>
                  <button
                    type="button"
                    className="btn btn-primary"
                    disabled={!canSend || !isCommandAllowed("set_tlb_entries")}
                    onClick={() =>
                      handleCommand({ name: "set_tlb_entries", tlb_entries: tlbEntries })
                    }
                  >
                    Apply TLB
                  </button>
                </>
              ) : null}
            </div>
          ) : null}

          {isCommandAllowed("set_disk_latency") || isCommandAllowed("set_terminal_latency") ? (
            <div className="control-row">
              {isCommandAllowed("set_disk_latency") ? (
                <>
                  <label>
                    Disk latency
                    <input
                      type="number"
                      min={1}
                      max={64}
                      value={diskLatency}
                      disabled={!canSend || !isCommandAllowed("set_disk_latency")}
                      onChange={(event) => setDiskLatency(Number(event.target.value))}
                    />
                  </label>
                  <button
                    type="button"
                    className="btn btn-primary"
                    disabled={!canSend || !isCommandAllowed("set_disk_latency")}
                    onClick={() =>
                      handleCommand({ name: "set_disk_latency", disk_latency: diskLatency })
                    }
                  >
                    Apply Disk Latency
                  </button>
                </>
              ) : null}

              {isCommandAllowed("set_terminal_latency") ? (
                <>
                  <label>
                    Terminal latency
                    <input
                      type="number"
                      min={1}
                      max={64}
                      value={terminalLatency}
                      disabled={!canSend || !isCommandAllowed("set_terminal_latency")}
                      onChange={(event) => setTerminalLatency(Number(event.target.value))}
                    />
                  </label>
                  <button
                    type="button"
                    className="btn btn-primary"
                    disabled={!canSend || !isCommandAllowed("set_terminal_latency")}
                    onClick={() =>
                      handleCommand({
                        name: "set_terminal_latency",
                        terminal_latency: terminalLatency,
                      })
                    }
                  >
                    Apply Terminal Latency
                  </button>
                </>
              ) : null}
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
                        - {describeCheck(item)} {failureGuidance(item)}
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

          {!runResult.passed && selectedStage ? (
            <section className="validator-group">
              <h3>Expected vs Actual</h3>
              <p className="lesson-outcome">
                - Expected: complete all pass conditions listed for this stage.
              </p>
              {(selectedStage.pass_conditions ?? []).map((item) => (
                <p key={item} className="lesson-outcome">
                  - Target: {item}
                </p>
              ))}
              {(selectedStage.expected_visual_cues ?? []).map((item) => (
                <p key={item} className="lesson-outcome">
                  - Visual cue: {item}
                </p>
              ))}
              <p className="lesson-outcome">
                - Actual: check failed items above, adjust actions, then re-run Check Result.
              </p>
            </section>
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
  const missingTrace = message.match(/^missing trace event\s+(.+)$/);
  if (missingTrace) {
    return `Trace is missing required event '${missingTrace[1]}'.`;
  }

  const metricEq = message.match(/^metric\s+(\S+)\s+got=([^\s]+)\s+want=([^\s]+)$/);
  if (metricEq) {
    const [, key, got, want] = metricEq;
    return `Metric '${key}' is ${got}, expected exactly ${want}.`;
  }

  const metricLTE = message.match(/^metric\s+(\S+)\s+got=([^\s]+)\s+want<=([^\s]+)$/);
  if (metricLTE) {
    const [, key, got, want] = metricLTE;
    return `Metric '${key}' is ${got}, expected <= ${want}.`;
  }

  const faultEq = message.match(/^fault\s+(\S+)\s+got=([^\s]+)\s+want=([^\s]+)$/);
  if (faultEq) {
    const [, key, got, want] = faultEq;
    return `Fault '${key}' is ${got}, expected exactly ${want}.`;
  }

  const faultLTE = message.match(/^fault\s+(\S+)\s+got=([^\s]+)\s+want<=([^\s]+)$/);
  if (faultLTE) {
    const [, key, got, want] = faultLTE;
    return `Fault '${key}' is ${got}, expected <= ${want}.`;
  }

  return message;
}

function failureGuidance(item: { type: string }): string {
  switch (item.type) {
    case "trace_contains_all":
      return "Inspect timeline ordering and trace events.";
    case "metric_eq":
    case "metric_lte":
      return "Inspect metrics panel values and scheduler behavior.";
    case "fault_eq":
    case "fault_lte":
      return "Inspect memory panel fault counters and recent access pattern.";
    case "fs_ok":
      return "Inspect filesystem-related trace events and invariant status.";
    default:
      return "";
  }
}
