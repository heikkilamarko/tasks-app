package ui

import (
	"context"
	"log/slog"
	"net/http"
	"tasks-app/internal/shared"
)

type DeleteUITask struct {
	TxManager                 shared.TxManager
	TaskRepository            shared.TaskRepository
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

	task, err := h.TaskRepository.GetByID(r.Context(), req.ID)
	if err != nil {
		h.Logger.Error("get task", "error", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	if task == nil {
		http.Error(w, "task not found", http.StatusNotFound)
		return
	}

	if err := h.TxManager.RunInTx(r.Context(), func(ctx context.Context) error {
		if err := h.TaskRepository.Delete(ctx, req.ID); err != nil {
			return err
		}

		if err := h.TaskAttachmentsRepository.DeleteTask(ctx, req.ID); err != nil {
			return err
		}

		return nil
	}); err != nil {
		h.Logger.Error("delete task", "error", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	tasks, err := h.TaskRepository.GetActive(r.Context(), 0, 50)
	if err != nil {
		h.Logger.Error("get tasks", "error", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	vm := NewTasksResponse(r, tasks)

	h.Renderer.Render(w, "active_tasks_table.html", vm)
}
