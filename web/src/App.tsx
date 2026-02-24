import { useEffect, useMemo, useReducer, useRef, useState } from "react";

import { ControlBar } from "./components/ControlBar";
import { EventLog } from "./components/EventLog";
import { LessonRunnerPanel } from "./components/LessonRunnerPanel";
import { StatusCards } from "./components/StatusCards";
import { VisualizationSuite } from "./components/VisualizationSuite";
import { createSession } from "./lib/api";
import type { LessonRunResponse } from "./lib/lessonApi";
import type { Command } from "./lib/types";
import { connectSessionSocket, type SessionSocket } from "./lib/ws";
import { snapshotFromLessonRun } from "./state/lessonSnapshot";
import { initialSessionState, sessionReducer } from "./state/sessionReducer";

export function App() {
  const [state, dispatch] = useReducer(sessionReducer, initialSessionState);
  const [baseURL, setBaseURL] = useState(defaultBaseURL());
  const [seed, setSeed] = useState(1);
  const [policy, setPolicy] = useState<"fifo" | "rr" | "mlfq">("rr");
  const [quantum, setQuantum] = useState(2);
  const [isCreating, setIsCreating] = useState(false);
  const [lessonSnapshot, setLessonSnapshot] = useState<ReturnType<typeof snapshotFromLessonRun> | null>(null);
  const [lessonTitle, setLessonTitle] = useState("");
  const [compareMode, setCompareMode] = useState(true);
  const socketRef = useRef<SessionSocket | null>(null);

  const canSend = useMemo(
    () => Boolean(state.sessionID && state.connected),
    [state.connected, state.sessionID],
  );

  useEffect(() => {
    return () => {
      socketRef.current?.close();
      socketRef.current = null;
    };
  }, []);

  async function handleCreateSession() {
    setIsCreating(true);
    dispatch({ type: "error", message: "" });
    try {
      socketRef.current?.close();
      socketRef.current = null;

      const created = await createSession(baseURL, {
        seed,
        policy,
        quantum,
      });

      dispatch({
        type: "session.created",
        sessionID: created.session_id,
        snapshot: created.snapshot,
      });

      const socket = connectSessionSocket(
        baseURL,
        created.session_id,
        (event) => {
          dispatch({ type: "event.received", event });
          dispatch({ type: "socket.connected" });
        },
        (error) => {
          dispatch({ type: "socket.disconnected" });
          dispatch({ type: "error", message: error.message });
        },
      );

      socketRef.current = socket;
    } catch (error) {
      dispatch({
        type: "error",
        message: error instanceof Error ? error.message : "failed to create session",
      });
    } finally {
      setIsCreating(false);
    }
  }

  function handleCommand(command: Command) {
    if (!canSend) {
      return;
    }
    socketRef.current?.sendCommand(command);
  }

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
          <button type="button" disabled={isCreating} onClick={handleCreateSession}>
            {isCreating ? "Creating..." : "Create Session"}
          </button>
        </div>
        {state.error ? <p className="error">{state.error}</p> : null}
      </section>

      <LessonRunnerPanel baseURL={baseURL} onRunResult={handleLessonRunResult} />

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

      <div className={compareMode && lessonSnapshot ? "viz-compare-grid" : "viz-single-grid"}>
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

function defaultBaseURL(): string {
  if (typeof window === "undefined") {
    return "http://localhost:8080";
  }
  const host = window.location.hostname || "localhost";
  return `http://${host}:8080`;
}
