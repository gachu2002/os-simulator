package realtime

import "fmt"

type ChallengeCommandPolicy struct {
	allowed          map[string]struct{}
	MaxSteps         int
	MaxPolicyChanges int
	MaxConfigChanges int
	usedSteps        int
	usedPolicyChange int
	usedConfigChange int
}

type ChallengeUsage struct {
	UsedSteps         int
	UsedPolicyChanges int
	UsedConfigChanges int
}

func NewChallengeCommandPolicy(allowedCommands []string, maxSteps, maxPolicyChanges, maxConfigChanges int) ChallengeCommandPolicy {
	allowed := make(map[string]struct{}, len(allowedCommands))
	for _, name := range allowedCommands {
		allowed[name] = struct{}{}
	}
	return ChallengeCommandPolicy{
		allowed:          allowed,
		MaxSteps:         maxSteps,
		MaxPolicyChanges: maxPolicyChanges,
		MaxConfigChanges: maxConfigChanges,
	}
}

func (p ChallengeCommandPolicy) Clone() ChallengeCommandPolicy {
	allowed := make(map[string]struct{}, len(p.allowed))
	for name := range p.allowed {
		allowed[name] = struct{}{}
	}
	p.allowed = allowed
	return p
}

func (p *ChallengeCommandPolicy) Validate(cmd Command) error {
	if _, ok := p.allowed[cmd.Name]; !ok {
		return fmt.Errorf("command %q is not allowed in challenge", cmd.Name)
	}

	switch cmd.Name {
	case "step":
		count := cmd.Count
		if count == 0 {
			count = 1
		}
		if p.MaxSteps > 0 && p.usedSteps+count > p.MaxSteps {
			return fmt.Errorf("step limit exceeded: used=%d requested=%d max=%d", p.usedSteps, count, p.MaxSteps)
		}
		p.usedSteps += count
	case "run":
		if p.MaxSteps > 0 && p.usedSteps+cmd.Count > p.MaxSteps {
			return fmt.Errorf("step limit exceeded: used=%d requested=%d max=%d", p.usedSteps, cmd.Count, p.MaxSteps)
		}
		p.usedSteps += cmd.Count
	case "policy":
		if p.MaxPolicyChanges > 0 && p.usedPolicyChange+1 > p.MaxPolicyChanges {
			return fmt.Errorf("policy change limit exceeded: used=%d max=%d", p.usedPolicyChange, p.MaxPolicyChanges)
		}
		p.usedPolicyChange++
	case "set_frames", "set_tlb_entries", "set_disk_latency", "set_terminal_latency":
		if p.MaxConfigChanges > 0 && p.usedConfigChange+1 > p.MaxConfigChanges {
			return fmt.Errorf("config change limit exceeded: used=%d max=%d", p.usedConfigChange, p.MaxConfigChanges)
		}
		p.usedConfigChange++
	}

	return nil
}

func (p ChallengeCommandPolicy) Usage() ChallengeUsage {
	return ChallengeUsage{UsedSteps: p.usedSteps, UsedPolicyChanges: p.usedPolicyChange, UsedConfigChanges: p.usedConfigChange}
}
