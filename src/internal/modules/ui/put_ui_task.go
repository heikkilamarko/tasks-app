package ui

import (
	"log/slog"
	"net/http"
	"tasks-app/internal/shared"
)

type PutUITask struct {
	TxManager                 shared.TxManager
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

		task.Update(req.Name, req.ExpiresAt)

		attachments := BuildAttachmentsUpdate(task.Attachments, req.Attachments.Names)

		if err := txc.TaskRepository.Update(r.Context(), task); err != nil {
			h.Logger.Error("update task", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return nil
		}

		if err := txc.TaskRepository.UpdateAttachments(r.Context(), task.ID, attachments.Inserted, attachments.Deleted); err != nil {
			h.Logger.Error("update task", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return nil
		}

		if err := h.TaskAttachmentsRepository.SaveAttachments(r.Context(), task.ID, req.Attachments.Files); err != nil {
			h.Logger.Error("update task", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return nil
		}

		if err := h.TaskAttachmentsRepository.DeleteAttachments(r.Context(), task.ID, attachments.Deleted); err != nil {
			h.Logger.Error("update task", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return nil
		}

		task, err = txc.TaskRepository.GetByID(r.Context(), req.ID)
		if err != nil {
			h.Logger.Error("get task", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return nil
		}

		vm := NewTaskResponse(r, task)

		return h.Renderer.Render(w, "active_tasks_table_row.html", vm)
	})
}
