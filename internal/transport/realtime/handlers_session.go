package realtime

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
)

func (s *Server) handleHealth(w http.ResponseWriter, _ *http.Request) {
	respondJSON(w, http.StatusOK, map[string]any{"status": "ok", "protocol_version": ProtocolVersion})
}

func (s *Server) handleSessions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, r, http.StatusMethodNotAllowed, "method_not_allowed", "method not allowed")
		return
	}
	defer func() { _ = r.Body.Close() }()
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	var cfg SessionConfig
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&cfg); err != nil && !errors.Is(err, io.EOF) {
		respondError(w, r, http.StatusBadRequest, "invalid_body", "invalid JSON body")
		return
	}

	session, err := s.manager.Create(cfg)
	if err != nil {
		respondError(w, r, http.StatusBadRequest, "invalid_session_config", err.Error())
		return
	}
	ev := session.SnapshotEvent("init")
	respondJSON(w, http.StatusCreated, CreateSessionResponse{SessionID: session.ID(), Snapshot: ev.Snapshot})
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
