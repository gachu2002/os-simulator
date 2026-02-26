import { useCallback, useEffect, useState } from "react";

import { LessonRunnerPanel } from "./components/LessonRunnerPanel";
import { OverviewPage } from "./components/OverviewPage";
import { SectionPage } from "./components/SectionPage";
import { VisualizationSuite } from "./components/VisualizationSuite";
import { useLessonsCatalog } from "./hooks/useLessonsCatalog";
import type { ChallengeGradeResponse } from "./lib/lessonApi";
import type { SnapshotDTO } from "./lib/types";
import { snapshotFromChallengeGrade } from "./state/lessonSnapshot";

type AppRoute =
  | { kind: "overview" }
  | { kind: "section"; sectionID: string }
  | { kind: "challenge"; lessonID?: string; stageIndex?: number };

export function App() {
  const [baseURL, setBaseURL] = useState(defaultBaseURL());
  const [routePath, setRoutePath] = useState(() =>
    typeof window === "undefined" ? "/" : window.location.pathname,
  );

  const [lessonSnapshot, setLessonSnapshot] = useState<SnapshotDTO | null>(null);
  const [lessonTitle, setLessonTitle] = useState("");
  const [hasGradedAttempt, setHasGradedAttempt] = useState(false);
  const { lessons, isLoading, errorMessage } = useLessonsCatalog({ baseURL });

  const route = parseRoute(routePath);

  const handleNavigate = useCallback((to: string) => {
    if (typeof window === "undefined") {
      return;
    }
    window.history.pushState({}, "", to);
    setRoutePath(window.location.pathname);
  }, []);

  useEffect(() => {
    if (typeof window === "undefined") {
      return;
    }
    const handlePopState = () => setRoutePath(window.location.pathname);
    window.addEventListener("popstate", handlePopState);
    return () => window.removeEventListener("popstate", handlePopState);
  }, []);

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
        <h1>OSTEP Simulator Course</h1>
        <p className="subtitle">
          Learn operating systems through sectioned lessons and interactive deterministic
          challenges.
        </p>
        <div className="control-row">
          <button type="button" onClick={() => handleNavigate("/")}>Overview</button>
          <button type="button" onClick={() => handleNavigate("/challenge")}>Challenge</button>
        </div>
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

      {route.kind === "overview" ? (
        <OverviewPage
          lessons={lessons}
          isLoading={isLoading}
          errorMessage={errorMessage}
          onNavigate={handleNavigate}
        />
      ) : null}

      {route.kind === "section" ? (
        <SectionPage lessons={lessons} sectionID={route.sectionID} onNavigate={handleNavigate} />
      ) : null}

      {route.kind === "challenge" ? (
        <>
          <LessonRunnerPanel
            baseURL={baseURL}
            onLiveSnapshot={handleChallengeLiveSnapshot}
            onGradeResult={handleChallengeGradeResult}
            preferredLessonID={route.lessonID}
            preferredStageIndex={route.stageIndex}
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
        </>
      ) : null}
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

function parseRoute(pathname: string): AppRoute {
  const normalized = pathname.length > 1 ? pathname.replace(/\/+$/, "") : pathname;
  if (normalized === "/") {
    return { kind: "overview" };
  }
  if (normalized === "/challenge") {
    return { kind: "challenge" };
  }

  const sectionPrefix = "/sections/";
  if (normalized.startsWith(sectionPrefix)) {
    const sectionID = decodeURIComponent(normalized.slice(sectionPrefix.length));
    if (sectionID !== "") {
      return { kind: "section", sectionID };
    }
  }

  const challengePrefix = "/challenge/";
  if (normalized.startsWith(challengePrefix)) {
    const parts = normalized.slice(challengePrefix.length).split("/");
    const lessonID = decodeURIComponent(parts[0] ?? "").trim();
    const parsedStage = Number(parts[1]);
    return {
      kind: "challenge",
      lessonID: lessonID === "" ? undefined : lessonID,
      stageIndex: Number.isFinite(parsedStage) ? parsedStage : 0,
    };
  }

  return { kind: "overview" };
}
