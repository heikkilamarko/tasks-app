package ui

import (
	"log/slog"
	"net/http"
	"tasks-app/internal/shared"
)

type GetUITask struct {
	TaskRepository shared.TaskRepository
	Logger         *slog.Logger
}

func (h *GetUITask) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	req, err := ParseTaskRequest(r)
	if err != nil {
		h.Logger.Error("parse request", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	task, err := h.TaskRepository.GetByID(r.Context(), req.ID)
	if err != nil {
		h.Logger.Error("get task", "error", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	if task == nil {
		http.Error(w, "task not found", http.StatusNotFound)
		return
	}

	if err := Templates.ExecuteTemplate(w, "active_tasks_table_row", task); err != nil {
		h.Logger.Error("execute template", "error", err)
		http.Error(w, "", http.StatusInternalServerError)
	}
}
