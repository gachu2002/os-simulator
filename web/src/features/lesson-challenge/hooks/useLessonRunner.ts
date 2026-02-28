import { useMutation, useQuery } from "@tanstack/react-query";
import { useEffect, useMemo, useRef, useState } from "react";

import { fetchCurriculumForLearner } from "../../curriculum/api/curriculumApi";
import { getOrCreateLearnerID } from "../../../lib/learner";
import type { Command } from "../../../lib/types";
import { connectSessionSocket, type SessionSocket } from "../../../lib/ws";
import { startChallenge, submitChallenge } from "../api/challengeApi";
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
  const connected = useSessionStore((state) => state.connected);
  const snapshot = useSessionStore((state) => state.snapshot);
  const liveError = useSessionStore((state) => state.error);
  const resetSession = useSessionStore((state) => state.reset);
  const onSocketConnected = useSessionStore((state) => state.onSocketConnected);
  const onSocketDisconnected = useSessionStore((state) => state.onSocketDisconnected);
  const onSessionEvent = useSessionStore((state) => state.onEvent);
  const setSessionError = useSessionStore((state) => state.setError);

  const attempt = useChallengeStore((state) => state.attempt);
  const runResult = useChallengeStore((state) => state.runResult);
  const runError = useChallengeStore((state) => state.runError);
  const policy = useChallengeStore((state) => state.policy);
  const quantum = useChallengeStore((state) => state.quantum);
  const frames = useChallengeStore((state) => state.frames);
  const tlbEntries = useChallengeStore((state) => state.tlbEntries);
  const diskLatency = useChallengeStore((state) => state.diskLatency);
  const terminalLatency = useChallengeStore((state) => state.terminalLatency);
  const setPolicy = useChallengeStore((state) => state.setPolicy);
  const setQuantum = useChallengeStore((state) => state.setQuantum);
  const setFrames = useChallengeStore((state) => state.setFrames);
  const setTLBEntries = useChallengeStore((state) => state.setTLBEntries);
  const setDiskLatency = useChallengeStore((state) => state.setDiskLatency);
  const setTerminalLatency = useChallengeStore((state) => state.setTerminalLatency);
  const setAttempt = useChallengeStore((state) => state.setAttempt);
  const setRunResult = useChallengeStore((state) => state.setRunResult);
  const setRunError = useChallengeStore((state) => state.setRunError);
  const clearRunError = useChallengeStore((state) => state.clearRunError);
  const socketRef = useRef<SessionSocket | null>(null);
  const [learnerID] = useState(() => getOrCreateLearnerID());

  const lessonsQuery = useQuery({
    queryKey: ["challenges", baseURL, learnerID],
    queryFn: () => fetchCurriculumForLearner(baseURL, learnerID),
  });

  const startChallengeMutation = useMutation({
    mutationFn: ({ lessonID, stageIndex }: { lessonID: string; stageIndex: number }) =>
      startChallenge(baseURL, lessonID, stageIndex, learnerID),
  });

  const gradeChallengeMutation = useMutation({
    mutationFn: ({ attemptID }: { attemptID: string }) =>
      submitChallenge(baseURL, attemptID, learnerID),
  });

  useEffect(() => {
    return () => {
      socketRef.current?.close();
      socketRef.current = null;
    };
  }, []);

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
    return (
      selectedLesson.stages.find((stage) => stage.unlocked !== false)?.index ??
      selectedLesson.stages[0]?.index ??
      0
    );
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

  const canSend = Boolean(attempt?.sessionId) && connected;
  const attemptID = attempt?.attemptId ?? "";

  const allowedCommandSet = useMemo(() => {
    return new Set(attempt?.allowedCommands ?? []);
  }, [attempt]);

  const isCommandAllowed = (name: Command["name"]) => {
    return allowedCommandSet.has(name);
  };

  const handleStart = async () => {
    if (!selectedLessonID) {
      return;
    }
    if (selectedStage?.unlocked === false) {
      setRunError("this stage is locked: complete prerequisites first");
      return;
    }

    clearRunError();
    setRunResult(null);
    resetSession();
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
        started.sessionId,
        (event) => {
          onSessionEvent(event);
          onSocketConnected();
        },
        (error) => {
          onSocketDisconnected();
          setSessionError(error.message);
        },
      );
      socketRef.current = socket;
    } catch (err) {
      setRunError(err instanceof Error ? err.message : "failed to start challenge");
    }
  };

  const handleCommand = (command: Command) => {
    if (!canSend || !isCommandAllowed(command.name)) {
      return;
    }
    socketRef.current?.sendCommand(command);
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
    selectedLessonID,
    selectedStageIndex,
    selectedStage,
    runResult,
    attempt,
    policy,
    quantum,
    frames,
    tlbEntries,
    diskLatency,
    terminalLatency,
    snapshot,
    liveError,
    canSend,
    errorMessage,
    isLessonsLoading: lessonsQuery.isLoading,
    isStartPending: startChallengeMutation.isPending,
    isGradePending: gradeChallengeMutation.isPending,
    setPolicy,
    setQuantum,
    setFrames,
    setTLBEntries,
    setDiskLatency,
    setTerminalLatency,
    handleStart,
    handleCommand,
    handleGrade,
    isCommandAllowed,
  };
}
