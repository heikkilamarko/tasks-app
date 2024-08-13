package ui

import (
	"log/slog"
	"net/http"
	"tasks-app/internal/shared"
)

type GetUITaskAttachment struct {
	TaskRepository            shared.TaskRepository
	TaskAttachmentsRepository shared.TaskAttachmentsRepository
	Logger                    *slog.Logger
}

func (h *GetUITaskAttachment) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	req, err := ParseTaskAttachmentRequest(r)
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

	data, err := h.TaskAttachmentsRepository.GetAttachment(r.Context(), req.ID, req.Name)
	if err != nil {
		h.Logger.Error("get task attachment", "error", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	if len(data) == 0 {
		http.Error(w, "task attachment not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", http.DetectContentType(data))
	w.Write(data)
}
