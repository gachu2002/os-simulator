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
	return withRequestID(withCORS(s.origins, router))
}
