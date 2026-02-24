import type { LogEntry } from "../state/sessionReducer";

interface EventLogProps {
  logs: LogEntry[];
}

export function EventLog({ logs }: EventLogProps) {
  return (
    <section className="panel log-panel">
      <h2>Event Log</h2>
      <ul>
        {logs.length === 0 ? (
          <li>No events yet.</li>
        ) : (
          logs
            .slice()
            .reverse()
            .map((entry) => (
              <li key={entry.id}>
                <span className="seq">#{entry.sequence}</span>
                <span>{entry.type}</span>
                <span>tick={entry.tick}</span>
                <span>{entry.detail}</span>
              </li>
            ))
        )}
      </ul>
    </section>
  );
}
