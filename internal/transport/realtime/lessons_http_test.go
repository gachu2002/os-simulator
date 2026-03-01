package realtime

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCurriculumV3Endpoint(t *testing.T) {
	ts := httptest.NewServer(NewServer(NewSessionManager()).Handler())
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/curriculum/v3")
	if err != nil {
		t.Fatalf("get curriculum v3 failed: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status=%d want=%d", resp.StatusCode, http.StatusOK)
	}

	var out struct {
		Version  string `json:"version"`
		Sections []struct {
			ID      string `json:"id"`
			Lessons []struct {
				ID string `json:"id"`
			} `json:"lessons"`
		} `json:"sections"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		t.Fatalf("decode curriculum v3 failed: %v", err)
	}
	if out.Version != "v3" {
		t.Fatalf("version=%q want=v3", out.Version)
	}
	if len(out.Sections) != 1 {
		t.Fatalf("sections=%d want=1", len(out.Sections))
	}
	if out.Sections[0].ID != sectionVirtualizationCPU {
		t.Fatalf("section id=%q want=%q", out.Sections[0].ID, sectionVirtualizationCPU)
	}
	if len(out.Sections[0].Lessons) != len(activeCPULessonOrder) {
		t.Fatalf("lessons=%d want=%d", len(out.Sections[0].Lessons), len(activeCPULessonOrder))
	}
}

func TestLessonLearnV3Endpoint(t *testing.T) {
	ts := httptest.NewServer(NewServer(NewSessionManager()).Handler())
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/lessons/l01-process-basics/learn/v3")
	if err != nil {
		t.Fatalf("get lesson learn v3 failed: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status=%d want=%d", resp.StatusCode, http.StatusOK)
	}

	var out struct {
		Version   string `json:"version"`
		SectionID string `json:"section_id"`
		Lesson    struct {
			ID        string `json:"id"`
			Title     string `json:"title"`
			Objective string `json:"objective"`
			Challenge struct {
				Actions []string `json:"actions"`
			} `json:"challenge"`
		} `json:"lesson"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		t.Fatalf("decode lesson learn v3 failed: %v", err)
	}
	if out.Version != "v3" {
		t.Fatalf("version=%q want=v3", out.Version)
	}
	if out.SectionID != sectionVirtualizationCPU {
		t.Fatalf("section id=%q want=%q", out.SectionID, sectionVirtualizationCPU)
	}
	if out.Lesson.ID != "l01-process-basics" {
		t.Fatalf("lesson id=%q want=%q", out.Lesson.ID, "l01-process-basics")
	}
	if out.Lesson.Title == "" || out.Lesson.Objective == "" || len(out.Lesson.Challenge.Actions) == 0 {
		t.Fatalf("expected lesson challenge metadata in v3 learn response")
	}
}

func TestChallengeManifestV3Endpoint(t *testing.T) {
	ts := httptest.NewServer(NewServer(NewSessionManager()).Handler())
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/lessons/l03-limited-direct-execution/challenge/v3")
	if err != nil {
		t.Fatalf("get challenge manifest v3 failed: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status=%d want=%d", resp.StatusCode, http.StatusOK)
	}

	var out struct {
		Version              string `json:"version"`
		SectionID            string `json:"section_id"`
		LessonID             string `json:"lesson_id"`
		ChallengeDescription string `json:"challenge_description"`
		Actions              []string
		PartRequired         bool `json:"part_required"`
		ActionCapabilities   struct {
			SupportedNow []string `json:"supported_now"`
			Planned      []string `json:"planned"`
		} `json:"action_capabilities"`
		ActionCapabilityNotes map[string]struct {
			Status         string `json:"status"`
			Reason         string `json:"reason"`
			FallbackAction string `json:"fallback_action"`
			MappedCommand  string `json:"mapped_command"`
		} `json:"action_capability_notes"`
		Visualizer           []string
		CrossCuttingFeatures []string `json:"cross_cutting_features"`
		Parts                []struct {
			ID string `json:"id"`
		}
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		t.Fatalf("decode challenge manifest v3 failed: %v", err)
	}
	if out.Version != "v3" {
		t.Fatalf("version=%q want=v3", out.Version)
	}
	if out.SectionID != sectionVirtualizationCPU {
		t.Fatalf("section id=%q want=%q", out.SectionID, sectionVirtualizationCPU)
	}
	if out.LessonID != "l03-limited-direct-execution" {
		t.Fatalf("lesson id=%q want=%q", out.LessonID, "l03-limited-direct-execution")
	}
	if out.ChallengeDescription == "" || len(out.Actions) == 0 || len(out.Visualizer) == 0 {
		t.Fatalf("expected challenge description/actions/visualizer")
	}
	if len(out.ActionCapabilities.SupportedNow) == 0 {
		t.Fatalf("expected supported_now capabilities")
	}
	if len(out.ActionCapabilities.Planned) != 0 {
		t.Fatalf("expected no planned capabilities, got=%v", out.ActionCapabilities.Planned)
	}
	note, ok := out.ActionCapabilityNotes["choose_next_process"]
	if !ok {
		t.Fatalf("expected note for choose_next_process")
	}
	if note.Status != "supported_now" || note.MappedCommand != "choose_next_process" {
		t.Fatalf("status=%q mapped=%q want supported_now/choose_next_process", note.Status, note.MappedCommand)
	}
	supportedNote, ok := out.ActionCapabilityNotes["execute_instruction"]
	if !ok {
		t.Fatalf("expected note for execute_instruction")
	}
	if supportedNote.Status != "supported_now" || supportedNote.MappedCommand != "step" {
		t.Fatalf("supported note=%+v want supported_now mapped step", supportedNote)
	}
	if len(out.CrossCuttingFeatures) == 0 {
		t.Fatalf("expected cross-cutting features")
	}
	if !out.PartRequired {
		t.Fatalf("expected part_required=true for lesson with A/B parts")
	}
	if len(out.Parts) != 2 {
		t.Fatalf("parts=%d want=2", len(out.Parts))
	}
}
