import { useState } from "react";

import { useLessonRunner } from "../hooks/useLessonRunner";
import type { Command } from "../lib/types";
import { VisualizationSuite } from "./VisualizationSuite";

interface LessonChallengePageProps {
  baseURL: string;
  lessonID: string;
  stageIndex?: number;
  onNavigate: (to: string) => void;
}

export function LessonChallengePage({
  baseURL,
  lessonID,
  stageIndex,
  onNavigate,
}: LessonChallengePageProps) {
  const [frames, setFrames] = useState(8);
  const [tlbEntries, setTLBEntries] = useState(4);
  const [diskLatency, setDiskLatency] = useState(3);
  const [terminalLatency, setTerminalLatency] = useState(1);
  const {
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
    preferredLessonID: lessonID,
    preferredStageIndex: stageIndex,
  });

  const challengeState = snapshot?.challenge;
  const remainingSteps = challengeState?.remaining_steps ?? attempt?.limits.max_steps ?? 0;
  const remainingPolicyChanges =
    challengeState?.remaining_policy_changes ?? attempt?.limits.max_policy_changes ?? 0;
  const remainingConfigChanges =
    challengeState?.remaining_config_changes ?? attempt?.limits.max_config_changes ?? 0;

  return (
    <section className="panel lesson-panel">
      <div className="top-nav">
        <button type="button" className="btn btn-ghost" onClick={() => onNavigate("/")}>Home</button>
        <button
          type="button"
          className="btn btn-secondary"
          onClick={() => onNavigate(`/lesson/${lessonID}/learn?stage=${stageIndex ?? 0}`)}
        >
          Back To Learn
        </button>
      </div>

      <h2>{attempt?.lesson_id ?? lessonID} Challenge</h2>
      <p className="subtitle">Challenge page has three sections: Actions, Visualization, Goal + Submit.</p>

      <section className="lesson-learn-block">
        <h3>1) Actions</h3>
        <div className="lesson-controls">
          <button
            type="button"
            className="btn btn-primary"
            disabled={isStartPending || isLessonsLoading || selectedStage?.unlocked === false}
            onClick={handleStart}
          >
            {isStartPending ? "Starting..." : "Start Challenge"}
          </button>
          <button
            type="button"
            className="btn btn-success"
            disabled={isGradePending || !attempt?.attempt_id}
            onClick={handleGrade}
          >
            {isGradePending ? "Submitting..." : "Submit"}
          </button>
        </div>

        {attempt ? (
          <p className="lesson-outcome">
            Remaining budget: steps {remainingSteps}, policy edits {remainingPolicyChanges}, config edits {remainingConfigChanges}
          </p>
        ) : null}

        {attempt ? <ActionControls canSend={canSend} isAllowed={isCommandAllowed} onCommand={handleCommand} /> : null}

        {isCommandAllowed("policy") ? (
          <div className="control-row">
            <label>
              Policy
              <select
                value={policy}
                disabled={!canSend || !isCommandAllowed("policy")}
                onChange={(event) => setPolicy(event.target.value as "fifo" | "rr" | "mlfq")}
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
                handleCommand({ name: "policy", policy, quantum: policy === "rr" ? quantum : 0 })
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
                  <input type="number" min={1} max={64} value={frames} onChange={(event) => setFrames(Number(event.target.value))} />
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
                  TLB Entries
                  <input
                    type="number"
                    min={1}
                    max={64}
                    value={tlbEntries}
                    onChange={(event) => setTLBEntries(Number(event.target.value))}
                  />
                </label>
                <button
                  type="button"
                  className="btn btn-primary"
                  disabled={!canSend || !isCommandAllowed("set_tlb_entries")}
                  onClick={() => handleCommand({ name: "set_tlb_entries", tlb_entries: tlbEntries })}
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
                  Disk Latency
                  <input
                    type="number"
                    min={1}
                    max={64}
                    value={diskLatency}
                    onChange={(event) => setDiskLatency(Number(event.target.value))}
                  />
                </label>
                <button
                  type="button"
                  className="btn btn-primary"
                  disabled={!canSend || !isCommandAllowed("set_disk_latency")}
                  onClick={() => handleCommand({ name: "set_disk_latency", disk_latency: diskLatency })}
                >
                  Apply Disk Latency
                </button>
              </>
            ) : null}
            {isCommandAllowed("set_terminal_latency") ? (
              <>
                <label>
                  Terminal Latency
                  <input
                    type="number"
                    min={1}
                    max={64}
                    value={terminalLatency}
                    onChange={(event) => setTerminalLatency(Number(event.target.value))}
                  />
                </label>
                <button
                  type="button"
                  className="btn btn-primary"
                  disabled={!canSend || !isCommandAllowed("set_terminal_latency")}
                  onClick={() =>
                    handleCommand({ name: "set_terminal_latency", terminal_latency: terminalLatency })
                  }
                >
                  Apply Terminal Latency
                </button>
              </>
            ) : null}
          </div>
        ) : null}
      </section>

      <section className="lesson-learn-block">
        <h3>2) Visualization</h3>
        <VisualizationSuite
          title="Live Challenge State"
          subtitle="Run actions and inspect trace, memory, process queues, and metrics."
          snapshot={snapshot}
        />
      </section>

      <section className="lesson-learn-block">
        <h3>3) Goal + Submit</h3>
        <p className="lesson-outcome">Goal: {attempt?.goal ?? selectedStage?.goal ?? selectedStage?.objective ?? "Pass all checks."}</p>

        {(attempt?.pass_conditions ?? selectedStage?.pass_conditions ?? []).map((item) => (
          <p key={item} className="lesson-outcome">- {item}</p>
        ))}

        {runResult ? (
          <>
            <div className="lesson-summary">
              <span className={runResult.passed ? "badge pass" : "badge fail"}>
                {runResult.passed ? "passed" : "failed"}
              </span>
              <span>feedback: {runResult.feedback_key}</span>
            </div>

            {runResult.validator_results?.map((item) => (
              <p key={`${item.name}.${item.type}`} className="lesson-outcome">
                - {item.passed ? "PASS" : "FAIL"}: {item.name} | expected {item.expected ?? "n/a"}, actual {item.actual ?? "n/a"}
              </p>
            ))}

            {!runResult.passed && runResult.hint ? (
              <p className="hint">
                Hint L{runResult.hint_level ?? 0}: {runResult.hint}
              </p>
            ) : null}
          </>
        ) : (
          <p className="empty">Start challenge actions, then submit to get pass/fail results.</p>
        )}
      </section>

      {errorMessage ? <p className="error">{errorMessage}</p> : null}
      {liveError ? <p className="error">{liveError}</p> : null}
    </section>
  );
}

function ActionControls({
  canSend,
  isAllowed,
  onCommand,
}: {
  canSend: boolean;
  isAllowed: (name: Command["name"]) => boolean;
  onCommand: (command: Command) => void;
}) {
  return (
    <div className="control-row">
      <button type="button" className="btn btn-secondary" disabled={!canSend || !isAllowed("run")} onClick={() => onCommand({ name: "run", count: 8 })}>Run 8</button>
      <button type="button" className="btn btn-secondary" disabled={!canSend || !isAllowed("step")} onClick={() => onCommand({ name: "step", count: 1 })}>Step</button>
      <button type="button" className="btn btn-ghost" disabled={!canSend || !isAllowed("pause")} onClick={() => onCommand({ name: "pause" })}>Pause</button>
      <button type="button" className="btn btn-danger" disabled={!canSend || !isAllowed("reset")} onClick={() => onCommand({ name: "reset" })}>Reset</button>
    </div>
  );
}
