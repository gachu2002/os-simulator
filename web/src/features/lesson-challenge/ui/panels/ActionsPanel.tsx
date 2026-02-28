import type { Command } from "../../../../lib/types";
import { Button } from "../../../../components/ui/button";
import { Input } from "../../../../components/ui/input";
import { Label } from "../../../../components/ui/label";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "../../../../components/ui/select";

interface ActionsPanelProps {
  canSend: boolean;
  isLessonsLoading: boolean;
  isStartPending: boolean;
  isGradePending: boolean;
  hasAttempt: boolean;
  isStageUnlocked: boolean;
  policy: "fifo" | "rr" | "mlfq";
  quantum: number;
  frames: number;
  tlbEntries: number;
  diskLatency: number;
  terminalLatency: number;
  remainingSteps: number;
  remainingPolicyChanges: number;
  remainingConfigChanges: number;
  isAllowed: (name: Command["name"]) => boolean;
  onPolicyChange: (value: "fifo" | "rr" | "mlfq") => void;
  onQuantumChange: (value: number) => void;
  onFramesChange: (value: number) => void;
  onTLBEntriesChange: (value: number) => void;
  onDiskLatencyChange: (value: number) => void;
  onTerminalLatencyChange: (value: number) => void;
  onStart: () => void;
  onSubmit: () => void;
  onCommand: (command: Command) => void;
}

export function ActionsPanel(props: ActionsPanelProps) {
  const {
    canSend,
    isLessonsLoading,
    isStartPending,
    isGradePending,
    hasAttempt,
    isStageUnlocked,
    policy,
    quantum,
    frames,
    tlbEntries,
    diskLatency,
    terminalLatency,
    remainingSteps,
    remainingPolicyChanges,
    remainingConfigChanges,
    isAllowed,
    onPolicyChange,
    onQuantumChange,
    onFramesChange,
    onTLBEntriesChange,
    onDiskLatencyChange,
    onTerminalLatencyChange,
    onStart,
    onSubmit,
    onCommand,
  } = props;

  return (
    <section className="mt-3 rounded-lg border border-slate-200 bg-slate-50 p-3">
      <h3 className="text-sm font-semibold text-slate-900">1) Actions</h3>
      <div className="mt-2 flex flex-wrap items-end gap-2.5">
        <Button
          type="button"
          disabled={isStartPending || isLessonsLoading || !isStageUnlocked}
          onClick={onStart}
        >
          {isStartPending ? "Starting..." : "Start Challenge"}
        </Button>
        <Button
          type="button"
          variant="success"
          disabled={isGradePending || !hasAttempt}
          onClick={onSubmit}
        >
          {isGradePending ? "Submitting..." : "Submit"}
        </Button>
      </div>

      {hasAttempt ? (
        <p className="mt-2 text-sm text-slate-600">
          Remaining budget: steps {remainingSteps}, policy edits {remainingPolicyChanges}, config edits {remainingConfigChanges}
        </p>
      ) : null}

      {hasAttempt ? <ActionControls canSend={canSend} isAllowed={isAllowed} onCommand={onCommand} /> : null}

      {isAllowed("policy") ? (
        <div className="mt-2 flex flex-wrap items-end gap-2.5">
          <Label>
            Policy
            <Select
              value={policy}
              disabled={!canSend || !isAllowed("policy")}
              onValueChange={(value) => onPolicyChange(value as "fifo" | "rr" | "mlfq")}
            >
              <SelectTrigger className="w-32">
                <SelectValue placeholder="Policy" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="fifo">FIFO</SelectItem>
                <SelectItem value="rr">RR</SelectItem>
                <SelectItem value="mlfq">MLFQ</SelectItem>
              </SelectContent>
            </Select>
          </Label>
          <Label>
            Quantum
            <Input
              type="number"
              min={1}
              max={16}
              value={quantum}
              disabled={!canSend || !isAllowed("policy") || policy !== "rr"}
              onChange={(event) => onQuantumChange(Number(event.target.value))}
            />
          </Label>
          <Button
            type="button"
            disabled={!canSend || !isAllowed("policy")}
            onClick={() =>
              onCommand({ name: "policy", policy, quantum: policy === "rr" ? quantum : 0 })
            }
          >
            Apply Policy
          </Button>
        </div>
      ) : null}

      {isAllowed("set_frames") || isAllowed("set_tlb_entries") ? (
        <div className="mt-2 flex flex-wrap items-end gap-2.5">
          {isAllowed("set_frames") ? (
            <>
              <Label>
                Frames
                <Input
                  type="number"
                  min={1}
                  max={64}
                  value={frames}
                  onChange={(event) => onFramesChange(Number(event.target.value))}
                />
              </Label>
              <Button
                type="button"
                disabled={!canSend || !isAllowed("set_frames")}
                onClick={() => onCommand({ name: "set_frames", frames })}
              >
                Apply Frames
              </Button>
            </>
          ) : null}
          {isAllowed("set_tlb_entries") ? (
            <>
              <Label>
                TLB Entries
                <Input
                  type="number"
                  min={1}
                  max={64}
                  value={tlbEntries}
                  onChange={(event) => onTLBEntriesChange(Number(event.target.value))}
                />
              </Label>
              <Button
                type="button"
                disabled={!canSend || !isAllowed("set_tlb_entries")}
                onClick={() => onCommand({ name: "set_tlb_entries", tlb_entries: tlbEntries })}
              >
                Apply TLB
              </Button>
            </>
          ) : null}
        </div>
      ) : null}

      {isAllowed("set_disk_latency") || isAllowed("set_terminal_latency") ? (
        <div className="mt-2 flex flex-wrap items-end gap-2.5">
          {isAllowed("set_disk_latency") ? (
            <>
              <Label>
                Disk Latency
                <Input
                  type="number"
                  min={1}
                  max={64}
                  value={diskLatency}
                  onChange={(event) => onDiskLatencyChange(Number(event.target.value))}
                />
              </Label>
              <Button
                type="button"
                disabled={!canSend || !isAllowed("set_disk_latency")}
                onClick={() => onCommand({ name: "set_disk_latency", disk_latency: diskLatency })}
              >
                Apply Disk Latency
              </Button>
            </>
          ) : null}
          {isAllowed("set_terminal_latency") ? (
            <>
              <Label>
                Terminal Latency
                <Input
                  type="number"
                  min={1}
                  max={64}
                  value={terminalLatency}
                  onChange={(event) => onTerminalLatencyChange(Number(event.target.value))}
                />
              </Label>
              <Button
                type="button"
                disabled={!canSend || !isAllowed("set_terminal_latency")}
                onClick={() =>
                  onCommand({ name: "set_terminal_latency", terminal_latency: terminalLatency })
                }
              >
                Apply Terminal Latency
              </Button>
            </>
          ) : null}
        </div>
      ) : null}
    </section>
  );
}

function ActionControls({
  canSend,
  isAllowed,
  onCommand,
}: {
  canSend: boolean;
  isAllowed: (name: Command["name"]) => boolean;
  onCommand: (command: Command) => void;
}) {
  return (
    <div className="mt-2 flex flex-wrap items-end gap-2.5">
      <Button
        type="button"
        variant="secondary"
        disabled={!canSend || !isAllowed("run")}
        onClick={() => onCommand({ name: "run", count: 8 })}
      >
        Run 8
      </Button>
      <Button
        type="button"
        variant="secondary"
        disabled={!canSend || !isAllowed("step")}
        onClick={() => onCommand({ name: "step", count: 1 })}
      >
        Step
      </Button>
      <Button
        type="button"
        variant="outline"
        disabled={!canSend || !isAllowed("pause")}
        onClick={() => onCommand({ name: "pause" })}
      >
        Pause
      </Button>
      <Button
        type="button"
        variant="destructive"
        disabled={!canSend || !isAllowed("reset")}
        onClick={() => onCommand({ name: "reset" })}
      >
        Reset
      </Button>
    </div>
  );
}
