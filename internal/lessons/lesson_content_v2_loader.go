package lessons

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
)

//go:embed lesson_content_v2/*.json
var lessonContentV2FS embed.FS

type lessonContentRecord struct {
	ID               string             `json:"id"`
	Difficulty       string             `json:"difficulty"`
	EstimatedMinutes int                `json:"estimated_minutes"`
	ChapterRefs      []string           `json:"chapter_refs"`
	Learn            lessonLearnContent `json:"learn"`
}

type lessonLearnContent struct {
	CoreIdea              string   `json:"core_idea"`
	MechanismSteps        []string `json:"mechanism_steps"`
	WorkedExample         string   `json:"worked_example"`
	CommonMistakes        []string `json:"common_mistakes"`
	PreChallengeChecklist []string `json:"pre_challenge_checklist"`
}

func loadLessonContentV2() (map[string]lessonContentRecord, error) {
	entries, err := fs.ReadDir(lessonContentV2FS, "lesson_content_v2")
	if err != nil {
		return nil, fmt.Errorf("read lesson_content_v2 dir: %w", err)
	}

	out := make(map[string]lessonContentRecord, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		file := "lesson_content_v2/" + entry.Name()
		b, err := lessonContentV2FS.ReadFile(file)
		if err != nil {
			return nil, fmt.Errorf("read %s: %w", file, err)
		}

		var item lessonContentRecord
		if err := json.Unmarshal(b, &item); err != nil {
			return nil, fmt.Errorf("decode %s: %w", file, err)
		}
		if item.ID == "" {
			return nil, fmt.Errorf("%s has empty id", file)
		}
		if _, exists := out[item.ID]; exists {
			return nil, fmt.Errorf("duplicate lesson content id %q", item.ID)
		}
		if item.Learn.CoreIdea == "" {
			return nil, fmt.Errorf("lesson %q missing core idea", item.ID)
		}
		if item.Difficulty == "" {
			return nil, fmt.Errorf("lesson %q missing difficulty", item.ID)
		}
		if item.EstimatedMinutes <= 0 {
			return nil, fmt.Errorf("lesson %q missing estimated_minutes", item.ID)
		}
		if len(item.ChapterRefs) == 0 {
			return nil, fmt.Errorf("lesson %q missing chapter_refs", item.ID)
		}
		if len(item.Learn.MechanismSteps) == 0 {
			return nil, fmt.Errorf("lesson %q missing mechanism_steps", item.ID)
		}
		if item.Learn.WorkedExample == "" {
			return nil, fmt.Errorf("lesson %q missing worked_example", item.ID)
		}
		if len(item.Learn.CommonMistakes) == 0 {
			return nil, fmt.Errorf("lesson %q missing common_mistakes", item.ID)
		}
		if len(item.Learn.PreChallengeChecklist) == 0 {
			return nil, fmt.Errorf("lesson %q missing pre_challenge_checklist", item.ID)
		}
		out[item.ID] = item
	}

	if len(out) == 0 {
		return nil, fmt.Errorf("lesson content v2 is empty")
	}

	return out, nil
}
