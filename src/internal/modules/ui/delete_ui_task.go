package ui

import (
	"log/slog"
	"net/http"
	"tasks-app/internal/shared"
)

type DeleteUITask struct {
	TxManager                 shared.TxManager
	TaskAttachmentsRepository shared.TaskAttachmentsRepository
	Renderer                  Renderer
	Logger                    *slog.Logger
}

func (h *DeleteUITask) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

		if err := txc.TaskRepository.Delete(r.Context(), req.ID); err != nil {
			h.Logger.Error("delete task", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return nil
		}

		if err := h.TaskAttachmentsRepository.DeleteTask(r.Context(), req.ID); err != nil {
			h.Logger.Error("delete task", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return nil
		}

		tasks, err := txc.TaskRepository.GetActive(r.Context(), 0, 50)
		if err != nil {
			h.Logger.Error("get tasks", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return nil
		}

		vm := NewTasksResponse(r, tasks)

		return h.Renderer.Render(w, "active_tasks_table.html", vm)
	})
}
