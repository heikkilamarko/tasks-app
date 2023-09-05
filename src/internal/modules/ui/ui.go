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

type UI struct {
	Config       *shared.Config
	Logger       *slog.Logger
	TaskRepo     shared.TaskRepository
	FileExporter shared.FileExporter
	server       *http.Server
}

func (*UI) Name() string { return "ui" }

func (s *UI) Run(ctx context.Context) error {
	router := chi.NewRouter()

	router.Use(middleware.Recoverer)
	router.Handle("/ui/static/*", http.StripPrefix("/ui", http.FileServer(http.FS(StaticFS))))
	router.Method(http.MethodGet, "/ui", &GetUI{s.TaskRepo, s.Logger})
	router.Method(http.MethodGet, "/ui/tasks", &GetUITasks{s.TaskRepo, s.Logger})
	router.Method(http.MethodGet, "/ui/tasks/export", &GetUITasksExport{s.TaskRepo, s.FileExporter, s.Logger})
	router.Method(http.MethodGet, "/ui/tasks/new", &GetUITasksNew{s.TaskRepo, s.Logger})
	router.Method(http.MethodGet, "/ui/tasks/{id}", &GetUITask{s.TaskRepo, s.Logger})
	router.Method(http.MethodGet, "/ui/tasks/{id}/edit", &GetUITaskEdit{s.TaskRepo, s.Logger})
	router.Method(http.MethodPost, "/ui/tasks", &PostUITasks{s.TaskRepo, s.Logger})
	router.Method(http.MethodPost, "/ui/tasks/{id}/complete", &PostUITaskComplete{s.TaskRepo, s.Logger})
	router.Method(http.MethodPut, "/ui/tasks/{id}", &PutUITask{s.TaskRepo, s.Logger})
	router.Method(http.MethodDelete, "/ui/tasks/{id}", &DeleteUITask{s.TaskRepo, s.Logger})
	router.Method(http.MethodGet, "/ui/completed", &GetUICompleted{s.TaskRepo, s.Logger})
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

func (s *UI) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return s.server.Shutdown(ctx)
}
