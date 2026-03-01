import { useQuery } from "@tanstack/react-query";
import { useEffect, useMemo, useState } from "react";

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

const SECTION_LINKS = [
  { id: "core-idea", label: "Core Idea" },
  { id: "mechanism", label: "Mechanism" },
  { id: "worked-example", label: "Worked Example" },
  { id: "challenge-actions", label: "Challenge Actions" },
  { id: "expected-visualization", label: "Expected Visualization" },
  { id: "common-mistakes", label: "Common Mistakes" },
];

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

  const stages = lesson?.stages ?? [];
  const [activeAnchor, setActiveAnchor] = useState(SECTION_LINKS[0].id);

  useEffect(() => {
    const onHashChange = () => {
      const currentHash = window.location.hash.replace("#", "");
      if (SECTION_LINKS.some((link) => link.id === currentHash)) {
        setActiveAnchor(currentHash);
      }
    };

    onHashChange();
    window.addEventListener("hashchange", onHashChange);
    return () => window.removeEventListener("hashchange", onHashChange);
  }, []);

  return (
    <Card className="shadow-sm">
      <CardHeader className="border-b border-slate-200 bg-gradient-to-r from-slate-50 to-cyan-50 pb-4">
        <CardTitle>{lesson?.title ?? "Lesson"}</CardTitle>
        <CardDescription>
          Review concept flow and then move to challenge from the sticky sidebar.
        </CardDescription>
      </CardHeader>
      <CardContent className="p-4 md:p-5">

      {query.isLoading ? <p className="mt-2 text-sm text-slate-600">Loading lesson theory...</p> : null}
      {query.error instanceof Error ? <p className="mt-2 text-sm text-red-700">{query.error.message}</p> : null}

      {stage ? (
        <div className="grid gap-4 lg:grid-cols-[280px_minmax(0,1fr)]">
          <aside className="h-fit rounded-lg border border-slate-200 bg-white p-3 shadow-sm lg:sticky lg:top-4">
            <div className="flex flex-wrap gap-2">
              <Button type="button" size="sm" variant="outline" onClick={() => onNavigate("/")}>
                Home
              </Button>
              <Button
                type="button"
                size="sm"
                onClick={() => onNavigate(buildChallengeRoute(lessonID, stage.index, stages.length > 1))}
              >
                Start Challenge
              </Button>
            </div>

            {stages.length > 1 ? (
              <section className="mt-4">
                <h3 className="text-xs font-semibold uppercase tracking-wide text-slate-700">Lesson Parts</h3>
                <div className="mt-2 grid gap-2">
                  {stages.map((item) => (
                    <Button
                      key={item.id}
                      type="button"
                      size="sm"
                      variant={item.index === stage.index ? "default" : "outline"}
                      className="justify-start"
                      onClick={() => onNavigate(buildLearnRoute(lessonID, item.index, true))}
                    >
                      {item.id}: {item.title}
                    </Button>
                  ))}
                </div>
              </section>
            ) : null}

            <section className="mt-4">
              <h3 className="text-xs font-semibold uppercase tracking-wide text-slate-700">On This Page</h3>
              <nav className="mt-2 grid gap-1.5">
                {SECTION_LINKS.map((link) => (
                  <a
                    key={link.id}
                    href={`#${link.id}`}
                    className={`rounded border px-2 py-1.5 text-sm transition ${
                      activeAnchor === link.id
                        ? "border-sky-200 bg-sky-50 font-medium text-sky-800"
                        : "border-transparent text-slate-600 hover:border-slate-200 hover:bg-slate-50 hover:text-slate-900"
                    }`}
                    onClick={() => setActiveAnchor(link.id)}
                  >
                    {link.label}
                  </a>
                ))}
              </nav>
            </section>
          </aside>

          <div className="space-y-3">
          <section id="core-idea" className="rounded-lg border border-slate-200 bg-slate-50 p-4 scroll-mt-4">
            <h3 className="text-sm font-semibold text-slate-900">Core Idea</h3>
            <p className="mt-2 text-sm text-slate-600">
              {stage.coreIdea ?? "Review the core operating systems concept for this stage."}
            </p>
          </section>

          <section id="mechanism" className="rounded-lg border border-slate-200 bg-slate-50 p-4 scroll-mt-4">
            <h3 className="text-sm font-semibold text-slate-900">Mechanism</h3>
            {(stage.mechanismSteps ?? []).length > 0 ? (
              <ul className="mt-2 grid gap-1.5">
                {(stage.mechanismSteps ?? []).map((item) => (
                  <li key={item} className="text-sm text-slate-600">
                    - {item}
                  </li>
                ))}
              </ul>
            ) : (
              <p className="mt-2 text-sm text-slate-600">
                Track cause-and-effect transitions in the trace and state panels.
              </p>
            )}
          </section>

          <section id="worked-example" className="rounded-lg border border-slate-200 bg-slate-50 p-4 scroll-mt-4">
            <h3 className="text-sm font-semibold text-slate-900">Worked Example</h3>
            <p className="mt-2 text-sm text-slate-600">
              {stage.workedExample ??
                `Stage objective: ${stage.objective ?? stage.title}. Focus on how state transitions satisfy the objective.`}
            </p>
          </section>

          <section id="challenge-actions" className="rounded-lg border border-slate-200 bg-slate-50 p-4 scroll-mt-4">
            <h3 className="text-sm font-semibold text-slate-900">Challenge Actions</h3>
            {(stage.preChallengeChecklist ?? []).length > 0 ? (
              <ul className="mt-2 grid gap-1.5">
                {(stage.preChallengeChecklist ?? []).map((item) => (
                  <li key={item} className="text-sm text-slate-600">
                    - {item}
                  </li>
                ))}
              </ul>
            ) : (
              <p className="mt-2 text-sm text-slate-600">No explicit actions listed for this part.</p>
            )}
          </section>

          <section id="expected-visualization" className="rounded-lg border border-slate-200 bg-slate-50 p-4 scroll-mt-4">
            <h3 className="text-sm font-semibold text-slate-900">Expected Visualization</h3>
            {(stage.expectedVisualCues ?? []).length > 0 ? (
              <ul className="mt-2 grid gap-1.5">
                {(stage.expectedVisualCues ?? []).map((item) => (
                  <li key={item} className="text-sm text-slate-600">
                    - {item}
                  </li>
                ))}
              </ul>
            ) : (
              <p className="mt-2 text-sm text-slate-600">Watch trace order, process progress, and objective-aligned metrics.</p>
            )}
          </section>

          <section id="common-mistakes" className="rounded-lg border border-slate-200 bg-slate-50 p-4 scroll-mt-4">
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
          </div>
        </div>
      ) : null}
      </CardContent>
    </Card>
  );
}

function buildLearnRoute(lessonID: string, stageIndex: number, hasMultipleStages: boolean): string {
  if (!hasMultipleStages || stageIndex <= 0) {
    return `/lesson/${lessonID}/learn`;
  }
  return `/lesson/${lessonID}/learn?stage=${stageIndex}`;
}

function buildChallengeRoute(lessonID: string, stageIndex: number, hasMultipleStages: boolean): string {
  if (!hasMultipleStages || stageIndex <= 0) {
    return `/lesson/${lessonID}/challenge`;
  }
  return `/lesson/${lessonID}/challenge?stage=${stageIndex}`;
}
