package realtime

import "strings"

type ActionCapabilities struct {
	SupportedNow []string `json:"supported_now"`
	Planned      []string `json:"planned"`
}

type ActionCapabilityNote struct {
	Status         string `json:"status"`
	Reason         string `json:"reason,omitempty"`
	FallbackAction string `json:"fallback_action,omitempty"`
	MappedCommand  string `json:"mapped_command,omitempty"`
}

func classifyActionCapabilities(actions []string) ActionCapabilities {
	supported := make([]string, 0, len(actions))
	planned := make([]string, 0, len(actions))
	for _, action := range actions {
		if supportsV3Action(action) {
			supported = append(supported, action)
			continue
		}
		planned = append(planned, action)
	}
	return ActionCapabilities{SupportedNow: supported, Planned: planned}
}

func buildActionCapabilityNotes(actions []string) map[string]ActionCapabilityNote {
	out := make(map[string]ActionCapabilityNote, len(actions))
	for _, action := range actions {
		out[action] = actionCapabilityNote(action)
	}
	return out
}

func actionCapabilityNote(action string) ActionCapabilityNote {
	trimmed := strings.TrimSpace(action)
	spec := actionSpecFor(trimmed)
	return ActionCapabilityNote{
		Status:         string(spec.status),
		Reason:         spec.reason,
		FallbackAction: spec.fallbackAction,
		MappedCommand:  spec.mappedCommand,
	}
}

func supportsV3Action(action string) bool {
	spec := actionSpecFor(action)
	return spec.status == actionStatusSupported
}
