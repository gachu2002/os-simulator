import { Button } from "../../../components/ui/button";
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from "../../../components/ui/card";
import { useLessonRunner } from "../hooks/useLessonRunner";
import { ActionsPanel } from "./panels/ActionsPanel";
import { GoalSubmitPanel } from "./panels/GoalSubmitPanel";
import { VisualizationPanel } from "./panels/VisualizationPanel";

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
  const {
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
    isLessonsLoading,
    isStartPending,
    isGradePending,
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

  return (
    <Card className="shadow-sm">
      <CardHeader className="pb-0">
      <div className="mb-3 flex flex-wrap justify-start gap-2">
        <Button
          type="button"
          variant="outline"
          onClick={() => onNavigate("/")}
        >
          Home
        </Button>
        <Button
          type="button"
          variant="secondary"
          onClick={() => onNavigate(`/lesson/${lessonID}/learn?stage=${stageIndex ?? 0}`)}
        >
          Back To Learn
        </Button>
      </div>

      <CardTitle>{attempt?.lessonId ?? lessonID} Challenge</CardTitle>
      <CardDescription>
        Challenge page has three sections: Actions, Visualization, Goal + Submit.
      </CardDescription>
      </CardHeader>
      <CardContent>

      <ActionsPanel
        canSend={canSend}
        isLessonsLoading={isLessonsLoading}
        isStartPending={isStartPending}
        isGradePending={isGradePending}
        hasAttempt={Boolean(attempt?.attemptId)}
        isStageUnlocked={selectedStage?.unlocked !== false}
        policy={policy}
        quantum={quantum}
        frames={frames}
        tlbEntries={tlbEntries}
        diskLatency={diskLatency}
        terminalLatency={terminalLatency}
        remainingSteps={remainingSteps}
        remainingPolicyChanges={remainingPolicyChanges}
        remainingConfigChanges={remainingConfigChanges}
        isAllowed={isCommandAllowed}
        onPolicyChange={setPolicy}
        onQuantumChange={setQuantum}
        onFramesChange={setFrames}
        onTLBEntriesChange={setTLBEntries}
        onDiskLatencyChange={setDiskLatency}
        onTerminalLatencyChange={setTerminalLatency}
        onStart={handleStart}
        onSubmit={handleGrade}
        onCommand={handleCommand}
      />

      <VisualizationPanel snapshot={snapshot} />

      <GoalSubmitPanel
        selectedStage={selectedStage}
        attemptGoal={attempt?.goal}
        attemptPassConditions={attempt?.passConditions}
        result={runResult}
      />

      {errorMessage ? <p className="mt-2 text-sm text-red-700">{errorMessage}</p> : null}
      {liveError ? <p className="mt-2 text-sm text-red-700">{liveError}</p> : null}
      </CardContent>
    </Card>
  );
}
