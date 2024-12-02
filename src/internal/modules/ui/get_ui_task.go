package ui

import (
	"log/slog"
	"net/http"
	"tasks-app/internal/shared"
)

type GetUITask struct {
	TxProvider shared.TxProvider
	Renderer   Renderer
	Logger     *slog.Logger
}

func (h *GetUITask) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	req, err := ParseTaskRequest(r)
	if err != nil {
		h.Logger.Error("parse request", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	h.TxProvider.Transact(func(adapters shared.TxAdapters) error {
		task, err := adapters.TaskRepository.GetByID(r.Context(), req.ID)
		if err != nil {
			h.Logger.Error("get task", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return nil
		}

		if task == nil {
			http.Error(w, "task not found", http.StatusNotFound)
			return nil
		}

		vm := NewTaskResponse(r, task)

		return h.Renderer.Render(w, "active_tasks_table_row.html", vm)
	})
}
