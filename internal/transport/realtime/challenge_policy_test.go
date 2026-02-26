package realtime

import "testing"

func TestChallengePolicyEnforcesConfigChangeLimit(t *testing.T) {
	policy := NewChallengeCommandPolicy(
		[]string{"step", "set_frames", "set_tlb_entries", "set_disk_latency", "set_terminal_latency", "reset"},
		20,
		0,
		2,
	)

	if err := policy.Validate(Command{Name: "set_frames", Frames: 6}); err != nil {
		t.Fatalf("first config change should be allowed: %v", err)
	}
	if err := policy.Validate(Command{Name: "set_tlb_entries", TLBEntries: 6}); err != nil {
		t.Fatalf("second config change should be allowed: %v", err)
	}
	if err := policy.Validate(Command{Name: "set_disk_latency", DiskLatency: 2}); err == nil {
		t.Fatalf("expected config change limit failure")
	}
}
