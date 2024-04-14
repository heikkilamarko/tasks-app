package ui

import (
	"context"
	"log/slog"
	"net/http"
	"tasks-app/internal/shared"
)

type PostUITasks struct {
	TxManager                 shared.TxManager
	TaskRepository            shared.TaskRepository
	TaskAttachmentsRepository shared.TaskAttachmentsRepository
	Renderer                  Renderer
	Logger                    *slog.Logger
}

func (h *PostUITasks) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	req, err := ParseNewTaskRequest(r)
	if err != nil {
		h.Logger.Error("parse request", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	task := shared.NewTask(req.Name, req.ExpiresAt)

	attachments := BuildAttachmentsUpdate(task.Attachments, req.Attachments.Names)

	if err := h.TxManager.RunInTx(r.Context(), func(ctx context.Context) error {
		if err := h.TaskRepository.Create(ctx, task); err != nil {
			return err
		}

		if err := h.TaskRepository.UpdateAttachments(ctx, task.ID, attachments.Inserted, attachments.Deleted); err != nil {
			return err
		}

		if err := h.TaskAttachmentsRepository.SaveAttachments(ctx, task.ID, req.Attachments.Files); err != nil {
			return err
		}

		return nil
	}); err != nil {
		h.Logger.Error("create task", "error", err)
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
