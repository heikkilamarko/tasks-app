package ui

import (
	"log/slog"
	"net/http"
)

type PostUITheme struct {
	Logger *slog.Logger
}

func (h *PostUITheme) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	req, err := ParseSetThemeRequest(r)
	if err != nil {
		h.Logger.Error("parse request", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	SetThemeCookie(w, req.Theme)

	w.Header().Add("HX-Redirect", "/ui")
}
