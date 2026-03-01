package realtime

func applyV3ActionPolicy(session *Session, baseAllowed []string, actions []string, maxSteps, maxPolicy, maxConfig int) []string {
	allowed := mergeAllowedCommands(baseAllowed, mappedAllowedCommandsFromActions(actions))
	session.SetChallengePolicy(NewChallengeCommandPolicy(allowed, maxSteps, maxPolicy, maxConfig))
	return allowed
}

func mappedAllowedCommandsFromActions(actions []string) []string {
	out := make([]string, 0, len(actions))
	for _, action := range actions {
		note := actionCapabilityNote(action)
		if note.Status != "supported_now" || note.MappedCommand == "" {
			continue
		}
		out = appendIfMissing(out, note.MappedCommand)
	}
	return out
}

func mergeAllowedCommands(base, extras []string) []string {
	out := make([]string, 0, len(base)+len(extras))
	for _, item := range base {
		out = appendIfMissing(out, item)
	}
	for _, item := range extras {
		out = appendIfMissing(out, item)
	}
	return out
}

func appendIfMissing(items []string, value string) []string {
	if value == "" {
		return items
	}
	for _, item := range items {
		if item == value {
			return items
		}
	}
	return append(items, value)
}
