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

	err = h.TxManager.RunInTx(func(txc shared.TxContext) error {
		if _, err := txc.TaskRepository.GetByID(r.Context(), req.ID); err != nil {
			return err
		}

		if err = txc.TaskRepository.Delete(r.Context(), req.ID); err != nil {
			return err
		}

		return h.TaskAttachmentsRepository.DeleteTask(r.Context(), req.ID)
	})

	if err != nil {
		if err == shared.ErrNotFound {
			http.Error(w, "task not found", http.StatusNotFound)
		} else {
			h.Logger.Error("delete task", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
		}
		return
	}

	var tasks []*shared.Task

	err = h.TxManager.RunInTx(func(txc shared.TxContext) error {
		tasks, err = txc.TaskRepository.GetActive(r.Context(), 0, 50)
		return err
	})

	if err != nil {
		h.Logger.Error("get tasks", "error", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	vm := NewTasksResponse(r, tasks)

	h.Renderer.Render(w, "active_tasks_table.html", vm)
}
