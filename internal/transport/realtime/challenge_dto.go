package realtime

type ChallengeLimitsDTO struct {
	MaxSteps         int `json:"max_steps,omitempty"`
	MaxPolicyChanges int `json:"max_policy_changes,omitempty"`
	MaxConfigChanges int `json:"max_config_changes,omitempty"`
}

type ValidatorResultDTO struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Key      string `json:"key,omitempty"`
	Passed   bool   `json:"passed"`
	Message  string `json:"message,omitempty"`
	Expected string `json:"expected,omitempty"`
	Actual   string `json:"actual,omitempty"`
}
