package ui

import (
	"log/slog"
	"net/http"
	"tasks-app/internal/shared"
)

type GetUI struct {
	TaskRepo shared.TaskRepository
	Logger   *slog.Logger
}

func (h *GetUI) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.TaskRepo.GetActive(r.Context())
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

	if err := Templates.ExecuteTemplate(w, "active_tasks.html", data); err != nil {
		h.Logger.Error("execute template", "error", err)
		http.Error(w, "", http.StatusInternalServerError)
	}
}
