package ui

import (
	"log/slog"
	"net/http"
	"tasks-app/internal/shared"
	"time"
)

type PostUITasks struct {
	TaskRepository shared.TaskRepository
	Logger         *slog.Logger
}

func (h *PostUITasks) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	task := shared.NewTask(name, expiresAt)

	err = h.TaskRepository.Create(r.Context(), task)
	if err != nil {
		h.Logger.Error("create task", "error", err)
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
