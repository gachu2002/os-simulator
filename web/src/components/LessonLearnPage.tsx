import { useQuery } from "@tanstack/react-query";
import { useMemo, useState } from "react";

import { fetchLessonLearn } from "../lib/lessonApi";
import { getOrCreateLearnerID } from "../lib/learner";

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
    <section className="panel lesson-panel">
      <div className="top-nav">
        <button type="button" className="btn btn-ghost" onClick={() => onNavigate("/")}>Home</button>
      </div>

      <h2>{lesson?.title ?? "Lesson"}</h2>
      <p className="subtitle">Learn page: theory only. Review this before opening the challenge.</p>

      {query.isLoading ? <p className="empty">Loading lesson theory...</p> : null}
      {query.error instanceof Error ? <p className="error">{query.error.message}</p> : null}

      {stage ? (
        <>
          <section className="lesson-learn-block">
            <h3>Core Idea</h3>
            <p className="lesson-outcome">
              {stage.core_idea ?? "Review the core operating systems concept for this stage."}
            </p>
          </section>

          <section className="lesson-learn-block">
            <h3>Mechanism</h3>
            {(stage.mechanism_steps ?? []).length > 0 ? (
              (stage.mechanism_steps ?? []).map((item) => (
                <p key={item} className="lesson-outcome">- {item}</p>
              ))
            ) : (
              <p className="lesson-outcome">
                Track cause-and-effect transitions in the trace and state panels.
              </p>
            )}
          </section>

          <section className="lesson-learn-block">
            <h3>Worked Example</h3>
            <p className="lesson-outcome">
              {stage.worked_example ??
                `Stage objective: ${stage.objective ?? stage.title}. Focus on how state transitions satisfy the objective.`}
            </p>
          </section>

          <section className="lesson-learn-block">
            <h3>Common Mistakes</h3>
            {(stage.common_mistakes ?? []).length > 0 ? (
              (stage.common_mistakes ?? []).map((item) => (
                <p key={item} className="lesson-outcome">- {item}</p>
              ))
            ) : (
              <>
                <p className="lesson-outcome">Do not infer correctness from one metric alone; always confirm trace and state evidence.</p>
                <p className="lesson-outcome">Do not ignore deterministic limits; step and config budgets are part of the challenge contract.</p>
                <p className="lesson-outcome">Do not submit before checking expected visual cues against observed behavior.</p>
              </>
            )}
          </section>

          <section className="lesson-learn-block">
            <h3>What To Watch In Challenge</h3>
            {(stage.pre_challenge_checklist ?? []).length ? (
              (stage.pre_challenge_checklist ?? []).map((item) => (
                <p key={item} className="lesson-outcome">- {item}</p>
              ))
            ) : (stage.expected_visual_cues ?? []).length ? (
              (stage.expected_visual_cues ?? []).map((item) => (
                <p key={item} className="lesson-outcome">- {item}</p>
              ))
            ) : (
              <p className="lesson-outcome">Watch trace order, process progress, and objective-aligned metrics.</p>
            )}
          </section>

          <div className="lesson-controls">
            <button
              type="button"
              className="btn btn-primary"
              onClick={() => onNavigate(`/lesson/${lessonID}/challenge?stage=${stage.index}`)}
            >
              Go To Challenge
            </button>
          </div>
        </>
      ) : null}
    </section>
  );
}
