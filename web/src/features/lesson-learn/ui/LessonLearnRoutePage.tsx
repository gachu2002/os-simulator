import { useNavigate, useParams, useSearchParams } from "react-router-dom";

import { LessonLearnPage } from "./LessonLearnPage";
import { apiBaseURL } from "../../../shared/config/env";

export function LessonLearnRoutePage() {
  const navigate = useNavigate();
  const { lessonID = "" } = useParams<{ lessonID: string }>();
  const [searchParams] = useSearchParams();
  const stageParam = searchParams.get("stage");
  const stage = stageParam === null ? undefined : Number(stageParam);

  return (
    <LessonLearnPage
      baseURL={apiBaseURL()}
      lessonID={lessonID}
      preferredStageIndex={Number.isFinite(stage) ? stage : undefined}
      onNavigate={(to) => navigate(to)}
    />
  );
}
