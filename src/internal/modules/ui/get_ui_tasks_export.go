package ui

import (
	"errors"
	"log/slog"
	"net/http"
	"tasks-app/internal/shared"
)

type GetUITasksExport struct {
	TxManager    shared.TxManager
	FileExporter shared.FileExporter
	Logger       *slog.Logger
}

func (h *GetUITasksExport) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.TxManager.RunInTx(func(txc shared.TxContext) error {
		var name string
		var tasks []*shared.Task
		var err error

		switch r.FormValue("filter") {
		case "active":
			name = "active_tasks"
			tasks, err = txc.TaskRepository.GetActive(r.Context(), 0, 10_000)
		case "completed":
			name = "completed_tasks"
			tasks, err = txc.TaskRepository.GetCompleted(r.Context(), 0, 10_000)
		default:
			name = "all_tasks"
			tasks1, err1 := txc.TaskRepository.GetActive(r.Context(), 0, 10_000)
			tasks2, err2 := txc.TaskRepository.GetCompleted(r.Context(), 0, 10_000)
			tasks, err = append(tasks1, tasks2...), errors.Join(err1, err2)
		}

		if err != nil {
			h.Logger.Error("get tasks", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return nil
		}

		if err := h.FileExporter.ExportTasks(w, tasks, name); err != nil {
			h.Logger.Error("export tasks", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return nil
		}

		return nil
	})
}
