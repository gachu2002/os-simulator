package realtime

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
)

type LessonLearnV3Response struct {
	Version   string `json:"version"`
	SectionID string `json:"section_id"`
	Lesson    any    `json:"lesson"`
}

func (s *Server) handleLessonLearnV3(w http.ResponseWriter, r *http.Request) {
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
		respondJSON(w, http.StatusOK, LessonLearnV3Response{Version: s.cpuCurriculumV3.Version, SectionID: section.ID, Lesson: lesson})
		return
	}

	respondError(w, r, http.StatusNotFound, "lesson_not_found", "lesson not found")
}
