package realtime

import (
	"errors"
	"fmt"
	"strings"

	"os-simulator-plan/internal/sim"
)

type actionStatus string

const (
	actionStatusSupported actionStatus = "supported_now"
	actionStatusPlanned   actionStatus = "planned"
)

type actionSpec struct {
	status         actionStatus
	mappedCommand  string
	reason         string
	fallbackAction string
	mapRequest     func(req ChallengeActionV3Request) (Command, error)
}

func actionSpecFor(action string) actionSpec {
	switch normalizeAction(action) {
	case "step", "execute_instruction", "issue_trap", "handle_syscall", "return_from_trap", "fire_timer_interrupt":
		return actionSpec{status: actionStatusSupported, mappedCommand: "step", mapRequest: mapStepAction}
	case "run", "run_quanta":
		return actionSpec{status: actionStatusSupported, mappedCommand: "run", mapRequest: mapRunAction}
	case "run_to_completion":
		return actionSpec{status: actionStatusSupported, mappedCommand: "run", mapRequest: mapRunToCompletionAction}
	case "create_process", "add_job", "submit_job", "fork":
		return actionSpec{status: actionStatusSupported, mappedCommand: "spawn", mapRequest: mapSpawnAction}
	case "exec":
		return actionSpec{
			status:        actionStatusSupported,
			mappedCommand: "spawn",
			reason:        "exec is approximated by spawning the selected program in current adapter",
			mapRequest:    mapExecAction,
		}
	case "wait":
		return actionSpec{
			status:        actionStatusSupported,
			mappedCommand: "run",
			reason:        "wait is approximated by advancing simulation until child completion window",
			mapRequest:    mapWaitAction,
		}
	case "set_quantum":
		return actionSpec{
			status:        actionStatusSupported,
			mappedCommand: "policy",
			reason:        "quantum updates are currently mapped to rr policy updates",
			mapRequest:    mapSetQuantumAction,
		}
	case "policy", "pause", "reset", "set_frames", "set_tlb_entries", "set_disk_latency", "set_terminal_latency":
		return actionSpec{status: actionStatusSupported, mappedCommand: normalizeAction(action), mapRequest: mapRuntimeConfigAction}
	case "set_policy_fifo_sjf_stcf":
		return actionSpec{status: actionStatusSupported, mappedCommand: "policy", mapRequest: mapPolicyAliasAction}
	case "block_process", "unblock_process", "kill_process", "preempt_current_job", "choose_next_process":
		return actionSpec{status: actionStatusSupported, mappedCommand: normalizeAction(action), mapRequest: mapProcessControlAction}
	case "skip_wait":
		return actionSpec{
			status:        actionStatusSupported,
			mappedCommand: "run",
			reason:        "skip_wait advances execution without issuing parent wait",
			mapRequest:    mapSkipWaitAction,
		}
	case "migrate_job", "toggle_work_stealing", "toggle_affinity_protection":
		return actionSpec{
			status:         actionStatusPlanned,
			reason:         "multi-cpu controls are not implemented in simulator core yet",
			fallbackAction: "step",
		}
	case "toggle_gaming_prevention", "set_boost_interval":
		return actionSpec{
			status:         actionStatusPlanned,
			reason:         "mlfq anti-gaming controls are not exposed in runtime adapter yet",
			fallbackAction: "set_policy_fifo_sjf_stcf",
		}
	case "transfer_tickets", "set_tickets", "set_mode_lottery_or_stride":
		return actionSpec{
			status:         actionStatusPlanned,
			reason:         "lottery/stride parameter controls are not implemented in action adapter yet",
			fallbackAction: "set_policy_fifo_sjf_stcf",
		}
	default:
		return actionSpec{status: actionStatusPlanned, reason: "action is not yet mapped in v3 adapter", fallbackAction: "step"}
	}
}

func normalizeAction(action string) string {
	return strings.ToLower(strings.TrimSpace(action))
}

func mapStepAction(req ChallengeActionV3Request) (Command, error) {
	count := req.Count
	if count <= 0 {
		count = 1
	}
	return Command{Name: "step", Count: count}, nil
}

func mapRunAction(req ChallengeActionV3Request) (Command, error) {
	if req.Count <= 0 {
		return Command{}, errors.New("run_quanta requires positive count")
	}
	return Command{Name: "run", Count: req.Count}, nil
}

func mapRunToCompletionAction(req ChallengeActionV3Request) (Command, error) {
	count := req.Count
	if count <= 0 {
		count = 25
	}
	return Command{Name: "run", Count: count}, nil
}

func mapSpawnAction(req ChallengeActionV3Request) (Command, error) {
	program := strings.TrimSpace(req.Program)
	if program == "" {
		program = "COMPUTE 2; EXIT"
	}
	return Command{Name: "spawn", Process: strings.TrimSpace(req.Process), Program: program}, nil
}

func mapExecAction(req ChallengeActionV3Request) (Command, error) {
	program := strings.TrimSpace(req.Program)
	if program == "" {
		return Command{}, errors.New("exec requires program")
	}
	return Command{Name: "spawn", Process: strings.TrimSpace(req.Process), Program: program}, nil
}

func mapWaitAction(req ChallengeActionV3Request) (Command, error) {
	count := req.Count
	if count <= 0 {
		count = 5
	}
	return Command{Name: "run", Count: count}, nil
}

func mapSkipWaitAction(req ChallengeActionV3Request) (Command, error) {
	count := req.Count
	if count <= 0 {
		count = 3
	}
	return Command{Name: "run", Count: count}, nil
}

func mapSetQuantumAction(req ChallengeActionV3Request) (Command, error) {
	if req.Quantum <= 0 {
		return Command{}, errors.New("set_quantum requires positive quantum")
	}
	return Command{Name: "policy", Policy: sim.PolicyRR, Quantum: req.Quantum}, nil
}

func mapRuntimeConfigAction(req ChallengeActionV3Request) (Command, error) {
	return Command{
		Name:            normalizeAction(req.Action),
		Policy:          strings.TrimSpace(req.Policy),
		Quantum:         req.Quantum,
		Frames:          req.Frames,
		TLBEntries:      req.TLBEntries,
		DiskLatency:     req.DiskLatency,
		TerminalLatency: req.TerminalLatency,
	}, nil
}

func mapPolicyAliasAction(req ChallengeActionV3Request) (Command, error) {
	policy := strings.ToLower(strings.TrimSpace(req.Policy))
	switch policy {
	case sim.PolicyFIFO, sim.PolicyRR, sim.PolicyMLFQ:
		return Command{Name: "policy", Policy: policy, Quantum: req.Quantum}, nil
	default:
		return Command{}, errors.New("policy must be one of fifo|rr|mlfq for current simulator")
	}
}

func mapProcessControlAction(req ChallengeActionV3Request) (Command, error) {
	return Command{Name: normalizeAction(req.Action), Process: strings.TrimSpace(req.Process)}, nil
}

func mapV3ActionToCommand(req ChallengeActionV3Request) (Command, error) {
	spec := actionSpecFor(req.Action)
	if spec.status == actionStatusPlanned {
		note := actionCapabilityNote(req.Action)
		return Command{}, fmt.Errorf("%s", unsupportedActionMessage(req.Action, note))
	}
	if spec.mapRequest == nil {
		return Command{}, errors.New("action is not supported by current simulator adapter")
	}
	return spec.mapRequest(req)
}
