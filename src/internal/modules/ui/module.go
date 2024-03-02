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
	Renderer                  Renderer
	TaskRepository            shared.TaskRepository
	TaskAttachmentsRepository shared.TaskAttachmentsRepository
	FileExporter              shared.FileExporter
}

func (m *Module) Run(ctx context.Context) error {
	if err := m.initAuth(ctx); err != nil {
		return fmt.Errorf("init auth: %w", err)
	}

	if err := m.initRenderer(); err != nil {
		return fmt.Errorf("init renderer: %w", err)
	}

	errorMW := ErrorRecoveryMiddleware(m.Logger)
	authnMW := m.Auth.Middleware.RequireAuthentication()
	userMW := UserContextMiddleware(m.Auth)
	loginMW := LoginMiddleware(m.Auth, m.Config)

	mux := http.NewServeMux()

	HandleWithMiddleware(mux, "GET /ui/auth/login", m.Auth.LoginHandler("/ui"), errorMW)
	HandleWithMiddleware(mux, "GET /ui/auth/callback", m.Auth.CallbackHandler(), errorMW)
	HandleWithMiddleware(mux, "GET /ui/auth/logout", m.Auth.LogoutHandler(), errorMW)
	HandleWithMiddleware(mux, "GET /ui/static/*", http.StripPrefix("/ui", http.FileServer(http.FS(StaticFS))), errorMW)
	HandleWithMiddleware(mux, "GET /ui", &GetUI{m.TaskRepository, m.Renderer, m.Logger}, errorMW, authnMW, userMW, loginMW)
	HandleWithMiddleware(mux, "POST /ui/theme", &PostUITheme{m.Logger}, errorMW, authnMW, userMW)
	HandleWithMiddleware(mux, "GET /ui/tasks", &GetUITasks{m.TaskRepository, m.Renderer, m.Logger}, errorMW, authnMW, userMW)
	HandleWithMiddleware(mux, "GET /ui/tasks/export", &GetUITasksExport{m.TaskRepository, m.FileExporter, m.Logger}, errorMW, authnMW, userMW)
	HandleWithMiddleware(mux, "GET /ui/tasks/new", &GetUITasksNew{m.TaskRepository, m.Renderer, m.Logger}, errorMW, authnMW, userMW)
	HandleWithMiddleware(mux, "GET /ui/tasks/{id}", &GetUITask{m.TaskRepository, m.Renderer, m.Logger}, errorMW, authnMW, userMW)
	HandleWithMiddleware(mux, "GET /ui/tasks/{id}/edit", &GetUITaskEdit{m.TaskRepository, m.Renderer, m.Logger}, errorMW, authnMW, userMW)
	HandleWithMiddleware(mux, "GET /ui/tasks/{id}/attachments/{name}", &GetUITaskAttachment{m.TaskAttachmentsRepository, m.Logger}, errorMW, authnMW, userMW)
	HandleWithMiddleware(mux, "POST /ui/tasks", &PostUITasks{m.TaskRepository, m.TaskAttachmentsRepository, m.Renderer, m.Logger}, errorMW, authnMW, userMW)
	HandleWithMiddleware(mux, "POST /ui/tasks/{id}/complete", &PostUITaskComplete{m.TaskRepository, m.Renderer, m.Logger}, errorMW, authnMW, userMW)
	HandleWithMiddleware(mux, "PUT /ui/tasks/{id}", &PutUITask{m.TaskRepository, m.TaskAttachmentsRepository, m.Renderer, m.Logger}, errorMW, authnMW, userMW)
	HandleWithMiddleware(mux, "DELETE /ui/tasks/{id}", &DeleteUITask{m.TaskRepository, m.TaskAttachmentsRepository, m.Renderer, m.Logger}, errorMW, authnMW, userMW)
	HandleWithMiddleware(mux, "GET /ui/completed", &GetUICompleted{m.TaskRepository, m.Renderer, m.Logger}, errorMW, authnMW, userMW)

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

func (m *Module) initRenderer() error {
	renderer, err := NewTemplateRenderer(m.Logger)
	if err != nil {
		return err
	}

	m.Renderer = renderer
	return nil
}
