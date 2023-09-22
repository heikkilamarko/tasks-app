package ui

import (
	"log/slog"
	"net/http"
	"tasks-app/internal/shared"

	"github.com/go-chi/chi/v5"
)

type DeleteUITask struct {
	TaskRepository shared.TaskRepository
	Logger         *slog.Logger
}

func (h *DeleteUITask) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	id, ok := ValidateID(chi.URLParam(r, "id"))
	if !ok {
		http.Error(w, "invalid id", http.StatusBadRequest)
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

	err = h.TaskRepository.Delete(r.Context(), id)
	if err != nil {
		h.Logger.Error("delete task", "error", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	tasks, err := h.TaskRepository.GetActive(r.Context())
	if err != nil {
		h.Logger.Error("get tasks", "error", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	vm := TasksViewModel{tasks, false}

	if err := Templates.ExecuteTemplate(w, "active_tasks_table", vm); err != nil {
		h.Logger.Error("execute template", "error", err)
		http.Error(w, "", http.StatusInternalServerError)
	}
}
