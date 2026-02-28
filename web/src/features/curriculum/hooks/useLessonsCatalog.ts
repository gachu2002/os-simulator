import { useQuery } from "@tanstack/react-query";
import { useMemo, useState } from "react";

import { getOrCreateLearnerID } from "../../../lib/learner";
import { fetchCurriculumForLearner } from "../api/curriculumApi";

interface UseLessonsCatalogOptions {
  baseURL: string;
}

export function useLessonsCatalog({ baseURL }: UseLessonsCatalogOptions) {
  const [learnerID] = useState(() => getOrCreateLearnerID());
  const query = useQuery({
    queryKey: ["lessons-catalog", baseURL, learnerID],
    queryFn: () => fetchCurriculumForLearner(baseURL, learnerID),
  });

  return {
    sections: useMemo(() => query.data ?? [], [query.data]),
    isLoading: query.isLoading,
    errorMessage: query.error instanceof Error ? query.error.message : "",
  };
}
