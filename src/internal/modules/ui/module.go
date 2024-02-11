package ui

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"tasks-app/internal/shared"
	"time"

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
	if err := m.initAuth(ctx); err != nil {
		return fmt.Errorf("init auth: %w", err)
	}

	errorMW := ErrorRecoveryMiddleware(m.Logger)
	authnMW := m.Auth.Middleware.RequireAuthentication()

	mux := http.NewServeMux()

	mux.Handle("GET /ui/auth/login", errorMW(m.Auth.LoginHandler("/ui")))
	mux.Handle("GET /ui/auth/callback", errorMW(m.Auth.CallbackHandler()))
	mux.Handle("GET /ui/auth/logout", errorMW(m.Auth.LogoutHandler()))
	mux.Handle("GET /ui/static/*", errorMW(http.StripPrefix("/ui", http.FileServer(http.FS(StaticFS)))))
	mux.Handle("GET /ui", errorMW(authnMW(&GetUI{m.TaskRepository, m.Logger, m.Auth})))
	mux.Handle("POST /ui/theme", errorMW(authnMW(&PostUITheme{m.Logger})))
	mux.Handle("GET /ui/tasks", errorMW(authnMW(&GetUITasks{m.TaskRepository, m.Logger})))
	mux.Handle("GET /ui/tasks/export", errorMW(authnMW(&GetUITasksExport{m.TaskRepository, m.FileExporter, m.Logger})))
	mux.Handle("GET /ui/tasks/new", errorMW(authnMW(&GetUITasksNew{m.TaskRepository, m.Logger})))
	mux.Handle("GET /ui/tasks/{id}", errorMW(authnMW(&GetUITask{m.TaskRepository, m.Logger})))
	mux.Handle("GET /ui/tasks/{id}/edit", errorMW(authnMW(&GetUITaskEdit{m.TaskRepository, m.Logger})))
	mux.Handle("GET /ui/tasks/{id}/attachments/{name}", errorMW(authnMW(&GetUITaskAttachment{m.TaskAttachmentsRepository, m.Logger})))
	mux.Handle("POST /ui/tasks", errorMW(authnMW(&PostUITasks{m.TaskRepository, m.TaskAttachmentsRepository, m.Logger})))
	mux.Handle("POST /ui/tasks/{id}/complete", errorMW(authnMW(&PostUITaskComplete{m.TaskRepository, m.Logger})))
	mux.Handle("PUT /ui/tasks/{id}", errorMW(authnMW(&PutUITask{m.TaskRepository, m.TaskAttachmentsRepository, m.Logger})))
	mux.Handle("DELETE /ui/tasks/{id}", errorMW(authnMW(&DeleteUITask{m.TaskRepository, m.TaskAttachmentsRepository, m.Logger})))
	mux.Handle("GET /ui/completed", errorMW(authnMW(&GetUICompleted{m.TaskRepository, m.Logger})))

	server := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
		Addr:         m.Config.UI.Addr,
		Handler:      mux,
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

func (m *Module) initAuth(ctx context.Context) error {
	auth, err := NewAuth(ctx, m.Config)
	if err != nil {
		return err
	}

	m.Auth = auth
	return nil
}
