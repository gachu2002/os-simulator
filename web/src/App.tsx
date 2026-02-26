import { useCallback, useState } from "react";

import { LessonRunnerPanel } from "./components/LessonRunnerPanel";
import { VisualizationSuite } from "./components/VisualizationSuite";
import type { ChallengeGradeResponse } from "./lib/lessonApi";
import type { SnapshotDTO } from "./lib/types";
import { snapshotFromChallengeGrade } from "./state/lessonSnapshot";

export function App() {
  const [baseURL, setBaseURL] = useState(defaultBaseURL());
  const [lessonSnapshot, setLessonSnapshot] = useState<SnapshotDTO | null>(null);
  const [lessonTitle, setLessonTitle] = useState("");
  const [hasGradedAttempt, setHasGradedAttempt] = useState(false);

  const handleChallengeLiveSnapshot = useCallback((snapshot: SnapshotDTO | null, title: string) => {
    setLessonSnapshot(snapshot);
    setLessonTitle(title);
    setHasGradedAttempt(false);
  }, []);

  const handleChallengeGradeResult = useCallback((result: ChallengeGradeResponse) => {
    setLessonSnapshot(snapshotFromChallengeGrade(result));
    setLessonTitle(`${result.lesson_id} stage ${result.stage_index + 1}`);
    setHasGradedAttempt(true);
  }, []);

  return (
    <main className="app-shell">
      <header className="hero">
        <p className="eyebrow">Challenge</p>
        <h1>Solve OSTEP Lesson Stages</h1>
        <p className="subtitle">
          Complete focused lesson stages with hints and deterministic grading feedback.
        </p>
      </header>

      <section className="panel challenge-setup-panel">
        <h2>Challenge Setup</h2>
        <label>
          Server URL
          <input
            value={baseURL}
            onChange={(event) => setBaseURL(event.target.value)}
            placeholder="http://127.0.0.1:8080"
          />
        </label>
        <p className="challenge-note">
          Use Learn to understand the concept, then switch to Exercise to run actions and
          validate with deterministic checks.
        </p>
      </section>

      <LessonRunnerPanel
        baseURL={baseURL}
        onLiveSnapshot={handleChallengeLiveSnapshot}
        onGradeResult={handleChallengeGradeResult}
      />

      <VisualizationSuite
        title="Lesson Stage Snapshot"
        subtitle={
          lessonSnapshot
            ? hasGradedAttempt
              ? `Latest graded result: ${lessonTitle}`
              : `Live lesson stage attempt: ${lessonTitle}`
            : "Start a lesson stage to render scheduler, memory, and process views"
        }
        snapshot={lessonSnapshot}
      />
    </main>
  );
}

function defaultBaseURL(): string {
  const envURL = import.meta.env.VITE_API_BASE_URL;
  if (typeof envURL === "string" && envURL.trim() !== "") {
    return envURL.trim();
  }
  if (typeof window === "undefined") {
    return "http://localhost:8080";
  }
  return window.location.origin;
}
