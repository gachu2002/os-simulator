import { useQuery } from "@tanstack/react-query";
import { useEffect, useMemo, useState } from "react";

import { ControlBar } from "./components/ControlBar";
import { EventLog } from "./components/EventLog";
import { LessonRunnerPanel } from "./components/LessonRunnerPanel";
import { StatusCards } from "./components/StatusCards";
import { VisualizationSuite } from "./components/VisualizationSuite";
import { useSession } from "./hooks/useSession";
import {
  fetchLessonProgress,
  type LessonRunResponse,
} from "./lib/lessonApi";
import { snapshotFromLessonRun } from "./state/lessonSnapshot";

type AppMode = "path" | "sandbox" | "challenge" | "progress";

const MODE_PATHS: Record<AppMode, string> = {
  path: "/path",
  sandbox: "/sandbox",
  challenge: "/challenge",
  progress: "/progress",
};

const MODE_META: Record<AppMode, { eyebrow: string; title: string; subtitle: string }> = {
  path: {
    eyebrow: "Guided Path",
    title: "Learn OSTEP Step by Step",
    subtitle:
      "Work through sequenced lessons, verify outcomes, and compare replay snapshots.",
  },
  sandbox: {
    eyebrow: "Sandbox",
    title: "Experiment Freely",
    subtitle:
      "Control a live deterministic session and inspect traces, queues, and memory behavior.",
  },
  challenge: {
    eyebrow: "Challenge",
    title: "Apply What You Learned",
    subtitle:
      "Run a live session while solving lesson goals under tighter constraints and less guidance.",
  },
  progress: {
    eyebrow: "Progress",
    title: "Track Mastery",
    subtitle:
      "Review your latest lesson outcomes and plan the next modules to strengthen.",
  },
};

function normalizeMode(pathname: string): AppMode {
  switch (pathname) {
    case "/sandbox":
      return "sandbox";
    case "/challenge":
      return "challenge";
    case "/progress":
      return "progress";
    case "/path":
    case "/":
    default:
      return "path";
  }
}

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
  const [mode, setMode] = useState<AppMode>(() =>
    normalizeMode(window.location.pathname),
  );
  const progressQuery = useQuery({
    queryKey: ["lesson-progress", baseURL],
    queryFn: () => fetchLessonProgress(baseURL),
    enabled: mode === "progress",
    refetchInterval: mode === "progress" ? 3000 : false,
  });

  useEffect(() => {
    function handlePopState() {
      setMode(normalizeMode(window.location.pathname));
    }

    window.addEventListener("popstate", handlePopState);
    return () => window.removeEventListener("popstate", handlePopState);
  }, []);

  const modeMeta = useMemo(() => MODE_META[mode], [mode]);

  function navigate(nextMode: AppMode) {
    const nextPath = MODE_PATHS[nextMode];
    if (window.location.pathname !== nextPath) {
      window.history.pushState({}, "", nextPath);
    }
    setMode(nextMode);
  }

  function handleLessonRunResult(result: LessonRunResponse) {
    setLessonSnapshot(snapshotFromLessonRun(result));
    setLessonTitle(`${result.lesson_id} stage ${result.stage_index}`);
  }

  return (
    <main className="app-shell">
      <header className="hero">
        <p className="eyebrow">{modeMeta.eyebrow}</p>
        <h1>{modeMeta.title}</h1>
        <p className="subtitle">{modeMeta.subtitle}</p>
      </header>

      <nav className="panel mode-nav" aria-label="Learning modes">
        {(Object.keys(MODE_PATHS) as AppMode[]).map((nextMode) => (
          <a
            key={nextMode}
            className={mode === nextMode ? "mode-link active" : "mode-link"}
            href={MODE_PATHS[nextMode]}
            onClick={(event) => {
              event.preventDefault();
              navigate(nextMode);
            }}
          >
            {nextMode}
          </a>
        ))}
      </nav>

      {mode === "path" ? (
        <>
          <LessonRunnerPanel
            baseURL={baseURL}
            onRunResult={handleLessonRunResult}
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
              compareMode && lessonSnapshot
                ? "viz-compare-grid"
                : "viz-single-grid"
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
        </>
      ) : null}

      {mode === "sandbox" || mode === "challenge" ? (
        <>
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

          {mode === "challenge" ? (
            <LessonRunnerPanel
              baseURL={baseURL}
              onRunResult={handleLessonRunResult}
            />
          ) : null}

          <VisualizationSuite
            title="Live Session Visualizations"
            subtitle="Realtime transport stream from active session controls"
            snapshot={state.snapshot}
          />

          <EventLog logs={state.logs} />
        </>
      ) : null}

      {mode === "progress" ? (
        <section className="panel progress-panel">
          <h2>Learning Progress Snapshot</h2>
          {progressQuery.error instanceof Error ? (
            <p className="error">{progressQuery.error.message}</p>
          ) : null}
          {progressQuery.data ? (
            <>
              <ul>
                <li>
                  Completed stages: {progressQuery.data.analytics.completed_stages}/
                  {progressQuery.data.analytics.total_stages}
                </li>
                <li>
                  Attempt coverage: {Math.round(progressQuery.data.analytics.attempt_coverage * 100)}%
                </li>
                <li>
                  Pilot checklist: {progressQuery.data.analytics.pilot_checklist_ok ? "ready" : "in progress"}
                </li>
                {lessonSnapshot ? <li>Latest lesson run: {lessonTitle}</li> : null}
              </ul>

              <h3>Weak concepts</h3>
              {progressQuery.data.analytics.weak_concepts.length > 0 ? (
                <ul>
                  {progressQuery.data.analytics.weak_concepts.map((item) => (
                    <li key={item.concept}>
                      {item.concept}: score {item.score.toFixed(1)} (fails {item.failed_attempts}, hints {item.high_hint_uses})
                    </li>
                  ))}
                </ul>
              ) : (
                <p className="empty">No weak concepts detected yet.</p>
              )}
            </>
          ) : (
            <p className="empty">Run lessons to populate progress analytics.</p>
          )}
        </section>
      ) : null}
    </main>
  );
}
