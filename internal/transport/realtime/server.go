package realtime

import (
	"net/http"
	"os"
	"sync"

	appchallenges "os-simulator-plan/internal/app/challenges"
	contentv3 "os-simulator-plan/internal/content/v3"
	"os-simulator-plan/internal/lessons"

	"github.com/go-chi/chi/v5"
)

type Server struct {
	manager           *SessionManager
	challengeAttempts *ChallengeAttemptStore
	challengeService  *appchallenges.Service
	lessonMu          sync.Mutex
	lessonCatalog     map[string]lessons.Lesson
	lessonEngines     map[string]*lessons.Engine
	cpuCurriculumV3   contentv3.Curriculum
	cpuCurriculumErr  error
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
	cpuCurriculumV3, cpuCurriculumErr := contentv3.LoadCPUCurriculum()
	server := &Server{
		manager:           manager,
		challengeAttempts: NewChallengeAttemptStore(),
		lessonCatalog:     catalog,
		lessonEngines:     map[string]*lessons.Engine{},
		cpuCurriculumV3:   cpuCurriculumV3,
		cpuCurriculumErr:  cpuCurriculumErr,
		origins:           origins,
	}
	server.challengeService = appchallenges.NewService(
		lessonEngineProviderAdapter{server: server},
		sessionStoreAdapter{manager: server.manager},
		challengeAttemptStoreAdapter{store: server.challengeAttempts},
	)
	return server
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
	router.HandleFunc("/curriculum/v3", s.handleCurriculumV3)
	router.HandleFunc("/lessons/{lessonID}/learn/v3", s.handleLessonLearnV3)
	router.HandleFunc("/lessons/{lessonID}/challenge/v3", s.handleChallengeManifestV3)
	router.HandleFunc("/challenges/start/v3", s.handleChallengeStartV3)
	router.HandleFunc("/challenges/action/v3", s.handleChallengeActionV3)
	router.HandleFunc("/challenges/submit/v3", s.handleChallengeSubmitV3)
	router.HandleFunc("/challenges/{attemptID}/replay/v3", s.handleChallengeReplayV3)
	return withRequestID(withCORS(s.origins, router))
}
