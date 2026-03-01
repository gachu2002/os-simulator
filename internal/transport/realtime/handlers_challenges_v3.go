package realtime

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
)

type ChallengeManifestV3Response struct {
	Version               string                          `json:"version"`
	SectionID             string                          `json:"section_id"`
	LessonID              string                          `json:"lesson_id"`
	LessonTitle           string                          `json:"lesson_title"`
	LessonObjective       string                          `json:"lesson_objective"`
	ChallengeDescription  string                          `json:"challenge_description"`
	Actions               []string                        `json:"actions"`
	ActionCapabilities    ActionCapabilities              `json:"action_capabilities"`
	ActionCapabilityNotes map[string]ActionCapabilityNote `json:"action_capability_notes,omitempty"`
	PartRequired          bool                            `json:"part_required"`
	Visualizer            []string                        `json:"visualizer"`
	CrossCuttingFeatures  []string                        `json:"cross_cutting_features"`
	Parts                 []struct {
		ID          string `json:"id"`
		Title       string `json:"title"`
		Objective   string `json:"objective"`
		Description string `json:"description"`
	} `json:"parts,omitempty"`
}

func (s *Server) handleChallengeManifestV3(w http.ResponseWriter, r *http.Request) {
	if !requireMethod(w, r, http.MethodGet) {
		return
	}
	if s.cpuCurriculumErr != nil {
		respondError(w, r, http.StatusInternalServerError, "curriculum_unavailable", "curriculum unavailable")
		return
	}

	lessonID := strings.TrimSpace(chi.URLParam(r, "lessonID"))
	if lessonID == "" {
		respondError(w, r, http.StatusBadRequest, "invalid_lesson_id", "lesson id is required")
		return
	}
	if !isActiveCPULesson(lessonID) {
		respondError(w, r, http.StatusNotFound, "lesson_not_found", "lesson not found")
		return
	}

	section, lesson, ok := s.lessonV3ByID(lessonID)
	if ok {
		out := ChallengeManifestV3Response{
			Version:               s.cpuCurriculumV3.Version,
			SectionID:             section.ID,
			LessonID:              lesson.ID,
			LessonTitle:           lesson.Title,
			LessonObjective:       lesson.Objective,
			ChallengeDescription:  lesson.Challenge.Description,
			Actions:               append([]string(nil), lesson.Challenge.Actions...),
			ActionCapabilities:    classifyActionCapabilities(lesson.Challenge.Actions),
			ActionCapabilityNotes: buildActionCapabilityNotes(lesson.Challenge.Actions),
			PartRequired:          len(lesson.Challenge.Parts) > 0,
			Visualizer:            append([]string(nil), lesson.Challenge.Visualizer...),
			CrossCuttingFeatures:  append([]string(nil), s.cpuCurriculumV3.CrossCuttingFeature...),
		}
		if len(lesson.Challenge.Parts) > 0 {
			out.Parts = make([]struct {
				ID          string `json:"id"`
				Title       string `json:"title"`
				Objective   string `json:"objective"`
				Description string `json:"description"`
			}, 0, len(lesson.Challenge.Parts))
			for _, part := range lesson.Challenge.Parts {
				out.Parts = append(out.Parts, struct {
					ID          string `json:"id"`
					Title       string `json:"title"`
					Objective   string `json:"objective"`
					Description string `json:"description"`
				}{ID: part.ID, Title: part.Title, Objective: part.Objective, Description: part.Description})
			}
		}
		respondJSON(w, http.StatusOK, out)
		return
	}

	respondError(w, r, http.StatusNotFound, "lesson_not_found", "lesson not found")
}
