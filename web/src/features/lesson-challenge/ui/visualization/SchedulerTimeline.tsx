import type { SnapshotDTO } from "../../../../lib/types";
import { cn } from "../../../../shared/lib/cn";
import { selectTimelineWindow } from "../../../../state/selectors";

interface SchedulerTimelineProps {
  snapshot: SnapshotDTO | null;
}

export function SchedulerTimeline({ snapshot }: SchedulerTimelineProps) {
  const cells = selectTimelineWindow(snapshot, 48);

  return (
    <section className="rounded-xl border border-slate-200 bg-white p-3">
      <h2 className="text-sm font-semibold text-slate-900">Scheduler Timeline</h2>
      {cells.length === 0 ? (
        <p className="mt-2 text-sm text-slate-600">Run or step a session to populate timeline slices.</p>
      ) : (
        <>
          <div
            className="mt-2 grid grid-cols-[repeat(auto-fill,minmax(56px,1fr))] gap-1.5"
            role="list"
            aria-label="scheduler timeline"
          >
            {cells.map((cell) => (
              <div
                key={`${cell.tick}-${cell.pid}`}
                className={cn(
                  "rounded-md border border-slate-200 bg-slate-50 p-1.5",
                  cell.pid === 0 ? "bg-slate-50" : "bg-slate-100",
                )}
                data-state={cell.pid === 0 ? "idle" : "active"}
                title={`tick=${cell.tick} ${cell.label}`}
                role="listitem"
              >
                <span className="text-xs text-slate-500">{cell.tick}</span>
                <span className="mt-0.5 block font-semibold text-slate-900">{cell.label}</span>
              </div>
            ))}
          </div>
          <div className="mt-2 flex flex-wrap gap-3 text-xs text-slate-500">
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
