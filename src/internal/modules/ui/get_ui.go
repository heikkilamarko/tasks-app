package ui

import (
	"log/slog"
	"net/http"
	"tasks-app/internal/shared"
)

type GetUI struct {
	TaskRepository shared.TaskRepository
	Renderer       Renderer
	Auth           *Auth
	Logger         *slog.Logger
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

	vm := NewTasksResponse(tasks).WithTheme(r)

	h.Renderer.Render(w, "active_tasks", vm)
}
