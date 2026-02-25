import { useLessonRunner } from "../hooks/useLessonRunner";
import type { LessonRunResponse } from "../lib/lessonApi";

interface LessonRunnerPanelProps {
  baseURL: string;
  onRunResult?: (result: LessonRunResponse) => void;
}

export function LessonRunnerPanel({
  baseURL,
  onRunResult,
}: LessonRunnerPanelProps) {
  const {
    lessons,
    selectedLesson,
    selectedLessonID,
    selectedStageIndex,
    runResult,
    determinismStatus,
    errorMessage,
    isLessonsLoading,
    isRunPending,
    setSelectedStageIndexState,
    handleLessonChange,
    handleRun,
  } = useLessonRunner({ baseURL, onRunResult });

  return (
    <section className="panel lesson-panel">
      <h2>Lesson Runner</h2>
      <div className="lesson-controls">
        <label>
          Lesson
          <select
            value={selectedLessonID}
            disabled={isLessonsLoading || lessons.length === 0}
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
            disabled={!selectedLesson}
            onChange={(event) =>
              setSelectedStageIndexState(Number(event.target.value))
            }
          >
            {(selectedLesson?.stages ?? []).map((stage) => (
              <option key={stage.id} value={stage.index}>
                {stage.title}
              </option>
            ))}
          </select>
        </label>

        <button
          type="button"
          disabled={isRunPending || !selectedLessonID}
          onClick={handleRun}
        >
          {isRunPending ? "Running..." : "Run Stage"}
        </button>
      </div>

      {errorMessage ? <p className="error">{errorMessage}</p> : null}

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
              Determinism check:{" "}
              {determinismStatus === "stable" ? "stable hash" : "hash changed"}
            </p>
          ) : null}

          <div className="analytics-grid">
            <article>
              <h3>Completion</h3>
              <p>
                {runResult.analytics.completed_stages}/
                {runResult.analytics.total_stages} (
                {formatPercent(runResult.analytics.completion_rate)})
              </p>
              <p>
                coverage: {runResult.analytics.attempted_stages}/
                {runResult.analytics.total_stages} (
                {formatPercent(runResult.analytics.attempt_coverage)})
              </p>
              <p>
                pilot checklist:{" "}
                {runResult.analytics.pilot_checklist_ok
                  ? "ready"
                  : "in progress"}
              </p>
            </article>

            <article>
              <h3>Modules</h3>
              <ul>
                {runResult.analytics.module_breakdown.map((mod) => (
                  <li key={mod.module}>
                    {mod.module}: {mod.completed_stage}/{mod.total_stages} (
                    {formatPercent(mod.completion_rate)})
                  </li>
                ))}
              </ul>
            </article>
          </div>
        </>
      ) : (
        <p className="empty">
          Run a lesson stage to view grading, hints, and analytics.
        </p>
      )}
    </section>
  );
}

function formatPercent(value: number): string {
  return `${Math.round(value * 100)}%`;
}
