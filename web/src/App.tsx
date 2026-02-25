import { useState } from "react";

import { ControlBar } from "./components/ControlBar";
import { EventLog } from "./components/EventLog";
import { LessonRunnerPanel } from "./components/LessonRunnerPanel";
import { StatusCards } from "./components/StatusCards";
import { VisualizationSuite } from "./components/VisualizationSuite";
import { useSession } from "./hooks/useSession";
import type { LessonRunResponse } from "./lib/lessonApi";
import { snapshotFromLessonRun } from "./state/lessonSnapshot";

export function App() {
  const {
    state,
    baseURL,
    seed,
    policy,
    quantum,
    canSend,
    isCreatingSession,
    setBaseURL,
    setSeed,
    setPolicy,
    setQuantum,
    handleCreateSession,
    handleCommand,
  } = useSession();
  const [lessonSnapshot, setLessonSnapshot] = useState<ReturnType<
    typeof snapshotFromLessonRun
  > | null>(null);
  const [lessonTitle, setLessonTitle] = useState("");
  const [compareMode, setCompareMode] = useState(true);

  function handleLessonRunResult(result: LessonRunResponse) {
    setLessonSnapshot(snapshotFromLessonRun(result));
    setLessonTitle(`${result.lesson_id} stage ${result.stage_index}`);
  }

  return (
    <main className="app-shell">
      <header className="hero">
        <p className="eyebrow">Milestone 14 Lesson Integration</p>
        <h1>OS Simulator Console</h1>
        <p className="subtitle">
          Simple, modern, classic control surface for deterministic sessions.
        </p>
      </header>

      <section className="panel session-panel">
        <div className="session-form">
          <label>
            Server URL
            <input
              value={baseURL}
              onChange={(event) => setBaseURL(event.target.value)}
              placeholder="http://localhost:8080"
            />
          </label>
          <label>
            Seed
            <input
              type="number"
              min={1}
              value={seed}
              onChange={(event) => setSeed(Number(event.target.value))}
            />
          </label>
          <button
            type="button"
            disabled={isCreatingSession}
            onClick={handleCreateSession}
          >
            {isCreatingSession ? "Creating..." : "Create Session"}
          </button>
        </div>
        {state.error ? <p className="error">{state.error}</p> : null}
      </section>

      <LessonRunnerPanel
        baseURL={baseURL}
        onRunResult={handleLessonRunResult}
      />

      <StatusCards
        connected={state.connected}
        sessionID={state.sessionID}
        snapshot={state.snapshot}
      />

      <ControlBar
        policy={policy}
        quantum={quantum}
        disabled={!canSend}
        onPolicyChange={setPolicy}
        onQuantumChange={setQuantum}
        onCommand={handleCommand}
      />

      <section className="panel compare-panel">
        <label>
          <input
            type="checkbox"
            checked={compareMode}
            onChange={(event) => setCompareMode(event.target.checked)}
          />
          Side-by-side lesson replay snapshot mode
        </label>
      </section>

      <div
        className={
          compareMode && lessonSnapshot ? "viz-compare-grid" : "viz-single-grid"
        }
      >
        <VisualizationSuite
          title="Live Session Visualizations"
          subtitle="Realtime transport stream from active session controls"
          snapshot={state.snapshot}
        />
        {compareMode && lessonSnapshot ? (
          <VisualizationSuite
            title="Lesson Replay Snapshot"
            subtitle={`Latest lesson run: ${lessonTitle}`}
            snapshot={lessonSnapshot}
          />
        ) : null}
      </div>

      <EventLog logs={state.logs} />
    </main>
  );
}
