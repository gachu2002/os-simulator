import type { SnapshotDTO } from "../../../../lib/types";
import { VisualizationSuite } from "../visualization/VisualizationSuite";

export function VisualizationPanel({ snapshot }: { snapshot: SnapshotDTO | null }) {
  return (
    <section className="mt-3 rounded-lg border border-slate-200 bg-slate-50 p-3">
      <h3 className="text-sm font-semibold text-slate-900">2) Visualization</h3>
      <VisualizationSuite
        title="Live Challenge State"
        subtitle="Run actions and inspect trace, memory, process queues, and metrics."
        snapshot={snapshot}
      />
    </section>
  );
}
