package ui

import (
	"log/slog"
	"net/http"
	"tasks-app/internal/shared"
)

type GetUITasksNew struct {
	TaskRepository shared.TaskRepository
	Logger         *slog.Logger
}

func (h *GetUITasksNew) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.TaskRepository.GetActive(r.Context())
	if err != nil {
		h.Logger.Error("get tasks", "error", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	data := struct {
		Tasks         []*shared.Task
		IsCreatingNew bool
	}{
		Tasks:         tasks,
		IsCreatingNew: true,
	}

	if err := Templates.ExecuteTemplate(w, "active_tasks_table.html", data); err != nil {
		h.Logger.Error("execute template", "error", err)
		http.Error(w, "", http.StatusInternalServerError)
	}
}
