package realtime

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLessonsListEndpoint(t *testing.T) {
	ts := httptest.NewServer(NewServer(NewSessionManager()).Handler())
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/lessons")
	if err != nil {
		t.Fatalf("get lessons failed: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status=%d want=%d", resp.StatusCode, http.StatusOK)
	}

	var out LessonsResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		t.Fatalf("decode lessons failed: %v", err)
	}
	if len(out.Lessons) != 28 {
		t.Fatalf("lesson count=%d want=28", len(out.Lessons))
	}
	if len(out.Lessons[0].Stages) == 0 {
		t.Fatalf("first lesson should include stage summaries")
	}
	if len(out.Lessons[0].Stages) != 3 {
		t.Fatalf("first lesson stage count=%d want=3", len(out.Lessons[0].Stages))
	}
	if out.Lessons[0].Stages[0].Title == "" {
		t.Fatalf("first stage title should be populated")
	}
	if out.Lessons[0].SectionID == "" || out.Lessons[0].SectionTitle == "" {
		t.Fatalf("expected section metadata on lesson summary")
	}
	if out.Lessons[0].Stages[0].Goal == "" {
		t.Fatalf("expected goal metadata on stage summary")
	}
	if len(out.Lessons[0].Stages[0].ActionDescriptions) == 0 {
		t.Fatalf("expected action descriptions on stage summary")
	}
}
