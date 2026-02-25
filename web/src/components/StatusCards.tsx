import type { SnapshotDTO } from "../lib/types";

interface StatusCardsProps {
  connected: boolean;
  sessionID: string;
  snapshot: SnapshotDTO | null;
}

export function StatusCards({
  connected,
  sessionID,
  snapshot,
}: StatusCardsProps) {
  return (
    <section className="status-grid">
      <article className="panel stat">
        <h2>Session</h2>
        <p>{sessionID || "not created"}</p>
      </article>
      <article className="panel stat">
        <h2>Connection</h2>
        <p>{connected ? "connected" : "disconnected"}</p>
      </article>
      <article className="panel stat">
        <h2>Tick</h2>
        <p>{snapshot?.tick ?? 0}</p>
      </article>
      <article className="panel stat">
        <h2>Trace Hash</h2>
        <p className="hash">{snapshot?.trace_hash ?? "-"}</p>
      </article>
      <article className="panel stat">
        <h2>Policy</h2>
        <p>{snapshot?.metrics.policy ?? "rr"}</p>
      </article>
      <article className="panel stat">
        <h2>Completed</h2>
        <p>{snapshot?.metrics.completed_processes ?? 0}</p>
      </article>
    </section>
  );
}
