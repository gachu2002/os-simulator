import type { ChallengeGrade } from "../../../../entities/challenge/model";
import type { LessonStageSummary } from "../../../../entities/lesson/model";
import { Badge } from "../../../../components/ui/badge";

interface GoalSubmitPanelProps {
  selectedStage: LessonStageSummary | null;
  attemptGoal?: string;
  attemptPassConditions?: string[];
  result: ChallengeGrade | null;
}

export function GoalSubmitPanel({
  selectedStage,
  attemptGoal,
  attemptPassConditions,
  result,
}: GoalSubmitPanelProps) {
  const passConditions = attemptPassConditions ?? selectedStage?.passConditions ?? [];

  return (
    <section className="mt-3 rounded-lg border border-slate-200 bg-slate-50 p-3">
      <h3 className="text-sm font-semibold text-slate-900">3) Goal + Submit</h3>
      <p className="mt-2 text-sm text-slate-600">
        Goal: {attemptGoal ?? selectedStage?.goal ?? selectedStage?.objective ?? "Pass all checks."}
      </p>

      {passConditions.map((item) => (
        <p key={item} className="mt-2 text-sm text-slate-600">
          - {item}
        </p>
      ))}

      {result ? (
        <>
          <div className="mt-3 flex flex-wrap gap-2 text-sm text-slate-600">
            <Badge variant={result.passed ? "success" : "destructive"}>
              {result.passed ? "passed" : "failed"}
            </Badge>
            <span>feedback: {result.feedbackKey}</span>
          </div>

          {result.validatorResults?.map((item) => (
            <p key={`${item.name}.${item.type}`} className="mt-2 text-sm text-slate-600">
              - {item.passed ? "PASS" : "FAIL"}: {item.name} | expected {item.expected ?? "n/a"},
              actual {item.actual ?? "n/a"}
            </p>
          ))}

          {!result.passed && result.hint ? (
            <p className="mt-2 text-sm text-orange-700">
              Hint L{result.hintLevel ?? 0}: {result.hint}
            </p>
          ) : null}
        </>
      ) : (
        <p className="mt-2 text-sm text-slate-600">
          Start challenge actions, then submit to get pass/fail results.
        </p>
      )}
    </section>
  );
}
