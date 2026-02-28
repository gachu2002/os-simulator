import { Navigate, Route, Routes } from "react-router-dom";

import { AppLayout } from "./layout/AppLayout";
import { CurriculumPage } from "../features/curriculum/ui/CurriculumPage";
import { LessonChallengeRoutePage } from "../features/lesson-challenge/ui/LessonChallengeRoutePage";
import { LessonLearnRoutePage } from "../features/lesson-learn/ui/LessonLearnRoutePage";

export function AppRouter() {
  return (
    <Routes>
      <Route element={<AppLayout />}>
        <Route index element={<CurriculumPage />} />
        <Route path="lesson/:lessonID/learn" element={<LessonLearnRoutePage />} />
        <Route path="lesson/:lessonID/challenge" element={<LessonChallengeRoutePage />} />
        <Route path="*" element={<Navigate to="/" replace />} />
      </Route>
    </Routes>
  );
}
