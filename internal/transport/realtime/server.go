package realtime

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
	"sync/atomic"

	"os-simulator-plan/internal/lessons"
	"os-simulator-plan/internal/sim"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
)

type Server struct {
	manager      *SessionManager
	lessonMu     sync.Mutex
	lessonEngine *lessons.Engine
	upgrader     websocket.Upgrader
	origins      map[string]struct{}
}

func NewServer(manager *SessionManager) *Server {
	return NewServerWithLessons(manager, lessons.NewEngine())
}

func NewServerWithLessons(manager *SessionManager, lessonEngine *lessons.Engine) *Server {
	origins := allowedOriginsFromEnv(os.Getenv("CORS_ALLOW_ORIGIN"))
	return &Server{
		manager:      manager,
		lessonEngine: lessonEngine,
		origins:      origins,
		upgrader: websocket.Upgrader{
			ReadBufferSize:  4096,
			WriteBufferSize: 4096,
			CheckOrigin: func(r *http.Request) bool {
				return isOriginAllowed(origins, r.Header.Get("Origin"))
			},
		},
	}
}

func (s *Server) Handler() http.Handler {
	router := chi.NewRouter()
	router.HandleFunc("/healthz", s.handleHealth)
	router.HandleFunc("/sessions", s.handleSessions)
	router.HandleFunc("/lessons", s.handleLessons)
	router.HandleFunc("/lessons/run", s.handleLessonRun)
	router.HandleFunc("/ws/{id}", s.handleWS)
	return withRequestID(withCORS(originsForMiddleware(s.origins), router))
}

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

func respondJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

type apiError struct {
	Code      string `json:"code"`
	Message   string `json:"message"`
	RequestID string `json:"request_id,omitempty"`
}

func respondError(w http.ResponseWriter, r *http.Request, status int, code, message string) {
	reqID, _ := requestIDFromContext(r.Context())
	respondJSON(w, status, apiError{Code: code, Message: message, RequestID: reqID})
}

func withCORS(allowed map[string]struct{}, next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if isOriginAllowed(allowed, origin) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-Request-ID")
		w.Header().Set("Vary", "Origin")

		if r.Method == http.MethodOptions {
			if origin != "" && !isOriginAllowed(allowed, origin) {
				w.WriteHeader(http.StatusForbidden)
				return
			}
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

type ctxKey string

const requestIDKey ctxKey = "request_id"

func withRequestID(next http.Handler) http.Handler {
	var seq atomic.Uint64
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := strings.TrimSpace(r.Header.Get("X-Request-ID"))
		if reqID == "" {
			reqID = fmt.Sprintf("req-%08d", seq.Add(1))
		}
		w.Header().Set("X-Request-ID", reqID)
		ctx := context.WithValue(r.Context(), requestIDKey, reqID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func requestIDFromContext(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(requestIDKey).(string)
	return v, ok
}

func allowedOriginsFromEnv(raw string) map[string]struct{} {
	out := map[string]struct{}{}
	if strings.TrimSpace(raw) == "" {
		out["http://localhost:5173"] = struct{}{}
		out["http://127.0.0.1:5173"] = struct{}{}
		out["https://localhost:5173"] = struct{}{}
		out["https://127.0.0.1:5173"] = struct{}{}
		return out
	}
	for _, part := range strings.Split(raw, ",") {
		origin := strings.TrimSpace(part)
		if origin == "" {
			continue
		}
		out[origin] = struct{}{}
	}
	return out
}

func originsForMiddleware(in map[string]struct{}) map[string]struct{} {
	out := make(map[string]struct{}, len(in))
	for k := range in {
		out[k] = struct{}{}
	}
	return out
}

func isOriginAllowed(allowed map[string]struct{}, origin string) bool {
	if origin == "" {
		return true
	}
	_, ok := allowed[origin]
	return ok
}

type LessonStageSummary struct {
	Index int    `json:"index"`
	ID    string `json:"id"`
	Title string `json:"title"`
}

type LessonSummary struct {
	ID     string               `json:"id"`
	Title  string               `json:"title"`
	Module string               `json:"module"`
	Stages []LessonStageSummary `json:"stages"`
}

type LessonsResponse struct {
	Lessons []LessonSummary `json:"lessons"`
}

type LessonRunRequest struct {
	LessonID   string `json:"lesson_id"`
	StageIndex int    `json:"stage_index"`
}

type LessonOutputDTO struct {
	Tick         sim.Tick              `json:"tick"`
	TraceHash    string                `json:"trace_hash"`
	TraceLength  int                   `json:"trace_length"`
	Processes    []sim.ProcessSnapshot `json:"processes"`
	Metrics      sim.SchedulingMetrics `json:"metrics"`
	Memory       sim.MemorySnapshot    `json:"memory"`
	FilesystemOK bool                  `json:"filesystem_ok"`
}

type ModuleAnalyticsDTO struct {
	Module         string  `json:"module"`
	TotalStages    int     `json:"total_stages"`
	CompletedStage int     `json:"completed_stage"`
	CompletionRate float64 `json:"completion_rate"`
}

type CompletionAnalyticsDTO struct {
	TotalStages      int                  `json:"total_stages"`
	CompletedStages  int                  `json:"completed_stages"`
	AttemptedStages  int                  `json:"attempted_stages"`
	CompletionRate   float64              `json:"completion_rate"`
	AttemptCoverage  float64              `json:"attempt_coverage"`
	ModuleBreakdown  []ModuleAnalyticsDTO `json:"module_breakdown"`
	PilotChecklist   []string             `json:"pilot_checklist"`
	PilotChecklistOK bool                 `json:"pilot_checklist_ok"`
}

type LessonRunResponse struct {
	LessonID    string                 `json:"lesson_id"`
	StageIndex  int                    `json:"stage_index"`
	Passed      bool                   `json:"passed"`
	FeedbackKey string                 `json:"feedback_key"`
	Hint        string                 `json:"hint,omitempty"`
	HintLevel   int                    `json:"hint_level,omitempty"`
	Output      LessonOutputDTO        `json:"output"`
	Analytics   CompletionAnalyticsDTO `json:"analytics"`
}

func (s *Server) handleLessons(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, r, http.StatusMethodNotAllowed, "method_not_allowed", "method not allowed")
		return
	}

	s.lessonMu.Lock()
	lessonsList := s.lessonEngine.Lessons()
	s.lessonMu.Unlock()

	out := make([]LessonSummary, 0, len(lessonsList))
	for _, lesson := range lessonsList {
		stages := make([]LessonStageSummary, 0, len(lesson.Stages))
		for idx, stage := range lesson.Stages {
			stages = append(stages, LessonStageSummary{Index: idx, ID: stage.ID, Title: stage.Title})
		}
		out = append(out, LessonSummary{ID: lesson.ID, Title: lesson.Title, Module: lesson.Module, Stages: stages})
	}

	respondJSON(w, http.StatusOK, LessonsResponse{Lessons: out})
}

func (s *Server) handleLessonRun(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, r, http.StatusMethodNotAllowed, "method_not_allowed", "method not allowed")
		return
	}

	defer func() { _ = r.Body.Close() }()
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	var req LessonRunRequest
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&req); err != nil {
		respondError(w, r, http.StatusBadRequest, "invalid_body", "invalid JSON body")
		return
	}
	if req.LessonID == "" {
		respondError(w, r, http.StatusBadRequest, "invalid_lesson_request", "lesson_id is required")
		return
	}

	s.lessonMu.Lock()
	result, err := s.lessonEngine.RunStage(req.LessonID, req.StageIndex)
	if err != nil {
		s.lessonMu.Unlock()
		respondError(w, r, http.StatusBadRequest, "lesson_run_failed", err.Error())
		return
	}
	analytics := s.lessonEngine.CompletionAnalytics()
	s.lessonMu.Unlock()

	respondJSON(w, http.StatusOK, LessonRunResponse{
		LessonID:    req.LessonID,
		StageIndex:  req.StageIndex,
		Passed:      result.Passed,
		FeedbackKey: result.FeedbackKey,
		Hint:        result.Hint,
		HintLevel:   result.HintLevel,
		Output: LessonOutputDTO{
			Tick:         result.Output.Metrics.TotalTicks,
			TraceHash:    sim.TraceHash(result.Output.Trace),
			TraceLength:  len(result.Output.Trace),
			Processes:    result.Output.Processes,
			Metrics:      result.Output.Metrics,
			Memory:       result.Output.Memory,
			FilesystemOK: result.Output.FilesystemOK,
		},
		Analytics: convertAnalytics(analytics),
	})
}

func convertAnalytics(in lessons.CompletionAnalytics) CompletionAnalyticsDTO {
	modules := make([]ModuleAnalyticsDTO, 0, len(in.ModuleBreakdown))
	for _, module := range in.ModuleBreakdown {
		modules = append(modules, ModuleAnalyticsDTO{
			Module:         module.Module,
			TotalStages:    module.TotalStages,
			CompletedStage: module.CompletedStage,
			CompletionRate: module.CompletionRate,
		})
	}
	return CompletionAnalyticsDTO{
		TotalStages:      in.TotalStages,
		CompletedStages:  in.CompletedStages,
		AttemptedStages:  in.AttemptedStages,
		CompletionRate:   in.CompletionRate,
		AttemptCoverage:  in.AttemptCoverage,
		ModuleBreakdown:  modules,
		PilotChecklist:   append([]string(nil), in.PilotChecklist...),
		PilotChecklistOK: in.PilotChecklistOK,
	}
}
