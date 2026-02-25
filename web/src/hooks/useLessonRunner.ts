import { useMutation, useQuery } from "@tanstack/react-query";
import { useCallback, useMemo, useState } from "react";

import {
  fetchLessons,
  runLessonStage,
  type LessonRunResponse,
} from "../lib/lessonApi";

interface UseLessonRunnerOptions {
  baseURL: string;
  onRunResult?: (result: LessonRunResponse) => void;
}

export function useLessonRunner({
  baseURL,
  onRunResult,
}: UseLessonRunnerOptions) {
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

  const runStageMutation = useMutation({
    mutationFn: ({
      lessonID,
      stageIndex,
    }: {
      lessonID: string;
      stageIndex: number;
    }) => runLessonStage(baseURL, lessonID, stageIndex),
  });

  const lessons = useMemo(() => lessonsQuery.data ?? [], [lessonsQuery.data]);

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
      selectedLesson.stages.some((stage) => stage.index === selectedStageIndexState)
    ) {
      return selectedStageIndexState;
    }
    return selectedLesson.stages[0]?.index ?? 0;
  }, [selectedLesson, selectedStageIndexState]);

  const errorMessage = useMemo(() => {
    if (runError !== "") {
      return runError;
    }
    return lessonsQuery.error instanceof Error ? lessonsQuery.error.message : "";
  }, [lessonsQuery.error, runError]);

  const handleLessonChange = useCallback(
    (lessonID: string) => {
      setSelectedLessonIDState(lessonID);
      const stage = lessons.find((lesson) => lesson.id === lessonID)?.stages[0];
      setSelectedStageIndexState(stage?.index ?? 0);
    },
    [lessons],
  );

  const handleRun = useCallback(async () => {
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
      setRunError(err instanceof Error ? err.message : "failed to run lesson stage");
    }
  }, [
    lastTraceHash,
    onRunResult,
    runStageMutation,
    selectedLessonID,
    selectedStageIndex,
  ]);

  return {
    lessons,
    selectedLesson,
    selectedLessonID,
    selectedStageIndex,
    runResult,
    determinismStatus,
    errorMessage,
    isLessonsLoading: lessonsQuery.isLoading,
    isRunPending: runStageMutation.isPending,
    setSelectedStageIndexState,
    handleLessonChange,
    handleRun,
  };
}
