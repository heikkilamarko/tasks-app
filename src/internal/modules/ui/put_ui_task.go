package ui

import (
	"log/slog"
	"net/http"
	"tasks-app/internal/shared"
)

type PutUITask struct {
	TaskRepository            shared.TaskRepository
	TaskAttachmentsRepository shared.TaskAttachmentsRepository
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

	attachmentNames, err := ParseTaskAttachments(r, task.ID, h.TaskAttachmentsRepository)
	if err != nil {
		h.Logger.Error("parse task attachments", "error", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	err = h.TaskRepository.UpdateAttachments(r.Context(), task.ID, attachmentNames)
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

	if err := Templates.ExecuteTemplate(w, "active_tasks_table_row", task); err != nil {
		h.Logger.Error("execute template", "error", err)
		http.Error(w, "", http.StatusInternalServerError)
	}
}
