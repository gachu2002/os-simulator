package lessons

import (
	_ "embed"
	"encoding/json"
	"fmt"
)

//go:embed catalog_content_v1.json
var defaultCatalogContentJSON []byte

type catalogContent struct {
	CPU         []cpuLessonSpec
	Memory      moduleContentMemory
	Concurrency moduleContentIO
	Persistence moduleContentIO
}

type moduleContentMemory struct {
	ModulePrerequisite string
	Lessons            []memoryLessonSpec
}

type moduleContentIO struct {
	ModulePrerequisite string
	Lessons            []ioLessonSpec
}

type catalogContentRaw struct {
	CPU         []cpuLessonRecord      `json:"cpu"`
	Memory      moduleContentMemoryRaw `json:"memory"`
	Concurrency moduleContentIORaw     `json:"concurrency"`
	Persistence moduleContentIORaw     `json:"persistence"`
}

type cpuLessonRecord struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Seed  uint64 `json:"seed"`
	P1    string `json:"p1"`
	P2    string `json:"p2"`
	Steps int    `json:"steps"`
}

type memoryLessonRecord struct {
	ID      string  `json:"id"`
	Title   string  `json:"title"`
	Seed    uint64  `json:"seed"`
	Frames  int     `json:"frames"`
	Program string  `json:"program"`
	Steps   int     `json:"steps"`
	Faults  float64 `json:"faults"`
}

type ioLessonRecord struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	Seed    uint64 `json:"seed"`
	Program string `json:"program"`
	Steps   int    `json:"steps"`
}

type moduleContentMemoryRaw struct {
	ModulePrerequisite string               `json:"module_prerequisite"`
	Lessons            []memoryLessonRecord `json:"lessons"`
}

type moduleContentIORaw struct {
	ModulePrerequisite string           `json:"module_prerequisite"`
	Lessons            []ioLessonRecord `json:"lessons"`
}

func loadDefaultCatalogContent() (catalogContent, error) {
	var raw catalogContentRaw
	if err := json.Unmarshal(defaultCatalogContentJSON, &raw); err != nil {
		return catalogContent{}, fmt.Errorf("decode catalog content: %w", err)
	}

	content := catalogContent{
		CPU:         mapCPULessons(raw.CPU),
		Memory:      mapMemoryModule(raw.Memory),
		Concurrency: mapIOModule(raw.Concurrency),
		Persistence: mapIOModule(raw.Persistence),
	}
	if err := validateCatalogContent(content); err != nil {
		return catalogContent{}, err
	}
	return content, nil
}

func mapCPULessons(raw []cpuLessonRecord) []cpuLessonSpec {
	out := make([]cpuLessonSpec, 0, len(raw))
	for _, item := range raw {
		out = append(out, cpuLessonSpec{
			id:    item.ID,
			title: item.Title,
			seed:  item.Seed,
			p1:    item.P1,
			p2:    item.P2,
			steps: item.Steps,
		})
	}
	return out
}

func mapMemoryModule(raw moduleContentMemoryRaw) moduleContentMemory {
	lessons := make([]memoryLessonSpec, 0, len(raw.Lessons))
	for _, item := range raw.Lessons {
		lessons = append(lessons, memoryLessonSpec{
			id:      item.ID,
			title:   item.Title,
			seed:    item.Seed,
			frames:  item.Frames,
			program: item.Program,
			steps:   item.Steps,
			faults:  item.Faults,
		})
	}
	return moduleContentMemory{ModulePrerequisite: raw.ModulePrerequisite, Lessons: lessons}
}

func mapIOModule(raw moduleContentIORaw) moduleContentIO {
	lessons := make([]ioLessonSpec, 0, len(raw.Lessons))
	for _, item := range raw.Lessons {
		lessons = append(lessons, ioLessonSpec{
			id:      item.ID,
			title:   item.Title,
			seed:    item.Seed,
			program: item.Program,
			steps:   item.Steps,
		})
	}
	return moduleContentIO{ModulePrerequisite: raw.ModulePrerequisite, Lessons: lessons}
}

func validateCatalogContent(content catalogContent) error {
	if len(content.CPU) == 0 {
		return fmt.Errorf("catalog content has no cpu lessons")
	}
	if len(content.Memory.Lessons) == 0 {
		return fmt.Errorf("catalog content has no memory lessons")
	}
	if len(content.Concurrency.Lessons) == 0 {
		return fmt.Errorf("catalog content has no concurrency lessons")
	}
	if len(content.Persistence.Lessons) == 0 {
		return fmt.Errorf("catalog content has no persistence lessons")
	}

	seen := make(map[string]struct{}, 64)
	for _, spec := range content.CPU {
		if err := validateCPUSpec(spec, seen); err != nil {
			return err
		}
	}
	for _, spec := range content.Memory.Lessons {
		if err := validateMemorySpec(spec, seen); err != nil {
			return err
		}
	}
	for _, spec := range content.Concurrency.Lessons {
		if err := validateIOSpec(spec, seen, "concurrency"); err != nil {
			return err
		}
	}
	for _, spec := range content.Persistence.Lessons {
		if err := validateIOSpec(spec, seen, "persistence"); err != nil {
			return err
		}
	}
	return nil
}

func validateCPUSpec(spec cpuLessonSpec, seen map[string]struct{}) error {
	if spec.id == "" || spec.title == "" || spec.p1 == "" || spec.p2 == "" || spec.steps <= 0 {
		return fmt.Errorf("invalid cpu lesson spec %q", spec.id)
	}
	if _, ok := seen[spec.id]; ok {
		return fmt.Errorf("duplicate lesson id %q", spec.id)
	}
	seen[spec.id] = struct{}{}
	return nil
}

func validateMemorySpec(spec memoryLessonSpec, seen map[string]struct{}) error {
	if spec.id == "" || spec.title == "" || spec.program == "" || spec.frames <= 0 || spec.steps <= 0 {
		return fmt.Errorf("invalid memory lesson spec %q", spec.id)
	}
	if _, ok := seen[spec.id]; ok {
		return fmt.Errorf("duplicate lesson id %q", spec.id)
	}
	seen[spec.id] = struct{}{}
	return nil
}

func validateIOSpec(spec ioLessonSpec, seen map[string]struct{}, module string) error {
	if spec.id == "" || spec.title == "" || spec.program == "" || spec.steps <= 0 {
		return fmt.Errorf("invalid %s lesson spec %q", module, spec.id)
	}
	if _, ok := seen[spec.id]; ok {
		return fmt.Errorf("duplicate lesson id %q", spec.id)
	}
	seen[spec.id] = struct{}{}
	return nil
}
