package ui

import (
	"context"
	"log/slog"
	"net/http"
	"tasks-app/internal/shared"
)

type PutUITask struct {
	TxManager                 shared.TxManager
	TaskRepository            shared.TaskRepository
	TaskAttachmentsRepository shared.TaskAttachmentsRepository
	Renderer                  Renderer
	Logger                    *slog.Logger
}

func (h *PutUITask) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	req, err := ParseUpdateTaskRequest(r)
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

	task.Update(req.Name, req.ExpiresAt)

	attachments := BuildAttachmentsUpdate(task.Attachments, req.Attachments.Names)

	if err := h.TxManager.RunInTx(r.Context(), func(ctx context.Context) error {
		if err := h.TaskRepository.Update(ctx, task); err != nil {
			return err
		}

		if err := h.TaskRepository.UpdateAttachments(ctx, task.ID, attachments.Inserted, attachments.Deleted); err != nil {
			return err
		}

		if err := h.TaskAttachmentsRepository.SaveAttachments(ctx, task.ID, req.Attachments.Files); err != nil {
			return err
		}

		if err := h.TaskAttachmentsRepository.DeleteAttachments(ctx, task.ID, attachments.Deleted); err != nil {
			return err
		}

		return nil
	}); err != nil {
		h.Logger.Error("update task", "error", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	task, err = h.TaskRepository.GetByID(r.Context(), req.ID)
	if err != nil {
		h.Logger.Error("get task", "error", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	vm := NewTaskResponse(r, task)

	h.Renderer.Render(w, "active_tasks_table_row.html", vm)
}
