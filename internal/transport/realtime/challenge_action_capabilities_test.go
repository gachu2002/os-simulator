package realtime

import "testing"

func TestClassifyActionCapabilities(t *testing.T) {
	out := classifyActionCapabilities([]string{"execute_instruction", "migrate_job", "set_quantum", "exec", "wait"})
	if len(out.SupportedNow) != 4 {
		t.Fatalf("supported_now=%d want=4", len(out.SupportedNow))
	}
	if len(out.Planned) != 1 || out.Planned[0] != "migrate_job" {
		t.Fatalf("planned=%v want=[migrate_job]", out.Planned)
	}
}

func TestBuildActionCapabilityNotes(t *testing.T) {
	notes := buildActionCapabilityNotes([]string{"execute_instruction", "migrate_job", "unknown_action"})

	stepNote, ok := notes["execute_instruction"]
	if !ok {
		t.Fatalf("missing execute_instruction note")
	}
	if stepNote.Status != "supported_now" || stepNote.MappedCommand != "step" {
		t.Fatalf("execute_instruction note=%+v", stepNote)
	}

	migrateNote, ok := notes["migrate_job"]
	if !ok {
		t.Fatalf("missing migrate_job note")
	}
	if migrateNote.Status != "planned" || migrateNote.FallbackAction == "" {
		t.Fatalf("migrate_job note=%+v", migrateNote)
	}

	unknownNote, ok := notes["unknown_action"]
	if !ok {
		t.Fatalf("missing unknown_action note")
	}
	if unknownNote.Status != "planned" {
		t.Fatalf("unknown_action note=%+v", unknownNote)
	}
}
