package ui

import (
	"context"
	"log/slog"
	"net/http"
	"tasks-app/internal/shared"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"golang.org/x/sync/errgroup"
)

type Module struct {
	Config                    *shared.Config
	Logger                    *slog.Logger
	Auth                      *Auth
	TaskRepository            shared.TaskRepository
	TaskAttachmentsRepository shared.TaskAttachmentsRepository
	FileExporter              shared.FileExporter
}

func (m *Module) Run(ctx context.Context) error {

	auth, err := NewAuth(ctx, m.Config)
	if err != nil {
		return err
	}

	m.Auth = auth

	router := chi.NewRouter()

	router.Use(middleware.Recoverer)

	m.Auth.RegisterRoutes(router)

	router.Handle("/ui/static/*", http.StripPrefix("/ui", http.FileServer(http.FS(StaticFS))))
	router.Method(http.MethodGet, "/ui", m.Auth.Middleware.RequireAuthentication()(&GetUI{m.TaskRepository, m.Logger, m.Auth}))
	router.Method(http.MethodGet, "/ui/tasks", m.Auth.Middleware.RequireAuthentication()(&GetUITasks{m.TaskRepository, m.Logger}))
	router.Method(http.MethodGet, "/ui/tasks/export", m.Auth.Middleware.RequireAuthentication()(&GetUITasksExport{m.TaskRepository, m.FileExporter, m.Logger}))
	router.Method(http.MethodGet, "/ui/tasks/new", m.Auth.Middleware.RequireAuthentication()(&GetUITasksNew{m.TaskRepository, m.Logger}))
	router.Method(http.MethodGet, "/ui/tasks/{id}", m.Auth.Middleware.RequireAuthentication()(&GetUITask{m.TaskRepository, m.Logger}))
	router.Method(http.MethodGet, "/ui/tasks/{id}/edit", m.Auth.Middleware.RequireAuthentication()(&GetUITaskEdit{m.TaskRepository, m.Logger}))
	router.Method(http.MethodGet, "/ui/tasks/{id}/attachments/{name}", m.Auth.Middleware.RequireAuthentication()(&GetUITaskAttachment{m.TaskAttachmentsRepository, m.Logger}))
	router.Method(http.MethodPost, "/ui/tasks", m.Auth.Middleware.RequireAuthentication()(&PostUITasks{m.TaskRepository, m.TaskAttachmentsRepository, m.Logger}))
	router.Method(http.MethodPost, "/ui/tasks/{id}/complete", m.Auth.Middleware.RequireAuthentication()(&PostUITaskComplete{m.TaskRepository, m.Logger}))
	router.Method(http.MethodPut, "/ui/tasks/{id}", m.Auth.Middleware.RequireAuthentication()(&PutUITask{m.TaskRepository, m.TaskAttachmentsRepository, m.Logger}))
	router.Method(http.MethodDelete, "/ui/tasks/{id}", m.Auth.Middleware.RequireAuthentication()(&DeleteUITask{m.TaskRepository, m.TaskAttachmentsRepository, m.Logger}))
	router.Method(http.MethodGet, "/ui/completed", m.Auth.Middleware.RequireAuthentication()(&GetUICompleted{m.TaskRepository, m.Logger}))
	router.NotFound(NotFound)

	server := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
		Addr:         m.Config.UI.Addr,
		Handler:      router,
	}

	g := &errgroup.Group{}

	g.Go(func() error {
		<-ctx.Done()
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return server.Shutdown(ctx)
	})

	m.Logger.Info("run http server", "addr", server.Addr)

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}

	return g.Wait()
}
