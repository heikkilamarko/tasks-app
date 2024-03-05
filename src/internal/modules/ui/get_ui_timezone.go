package ui

import (
	"log/slog"
	"net/http"
)

type GetUITimezone struct {
	Logger *slog.Logger
}

func (h *GetUITimezone) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	req, err := ParseSetTimezoneRequest(r)
	if err != nil {
		h.Logger.Error("parse request", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	SetTimezoneCookie(w, req.Timezone)

	w.Header().Add("HX-Redirect", GetRedirectURL(r))
}
