import type { SnapshotDTO } from "../lib/types";
import { selectFrameRows } from "../state/selectors";

interface MemoryPanelProps {
  snapshot: SnapshotDTO | null;
}

export function MemoryPanel({ snapshot }: MemoryPanelProps) {
  const rows = selectFrameRows(snapshot);
  const faults = snapshot?.memory.faults;
  const tlbSlots = snapshot?.memory.tlb?.length ?? 0;

  return (
    <section className="panel memory-panel">
      <h2>Memory View</h2>
      <div className="memory-stats">
        <span>Frames: {snapshot?.memory.total_frames ?? 0}</span>
        <span>TLB slots: {tlbSlots}</span>
        <span>Faults NP: {faults?.not_present ?? 0}</span>
        <span>Faults Perm: {faults?.permission ?? 0}</span>
        <span>
          TLB hit/miss: {faults?.tlb_hit ?? 0}/{faults?.tlb_miss ?? 0}
        </span>
      </div>

      {rows.length === 0 ? (
        <p className="empty">No frame allocations yet.</p>
      ) : (
        <div className="table-wrap">
          <table>
            <thead>
              <tr>
                <th>Frame</th>
                <th>Owner</th>
                <th>VPN</th>
              </tr>
            </thead>
            <tbody>
              {rows.map((row) => (
                <tr key={row.frame}>
                  <td>{row.frame}</td>
                  <td>{row.owner}</td>
                  <td>{row.vpn}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </section>
  );
}
