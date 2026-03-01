import { useQuery } from "@tanstack/react-query";
import { useEffect, useMemo, useState } from "react";

import { Button } from "../../../components/ui/button";
import { Badge } from "../../../components/ui/badge";
import { getOrCreateLearnerID } from "../../../lib/learner";
import { isLessonCompleted, markLessonCompleted } from "../../../shared/lib/lessonProgress";
import { fetchLessonLearn } from "../../lesson-learn/api/lessonLearnApi";
import { getActionPurposeMap, getLessonBlueprint } from "../model/lessonBlueprints";
import { useLessonRunner } from "../hooks/useLessonRunner";
import { ActionsPanel } from "./panels/ActionsPanel";
import { ChallengeBriefPanel } from "./panels/ChallengeBriefPanel";
import { GoalSubmitPanel } from "./panels/GoalSubmitPanel";
import { VisualizationPanel } from "./panels/VisualizationPanel";
import { TheoryModal } from "./TheoryModal";

interface LessonChallengePageProps {
  baseURL: string;
  lessonID: string;
  stageIndex?: number;
  onNavigate: (to: string) => void;
}

export function LessonChallengePage({
  baseURL,
  lessonID,
  stageIndex,
  onNavigate,
}: LessonChallengePageProps) {
  const [learnerID] = useState(() => getOrCreateLearnerID());
  const [closedAutoTheoryByLesson, setClosedAutoTheoryByLesson] = useState<Record<string, boolean>>({});
  const [manualOpenLessonID, setManualOpenLessonID] = useState<string | null>(null);

  const theoryQuery = useQuery({
    queryKey: ["lesson-theory", baseURL, learnerID, lessonID],
    queryFn: () => fetchLessonLearn(baseURL, lessonID, learnerID),
  });

  const {
    selectedLesson,
    selectedStage,
    selectedStageIndex,
    runResult,
    attempt,
    snapshot,
    liveError,
    lastLessonAction,
    canSend,
    errorMessage,
    isLessonsLoading,
    isStartPending,
    isGradePending,
    handleStart,
    handleLessonAction,
    handleGrade,
  } = useLessonRunner({
    baseURL,
    preferredLessonID: lessonID,
    preferredStageIndex: stageIndex,
  });

  const challengeState = snapshot?.challenge;
  const remainingSteps = challengeState?.remaining_steps ?? attempt?.limits.maxSteps ?? 0;
  const remainingPolicyChanges =
    challengeState?.remaining_policy_changes ?? attempt?.limits.maxPolicyChanges ?? 0;
  const remainingConfigChanges =
    challengeState?.remaining_config_changes ?? attempt?.limits.maxConfigChanges ?? 0;

  const theoryStage = useMemo(() => {
    const stages = theoryQuery.data?.stages ?? [];
    return stages.find((item) => item.index === selectedStageIndex) ?? stages[0] ?? null;
  }, [theoryQuery.data?.stages, selectedStageIndex]);

  const blueprint = useMemo(() => {
    return getLessonBlueprint(lessonID, selectedStage?.id);
  }, [lessonID, selectedStage?.id]);

  const actionPurposeMap = useMemo(() => {
    return getActionPurposeMap(lessonID, selectedStage?.id);
  }, [lessonID, selectedStage?.id]);

  const isTheoryAlreadySeen = wasTheorySeen(lessonID);
  const lessonCompleted = isLessonCompleted(learnerID, attempt?.lessonId ?? lessonID);
  const isTheoryOpen =
    manualOpenLessonID === lessonID ||
    (!isTheoryAlreadySeen && closedAutoTheoryByLesson[lessonID] !== true);

  const closeTheory = () => {
    window.localStorage.setItem(theorySeenKey(lessonID), "1");
    setManualOpenLessonID(null);
    setClosedAutoTheoryByLesson((current) => ({ ...current, [lessonID]: true }));
  };

  useEffect(() => {
    if (!runResult?.passed) {
      return;
    }
    const completedLessonID = attempt?.lessonId ?? lessonID;
    markLessonCompleted(learnerID, completedLessonID);
  }, [attempt?.lessonId, learnerID, lessonID, runResult?.passed]);

  return (
    <>
      <section className="grid min-h-[calc(100vh-2rem)] grid-rows-[auto_1fr] gap-3">
        <header className="rounded-xl border border-slate-200 bg-white px-4 py-3 shadow-sm">
          <div className="mb-2 flex flex-wrap justify-start gap-2">
            <Button type="button" variant="outline" onClick={() => onNavigate("/")}>
              Home
            </Button>
            <Button type="button" variant="secondary" onClick={() => setManualOpenLessonID(lessonID)}>
              Theory
            </Button>
          </div>

          <div className="flex items-center gap-2">
            <h1 className="text-xl font-semibold text-slate-900">{selectedLesson?.title ?? lessonID}</h1>
            {lessonCompleted ? <Badge variant="success">Completed</Badge> : null}
          </div>
          <p className="mt-1 text-sm text-slate-600">
            Run actions on the left, monitor visualization on the right, and open Theory anytime.
          </p>

          {(selectedLesson?.stages.length ?? 0) > 1 ? (
            <div className="mt-3 flex flex-wrap gap-2">
              {(selectedLesson?.stages ?? []).map((item) => (
                <Button
                  key={item.id}
                  type="button"
                  size="sm"
                  variant={item.index === selectedStageIndex ? "default" : "outline"}
                  onClick={() => onNavigate(buildLessonRoute(lessonID, item.index, true))}
                >
                  {item.id}: {item.title}
                </Button>
              ))}
            </div>
          ) : null}
        </header>

        <div className="grid min-h-0 gap-3 lg:grid-cols-[360px_minmax(0,1fr)]">
          <div className="grid min-h-0 content-start gap-3">
            <ChallengeBriefPanel
              objective={blueprint?.objective ?? selectedStage?.objective ?? attempt?.objective}
              description={blueprint?.description ?? attempt?.goal ?? selectedStage?.goal}
              successCriteria={blueprint?.successCriteria ?? attempt?.passConditions ?? selectedStage?.passConditions ?? []}
              visualChecks={blueprint?.visualChecks ?? selectedStage?.expectedVisualCues ?? []}
            />

            <ActionsPanel
              canSend={canSend}
              isLessonsLoading={isLessonsLoading}
              isStartPending={isStartPending}
              isGradePending={isGradePending}
              hasAttempt={Boolean(attempt?.attemptId)}
              isStageUnlocked={selectedStage?.unlocked !== false}
              remainingSteps={remainingSteps}
              remainingPolicyChanges={remainingPolicyChanges}
              remainingConfigChanges={remainingConfigChanges}
              lessonActions={selectedStage?.actionDescriptions?.map((item) => item.command) ?? []}
              actionPurposeMap={actionPurposeMap}
              actionCapabilities={attempt?.actionCapabilities}
              onStart={handleStart}
              onSubmit={handleGrade}
              onLessonAction={handleLessonAction}
            />

            <GoalSubmitPanel
              selectedStage={selectedStage}
              attemptGoal={blueprint?.objective ?? attempt?.goal}
              attemptPassConditions={runResult?.passConditions ?? blueprint?.successCriteria ?? attempt?.passConditions}
              result={runResult}
            />
          </div>

          <VisualizationPanel
            lessonID={attempt?.lessonId ?? lessonID}
            snapshot={snapshot}
            lastLessonAction={lastLessonAction}
          />
        </div>

        {errorMessage ? <p className="text-sm text-red-700">{errorMessage}</p> : null}
        {liveError ? <p className="text-sm text-red-700">{liveError}</p> : null}
      </section>

      {isTheoryOpen ? (
        <TheoryModal
          lessonTitle={selectedLesson?.title ?? lessonID}
          stage={theoryStage}
          blueprint={blueprint}
          isLoading={theoryQuery.isLoading}
          error={theoryQuery.error instanceof Error ? theoryQuery.error.message : ""}
          onClose={closeTheory}
        />
      ) : null}
    </>
  );
}

function buildLessonRoute(lessonID: string, stageIndex: number, hasMultipleStages: boolean): string {
  if (!hasMultipleStages || stageIndex <= 0) {
    return `/lesson/${lessonID}`;
  }
  return `/lesson/${lessonID}?stage=${stageIndex}`;
}

function theorySeenKey(lessonID: string): string {
  return `lesson-theory-seen:${lessonID}`;
}

function wasTheorySeen(lessonID: string): boolean {
  return window.localStorage.getItem(theorySeenKey(lessonID)) === "1";
}
