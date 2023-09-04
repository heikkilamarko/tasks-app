package ui

import (
	"context"
	"log/slog"
	"net/http"
	"tasks-app/internal/shared"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type HTTPService struct {
	Config       *shared.Config
	Logger       *slog.Logger
	TaskRepo     shared.TaskRepository
	FileExporter shared.FileExporter
	server       *http.Server
}

func (*HTTPService) Name() string { return "ui" }

func (s *HTTPService) Run(ctx context.Context) error {
	router := chi.NewRouter()

	router.Use(middleware.Recoverer)
	router.Handle("/ui/static/*", http.StripPrefix("/ui", http.FileServer(http.FS(StaticFS))))
	router.Method(http.MethodGet, "/ui", &GetUIHandler{s.TaskRepo, s.Logger})
	router.Method(http.MethodGet, "/ui/tasks", &GetUITasksHandler{s.TaskRepo, s.Logger})
	router.Method(http.MethodGet, "/ui/tasks/export", &GetUITasksExportHandler{s.TaskRepo, s.FileExporter, s.Logger})
	router.Method(http.MethodGet, "/ui/tasks/new", &GetUITaskNewHandler{s.TaskRepo, s.Logger})
	router.Method(http.MethodGet, "/ui/tasks/{id}", &GetUITaskHandler{s.TaskRepo, s.Logger})
	router.Method(http.MethodGet, "/ui/tasks/{id}/edit", &GetUITaskEditHandler{s.TaskRepo, s.Logger})
	router.Method(http.MethodPost, "/ui/tasks", &PostUITaskHandler{s.TaskRepo, s.Logger})
	router.Method(http.MethodPost, "/ui/tasks/{id}/complete", &PostUITaskCompleteHandler{s.TaskRepo, s.Logger})
	router.Method(http.MethodPut, "/ui/tasks/{id}", &PutUITaskHandler{s.TaskRepo, s.Logger})
	router.Method(http.MethodDelete, "/ui/tasks/{id}", &DeleteUITaskHandler{s.TaskRepo, s.Logger})
	router.Method(http.MethodGet, "/ui/completed", &GetUICompletedHandler{s.TaskRepo, s.Logger})
	router.NotFound(NotFound)

	s.server = &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
		Addr:         s.Config.Addr,
		Handler:      router,
	}

	s.Logger.Info("http server is running", "port", s.server.Addr)

	if err := s.server.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}

	return nil
}

func (s *HTTPService) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return s.server.Shutdown(ctx)
}
