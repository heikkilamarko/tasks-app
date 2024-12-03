package ui

import (
	"log/slog"
	"net/http"
	"tasks-app/internal/shared"
)

type GetUITaskEdit struct {
	TxManager shared.TxManager
	Renderer  Renderer
	Logger    *slog.Logger
}

func (h *GetUITaskEdit) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	req, err := ParseTaskRequest(r)
	if err != nil {
		h.Logger.Error("parse request", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	h.TxManager.RunInTx(func(txc shared.TxContext) error {
		task, err := txc.TaskRepository.GetByID(r.Context(), req.ID)
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

		return h.Renderer.Render(w, "active_tasks_table_row_edit.html", vm)
	})
}
