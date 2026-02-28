import type { SnapshotDTO } from "../../../../lib/types";
import { selectProcessMetricRows } from "../../../../state/selectors";

interface ProcessMetricsPanelProps {
  snapshot: SnapshotDTO | null;
}

export function ProcessMetricsPanel({ snapshot }: ProcessMetricsPanelProps) {
  const rows = selectProcessMetricRows(snapshot);

  return (
    <section className="rounded-xl border border-slate-200 bg-white p-3 md:col-span-2">
      <h2 className="text-sm font-semibold text-slate-900">Process Metrics</h2>
      {rows.length === 0 ? (
        <p className="mt-2 text-sm text-slate-600">No process metrics yet.</p>
      ) : (
        <div className="mt-2 overflow-auto">
          <table className="w-full border-collapse text-sm">
            <thead>
              <tr>
                <th className="border-b border-slate-200 p-1.5 text-left font-semibold">PID</th>
                <th className="border-b border-slate-200 p-1.5 text-left font-semibold">Name</th>
                <th className="border-b border-slate-200 p-1.5 text-left font-semibold">State</th>
                <th className="border-b border-slate-200 p-1.5 text-left font-semibold">PC</th>
                <th className="border-b border-slate-200 p-1.5 text-left font-semibold">Resp</th>
                <th className="border-b border-slate-200 p-1.5 text-left font-semibold">Turn</th>
                <th className="border-b border-slate-200 p-1.5 text-left font-semibold">Run</th>
                <th className="border-b border-slate-200 p-1.5 text-left font-semibold">Wait</th>
              </tr>
            </thead>
            <tbody>
              {rows.map((row) => (
                <tr key={row.pid}>
                  <td className="border-b border-slate-100 p-1.5">{row.pid}</td>
                  <td className="border-b border-slate-100 p-1.5">{row.name}</td>
                  <td className="border-b border-slate-100 p-1.5">{row.state}</td>
                  <td className="border-b border-slate-100 p-1.5">{row.pc}</td>
                  <td className="border-b border-slate-100 p-1.5">{row.responseTime}</td>
                  <td className="border-b border-slate-100 p-1.5">{row.turnaround}</td>
                  <td className="border-b border-slate-100 p-1.5">{row.runTicks}</td>
                  <td className="border-b border-slate-100 p-1.5">{row.waitTicks}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </section>
  );
}
