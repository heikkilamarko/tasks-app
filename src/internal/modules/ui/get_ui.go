package ui

import (
	"log/slog"
	"net/http"
	"tasks-app/internal/shared"
)

type GetUI struct {
	TaskRepository shared.TaskRepository
	Logger         *slog.Logger
	Auth           *Auth
}

func (h *GetUI) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	authInfo := h.Auth.GetAuthInfo(r)

	h.Logger.Info("auth", slog.Any("info", authInfo))

	tasks, err := h.TaskRepository.GetActive(r.Context(), 0, 50)
	if err != nil {
		h.Logger.Error("get tasks", "error", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	vm := TasksResponse{tasks, false}

	if err := Templates.ExecuteTemplate(w, "active_tasks", vm); err != nil {
		h.Logger.Error("execute template", "error", err)
		http.Error(w, "", http.StatusInternalServerError)
	}
}
