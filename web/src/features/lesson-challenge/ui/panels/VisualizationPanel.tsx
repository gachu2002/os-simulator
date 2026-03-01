import type { SnapshotDTO } from "../../../../lib/types";
import { LessonSpecificVisualization } from "../visualization/LessonSpecificVisualization";

export function VisualizationPanel({
  lessonID,
  snapshot,
  lastLessonAction,
}: {
  lessonID: string;
  snapshot: SnapshotDTO | null;
  lastLessonAction: string;
}) {
  return (
    <section className="min-h-0 rounded-lg border border-slate-200 bg-white p-3">
      <h3 className="text-sm font-semibold text-slate-900">Visualization</h3>
      <LessonSpecificVisualization
        lessonID={lessonID}
        snapshot={snapshot}
        lastLessonAction={lastLessonAction}
      />
    </section>
  );
}
