import type { SnapshotDTO } from "../../../../lib/types";

import { MemoryPanel } from "./MemoryPanel";
import { ProcessMetricsPanel } from "./ProcessMetricsPanel";
import { ProcessQueuesPanel } from "./ProcessQueuesPanel";
import { SchedulerTimeline } from "./SchedulerTimeline";

interface VisualizationSuiteProps {
  title: string;
  subtitle: string;
  snapshot: SnapshotDTO | null;
}

export function VisualizationSuite({
  title,
  subtitle,
  snapshot,
}: VisualizationSuiteProps) {
  return (
    <section className="mt-2 grid gap-3 rounded-xl border border-slate-200 bg-white p-4">
      <div className="border-b border-slate-200 pb-2">
        <h2 className="text-sm font-semibold text-slate-900">{title}</h2>
        <p className="mt-1 text-sm text-slate-600">{subtitle}</p>
      </div>
      <div className="grid gap-3 md:grid-cols-[1.45fr_1fr]">
        <SchedulerTimeline snapshot={snapshot} />
        <MemoryPanel snapshot={snapshot} />
        <ProcessQueuesPanel snapshot={snapshot} />
        <ProcessMetricsPanel snapshot={snapshot} />
      </div>
    </section>
  );
}
