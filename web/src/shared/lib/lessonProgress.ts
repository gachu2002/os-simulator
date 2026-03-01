interface LessonProgressStore {
  completedLessonIDs: string[];
}

const LESSON_PROGRESS_PREFIX = "lesson-progress-v1:";

export function isLessonCompleted(learnerID: string, lessonID: string): boolean {
  const progress = loadLessonProgress(learnerID);
  return progress.completedLessonIDs.includes(lessonID);
}

export function markLessonCompleted(learnerID: string, lessonID: string): void {
  const progress = loadLessonProgress(learnerID);
  if (progress.completedLessonIDs.includes(lessonID)) {
    return;
  }
  const next: LessonProgressStore = {
    completedLessonIDs: [...progress.completedLessonIDs, lessonID].sort(),
  };
  window.localStorage.setItem(storeKey(learnerID), JSON.stringify(next));
}

export function completedLessonCount(learnerID: string): number {
  return loadLessonProgress(learnerID).completedLessonIDs.length;
}

function loadLessonProgress(learnerID: string): LessonProgressStore {
  const raw = window.localStorage.getItem(storeKey(learnerID));
  if (!raw) {
    return { completedLessonIDs: [] };
  }
  try {
    const parsed = JSON.parse(raw) as LessonProgressStore;
    if (!Array.isArray(parsed.completedLessonIDs)) {
      return { completedLessonIDs: [] };
    }
    return {
      completedLessonIDs: parsed.completedLessonIDs.filter((item) => typeof item === "string"),
    };
  } catch {
    return { completedLessonIDs: [] };
  }
}

function storeKey(learnerID: string): string {
  return `${LESSON_PROGRESS_PREFIX}${learnerID}`;
}
