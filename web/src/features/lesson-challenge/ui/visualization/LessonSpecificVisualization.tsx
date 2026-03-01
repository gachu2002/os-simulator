import { useMemo, useState } from "react";

import type { SnapshotDTO, ProcessSnapshot } from "../../../../lib/types";

interface LessonSpecificVisualizationProps {
  lessonID: string;
  snapshot: SnapshotDTO | null;
  lastLessonAction: string;
}

export function LessonSpecificVisualization({
  lessonID,
  snapshot,
  lastLessonAction,
}: LessonSpecificVisualizationProps) {
  if (lessonID === "l01-process-basics") {
    return <ProcessStateLanes snapshot={snapshot} />;
  }
  if (lessonID === "l02-process-api-fork-exec-wait") {
    return <ProcessAPIPanel snapshot={snapshot} lastLessonAction={lastLessonAction} />;
  }
  if (lessonID === "l03-limited-direct-execution") {
    return <LimitedDirectExecutionPanel snapshot={snapshot} lastLessonAction={lastLessonAction} />;
  }
  if (lessonID === "l04-cpu-scheduling-basics") {
    return <SchedulingBasicsPanel snapshot={snapshot} />;
  }
  if (lessonID === "l05-round-robin") {
    return <RoundRobinTuningPanel snapshot={snapshot} />;
  }
  if (lessonID === "l06-mlfq") {
    return <MLFQPanel snapshot={snapshot} />;
  }
  if (lessonID === "l07-lottery-stride") {
    return <LotteryStridePanel snapshot={snapshot} />;
  }
  if (lessonID === "l08-multi-cpu-scheduling") {
    return <MultiCPUPanel snapshot={snapshot} lastLessonAction={lastLessonAction} />;
  }
  return null;
}

function ProcessStateLanes({ snapshot }: { snapshot: SnapshotDTO | null }) {
  const processes = useMemo(() => snapshot?.processes ?? [], [snapshot?.processes]);
  const running = processes.filter((p) => p.state === "running");
  const blocked = processes.filter((p) => p.state === "blocked");
  const ready = processes.filter((p) => p.state !== "running" && p.state !== "blocked");
  const [selectedPID, setSelectedPID] = useState<number | null>(null);
  const selectedProcess = useMemo(() => {
    if (processes.length === 0) {
      return null;
    }
    if (selectedPID !== null) {
      const matched = processes.find((item) => item.pid === selectedPID);
      if (matched) {
        return matched;
      }
    }
    return [...processes].sort((a, b) => a.pid - b.pid)[0];
  }, [processes, selectedPID]);
  const stateSummary = {
    ready: ready.length,
    running: running.length,
    blocked: blocked.length,
  };

  return (
    <section className="rounded-xl border border-cyan-200 bg-cyan-50 p-3">
      <h4 className="text-sm font-semibold text-slate-900">Process State Swim Lanes</h4>
      <p className="mt-1 text-xs text-slate-600">Ready Queue, Running, and Blocked lanes update each action.</p>
      <div className="mt-3 grid gap-2 lg:grid-cols-[1.6fr_1fr]">
        <div className="grid gap-2 md:grid-cols-3">
          <Lane title="Ready" items={ready} selectedPID={selectedProcess?.pid ?? null} onSelect={setSelectedPID} />
          <Lane title="Running" items={running} selectedPID={selectedProcess?.pid ?? null} onSelect={setSelectedPID} />
          <Lane title="Blocked" items={blocked} selectedPID={selectedProcess?.pid ?? null} onSelect={setSelectedPID} />
        </div>

        <article className="rounded-md border border-cyan-200 bg-white p-2">
          <p className="text-xs font-semibold text-slate-800">PCB Inspector</p>
          {selectedProcess ? (
            <div className="mt-2 grid gap-1 text-xs text-slate-700">
              <p>PID: {selectedProcess.pid}</p>
              <p>Name: {selectedProcess.name}</p>
              <p>State: {selectedProcess.state}</p>
              <p>PC: {selectedProcess.pc}</p>
              <p>Registers: R0={selectedProcess.pid * 3}, R1={selectedProcess.pc * 2}, SP=0x{(0x1000 + selectedProcess.pid * 32).toString(16)}</p>
              <p>Open files: /dev/stdin, /dev/stdout, /tmp/p{selectedProcess.pid}.log</p>
            </div>
          ) : (
            <p className="mt-2 text-xs text-slate-500">No process selected.</p>
          )}
        </article>
      </div>

      <div className="mt-3 rounded-md border border-cyan-200 bg-white p-2">
        <p className="text-xs font-semibold text-slate-800">Timeline</p>
        <div className="mt-1 flex flex-wrap gap-2 text-[11px] text-slate-700">
          <span className="rounded-full border border-slate-200 bg-slate-50 px-2 py-0.5">Ready={stateSummary.ready}</span>
          <span className="rounded-full border border-slate-200 bg-slate-50 px-2 py-0.5">Running={stateSummary.running}</span>
          <span className="rounded-full border border-slate-200 bg-slate-50 px-2 py-0.5">Blocked={stateSummary.blocked}</span>
        </div>
        <p className="mt-1 text-xs text-slate-600">tick={snapshot?.tick ?? 0}, trace_len={snapshot?.trace_length ?? 0}, last_command={snapshot?.last_command ?? "-"}</p>
      </div>
    </section>
  );
}

function ProcessAPIPanel({
  snapshot,
  lastLessonAction,
}: {
  snapshot: SnapshotDTO | null;
  lastLessonAction: string;
}) {
  const processes = snapshot?.processes ?? [];
  const sorted = [...processes].sort((a, b) => a.pid - b.pid);
  const parentPID = sorted[0]?.pid ?? 0;
  const showZombie = lastLessonAction === "skip_wait";
  return (
    <section className="rounded-xl border border-indigo-200 bg-indigo-50 p-3">
      <h4 className="text-sm font-semibold text-slate-900">Process API View</h4>
      <p className="mt-1 text-xs text-slate-600">Track process roster while practicing fork/exec/wait ordering.</p>
      <div className="mt-2 grid gap-2 lg:grid-cols-[1.2fr_1fr]">
        <div className="rounded-md border border-indigo-200 bg-white p-2">
          <p className="text-xs font-semibold text-slate-800">Process family tree</p>
          <div className="mt-2 grid gap-1">
            {sorted.length === 0 ? <p className="text-xs text-slate-500">No processes yet.</p> : null}
            {sorted.map((p, index) => (
              <p key={p.pid} className="text-xs text-slate-700">
                {index === 0 ? "root" : `child of PID ${parentPID}`} - PID {p.pid} ({p.name}) state={p.state}
                {showZombie && p.state === "terminated" ? " | ZOMBIE" : ""}
              </p>
            ))}
          </div>
        </div>

        <div className="rounded-md border border-indigo-200 bg-white p-2">
          <p className="text-xs font-semibold text-slate-800">Pseudo-code</p>
          <pre className="mt-1 whitespace-pre-wrap text-[11px] text-slate-700">
{`pid = fork();
if (pid == 0) {
  exec(program);
  exit();
}
${showZombie ? "// skip wait(): zombie risk" : "wait(pid);"}`}
          </pre>
          <p className="mt-2 text-xs font-semibold text-slate-800">Console Output</p>
          <div className="mt-1 rounded border border-slate-200 bg-slate-50 p-1.5 text-[11px] text-slate-700">
            [{snapshot?.tick ?? 0}] shell$ run child command
            <br />[{snapshot?.tick ?? 0}] child: program output...
          </div>
        </div>
      </div>
      <div className="mt-2 overflow-x-auto">
        <table className="w-full min-w-[420px] text-left text-xs">
          <thead>
            <tr className="border-b border-indigo-200 text-slate-700">
              <th className="py-1 pr-2">PID</th>
              <th className="py-1 pr-2">Name</th>
              <th className="py-1 pr-2">State</th>
              <th className="py-1 pr-2">PC</th>
            </tr>
          </thead>
          <tbody>
            {processes.map((p) => (
              <tr key={p.pid} className="border-b border-indigo-100 text-slate-700">
                <td className="py-1 pr-2">{p.pid}</td>
                <td className="py-1 pr-2">{p.name}</td>
                <td className="py-1 pr-2">{p.state}</td>
                <td className="py-1 pr-2">{p.pc}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </section>
  );
}

function Lane({
  title,
  items,
  selectedPID,
  onSelect,
}: {
  title: string;
  items: ProcessSnapshot[];
  selectedPID: number | null;
  onSelect: (pid: number) => void;
}) {
  return (
    <article className="rounded-md border border-slate-200 bg-white p-2">
      <h5 className="text-xs font-semibold text-slate-800">{title}</h5>
      <div className="mt-2 grid gap-1">
        {items.length === 0 ? <p className="text-xs text-slate-500">(empty)</p> : null}
        {items.map((p) => (
          <button
            key={p.pid}
            type="button"
            onClick={() => onSelect(p.pid)}
            className={`rounded border px-2 py-1 text-left transition ${
              selectedPID === p.pid
                ? "border-cyan-300 bg-cyan-50"
                : "border-slate-200 bg-slate-50 hover:border-slate-300"
            }`}
          >
            <p className="text-xs font-medium text-slate-700">PID {p.pid} - {p.name}</p>
            <p className="text-[11px] text-slate-500">state={p.state}, pc={p.pc}</p>
          </button>
        ))}
      </div>
    </article>
  );
}

function LimitedDirectExecutionPanel({
  snapshot,
  lastLessonAction,
}: {
  snapshot: SnapshotDTO | null;
  lastLessonAction: string;
}) {
  const normalized = lastLessonAction.trim().toLowerCase();
  const kernelActions = new Set(["issue_trap", "handle_syscall", "fire_timer_interrupt"]);
  const mode = kernelActions.has(normalized) ? "KERNEL" : "USER";
  const trapRows = [
    { id: "0x01", name: "read" },
    { id: "0x02", name: "write" },
    { id: "0x03", name: "open" },
  ];
  const activeTrap = normalized === "issue_trap" || normalized === "handle_syscall" ? "0x01" : "";

  return (
    <section className="rounded-xl border border-violet-200 bg-violet-50 p-3">
      <h4 className="text-sm font-semibold text-slate-900">Limited Direct Execution View</h4>
      <p className="mt-1 text-xs text-slate-600">
        User mode and kernel mode are shown as separate execution zones.
      </p>

      <div className="mt-3 grid gap-2 lg:grid-cols-[1.2fr_1fr]">
        <article className="rounded-md border border-violet-200 bg-white p-2">
          <p className="text-xs font-semibold text-slate-800">Execution zone</p>
          <div className="mt-2 grid gap-2 md:grid-cols-2">
            <div className={`rounded border px-2 py-2 text-xs ${mode === "USER" ? "border-cyan-300 bg-cyan-50" : "border-slate-200 bg-slate-50"}`}>
              <p className="font-semibold text-slate-800">User Space</p>
              <p className="mt-1 text-slate-600">PC tick: {snapshot?.tick ?? 0}</p>
            </div>
            <div className={`rounded border px-2 py-2 text-xs ${mode === "KERNEL" ? "border-amber-300 bg-amber-50" : "border-slate-200 bg-slate-50"}`}>
              <p className="font-semibold text-slate-800">Kernel Space</p>
              <p className="mt-1 text-slate-600">mode bit: {mode === "KERNEL" ? "K" : "U"}</p>
            </div>
          </div>
          <p className="mt-2 text-xs text-slate-600">Last action: {lastLessonAction || "-"}</p>
        </article>

        <article className="rounded-md border border-violet-200 bg-white p-2">
          <p className="text-xs font-semibold text-slate-800">Trap table</p>
          <div className="mt-2 grid gap-1">
            {trapRows.map((row) => (
              <div
                key={row.id}
                className={`rounded border px-2 py-1 text-xs ${activeTrap === row.id ? "border-violet-400 bg-violet-100" : "border-slate-200 bg-slate-50"}`}
              >
                <span className="font-medium text-slate-700">{row.id}</span>
                <span className="ml-2 text-slate-600">{row.name}</span>
              </div>
            ))}
          </div>
          <p className="mt-2 text-xs text-slate-600">
            Timer/interrupt hint: trigger `fire_timer_interrupt` to force kernel re-entry.
          </p>
        </article>
      </div>
    </section>
  );
}

function SchedulingBasicsPanel({ snapshot }: { snapshot: SnapshotDTO | null }) {
  const metrics = snapshot?.metrics;
  const policy = (metrics?.policy ?? "-").toUpperCase();
  const avgTurnaround = metrics?.avg_turnaround_time ?? 0;
  const avgResponse = metrics?.avg_response_time ?? 0;
  const fairness = metrics?.fairness_jain_index ?? 0;

  return (
    <section className="rounded-xl border border-emerald-200 bg-emerald-50 p-3">
      <h4 className="text-sm font-semibold text-slate-900">Scheduling Comparison Focus</h4>
      <p className="mt-1 text-xs text-slate-600">Use policy toggles to compare turnaround, response, and fairness on one workload.</p>
      <div className="mt-3 grid gap-2 md:grid-cols-4">
        <MetricCard label="Policy" value={policy} />
        <MetricCard label="Avg Turnaround" value={formatNumber(avgTurnaround)} />
        <MetricCard label="Avg Response" value={formatNumber(avgResponse)} />
        <MetricCard label="Fairness" value={formatNumber(fairness)} />
      </div>
    </section>
  );
}

function RoundRobinTuningPanel({ snapshot }: { snapshot: SnapshotDTO | null }) {
  const metrics = snapshot?.metrics;
  const quantum = metrics?.quantum ?? 0;
  const avgResponse = metrics?.avg_response_time ?? 0;
  const avgTurnaround = metrics?.avg_turnaround_time ?? 0;
  const throughput = metrics?.throughput_per_100_ticks ?? 0;

  const responseScore = Math.max(0, 100 - avgResponse * 8);
  const turnaroundScore = Math.max(0, 100 - avgTurnaround * 6);

  return (
    <section className="rounded-xl border border-blue-200 bg-blue-50 p-3">
      <h4 className="text-sm font-semibold text-slate-900">Round Robin Tradeoff Lens</h4>
      <p className="mt-1 text-xs text-slate-600">Tune quantum, then compare response and turnaround trend bars.</p>
      <div className="mt-2 text-xs text-slate-700">Current quantum: {quantum || "-"}</div>
      <div className="mt-3 grid gap-2">
        <TrendBar label="Response trend" value={responseScore} meta={`avg=${formatNumber(avgResponse)}`} />
        <TrendBar label="Turnaround trend" value={turnaroundScore} meta={`avg=${formatNumber(avgTurnaround)}`} />
      </div>
      <p className="mt-3 text-xs text-slate-700">Throughput/100 ticks: {formatNumber(throughput)}</p>
    </section>
  );
}

function MLFQPanel({ snapshot }: { snapshot: SnapshotDTO | null }) {
  const processes = snapshot?.processes ?? [];
  const topQueue = processes.filter((p) => p.state === "ready").slice(0, 2);
  const midQueue = processes.filter((p) => p.state === "running").slice(0, 2);
  const lowQueue = processes.filter((p) => p.state === "blocked").slice(0, 2);
  const fairness = snapshot?.metrics.fairness_jain_index ?? 0;

  return (
    <section className="rounded-xl border border-amber-200 bg-amber-50 p-3">
      <h4 className="text-sm font-semibold text-slate-900">MLFQ Queue View</h4>
      <p className="mt-1 text-xs text-slate-600">Queue lanes approximate top/mid/low priority behavior while you tune controls.</p>
      <div className="mt-3 grid gap-2 md:grid-cols-3">
        <QueueLane title="Q0 (top)" items={topQueue} tone="blue" />
        <QueueLane title="Q1 (mid)" items={midQueue} tone="yellow" />
        <QueueLane title="Q2 (low)" items={lowQueue} tone="red" />
      </div>
      <p className="mt-3 text-xs text-slate-700">Fairness score: {formatNumber(fairness)}</p>
    </section>
  );
}

function LotteryStridePanel({ snapshot }: { snapshot: SnapshotDTO | null }) {
  const rows = snapshot?.metrics.processes ?? [];
  const totalRun = rows.reduce((sum, row) => sum + row.run_ticks, 0);

  return (
    <section className="rounded-xl border border-fuchsia-200 bg-fuchsia-50 p-3">
      <h4 className="text-sm font-semibold text-slate-900">Proportional Share View</h4>
      <p className="mt-1 text-xs text-slate-600">Use CPU share bars as a stand-in for lottery/stride distribution feedback.</p>
      <div className="mt-3 grid gap-2">
        {rows.length === 0 ? <p className="text-xs text-slate-500">No process share data yet.</p> : null}
        {rows.map((row) => {
          const share = totalRun > 0 ? (row.run_ticks / totalRun) * 100 : 0;
          return (
            <article key={row.pid}>
              <div className="mb-1 flex items-center justify-between text-xs text-slate-700">
                <span>{row.name} (PID {row.pid})</span>
                <span>{formatNumber(share)}%</span>
              </div>
              <div className="h-2 rounded bg-fuchsia-100">
                <div className="h-2 rounded bg-fuchsia-500" style={{ width: `${Math.max(0, Math.min(100, share))}%` }} />
              </div>
            </article>
          );
        })}
      </div>
    </section>
  );
}

function MultiCPUPanel({
  snapshot,
  lastLessonAction,
}: {
  snapshot: SnapshotDTO | null;
  lastLessonAction: string;
}) {
  const processes = snapshot?.processes ?? [];
  const leftCPU = processes.filter((p) => p.pid%2 === 0);
  const rightCPU = processes.filter((p) => p.pid%2 !== 0);
  const imbalance = Math.abs(leftCPU.length - rightCPU.length);

  return (
    <section className="rounded-xl border border-rose-200 bg-rose-50 p-3">
      <h4 className="text-sm font-semibold text-slate-900">Multi-CPU Scheduler View</h4>
      <p className="mt-1 text-xs text-slate-600">Dual CPU columns approximate queue split and migration pressure.</p>
      <div className="mt-3 grid gap-2 md:grid-cols-2">
        <CPUColumn title="CPU 0" items={leftCPU} />
        <CPUColumn title="CPU 1" items={rightCPU} />
      </div>
      <p className="mt-3 text-xs text-slate-700">Imbalance indicator: {imbalance} | last action: {lastLessonAction || "-"}</p>
    </section>
  );
}

function MetricCard({ label, value }: { label: string; value: string }) {
  return (
    <article className="rounded-md border border-slate-200 bg-white px-2 py-2">
      <p className="text-[11px] text-slate-500">{label}</p>
      <p className="mt-1 text-sm font-semibold text-slate-800">{value}</p>
    </article>
  );
}

function TrendBar({ label, value, meta }: { label: string; value: number; meta: string }) {
  const clamped = Math.max(0, Math.min(100, value));
  return (
    <article>
      <div className="mb-1 flex items-center justify-between text-xs text-slate-700">
        <span>{label}</span>
        <span>{meta}</span>
      </div>
      <div className="h-2 rounded bg-blue-100">
        <div className="h-2 rounded bg-blue-500" style={{ width: `${clamped}%` }} />
      </div>
    </article>
  );
}

function QueueLane({
  title,
  items,
  tone,
}: {
  title: string;
  items: ProcessSnapshot[];
  tone: "blue" | "yellow" | "red";
}) {
  const toneClasses =
    tone === "blue"
      ? "border-blue-200 bg-blue-50"
      : tone === "yellow"
        ? "border-amber-200 bg-amber-50"
        : "border-rose-200 bg-rose-50";
  return (
    <article className={`rounded-md border p-2 ${toneClasses}`}>
      <p className="text-xs font-semibold text-slate-800">{title}</p>
      <div className="mt-2 grid gap-1">
        {items.length === 0 ? <p className="text-xs text-slate-500">(empty)</p> : null}
        {items.map((item) => (
          <p key={item.pid} className="rounded border border-slate-200 bg-white px-2 py-1 text-xs text-slate-700">
            PID {item.pid} - {item.name}
          </p>
        ))}
      </div>
    </article>
  );
}

function CPUColumn({ title, items }: { title: string; items: ProcessSnapshot[] }) {
  return (
    <article className="rounded-md border border-rose-200 bg-white p-2">
      <p className="text-xs font-semibold text-slate-800">{title}</p>
      <div className="mt-2 grid gap-1">
        {items.length === 0 ? <p className="text-xs text-slate-500">(idle)</p> : null}
        {items.map((item) => (
          <p key={item.pid} className="rounded border border-slate-200 bg-slate-50 px-2 py-1 text-xs text-slate-700">
            PID {item.pid} - {item.name} [{item.state}]
          </p>
        ))}
      </div>
    </article>
  );
}

function formatNumber(value: number): string {
  if (!Number.isFinite(value)) {
    return "-";
  }
  return value.toFixed(2);
}
