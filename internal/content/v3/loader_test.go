package v3

import "testing"

func TestLoadCPUCurriculum(t *testing.T) {
	curriculum, err := LoadCPUCurriculum()
	if err != nil {
		t.Fatalf("load cpu curriculum: %v", err)
	}
	if curriculum.Version != "v3" {
		t.Fatalf("version=%q want=v3", curriculum.Version)
	}
	if len(curriculum.Sections) != 1 {
		t.Fatalf("sections=%d want=1", len(curriculum.Sections))
	}
	section := curriculum.Sections[0]
	if section.ID != "virtualization-cpu" {
		t.Fatalf("section id=%q want=virtualization-cpu", section.ID)
	}
	if len(section.Lessons) != 8 {
		t.Fatalf("lessons=%d want=8", len(section.Lessons))
	}
	if section.Lessons[0].Title == "" || len(section.Lessons[0].Challenge.Actions) == 0 {
		t.Fatalf("expected lesson details in first lesson")
	}
}

func TestValidateRejectsInvalidCurriculum(t *testing.T) {
	err := validate(Curriculum{})
	if err == nil {
		t.Fatalf("expected validate error for empty curriculum")
	}
}
