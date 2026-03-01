package realtime

import contentv3 "os-simulator-plan/internal/content/v3"

func (s *Server) lessonV3ByID(lessonID string) (contentv3.Section, contentv3.Lesson, bool) {
	if len(s.cpuCurriculumV3.Sections) == 0 {
		return contentv3.Section{}, contentv3.Lesson{}, false
	}
	section := s.cpuCurriculumV3.Sections[0]
	for _, lesson := range section.Lessons {
		if lesson.ID == lessonID {
			return section, lesson, true
		}
	}
	return contentv3.Section{}, contentv3.Lesson{}, false
}

func lessonPartByStageIndex(lesson contentv3.Lesson, stageIndex int) *contentv3.ChallengePart {
	if stageIndex < 0 || stageIndex >= len(lesson.Challenge.Parts) {
		return nil
	}
	return &lesson.Challenge.Parts[stageIndex]
}
