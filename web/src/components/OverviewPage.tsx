import { useMemo } from "react";

import type { CurriculumSection } from "../lib/lessonApi";

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
        locked: section.coming_soon,
        comingSoon: section.coming_soon,
        lessonNodes,
        completedStages: section.completed_stages ?? completedStages,
        totalStages: section.total_stages ?? totalStages,
      };
    });
  }, [sections]);

  return (
    <section className="panel section-overview-panel">
      <h2>Course Overview</h2>
      <p className="subtitle">Follow OSTEP subjects from top to bottom.</p>

      {isLoading ? <p className="empty">Loading sections...</p> : null}
      {errorMessage ? <p className="error">{errorMessage}</p> : null}

      {!isLoading && !errorMessage ? (
        <div className="subject-stack">
          {subjectBlocks.map((subject) => {
            const completionRate =
              subject.totalStages > 0
                ? Math.round((subject.completedStages / subject.totalStages) * 100)
                : 0;
            return (
              <article
                key={subject.id}
                className={
                  subject.locked
                    ? `subject-block subject-${subject.id} is-locked`
                    : `subject-block subject-${subject.id}`
                }
              >
                <div className="subject-header">
                  <h3>{subject.title}</h3>
                  {subject.comingSoon ? <span className="badge fail">Coming Soon</span> : null}
                </div>
                <p className="lesson-outcome">{subject.subtitle}</p>
                {subject.totalStages > 0 ? (
                  <p className="lesson-outcome subject-progress">
                    Progress: {subject.completedStages}/{subject.totalStages} ({completionRate}%)
                  </p>
                ) : (
                  <p className="lesson-outcome">Lesson nodes will appear here after implementation.</p>
                )}

                <div className="lesson-node-grid">
                  {subject.lessonNodes.map((node) => (
                    <div key={node.id} className="lesson-node-wrap">
                      <button
                        type="button"
                        className={`lesson-node lesson-node--${node.status}`}
                        disabled={node.status === "locked"}
                        onClick={() => onNavigate(`/lesson/${node.id}/learn?stage=${node.firstUnlockedStage}`)}
                        aria-label={node.title}
                      />
                      <span className="lesson-node-title">{node.title}</span>
                    </div>
                  ))}
                </div>
              </article>
            );
          })}
        </div>
      ) : null}
    </section>
  );
}
