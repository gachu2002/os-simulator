import type { SnapshotDTO } from "../lib/types";
import { selectProcessQueues } from "../state/selectors";

interface ProcessQueuesPanelProps {
  snapshot: SnapshotDTO | null;
}

export function ProcessQueuesPanel({ snapshot }: ProcessQueuesPanelProps) {
  const queues = selectProcessQueues(snapshot);

  return (
    <section className="panel process-panel">
      <h2>Process Queues</h2>
      <div className="queue-grid">
        <QueueColumn title="Running" items={queues.running} tone="running" />
        <QueueColumn title="Ready" items={queues.ready} tone="ready" />
        <QueueColumn title="Blocked" items={queues.blocked} tone="blocked" />
        <QueueColumn
          title="Terminated"
          items={queues.terminated}
          tone="terminated"
        />
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
    <article className="queue-column" data-tone={tone}>
      <header>
        <h3>{title}</h3>
        <span>{items.length}</span>
      </header>
      {items.length === 0 ? (
        <p className="empty">None</p>
      ) : (
        <ul>
          {items.map((item) => (
            <li key={item}>{item}</li>
          ))}
        </ul>
      )}
    </article>
  );
}
