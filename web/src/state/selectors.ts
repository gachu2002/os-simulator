import type { SnapshotDTO } from "../lib/types";

export interface TimelineCell {
  tick: number;
  pid: number;
  label: string;
}

export function selectTimelineWindow(
  snapshot: SnapshotDTO | null,
  width = 40,
): TimelineCell[] {
  if (!snapshot) {
    return [];
  }
  const full = asArray(snapshot.metrics.gantt);
  const sliced = full.slice(Math.max(0, full.length - width));
  return sliced.map((slice) => ({
    tick: slice.tick,
    pid: slice.pid,
    label: slice.pid === 0 ? "idle" : `P${slice.pid}`,
  }));
}

export interface MemoryRow {
  frame: number;
  owner: string;
  vpn: string;
}

export function selectFrameRows(snapshot: SnapshotDTO | null): MemoryRow[] {
  if (!snapshot) {
    return [];
  }
  return asArray(snapshot.memory.frames)
    .slice()
    .sort((a, b) => a.frame - b.frame)
    .map((frame) => ({
      frame: frame.frame,
      owner: frame.pid ? `P${frame.pid}` : "free",
      vpn: `0x${frame.vpn.toString(16)}`,
    }));
}

export interface ProcessQueues {
  running: string[];
  ready: string[];
  blocked: string[];
  terminated: string[];
}

export function selectProcessQueues(snapshot: SnapshotDTO | null): ProcessQueues {
  const queues: ProcessQueues = {
    running: [],
    ready: [],
    blocked: [],
    terminated: [],
  };
  if (!snapshot) {
    return queues;
  }

  const ordered = asArray(snapshot.processes).slice().sort((a, b) => a.pid - b.pid);
  for (const proc of ordered) {
    const label = `P${proc.pid} ${proc.name}`;
    if (proc.state === "running") {
      queues.running.push(label);
      continue;
    }
    if (proc.state === "ready") {
      queues.ready.push(label);
      continue;
    }
    if (proc.state === "blocked") {
      const blockedSuffix = proc.blocked_until ? ` (until ${proc.blocked_until})` : "";
      queues.blocked.push(`${label}${blockedSuffix}`);
      continue;
    }
    if (proc.state === "terminated") {
      queues.terminated.push(label);
    }
  }
  return queues;
}

export interface ProcessMetricRow {
  pid: number;
  name: string;
  state: string;
  pc: number;
  responseTime: number;
  turnaround: number;
  runTicks: number;
  waitTicks: number;
}

export function selectProcessMetricRows(
  snapshot: SnapshotDTO | null,
): ProcessMetricRow[] {
  if (!snapshot) {
    return [];
  }

  const byPID = new Map(asArray(snapshot.processes).map((proc) => [proc.pid, proc]));
  return asArray(snapshot.metrics.processes)
    .slice()
    .sort((a, b) => a.pid - b.pid)
    .map((metric) => {
      const proc = byPID.get(metric.pid);
      return {
        pid: metric.pid,
        name: metric.name,
        state: proc?.state ?? "unknown",
        pc: proc?.pc ?? 0,
        responseTime: metric.response_time,
        turnaround: metric.turnaround,
        runTicks: metric.run_ticks,
        waitTicks: metric.wait_ticks,
      };
    });
}

function asArray<T>(value: T[] | null | undefined): T[] {
  return Array.isArray(value) ? value : [];
}
