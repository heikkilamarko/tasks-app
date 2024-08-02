package ui

import (
	"log/slog"
	"net/http"
	"tasks-app/internal/shared"
)

type PostUILanguage struct {
	Config *shared.Config
	Logger *slog.Logger
}

func (h *PostUILanguage) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	redirectURL, err := GetRedirectURL(r, h.Config.UI.TrustedHosts)
	if err != nil {
		h.Logger.Error("get redirect url", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	req, err := ParseSetLanguageRequest(r)
	if err != nil {
		h.Logger.Error("parse request", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	SetLanguageCookie(w, req.Language)

	http.Redirect(w, r, redirectURL, http.StatusFound)
}
