package ui

import (
	"log/slog"
	"net/http"
)

type GetUILanguage struct {
	Logger *slog.Logger
}

func (h *GetUILanguage) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	req, err := ParseSetLanguageRequest(r)
	if err != nil {
		h.Logger.Error("parse request", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	SetLanguageCookie(w, req.Language)

	w.Header().Add("HX-Redirect", GetRedirectURL(r))
}
