package ui

import (
	"log/slog"
	"net/http"
)

type GetUITheme struct {
	Logger *slog.Logger
}

func (h *GetUITheme) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	req, err := ParseSetThemeRequest(r)
	if err != nil {
		h.Logger.Error("parse request", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	SetThemeCookie(w, req.Theme)

	w.Header().Add("HX-Redirect", GetRedirectURL(r))
}
