package ui

import (
	"log/slog"
	"net/http"
	"tasks-app/internal/shared"
)

type GetUITasks struct {
	TaskRepository shared.TaskRepository
	Renderer       Renderer
	Logger         *slog.Logger
}

func (h *GetUITasks) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.TaskRepository.GetActive(r.Context(), 0, 50)
	if err != nil {
		h.Logger.Error("get tasks", "error", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	vm := NewTasksResponse(r, tasks)

	h.Renderer.Render(w, "active_tasks_table", vm)
}
