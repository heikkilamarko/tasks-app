package ui

import (
	"errors"
	"log/slog"
	"net/http"
	"tasks-app/internal/shared"
)

type GetUITaskAttachment struct {
	TxProvider                shared.TxProvider
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

	h.TxProvider.Transact(func(adapters shared.TxAdapters) error {
		task, err := adapters.TaskRepository.GetByID(r.Context(), req.ID)
		if err != nil {
			h.Logger.Error("get task", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return err
		}

		if task == nil {
			http.Error(w, "task not found", http.StatusNotFound)
			return errors.New("task not found")
		}

		return nil
	})
	if err != nil {
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
