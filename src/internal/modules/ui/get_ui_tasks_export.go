package ui

import (
	"errors"
	"log/slog"
	"net/http"
	"tasks-app/internal/shared"
)

type GetUITasksExport struct {
	TaskRepo     shared.TaskRepository
	FileExporter shared.FileExporter
	Logger       *slog.Logger
}

func (h *GetUITasksExport) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var tasks []*shared.Task
	var err error

	switch r.FormValue("filter") {
	case "active":
		tasks, err = h.TaskRepo.GetActive(r.Context())
	case "completed":
		tasks, err = h.TaskRepo.GetCompleted(r.Context())
	default:
		tasks1, err1 := h.TaskRepo.GetActive(r.Context())
		tasks2, err2 := h.TaskRepo.GetCompleted(r.Context())
		tasks, err = append(tasks1, tasks2...), errors.Join(err1, err2)
	}

	if err != nil {
		h.Logger.Error("get tasks", "error", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	data, err := h.FileExporter.Export(tasks)
	if err != nil {
		h.Logger.Error("export tasks", "error", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", "attachment; filename=tasks.xlsx")

	if _, err := w.Write(data); err != nil {
		h.Logger.Error("write response", "error", err)
		http.Error(w, "", http.StatusInternalServerError)
	}
}
