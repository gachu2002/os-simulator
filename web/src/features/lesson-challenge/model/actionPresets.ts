export type SchedulerPolicy = "fifo" | "rr" | "mlfq";

export interface LessonActionOptions {
  count?: number;
  process?: string;
  program?: string;
  policy?: SchedulerPolicy;
  quantum?: number;
  frames?: number;
  tlbEntries?: number;
  diskLatency?: number;
  terminalLatency?: number;
}

interface ActionPreset {
  label: string;
  options?: LessonActionOptions;
}

export function toActionPreset(action: string): ActionPreset {
  switch (action) {
    case "create_process":
      return { label: "Create Process", options: { program: "COMPUTE 3; EXIT" } };
    case "block_process":
      return { label: "Block Process" };
    case "unblock_process":
      return { label: "Unblock Process" };
    case "kill_process":
      return { label: "Kill Process" };
    case "step":
    case "execute_instruction":
      return { label: "Step (1 tick)", options: { count: 1 } };
    case "run":
      return { label: "Run", options: { count: 8 } };
    case "run_to_completion":
      return { label: "Run To Completion", options: { count: 20 } };
    case "run_quanta":
      return { label: "Run 20 Quanta", options: { count: 20 } };
    case "fork":
      return { label: "fork()", options: { program: "COMPUTE 2; EXIT" } };
    case "exec":
      return { label: "exec(ls)", options: { process: "child", program: "ls" } };
    case "wait":
      return { label: "wait()", options: { count: 5 } };
    case "skip_wait":
      return { label: "Skip wait()" };
    case "issue_trap":
      return { label: "Issue Trap" };
    case "handle_syscall":
      return { label: "Handle Syscall" };
    case "return_from_trap":
      return { label: "Return From Trap" };
    case "fire_timer_interrupt":
      return { label: "Fire Timer Interrupt" };
    case "choose_next_process":
      return { label: "Choose Next Process" };
    case "set_policy_fifo_sjf_stcf":
      return { label: "Set Policy (FIFO)", options: { policy: "fifo" } };
    case "preempt_current_job":
      return { label: "Preempt Current Job" };
    case "set_quantum":
      return { label: "Set Quantum = 4", options: { quantum: 4 } };
    case "add_job":
    case "submit_job":
      return { label: "Submit Job", options: { program: "COMPUTE 8; EXIT" } };
    case "toggle_context_switch_cost":
      return { label: "Toggle Context-Switch Cost" };
    case "trigger_priority_boost":
      return { label: "Trigger Priority Boost" };
    case "toggle_gaming_prevention":
      return { label: "Toggle Gaming Prevention" };
    case "set_boost_interval":
      return { label: "Set Boost Interval" };
    case "set_tickets":
      return { label: "Set Tickets" };
    case "set_mode_lottery_or_stride":
      return { label: "Switch Lottery/Stride" };
    case "transfer_tickets":
      return { label: "Transfer Tickets" };
    case "migrate_job":
      return { label: "Migrate Job" };
    case "toggle_work_stealing":
      return { label: "Toggle Work Stealing" };
    case "toggle_affinity_protection":
      return { label: "Toggle Affinity Protection" };
    default:
      return { label: action };
  }
}
