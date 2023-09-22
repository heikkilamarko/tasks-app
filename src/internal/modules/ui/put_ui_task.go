package ui

import (
	"log/slog"
	"net/http"
	"tasks-app/internal/shared"

	"github.com/go-chi/chi/v5"
)

type PutUITask struct {
	TaskRepository shared.TaskRepository
	Logger         *slog.Logger
}

func (h *PutUITask) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	id, ok := ValidateID(chi.URLParam(r, "id"))
	if !ok {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	name, ok := ValidateName(r.FormValue("name"))
	if !ok {
		http.Error(w, "invalid name", http.StatusBadRequest)
		return
	}

	expiresAt, ok := ValidateExpiresAt(r.FormValue("expires_at"))
	if !ok {
		http.Error(w, "invalid expires_at", http.StatusBadRequest)
		return
	}

	task, err := h.TaskRepository.GetByID(r.Context(), id)
	if err != nil {
		h.Logger.Error("get task", "error", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	if task == nil {
		http.Error(w, "task not found", http.StatusNotFound)
		return
	}

	task.Update(name, expiresAt)

	err = h.TaskRepository.Update(r.Context(), task)
	if err != nil {
		h.Logger.Error("update task", "error", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	if err := Templates.ExecuteTemplate(w, "active_tasks_table_row", task); err != nil {
		h.Logger.Error("execute template", "error", err)
		http.Error(w, "", http.StatusInternalServerError)
	}
}
