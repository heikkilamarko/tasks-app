package ui

import (
	"log/slog"
	"net/http"
	"strconv"
	"tasks-app/internal/shared"

	"github.com/go-chi/chi/v5"
)

type GetUITaskEdit struct {
	TaskRepository shared.TaskRepository
	Logger         *slog.Logger
}

func (h *GetUITaskEdit) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
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

	if err := Templates.ExecuteTemplate(w, "active_tasks_table_row_edit.html", task); err != nil {
		h.Logger.Error("execute template", "error", err)
		http.Error(w, "", http.StatusInternalServerError)
	}
}
