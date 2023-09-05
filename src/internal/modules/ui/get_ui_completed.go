package ui

import (
	"log/slog"
	"net/http"
	"tasks-app/internal/shared"
)

type GetUICompleted struct {
	TaskRepo shared.TaskRepository
	Logger   *slog.Logger
}

func (h *GetUICompleted) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.TaskRepo.GetCompleted(r.Context())
	if err != nil {
		h.Logger.Error("get tasks", "error", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	data := struct {
		Tasks []*shared.Task
	}{
		Tasks: tasks,
	}

	if err := Templates.ExecuteTemplate(w, "completed_tasks.html", data); err != nil {
		h.Logger.Error("execute template", "error", err)
		http.Error(w, "", http.StatusInternalServerError)
	}
}
