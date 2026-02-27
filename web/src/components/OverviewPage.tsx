import { useMemo } from "react";

import type { LessonSummary } from "../lib/lessonApi";

interface OverviewPageProps {
  lessons: LessonSummary[];
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

const OSTEP_SUBJECTS: Array<{
  id: string;
  title: string;
  subtitle: string;
  sectionIDs: string[];
  comingSoon?: boolean;
}> = [
  {
    id: "introduction",
    title: "Introduction",
    subtitle: "OSTEP setup and foundational framing",
    sectionIDs: [],
    comingSoon: true,
  },
  {
    id: "virtualization",
    title: "Virtualization",
    subtitle: "CPU and memory virtualization lessons",
    sectionIDs: ["virtualization"],
  },
  {
    id: "concurrency",
    title: "Concurrency",
    subtitle: "Threads, wakeups, and interrupt-driven progress",
    sectionIDs: ["concurrency"],
  },
  {
    id: "persistence",
    title: "Persistence",
    subtitle: "Storage and filesystem correctness",
    sectionIDs: ["persistence"],
  },
  {
    id: "security",
    title: "Security",
    subtitle: "Authentication, access control, and protection",
    sectionIDs: [],
    comingSoon: true,
  },
];

export function OverviewPage({
  lessons,
  isLoading,
  errorMessage,
  onNavigate,
}: OverviewPageProps) {
  const subjectBlocks = useMemo(() => {
    return OSTEP_SUBJECTS.map((subject): SubjectBlock => {
      const matchedLessons = lessons.filter((lesson) => {
        const sectionID = lesson.section_id ?? lesson.module;
        return subject.sectionIDs.includes(sectionID);
      });
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
        id: subject.id,
        title: subject.title,
        subtitle: subject.subtitle,
        locked: subject.comingSoon === true,
        comingSoon: subject.comingSoon === true,
        lessonNodes,
        completedStages,
        totalStages,
      };
    });
  }, [lessons]);

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
                        onClick={() => onNavigate(`/challenge/${node.id}/${node.firstUnlockedStage}`)}
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
