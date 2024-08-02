package ui

import (
	"log/slog"
	"net/http"
)

type PostUITimezone struct {
	Logger *slog.Logger
}

func (h *PostUITimezone) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	req, err := ParseSetTimezoneRequest(r)
	if err != nil {
		h.Logger.Error("parse request", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	SetTimezoneCookie(w, req.Timezone)

	http.Redirect(w, r, GetRedirectURL(r), http.StatusFound)
}
