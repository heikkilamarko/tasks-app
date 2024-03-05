package ui

import (
	"log/slog"
	"net/http"
	"tasks-app/internal/shared"
)

type GetUICompleted struct {
	TaskRepository shared.TaskRepository
	Renderer       Renderer
	Logger         *slog.Logger
}

func (h *GetUICompleted) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.TaskRepository.GetCompleted(r.Context(), 0, 50)
	if err != nil {
		h.Logger.Error("get tasks", "error", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	vm := NewTasksResponse(r, tasks).
		WithTheme(r).
		WithUser(r).
		WithHubURL()
	vm.Title = "Completed"

	h.Renderer.Render(w, "completed_tasks", vm)
}
