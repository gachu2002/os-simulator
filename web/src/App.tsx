import { useCallback, useEffect, useState } from "react";

import { LessonRunnerPanel } from "./components/LessonRunnerPanel";
import { OverviewPage } from "./components/OverviewPage";
import { VisualizationSuite } from "./components/VisualizationSuite";
import { useLessonsCatalog } from "./hooks/useLessonsCatalog";
import type { ChallengeGradeResponse } from "./lib/lessonApi";
import type { SnapshotDTO } from "./lib/types";
import { snapshotFromChallengeGrade } from "./state/lessonSnapshot";

type AppRoute =
  | { kind: "overview" }
  | { kind: "challenge"; lessonID?: string; stageIndex?: number };

export function App() {
  const [baseURL] = useState(defaultBaseURL());
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
      {route.kind === "challenge" ? (
        <div className="top-nav">
          <button type="button" className="btn btn-ghost" onClick={() => handleNavigate("/")}>
            Home
          </button>
        </div>
      ) : null}

      <header className="hero">
        <p className="eyebrow">Challenge</p>
        <h1>OSTEP Simulator Course</h1>
        <p className="subtitle">
          Learn operating systems through sectioned lessons and interactive deterministic
          challenges.
        </p>
      </header>

      {route.kind === "overview" ? (
        <OverviewPage
          lessons={lessons}
          isLoading={isLoading}
          errorMessage={errorMessage}
          onNavigate={handleNavigate}
        />
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

  const challengePrefix = "/challenge/";
  if (normalized.startsWith(challengePrefix)) {
    const parts = normalized.slice(challengePrefix.length).split("/");
    const lessonID = decodeURIComponent(parts[0] ?? "").trim();
    if (lessonID === "") {
      return { kind: "overview" };
    }
    const parsedStage = Number(parts[1]);
    return {
      kind: "challenge",
      lessonID,
      stageIndex: Number.isFinite(parsedStage) ? parsedStage : 0,
    };
  }

  return { kind: "overview" };
}
