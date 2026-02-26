export interface ProcessSnapshot {
  pid: number;
  name: string;
  state: string;
  pc: number;
  blocked_until?: number;
}

export interface ProcessMetric {
  pid: number;
  name: string;
  response_time: number;
  turnaround: number;
  run_ticks: number;
  wait_ticks: number;
}

export interface SchedulingMetrics {
  policy: string;
  quantum?: number;
  total_ticks: number;
  completed_processes: number;
  avg_response_time: number;
  avg_turnaround_time: number;
  throughput_per_100_ticks: number;
  fairness_jain_index: number;
  processes: ProcessMetric[] | null;
  gantt: Array<{
    tick: number;
    pid: number;
  }> | null;
}

export interface MemorySnapshot {
  page_size: number;
  total_frames: number;
  frames: Array<{
    frame: number;
    pid?: number;
    vpn: number;
  }> | null;
  tlb: Array<{
    slot: number;
    pid: number;
    vpn: number;
    frame: number;
  }> | null;
  faults: {
    not_present: number;
    permission: number;
    tlb_hit: number;
    tlb_miss: number;
  };
}

export interface SnapshotDTO {
  protocol_version: string;
  session_id: string;
  tick: number;
  trace_hash: string;
  trace_length: number;
  processes: ProcessSnapshot[] | null;
  metrics: SchedulingMetrics;
  memory: MemorySnapshot;
  last_command?: string;
  challenge?: {
    max_steps?: number;
    max_policy_changes?: number;
    used_steps?: number;
    used_policy_changes?: number;
    remaining_steps?: number;
    remaining_policy_changes?: number;
  };
}

export interface Command {
  name: "spawn" | "step" | "run" | "pause" | "policy" | "reset";
  count?: number;
  process?: string;
  program?: string;
  policy?: "fifo" | "rr" | "mlfq";
  quantum?: number;
}

export interface CommandEnvelope {
  type: "command";
  command: Command;
}

export interface SessionEvent {
  type: "session.snapshot" | "session.error";
  sequence: number;
  session_id: string;
  snapshot?: SnapshotDTO;
  error?: string;
}
