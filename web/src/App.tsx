import { useCallback, useEffect, useState } from "react";

import { LessonChallengePage } from "./components/LessonChallengePage";
import { LessonLearnPage } from "./components/LessonLearnPage";
import { OverviewPage } from "./components/OverviewPage";
import { useLessonsCatalog } from "./hooks/useLessonsCatalog";
import { parseRoute } from "./lib/routes";

export function App() {
  const [baseURL] = useState(defaultBaseURL());
  const [routePath, setRoutePath] = useState(() =>
    typeof window === "undefined"
      ? "/"
      : `${window.location.pathname}${window.location.search}`,
  );

  const { sections, isLoading, errorMessage } = useLessonsCatalog({ baseURL });

  const route = parseRoute(routePath);

  const handleNavigate = useCallback((to: string) => {
    if (typeof window === "undefined") {
      return;
    }
    window.history.pushState({}, "", to);
    setRoutePath(`${window.location.pathname}${window.location.search}`);
  }, []);

  useEffect(() => {
    if (typeof window === "undefined") {
      return;
    }
    const handlePopState = () =>
      setRoutePath(`${window.location.pathname}${window.location.search}`);
    window.addEventListener("popstate", handlePopState);
    return () => window.removeEventListener("popstate", handlePopState);
  }, []);

  return (
    <main className="app-shell">
      {route.kind !== "overview" ? (
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
          sections={sections}
          isLoading={isLoading}
          errorMessage={errorMessage}
          onNavigate={handleNavigate}
        />
      ) : null}

      {route.kind === "learn" ? (
        <LessonLearnPage
          baseURL={baseURL}
          lessonID={route.lessonID}
          preferredStageIndex={route.stageIndex}
          onNavigate={handleNavigate}
        />
      ) : null}

      {route.kind === "challenge" ? (
        <LessonChallengePage
          baseURL={baseURL}
          lessonID={route.lessonID}
          stageIndex={route.stageIndex}
          onNavigate={handleNavigate}
        />
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
