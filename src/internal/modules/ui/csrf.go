package ui

import (
	"net/http"
	"tasks-app/internal/shared"

	"github.com/gorilla/csrf"
)

type CSRF struct {
	Config     *shared.Config
	Middleware func(http.Handler) http.Handler
}

func NewCSRF(config *shared.Config) *CSRF {
	mw := csrf.Protect(
		[]byte(config.UI.AuthEncryptionKey),
		csrf.SameSite(csrf.SameSiteStrictMode),
		csrf.CookieName("csrf.token"),
		csrf.FieldName("csrf-token"),
		csrf.RequestHeader("X-CSRF-Token"),
	)
	return &CSRF{config, mw}
}
