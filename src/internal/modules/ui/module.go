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
	TxManager                 shared.TxManager
	TaskRepository            shared.TaskRepository
	TaskAttachmentsRepository shared.TaskAttachmentsRepository
	FileExporter              shared.FileExporter
}

var _ shared.AppModule = (*Module)(nil)

func (m *Module) Run(ctx context.Context) error {
	if err := m.initAuth(ctx); err != nil {
		return fmt.Errorf("init auth: %w", err)
	}

	if err := m.initRenderer(); err != nil {
		return fmt.Errorf("init renderer: %w", err)
	}

	errorMW := ErrorRecoveryMiddleware(m.Logger)
	csrfMW := NewCSRF(m.Config).Middleware
	authnMW := m.Auth.Middleware.RequireAuthentication()
	userMW := UserContextMiddleware(m.Auth)
	loginMW := LoginMiddleware(m.Auth)

	mux := http.NewServeMux()

	HandleWithMiddleware(mux, "GET /ui/auth/login", m.Auth.LoginHandler("/ui"))
	HandleWithMiddleware(mux, "GET /ui/auth/callback", m.Auth.CallbackHandler())
	HandleWithMiddleware(mux, "GET /ui/auth/logout", m.Auth.LogoutHandler())
	HandleWithMiddleware(mux, "GET /ui/static/", http.StripPrefix("/ui", http.FileServer(http.FS(StaticFS))))
	HandleWithMiddleware(mux, "GET /ui", &GetUI{m.TaskRepository, m.Renderer, m.Logger}, authnMW, userMW, loginMW)
	HandleWithMiddleware(mux, "POST /ui/language", &PostUILanguage{m.Logger}, authnMW, userMW)
	HandleWithMiddleware(mux, "POST /ui/theme", &PostUITheme{m.Logger}, authnMW, userMW)
	HandleWithMiddleware(mux, "POST /ui/timezone", &PostUITimezone{m.Logger}, authnMW, userMW)
	HandleWithMiddleware(mux, "GET /ui/tasks", &GetUITasks{m.TaskRepository, m.Renderer, m.Logger}, authnMW, userMW)
	HandleWithMiddleware(mux, "GET /ui/tasks/export", &GetUITasksExport{m.TaskRepository, m.FileExporter, m.Logger}, authnMW, userMW)
	HandleWithMiddleware(mux, "GET /ui/tasks/new", &GetUITasksNew{m.TaskRepository, m.Renderer, m.Logger}, authnMW, userMW)
	HandleWithMiddleware(mux, "GET /ui/tasks/{id}", &GetUITask{m.TaskRepository, m.Renderer, m.Logger}, authnMW, userMW)
	HandleWithMiddleware(mux, "GET /ui/tasks/{id}/edit", &GetUITaskEdit{m.TaskRepository, m.Renderer, m.Logger}, authnMW, userMW)
	HandleWithMiddleware(mux, "GET /ui/tasks/{id}/attachments/{name}", &GetUITaskAttachment{m.TaskAttachmentsRepository, m.Logger}, authnMW, userMW)
	HandleWithMiddleware(mux, "POST /ui/tasks", &PostUITasks{m.TxManager, m.TaskRepository, m.TaskAttachmentsRepository, m.Renderer, m.Logger}, authnMW, userMW)
	HandleWithMiddleware(mux, "POST /ui/tasks/{id}/complete", &PostUITaskComplete{m.TaskRepository, m.Renderer, m.Logger}, authnMW, userMW)
	HandleWithMiddleware(mux, "PUT /ui/tasks/{id}", &PutUITask{m.TxManager, m.TaskRepository, m.TaskAttachmentsRepository, m.Renderer, m.Logger}, authnMW, userMW)
	HandleWithMiddleware(mux, "DELETE /ui/tasks/{id}", &DeleteUITask{m.TxManager, m.TaskRepository, m.TaskAttachmentsRepository, m.Renderer, m.Logger}, authnMW, userMW)
	HandleWithMiddleware(mux, "GET /ui/completed", &GetUICompleted{m.TaskRepository, m.Renderer, m.Logger}, authnMW, userMW)
	HandleWithMiddleware(mux, "GET /ui/completed/tasks", &GetUICompletedTasks{m.TaskRepository, m.Renderer, m.Logger}, authnMW, userMW)

	server := &http.Server{
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  120 * time.Second,
		Addr:         m.Config.UI.Addr,
		Handler:      errorMW(csrfMW(mux)),
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
