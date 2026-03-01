export interface LessonBlueprintAction {
  command: string;
  label: string;
  purpose: string;
}

export interface LessonBlueprintPart {
  partID?: string;
  title: string;
  objective: string;
  description: string;
  theory: string[];
  successCriteria: string[];
  actions: LessonBlueprintAction[];
  visualChecks: string[];
  commonPitfalls: string[];
}

const LESSON_BLUEPRINTS: Record<string, LessonBlueprintPart[]> = {
  "l01-process-basics": [
    {
      title: "What is a Process?",
      objective: "Understand process state and state transitions.",
      description: "You control 3 processes and observe Ready/Running/Blocked transitions over time.",
      theory: [
        "A process is a running program abstracted by the OS.",
        "Process state includes code/heap/stack, registers, program counter, and open files.",
        "Running, Ready, and Blocked form a state machine driven by scheduling and I/O.",
        "The PCB stores per-process execution and scheduling metadata.",
      ],
      successCriteria: [
        "Create and run processes to show at least one transition into each lane.",
        "Explain why a process moved into Blocked or back into Ready.",
        "Submit only after transition timeline shows deterministic state flow.",
      ],
      actions: [
        { command: "create_process", label: "Create Process", purpose: "Introduce a runnable process." },
        { command: "block_process", label: "Block Process", purpose: "Simulate waiting on I/O." },
        { command: "unblock_process", label: "Unblock Process", purpose: "Return process to ready state." },
        { command: "kill_process", label: "Kill Process", purpose: "End a process lifecycle." },
        { command: "step", label: "Step", purpose: "Advance exactly one tick." },
      ],
      visualChecks: [
        "Ready/Running/Blocked lanes update after each action.",
        "Process cards move between lanes as transitions occur.",
        "Timeline grows by one transition per tick-level change.",
      ],
      commonPitfalls: [
        "Only watching one process and missing global queue behavior.",
        "Submitting before observing blocked-to-ready wakeup behavior.",
      ],
    },
  ],
  "l02-process-api-fork-exec-wait": [
    {
      title: "fork/exec/wait Lifecycle",
      objective: "Build a correct process family with parent/child synchronization.",
      description: "Act as a shell running commands in sequence and handle child reaping correctly.",
      theory: [
        "fork duplicates calling process state at the split point.",
        "exec replaces current process image with a new program.",
        "wait blocks parent until child completion and status reap.",
        "Skipping wait produces zombie-like outcomes in process lifecycle views.",
      ],
      successCriteria: [
        "Create at least one parent-child pair and run child workload.",
        "Use wait sequencing to avoid zombie state outcomes.",
        "Demonstrate effect difference between wait and skip_wait.",
      ],
      actions: [
        { command: "fork", label: "fork()", purpose: "Create child process from parent." },
        { command: "exec", label: "exec()", purpose: "Replace child image with selected program." },
        { command: "wait", label: "wait()", purpose: "Block parent until child finishes." },
        { command: "run_to_completion", label: "Run To Completion", purpose: "Execute full lifecycle quickly." },
        { command: "skip_wait", label: "Skip wait()", purpose: "Observe zombie-style lifecycle issue." },
      ],
      visualChecks: [
        "Family graph adds parent-child edge on fork.",
        "Node labels reflect PID/PPID/state/program changes.",
        "Zombie indicator appears when child exits unreaped.",
      ],
      commonPitfalls: [
        "Calling exec without first establishing child context.",
        "Ignoring parent wait semantics and misreading zombie outcome.",
      ],
    },
  ],
  "l03-limited-direct-execution": [
    {
      partID: "A",
      title: "System Call Round Trip",
      objective: "Complete a user-to-kernel-to-user trap cycle.",
      description: "Walk through syscall trap entry, handling, and return transitions.",
      theory: [
        "User mode executes unprivileged instructions; kernel mode handles privileged work.",
        "System calls trap into kernel using trap table entries.",
        "Correct trap handling preserves process execution context and mode bits.",
      ],
      successCriteria: [
        "Execute instruction in user mode before issuing trap.",
        "Handle syscall in kernel mode and return to user mode.",
        "Show trap table row highlight and mode-bit flip during flow.",
      ],
      actions: [
        { command: "execute_instruction", label: "Execute Instruction", purpose: "Advance user-space PC." },
        { command: "issue_trap", label: "Issue Trap", purpose: "Enter kernel via trap path." },
        { command: "handle_syscall", label: "Handle Syscall", purpose: "Run kernel handler logic." },
        { command: "return_from_trap", label: "Return From Trap", purpose: "Resume user execution." },
      ],
      visualChecks: [
        "Execution arrow drops from user to kernel zone on trap.",
        "Trap table entry highlight updates during handling.",
        "Mode bit toggles U -> K -> U across full round-trip.",
      ],
      commonPitfalls: [
        "Returning from trap before syscall handling step.",
        "Treating trap as regular compute instruction path.",
      ],
    },
    {
      partID: "B",
      title: "Timer Interrupt and Context Switch",
      objective: "Observe forced OS regain and process switch.",
      description: "Two processes compete while timer interrupts trigger scheduler intervention.",
      theory: [
        "Timer interrupts guarantee OS regains CPU control.",
        "Context switch saves current register set and restores next process state.",
        "Preemption timing drives fairness and response behavior.",
      ],
      successCriteria: [
        "Fire timer interrupt while process is running.",
        "Choose next process and confirm CPU ownership changes.",
        "Verify timeline records switch boundary at interrupt tick.",
      ],
      actions: [
        { command: "step", label: "Step", purpose: "Advance execution one tick." },
        { command: "fire_timer_interrupt", label: "Fire Timer Interrupt", purpose: "Force kernel re-entry." },
        { command: "choose_next_process", label: "Choose Next Process", purpose: "Select process after preemption." },
      ],
      visualChecks: [
        "Timer countdown/trigger aligns with switch event.",
        "Kernel stack save/restore panel reflects switch.",
        "CPU timeline shows process ownership handoff.",
      ],
      commonPitfalls: [
        "Stepping too far and missing interrupt boundary.",
        "Not validating switch in timeline after choosing next process.",
      ],
    },
  ],
  "l04-cpu-scheduling-basics": [
    {
      title: "FIFO vs SJF vs STCF",
      objective: "Minimize turnaround while understanding response/fairness tradeoffs.",
      description: "Run the same workload under FIFO, SJF, and STCF and compare metric outcomes.",
      theory: [
        "FIFO is easy to implement but can trigger convoy effects.",
        "SJF is optimal for turnaround when full job length info is known.",
        "STCF preempts and tracks shortest remaining time to adapt to arrivals.",
        "No single policy dominates all metrics for all workloads.",
      ],
      successCriteria: [
        "Run at least two different policies on the same workload shape.",
        "Explain one metric that improves and one that worsens after policy switch.",
        "Submit with evidence from timeline and metrics, not just final pass/fail.",
      ],
      actions: [
        { command: "set_arrival_time", label: "Set Arrival", purpose: "Change when jobs enter ready queue." },
        { command: "set_burst_length", label: "Set Burst Length", purpose: "Control job runtime demand." },
        { command: "set_policy_fifo_sjf_stcf", label: "Set Policy", purpose: "Switch between FIFO/SJF/STCF behavior." },
        { command: "step", label: "Step", purpose: "Advance by one tick to inspect queue evolution." },
        { command: "run_to_completion", label: "Run To Completion", purpose: "Finish workload quickly for metric comparison." },
        { command: "preempt_current_job", label: "Preempt Current Job", purpose: "Force preemption in preemptive scenarios." },
      ],
      visualChecks: [
        "Gantt chart reflects changed dispatch ordering after policy change.",
        "Ready queue ordering matches expected policy behavior.",
        "Turnaround and response metrics change consistently with observed trace.",
      ],
      commonPitfalls: [
        "Changing policy without resetting workload assumptions.",
        "Comparing metrics across different workloads instead of same workload.",
      ],
    },
  ],
  "l05-round-robin": [
    {
      title: "Round Robin Quantum Tuning",
      objective: "Find a quantum that balances response and throughput.",
      description: "Tune quantum across interactive jobs and observe the response-turnaround tension directly.",
      theory: [
        "Round Robin rotates jobs at fixed quantum boundaries.",
        "Small quantum reduces response latency but raises switch overhead.",
        "Large quantum reduces overhead but can hurt interactive responsiveness.",
        "Context switch cost matters in real scheduling performance.",
      ],
      successCriteria: [
        "Run at least two quantum values and compare trends.",
        "Explain tradeoff using both response and turnaround graphs.",
        "Use context switch cost toggle to justify chosen quantum.",
      ],
      actions: [
        { command: "set_quantum", label: "Set Quantum", purpose: "Choose time slice length." },
        { command: "add_job", label: "Add Job", purpose: "Increase contention in scheduler." },
        { command: "remove_job", label: "Remove Job", purpose: "Reduce contention and compare behavior." },
        { command: "step", label: "Step", purpose: "Inspect tick-level dispatch pattern." },
        { command: "run", label: "Run", purpose: "Advance multiple ticks under current tuning." },
        { command: "toggle_context_switch_cost", label: "Toggle Switch Cost", purpose: "Include/exclude overhead model." },
      ],
      visualChecks: [
        "Gantt chart shows shorter slices with lower quantum.",
        "Response trend and turnaround trend move in opposite directions as quantum changes.",
        "Context-switch markers increase when quantum is very small.",
      ],
      commonPitfalls: [
        "Optimizing only one metric and ignoring the other.",
        "Using too few ticks to infer long-run trend.",
      ],
    },
  ],
  "l06-mlfq": [
    {
      partID: "A",
      title: "Basic MLFQ Demotion",
      objective: "Understand queue demotion and boost behavior.",
      description: "Submit CPU-bound and I/O-bound jobs to observe how queue levels evolve.",
      theory: [
        "New jobs start at highest priority queue.",
        "Jobs that consume full quantum are demoted.",
        "Jobs that relinquish early can stay higher.",
        "Periodic boosts prevent starvation of lower queues.",
      ],
      successCriteria: [
        "Show at least one job demotion path from top to lower queue.",
        "Trigger boost and observe queue redistribution.",
        "Differentiate CPU-bound and I/O-bound queue behavior.",
      ],
      actions: [
        { command: "submit_job", label: "Submit Job", purpose: "Introduce CPU or I/O style workload." },
        { command: "choose_job_type_cpu_or_io_bound", label: "Choose Job Type", purpose: "Set expected queue behavior pattern." },
        { command: "step", label: "Step", purpose: "Observe precise queue movement." },
        { command: "trigger_priority_boost", label: "Trigger Boost", purpose: "Reset starvation pressure." },
      ],
      visualChecks: [
        "Queue rows show position changes over time.",
        "Demotion colors reflect queue depth.",
        "History records full job queue path.",
      ],
      commonPitfalls: [
        "Assuming every job must demote equally.",
        "Ignoring boost timing when interpreting fairness.",
      ],
    },
    {
      partID: "B",
      title: "Catch the Gamer",
      objective: "Identify and mitigate scheduler gaming behavior.",
      description: "Compare behavior with anti-gaming rule off and on, then adjust boost interval.",
      theory: [
        "A process can attempt to avoid demotion by yielding before quantum end.",
        "Anti-gaming tracks total CPU consumed per level, not just one slice.",
        "Boost interval tuning affects fairness vs responsiveness.",
      ],
      successCriteria: [
        "Demonstrate gamer persistence when prevention is off.",
        "Enable prevention and confirm expected demotion behavior.",
        "Use fairness score change to justify the fix.",
      ],
      actions: [
        { command: "step", label: "Step", purpose: "Advance and inspect exploit behavior." },
        { command: "toggle_gaming_prevention", label: "Toggle Prevention", purpose: "Switch anti-gaming rule mode." },
        { command: "set_boost_interval", label: "Set Boost Interval", purpose: "Tune starvation prevention cadence." },
      ],
      visualChecks: [
        "CPU hog meter grows on gaming process.",
        "Warning state appears when gamer remains high with prevention off.",
        "Fairness score improves once prevention is enabled.",
      ],
      commonPitfalls: [
        "Comparing one short window and concluding policy quality.",
        "Changing boost interval without tracking fairness trend.",
      ],
    },
  ],
  "l07-lottery-stride": [
    {
      title: "Proportional Share Allocation",
      objective: "Achieve and verify 50/30/20 CPU share targets.",
      description: "Configure tickets and compare Lottery vs Stride allocation behavior over multiple quanta.",
      theory: [
        "Lottery gives probabilistic proportional share via random draws.",
        "Stride gives deterministic proportional share via pass values.",
        "Ticket transfer changes effective CPU share weights.",
        "Long-run measurements should approach target allocation.",
      ],
      successCriteria: [
        "Configure tickets to represent 50/30/20 target split.",
        "Run enough quanta to compare actual vs target distribution.",
        "Switch mode and explain deterministic vs probabilistic behavior.",
      ],
      actions: [
        { command: "set_tickets", label: "Set Tickets", purpose: "Assign CPU share weights." },
        { command: "set_mode_lottery_or_stride", label: "Switch Mode", purpose: "Choose lottery or stride behavior." },
        { command: "run_quanta", label: "Run Quanta", purpose: "Collect distribution evidence over time." },
        { command: "transfer_tickets", label: "Transfer Tickets", purpose: "Reallocate shares mid-run." },
      ],
      visualChecks: [
        "Actual share converges toward target percentages.",
        "Lottery draw visualization changes each quantum.",
        "Stride pass ordering picks lowest pass process deterministically.",
      ],
      commonPitfalls: [
        "Judging lottery fairness from too few quanta.",
        "Changing tickets without resetting target expectations.",
      ],
    },
  ],
  "l08-multi-cpu-scheduling": [
    {
      title: "Multi-CPU Load Balancing and Affinity",
      objective: "Balance two CPUs without destroying cache affinity.",
      description: "Manage migration and work stealing under load imbalance while tracking cache warmth.",
      theory: [
        "Per-CPU queues improve locality but risk load imbalance.",
        "Migration raises utilization but can hurt cache warmth.",
        "Work stealing helps idle CPUs find runnable work.",
        "Affinity-aware scheduling trades balance against cache cost.",
      ],
      successCriteria: [
        "Resolve at least one visible imbalance event.",
        "Compare behavior with work stealing off vs on.",
        "Explain one migration decision using warmth and utilization evidence.",
      ],
      actions: [
        { command: "step", label: "Step", purpose: "Advance scheduler by one tick." },
        { command: "migrate_job", label: "Migrate Job", purpose: "Move task between CPU queues." },
        { command: "toggle_work_stealing", label: "Toggle Work Stealing", purpose: "Enable/disable auto-balance." },
        { command: "toggle_affinity_protection", label: "Toggle Affinity", purpose: "Protect or relax cache locality rule." },
      ],
      visualChecks: [
        "Dual CPU columns show queue imbalance clearly.",
        "Migration arrows align with balancing decisions.",
        "Per-CPU timeline indicates improved utilization after balancing.",
      ],
      commonPitfalls: [
        "Over-migrating and losing cache warmth benefits.",
        "Leaving one CPU idle while work remains queued elsewhere.",
      ],
    },
  ],
};

export function getLessonBlueprint(lessonID: string, partID?: string): LessonBlueprintPart | null {
  const parts = LESSON_BLUEPRINTS[lessonID] ?? [];
  if (parts.length === 0) {
    return null;
  }
  if (partID) {
    const matched = parts.find((item) => item.partID === partID);
    if (matched) {
      return matched;
    }
  }
  return parts[0];
}

export function getActionPurposeMap(lessonID: string, partID?: string): Record<string, string> {
  const blueprint = getLessonBlueprint(lessonID, partID);
  if (!blueprint) {
    return {};
  }
  return Object.fromEntries(blueprint.actions.map((item) => [item.command, item.purpose]));
}
