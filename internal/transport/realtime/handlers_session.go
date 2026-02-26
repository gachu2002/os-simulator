package realtime

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
)

func (s *Server) handleHealth(w http.ResponseWriter, _ *http.Request) {
	respondJSON(w, http.StatusOK, map[string]any{"status": "ok", "protocol_version": ProtocolVersion})
}

func (s *Server) handleWS(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(chi.URLParam(r, "id"))
	if id == "" {
		respondError(w, r, http.StatusBadRequest, "missing_session_id", "missing session id")
		return
	}

	session, ok := s.manager.Get(id)
	if !ok {
		respondError(w, r, http.StatusNotFound, "session_not_found", "session not found")
		return
	}

	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer func() { _ = conn.Close() }()

	if err := conn.WriteJSON(session.SnapshotEvent("connected")); err != nil {
		return
	}

	for {
		var req CommandEnvelope
		if err := conn.ReadJSON(&req); err != nil {
			return
		}
		if req.Type != "command" {
			if err := conn.WriteJSON(session.EmitError("unsupported message type")); err != nil {
				return
			}
			continue
		}
		if err := conn.WriteJSON(session.Apply(req.Command)); err != nil {
			return
		}
	}
}
