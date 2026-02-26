import { useCallback, useEffect, useMemo, useState } from "react";

import { ControlBar } from "./components/ControlBar";
import { EventLog } from "./components/EventLog";
import { LessonRunnerPanel } from "./components/LessonRunnerPanel";
import { StatusCards } from "./components/StatusCards";
import { VisualizationSuite } from "./components/VisualizationSuite";
import { useSession } from "./hooks/useSession";
import type { ChallengeGradeResponse } from "./lib/lessonApi";
import type { SnapshotDTO } from "./lib/types";
import { snapshotFromChallengeGrade } from "./state/lessonSnapshot";

type AppMode = "sandbox" | "challenge";

const MODE_PATHS: Record<AppMode, string> = {
  sandbox: "/sandbox",
  challenge: "/challenge",
};

const MODE_META: Record<AppMode, { eyebrow: string; title: string; subtitle: string }> = {
  sandbox: {
    eyebrow: "Sandbox",
    title: "Experiment with OSTEP Concepts",
    subtitle:
      "Build intuition by running your own deterministic workloads and inspecting what changes.",
  },
  challenge: {
    eyebrow: "Challenge",
    title: "Solve Small OSTEP Steps",
    subtitle:
      "Complete focused challenge steps with hints and deterministic grading feedback.",
  },
};

function normalizeMode(pathname: string): AppMode {
  switch (pathname) {
    case "/challenge":
      return "challenge";
    case "/sandbox":
    case "/":
    default:
      return "sandbox";
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
  const [lessonSnapshot, setLessonSnapshot] = useState<SnapshotDTO | null>(null);
  const [lessonTitle, setLessonTitle] = useState("");
  const [hasGradedAttempt, setHasGradedAttempt] = useState(false);
  const [mode, setMode] = useState<AppMode>(() =>
    normalizeMode(window.location.pathname),
  );

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

  const handleChallengeLiveSnapshot = useCallback((snapshot: SnapshotDTO | null, title: string) => {
    setLessonSnapshot(snapshot);
    setLessonTitle(title);
    setHasGradedAttempt(false);
  }, []);

  const handleChallengeGradeResult = useCallback((result: ChallengeGradeResponse) => {
    setLessonSnapshot(snapshotFromChallengeGrade(result));
    setLessonTitle(`${result.lesson_id} step ${result.stage_index + 1}`);
    setHasGradedAttempt(true);
  }, []);

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

      {mode === "sandbox" ? (
        <>
          <section className="panel session-panel">
            <h2>Session Setup</h2>
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

          <VisualizationSuite
            title="Live Session Visualizations"
            subtitle="Realtime transport stream from active session controls"
            snapshot={state.snapshot}
          />

          <EventLog logs={state.logs} />
        </>
      ) : null}

      {mode === "challenge" ? (
        <>
          <section className="panel challenge-setup-panel">
            <h2>Challenge Setup</h2>
            <label>
              Server URL
              <input
                value={baseURL}
                onChange={(event) => setBaseURL(event.target.value)}
                placeholder="http://127.0.0.1:8080"
              />
            </label>
            <p className="challenge-note">
              Each challenge step is short and focused. Run the step, inspect
              the result, and use the hint to refine your reasoning.
            </p>
          </section>

          <LessonRunnerPanel
            baseURL={baseURL}
            onLiveSnapshot={handleChallengeLiveSnapshot}
            onGradeResult={handleChallengeGradeResult}
          />

          <VisualizationSuite
            title="Challenge Step Snapshot"
            subtitle={
              lessonSnapshot
                ? hasGradedAttempt
                  ? `Latest graded result: ${lessonTitle}`
                  : `Live challenge attempt: ${lessonTitle}`
                : "Start a challenge to render scheduler, memory, and process views"
            }
            snapshot={lessonSnapshot}
          />
        </>
      ) : null}
    </main>
  );
}
