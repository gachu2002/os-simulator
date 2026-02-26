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
    errorMessage,
    isLessonsLoading,
    isRunPending,
    setSelectedStageIndexState,
    handleLessonChange,
    handleRun,
  } = useLessonRunner({ baseURL, onRunResult });

  return (
    <section className="panel lesson-panel">
      <h2>Challenge Runner</h2>
      <div className="lesson-controls">
        <label>
          Challenge
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
          Step
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
          {isRunPending ? "Running..." : "Run Step"}
        </button>
      </div>

      {errorMessage ? <p className="error">{errorMessage}</p> : null}

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
      ) : (
        <p className="empty">
          Pick a challenge and run a step to get grading feedback and hints.
        </p>
      )}
    </section>
  );
}

function formatPercent(value: number): string {
  return `${Math.round(value * 100)}%`;
}
