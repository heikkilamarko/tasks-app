package ui

import (
	"log/slog"
	"net/http"
	"tasks-app/internal/shared"
)

type PutUITask struct {
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

	err = h.TaskRepository.Update(r.Context(), task)
	if err != nil {
		h.Logger.Error("update task", "error", err)
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

	err = h.TaskAttachmentsRepository.DeleteAttachments(r.Context(), task.ID, u.Deleted)
	if err != nil {
		h.Logger.Error("update task attachments", "error", err)
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
