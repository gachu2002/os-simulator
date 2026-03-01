import { Badge } from "../../../../components/ui/badge";
import { Button } from "../../../../components/ui/button";
import type { ActionCapabilities } from "../../../../entities/challenge/actionCapabilities";
import { type LessonActionOptions, toActionPreset } from "../../model/actionPresets";

interface ActionsPanelProps {
  canSend: boolean;
  isLessonsLoading: boolean;
  isStartPending: boolean;
  isGradePending: boolean;
  hasAttempt: boolean;
  isStageUnlocked: boolean;
  remainingSteps: number;
  remainingPolicyChanges: number;
  remainingConfigChanges: number;
  lessonActions: string[];
  actionPurposeMap?: Record<string, string>;
  actionCapabilities?: ActionCapabilities;
  onStart: () => void;
  onSubmit: () => void;
  onLessonAction: (
    action: string,
    options?: LessonActionOptions,
  ) => void;
}

export function ActionsPanel(props: ActionsPanelProps) {
  const {
    canSend,
    isLessonsLoading,
    isStartPending,
    isGradePending,
    hasAttempt,
    isStageUnlocked,
    remainingSteps,
    remainingPolicyChanges,
    remainingConfigChanges,
    lessonActions,
    actionPurposeMap,
    actionCapabilities,
    onStart,
    onSubmit,
    onLessonAction,
  } = props;

  const supportedNow = actionCapabilities?.supportedNow ?? [];
  const visibleActions =
    supportedNow.length > 0
      ? lessonActions.filter((action) => supportedNow.includes(action))
      : lessonActions;

  return (
    <section className="rounded-lg border border-slate-200 bg-white p-3">
      <div className="flex items-center justify-between gap-2">
        <h3 className="text-sm font-semibold text-slate-900">Actions</h3>
        {actionCapabilities ? (
          <Badge variant="success">{actionCapabilities.supportedNow.length} usable</Badge>
        ) : null}
      </div>

      <div className="mt-2 flex flex-wrap items-end gap-2">
        <Button
          type="button"
          disabled={isStartPending || isLessonsLoading || !isStageUnlocked}
          onClick={onStart}
        >
          {isStartPending ? "Starting..." : "Start Challenge"}
        </Button>
        <Button type="button" variant="success" disabled={isGradePending || !hasAttempt} onClick={onSubmit}>
          {isGradePending ? "Submitting..." : "Submit"}
        </Button>
      </div>

      {hasAttempt ? (
        <p className="mt-2 text-xs text-slate-600">
          Remaining budget: steps {remainingSteps}, policy edits {remainingPolicyChanges}, config edits {remainingConfigChanges}
        </p>
      ) : null}

      {hasAttempt ? (
        <div className="mt-3 grid gap-2 sm:grid-cols-2">
          {visibleActions.map((action) => {
            const preset = toActionPreset(action);
            return (
              <article key={action} className="rounded-md border border-slate-200 bg-slate-50 p-2">
                <Button
                  type="button"
                  variant="outline"
                  className="w-full justify-start"
                  disabled={!canSend}
                  onClick={() => onLessonAction(action, preset.options)}
                >
                  {preset.label}
                </Button>
                <p className="mt-1 text-xs text-slate-600">
                  {actionPurposeMap?.[action] ?? "Run this action and inspect state change in visualization."}
                </p>
              </article>
            );
          })}
        </div>
      ) : null}

      {hasAttempt && visibleActions.length === 0 ? (
        <p className="mt-3 text-sm text-slate-600">
          No executable actions are exposed for this stage yet.
        </p>
      ) : null}
    </section>
  );
}
