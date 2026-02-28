import { useMemo } from "react";

import type { CurriculumSection } from "../../../entities/lesson/model";
import { cn } from "../../../shared/lib/cn";
import { Badge } from "../../../components/ui/badge";
import { Button } from "../../../components/ui/button";
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
  const subjectBlocks = useMemo(() => {
    return sections.map((section): SubjectBlock => {
      const matchedLessons = section.lessons ?? [];
      const lessonNodes: LessonNode[] = matchedLessons.map((lesson) => {
        const completedStages = lesson.stages.filter((stage) => stage.completed).length;
        const unlockedStages = lesson.stages.filter((stage) => stage.unlocked !== false);
        const allCompleted = completedStages === lesson.stages.length && lesson.stages.length > 0;
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
  }, [sections]);

  return (
    <Card className="shadow-sm">
      <CardHeader className="pb-0">
      <CardTitle>Course Overview</CardTitle>
      <CardDescription>Follow OSTEP subjects from top to bottom.</CardDescription>
      </CardHeader>
      <CardContent>

      {isLoading ? <p className="mt-2 text-sm text-slate-600">Loading sections...</p> : null}
      {errorMessage ? <p className="mt-2 text-sm text-red-700">{errorMessage}</p> : null}

      {!isLoading && !errorMessage ? (
        <div className="mt-3 grid gap-4">
          {subjectBlocks.map((subject) => {
            const completionRate =
              subject.totalStages > 0
                ? Math.round((subject.completedStages / subject.totalStages) * 100)
                : 0;
            return (
              <article
                key={subject.id}
                className={cn(
                  "rounded-xl border p-4 shadow-inner",
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
                <p className="mt-2 text-sm text-slate-600">{subject.subtitle}</p>
                {subject.totalStages > 0 ? (
                  <p className="mt-2 inline-block rounded-full border border-blue-200 bg-white/70 px-2 py-0.5 text-xs text-slate-700">
                    Progress: {subject.completedStages}/{subject.totalStages} ({completionRate}%)
                  </p>
                ) : (
                  <p className="mt-2 text-sm text-slate-600">
                    Lesson nodes will appear here after implementation.
                  </p>
                )}

                <div className="mt-3 grid grid-cols-[repeat(auto-fill,minmax(130px,1fr))] gap-3 max-[720px]:grid-cols-[repeat(auto-fill,minmax(105px,1fr))]">
                  {subject.lessonNodes.map((node) => (
                    <div key={node.id} className="grid justify-items-center gap-1.5">
                      <Button
                        type="button"
                        size="sm"
                        variant="outline"
                        className={cn(
                          "h-11 w-11 rounded-full border-2 transition hover:-translate-y-0.5 hover:shadow-sm",
                          getNodeTone(node.status),
                        )}
                        disabled={node.status === "locked"}
                        onClick={() =>
                          onNavigate(`/lesson/${node.id}/learn?stage=${node.firstUnlockedStage}`)
                        }
                        aria-label={node.title}
                      />
                      <span className="min-h-[2.1em] text-center text-xs leading-tight text-slate-600">
                        {node.title}
                      </span>
                    </div>
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

function getNodeTone(status: LessonNode["status"]): string {
  if (status === "passed") {
    return "border-emerald-500 bg-emerald-100";
  }
  if (status === "ready") {
    return "border-slate-500 bg-slate-200";
  }
  return "border-slate-300 bg-slate-50";
}

function getSubjectTone(sectionID: string): string {
  switch (sectionID) {
    case "introduction":
      return "border-slate-200 bg-gradient-to-b from-slate-50 to-slate-100";
    case "virtualization":
      return "border-sky-200 bg-gradient-to-b from-sky-50 to-blue-50";
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
