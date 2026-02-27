package lessons

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
)

//go:embed lesson_stage_content_v2/*.json
var lessonStageContentV2FS embed.FS

type lessonStageContentRecord struct {
	LessonID string                    `json:"lesson_id"`
	Stages   []lessonStageContentStage `json:"stages"`
}

type lessonStageContentStage struct {
	ID             string             `json:"id"`
	Objective      string             `json:"objective"`
	Goal           string             `json:"goal"`
	TheoryDetail   string             `json:"theory_detail"`
	Hints          HintSet            `json:"hints"`
	ValidatorHints []validatorHintRaw `json:"validator_hints"`
}

type validatorHintRaw struct {
	Validator string  `json:"validator"`
	Hints     HintSet `json:"hints"`
}

func loadLessonStageContentV2() (map[string]map[string]lessonStageContentStage, error) {
	entries, err := fs.ReadDir(lessonStageContentV2FS, "lesson_stage_content_v2")
	if err != nil {
		return nil, fmt.Errorf("read lesson_stage_content_v2 dir: %w", err)
	}
	out := make(map[string]map[string]lessonStageContentStage, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		file := "lesson_stage_content_v2/" + entry.Name()
		b, err := lessonStageContentV2FS.ReadFile(file)
		if err != nil {
			return nil, fmt.Errorf("read %s: %w", file, err)
		}
		var record lessonStageContentRecord
		if err := json.Unmarshal(b, &record); err != nil {
			return nil, fmt.Errorf("decode %s: %w", file, err)
		}
		if record.LessonID == "" {
			return nil, fmt.Errorf("%s has empty lesson_id", file)
		}
		if _, exists := out[record.LessonID]; exists {
			return nil, fmt.Errorf("duplicate lesson stage content for %s", record.LessonID)
		}
		stageMap := make(map[string]lessonStageContentStage, len(record.Stages))
		for _, stage := range record.Stages {
			if stage.ID == "" {
				return nil, fmt.Errorf("%s has stage with empty id", file)
			}
			if _, ok := stageMap[stage.ID]; ok {
				return nil, fmt.Errorf("%s duplicate stage id %s", file, stage.ID)
			}
			stageMap[stage.ID] = stage
		}
		out[record.LessonID] = stageMap
	}
	return out, nil
}

func toValidatorHints(raw []validatorHintRaw) []ValidatorHint {
	if len(raw) == 0 {
		return nil
	}
	out := make([]ValidatorHint, 0, len(raw))
	for _, item := range raw {
		if item.Validator == "" {
			continue
		}
		out = append(out, ValidatorHint{Validator: item.Validator, Hints: item.Hints})
	}
	return out
}
