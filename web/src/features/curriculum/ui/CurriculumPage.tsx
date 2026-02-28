import { useNavigate } from "react-router-dom";

import { OverviewPage } from "./OverviewPage";
import { useLessonsCatalog } from "../hooks/useLessonsCatalog";
import { apiBaseURL } from "../../../shared/config/env";

export function CurriculumPage() {
  const navigate = useNavigate();
  const { sections, isLoading, errorMessage } = useLessonsCatalog({ baseURL: apiBaseURL() });

  return (
    <OverviewPage
      sections={sections}
      isLoading={isLoading}
      errorMessage={errorMessage}
      onNavigate={(to) => navigate(to)}
    />
  );
}
