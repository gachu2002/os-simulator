import { useNavigate, useParams, useSearchParams } from "react-router-dom";

import { LessonChallengePage } from "./LessonChallengePage";
import { apiBaseURL } from "../../../shared/config/env";

export function LessonChallengeRoutePage() {
  const navigate = useNavigate();
  const { lessonID = "" } = useParams<{ lessonID: string }>();
  const [searchParams] = useSearchParams();
  const stage = Number(searchParams.get("stage") ?? 0);

  return (
    <LessonChallengePage
      baseURL={apiBaseURL()}
      lessonID={lessonID}
      stageIndex={Number.isFinite(stage) ? stage : 0}
      onNavigate={(to) => navigate(to)}
    />
  );
}
