import { useNavigate, useParams, useSearchParams } from "react-router-dom";

import { LessonChallengePage } from "./LessonChallengePage";
import { apiBaseURL } from "../../../shared/config/env";

export function LessonChallengeRoutePage() {
  const navigate = useNavigate();
  const { lessonID = "" } = useParams<{ lessonID: string }>();
  const [searchParams] = useSearchParams();
  const stageParam = searchParams.get("stage");
  const stage = stageParam === null ? undefined : Number(stageParam);

  return (
    <LessonChallengePage
      baseURL={apiBaseURL()}
      lessonID={lessonID}
      stageIndex={Number.isFinite(stage) ? stage : undefined}
      onNavigate={(to) => navigate(to)}
    />
  );
}
