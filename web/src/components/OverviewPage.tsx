import { useMemo } from "react";

import type { LessonSummary } from "../lib/lessonApi";

interface OverviewPageProps {
  lessons: LessonSummary[];
  isLoading: boolean;
  errorMessage: string;
  onNavigate: (to: string) => void;
}

interface SectionSummary {
  id: string;
  title: string;
  lessons: number;
  completedStages: number;
  totalStages: number;
  nextLessonID: string;
}

export function OverviewPage({
  lessons,
  isLoading,
  errorMessage,
  onNavigate,
}: OverviewPageProps) {
  const sectionSummaries = useMemo(() => {
    const map = new Map<string, SectionSummary>();
    for (const lesson of lessons) {
      const sectionID = lesson.section_id ?? lesson.module;
      const sectionTitle = lesson.section_title ?? lesson.module;
      const totalStages = lesson.stages.length;
      const completedStages = lesson.stages.filter((stage) => stage.completed).length;
      const existing = map.get(sectionID);
      if (!existing) {
        map.set(sectionID, {
          id: sectionID,
          title: sectionTitle,
          lessons: 1,
          completedStages,
          totalStages,
          nextLessonID: lesson.id,
        });
        continue;
      }

      const nextLessonID =
        existing.nextLessonID === "" && completedStages < totalStages
          ? lesson.id
          : existing.nextLessonID;

      map.set(sectionID, {
        ...existing,
        lessons: existing.lessons + 1,
        completedStages: existing.completedStages + completedStages,
        totalStages: existing.totalStages + totalStages,
        nextLessonID,
      });
    }

    return Array.from(map.values()).sort((left, right) => left.id.localeCompare(right.id));
  }, [lessons]);

  return (
    <section className="panel section-overview-panel">
      <h2>Course Overview</h2>
      <p className="subtitle">
        Follow OSTEP sections in order: CPU, Memory, Concurrency, then Persistence.
      </p>

      {isLoading ? <p className="empty">Loading sections...</p> : null}
      {errorMessage ? <p className="error">{errorMessage}</p> : null}

      {!isLoading && !errorMessage ? (
        <div className="section-grid">
          {sectionSummaries.map((section) => {
            const completionRate =
              section.totalStages > 0
                ? Math.round((section.completedStages / section.totalStages) * 100)
                : 0;
            return (
              <article key={section.id} className="section-card">
                <h3>{section.title}</h3>
                <p className="lesson-outcome">Lessons: {section.lessons}</p>
                <p className="lesson-outcome">
                  Progress: {section.completedStages}/{section.totalStages} ({completionRate}%)
                </p>
                <div className="control-row">
                  <button type="button" onClick={() => onNavigate(`/sections/${section.id}`)}>
                    View Section
                  </button>
                  <button
                    type="button"
                    onClick={() => onNavigate(`/challenge/${section.nextLessonID}/0`)}
                  >
                    Continue
                  </button>
                </div>
              </article>
            );
          })}
        </div>
      ) : null}
    </section>
  );
}
