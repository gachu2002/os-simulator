import type { Command } from "../lib/types";

interface ControlBarProps {
  policy: "fifo" | "rr" | "mlfq";
  quantum: number;
  disabled: boolean;
  onPolicyChange: (policy: "fifo" | "rr" | "mlfq") => void;
  onQuantumChange: (quantum: number) => void;
  onCommand: (command: Command) => void;
}

export function ControlBar(props: ControlBarProps) {
  return (
    <section className="panel control-panel">
      <div className="control-row">
        <button
          type="button"
          disabled={props.disabled}
          onClick={() => props.onCommand({ name: "run", count: 8 })}
        >
          Run 8
        </button>
        <button
          type="button"
          disabled={props.disabled}
          onClick={() => props.onCommand({ name: "step", count: 1 })}
        >
          Step
        </button>
        <button
          type="button"
          disabled={props.disabled}
          onClick={() => props.onCommand({ name: "pause" })}
        >
          Pause
        </button>
        <button
          type="button"
          disabled={props.disabled}
          onClick={() => props.onCommand({ name: "reset" })}
        >
          Reset
        </button>
        <button
          type="button"
          disabled={props.disabled}
          onClick={() =>
            props.onCommand({
              name: "spawn",
              process: "demo",
              program: "COMPUTE 5; SYSCALL sleep 2; COMPUTE 3; EXIT",
            })
          }
        >
          Spawn Demo
        </button>
      </div>

      <div className="control-row">
        <label>
          Policy
          <select
            value={props.policy}
            disabled={props.disabled}
            onChange={(event) =>
              props.onPolicyChange(event.target.value as "fifo" | "rr" | "mlfq")
            }
          >
            <option value="fifo">FIFO</option>
            <option value="rr">RR</option>
            <option value="mlfq">MLFQ</option>
          </select>
        </label>
        <label>
          Quantum
          <input
            type="number"
            min={1}
            max={16}
            value={props.quantum}
            disabled={props.disabled || props.policy !== "rr"}
            onChange={(event) =>
              props.onQuantumChange(Number(event.target.value))
            }
          />
        </label>
        <button
          type="button"
          disabled={props.disabled}
          onClick={() =>
            props.onCommand({
              name: "policy",
              policy: props.policy,
              quantum: props.policy === "rr" ? props.quantum : 0,
            })
          }
        >
          Apply Policy
        </button>
      </div>
    </section>
  );
}
