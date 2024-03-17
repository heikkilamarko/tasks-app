package ui

import (
	"log/slog"
	"net/http"
	"tasks-app/internal/shared"
)

type PostUITasks struct {
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

	err = h.TaskRepository.Create(r.Context(), task)
	if err != nil {
		h.Logger.Error("create task", "error", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	err = h.TaskAttachmentsRepository.SaveAttachments(r.Context(), task.ID, req.Attachments.Files)
	if err != nil {
		h.Logger.Error("save attachments", "error", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	u := BuildAttachmentsUpdate(task.Attachments, req.Attachments.Names)
	err = h.TaskRepository.UpdateAttachments(r.Context(), task.ID, u.Inserted, u.Deleted)
	if err != nil {
		h.Logger.Error("update task attachments", "error", err)
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
