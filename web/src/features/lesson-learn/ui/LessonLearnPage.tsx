import { useQuery } from "@tanstack/react-query";
import { useMemo, useState } from "react";

import { Button } from "../../../components/ui/button";
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from "../../../components/ui/card";
import { getOrCreateLearnerID } from "../../../lib/learner";
import { fetchLessonLearn } from "../api/lessonLearnApi";

interface LessonLearnPageProps {
  baseURL: string;
  lessonID: string;
  onNavigate: (to: string) => void;
  preferredStageIndex?: number;
}

export function LessonLearnPage({
  baseURL,
  lessonID,
  onNavigate,
  preferredStageIndex,
}: LessonLearnPageProps) {
  const [learnerID] = useState(() => getOrCreateLearnerID());
  const query = useQuery({
    queryKey: ["lesson-learn", baseURL, learnerID, lessonID],
    queryFn: () => fetchLessonLearn(baseURL, lessonID, learnerID),
  });

  const stageIndex = Number.isFinite(preferredStageIndex) ? preferredStageIndex ?? 0 : 0;
  const lesson = query.data;
  const stage = useMemo(() => {
    if (!lesson) {
      return null;
    }
    return lesson.stages.find((item) => item.index === stageIndex) ?? lesson.stages[0] ?? null;
  }, [lesson, stageIndex]);

  return (
    <Card className="shadow-sm">
      <CardHeader className="pb-0">
      <div className="mb-3 flex justify-start">
        <Button
          type="button"
          variant="outline"
          onClick={() => onNavigate("/")}
        >
          Home
        </Button>
      </div>

      <CardTitle>{lesson?.title ?? "Lesson"}</CardTitle>
      <CardDescription>
        Learn page: theory only. Review this before opening the challenge.
      </CardDescription>
      </CardHeader>
      <CardContent>

      {query.isLoading ? <p className="mt-2 text-sm text-slate-600">Loading lesson theory...</p> : null}
      {query.error instanceof Error ? <p className="mt-2 text-sm text-red-700">{query.error.message}</p> : null}

      {stage ? (
        <>
          <section className="mt-3 rounded-lg border border-slate-200 bg-slate-50 p-3">
            <h3 className="text-sm font-semibold text-slate-900">Core Idea</h3>
            <p className="mt-2 text-sm text-slate-600">
              {stage.coreIdea ?? "Review the core operating systems concept for this stage."}
            </p>
          </section>

          <section className="mt-3 rounded-lg border border-slate-200 bg-slate-50 p-3">
            <h3 className="text-sm font-semibold text-slate-900">Mechanism</h3>
            {(stage.mechanismSteps ?? []).length > 0 ? (
              (stage.mechanismSteps ?? []).map((item) => (
                <p key={item} className="mt-2 text-sm text-slate-600">
                  - {item}
                </p>
              ))
            ) : (
              <p className="mt-2 text-sm text-slate-600">
                Track cause-and-effect transitions in the trace and state panels.
              </p>
            )}
          </section>

          <section className="mt-3 rounded-lg border border-slate-200 bg-slate-50 p-3">
            <h3 className="text-sm font-semibold text-slate-900">Worked Example</h3>
            <p className="mt-2 text-sm text-slate-600">
              {stage.workedExample ??
                `Stage objective: ${stage.objective ?? stage.title}. Focus on how state transitions satisfy the objective.`}
            </p>
          </section>

          <section className="mt-3 rounded-lg border border-slate-200 bg-slate-50 p-3">
            <h3 className="text-sm font-semibold text-slate-900">Common Mistakes</h3>
            {(stage.commonMistakes ?? []).length > 0 ? (
              (stage.commonMistakes ?? []).map((item) => (
                <p key={item} className="mt-2 text-sm text-slate-600">
                  - {item}
                </p>
              ))
            ) : (
              <>
                <p className="mt-2 text-sm text-slate-600">
                  Do not infer correctness from one metric alone; always confirm trace and state evidence.
                </p>
                <p className="mt-2 text-sm text-slate-600">
                  Do not ignore deterministic limits; step and config budgets are part of the challenge contract.
                </p>
                <p className="mt-2 text-sm text-slate-600">
                  Do not submit before checking expected visual cues against observed behavior.
                </p>
              </>
            )}
          </section>

          <section className="mt-3 rounded-lg border border-slate-200 bg-slate-50 p-3">
            <h3 className="text-sm font-semibold text-slate-900">What To Watch In Challenge</h3>
            {(stage.preChallengeChecklist ?? []).length ? (
              (stage.preChallengeChecklist ?? []).map((item) => (
                <p key={item} className="mt-2 text-sm text-slate-600">
                  - {item}
                </p>
              ))
            ) : (stage.expectedVisualCues ?? []).length ? (
              (stage.expectedVisualCues ?? []).map((item) => (
                <p key={item} className="mt-2 text-sm text-slate-600">
                  - {item}
                </p>
              ))
            ) : (
              <p className="mt-2 text-sm text-slate-600">
                Watch trace order, process progress, and objective-aligned metrics.
              </p>
            )}
          </section>

          <div className="mt-3 flex flex-wrap items-end gap-2.5">
            <Button
              type="button"
              onClick={() => onNavigate(`/lesson/${lessonID}/challenge?stage=${stage.index}`)}
            >
              Go To Challenge
            </Button>
          </div>
        </>
      ) : null}
      </CardContent>
    </Card>
  );
}
