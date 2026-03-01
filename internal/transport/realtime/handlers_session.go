package realtime

import (
	"net/http"
)

func (s *Server) handleHealth(w http.ResponseWriter, _ *http.Request) {
	respondJSON(w, http.StatusOK, map[string]any{"status": "ok", "protocol_version": ProtocolVersion})
}
