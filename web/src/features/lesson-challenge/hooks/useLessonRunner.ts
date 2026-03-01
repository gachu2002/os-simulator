import { useMutation, useQuery } from "@tanstack/react-query";
import { useMemo, useState } from "react";

import { fetchCurriculumForLearner } from "../../curriculum/api/curriculumApi";
import { getOrCreateLearnerID } from "../../../lib/learner";
import { actionChallengeV3, startChallenge, submitChallenge } from "../api/challengeApi";
import type { LessonActionOptions } from "../model/actionPresets";
import { useChallengeStore } from "../state/challengeStore";
import { useSessionStore } from "../state/sessionStore";

interface UseLessonRunnerOptions {
  baseURL: string;
  preferredLessonID?: string;
  preferredStageIndex?: number;
}

export function useLessonRunner({
  baseURL,
  preferredLessonID,
  preferredStageIndex,
}: UseLessonRunnerOptions) {
  const snapshot = useSessionStore((state) => state.snapshot);
  const liveError = useSessionStore((state) => state.error);
  const resetSession = useSessionStore((state) => state.reset);
  const onSessionEvent = useSessionStore((state) => state.onEvent);
  const setSessionError = useSessionStore((state) => state.setError);

  const attempt = useChallengeStore((state) => state.attempt);
  const runResult = useChallengeStore((state) => state.runResult);
  const runError = useChallengeStore((state) => state.runError);
  const setAttempt = useChallengeStore((state) => state.setAttempt);
  const setRunResult = useChallengeStore((state) => state.setRunResult);
  const setRunError = useChallengeStore((state) => state.setRunError);
  const clearRunError = useChallengeStore((state) => state.clearRunError);
  const [learnerID] = useState(() => getOrCreateLearnerID());
  const [lastLessonAction, setLastLessonAction] = useState("");

  const lessonsQuery = useQuery({
    queryKey: ["challenges", baseURL, learnerID],
    queryFn: () => fetchCurriculumForLearner(baseURL, learnerID),
  });

  const startChallengeMutation = useMutation({
    mutationFn: ({
      lessonID,
      stageID,
    }: {
      lessonID: string;
      stageID?: string;
    }) => startChallenge(baseURL, lessonID, stageID, learnerID),
  });

  const gradeChallengeMutation = useMutation({
    mutationFn: ({ attemptID }: { attemptID: string }) => submitChallenge(baseURL, attemptID, learnerID),
  });

  const lessons = useMemo(() => {
    const sections = lessonsQuery.data ?? [];
    return sections.flatMap((section) => section.lessons ?? []);
  }, [lessonsQuery.data]);

  const selectedLessonID = useMemo(() => {
    if (lessons.length === 0) {
      return "";
    }
    if (preferredLessonID && lessons.some((lesson) => lesson.id === preferredLessonID)) {
      return preferredLessonID;
    }
    return lessons[0].id;
  }, [lessons, preferredLessonID]);

  const selectedLesson = useMemo(
    () => lessons.find((lesson) => lesson.id === selectedLessonID) ?? null,
    [lessons, selectedLessonID],
  );

  const selectedStageIndex = useMemo(() => {
    if (!selectedLesson) {
      return 0;
    }
    if (
      typeof preferredStageIndex === "number" &&
      selectedLesson.id === preferredLessonID &&
      selectedLesson.stages.some((stage) => stage.index === preferredStageIndex)
    ) {
      return preferredStageIndex;
    }
    return selectedLesson.stages[0]?.index ?? 0;
  }, [preferredLessonID, preferredStageIndex, selectedLesson]);

  const selectedStage = useMemo(() => {
    if (!selectedLesson) {
      return null;
    }
    return selectedLesson.stages.find((stage) => stage.index === selectedStageIndex) ?? null;
  }, [selectedLesson, selectedStageIndex]);

  const errorMessage = useMemo(() => {
    if (runError !== "") {
      return runError;
    }
    return lessonsQuery.error instanceof Error ? lessonsQuery.error.message : "";
  }, [lessonsQuery.error, runError]);

  const canSend = Boolean(attempt?.sessionId);
  const attemptID = attempt?.attemptId ?? "";

  const handleStart = async () => {
    if (!selectedLessonID) {
      return;
    }

    clearRunError();
    setLastLessonAction("");
    setRunResult(null);
    resetSession();

    try {
      const started = await startChallengeMutation.mutateAsync({
        lessonID: selectedLessonID,
        stageID: selectedStage?.id,
      });
      setAttempt(started);
    } catch (err) {
      setRunError(err instanceof Error ? err.message : "failed to start challenge");
    }
  };

  const handleLessonAction = async (action: string, options?: LessonActionOptions) => {
    if (!canSend || attemptID === "") {
      return;
    }
    try {
      const response = await actionChallengeV3(baseURL, {
        attemptID,
        learnerID,
        action,
        count: options?.count,
        process: options?.process,
        program: options?.program,
        policy: options?.policy,
        quantum: options?.quantum,
        frames: options?.frames,
        tlbEntries: options?.tlbEntries,
        diskLatency: options?.diskLatency,
        terminalLatency: options?.terminalLatency,
      });
      setLastLessonAction(action);
      onSessionEvent(response.event);
    } catch (err) {
      setSessionError(err instanceof Error ? err.message : "failed to execute action");
    }
  };

  const handleGrade = async () => {
    if (attemptID === "") {
      return;
    }
    clearRunError();
    try {
      const result = await gradeChallengeMutation.mutateAsync({
        attemptID,
      });
      setRunResult(result);
    } catch (err) {
      setRunError(err instanceof Error ? err.message : "failed to check challenge");
    }
  };

  return {
    lessons,
    selectedLesson,
    selectedStageIndex,
    selectedStage,
    runResult,
    attempt,
    snapshot,
    liveError,
    lastLessonAction,
    canSend,
    errorMessage,
    isLessonsLoading: lessonsQuery.isLoading,
    isStartPending: startChallengeMutation.isPending,
    isGradePending: gradeChallengeMutation.isPending,
    handleStart,
    handleLessonAction,
    handleGrade,
  };
}
