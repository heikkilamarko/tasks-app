package ui

import (
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
	tasks, err := h.TaskRepo.GetActive(r.Context())
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
	w.Header().Set("Content-Disposition", "attachment; filename=export.xlsx")

	if _, err := w.Write(data); err != nil {
		h.Logger.Error("write response", "error", err)
		http.Error(w, "", http.StatusInternalServerError)
	}
}
