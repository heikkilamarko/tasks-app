package ui

import (
	"errors"
	"log/slog"
	"net/http"
	"tasks-app/internal/shared"
)

type GetUITasksExport struct {
	TaskRepository shared.TaskRepository
	FileExporter   shared.FileExporter
	Logger         *slog.Logger
}

func (h *GetUITasksExport) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var name string
	var tasks []*shared.Task
	var err error

	switch r.FormValue("filter") {
	case "active":
		name = "active_tasks"
		tasks, err = h.TaskRepository.GetActive(r.Context())
	case "completed":
		name = "completed_tasks"
		tasks, err = h.TaskRepository.GetCompleted(r.Context())
	default:
		name = "all_tasks"
		tasks1, err1 := h.TaskRepository.GetActive(r.Context())
		tasks2, err2 := h.TaskRepository.GetCompleted(r.Context())
		tasks, err = append(tasks1, tasks2...), errors.Join(err1, err2)
	}

	if err != nil {
		h.Logger.Error("get tasks", "error", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	if err := h.FileExporter.ExportTasks(w, tasks, name); err != nil {
		h.Logger.Error("export tasks", "error", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
}
