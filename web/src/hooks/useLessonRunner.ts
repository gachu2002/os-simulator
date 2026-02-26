import { useMutation, useQuery } from "@tanstack/react-query";
import { useCallback, useEffect, useMemo, useReducer, useRef, useState } from "react";

import {
  fetchLessonsForLearner,
  gradeChallenge,
  startChallenge,
  type ChallengeGradeResponse,
  type ChallengeStartResponse,
} from "../lib/lessonApi";
import { getOrCreateLearnerID } from "../lib/learner";
import type { Command } from "../lib/types";
import { connectSessionSocket, type SessionSocket } from "../lib/ws";
import { initialSessionState, sessionReducer } from "../state/sessionReducer";

interface UseLessonRunnerOptions {
  baseURL: string;
  onGradeResult?: (result: ChallengeGradeResponse) => void;
  preferredLessonID?: string;
  preferredStageIndex?: number;
}

export function useLessonRunner({
  baseURL,
  onGradeResult,
  preferredLessonID,
  preferredStageIndex,
}: UseLessonRunnerOptions) {
  const [liveState, dispatch] = useReducer(sessionReducer, initialSessionState);
  const [selectedLessonIDState, setSelectedLessonIDState] = useState("");
  const [selectedStageIndexState, setSelectedStageIndexState] = useState(0);
  const [runResult, setRunResult] = useState<ChallengeGradeResponse | null>(null);
  const [runError, setRunError] = useState("");
  const [policy, setPolicy] = useState<"fifo" | "rr" | "mlfq">("rr");
  const [quantum, setQuantum] = useState(2);
  const socketRef = useRef<SessionSocket | null>(null);
  const [learnerID] = useState(() => getOrCreateLearnerID());

  const lessonsQuery = useQuery({
    queryKey: ["challenges", baseURL, learnerID],
    queryFn: () => fetchLessonsForLearner(baseURL, learnerID),
  });

  const startChallengeMutation = useMutation({
    mutationFn: ({ lessonID, stageIndex }: { lessonID: string; stageIndex: number }) =>
      startChallenge(baseURL, lessonID, stageIndex, learnerID),
  });

  const gradeChallengeMutation = useMutation({
    mutationFn: ({ attemptID }: { attemptID: string }) =>
      gradeChallenge(baseURL, attemptID, learnerID),
  });

  const [attempt, setAttempt] = useState<ChallengeStartResponse | null>(null);

  useEffect(() => {
    return () => {
      socketRef.current?.close();
      socketRef.current = null;
    };
  }, []);

  const lessons = useMemo(() => lessonsQuery.data ?? [], [lessonsQuery.data]);

  const selectedLessonID = useMemo(() => {
    if (lessons.length === 0) {
      return "";
    }
    if (preferredLessonID && lessons.some((lesson) => lesson.id === preferredLessonID)) {
      return preferredLessonID;
    }
    if (lessons.some((lesson) => lesson.id === selectedLessonIDState)) {
      return selectedLessonIDState;
    }
    return lessons[0].id;
  }, [lessons, preferredLessonID, selectedLessonIDState]);

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
    if (
      selectedLesson.stages.some((stage) => stage.index === selectedStageIndexState)
    ) {
      return selectedStageIndexState;
    }
    return selectedLesson.stages[0]?.index ?? 0;
  }, [preferredLessonID, preferredStageIndex, selectedLesson, selectedStageIndexState]);

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

  const canSend = Boolean(attempt?.session_id) && liveState.connected;
  const attemptID = attempt?.attempt_id ?? "";

  const allowedCommandSet = useMemo(() => {
    return new Set(attempt?.allowed_commands ?? []);
  }, [attempt]);

  const isCommandAllowed = useCallback(
    (name: Command["name"]) => {
      return allowedCommandSet.has(name);
    },
    [allowedCommandSet],
  );

  const handleLessonChange = useCallback(
    (lessonID: string) => {
      setSelectedLessonIDState(lessonID);
      const lesson = lessons.find((item) => item.id === lessonID);
      const stage = lesson?.stages.find((item) => item.unlocked !== false) ?? lesson?.stages[0];
      setSelectedStageIndexState(stage?.index ?? 0);
    },
    [lessons],
  );

  const handleStart = useCallback(async () => {
    if (!selectedLessonID) {
      return;
    }
    if (selectedStage?.unlocked === false) {
      setRunError("this stage is locked: complete prerequisites first");
      return;
    }

    setRunError("");
    setRunResult(null);
    dispatch({ type: "session.reset" });
    socketRef.current?.close();
    socketRef.current = null;

    try {
      const started = await startChallengeMutation.mutateAsync({
        lessonID: selectedLessonID,
        stageIndex: selectedStageIndex,
      });
      setAttempt(started);

      const socket = connectSessionSocket(
        baseURL,
        started.session_id,
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
    } catch (err) {
      setRunError(err instanceof Error ? err.message : "failed to start challenge");
    }
  }, [
    baseURL,
    selectedLessonID,
    selectedStageIndex,
    selectedStage,
    startChallengeMutation,
  ]);

  const handleCommand = useCallback(
    (command: Command) => {
      if (!canSend || !isCommandAllowed(command.name)) {
        return;
      }
      socketRef.current?.sendCommand(command);
    },
    [canSend, isCommandAllowed],
  );

  const handleGrade = useCallback(async () => {
    if (attemptID === "") {
      return;
    }
    setRunError("");
    try {
      const result = await gradeChallengeMutation.mutateAsync({
        attemptID,
      });
      setRunResult(result);
      onGradeResult?.(result);
    } catch (err) {
      setRunError(err instanceof Error ? err.message : "failed to check challenge");
    }
  }, [attemptID, gradeChallengeMutation, onGradeResult]);

  return {
    lessons,
    selectedLesson,
    selectedLessonID,
    selectedStageIndex,
    selectedStage,
    runResult,
    attempt,
    policy,
    quantum,
    snapshot: liveState.snapshot,
    liveError: liveState.error,
    canSend,
    errorMessage,
    isLessonsLoading: lessonsQuery.isLoading,
    isStartPending: startChallengeMutation.isPending,
    isGradePending: gradeChallengeMutation.isPending,
    setPolicy,
    setQuantum,
    setSelectedStageIndexState,
    handleLessonChange,
    handleStart,
    handleCommand,
    handleGrade,
    isCommandAllowed,
  };
}
