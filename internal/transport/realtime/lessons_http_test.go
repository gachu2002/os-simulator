package realtime

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCurriculumEndpoint(t *testing.T) {
	ts := httptest.NewServer(NewServer(NewSessionManager()).Handler())
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/curriculum")
	if err != nil {
		t.Fatalf("get curriculum failed: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status=%d want=%d", resp.StatusCode, http.StatusOK)
	}

	var out CurriculumResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		t.Fatalf("decode curriculum failed: %v", err)
	}
	if len(out.Sections) < 5 {
		t.Fatalf("section count=%d want>=5", len(out.Sections))
	}
	if out.Sections[0].ID != "introduction" || !out.Sections[0].ComingSoon {
		t.Fatalf("expected introduction section as coming soon")
	}

	var foundVirtualization bool
	for _, section := range out.Sections {
		if section.ID == "virtualization" {
			foundVirtualization = true
			if len(section.Lessons) == 0 {
				t.Fatalf("expected virtualization lessons")
			}
		}
	}
	if !foundVirtualization {
		t.Fatalf("expected virtualization section")
	}
}

func TestLessonLearnEndpoint(t *testing.T) {
	ts := httptest.NewServer(NewServer(NewSessionManager()).Handler())
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/lessons/l01-sched-rr-basics/learn")
	if err != nil {
		t.Fatalf("get lesson learn failed: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status=%d want=%d", resp.StatusCode, http.StatusOK)
	}

	var out LessonLearnResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		t.Fatalf("decode lesson learn failed: %v", err)
	}
	if out.Lesson.ID != "l01-sched-rr-basics" {
		t.Fatalf("lesson id=%q want=%q", out.Lesson.ID, "l01-sched-rr-basics")
	}
	if len(out.Lesson.Stages) == 0 {
		t.Fatalf("expected learn stages")
	}
	if strings.TrimSpace(out.Lesson.Stages[0].CoreIdea) == "" {
		t.Fatalf("expected core idea content for first stage")
	}
}
