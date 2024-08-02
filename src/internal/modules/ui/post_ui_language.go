package ui

import (
	"log/slog"
	"net/http"
)

type PostUILanguage struct {
	Logger *slog.Logger
}

func (h *PostUILanguage) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	req, err := ParseSetLanguageRequest(r)
	if err != nil {
		h.Logger.Error("parse request", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	SetLanguageCookie(w, req.Language)

	http.Redirect(w, r, GetRedirectURL(r), http.StatusFound)
}
