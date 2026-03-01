package realtime

import "net/http"

func (s *Server) handleCurriculumV3(w http.ResponseWriter, r *http.Request) {
	if !requireMethod(w, r, http.MethodGet) {
		return
	}
	if s.cpuCurriculumErr != nil {
		respondError(w, r, http.StatusInternalServerError, "curriculum_unavailable", "curriculum unavailable")
		return
	}
	if err := validateCPUCurriculumScope(s.cpuCurriculumV3); err != nil {
		respondError(w, r, http.StatusInternalServerError, "curriculum_invalid", "curriculum is invalid")
		return
	}
	respondJSON(w, http.StatusOK, s.cpuCurriculumV3)
}
