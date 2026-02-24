import { useEffect, useMemo, useState } from "react";

import {
  fetchLessons,
  runLessonStage,
  type LessonRunResponse,
  type LessonSummary,
} from "../lib/lessonApi";

interface LessonRunnerPanelProps {
  baseURL: string;
  onRunResult?: (result: LessonRunResponse) => void;
}

export function LessonRunnerPanel({ baseURL, onRunResult }: LessonRunnerPanelProps) {
  const [lessons, setLessons] = useState<LessonSummary[]>([]);
  const [selectedLessonID, setSelectedLessonID] = useState("");
  const [selectedStageIndex, setSelectedStageIndex] = useState(0);
  const [runResult, setRunResult] = useState<LessonRunResponse | null>(null);
  const [lastTraceHash, setLastTraceHash] = useState("");
  const [determinismStatus, setDeterminismStatus] = useState<"" | "stable" | "changed">("");
  const [loading, setLoading] = useState(false);
  const [running, setRunning] = useState(false);
  const [error, setError] = useState("");

  const selectedLesson = useMemo(
    () => lessons.find((lesson) => lesson.id === selectedLessonID) ?? null,
    [lessons, selectedLessonID],
  );

  useEffect(() => {
    let active = true;
    setLoading(true);
    setError("");
    fetchLessons(baseURL)
      .then((loaded) => {
        if (!active) {
          return;
        }
        setLessons(loaded);
        if (loaded.length > 0) {
          setSelectedLessonID(loaded[0].id);
          setSelectedStageIndex(loaded[0].stages[0]?.index ?? 0);
        }
      })
      .catch((err: unknown) => {
        if (!active) {
          return;
        }
        setError(err instanceof Error ? err.message : "failed to load lessons");
      })
      .finally(() => {
        if (active) {
          setLoading(false);
        }
      });

    return () => {
      active = false;
    };
  }, [baseURL]);

  async function handleRun() {
    if (!selectedLessonID) {
      return;
    }
    setRunning(true);
    setError("");
    try {
      const result = await runLessonStage(baseURL, selectedLessonID, selectedStageIndex);
      if (lastTraceHash !== "") {
        setDeterminismStatus(lastTraceHash === result.output.trace_hash ? "stable" : "changed");
      }
      setLastTraceHash(result.output.trace_hash);
      setRunResult(result);
      onRunResult?.(result);
    } catch (err) {
      setError(err instanceof Error ? err.message : "failed to run lesson stage");
    } finally {
      setRunning(false);
    }
  }

  return (
    <section className="panel lesson-panel">
      <h2>Lesson Runner</h2>
      <div className="lesson-controls">
        <label>
          Lesson
          <select
            value={selectedLessonID}
            disabled={loading || lessons.length === 0}
            onChange={(event) => {
              const nextID = event.target.value;
              setSelectedLessonID(nextID);
              const stage = lessons.find((l) => l.id === nextID)?.stages[0];
              setSelectedStageIndex(stage?.index ?? 0);
            }}
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
            disabled={!selectedLesson}
            onChange={(event) => setSelectedStageIndex(Number(event.target.value))}
          >
            {(selectedLesson?.stages ?? []).map((stage) => (
              <option key={stage.id} value={stage.index}>
                {stage.title}
              </option>
            ))}
          </select>
        </label>

        <button type="button" disabled={running || !selectedLessonID} onClick={handleRun}>
          {running ? "Running..." : "Run Stage"}
        </button>
      </div>

      {error ? <p className="error">{error}</p> : null}

      {runResult ? (
        <>
          <div className="lesson-summary">
            <span className={runResult.passed ? "badge pass" : "badge fail"}>
              {runResult.passed ? "passed" : "failed"}
            </span>
            <span>feedback: {runResult.feedback_key}</span>
            <span>trace hash: {runResult.output.trace_hash}</span>
            <span>trace length: {runResult.output.trace_length}</span>
            <span>fs: {runResult.output.filesystem_ok ? "ok" : "failed"}</span>
          </div>

          {!runResult.passed && runResult.hint ? (
            <p className="hint">
              Hint L{runResult.hint_level ?? 0}: {runResult.hint}
            </p>
          ) : null}

          {determinismStatus ? (
            <p className="determinism">
              Determinism check: {determinismStatus === "stable" ? "stable hash" : "hash changed"}
            </p>
          ) : null}

          <div className="analytics-grid">
            <article>
              <h3>Completion</h3>
              <p>
                {runResult.analytics.completed_stages}/{runResult.analytics.total_stages} (
                {formatPercent(runResult.analytics.completion_rate)})
              </p>
              <p>
                coverage: {runResult.analytics.attempted_stages}/{runResult.analytics.total_stages} (
                {formatPercent(runResult.analytics.attempt_coverage)})
              </p>
              <p>pilot checklist: {runResult.analytics.pilot_checklist_ok ? "ready" : "in progress"}</p>
            </article>

            <article>
              <h3>Modules</h3>
              <ul>
                {runResult.analytics.module_breakdown.map((mod) => (
                  <li key={mod.module}>
                    {mod.module}: {mod.completed_stage}/{mod.total_stages} ({formatPercent(mod.completion_rate)})
                  </li>
                ))}
              </ul>
            </article>
          </div>
        </>
      ) : (
        <p className="empty">Run a lesson stage to view grading, hints, and analytics.</p>
      )}
    </section>
  );
}

function formatPercent(value: number): string {
  return `${Math.round(value * 100)}%`;
}
