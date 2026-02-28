import type { SnapshotDTO } from "../../../../lib/types";
import { cn } from "../../../../shared/lib/cn";
import { selectProcessQueues } from "../../../../state/selectors";

interface ProcessQueuesPanelProps {
  snapshot: SnapshotDTO | null;
}

export function ProcessQueuesPanel({ snapshot }: ProcessQueuesPanelProps) {
  const queues = selectProcessQueues(snapshot);

  return (
    <section className="rounded-xl border border-slate-200 bg-white p-3 md:col-span-2">
      <h2 className="text-sm font-semibold text-slate-900">Process Queues</h2>
      <div className="mt-2 grid gap-2 md:grid-cols-4">
        <QueueColumn title="Running" items={queues.running} tone="running" />
        <QueueColumn title="Ready" items={queues.ready} tone="ready" />
        <QueueColumn title="Blocked" items={queues.blocked} tone="blocked" />
        <QueueColumn title="Terminated" items={queues.terminated} tone="terminated" />
      </div>
    </section>
  );
}

interface QueueColumnProps {
  title: string;
  items: string[];
  tone: "running" | "ready" | "blocked" | "terminated";
}

function QueueColumn({ title, items, tone }: QueueColumnProps) {
  return (
    <article
      className={cn(
        "rounded-md border p-2",
        tone === "running" && "border-emerald-200 bg-emerald-50/40",
        tone === "ready" && "border-blue-200 bg-blue-50/40",
        tone === "blocked" && "border-amber-200 bg-amber-50/40",
        tone === "terminated" && "border-slate-200 bg-slate-50/70",
      )}
      data-tone={tone}
    >
      <header className="mb-2 flex items-center justify-between">
        <h3 className="text-sm font-semibold text-slate-900">{title}</h3>
        <span className="text-xs text-slate-500">{items.length}</span>
      </header>
      {items.length === 0 ? (
        <p className="text-sm text-slate-600">None</p>
      ) : (
        <ul className="space-y-1">
          {items.map((item) => (
            <li key={item} className="rounded border border-slate-200 bg-white px-2 py-1 text-sm">
              {item}
            </li>
          ))}
        </ul>
      )}
    </article>
  );
}
