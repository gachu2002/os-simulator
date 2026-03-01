package v3

import (
	_ "embed"
	"encoding/json"
	"fmt"
)

//go:embed cpu_curriculum.json
var cpuCurriculumJSON []byte

type Curriculum struct {
	Version             string    `json:"version"`
	Sections            []Section `json:"sections"`
	CrossCuttingFeature []string  `json:"cross_cutting_features"`
}

type Section struct {
	ID       string   `json:"id"`
	Title    string   `json:"title"`
	Subtitle string   `json:"subtitle"`
	Order    int      `json:"order"`
	Lessons  []Lesson `json:"lessons"`
}

type Lesson struct {
	ID            string    `json:"id"`
	Order         int       `json:"order"`
	Title         string    `json:"title"`
	Objective     string    `json:"objective"`
	Theory        Theory    `json:"theory"`
	Challenge     Challenge `json:"challenge"`
	Prerequisites []string  `json:"prerequisites,omitempty"`
}

type Theory struct {
	Concepts []string `json:"concepts"`
}

type Challenge struct {
	Description string          `json:"description"`
	Actions     []string        `json:"actions"`
	Visualizer  []string        `json:"visualizer"`
	Parts       []ChallengePart `json:"parts,omitempty"`
}

type ChallengePart struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Objective   string `json:"objective"`
	Description string `json:"description"`
}

func LoadCPUCurriculum() (Curriculum, error) {
	var c Curriculum
	if err := json.Unmarshal(cpuCurriculumJSON, &c); err != nil {
		return Curriculum{}, fmt.Errorf("decode v3 cpu curriculum: %w", err)
	}
	if err := validate(c); err != nil {
		return Curriculum{}, err
	}
	return c, nil
}

func validate(c Curriculum) error {
	if c.Version == "" {
		return fmt.Errorf("v3 curriculum missing version")
	}
	if len(c.Sections) != 1 {
		return fmt.Errorf("v3 curriculum expects exactly one section, got %d", len(c.Sections))
	}
	section := c.Sections[0]
	if section.ID == "" || section.Title == "" || section.Order <= 0 {
		return fmt.Errorf("v3 section metadata is incomplete")
	}
	if len(section.Lessons) == 0 {
		return fmt.Errorf("v3 section has no lessons")
	}

	seenIDs := make(map[string]struct{}, len(section.Lessons))
	seenOrders := make(map[int]struct{}, len(section.Lessons))
	for _, lesson := range section.Lessons {
		if lesson.ID == "" || lesson.Title == "" || lesson.Objective == "" || lesson.Order <= 0 {
			return fmt.Errorf("v3 lesson metadata is incomplete for %q", lesson.ID)
		}
		if _, ok := seenIDs[lesson.ID]; ok {
			return fmt.Errorf("duplicate v3 lesson id %q", lesson.ID)
		}
		seenIDs[lesson.ID] = struct{}{}
		if _, ok := seenOrders[lesson.Order]; ok {
			return fmt.Errorf("duplicate v3 lesson order %d", lesson.Order)
		}
		seenOrders[lesson.Order] = struct{}{}

		if len(lesson.Theory.Concepts) == 0 {
			return fmt.Errorf("v3 lesson %q missing theory concepts", lesson.ID)
		}
		if lesson.Challenge.Description == "" || len(lesson.Challenge.Actions) == 0 || len(lesson.Challenge.Visualizer) == 0 {
			return fmt.Errorf("v3 lesson %q missing challenge details", lesson.ID)
		}
		for _, part := range lesson.Challenge.Parts {
			if part.ID == "" || part.Title == "" || part.Objective == "" || part.Description == "" {
				return fmt.Errorf("v3 lesson %q has incomplete challenge part", lesson.ID)
			}
		}
	}
	if len(c.CrossCuttingFeature) == 0 {
		return fmt.Errorf("v3 curriculum missing cross-cutting features")
	}
	return nil
}
