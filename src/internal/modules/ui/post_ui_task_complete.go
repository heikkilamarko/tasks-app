package ui

import (
	"log/slog"
	"net/http"
	"strconv"
	"tasks-app/internal/shared"

	"github.com/go-chi/chi/v5"
)

type PostUITaskComplete struct {
	TaskRepository shared.TaskRepository
	Logger         *slog.Logger
}

func (h *PostUITaskComplete) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	task.SetCompleted()

	err = h.TaskRepository.Update(r.Context(), task)
	if err != nil {
		h.Logger.Error("update task", "error", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

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
		IsCreatingNew: false,
	}

	if err := Templates.ExecuteTemplate(w, "active_tasks_table.html", data); err != nil {
		h.Logger.Error("execute template", "error", err)
		http.Error(w, "", http.StatusInternalServerError)
	}
}
