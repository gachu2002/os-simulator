import type { LessonLearnStage } from "../../../entities/lesson/model";
import { Button } from "../../../components/ui/button";
import type { LessonBlueprintPart } from "../model/lessonBlueprints";

interface TheoryModalProps {
  lessonTitle: string;
  stage: LessonLearnStage | null;
  blueprint: LessonBlueprintPart | null;
  isLoading: boolean;
  error: string;
  onClose: () => void;
}

export function TheoryModal({ lessonTitle, stage, blueprint, isLoading, error, onClose }: TheoryModalProps) {
  return (
    <div
      className="fixed inset-0 z-50 grid place-items-center bg-slate-900/40 px-3 py-4"
      onClick={onClose}
      role="presentation"
    >
      <section
        className="max-h-[92vh] w-full max-w-3xl overflow-y-auto rounded-xl border border-slate-200 bg-white p-4 shadow-xl"
        onClick={(event) => event.stopPropagation()}
      >
        <div className="flex items-start justify-between gap-3">
          <div>
            <p className="text-xs font-semibold uppercase tracking-[0.14em] text-sky-700">Theory</p>
            <h2 className="mt-1 text-lg font-semibold text-slate-900">{lessonTitle}</h2>
          </div>
          <Button type="button" variant="outline" size="sm" onClick={onClose}>
            Close
          </Button>
        </div>

        {isLoading ? <p className="mt-3 text-sm text-slate-600">Loading theory...</p> : null}
        {error ? <p className="mt-3 text-sm text-red-700">{error}</p> : null}

        {stage ? (
          <div className="mt-3 grid gap-3">
            <section className="rounded-lg border border-slate-200 bg-slate-50 p-3">
              <h3 className="text-sm font-semibold text-slate-900">Core Idea</h3>
              <p className="mt-1 text-sm text-slate-600">{blueprint?.objective ?? stage.coreIdea ?? "Review this lesson concept."}</p>
            </section>

            <section className="rounded-lg border border-slate-200 bg-slate-50 p-3">
              <h3 className="text-sm font-semibold text-slate-900">Mechanism</h3>
              {(blueprint?.theory ?? stage.mechanismSteps ?? []).length > 0 ? (
                <ul className="mt-1 grid gap-1.5">
                  {(blueprint?.theory ?? stage.mechanismSteps ?? []).map((item) => (
                    <li key={item} className="text-sm text-slate-600">
                      - {item}
                    </li>
                  ))}
                </ul>
              ) : (
                <p className="mt-1 text-sm text-slate-600">Track cause-and-effect transitions in the simulator.</p>
              )}
            </section>

            <section className="rounded-lg border border-slate-200 bg-slate-50 p-3">
              <h3 className="text-sm font-semibold text-slate-900">Challenge Goal</h3>
              <p className="mt-1 text-sm text-slate-600">
                {blueprint?.description ?? stage.goal ?? stage.workedExample ?? "Apply theory to challenge steps."}
              </p>
              <h4 className="mt-2 text-xs font-semibold uppercase tracking-wide text-slate-700">Success Criteria</h4>
              {(blueprint?.successCriteria ?? stage.preChallengeChecklist ?? []).length > 0 ? (
                <ul className="mt-1 grid gap-1.5">
                  {(blueprint?.successCriteria ?? stage.preChallengeChecklist ?? []).map((item) => (
                    <li key={item} className="text-sm text-slate-600">
                      - {item}
                    </li>
                  ))}
                </ul>
              ) : (
                <p className="mt-1 text-sm text-slate-600">No extra checklist for this lesson.</p>
              )}
            </section>

            <section className="rounded-lg border border-slate-200 bg-slate-50 p-3">
              <h3 className="text-sm font-semibold text-slate-900">Expected Visualization</h3>
              {(blueprint?.visualChecks ?? stage.expectedVisualCues ?? []).length > 0 ? (
                <ul className="mt-1 grid gap-1.5">
                  {(blueprint?.visualChecks ?? stage.expectedVisualCues ?? []).map((item) => (
                    <li key={item} className="text-sm text-slate-600">
                      - {item}
                    </li>
                  ))}
                </ul>
              ) : (
                <p className="mt-1 text-sm text-slate-600">Watch trace, process states, and key metrics.</p>
              )}
            </section>

            {(blueprint?.commonPitfalls ?? stage.commonMistakes ?? []).length > 0 ? (
              <section className="rounded-lg border border-amber-200 bg-amber-50 p-3">
                <h3 className="text-sm font-semibold text-slate-900">Common Pitfalls</h3>
                <ul className="mt-1 grid gap-1.5">
                  {(blueprint?.commonPitfalls ?? stage.commonMistakes ?? []).map((item) => (
                    <li key={item} className="text-sm text-slate-700">
                      - {item}
                    </li>
                  ))}
                </ul>
              </section>
            ) : null}
          </div>
        ) : null}
      </section>
    </div>
  );
}
