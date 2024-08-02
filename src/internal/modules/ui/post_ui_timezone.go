package ui

import (
	"log/slog"
	"net/http"
	"tasks-app/internal/shared"
)

type PostUITimezone struct {
	Config *shared.Config
	Logger *slog.Logger
}

func (h *PostUITimezone) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	redirectURL, err := GetRedirectURL(r, h.Config.UI.TrustedHosts)
	if err != nil {
		h.Logger.Error("get redirect url", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	req, err := ParseSetTimezoneRequest(r)
	if err != nil {
		h.Logger.Error("parse request", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	SetTimezoneCookie(w, req.Timezone)

	http.Redirect(w, r, redirectURL, http.StatusFound)
}
