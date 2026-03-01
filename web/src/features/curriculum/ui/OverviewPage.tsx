import { useMemo, useState } from "react";

import type { CurriculumSection } from "../../../entities/lesson/model";
import { cn } from "../../../shared/lib/cn";
import { getOrCreateLearnerID } from "../../../lib/learner";
import { completedLessonCount, isLessonCompleted } from "../../../shared/lib/lessonProgress";
import { Badge } from "../../../components/ui/badge";
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from "../../../components/ui/card";

interface OverviewPageProps {
  sections: CurriculumSection[];
  isLoading: boolean;
  errorMessage: string;
  onNavigate: (to: string) => void;
}

interface LessonNode {
  id: string;
  title: string;
  status: "locked" | "ready" | "passed";
  firstUnlockedStage: number;
  hasMultipleStages: boolean;
}

interface SubjectBlock {
  id: string;
  title: string;
  subtitle: string;
  locked: boolean;
  comingSoon: boolean;
  lessonNodes: LessonNode[];
  completedStages: number;
  totalStages: number;
}

export function OverviewPage({
  sections,
  isLoading,
  errorMessage,
  onNavigate,
}: OverviewPageProps) {
  const [learnerID] = useState(() => getOrCreateLearnerID());
  const subjectBlocks = useMemo(() => {
    return sections.map((section): SubjectBlock => {
      const matchedLessons = section.lessons ?? [];
      const lessonNodes: LessonNode[] = matchedLessons.map((lesson) => {
        const completedStages = lesson.stages.filter((stage) => stage.completed).length;
        const unlockedStages = lesson.stages.filter((stage) => stage.unlocked !== false);
        const allCompleted =
          isLessonCompleted(learnerID, lesson.id) ||
          (completedStages === lesson.stages.length && lesson.stages.length > 0);
        const status: LessonNode["status"] = allCompleted
          ? "passed"
          : unlockedStages.length === 0
            ? "locked"
            : "ready";
        const firstUnlockedStage = unlockedStages[0]?.index ?? 0;

        return {
          id: lesson.id,
          title: lesson.title,
          status,
          firstUnlockedStage,
          hasMultipleStages: lesson.stages.length > 1,
        };
      });

      const totalStages = matchedLessons.reduce((sum, lesson) => sum + lesson.stages.length, 0);
      const completedStages = matchedLessons.reduce(
        (sum, lesson) => sum + lesson.stages.filter((stage) => stage.completed).length,
        0,
      );

      return {
        id: section.id,
        title: section.title,
        subtitle: section.subtitle ?? "",
        locked: section.comingSoon,
        comingSoon: section.comingSoon,
        lessonNodes,
        completedStages: section.completedStages ?? completedStages,
        totalStages: section.totalStages ?? totalStages,
      };
    });
  }, [learnerID, sections]);

  const completedLessons = completedLessonCount(learnerID);

  return (
    <Card className="overflow-hidden border-slate-200 shadow-sm">
      <CardHeader className="border-b border-slate-200 bg-gradient-to-r from-cyan-50 via-sky-50 to-blue-50 pb-4">
      <p className="text-xs font-semibold uppercase tracking-[0.14em] text-sky-700">Course Overview</p>
      <CardTitle className="mt-1">Section 1: Virtualization - CPU</CardTitle>
      <CardDescription>Pick a lesson card to open the lesson workspace.</CardDescription>
      <p className="mt-2 inline-block w-fit rounded-full border border-emerald-200 bg-emerald-50 px-2 py-0.5 text-xs text-emerald-700">
        Completed lessons: {completedLessons}
      </p>
      </CardHeader>
      <CardContent>

      {isLoading ? <p className="mt-2 text-sm text-slate-600">Loading sections...</p> : null}
      {errorMessage ? <p className="mt-2 text-sm text-red-700">{errorMessage}</p> : null}

      {!isLoading && !errorMessage ? (
        <div className="mt-4 grid gap-4">
          {subjectBlocks.map((subject) => {
            const completionRate =
              subject.totalStages > 0
                ? Math.round((subject.completedStages / subject.totalStages) * 100)
                : 0;
            return (
              <article
                key={subject.id}
                className={cn(
                  "rounded-xl border p-4",
                  subject.locked && "opacity-80",
                  getSubjectTone(subject.id),
                )}
              >
                <div className="flex items-center justify-between gap-2">
                  <h3 className="text-base font-semibold text-slate-900">{subject.title}</h3>
                  {subject.comingSoon ? (
                    <Badge variant="destructive">Coming Soon</Badge>
                  ) : null}
                </div>
                {subject.subtitle ? <p className="mt-2 text-sm text-slate-600">{subject.subtitle}</p> : null}
                {subject.totalStages > 0 ? (
                  <p className="mt-2 inline-block rounded-full border border-blue-200 bg-white/70 px-2 py-0.5 text-xs text-slate-700">
                    Progress: {subject.completedStages}/{subject.totalStages} ({completionRate}%)
                  </p>
                ) : (
                  <p className="mt-2 text-sm text-slate-600">
                    Lesson nodes will appear here after implementation.
                  </p>
                )}

                <div className="mt-4 grid gap-2 sm:grid-cols-2 xl:grid-cols-3">
                  {subject.lessonNodes.map((node) => (
                    <button
                      key={node.id}
                      type="button"
                      disabled={node.status === "locked"}
                      onClick={() =>
                        onNavigate(buildLessonRoute(node.id, node.firstUnlockedStage, node.hasMultipleStages))
                      }
                      className={cn(
                        "grid gap-1.5 rounded-lg border border-slate-200 bg-white p-3 text-left shadow-sm transition disabled:cursor-not-allowed disabled:opacity-70",
                        node.status !== "locked" && "hover:-translate-y-0.5 hover:border-sky-300 hover:shadow",
                      )}
                    >
                      <div className="flex items-center justify-between gap-2">
                        <p className="line-clamp-2 text-sm font-semibold text-slate-900">{node.title}</p>
                        <Badge className={cn("capitalize", getNodeTone(node.status))}>{node.status}</Badge>
                      </div>
                      <p className="text-xs text-slate-500">
                        {node.status === "locked" ? "Locked until previous lesson is passed" : "Open lesson"}
                      </p>
                    </button>
                  ))}
                </div>
              </article>
            );
          })}
        </div>
      ) : null}
      </CardContent>
    </Card>
  );
}

function buildLessonRoute(lessonID: string, stageIndex: number, hasMultipleStages: boolean): string {
  if (!hasMultipleStages || stageIndex <= 0) {
    return `/lesson/${lessonID}`;
  }
  return `/lesson/${lessonID}?stage=${stageIndex}`;
}

function getNodeTone(status: LessonNode["status"]): string {
  if (status === "passed") {
    return "border-emerald-200 bg-emerald-50 text-emerald-700";
  }
  if (status === "ready") {
    return "border-sky-200 bg-sky-50 text-sky-700";
  }
  return "border-slate-200 bg-slate-100 text-slate-600";
}

function getSubjectTone(sectionID: string): string {
  switch (sectionID) {
    case "introduction":
      return "border-slate-200 bg-gradient-to-b from-slate-50 to-slate-100";
    case "virtualization":
      return "border-sky-200 bg-gradient-to-b from-sky-50 to-blue-50";
    case "virtualization-cpu":
      return "border-cyan-200 bg-gradient-to-b from-cyan-50 to-sky-50";
    case "concurrency":
      return "border-emerald-200 bg-gradient-to-b from-emerald-50 to-green-50";
    case "persistence":
      return "border-orange-200 bg-gradient-to-b from-orange-50 to-amber-50";
    case "security":
      return "border-slate-200 bg-gradient-to-b from-slate-50 to-zinc-100";
    default:
      return "border-slate-200 bg-gradient-to-b from-slate-50 to-blue-50";
  }
}
