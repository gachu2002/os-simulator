package realtime

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"sync"

	"os-simulator-plan/internal/lessons"
	"os-simulator-plan/internal/sim"

	"github.com/gorilla/websocket"
)

type Server struct {
	manager      *SessionManager
	lessonMu     sync.Mutex
	lessonEngine *lessons.Engine
	upgrader     websocket.Upgrader
}

func NewServer(manager *SessionManager) *Server {
	return NewServerWithLessons(manager, lessons.NewEngine())
}

func NewServerWithLessons(manager *SessionManager, lessonEngine *lessons.Engine) *Server {
	return &Server{
		manager:      manager,
		lessonEngine: lessonEngine,
		upgrader: websocket.Upgrader{
			ReadBufferSize:  4096,
			WriteBufferSize: 4096,
			CheckOrigin: func(*http.Request) bool {
				return true
			},
		},
	}
}

func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", s.handleHealth)
	mux.HandleFunc("/sessions", s.handleSessions)
	mux.HandleFunc("/lessons", s.handleLessons)
	mux.HandleFunc("/lessons/run", s.handleLessonRun)
	mux.HandleFunc("/ws/", s.handleWS)
	return withCORS(mux)
}

func (s *Server) handleHealth(w http.ResponseWriter, _ *http.Request) {
	respondJSON(w, http.StatusOK, map[string]any{"status": "ok", "protocol_version": ProtocolVersion})
}

func (s *Server) handleSessions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}
	defer r.Body.Close()
	var cfg SessionConfig
	if err := json.NewDecoder(r.Body).Decode(&cfg); err != nil && err.Error() != "EOF" {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON body"})
		return
	}
	session, err := s.manager.Create(cfg)
	if err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	ev := session.SnapshotEvent("init")
	respondJSON(w, http.StatusCreated, CreateSessionResponse{SessionID: session.ID(), Snapshot: ev.Snapshot})
}

func (s *Server) handleWS(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/ws/")
	if id == "" {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "missing session id"})
		return
	}
	session, ok := s.manager.Get(id)
	if !ok {
		respondJSON(w, http.StatusNotFound, map[string]string{"error": "session not found"})
		return
	}
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	if err := conn.WriteJSON(session.SnapshotEvent("connected")); err != nil {
		return
	}

	for {
		var req CommandEnvelope
		if err := conn.ReadJSON(&req); err != nil {
			return
		}
		if req.Type != "command" {
			if err := conn.WriteJSON(Event{Type: "session.error", SessionID: session.ID(), Error: "unsupported message type"}); err != nil {
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

func withCORS(next http.Handler) http.Handler {
	origin := os.Getenv("CORS_ALLOW_ORIGIN")
	if origin == "" {
		origin = "*"
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Vary", "Origin")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
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
		respondJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
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
		respondJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}

	defer r.Body.Close()
	var req LessonRunRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON body"})
		return
	}
	if req.LessonID == "" {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "lesson_id is required"})
		return
	}

	s.lessonMu.Lock()
	result, err := s.lessonEngine.RunStage(req.LessonID, req.StageIndex)
	if err != nil {
		s.lessonMu.Unlock()
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
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
