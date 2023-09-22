package ui

import (
	"log/slog"
	"net/http"
	"tasks-app/internal/shared"
)

type PostUITasks struct {
	TaskRepository shared.TaskRepository
	Logger         *slog.Logger
}

func (h *PostUITasks) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	task := shared.NewTask(name, expiresAt)

	err := h.TaskRepository.Create(r.Context(), task)
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

	vm := TasksViewModel{tasks, false}

	if err := Templates.ExecuteTemplate(w, "active_tasks_table", vm); err != nil {
		h.Logger.Error("execute template", "error", err)
		http.Error(w, "", http.StatusInternalServerError)
	}
}
