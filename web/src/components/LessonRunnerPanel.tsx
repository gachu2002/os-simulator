import { useMutation, useQuery } from "@tanstack/react-query";
import { useMemo, useState } from "react";

import {
  fetchLessons,
  runLessonStage,
  type LessonRunResponse,
} from "../lib/lessonApi";

interface LessonRunnerPanelProps {
  baseURL: string;
  onRunResult?: (result: LessonRunResponse) => void;
}

export function LessonRunnerPanel({
  baseURL,
  onRunResult,
}: LessonRunnerPanelProps) {
  const [selectedLessonIDState, setSelectedLessonIDState] = useState("");
  const [selectedStageIndexState, setSelectedStageIndexState] = useState(0);
  const [runResult, setRunResult] = useState<LessonRunResponse | null>(null);
  const [lastTraceHash, setLastTraceHash] = useState("");
  const [determinismStatus, setDeterminismStatus] = useState<
    "" | "stable" | "changed"
  >("");
  const [runError, setRunError] = useState("");

  const lessonsQuery = useQuery({
    queryKey: ["lessons", baseURL],
    queryFn: () => fetchLessons(baseURL),
  });

  const lessons = useMemo(() => lessonsQuery.data ?? [], [lessonsQuery.data]);

  const runStageMutation = useMutation({
    mutationFn: ({
      lessonID,
      stageIndex,
    }: {
      lessonID: string;
      stageIndex: number;
    }) => runLessonStage(baseURL, lessonID, stageIndex),
  });

  const selectedLessonID = useMemo(() => {
    if (lessons.length === 0) {
      return "";
    }
    if (lessons.some((lesson) => lesson.id === selectedLessonIDState)) {
      return selectedLessonIDState;
    }
    return lessons[0].id;
  }, [lessons, selectedLessonIDState]);

  const selectedLesson = useMemo(
    () => lessons.find((lesson) => lesson.id === selectedLessonID) ?? null,
    [lessons, selectedLessonID],
  );

  const selectedStageIndex = useMemo(() => {
    if (!selectedLesson) {
      return 0;
    }
    if (
      selectedLesson.stages.some(
        (stage) => stage.index === selectedStageIndexState,
      )
    ) {
      return selectedStageIndexState;
    }
    return selectedLesson.stages[0]?.index ?? 0;
  }, [selectedLesson, selectedStageIndexState]);

  const errorMessage = useMemo(() => {
    if (runError !== "") {
      return runError;
    }
    return lessonsQuery.error instanceof Error
      ? lessonsQuery.error.message
      : "";
  }, [lessonsQuery.error, runError]);

  async function handleRun() {
    if (!selectedLessonID) {
      return;
    }
    setRunError("");
    try {
      const result = await runStageMutation.mutateAsync({
        lessonID: selectedLessonID,
        stageIndex: selectedStageIndex,
      });
      if (lastTraceHash !== "") {
        setDeterminismStatus(
          lastTraceHash === result.output.trace_hash ? "stable" : "changed",
        );
      }
      setLastTraceHash(result.output.trace_hash);
      setRunResult(result);
      onRunResult?.(result);
    } catch (err) {
      setRunError(
        err instanceof Error ? err.message : "failed to run lesson stage",
      );
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
            disabled={lessonsQuery.isLoading || lessons.length === 0}
            onChange={(event) => {
              const nextID = event.target.value;
              setSelectedLessonIDState(nextID);
              const stage = lessons.find((l) => l.id === nextID)?.stages[0];
              setSelectedStageIndexState(stage?.index ?? 0);
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
          disabled={runStageMutation.isPending || !selectedLessonID}
          onClick={handleRun}
        >
          {runStageMutation.isPending ? "Running..." : "Run Stage"}
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
