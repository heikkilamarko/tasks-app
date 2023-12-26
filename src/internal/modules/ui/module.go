package ui

import (
	"context"
	"crypto/tls"
	"log/slog"
	"net/http"
	"tasks-app/internal/shared"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httphelper "github.com/zitadel/oidc/v3/pkg/http"
	"github.com/zitadel/oidc/v3/pkg/oidc"
	"github.com/zitadel/zitadel-go/v3/pkg/authentication"
	aoidc "github.com/zitadel/zitadel-go/v3/pkg/authentication/oidc"
	"github.com/zitadel/zitadel-go/v3/pkg/zitadel"
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

	httphelper.DefaultHTTPClient.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	authN, err := authentication.New(
		ctx,
		zitadel.New(m.Config.UI.AuthDomain),
		m.Config.UI.AuthEncryptionKey,
		aoidc.WithCodeFlow[*aoidc.UserInfoContext[*oidc.IDTokenClaims, *oidc.UserInfo], *oidc.IDTokenClaims, *oidc.UserInfo](
			aoidc.PKCEAuthentication(
				m.Config.UI.AuthClientId,
				m.Config.UI.AuthRedirectURI,
				[]string{
					oidc.ScopeOpenID,
					oidc.ScopeProfile,
					oidc.ScopeEmail,
				},
				httphelper.NewCookieHandler(
					[]byte(m.Config.UI.AuthEncryptionKey),
					[]byte(m.Config.UI.AuthEncryptionKey),
				),
			),
		),
	)
	if err != nil {
		return err
	}

	mw := authentication.Middleware(authN)

	getUserName := func(r *http.Request) string {
		if actx := mw.Context(r.Context()); actx != nil {
			return actx.UserInfo.Name
		}
		return ""
	}

	router := chi.NewRouter()

	router.Use(middleware.Recoverer)

	router.Handle(m.Config.UI.AuthPath+"/login", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authN.Authenticate(w, r, "/ui")
	}))
	router.Handle(m.Config.UI.AuthPath+"/callback", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authN.Callback(w, r)
	}))
	router.Handle(m.Config.UI.AuthPath+"/logout", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authN.Logout(w, r)
	}))

	router.Handle("/ui/static/*", http.StripPrefix("/ui", http.FileServer(http.FS(StaticFS))))
	router.Method(http.MethodGet, "/ui", mw.RequireAuthentication()(&GetUI{m.TaskRepository, m.Logger, getUserName}))
	router.Method(http.MethodGet, "/ui/tasks", mw.RequireAuthentication()(&GetUITasks{m.TaskRepository, m.Logger}))
	router.Method(http.MethodGet, "/ui/tasks/export", mw.RequireAuthentication()(&GetUITasksExport{m.TaskRepository, m.FileExporter, m.Logger}))
	router.Method(http.MethodGet, "/ui/tasks/new", mw.RequireAuthentication()(&GetUITasksNew{m.TaskRepository, m.Logger}))
	router.Method(http.MethodGet, "/ui/tasks/{id}", mw.RequireAuthentication()(&GetUITask{m.TaskRepository, m.Logger}))
	router.Method(http.MethodGet, "/ui/tasks/{id}/edit", mw.RequireAuthentication()(&GetUITaskEdit{m.TaskRepository, m.Logger}))
	router.Method(http.MethodGet, "/ui/tasks/{id}/attachments/{name}", mw.RequireAuthentication()(&GetUITaskAttachment{m.TaskAttachmentsRepository, m.Logger}))
	router.Method(http.MethodPost, "/ui/tasks", mw.RequireAuthentication()(&PostUITasks{m.TaskRepository, m.TaskAttachmentsRepository, m.Logger}))
	router.Method(http.MethodPost, "/ui/tasks/{id}/complete", mw.RequireAuthentication()(&PostUITaskComplete{m.TaskRepository, m.Logger}))
	router.Method(http.MethodPut, "/ui/tasks/{id}", mw.RequireAuthentication()(&PutUITask{m.TaskRepository, m.TaskAttachmentsRepository, m.Logger}))
	router.Method(http.MethodDelete, "/ui/tasks/{id}", mw.RequireAuthentication()(&DeleteUITask{m.TaskRepository, m.TaskAttachmentsRepository, m.Logger}))
	router.Method(http.MethodGet, "/ui/completed", mw.RequireAuthentication()(&GetUICompleted{m.TaskRepository, m.Logger}))
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
