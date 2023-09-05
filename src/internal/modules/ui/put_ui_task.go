package ui

import (
	"log/slog"
	"net/http"
	"strconv"
	"tasks-app/internal/shared"
	"time"

	"github.com/go-chi/chi/v5"
)

type PutUITask struct {
	TaskRepo shared.TaskRepository
	Logger   *slog.Logger
}

func (h *PutUITask) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "failed to parse form data", http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	if len(name) < 1 {
		http.Error(w, "invalid name", http.StatusBadRequest)
		return
	}

	var expiresAt *time.Time
	expiresAtStr := r.FormValue("expires_at")
	if expiresAtStr != "" {
		expiresAtTemp, err := ParseUITime(expiresAtStr)
		if err != nil {
			http.Error(w, "invalid expires_at format", http.StatusBadRequest)
			return
		}
		expiresAt = &expiresAtTemp
	}

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	task, err := h.TaskRepo.GetByID(r.Context(), id)
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

	err = h.TaskRepo.Update(r.Context(), task)
	if err != nil {
		h.Logger.Error("update task", "error", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	if err := Templates.ExecuteTemplate(w, "active_tasks_table_row.html", task); err != nil {
		h.Logger.Error("execute template", "error", err)
		http.Error(w, "", http.StatusInternalServerError)
	}
}
