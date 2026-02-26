import { useQuery } from "@tanstack/react-query";
import { useMemo, useState } from "react";

import { fetchLessonsForLearner } from "../lib/lessonApi";
import { getOrCreateLearnerID } from "../lib/learner";

interface UseLessonsCatalogOptions {
  baseURL: string;
}

export function useLessonsCatalog({ baseURL }: UseLessonsCatalogOptions) {
  const [learnerID] = useState(() => getOrCreateLearnerID());
  const query = useQuery({
    queryKey: ["lessons-catalog", baseURL, learnerID],
    queryFn: () => fetchLessonsForLearner(baseURL, learnerID),
  });

  return {
    lessons: useMemo(() => query.data ?? [], [query.data]),
    isLoading: query.isLoading,
    errorMessage: query.error instanceof Error ? query.error.message : "",
  };
}
