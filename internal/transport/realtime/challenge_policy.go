package realtime

import "fmt"

type ChallengeCommandPolicy struct {
	allowed          map[string]struct{}
	MaxSteps         int
	MaxPolicyChanges int
	usedSteps        int
	usedPolicyChange int
}

type ChallengeUsage struct {
	UsedSteps         int
	UsedPolicyChanges int
}

func NewChallengeCommandPolicy(allowedCommands []string, maxSteps, maxPolicyChanges int) ChallengeCommandPolicy {
	allowed := make(map[string]struct{}, len(allowedCommands))
	for _, name := range allowedCommands {
		allowed[name] = struct{}{}
	}
	return ChallengeCommandPolicy{
		allowed:          allowed,
		MaxSteps:         maxSteps,
		MaxPolicyChanges: maxPolicyChanges,
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
	}

	return nil
}

func (p ChallengeCommandPolicy) Usage() ChallengeUsage {
	return ChallengeUsage{UsedSteps: p.usedSteps, UsedPolicyChanges: p.usedPolicyChange}
}
