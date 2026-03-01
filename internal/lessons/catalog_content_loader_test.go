package lessons

import (
	"strings"
	"testing"
)

func TestLoadDefaultCatalogContent(t *testing.T) {
	content, err := loadDefaultCatalogContent()
	if err != nil {
		t.Fatalf("load default catalog content: %v", err)
	}

	if got := len(content.CPU); got != 8 {
		t.Fatalf("cpu lesson specs=%d want=8", got)
	}
	if got := len(content.Memory.Lessons); got != 7 {
		t.Fatalf("memory lesson specs=%d want=7", got)
	}
	if got := len(content.Concurrency.Lessons); got != 6 {
		t.Fatalf("concurrency lesson specs=%d want=6", got)
	}
	if got := len(content.Persistence.Lessons); got != 7 {
		t.Fatalf("persistence lesson specs=%d want=7", got)
	}
}

func TestValidateCatalogContentRejectsMissingModulesAndDuplicates(t *testing.T) {
	base := catalogContent{
		CPU: []cpuLessonSpec{{id: "cpu-1", title: "CPU", seed: 1, p1: "COMPUTE 1; EXIT", p2: "COMPUTE 1; EXIT", steps: 2}},
		Memory: moduleContentMemory{Lessons: []memoryLessonSpec{{
			id:      "mem-1",
			title:   "Memory",
			seed:    2,
			frames:  2,
			program: "ACCESS 0x0 r; EXIT",
			steps:   2,
			faults:  1,
		}}},
		Concurrency: moduleContentIO{Lessons: []ioLessonSpec{{
			id:      "io-1",
			title:   "Concurrency",
			seed:    3,
			program: "BLOCK 1; EXIT",
			steps:   2,
		}}},
		Persistence: moduleContentIO{Lessons: []ioLessonSpec{{
			id:      "fs-1",
			title:   "Persistence",
			seed:    4,
			program: "SYSCALL open /docs/readme.txt; EXIT",
			steps:   2,
		}}},
	}

	if err := validateCatalogContent(base); err != nil {
		t.Fatalf("expected valid baseline catalog, got error: %v", err)
	}

	tests := []struct {
		name    string
		content catalogContent
		wantErr string
	}{
		{
			name: "missing cpu lessons",
			content: func() catalogContent {
				c := base
				c.CPU = nil
				return c
			}(),
			wantErr: "no cpu lessons",
		},
		{
			name: "missing memory lessons",
			content: func() catalogContent {
				c := base
				c.Memory.Lessons = nil
				return c
			}(),
			wantErr: "no memory lessons",
		},
		{
			name: "missing concurrency lessons",
			content: func() catalogContent {
				c := base
				c.Concurrency.Lessons = nil
				return c
			}(),
			wantErr: "no concurrency lessons",
		},
		{
			name: "missing persistence lessons",
			content: func() catalogContent {
				c := base
				c.Persistence.Lessons = nil
				return c
			}(),
			wantErr: "no persistence lessons",
		},
		{
			name: "duplicate ids across modules",
			content: func() catalogContent {
				c := base
				c.Persistence.Lessons[0].id = "cpu-1"
				return c
			}(),
			wantErr: "duplicate lesson id",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := validateCatalogContent(tc.content)
			if err == nil || !strings.Contains(err.Error(), tc.wantErr) {
				t.Fatalf("err=%v want contains %q", err, tc.wantErr)
			}
		})
	}
}

func TestValidateLessonSpecsRejectInvalidValues(t *testing.T) {
	if err := validateCPUSpec(cpuLessonSpec{}, map[string]struct{}{}); err == nil {
		t.Fatalf("expected invalid cpu lesson error")
	}
	if err := validateMemorySpec(memoryLessonSpec{}, map[string]struct{}{}); err == nil {
		t.Fatalf("expected invalid memory lesson error")
	}
	if err := validateIOSpec(ioLessonSpec{}, map[string]struct{}{}, "concurrency"); err == nil {
		t.Fatalf("expected invalid io lesson error")
	}
}
