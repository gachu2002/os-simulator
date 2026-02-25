import type { SnapshotDTO } from "../lib/types";
import { selectTimelineWindow } from "../state/selectors";

interface SchedulerTimelineProps {
  snapshot: SnapshotDTO | null;
}

export function SchedulerTimeline({ snapshot }: SchedulerTimelineProps) {
  const cells = selectTimelineWindow(snapshot, 48);

  return (
    <section className="panel timeline-panel">
      <h2>Scheduler Timeline</h2>
      {cells.length === 0 ? (
        <p className="empty">
          Run or step a session to populate timeline slices.
        </p>
      ) : (
        <>
          <div
            className="timeline-grid"
            role="list"
            aria-label="scheduler timeline"
          >
            {cells.map((cell) => (
              <div
                key={`${cell.tick}-${cell.pid}`}
                className="timeline-cell"
                data-state={cell.pid === 0 ? "idle" : "active"}
                title={`tick=${cell.tick} ${cell.label}`}
                role="listitem"
              >
                <span className="tick">{cell.tick}</span>
                <span className="pid">{cell.label}</span>
              </div>
            ))}
          </div>
          <div className="timeline-legend">
            <span>
              <strong>Active</strong>: running PID slice
            </span>
            <span>
              <strong>Idle</strong>: no runnable process
            </span>
          </div>
        </>
      )}
    </section>
  );
}
