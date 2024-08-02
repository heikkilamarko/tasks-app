package ui

import (
	"log/slog"
	"net/http"
	"tasks-app/internal/shared"
)

type PostUITheme struct {
	Config *shared.Config
	Logger *slog.Logger
}

func (h *PostUITheme) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	redirectURL, err := GetRedirectURL(r, h.Config.UI.TrustedHosts)
	if err != nil {
		h.Logger.Error("get redirect url", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	req, err := ParseSetThemeRequest(r)
	if err != nil {
		h.Logger.Error("parse request", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	SetThemeCookie(w, req.Theme)

	http.Redirect(w, r, redirectURL, http.StatusFound)
}
