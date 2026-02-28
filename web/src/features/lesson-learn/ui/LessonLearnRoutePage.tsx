import { useNavigate, useParams, useSearchParams } from "react-router-dom";

import { LessonLearnPage } from "./LessonLearnPage";
import { apiBaseURL } from "../../../shared/config/env";

export function LessonLearnRoutePage() {
  const navigate = useNavigate();
  const { lessonID = "" } = useParams<{ lessonID: string }>();
  const [searchParams] = useSearchParams();
  const stage = Number(searchParams.get("stage") ?? 0);

  return (
    <LessonLearnPage
      baseURL={apiBaseURL()}
      lessonID={lessonID}
      preferredStageIndex={Number.isFinite(stage) ? stage : 0}
      onNavigate={(to) => navigate(to)}
    />
  );
}
