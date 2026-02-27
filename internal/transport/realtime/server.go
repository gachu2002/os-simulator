package realtime

import (
	"net/http"
	"os"
	"sync"

	"os-simulator-plan/internal/lessons"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
)

type Server struct {
	manager           *SessionManager
	challengeAttempts *ChallengeAttemptStore
	lessonMu          sync.Mutex
	lessonCatalog     map[string]lessons.Lesson
	lessonEngines     map[string]*lessons.Engine
	upgrader          websocket.Upgrader
	origins           map[string]struct{}
}

func NewServer(manager *SessionManager) *Server {
	return NewServerWithLessons(manager, lessons.NewEngine())
}

func NewServerWithLessons(manager *SessionManager, lessonEngine *lessons.Engine) *Server {
	origins := allowedOriginsFromEnv(os.Getenv("CORS_ALLOW_ORIGIN"))
	catalog := make(map[string]lessons.Lesson)
	for _, lesson := range lessonEngine.Lessons() {
		catalog[lesson.ID] = lesson
	}
	return &Server{
		manager:           manager,
		challengeAttempts: NewChallengeAttemptStore(),
		lessonCatalog:     catalog,
		lessonEngines:     map[string]*lessons.Engine{},
		origins:           origins,
		upgrader: websocket.Upgrader{
			ReadBufferSize:  4096,
			WriteBufferSize: 4096,
			CheckOrigin: func(r *http.Request) bool {
				return isOriginAllowed(origins, r.Header.Get("Origin"))
			},
		},
	}
}

func (s *Server) lessonEngineForLearner(learnerID string) *lessons.Engine {
	s.lessonMu.Lock()
	defer s.lessonMu.Unlock()
	if learnerID == "" {
		learnerID = "anonymous"
	}
	engine, ok := s.lessonEngines[learnerID]
	if ok {
		return engine
	}
	engine = lessons.NewEngineWithCatalog(s.lessonCatalog)
	s.lessonEngines[learnerID] = engine
	return engine
}

func (s *Server) Handler() http.Handler {
	router := chi.NewRouter()
	router.HandleFunc("/healthz", s.handleHealth)
	router.HandleFunc("/curriculum", s.handleCurriculum)
	router.HandleFunc("/lessons/{lessonID}/learn", s.handleLessonLearn)
	router.HandleFunc("/challenges/start", s.handleChallengeStart)
	router.HandleFunc("/challenges/submit", s.handleChallengeSubmit)
	router.HandleFunc("/ws/{id}", s.handleWS)
	return withRequestID(withCORS(s.origins, router))
}
