import type { SnapshotDTO } from "../../../../lib/types";
import { selectFrameRows } from "../../../../state/selectors";

interface MemoryPanelProps {
  snapshot: SnapshotDTO | null;
}

export function MemoryPanel({ snapshot }: MemoryPanelProps) {
  const rows = selectFrameRows(snapshot);
  const faults = snapshot?.memory.faults;
  const tlbSlots = snapshot?.memory.tlb?.length ?? 0;

  return (
    <section className="rounded-xl border border-slate-200 bg-white p-3">
      <h2 className="text-sm font-semibold text-slate-900">Memory View</h2>
      <div className="mt-2 flex flex-wrap gap-3 text-xs text-slate-500">
        <span>Frames: {snapshot?.memory.total_frames ?? 0}</span>
        <span>TLB slots: {tlbSlots}</span>
        <span>Faults NP: {faults?.not_present ?? 0}</span>
        <span>Faults Perm: {faults?.permission ?? 0}</span>
        <span>
          TLB hit/miss: {faults?.tlb_hit ?? 0}/{faults?.tlb_miss ?? 0}
        </span>
      </div>

      {rows.length === 0 ? (
        <p className="mt-2 text-sm text-slate-600">No frame allocations yet.</p>
      ) : (
        <div className="mt-2 overflow-auto">
          <table className="w-full border-collapse text-sm">
            <thead>
              <tr>
                <th className="border-b border-slate-200 p-1.5 text-left font-semibold">Frame</th>
                <th className="border-b border-slate-200 p-1.5 text-left font-semibold">Owner</th>
                <th className="border-b border-slate-200 p-1.5 text-left font-semibold">VPN</th>
              </tr>
            </thead>
            <tbody>
              {rows.map((row) => (
                <tr key={row.frame}>
                  <td className="border-b border-slate-100 p-1.5">{row.frame}</td>
                  <td className="border-b border-slate-100 p-1.5">{row.owner}</td>
                  <td className="border-b border-slate-100 p-1.5">{row.vpn}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </section>
  );
}
