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
	TaskRepository            shared.TaskRepository
	TaskAttachmentsRepository shared.TaskAttachmentsRepository
	FileExporter              shared.FileExporter
}

func (m *Module) Run(ctx context.Context) error {
	router := chi.NewRouter()

	router.Use(middleware.Recoverer)
	router.Use(SessionMiddleware)

	router.Handle("/ui/static/*", http.StripPrefix("/ui", http.FileServer(http.FS(StaticFS))))
	router.Handle("/ui/attachments/*", http.StripPrefix("/ui/attachments/", http.FileServer(http.Dir(m.Config.Shared.AttachmentsPath))))
	router.Method(http.MethodGet, "/ui", &GetUI{m.TaskRepository, m.Logger})
	router.Method(http.MethodGet, "/ui/tasks", &GetUITasks{m.TaskRepository, m.Logger})
	router.Method(http.MethodGet, "/ui/tasks/export", &GetUITasksExport{m.TaskRepository, m.FileExporter, m.Logger})
	router.Method(http.MethodGet, "/ui/tasks/new", &GetUITasksNew{m.TaskRepository, m.Logger})
	router.Method(http.MethodGet, "/ui/tasks/{id}", &GetUITask{m.TaskRepository, m.Logger})
	router.Method(http.MethodGet, "/ui/tasks/{id}/edit", &GetUITaskEdit{m.TaskRepository, m.Logger})
	router.Method(http.MethodPost, "/ui/tasks", &PostUITasks{m.TaskRepository, m.TaskAttachmentsRepository, m.Logger})
	router.Method(http.MethodPost, "/ui/tasks/{id}/complete", &PostUITaskComplete{m.TaskRepository, m.Logger})
	router.Method(http.MethodPut, "/ui/tasks/{id}", &PutUITask{m.TaskRepository, m.TaskAttachmentsRepository, m.Logger})
	router.Method(http.MethodDelete, "/ui/tasks/{id}", &DeleteUITask{m.TaskRepository, m.TaskAttachmentsRepository, m.Logger})
	router.Method(http.MethodGet, "/ui/completed", &GetUICompleted{m.TaskRepository, m.Logger})
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
