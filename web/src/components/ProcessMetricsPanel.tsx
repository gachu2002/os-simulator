import type { SnapshotDTO } from "../lib/types";
import { selectProcessMetricRows } from "../state/selectors";

interface ProcessMetricsPanelProps {
  snapshot: SnapshotDTO | null;
}

export function ProcessMetricsPanel({ snapshot }: ProcessMetricsPanelProps) {
  const rows = selectProcessMetricRows(snapshot);

  return (
    <section className="panel metrics-panel">
      <h2>Process Metrics</h2>
      {rows.length === 0 ? (
        <p className="empty">No process metrics yet.</p>
      ) : (
        <div className="table-wrap">
          <table>
            <thead>
              <tr>
                <th>PID</th>
                <th>Name</th>
                <th>State</th>
                <th>PC</th>
                <th>Resp</th>
                <th>Turn</th>
                <th>Run</th>
                <th>Wait</th>
              </tr>
            </thead>
            <tbody>
              {rows.map((row) => (
                <tr key={row.pid}>
                  <td>{row.pid}</td>
                  <td>{row.name}</td>
                  <td>{row.state}</td>
                  <td>{row.pc}</td>
                  <td>{row.responseTime}</td>
                  <td>{row.turnaround}</td>
                  <td>{row.runTicks}</td>
                  <td>{row.waitTicks}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </section>
  );
}
