package ui

import (
	"log/slog"
	"net/http"
	"tasks-app/internal/shared"
)

type PostUITasks struct {
	TxManager                 shared.TxManager
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

	h.TxManager.RunInTx(func(txc shared.TxContext) error {
		if err := txc.TaskRepository.Create(r.Context(), task); err != nil {
			h.Logger.Error("create task", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return nil
		}

		if err := txc.TaskRepository.UpdateAttachments(r.Context(), task.ID, attachments.Inserted, attachments.Deleted); err != nil {
			h.Logger.Error("create task", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return nil
		}

		if err := h.TaskAttachmentsRepository.SaveAttachments(r.Context(), task.ID, req.Attachments.Files); err != nil {
			h.Logger.Error("create task", "error", err)
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
