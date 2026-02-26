import { useMemo } from "react";

import type { LessonSummary } from "../lib/lessonApi";

interface SectionPageProps {
  lessons: LessonSummary[];
  sectionID: string;
  onNavigate: (to: string) => void;
}

export function SectionPage({ lessons, sectionID, onNavigate }: SectionPageProps) {
  const sectionLessons = useMemo(() => {
    return lessons.filter((lesson) => (lesson.section_id ?? lesson.module) === sectionID);
  }, [lessons, sectionID]);

  const sectionTitle = sectionLessons[0]?.section_title ?? sectionID;

  return (
    <section className="panel section-page-panel">
      <div className="control-row">
        <button type="button" onClick={() => onNavigate("/")}>Back to Overview</button>
      </div>
      <h2>{sectionTitle}</h2>
      <p className="subtitle">Select a lesson and launch the challenge workflow.</p>

      {sectionLessons.length === 0 ? (
        <p className="empty">No lessons found for this section.</p>
      ) : (
        <div className="lesson-card-grid">
          {sectionLessons.map((lesson) => {
            const unlockedStages = lesson.stages.filter((stage) => stage.unlocked !== false);
            const firstUnlocked = unlockedStages[0]?.index ?? lesson.stages[0]?.index ?? 0;
            const completedStages = lesson.stages.filter((stage) => stage.completed).length;
            return (
              <article key={lesson.id} className="lesson-card">
                <h3>{lesson.title}</h3>
                <p className="lesson-outcome">Difficulty: {lesson.difficulty ?? "intermediate"}</p>
                <p className="lesson-outcome">
                  Estimated time: {lesson.estimated_minutes ?? 20} minutes
                </p>
                <p className="lesson-outcome">
                  Stages completed: {completedStages}/{lesson.stages.length}
                </p>
                <p className="lesson-outcome">
                  Chapter refs: {(lesson.chapter_refs ?? []).join(", ") || "n/a"}
                </p>
                <div className="control-row">
                  <button
                    type="button"
                    onClick={() => onNavigate(`/challenge/${lesson.id}/${firstUnlocked}`)}
                  >
                    Open Challenge
                  </button>
                </div>
              </article>
            );
          })}
        </div>
      )}
    </section>
  );
}
