import type { SnapshotDTO } from "../lib/types";

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
    <section className="viz-suite">
      <div className="viz-suite-header">
        <h2>{title}</h2>
        <p>{subtitle}</p>
      </div>
      <div className="viz-grid">
        <SchedulerTimeline snapshot={snapshot} />
        <MemoryPanel snapshot={snapshot} />
        <ProcessQueuesPanel snapshot={snapshot} />
        <ProcessMetricsPanel snapshot={snapshot} />
      </div>
    </section>
  );
}
